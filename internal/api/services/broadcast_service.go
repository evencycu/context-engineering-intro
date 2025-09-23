package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/repositories"
	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/google/uuid"
)

// BroadcastService handles sending messages to multiple Teams targets
type BroadcastService interface {
	SendToAllActive(ctx context.Context, botID uuid.UUID, botType database.BotType, tenantID string, message string) (*BroadcastResult, error)
	SendToScope(ctx context.Context, botID uuid.UUID, botType database.BotType, tenantID string, scope string, message string) (*BroadcastResult, error)
	SendToDestinations(ctx context.Context, destinationIDs []uuid.UUID, message string) (*BroadcastResult, error)
	SendToTargets(ctx context.Context, targets []database.TeamsTarget, message string) (*BroadcastResult, error)
}

type broadcastService struct {
	platformRepo    repositories.PlatformBotRepository
	thirdPartyRepo  repositories.ThirdPartyBotRepository
	installRepo     repositories.BotInstallationRepository
	destinationRepo repositories.DestinationRepository
}

// BroadcastResult represents the result of a broadcast operation
type BroadcastResult struct {
	TotalTargets int            `json:"total_targets"`
	Successful   int            `json:"successful"`
	Failed       int            `json:"failed"`
	Results      []TargetResult `json:"results"`
	Duration     time.Duration  `json:"duration"`
	Error        string         `json:"error,omitempty"`
}

// TargetResult represents the result for a single target
type TargetResult struct {
	TargetID       string        `json:"target_id"`
	ConversationID string        `json:"conversation_id"`
	Scope          string        `json:"scope"`
	Success        bool          `json:"success"`
	Error          string        `json:"error,omitempty"`
	ResponseTime   time.Duration `json:"response_time"`
	TeamsMessageID string        `json:"teams_message_id,omitempty"`
}

// NewBroadcastService creates a new broadcast service
func NewBroadcastService(
	platformRepo repositories.PlatformBotRepository,
	thirdPartyRepo repositories.ThirdPartyBotRepository,
	installRepo repositories.BotInstallationRepository,
	destinationRepo repositories.DestinationRepository,
) BroadcastService {
	return &broadcastService{
		platformRepo:    platformRepo,
		thirdPartyRepo:  thirdPartyRepo,
		installRepo:     installRepo,
		destinationRepo: destinationRepo,
	}
}

// SendToAllActive sends a message to all active installations for a bot and tenant
func (s *broadcastService) SendToAllActive(ctx context.Context, botID uuid.UUID, botType database.BotType, tenantID string, message string) (*BroadcastResult, error) {
	start := time.Now()

	// Get all active installations
	installations, err := s.installRepo.GetActiveInstallations(ctx, botID, botType, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active installations: %w", err)
	}

	if len(installations) == 0 {
		return &BroadcastResult{
			TotalTargets: 0,
			Successful:   0,
			Failed:       0,
			Results:      []TargetResult{},
			Duration:     time.Since(start),
		}, nil
	}

	// Get bot credentials
	var appID, appPassword string
	if botType == database.BotTypePlatform {
		bot, err := s.platformRepo.GetByID(ctx, botID)
		if err != nil {
			return nil, fmt.Errorf("failed to get platform bot: %w", err)
		}
		appID = bot.AppID
		// Note: In production, you'd need to decrypt the password hash
		// For now, we'll assume it's stored in plain text (not recommended for production)
		appPassword = bot.AppPasswordHash
	} else {
		bot, err := s.thirdPartyRepo.GetByID(ctx, botID)
		if err != nil {
			return nil, fmt.Errorf("failed to get third-party bot: %w", err)
		}
		appID = bot.AppID
		appPassword = bot.AppPasswordHash
	}

	// Send to all installations
	results := make([]TargetResult, 0, len(installations))
	successful := 0
	failed := 0

	for _, inst := range installations {
		result := s.sendToInstallation(ctx, inst, appID, appPassword, message)
		results = append(results, result)

		if result.Success {
			successful++
			// Update activity timestamp
			s.installRepo.UpdateActivity(ctx, inst.ID)
		} else {
			failed++
			// Mark as stale if it fails consistently
			if strings.Contains(result.Error, "unauthorized") || strings.Contains(result.Error, "forbidden") {
				s.installRepo.MarkAsStale(ctx, inst.ID)
			}
		}
	}

	return &BroadcastResult{
		TotalTargets: len(installations),
		Successful:   successful,
		Failed:       failed,
		Results:      results,
		Duration:     time.Since(start),
	}, nil
}

