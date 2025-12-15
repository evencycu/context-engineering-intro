package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles file upload and management
type Handler struct {
	fileService services.FileService
}

// NewHandler creates a new file handler
func NewHandler(fileService services.FileService) *Handler {
	return &Handler{
		fileService: fileService,
	}
}

// RegisterRoutes registers file routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	files := rg.Group("/files")
	{
		// File upload
		files.POST("/upload", h.UploadFile)
		files.POST("/upload/multiple", h.UploadMultipleFiles)

		// File management
		files.GET("/:id", h.GetFile)
		files.GET("/:id/download", h.DownloadFile)
		files.DELETE("/:id", h.DeleteFile)
		files.GET("", h.ListFiles)

		// File validation
		files.POST("/validate", h.ValidateFile)
	}
}

// UploadFileRequest represents a file upload request
type UploadFileRequest struct {
	ProjectID   uuid.UUID `form:"project_id" binding:"required"`
	Description string    `form:"description"`
	Tags        []string  `form:"tags"`
	IsPublic    bool      `form:"is_public"`
	ExpiresAt   *string   `form:"expires_at"`
}

// UploadFileResponse represents the response for file upload
type UploadFileResponse struct {
	Success bool               `json:"success"`
	Data    *services.FileInfo `json:"data,omitempty"`
	Error   string             `json:"error,omitempty"`
}

// UploadFile handles single file upload
func (h *Handler) UploadFile(c *gin.Context) {
	var req UploadFileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, UploadFileResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, UploadFileResponse{
			Success: false,
			Error:   "No file uploaded: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// Validate file
	if err := h.validateFile(header); err != nil {
		c.JSON(http.StatusBadRequest, UploadFileResponse{
			Success: false,
			Error:   "File validation failed: " + err.Error(),
		})
		return
	}

	// Upload file
	fileInfo, err := h.fileService.UploadFile(c.Request.Context(), &services.UploadFileRequest{
		ProjectID:   req.ProjectID,
		FileName:    header.Filename,
		FileSize:    header.Size,
		ContentType: header.Header.Get("Content-Type"),
		Reader:      file,
		Description: req.Description,
		Tags:        req.Tags,
		IsPublic:    req.IsPublic,
		ExpiresAt:   req.ExpiresAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, UploadFileResponse{
			Success: false,
			Error:   "Failed to upload file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, UploadFileResponse{
		Success: true,
		Data:    fileInfo,
	})
}

// UploadMultipleFilesRequest represents a multiple file upload request
type UploadMultipleFilesRequest struct {
	ProjectID   uuid.UUID `form:"project_id" binding:"required"`
	Description string    `form:"description"`
	Tags        []string  `form:"tags"`
	IsPublic    bool      `form:"is_public"`
	ExpiresAt   *string   `form:"expires_at"`
}

// UploadMultipleFilesResponse represents the response for multiple file upload
type UploadMultipleFilesResponse struct {
	Success bool                 `json:"success"`
	Data    []*services.FileInfo `json:"data,omitempty"`
	Errors  []string             `json:"errors,omitempty"`
}

// UploadMultipleFiles handles multiple file upload
func (h *Handler) UploadMultipleFiles(c *gin.Context) {
	var req UploadMultipleFilesRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, UploadMultipleFilesResponse{
			Success: false,
			Errors:  []string{"Invalid request: " + err.Error()},
		})
		return
	}

	// Get form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, UploadMultipleFilesResponse{
			Success: false,
			Errors:  []string{"Failed to parse multipart form: " + err.Error()},
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, UploadMultipleFilesResponse{
			Success: false,
			Errors:  []string{"No files uploaded"},
		})
		return
	}

	var uploadedFiles []*services.FileInfo
	var errors []string

	for _, fileHeader := range files {
		// Validate file
		if err := h.validateFile(fileHeader); err != nil {
			errors = append(errors, fmt.Sprintf("File %s: %s", fileHeader.Filename, err.Error()))
			continue
		}

		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("File %s: Failed to open: %s", fileHeader.Filename, err.Error()))
			continue
		}

		// Upload file
		fileInfo, err := h.fileService.UploadFile(c.Request.Context(), &services.UploadFileRequest{
			ProjectID:   req.ProjectID,
			FileName:    fileHeader.Filename,
			FileSize:    fileHeader.Size,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Reader:      file,
			Description: req.Description,
			Tags:        req.Tags,
			IsPublic:    req.IsPublic,
			ExpiresAt:   req.ExpiresAt,
		})
		file.Close()

		if err != nil {
			errors = append(errors, fmt.Sprintf("File %s: Upload failed: %s", fileHeader.Filename, err.Error()))
			continue
		}

		uploadedFiles = append(uploadedFiles, fileInfo)
	}

	success := len(uploadedFiles) > 0
	c.JSON(http.StatusCreated, UploadMultipleFilesResponse{
		Success: success,
		Data:    uploadedFiles,
		Errors:  errors,
	})
}

