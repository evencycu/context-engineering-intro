package actor

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/google/uuid"
)

// BillingService interface for recording usage
type BillingService interface {
	RecordUsage(ctx context.Context, req *RecordUsageRequest) error
}

// RecordUsageRequest represents a usage record request
type RecordUsageRequest struct {
	CompanyID      uuid.UUID
	ProjectID      *uuid.UUID
	UserID         *uuid.UUID
	NotificationID *uuid.UUID
	BotID          *uuid.UUID
	BotType        *models.BotType
	RecordType     string
	Quantity       int
	UnitPrice      float64
	TotalCost      float64
	BillingPeriod  time.Time
}

// WorkerDB defines the operations needed by notification workers
type WorkerDB interface {
	UpdateNotificationDestination(ctx context.Context, id uuid.UUID, update *NotificationDestinationUpdate) error
	GetBotInstallation(ctx context.Context, conversationID string) (*models.BotInstallation, error)
	GetNotificationDestinationByID(ctx context.Context, id uuid.UUID) (*models.NotificationDestination, error)
	GetTeamsBot(ctx context.Context, botID uuid.UUID) (*models.TeamsBot, error)
	GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error)
}

// NotificationWorker handles sending a single notification with retry logic
type NotificationWorker struct {
	ID                 string
	NotificationDestID uuid.UUID
	NotificationID     uuid.UUID
	ConversationID     string
	Installation       *models.BotInstallation
	DB                 WorkerDB
	Connector          TeamsConnector
	BillingService     BillingService
	TaskQueue          TaskQueue
	MaxRetries         int
	CurrentRetry       int
	Priority           string
	StopChan           chan struct{}
	DoneChan           chan *WorkerResult
}

// WorkerResult represents the result of a worker's execution
type WorkerResult struct {
	WorkerID           string
	NotificationDestID uuid.UUID
	Success            bool
	Error              error
	Retrying           bool
	NextRetryAt        *time.Time
	TeamsMessageID     string
}

// NewNotificationWorker creates a new notification worker
func NewNotificationWorker(
	notificationDestID uuid.UUID,
	notificationID uuid.UUID,
	conversationID string,
	installation *models.BotInstallation,
	db WorkerDB,
	connector TeamsConnector,
	billingService BillingService,
	taskQueue TaskQueue,
) *NotificationWorker {
	return &NotificationWorker{
		ID:                 fmt.Sprintf("worker-%s", notificationDestID.String()[:8]),
		NotificationDestID: notificationDestID,
		NotificationID:     notificationID,
		ConversationID:     conversationID,
		Installation:       installation,
		DB:                 db,
		Connector:          connector,
		BillingService:     billingService,
		TaskQueue:          taskQueue,
		MaxRetries:         5,
		CurrentRetry:       0,
		StopChan:           make(chan struct{}),
		DoneChan:           make(chan *WorkerResult, 1),
	}
}

// Start begins the worker's execution
func (w *NotificationWorker) Start(ctx context.Context) {
	go w.run(ctx)
}

// Stop stops the worker
func (w *NotificationWorker) Stop() {
	close(w.StopChan)
}

// Wait waits for the worker to complete and returns the result
func (w *NotificationWorker) Wait() *WorkerResult {
	return <-w.DoneChan
}

