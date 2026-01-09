package actor

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func TestTaskQueue_EnqueueDequeue(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("redis not available")
	}
	rdb.Del(ctx, "queue:notifications:high", "queue:notifications:normal", "queue:notifications:low")

	q := NewTaskQueue(rdb)
	id := uuid.New()
	
	if err := q.Enqueue(ctx, id, PriorityHigh); err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	got, err := q.DequeueBlocking(ctx)
	if err != nil {
		t.Fatalf("dequeue failed: %v", err)
	}
	if got != id {
		t.Fatalf("want %s, got %s", id, got)
	}
}