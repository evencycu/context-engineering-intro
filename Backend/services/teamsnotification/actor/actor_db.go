package actor

import (
	"context"
	"time"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/repositories"
	"github.com/google/uuid"
)

// NotificationDestinationUpdate represents fields to update in a notification destination record
type NotificationDestinationUpdate struct {
	Status         *string
	ErrorMessage   *string
	TeamsMessageID *string
	SentAt         *time.Time
	RetryCount     *int
	NextRetryAt    *time.Time
	FirstAttemptAt *time.Time
	LastAttemptAt  *time.Time
	FailureReason  *string
	RetryAfter     *int
	ActorID        *string
}

// workerDB implements WorkerDB interface
type workerDB struct {
	notificationDestRepo repositories.NotificationDestinationRepository
	botInstallationRepo  repositories.BotInstallationRepository
	teamsBotRepo         repositories.TeamsBotRepository
	notificationRepo     repositories.NotificationRepository
	projectRepo          repositories.ProjectRepository
}

// NewWorkerDB creates a new worker DB implementation
func NewWorkerDB(notificationDestRepo repositories.NotificationDestinationRepository, botInstallationRepo repositories.BotInstallationRepository, teamsBotRepo repositories.TeamsBotRepository, notificationRepo repositories.NotificationRepository, projectRepo repositories.ProjectRepository) WorkerDB {
	return &workerDB{
		notificationDestRepo: notificationDestRepo,
		botInstallationRepo:  botInstallationRepo,
		teamsBotRepo:         teamsBotRepo,
		notificationRepo:     notificationRepo,
		projectRepo:          projectRepo,
	}
}

func (db *workerDB) UpdateNotificationDestination(ctx context.Context, id uuid.UUID, update *NotificationDestinationUpdate) error {
	entity, err := db.notificationDestRepo.GetByID(ctx, id)
	if err != nil || entity == nil {
		return err
	}

	if update.Status != nil { entity.Status = *update.Status }
	if update.ErrorMessage != nil { em := *update.ErrorMessage; entity.ErrorMessage = &em }
	if update.TeamsMessageID != nil { tm := *update.TeamsMessageID; entity.TeamsMessageID = &tm }
	if update.SentAt != nil { entity.SentAt = update.SentAt }
	if update.RetryCount != nil { entity.RetryCount = *update.RetryCount }
	if update.NextRetryAt != nil { entity.NextRetryAt = update.NextRetryAt }
	if update.FirstAttemptAt != nil { entity.FirstAttemptAt = update.FirstAttemptAt }
	if update.LastAttemptAt != nil { entity.LastAttemptAt = update.LastAttemptAt }
	if update.FailureReason != nil { fr := *update.FailureReason; entity.FailureReason = &fr }
	if update.RetryAfter != nil { entity.RetryAfter = update.RetryAfter }
	if update.ActorID != nil { aid := *update.ActorID; entity.ActorID = &aid }

	entity.UpdatedAt = time.Now()
	return db.notificationDestRepo.Update(ctx, entity)
}

func (db *workerDB) GetBotInstallation(ctx context.Context, conversationID string) (*models.BotInstallation, error) {
	return db.botInstallationRepo.GetByConversationID(ctx, conversationID)
}

func (db *workerDB) GetNotificationDestinationByID(ctx context.Context, id uuid.UUID) (*models.NotificationDestination, error) {
	return db.notificationDestRepo.GetByID(ctx, id)
}

func (db *workerDB) GetTeamsBot(ctx context.Context, botID uuid.UUID) (*models.TeamsBot, error) {
	return db.teamsBotRepo.GetByID(ctx, botID)
}

func (db *workerDB) GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	return db.notificationRepo.GetByID(ctx, id)
}

func (db *workerDB) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return db.projectRepo.GetByID(ctx, id)
}