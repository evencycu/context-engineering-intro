package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository interface defines common models operations
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id uuid.UUID) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*T, error)
	Count(ctx context.Context) (int64, error)
}

// CompanyRepository defines company-specific operations
type CompanyRepository interface {
	Repository[models.Company]
	GetByEmail(ctx context.Context, email string) (*models.Company, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Company, error)
}

// BotRepository defines bot-specific operations
type BotRepository interface {
	Repository[models.TeamsBot]
	GetByAppID(ctx context.Context, appID string) (*models.TeamsBot, error)
	GetByStatus(ctx context.Context, status string) ([]*models.TeamsBot, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*models.TeamsBot, error)
}

// DestinationRepository defines destination-specific operations
type DestinationRepository interface {
	Repository[models.Destination]
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Destination, error)
	GetByBotID(ctx context.Context, botID uuid.UUID) ([]*models.Destination, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Destination, error)
	GetByValidationStatus(ctx context.Context, validationStatus string) ([]*models.Destination, error)
	SearchDestinations(ctx context.Context, query string) ([]*models.Destination, error)
}

// NotificationRepository defines notification-specific operations
type NotificationRepository interface {
	Repository[models.Notification]
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Notification, error)
	GetBySenderID(ctx context.Context, senderID uuid.UUID) ([]*models.Notification, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Notification, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Notification, error)
	UpdateStatus(ctx context.Context, notificationID uuid.UUID, status models.NotificationStatus, errorMessage string) error
}

// NotificationDestinationRepository defines operations for notification_destinations
type NotificationDestinationRepository interface {
	Create(ctx context.Context, entity *models.NotificationDestination) error
	CreateBatch(ctx context.Context, entities []*models.NotificationDestination) error
	Update(ctx context.Context, entity *models.NotificationDestination) error
	UpdateStatusAndRetry(ctx context.Context, id uuid.UUID, status string, retryCount int, nextRetryAt *time.Time, failureReason string, errorMessage string, retryAfter *int) error
	MarkSent(ctx context.Context, id uuid.UUID, teamsMessageID string, sentAt time.Time) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.NotificationDestination, error)
	GetByNotificationID(ctx context.Context, notificationID uuid.UUID) ([]*models.NotificationDestination, error)
	GetRetryReady(ctx context.Context, limit int) ([]*models.NotificationDestination, error)
	GetPending(ctx context.Context, limit int) ([]*models.NotificationDestination, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

// BaseRepository provides common repository functionality
type BaseRepository struct {
	// This will be implemented with actual models connection
	// For now, it's a placeholder
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Limit  int `form:"limit" binding:"min=1,max=100"`
	Offset int `form:"offset" binding:"min=0"`
}

// SortParams represents sorting parameters
type SortParams struct {
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order" binding:"oneof=asc desc"`
}

// FilterParams represents filtering parameters
type FilterParams struct {
	Status string `form:"status"`
	Search string `form:"search"`
}

// Response represents a standard API response
type Response[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Pagination struct {
		Total  int64 `json:"total"`
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
		Pages  int   `json:"pages"`
	} `json:"pagination"`
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository(db *sqlx.DB) CompanyRepository {
	return &companyRepository{db: db}
}

// UserRepository defines the interface for user operations
type UserRepository interface {
	Repository[models.User]
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.User, error)
	GetByRole(ctx context.Context, role string) ([]*models.User, error)
	GetByStatus(ctx context.Context, status string) ([]*models.User, error)
}

// userRepository implements UserRepository
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, entity *models.User) error {
	query := `INSERT INTO users (id, company_id, email, name, password_hash, role, status, api_key_hash, api_key_expires_at, last_login_at, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.Email, entity.Name,
		entity.PasswordHash, entity.Role, entity.Status, entity.APIKeyHash,
		entity.APIKeyExpiresAt, entity.LastLoginAt, entity.CreatedAt, entity.UpdatedAt)
	return err
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, entity *models.User) error {
	query := `UPDATE users SET company_id = $2, email = $3, name = $4, password_hash = $5, 
			  role = $6, status = $7, api_key_hash = $8, api_key_expires_at = $9, last_login_at = $10, updated_at = $11 
			  WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.Email, entity.Name,
		entity.PasswordHash, entity.Role, entity.Status, entity.APIKeyHash,
		entity.APIKeyExpiresAt, entity.LastLoginAt, entity.UpdatedAt)
	return err
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	return users, err
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByCompanyID retrieves users by company ID
func (r *userRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	query := `SELECT * FROM users WHERE company_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &users, query, companyID)
	return users, err
}

// GetByRole retrieves users by role
func (r *userRepository) GetByRole(ctx context.Context, role string) ([]*models.User, error) {
	var users []*models.User
	query := `SELECT * FROM users WHERE role = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &users, query, role)
	return users, err
}

