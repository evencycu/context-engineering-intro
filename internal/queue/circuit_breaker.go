package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// CircuitBreakerState represents the state of the circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota // 正常狀態，允許請求
	StateOpen                              // 熔斷狀態，拒絕請求
	StateHalfOpen                          // 半開狀態，允許部分請求測試
)

// CircuitBreaker implements the circuit breaker pattern for Teams API calls
type CircuitBreaker struct {
	mu sync.RWMutex

	// Configuration
	maxFailures     uint32        // 觸發熔斷的失敗次數閾值
	timeout         time.Duration // 熔斷後等待時間
	halfOpenMaxReqs uint32        // 半開狀態允許的最大測試請求數

	// State
	state           CircuitBreakerState
	failures        uint32
	lastFailureTime time.Time
	lastStateChange time.Time
	halfOpenReqs    uint32

	// Metrics
	totalRequests   uint64
	successRequests uint64
	failedRequests  uint64
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures uint32, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:     maxFailures,
		timeout:         timeout,
		halfOpenMaxReqs: 3, // 預設半開狀態允許 3 個測試請求
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

var (
	ErrCircuitOpen     = errors.New("circuit breaker is open")
	ErrTooManyRequests = errors.New("too many requests in half-open state")
)

// Execute executes a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	cb.mu.Lock()
	cb.totalRequests++

	// Check if we can execute the request
	if err := cb.canExecute(); err != nil {
		cb.mu.Unlock()
		return err
	}

	// If half-open, increment test request counter
	if cb.state == StateHalfOpen {
		cb.halfOpenReqs++
	}
	cb.mu.Unlock()

	// Execute the function
	err := fn()

	// Record the result
	cb.recordResult(err)

	return err
}

// canExecute checks if the circuit breaker allows the request
func (cb *CircuitBreaker) canExecute() error {
	switch cb.state {
	case StateClosed:
		return nil

	case StateOpen:
		// Check if timeout has passed
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.setState(StateHalfOpen)
			return nil
		}
		return ErrCircuitOpen

	case StateHalfOpen:
		// Limit concurrent requests in half-open state
		if cb.halfOpenReqs >= cb.halfOpenMaxReqs {
			return ErrTooManyRequests
		}
		return nil

	default:
		return nil
	}
}

// recordResult records the result of a request
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err == nil {
		cb.onSuccess()
	} else {
		cb.onFailure()
	}
}

// onSuccess handles successful request
func (cb *CircuitBreaker) onSuccess() {
	cb.successRequests++
	cb.failures = 0

	switch cb.state {
	case StateHalfOpen:
		// If enough successful requests in half-open state, close the circuit
		if cb.halfOpenReqs >= cb.halfOpenMaxReqs {
			cb.setState(StateClosed)
		}
	}
}

// onFailure handles failed request
func (cb *CircuitBreaker) onFailure() {
	cb.failedRequests++
	cb.failures++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		// If failures exceed threshold, open the circuit
		if cb.failures >= cb.maxFailures {
			cb.setState(StateOpen)
		}

	case StateHalfOpen:
		// If any failure in half-open state, reopen the circuit
		cb.setState(StateOpen)
	}
}

// setState changes the circuit breaker state
func (cb *CircuitBreaker) setState(newState CircuitBreakerState) {
	if cb.state == newState {
		return
	}

	cb.state = newState
	cb.lastStateChange = time.Now()

	// Reset counters when changing state
	if newState == StateHalfOpen {
		cb.halfOpenReqs = 0
	} else if newState == StateClosed {
		cb.failures = 0
		cb.halfOpenReqs = 0
	}
}

// GetState returns the current state
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetMetrics returns circuit breaker metrics
func (cb *CircuitBreaker) GetMetrics() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"state":            cb.state.String(),
		"total_requests":   cb.totalRequests,
		"success_requests": cb.successRequests,
		"failed_requests":  cb.failedRequests,
		"failures":         cb.failures,
		"last_state_change": cb.lastStateChange,
	}
}

// String returns string representation of state
func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.setState(StateClosed)
	cb.failures = 0
	cb.halfOpenReqs = 0
}

