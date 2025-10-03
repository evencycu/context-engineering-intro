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

	q := NewRedisQueue(rdb)

	id := uuid.New()
	if err := q.Enqueue(ctx, id, PriorityHigh); err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	// Use separate context for blocking pop
	popCtx, cancelPop := context.WithTimeout(context.Background(), time.Second)
	defer cancelPop()
	got, err := q.DequeueBlocking(popCtx)
	if err != nil {
		t.Fatalf("dequeue failed: %v", err)
	}
	if got != id {
		t.Fatalf("want %s, got %s", id, got)
	}
}
