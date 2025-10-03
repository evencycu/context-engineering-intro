package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// Base model with common fields
type BaseModel struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// BotType represents the type of bot
type BotType string

const (
	BotTypePlatform   BotType = "platform"
	BotTypeThirdParty BotType = "third_party"
)

// BotStatus represents the status of a bot
type BotStatus string

const (
	BotStatusActive      BotStatus = "active"
	BotStatusInactive    BotStatus = "inactive"
	BotStatusSuspended   BotStatus = "suspended"
	BotStatusMaintenance BotStatus = "maintenance"
)

// =============================================

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
	NotifyKey    string    `json:"notify_key" db:"notify_key"`
	Description  string    `json:"description" db:"description"`
	Status       string    `json:"status" db:"status"`
	DailyLimit   int       `json:"daily_limit" db:"daily_limit"`
	MonthlyLimit int       `json:"monthly_limit" db:"monthly_limit"`
	Priority     string    `json:"priority" db:"priority"`
	CreatedBy    uuid.UUID `json:"created_by" db:"created_by"`
}

// =============================================
// Bot Management Models
// =============================================

// TeamsBot represents a platform bot or third-party bot
type TeamsBot struct {
	BaseModel
	Type                  string       `json:"type" db:"type"`             // 'platform' or 'third_party'
	CompanyID             *uuid.UUID   `json:"company_id" db:"company_id"` // Only for third_party bots
	Name                  string       `json:"name" db:"name"`
	Description           *string      `json:"description" db:"description"`
	AppID                 string       `json:"app_id" db:"app_id"`
	AppPasswordHash       string       `json:"-" db:"app_password_hash"`
	TenantID              *string      `json:"tenant_id" db:"tenant_id"`
	Status                *BotStatus   `json:"status" db:"status"`
	WebhookURL            *string      `json:"webhook_url" db:"webhook_url"`
	Capabilities          *JSONBObject `json:"capabilities" db:"capabilities"`
	RateLimitPerMinute    *int         `json:"rate_limit_per_minute" db:"rate_limit_per_minute"`
	MaxConcurrentRequests *int         `json:"max_concurrent_requests" db:"max_concurrent_requests"`
	APIEndpoint           *string      `json:"api_endpoint" db:"api_endpoint"`   // Only for third_party bots
	APIKeyHash            *string      `json:"-" db:"api_key_hash"`              // Only for third_party bots
	ContactEmail          *string      `json:"contact_email" db:"contact_email"` // Only for third_party bots
	ContactPhone          *string      `json:"contact_phone" db:"contact_phone"` // Only for third_party bots
	CreatedBy             *uuid.UUID   `json:"created_by" db:"created_by"`       // Only for third_party bots
}

// BotInstallation represents a bot installation
type BotInstallation struct {
	BaseModel
	BotID         uuid.UUID `json:"bot_id" db:"bot_id"`
	BotType       BotType   `json:"bot_type" db:"bot_type"`
	TeamsTenantID string    `json:"teams_tenant_id" db:"teams_tenant_id"`
	// Conversation type and identification
	ConversationType string `json:"conversation_type" db:"conversation_type"` // personal, channel, groupChat
	ConversationID   string `json:"conversation_id" db:"conversation_id"`
	ServiceURL       string `json:"service_url" db:"service_url"`
	// Bot and user identification
	RecipientID     string `json:"recipient_id" db:"recipient_id"`
	RecipientName   string `json:"recipient_name" db:"recipient_name"`
	FromID          string `json:"from_id" db:"from_id"`
	FromName        string `json:"from_name" db:"from_name"`
	FromAADObjectID string `json:"from_aad_object_id" db:"from_aad_object_id"`
	Email           string `json:"email" db:"email"`                       // User email address
	DescriptionName string `json:"description_name" db:"description_name"` // User description or display name
	// Status and lifecycle
	InstallationStatus string      `json:"installation_status" db:"installation_status"`
	InstalledAt        time.Time   `json:"installed_at" db:"installed_at"`
	UninstalledAt      *time.Time  `json:"uninstalled_at" db:"uninstalled_at"`
	LastActivityAt     *time.Time  `json:"last_activity_at" db:"last_activity_at"`
	Metadata           JSONBObject `json:"metadata" db:"metadata"`
}

