package repositories

import (
	"context"
	"encoding/json"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AudienceListRepository defines the interface for audience list operations
type AudienceListRepository interface {
	Create(ctx context.Context, entity *models.AudienceList) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.AudienceList, error)
	Update(ctx context.Context, entity *models.AudienceList) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.AudienceList, error)
	List(ctx context.Context, limit, offset int) ([]*models.AudienceList, error)
	Count(ctx context.Context) (int64, error)
}

// audienceListRepository implements AudienceListRepository
type audienceListRepository struct {
	db *sqlx.DB
}

// NewAudienceListRepository creates a new audience list repository
func NewAudienceListRepository(db *sqlx.DB) AudienceListRepository {
	return &audienceListRepository{db: db}
}

// Create creates a new audience list
func (r *audienceListRepository) Create(ctx context.Context, entity *models.AudienceList) error {
	query := `INSERT INTO audience_lists (
		id, project_id, name, type, count, description, metadata, last_updated, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	)`
	
	metadataJSON, _ := json.Marshal(entity.Metadata)
	if entity.Metadata == nil {
		metadataJSON = []byte("{}")
	}
	
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.Name, entity.Type, entity.Count,
		entity.Description, metadataJSON, entity.LastUpdated,
		entity.CreatedAt, entity.UpdatedAt)
	return err
}

// GetByID retrieves an audience list by ID
func (r *audienceListRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AudienceList, error) {
	var list models.AudienceList
	query := `SELECT * FROM audience_lists WHERE id = $1`
	err := r.db.GetContext(ctx, &list, query, id)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Update updates an audience list
func (r *audienceListRepository) Update(ctx context.Context, entity *models.AudienceList) error {
	query := `UPDATE audience_lists SET
		name = $2, type = $3, count = $4, description = $5, metadata = $6,
		last_updated = $7, updated_at = $8
		WHERE id = $1`
	
	metadataJSON, _ := json.Marshal(entity.Metadata)
	if entity.Metadata == nil {
		metadataJSON = []byte("{}")
	}
	
	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.Name, entity.Type, entity.Count, entity.Description,
		metadataJSON, entity.LastUpdated, entity.UpdatedAt)
	return err
}

// Delete deletes an audience list
func (r *audienceListRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM audience_lists WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByProjectID retrieves all audience lists for a project
func (r *audienceListRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.AudienceList, error) {
	var lists []*models.AudienceList
	query := `SELECT * FROM audience_lists WHERE project_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &lists, query, projectID)
	return lists, err
}

// List lists all audience lists with pagination
func (r *audienceListRepository) List(ctx context.Context, limit, offset int) ([]*models.AudienceList, error) {
	var lists []*models.AudienceList
	query := `SELECT * FROM audience_lists ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &lists, query, limit, offset)
	return lists, err
}

// Count counts all audience lists
func (r *audienceListRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM audience_lists`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}
