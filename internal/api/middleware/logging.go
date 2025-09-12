package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware handles request logging
type LoggingMiddleware struct {
	logger *logrus.Logger
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger *logrus.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// RequestLogger logs HTTP requests
func (m *LoggingMiddleware) RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Create structured log entry
		entry := m.logger.WithFields(logrus.Fields{
			"timestamp":  param.TimeStamp.Format(time.RFC3339),
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"method":     param.Method,
			"path":       param.Path,
			"user_agent": param.Request.UserAgent(),
			"request_id": param.Keys["request_id"],
		})

		// Set log level based on status code
		switch {
		case param.StatusCode >= 500:
			entry.Error("HTTP Request")
		case param.StatusCode >= 400:
			entry.Warn("HTTP Request")
		default:
			entry.Info("HTTP Request")
		}

		return ""
	})
}

// RequestBodyLogger logs request body for debugging
func (m *LoggingMiddleware) RequestBodyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only log for POST, PUT, PATCH requests
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			// Read request body
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				m.logger.WithError(err).Error("Failed to read request body")
				c.Next()
				return
			}

			// Restore request body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// Log request body (be careful with sensitive data)
			if len(body) > 0 {
				var bodyData interface{}
				if err := json.Unmarshal(body, &bodyData); err == nil {
					m.logger.WithFields(logrus.Fields{
						"request_id":   c.GetString("request_id"),
						"method":       c.Request.Method,
						"path":         c.Request.URL.Path,
						"request_body": bodyData,
					}).Debug("Request body")
				}
			}
		}

		c.Next()
	}
}

// ResponseLogger logs response data
func (m *LoggingMiddleware) ResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture response
		blw := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		// Log response
		statusCode := c.Writer.Status()
		responseBody := blw.body.String()

		// Only log response body for certain status codes or in debug mode
		if statusCode >= 400 || m.logger.Level == logrus.DebugLevel {
			var responseData interface{}
			if err := json.Unmarshal([]byte(responseBody), &responseData); err == nil {
				m.logger.WithFields(logrus.Fields{
					"request_id":    c.GetString("request_id"),
					"status_code":   statusCode,
					"response_body": responseData,
				}).Debug("Response body")
			}
		}

		// Log response summary
		entry := m.logger.WithFields(logrus.Fields{
			"request_id":    c.GetString("request_id"),
			"status_code":   statusCode,
			"latency":       time.Since(c.GetTime("start_time")),
			"response_size": len(responseBody),
		})

		switch {
		case statusCode >= 500:
			entry.Error("Response sent")
		case statusCode >= 400:
			entry.Warn("Response sent")
		default:
			entry.Info("Response sent")
		}
	}
}

// bodyLogWriter wraps gin.ResponseWriter to capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// StartTimeMiddleware records request start time
func (m *LoggingMiddleware) StartTimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("start_time", time.Now())
		c.Next()
	}
}

// ErrorLogger logs errors
func (m *LoggingMiddleware) ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				m.logger.WithFields(logrus.Fields{
					"request_id": c.GetString("request_id"),
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"error":      err.Error(),
					"type":       err.Type,
				}).Error("Request error")
			}
		}
	}
}

// SecurityHeaders adds security headers
func (m *LoggingMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

// RateLimitLogger logs rate limit events
func (m *LoggingMiddleware) RateLimitLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if rate limited
		if c.Writer.Status() == 429 {
			m.logger.WithFields(logrus.Fields{
				"request_id": c.GetString("request_id"),
				"client_ip":  c.ClientIP(),
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"user_agent": c.Request.UserAgent(),
			}).Warn("Rate limit exceeded")
		}
	}
}