// GetByStatus retrieves users by status
func (r *userRepository) GetByStatus(ctx context.Context, status string) ([]*models.User, error) {
	var users []*models.User
	query := `SELECT * FROM users WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &users, query, status)
	return users, err
}

// companyRepository implements CompanyRepository
type companyRepository struct {
	db *sqlx.DB
}

func (r *companyRepository) Create(ctx context.Context, entity *models.Company) error {
	query := `INSERT INTO companies (id, name, contact_email, contact_phone, address, status, billing_enabled, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.Name, entity.ContactEmail, entity.ContactPhone,
		entity.Address, entity.Status, entity.BillingEnabled, entity.CreatedAt, entity.UpdatedAt)
	return err
}

func (r *companyRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Company, error) {
	var company models.Company
	query := `SELECT * FROM companies WHERE id = $1`
	err := r.db.GetContext(ctx, &company, query, id)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) Update(ctx context.Context, entity *models.Company) error {
	query := `UPDATE companies SET name = $2, contact_email = $3, contact_phone = $4, 
			  address = $5, status = $6, billing_enabled = $7, updated_at = $8 
			  WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.Name, entity.ContactEmail, entity.ContactPhone,
		entity.Address, entity.Status, entity.BillingEnabled, entity.UpdatedAt)
	return err
}

func (r *companyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM companies WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *companyRepository) List(ctx context.Context, limit, offset int) ([]*models.Company, error) {
	var companies []*models.Company
	query := `SELECT * FROM companies ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &companies, query, limit, offset)
	return companies, err
}

func (r *companyRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM companies`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

func (r *companyRepository) GetByEmail(ctx context.Context, email string) (*models.Company, error) {
	var company models.Company
	query := `SELECT * FROM companies WHERE contact_email = $1`
	err := r.db.GetContext(ctx, &company, query, email)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) GetByStatus(ctx context.Context, status string) ([]*models.Company, error) {
	var companies []*models.Company
	query := `SELECT * FROM companies WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &companies, query, status)
	return companies, err
}

// ProjectRepository defines the interface for project operations
type ProjectRepository interface {
	Repository[models.Project]
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.Project, error)
	GetByKeyName(ctx context.Context, keyName string) (*models.Project, error)
	GetByNotifyKey(ctx context.Context, notifyKey string) (*models.Project, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Project, error)
}

// projectRepository implements ProjectRepository
type projectRepository struct {
	db *sqlx.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *sqlx.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// Create creates a new project
func (r *projectRepository) Create(ctx context.Context, entity *models.Project) error {
	query := `INSERT INTO projects (id, company_id, notify_key, description, status, daily_limit, monthly_limit, priority, created_by, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.NotifyKey, entity.Description,
		entity.Status, entity.DailyLimit, entity.MonthlyLimit, entity.Priority,
		entity.CreatedBy, entity.CreatedAt, entity.UpdatedAt)
	return err
}

// GetByID retrieves a project by ID
func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	query := `SELECT * FROM projects WHERE id = $1`
	err := r.db.GetContext(ctx, &project, query, id)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// Update updates an existing project
func (r *projectRepository) Update(ctx context.Context, entity *models.Project) error {
	query := `UPDATE projects SET company_id = $2, notify_key = $3, description = $4, 
			  status = $5, daily_limit = $6, monthly_limit = $7, priority = $8, updated_at = $9 
			  WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.NotifyKey, entity.Description,
		entity.Status, entity.DailyLimit, entity.MonthlyLimit, entity.Priority, entity.UpdatedAt)
	return err
}

// Delete deletes a project by ID
func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves projects with pagination
func (r *projectRepository) List(ctx context.Context, limit, offset int) ([]*models.Project, error) {
	var projects []*models.Project
	query := `SELECT * FROM projects ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &projects, query, limit, offset)
	return projects, err
}

// Count returns the total number of projects
func (r *projectRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM projects`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

// GetByCompanyID retrieves projects by company ID
func (r *projectRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.Project, error) {
	var projects []*models.Project
	query := `SELECT * FROM projects WHERE company_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &projects, query, companyID)
	return projects, err
}

// GetByKeyName retrieves a project by key name
func (r *projectRepository) GetByKeyName(ctx context.Context, keyName string) (*models.Project, error) {
	var project models.Project
	query := `SELECT * FROM projects WHERE notify_key = $1`
	err := r.db.GetContext(ctx, &project, query, keyName)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByNotifyKey retrieves a project by notify_key
func (r *projectRepository) GetByNotifyKey(ctx context.Context, notifyKey string) (*models.Project, error) {
	var project models.Project
	query := `SELECT * FROM projects WHERE notify_key = $1`
	err := r.db.GetContext(ctx, &project, query, notifyKey)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByStatus retrieves projects by status
func (r *projectRepository) GetByStatus(ctx context.Context, status string) ([]*models.Project, error) {
	var projects []*models.Project
	query := `SELECT * FROM projects WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &projects, query, status)
	return projects, err
}

// TeamsBotRepository defines the interface for bot operations
type TeamsBotRepository interface {
	Repository[models.TeamsBot]
	GetByAppID(ctx context.Context, appID string) (*models.TeamsBot, error)
	GetByStatus(ctx context.Context, status string) ([]*models.TeamsBot, error)
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.TeamsBot, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*models.TeamsBot, error)
}

