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

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/repositories"
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
	Service[models.Company]
	GetByEmail(ctx context.Context, email string) (*models.Company, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Company, error)
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
	Service[models.User]
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.User, error)
	GetByRole(ctx context.Context, role string) ([]*models.User, error)
	GetByStatus(ctx context.Context, status string) ([]*models.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
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
func (s *userService) Create(ctx context.Context, req *CreateRequest[models.User]) (*models.User, error) {
	// Set default values
	req.Data.ID = uuid.New()
	req.Data.CreatedAt = time.Now()
	req.Data.UpdatedAt = time.Now()
	req.Data.Status = "active"

	// Hash password if provided
	if req.Data.PasswordHash != nil && *req.Data.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Data.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashedStr := string(hashedPassword)
		req.Data.PasswordHash = &hashedStr
	}

	err := s.repo.Create(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &req.Data, nil
}

// GetByID retrieves a user by ID
func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates an existing user
func (s *userService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[models.User]) (*models.User, error) {
	req.Data.ID = id
	req.Data.UpdatedAt = time.Now()

	// Hash password if provided
	if req.Data.PasswordHash != nil && *req.Data.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Data.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashedStr := string(hashedPassword)
		req.Data.PasswordHash = &hashedStr
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
func (s *userService) List(ctx context.Context, req *ListRequest) ([]*models.User, error) {
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
func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// GetByCompanyID retrieves users by company ID
func (s *userService) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.User, error) {
	return s.repo.GetByCompanyID(ctx, companyID)
}

// GetByRole retrieves users by role
func (s *userService) GetByRole(ctx context.Context, role string) ([]*models.User, error) {
	return s.repo.GetByRole(ctx, role)
}

// GetByStatus retrieves users by status
func (s *userService) GetByStatus(ctx context.Context, status string) ([]*models.User, error) {
	return s.repo.GetByStatus(ctx, status)
}

// ChangePassword changes a user's password
func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify old password
	if user.PasswordHash == nil {
		return fmt.Errorf("user has no password set")
	}
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return fmt.Errorf("invalid old password: %w", err)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	hashedStr := string(hashedPassword)
	user.PasswordHash = &hashedStr
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
func (s *companyService) Create(ctx context.Context, req *CreateRequest[models.Company]) (*models.Company, error) {
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
func (s *companyService) GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &NotFoundError{Resource: "Company", ID: id}
	}
	return company, nil
}

// Update updates a company
func (s *companyService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[models.Company]) (*models.Company, error) {
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
func (s *companyService) List(ctx context.Context, req *ListRequest) ([]*models.Company, error) {
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
func (s *companyService) GetByEmail(ctx context.Context, email string) (*models.Company, error) {
	company, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, &NotFoundError{Resource: "Company", ID: uuid.Nil}
	}
	return company, nil
}

// GetByStatus gets companies by status
func (s *companyService) GetByStatus(ctx context.Context, status string) ([]*models.Company, error) {
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
	Service[models.Project]
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.Project, error)
	GetByKeyName(ctx context.Context, keyName string) (*models.Project, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Project, error)
	UpdateLimits(ctx context.Context, id uuid.UUID, dailyLimit, monthlyLimit int) error
	CheckLimit(ctx context.Context, id uuid.UUID) (bool, error)
}

// BotService defines bot-specific operations
type TeamsBotService interface {
	Service[models.TeamsBot]
	GetByAppID(ctx context.Context, appID string) (*models.TeamsBot, error)
	GetByStatus(ctx context.Context, status string) ([]*models.TeamsBot, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*models.TeamsBot, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateCapabilities(ctx context.Context, id uuid.UUID, capabilities map[string]any) error
	TestConnection(ctx context.Context, id uuid.UUID) error
}

// DestinationService defines destination-specific operations
type DestinationService interface {
	Service[models.Destination]
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Destination, error)
	GetByBotID(ctx context.Context, botID uuid.UUID) ([]*models.Destination, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Destination, error)
	GetByValidationStatus(ctx context.Context, validationStatus string) ([]*models.Destination, error)
	SearchDestinations(ctx context.Context, query string) ([]*models.Destination, error)
	UpdateTargets(ctx context.Context, id uuid.UUID, targets models.JSONBTargets) (*models.Destination, error)
	ValidateTargets(ctx context.Context, id uuid.UUID) (*models.Destination, error)
}

// NotificationService defines notification-specific operations
type NotificationService interface {
	Service[models.Notification]
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Notification, error)
	GetBySenderID(ctx context.Context, senderID uuid.UUID) ([]*models.Notification, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Notification, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Notification, error)
	SendNotification(ctx context.Context, req *SendNotificationRequest) (*models.Notification, error)
	RetryNotification(ctx context.Context, id uuid.UUID) error
	CancelNotification(ctx context.Context, id uuid.UUID) error
}

