package actor

import (
	"context"
	"log"
	"time"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/google/uuid"
)

// TaskRepository subset needed by the producer
type TaskRepository interface {
	GetPending(ctx context.Context, limit int) ([]*models.NotificationDestination, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

// QueueProducer scans the database for pending notifications and pushes them to the task queue
type QueueProducer struct {
	repo     TaskRepository
	queue    TaskQueue
	interval time.Duration
	cancelFn context.CancelFunc
}

// NewQueueProducer creates a new queue producer
func NewQueueProducer(repo TaskRepository, queue TaskQueue, interval time.Duration) *QueueProducer {
	return &QueueProducer{
		repo:     repo,
		queue:    queue,
		interval: interval,
	}
}

// Start begins the producer's scan loop
func (p *QueueProducer) Start(ctx context.Context) {
	runCtx, cancel := context.WithCancel(ctx)
	p.cancelFn = cancel
	go p.loop(runCtx)
}

// Stop stops the producer
func (p *QueueProducer) Stop() {
	if p.cancelFn != nil {
		p.cancelFn()
	}
}

func (p *QueueProducer) loop(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.scanOnce(ctx)
		}
	}
}

func (p *QueueProducer) scanOnce(ctx context.Context) {
	rows, err := p.repo.GetPending(ctx, 100)
	if err != nil {
		log.Printf("QueueProducer: failed to get pending tasks: %v", err)
		return
	}

	for _, row := range rows {
		prio := PriorityNormal
		switch row.Priority {
		case "high":
			prio = PriorityHigh
		case "low":
			prio = PriorityLow
		}

		if err := p.queue.Enqueue(ctx, row.ID, prio); err != nil {
			log.Printf("QueueProducer: failed to enqueue task %s: %v", row.ID, err)
			continue
		}

		if err := p.repo.UpdateStatus(ctx, row.ID, string(models.NotificationStatusEnqueued)); err != nil {
			log.Printf("QueueProducer: failed to update status for %s: %v", row.ID, err)
		}
	}
}