// teamsBotRepository implements TeamsBotRepository
type teamsBotRepository struct {
	db *sqlx.DB
}

// NewTeamsBotRepository creates a new bot repository
func NewTeamsBotRepository(db *sqlx.DB) TeamsBotRepository {
	return &teamsBotRepository{db: db}
}

// Create creates a new platform bot
func (r *teamsBotRepository) Create(ctx context.Context, entity *models.TeamsBot) error {
	query := `INSERT INTO teams_bots (id, type, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests, created_at, updated_at) 
			  VALUES ($1, 'platform', $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	// Convert capabilities to JSONB
	var capabilitiesJSON []byte
	var err error
	if entity.Capabilities != nil {
		capabilitiesJSON, err = json.Marshal(map[string]any(*entity.Capabilities))
		if err != nil {
			return fmt.Errorf("failed to marshal capabilities: %w", err)
		}
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.Name, entity.Description, entity.AppID, entity.AppPasswordHash,
		entity.TenantID, entity.Status, entity.WebhookURL, capabilitiesJSON,
		entity.RateLimitPerMinute, entity.MaxConcurrentRequests, entity.CreatedAt, entity.UpdatedAt)
	return err
}

// GetByID retrieves a platform bot by ID
func (r *teamsBotRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TeamsBot, error) {
	var bot models.TeamsBot
	var capabilitiesJSON []byte
	query := `SELECT id, created_at, updated_at, name, description, app_id, app_password_hash, 
			  tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests 
			  FROM teams_bots WHERE type = 'platform' AND id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&bot.ID, &bot.CreatedAt, &bot.UpdatedAt, &bot.Name, &bot.Description,
		&bot.AppID, &bot.AppPasswordHash, &bot.TenantID, &bot.Status, &bot.WebhookURL,
		&capabilitiesJSON, &bot.RateLimitPerMinute, &bot.MaxConcurrentRequests,
	)
	if err != nil {
		return nil, err
	}

	// Parse capabilities JSON
	if len(capabilitiesJSON) > 0 {
		var capabilities map[string]any
		if err := json.Unmarshal(capabilitiesJSON, &capabilities); err != nil {
			return nil, fmt.Errorf("failed to unmarshal capabilities: %w", err)
		}
		capabilitiesObj := models.JSONBObject(capabilities)
		bot.Capabilities = &capabilitiesObj
	}

	return &bot, nil
}

// Update updates an existing platform bot
func (r *teamsBotRepository) Update(ctx context.Context, entity *models.TeamsBot) error {
	query := `UPDATE teams_bots SET name = $2, description = $3, app_id = $4, app_password_hash = $5, 
			  tenant_id = $6, status = $7, webhook_url = $8, capabilities = $9, rate_limit_per_minute = $10, 
			  max_concurrent_requests = $11, updated_at = $12 
			  WHERE id = $1 AND type = 'platform'`

	// Convert capabilities to JSONB
	var capabilitiesJSON []byte
	var err error
	if entity.Capabilities != nil {
		capabilitiesJSON, err = json.Marshal(map[string]any(*entity.Capabilities))
		if err != nil {
			return fmt.Errorf("failed to marshal capabilities: %w", err)
		}
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.Name, entity.Description, entity.AppID, entity.AppPasswordHash,
		entity.TenantID, entity.Status, entity.WebhookURL, capabilitiesJSON,
		entity.RateLimitPerMinute, entity.MaxConcurrentRequests, entity.UpdatedAt)
	return err
}

// Delete deletes a platform bot by ID
func (r *teamsBotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM teams_bots WHERE id = $1 AND type = 'platform'`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves platform bots with pagination
func (r *teamsBotRepository) List(ctx context.Context, limit, offset int) ([]*models.TeamsBot, error) {
	query := `SELECT id, created_at, updated_at, name, description, app_id, app_password_hash, 
			  tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests 
			  FROM teams_bots WHERE type = 'platform' ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bots []*models.TeamsBot
	for rows.Next() {
		var bot models.TeamsBot
		var capabilitiesJSON []byte
		err := rows.Scan(
			&bot.ID, &bot.CreatedAt, &bot.UpdatedAt, &bot.Name, &bot.Description,
			&bot.AppID, &bot.AppPasswordHash, &bot.TenantID, &bot.Status, &bot.WebhookURL,
			&capabilitiesJSON, &bot.RateLimitPerMinute, &bot.MaxConcurrentRequests,
		)
		if err != nil {
			return nil, err
		}

		// Parse capabilities JSON
		if len(capabilitiesJSON) > 0 {
			var capabilities map[string]any
			if err := json.Unmarshal(capabilitiesJSON, &capabilities); err != nil {
				return nil, fmt.Errorf("failed to unmarshal capabilities: %w", err)
			}
			capabilitiesObj := models.JSONBObject(capabilities)
			bot.Capabilities = &capabilitiesObj
		}

		bots = append(bots, &bot)
	}

	return bots, nil
}

