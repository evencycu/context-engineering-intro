package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository interface defines common database operations
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
	Repository[database.Company]
	GetByEmail(ctx context.Context, email string) (*database.Company, error)
	GetByStatus(ctx context.Context, status string) ([]*database.Company, error)
}

// BotRepository defines bot-specific operations
type BotRepository interface {
	Repository[database.PlatformBot]
	GetByAppID(ctx context.Context, appID string) (*database.PlatformBot, error)
	GetByStatus(ctx context.Context, status string) ([]*database.PlatformBot, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*database.PlatformBot, error)
}

// DestinationRepository defines destination-specific operations
type DestinationRepository interface {
	Repository[database.Destination]
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*database.Destination, error)
	GetByBotID(ctx context.Context, botID uuid.UUID) ([]*database.Destination, error)
	GetByStatus(ctx context.Context, status string) ([]*database.Destination, error)
	GetByValidationStatus(ctx context.Context, validationStatus string) ([]*database.Destination, error)
	SearchDestinations(ctx context.Context, query string) ([]*database.Destination, error)
}

// NotificationRepository defines notification-specific operations
type NotificationRepository interface {
	Repository[database.Notification]
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*database.Notification, error)
	GetBySenderID(ctx context.Context, senderID uuid.UUID) ([]*database.Notification, error)
	GetByStatus(ctx context.Context, status string) ([]*database.Notification, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*database.Notification, error)
}

// BaseRepository provides common repository functionality
type BaseRepository struct {
	// This will be implemented with actual database connection
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
	Repository[database.User]
	GetByEmail(ctx context.Context, email string) (*database.User, error)
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.User, error)
	GetByRole(ctx context.Context, role string) ([]*database.User, error)
	GetByStatus(ctx context.Context, status string) ([]*database.User, error)
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
func (r *userRepository) Create(ctx context.Context, entity *database.User) error {
	query := `INSERT INTO users (id, company_id, email, name, password_hash, role, status, api_key_hash, api_key_expires_at, last_login_at, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.Email, entity.Name,
		entity.PasswordHash, entity.Role, entity.Status, entity.APIKeyHash,
		entity.APIKeyExpiresAt, entity.LastLoginAt, entity.CreatedAt, entity.UpdatedAt)
	return err
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.User, error) {
	var user database.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, entity *database.User) error {
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
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*database.User, error) {
	var users []*database.User
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
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*database.User, error) {
	var user database.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByCompanyID retrieves users by company ID
func (r *userRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.User, error) {
	var users []*database.User
	query := `SELECT * FROM users WHERE company_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &users, query, companyID)
	return users, err
}

// GetByRole retrieves users by role
func (r *userRepository) GetByRole(ctx context.Context, role string) ([]*database.User, error) {
	var users []*database.User
	query := `SELECT * FROM users WHERE role = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &users, query, role)
	return users, err
}

// GetByStatus retrieves users by status
func (r *userRepository) GetByStatus(ctx context.Context, status string) ([]*database.User, error) {
	var users []*database.User
	query := `SELECT * FROM users WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &users, query, status)
	return users, err
}

// companyRepository implements CompanyRepository
type companyRepository struct {
	db *sqlx.DB
}

func (r *companyRepository) Create(ctx context.Context, entity *database.Company) error {
	query := `INSERT INTO companies (id, name, contact_email, contact_phone, address, status, billing_enabled, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.Name, entity.ContactEmail, entity.ContactPhone,
		entity.Address, entity.Status, entity.BillingEnabled, entity.CreatedAt, entity.UpdatedAt)
	return err
}

func (r *companyRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Company, error) {
	var company database.Company
	query := `SELECT * FROM companies WHERE id = $1`
	err := r.db.GetContext(ctx, &company, query, id)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) Update(ctx context.Context, entity *database.Company) error {
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

func (r *companyRepository) List(ctx context.Context, limit, offset int) ([]*database.Company, error) {
	var companies []*database.Company
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

func (r *companyRepository) GetByEmail(ctx context.Context, email string) (*database.Company, error) {
	var company database.Company
	query := `SELECT * FROM companies WHERE contact_email = $1`
	err := r.db.GetContext(ctx, &company, query, email)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) GetByStatus(ctx context.Context, status string) ([]*database.Company, error) {
	var companies []*database.Company
	query := `SELECT * FROM companies WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &companies, query, status)
	return companies, err
}

