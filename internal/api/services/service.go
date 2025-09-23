package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
    "log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/repositories"
	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service interface defines common service operations
type Service[T any] interface {
	Create(ctx context.Context, req *CreateRequest[T]) (*T, error)
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[T]) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, req *ListRequest) ([]*T, error)
	Count(ctx context.Context, req *CountRequest) (int64, error)
}

// CreateRequest represents a create request
type CreateRequest[T any] struct {
	Data T `json:"data" validate:"required"`
}

// UpdateRequest represents an update request
type UpdateRequest[T any] struct {
	Data T `json:"data" validate:"required"`
}

// ListRequest represents a list request
type ListRequest struct {
	Limit  int    `form:"limit" binding:"min=1,max=100"`
	Offset int    `form:"offset" binding:"min=0"`
	SortBy string `form:"sort_by"`
	Order  string `form:"order" binding:"oneof=asc desc"`
	Search string `form:"search"`
	Status string `form:"status"`
}

// CountRequest represents a count request
type CountRequest struct {
	Search string `form:"search"`
	Status string `form:"status"`
}

// CompanyService defines company-specific operations
type CompanyService interface {
	Service[database.Company]
	GetByEmail(ctx context.Context, email string) (*database.Company, error)
	GetByStatus(ctx context.Context, status string) ([]*database.Company, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	EnableBilling(ctx context.Context, id uuid.UUID) error
	DisableBilling(ctx context.Context, id uuid.UUID) error
}

// companyService implements CompanyService
type companyService struct {
	repo repositories.CompanyRepository
}

// NewCompanyService creates a new company service
func NewCompanyService(repo repositories.CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

// UserService defines the interface for user operations
type UserService interface {
	Service[database.User]
	GetByEmail(ctx context.Context, email string) (*database.User, error)
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.User, error)
	GetByRole(ctx context.Context, role string) ([]*database.User, error)
	GetByStatus(ctx context.Context, status string) ([]*database.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	GenerateAPIKey(ctx context.Context, userID uuid.UUID) (string, error)
	RevokeAPIKey(ctx context.Context, userID uuid.UUID) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}

// userService implements UserService
type userService struct {
	repo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, req *CreateRequest[database.User]) (*database.User, error) {
	// Set default values
	req.Data.ID = uuid.New()
	req.Data.CreatedAt = time.Now()
	req.Data.UpdatedAt = time.Now()
	req.Data.Status = "active"

	// Hash password if provided
	if req.Data.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Data.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		req.Data.PasswordHash = string(hashedPassword)
	}

	err := s.repo.Create(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &req.Data, nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*database.User, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[database.User]) (*database.User, error) {
	req.Data.ID = id
	req.Data.UpdatedAt = time.Now()

	// Hash password if provided
	if req.Data.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Data.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		req.Data.PasswordHash = string(hashedPassword)
	}

	err := s.repo.Update(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &req.Data, nil
}

// Delete deletes a user by ID
func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves users with pagination
func (s *userService) List(ctx context.Context, req *ListRequest) ([]*database.User, error) {
	users, err := s.repo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// Count returns the total number of users
func (s *userService) Count(ctx context.Context, req *CountRequest) (int64, error) {
	return s.repo.Count(ctx)
}

// GetByEmail retrieves a user by email
func (s *userService) GetByEmail(ctx context.Context, email string) (*database.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// GetByCompanyID retrieves users by company ID
func (s *userService) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.User, error) {
	return s.repo.GetByCompanyID(ctx, companyID)
}

// GetByRole retrieves users by role
func (s *userService) GetByRole(ctx context.Context, role string) ([]*database.User, error) {
	return s.repo.GetByRole(ctx, role)
}

// GetByStatus retrieves users by status
func (s *userService) GetByStatus(ctx context.Context, status string) ([]*database.User, error) {
	return s.repo.GetByStatus(ctx, status)
}

// ChangePassword changes a user's password
func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return fmt.Errorf("invalid old password: %w", err)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.repo.Update(ctx, user)
}

// GenerateAPIKey generates a new API key for a user
func (s *userService) GenerateAPIKey(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Generate new API key
	apiKey := uuid.New().String()
	expiresAt := time.Now().Add(365 * 24 * time.Hour) // 1 year

	user.APIKeyHash = apiKey
	user.APIKeyExpiresAt = &expiresAt
	user.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to update user: %w", err)
	}

	return apiKey, nil
}

// RevokeAPIKey revokes a user's API key
func (s *userService) RevokeAPIKey(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.APIKeyHash = ""
	user.APIKeyExpiresAt = nil
	user.UpdatedAt = time.Now()

	return s.repo.Update(ctx, user)
}

// UpdateLastLogin updates the user's last login time
func (s *userService) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now

	return s.repo.Update(ctx, user)
}

