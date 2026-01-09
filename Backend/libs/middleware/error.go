package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandler handles errors globally
type ErrorHandler struct {
	logger *logrus.Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *logrus.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// ErrorHandlerMiddleware handles errors globally
func (h *ErrorHandler) ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Log error
			h.logger.WithFields(logrus.Fields{
				"request_id": c.GetString("request_id"),
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"error":      err.Error(),
				"type":       err.Type,
			}).Error("Request error")

			// Return appropriate error response
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid request data",
					"details": err.Error(),
				})
			case gin.ErrorTypeRender:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to render response",
					"details": err.Error(),
				})
			case gin.ErrorTypePublic:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"details": err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}
	}
}

// NotFoundHandler handles 404 errors
func (h *ErrorHandler) NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Resource not found",
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})
	}
}

// MethodNotAllowedHandler handles 405 errors
func (h *ErrorHandler) MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error":  "Method not allowed",
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})
	}
}

// PanicRecovery recovers from panics
func (h *ErrorHandler) PanicRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log panic
		h.logger.WithFields(logrus.Fields{
			"request_id": c.GetString("request_id"),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"panic":      recovered,
		}).Error("Panic recovered")

		// Return error response
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrorResponse represents a validation error response
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Details []ValidationError `json:"details"`
}

// HandleValidationError handles validation errors
func (h *ErrorHandler) HandleValidationError(c *gin.Context, err error) {
	// TODO: Implement proper validation error parsing
	// This would typically parse validation errors from a validator library

	c.JSON(http.StatusBadRequest, ValidationErrorResponse{
		Error: "Validation failed",
		Details: []ValidationError{
			{
				Field:   "unknown",
				Message: err.Error(),
			},
		},
	})
}

// HandleServiceError handles service layer errors
func (h *ErrorHandler) HandleServiceError(c *gin.Context, err error) {
	// Check error type and return appropriate response
	switch e := err.(type) {
	case *NotFoundError:
		c.JSON(http.StatusNotFound, gin.H{
			"error": e.Error(),
		})
	case *ConflictError:
		c.JSON(http.StatusConflict, gin.H{
			"error": e.Error(),
		})
	case *BadRequestError:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": e.Error(),
		})
	default:
		// Log unexpected error
		h.logger.WithFields(logrus.Fields{
			"request_id": c.GetString("request_id"),
			"error":      err.Error(),
		}).Error("Unexpected service error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	}
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return e.Resource + " not found"
}

// ConflictError represents a conflict error
type ConflictError struct {
	Resource string
	Field    string
	Value    string
}

func (e *ConflictError) Error() string {
	return e.Resource + " with " + e.Field + " " + e.Value + " already exists"
}

// UnauthorizedError represents an unauthorized error
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

// ForbiddenError represents a forbidden error
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

// BadRequestError represents a bad request error
type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

// InternalServerError represents an internal server error
type InternalServerError struct {
	Message string
}

func (e *InternalServerError) Error() string {
	return e.Message
}