// Count returns the total number of platform bots
func (r *teamsBotRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM teams_bots WHERE type = 'platform'`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

// GetByAppID retrieves a bot by app ID (any type)
func (r *teamsBotRepository) GetByAppID(ctx context.Context, appID string) (*models.TeamsBot, error) {
	var bot models.TeamsBot
	query := `SELECT * FROM teams_bots WHERE app_id = $1`
	err := r.db.GetContext(ctx, &bot, query, appID)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

// GetByStatus retrieves platform bots by status
func (r *teamsBotRepository) GetByStatus(ctx context.Context, status string) ([]*models.TeamsBot, error) {
	var bots []*models.TeamsBot
	query := `SELECT * FROM teams_bots WHERE type = 'platform' AND status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &bots, query, status)
	return bots, err
}

// GetByCompanyID retrieves platform bots by company ID
func (r *teamsBotRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*models.TeamsBot, error) {
	var bots []*models.TeamsBot
	query := `SELECT * FROM teams_bots WHERE type = 'platform' AND company_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &bots, query, companyID)
	return bots, err
}

// GetByTenantID retrieves platform bots by tenant ID
func (r *teamsBotRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*models.TeamsBot, error) {
	var bots []*models.TeamsBot
	query := `SELECT * FROM teams_bots WHERE type = 'platform' AND tenant_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &bots, query, tenantID)
	return bots, err
}

// BotInstallationRepository defines operations for bot_installations
type BotInstallationRepository interface {
	Upsert(ctx context.Context, entity *models.BotInstallation) error
	GetByBotAndTenant(ctx context.Context, botID uuid.UUID, botType models.BotType, tenantID string) ([]*models.BotInstallation, error)
	GetByConversationType(ctx context.Context, conversationType string, status string) ([]*models.BotInstallation, error)
	GetByConversationID(ctx context.Context, conversationID string) (*models.BotInstallation, error)
	GetActiveInstallations(ctx context.Context, botID uuid.UUID, botType models.BotType, tenantID string) ([]*models.BotInstallation, error)
	GetActiveInstallationsByTenant(ctx context.Context, tenantID string) ([]*models.BotInstallation, error)
	GetActivePersonalByEmail(ctx context.Context, tenantID string, email string) ([]*models.BotInstallation, error)
	UpdateActivity(ctx context.Context, id uuid.UUID) error
	MarkAsStale(ctx context.Context, id uuid.UUID) error
	MarkAsUninstalled(ctx context.Context, id uuid.UUID) error
}

type botInstallationRepository struct {
	db *sqlx.DB
}

// NewBotInstallationRepository creates a new bot installation repository
func NewBotInstallationRepository(db *sqlx.DB) BotInstallationRepository {
	return &botInstallationRepository{db: db}
}

