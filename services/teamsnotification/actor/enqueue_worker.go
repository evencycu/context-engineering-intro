package actor

import (
	"context"
	"log"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/libs/models"
	"github.com/google/uuid"
)

// NotificationDestinationRepository subset needed by worker
type NotificationDestinationRepository interface {
	GetPending(ctx context.Context, limit int) ([]*models.NotificationDestination, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

// EnqueueWorker moves ND from pending -> enqueued by pushing to Redis
type EnqueueWorker struct {
	repo     NotificationDestinationRepository
	queue    RedisQueue
	interval time.Duration
	cancelFn context.CancelFunc
}

func NewEnqueueWorker(repo NotificationDestinationRepository, queue RedisQueue, interval time.Duration) *EnqueueWorker {
	return &EnqueueWorker{repo: repo, queue: queue, interval: interval}
}

func (w *EnqueueWorker) Start(ctx context.Context) {
	runCtx, cancel := context.WithCancel(ctx)
	w.cancelFn = cancel
	go w.loop(runCtx)
}

func (w *EnqueueWorker) Stop() {
	if w.cancelFn != nil {
		w.cancelFn()
	}
}

func (w *EnqueueWorker) loop(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.scanOnce(ctx)
		}
	}
}

func (w *EnqueueWorker) scanOnce(ctx context.Context) {
	// pull small batches to avoid spikes
	rows, err := w.repo.GetPending(ctx, 100)
	if err != nil {
		log.Printf("EnqueueWorker: GetPending error: %v", err)
		return
	}
	for _, nd := range rows {
		// default to normal priority when unknown
		prio := PriorityNormal
		switch nd.Priority {
		case "high":
			prio = PriorityHigh
		case "low":
			prio = PriorityLow
		}
		if err := w.queue.Enqueue(ctx, nd.ID, prio); err != nil {
			log.Printf("EnqueueWorker: enqueue %s failed: %v", nd.ID, err)
			continue // keep pending; no retry here per spec
		}
		if err := w.repo.UpdateStatus(ctx, nd.ID, string(models.NotificationStatusEnqueued)); err != nil {
			log.Printf("EnqueueWorker: update status for %s failed: %v", nd.ID, err)
		}
	}
}