// =============================================
// Destination Management Models
// =============================================

// TeamsTarget represents a single Teams target
type TeamsTarget struct {
	Type           string `json:"type"` // person, channel, chatgroup
	Email          string `json:"email,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"` // Direct conversation ID
	DisplayName    string `json:"display_name,omitempty"`
	TenantID       string `json:"tenant_id,omitempty"` // AAD Tenant ID
}

// JSONBTargets is a JSONB-backed slice of TeamsTarget for DB scanning/valuing
type JSONBTargets []TeamsTarget

// JSONBStringArray is a JSONB-backed slice of strings for DB scanning/valuing
type JSONBStringArray []string

// JSONBObject is a JSONB-backed object for DB scanning/valuing
type JSONBObject map[string]any

// JSONBNullableObject is a JSONB-backed nullable object for DB scanning/valuing
type JSONBNullableObject map[string]any

// Value implements driver.Valuer to convert JSONBTargets to JSON bytes
func (t JSONBTargets) Value() (driver.Value, error) {
	if t == nil {
		return []byte("[]"), nil
	}
	b, err := json.Marshal([]TeamsTarget(t))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONBTargets: %w", err)
	}
	log.Println("Value() JSONBTargets: ", string(b))
	return b, nil
}

// Scan implements sql.Scanner to convert JSON bytes into JSONBTargets
func (t *JSONBTargets) Scan(src any) error {
	if src == nil {
		*t = JSONBTargets{}
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONBTargets Scan: %T", src)
	}
	log.Println("Scan() JSONBTargets: ", string(data))
	var arr []TeamsTarget
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("failed to unmarshal JSONBTargets: %w", err)
	}
	*t = JSONBTargets(arr)
	return nil
}

// Value implements driver.Valuer to convert JSONBStringArray to JSON bytes
func (s JSONBStringArray) Value() (driver.Value, error) {
	if s == nil {
		return []byte("[]"), nil
	}
	b, err := json.Marshal([]string(s))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONBStringArray: %w", err)
	}
	return b, nil
}

// Scan implements sql.Scanner to convert JSON bytes into JSONBStringArray
func (s *JSONBStringArray) Scan(src any) error {
	if src == nil {
		*s = JSONBStringArray{}
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONBStringArray Scan: %T", src)
	}
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("failed to unmarshal JSONBStringArray: %w", err)
	}
	*s = JSONBStringArray(arr)
	return nil
}

// Value implements driver.Valuer to convert JSONBObject to JSON bytes
func (o JSONBObject) Value() (driver.Value, error) {
	if o == nil {
		return []byte("{}"), nil
	}
	b, err := json.Marshal(map[string]any(o))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONBObject: %w", err)
	}
	return b, nil
}

// Scan implements sql.Scanner to convert JSON bytes into JSONBObject
func (o *JSONBObject) Scan(src any) error {
	if src == nil {
		*o = JSONBObject{}
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONBObject Scan: %T", src)
	}
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("failed to unmarshal JSONBObject: %w", err)
	}
	*o = JSONBObject(obj)
	return nil
}

// Value implements driver.Valuer to convert JSONBNullableObject to JSON bytes
func (o JSONBNullableObject) Value() (driver.Value, error) {
	if o == nil {
		return nil, nil
	}
	b, err := json.Marshal(map[string]any(o))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONBNullableObject: %w", err)
	}
	return b, nil
}