// Upsert inserts or updates a bot installation record based on unique keys
func (r *botInstallationRepository) Upsert(ctx context.Context, entity *models.BotInstallation) error {
	metadataJSON, err := json.Marshal(entity.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
        INSERT INTO bot_installations (
            bot_id, bot_type, teams_tenant_id, conversation_type, conversation_id, service_url,
            recipient_id, recipient_name, from_id, from_name, from_aad_object_id, email, description_name,
            installation_status, installed_at, uninstalled_at, metadata, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6,
            $7, $8, $9, $10, $11, $12, $13,
            $14, COALESCE($15, NOW()), $16, $17, COALESCE($18, NOW()), COALESCE($19, NOW())
        )
        ON CONFLICT (bot_id, bot_type, teams_tenant_id, conversation_id)
        DO UPDATE SET
            conversation_type = EXCLUDED.conversation_type,
            service_url = EXCLUDED.service_url,
            recipient_id = EXCLUDED.recipient_id,
            recipient_name = EXCLUDED.recipient_name,
            from_id = EXCLUDED.from_id,
            from_name = EXCLUDED.from_name,
            from_aad_object_id = EXCLUDED.from_aad_object_id,
            email = EXCLUDED.email,
            description_name = EXCLUDED.description_name,
            installation_status = EXCLUDED.installation_status,
            uninstalled_at = EXCLUDED.uninstalled_at,
            metadata = EXCLUDED.metadata,
            updated_at = NOW()
        RETURNING id;
    `

	return r.db.QueryRowContext(
		ctx,
		query,
		entity.BotID,
		entity.BotType,
		entity.TeamsTenantID,
		entity.ConversationType,
		entity.ConversationID,
		entity.ServiceURL,
		entity.RecipientID,
		entity.RecipientName,
		entity.FromID,
		entity.FromName,
		entity.FromAADObjectID,
		entity.Email,
		entity.DescriptionName,
		entity.InstallationStatus,
		entity.InstalledAt,
		entity.UninstalledAt,
		metadataJSON,
		entity.CreatedAt,
		entity.UpdatedAt,
	).Scan(&entity.ID)
}

// GetByBotAndTenant retrieves installations for a specific bot and tenant
func (r *botInstallationRepository) GetByBotAndTenant(ctx context.Context, botID uuid.UUID, botType models.BotType, tenantID string) ([]*models.BotInstallation, error) {
	query := `
        SELECT id, bot_id, bot_type, teams_tenant_id, conversation_type, conversation_id, service_url,
               recipient_id, recipient_name, from_id, from_name, from_aad_object_id,
               installation_status, installed_at, uninstalled_at, metadata, created_at, updated_at
        FROM bot_installations
        WHERE bot_id = $1 AND bot_type = $2 AND teams_tenant_id = $3
        ORDER BY installed_at DESC
    `

	var installations []*models.BotInstallation
	err := r.db.SelectContext(ctx, &installations, query, botID, botType, tenantID)
	return installations, err
}

// GetByConversationType retrieves installations by conversation type and status
func (r *botInstallationRepository) GetByConversationType(ctx context.Context, conversationType string, status string) ([]*models.BotInstallation, error) {
	query := `
        SELECT id, bot_id, bot_type, teams_tenant_id, conversation_type, conversation_id, service_url,
               recipient_id, recipient_name, from_id, from_name, from_aad_object_id,
               installation_status, installed_at, uninstalled_at, metadata, created_at, updated_at
        FROM bot_installations
        WHERE conversation_type = $1 AND installation_status = $2
        ORDER BY installed_at DESC
    `

	var installations []*models.BotInstallation
	err := r.db.SelectContext(ctx, &installations, query, conversationType, status)
	return installations, err
}

// GetByConversationID retrieves installation by conversation ID
func (r *botInstallationRepository) GetByConversationID(ctx context.Context, conversationID string) (*models.BotInstallation, error) {
	query := `
        SELECT id, bot_id, bot_type, teams_tenant_id, conversation_type, conversation_id, service_url,
               recipient_id, recipient_name, from_id, from_name, from_aad_object_id,
               installation_status, installed_at, uninstalled_at, metadata, created_at, updated_at
        FROM bot_installations
        WHERE conversation_id = $1
    `

	var installation models.BotInstallation
	err := r.db.GetContext(ctx, &installation, query, conversationID)
	if err != nil {
		return nil, err
	}
	return &installation, nil
}

// GetActiveInstallations retrieves active installations for a bot and tenant
func (r *botInstallationRepository) GetActiveInstallations(ctx context.Context, botID uuid.UUID, botType models.BotType, tenantID string) ([]*models.BotInstallation, error) {
	query := `
        SELECT id, bot_id, bot_type, teams_tenant_id, conversation_type, conversation_id, service_url,
               recipient_id, recipient_name, from_id, from_name, from_aad_object_id,
               installation_status, installed_at, uninstalled_at, metadata, created_at, updated_at
        FROM bot_installations
        WHERE bot_id = $1 AND bot_type = $2 AND teams_tenant_id = $3 AND installation_status = 'active'
        ORDER BY installed_at DESC
    `

	var installations []*models.BotInstallation
	err := r.db.SelectContext(ctx, &installations, query, botID, botType, tenantID)
	return installations, err
}

// GetActiveInstallationsByTenant retrieves all active installations for a tenant
func (r *botInstallationRepository) GetActiveInstallationsByTenant(ctx context.Context, tenantID string) ([]*models.BotInstallation, error) {
	query := `
        SELECT id, bot_id, bot_type, teams_tenant_id, conversation_type, conversation_id, service_url,
               recipient_id, recipient_name, from_id, from_name, from_aad_object_id,
               installation_status, installed_at, uninstalled_at, metadata, created_at, updated_at
        FROM bot_installations
        WHERE teams_tenant_id = $1 AND installation_status = 'active'
        ORDER BY installed_at DESC
    `

	var installations []*models.BotInstallation
	err := r.db.SelectContext(ctx, &installations, query, tenantID)
	return installations, err
}

// GetActivePersonalByEmail retrieves active personal installations for a tenant by email
func (r *botInstallationRepository) GetActivePersonalByEmail(ctx context.Context, tenantID string, email string) ([]*models.BotInstallation, error) {
	query := `
        SELECT id, bot_id, bot_type, teams_tenant_id, conversation_type, conversation_id, service_url,
               recipient_id, recipient_name, from_id, from_name, from_aad_object_id,
               installation_status, installed_at, uninstalled_at, metadata, created_at, updated_at, email, description_name
        FROM bot_installations
        WHERE teams_tenant_id = $1
          AND conversation_type = 'personal'
          AND installation_status = 'active'
          AND LOWER(email) = LOWER($2)
        ORDER BY installed_at DESC
    `

	var installations []*models.BotInstallation
	err := r.db.SelectContext(ctx, &installations, query, tenantID, email)
	return installations, err
}

// UpdateActivity updates the last activity timestamp
func (r *botInstallationRepository) UpdateActivity(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE bot_installations SET last_activity_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// MarkAsStale marks an installation as stale
func (r *botInstallationRepository) MarkAsStale(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE bot_installations SET installation_status = 'stale', updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// MarkAsUninstalled marks an installation as uninstalled and sets uninstalled_at
func (r *botInstallationRepository) MarkAsUninstalled(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE bot_installations SET installation_status = 'uninstalled', uninstalled_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// =============================================
// Destination Repository Implementation
// =============================================

// NewDestinationRepository creates a new destination repository
func NewDestinationRepository(db *sqlx.DB) DestinationRepository {
	return &destinationRepository{db: db}
}

type destinationRepository struct {
	db *sqlx.DB
}

// Create creates a new destination
func (r *destinationRepository) Create(ctx context.Context, entity *models.Destination) error {
	// Marshal targets to JSONB
	targetsJSON, err := json.Marshal(entity.Targets)
	if err != nil {
		return fmt.Errorf("failed to marshal targets: %w", err)
	}

	query := `INSERT INTO destinations (id, project_id, name, description, teams_tenant_id, targets, bot_id, status, validation_status, last_validated_at, created_by, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.Name, entity.Description, entity.TeamsTenantID,
		targetsJSON, entity.BotID, entity.Status, entity.ValidationStatus,
		entity.LastValidatedAt, entity.CreatedBy, entity.CreatedAt, entity.UpdatedAt)

	return err
}

// GetByID retrieves a destination by ID
func (r *destinationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Destination, error) {
	var destination models.Destination
	query := `SELECT * FROM destinations WHERE id = $1`
	err := r.db.GetContext(ctx, &destination, query, id)
	if err != nil {
		return nil, err
	}
	return &destination, nil
}

// Update updates a destination
func (r *destinationRepository) Update(ctx context.Context, entity *models.Destination) error {
	// Marshal targets to JSONB
	targetsJSON, err := json.Marshal(entity.Targets)
	if err != nil {
		return fmt.Errorf("failed to marshal targets: %w", err)
	}

	query := `UPDATE destinations SET 
              project_id = $2, name = $3, description = $4, teams_tenant_id = $5, 
              targets = $6, bot_id = $7, status = $8, 
              validation_status = $9, last_validated_at = $10, updated_at = $11
              WHERE id = $1`

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.Name, entity.Description, entity.TeamsTenantID,
		targetsJSON, entity.BotID, entity.Status, entity.ValidationStatus,
		entity.LastValidatedAt, entity.UpdatedAt)

	return err
}

