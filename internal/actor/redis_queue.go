package actor

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Priority is a string alias for queue priority levels
type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityNormal Priority = "normal"
	PriorityLow    Priority = "low"
)

// RedisQueue provides enqueue/dequeue operations backed by Redis Lists
// with SETNX used to avoid duplicate processing.
type RedisQueue interface {
	Enqueue(ctx context.Context, ndID uuid.UUID, priority Priority) error
	DequeueBlocking(ctx context.Context) (uuid.UUID, error)
}

type redisQueue struct {
	redis *redis.Client
}

func NewRedisQueue(redis *redis.Client) RedisQueue {
	return &redisQueue{redis: redis}
}

func (q *redisQueue) keyForPriority(p Priority) string {
	return fmt.Sprintf("queue:notifications:%s", string(p))
}

// Enqueue pushes the notification destination ID into the priority queue.
// It also sets a short-lived dedup key to reduce duplicate enqueues.
func (q *redisQueue) Enqueue(ctx context.Context, ndID uuid.UUID, priority Priority) error {
	// Dedup on enqueue (best-effort): if recently enqueued, skip.
	dedupKey := fmt.Sprintf("queue:notifications:dedup:%s", ndID)
	set, err := q.redis.SetNX(ctx, dedupKey, 1, 60*time.Second).Result()
	if err == nil && !set {
		// Already enqueued recently, skip silently
		return nil
	}
	// Push to the tail (RPUSH) of the priority list
	return q.redis.RPush(ctx, q.keyForPriority(priority), ndID.String()).Err()
}

// DequeueBlocking pops one item from the highest available priority list.
// It uses BLPOP across [high, normal, low] and returns the ndID.
// A processing key is SETNX to avoid duplicate concurrent processing.
func (q *redisQueue) DequeueBlocking(ctx context.Context) (uuid.UUID, error) {
	// BLPOP returns [list, element]
	res, err := q.redis.BLPop(ctx, 0,
		q.keyForPriority(PriorityHigh),
		q.keyForPriority(PriorityNormal),
		q.keyForPriority(PriorityLow),
	).Result()
	if err != nil {
		return uuid.Nil, err
	}
	if len(res) < 2 {
		return uuid.Nil, fmt.Errorf("invalid BLPOP response")
	}
	ndStr := res[1]
	ndID, err := uuid.Parse(ndStr)
	if err != nil {
		return uuid.Nil, err
	}
	// SETNX a processing key to prevent duplicate processing
	processingKey := fmt.Sprintf("queue:notifications:processing:%s", ndID)
	ok, err := q.redis.SetNX(ctx, processingKey, 1, 30*time.Minute).Result()
	if err != nil {
		return uuid.Nil, err
	}
	if !ok {
		// Already being processed by another consumer; skip
		return uuid.Nil, nil
	}
	return ndID, nil
}
