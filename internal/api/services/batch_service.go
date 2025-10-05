package services

import (
	"context"
	"fmt"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/repositories"
	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/google/uuid"
)

// BatchService defines the interface for batch operations
type BatchService interface {
	// Batch notifications
	SendBatchNotifications(ctx context.Context, req *SendBatchNotificationsRequest) (*BatchResult, error)
	GetBatchStatus(ctx context.Context, batchID uuid.UUID) (*BatchStatus, error)
	GetBatchResults(ctx context.Context, batchID uuid.UUID) (*BatchResults, error)
	CancelBatch(ctx context.Context, batchID uuid.UUID) error

	// Batch templates
	CreateBatchTemplate(ctx context.Context, req *CreateBatchTemplateRequest) (*BatchTemplate, error)
	GetBatchTemplate(ctx context.Context, templateID uuid.UUID) (*BatchTemplate, error)
	ListBatchTemplates(ctx context.Context, projectID *uuid.UUID) ([]*BatchTemplate, int64, error)
	UpdateBatchTemplate(ctx context.Context, templateID uuid.UUID, req *UpdateBatchTemplateRequest) (*BatchTemplate, error)
	DeleteBatchTemplate(ctx context.Context, templateID uuid.UUID) error
}

// batchService implements BatchService
type batchService struct {
	batchRepo           repositories.BatchRepository
	notificationService NotificationService
	projectRepo         repositories.ProjectRepository
	destinationRepo     repositories.DestinationRepository
}

// NewBatchService creates a new batch service
func NewBatchService(
	batchRepo repositories.BatchRepository,
	notificationService NotificationService,
	projectRepo repositories.ProjectRepository,
	destinationRepo repositories.DestinationRepository,
) BatchService {
	return &batchService{
		batchRepo:           batchRepo,
		notificationService: notificationService,
		projectRepo:         projectRepo,
		destinationRepo:     destinationRepo,
	}
}

// Request/Response types for batch operations

type SendBatchNotificationsRequest struct {
	ProjectID   uuid.UUID
	SenderID    *uuid.UUID
	MessageType string
	Content     string
	Priority    string
	Mentions    []string
	Metadata    map[string]any
	Targets     []BatchTarget
	TemplateID  *uuid.UUID
	ScheduleAt  *string
	ExpiresAt   *string
	MaxRetries  int
	RetryDelay  *string
}

type BatchTarget struct {
	DestinationID  *uuid.UUID
	ConversationID *string
	UserID         *string
	Email          *string
	CustomData     map[string]any
}

type BatchResult struct {
	BatchID             uuid.UUID `json:"batch_id"`
	Status              string    `json:"status"`
	TotalTargets        int       `json:"total_targets"`
	CreatedAt           time.Time `json:"created_at"`
	EstimatedCompletion time.Time `json:"estimated_completion"`
}

type BatchStatus struct {
	BatchID             uuid.UUID  `json:"batch_id"`
	Status              string     `json:"status"`
	TotalTargets        int        `json:"total_targets"`
	CompletedTargets    int        `json:"completed_targets"`
	FailedTargets       int        `json:"failed_targets"`
	PendingTargets      int        `json:"pending_targets"`
	Progress            float64    `json:"progress"`
	CreatedAt           time.Time  `json:"created_at"`
	StartedAt           *time.Time `json:"started_at,omitempty"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	EstimatedCompletion *time.Time `json:"estimated_completion,omitempty"`
}

type BatchResults struct {
	BatchID      uuid.UUID           `json:"batch_id"`
	Status       string              `json:"status"`
	TotalTargets int                 `json:"total_targets"`
	Results      []BatchTargetResult `json:"results"`
	Summary      BatchSummary        `json:"summary"`
	CreatedAt    time.Time           `json:"created_at"`
	CompletedAt  *time.Time          `json:"completed_at,omitempty"`
}

type BatchTargetResult struct {
	TargetID       uuid.UUID  `json:"target_id"`
	DestinationID  *uuid.UUID `json:"destination_id,omitempty"`
	ConversationID *string    `json:"conversation_id,omitempty"`
	Status         string     `json:"status"`
	Success        bool       `json:"success"`
	ErrorMessage   *string    `json:"error_message,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	RetryCount     int        `json:"retry_count"`
}

type BatchSummary struct {
	TotalSent      int     `json:"total_sent"`
	TotalFailed    int     `json:"total_failed"`
	TotalPending   int     `json:"total_pending"`
	SuccessRate    float64 `json:"success_rate"`
	AverageLatency float64 `json:"average_latency"`
}

type CreateBatchTemplateRequest struct {
	Name        string
	Description string
	ProjectID   uuid.UUID
	MessageType string
	Content     string
	Priority    string
	Mentions    []string
	Metadata    map[string]any
	Targets     []BatchTarget
	IsActive    bool
}

type UpdateBatchTemplateRequest struct {
	Name        *string
	Description *string
	Content     *string
	Priority    *string
	Mentions    []string
	Metadata    map[string]any
	Targets     []BatchTarget
	IsActive    *bool
}

