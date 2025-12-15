package actor

import (
	"context"
	"time"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/repositories"
	"github.com/google/uuid"
)

// ActorDB defines the models operations needed by actors
type ActorDB interface {
	UpdateNotificationDestination(ctx context.Context, id uuid.UUID, update *NotificationDestinationUpdate) error
	GetBotInstallation(ctx context.Context, conversationID string) (*models.BotInstallation, error)
	GetNotificationDestinationByID(ctx context.Context, id uuid.UUID) (*models.NotificationDestination, error)
	GetRetryReadyNotificationDestinations(ctx context.Context, limit int) ([]*models.NotificationDestination, error)
	GetTeamsBot(ctx context.Context, botID uuid.UUID) (*models.TeamsBot, error)
	GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
}

// actorDB implements ActorDB interface for actors
type actorDB struct {
	notificationDestRepo repositories.NotificationDestinationRepository
	botInstallationRepo  repositories.BotInstallationRepository
	teamsBotRepo         repositories.TeamsBotRepository
	notificationRepo     repositories.NotificationRepository
	projectRepo          repositories.ProjectRepository
}

// NewActorDB creates a new actor DB implementation
func NewActorDB(notificationDestRepo repositories.NotificationDestinationRepository, botInstallationRepo repositories.BotInstallationRepository, teamsBotRepo repositories.TeamsBotRepository, notificationRepo repositories.NotificationRepository, projectRepo repositories.ProjectRepository) ActorDB {
	return &actorDB{
		notificationDestRepo: notificationDestRepo,
		botInstallationRepo:  botInstallationRepo,
		teamsBotRepo:         teamsBotRepo,
		notificationRepo:     notificationRepo,
		projectRepo:          projectRepo,
	}
}

// UpdateNotificationDestination updates a notification destination
func (db *actorDB) UpdateNotificationDestination(ctx context.Context, id uuid.UUID, update *NotificationDestinationUpdate) error {
	// Get current record
	entity, err := db.notificationDestRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entity == nil {
		return nil // Not found
	}

	if update.Status != nil {
		entity.Status = *update.Status
	}
	if update.ErrorMessage != nil {
		em := *update.ErrorMessage
		entity.ErrorMessage = &em
	}
	if update.TeamsMessageID != nil {
		tm := *update.TeamsMessageID
		entity.TeamsMessageID = &tm
	}
	if update.SentAt != nil {
		entity.SentAt = update.SentAt
	}
	if update.RetryCount != nil {
		entity.RetryCount = *update.RetryCount
	}
	if update.NextRetryAt != nil {
		entity.NextRetryAt = update.NextRetryAt
	}
	if update.FirstAttemptAt != nil {
		entity.FirstAttemptAt = update.FirstAttemptAt
	}
	if update.LastAttemptAt != nil {
		entity.LastAttemptAt = update.LastAttemptAt
	}
	if update.FailureReason != nil {
		fr := *update.FailureReason
		entity.FailureReason = &fr
	}
	if update.RetryAfter != nil {
		entity.RetryAfter = update.RetryAfter
	}
	if update.ActorID != nil {
		aid := *update.ActorID
		entity.ActorID = &aid
	}

	entity.UpdatedAt = time.Now()

	return db.notificationDestRepo.Update(ctx, entity)
}

// GetBotInstallation gets a bot installation by conversation ID
func (db *actorDB) GetBotInstallation(ctx context.Context, conversationID string) (*models.BotInstallation, error) {
	return db.botInstallationRepo.GetByConversationID(ctx, conversationID)
}

// GetNotificationDestinationByID gets a notification destination by ID
func (db *actorDB) GetNotificationDestinationByID(ctx context.Context, id uuid.UUID) (*models.NotificationDestination, error) {
	return db.notificationDestRepo.GetByID(ctx, id)
}

// GetRetryReadyNotificationDestinations gets retry-ready notification destinations
func (db *actorDB) GetRetryReadyNotificationDestinations(ctx context.Context, limit int) ([]*models.NotificationDestination, error) {
	return db.notificationDestRepo.GetRetryReady(ctx, limit)
}

// GetTeamsBot gets a teams bot by ID
func (db *actorDB) GetTeamsBot(ctx context.Context, botID uuid.UUID) (*models.TeamsBot, error) {
	return db.teamsBotRepo.GetByID(ctx, botID)
}

// GetNotificationByID gets a notification by ID
func (db *actorDB) GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	return db.notificationRepo.GetByID(ctx, id)
}

// GetProjectByID gets a project by ID
func (db *actorDB) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return db.projectRepo.GetByID(ctx, id)
}