// Create creates a new company
func (s *companyService) Create(ctx context.Context, req *CreateRequest[database.Company]) (*database.Company, error) {
	company := &req.Data
	company.ID = uuid.New()
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}

// GetByID gets a company by ID
func (s *companyService) GetByID(ctx context.Context, id uuid.UUID) (*database.Company, error) {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &NotFoundError{Resource: "Company", ID: id}
	}
	return company, nil
}

// Update updates a company
func (s *companyService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[database.Company]) (*database.Company, error) {
	// First get the existing company
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &NotFoundError{Resource: "Company", ID: id}
	}

	// Update fields
	company.Name = req.Data.Name
	company.ContactEmail = req.Data.ContactEmail
	company.ContactPhone = req.Data.ContactPhone
	company.Address = req.Data.Address
	company.Status = req.Data.Status
	company.BillingEnabled = req.Data.BillingEnabled
	company.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}

// Delete deletes a company
func (s *companyService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if company exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return &NotFoundError{Resource: "Company", ID: id}
	}

	return s.repo.Delete(ctx, id)
}

// List lists companies
func (s *companyService) List(ctx context.Context, req *ListRequest) ([]*database.Company, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}

// Count counts companies
func (s *companyService) Count(ctx context.Context, req *CountRequest) (int64, error) {
	return s.repo.Count(ctx)
}

// GetByEmail gets a company by email
func (s *companyService) GetByEmail(ctx context.Context, email string) (*database.Company, error) {
	company, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, &NotFoundError{Resource: "Company", ID: uuid.Nil}
	}
	return company, nil
}

// GetByStatus gets companies by status
func (s *companyService) GetByStatus(ctx context.Context, status string) ([]*database.Company, error) {
	return s.repo.GetByStatus(ctx, status)
}

// UpdateStatus updates company status
func (s *companyService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return &NotFoundError{Resource: "Company", ID: id}
	}

	company.Status = status
	company.UpdatedAt = time.Now()

	return s.repo.Update(ctx, company)
}

// EnableBilling enables billing for a company
func (s *companyService) EnableBilling(ctx context.Context, id uuid.UUID) error {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return &NotFoundError{Resource: "Company", ID: id}
	}

	company.BillingEnabled = true
	company.UpdatedAt = time.Now()

	return s.repo.Update(ctx, company)
}

// DisableBilling disables billing for a company
func (s *companyService) DisableBilling(ctx context.Context, id uuid.UUID) error {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return &NotFoundError{Resource: "Company", ID: id}
	}

	company.BillingEnabled = false
	company.UpdatedAt = time.Now()

	return s.repo.Update(ctx, company)
}

// ProjectService defines project-specific operations
type ProjectService interface {
	Service[database.Project]
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.Project, error)
	GetByKeyName(ctx context.Context, keyName string) (*database.Project, error)
	GetByStatus(ctx context.Context, status string) ([]*database.Project, error)
	UpdateLimits(ctx context.Context, id uuid.UUID, dailyLimit, monthlyLimit int) error
	CheckLimit(ctx context.Context, id uuid.UUID) (bool, error)
}

// BotService defines bot-specific operations
type BotService interface {
	Service[database.PlatformBot]
	GetByAppID(ctx context.Context, appID string) (*database.PlatformBot, error)
	GetByStatus(ctx context.Context, status string) ([]*database.PlatformBot, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*database.PlatformBot, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateCapabilities(ctx context.Context, id uuid.UUID, capabilities map[string]any) error
	TestConnection(ctx context.Context, id uuid.UUID) error
}

// ThirdPartyBotService defines third-party bot operations
type ThirdPartyBotService interface {
	Service[database.ThirdPartyBot]
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.ThirdPartyBot, error)
	GetByAppID(ctx context.Context, appID string) (*database.ThirdPartyBot, error)
	GetByStatus(ctx context.Context, status string) ([]*database.ThirdPartyBot, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateAPIKey(ctx context.Context, id uuid.UUID, apiKey string) error
	TestConnection(ctx context.Context, id uuid.UUID) error
}

// DestinationService defines destination-specific operations
type DestinationService interface {
	Service[database.Destination]
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*database.Destination, error)
	GetByBotID(ctx context.Context, botID uuid.UUID) ([]*database.Destination, error)
	GetByStatus(ctx context.Context, status string) ([]*database.Destination, error)
	GetByValidationStatus(ctx context.Context, validationStatus string) ([]*database.Destination, error)
	SearchDestinations(ctx context.Context, query string) ([]*database.Destination, error)
	UpdateTargets(ctx context.Context, id uuid.UUID, targets database.JSONBTargets) (*database.Destination, error)
	ValidateTargets(ctx context.Context, id uuid.UUID) (*database.Destination, error)
}