// Delete deletes a destination by ID
func (r *destinationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM destinations WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves destinations with pagination
func (r *destinationRepository) List(ctx context.Context, limit, offset int) ([]*models.Destination, error) {
	var destinations []*models.Destination
	query := `SELECT * FROM destinations ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &destinations, query, limit, offset)
	return destinations, err
}

// Count returns the total number of destinations
func (r *destinationRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM destinations`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

// GetByProjectID retrieves destinations by project ID
func (r *destinationRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Destination, error) {
	var destinations []*models.Destination
	query := `SELECT * FROM destinations WHERE project_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &destinations, query, projectID)
	return destinations, err
}

// GetByBotID retrieves destinations by bot ID
func (r *destinationRepository) GetByBotID(ctx context.Context, botID uuid.UUID) ([]*models.Destination, error) {
	var destinations []*models.Destination
	query := `SELECT * FROM destinations WHERE bot_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &destinations, query, botID)
	return destinations, err
}

// GetByStatus retrieves destinations by status
func (r *destinationRepository) GetByStatus(ctx context.Context, status string) ([]*models.Destination, error) {
	var destinations []*models.Destination
	query := `SELECT * FROM destinations WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &destinations, query, status)
	return destinations, err
}

// GetByValidationStatus retrieves destinations by validation status
func (r *destinationRepository) GetByValidationStatus(ctx context.Context, validationStatus string) ([]*models.Destination, error) {
	var destinations []*models.Destination
	query := `SELECT * FROM destinations WHERE validation_status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &destinations, query, validationStatus)
	return destinations, err
}

// SearchDestinations searches destinations by name or description
func (r *destinationRepository) SearchDestinations(ctx context.Context, query string) ([]*models.Destination, error) {
	var destinations []*models.Destination
	sql := `SELECT * FROM destinations WHERE name ILIKE $1 OR description ILIKE $1 ORDER BY created_at DESC`
	searchQuery := "%" + query + "%"
	err := r.db.SelectContext(ctx, &destinations, sql, searchQuery)
	return destinations, err
}

// =============================================
// Notification Repository Implementation
// =============================================

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *sqlx.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

type notificationRepository struct {
	db *sqlx.DB
}

