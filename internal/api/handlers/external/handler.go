package external

import (
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/services"
	"github.com/gin-gonic/gin"
)

// Handler handles external user API requests
type Handler struct {
	externalService services.ExternalService
}

// NewHandler creates a new external handler
func NewHandler(externalService services.ExternalService) *Handler {
	return &Handler{
		externalService: externalService,
	}
}

// RegisterRoutes registers external API routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// External API routes (no authentication required)
	external := rg.Group("/external")
	{
		// Send notification
		external.POST("/notify", h.SendNotification)


		// Get project destinations
		external.GET("/destinations/:notifyKey", h.GetProjectDestinations)

		// Health check for external API
		external.GET("/health", h.HealthCheck)
	}
}

// SendNotificationRequest represents the request to send notification
type SendNotificationRequest struct {
	NotifyKey   string         `json:"notify_key" binding:"required" example:"my-project-key"`
	Message     string         `json:"message" binding:"required" example:"Hello from external API"`
	MessageType string         `json:"message_type" example:"text" enums:"text,file,adaptive_card"`
	Priority    string         `json:"priority" example:"normal" enums:"low,normal,high,urgent"`
	Targets     []string       `json:"targets" example:"all" description:"List of target Emails or Conversation IDs or 'all'"`
	Mentions    []string       `json:"mentions" example:"@user1,@user2"`
	Metadata    map[string]any `json:"metadata" example:"{\"source\":\"external\",\"version\":\"1.0\"}"`
}

// SendNotificationResponse represents the response for send notification
type SendNotificationResponse struct {
	Success bool                                   `json:"success"`
	Data    *services.ExternalNotificationResponse `json:"data,omitempty"`
	Error   string                                 `json:"error,omitempty"`
}

// SendNotification handles external notification sending
func (h *Handler) SendNotification(c *gin.Context) {
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SendNotificationResponse{
			Success: false,
			Error:   "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if req.NotifyKey == "" {
		c.JSON(http.StatusBadRequest, SendNotificationResponse{
			Success: false,
			Error:   "notify_key is required",
		})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusBadRequest, SendNotificationResponse{
			Success: false,
			Error:   "message is required",
		})
		return
	}

	// Set defaults
	if req.MessageType == "" {
		req.MessageType = "text"
	}
	if req.Priority == "" {
		req.Priority = "normal"
	}

	// Convert to service request
	serviceReq := &services.ExternalNotificationRequest{
		NotifyKey:   req.NotifyKey,
		Message:     req.Message,
		MessageType: req.MessageType,
		Priority:    req.Priority,
		TargetIDs:   req.Targets,
		Mentions:    req.Mentions,
		Metadata:    req.Metadata,
	}

	// Send notification
	response, err := h.externalService.SendNotification(c.Request.Context(), serviceReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, SendNotificationResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SendNotificationResponse{
		Success: true,
		Data:    response,
	})
}


// GetProjectDestinationsRequest represents the request to get project destinations
type GetProjectDestinationsRequest struct {
	NotifyKey string `uri:"notifyKey" binding:"required" example:"my-project-key"`
}

// GetProjectDestinationsResponse represents the response for get project destinations
type GetProjectDestinationsResponse struct {
	Success      bool                               `json:"success"`
	Destinations []services.ExternalDestinationInfo `json:"destinations,omitempty"`
	Error        string                             `json:"error,omitempty"`
}

// GetProjectDestinations gets available destinations for a project
func (h *Handler) GetProjectDestinations(c *gin.Context) {
	var req GetProjectDestinationsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetProjectDestinationsResponse{
			Success: false,
			Error:   "Invalid notify key format",
		})
		return
	}

	// For now, return empty destinations - TODO: Implement proper destinations endpoint
	destinations := []services.ExternalDestinationInfo{}

	c.JSON(http.StatusOK, GetProjectDestinationsResponse{
		Success:      true,
		Destinations: destinations,
	})
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

// HealthCheck provides health check for external API
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckResponse{
		Status:    "healthy",
		Service:   "external-api",
		Timestamp: "2025-09-26T03:51:33Z",
	})
}
