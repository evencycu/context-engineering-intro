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

// TaskQueue provides enqueue/dequeue operations for notification tasks
type TaskQueue interface {
	Enqueue(ctx context.Context, id uuid.UUID, priority Priority) error
	DequeueBlocking(ctx context.Context) (uuid.UUID, error)
	Requeue(ctx context.Context, id uuid.UUID) error
	RequeueToFront(ctx context.Context, id uuid.UUID) error
}

type redisTaskQueue struct {
	redis *redis.Client
}

// NewTaskQueue creates a new Redis-backed task queue
func NewTaskQueue(redis *redis.Client) TaskQueue {
	return &redisTaskQueue{redis: redis}
}

func (q *redisTaskQueue) keyForPriority(p Priority) string {
	return fmt.Sprintf("queue:notifications:%s", string(p))
}

func (q *redisTaskQueue) Enqueue(ctx context.Context, id uuid.UUID, priority Priority) error {
	dedupKey := fmt.Sprintf("queue:notifications:dedup:%s", id)
	set, err := q.redis.SetNX(ctx, dedupKey, 1, 60*time.Second).Result()
	if err == nil && !set {
		return nil
	}
	return q.redis.RPush(ctx, q.keyForPriority(priority), id.String()).Err()
}

func (q *redisTaskQueue) DequeueBlocking(ctx context.Context) (uuid.UUID, error) {
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
	id, err := uuid.Parse(res[1])
	if err != nil {
		return uuid.Nil, err
	}

	processingKey := fmt.Sprintf("queue:notifications:processing:%s", id)
	ok, err := q.redis.SetNX(ctx, processingKey, 1, 30*time.Minute).Result()
	if err != nil || !ok {
		return uuid.Nil, nil
	}
	return id, nil
}

func (q *redisTaskQueue) Requeue(ctx context.Context, id uuid.UUID) error {
	processingKey := fmt.Sprintf("queue:notifications:processing:%s", id)
	q.redis.Del(ctx, processingKey)
	return q.redis.RPush(ctx, q.keyForPriority(PriorityNormal), id.String()).Err()
}

func (q *redisTaskQueue) RequeueToFront(ctx context.Context, id uuid.UUID) error {
	processingKey := fmt.Sprintf("queue:notifications:processing:%s", id)
	q.redis.Del(ctx, processingKey)
	return q.redis.LPush(ctx, q.keyForPriority(PriorityHigh), id.String()).Err()
}