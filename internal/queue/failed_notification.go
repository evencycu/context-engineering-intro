package queue

import (
	"time"

	"github.com/google/uuid"
)

// FailureReason represents why a notification failed
type FailureReason string

const (
	FailureReasonRateLimit   FailureReason = "rate_limit"   // 429 Too Many Requests
	FailureReasonTimeout     FailureReason = "timeout"      // Request timeout
	FailureReasonUnauthorized FailureReason = "unauthorized" // 401/403
	FailureReasonServerError FailureReason = "server_error" // 5xx errors
	FailureReasonNetworkError FailureReason = "network_error" // Network issues
	FailureReasonUnknown     FailureReason = "unknown"      // Other errors
)

// FailedNotification represents a notification that failed to send
type FailedNotification struct {
	ID             uuid.UUID     `json:"id" db:"id"`
	NotificationID uuid.UUID     `json:"notification_id" db:"notification_id"`
	ProjectID      uuid.UUID     `json:"project_id" db:"project_id"`
	TargetID       string        `json:"target_id" db:"target_id"` // conversation_id or email
	Message        string        `json:"message" db:"message"`
	Reason         FailureReason `json:"reason" db:"reason"`
	ErrorMessage   string        `json:"error_message" db:"error_message"`
	RetryCount     int           `json:"retry_count" db:"retry_count"`
	MaxRetries     int           `json:"max_retries" db:"max_retries"`
	NextRetryAt    time.Time     `json:"next_retry_at" db:"next_retry_at"`
	RetryAfter     *int          `json:"retry_after,omitempty" db:"retry_after"` // From Retry-After header (seconds)
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
	LastAttemptAt  *time.Time    `json:"last_attempt_at,omitempty" db:"last_attempt_at"`
	
	// Additional context
	InstallationID *uuid.UUID             `json:"installation_id,omitempty" db:"installation_id"`
	BotID          *uuid.UUID             `json:"bot_id,omitempty" db:"bot_id"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
}

// IsRetryable returns true if the notification can be retried
func (fn *FailedNotification) IsRetryable() bool {
	return fn.RetryCount < fn.MaxRetries
}

// ShouldRetryNow returns true if the notification should be retried now
func (fn *FailedNotification) ShouldRetryNow() bool {
	return fn.IsRetryable() && time.Now().After(fn.NextRetryAt)
}

// IncrementRetry increments the retry count and updates next retry time
func (fn *FailedNotification) IncrementRetry(nextRetryAt time.Time) {
	fn.RetryCount++
	fn.NextRetryAt = nextRetryAt
	fn.UpdatedAt = time.Now()
	now := time.Now()
	fn.LastAttemptAt = &now
}

// NotificationQueueStatus represents the status of a notification in the queue
type NotificationQueueStatus struct {
	TotalPending     int       `json:"total_pending"`
	TotalRetrying    int       `json:"total_retrying"`
	TotalFailed      int       `json:"total_failed"`
	OldestPending    *time.Time `json:"oldest_pending,omitempty"`
	NextRetryDue     *time.Time `json:"next_retry_due,omitempty"`
	CircuitState     string    `json:"circuit_state"`
	LastUpdated      time.Time `json:"last_updated"`
}

