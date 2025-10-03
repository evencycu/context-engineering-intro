package actor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ActorPoolV2 manages multiple notification actors for the new async architecture
type ActorPoolV2 struct {
	redis        *redis.Client
	db           ActorDB
	actors       map[uuid.UUID]*NotificationActor
	mu           sync.RWMutex
	maxActors    int
	workerTicker *time.Ticker
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// NewActorPoolV2 creates a new actor pool for the async architecture
func NewActorPoolV2(redis *redis.Client, db ActorDB, maxActors int) *ActorPoolV2 {
	return &ActorPoolV2{
		redis:     redis,
		db:        db,
		actors:    make(map[uuid.UUID]*NotificationActor),
		maxActors: maxActors,
		stopChan:  make(chan struct{}),
	}
}

// Start begins the actor pool's operation
func (p *ActorPoolV2) Start(ctx context.Context) {
	log.Printf("Starting Actor Pool V2 (max actors: %d)", p.maxActors)

	// Start worker that polls for retry-ready notifications
	p.workerTicker = time.NewTicker(10 * time.Second)
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case <-p.workerTicker.C:
				p.pollAndSpawnActors(ctx)
			case <-p.stopChan:
				log.Println("Actor Pool V2 worker stopped")
				return
			}
		}
	}()

	log.Println("Actor Pool V2 started successfully")
}

// Stop gracefully shuts down the actor pool
func (p *ActorPoolV2) Stop() {
	log.Println("Stopping Actor Pool V2...")
	close(p.stopChan)
	p.workerTicker.Stop()
	p.wg.Wait()

	// Stop all active actors
	p.mu.Lock()
	for _, actor := range p.actors {
		actor.Stop()
	}
	p.mu.Unlock()

	log.Println("Actor Pool V2 stopped")
}

// EnqueueNotificationDestination adds a new notification destination to be processed by an actor
func (p *ActorPoolV2) EnqueueNotificationDestination(ctx context.Context, notificationDestID uuid.UUID) error {
	p.mu.RLock()
	currentActors := len(p.actors)
	p.mu.RUnlock()

	if currentActors >= p.maxActors {
		return fmt.Errorf("actor pool is full, cannot enqueue notification_dest %s", notificationDestID)
	}

	// Spawn an actor immediately for new notifications
	p.spawnActor(ctx, notificationDestID)
	return nil
}

// spawnActor creates and starts a new NotificationActor
func (p *ActorPoolV2) spawnActor(ctx context.Context, notificationDestID uuid.UUID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.actors[notificationDestID]; exists {
		log.Printf("Actor for notification_dest %s already exists, skipping spawn", notificationDestID)
		return
	}

	// Get notification destination details from database
	nd, err := p.db.GetNotificationDestinationByID(ctx, notificationDestID)
	if err != nil {
		log.Printf("Failed to get notification destination %s: %v", notificationDestID, err)
		return
	}

	// Get bot installation for this conversation
	if nd.ConversationID == nil {
		log.Printf("Notification destination %s has no conversation ID", notificationDestID)
		return
	}
	
	installation, err := p.db.GetBotInstallation(ctx, *nd.ConversationID)
	if err != nil {
		log.Printf("Failed to get bot installation for conversation %s: %v", *nd.ConversationID, err)
		return
	}

	// Create actor with required parameters
	actor := NewNotificationActor(
		notificationDestID,
		nd.NotificationID,
		*nd.ConversationID,
		"", // Message will be fetched by actor
		installation,
		"", // BotAppID will be set by actor
		"", // BotAppPassword will be set by actor
		p.redis,
		p.db,
	)

	p.actors[notificationDestID] = actor
	actor.Start(ctx)

	log.Printf("Spawned actor %s for notification_dest=%s (pool: %d/%d)",
		actor.ID, notificationDestID, len(p.actors), p.maxActors)
}

// pollAndSpawnActors checks for retry-ready notifications and spawns actors
func (p *ActorPoolV2) pollAndSpawnActors(ctx context.Context) {
	p.mu.RLock()
	currentActors := len(p.actors)
	p.mu.RUnlock()

	if currentActors >= p.maxActors {
		log.Printf("Actor pool is full (%d/%d), skipping poll for new retries", currentActors, p.maxActors)
		return
	}

	// Query notification_destinations WHERE status='pending' OR status='failed' AND next_retry_at <= NOW()
	limit := p.maxActors - currentActors
	if limit <= 0 {
		return
	}

	log.Printf("Checking for retry-ready notifications (current actors: %d/%d, limit: %d)", currentActors, p.maxActors, limit)

	retryReady, err := p.db.GetRetryReadyNotificationDestinations(ctx, limit)
	if err != nil {
		log.Printf("Failed to query retry-ready notifications: %v", err)
		return
	}

	for _, nd := range retryReady {
		p.spawnActor(ctx, nd.ID)
	}
}

// GetPoolStatus returns current status of the actor pool
func (p *ActorPoolV2) GetPoolStatus() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	activeActors := make(map[uuid.UUID]string)
	for id, actor := range p.actors {
		activeActors[id] = actor.NotificationDestID.String()
	}

	return map[string]interface{}{
		"max_actors":      p.maxActors,
		"active_actors":   len(p.actors),
		"available_slots": p.maxActors - len(p.actors),
		"actors":          activeActors,
	}
}
