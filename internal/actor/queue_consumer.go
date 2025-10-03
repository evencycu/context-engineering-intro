package actor

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// QueueConsumer continuously consumes from RedisQueue and spawns actors
// for each dequeued notification destination ID.
type spawnPool interface {
	spawnActor(ctx context.Context, id uuid.UUID)
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

		// Spawn actor for this notification destination
		c.pool.spawnActor(ctx, ndID)
	}
}
