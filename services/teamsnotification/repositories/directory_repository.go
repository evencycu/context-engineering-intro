package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/jmoiron/sqlx"
)

// DirectoryRepository handles database operations for Azure AD directory
type DirectoryRepository interface {
	UpsertUser(ctx context.Context, user *models.AzureADUser) error
	SearchUsers(ctx context.Context, query string) ([]models.AzureADUser, error)
	UpsertGroup(ctx context.Context, group *models.AzureADGroup) error
	ListGroups(ctx context.Context) ([]models.AzureADGroup, error)
}

type directoryRepository struct {
	db *sqlx.DB
}

// NewDirectoryRepository creates a new directory repository
func NewDirectoryRepository(db *sqlx.DB) DirectoryRepository {
	return &directoryRepository{db: db}
}

func (r *directoryRepository) UpsertUser(ctx context.Context, user *models.AzureADUser) error {
	query := `
		INSERT INTO azure_ad_users (
			id, azure_ad_id, display_name, email, job_title, department, synced_at, created_at, updated_at
		) VALUES (
			COALESCE(NULLIF(:id, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()), 
			:azure_ad_id, :display_name, :email, :job_title, :department, :synced_at, :created_at, :updated_at
		)
		ON CONFLICT (azure_ad_id) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			email = EXCLUDED.email,
			job_title = EXCLUDED.job_title,
			department = EXCLUDED.department,
			synced_at = EXCLUDED.synced_at,
			updated_at = EXCLUDED.updated_at
	`
	user.SyncedAt = timePtr(time.Now().UTC())
	user.UpdatedAt = time.Now().UTC()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}
	return nil
}

func (r *directoryRepository) SearchUsers(ctx context.Context, q string) ([]models.AzureADUser, error) {
	query := `
		SELECT * FROM azure_ad_users 
		WHERE display_name ILIKE $1 OR email ILIKE $1 
		ORDER BY display_name 
		LIMIT 20
	`
	searchQuery := "%" + q + "%"
	var users []models.AzureADUser
	err := r.db.SelectContext(ctx, &users, query, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	return users, nil
}

func (r *directoryRepository) UpsertGroup(ctx context.Context, group *models.AzureADGroup) error {
	query := `
		INSERT INTO azure_ad_groups (
			id, azure_ad_id, display_name, description, group_types, synced_at, created_at, updated_at
		) VALUES (
			COALESCE(NULLIF(:id, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()),
			:azure_ad_id, :display_name, :description, :group_types, :synced_at, :created_at, :updated_at
		)
		ON CONFLICT (azure_ad_id) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			description = EXCLUDED.description,
			group_types = EXCLUDED.group_types,
			synced_at = EXCLUDED.synced_at,
			updated_at = EXCLUDED.updated_at
	`
	group.SyncedAt = timePtr(time.Now().UTC())
	group.UpdatedAt = time.Now().UTC()
	if group.CreatedAt.IsZero() {
		group.CreatedAt = time.Now().UTC()
	}

	_, err := r.db.NamedExecContext(ctx, query, group)
	if err != nil {
		return fmt.Errorf("failed to upsert group: %w", err)
	}
	return nil
}

func (r *directoryRepository) ListGroups(ctx context.Context) ([]models.AzureADGroup, error) {
	query := `SELECT * FROM azure_ad_groups ORDER BY display_name`
	var groups []models.AzureADGroup
	err := r.db.SelectContext(ctx, &groups, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}
	return groups, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}
