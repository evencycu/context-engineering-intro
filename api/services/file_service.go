package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/api/repositories"
	"github.com/evencycu/TeamsNotifyGoV2/database"
	"github.com/google/uuid"
)

// FileService defines the interface for file operations
type FileService interface {
	// File upload
	UploadFile(ctx context.Context, req *UploadFileRequest) (*FileInfo, error)
	UploadMultipleFiles(ctx context.Context, req *UploadMultipleFilesRequest) ([]*FileInfo, error)

	// File management
	GetFile(ctx context.Context, fileID uuid.UUID) (*FileInfo, error)
	DownloadFile(ctx context.Context, fileID uuid.UUID) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, fileID uuid.UUID) error
	ListFiles(ctx context.Context, req *ListFilesRequest) ([]*FileInfo, int64, error)

	// File validation
	ValidateFile(ctx context.Context, req *ValidateFileRequest) (*ValidateFileResponse, error)

	// File cleanup
	CleanupExpiredFiles(ctx context.Context) error
}

// fileService implements FileService
type fileService struct {
	fileRepo repositories.FileRepository
	storage  FileStorage
}

// NewFileService creates a new file service
func NewFileService(fileRepo repositories.FileRepository, storage FileStorage) FileService {
	return &fileService{
		fileRepo: fileRepo,
		storage:  storage,
	}
}

// Request/Response types for file operations

type UploadFileRequest struct {
	ProjectID   uuid.UUID
	FileName    string
	FileSize    int64
	ContentType string
	Reader      io.Reader
	Description string
	Tags        []string
	IsPublic    bool
	ExpiresAt   *string
}

type UploadMultipleFilesRequest struct {
	ProjectID   uuid.UUID
	Files       []*UploadFileRequest
	Description string
	Tags        []string
	IsPublic    bool
	ExpiresAt   *string
}

