package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// JSONResponse is the standard API response structure
type JSONResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// Success responds with a success message and optional data
func Success(c *gin.Context, httpStatus int, message string, data interface{}) {
	c.JSON(httpStatus, JSONResponse{
		Code:    httpStatus,
		Message: message,
		Data:    data,
	})
}

// Error responds with an error message and optional details
func Error(c *gin.Context, httpStatus int, message string, err error, details interface{}) {
	logrus.WithFields(logrus.Fields{
		"httpStatus": httpStatus,
		"message":    message,
		"error":      err,
		"details":    details,
		"request_id": c.GetString("request_id"), // Assuming request_id middleware is used
	}).Error("API Error")

	c.JSON(httpStatus, JSONResponse{
		Code:    httpStatus,
		Message: message,
		Error:   err.Error(),
		Details: details,
	})
}

// BadRequest responds with a 400 Bad Request error
func BadRequest(c *gin.Context, message string, err error, details interface{}) {
	Error(c, http.StatusBadRequest, message, err, details)
}

// Unauthorized responds with a 401 Unauthorized error
func Unauthorized(c *gin.Context, message string, err error, details interface{}) {
	Error(c, http.StatusUnauthorized, message, err, details)
}

// Forbidden responds with a 403 Forbidden error
func Forbidden(c *gin.Context, message string, err error, details interface{}) {
	Error(c, http.StatusForbidden, message, err, details)
}

// NotFound responds with a 404 Not Found error
func NotFound(c *gin.Context, message string, err error, details interface{}) {
	Error(c, http.StatusNotFound, message, err, details)
}

// InternalServerError responds with a 500 Internal Server Error
func InternalServerError(c *gin.Context, message string, err error, details interface{}) {
	Error(c, http.StatusInternalServerError, message, err, details)
}
