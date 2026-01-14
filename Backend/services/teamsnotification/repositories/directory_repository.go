package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// DirectoryRepository handles database operations for Azure AD directory
type DirectoryRepository interface {
	UpsertUser(ctx context.Context, user *models.AzureADUser) error
	SearchUsers(ctx context.Context, query string) ([]models.AzureADUser, error)
	UpsertGroup(ctx context.Context, group *models.AzureADGroup) error
	ListGroups(ctx context.Context) ([]models.AzureADGroup, error)
	UpsertChannel(ctx context.Context, channel *models.AzureADChannel) error
	ListChannels(ctx context.Context, teamID string) ([]models.AzureADChannel, error)
	CountChannels(ctx context.Context) (int, error)
	// ChatGroup operations
	CreateChatGroup(ctx context.Context, cg *models.ChatGroup) error
	ListChatGroups(ctx context.Context, projectID uuid.UUID) ([]models.ChatGroup, error)
	DeleteChatGroup(ctx context.Context, id uuid.UUID) error
	// Sync status operations
	CountUsers(ctx context.Context) (int, error)
	CountGroups(ctx context.Context) (int, error)
	GetLastSyncTime(ctx context.Context) (*time.Time, error)
}

type directoryRepository struct {
	db *sqlx.DB
}

// NewDirectoryRepository creates a new directory repository
func NewDirectoryRepository(db *sqlx.DB) DirectoryRepository {
	return &directoryRepository{db: db}
}

func (r *directoryRepository) CreateChatGroup(ctx context.Context, cg *models.ChatGroup) error {
	query := `
		INSERT INTO chat_groups (id, project_id, name, chat_id, created_at, updated_at)
		VALUES (COALESCE(NULLIF(:id, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()), 
		        :project_id, :name, :chat_id, :created_at, :updated_at)
		ON CONFLICT (project_id, chat_id) DO UPDATE SET
			name = EXCLUDED.name,
			updated_at = EXCLUDED.updated_at
	`
	if cg.CreatedAt.IsZero() {
		cg.CreatedAt = time.Now().UTC()
	}
	cg.UpdatedAt = time.Now().UTC()

	_, err := r.db.NamedExecContext(ctx, query, cg)
	return err
}

func (r *directoryRepository) ListChatGroups(ctx context.Context, projectID uuid.UUID) ([]models.ChatGroup, error) {
	query := `SELECT * FROM chat_groups WHERE project_id = $1 ORDER BY name`
	var chatGroups []models.ChatGroup
	err := r.db.SelectContext(ctx, &chatGroups, query, projectID)
	return chatGroups, err
}

func (r *directoryRepository) DeleteChatGroup(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM chat_groups WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *directoryRepository) UpsertUser(ctx context.Context, user *models.AzureADUser) error {
	query := `
		INSERT INTO azure_ad_users (
			azure_ad_id, display_name, email, job_title, department, synced_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (azure_ad_id) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			email = EXCLUDED.email,
			job_title = EXCLUDED.job_title,
			department = EXCLUDED.department,
			synced_at = EXCLUDED.synced_at,
			updated_at = EXCLUDED.updated_at
	`
	now := time.Now().UTC()
	user.SyncedAt = timePtr(now)
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.AzureADID,
		user.DisplayName,
		user.Email,
		user.JobTitle,
		user.Department,
		user.SyncedAt,
		user.CreatedAt,
		user.UpdatedAt,
	)
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
			azure_ad_id, display_name, description, group_types, synced_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		ON CONFLICT (azure_ad_id) DO UPDATE SET
			display_name = EXCLUDED.display_name,
			description = EXCLUDED.description,
			group_types = EXCLUDED.group_types,
			synced_at = EXCLUDED.synced_at,
			updated_at = EXCLUDED.updated_at
	`
	now := time.Now().UTC()
	group.SyncedAt = timePtr(now)
	if group.CreatedAt.IsZero() {
		group.CreatedAt = now
	}
	group.UpdatedAt = now

	_, err := r.db.ExecContext(
		ctx,
		query,
		group.AzureADID,
		group.DisplayName,
		group.Description,
		group.GroupTypes,
		group.SyncedAt,
		group.CreatedAt,
		group.UpdatedAt,
	)
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

func (r *directoryRepository) CountUsers(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM azure_ad_users`
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

func (r *directoryRepository) CountGroups(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM azure_ad_groups`
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count groups: %w", err)
	}
	return count, nil
}

func (r *directoryRepository) UpsertChannel(ctx context.Context, channel *models.AzureADChannel) error {
	query := `
		INSERT INTO azure_ad_channels (
			azure_ad_id, team_id, display_name, description, membership_type, conversation_id, synced_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		ON CONFLICT (azure_ad_id) DO UPDATE SET
			team_id = EXCLUDED.team_id,
			display_name = EXCLUDED.display_name,
			description = EXCLUDED.description,
			membership_type = EXCLUDED.membership_type,
			conversation_id = EXCLUDED.conversation_id,
			synced_at = EXCLUDED.synced_at,
			updated_at = EXCLUDED.updated_at
	`
	now := time.Now().UTC()
	channel.SyncedAt = timePtr(now)
	if channel.CreatedAt.IsZero() {
		channel.CreatedAt = now
	}
	channel.UpdatedAt = now

	_, err := r.db.ExecContext(
		ctx,
		query,
		channel.AzureADID,
		channel.TeamID,
		channel.DisplayName,
		channel.Description,
		channel.MembershipType,
		channel.ConversationID,
		channel.SyncedAt,
		channel.CreatedAt,
		channel.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert channel: %w", err)
	}
	return nil
}

func (r *directoryRepository) ListChannels(ctx context.Context, teamID string) ([]models.AzureADChannel, error) {
	var query string
	var args []interface{}
	if teamID != "" {
		query = `SELECT * FROM azure_ad_channels WHERE team_id = $1 ORDER BY display_name`
		args = []interface{}{teamID}
	} else {
		query = `SELECT * FROM azure_ad_channels ORDER BY display_name`
		args = []interface{}{}
	}
	var channels []models.AzureADChannel
	err := r.db.SelectContext(ctx, &channels, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}
	return channels, nil
}

func (r *directoryRepository) CountChannels(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM azure_ad_channels`
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count channels: %w", err)
	}
	return count, nil
}

func (r *directoryRepository) GetLastSyncTime(ctx context.Context) (*time.Time, error) {
	// Get the most recent synced_at from users, groups, or channels
	var lastSync time.Time
	query := `
		SELECT GREATEST(
			COALESCE((SELECT MAX(synced_at) FROM azure_ad_users WHERE synced_at IS NOT NULL), '1970-01-01'::timestamp),
			COALESCE((SELECT MAX(synced_at) FROM azure_ad_groups WHERE synced_at IS NOT NULL), '1970-01-01'::timestamp),
			COALESCE((SELECT MAX(synced_at) FROM azure_ad_channels WHERE synced_at IS NOT NULL), '1970-01-01'::timestamp)
		) as last_sync
	`
	err := r.db.GetContext(ctx, &lastSync, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get last sync time: %w", err)
	}
	
	// Return nil if no sync has occurred (1970-01-01)
	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	if lastSync.Equal(epoch) || lastSync.Before(epoch) {
		return nil, nil
	}
	
	return &lastSync, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}
