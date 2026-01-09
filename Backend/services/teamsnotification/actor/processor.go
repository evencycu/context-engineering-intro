package actor

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// NotificationProcessor unifies all actor components into a single system
type NotificationProcessor struct {
	pool     *WorkerPool
	consumer *QueueConsumer
	producer *QueueProducer
}

// NewNotificationProcessor initializes the entire notification actor system
func NewNotificationProcessor(
	redisClient *redis.Client,
	db WorkerDB,
	tokenManager TokenManager,
	billingService BillingService,
	maxWorkers int,
	taskRepo TaskRepository,
) *NotificationProcessor {
	// 1. Infrastructure
	taskQueue := NewTaskQueue(redisClient)
	circuitBreaker := NewRedisCircuitBreaker(redisClient)
	
	// 2. Execution Layer
	teamsConnector := NewTeamsConnector(tokenManager, circuitBreaker)
	
	// 3. Orchestration Layer
	workerPool := NewWorkerPool(db, teamsConnector, billingService, taskQueue, maxWorkers)
	queueConsumer := NewQueueConsumer(taskQueue, workerPool)
	queueProducer := NewQueueProducer(taskRepo, taskQueue, 500*time.Millisecond)

	return &NotificationProcessor{
		pool:     workerPool,
		consumer: queueConsumer,
		producer: queueProducer,
	}
}

// Start launches all components
func (p *NotificationProcessor) Start(ctx context.Context) {
	log.Println("NotificationProcessor: starting system...")
	p.pool.Start(ctx)
	p.consumer.Start(ctx)
	p.producer.Start(ctx)
}

// Stop gracefully shuts down the system
func (p *NotificationProcessor) Stop() {
	log.Println("NotificationProcessor: stopping system...")
	p.producer.Stop()
	p.consumer.Stop()
	p.pool.Stop()
}

// GetStatus returns the operational status of the system
func (p *NotificationProcessor) GetStatus() map[string]any {
	return map[string]any{
		"pool": p.pool.GetStatus(),
	}
}