// ProjectRepository defines the interface for project operations
type ProjectRepository interface {
	Repository[database.Project]
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.Project, error)
	GetByKeyName(ctx context.Context, keyName string) (*database.Project, error)
	GetByStatus(ctx context.Context, status string) ([]*database.Project, error)
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
func (r *projectRepository) Create(ctx context.Context, entity *database.Project) error {
	query := `INSERT INTO projects (id, company_id, key_name, description, status, daily_limit, monthly_limit, priority, created_by, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.KeyName, entity.Description,
		entity.Status, entity.DailyLimit, entity.MonthlyLimit, entity.Priority,
		entity.CreatedBy, entity.CreatedAt, entity.UpdatedAt)
	return err
}

// GetByID retrieves a project by ID
func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Project, error) {
	var project database.Project
	query := `SELECT * FROM projects WHERE id = $1`
	err := r.db.GetContext(ctx, &project, query, id)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// Update updates an existing project
func (r *projectRepository) Update(ctx context.Context, entity *database.Project) error {
	query := `UPDATE projects SET company_id = $2, key_name = $3, description = $4, 
			  status = $5, daily_limit = $6, monthly_limit = $7, priority = $8, updated_at = $9 
			  WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.KeyName, entity.Description,
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
func (r *projectRepository) List(ctx context.Context, limit, offset int) ([]*database.Project, error) {
	var projects []*database.Project
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
func (r *projectRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.Project, error) {
	var projects []*database.Project
	query := `SELECT * FROM projects WHERE company_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &projects, query, companyID)
	return projects, err
}

// GetByKeyName retrieves a project by key name
func (r *projectRepository) GetByKeyName(ctx context.Context, keyName string) (*database.Project, error) {
	var project database.Project
	query := `SELECT * FROM projects WHERE key_name = $1`
	err := r.db.GetContext(ctx, &project, query, keyName)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByStatus retrieves projects by status
func (r *projectRepository) GetByStatus(ctx context.Context, status string) ([]*database.Project, error) {
	var projects []*database.Project
	query := `SELECT * FROM projects WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &projects, query, status)
	return projects, err
}

// PlatformBotRepository defines the interface for platform bot operations
type PlatformBotRepository interface {
	Repository[database.PlatformBot]
	GetByAppID(ctx context.Context, appID string) (*database.PlatformBot, error)
	GetByStatus(ctx context.Context, status string) ([]*database.PlatformBot, error)
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.PlatformBot, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*database.PlatformBot, error)
}

// platformBotRepository implements PlatformBotRepository
type platformBotRepository struct {
	db *sqlx.DB
}

// NewPlatformBotRepository creates a new platform bot repository
func NewPlatformBotRepository(db *sqlx.DB) PlatformBotRepository {
	return &platformBotRepository{db: db}
}

