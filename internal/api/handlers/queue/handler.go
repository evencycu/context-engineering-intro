package queue

import (
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV2/internal/queue"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles queue management endpoints
type Handler struct {
	queueManager *queue.Manager
}

// NewHandler creates a new queue handler
func NewHandler(queueManager *queue.Manager) *Handler {
	return &Handler{
		queueManager: queueManager,
	}
}

// GetStatus returns the current queue status
// @Summary Get queue status
// @Description Get the current status of the notification retry queue
// @Tags queue
// @Produce json
// @Success 200 {object} queue.NotificationQueueStatus
// @Router /queue/status [get]
func (h *Handler) GetStatus(c *gin.Context) {
	status, err := h.queueManager.GetQueueStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetCircuitBreakerMetrics returns circuit breaker metrics
// @Summary Get circuit breaker metrics
// @Description Get metrics and statistics for the circuit breaker
// @Tags queue
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /queue/circuit-breaker/metrics [get]
func (h *Handler) GetCircuitBreakerMetrics(c *gin.Context) {
	metrics := h.queueManager.GetCircuitBreakerMetrics()
	c.JSON(http.StatusOK, metrics)
}

// ResetCircuitBreaker manually resets the circuit breaker
// @Summary Reset circuit breaker
// @Description Manually reset the circuit breaker to closed state
// @Tags queue
// @Success 200 {object} map[string]string
// @Router /queue/circuit-breaker/reset [post]
func (h *Handler) ResetCircuitBreaker(c *gin.Context) {
	h.queueManager.ResetCircuitBreaker()
	c.JSON(http.StatusOK, gin.H{"message": "Circuit breaker reset successfully"})
}

// GetFailedNotifications returns failed notifications by notification ID
// @Summary Get failed notifications
// @Description Get failed notifications for a specific notification ID
// @Tags queue
// @Param notification_id path string true "Notification ID"
// @Produce json
// @Success 200 {array} queue.FailedNotification
// @Router /queue/failed/{notification_id} [get]
func (h *Handler) GetFailedNotifications(c *gin.Context) {
	notificationIDStr := c.Param("notification_id")
	_, err := uuid.Parse(notificationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// This requires adding a method to queue manager
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}

// RegisterRoutes registers queue routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/queue")
	{
		g.GET("/status", h.GetStatus)
		g.GET("/circuit-breaker/metrics", h.GetCircuitBreakerMetrics)
		g.POST("/circuit-breaker/reset", h.ResetCircuitBreaker)
		// g.GET("/failed/:notification_id", h.GetFailedNotifications)
	}
}