// =============================================
// Messages (Bot Installation) Service
// =============================================

// Activity represents Bot Framework Activity payload
type Activity struct {
	ID             string           `json:"id"`
	Type           string           `json:"type"`
	Name           string           `json:"name,omitempty"`
	ServiceURL     string           `json:"serviceUrl"`
	From           ChannelAccount   `json:"from"`
	Recipient      ChannelAccount   `json:"recipient"`
	Conversation   Conversation     `json:"conversation"`
	ChannelID      string           `json:"channelId"`
	Locale         string           `json:"locale,omitempty"`
	LocalTimestamp string           `json:"localTimestamp,omitempty"`
	LocalTimezone  string           `json:"localTimezone,omitempty"`
	Timestamp      string           `json:"timestamp,omitempty"`
	Text           string           `json:"text,omitempty"`
	TextFormat     string           `json:"textFormat,omitempty"`
	ChannelData    map[string]any   `json:"channelData,omitempty"`
	Attachments    []Attachment     `json:"attachments,omitempty"`
	Entities       []map[string]any `json:"entities,omitempty"`
	ReplyToID      string           `json:"replyToId,omitempty"`
	Value          json.RawMessage  `json:"value,omitempty"`
}

type Conversation struct {
	ID               string `json:"id,omitempty"`
	ConversationType string `json:"conversationType,omitempty"`
	TenantID         string `json:"tenantId,omitempty"`
}
type ChannelAccount struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// OAuth Card 內容
type OAuthCard struct {
	Text           string       `json:"text,omitempty"`
	ConnectionName string       `json:"connectionName"`
	Buttons        []CardAction `json:"buttons,omitempty"`
}

// OAuth Card 按鈕
type CardAction struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Value string `json:"value,omitempty"`
}

type Attachment struct {
	Content     interface{} `json:"content"`
	ContentType string      `json:"contentType"`
}

// Token Response（從 OAuth 回調中取得）
type TokenResponse struct {
	Token string `json:"token"`
}
type MessagesService interface {
	HandleActivity(ctx context.Context, act *Activity, rawPayload map[string]any) error
	SendProactiveTest(ctx context.Context, activity map[string]any, text string) error
}

type messagesService struct {
	teamsBotRepo    repositories.TeamsBotRepository
	installRepo     repositories.BotInstallationRepository
	defaultTenantID string
}

func NewMessagesService(teamsBotRepo repositories.TeamsBotRepository, installRepo repositories.BotInstallationRepository, defaultTenantID string) MessagesService {
	return &messagesService{teamsBotRepo: teamsBotRepo, installRepo: installRepo, defaultTenantID: defaultTenantID}
}

func (s *messagesService) HandleActivity(ctx context.Context, act *Activity, rawPayload map[string]any) error {
	if act == nil {
		return errors.New("nil activity")
	}
	activityType := strings.ToLower(strings.TrimSpace(act.Type))
	switch activityType {
	case "installationupdate":
		return s.handleInstallationUpdate(ctx, act, rawPayload)
	case "message":
		return s.handleMessage(ctx, act, rawPayload)
	case "conversationupdate":
		return s.handleConversationUpdate(ctx, act, rawPayload)
	case "event":
		return s.handleEvent(ctx, act, rawPayload)
	default:
		log.Printf("messages: unhandled activity type=%s", activityType)
		return nil
	}
}

