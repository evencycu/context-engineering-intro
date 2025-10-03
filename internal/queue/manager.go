package queue

import (
	"context"
	"log"
	"sync"
	"time"
)

// NotificationSender defines the interface for sending notifications
type NotificationSender interface {
	SendToTarget(ctx context.Context, targetID, message string, metadata map[string]interface{}) error
}

// Manager manages the failed notification queue and retry mechanism
type Manager struct {
	repo           Repository
	retryPolicy    *RetryPolicy
	circuitBreaker *CircuitBreaker
	sender         NotificationSender

	// Worker control
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	workerPool int

	// Configuration
	pollInterval  time.Duration
	batchSize     int
	cleanupPeriod time.Duration
}

// Config holds configuration for queue manager
type Config struct {
	WorkerPool    int           // 並發工作數
	PollInterval  time.Duration // 輪詢間隔
	BatchSize     int           // 每次處理的批次大小
	CleanupPeriod time.Duration // 清理舊記錄的週期
	MaxFailures   uint32        // Circuit breaker 失敗閾值
	CBTimeout     time.Duration // Circuit breaker 超時時間
}

// DefaultConfig returns default queue manager configuration
func DefaultConfig() *Config {
	return &Config{
		WorkerPool:    3,
		PollInterval:  10 * time.Second,
		BatchSize:     10,
		CleanupPeriod: 24 * time.Hour,
		MaxFailures:   5,
		CBTimeout:     1 * time.Minute,
	}
}

// NewManager creates a new queue manager
func NewManager(repo Repository, sender NotificationSender, config *Config) *Manager {
	if config == nil {
		config = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		repo:           repo,
		retryPolicy:    DefaultRetryPolicy(),
		circuitBreaker: NewCircuitBreaker(config.MaxFailures, config.CBTimeout),
		sender:         sender,
		ctx:            ctx,
		cancel:         cancel,
		workerPool:     config.WorkerPool,
		pollInterval:   config.PollInterval,
		batchSize:      config.BatchSize,
		cleanupPeriod:  config.CleanupPeriod,
	}
}

// Start starts the queue workers
func (m *Manager) Start() {
	log.Printf("Starting notification queue manager with %d workers", m.workerPool)

	// Start worker pool
	for i := 0; i < m.workerPool; i++ {
		m.wg.Add(1)
		go m.worker(i)
	}

	// Start cleanup worker
	m.wg.Add(1)
	go m.cleanupWorker()

	log.Printf("Queue manager started successfully")
}

// Stop stops the queue workers gracefully
func (m *Manager) Stop() {
	log.Printf("Stopping queue manager...")
	m.cancel()
	m.wg.Wait()
	log.Printf("Queue manager stopped")
}

// worker processes notifications from the queue
func (m *Manager) worker(workerID int) {
	defer m.wg.Done()

	log.Printf("Worker %d started", workerID)
	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			log.Printf("Worker %d stopping", workerID)
			return

		case <-ticker.C:
			m.processQueue(workerID)
		}
	}
}

// processQueue processes a batch of notifications from the queue
func (m *Manager) processQueue(workerID int) {
	ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()

	// Check circuit breaker state
	state := m.circuitBreaker.GetState()
	if state == StateOpen {
		log.Printf("Worker %d: Circuit breaker is open, skipping queue processing", workerID)
		return
	}

	// Dequeue notifications ready for retry
	notifications, err := m.repo.Dequeue(ctx, m.batchSize)
	if err != nil {
		log.Printf("Worker %d: Failed to dequeue notifications: %v", workerID, err)
		return
	}

	if len(notifications) == 0 {
		return
	}

	log.Printf("Worker %d: Processing %d notifications from queue", workerID, len(notifications))

	for _, fn := range notifications {
		if err := m.ctx.Err(); err != nil {
			return
		}

		m.retryNotification(ctx, fn)
	}
}

