package actor

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCircuitBreaker implements CircuitBreaker interface using Redis
type redisCircuitBreaker struct {
	redis     *redis.Client
	key       string
	timeout   time.Duration
	threshold int
}

// CircuitBreaker defines the interface for circuit breaker operations
type CircuitBreaker interface {
	GetState(ctx context.Context) (string, error)
	RecordSuccess(ctx context.Context) error
	RecordFailure(ctx context.Context, retryAfter int) error
}

// NewRedisCircuitBreaker creates a new Redis-backed circuit breaker
func NewRedisCircuitBreaker(redis *redis.Client) CircuitBreaker {
	return &redisCircuitBreaker{
		redis:     redis,
		key:       "circuit_breaker:teams_api",
		timeout:   5 * time.Minute, // Circuit breaker timeout
		threshold: 5,               // Failure threshold
	}
}

// GetState returns the current circuit breaker state
func (cb *redisCircuitBreaker) GetState(ctx context.Context) (string, error) {
	state, err := cb.redis.Get(ctx, cb.key+":state").Result()
	if err == redis.Nil {
		return "closed", nil // Default to closed
	}
	if err != nil {
		return "", err
	}
	return state, nil
}

// RecordSuccess records a successful operation
func (cb *redisCircuitBreaker) RecordSuccess(ctx context.Context) error {
	// Reset failure count and set state to closed
	pipe := cb.redis.Pipeline()
	pipe.Del(ctx, cb.key+":failures")
	pipe.Set(ctx, cb.key+":state", "closed", 0)
	pipe.Expire(ctx, cb.key+":state", cb.timeout)
	_, err := pipe.Exec(ctx)
	return err
}

// RecordFailure records a failed operation
func (cb *redisCircuitBreaker) RecordFailure(ctx context.Context, retryAfter int) error {
	// Increment failure count
	failures, err := cb.redis.Incr(ctx, cb.key+":failures").Result()
	if err != nil {
		return err
	}

	// Set expiration for failure count
	cb.redis.Expire(ctx, cb.key+":failures", cb.timeout)

	// Check if threshold is reached
	if failures >= int64(cb.threshold) {
		// Open circuit breaker
		openUntil := time.Now().Add(cb.timeout)
		if retryAfter > 0 {
			// Use Retry-After if provided
			openUntil = time.Now().Add(time.Duration(retryAfter) * time.Second)
		}

		pipe := cb.redis.Pipeline()
		pipe.Set(ctx, cb.key+":state", "open", 0)
		pipe.Set(ctx, cb.key+":open_until", openUntil.Format(time.RFC3339), 0)
		pipe.Expire(ctx, cb.key+":open_until", cb.timeout)
		_, err = pipe.Exec(ctx)

		log.Printf("Circuit breaker opened due to %d failures (open until %v)", failures, openUntil)
		return err
	}

	return nil
}

// IsOpen checks if circuit breaker is currently open
func (cb *redisCircuitBreaker) IsOpen(ctx context.Context) (bool, error) {
	state, err := cb.GetState(ctx)
	if err != nil {
		return false, err
	}

	if state != "open" {
		return false, nil
	}

	// Check if timeout has passed
	openUntilStr, err := cb.redis.Get(ctx, cb.key+":open_until").Result()
	if err == redis.Nil {
		// No open_until timestamp, consider it closed
		cb.redis.Set(ctx, cb.key+":state", "closed", 0)
		return false, nil
	}
	if err != nil {
		return false, err
	}

	openUntil, err := time.Parse(time.RFC3339, openUntilStr)
	if err != nil {
		return false, err
	}

	if time.Now().After(openUntil) {
		// Timeout passed, move to half-open
		cb.redis.Set(ctx, cb.key+":state", "half-open", 0)
		return false, nil
	}

	return true, nil
}

// GetMetrics returns circuit breaker metrics
func (cb *redisCircuitBreaker) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	state, err := cb.GetState(ctx)
	if err != nil {
		return nil, err
	}

	failures, err := cb.redis.Get(ctx, cb.key+":failures").Result()
	if err == redis.Nil {
		failures = "0"
	} else if err != nil {
		return nil, err
	}

	failureCount, _ := strconv.Atoi(failures)

	metrics := map[string]interface{}{
		"state":         state,
		"failure_count": failureCount,
		"threshold":     cb.threshold,
		"timeout":       cb.timeout.String(),
	}

	if state == "open" {
		openUntilStr, err := cb.redis.Get(ctx, cb.key+":open_until").Result()
		if err == nil {
			openUntil, err := time.Parse(time.RFC3339, openUntilStr)
			if err == nil {
				metrics["open_until"] = openUntil.Format(time.RFC3339)
			}
		}
	}

	return metrics, nil
}
