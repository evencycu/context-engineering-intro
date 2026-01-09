package actor

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type mockPool struct{ calls int32 }

func (m *mockPool) SpawnWorker(ctx context.Context, id uuid.UUID) bool {
	atomic.AddInt32(&m.calls, 1)
	return true
}

func (m *mockPool) HasCapacity() bool {
	return true
}

func TestQueueConsumer_Basic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("redis not available: ", err)
	}
	rdb.Del(ctx, "queue:notifications:high", "queue:notifications:normal", "queue:notifications:low")

	q := NewTaskQueue(rdb)
	pool := &mockPool{}
	consumer := NewQueueConsumer(q, pool)
	consumer.Start(ctx)
	defer consumer.Stop()

	id := uuid.New()
	_ = q.Enqueue(ctx, id, PriorityHigh)

	time.Sleep(300 * time.Millisecond)

	if atomic.LoadInt32(&pool.calls) < 1 {
		t.Fatalf("expected at least 1 spawn call")
	}
}