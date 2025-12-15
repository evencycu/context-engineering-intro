package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// TemplateRepository handles database operations for templates
type TemplateRepository interface {
	Create(ctx context.Context, template *models.Template) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Template, error)
	Update(ctx context.Context, template *models.Template) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type templateRepository struct {
	db *sqlx.DB
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository(db *sqlx.DB) TemplateRepository {
	return &templateRepository{db: db}
}

// Create creates a new template
func (r *templateRepository) Create(ctx context.Context, template *models.Template) error {
	query := `
		INSERT INTO templates (
			id, project_id, name, description, variables, default_json_structure, created_at, updated_at
		) VALUES (
			:id, :project_id, :name, :description, :variables, :default_json_structure, :created_at, :updated_at
		)
	`
	template.CreatedAt = time.Now().UTC()
	template.UpdatedAt = time.Now().UTC()

	if template.ID == uuid.Nil {
		template.ID = uuid.New()
	}

	_, err := r.db.NamedExecContext(ctx, query, template)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}
	return nil
}

// GetByID gets a template by ID
func (r *templateRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	query := `SELECT * FROM templates WHERE id = $1`
	var template models.Template
	err := r.db.GetContext(ctx, &template, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("template not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return &template, nil
}

// GetByProjectID gets templates by project ID
func (r *templateRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Template, error) {
	query := `SELECT * FROM templates WHERE project_id = $1 ORDER BY created_at DESC`
	var templates []models.Template
	err := r.db.SelectContext(ctx, &templates, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	return templates, nil
}

// Update updates a template
func (r *templateRepository) Update(ctx context.Context, template *models.Template) error {
	query := `
		UPDATE templates SET
			name = :name,
			description = :description,
			variables = :variables,
			default_json_structure = :default_json_structure,
			updated_at = :updated_at
		WHERE id = :id
	`
	template.UpdatedAt = time.Now().UTC()

	result, err := r.db.NamedExecContext(ctx, query, template)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return errors.New("template not found or no changes made")
	}

	return nil
}

// Delete deletes a template
func (r *templateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM templates WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return errors.New("template not found")
	}

	return nil
}