// NotificationService defines notification-specific operations
type NotificationService interface {
	Service[database.Notification]
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*database.Notification, error)
	GetBySenderID(ctx context.Context, senderID uuid.UUID) ([]*database.Notification, error)
	GetByStatus(ctx context.Context, status string) ([]*database.Notification, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*database.Notification, error)
	SendNotification(ctx context.Context, req *SendNotificationRequest) (*database.Notification, error)
	RetryNotification(ctx context.Context, id uuid.UUID) error
	CancelNotification(ctx context.Context, id uuid.UUID) error
}

// =============================================
// Messages (Bot Installation) Service
// =============================================

// Activity minimal for Bot Framework events
type Activity struct {
	Type       string `json:"type"`
	ServiceURL string `json:"serviceUrl"`
	From       struct {
		ID string `json:"id"`
	} `json:"from"`
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Conversation struct {
		ID string `json:"id"`
	} `json:"conversation"`
	ChannelID   string         `json:"channelId"`
	Locale      string         `json:"locale"`
	ChannelData map[string]any `json:"channelData"`
}

type MessagesService interface {
	HandleActivity(ctx context.Context, act *Activity, rawPayload map[string]any) error
	SendProactiveTest(ctx context.Context, activity map[string]any, text string) error
}

type messagesService struct {
	platformRepo    repositories.PlatformBotRepository
    thirdPartyRepo  repositories.ThirdPartyBotRepository
	installRepo     repositories.BotInstallationRepository
	defaultTenantID string
}

func NewMessagesService(platformRepo repositories.PlatformBotRepository, installRepo repositories.BotInstallationRepository, defaultTenantID string) MessagesService {
    // Backward-compat for callers not yet updated; thirdPartyRepo can be set later if needed
    return &messagesService{platformRepo: platformRepo, installRepo: installRepo, defaultTenantID: defaultTenantID}
}

func (s *messagesService) HandleActivity(ctx context.Context, act *Activity, rawPayload map[string]any) error {
	if act == nil {
		return errors.New("nil activity")
	}
    // Only handle installationUpdate events here
    if strings.TrimSpace(strings.ToLower(act.Type)) != "installationupdate" {
        return nil
    }
	// Extract appId from recipient.id like "28:<APPID>"
	appID := ""
	if rid := act.Recipient.ID; rid != "" {
		if strings.HasPrefix(rid, "28:") && len(rid) > 3 {
			appID = rid[3:]
		}
	}
	if appID == "" {
		return errors.New("cannot determine app_id from recipient.id")
	}

    // Find platform bot by app_id; if missing, log and return nil (no creation here)
    bot, err := s.platformRepo.GetByAppID(ctx, appID)
    if err != nil {
        log.Printf("installationUpdate: platform bot not found for app_id=%s: %v", appID, err)
        return nil
    }

    // Also check if any third-party bot shares this app_id (for observability only)
    if s.thirdPartyRepo != nil {
        if _, err := s.thirdPartyRepo.GetByAppID(ctx, appID); err != nil {
            log.Printf("installationUpdate: third_party_bots no record for app_id=%s (expected provision-managed). err=%v", appID, err)
        }
    }

	// Determine tenant id
	tenantID := s.defaultTenantID
	if t, ok := rawPayload["tenant"].(map[string]any); ok {
		if tid, ok2 := t["id"].(string); ok2 && tid != "" {
			tenantID = tid
		}
	}
	if tenantID == "" {
		if cd, ok := act.ChannelData["tenant"].(map[string]any); ok {
			if tid, ok2 := cd["id"].(string); ok2 && tid != "" {
				tenantID = tid
			}
		}
	}
	if tenantID == "" {
		tenantID = "unknown-tenant"
	}

	// Determine conversation details
	conversationID := act.Conversation.ID
	serviceURL := ""

	// Extract service URL from raw payload
	if su, ok := rawPayload["serviceUrl"].(string); ok && su != "" {
		serviceURL = su
	}

	// Determine conversation type from Teams event
	conversationType := "personal"
	if convType, ok := rawPayload["conversation"].(map[string]any); ok {
		if ct, ok := convType["conversationType"].(string); ok {
			switch ct {
			case "personal":
				conversationType = "personal"
			case "channel":
				conversationType = "channel"
			case "groupChat":
				conversationType = "groupChat"
			}
		}
	}

	// Extract AAD Object ID and names
	aadObjectID := ""
	fromName := ""
	recipientName := ""
	if from, ok := rawPayload["from"].(map[string]any); ok {
		if aad, ok := from["aadObjectId"].(string); ok {
			aadObjectID = aad
		}
		if name, ok := from["name"].(string); ok {
			fromName = name
		}
	}
	if recipient, ok := rawPayload["recipient"].(map[string]any); ok {
		if name, ok := recipient["name"].(string); ok {
			recipientName = name
		}
	}

    // Determine action (add/remove)
    action := "add"
    if a, ok := rawPayload["action"].(string); ok && a != "" {
        action = a
    }

    // Prepare installation entity
	now := time.Now()
	inst := &database.BotInstallation{
		BotID:              bot.ID,
		BotType:            database.BotType("platform"),
		TeamsTenantID:      tenantID,
		ConversationType:   conversationType,
		ConversationID:     conversationID,
		ServiceURL:         serviceURL,
		RecipientID:        act.Recipient.ID,
		RecipientName:      recipientName,
		FromID:             act.From.ID,
		FromName:           fromName,
		FromAADObjectID:    aadObjectID,
        InstallationStatus: "active",
        InstalledAt:        now,
        UninstalledAt:      nil,
		LastActivityAt:     &now,
		Metadata:           rawPayload,
	}
    if strings.EqualFold(action, "remove") || strings.EqualFold(action, "uninstall") {
        inst.InstallationStatus = "uninstalled"
        inst.UninstalledAt = &now
    }

    return s.installRepo.Upsert(ctx, inst)
}

