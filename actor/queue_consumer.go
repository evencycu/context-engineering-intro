package actor

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

// QueueConsumer continuously consumes from RedisQueue and spawns actors
// for each dequeued notification destination ID.
// This is now the ONLY source of tasks for ActorPool.
type spawnPool interface {
	SpawnActor(ctx context.Context, id uuid.UUID) bool
	HasCapacity() bool
}

type QueueConsumer struct {
	queue    RedisQueue
	pool     spawnPool
	cancelFn context.CancelFunc
}

func NewQueueConsumer(queue RedisQueue, pool spawnPool) *QueueConsumer {
	return &QueueConsumer{queue: queue, pool: pool}
}

func (c *QueueConsumer) Start(ctx context.Context) {
	runCtx, cancel := context.WithCancel(ctx)
	c.cancelFn = cancel
	go c.loop(runCtx)
}

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

		// Check if pool has capacity before dequeuing
		if !c.pool.HasCapacity() {
			// Pool is full, wait a bit before checking again
			// This prevents dequeuing tasks that can't be processed
			log.Printf("Actor pool full, waiting before dequeuing...")
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				continue
			}
		}

		ndID, err := c.queue.DequeueBlocking(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("QueueConsumer dequeue error: %v", err)
			continue
		}
		if ndID == uuid.Nil {
			continue
		}

		// Double-check capacity before spawning (race condition protection)
		success := c.pool.SpawnActor(ctx, ndID)
		if !success {
			// Pool became full between check and spawn
			// Re-queue the task to the front of the queue (high priority)
			log.Printf("Actor pool became full, re-queuing task %s to front", ndID)
			if requeueErr := c.queue.RequeueToFront(ctx, ndID); requeueErr != nil {
				log.Printf("Failed to re-queue task %s: %v", ndID, requeueErr)
			}
		}
	}
}