// Create creates a new notification
func (r *notificationRepository) Create(ctx context.Context, entity *models.Notification) error {
	// Marshal JSONB fields
	mentionsJSON, err := json.Marshal([]string(entity.Mentions))
	if err != nil {
		return fmt.Errorf("failed to marshal mentions: %w", err)
	}

	attachmentsJSON, err := json.Marshal([]map[string]any(entity.Attachments))
	if err != nil {
		return fmt.Errorf("failed to marshal attachments: %w", err)
	}

	metadataJSON, err := json.Marshal(map[string]any(entity.Metadata))
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `INSERT INTO notifications (id, project_id, sender_id, message_type, content, mentions, attachments, priority, status, error_message, metadata, sent_at, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	// Handle nil values for JSONB fields
	var attachmentsJSONValue interface{} = attachmentsJSON
	if len(attachmentsJSON) == 0 {
		attachmentsJSONValue = []byte("[]")
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.SenderID, entity.MessageType, entity.Content,
		mentionsJSON, attachmentsJSONValue, entity.Priority, entity.Status,
		entity.ErrorMessage, metadataJSON, entity.SentAt, entity.CreatedAt, entity.UpdatedAt)

	return err
}

// GetByID retrieves a notification by ID
func (r *notificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	query := `SELECT * FROM notifications WHERE id = $1`
	err := r.db.GetContext(ctx, &notification, query, id)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// Update updates a notification
func (r *notificationRepository) Update(ctx context.Context, entity *models.Notification) error {
	// Marshal JSONB fields
	mentionsJSON, err := json.Marshal([]string(entity.Mentions))
	if err != nil {
		return fmt.Errorf("failed to marshal mentions: %w", err)
	}

	attachmentsJSON, err := json.Marshal([]map[string]any(entity.Attachments))
	if err != nil {
		return fmt.Errorf("failed to marshal attachments: %w", err)
	}

	metadataJSON, err := json.Marshal(map[string]any(entity.Metadata))
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `UPDATE notifications SET 
              project_id = $2, sender_id = $3, message_type = $4, content = $5, 
              mentions = $6, attachments = $7, priority = $8, 
              status = $9, error_message = $10, metadata = $11, sent_at = $12, updated_at = $13
              WHERE id = $1`

	// Handle nil values for JSONB fields
	var attachmentsJSONValue interface{} = attachmentsJSON
	if len(attachmentsJSON) == 0 {
		attachmentsJSONValue = []byte("[]")
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.SenderID, entity.MessageType, entity.Content,
		mentionsJSON, attachmentsJSONValue, entity.Priority, entity.Status,
		entity.ErrorMessage, metadataJSON, entity.SentAt, entity.UpdatedAt)

	return err
}

// Delete deletes a notification by ID
func (r *notificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves notifications with pagination
func (r *notificationRepository) List(ctx context.Context, limit, offset int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := `SELECT * FROM notifications ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &notifications, query, limit, offset)
	return notifications, err
}

// Count returns the total number of notifications
func (r *notificationRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM notifications`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

// =============================================
// NotificationDestination Repository Implementation
// =============================================

// NewNotificationDestinationRepository creates a new repository
func NewNotificationDestinationRepository(db *sqlx.DB) NotificationDestinationRepository {
	return &notificationDestinationRepository{db: db}
}

type notificationDestinationRepository struct {
	db *sqlx.DB
}

// Create inserts a notification_destination record
func (r *notificationDestinationRepository) Create(ctx context.Context, entity *models.NotificationDestination) error {
	query := `INSERT INTO notification_destinations (
        id, notification_id, destination_id, priority, conversation_id, bot_id, bot_type, status,
        error_message, teams_message_id, sent_at, retry_count, max_retries,
        next_retry_at, first_attempt_at, last_attempt_at, failure_reason, retry_after, actor_id,
        created_at, updated_at
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8,
        $9, $10, $11, $12, $13,
        $14, $15, $16, $17, $18, $19,
        $20, $21
    )`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.NotificationID, entity.DestinationID, entity.Priority, entity.ConversationID, entity.BotID, entity.BotType, entity.Status,
		entity.ErrorMessage, entity.TeamsMessageID, entity.SentAt, entity.RetryCount, entity.MaxRetries,
		entity.NextRetryAt, entity.FirstAttemptAt, entity.LastAttemptAt, entity.FailureReason, entity.RetryAfter, entity.ActorID,
		entity.CreatedAt, entity.UpdatedAt,
	)
	return err
}

// CreateBatch inserts multiple notification_destination records
func (r *notificationDestinationRepository) CreateBatch(ctx context.Context, entities []*models.NotificationDestination) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	query := `INSERT INTO notification_destinations (
        id, notification_id, destination_id, priority, conversation_id, bot_id, bot_type, status,
        error_message, teams_message_id, sent_at, retry_count, max_retries,
        next_retry_at, first_attempt_at, last_attempt_at, failure_reason, retry_after, actor_id,
        created_at, updated_at
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8,
        $9, $10, $11, $12, $13,
        $14, $15, $16, $17, $18, $19,
        $20, $21
    )`

	for _, e := range entities {
		if _, err = tx.ExecContext(ctx, query,
			e.ID, e.NotificationID, e.DestinationID, e.Priority, e.ConversationID, e.BotID, e.BotType, e.Status,
			e.ErrorMessage, e.TeamsMessageID, e.SentAt, e.RetryCount, e.MaxRetries,
			e.NextRetryAt, e.FirstAttemptAt, e.LastAttemptAt, e.FailureReason, e.RetryAfter, e.ActorID,
			e.CreatedAt, e.UpdatedAt,
		); err != nil {
			return err
		}
	}
	return nil
}

// Update updates a notification_destination record
func (r *notificationDestinationRepository) Update(ctx context.Context, entity *models.NotificationDestination) error {
	query := `UPDATE notification_destinations SET 
        destination_id = $2, priority = $3, conversation_id = $4, bot_id = $5, bot_type = $6, status = $7,
        error_message = $8, teams_message_id = $9, sent_at = $10, retry_count = $11, max_retries = $12,
        next_retry_at = $13, first_attempt_at = $14, last_attempt_at = $15, failure_reason = $16, retry_after = $17, actor_id = $18,
        updated_at = $19
        WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.DestinationID, entity.Priority, entity.ConversationID, entity.BotID, entity.BotType, entity.Status,
		entity.ErrorMessage, entity.TeamsMessageID, entity.SentAt, entity.RetryCount, entity.MaxRetries,
		entity.NextRetryAt, entity.FirstAttemptAt, entity.LastAttemptAt, entity.FailureReason, entity.RetryAfter, entity.ActorID,
		time.Now(),
	)
	return err
}

