package actor

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
)

// WorkerPool manages the execution lifecycle of notification workers
type WorkerPool struct {
	db             WorkerDB
	connector      TeamsConnector
	billingService BillingService
	taskQueue      TaskQueue
	workers        map[uuid.UUID]*NotificationWorker
	mu             sync.RWMutex
	maxWorkers     int
	stopChan       chan struct{}
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(db WorkerDB, connector TeamsConnector, billingService BillingService, taskQueue TaskQueue, maxWorkers int) *WorkerPool {
	return &WorkerPool{
		db:             db,
		connector:      connector,
		billingService: billingService,
		taskQueue:      taskQueue,
		workers:        make(map[uuid.UUID]*NotificationWorker),
		maxWorkers:     maxWorkers,
		stopChan:       make(chan struct{}),
	}
}

// Start initializes the pool
func (p *WorkerPool) Start(ctx context.Context) {
	log.Printf("WorkerPool: started with capacity=%d", p.maxWorkers)
}

// Stop shuts down all workers
func (p *WorkerPool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, w := range p.workers {
		w.Stop()
	}
	log.Println("WorkerPool: stopped")
}

// SpawnWorker creates and starts a new worker for a task
func (p *WorkerPool) SpawnWorker(ctx context.Context, id uuid.UUID) bool {
	p.mu.Lock()
	if len(p.workers) >= p.maxWorkers {
		p.mu.Unlock()
		return false
	}
	
	if _, exists := p.workers[id]; exists {
		p.mu.Unlock()
		return true
	}

	// Fetch task details
	nd, err := p.db.GetNotificationDestinationByID(ctx, id)
	if err != nil || nd == nil {
		p.mu.Unlock()
		log.Printf("WorkerPool: task %s not found: %v", id, err)
		return true // Consider handled if not found
	}

	installation, err := p.db.GetBotInstallation(ctx, *nd.ConversationID)
	if err != nil {
		p.mu.Unlock()
		log.Printf("WorkerPool: installation for %s not found", id)
		return true
	}

	worker := NewNotificationWorker(
		id,
		nd.NotificationID,
		*nd.ConversationID,
		installation,
		p.db,
		p.connector,
		p.billingService,
		p.taskQueue,
	)
	worker.Priority = nd.Priority

	p.workers[id] = worker
	p.mu.Unlock()

	worker.Start(ctx)
	
	go func() {
		_ = worker.Wait()
		p.mu.Lock()
		delete(p.workers, id)
		p.mu.Unlock()
	}()

	return true
}

func (p *WorkerPool) HasCapacity() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.workers) < p.maxWorkers
}

func (p *WorkerPool) GetStatus() map[string]any {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return map[string]any{
		"max_workers":    p.maxWorkers,
		"active_workers": len(p.workers),
	}
}