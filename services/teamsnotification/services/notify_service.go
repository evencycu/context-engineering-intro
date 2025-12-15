package services

import (
	"context"
	"fmt"
	"log"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/repositories"
)

// NotifyService handles external user notification requests
type NotifyService interface {
	SendNotification(ctx context.Context, req *NotifyRequest) (*NotifyResponse, error)
	GetProjectDestinations(ctx context.Context, notifyKey string) ([]DestinationInfo, error)
}

// notifyService implements NotifyService
type notifyService struct {
	projectRepo         repositories.ProjectRepository
	destinationRepo     repositories.DestinationRepository
	notificationService NotificationService
}

// NewNotifyService creates a new notify service
func NewNotifyService(
	projectRepo repositories.ProjectRepository,
	destinationRepo repositories.DestinationRepository,
	notificationService NotificationService,
) NotifyService {
	return &notifyService{
		projectRepo:         projectRepo,
		destinationRepo:     destinationRepo,
		notificationService: notificationService,
	}
}

// NotifyRequest represents an external notification request
type NotifyRequest struct {
	NotifyKey   string           `json:"notify_key" validate:"required,min=3,max=50"`
	Message     string           `json:"message" validate:"required,min=1,max=4000"`
	MessageType string           `json:"message_type" validate:"omitempty,oneof=text file adaptive_card"`
	Priority    string           `json:"priority" validate:"omitempty,oneof=low normal high"`
	TargetIDs   []string         `json:"target_ids" validate:"omitempty"` // If empty, send to all, personal emails, or channel/groupChats IDs
	Mentions    []string         `json:"mentions" validate:"omitempty"`
	Metadata    map[string]any   `json:"metadata" validate:"omitempty"`
	Attachments []map[string]any `json:"attachments" validate:"omitempty"`
}

// NotifyResponse represents the response for external notification
type NotifyResponse struct {
	NotificationID    string              `json:"notification_id"`
	Status            string              `json:"status"`
	Message           string              `json:"message"`
	ProjectID         string              `json:"project_id"`
	ProjectName       string              `json:"project_name"`
	DestinationsCount int                 `json:"destinations_count"`
	Results           []DestinationResult `json:"results"`
	EstimatedDelivery string              `json:"estimated_delivery"`
}

// DestinationResult represents the result for each destination
type DestinationResult struct {
	DestinationID   string `json:"destination_id"`
	DestinationName string `json:"destination_name"`
	Success         bool   `json:"success"`
	Error           string `json:"error,omitempty"`
	Message         string `json:"message,omitempty"`
}

// SendNotification sends notification to external users
func (s *notifyService) SendNotification(ctx context.Context, req *NotifyRequest) (*NotifyResponse, error) {
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
		ProjectID:   project.ID,
		SenderID:    nil, // External notifications don't have a sender
		MessageType: req.MessageType,
		Content:     req.Message,
		Mentions:    req.Mentions,
		Attachments: req.Attachments,
		Priority:    req.Priority,
		Metadata:    req.Metadata,
		Targets:     req.TargetIDs, // Empty means send to all destinations
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
	destinationResults := make([]DestinationResult, len(destinations))
	for i, dest := range destinations {
		destinationResults[i] = DestinationResult{
			DestinationID:   dest.ID.String(),
			DestinationName: dest.Name,
			Success:         true, // Will be updated by actors
			Error:           "",
			Message:         "Queued for processing",
		}
	}

	response := &NotifyResponse{
		NotificationID:    notification.ID.String(),
		Status:            string(models.NotificationStatusPending), // Always pending for async processing
		Message:           "Notification queued for processing",
		ProjectID:         project.ID.String(),
		ProjectName:       project.NotifyKey,
		DestinationsCount: len(destinations),
		Results:           destinationResults,
		EstimatedDelivery: estimatedTime,
	}

	// Log the external notification
	log.Printf("External notification queued - Project: %s, Destinations: %d, NotificationID: %s",
		req.NotifyKey, len(destinations), notification.ID.String())

	return response, nil
}

// GetProjectDestinations returns available destinations for a project
func (s *notifyService) GetProjectDestinations(ctx context.Context, notifyKey string) ([]DestinationInfo, error) {
	project, err := s.projectRepo.GetByNotifyKey(ctx, notifyKey)
	if err != nil {
		return nil, fmt.Errorf("project not found for notify_key: %s", notifyKey)
	}

	destinations, err := s.destinationRepo.GetByProjectID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get destinations: %w", err)
	}

	result := make([]DestinationInfo, len(destinations))
	for i, dest := range destinations {
		result[i] = DestinationInfo{
			ID:          dest.ID.String(),
			Name:        dest.Name,
			Description: dest.Description,
			Status:      dest.Status,
			Targets:     dest.Targets,
		}
	}

	return result, nil
}

// DestinationInfo represents destination information for external users
type DestinationInfo struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Status      string              `json:"status"`
	Targets     models.JSONBTargets `json:"targets"`
}