// UpdateStatusAndRetry updates status and retry-related fields
func (r *notificationDestinationRepository) UpdateStatusAndRetry(ctx context.Context, id uuid.UUID, status string, retryCount int, nextRetryAt *time.Time, failureReason string, errorMessage string, retryAfter *int) error {
	query := `UPDATE notification_destinations SET 
        status = $2, retry_count = $3, next_retry_at = $4, failure_reason = $5, error_message = $6, retry_after = $7, updated_at = NOW()
        WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, status, retryCount, nextRetryAt, failureReason, errorMessage, retryAfter)
	return err
}

// MarkSent marks a destination as sent
func (r *notificationDestinationRepository) MarkSent(ctx context.Context, id uuid.UUID, teamsMessageID string, sentAt time.Time) error {
	query := `UPDATE notification_destinations SET status = 'sent', teams_message_id = $2, sent_at = $3, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, teamsMessageID, sentAt)
	return err
}

// GetByNotificationID retrieves destinations for a notification
func (r *notificationDestinationRepository) GetByNotificationID(ctx context.Context, notificationID uuid.UUID) ([]*models.NotificationDestination, error) {
	var rows []*models.NotificationDestination
	query := `SELECT * FROM notification_destinations WHERE notification_id = $1 ORDER BY created_at ASC`
	err := r.db.SelectContext(ctx, &rows, query, notificationID)
	return rows, err
}

// GetByID retrieves a notification destination by ID
func (r *notificationDestinationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.NotificationDestination, error) {
	var row models.NotificationDestination
	query := `SELECT * FROM notification_destinations WHERE id = $1`
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// GetRetryReady retrieves retry-ready destinations
func (r *notificationDestinationRepository) GetRetryReady(ctx context.Context, limit int) ([]*models.NotificationDestination, error) {
	var rows []*models.NotificationDestination
	query := `SELECT * FROM notification_destinations 
              WHERE status IN ('pending','failed') 
                AND (next_retry_at IS NULL OR next_retry_at <= NOW())
                AND retry_count < max_retries
              ORDER BY COALESCE(next_retry_at, NOW()) ASC
              LIMIT $1`
	err := r.db.SelectContext(ctx, &rows, query, limit)
	return rows, err
}

// GetPending retrieves pending destinations without retry logic (for enqueue worker)
func (r *notificationDestinationRepository) GetPending(ctx context.Context, limit int) ([]*models.NotificationDestination, error) {
	var rows []*models.NotificationDestination
	query := `SELECT * FROM notification_destinations
              WHERE status = 'pending'
              ORDER BY created_at ASC
              LIMIT $1`
	err := r.db.SelectContext(ctx, &rows, query, limit)
	return rows, err
}

// UpdateStatus updates only the status field for a destination
func (r *notificationDestinationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE notification_destinations SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// GetByProjectID retrieves notifications by project ID
func (r *notificationRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := `SELECT * FROM notifications WHERE project_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &notifications, query, projectID)
	return notifications, err
}

// GetBySenderID retrieves notifications by sender ID
func (r *notificationRepository) GetBySenderID(ctx context.Context, senderID uuid.UUID) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := `SELECT * FROM notifications WHERE sender_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &notifications, query, senderID)
	return notifications, err
}

// GetByStatus retrieves notifications by status
func (r *notificationRepository) GetByStatus(ctx context.Context, status string) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := `SELECT * FROM notifications WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &notifications, query, status)
	return notifications, err
}

// GetByDateRange retrieves notifications by date range
func (r *notificationRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := `SELECT * FROM notifications WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &notifications, query, start, end)
	return notifications, err
}

// UpdateStatus updates the status of a notification
func (r *notificationRepository) UpdateStatus(ctx context.Context, notificationID uuid.UUID, status models.NotificationStatus, errorMessage string) error {
	query := `UPDATE notifications SET status = $1, error_message = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, string(status), errorMessage, notificationID)
	return err
}

// (Removed) GetPendingNotifications was used by deprecated NotificationScanner