// retryNotification retries sending a single notification
func (m *Manager) retryNotification(ctx context.Context, fn *FailedNotification) {
	log.Printf("Retrying notification %s (attempt %d/%d)", fn.ID, fn.RetryCount+1, fn.MaxRetries)

	// Execute with circuit breaker protection
	err := m.circuitBreaker.Execute(ctx, func() error {
		return m.sender.SendToTarget(ctx, fn.TargetID, fn.Message, fn.Metadata)
	})

	if err != nil {
		m.handleRetryFailure(ctx, fn, err)
	} else {
		m.handleRetrySuccess(ctx, fn)
	}
}

// handleRetrySuccess handles successful retry
func (m *Manager) handleRetrySuccess(ctx context.Context, fn *FailedNotification) {
	log.Printf("Notification %s retry succeeded", fn.ID)

	// Remove from queue
	if err := m.repo.Delete(ctx, fn.ID); err != nil {
		log.Printf("Failed to remove notification %s from queue: %v", fn.ID, err)
	}
}

// handleRetryFailure handles retry failure
func (m *Manager) handleRetryFailure(ctx context.Context, fn *FailedNotification, err error) {
	log.Printf("Notification %s retry failed: %v", fn.ID, err)

	// Check if we should retry again
	if !m.retryPolicy.ShouldRetry(fn.RetryCount + 1) {
		log.Printf("Notification %s exhausted all retries", fn.ID)
		if err := m.repo.MarkExhausted(ctx, fn.ID); err != nil {
			log.Printf("Failed to mark notification %s as exhausted: %v", fn.ID, err)
		}
		return
	}

	// Calculate next retry time
	backoff := m.retryPolicy.CalculateBackoff(fn.RetryCount + 1)
	nextRetryAt := time.Now().Add(backoff)

	// Update retry information
	if err := m.repo.UpdateRetry(ctx, fn.ID, false, err.Error(), nextRetryAt); err != nil {
		log.Printf("Failed to update retry info for notification %s: %v", fn.ID, err)
	} else {
		log.Printf("Notification %s scheduled for retry at %s (backoff: %s)", 
			fn.ID, nextRetryAt.Format(time.RFC3339), backoff)
	}
}

// cleanupWorker periodically cleans up old exhausted notifications
func (m *Manager) cleanupWorker() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.cleanupPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return

		case <-ticker.C:
			m.cleanup()
		}
	}
}

// cleanup removes old exhausted notifications
func (m *Manager) cleanup() {
	ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()

	// Remove notifications older than 7 days
	deleted, err := m.repo.CleanupOld(ctx, 7*24*time.Hour)
	if err != nil {
		log.Printf("Cleanup failed: %v", err)
		return
	}

	if deleted > 0 {
		log.Printf("Cleaned up %d old exhausted notifications", deleted)
	}
}

// EnqueueFailedNotification adds a failed notification to the retry queue
func (m *Manager) EnqueueFailedNotification(ctx context.Context, fn *FailedNotification) error {
	// Set defaults
	if fn.MaxRetries == 0 {
		fn.MaxRetries = m.retryPolicy.MaxRetries
	}

	// Calculate initial retry time
	if fn.NextRetryAt.IsZero() {
		backoff := m.retryPolicy.CalculateBackoff(fn.RetryCount)
		fn.NextRetryAt = time.Now().Add(backoff)
	}

	log.Printf("Enqueueing failed notification for target %s, next retry at %s", 
		fn.TargetID, fn.NextRetryAt.Format(time.RFC3339))

	return m.repo.Enqueue(ctx, fn)
}

// GetQueueStatus returns the current status of the queue
func (m *Manager) GetQueueStatus(ctx context.Context) (*NotificationQueueStatus, error) {
	status, err := m.repo.GetStatus(ctx)
	if err != nil {
		return nil, err
	}

	// Add circuit breaker state
	status.CircuitState = m.circuitBreaker.GetState().String()

	return status, nil
}

// GetCircuitBreakerMetrics returns circuit breaker metrics
func (m *Manager) GetCircuitBreakerMetrics() map[string]interface{} {
	return m.circuitBreaker.GetMetrics()
}

// ResetCircuitBreaker manually resets the circuit breaker
func (m *Manager) ResetCircuitBreaker() {
	log.Printf("Manually resetting circuit breaker")
	m.circuitBreaker.Reset()
}

