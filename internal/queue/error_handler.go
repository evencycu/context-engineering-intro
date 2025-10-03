package queue

import (
	"net/http"
	"strconv"
	"strings"
)

// IsRateLimitError checks if an HTTP response is a rate limit error
func IsRateLimitError(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests // 429
}

// IsRetryableError checks if an error should be retried
func IsRetryableError(statusCode int) bool {
	// Retry on:
	// - 429 (Rate Limit)
	// - 408 (Request Timeout)
	// - 5xx (Server Errors)
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusRequestTimeout ||
		(statusCode >= 500 && statusCode < 600)
}

// ExtractRetryAfter extracts the Retry-After header value (in seconds)
func ExtractRetryAfter(headers http.Header) *int {
	retryAfter := headers.Get("Retry-After")
	if retryAfter == "" {
		return nil
	}

	// Try to parse as integer seconds
	if seconds, err := strconv.Atoi(strings.TrimSpace(retryAfter)); err == nil {
		return &seconds
	}

	// If it's an HTTP-date format, we'll handle it in GetRetryAfterDuration
	return nil
}

// ClassifyFailureReason determines the failure reason from HTTP status code
func ClassifyFailureReason(statusCode int) FailureReason {
	switch {
	case statusCode == http.StatusTooManyRequests:
		return FailureReasonRateLimit
	case statusCode == http.StatusRequestTimeout:
		return FailureReasonTimeout
	case statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden:
		return FailureReasonUnauthorized
	case statusCode >= 500 && statusCode < 600:
		return FailureReasonServerError
	case statusCode == 0:
		return FailureReasonNetworkError
	default:
		return FailureReasonUnknown
	}
}