// SendProactiveTest sends a simple proactive text message using provided activity payload (no DB)
func (s *messagesService) SendProactiveTest(ctx context.Context, activity map[string]any, text string) error {
	// Extract required fields
	serviceURL, _ := activity["serviceUrl"].(string)
	if serviceURL == "" {
		return errors.New("serviceUrl missing in activity")
	}
	conv, _ := activity["conversation"].(map[string]any)
	convID, _ := conv["id"].(string)
	if convID == "" {
		return errors.New("conversation.id missing in activity")
	}
	// tenant id from conversation.tenantId or channelData.tenant.id
	tenantID, _ := conv["tenantId"].(string)
	if tenantID == "" {
		if cd, ok := activity["channelData"].(map[string]any); ok {
			if t, ok2 := cd["tenant"].(map[string]any); ok2 {
				if tid, ok3 := t["id"].(string); ok3 {
					tenantID = tid
				}
			}
		}
	}
	if tenantID == "" {
		tenantID = s.defaultTenantID
	}
	recipient, _ := activity["recipient"].(map[string]any)
	recipientID, _ := recipient["id"].(string) // e.g., "28:<APPID>"
	from, _ := activity["from"].(map[string]any)
	fromID, _ := from["id"].(string)
	// In proactive sample: from = botID, recipient = userID
	botID := recipientID
	userID := fromID

	// Acquire token from Bot Framework (using env vars)
	clientID := getEnv("TEAMS_BOT_APP_ID", "")
	clientSecret := getEnv("TEAMS_BOT_APP_PASSWORD", "")
	if clientID == "" || clientSecret == "" {
		return errors.New("TEAMS_BOT_APP_ID/TEAMS_BOT_APP_PASSWORD are required in env")
	}
	token, err := fetchBotFrameworkToken(ctx, tenantID, clientID, clientSecret)
	if err != nil {
		return fmt.Errorf("failed to fetch token: %w", err)
	}

	// Build activity payload
	msg := map[string]any{
		"type":         "message",
		"text":         text,
		"from":         map[string]any{"id": botID},
		"recipient":    map[string]any{"id": userID},
		"conversation": map[string]any{"id": convID},
	}
	body, _ := jsonMarshal(msg)

	// POST to {serviceUrl}/v3/conversations/{conversationId}/activities
	endpoint := strings.TrimRight(serviceURL, "/") + "/v3/conversations/" + url.PathEscape(convID) + "/activities"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("botframework send failed: %s: %s", resp.Status, string(b))
	}
	return nil
}

// fetchBotFrameworkToken performs client credentials to get BF token
func fetchBotFrameworkToken(ctx context.Context, tenantID, clientID, clientSecret string) (string, error) {
	if tenantID == "" {
		tenantID = "botframework.com"
	}
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("scope", "https://api.botframework.com/.default")
	tokenURL := "https://login.microsoftonline.com/" + url.PathEscape(tenantID) + "/oauth2/v2.0/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed: %s: %s", resp.Status, string(b))
	}
	var tr struct {
		AccessToken string `json:"access_token"`
	}
	if err := jsonNewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}
	if tr.AccessToken == "" {
		return "", errors.New("empty access_token")
	}
	return tr.AccessToken, nil
}

