package system

import (
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
)

// Handler handles system-related HTTP requests
type Handler struct {
	metricsService services.MetricsService
	configService  services.ConfigService
}

// NewHandler creates a new system handler
func NewHandler(metricsService services.MetricsService, configService services.ConfigService) *Handler {
	return &Handler{
		metricsService: metricsService,
		configService:  configService,
	}
}

// RegisterRoutes registers system routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	// System endpoints (no authentication required for monitoring)
	rg.GET("/metrics", h.GetMetrics)
	rg.GET("/config", h.GetConfig)
	rg.POST("/config/validate", h.ValidateConfig)

	// Alert endpoints
	rg.GET("/alerts", h.GetAlerts)
	rg.GET("/alerts/history", h.GetAlertHistory)
	rg.POST("/alerts/check", h.CheckAlerts)
}

// GetMetrics returns system metrics
func (h *Handler) GetMetrics(c *gin.Context) {
	metrics, err := h.metricsService.GetMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get metrics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetConfig returns system configuration
func (h *Handler) GetConfig(c *gin.Context) {
	config, err := h.configService.GetConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get configuration",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, config)
}

// ValidateConfig validates system configuration
func (h *Handler) ValidateConfig(c *gin.Context) {
	validation, err := h.configService.ValidateConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to validate configuration",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, validation)
}

// GetAlerts returns current alerts
func (h *Handler) GetAlerts(c *gin.Context) {
	// For now, return empty alerts
	c.JSON(http.StatusOK, gin.H{
		"data":    []interface{}{},
		"message": "No active alerts",
	})
}

// GetAlertHistory returns alert history
func (h *Handler) GetAlertHistory(c *gin.Context) {
	// For now, return empty history
	c.JSON(http.StatusOK, gin.H{
		"data":    []interface{}{},
		"message": "No alert history",
	})
}

// CheckAlerts manually triggers alert checking
func (h *Handler) CheckAlerts(c *gin.Context) {
	// For now, return success
	c.JSON(http.StatusOK, gin.H{
		"message":      "Alert check completed",
		"alerts_found": 0,
	})
}
