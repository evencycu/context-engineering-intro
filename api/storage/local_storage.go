package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage implements FileStorage interface using local file system
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
	}
}

// Store stores a file to local storage
func (s *LocalStorage) Store(ctx context.Context, key string, reader io.Reader) error {
	// Create full path
	fullPath := filepath.Join(s.basePath, key)

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy content
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// Retrieve retrieves a file from local storage
func (s *LocalStorage) Retrieve(ctx context.Context, key string) (io.ReadCloser, error) {
	// Create full path
	fullPath := filepath.Join(s.basePath, key)

	// Open file
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete deletes a file from local storage
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	// Create full path
	fullPath := filepath.Join(s.basePath, key)

	// Delete file
	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetURL returns the URL for a file in local storage
func (s *LocalStorage) GetURL(ctx context.Context, key string) (string, error) {
	// For local storage, return a file:// URL or a relative path
	// In production, this would typically return a public URL
	return fmt.Sprintf("file://%s", filepath.Join(s.basePath, key)), nil
}

// Cleanup removes old files from storage
func (s *LocalStorage) Cleanup(ctx context.Context, olderThanDays int) error {
	// Implementation for cleaning up old files
	// This would typically be called by a background job
	return nil
}

// GetStats returns storage statistics
func (s *LocalStorage) GetStats(ctx context.Context) (*StorageStats, error) {
	// Implementation for getting storage statistics
	return &StorageStats{}, nil
}

// StorageStats represents storage statistics
type StorageStats struct {
	TotalFiles int64 `json:"total_files"`
	TotalSize  int64 `json:"total_size"`
	UsedSpace  int64 `json:"used_space"`
	FreeSpace  int64 `json:"free_space"`
}

// Helper function to sanitize file paths
func sanitizePath(path string) string {
	// Remove any path traversal attempts
	path = strings.ReplaceAll(path, "..", "")
	path = strings.ReplaceAll(path, "//", "/")
	return strings.TrimPrefix(path, "/")
}