// run is the main worker loop
func (w *NotificationWorker) run(ctx context.Context) {
	defer close(w.DoneChan)

	log.Printf("[%s] Starting worker for notification_dest=%s", w.ID, w.NotificationDestID)

	// Mark as processing
	now := time.Now()
	processing := string(models.NotificationStatusProcessing)
	_ = w.DB.UpdateNotificationDestination(ctx, w.NotificationDestID, &NotificationDestinationUpdate{
		Status:         &processing,
		ActorID:        &w.ID, // Keep field name in DB for now
		FirstAttemptAt: &now,
		LastAttemptAt:  &now,
	})

	// Fetch notification details
	notification, err := w.DB.GetNotificationByID(ctx, w.NotificationID)
	if err != nil {
		w.handleFatalError(ctx, fmt.Errorf("failed to fetch notification: %w", err))
		return
	}

	// Prepare Payload
	payload := &TeamsSendPayload{
		ConversationID: w.ConversationID,
		ServiceURL:     w.Installation.ServiceURL,
		Message:        notification.Content,
		Attachments:    notification.Attachments,
		RecipientID:    w.Installation.RecipientID,
		FromID:         w.Installation.FromID,
		TenantID:       w.Installation.TeamsTenantID,
	}

	// Execution attempt
	result, err := w.Connector.Send(ctx, payload)
	
	// Update attempt time in DB
	attemptTime := time.Now()
	_ = w.DB.UpdateNotificationDestination(ctx, w.NotificationDestID, &NotificationDestinationUpdate{
		LastAttemptAt: &attemptTime,
	})

	if err == nil && result.Success {
		w.handleSuccess(ctx, result)
	} else {
		w.handleFailure(ctx, result, err)
	}
}

func (w *NotificationWorker) handleSuccess(ctx context.Context, res *TeamsSendResult) {
	log.Printf("[%s] Successfully sent notification", w.ID)
	sentStatus := string(models.NotificationStatusSent)
	sentAt := time.Now()
	
	_ = w.DB.UpdateNotificationDestination(ctx, w.NotificationDestID, &NotificationDestinationUpdate{
		Status:         &sentStatus,
		TeamsMessageID: &res.TeamsMessageID,
		SentAt:         &sentAt,
		ErrorMessage:   nil,
	})

	// Record usage
	if w.BillingService != nil {
		_ = w.recordUsage(ctx)
	}

	w.DoneChan <- &WorkerResult{
		WorkerID:           w.ID,
		NotificationDestID: w.NotificationDestID,
		Success:            true,
		TeamsMessageID:     res.TeamsMessageID,
	}
}

func (w *NotificationWorker) handleFailure(ctx context.Context, res *TeamsSendResult, err error) {
	w.CurrentRetry++
	
	statusCode := 0
	if res != nil {
		statusCode = res.StatusCode
	}

	failureReason := w.classifyFailure(statusCode, err)
	isRetryable := w.isRetryableError(statusCode, err)

	if !isRetryable || w.CurrentRetry >= w.MaxRetries {
		log.Printf("[%s] Permanent failure or max retries reached: %v", w.ID, err)
		failedStatus := string(models.NotificationStatusFailed)
		errMsg := "unknown error"
		if err != nil {
			errMsg = err.Error()
		}
		
		_ = w.DB.UpdateNotificationDestination(ctx, w.NotificationDestID, &NotificationDestinationUpdate{
			Status:        &failedStatus,
			ErrorMessage:  &errMsg,
			RetryCount:    &w.CurrentRetry,
			FailureReason: &failureReason,
		})

		w.DoneChan <- &WorkerResult{
			WorkerID:           w.ID,
			NotificationDestID: w.NotificationDestID,
			Success:            false,
			Error:              err,
		}
		return
	}

	// Schedule Retry
	var retryAfter *int
	if res != nil && res.RetryAfter > 0 {
		retryAfter = &res.RetryAfter
	}
	
	w.scheduleRetry(ctx, failureReason, retryAfter, err)

	w.DoneChan <- &WorkerResult{
		WorkerID:           w.ID,
		NotificationDestID: w.NotificationDestID,
		Success:            false,
		Retrying:           true,
		NextRetryAt:        w.calculateNextRetryTime(w.CurrentRetry, retryAfter),
	}
}