// Helpers for local package to avoid importing encoding/json directly at top
func jsonMarshal(v any) ([]byte, error)        { return json.Marshal(v) }
func jsonNewDecoder(r io.Reader) *json.Decoder { return json.NewDecoder(r) }
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// SendNotificationRequest represents a send notification request
type SendNotificationRequest struct {
	ProjectID    uuid.UUID              `json:"project_id" validate:"required"`
	SenderID     *uuid.UUID             `json:"sender_id"`
	MessageType  string                 `json:"message_type" validate:"required,oneof=text file adaptive_card"`
	Content      string                 `json:"content" validate:"required,max=4000"`
	Mentions     []string               `json:"mentions"`
	Attachment   *database.Attachment   `json:"attachment"`
	AdaptiveCard *database.AdaptiveCard `json:"adaptive_card"`
	Priority     string                 `json:"priority" validate:"oneof=low normal high urgent"`
	Metadata     map[string]any         `json:"metadata"`
	Destinations []uuid.UUID            `json:"destinations" validate:"required,min=1"`
}

// BaseService provides common service functionality
type BaseService struct {
	// This will be implemented with actual repository
	// For now, it's a placeholder
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Resource string
	ID       uuid.UUID
}

func (e NotFoundError) Error() string {
	return e.Resource + " not found"
}

// ConflictError represents a conflict error
type ConflictError struct {
	Resource string
	Field    string
	Value    string
}

func (e ConflictError) Error() string {
	return e.Resource + " with " + e.Field + " " + e.Value + " already exists"
}

// projectService implements ProjectService
type projectService struct {
	repo repositories.ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(repo repositories.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

// Create creates a new project
func (s *projectService) Create(ctx context.Context, req *CreateRequest[database.Project]) (*database.Project, error) {
	// Set default values
	req.Data.ID = uuid.New()
	req.Data.CreatedAt = time.Now()
	req.Data.UpdatedAt = time.Now()
	req.Data.Status = "active"

	// Set default priority if not provided
	if req.Data.Priority == "" {
		req.Data.Priority = "normal"
	}

	err := s.repo.Create(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return &req.Data, nil
}

// GetByID retrieves a project by ID
func (s *projectService) GetByID(ctx context.Context, id uuid.UUID) (*database.Project, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates an existing project
func (s *projectService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[database.Project]) (*database.Project, error) {
	req.Data.ID = id
	req.Data.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return &req.Data, nil
}

// Delete deletes a project by ID
func (s *projectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves projects with pagination
func (s *projectService) List(ctx context.Context, req *ListRequest) ([]*database.Project, error) {
	projects, err := s.repo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return projects, nil
}

// Count returns the total number of projects
func (s *projectService) Count(ctx context.Context, req *CountRequest) (int64, error) {
	return s.repo.Count(ctx)
}

// GetByCompanyID retrieves projects by company ID
func (s *projectService) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.Project, error) {
	return s.repo.GetByCompanyID(ctx, companyID)
}

// GetByKeyName retrieves a project by key name
func (s *projectService) GetByKeyName(ctx context.Context, keyName string) (*database.Project, error) {
	return s.repo.GetByKeyName(ctx, keyName)
}

// GetByStatus retrieves projects by status
func (s *projectService) GetByStatus(ctx context.Context, status string) ([]*database.Project, error) {
	return s.repo.GetByStatus(ctx, status)
}

// UpdateLimits updates project limits
func (s *projectService) UpdateLimits(ctx context.Context, id uuid.UUID, dailyLimit, monthlyLimit int) error {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Update limits
	project.DailyLimit = dailyLimit
	project.MonthlyLimit = monthlyLimit
	project.UpdatedAt = time.Now()

	return s.repo.Update(ctx, project)
}

// CheckLimit checks if project is within limits
func (s *projectService) CheckLimit(ctx context.Context, id uuid.UUID) (bool, error) {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to get project: %w", err)
	}

	// TODO: Implement actual limit checking logic
	// This would typically check usage records against the limits
	return true, nil
}

// platformBotService implements BotService
type platformBotService struct {
	repo repositories.PlatformBotRepository
}

// NewPlatformBotService creates a new platform bot service
func NewPlatformBotService(repo repositories.PlatformBotRepository) BotService {
	return &platformBotService{repo: repo}
}

// Create creates a new platform bot
func (s *platformBotService) Create(ctx context.Context, req *CreateRequest[database.PlatformBot]) (*database.PlatformBot, error) {
	// Set default values
	req.Data.ID = uuid.New()
	req.Data.CreatedAt = time.Now()
	req.Data.UpdatedAt = time.Now()
	req.Data.Status = database.BotStatusActive

	err := s.repo.Create(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform bot: %w", err)
	}

	return &req.Data, nil
}

// GetByID retrieves a platform bot by ID
func (s *platformBotService) GetByID(ctx context.Context, id uuid.UUID) (*database.PlatformBot, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates an existing platform bot
func (s *platformBotService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[database.PlatformBot]) (*database.PlatformBot, error) {
	req.Data.ID = id
	req.Data.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to update platform bot: %w", err)
	}

	return &req.Data, nil
}

// Delete deletes a platform bot by ID
func (s *platformBotService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves platform bots with pagination
func (s *platformBotService) List(ctx context.Context, req *ListRequest) ([]*database.PlatformBot, error) {
	bots, err := s.repo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list platform bots: %w", err)
	}

	return bots, nil
}

// Count returns the total number of platform bots
func (s *platformBotService) Count(ctx context.Context, req *CountRequest) (int64, error) {
	return s.repo.Count(ctx)
}

// GetByAppID retrieves a platform bot by app ID
func (s *platformBotService) GetByAppID(ctx context.Context, appID string) (*database.PlatformBot, error) {
	return s.repo.GetByAppID(ctx, appID)
}

// GetByStatus retrieves platform bots by status
func (s *platformBotService) GetByStatus(ctx context.Context, status string) ([]*database.PlatformBot, error) {
	return s.repo.GetByStatus(ctx, status)
}

// GetByCompanyID retrieves platform bots by company ID
func (s *platformBotService) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.PlatformBot, error) {
	return s.repo.GetByCompanyID(ctx, companyID)
}

// GetByTenantID retrieves platform bots by tenant ID
func (s *platformBotService) GetByTenantID(ctx context.Context, tenantID string) ([]*database.PlatformBot, error) {
	return s.repo.GetByTenantID(ctx, tenantID)
}

// UpdateStatus updates platform bot status
func (s *platformBotService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	bot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get platform bot: %w", err)
	}

	bot.Status = database.BotStatus(status)
	bot.UpdatedAt = time.Now()

	return s.repo.Update(ctx, bot)
}

// UpdateCapabilities updates platform bot capabilities
func (s *platformBotService) UpdateCapabilities(ctx context.Context, id uuid.UUID, capabilities map[string]interface{}) error {
	bot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get platform bot: %w", err)
	}

	bot.Capabilities = capabilities
	bot.UpdatedAt = time.Now()

	return s.repo.Update(ctx, bot)
}

