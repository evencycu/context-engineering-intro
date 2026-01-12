package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/repositories"
	"github.com/google/uuid"
)

// TemplateService defines template operations
type TemplateService interface {
	Create(ctx context.Context, req *CreateTemplateRequest) (*models.Template, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Template, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateTemplateRequest) (*models.Template, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type templateService struct {
	repo repositories.TemplateRepository
}

// NewTemplateService creates a new template service
func NewTemplateService(repo repositories.TemplateRepository) TemplateService {
	return &templateService{repo: repo}
}

// DTOs
type CreateTemplateRequest struct {
	ProjectID            uuid.UUID                     `json:"projectId"`            // Set by handler from path param, not from JSON
	Name                 string                        `json:"name"`
	Description          string                        `json:"description"`
	Variables            models.JSONBTemplateVariables `json:"variables"`
	DefaultJsonStructure string                        `json:"defaultJsonStructure"`
}

type UpdateTemplateRequest struct {
	Name                 string                        `json:"name"`
	Description          string                        `json:"description"`
	Variables            models.JSONBTemplateVariables `json:"variables"`
	DefaultJsonStructure string                        `json:"defaultJsonStructure"`
}

func (s *templateService) Create(ctx context.Context, req *CreateTemplateRequest) (*models.Template, error) {
	// Validate JSON structure
	if req.DefaultJsonStructure != "" && !isValidJSON(req.DefaultJsonStructure) {
		return nil, errors.New("invalid JSON structure")
	}

	template := &models.Template{
		ProjectID:            req.ProjectID,
		Name:                 req.Name,
		Description:          req.Description,
		Variables:            req.Variables,
		DefaultJsonStructure: req.DefaultJsonStructure,
	}

	if err := s.repo.Create(ctx, template); err != nil {
		return nil, err
	}
	return template, nil
}

func (s *templateService) GetByID(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *templateService) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Template, error) {
	return s.repo.GetByProjectID(ctx, projectID)
}

func (s *templateService) Update(ctx context.Context, id uuid.UUID, req *UpdateTemplateRequest) (*models.Template, error) {
	// Validate JSON structure
	if req.DefaultJsonStructure != "" && !isValidJSON(req.DefaultJsonStructure) {
		return nil, errors.New("invalid JSON structure")
	}

	template, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	template.Name = req.Name
	template.Description = req.Description
	template.Variables = req.Variables
	template.DefaultJsonStructure = req.DefaultJsonStructure

	if err := s.repo.Update(ctx, template); err != nil {
		return nil, err
	}
	return template, nil
}

func (s *templateService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func isValidJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
