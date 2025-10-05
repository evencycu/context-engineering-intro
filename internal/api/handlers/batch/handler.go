package batch

import (
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles batch notification requests
type Handler struct {
	batchService services.BatchService
}

// NewHandler creates a new batch handler
func NewHandler(batchService services.BatchService) *Handler {
	return &Handler{
		batchService: batchService,
	}
}

// RegisterRoutes registers batch routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	batch := rg.Group("/batch")
	{
		// Batch notifications
		batch.POST("/notifications", h.SendBatchNotifications)
		batch.GET("/notifications/:id", h.GetBatchStatus)
		batch.GET("/notifications/:id/results", h.GetBatchResults)
		batch.DELETE("/notifications/:id", h.CancelBatch)

		// Batch templates
		batch.POST("/templates", h.CreateBatchTemplate)
		batch.GET("/templates", h.ListBatchTemplates)
		batch.GET("/templates/:id", h.GetBatchTemplate)
		batch.PUT("/templates/:id", h.UpdateBatchTemplate)
		batch.DELETE("/templates/:id", h.DeleteBatchTemplate)
	}
}

// SendBatchNotificationsRequest represents a batch notification request
type SendBatchNotificationsRequest struct {
	ProjectID   uuid.UUID      `json:"project_id" binding:"required"`
	SenderID    *uuid.UUID     `json:"sender_id"`
	MessageType string         `json:"message_type" binding:"required,oneof=text file adaptive_card"`
	Content     string         `json:"content" binding:"required"`
	Priority    string         `json:"priority" binding:"required,oneof=low normal high"`
	Mentions    []string       `json:"mentions"`
	Metadata    map[string]any `json:"metadata"`
	Targets     []BatchTarget  `json:"targets" binding:"required,min=1"`
	TemplateID  *uuid.UUID     `json:"template_id"`
	ScheduleAt  *string        `json:"schedule_at"`
	ExpiresAt   *string        `json:"expires_at"`
	MaxRetries  int            `json:"max_retries"`
	RetryDelay  *string        `json:"retry_delay"`
}

// BatchTarget represents a target for batch notification
type BatchTarget struct {
	DestinationID  *uuid.UUID     `json:"destination_id"`
	ConversationID *string        `json:"conversation_id"`
	UserID         *string        `json:"user_id"`
	Email          *string        `json:"email"`
	CustomData     map[string]any `json:"custom_data"`
}

// SendBatchNotificationsResponse represents the response for batch notification
type SendBatchNotificationsResponse struct {
	Success bool                  `json:"success"`
	Data    *services.BatchResult `json:"data,omitempty"`
	Error   string                `json:"error,omitempty"`
}

// SendBatchNotifications handles batch notification sending
func (h *Handler) SendBatchNotifications(c *gin.Context) {
	var req SendBatchNotificationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SendBatchNotificationsResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Convert to service request
	serviceReq := &services.SendBatchNotificationsRequest{
		ProjectID:   req.ProjectID,
		SenderID:    req.SenderID,
		MessageType: req.MessageType,
		Content:     req.Content,
		Priority:    req.Priority,
		Mentions:    req.Mentions,
		Metadata:    req.Metadata,
		TemplateID:  req.TemplateID,
		ScheduleAt:  req.ScheduleAt,
		ExpiresAt:   req.ExpiresAt,
		MaxRetries:  req.MaxRetries,
		RetryDelay:  req.RetryDelay,
	}

	// Convert targets
	for _, target := range req.Targets {
		serviceTarget := services.BatchTarget{
			DestinationID:  target.DestinationID,
			ConversationID: target.ConversationID,
			UserID:         target.UserID,
			Email:          target.Email,
			CustomData:     target.CustomData,
		}
		serviceReq.Targets = append(serviceReq.Targets, serviceTarget)
	}

	// Send batch notification
	result, err := h.batchService.SendBatchNotifications(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SendBatchNotificationsResponse{
			Success: false,
			Error:   "Failed to send batch notifications: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SendBatchNotificationsResponse{
		Success: true,
		Data:    result,
	})
}

// GetBatchStatusResponse represents the response for get batch status
type GetBatchStatusResponse struct {
	Success bool                  `json:"success"`
	Data    *services.BatchStatus `json:"data,omitempty"`
	Error   string                `json:"error,omitempty"`
}