// Scan implements sql.Scanner to convert JSON bytes into JSONBNullableObject
func (o *JSONBNullableObject) Scan(src any) error {
	if src == nil {
		*o = nil
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("unsupported type for JSONBNullableObject Scan: %T", src)
	}
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("failed to unmarshal JSONBNullableObject: %w", err)
	}
	*o = JSONBNullableObject(obj)
	return nil
}

// Destination represents a Teams target group (can contain multiple targets)
type Destination struct {
	BaseModel
	ProjectID        uuid.UUID    `json:"project_id" db:"project_id"`
	Name             string       `json:"name" db:"name"`
	Description      string       `json:"description" db:"description"`
	TeamsTenantID    string       `json:"teams_tenant_id" db:"teams_tenant_id"`
	Targets          JSONBTargets `json:"targets" db:"targets"`
	BotID            *uuid.UUID   `json:"bot_id" db:"bot_id"`
	Status           string       `json:"status" db:"status"`
	ValidationStatus string       `json:"validation_status" db:"validation_status"`
	LastValidatedAt  *time.Time   `json:"last_validated_at" db:"last_validated_at"`
	CreatedBy        uuid.UUID    `json:"created_by" db:"created_by"`
}

// =============================================
// Notification System Models
// =============================================

// Notification represents a notification message
type Notification struct {
	BaseModel
	ProjectID    uuid.UUID            `json:"project_id" db:"project_id"`
	SenderID     *uuid.UUID           `json:"sender_id" db:"sender_id"`
	MessageType  string               `json:"message_type" db:"message_type"` // text, file, adaptive_card
	Content      string               `json:"content" db:"content"`
	Mentions     JSONBStringArray     `json:"mentions" db:"mentions"`
	Attachment   *JSONBNullableObject `json:"attachment" db:"attachment"`
	AdaptiveCard *JSONBNullableObject `json:"adaptive_card" db:"adaptive_card"`
	Priority     string               `json:"priority" db:"priority"`
	Status       string               `json:"status" db:"status"`
	ErrorMessage string               `json:"error_message" db:"error_message"`
	Metadata     JSONBObject          `json:"metadata" db:"metadata"`
	SentAt       *time.Time           `json:"sent_at" db:"sent_at"`
}

// Attachment represents a file attachment
type Attachment struct {
	FileName string `json:"file_name" db:"file_name"`
	FileURL  string `json:"file_url" db:"file_url"`
	FileSize int64  `json:"file_size" db:"file_size"`
	MimeType string `json:"mime_type" db:"mime_type"`
}

// AdaptiveCard represents an adaptive card
type AdaptiveCard struct {
	Type    string        `json:"type" db:"type"`
	Version string        `json:"version" db:"version"`
	Body    []CardElement `json:"body" db:"body"`
	Actions []CardAction  `json:"actions" db:"actions"`
}

// CardElement represents an adaptive card element
type CardElement struct {
	Type    string        `json:"type" db:"type"`
	Text    string        `json:"text,omitempty" db:"text"`
	Size    string        `json:"size,omitempty" db:"size"`
	Weight  string        `json:"weight,omitempty" db:"weight"`
	Color   string        `json:"color,omitempty" db:"color"`
	Items   []CardElement `json:"items,omitempty" db:"items"`
	Columns []CardColumn  `json:"columns,omitempty" db:"columns"`
}

// CardColumn represents an adaptive card column
type CardColumn struct {
	Type  string        `json:"type" db:"type"`
	Width string        `json:"width,omitempty" db:"width"`
	Items []CardElement `json:"items" db:"items"`
}

// CardAction represents an adaptive card action
type CardAction struct {
	Type  string `json:"type" db:"type"`
	Title string `json:"title" db:"title"`
	URL   string `json:"url,omitempty" db:"url"`
	Data  string `json:"data,omitempty" db:"data"`
}

