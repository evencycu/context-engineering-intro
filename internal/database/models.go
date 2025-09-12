package database

import (
	"time"

	"github.com/google/uuid"
)

// Base model with common fields
type BaseModel struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Company represents a company/organization
type Company struct {
	BaseModel
	Name           string `json:"name" db:"name"`
	ContactEmail   string `json:"contact_email" db:"contact_email"`
	ContactPhone   string `json:"contact_phone" db:"contact_phone"`
	Address        string `json:"address" db:"address"`
	Status         string `json:"status" db:"status"`
	BillingEnabled bool   `json:"billing_enabled" db:"billing_enabled"`
}

// User represents a system user
type User struct {
	BaseModel
	CompanyID       uuid.UUID  `json:"company_id" db:"company_id"`
	Email           string     `json:"email" db:"email"`
	Name            string     `json:"name" db:"name"`
	Role            string     `json:"role" db:"role"`
	Status          string     `json:"status" db:"status"`
	LastLoginAt     *time.Time `json:"last_login_at" db:"last_login_at"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	APIKeyHash      string     `json:"-" db:"api_key_hash"`
	APIKeyExpiresAt *time.Time `json:"api_key_expires_at" db:"api_key_expires_at"`
}

// Project represents a project/notification key
type Project struct {
	BaseModel
	CompanyID    uuid.UUID `json:"company_id" db:"company_id"`
	KeyName      string    `json:"key_name" db:"key_name"`
	Description  string    `json:"description" db:"description"`
	Status       string    `json:"status" db:"status"`
	DailyLimit   int       `json:"daily_limit" db:"daily_limit"`
	MonthlyLimit int       `json:"monthly_limit" db:"monthly_limit"`
	Priority     string    `json:"priority" db:"priority"`
	CreatedBy    uuid.UUID `json:"created_by" db:"created_by"`
}

// Destination represents a Teams target
type Destination struct {
	BaseModel
	ProjectID        uuid.UUID  `json:"project_id" db:"project_id"`
	Name             string     `json:"name" db:"name"`
	Type             string     `json:"type" db:"type"` // person, channel, chatgroup
	TeamsUserID      string     `json:"teams_user_id" db:"teams_user_id"`
	TeamsTeamID      string     `json:"teams_team_id" db:"teams_team_id"`
	TeamsChannelID   string     `json:"teams_channel_id" db:"teams_channel_id"`
	TeamsGroupID     string     `json:"teams_group_id" db:"teams_group_id"`
	Status           string     `json:"status" db:"status"`
	ValidationStatus string     `json:"validation_status" db:"validation_status"`
	LastValidatedAt  *time.Time `json:"last_validated_at" db:"last_validated_at"`
	CreatedBy        uuid.UUID  `json:"created_by" db:"created_by"`
}

// Notification represents a notification message
type Notification struct {
	BaseModel
	ProjectID    uuid.UUID      `json:"project_id" db:"project_id"`
	SenderID     uuid.UUID      `json:"sender_id" db:"sender_id"`
	MessageType  string         `json:"message_type" db:"message_type"` // text, file
	Content      string         `json:"content" db:"content"`
	Mentions     []string       `json:"mentions" db:"mentions"`
	Attachment   *Attachment    `json:"attachment" db:"attachment"`
	Priority     string         `json:"priority" db:"priority"`
	Status       string         `json:"status" db:"status"`
	ErrorMessage string         `json:"error_message" db:"error_message"`
	Metadata     map[string]any `json:"metadata" db:"metadata"`
	SentAt       *time.Time     `json:"sent_at" db:"sent_at"`
}

// Attachment represents a file attachment
type Attachment struct {
	FileName string `json:"file_name" db:"file_name"`
	FileURL  string `json:"file_url" db:"file_url"`
	FileSize int64  `json:"file_size" db:"file_size"`
	MimeType string `json:"mime_type" db:"mime_type"`
}

// NotificationDestination represents the relationship between notifications and destinations
type NotificationDestination struct {
	BaseModel
	NotificationID uuid.UUID  `json:"notification_id" db:"notification_id"`
	DestinationID  uuid.UUID  `json:"destination_id" db:"destination_id"`
	Status         string     `json:"status" db:"status"`
	ErrorMessage   string     `json:"error_message" db:"error_message"`
	TeamsMessageID string     `json:"teams_message_id" db:"teams_message_id"`
	SentAt         *time.Time `json:"sent_at" db:"sent_at"`
}

// BillingPlan represents a billing plan
type BillingPlan struct {
	BaseModel
	Name                 string  `json:"name" db:"name"`
	Description          string  `json:"description" db:"description"`
	PricePerNotification float64 `json:"price_per_notification" db:"price_per_notification"`
	PricePerAttachment   float64 `json:"price_per_attachment" db:"price_per_attachment"`
	PricePerMention      float64 `json:"price_per_mention" db:"price_per_mention"`
	DailyLimit           *int    `json:"daily_limit" db:"daily_limit"`
	MonthlyLimit         *int    `json:"monthly_limit" db:"monthly_limit"`
	Status               string  `json:"status" db:"status"`
}

// CompanyBilling represents company billing information
type CompanyBilling struct {
	BaseModel
	CompanyID     uuid.UUID `json:"company_id" db:"company_id"`
	BillingPlanID uuid.UUID `json:"billing_plan_id" db:"billing_plan_id"`
	Status        string    `json:"status" db:"status"`
	BillingEmail  string    `json:"billing_email" db:"billing_email"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	Currency      string    `json:"currency" db:"currency"`
}

