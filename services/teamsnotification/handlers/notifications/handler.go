package notifications

import (
	"net/http"
	"strconv"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/services/teamsnotification/services"
	"github.com/evencycu/TeamsNotifyGoV2/libs/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles notification-related HTTP requests
type Handler struct {
	notificationService services.NotificationService
}

// NewHandler creates a new notification handler
func NewHandler(notificationService services.NotificationService) *Handler {
	return &Handler{
		notificationService: notificationService,
	}
}

// RegisterRoutes registers notification routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	notifications := rg.Group("/notifications")
	{
		notifications.POST("", h.SendNotification)
		notifications.GET("", h.ListNotifications)
		notifications.GET("/:id", h.GetNotification)
		notifications.POST("/:id/retry", h.RetryNotification)
		notifications.DELETE("/:id", h.CancelNotification)
		notifications.GET("/project/:projectId", h.GetNotificationsByProject)
		notifications.GET("/sender/:senderId", h.GetNotificationsBySender)
		notifications.GET("/status/:status", h.GetNotificationsByStatus)
		notifications.GET("/date-range", h.GetNotificationsByDateRange)
	}
}

// SendNotificationRequest represents a send notification request
type SendNotificationRequest struct {
	ProjectID    uuid.UUID              `json:"projectId" validate:"required"`
	SenderID     *uuid.UUID             `json:"senderId"`
	MessageType  string                 `json:"messageType" validate:"required,oneof=text file adaptive_card"`
	Content      string                 `json:"content" validate:"required,max=4000"`
	Mentions     []string               `json:"mentions"`
	Attachment   *database.Attachment   `json:"attachment"`
	AdaptiveCard *database.AdaptiveCard `json:"adaptiveCard"`
	Priority     string                 `json:"priority" validate:"oneof=low normal high"`
	Metadata     map[string]interface{} `json:"metadata"`
	Targets      []string               `json:"targets" validate:"omitempty,min=1"`
}

// SendNotificationResponse represents a send notification response
type SendNotificationResponse struct {
	NotificationID uuid.UUID `json:"notificationId"`
	Status         string    `json:"status"`
	Message        string    `json:"message"`
	Destinations   int       `json:"destinationsCount"`
	EstimatedTime  string    `json:"estimated_delivery_time,omitempty"`
}

// NotificationStatusResponse represents a notification status response
type NotificationStatusResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Status             string     `json:"status"`
	Message            string     `json:"message"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
	SentAt             *time.Time `json:"sent_at,omitempty"`
	ErrorMessage       string     `json:"error_message,omitempty"`
	DestinationsSent   int        `json:"destinationsSent"`
	DestinationsFailed int        `json:"destinationsFailed"`
	TotalDestinations  int        `json:"totalDestinations"`
}

// DateRangeRequest represents a date range request
type DateRangeRequest struct {
	StartDate string `form:"start_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate   string `form:"end_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

// SendNotification sends a new notification
func (h *Handler) SendNotification(c *gin.Context) {
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert request to service request
	serviceReq := &services.SendNotificationRequest{
		ProjectID:    req.ProjectID,
		SenderID:     req.SenderID,
		MessageType:  req.MessageType,
		Content:      req.Content,
		Mentions:     req.Mentions,
		Attachment:   req.Attachment,
		AdaptiveCard: req.AdaptiveCard,
		Priority:     req.Priority,
		Metadata:     req.Metadata,
		Targets:      req.Targets,
	}

	// Send notification
	notification, err := h.notificationService.SendNotification(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send notification",
			"details": err.Error(),
		})
		return
	}

	// Create response
	response := SendNotificationResponse{
		NotificationID: notification.ID,
		Status:         string(notification.Status),
		Message:        "Notification queued successfully",
		Destinations:   len(req.Targets),
	}

	// Add estimated delivery time based on priority
	switch req.Priority {
	case "high":
		response.EstimatedTime = "2-5 minutes"
	case "normal":
		response.EstimatedTime = "5-10 minutes"
	case "low":
		response.EstimatedTime = "10-30 minutes"
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data": response,
	})
}

// ListNotifications lists notifications with pagination
func (h *Handler) ListNotifications(c *gin.Context) {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	search := c.Query("search")
	status := c.Query("status")

	// Create list request
	req := &services.ListRequest{
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
		Search: search,
		Status: status,
	}

	// Get notifications
	notifications, err := h.notificationService.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list notifications",
			"details": err.Error(),
		})
		return
	}

	// Get total count
	countReq := &services.CountRequest{
		Search: search,
		Status: status,
	}
	total, err := h.notificationService.Count(c.Request.Context(), countReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to count notifications",
			"details": err.Error(),
		})
		return
	}

	// Calculate pages
	pages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data": notifications,
		"pagination": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"pages":  pages,
		},
	})
}

// GetNotification gets a notification by ID
func (h *Handler) GetNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid notification ID",
		})
		return
	}

	notification, err := h.notificationService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Notification not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": notification,
	})
}

// RetryNotification retries a failed notification
func (h *Handler) RetryNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid notification ID",
		})
		return
	}

	err = h.notificationService.RetryNotification(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Notification not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retry notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification retry initiated successfully",
	})
}

// CancelNotification cancels a pending notification
func (h *Handler) CancelNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid notification ID",
		})
		return
	}

	err = h.notificationService.CancelNotification(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Notification not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to cancel notification",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification cancelled successfully",
	})
}

// GetNotificationsByProject gets notifications by project ID
func (h *Handler) GetNotificationsByProject(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	notifications, err := h.notificationService.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get notifications by project",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": notifications,
	})
}

// GetNotificationsBySender gets notifications by sender ID
func (h *Handler) GetNotificationsBySender(c *gin.Context) {
	senderIDStr := c.Param("senderId")
	senderID, err := uuid.Parse(senderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid sender ID",
		})
		return
	}

	notifications, err := h.notificationService.GetBySenderID(c.Request.Context(), senderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get notifications by sender",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": notifications,
	})
}

// GetNotificationsByStatus gets notifications by status
func (h *Handler) GetNotificationsByStatus(c *gin.Context) {
	status := c.Param("status")

	notifications, err := h.notificationService.GetByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get notifications by status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": notifications,
	})
}

// GetNotificationsByDateRange gets notifications by date range
func (h *Handler) GetNotificationsByDateRange(c *gin.Context) {
	var req DateRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Parse dates
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start date format",
		})
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end date format",
		})
		return
	}

	// Validate date range
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Start date must be before end date",
		})
		return
	}

	notifications, err := h.notificationService.GetByDateRange(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get notifications by date range",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  notifications,
		"count": len(notifications),
		"date_range": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
	})
}