// GetFileResponse represents the response for get file
type GetFileResponse struct {
	Success bool               `json:"success"`
	Data    *services.FileInfo `json:"data,omitempty"`
	Error   string             `json:"error,omitempty"`
}

// GetFile gets file information
func (h *Handler) GetFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetFileResponse{
			Success: false,
			Error:   "Invalid file ID format",
		})
		return
	}

	fileInfo, err := h.fileService.GetFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, GetFileResponse{
			Success: false,
			Error:   "File not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetFileResponse{
		Success: true,
		Data:    fileInfo,
	})
}

// DownloadFile downloads a file
func (h *Handler) DownloadFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file ID format",
		})
		return
	}

	// Get file info first
	fileInfo, err := h.fileService.GetFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found: " + err.Error(),
		})
		return
	}

	// Get file content
	content, err := h.fileService.DownloadFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to download file: " + err.Error(),
		})
		return
	}
	defer content.Close()

	// Set headers
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.FileName))
	c.Header("Content-Type", fileInfo.ContentType)
	c.Header("Content-Length", strconv.FormatInt(fileInfo.FileSize, 10))

	// Stream file content
	_, err = io.Copy(c.Writer, content)
	if err != nil {
		// Log error but don't send response as headers are already sent
		fmt.Printf("Error streaming file: %v\n", err)
	}
}

// DeleteFileResponse represents the response for delete file
type DeleteFileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// DeleteFile deletes a file
func (h *Handler) DeleteFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, DeleteFileResponse{
			Success: false,
			Error:   "Invalid file ID format",
		})
		return
	}

	err = h.fileService.DeleteFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, DeleteFileResponse{
			Success: false,
			Error:   "Failed to delete file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, DeleteFileResponse{
		Success: true,
		Message: "File deleted successfully",
	})
}

// ListFilesRequest represents a list files request
type ListFilesRequest struct {
	ProjectID *uuid.UUID `form:"project_id"`
	Tags      []string   `form:"tags"`
	IsPublic  *bool      `form:"is_public"`
	Page      int        `form:"page,default=1"`
	PageSize  int        `form:"page_size,default=50"`
	SortBy    string     `form:"sort_by,default=created_at"`
	SortOrder string     `form:"sort_order,default=desc"`
}

// ListFilesResponse represents the response for list files
type ListFilesResponse struct {
	Success  bool                 `json:"success"`
	Data     []*services.FileInfo `json:"data,omitempty"`
	Total    int64                `json:"total,omitempty"`
	Page     int                  `json:"page,omitempty"`
	PageSize int                  `json:"page_size,omitempty"`
	Error    string               `json:"error,omitempty"`
}

// ListFiles lists files with filtering and pagination
func (h *Handler) ListFiles(c *gin.Context) {
	var req ListFilesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ListFilesResponse{
			Success: false,
			Error:   "Invalid query parameters: " + err.Error(),
		})
		return
	}

	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 50
	}
	if req.PageSize > 1000 {
		req.PageSize = 1000
	}

	files, total, err := h.fileService.ListFiles(c.Request.Context(), &services.ListFilesRequest{
		ProjectID: req.ProjectID,
		Tags:      req.Tags,
		IsPublic:  req.IsPublic,
		Page:      req.Page,
		PageSize:  req.PageSize,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ListFilesResponse{
			Success: false,
			Error:   "Failed to list files: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListFilesResponse{
		Success:  true,
		Data:     files,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// ValidateFileRequest represents a file validation request
type ValidateFileRequest struct {
	FileName    string `json:"fileName" binding:"required"`
	FileSize    int64  `json:"fileSize" binding:"required"`
	ContentType string `json:"contentType" binding:"required"`
}

// ValidateFileResponse represents the response for file validation
type ValidateFileResponse struct {
	Success bool   `json:"success"`
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ValidateFile validates a file before upload
func (h *Handler) ValidateFile(c *gin.Context) {
	var req ValidateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ValidateFileResponse{
			Success: false,
			Valid:   false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Create a mock file header for validation
	header := &multipart.FileHeader{
		Filename: req.FileName,
		Size:     req.FileSize,
		Header:   make(map[string][]string),
	}
	header.Header.Set("Content-Type", req.ContentType)

	// Validate file
	err := h.validateFile(header)
	if err != nil {
		c.JSON(http.StatusOK, ValidateFileResponse{
			Success: true,
			Valid:   false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ValidateFileResponse{
		Success: true,
		Valid:   true,
		Message: "File is valid",
	})
}

// validateFile validates a file header
func (h *Handler) validateFile(header *multipart.FileHeader) error {
	// Check file size (10MB limit)
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if header.Size > maxFileSize {
		return fmt.Errorf("file size exceeds 10MB limit")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".txt", ".jpg", ".jpeg", ".png", ".gif", ".zip", ".rar",
	}

	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file type not allowed: %s", ext)
	}

	// Check content type
	contentType := header.Header.Get("Content-Type")
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

	allowed = false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("content type not allowed: %s", contentType)
	}

	return nil
}