// handleInstallationUpdate processes Teams installation add/remove events
func (s *messagesService) handleInstallationUpdate(ctx context.Context, act *Activity, rawPayload map[string]any) error {
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

	// Resolve bot presence: proceed if either platform or third-party bot exists
	var bot *models.TeamsBot
	var botErr error
	bot, botErr = s.teamsBotRepo.GetByAppID(ctx, appID)

	if bot == nil || botErr != nil {
		log.Printf("installationUpdate: teams bot not found for app_id=%s, skip upsert. err=%v", appID, botErr)
		return nil
	}
	// Determine tenant id
	tenantID := s.defaultTenantID
	if t, ok := rawPayload["tenant"].(map[string]any); ok {
		if tid, ok2 := t["id"].(string); ok2 && tid != "" {
			tenantID = tid
		}
	}
	// Try channelData.tenant.id
	if tenantID == "" {
		if cd, ok := act.ChannelData["tenant"].(map[string]any); ok {
			if tid, ok2 := cd["id"].(string); ok2 && tid != "" {
				tenantID = tid
			}
		}
	}
	// Try conversation.tenantId
	if tenantID == "" {
		if convAny, ok := rawPayload["conversation"].(map[string]any); ok {
			if tid, ok2 := convAny["tenantId"].(string); ok2 && tid != "" {
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

	// Query additional information based on conversation type
	email := ""
	descriptionName := ""

	switch conversationType {
	case "personal":
		// For personal chat, get user email from Graph API
		if aadObjectID != "" {
			email = s.getUserEmailFromGraphAPI(ctx, aadObjectID, tenantID)
		}
		descriptionName = fromName // Use from name as description for personal

	case "groupChat":
		// For group chat, get topic from conversation
		if conv, ok := rawPayload["conversation"].(map[string]any); ok {
			if topic, ok := conv["name"].(string); ok {
				descriptionName = topic
			}
		}

	case "channel":
		// For channel, get team name from channel data
		if channelData, ok := rawPayload["channelData"].(map[string]any); ok {
			if team, ok := channelData["team"].(map[string]any); ok {
				if teamName, ok := team["name"].(string); ok {
					descriptionName = teamName
				}
			}
		}
	}

	// Determine action (add/remove)
	action := "add"
	if a, ok := rawPayload["action"].(string); ok && a != "" {
		action = a
	}

	// Prepare installation entity
	now := time.Now()
	inst := &models.BotInstallation{
		BotID:              bot.ID,
		BotType:            models.BotType("platform"),
		TeamsTenantID:      tenantID,
		ConversationType:   conversationType,
		ConversationID:     conversationID,
		ServiceURL:         serviceURL,
		RecipientID:        act.Recipient.ID,
		RecipientName:      recipientName,
		FromID:             act.From.ID,
		FromName:           fromName,
		FromAADObjectID:    aadObjectID,
		Email:              email,
		DescriptionName:    descriptionName,
		InstallationStatus: "active",
		InstalledAt:        now,
		UninstalledAt:      nil,
		LastActivityAt:     &now,
		Metadata:           rawPayload,
	}
	if strings.EqualFold(action, "remove") || strings.EqualFold(action, "uninstall") {
		// Prefer marking existing row as uninstalled to preserve original installed_at
		// We don't necessarily know the ID here (Upsert assigns), so try update by unique keys via Upsert fallback
		inst.InstallationStatus = "uninstalled"
		inst.UninstalledAt = &now
		return s.installRepo.Upsert(ctx, inst)
	}
	return s.installRepo.Upsert(ctx, inst)
}

// handleMessage processes regular chat messages (stub for future use)
func (s *messagesService) handleMessage(ctx context.Context, act *Activity, rawPayload map[string]any) error {
	// Log basic message info only
	convID := ""
	if act != nil {
		convID = act.Conversation.ID
	}
	log.Printf("messages: received message event. conversation_id=%s channel_id=%s", convID, act.ChannelID)
	if strings.Contains(act.Text, "email") {
		reply := createOAuthCardActivity(act)

		// 透過 Bot Connector REST API 回覆
		if err := sendReply(act.ServiceURL, act.Conversation.ID, act.ID, reply); err != nil {
			log.Printf("發送回覆失敗: %v", err)

			return errors.New("send Outh Reply failed")
		}
	}
	return nil
}

// 透過 Bot Connector REST API 發送回覆
func sendReply(serviceURL, conversationID, activityID string, reply Activity) error {
	// Bot Connector REST API 端點
	// POST {serviceUrl}/v3/conversations/{conversationId}/activities/{activityId}
	url := fmt.Sprintf("%s/v3/conversations/%s/activities/%s", serviceURL, conversationID, activityID)

	// 建立 HTTP 請求
	client := &http.Client{Timeout: 10 * time.Second}

	jsonData, err := json.Marshal(reply)
	if err != nil {
		return fmt.Errorf("序列化回覆失敗: %w", err)
	}

	// Debug: log the actual JSON being sent
	log.Printf("Sending reply to Bot Framework:\n%s", string(jsonData))

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("建立請求失敗: %w", err)
	}

	clientID := getEnv("TEAMS_BOT_APP_ID", "")
	clientSecret := getEnv("TEAMS_BOT_APP_PASSWORD", "")
	if clientID == "" || clientSecret == "" {
		return errors.New("TEAMS_BOT_APP_ID/TEAMS_BOT_APP_PASSWORD are required in env")
	}
	token, err := fetchBotFrameworkToken(context.Background(), "051cece0-e4dc-4aed-b471-bf29824e1ee6", clientID, clientSecret)

	if err != nil {
		return fmt.Errorf("failed to fetch token: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// 注意：實際應該使用 Bot Framework 的認證機制
	// 這裡簡化處理，實際應該加上 JWT 認證
	// 參考: https://learn.microsoft.com/en-us/azure/bot-service/rest-api/bot-framework-rest-connector-authentication

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("發送請求失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API 回應錯誤: %d", resp.StatusCode)
	}

	return nil
}

// 處理事件類型的 Activity（OAuth 回調）
func (s *messagesService) handleEvent(ctx context.Context, act *Activity, rawPayload map[string]any) error {
	// 檢查是否為 OAuth Token 回調
	if act.Name == "tokens/response" {
		log.Printf("收到 OAuth Token 回調")

		var tokenResp TokenResponse
		if err := json.Unmarshal(act.Value, &tokenResp); err != nil {
			log.Printf("解析 Token Response 失敗: %v", err)
			return errors.New(fmt.Sprintf("解析 Token Response 失敗: %v", err))

		}

		if tokenResp.Token == "" {
			return errors.New("Token 為空")
		}

		log.Printf("Token: %s", tokenResp.Token)
		// 使用 Token 查詢郵件
		// emails, err := getTodayEmails(tokenResp.Token)
		// if err != nil {
		// 	log.Printf("查詢郵件失敗: %v", err)
		// 	// 回覆錯誤訊息
		// 	reply := Activity{
		// 		Type:         "message",
		// 		From:         activity.Recipient,
		// 		Recipient:    activity.From,
		// 		Conversation: activity.Conversation,
		// 		ServiceURL:   activity.ServiceURL,
		// 		ReplyToID:    activity.ID,
		// 		Text:         "查詢郵件時發生錯誤，請稍後再試。",
		// 	}
		// 	sendReply(activity.ServiceURL, activity.Conversation.ID, activity.ID, reply)
		// 	c.JSON(http.StatusOK, gin.H{})
		// 	return
		// }

		// // 整理郵件結果並回覆
		// emailSummary := formatEmailSummary(emails)
		// reply := Activity{
		// 	Type:         "message",
		// 	From:         activity.Recipient,
		// 	Recipient:    activity.From,
		// 	Conversation: activity.Conversation,
		// 	ServiceURL:   activity.ServiceURL,
		// 	ReplyToID:    activity.ID,
		// 	Text:         emailSummary,
		// }

		// if err := sendReply(activity.ServiceURL, activity.Conversation.ID, activity.ID, reply); err != nil {
		// 	log.Printf("發送回覆失敗: %v", err)
		// }

		// c.JSON(http.StatusOK, gin.H{})

	}

	return nil
}

func createOAuthCardActivity(act *Activity) Activity {
	oauthCard := OAuthCard{
		Text:           "請先登入，讓我幫你查詢今天的 email。",
		ConnectionName: "lab-test-teams-notify-login",
	}

	attachment := Attachment{
		ContentType: "application/vnd.microsoft.card.oauth",
		Content:     oauthCard,
	}

	reply := Activity{
		Type:         "message",
		ChannelID:    act.ChannelID,
		ServiceURL:   act.ServiceURL,
		Conversation: act.Conversation,
		ReplyToID:    act.ID,
		Attachments:  []Attachment{attachment},
	}

	// Copy From/Recipient (swap them for reply)
	reply.From = act.Recipient
	reply.Recipient = act.From

	// Debug: log the OAuth card payload
	if jsonData, err := json.MarshalIndent(reply, "", "  "); err == nil {
		log.Printf("OAuth Card Activity payload:\n%s", string(jsonData))
	}

	return reply
}

// handleConversationUpdate processes member added/removed, etc. (stub)
func (s *messagesService) handleConversationUpdate(ctx context.Context, act *Activity, rawPayload map[string]any) error {
	// Intentionally no-op for now
	return nil
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
	ProjectID   uuid.UUID        `json:"project_id" validate:"required"`
	SenderID    *uuid.UUID       `json:"sender_id"`
	MessageType string           `json:"message_type" validate:"required,oneof=text file adaptive_card"`
	Content     string           `json:"content" validate:"required,max=4000"`
	Mentions    []string         `json:"mentions"`
	Attachments []map[string]any `json:"attachments"`
	Priority    string           `json:"priority" validate:"oneof=low normal high"`
	Metadata    map[string]any   `json:"metadata"`
	Targets     []string         `json:"targets" validate:"omitempty,min=1"`
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
func (s *projectService) Create(ctx context.Context, req *CreateRequest[models.Project]) (*models.Project, error) {
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
func (s *projectService) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates an existing project
func (s *projectService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[models.Project]) (*models.Project, error) {
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
func (s *projectService) List(ctx context.Context, req *ListRequest) ([]*models.Project, error) {
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
func (s *projectService) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.Project, error) {
	return s.repo.GetByCompanyID(ctx, companyID)
}

// GetByKeyName retrieves a project by key name
func (s *projectService) GetByKeyName(ctx context.Context, keyName string) (*models.Project, error) {
	return s.repo.GetByKeyName(ctx, keyName)
}

// GetByStatus retrieves projects by status
func (s *projectService) GetByStatus(ctx context.Context, status string) ([]*models.Project, error) {
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

// teamsBotService implements BotService
type teamsBotService struct {
	repo repositories.TeamsBotRepository
}

// NewteamsBotService creates a new platform bot service
func NewteamsBotService(repo repositories.TeamsBotRepository) TeamsBotService {
	return &teamsBotService{repo: repo}
}

// Create creates a new platform bot
func (s *teamsBotService) Create(ctx context.Context, req *CreateRequest[models.TeamsBot]) (*models.TeamsBot, error) {
	// Set default values
	req.Data.ID = uuid.New()
	req.Data.CreatedAt = time.Now()
	req.Data.UpdatedAt = time.Now()
	status := models.BotStatusActive
	req.Data.Status = &status

	err := s.repo.Create(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform bot: %w", err)
	}

	return &req.Data, nil
}

// GetByID retrieves a platform bot by ID
func (s *teamsBotService) GetByID(ctx context.Context, id uuid.UUID) (*models.TeamsBot, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates an existing platform bot
func (s *teamsBotService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[models.TeamsBot]) (*models.TeamsBot, error) {
	req.Data.ID = id
	req.Data.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, &req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to update platform bot: %w", err)
	}

	return &req.Data, nil
}

// Delete deletes a platform bot by ID
func (s *teamsBotService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves platform bots with pagination
func (s *teamsBotService) List(ctx context.Context, req *ListRequest) ([]*models.TeamsBot, error) {
	bots, err := s.repo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list platform bots: %w", err)
	}

	return bots, nil
}

// Count returns the total number of platform bots
func (s *teamsBotService) Count(ctx context.Context, req *CountRequest) (int64, error) {
	return s.repo.Count(ctx)
}

// GetByAppID retrieves a platform bot by app ID
func (s *teamsBotService) GetByAppID(ctx context.Context, appID string) (*models.TeamsBot, error) {
	return s.repo.GetByAppID(ctx, appID)
}

// GetByStatus retrieves platform bots by status
func (s *teamsBotService) GetByStatus(ctx context.Context, status string) ([]*models.TeamsBot, error) {
	return s.repo.GetByStatus(ctx, status)
}

// GetByCompanyID retrieves platform bots by company ID
func (s *teamsBotService) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.TeamsBot, error) {
	return s.repo.GetByCompanyID(ctx, companyID)
}

// GetByTenantID retrieves platform bots by tenant ID
func (s *teamsBotService) GetByTenantID(ctx context.Context, tenantID string) ([]*models.TeamsBot, error) {
	return s.repo.GetByTenantID(ctx, tenantID)
}

// UpdateStatus updates platform bot status
func (s *teamsBotService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	bot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get platform bot: %w", err)
	}

	botStatus := models.BotStatus(status)
	bot.Status = &botStatus
	bot.UpdatedAt = time.Now()

	return s.repo.Update(ctx, bot)
}

// UpdateCapabilities updates platform bot capabilities
func (s *teamsBotService) UpdateCapabilities(ctx context.Context, id uuid.UUID, capabilities map[string]interface{}) error {
	bot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get platform bot: %w", err)
	}

	capabilitiesObj := models.JSONBObject(capabilities)
	bot.Capabilities = &capabilitiesObj
	bot.UpdatedAt = time.Now()

	return s.repo.Update(ctx, bot)
}

// TestConnection tests platform bot connection
func (s *teamsBotService) TestConnection(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get platform bot: %w", err)
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
func (s *destinationService) Create(ctx context.Context, req *CreateRequest[models.Destination]) (*models.Destination, error) {
	// Set default values
	destination := req.Data
	// validate targets before create
	if err := validateTeamsTargets(destination.Targets); err != nil {
		return nil, fmt.Errorf("invalid targets: %w", err)
	}
	destination.ID = uuid.New()
	destination.CreatedAt = time.Now()
	destination.UpdatedAt = time.Now()
	// Only set default status if not already set
	if destination.Status == "" {
		destination.Status = "active"
	}
	if destination.ValidationStatus == "" {
		destination.ValidationStatus = string(models.NotificationStatusPending)
	}

	err := s.repo.Create(ctx, &destination)
	if err != nil {
		return nil, err
	}

	return &destination, nil
}

// GetByID retrieves a destination by ID
func (s *destinationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Destination, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates a destination
func (s *destinationService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[models.Destination]) (*models.Destination, error) {
	destination := req.Data
	destination.ID = id
	// validate if targets provided (Update via this path might include targets)
	if len(destination.Targets) > 0 {
		if err := validateTeamsTargets(destination.Targets); err != nil {
			return nil, fmt.Errorf("invalid targets: %w", err)
		}
	}
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
func (s *destinationService) List(ctx context.Context, req *ListRequest) ([]*models.Destination, error) {
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
func (s *destinationService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Destination, error) {
	return s.repo.GetByProjectID(ctx, projectID)
}

// GetByBotID retrieves destinations by bot ID
func (s *destinationService) GetByBotID(ctx context.Context, botID uuid.UUID) ([]*models.Destination, error) {
	return s.repo.GetByBotID(ctx, botID)
}

// GetByStatus retrieves destinations by status
func (s *destinationService) GetByStatus(ctx context.Context, status string) ([]*models.Destination, error) {
	return s.repo.GetByStatus(ctx, status)
}

// GetByValidationStatus retrieves destinations by validation status
func (s *destinationService) GetByValidationStatus(ctx context.Context, validationStatus string) ([]*models.Destination, error) {
	return s.repo.GetByValidationStatus(ctx, validationStatus)
}

// SearchDestinations searches destinations by name or description
func (s *destinationService) SearchDestinations(ctx context.Context, query string) ([]*models.Destination, error) {
	return s.repo.SearchDestinations(ctx, query)
}

// UpdateTargets updates destination targets
func (s *destinationService) UpdateTargets(ctx context.Context, id uuid.UUID, targets models.JSONBTargets) (*models.Destination, error) {
	destination, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination: %w", err)
	}

	// validate targets before update
	if err := validateTeamsTargets(targets); err != nil {
		return nil, fmt.Errorf("invalid targets: %w", err)
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
func (s *destinationService) ValidateTargets(ctx context.Context, id uuid.UUID) (*models.Destination, error) {
	destination, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get destination: %w", err)
	}

	// perform validation
	if err := validateTeamsTargets(destination.Targets); err != nil {
		return nil, fmt.Errorf("invalid targets: %w", err)
	}
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

// validateTeamsTargets enforces per-type requirements for TeamsTarget
func validateTeamsTargets(targets models.JSONBTargets) error {
	if len(targets) == 0 {
		return fmt.Errorf("targets cannot be empty")
	}
	for idx, t := range targets {
		tt := strings.TrimSpace(strings.ToLower(t.Type))
		switch tt {
		case "personal":
			// personal allows either conversation_id or email
			if strings.TrimSpace(t.ConversationID) == "" && strings.TrimSpace(t.Email) == "" {
				return fmt.Errorf("targets[%d]: personal requires conversation_id or email", idx)
			}
		case "channel", "groupchat":
			// channel and groupChat must provide conversation_id
			if strings.TrimSpace(t.ConversationID) == "" {
				return fmt.Errorf("targets[%d]: %s requires conversation_id", idx, tt)
			}
		default:
			return fmt.Errorf("targets[%d]: unsupported type %q", idx, t.Type)
		}
		// tenant id recommended for routing
		if strings.TrimSpace(t.TenantID) == "" {
			return fmt.Errorf("targets[%d]: tenant_id is required", idx)
		}
	}
	return nil
}

// =============================================
// Notification Service Implementation
// =============================================

// NewNotificationService creates a new notification service
func NewNotificationService(repo repositories.NotificationRepository, destinationRepo repositories.DestinationRepository, notificationDestRepo repositories.NotificationDestinationRepository, broadcaster BroadcastService, metricsService MetricsService) NotificationService {
	return &notificationService{
		repo:                 repo,
		destinationRepo:      destinationRepo,
		notificationDestRepo: notificationDestRepo,
		broadcaster:          broadcaster,
		metricsService:       metricsService,
	}
}

type notificationService struct {
	repo                 repositories.NotificationRepository
	destinationRepo      repositories.DestinationRepository
	notificationDestRepo repositories.NotificationDestinationRepository
	broadcaster          BroadcastService
	metricsService       MetricsService
}

// Create creates a new notification
func (s *notificationService) Create(ctx context.Context, req *CreateRequest[models.Notification]) (*models.Notification, error) {
	notification := req.Data
	notification.ID = uuid.New()
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()
	if notification.Status == models.NotificationStatus("") {
		notification.Status = models.NotificationStatusPending
	}

	err := s.repo.Create(ctx, &notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// GetByID retrieves a notification by ID
func (s *notificationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates a notification
func (s *notificationService) Update(ctx context.Context, id uuid.UUID, req *UpdateRequest[models.Notification]) (*models.Notification, error) {
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
func (s *notificationService) List(ctx context.Context, req *ListRequest) ([]*models.Notification, error) {
	return s.repo.List(ctx, req.Limit, req.Offset)
}

// Count returns the total number of notifications
func (s *notificationService) Count(ctx context.Context, req *CountRequest) (int64, error) {
	return s.repo.Count(ctx)
}

// GetByProjectID retrieves notifications by project ID
func (s *notificationService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Notification, error) {
	return s.repo.GetByProjectID(ctx, projectID)
}

// GetBySenderID retrieves notifications by sender ID
func (s *notificationService) GetBySenderID(ctx context.Context, senderID uuid.UUID) ([]*models.Notification, error) {
	return s.repo.GetBySenderID(ctx, senderID)
}

// GetByStatus retrieves notifications by status
func (s *notificationService) GetByStatus(ctx context.Context, status string) ([]*models.Notification, error) {
	return s.repo.GetByStatus(ctx, status)
}

// GetByDateRange retrieves notifications by date range
func (s *notificationService) GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Notification, error) {
	return s.repo.GetByDateRange(ctx, start, end)
}

// SendNotification sends a notification
func (s *notificationService) SendNotification(ctx context.Context, req *SendNotificationRequest) (*models.Notification, error) {
	// Convert attachments slice
	var attachments models.JSONBObjectArray
	if req.Attachments != nil {
		attachments = models.JSONBObjectArray(req.Attachments)
	} else {
		attachments = models.JSONBObjectArray{}
	}

	notification := &models.Notification{
		ProjectID:   req.ProjectID,
		SenderID:    req.SenderID,
		MessageType: req.MessageType,
		Content:     req.Content,
		Mentions:    models.JSONBStringArray(req.Mentions),
		Attachments: attachments,
		Priority:    req.Priority,
		Status:      models.NotificationStatusSent,
		Metadata:    models.JSONBObject(req.Metadata),
	}

	// Get destinations for the project
	destinations, err := s.destinationRepo.GetByProjectID(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get destinations: %w", err)
	}

	// ASYNC: Persist notification first, then create notification_destinations for async actors
	notification.ID = uuid.New()
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()
	notification.Status = models.NotificationStatusPending // Start as pending, actors will update status

	err = s.repo.Create(ctx, notification)
	if err != nil {
		return nil, err
	}
	// handle filter
	filterFlag := len(req.Targets) > 0
	filterList := make([]string, 0)
	if filterFlag {
		for _, rt := range req.Targets {
			if strings.EqualFold(rt, "all") {
				filterFlag = false
				break
			}
			filterList = append(filterList, rt)
		}
	}

	if filterFlag {
		destinations = filterDestinations(destinations, filterList)
	}

	// Create notification_destinations entries for each target in destinations
	// Each Teams target gets its own notification_destination record
	if len(destinations) > 0 {
		var nds []*models.NotificationDestination
		for _, d := range destinations {
			for _, t := range d.Targets {
				conversationID := t.ConversationID
				nd := &models.NotificationDestination{
					BaseModel:      models.BaseModel{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
					NotificationID: notification.ID,
					DestinationID:  d.ID,
					Priority:       req.Priority,
					ConversationID: &conversationID,
					BotID:          d.BotID,
					BotType:        nil, // Will be set by actor when it gets bot info
					Status:         string(models.NotificationStatusPending),
					RetryCount:     0,
					MaxRetries:     5,
				}
				nds = append(nds, nd)
			}
		}

		if len(nds) > 0 {
			if err := s.notificationDestRepo.CreateBatch(ctx, nds); err != nil {
				return nil, fmt.Errorf("failed to create notification_destinations: %w", err)
			}
		}
	}

	// Update metrics
	if s.metricsService != nil {
		// Increment pending notifications
		s.metricsService.IncrementNotificationPending(ctx)

		// Update active projects count (simplified - in real implementation, query from DB)
		// s.metricsService.UpdateActiveProjects(ctx, int64(len(destinations)))
	}

	// TODO: Spawn actors for each notification_destination
	// This will be implemented when ActorPool is wired in

	return notification, nil
}

func filterDestinations(destinations []*models.Destination, targets []string) []*models.Destination {
	if len(targets) == 0 {
		return destinations
	}

	// Preprocess targets into sets for O(1) lookup
	targetEmails := make(map[string]struct{})
	targetConvIDs := make(map[string]struct{})

	for _, t := range targets {
		t = strings.TrimSpace(t)
		if strings.Contains(t, "@") {
			targetEmails[strings.ToLower(t)] = struct{}{}
		} else {
			targetConvIDs[strings.ToLower(t)] = struct{}{}
		}
	}

	filtered := make([]*models.Destination, 0)

	for _, d := range destinations {
		if isDestinationMatched(d, targetEmails, targetConvIDs) {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

func isDestinationMatched(d *models.Destination, targetEmails, targetConvIDs map[string]struct{}) bool {
	for _, tgt := range d.Targets {
		// Check Email
		if tgt.Email != "" {
			if _, ok := targetEmails[strings.ToLower(tgt.Email)]; ok {
				return true
			}
		}
		// Check ConversationID
		if tgt.ConversationID != "" {
			if _, ok := targetConvIDs[strings.ToLower(tgt.ConversationID)]; ok {
				return true
			}
		}
	}
	return false
}

// RetryNotification retries a failed notification
func (s *notificationService) RetryNotification(ctx context.Context, id uuid.UUID) error {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if notification.Status != models.NotificationStatusFailed {
		return fmt.Errorf("notification is not in failed status")
	}

	// Reset status and retry
	notification.Status = models.NotificationStatusPending
	notification.ErrorMessage = ""
	notification.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, notification)
	if err != nil {
		return err
	}

	return nil
}

// CancelNotification cancels a notification
func (s *notificationService) CancelNotification(ctx context.Context, id uuid.UUID) error {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if notification.Status == models.NotificationStatusSent {
		return fmt.Errorf("cannot cancel already sent notification")
	}

	notification.Status = models.NotificationStatusCancelled
	notification.UpdatedAt = time.Now()

	return s.repo.Update(ctx, notification)
}

// getUserEmailFromGraphAPI queries Microsoft Graph API to get user email
func (s *messagesService) getUserEmailFromGraphAPI(ctx context.Context, aadObjectID, tenantID string) string {
	// Get access token for Graph API
	accessToken, err := s.getGraphAPIAccessToken(ctx, tenantID)
	if err != nil {
		log.Printf("messages: failed to get Graph API access token: %v", err)
		return ""
	}

	// Make HTTP request to Graph API
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", aadObjectID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("messages: failed to create Graph API request: %v", err)
		return ""
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("messages: failed to call Graph API: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("messages: Graph API returned status %d", resp.StatusCode)
		return ""
	}

	// Parse response
	var userInfo struct {
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
		DisplayName       string `json:"displayName"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("messages: failed to read Graph API response: %v", err)
		return ""
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		log.Printf("messages: failed to parse Graph API response: %v", err)
		return ""
	}

	// Return mail if available, otherwise userPrincipalName
	if userInfo.Mail != "" {
		return userInfo.Mail
	}
	return userInfo.UserPrincipalName
}

// getGraphAPIAccessToken gets access token for Microsoft Graph API
func (s *messagesService) getGraphAPIAccessToken(ctx context.Context, tenantID string) (string, error) {
	// Use the same bot credentials to get Graph API access token
	// This uses Client Credentials flow for application permissions

	// Get bot credentials from environment or models
	appID := os.Getenv("TEAMS_BOT_APP_ID")
	appPassword := os.Getenv("TEAMS_BOT_APP_PASSWORD")

	if appID == "" || appPassword == "" {
		return "", fmt.Errorf("TEAMS_BOT_APP_ID or TEAMS_BOT_APP_PASSWORD not set")
	}

	// Microsoft Graph API token endpoint
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	// Prepare token request
	data := url.Values{}
	data.Set("client_id", appID)
	data.Set("client_secret", appPassword)
	data.Set("scope", "https://graph.microsoft.com/.default")
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed: %d %s", resp.StatusCode, string(body))
	}

	// Parse token response
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}