// Create creates a new platform bot
func (r *platformBotRepository) Create(ctx context.Context, entity *database.PlatformBot) error {
	query := `INSERT INTO platform_bots (id, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	// Convert capabilities to JSONB
	capabilitiesJSON, err := json.Marshal(entity.Capabilities)
	if err != nil {
		return fmt.Errorf("failed to marshal capabilities: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.Name, entity.Description, entity.AppID, entity.AppPasswordHash,
		entity.TenantID, entity.Status, entity.WebhookURL, capabilitiesJSON,
		entity.RateLimitPerMinute, entity.MaxConcurrentRequests, entity.CreatedAt, entity.UpdatedAt)
	return err
}

// GetByID retrieves a platform bot by ID
func (r *platformBotRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.PlatformBot, error) {
	var bot database.PlatformBot
	query := `SELECT * FROM platform_bots WHERE id = $1`
	err := r.db.GetContext(ctx, &bot, query, id)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

// Update updates an existing platform bot
func (r *platformBotRepository) Update(ctx context.Context, entity *database.PlatformBot) error {
	query := `UPDATE platform_bots SET name = $2, description = $3, app_id = $4, app_password_hash = $5, 
			  tenant_id = $6, status = $7, webhook_url = $8, capabilities = $9, rate_limit_per_minute = $10, 
			  max_concurrent_requests = $11, updated_at = $12 
			  WHERE id = $1`

	// Convert capabilities to JSONB
	capabilitiesJSON, err := json.Marshal(entity.Capabilities)
	if err != nil {
		return fmt.Errorf("failed to marshal capabilities: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.Name, entity.Description, entity.AppID, entity.AppPasswordHash,
		entity.TenantID, entity.Status, entity.WebhookURL, capabilitiesJSON,
		entity.RateLimitPerMinute, entity.MaxConcurrentRequests, entity.UpdatedAt)
	return err
}

// Delete deletes a platform bot by ID
func (r *platformBotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM platform_bots WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves platform bots with pagination
func (r *platformBotRepository) List(ctx context.Context, limit, offset int) ([]*database.PlatformBot, error) {
	var bots []*database.PlatformBot
	query := `SELECT * FROM platform_bots ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &bots, query, limit, offset)
	return bots, err
}

// Count returns the total number of platform bots
func (r *platformBotRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM platform_bots`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

// GetByAppID retrieves a platform bot by app ID
func (r *platformBotRepository) GetByAppID(ctx context.Context, appID string) (*database.PlatformBot, error) {
	var bot database.PlatformBot
	query := `SELECT * FROM platform_bots WHERE app_id = $1`
	err := r.db.GetContext(ctx, &bot, query, appID)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

// GetByStatus retrieves platform bots by status
func (r *platformBotRepository) GetByStatus(ctx context.Context, status string) ([]*database.PlatformBot, error) {
	var bots []*database.PlatformBot
	query := `SELECT * FROM platform_bots WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &bots, query, status)
	return bots, err
}

// GetByCompanyID retrieves platform bots by company ID
func (r *platformBotRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.PlatformBot, error) {
	var bots []*database.PlatformBot
	query := `SELECT * FROM platform_bots WHERE company_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &bots, query, companyID)
	return bots, err
}

// GetByTenantID retrieves platform bots by tenant ID
func (r *platformBotRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*database.PlatformBot, error) {
	var bots []*database.PlatformBot
	query := `SELECT * FROM platform_bots WHERE tenant_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &bots, query, tenantID)
	return bots, err
}

// ThirdPartyBotRepository defines the interface for third-party bot operations
type ThirdPartyBotRepository interface {
	Repository[database.ThirdPartyBot]
	GetByAppID(ctx context.Context, appID string) (*database.ThirdPartyBot, error)
	GetByStatus(ctx context.Context, status string) ([]*database.ThirdPartyBot, error)
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.ThirdPartyBot, error)
}

// thirdPartyBotRepository implements ThirdPartyBotRepository
type thirdPartyBotRepository struct {
	db *sqlx.DB
}

// NewThirdPartyBotRepository creates a new third-party bot repository
func NewThirdPartyBotRepository(db *sqlx.DB) ThirdPartyBotRepository {
	return &thirdPartyBotRepository{db: db}
}