// UsageRecord represents a usage record for billing
type UsageRecord struct {
	BaseModel
	CompanyID      uuid.UUID  `json:"company_id" db:"company_id"`
	ProjectID      *uuid.UUID `json:"project_id" db:"project_id"`
	UserID         *uuid.UUID `json:"user_id" db:"user_id"`
	NotificationID *uuid.UUID `json:"notification_id" db:"notification_id"`
	RecordType     string     `json:"record_type" db:"record_type"` // notification, attachment, mention
	Quantity       int        `json:"quantity" db:"quantity"`
	UnitPrice      float64    `json:"unit_price" db:"unit_price"`
	TotalCost      float64    `json:"total_cost" db:"total_cost"`
	BillingPeriod  time.Time  `json:"billing_period" db:"billing_period"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	BaseModel
	CompanyID    *uuid.UUID     `json:"company_id" db:"company_id"`
	UserID       *uuid.UUID     `json:"user_id" db:"user_id"`
	Action       string         `json:"action" db:"action"`
	ResourceType string         `json:"resource_type" db:"resource_type"`
	ResourceID   *uuid.UUID     `json:"resource_id" db:"resource_id"`
	OldValues    map[string]any `json:"old_values" db:"old_values"`
	NewValues    map[string]any `json:"new_values" db:"new_values"`
	IPAddress    string         `json:"ip_address" db:"ip_address"`
	UserAgent    string         `json:"user_agent" db:"user_agent"`
}

// SystemLog represents a system log entry
type SystemLog struct {
	BaseModel
	Level   string         `json:"level" db:"level"`
	Service string         `json:"service" db:"service"`
	Message string         `json:"message" db:"message"`
	Context map[string]any `json:"context" db:"context"`
	TraceID string         `json:"trace_id" db:"trace_id"`
	SpanID  string         `json:"span_id" db:"span_id"`
}

// SystemSetting represents a system setting
type SystemSetting struct {
	BaseModel
	Key         string         `json:"key" db:"key"`
	Value       map[string]any `json:"value" db:"value"`
	Description string         `json:"description" db:"description"`
	IsPublic    bool           `json:"is_public" db:"is_public"`
}

// FeatureFlag represents a feature flag
type FeatureFlag struct {
	BaseModel
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Enabled     bool       `json:"enabled" db:"enabled"`
	CompanyID   *uuid.UUID `json:"company_id" db:"company_id"`
}

// =============================================
// View Models
// =============================================

// NotificationSummary represents a notification summary view
type NotificationSummary struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	ProjectID        uuid.UUID  `json:"project_id" db:"project_id"`
	KeyName          string     `json:"key_name" db:"key_name"`
	CompanyID        uuid.UUID  `json:"company_id" db:"company_id"`
	CompanyName      string     `json:"company_name" db:"company_name"`
	SenderEmail      string     `json:"sender_email" db:"sender_email"`
	MessageType      string     `json:"message_type" db:"message_type"`
	Content          string     `json:"content" db:"content"`
	Priority         string     `json:"priority" db:"priority"`
	Status           string     `json:"status" db:"status"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	SentAt           *time.Time `json:"sent_at" db:"sent_at"`
	DestinationCount int        `json:"destination_count" db:"destination_count"`
	SentCount        int        `json:"sent_count" db:"sent_count"`
	FailedCount      int        `json:"failed_count" db:"failed_count"`
}

// UsageSummary represents a usage summary view
type UsageSummary struct {
	CompanyID     uuid.UUID `json:"company_id" db:"company_id"`
	CompanyName   string    `json:"company_name" db:"company_name"`
	BillingPeriod time.Time `json:"billing_period" db:"billing_period"`
	RecordType    string    `json:"record_type" db:"record_type"`
	TotalQuantity int       `json:"total_quantity" db:"total_quantity"`
	TotalCost     float64   `json:"total_cost" db:"total_cost"`
	ProjectsUsed  int       `json:"projects_used" db:"projects_used"`
}

// =============================================
// Request/Response Models
// =============================================

// CreateNotificationRequest represents a request to create a notification
type CreateNotificationRequest struct {
	ProjectID   uuid.UUID      `json:"project_id" validate:"required"`
	MessageType string         `json:"message_type" validate:"required,oneof=text file"`
	Content     string         `json:"content" validate:"required,max=4000"`
	Mentions    []string       `json:"mentions"`
	Attachment  *Attachment    `json:"attachment"`
	Priority    string         `json:"priority" validate:"oneof=low normal high urgent"`
	Metadata    map[string]any `json:"metadata"`
}

// NotificationResponse represents a notification response
type NotificationResponse struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// NotificationStatusResponse represents a notification status response
type NotificationStatusResponse struct {
	ID           uuid.UUID  `json:"id"`
	Status       string     `json:"status"`
	Message      string     `json:"message"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	SentAt       *time.Time `json:"sent_at"`
	ErrorMessage string     `json:"error_message"`
}

// BillingSummaryResponse represents a billing summary response
type BillingSummaryResponse struct {
	Period             string             `json:"period"`
	TotalNotifications int                `json:"total_notifications"`
	TotalCost          float64            `json:"total_cost"`
	Breakdown          []BillingBreakdown `json:"breakdown"`
}

// BillingBreakdown represents billing breakdown by project
type BillingBreakdown struct {
	ProjectKey string  `json:"project_key"`
	Count      int     `json:"count"`
	Cost       float64 `json:"cost"`
}