func (w *NotificationWorker) handleFatalError(ctx context.Context, err error) {
	log.Printf("[%s] Fatal error: %v", w.ID, err)
	failedStatus := string(models.NotificationStatusFailed)
	errMsg := err.Error()
	_ = w.DB.UpdateNotificationDestination(ctx, w.NotificationDestID, &NotificationDestinationUpdate{
		Status:       &failedStatus,
		ErrorMessage: &errMsg,
	})
	w.DoneChan <- &WorkerResult{
		WorkerID:           w.ID,
		NotificationDestID: w.NotificationDestID,
		Success:            false,
		Error:              err,
	}
}

func (w *NotificationWorker) isRetryableError(statusCode int, err error) bool {
	if err != nil && strings.Contains(err.Error(), "circuit breaker is open") {
		return true
	}
	return statusCode == 429 || statusCode == 408 || (statusCode >= 500 && statusCode < 600)
}

func (w *NotificationWorker) classifyFailure(statusCode int, err error) string {
	if err != nil && strings.Contains(err.Error(), "circuit breaker is open") {
		return "circuit_breaker"
	}
	switch {
	case statusCode == 429:
		return "rate_limit"
	case statusCode == 408:
		return "timeout"
	case statusCode >= 500 && statusCode < 600:
		return "server_error"
	default:
		return "network_error"
	}
}

func (w *NotificationWorker) scheduleRetry(ctx context.Context, reason string, retryAfter *int, err error) {
	nextRetryAt := w.calculateNextRetryTime(w.CurrentRetry, retryAfter)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	_ = w.DB.UpdateNotificationDestination(ctx, w.NotificationDestID, &NotificationDestinationUpdate{
		ErrorMessage:  &errMsg,
		RetryCount:    &w.CurrentRetry,
		NextRetryAt:   nextRetryAt,
		FailureReason: &reason,
		RetryAfter:    retryAfter,
	})

	// Put back to queue
	prio := PriorityNormal
	if w.Priority == "high" || reason == "rate_limit" {
		prio = PriorityHigh
	} else if w.Priority == "low" {
		prio = PriorityLow
	}

	// Use background context for re-queueing to ensure it happens even if request context is done
	go func(id uuid.UUID, delay time.Duration, p Priority) {
		time.Sleep(delay)
		_ = w.TaskQueue.Requeue(context.Background(), id) 
	}(w.NotificationDestID, time.Until(*nextRetryAt), prio)
}

func (w *NotificationWorker) calculateNextRetryTime(retryCount int, retryAfter *int) *time.Time {
	var backoff time.Duration
	if retryAfter != nil {
		backoff = time.Duration(*retryAfter+1) * time.Second
	} else {
		backoff = time.Duration(math.Pow(2, float64(retryCount))) * time.Second
		if backoff > 10*time.Minute {
			backoff = 10 * time.Minute
		}
		jitter := backoff / 10
		backoff += time.Duration(rand.Int63n(int64(2*jitter))) - jitter
	}
	next := time.Now().Add(backoff)
	return &next
}

func (w *NotificationWorker) recordUsage(ctx context.Context) error {
	notification, err := w.DB.GetNotificationByID(ctx, w.NotificationID)
	if err != nil {
		return err
	}
	project, err := w.DB.GetProjectByID(ctx, notification.ProjectID)
	if err != nil {
		return err
	}

	var botType *models.BotType
	bot, err := w.DB.GetTeamsBot(ctx, w.Installation.BotID)
	if err == nil && bot != nil {
		bt := models.BotType(bot.Type)
		botType = &bt
	}

	now := time.Now()
	return w.BillingService.RecordUsage(ctx, &RecordUsageRequest{
		CompanyID:      project.CompanyID,
		ProjectID:      &notification.ProjectID,
		UserID:         notification.SenderID,
		NotificationID: &w.NotificationID,
		BotID:          &w.Installation.BotID,
		BotType:        botType,
		RecordType:     "notification",
		Quantity:       1,
		UnitPrice:      0.001,
		TotalCost:      0.001,
		BillingPeriod:  time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC),
	})
}