// GetBatchStatus gets the status of a batch notification
func (h *Handler) GetBatchStatus(c *gin.Context) {
	batchIDStr := c.Param("id")
	batchID, err := uuid.Parse(batchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetBatchStatusResponse{
			Success: false,
			Error:   "Invalid batch ID format",
		})
		return
	}

	status, err := h.batchService.GetBatchStatus(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, GetBatchStatusResponse{
			Success: false,
			Error:   "Batch not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetBatchStatusResponse{
		Success: true,
		Data:    status,
	})
}

// GetBatchResultsResponse represents the response for get batch results
type GetBatchResultsResponse struct {
	Success bool                   `json:"success"`
	Data    *services.BatchResults `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// GetBatchResults gets the results of a batch notification
func (h *Handler) GetBatchResults(c *gin.Context) {
	batchIDStr := c.Param("id")
	batchID, err := uuid.Parse(batchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetBatchResultsResponse{
			Success: false,
			Error:   "Invalid batch ID format",
		})
		return
	}

	results, err := h.batchService.GetBatchResults(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, GetBatchResultsResponse{
			Success: false,
			Error:   "Batch results not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetBatchResultsResponse{
		Success: true,
		Data:    results,
	})
}

// CancelBatchResponse represents the response for cancel batch
type CancelBatchResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// CancelBatch cancels a batch notification
func (h *Handler) CancelBatch(c *gin.Context) {
	batchIDStr := c.Param("id")
	batchID, err := uuid.Parse(batchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, CancelBatchResponse{
			Success: false,
			Error:   "Invalid batch ID format",
		})
		return
	}

	err = h.batchService.CancelBatch(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CancelBatchResponse{
			Success: false,
			Error:   "Failed to cancel batch: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CancelBatchResponse{
		Success: true,
		Message: "Batch cancelled successfully",
	})
}

// CreateBatchTemplateRequest represents a create batch template request
type CreateBatchTemplateRequest struct {
	Name        string         `json:"name" binding:"required,min=3,max=100"`
	Description string         `json:"description" binding:"required,min=10,max=500"`
	ProjectID   uuid.UUID      `json:"project_id" binding:"required"`
	MessageType string         `json:"message_type" binding:"required,oneof=text file adaptive_card"`
	Content     string         `json:"content" binding:"required"`
	Priority    string         `json:"priority" binding:"required,oneof=low normal high"`
	Mentions    []string       `json:"mentions"`
	Metadata    map[string]any `json:"metadata"`
	Targets     []BatchTarget  `json:"targets" binding:"required,min=1"`
	IsActive    bool           `json:"is_active"`
}

// CreateBatchTemplateResponse represents the response for create batch template
type CreateBatchTemplateResponse struct {
	Success bool                    `json:"success"`
	Data    *services.BatchTemplate `json:"data,omitempty"`
	Error   string                  `json:"error,omitempty"`
}

// CreateBatchTemplate creates a new batch template
func (h *Handler) CreateBatchTemplate(c *gin.Context) {
	var req CreateBatchTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CreateBatchTemplateResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Convert to service request
	serviceReq := &services.CreateBatchTemplateRequest{
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		MessageType: req.MessageType,
		Content:     req.Content,
		Priority:    req.Priority,
		Mentions:    req.Mentions,
		Metadata:    req.Metadata,
		IsActive:    req.IsActive,
	}

	// Convert targets
	for _, target := range req.Targets {
		serviceTarget := services.BatchTarget{
			DestinationID:  target.DestinationID,
			ConversationID: target.ConversationID,
			UserID:         target.UserID,
			Email:          target.Email,
			CustomData:     target.CustomData,
		}
		serviceReq.Targets = append(serviceReq.Targets, serviceTarget)
	}

	template, err := h.batchService.CreateBatchTemplate(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CreateBatchTemplateResponse{
			Success: false,
			Error:   "Failed to create batch template: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, CreateBatchTemplateResponse{
		Success: true,
		Data:    template,
	})
}

// ListBatchTemplatesResponse represents the response for list batch templates
type ListBatchTemplatesResponse struct {
	Success bool                      `json:"success"`
	Data    []*services.BatchTemplate `json:"data,omitempty"`
	Total   int64                     `json:"total,omitempty"`
	Error   string                    `json:"error,omitempty"`
}

// ListBatchTemplates lists batch templates
func (h *Handler) ListBatchTemplates(c *gin.Context) {
	projectIDStr := c.Query("project_id")
	var projectID *uuid.UUID
	if projectIDStr != "" {
		if id, err := uuid.Parse(projectIDStr); err == nil {
			projectID = &id
		}
	}

	templates, total, err := h.batchService.ListBatchTemplates(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ListBatchTemplatesResponse{
			Success: false,
			Error:   "Failed to list batch templates: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ListBatchTemplatesResponse{
		Success: true,
		Data:    templates,
		Total:   total,
	})
}

// GetBatchTemplateResponse represents the response for get batch template
type GetBatchTemplateResponse struct {
	Success bool                    `json:"success"`
	Data    *services.BatchTemplate `json:"data,omitempty"`
	Error   string                  `json:"error,omitempty"`
}

// GetBatchTemplate gets a batch template
func (h *Handler) GetBatchTemplate(c *gin.Context) {
	templateIDStr := c.Param("id")
	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetBatchTemplateResponse{
			Success: false,
			Error:   "Invalid template ID format",
		})
		return
	}

	template, err := h.batchService.GetBatchTemplate(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, GetBatchTemplateResponse{
			Success: false,
			Error:   "Template not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetBatchTemplateResponse{
		Success: true,
		Data:    template,
	})
}

// UpdateBatchTemplateRequest represents an update batch template request
type UpdateBatchTemplateRequest struct {
	Name        *string        `json:"name" validate:"omitempty,min=3,max=100"`
	Description *string        `json:"description" validate:"omitempty,min=10,max=500"`
	Content     *string        `json:"content"`
	Priority    *string        `json:"priority" validate:"omitempty,oneof=low normal high"`
	Mentions    []string       `json:"mentions"`
	Metadata    map[string]any `json:"metadata"`
	Targets     []BatchTarget  `json:"targets"`
	IsActive    *bool          `json:"is_active"`
}

// UpdateBatchTemplateResponse represents the response for update batch template
type UpdateBatchTemplateResponse struct {
	Success bool                    `json:"success"`
	Data    *services.BatchTemplate `json:"data,omitempty"`
	Error   string                  `json:"error,omitempty"`
}

// UpdateBatchTemplate updates a batch template
func (h *Handler) UpdateBatchTemplate(c *gin.Context) {
	templateIDStr := c.Param("id")
	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, UpdateBatchTemplateResponse{
			Success: false,
			Error:   "Invalid template ID format",
		})
		return
	}

	var req UpdateBatchTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, UpdateBatchTemplateResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	// Convert to service request
	serviceReq := &services.UpdateBatchTemplateRequest{
		Name:        req.Name,
		Description: req.Description,
		Content:     req.Content,
		Priority:    req.Priority,
		Mentions:    req.Mentions,
		Metadata:    req.Metadata,
		IsActive:    req.IsActive,
	}

	// Convert targets
	for _, target := range req.Targets {
		serviceTarget := services.BatchTarget{
			DestinationID:  target.DestinationID,
			ConversationID: target.ConversationID,
			UserID:         target.UserID,
			Email:          target.Email,
			CustomData:     target.CustomData,
		}
		serviceReq.Targets = append(serviceReq.Targets, serviceTarget)
	}

	template, err := h.batchService.UpdateBatchTemplate(c.Request.Context(), templateID, serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateBatchTemplateResponse{
			Success: false,
			Error:   "Failed to update batch template: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UpdateBatchTemplateResponse{
		Success: true,
		Data:    template,
	})
}

// DeleteBatchTemplateResponse represents the response for delete batch template
type DeleteBatchTemplateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// DeleteBatchTemplate deletes a batch template
func (h *Handler) DeleteBatchTemplate(c *gin.Context) {
	templateIDStr := c.Param("id")
	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, DeleteBatchTemplateResponse{
			Success: false,
			Error:   "Invalid template ID format",
		})
		return
	}

	err = h.batchService.DeleteBatchTemplate(c.Request.Context(), templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, DeleteBatchTemplateResponse{
			Success: false,
			Error:   "Failed to delete batch template: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, DeleteBatchTemplateResponse{
		Success: true,
		Message: "Batch template deleted successfully",
	})
}
