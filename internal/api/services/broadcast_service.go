package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	SendToDestinations(ctx context.Context, destinationIDs []*database.Destination, targets []string, message string) (*BroadcastResult, error)
	SendToTargets(ctx context.Context, targets []database.TeamsTarget, message string) (*BroadcastResult, error)
}

type broadcastService struct {
	teamsBotRepo    repositories.TeamsBotRepository
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
	teamsBotRepo repositories.TeamsBotRepository,
	installRepo repositories.BotInstallationRepository,
	destinationRepo repositories.DestinationRepository,
) BroadcastService {
	return &broadcastService{
		teamsBotRepo:    teamsBotRepo,
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
	bot, err := s.teamsBotRepo.GetByID(ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform bot: %w", err)
	}
	appID = bot.AppID
	// Note: In production, you'd need to decrypt the password hash
	// For now, we'll assume it's stored in plain text (not recommended for production)
	appPassword = bot.AppPasswordHash

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
	bot, err := s.teamsBotRepo.GetByID(ctx, botID)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform bot: %w", err)
	}
	appID = bot.AppID
	appPassword = bot.AppPasswordHash

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
func (s *broadcastService) SendToDestinations(ctx context.Context, destinations []*database.Destination, targets []string, message string) (*BroadcastResult, error) {
	start := time.Now()

	if len(destinations) == 0 {
		return &BroadcastResult{
			TotalTargets: 0,
			Successful:   0,
			Failed:       0,
			Results:      []TargetResult{},
			Duration:     time.Since(start),
		}, nil
	}

	// white list filter targets
	targetsMap := make(map[string]struct{})
	for _, target := range targets {
		targetsMap[target] = struct{}{}
	}

	// Convert destinations to targets and send
	allTargets := make([]database.TeamsTarget, 0)
	log.Printf("Targets map: %v", targetsMap)
	for _, dest := range destinations {
		for _, target := range dest.Targets {
			if len(targetsMap) > 0 {
				if _, ok := targetsMap[target.Email]; ok {
					allTargets = append(allTargets, target)
					log.Printf("Added email target: %s", target.Email)
					continue
				}
				if _, ok := targetsMap[target.ConversationID]; ok {
					log.Printf("Added conversation target: %s", target.ConversationID)
					allTargets = append(allTargets, target)
				}
			} else {
				log.Printf("Original target: %s", target.ConversationID)
				allTargets = append(allTargets, target)
			}
		}
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

	// Basic validation and early failures
	validated := make([]database.TeamsTarget, 0, len(targets))
	earlyResults := make([]TargetResult, 0)
	for _, t := range targets {
		t.Type = strings.TrimSpace(t.Type)
		switch t.Type {
		case "personal":
			if t.ConversationID == "" && t.Email == "" {
				earlyResults = append(earlyResults, TargetResult{
					TargetID:       "",
					ConversationID: "",
					Scope:          t.Type,
					Success:        false,
					Error:          "personal target requires conversation_id or email",
				})
				continue
			}
		case "channel", "groupChat":
			if t.ConversationID == "" {
				earlyResults = append(earlyResults, TargetResult{
					TargetID:       "",
					ConversationID: "",
					Scope:          t.Type,
					Success:        false,
					Error:          "channel/groupChat target requires conversation_id",
				})
				continue
			}
		default:
			earlyResults = append(earlyResults, TargetResult{
				TargetID:       "",
				ConversationID: "",
				Scope:          t.Type,
				Success:        false,
				Error:          "unsupported target type",
			})
			continue
		}
		validated = append(validated, t)
	}

	if len(validated) == 0 {
		return &BroadcastResult{
			TotalTargets: len(targets),
			Successful:   0,
			Failed:       len(earlyResults),
			Results:      earlyResults,
			Duration:     time.Since(start),
		}, nil
	}

	// Since we only have a single tenant, no need to group by tenant
	// Get tenant ID from first target (all targets should have same tenant)
	tenantID := ""
	if len(validated) > 0 {
		tenantID = validated[0].TenantID
	}

	// Get all active installations for this tenant
	installations, err := s.installRepo.GetActiveInstallationsByTenant(ctx, tenantID)
	if err != nil {
		log.Printf("Failed to get installations for tenant %s: %v", tenantID, err)
		return &BroadcastResult{
			TotalTargets: len(targets),
			Successful:   0,
			Failed:       len(targets),
			Results:      earlyResults,
			Duration:     time.Since(start),
		}, fmt.Errorf("failed to get installations: %w", err)
	}
	log.Printf("Found %d installations for tenant %s", len(installations), tenantID)

	allResults := make([]TargetResult, 0)
	allResults = append(allResults, earlyResults...)
	totalSuccessful := 0
	totalFailed := len(earlyResults)

	// Match targets to installations and send
	for _, target := range validated {
		// For personal email-only target, try resolve to existing personal installation by email
		if target.Type == "personal" && target.ConversationID == "" && target.Email != "" {
			if emailInstalls, err := s.installRepo.GetActivePersonalByEmail(ctx, tenantID, target.Email); err == nil {
				installations = append(installations, emailInstalls...)
			} else {
				log.Printf("Lookup personal by email failed: %v", err)
			}
		}
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
		// 優先從資料庫獲取，如果沒有則使用環境變數
		appID := ""
		appPassword := ""
		bot, err := s.teamsBotRepo.GetByID(ctx, matchedInstallation.BotID)
		if err != nil {
			log.Printf("Failed to get teams bot: %v", err)
		}
		appID = bot.AppID
		appPassword = bot.AppPasswordHash
		log.Printf("Using DB credentials: appID=%s, password length=%d", appID, len(appPassword))

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
	log.Printf("sendToInstallation: Starting for installation %s, conversation %s", inst.ID, inst.ConversationID)

	// Get token
	log.Printf("sendToInstallation: Getting token for tenant %s, appID %s", inst.TeamsTenantID, appID)
	token, err := s.getConnectorToken(ctx, inst.TeamsTenantID, appID, appPassword)
	if err != nil {
		log.Printf("sendToInstallation: Failed to get token: %v", err)
		return TargetResult{
			TargetID:       inst.ConversationID,
			ConversationID: inst.ConversationID,
			Scope:          inst.ConversationType,
			Success:        false,
			Error:          fmt.Sprintf("Failed to get token: %v", err),
			ResponseTime:   time.Since(start),
		}
	}
	log.Printf("sendToInstallation: Successfully got token")

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
	log.Printf("sendToInstallation: Prepared payload: %s", string(reqBody))

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		log.Printf("sendToInstallation: Failed to create request: %v", err)
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
	log.Printf("sendToInstallation: Sending request to %s, conversation_id=%s", url, inst.ConversationID)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("sendToInstallation: Failed to send request: %v", err)
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

	log.Printf("sendToInstallation: Received response status %d: %s", resp.StatusCode, resp.Status)

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("sendToInstallation: Error response body: %s", string(body))
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

	// Use the same implementation as proactive message test
	if tenantID == "" {
		tenantID = "botframework.com"
	}

	log.Printf("getConnectorToken: Requesting token from tenant %s, appID %s", tenantID, appID)

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", appID)
	form.Set("client_secret", appPassword)
	form.Set("scope", "https://api.botframework.com/.default")

	tokenURL := "https://login.microsoftonline.com/" + url.PathEscape(tenantID) + "/oauth2/v2.0/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		log.Printf("getConnectorToken: Failed to create request: %v", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("getConnectorToken: HTTP request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("getConnectorToken: Response status: %d %s", resp.StatusCode, resp.Status)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("getConnectorToken: Error response body: %s", string(body))
		return "", fmt.Errorf("token request failed: %s: %s", resp.Status, string(body))
	}

	var tr struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		log.Printf("getConnectorToken: Failed to decode response: %v", err)
		return "", err
	}

	if tr.AccessToken == "" {
		log.Printf("getConnectorToken: Empty access token in response")
		return "", fmt.Errorf("empty access_token")
	}

	log.Printf("getConnectorToken: Successfully got token")
	return tr.AccessToken, nil
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

	// Primary: conversation_type + conversation_id match
	if inst.ConversationType == desiredType && target.ConversationID != "" && inst.ConversationID == target.ConversationID {
		return true
	}
	// Fallback for personal email-only targets: match by email
	if desiredType == "personal" && target.Email != "" && !strings.EqualFold(target.Email, "") {
		return inst.Email != "" && strings.EqualFold(inst.Email, target.Email)
	}
	return false
}
