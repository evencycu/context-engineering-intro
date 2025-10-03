package queue

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryPolicy defines the retry behavior for failed notifications
type RetryPolicy struct {
	MaxRetries     int           // 最大重試次數
	InitialBackoff time.Duration // 初始退避時間
	MaxBackoff     time.Duration // 最大退避時間
	Multiplier     float64       // 退避倍數
	JitterFraction float64       // 抖動比例 (0.0-1.0)
}

// DefaultRetryPolicy returns a default retry policy
// 使用指數退避 + 隨機抖動
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries:     5,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     5 * time.Minute,
		Multiplier:     2.0,
		JitterFraction: 0.1, // 10% 隨機抖動
	}
}

// CalculateBackoff calculates the backoff duration for a given retry attempt
// Uses exponential backoff with jitter
func (p *RetryPolicy) CalculateBackoff(retryCount int) time.Duration {
	if retryCount <= 0 {
		return p.InitialBackoff
	}

	// Calculate exponential backoff
	backoff := float64(p.InitialBackoff) * math.Pow(p.Multiplier, float64(retryCount))

	// Cap at max backoff
	if backoff > float64(p.MaxBackoff) {
		backoff = float64(p.MaxBackoff)
	}

	// Add jitter to avoid thundering herd
	jitter := backoff * p.JitterFraction * (rand.Float64()*2 - 1) // -jitter to +jitter
	backoffWithJitter := backoff + jitter

	// Ensure non-negative
	if backoffWithJitter < 0 {
		backoffWithJitter = backoff
	}

	return time.Duration(backoffWithJitter)
}

// ShouldRetry determines if a notification should be retried
func (p *RetryPolicy) ShouldRetry(retryCount int) bool {
	return retryCount < p.MaxRetries
}

// GetRetryAfterDuration parses Retry-After header and returns duration
// Retry-After can be in seconds (integer) or HTTP-date format
func GetRetryAfterDuration(retryAfter string) time.Duration {
	if retryAfter == "" {
		return 0
	}

	// Try to parse as seconds
	var seconds int
	if _, err := time.Parse(time.RFC1123, retryAfter); err == nil {
		// HTTP-date format
		t, _ := time.Parse(time.RFC1123, retryAfter)
		return time.Until(t)
	}

	// Try as integer seconds
	if n, err := time.ParseDuration(retryAfter + "s"); err == nil {
		return n
	}

	// Parse as integer
	if _, err := fmt.Sscanf(retryAfter, "%d", &seconds); err == nil {
		return time.Duration(seconds) * time.Second
	}

	return 0
}

// RateLimitBackoff calculates backoff for rate limit (429) errors
// Respects Retry-After header if present
func RateLimitBackoff(retryAfterHeader string, defaultBackoff time.Duration) time.Duration {
	if retryAfter := GetRetryAfterDuration(retryAfterHeader); retryAfter > 0 {
		// Add small buffer to retry-after value
		return retryAfter + (time.Second * 1)
	}
	return defaultBackoff
}