type BatchTemplate struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	ProjectID   uuid.UUID      `json:"project_id"`
	MessageType string         `json:"message_type"`
	Content     string         `json:"content"`
	Priority    string         `json:"priority"`
	Mentions    []string       `json:"mentions"`
	Metadata    map[string]any `json:"metadata"`
	Targets     []BatchTarget  `json:"targets"`
	IsActive    bool           `json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Implementation methods

func (s *batchService) SendBatchNotifications(ctx context.Context, req *SendBatchNotificationsRequest) (*BatchResult, error) {
	// Validate project exists
	_, err := s.projectRepo.GetByID(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Create batch record
	batch := &database.BatchNotification{
		BaseModel: database.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ProjectID:    req.ProjectID,
		SenderID:     req.SenderID,
		MessageType:  req.MessageType,
		Content:      req.Content,
		Priority:     req.Priority,
		Mentions:     database.JSONBStringArray(req.Mentions),
		Metadata:     database.JSONBObject(req.Metadata),
		Status:       "pending",
		TotalTargets: len(req.Targets),
		MaxRetries:   req.MaxRetries,
	}

	// Parse schedule time if provided
	if req.ScheduleAt != nil && *req.ScheduleAt != "" {
		if t, err := time.Parse(time.RFC3339, *req.ScheduleAt); err == nil {
			batch.ScheduledAt = &t
		}
	}

	// Parse expires time if provided
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		if t, err := time.Parse(time.RFC3339, *req.ExpiresAt); err == nil {
			batch.ExpiresAt = &t
		}
	}

	// Save batch to database
	if err := s.batchRepo.Create(ctx, batch); err != nil {
		return nil, fmt.Errorf("failed to create batch: %w", err)
	}

	// Create batch targets
	var batchTargets []*database.BatchTarget
	for i, target := range req.Targets {
		batchTarget := &database.BatchTarget{
			BaseModel: database.BaseModel{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			BatchID:        batch.ID,
			TargetIndex:    i,
			DestinationID:  target.DestinationID,
			ConversationID: target.ConversationID,
			UserID:         target.UserID,
			Email:          target.Email,
			CustomData:     database.JSONBObject(target.CustomData),
			Status:         "pending",
		}
		batchTargets = append(batchTargets, batchTarget)
	}

	// Save batch targets
	if err := s.batchRepo.CreateBatchTargets(ctx, batchTargets); err != nil {
		return nil, fmt.Errorf("failed to create batch targets: %w", err)
	}

	// Start processing batch in background
	go s.processBatch(context.Background(), batch.ID)

	// Calculate estimated completion time
	estimatedCompletion := time.Now().Add(time.Duration(len(req.Targets)) * 2 * time.Second)

	return &BatchResult{
		BatchID:             batch.ID,
		Status:              batch.Status,
		TotalTargets:        len(req.Targets),
		CreatedAt:           batch.CreatedAt,
		EstimatedCompletion: estimatedCompletion,
	}, nil
}

func (s *batchService) GetBatchStatus(ctx context.Context, batchID uuid.UUID) (*BatchStatus, error) {
	batch, err := s.batchRepo.GetByID(ctx, batchID)
	if err != nil {
		return nil, fmt.Errorf("batch not found: %w", err)
	}

	// Get batch targets status
	targets, err := s.batchRepo.GetBatchTargets(ctx, batchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch targets: %w", err)
	}

	// Calculate status
	completed := 0
	failed := 0
	pending := 0

	for _, target := range targets {
		switch target.Status {
		case "completed":
			completed++
		case "failed":
			failed++
		case "pending":
			pending++
		}
	}

	progress := 0.0
	if batch.TotalTargets > 0 {
		progress = float64(completed+failed) / float64(batch.TotalTargets) * 100
	}

	status := &BatchStatus{
		BatchID:          batch.ID,
		Status:           batch.Status,
		TotalTargets:     batch.TotalTargets,
		CompletedTargets: completed,
		FailedTargets:    failed,
		PendingTargets:   pending,
		Progress:         progress,
		CreatedAt:        batch.CreatedAt,
		StartedAt:        batch.StartedAt,
		CompletedAt:      batch.CompletedAt,
	}

	// Calculate estimated completion if still processing
	if batch.Status == "processing" && pending > 0 {
		estimatedCompletion := time.Now().Add(time.Duration(pending) * 2 * time.Second)
		status.EstimatedCompletion = &estimatedCompletion
	}

	return status, nil
}

func (s *batchService) GetBatchResults(ctx context.Context, batchID uuid.UUID) (*BatchResults, error) {
	batch, err := s.batchRepo.GetByID(ctx, batchID)
	if err != nil {
		return nil, fmt.Errorf("batch not found: %w", err)
	}

	// Get batch targets
	targets, err := s.batchRepo.GetBatchTargets(ctx, batchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch targets: %w", err)
	}

	// Convert to results
	var results []BatchTargetResult
	summary := BatchSummary{}

	for _, target := range targets {
		result := BatchTargetResult{
			TargetID:       target.ID,
			DestinationID:  target.DestinationID,
			ConversationID: target.ConversationID,
			Status:         target.Status,
			Success:        target.Status == "completed",
			ErrorMessage:   target.ErrorMessage,
			SentAt:         target.SentAt,
			RetryCount:     target.RetryCount,
		}
		results = append(results, result)

		// Update summary
		if target.Status == "completed" {
			summary.TotalSent++
		} else if target.Status == "failed" {
			summary.TotalFailed++
		} else {
			summary.TotalPending++
		}
	}

	// Calculate success rate
	if batch.TotalTargets > 0 {
		summary.SuccessRate = float64(summary.TotalSent) / float64(batch.TotalTargets) * 100
	}

	return &BatchResults{
		BatchID:      batch.ID,
		Status:       batch.Status,
		TotalTargets: batch.TotalTargets,
		Results:      results,
		Summary:      summary,
		CreatedAt:    batch.CreatedAt,
		CompletedAt:  batch.CompletedAt,
	}, nil
}

func (s *batchService) CancelBatch(ctx context.Context, batchID uuid.UUID) error {
	batch, err := s.batchRepo.GetByID(ctx, batchID)
	if err != nil {
		return fmt.Errorf("batch not found: %w", err)
	}

	// Only cancel if batch is pending or processing
	if batch.Status != "pending" && batch.Status != "processing" {
		return fmt.Errorf("batch cannot be cancelled in status: %s", batch.Status)
	}

	// Update batch status
	batch.Status = "cancelled"
	batch.UpdatedAt = time.Now()
	if err := s.batchRepo.Update(ctx, batch); err != nil {
		return fmt.Errorf("failed to cancel batch: %w", err)
	}

	return nil
}

func (s *batchService) CreateBatchTemplate(ctx context.Context, req *CreateBatchTemplateRequest) (*BatchTemplate, error) {
	// Implementation will be added when template functionality is needed
	return nil, fmt.Errorf("not implemented")
}

func (s *batchService) GetBatchTemplate(ctx context.Context, templateID uuid.UUID) (*BatchTemplate, error) {
	// Implementation will be added when template functionality is needed
	return nil, fmt.Errorf("not implemented")
}

func (s *batchService) ListBatchTemplates(ctx context.Context, projectID *uuid.UUID) ([]*BatchTemplate, int64, error) {
	// Implementation will be added when template functionality is needed
	return []*BatchTemplate{}, 0, nil
}

func (s *batchService) UpdateBatchTemplate(ctx context.Context, templateID uuid.UUID, req *UpdateBatchTemplateRequest) (*BatchTemplate, error) {
	// Implementation will be added when template functionality is needed
	return nil, fmt.Errorf("not implemented")
}

func (s *batchService) DeleteBatchTemplate(ctx context.Context, templateID uuid.UUID) error {
	// Implementation will be added when template functionality is needed
	return fmt.Errorf("not implemented")
}

// processBatch processes a batch notification in the background
func (s *batchService) processBatch(ctx context.Context, batchID uuid.UUID) {
	// Update batch status to processing
	batch, err := s.batchRepo.GetByID(ctx, batchID)
	if err != nil {
		fmt.Printf("Failed to get batch %s: %v\n", batchID, err)
		return
	}

	batch.Status = "processing"
	batch.StartedAt = &[]time.Time{time.Now()}[0]
	batch.UpdatedAt = time.Now()
	if err := s.batchRepo.Update(ctx, batch); err != nil {
		fmt.Printf("Failed to update batch status: %v\n", err)
		return
	}

	// Get batch targets
	targets, err := s.batchRepo.GetBatchTargets(ctx, batchID)
	if err != nil {
		fmt.Printf("Failed to get batch targets: %v\n", err)
		return
	}

	// Process each target
	for _, target := range targets {
		// Create notification for this target
		notificationReq := &SendNotificationRequest{
			ProjectID:    batch.ProjectID,
			SenderID:     batch.SenderID,
			MessageType:  batch.MessageType,
			Content:      batch.Content,
			Priority:     batch.Priority,
			Mentions:     []string(batch.Mentions),
			Metadata:     map[string]any(batch.Metadata),
			Destinations: []uuid.UUID{}, // Will be determined by target
		}

		// If target has destination ID, use it
		if target.DestinationID != nil {
			notificationReq.Destinations = []uuid.UUID{*target.DestinationID}
		}

		// Send notification
		_, err := s.notificationService.SendNotification(ctx, notificationReq)
		if err != nil {
			// Update target status to failed
			target.Status = "failed"
			errorMsg := err.Error()
			target.ErrorMessage = &errorMsg
			target.UpdatedAt = time.Now()
			s.batchRepo.UpdateBatchTarget(ctx, target)
		} else {
			// Update target status to completed
			target.Status = "completed"
			sentAt := time.Now()
			target.SentAt = &sentAt
			target.UpdatedAt = time.Now()
			s.batchRepo.UpdateBatchTarget(ctx, target)
		}
	}

	// Update batch status to completed
	batch.Status = "completed"
	batch.CompletedAt = &[]time.Time{time.Now()}[0]
	batch.UpdatedAt = time.Now()
	s.batchRepo.Update(ctx, batch)
}
