package actor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ActorPool manages multiple notification actors
type ActorPool struct {
	redis        *redis.Client
	db           ActorDB
	actors       map[uuid.UUID]*NotificationActor
	mu           sync.RWMutex
	maxActors    int
	workerTicker *time.Ticker
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// NewActorPool creates a new actor pool
func NewActorPool(redis *redis.Client, db ActorDB, maxActors int) *ActorPool {
	return &ActorPool{
		redis:     redis,
		db:        db,
		actors:    make(map[uuid.UUID]*NotificationActor),
		maxActors: maxActors,
		stopChan:  make(chan struct{}),
	}
}

// Start starts the actor pool and begins processing
func (p *ActorPool) Start(ctx context.Context) {
	log.Printf("Starting Actor Pool (max actors: %d)", p.maxActors)

	// Start worker that polls for retry-ready notifications
	p.workerTicker = time.NewTicker(10 * time.Second)
	p.wg.Add(1)
	go p.worker(ctx)
}

// Stop stops the actor pool
func (p *ActorPool) Stop() {
	log.Println("Stopping Actor Pool")
	close(p.stopChan)
	if p.workerTicker != nil {
		p.workerTicker.Stop()
	}
	p.wg.Wait()

	// Stop all active actors
	p.mu.Lock()
	for _, actor := range p.actors {
		actor.Stop()
	}
	p.mu.Unlock()

	log.Println("Actor Pool stopped")
}

// SpawnActor spawns a new actor for a notification destination
func (p *ActorPool) SpawnActor(
	ctx context.Context,
	notificationDestID uuid.UUID,
	notificationID uuid.UUID,
	conversationID string,
	message string,
	installation *database.BotInstallation,
	botAppID string,
	botAppPassword string,
) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if actor already exists
	if _, exists := p.actors[notificationDestID]; exists {
		log.Printf("Actor already exists for notification_dest=%s", notificationDestID)
		return nil
	}

	// Check pool capacity
	if len(p.actors) >= p.maxActors {
		log.Printf("Actor pool at capacity (%d/%d), queuing for later", len(p.actors), p.maxActors)
		return nil
	}

	// Create and start actor
	actor := NewNotificationActor(
		notificationDestID,
		notificationID,
		conversationID,
		message,
		installation,
		botAppID,
		botAppPassword,
		p.redis,
		p.db,
	)

	p.actors[notificationDestID] = actor
	actor.Start(ctx)

	// Monitor actor completion
	go func() {
		result := actor.Wait()
		p.handleActorCompletion(result)
	}()

	log.Printf("Spawned actor %s for notification_dest=%s (pool: %d/%d)",
		actor.ID, notificationDestID, len(p.actors), p.maxActors)

	return nil
}

// handleActorCompletion handles actor completion
func (p *ActorPool) handleActorCompletion(result *ActorResult) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.actors, result.NotificationDestID)

	if result.Success {
		log.Printf("Actor completed successfully: %s (notification_dest=%s)",
			result.ActorID, result.NotificationDestID)
	} else if result.Retrying {
		log.Printf("Actor scheduled retry: %s (notification_dest=%s, next_retry=%v)",
			result.ActorID, result.NotificationDestID, result.NextRetryAt)
	} else {
		log.Printf("Actor failed permanently: %s (notification_dest=%s, error=%v)",
			result.ActorID, result.NotificationDestID, result.Error)
	}
}

// worker periodically checks for retry-ready notifications
func (p *ActorPool) worker(ctx context.Context) {
	defer p.wg.Done()

	log.Println("Actor Pool worker started")

	for {
		select {
		case <-p.stopChan:
			log.Println("Actor Pool worker stopped")
			return
		case <-ctx.Done():
			log.Println("Actor Pool worker context cancelled")
			return
		case <-p.workerTicker.C:
			p.processRetryReady(ctx)
		}
	}
}

// processRetryReady processes retry-ready notifications
func (p *ActorPool) processRetryReady(ctx context.Context) {
	p.mu.RLock()
	currentActors := len(p.actors)
	p.mu.RUnlock()

	if currentActors >= p.maxActors {
		log.Printf("Pool at capacity (%d/%d), skipping retry processing", currentActors, p.maxActors)
		return
	}

	// This would query the database for retry-ready notifications
	// For now, we'll leave this as a placeholder
	// In a real implementation, this would:
	// 1. Query notification_destinations WHERE status='failed' AND next_retry_at <= NOW()
	// 2. Spawn actors for each one (up to pool capacity)

	log.Printf("Checking for retry-ready notifications (current actors: %d/%d)", currentActors, p.maxActors)
}

// GetPoolStatus returns the current pool status
func (p *ActorPool) GetPoolStatus() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	activeActors := make([]string, 0, len(p.actors))
	for destID, actor := range p.actors {
		activeActors = append(activeActors, fmt.Sprintf("%s:%s", destID, actor.ID))
	}

	return map[string]interface{}{
		"max_actors":      p.maxActors,
		"active_actors":   len(p.actors),
		"available_slots": p.maxActors - len(p.actors),
		"actors":          activeActors,
	}
}
