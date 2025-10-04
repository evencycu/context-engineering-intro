package actor

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ActorPool manages multiple notification actors for the async architecture
type ActorPool struct {
	redis        *redis.Client
	db           ActorDB
	actors       map[uuid.UUID]*NotificationActor
	mu           sync.RWMutex
	maxActors    int
	pollInterval time.Duration
	workerTicker *time.Ticker
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// NewActorPool creates a new actor pool for the async architecture
// pollInterval controls how often the pool polls for retry-ready notifications
func NewActorPool(redis *redis.Client, db ActorDB, maxActors int, pollInterval time.Duration) *ActorPool {
	return &ActorPool{
		redis:        redis,
		db:           db,
		actors:       make(map[uuid.UUID]*NotificationActor),
		maxActors:    maxActors,
		pollInterval: pollInterval,
		stopChan:     make(chan struct{}),
	}
}

// Start begins the actor pool's operation
func (p *ActorPool) Start(ctx context.Context) {
	log.Printf("Starting Actor Pool (max actors: %d)", p.maxActors)

	// Start worker that polls for retry-ready notifications
	if p.pollInterval <= 0 {
		p.pollInterval = 100 * time.Millisecond
	}
	p.workerTicker = time.NewTicker(p.pollInterval)
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case <-p.workerTicker.C:
				p.pollAndSpawnActors(ctx)
			case <-p.stopChan:
				log.Println("Actor Pool worker stopped")
				return
			}
		}
	}()

	log.Println("Actor Pool started successfully")
}

// Stop gracefully shuts down the actor pool
func (p *ActorPool) Stop() {
	log.Println("Stopping Actor Pool...")
	close(p.stopChan)
	p.workerTicker.Stop()
	p.wg.Wait()

	// Stop all active actors
	p.mu.Lock()
	for _, actor := range p.actors {
		actor.Stop()
	}
	p.mu.Unlock()

	log.Println("Actor Pool stopped")
}

// spawnActor creates and starts a new NotificationActor
func (p *ActorPool) spawnActor(ctx context.Context, notificationDestID uuid.UUID) {
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

	// Get TeamsBot to get AppID
	teamsBot, err := p.db.GetTeamsBot(ctx, installation.BotID)
	if err != nil {
		log.Printf("Failed to get teams bot %s: %v", installation.BotID, err)
		return
	}

	// Get password from environment variables
	appPassword := os.Getenv("TEAMS_BOT_APP_PASSWORD")
	if appPassword == "" {
		log.Printf("TEAMS_BOT_APP_PASSWORD environment variable not set")
		return
	}

	// Create actor with required parameters
	actor := NewNotificationActor(
		notificationDestID,
		nd.NotificationID,
		*nd.ConversationID,
		"", // Message will be fetched by actor
		installation,
		teamsBot.AppID, // Get from TeamsBot
		appPassword,    // Get from environment
		p.redis,
		p.db,
	)

	// Override tenant ID from TeamsBot if available
	if teamsBot.TenantID != nil && *teamsBot.TenantID != "" {
		actor.TenantID = *teamsBot.TenantID
	}

	p.actors[notificationDestID] = actor
	actor.Start(ctx)

	log.Printf("Spawned actor %s for notification_dest=%s (pool: %d/%d)",
		actor.ID, notificationDestID, len(p.actors), p.maxActors)
}

// pollAndSpawnActors checks for retry-ready notifications and spawns actors
func (p *ActorPool) pollAndSpawnActors(ctx context.Context) {
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
func (p *ActorPool) GetPoolStatus() map[string]interface{} {
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
