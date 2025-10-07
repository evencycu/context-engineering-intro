package actor

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestRedisQueue_EnqueueDequeue(t *testing.T) {
	// This test requires a local Redis at redis://localhost:6379
	// Skip if not available quickly
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("redis not available: ", err)
	}

	// Clean up any existing keys
	rdb.Del(ctx, "queue:notifications:high", "queue:notifications:normal", "queue:notifications:low")

	q := NewRedisQueue(rdb)

	id := uuid.New()
	if err := q.Enqueue(ctx, id, PriorityHigh); err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	// Use separate context for blocking pop with shorter timeout
	popCtx, cancelPop := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelPop()
	got, err := q.DequeueBlocking(popCtx)
	if err != nil {
		if ctx.Err() != nil {
			t.Skip("redis connection lost during test: ", err)
		}
		t.Fatalf("dequeue failed: %v", err)
	}
	if got != id {
		t.Fatalf("want %s, got %s", id, got)
	}
}
