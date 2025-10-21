package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/libs/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// FileRepository defines the interface for file operations
type FileRepository interface {
	Create(ctx context.Context, entity *database.File) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.File, error)
	Update(ctx context.Context, entity *database.File) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, req *ListFilesRequest) ([]*database.File, int64, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*database.File, error)
	GetExpiredFiles(ctx context.Context) ([]*database.File, error)
	GetByTags(ctx context.Context, tags []string) ([]*database.File, error)
	GetPublicFiles(ctx context.Context) ([]*database.File, error)
}

// ListFilesRequest represents a request to list files
type ListFilesRequest struct {
	ProjectID *uuid.UUID
	Tags      []string
	IsPublic  *bool
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

// fileRepository implements FileRepository
type fileRepository struct {
	db *sqlx.DB
}

// NewFileRepository creates a new file repository
func NewFileRepository(db *sqlx.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, entity *database.File) error {
	query := `INSERT INTO files (
		id, project_id, file_name, file_size, content_type, file_key, file_url,
		description, tags, is_public, expires_at, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
	)`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.FileName, entity.FileSize, entity.ContentType,
		entity.FileKey, entity.FileURL, entity.Description, entity.Tags, entity.IsPublic,
		entity.ExpiresAt, entity.CreatedAt, entity.UpdatedAt,
	)
	return err
}

func (r *fileRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.File, error) {
	query := `SELECT id, project_id, file_name, file_size, content_type, file_key, file_url,
		description, tags, is_public, expires_at, created_at, updated_at
		FROM files WHERE id = $1`

	var file database.File
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&file.ID, &file.ProjectID, &file.FileName, &file.FileSize, &file.ContentType,
		&file.FileKey, &file.FileURL, &file.Description, &file.Tags, &file.IsPublic,
		&file.ExpiresAt, &file.CreatedAt, &file.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) Update(ctx context.Context, entity *database.File) error {
	query := `UPDATE files SET
		file_name = $2, file_size = $3, content_type = $4, file_key = $5, file_url = $6,
		description = $7, tags = $8, is_public = $9, expires_at = $10, updated_at = $11
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.FileName, entity.FileSize, entity.ContentType, entity.FileKey,
		entity.FileURL, entity.Description, entity.Tags, entity.IsPublic, entity.ExpiresAt,
		entity.UpdatedAt,
	)
	return err
}

func (r *fileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *fileRepository) List(ctx context.Context, req *ListFilesRequest) ([]*database.File, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if req.ProjectID != nil {
		whereClause += fmt.Sprintf(" AND project_id = $%d", argIndex)
		args = append(args, *req.ProjectID)
		argIndex++
	}

	if req.IsPublic != nil {
		whereClause += fmt.Sprintf(" AND is_public = $%d", argIndex)
		args = append(args, *req.IsPublic)
		argIndex++
	}

	if len(req.Tags) > 0 {
		whereClause += fmt.Sprintf(" AND tags && $%d", argIndex)
		args = append(args, req.Tags)
		argIndex++
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM files %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause with whitelist validation
	orderBy := "ORDER BY created_at DESC"
	if req.SortBy != "" {
		// Whitelist allowed sort columns to prevent SQL injection
		allowedSortColumns := map[string]bool{
			"created_at":   true,
			"updated_at":   true,
			"file_name":    true,
			"file_size":    true,
			"content_type": true,
		}

		if allowedSortColumns[req.SortBy] {
			orderBy = fmt.Sprintf("ORDER BY %s", req.SortBy)
			if req.SortOrder == "asc" {
				orderBy += " ASC"
			} else {
				orderBy += " DESC"
			}
		}
	}

	// Build LIMIT and OFFSET
	limit := req.PageSize
	offset := (req.Page - 1) * req.PageSize
	limitClause := fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	// Execute query
	query := fmt.Sprintf(`SELECT id, project_id, file_name, file_size, content_type, file_key, file_url,
		description, tags, is_public, expires_at, created_at, updated_at
		FROM files %s %s %s`, whereClause, orderBy, limitClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var files []*database.File
	for rows.Next() {
		var file database.File
		err := rows.Scan(
			&file.ID, &file.ProjectID, &file.FileName, &file.FileSize, &file.ContentType,
			&file.FileKey, &file.FileURL, &file.Description, &file.Tags, &file.IsPublic,
			&file.ExpiresAt, &file.CreatedAt, &file.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		files = append(files, &file)
	}

	return files, total, nil
}

func (r *fileRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*database.File, error) {
	query := `SELECT id, project_id, file_name, file_size, content_type, file_key, file_url,
		description, tags, is_public, expires_at, created_at, updated_at
		FROM files WHERE project_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*database.File
	for rows.Next() {
		var file database.File
		err := rows.Scan(
			&file.ID, &file.ProjectID, &file.FileName, &file.FileSize, &file.ContentType,
			&file.FileKey, &file.FileURL, &file.Description, &file.Tags, &file.IsPublic,
			&file.ExpiresAt, &file.CreatedAt, &file.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}

func (r *fileRepository) GetExpiredFiles(ctx context.Context) ([]*database.File, error) {
	query := `SELECT id, project_id, file_name, file_size, content_type, file_key, file_url,
		description, tags, is_public, expires_at, created_at, updated_at
		FROM files WHERE expires_at IS NOT NULL AND expires_at < $1`

	rows, err := r.db.QueryContext(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*database.File
	for rows.Next() {
		var file database.File
		err := rows.Scan(
			&file.ID, &file.ProjectID, &file.FileName, &file.FileSize, &file.ContentType,
			&file.FileKey, &file.FileURL, &file.Description, &file.Tags, &file.IsPublic,
			&file.ExpiresAt, &file.CreatedAt, &file.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}

func (r *fileRepository) GetByTags(ctx context.Context, tags []string) ([]*database.File, error) {
	query := `SELECT id, project_id, file_name, file_size, content_type, file_key, file_url,
		description, tags, is_public, expires_at, created_at, updated_at
		FROM files WHERE tags && $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, tags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*database.File
	for rows.Next() {
		var file database.File
		err := rows.Scan(
			&file.ID, &file.ProjectID, &file.FileName, &file.FileSize, &file.ContentType,
			&file.FileKey, &file.FileURL, &file.Description, &file.Tags, &file.IsPublic,
			&file.ExpiresAt, &file.CreatedAt, &file.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}

func (r *fileRepository) GetPublicFiles(ctx context.Context) ([]*database.File, error) {
	query := `SELECT id, project_id, file_name, file_size, content_type, file_key, file_url,
		description, tags, is_public, expires_at, created_at, updated_at
		FROM files WHERE is_public = true ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*database.File
	for rows.Next() {
		var file database.File
		err := rows.Scan(
			&file.ID, &file.ProjectID, &file.FileName, &file.FileSize, &file.ContentType,
			&file.FileKey, &file.FileURL, &file.Description, &file.Tags, &file.IsPublic,
			&file.ExpiresAt, &file.CreatedAt, &file.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}