// Create creates a new third-party bot
func (r *thirdPartyBotRepository) Create(ctx context.Context, entity *database.ThirdPartyBot) error {
	query := `INSERT INTO third_party_bots (id, company_id, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, api_endpoint, api_key_hash, capabilities, rate_limit_per_minute, max_concurrent_requests, contact_email, contact_phone, created_by, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	// Convert capabilities to JSONB
	capabilitiesJSON, err := json.Marshal(entity.Capabilities)
	if err != nil {
		return fmt.Errorf("failed to marshal capabilities: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.Name, entity.Description, entity.AppID, entity.AppPasswordHash,
		entity.TenantID, entity.Status, entity.WebhookURL, entity.APIEndpoint, entity.APIKeyHash,
		capabilitiesJSON, entity.RateLimitPerMinute, entity.MaxConcurrentRequests,
		entity.ContactEmail, entity.ContactPhone, entity.CreatedBy, entity.CreatedAt, entity.UpdatedAt)
	return err
}

// GetByID retrieves a third-party bot by ID
func (r *thirdPartyBotRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.ThirdPartyBot, error) {
	var bot database.ThirdPartyBot
	query := `SELECT * FROM third_party_bots WHERE id = $1`
	err := r.db.GetContext(ctx, &bot, query, id)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

// Update updates an existing third-party bot
func (r *thirdPartyBotRepository) Update(ctx context.Context, entity *database.ThirdPartyBot) error {
	query := `UPDATE third_party_bots SET company_id = $2, name = $3, description = $4, app_id = $5, app_password_hash = $6, 
			  tenant_id = $7, status = $8, webhook_url = $9, api_endpoint = $10, api_key_hash = $11, capabilities = $12, 
			  rate_limit_per_minute = $13, max_concurrent_requests = $14, contact_email = $15, contact_phone = $16, updated_at = $17 
			  WHERE id = $1`

	// Convert capabilities to JSONB
	capabilitiesJSON, err := json.Marshal(entity.Capabilities)
	if err != nil {
		return fmt.Errorf("failed to marshal capabilities: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.Name, entity.Description, entity.AppID, entity.AppPasswordHash,
		entity.TenantID, entity.Status, entity.WebhookURL, entity.APIEndpoint, entity.APIKeyHash,
		capabilitiesJSON, entity.RateLimitPerMinute, entity.MaxConcurrentRequests,
		entity.ContactEmail, entity.ContactPhone, entity.UpdatedAt)
	return err
}

// Delete deletes a third-party bot by ID
func (r *thirdPartyBotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM third_party_bots WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves third-party bots with pagination
func (r *thirdPartyBotRepository) List(ctx context.Context, limit, offset int) ([]*database.ThirdPartyBot, error) {
	var bots []*database.ThirdPartyBot
	query := `SELECT * FROM third_party_bots ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &bots, query, limit, offset)
	return bots, err
}

// Count returns the total number of third-party bots
func (r *thirdPartyBotRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM third_party_bots`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

// GetByAppID retrieves a third-party bot by app ID
func (r *thirdPartyBotRepository) GetByAppID(ctx context.Context, appID string) (*database.ThirdPartyBot, error) {
	var bot database.ThirdPartyBot
	query := `SELECT * FROM third_party_bots WHERE app_id = $1`
	err := r.db.GetContext(ctx, &bot, query, appID)
	if err != nil {
		return nil, err
	}
	return &bot, nil
}

// GetByStatus retrieves third-party bots by status
func (r *thirdPartyBotRepository) GetByStatus(ctx context.Context, status string) ([]*database.ThirdPartyBot, error) {
	var bots []*database.ThirdPartyBot
	query := `SELECT * FROM third_party_bots WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &bots, query, status)
	return bots, err
}

// GetByCompanyID retrieves third-party bots by company ID
func (r *thirdPartyBotRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) ([]*database.ThirdPartyBot, error) {
	var bots []*database.ThirdPartyBot
	query := `SELECT * FROM third_party_bots WHERE company_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &bots, query, companyID)
	return bots, err
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
func (r *destinationRepository) Create(ctx context.Context, entity *database.Destination) error {
	// Marshal targets to JSONB
	targetsJSON, err := json.Marshal(entity.Targets)
	if err != nil {
		return fmt.Errorf("failed to marshal targets: %w", err)
	}

	query := `INSERT INTO destinations (id, project_id, name, description, teams_tenant_id, targets, bot_id, bot_type, status, validation_status, last_validated_at, created_by, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.Name, entity.Description, entity.TeamsTenantID,
		targetsJSON, entity.BotID, entity.BotType, entity.Status, entity.ValidationStatus,
		entity.LastValidatedAt, entity.CreatedBy, entity.CreatedAt, entity.UpdatedAt)

	return err
}