// NotificationDestination represents the relationship between notifications and destinations
// Used for actor-based async notification sending with queue management
type NotificationDestination struct {
	BaseModel
	NotificationID uuid.UUID  `json:"notification_id" db:"notification_id"`
	DestinationID  uuid.UUID  `json:"destination_id" db:"destination_id"`
	Priority       string     `json:"priority" db:"priority"`
	ConversationID *string    `json:"conversation_id" db:"conversation_id"` // Teams conversation ID
	BotID          *uuid.UUID `json:"bot_id" db:"bot_id"`
	BotType        *BotType   `json:"bot_type" db:"bot_type"`
	Status         string     `json:"status" db:"status"` // pending, processing, sent, failed, cancelled
	ErrorMessage   *string    `json:"error_message" db:"error_message"`
	TeamsMessageID *string    `json:"teams_message_id" db:"teams_message_id"`
	SentAt         *time.Time `json:"sent_at" db:"sent_at"`
	RetryCount     int        `json:"retry_count" db:"retry_count"`
	MaxRetries     int        `json:"max_retries" db:"max_retries"`
	NextRetryAt    *time.Time `json:"next_retry_at" db:"next_retry_at"` // When to retry (NULL = immediate)
	FirstAttemptAt *time.Time `json:"first_attempt_at" db:"first_attempt_at"`
	LastAttemptAt  *time.Time `json:"last_attempt_at" db:"last_attempt_at"`
	FailureReason  *string    `json:"failure_reason" db:"failure_reason"` // rate_limit, timeout, server_error, etc
	RetryAfter     *int       `json:"retry_after" db:"retry_after"`       // Retry-After from 429 (seconds)
	ActorID        *string    `json:"actor_id" db:"actor_id"`             // ID of actor handling this
}

// =============================================
// Bot Routing and Load Balancing Models
// =============================================

// BotRoutingRule represents a bot routing rule
type BotRoutingRule struct {
	BaseModel
	Name          string         `json:"name" db:"name"`
	Description   string         `json:"description" db:"description"`
	Priority      int            `json:"priority" db:"priority"`
	Conditions    map[string]any `json:"conditions" db:"conditions"`
	TargetBotID   *uuid.UUID     `json:"target_bot_id" db:"target_bot_id"`
	TargetBotType *BotType       `json:"target_bot_type" db:"target_bot_type"`
	IsActive      bool           `json:"is_active" db:"is_active"`
}

// BotHealthStatus represents bot health status
type BotHealthStatus struct {
	BaseModel
	BotID          uuid.UUID      `json:"bot_id" db:"bot_id"`
	BotType        BotType        `json:"bot_type" db:"bot_type"`
	Status         string         `json:"status" db:"status"`
	LastCheckAt    time.Time      `json:"last_check_at" db:"last_check_at"`
	ResponseTimeMs int            `json:"response_time_ms" db:"response_time_ms"`
	ErrorCount     int            `json:"error_count" db:"error_count"`
	SuccessCount   int            `json:"success_count" db:"success_count"`
	ErrorMessage   string         `json:"error_message" db:"error_message"`
	Metadata       map[string]any `json:"metadata" db:"metadata"`
}

// =============================================
// Third Party Bot API Management Models
// =============================================

// ThirdPartyBotAPIKey represents an API key for third-party bots
type ThirdPartyBotAPIKey struct {
	BaseModel
	BotID              uuid.UUID      `json:"bot_id" db:"bot_id"`
	NotifyKey          string         `json:"notify_key" db:"notify_key"`
	APIKeyHash         string         `json:"-" db:"api_key_hash"`
	Permissions        map[string]any `json:"permissions" db:"permissions"`
	RateLimitPerMinute int            `json:"rate_limit_per_minute" db:"rate_limit_per_minute"`
	ExpiresAt          *time.Time     `json:"expires_at" db:"expires_at"`
	IsActive           bool           `json:"is_active" db:"is_active"`
	CreatedBy          uuid.UUID      `json:"created_by" db:"created_by"`
}

