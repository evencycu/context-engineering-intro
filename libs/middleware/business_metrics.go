package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// BusinessMetricsMiddleware logs business metrics events
type BusinessMetricsMiddleware struct {
	logger *logrus.Logger
}

// NewBusinessMetricsMiddleware creates a new business metrics middleware
func NewBusinessMetricsMiddleware(logger *logrus.Logger) *BusinessMetricsMiddleware {
	return &BusinessMetricsMiddleware{
		logger: logger,
	}
}

// BusinessMetricsLogger logs business metrics events
func (m *BusinessMetricsMiddleware) BusinessMetricsLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// Log business metrics for specific endpoints
		path := c.Request.URL.Path
		method := c.Request.Method
		status := c.Writer.Status()

		// Log notification-related metrics
		if method == "POST" && (path == "/api/v1/notifications" || path == "/api/v1/provision") {
			latency := time.Since(start)
			m.logger.WithFields(logrus.Fields{
				"request_id":     c.GetString("request_id"),
				"endpoint":       path,
				"method":         method,
				"status":         status,
				"latency_ms":     latency.Milliseconds(),
				"business_event": "notification_created",
			}).Info("Business Event - Notification Created")
		}

		// Log metrics endpoint access
		if path == "/api/v1/metrics" {
			m.logger.WithFields(logrus.Fields{
				"request_id":     c.GetString("request_id"),
				"endpoint":       path,
				"method":         method,
				"status":         status,
				"latency_ms":     time.Since(start).Milliseconds(),
				"business_event": "metrics_accessed",
			}).Info("Business Event - Metrics Accessed")
		}

		// Log monitoring endpoint access
		if path == "/api/v1/monitoring/health" || path == "/api/v1/monitoring/performance" ||
			path == "/api/v1/monitoring/business" || path == "/api/v1/monitoring/alerts" ||
			path == "/api/v1/monitoring/dashboard" {
			m.logger.WithFields(logrus.Fields{
				"request_id":     c.GetString("request_id"),
				"endpoint":       path,
				"method":         method,
				"status":         status,
				"latency_ms":     time.Since(start).Milliseconds(),
				"business_event": "monitoring_accessed",
			}).Info("Business Event - Monitoring Accessed")
		}
	}
}