// GetByID retrieves a destination by ID
func (r *destinationRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Destination, error) {
	var destination database.Destination
	query := `SELECT * FROM destinations WHERE id = $1`
	err := r.db.GetContext(ctx, &destination, query, id)
	if err != nil {
		return nil, err
	}
	return &destination, nil
}

// Update updates a destination
func (r *destinationRepository) Update(ctx context.Context, entity *database.Destination) error {
	// Marshal targets to JSONB
	targetsJSON, err := json.Marshal(entity.Targets)
	if err != nil {
		return fmt.Errorf("failed to marshal targets: %w", err)
	}

	query := `UPDATE destinations SET 
			  project_id = $2, name = $3, description = $4, teams_tenant_id = $5, 
			  targets = $6, bot_id = $7, bot_type = $8, status = $9, 
			  validation_status = $10, last_validated_at = $11, updated_at = $12
			  WHERE id = $1`

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.Name, entity.Description, entity.TeamsTenantID,
		targetsJSON, entity.BotID, entity.BotType, entity.Status, entity.ValidationStatus,
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
func (r *destinationRepository) List(ctx context.Context, limit, offset int) ([]*database.Destination, error) {
	var destinations []*database.Destination
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
func (r *destinationRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*database.Destination, error) {
	var destinations []*database.Destination
	query := `SELECT * FROM destinations WHERE project_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &destinations, query, projectID)
	return destinations, err
}

// GetByBotID retrieves destinations by bot ID
func (r *destinationRepository) GetByBotID(ctx context.Context, botID uuid.UUID) ([]*database.Destination, error) {
	var destinations []*database.Destination
	query := `SELECT * FROM destinations WHERE bot_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &destinations, query, botID)
	return destinations, err
}

// GetByStatus retrieves destinations by status
func (r *destinationRepository) GetByStatus(ctx context.Context, status string) ([]*database.Destination, error) {
	var destinations []*database.Destination
	query := `SELECT * FROM destinations WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &destinations, query, status)
	return destinations, err
}

// GetByValidationStatus retrieves destinations by validation status
func (r *destinationRepository) GetByValidationStatus(ctx context.Context, validationStatus string) ([]*database.Destination, error) {
	var destinations []*database.Destination
	query := `SELECT * FROM destinations WHERE validation_status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &destinations, query, validationStatus)
	return destinations, err
}

