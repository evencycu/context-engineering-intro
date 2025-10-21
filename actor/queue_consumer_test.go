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

func (m *mockPool) SpawnActor(ctx context.Context, id uuid.UUID) bool {
	atomic.AddInt32(&m.calls, 1)
	return true // Always succeed for testing
}

func (m *mockPool) HasCapacity() bool {
	return true // Always has capacity for testing
}

func TestQueueConsumer_Basic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("redis not available: ", err)
	}
	// Clean keys
	rdb.Del(ctx, "queue:notifications:high", "queue:notifications:normal", "queue:notifications:low")

	q := NewRedisQueue(rdb)
	pool := &mockPool{}
	consumer := NewQueueConsumer(q, pool)
	consumer.Start(ctx)
	defer consumer.Stop()

	// Enqueue three ids with different priorities
	idHigh := uuid.New()
	idNormal := uuid.New()
	idLow := uuid.New()
	_ = q.Enqueue(ctx, idLow, PriorityLow)
	_ = q.Enqueue(ctx, idNormal, PriorityNormal)
	_ = q.Enqueue(ctx, idHigh, PriorityHigh)

	// Give consumer time to process
	time.Sleep(300 * time.Millisecond)

	if atomic.LoadInt32(&pool.calls) < 1 {
		t.Fatalf("expected at least 1 spawn call")
	}
}
