package repositories

import (
	"context"
	"fmt"

	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BatchRepository defines the interface for batch operations
type BatchRepository interface {
	// Batch notifications
	Create(ctx context.Context, entity *database.BatchNotification) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.BatchNotification, error)
	Update(ctx context.Context, entity *database.BatchNotification) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, projectID *uuid.UUID, status *string, page, pageSize int) ([]*database.BatchNotification, int64, error)

	// Batch targets
	CreateBatchTargets(ctx context.Context, entities []*database.BatchTarget) error
	GetBatchTargets(ctx context.Context, batchID uuid.UUID) ([]*database.BatchTarget, error)
	UpdateBatchTarget(ctx context.Context, entity *database.BatchTarget) error
	GetBatchTargetsByStatus(ctx context.Context, batchID uuid.UUID, status string) ([]*database.BatchTarget, error)

	// Batch templates
	CreateTemplate(ctx context.Context, entity *database.BatchTemplate) error
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*database.BatchTemplate, error)
	UpdateTemplate(ctx context.Context, entity *database.BatchTemplate) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
	ListTemplates(ctx context.Context, projectID *uuid.UUID, page, pageSize int) ([]*database.BatchTemplate, int64, error)
}

// batchRepository implements BatchRepository
type batchRepository struct {
	db *sqlx.DB
}

// NewBatchRepository creates a new batch repository
func NewBatchRepository(db *sqlx.DB) BatchRepository {
	return &batchRepository{db: db}
}

func (r *batchRepository) Create(ctx context.Context, entity *database.BatchNotification) error {
	query := `INSERT INTO batch_notifications (
		id, project_id, sender_id, message_type, content, priority, mentions, metadata,
		status, total_targets, max_retries, scheduled_at, started_at, completed_at, expires_at,
		created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
	)`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.SenderID, entity.MessageType, entity.Content,
		entity.Priority, entity.Mentions, entity.Metadata, entity.Status, entity.TotalTargets,
		entity.MaxRetries, entity.ScheduledAt, entity.StartedAt, entity.CompletedAt, entity.ExpiresAt,
		entity.CreatedAt, entity.UpdatedAt,
	)
	return err
}