// SearchDestinations searches destinations by name or description
func (r *destinationRepository) SearchDestinations(ctx context.Context, query string) ([]*database.Destination, error) {
	var destinations []*database.Destination
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
func (r *notificationRepository) Create(ctx context.Context, entity *database.Notification) error {
	// Marshal JSONB fields
	mentionsJSON, err := json.Marshal([]string(entity.Mentions))
	if err != nil {
		return fmt.Errorf("failed to marshal mentions: %w", err)
	}

	var attachmentJSON []byte
	if entity.Attachment != nil {
		attachmentJSON, err = json.Marshal(*entity.Attachment)
		if err != nil {
			return fmt.Errorf("failed to marshal attachment: %w", err)
		}
	}

	var adaptiveCardJSON []byte
	if entity.AdaptiveCard != nil {
		adaptiveCardJSON, err = json.Marshal(*entity.AdaptiveCard)
		if err != nil {
			return fmt.Errorf("failed to marshal adaptive_card: %w", err)
		}
	}

	metadataJSON, err := json.Marshal(map[string]any(entity.Metadata))
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `INSERT INTO notifications (id, project_id, sender_id, message_type, content, mentions, attachment, adaptive_card, priority, status, error_message, metadata, sent_at, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	// Handle nil values for JSONB fields
	var attachmentJSONValue interface{} = attachmentJSON
	if len(attachmentJSON) == 0 {
		attachmentJSONValue = nil
	}

	var adaptiveCardJSONValue interface{} = adaptiveCardJSON
	if len(adaptiveCardJSON) == 0 {
		adaptiveCardJSONValue = nil
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.SenderID, entity.MessageType, entity.Content,
		mentionsJSON, attachmentJSONValue, adaptiveCardJSONValue, entity.Priority, entity.Status,
		entity.ErrorMessage, metadataJSON, entity.SentAt, entity.CreatedAt, entity.UpdatedAt)

	return err
}

// GetByID retrieves a notification by ID
func (r *notificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.Notification, error) {
	var notification database.Notification
	query := `SELECT * FROM notifications WHERE id = $1`
	err := r.db.GetContext(ctx, &notification, query, id)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// Update updates a notification
func (r *notificationRepository) Update(ctx context.Context, entity *database.Notification) error {
	// Marshal JSONB fields
	mentionsJSON, err := json.Marshal([]string(entity.Mentions))
	if err != nil {
		return fmt.Errorf("failed to marshal mentions: %w", err)
	}

	var attachmentJSON []byte
	if entity.Attachment != nil {
		attachmentJSON, err = json.Marshal(*entity.Attachment)
		if err != nil {
			return fmt.Errorf("failed to marshal attachment: %w", err)
		}
	}

	var adaptiveCardJSON []byte
	if entity.AdaptiveCard != nil {
		adaptiveCardJSON, err = json.Marshal(*entity.AdaptiveCard)
		if err != nil {
			return fmt.Errorf("failed to marshal adaptive_card: %w", err)
		}
	}

	metadataJSON, err := json.Marshal(map[string]any(entity.Metadata))
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `UPDATE notifications SET 
			  project_id = $2, sender_id = $3, message_type = $4, content = $5, 
			  mentions = $6, attachment = $7, adaptive_card = $8, priority = $9, 
			  status = $10, error_message = $11, metadata = $12, sent_at = $13, updated_at = $14
			  WHERE id = $1`

	// Handle nil values for JSONB fields
	var attachmentJSONValue interface{} = attachmentJSON
	if len(attachmentJSON) == 0 {
		attachmentJSONValue = nil
	}

	var adaptiveCardJSONValue interface{} = adaptiveCardJSON
	if len(adaptiveCardJSON) == 0 {
		adaptiveCardJSONValue = nil
	}

	_, err = r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.SenderID, entity.MessageType, entity.Content,
		mentionsJSON, attachmentJSONValue, adaptiveCardJSONValue, entity.Priority, entity.Status,
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
func (r *notificationRepository) List(ctx context.Context, limit, offset int) ([]*database.Notification, error) {
	var notifications []*database.Notification
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

// GetByProjectID retrieves notifications by project ID
func (r *notificationRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*database.Notification, error) {
	var notifications []*database.Notification
	query := `SELECT * FROM notifications WHERE project_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &notifications, query, projectID)
	return notifications, err
}

// GetBySenderID retrieves notifications by sender ID
func (r *notificationRepository) GetBySenderID(ctx context.Context, senderID uuid.UUID) ([]*database.Notification, error) {
	var notifications []*database.Notification
	query := `SELECT * FROM notifications WHERE sender_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &notifications, query, senderID)
	return notifications, err
}

// GetByStatus retrieves notifications by status
func (r *notificationRepository) GetByStatus(ctx context.Context, status string) ([]*database.Notification, error) {
	var notifications []*database.Notification
	query := `SELECT * FROM notifications WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &notifications, query, status)
	return notifications, err
}

// GetByDateRange retrieves notifications by date range
func (r *notificationRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*database.Notification, error) {
	var notifications []*database.Notification
	query := `SELECT * FROM notifications WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &notifications, query, start, end)
	return notifications, err
}
