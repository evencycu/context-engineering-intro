package actor

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NotificationProcessor manages the complete notification processing pipeline
// It combines ActorPool, QueueConsumer, and EnqueueWorker into a single, cohesive unit
type NotificationProcessor struct {
	pool     *ActorPool
	consumer *QueueConsumer
	enqueuer *EnqueueWorker
}

// NewNotificationProcessor creates a new notification processor
// This is the unified entry point for the notification processing system
func NewNotificationProcessor(
	redis *redis.Client,
	queue RedisQueue,
	db ActorDB,
	tokenManager TokenManager,
	maxActors int,
	ndRepo NotificationDestinationRepository,
) *NotificationProcessor {
	// Create actor pool first
	pool := NewActorPool(redis, db, tokenManager, maxActors)

	// Create queue consumer with the pool
	consumer := NewQueueConsumer(queue, pool)

	// Create enqueue worker: pending ND -> enqueue -> mark enqueued
	enqueuer := NewEnqueueWorker(ndRepo, queue, 200*time.Millisecond)

	return &NotificationProcessor{
		pool:     pool,
		consumer: consumer,
		enqueuer: enqueuer,
	}
}

// Start initializes and starts the notification processor
func (np *NotificationProcessor) Start(ctx context.Context) {
	log.Printf("Starting Notification Processor (max actors: %d)", np.pool.maxActors)

	// Start actor pool
	np.pool.Start(ctx)

	// Start queue consumer
	np.consumer.Start(ctx)

	// Start enqueue worker
	np.enqueuer.Start(ctx)

	log.Println("Notification Processor started successfully")
}

// Stop gracefully shuts down the notification processor
func (np *NotificationProcessor) Stop() {
	log.Println("Stopping Notification Processor...")

	// Stop enqueuer
	np.enqueuer.Stop()

	// Stop consumer
	np.consumer.Stop()

	// Then stop actor pool
	np.pool.Stop()

	log.Println("Notification Processor stopped")
}

// GetStatus returns the current status of the notification processor
func (np *NotificationProcessor) GetStatus() map[string]interface{} {
	poolStatus := np.pool.GetPoolStatus()

	return map[string]interface{}{
		"processor": "notification",
		"status":    "running",
		"pool":      poolStatus,
	}
}
