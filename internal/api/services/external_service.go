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
	projectRepo      repositories.ProjectRepository
	destinationRepo  repositories.DestinationRepository
	broadcastService BroadcastService
}

// NewExternalService creates a new external service
func NewExternalService(
	projectRepo repositories.ProjectRepository,
	destinationRepo repositories.DestinationRepository,
	broadcastService BroadcastService,
) ExternalService {
	return &externalService{
		projectRepo:      projectRepo,
		destinationRepo:  destinationRepo,
		broadcastService: broadcastService,
	}
}

// ExternalNotificationRequest represents an external notification request
type ExternalNotificationRequest struct {
	NotifyKey    string         `json:"notify_key" validate:"required,min=3,max=50"`
	Message      string         `json:"message" validate:"required,min=1,max=4000"`
	MessageType  string         `json:"message_type" validate:"omitempty,oneof=text file adaptive_card"`
	Priority     string         `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	Destinations []string       `json:"destinations" validate:"omitempty"` // If empty, send to all
	Mentions     []string       `json:"mentions" validate:"omitempty"`
	Metadata     map[string]any `json:"metadata" validate:"omitempty"`
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

	// Filter destinations if specific ones are requested
	var targetDestinations []*database.Destination
	if len(req.Destinations) == 0 || (len(req.Destinations) == 1 && req.Destinations[0] == "all") {
		// Send to all destinations
		targetDestinations = destinations
	} else {
		// Send to specific destinations
		destinationMap := make(map[string]*database.Destination)
		for _, dest := range destinations {
			destinationMap[dest.ID.String()] = dest
		}

		for _, destID := range req.Destinations {
			if dest, exists := destinationMap[destID]; exists {
				targetDestinations = append(targetDestinations, dest)
			} else {
				log.Printf("Warning: Destination %s not found for project %s", destID, req.NotifyKey)
			}
		}
	}

	if len(targetDestinations) == 0 {
		return nil, fmt.Errorf("no valid destinations found")
	}

	// Prepare destination IDs for broadcast service
	destinationIDs := make([]uuid.UUID, len(targetDestinations))
	for i, dest := range targetDestinations {
		destinationIDs[i] = dest.ID
	}

	// Send notification using broadcast service
	results, err := s.broadcastService.SendToDestinations(ctx, destinationIDs, req.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	// Map results to external format
	externalResults := make([]ExternalDestinationResult, len(results.Results))
	successCount := 0

	for i, result := range results.Results {
		// Find corresponding destination
		var destName string
		for _, dest := range targetDestinations {
			if dest.ID.String() == result.TargetID {
				destName = dest.Name
				break
			}
		}

		externalResults[i] = ExternalDestinationResult{
			DestinationID:   result.TargetID,
			DestinationName: destName,
			Success:         result.Success,
			Error:           result.Error,
			Message:         "",
		}

		if result.Success {
			successCount++
		}
	}

	// Determine overall status
	status := "sent"
	if successCount == 0 {
		status = "failed"
	} else if successCount < len(results.Results) {
		status = "partial"
	}

	// Generate notification ID
	notificationID := uuid.New().String()

	response := &ExternalNotificationResponse{
		NotificationID:    notificationID,
		Status:            status,
		Message:           "Notification processed successfully",
		ProjectID:         project.ID.String(),
		ProjectName:       project.NotifyKey,
		DestinationsCount: len(targetDestinations),
		Results:           externalResults,
		EstimatedDelivery: "5-10 minutes",
	}

	// Log the external notification
	log.Printf("External notification sent - Project: %s, Destinations: %d, Success: %d/%d",
		req.NotifyKey, len(targetDestinations), successCount, len(results.Results))

	return response, nil
}

// ValidateNotifyKey validates if a notify key exists and is active
func (s *externalService) ValidateNotifyKey(ctx context.Context, notifyKey string) error {
	project, err := s.projectRepo.GetByNotifyKey(ctx, notifyKey)
	if err != nil {
		return fmt.Errorf("project not found for notify_key: %s", notifyKey)
	}

	if project.Status != "active" {
		return fmt.Errorf("project is not active: %s", project.Status)
	}

	return nil
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