// ThirdPartyBotAPIUsage represents API usage for third-party bots
type ThirdPartyBotAPIUsage struct {
	BaseModel
	BotID             uuid.UUID  `json:"bot_id" db:"bot_id"`
	APIKeyID          *uuid.UUID `json:"api_key_id" db:"api_key_id"`
	Endpoint          string     `json:"endpoint" db:"endpoint"`
	Method            string     `json:"method" db:"method"`
	StatusCode        *int       `json:"status_code" db:"status_code"`
	ResponseTimeMs    *int       `json:"response_time_ms" db:"response_time_ms"`
	RequestSizeBytes  *int       `json:"request_size_bytes" db:"request_size_bytes"`
	ResponseSizeBytes *int       `json:"response_size_bytes" db:"response_size_bytes"`
	IPAddress         string     `json:"ip_address" db:"ip_address"`
	UserAgent         string     `json:"user_agent" db:"user_agent"`
}

// =============================================
// Billing and Usage Models
// =============================================

// BillingPlan represents a billing plan
type BillingPlan struct {
	BaseModel
	Name                 string  `json:"name" db:"name"`
	Description          string  `json:"description" db:"description"`
	PricePerNotification float64 `json:"price_per_notification" db:"price_per_notification"`
	PricePerAttachment   float64 `json:"price_per_attachment" db:"price_per_attachment"`
	PricePerMention      float64 `json:"price_per_mention" db:"price_per_mention"`
	PricePerAdaptiveCard float64 `json:"price_per_adaptive_card" db:"price_per_adaptive_card"`
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
	BotID          *uuid.UUID `json:"bot_id" db:"bot_id"`
	BotType        *BotType   `json:"bot_type" db:"bot_type"`
	RecordType     string     `json:"record_type" db:"record_type"` // notification, attachment, mention, adaptive_card
	Quantity       int        `json:"quantity" db:"quantity"`
	UnitPrice      float64    `json:"unit_price" db:"unit_price"`
	TotalCost      float64    `json:"total_cost" db:"total_cost"`
	BillingPeriod  time.Time  `json:"billing_period" db:"billing_period"`
}

// =============================================
// Audit and Logging Models
// =============================================

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
	BotID   *uuid.UUID     `json:"bot_id" db:"bot_id"`
	BotType *BotType       `json:"bot_type" db:"bot_type"`
}

// =============================================
// Configuration Models
// =============================================

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
	NotifyKey        string     `json:"notify_key" db:"notify_key"`
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

// BotStatusSummary represents a bot status summary view
type BotStatusSummary struct {
	BotID               uuid.UUID  `json:"bot_id" db:"bot_id"`
	BotName             string     `json:"bot_name" db:"bot_name"`
	Status              string     `json:"status" db:"status"`
	BotType             BotType    `json:"bot_type" db:"bot_type"`
	AppID               string     `json:"app_id" db:"app_id"`
	InstallationCount   int        `json:"installation_count" db:"installation_count"`
	ActiveInstallations int        `json:"active_installations" db:"active_installations"`
	HealthStatus        string     `json:"health_status" db:"health_status"`
	LastCheckAt         *time.Time `json:"last_check_at" db:"last_check_at"`
	ResponseTimeMs      *int       `json:"response_time_ms" db:"response_time_ms"`
}

// UsageSummary represents a usage summary view
type UsageSummary struct {
	CompanyID     uuid.UUID `json:"company_id" db:"company_id"`
	CompanyName   string    `json:"company_name" db:"company_name"`
	BillingPeriod time.Time `json:"billing_period" db:"billing_period"`
	RecordType    string    `json:"record_type" db:"record_type"`
	BotType       *BotType  `json:"bot_type" db:"bot_type"`
	TotalQuantity int       `json:"total_quantity" db:"total_quantity"`
	TotalCost     float64   `json:"total_cost" db:"total_cost"`
	ProjectsUsed  int       `json:"projects_used" db:"projects_used"`
	BotsUsed      int       `json:"bots_used" db:"bots_used"`
}