func (r *batchRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.BatchNotification, error) {
	query := `SELECT id, project_id, sender_id, message_type, content, priority, mentions, metadata,
		status, total_targets, max_retries, scheduled_at, started_at, completed_at, expires_at,
		created_at, updated_at
		FROM batch_notifications WHERE id = $1`

	var batch database.BatchNotification
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&batch.ID, &batch.ProjectID, &batch.SenderID, &batch.MessageType, &batch.Content,
		&batch.Priority, &batch.Mentions, &batch.Metadata, &batch.Status, &batch.TotalTargets,
		&batch.MaxRetries, &batch.ScheduledAt, &batch.StartedAt, &batch.CompletedAt, &batch.ExpiresAt,
		&batch.CreatedAt, &batch.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *batchRepository) Update(ctx context.Context, entity *database.BatchNotification) error {
	query := `UPDATE batch_notifications SET
		project_id = $2, sender_id = $3, message_type = $4, content = $5, priority = $6,
		mentions = $7, metadata = $8, status = $9, total_targets = $10, max_retries = $11,
		scheduled_at = $12, started_at = $13, completed_at = $14, expires_at = $15, updated_at = $16
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.SenderID, entity.MessageType, entity.Content,
		entity.Priority, entity.Mentions, entity.Metadata, entity.Status, entity.TotalTargets,
		entity.MaxRetries, entity.ScheduledAt, entity.StartedAt, entity.CompletedAt, entity.ExpiresAt,
		entity.UpdatedAt,
	)
	return err
}

func (r *batchRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM batch_notifications WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *batchRepository) List(ctx context.Context, projectID *uuid.UUID, status *string, page, pageSize int) ([]*database.BatchNotification, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if projectID != nil {
		whereClause += fmt.Sprintf(" AND project_id = $%d", argIndex)
		args = append(args, *projectID)
		argIndex++
	}

	if status != nil {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *status)
		argIndex++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM batch_notifications %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY and LIMIT clauses
	orderBy := "ORDER BY created_at DESC"
	limit := pageSize
	offset := (page - 1) * pageSize
	limitClause := fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	// Execute query
	query := fmt.Sprintf(`SELECT id, project_id, sender_id, message_type, content, priority, mentions, metadata,
		status, total_targets, max_retries, scheduled_at, started_at, completed_at, expires_at,
		created_at, updated_at
		FROM batch_notifications %s %s %s`, whereClause, orderBy, limitClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var batches []*database.BatchNotification
	for rows.Next() {
		var batch database.BatchNotification
		err := rows.Scan(
			&batch.ID, &batch.ProjectID, &batch.SenderID, &batch.MessageType, &batch.Content,
			&batch.Priority, &batch.Mentions, &batch.Metadata, &batch.Status, &batch.TotalTargets,
			&batch.MaxRetries, &batch.ScheduledAt, &batch.StartedAt, &batch.CompletedAt, &batch.ExpiresAt,
			&batch.CreatedAt, &batch.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		batches = append(batches, &batch)
	}

	return batches, total, nil
}

func (r *batchRepository) CreateBatchTargets(ctx context.Context, entities []*database.BatchTarget) error {
	if len(entities) == 0 {
		return nil
	}

	query := `INSERT INTO batch_targets (
		id, batch_id, target_index, destination_id, conversation_id, user_id, email,
		custom_data, status, error_message, sent_at, retry_count, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
	)`

	for _, entity := range entities {
		_, err := r.db.ExecContext(ctx, query,
			entity.ID, entity.BatchID, entity.TargetIndex, entity.DestinationID, entity.ConversationID,
			entity.UserID, entity.Email, entity.CustomData, entity.Status, entity.ErrorMessage,
			entity.SentAt, entity.RetryCount, entity.CreatedAt, entity.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *batchRepository) GetBatchTargets(ctx context.Context, batchID uuid.UUID) ([]*database.BatchTarget, error) {
	query := `SELECT id, batch_id, target_index, destination_id, conversation_id, user_id, email,
		custom_data, status, error_message, sent_at, retry_count, created_at, updated_at
		FROM batch_targets WHERE batch_id = $1 ORDER BY target_index`

	rows, err := r.db.QueryContext(ctx, query, batchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []*database.BatchTarget
	for rows.Next() {
		var target database.BatchTarget
		err := rows.Scan(
			&target.ID, &target.BatchID, &target.TargetIndex, &target.DestinationID, &target.ConversationID,
			&target.UserID, &target.Email, &target.CustomData, &target.Status, &target.ErrorMessage,
			&target.SentAt, &target.RetryCount, &target.CreatedAt, &target.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		targets = append(targets, &target)
	}

	return targets, nil
}

func (r *batchRepository) UpdateBatchTarget(ctx context.Context, entity *database.BatchTarget) error {
	query := `UPDATE batch_targets SET
		destination_id = $2, conversation_id = $3, user_id = $4, email = $5, custom_data = $6,
		status = $7, error_message = $8, sent_at = $9, retry_count = $10, updated_at = $11
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.DestinationID, entity.ConversationID, entity.UserID, entity.Email,
		entity.CustomData, entity.Status, entity.ErrorMessage, entity.SentAt, entity.RetryCount,
		entity.UpdatedAt,
	)
	return err
}

func (r *batchRepository) GetBatchTargetsByStatus(ctx context.Context, batchID uuid.UUID, status string) ([]*database.BatchTarget, error) {
	query := `SELECT id, batch_id, target_index, destination_id, conversation_id, user_id, email,
		custom_data, status, error_message, sent_at, retry_count, created_at, updated_at
		FROM batch_targets WHERE batch_id = $1 AND status = $2 ORDER BY target_index`

	rows, err := r.db.QueryContext(ctx, query, batchID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []*database.BatchTarget
	for rows.Next() {
		var target database.BatchTarget
		err := rows.Scan(
			&target.ID, &target.BatchID, &target.TargetIndex, &target.DestinationID, &target.ConversationID,
			&target.UserID, &target.Email, &target.CustomData, &target.Status, &target.ErrorMessage,
			&target.SentAt, &target.RetryCount, &target.CreatedAt, &target.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		targets = append(targets, &target)
	}

	return targets, nil
}

// Template methods (placeholder implementations)
func (r *batchRepository) CreateTemplate(ctx context.Context, entity *database.BatchTemplate) error {
	// Implementation will be added when template functionality is needed
	return fmt.Errorf("not implemented")
}

func (r *batchRepository) GetTemplateByID(ctx context.Context, id uuid.UUID) (*database.BatchTemplate, error) {
	// Implementation will be added when template functionality is needed
	return nil, fmt.Errorf("not implemented")
}

func (r *batchRepository) UpdateTemplate(ctx context.Context, entity *database.BatchTemplate) error {
	// Implementation will be added when template functionality is needed
	return fmt.Errorf("not implemented")
}

func (r *batchRepository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	// Implementation will be added when template functionality is needed
	return fmt.Errorf("not implemented")
}

func (r *batchRepository) ListTemplates(ctx context.Context, projectID *uuid.UUID, page, pageSize int) ([]*database.BatchTemplate, int64, error) {
	// Implementation will be added when template functionality is needed
	return []*database.BatchTemplate{}, 0, nil
}