// SendToScope sends a message to all active installations of a specific scope
func (s *broadcastService) SendToScope(ctx context.Context, botID uuid.UUID, botType database.BotType, tenantID string, scope string, message string) (*BroadcastResult, error) {
	start := time.Now()

	// Get installations by scope
	installations, err := s.installRepo.GetByConversationType(ctx, scope, "active")
	if err != nil {
		return nil, fmt.Errorf("failed to get installations by scope: %w", err)
	}

	// Filter by bot and tenant
	filtered := make([]*database.BotInstallation, 0)
	for _, inst := range installations {
		if inst.BotID == botID && inst.BotType == botType && inst.TeamsTenantID == tenantID {
			filtered = append(filtered, inst)
		}
	}

	if len(filtered) == 0 {
		return &BroadcastResult{
			TotalTargets: 0,
			Successful:   0,
			Failed:       0,
			Results:      []TargetResult{},
			Duration:     time.Since(start),
		}, nil
	}

	// Get bot credentials
	var appID, appPassword string
	if botType == database.BotTypePlatform {
		bot, err := s.platformRepo.GetByID(ctx, botID)
		if err != nil {
			return nil, fmt.Errorf("failed to get platform bot: %w", err)
		}
		appID = bot.AppID
		appPassword = bot.AppPasswordHash
	} else {
		bot, err := s.thirdPartyRepo.GetByID(ctx, botID)
		if err != nil {
			return nil, fmt.Errorf("failed to get third-party bot: %w", err)
		}
		appID = bot.AppID
		appPassword = bot.AppPasswordHash
	}

	// Send to filtered installations
	results := make([]TargetResult, 0, len(filtered))
	successful := 0
	failed := 0

	for _, inst := range filtered {
		result := s.sendToInstallation(ctx, inst, appID, appPassword, message)
		results = append(results, result)

		if result.Success {
			successful++
			s.installRepo.UpdateActivity(ctx, inst.ID)
		} else {
			failed++
			if strings.Contains(result.Error, "unauthorized") || strings.Contains(result.Error, "forbidden") {
				s.installRepo.MarkAsStale(ctx, inst.ID)
			}
		}
	}

	return &BroadcastResult{
		TotalTargets: len(filtered),
		Successful:   successful,
		Failed:       failed,
		Results:      results,
		Duration:     time.Since(start),
	}, nil
}

// SendToDestinations sends a message to specific destinations
func (s *broadcastService) SendToDestinations(ctx context.Context, destinationIDs []uuid.UUID, message string) (*BroadcastResult, error) {
	start := time.Now()

	// Get destinations
	destinations := make([]*database.Destination, 0, len(destinationIDs))
	for _, id := range destinationIDs {
		dest, err := s.destinationRepo.GetByID(ctx, id)
		if err != nil {
			log.Printf("Failed to get destination %s: %v", id, err)
			continue
		}
		destinations = append(destinations, dest)
	}

	if len(destinations) == 0 {
		return &BroadcastResult{
			TotalTargets: 0,
			Successful:   0,
			Failed:       0,
			Results:      []TargetResult{},
			Duration:     time.Since(start),
		}, nil
	}

	// Convert destinations to targets and send
	allTargets := make([]database.TeamsTarget, 0)
	for _, dest := range destinations {
		// dest.Targets is already a JSONB-backed slice of TeamsTarget
		allTargets = append(allTargets, dest.Targets...)
	}

	return s.SendToTargets(ctx, allTargets, message)
}

// SendToTargets sends a message to specific Teams targets
func (s *broadcastService) SendToTargets(ctx context.Context, targets []database.TeamsTarget, message string) (*BroadcastResult, error) {
	start := time.Now()

	if len(targets) == 0 {
		return &BroadcastResult{
			TotalTargets: 0,
			Successful:   0,
			Failed:       0,
			Results:      []TargetResult{},
			Duration:     time.Since(start),
		}, nil
	}

	// Group targets by tenant (no per-bot grouping due to schema)
	tenantGroups := make(map[string][]database.TeamsTarget)
	for _, target := range targets {
		if target.TenantID == "" {
			continue
		}
		tenantGroups[target.TenantID] = append(tenantGroups[target.TenantID], target)
	}

	// Send to each tenant group
	allResults := make([]TargetResult, 0)
	totalSuccessful := 0
	totalFailed := 0

	for tenantID, botTargets := range tenantGroups {
		// Get all active installations for this tenant
		installations, err := s.installRepo.GetActiveInstallationsByTenant(ctx, tenantID)
		if err != nil {
			log.Printf("Failed to get installations for tenant %s: %v", tenantID, err)
			continue
		}
		log.Printf("Found %d installations for tenant %s", len(installations), tenantID)
		// Match targets to installations
		for _, target := range botTargets {
			log.Printf("Matching target: type=%s, conversation_id=%s", target.Type, target.ConversationID)
			var matchedInstallation *database.BotInstallation
			for _, inst := range installations {
				log.Printf("Checking installation: type=%s, conversation_id=%s", inst.ConversationType, inst.ConversationID)
				if s.matchesTarget(inst, target) {
					matchedInstallation = inst
					log.Printf("Found matching installation: %s", inst.ID)
					break
				}
			}

			if matchedInstallation == nil {
				log.Printf("No matching installation found for target: %s", target.ConversationID)
				allResults = append(allResults, TargetResult{
					TargetID:       target.ConversationID,
					ConversationID: target.ConversationID,
					Scope:          target.Type,
					Success:        false,
					Error:          "No matching installation found",
				})
				totalFailed++
				continue
			}

			// Get bot credentials
			// 若 DB 取不到，退回使用環境變數
			appID := ""
			appPassword := ""
			if bot, err := s.platformRepo.GetByID(ctx, matchedInstallation.BotID); err == nil {
				appID = bot.AppID
				appPassword = bot.AppPasswordHash
			} else {
				appID = os.Getenv("TEAMS_BOT_APP_ID")
				appPassword = os.Getenv("TEAMS_BOT_APP_PASSWORD")
			}

			// Send message
			result := s.sendToInstallation(ctx, matchedInstallation, appID, appPassword, message)
			allResults = append(allResults, result)

			if result.Success {
				totalSuccessful++
				s.installRepo.UpdateActivity(ctx, matchedInstallation.ID)
			} else {
				totalFailed++
				if strings.Contains(result.Error, "unauthorized") || strings.Contains(result.Error, "forbidden") {
					s.installRepo.MarkAsStale(ctx, matchedInstallation.ID)
				}
			}
		}
	}

	return &BroadcastResult{
		TotalTargets: len(targets),
		Successful:   totalSuccessful,
		Failed:       totalFailed,
		Results:      allResults,
		Duration:     time.Since(start),
	}, nil
}

