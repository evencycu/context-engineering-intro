package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents error codes
type ErrorCode string

const (
	// Authentication errors
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrCodeInvalidToken  ErrorCode = "INVALID_TOKEN"
	ErrCodeTokenExpired  ErrorCode = "TOKEN_EXPIRED"
	ErrCodeInvalidAPIKey ErrorCode = "INVALID_API_KEY"

	// Teams API errors
	ErrCodeTeamsAuthFailed ErrorCode = "TEAMS_AUTH_FAILED"
	ErrCodeTeamsAPIError   ErrorCode = "TEAMS_API_ERROR"
	ErrCodeTeamsRateLimit  ErrorCode = "TEAMS_RATE_LIMIT"

	// Database errors
	ErrCodeDatabaseError   ErrorCode = "DATABASE_ERROR"
	ErrCodeRecordNotFound  ErrorCode = "RECORD_NOT_FOUND"
	ErrCodeDuplicateRecord ErrorCode = "DUPLICATE_RECORD"

	// Validation errors
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField     ErrorCode = "MISSING_FIELD"

	// System errors
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeTimeout            ErrorCode = "TIMEOUT"
)

// AppError represents application error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
	Cause      error     `json:"-"`
}

// Error implements error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, cause error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    "",
		HTTPStatus: getHTTPStatus(code),
		Cause:      cause,
	}
}

// NewAppErrorWithDetails creates a new application error with details
func NewAppErrorWithDetails(code ErrorCode, message, details string, cause error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: getHTTPStatus(code),
		Cause:      cause,
	}
}

// getHTTPStatus returns HTTP status code for error code
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeUnauthorized, ErrCodeInvalidToken, ErrCodeTokenExpired, ErrCodeInvalidAPIKey:
		return http.StatusUnauthorized
	case ErrCodeValidationFailed, ErrCodeInvalidInput, ErrCodeMissingField:
		return http.StatusBadRequest
	case ErrCodeRecordNotFound:
		return http.StatusNotFound
	case ErrCodeDuplicateRecord:
		return http.StatusConflict
	case ErrCodeTeamsRateLimit:
		return http.StatusTooManyRequests
	case ErrCodeServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrCodeTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}

// Predefined error constructors
func NewUnauthorizedError(message string, cause error) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, cause)
}

func NewInvalidTokenError(cause error) *AppError {
	return NewAppError(ErrCodeInvalidToken, "Invalid or malformed token", cause)
}

func NewTokenExpiredError() *AppError {
	return NewAppError(ErrCodeTokenExpired, "Token has expired", nil)
}

func NewTeamsAuthFailedError(cause error) *AppError {
	return NewAppError(ErrCodeTeamsAuthFailed, "Teams authentication failed", cause)
}

func NewTeamsAPIError(message string, cause error) *AppError {
	return NewAppError(ErrCodeTeamsAPIError, message, cause)
}

func NewTeamsRateLimitError(retryAfter int) *AppError {
	return NewAppErrorWithDetails(
		ErrCodeTeamsRateLimit,
		"Teams API rate limit exceeded",
		fmt.Sprintf("Retry after %d seconds", retryAfter),
		nil,
	)
}

func NewDatabaseError(operation string, cause error) *AppError {
	return NewAppError(ErrCodeDatabaseError, fmt.Sprintf("Database operation failed: %s", operation), cause)
}

func NewRecordNotFoundError(resource string) *AppError {
	return NewAppError(ErrCodeRecordNotFound, fmt.Sprintf("%s not found", resource), nil)
}

func NewValidationError(field, message string) *AppError {
	return NewAppErrorWithDetails(
		ErrCodeValidationFailed,
		"Validation failed",
		fmt.Sprintf("Field '%s': %s", field, message),
		nil,
	)
}

func NewInternalError(operation string, cause error) *AppError {
	return NewAppError(ErrCodeInternalError, fmt.Sprintf("Internal error in %s", operation), cause)
}
