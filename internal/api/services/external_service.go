package services

import (
	"context"
	"fmt"
	"log"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/repositories"
	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/google/uuid"
)

// ExternalService handles external user notification requests
type ExternalService interface {
	SendNotification(ctx context.Context, req *ExternalNotificationRequest) (*ExternalNotificationResponse, error)
}

// externalService implements ExternalService
type externalService struct {
	projectRepo         repositories.ProjectRepository
	destinationRepo     repositories.DestinationRepository
	notificationService NotificationService
}

// NewExternalService creates a new external service
func NewExternalService(
	projectRepo repositories.ProjectRepository,
	destinationRepo repositories.DestinationRepository,
	notificationService NotificationService,
) ExternalService {
	return &externalService{
		projectRepo:         projectRepo,
		destinationRepo:     destinationRepo,
		notificationService: notificationService,
	}
}

// ExternalNotificationRequest represents an external notification request
type ExternalNotificationRequest struct {
	NotifyKey   string         `json:"notify_key" validate:"required,min=3,max=50"`
	Message     string         `json:"message" validate:"required,min=1,max=4000"`
	MessageType string         `json:"message_type" validate:"omitempty,oneof=text file adaptive_card"`
	Priority    string         `json:"priority" validate:"omitempty,oneof=low normal high"`
	TargetIDs   []string       `json:"target_ids" validate:"omitempty"` // If empty, send to all, personal emails, or channel/groupChats IDs
	Mentions    []string       `json:"mentions" validate:"omitempty"`
	Metadata    map[string]any `json:"metadata" validate:"omitempty"`
}

// ExternalNotificationResponse represents the response for external notification
type ExternalNotificationResponse struct {
	NotificationID    string                      `json:"notification_id"`
	Status            string                      `json:"status"`
	Message           string                      `json:"message"`
	ProjectID         string                      `json:"project_id"`
	ProjectName       string                      `json:"project_name"`
	DestinationsCount int                         `json:"destinations_count"`
	Results           []ExternalDestinationResult `json:"results"`
	EstimatedDelivery string                      `json:"estimated_delivery"`
}

// ExternalDestinationResult represents the result for each destination
type ExternalDestinationResult struct {
	DestinationID   string `json:"destination_id"`
	DestinationName string `json:"destination_name"`
	Success         bool   `json:"success"`
	Error           string `json:"error,omitempty"`
	Message         string `json:"message,omitempty"`
}

// SendNotification sends notification to external users
func (s *externalService) SendNotification(ctx context.Context, req *ExternalNotificationRequest) (*ExternalNotificationResponse, error) {
	// Validate and set defaults
	if req.MessageType == "" {
		req.MessageType = "text"
	}
	if req.Priority == "" {
		req.Priority = "normal"
	}

	// Get project by notify_key
	project, err := s.projectRepo.GetByNotifyKey(ctx, req.NotifyKey)
	if err != nil {
		return nil, fmt.Errorf("project not found for notify_key: %s", req.NotifyKey)
	}

	// Check if project is active
	if project.Status != "active" {
		return nil, fmt.Errorf("project is not active: %s", project.Status)
	}

	// Get destinations for the project
	destinations, err := s.destinationRepo.GetByProjectID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get destinations: %w", err)
	}

	if len(destinations) == 0 {
		return nil, fmt.Errorf("no destinations found for project: %s", req.NotifyKey)
	}

	// Note: results may target a subset based on req.TargetIDs

	// Create notification request for NotificationService
	notificationReq := &SendNotificationRequest{
		ProjectID:    project.ID,
		SenderID:     nil, // External notifications don't have a sender
		MessageType:  req.MessageType,
		Content:      req.Message,
		Mentions:     req.Mentions,
		Priority:     req.Priority,
		Metadata:     req.Metadata,
		Destinations: []uuid.UUID{}, // Empty means send to all destinations
	}

	// Send notification using NotificationService (async)
	notification, err := s.notificationService.SendNotification(ctx, notificationReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	// Calculate estimated delivery time based on priority
	estimatedTime := "5-10 minutes"
	switch req.Priority {
	case "high":
		estimatedTime = "2-5 minutes"
	case "low":
		estimatedTime = "10-30 minutes"
	}

	// Create external results based on destinations
	externalResults := make([]ExternalDestinationResult, len(destinations))
	for i, dest := range destinations {
		externalResults[i] = ExternalDestinationResult{
			DestinationID:   dest.ID.String(),
			DestinationName: dest.Name,
			Success:         true, // Will be updated by actors
			Error:           "",
			Message:         "Queued for processing",
		}
	}

	response := &ExternalNotificationResponse{
		NotificationID:    notification.ID.String(),
		Status:            "pending", // Always pending for async processing
		Message:           "Notification queued for processing",
		ProjectID:         project.ID.String(),
		ProjectName:       project.NotifyKey,
		DestinationsCount: len(destinations),
		Results:           externalResults,
		EstimatedDelivery: estimatedTime,
	}

	// Log the external notification
	log.Printf("External notification queued - Project: %s, Destinations: %d, NotificationID: %s",
		req.NotifyKey, len(destinations), notification.ID.String())

	return response, nil
}

// GetProjectDestinations returns available destinations for a project
func (s *externalService) GetProjectDestinations(ctx context.Context, notifyKey string) ([]ExternalDestinationInfo, error) {
	project, err := s.projectRepo.GetByNotifyKey(ctx, notifyKey)
	if err != nil {
		return nil, fmt.Errorf("project not found for notify_key: %s", notifyKey)
	}

	destinations, err := s.destinationRepo.GetByProjectID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get destinations: %w", err)
	}

	result := make([]ExternalDestinationInfo, len(destinations))
	for i, dest := range destinations {
		result[i] = ExternalDestinationInfo{
			ID:          dest.ID.String(),
			Name:        dest.Name,
			Description: dest.Description,
			Status:      dest.Status,
			Targets:     dest.Targets,
		}
	}

	return result, nil
}

// ExternalDestinationInfo represents destination information for external users
type ExternalDestinationInfo struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Status      string                `json:"status"`
	Targets     database.JSONBTargets `json:"targets"`
}
