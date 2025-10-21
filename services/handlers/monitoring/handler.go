package monitoring

import (
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV2/services/services"
	"github.com/gin-gonic/gin"
)

// MonitoringHandler handles monitoring-related requests
type MonitoringHandler struct {
	monitoringService services.MonitoringService
}

// NewMonitoringHandler creates a new monitoring handler
func NewMonitoringHandler(monitoringService services.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{
		monitoringService: monitoringService,
	}
}

// GetSystemHealth returns system health status
func (h *MonitoringHandler) GetSystemHealth(c *gin.Context) {
	health, err := h.monitoringService.GetSystemHealth(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get system health",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    health,
	})
}

// GetPerformanceMetrics returns performance metrics
func (h *MonitoringHandler) GetPerformanceMetrics(c *gin.Context) {
	metrics, err := h.monitoringService.GetPerformanceMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get performance metrics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics,
	})
}

// GetBusinessHealth returns business health metrics
func (h *MonitoringHandler) GetBusinessHealth(c *gin.Context) {
	health, err := h.monitoringService.GetBusinessHealth(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get business health",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    health,
	})
}

// GetAlertStatus returns current alert status
func (h *MonitoringHandler) GetAlertStatus(c *gin.Context) {
	alerts, err := h.monitoringService.GetAlertStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get alert status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    alerts,
	})
}

// GetDashboard returns a comprehensive monitoring dashboard
func (h *MonitoringHandler) GetDashboard(c *gin.Context) {
	// Get all monitoring data
	systemHealth, err1 := h.monitoringService.GetSystemHealth(c.Request.Context())
	performanceMetrics, err2 := h.monitoringService.GetPerformanceMetrics(c.Request.Context())
	businessHealth, err3 := h.monitoringService.GetBusinessHealth(c.Request.Context())
	alertStatus, err4 := h.monitoringService.GetAlertStatus(c.Request.Context())

	// Check for any errors
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get monitoring dashboard data",
			"details": map[string]string{
				"system_health": getErrorString(err1),
				"performance":   getErrorString(err2),
				"business":      getErrorString(err3),
				"alerts":        getErrorString(err4),
			},
		})
		return
	}

	dashboard := gin.H{
		"system_health": systemHealth,
		"performance":   performanceMetrics,
		"business":      businessHealth,
		"alerts":        alertStatus,
		"timestamp":     systemHealth.Timestamp,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dashboard,
	})
}

// RegisterRoutes registers monitoring routes
func (h *MonitoringHandler) RegisterRoutes(rg *gin.RouterGroup) {
	monitoring := rg.Group("/monitoring")
	{
		monitoring.GET("/health", h.GetSystemHealth)
		monitoring.GET("/performance", h.GetPerformanceMetrics)
		monitoring.GET("/business", h.GetBusinessHealth)
		monitoring.GET("/alerts", h.GetAlertStatus)
		monitoring.GET("/dashboard", h.GetDashboard)
	}
}

// Helper function to get error string
func getErrorString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
