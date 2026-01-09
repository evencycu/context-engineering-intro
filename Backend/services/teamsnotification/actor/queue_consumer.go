package actor

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

// WorkerPoolInterface defines the methods needed by QueueConsumer to process tasks
type WorkerPoolInterface interface {
	SpawnWorker(ctx context.Context, id uuid.UUID) bool
	HasCapacity() bool
}

// QueueConsumer listens to the task queue and requests workers from the pool
type QueueConsumer struct {
	queue    TaskQueue
	pool     WorkerPoolInterface
	cancelFn context.CancelFunc
}

// NewQueueConsumer creates a new queue consumer
func NewQueueConsumer(queue TaskQueue, pool WorkerPoolInterface) *QueueConsumer {
	return &QueueConsumer{
		queue: queue,
		pool:  pool,
	}
}

// Start begins the consumer's loop
func (c *QueueConsumer) Start(ctx context.Context) {
	runCtx, cancel := context.WithCancel(ctx)
	c.cancelFn = cancel
	go c.loop(runCtx)
}

// Stop stops the consumer
func (c *QueueConsumer) Stop() {
	if c.cancelFn != nil {
		c.cancelFn()
	}
}

func (c *QueueConsumer) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if !c.pool.HasCapacity() {
			log.Printf("QueueConsumer: worker pool is full, waiting...")
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				continue
			}
		}

		id, err := c.queue.DequeueBlocking(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("QueueConsumer: dequeue error: %v", err)
			continue
		}

		if id == uuid.Nil {
			continue
		}

		// Try to spawn worker
		if !c.pool.SpawnWorker(ctx, id) {
			log.Printf("QueueConsumer: failed to spawn worker for %s, requeueing...", id)
			_ = c.queue.RequeueToFront(ctx, id)
		}
	}
}