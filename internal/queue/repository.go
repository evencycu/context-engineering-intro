package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for failed notifications queue
type Repository interface {
	// Enqueue adds a failed notification to the retry queue
	Enqueue(ctx context.Context, fn *FailedNotification) error

	// Dequeue retrieves notifications ready for retry
	Dequeue(ctx context.Context, limit int) ([]*FailedNotification, error)

	// UpdateRetry updates retry information after an attempt
	UpdateRetry(ctx context.Context, id uuid.UUID, success bool, errorMsg string, nextRetryAt time.Time) error

	// MarkExhausted marks a notification as exhausted (max retries reached)
	MarkExhausted(ctx context.Context, id uuid.UUID) error

	// Delete removes a notification from the queue
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByNotificationID retrieves failed notifications by notification ID
	GetByNotificationID(ctx context.Context, notificationID uuid.UUID) ([]*FailedNotification, error)

	// GetStatus returns the current queue status
	GetStatus(ctx context.Context) (*NotificationQueueStatus, error)

	// CleanupOld removes old exhausted notifications
	CleanupOld(ctx context.Context, olderThan time.Duration) (int64, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new queue repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Enqueue adds a failed notification to the retry queue
func (r *repository) Enqueue(ctx context.Context, fn *FailedNotification) error {
	if fn.ID == uuid.Nil {
		fn.ID = uuid.New()
	}
	if fn.CreatedAt.IsZero() {
		fn.CreatedAt = time.Now()
	}
	fn.UpdatedAt = time.Now()

	metadataJSON, err := json.Marshal(fn.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	query := `
		INSERT INTO failed_notifications (
			id, notification_id, project_id, target_id, message,
			reason, error_message, retry_count, max_retries, next_retry_at,
			retry_after, created_at, updated_at, installation_id, bot_id, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err = r.db.ExecContext(ctx, query,
		fn.ID, fn.NotificationID, fn.ProjectID, fn.TargetID, fn.Message,
		fn.Reason, fn.ErrorMessage, fn.RetryCount, fn.MaxRetries, fn.NextRetryAt,
		fn.RetryAfter, fn.CreatedAt, fn.UpdatedAt, fn.InstallationID, fn.BotID, metadataJSON,
	)

	return err
}

// Dequeue retrieves notifications ready for retry
func (r *repository) Dequeue(ctx context.Context, limit int) ([]*FailedNotification, error) {
	query := `
		SELECT 
			id, notification_id, project_id, target_id, message,
			reason, error_message, retry_count, max_retries, next_retry_at,
			retry_after, created_at, updated_at, last_attempt_at,
			installation_id, bot_id, metadata
		FROM failed_notifications
		WHERE retry_count < max_retries 
			AND next_retry_at <= $1
		ORDER BY next_retry_at ASC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, time.Now(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*FailedNotification
	for rows.Next() {
		fn := &FailedNotification{}
		var metadataJSON []byte

		err := rows.Scan(
			&fn.ID, &fn.NotificationID, &fn.ProjectID, &fn.TargetID, &fn.Message,
			&fn.Reason, &fn.ErrorMessage, &fn.RetryCount, &fn.MaxRetries, &fn.NextRetryAt,
			&fn.RetryAfter, &fn.CreatedAt, &fn.UpdatedAt, &fn.LastAttemptAt,
			&fn.InstallationID, &fn.BotID, &metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &fn.Metadata)
		}

		notifications = append(notifications, fn)
	}

	return notifications, rows.Err()
}

// UpdateRetry updates retry information after an attempt
func (r *repository) UpdateRetry(ctx context.Context, id uuid.UUID, success bool, errorMsg string, nextRetryAt time.Time) error {
	if success {
		// If successful, remove from queue
		return r.Delete(ctx, id)
	}

	query := `
		UPDATE failed_notifications
		SET retry_count = retry_count + 1,
			error_message = $2,
			next_retry_at = $3,
			last_attempt_at = $4,
			updated_at = $5
		WHERE id = $1
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, id, errorMsg, nextRetryAt, now, now)
	return err
}

// MarkExhausted marks a notification as exhausted (max retries reached)
func (r *repository) MarkExhausted(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE failed_notifications
		SET retry_count = max_retries,
			updated_at = $2
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, time.Now())
	return err
}

// Delete removes a notification from the queue
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM failed_notifications WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByNotificationID retrieves failed notifications by notification ID
func (r *repository) GetByNotificationID(ctx context.Context, notificationID uuid.UUID) ([]*FailedNotification, error) {
	query := `
		SELECT 
			id, notification_id, project_id, target_id, message,
			reason, error_message, retry_count, max_retries, next_retry_at,
			retry_after, created_at, updated_at, last_attempt_at,
			installation_id, bot_id, metadata
		FROM failed_notifications
		WHERE notification_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, notificationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*FailedNotification
	for rows.Next() {
		fn := &FailedNotification{}
		var metadataJSON []byte

		err := rows.Scan(
			&fn.ID, &fn.NotificationID, &fn.ProjectID, &fn.TargetID, &fn.Message,
			&fn.Reason, &fn.ErrorMessage, &fn.RetryCount, &fn.MaxRetries, &fn.NextRetryAt,
			&fn.RetryAfter, &fn.CreatedAt, &fn.UpdatedAt, &fn.LastAttemptAt,
			&fn.InstallationID, &fn.BotID, &metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &fn.Metadata)
		}

		notifications = append(notifications, fn)
	}

	return notifications, rows.Err()
}

// GetStatus returns the current queue status
func (r *repository) GetStatus(ctx context.Context) (*NotificationQueueStatus, error) {
	query := `
		SELECT 
			pending_count,
			ready_for_retry_count,
			exhausted_count,
			next_retry_due,
			oldest_pending,
			last_activity
		FROM failed_notifications_queue_status
	`

	status := &NotificationQueueStatus{
		LastUpdated: time.Now(),
	}

	var nextRetryDue, oldestPending, lastActivity sql.NullTime

	err := r.db.QueryRowContext(ctx, query).Scan(
		&status.TotalPending,
		&status.TotalRetrying,
		&status.TotalFailed,
		&nextRetryDue,
		&oldestPending,
		&lastActivity,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if nextRetryDue.Valid {
		status.NextRetryDue = &nextRetryDue.Time
	}
	if oldestPending.Valid {
		status.OldestPending = &oldestPending.Time
	}

	return status, nil
}

// CleanupOld removes old exhausted notifications
func (r *repository) CleanupOld(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `
		DELETE FROM failed_notifications
		WHERE retry_count >= max_retries
			AND updated_at < $1
	`

	cutoff := time.Now().Add(-olderThan)
	result, err := r.db.ExecContext(ctx, query, cutoff)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// BatchEnqueue adds multiple failed notifications to the queue
func (r *repository) BatchEnqueue(ctx context.Context, notifications []*FailedNotification) error {
	if len(notifications) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO failed_notifications (
			id, notification_id, project_id, target_id, message,
			reason, error_message, retry_count, max_retries, next_retry_at,
			retry_after, created_at, updated_at, installation_id, bot_id, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, fn := range notifications {
		if fn.ID == uuid.Nil {
			fn.ID = uuid.New()
		}
		if fn.CreatedAt.IsZero() {
			fn.CreatedAt = time.Now()
		}
		fn.UpdatedAt = time.Now()

		metadataJSON, err := json.Marshal(fn.Metadata)
		if err != nil {
			metadataJSON = []byte("{}")
		}

		_, err = stmt.ExecContext(ctx,
			fn.ID, fn.NotificationID, fn.ProjectID, fn.TargetID, fn.Message,
			fn.Reason, fn.ErrorMessage, fn.RetryCount, fn.MaxRetries, fn.NextRetryAt,
			fn.RetryAfter, fn.CreatedAt, fn.UpdatedAt, fn.InstallationID, fn.BotID, metadataJSON,
		)
		if err != nil {
			return fmt.Errorf("failed to insert notification %s: %w", fn.ID, err)
		}
	}

	return tx.Commit()
}