type ListFilesRequest struct {
	ProjectID *uuid.UUID
	Tags      []string
	IsPublic  *bool
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

type ValidateFileRequest struct {
	FileName    string
	FileSize    int64
	ContentType string
}

type ValidateFileResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

type FileInfo struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	FileName    string     `json:"file_name"`
	FileSize    int64      `json:"file_size"`
	ContentType string     `json:"content_type"`
	FileURL     string     `json:"file_url"`
	Description string     `json:"description"`
	Tags        []string   `json:"tags"`
	IsPublic    bool       `json:"is_public"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// FileStorage defines the interface for file storage operations
type FileStorage interface {
	Store(ctx context.Context, key string, reader io.Reader) error
	Retrieve(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) (string, error)
}

// Implementation methods

func (s *fileService) UploadFile(ctx context.Context, req *UploadFileRequest) (*FileInfo, error) {
	// Generate unique file key
	fileKey := fmt.Sprintf("files/%s/%s", req.ProjectID.String(), uuid.New().String())

	// Store file in storage
	if err := s.storage.Store(ctx, fileKey, req.Reader); err != nil {
		return nil, fmt.Errorf("failed to store file: %w", err)
	}

	// Get file URL
	fileURL, err := s.storage.GetURL(ctx, fileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get file URL: %w", err)
	}

	// Parse expires at if provided
	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		if t, err := time.Parse(time.RFC3339, *req.ExpiresAt); err == nil {
			expiresAt = &t
		}
	}

	// Create file record
	file := &database.File{
		BaseModel: database.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ProjectID:   req.ProjectID,
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		ContentType: req.ContentType,
		FileKey:     fileKey,
		FileURL:     fileURL,
		Description: req.Description,
		Tags:        database.JSONBStringArray(req.Tags),
		IsPublic:    req.IsPublic,
		ExpiresAt:   expiresAt,
	}

	// Save to database
	if err := s.fileRepo.Create(ctx, file); err != nil {
		// Cleanup stored file on database error
		s.storage.Delete(ctx, fileKey)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	// Convert to FileInfo
	fileInfo := &FileInfo{
		ID:          file.ID,
		ProjectID:   file.ProjectID,
		FileName:    file.FileName,
		FileSize:    file.FileSize,
		ContentType: file.ContentType,
		FileURL:     file.FileURL,
		Description: file.Description,
		Tags:        req.Tags,
		IsPublic:    file.IsPublic,
		ExpiresAt:   file.ExpiresAt,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
	}

	return fileInfo, nil
}

func (s *fileService) UploadMultipleFiles(ctx context.Context, req *UploadMultipleFilesRequest) ([]*FileInfo, error) {
	var fileInfos []*FileInfo
	var errors []error

	for _, fileReq := range req.Files {
		// Set common properties
		fileReq.ProjectID = req.ProjectID
		fileReq.Description = req.Description
		fileReq.Tags = req.Tags
		fileReq.IsPublic = req.IsPublic
		fileReq.ExpiresAt = req.ExpiresAt

		fileInfo, err := s.UploadFile(ctx, fileReq)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to upload %s: %w", fileReq.FileName, err))
			continue
		}

		fileInfos = append(fileInfos, fileInfo)
	}

	// If all files failed, return error
	if len(fileInfos) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("all files failed to upload: %v", errors)
	}

	return fileInfos, nil
}

func (s *fileService) GetFile(ctx context.Context, fileID uuid.UUID) (*FileInfo, error) {
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	// Convert to FileInfo
	fileInfo := &FileInfo{
		ID:          file.ID,
		ProjectID:   file.ProjectID,
		FileName:    file.FileName,
		FileSize:    file.FileSize,
		ContentType: file.ContentType,
		FileURL:     file.FileURL,
		Description: file.Description,
		Tags:        []string(file.Tags),
		IsPublic:    file.IsPublic,
		ExpiresAt:   file.ExpiresAt,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
	}

	return fileInfo, nil
}

func (s *fileService) DownloadFile(ctx context.Context, fileID uuid.UUID) (io.ReadCloser, error) {
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	// Check if file is expired
	if file.ExpiresAt != nil && file.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("file has expired")
	}

	// Retrieve file from storage
	reader, err := s.storage.Retrieve(ctx, file.FileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file: %w", err)
	}

	return reader, nil
}

func (s *fileService) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}

	// Delete from storage
	if err := s.storage.Delete(ctx, file.FileKey); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete from database
	if err := s.fileRepo.Delete(ctx, fileID); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

func (s *fileService) ListFiles(ctx context.Context, req *ListFilesRequest) ([]*FileInfo, int64, error) {
	files, total, err := s.fileRepo.List(ctx, &repositories.ListFilesRequest{
		ProjectID: req.ProjectID,
		Tags:      req.Tags,
		IsPublic:  req.IsPublic,
		Page:      req.Page,
		PageSize:  req.PageSize,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list files: %w", err)
	}

	// Convert to FileInfo
	var fileInfos []*FileInfo
	for _, file := range files {
		fileInfo := &FileInfo{
			ID:          file.ID,
			ProjectID:   file.ProjectID,
			FileName:    file.FileName,
			FileSize:    file.FileSize,
			ContentType: file.ContentType,
			FileURL:     file.FileURL,
			Description: file.Description,
			Tags:        []string(file.Tags),
			IsPublic:    file.IsPublic,
			ExpiresAt:   file.ExpiresAt,
			CreatedAt:   file.CreatedAt,
			UpdatedAt:   file.UpdatedAt,
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, total, nil
}

func (s *fileService) ValidateFile(ctx context.Context, req *ValidateFileRequest) (*ValidateFileResponse, error) {
	// Check file name
	if req.FileName == "" {
		return &ValidateFileResponse{
			Valid:   false,
			Message: "File name cannot be empty",
		}, nil
	}

	// Check file size (10MB limit)
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if req.FileSize <= 0 {
		return &ValidateFileResponse{
			Valid:   false,
			Message: "File size must be greater than 0",
		}, nil
	}

	if req.FileSize > maxFileSize {
		return &ValidateFileResponse{
			Valid:   false,
			Message: "File size exceeds 10MB limit",
		}, nil
	}

	// Check content type
	allowedTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text/plain",
		"image/jpeg",
		"image/png",
		"image/gif",
		"application/zip",
		"application/x-rar-compressed",
	}

	allowed := false
	for _, allowedType := range allowedTypes {
		if req.ContentType == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		return &ValidateFileResponse{
			Valid:   false,
			Message: fmt.Sprintf("Content type not allowed: %s", req.ContentType),
		}, nil
	}

	return &ValidateFileResponse{
		Valid:   true,
		Message: "File is valid",
	}, nil
}

func (s *fileService) CleanupExpiredFiles(ctx context.Context) error {
	// Get expired files
	expiredFiles, err := s.fileRepo.GetExpiredFiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired files: %w", err)
	}

	// Delete expired files
	for _, file := range expiredFiles {
		// Delete from storage
		if err := s.storage.Delete(ctx, file.FileKey); err != nil {
			// Log error but continue
			fmt.Printf("Failed to delete expired file from storage: %v\n", err)
		}

		// Delete from database
		if err := s.fileRepo.Delete(ctx, file.ID); err != nil {
			// Log error but continue
			fmt.Printf("Failed to delete expired file record: %v\n", err)
		}
	}

	return nil
}