// sendToInstallation sends a message to a specific installation
func (s *broadcastService) sendToInstallation(ctx context.Context, inst *database.BotInstallation, appID, appPassword, message string) TargetResult {
	start := time.Now()

	// Get token
	token, err := s.getConnectorToken(ctx, inst.TeamsTenantID, appID, appPassword)
	if err != nil {
		return TargetResult{
			TargetID:       inst.ConversationID,
			ConversationID: inst.ConversationID,
			Scope:          inst.ConversationType,
			Success:        false,
			Error:          fmt.Sprintf("Failed to get token: %v", err),
			ResponseTime:   time.Since(start),
		}
	}

	// Prepare message payload
	payload := map[string]any{
		"type": "message",
		"text": message,
		"from": map[string]any{
			"id":   inst.RecipientID,
			"name": "Broadcast Bot",
		},
		"recipient": map[string]any{
			"id": inst.FromID,
		},
		"conversation": map[string]any{
			"id": inst.ConversationID,
		},
		"channelId":  "msteams",
		"serviceUrl": inst.ServiceURL,
	}

	// Send message
	url := inst.ServiceURL + "/v3/conversations/" + inst.ConversationID + "/activities"
	reqBody, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return TargetResult{
			TargetID:       inst.ConversationID,
			ConversationID: inst.ConversationID,
			Scope:          inst.ConversationType,
			Success:        false,
			Error:          fmt.Sprintf("Failed to create request: %v", err),
			ResponseTime:   time.Since(start),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return TargetResult{
			TargetID:       inst.ConversationID,
			ConversationID: inst.ConversationID,
			Scope:          inst.ConversationType,
			Success:        false,
			Error:          fmt.Sprintf("Failed to send request: %v", err),
			ResponseTime:   time.Since(start),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return TargetResult{
			TargetID:       inst.ConversationID,
			ConversationID: inst.ConversationID,
			Scope:          inst.ConversationType,
			Success:        false,
			Error:          fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
			ResponseTime:   time.Since(start),
		}
	}

	// Parse response to get message ID
	var response map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&response); err == nil {
		if id, ok := response["id"].(string); ok {
			return TargetResult{
				TargetID:       inst.ConversationID,
				ConversationID: inst.ConversationID,
				Scope:          inst.ConversationType,
				Success:        true,
				ResponseTime:   time.Since(start),
				TeamsMessageID: id,
			}
		}
	}

	return TargetResult{
		TargetID:       inst.ConversationID,
		ConversationID: inst.ConversationID,
		Scope:          inst.ConversationType,
		Success:        true,
		ResponseTime:   time.Since(start),
	}
}

// getConnectorToken gets a Bot Framework connector token
func (s *broadcastService) getConnectorToken(ctx context.Context, tenantID, appID, appPassword string) (string, error) {
	// This is a simplified implementation
	// In production, you'd implement proper OAuth2 Client Credentials flow
	// For now, we'll use the same logic as in the existing MessagesService

	// Use tenant-specific endpoint
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	// Prepare request
	data := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     appID,
		"client_secret": appPassword,
		"scope":         "https://api.botframework.com/.default",
	}

	// Send request
	reqBody := make([]string, 0, len(data))
	for k, v := range data {
		reqBody = append(reqBody, fmt.Sprintf("%s=%s", k, v))
	}

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(strings.Join(reqBody, "&")))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("token request failed: %d %s", resp.StatusCode, resp.Status)
	}

	var tokenResp map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	accessToken, ok := tokenResp["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access token in response")
	}

	return accessToken, nil
}

// matchesTarget checks if an installation matches a target
func (s *broadcastService) matchesTarget(inst *database.BotInstallation, target database.TeamsTarget) bool {
	// Map target.Type to conversation_type used in DB
	desiredType := ""
	switch target.Type {
	case "personal":
		desiredType = "personal"
	case "channel":
		desiredType = "channel"
	case "groupChat":
		desiredType = "groupChat"
	default:
		return false
	}

	// Only rely on conversation_type + conversation_id
	return inst.ConversationType == desiredType && inst.ConversationID == target.ConversationID
}
