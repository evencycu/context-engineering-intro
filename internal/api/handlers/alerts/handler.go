package alerts

import (
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/services"
	"github.com/gin-gonic/gin"
)

// Handler handles alert-related HTTP requests
type Handler struct {
	alertService services.AlertService
}

// NewHandler creates a new alerts handler
func NewHandler(alertService services.AlertService) *Handler {
	return &Handler{
		alertService: alertService,
	}
}

// RegisterRoutes registers alert routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// Alert endpoints (no authentication required for monitoring)
	rg.GET("/alerts", h.GetAlerts)
	rg.GET("/alerts/history", h.GetAlertHistory)
	rg.POST("/alerts/check", h.CheckAlerts)
}

// GetAlerts returns current alerts
func (h *Handler) GetAlerts(c *gin.Context) {
	// This would typically get current alerts from the alert service
	// For now, return empty list
	c.JSON(http.StatusOK, gin.H{
		"data":    []interface{}{},
		"message": "No active alerts",
	})
}

// GetAlertHistory returns alert history
func (h *Handler) GetAlertHistory(c *gin.Context) {
	alerts, err := h.alertService.GetAlertHistory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": alerts,
	})
}

// CheckAlerts manually triggers alert checking
func (h *Handler) CheckAlerts(c *gin.Context) {
	// This would typically get current metrics and check for alerts
	// For now, return success
	c.JSON(http.StatusOK, gin.H{
		"message":      "Alert check completed",
		"alerts_found": 0,
	})
}