// TestConnection tests platform bot connection
func (s *platformBotService) TestConnection(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get platform bot: %w", err)
	}

	// TODO: Implement actual connection testing logic
	// This would typically make a test API call to the bot's webhook URL
	return nil
}

// thirdPartyBotService implements ThirdPartyBotService
type thirdPartyBotService struct {
	repo repositories.ThirdPartyBotRepository
}

// NewThirdPartyBotService creates a new third-party bot service
func NewThirdPartyBotService(repo repositories.ThirdPartyBotRepository) ThirdPartyBotService {
	return &thirdPartyBotService{repo: repo}
}

// Create creates a new third-party bot
func (s *thirdPartyBotService) Create(ctx context.Context, req *CreateRequest[database.ThirdPartyBot]) (*database.ThirdPartyBot, error) {
	// Set default values
	req.Data.ID = uuid.New()
	req.Data.CreatedAt = time.Now()
	req.Data.UpdatedAt = time.Now()
	req.Data.Status = database.BotStatusActive

	err := s.repo.Create(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to create third-party bot: %w", err)
	}

	return &req.Data, nil
}

// GetByID retrieves a third-party bot by ID
func (s *thirdPartyBotService) GetByID(ctx context.Context, id uuid.UUID) (*database.ThirdPartyBot, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates an existing third-party bot
func (s *thirdPartyBotService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[database.ThirdPartyBot]) (*database.ThirdPartyBot, error) {
	req.Data.ID = id
	req.Data.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to update third-party bot: %w", err)
	}

	return &req.Data, nil
}

// Delete deletes a third-party bot by ID
func (s *thirdPartyBotService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves third-party bots with pagination
func (s *thirdPartyBotService) List(ctx context.Context, req *ListRequest) ([]*database.ThirdPartyBot, error) {
	bots, err := s.repo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list third-party bots: %w", err)
	}

	return bots, nil
}

// Count returns the total number of third-party bots
func (s *thirdPartyBotService) Count(ctx context.Context, req *CountRequest) (int64, error) {
	return s.repo.Count(ctx)
}

// GetByAppID retrieves a third-party bot by app ID
func (s *thirdPartyBotService) GetByAppID(ctx context.Context, appID string) (*database.ThirdPartyBot, error) {
	return s.repo.GetByAppID(ctx, appID)
}

// GetByStatus retrieves third-party bots by status
func (s *thirdPartyBotService) GetByStatus(ctx context.Context, status string) ([]*database.ThirdPartyBot, error) {
	return s.repo.GetByStatus(ctx, status)
}

// GetByCompanyID retrieves third-party bots by company ID
func (s *thirdPartyBotService) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.ThirdPartyBot, error) {
	return s.repo.GetByCompanyID(ctx, companyID)
}

// UpdateStatus updates third-party bot status
func (s *thirdPartyBotService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	bot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get third-party bot: %w", err)
	}

	bot.Status = database.BotStatus(status)
	bot.UpdatedAt = time.Now()

	return s.repo.Update(ctx, bot)
}

