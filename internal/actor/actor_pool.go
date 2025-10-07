package actor

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ActorPool manages multiple notification actors for the async architecture
type ActorPool struct {
	redis        *redis.Client
	db           ActorDB
	tokenManager TokenManager // New field for token management
	actors       map[uuid.UUID]*NotificationActor
	mu           sync.RWMutex
	maxActors    int
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// NewActorPool creates a new actor pool for the async architecture
// Tasks are now provided by QueueConsumer instead of polling DB
func NewActorPool(redis *redis.Client, db ActorDB, tokenManager TokenManager, maxActors int) *ActorPool {
	return &ActorPool{
		redis:        redis,
		db:           db,
		tokenManager: tokenManager,
		actors:       make(map[uuid.UUID]*NotificationActor),
		maxActors:    maxActors,
		stopChan:     make(chan struct{}),
	}
}

// Start begins the actor pool's operation
// Now just initializes the pool - tasks come from QueueConsumer via SpawnActor
func (p *ActorPool) Start(ctx context.Context) {
	log.Printf("Starting Actor Pool (max actors: %d)", p.maxActors)
	log.Println("Actor Pool ready to receive tasks from QueueConsumer")
}

// Stop gracefully shuts down the actor pool
func (p *ActorPool) Stop() {
	log.Println("Stopping Actor Pool...")
	close(p.stopChan)
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
		p.tokenManager, // Pass Token Manager
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

// SpawnActor spawns a new actor for the given notification destination ID
// Returns true if actor was successfully spawned, false if pool is full
// This method is now called by QueueConsumer instead of polling DB
func (p *ActorPool) SpawnActor(ctx context.Context, ndID uuid.UUID) bool {
	p.mu.RLock()
	currentActors := len(p.actors)
	p.mu.RUnlock()

	if currentActors >= p.maxActors {
		log.Printf("Actor pool is full (%d/%d), cannot spawn actor for %s", currentActors, p.maxActors, ndID)
		return false
	}

	log.Printf("Spawning actor for notification_dest=%s (pool: %d/%d)", ndID, currentActors, p.maxActors)
	p.spawnActor(ctx, ndID)
	return true
}

// HasCapacity checks if the actor pool has available capacity
func (p *ActorPool) HasCapacity() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.actors) < p.maxActors
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