// UpdateAPIKey updates third-party bot API key
func (s *thirdPartyBotService) UpdateAPIKey(ctx context.Context, id uuid.UUID, apiKey string) error {
	bot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get third-party bot: %w", err)
	}

	// Hash the API key
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash API key: %w", err)
	}

	bot.APIKeyHash = string(hashedKey)
	bot.UpdatedAt = time.Now()

	return s.repo.Update(ctx, bot)
}

// TestConnection tests third-party bot connection
func (s *thirdPartyBotService) TestConnection(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get third-party bot: %w", err)
	}

	// TODO: Implement actual connection testing logic
	// This would typically make a test API call to the bot's webhook URL
	return nil
}

// =============================================
// Destination Service Implementation
// =============================================

// NewDestinationService creates a new destination service
func NewDestinationService(repo repositories.DestinationRepository) DestinationService {
	return &destinationService{repo: repo}
}

type destinationService struct {
	repo repositories.DestinationRepository
}

// Create creates a new destination
func (s *destinationService) Create(ctx context.Context, req *CreateRequest[database.Destination]) (*database.Destination, error) {
	// Set default values
	destination := req.Data
	destination.ID = uuid.New()
	destination.CreatedAt = time.Now()
	destination.UpdatedAt = time.Now()
	// Only set default status if not already set
	if destination.Status == "" {
		destination.Status = "active"
	}
	if destination.ValidationStatus == "" {
		destination.ValidationStatus = "pending"
	}

	err := s.repo.Create(ctx, &destination)
	if err != nil {
		return nil, err
	}

	return &destination, nil
}

// GetByID retrieves a destination by ID
func (s *destinationService) GetByID(ctx context.Context, id uuid.UUID) (*database.Destination, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates a destination
func (s *destinationService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[database.Destination]) (*database.Destination, error) {
	destination := req.Data
	destination.ID = id
	destination.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, &destination)
	if err != nil {
		return nil, err
	}

	return &destination, nil
}

// Delete deletes a destination by ID
func (s *destinationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves destinations with pagination
func (s *destinationService) List(ctx context.Context, req *ListRequest) ([]*database.Destination, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}

// Count returns the total number of destinations
func (s *destinationService) Count(ctx context.Context, req *CountRequest) (int64, error) {
	return s.repo.Count(ctx)
}

// GetByProjectID retrieves destinations by project ID
func (s *destinationService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*database.Destination, error) {
	return s.repo.GetByProjectID(ctx, projectID)
}

// GetByBotID retrieves destinations by bot ID
func (s *destinationService) GetByBotID(ctx context.Context, botID uuid.UUID) ([]*database.Destination, error) {
	return s.repo.GetByBotID(ctx, botID)
}

// GetByStatus retrieves destinations by status
func (s *destinationService) GetByStatus(ctx context.Context, status string) ([]*database.Destination, error) {
	return s.repo.GetByStatus(ctx, status)
}

// GetByValidationStatus retrieves destinations by validation status
func (s *destinationService) GetByValidationStatus(ctx context.Context, validationStatus string) ([]*database.Destination, error) {
	return s.repo.GetByValidationStatus(ctx, validationStatus)
}

// SearchDestinations searches destinations by name or description
func (s *destinationService) SearchDestinations(ctx context.Context, query string) ([]*database.Destination, error) {
	return s.repo.SearchDestinations(ctx, query)
}

// UpdateTargets updates destination targets
func (s *destinationService) UpdateTargets(ctx context.Context, id uuid.UUID, targets database.JSONBTargets) (*database.Destination, error) {
	destination, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination: %w", err)
	}

	destination.Targets = targets
	destination.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, destination)
	if err != nil {
		return nil, err
	}

	return destination, nil
}

// ValidateTargets validates destination targets
func (s *destinationService) ValidateTargets(ctx context.Context, id uuid.UUID) (*database.Destination, error) {
	destination, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination: %w", err)
	}

	// TODO: Implement actual target validation logic
	// This would typically validate Teams channels, chat groups, or users
	destination.ValidationStatus = "validated"
	destination.LastValidatedAt = &time.Time{}
	*destination.LastValidatedAt = time.Now()
	destination.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, destination)
	if err != nil {
		return nil, err
	}

	return destination, nil
}

// =============================================
// Notification Service Implementation
// =============================================

// NewNotificationService creates a new notification service
func NewNotificationService(repo repositories.NotificationRepository, broadcaster BroadcastService) NotificationService {
	return &notificationService{repo: repo, broadcaster: broadcaster}
}

type notificationService struct {
	repo        repositories.NotificationRepository
	broadcaster BroadcastService
}

// Create creates a new notification
func (s *notificationService) Create(ctx context.Context, req *CreateRequest[database.Notification]) (*database.Notification, error) {
	notification := req.Data
	notification.ID = uuid.New()
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()
	if notification.Status == "" {
		notification.Status = "pending"
	}

	err := s.repo.Create(ctx, &notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// GetByID retrieves a notification by ID
func (s *notificationService) GetByID(ctx context.Context, id uuid.UUID) (*database.Notification, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates a notification
func (s *notificationService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[database.Notification]) (*database.Notification, error) {
	notification := req.Data
	notification.ID = id
	notification.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, &notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// Delete deletes a notification by ID
func (s *notificationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves notifications with pagination
func (s *notificationService) List(ctx context.Context, req *ListRequest) ([]*database.Notification, error) {
	return s.repo.List(ctx, req.Limit, req.Offset)
}

// Count returns the total number of notifications
func (s *notificationService) Count(ctx context.Context, req *CountRequest) (int64, error) {
	return s.repo.Count(ctx)
}

// GetByProjectID retrieves notifications by project ID
func (s *notificationService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*database.Notification, error) {
	return s.repo.GetByProjectID(ctx, projectID)
}

// GetBySenderID retrieves notifications by sender ID
func (s *notificationService) GetBySenderID(ctx context.Context, senderID uuid.UUID) ([]*database.Notification, error) {
	return s.repo.GetBySenderID(ctx, senderID)
}

// GetByStatus retrieves notifications by status
func (s *notificationService) GetByStatus(ctx context.Context, status string) ([]*database.Notification, error) {
	return s.repo.GetByStatus(ctx, status)
}

// GetByDateRange retrieves notifications by date range
func (s *notificationService) GetByDateRange(ctx context.Context, start, end time.Time) ([]*database.Notification, error) {
	return s.repo.GetByDateRange(ctx, start, end)
}

// SendNotification sends a notification
func (s *notificationService) SendNotification(ctx context.Context, req *SendNotificationRequest) (*database.Notification, error) {
	// Convert Attachment and AdaptiveCard to JSONBNullableObject
	var attachment *database.JSONBNullableObject
	if req.Attachment != nil {
		// Convert Attachment struct to map[string]any
		attachmentData := map[string]any{
			"file_name": req.Attachment.FileName,
			"file_url":  req.Attachment.FileURL,
			"file_size": req.Attachment.FileSize,
			"mime_type": req.Attachment.MimeType,
		}
		attachment = (*database.JSONBNullableObject)(&attachmentData)
	}

	var adaptiveCard *database.JSONBNullableObject
	if req.AdaptiveCard != nil {
		// Convert AdaptiveCard struct to map[string]any
		adaptiveCardData := map[string]any{
			"type":    req.AdaptiveCard.Type,
			"version": req.AdaptiveCard.Version,
			"body":    req.AdaptiveCard.Body,
			"actions": req.AdaptiveCard.Actions,
		}
		adaptiveCard = (*database.JSONBNullableObject)(&adaptiveCardData)
	}

	notification := &database.Notification{
		ProjectID:    req.ProjectID,
		SenderID:     req.SenderID,
		MessageType:  req.MessageType,
		Content:      req.Content,
		Mentions:     database.JSONBStringArray(req.Mentions),
		Attachment:   attachment,
		AdaptiveCard: adaptiveCard,
		Priority:     req.Priority,
		Status:       "sent",
		Metadata:     database.JSONBObject(req.Metadata),
	}

	// Send to destinations first
	if s.broadcaster != nil && len(req.Destinations) > 0 {
		if _, err := s.broadcaster.SendToDestinations(ctx, req.Destinations, req.Content); err != nil {
			return nil, fmt.Errorf("failed to send message to destinations: %w", err)
		}
	}

	// After successful send, persist the record
	notification.ID = uuid.New()
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	err := s.repo.Create(ctx, notification)
	if err != nil {
		return nil, err
	}

	// TODO: Queue notification for actual sending
	// This would typically involve:
	// 1. Finding destinations for the project
	// 2. Creating notification_destinations records
	// 3. Queuing the notification for processing

	return notification, nil
}

// RetryNotification retries a failed notification
func (s *notificationService) RetryNotification(ctx context.Context, id uuid.UUID) error {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if notification.Status != "failed" {
		return fmt.Errorf("notification is not in failed status")
	}

	// Reset status and retry
	notification.Status = "pending"
	notification.ErrorMessage = ""
	notification.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, notification)
	if err != nil {
		return err
	}

	// TODO: Queue notification for retry processing

	return nil
}

// CancelNotification cancels a notification
func (s *notificationService) CancelNotification(ctx context.Context, id uuid.UUID) error {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if notification.Status == "sent" {
		return fmt.Errorf("cannot cancel already sent notification")
	}

	notification.Status = "cancelled"
	notification.UpdatedAt = time.Now()

	return s.repo.Update(ctx, notification)
}
