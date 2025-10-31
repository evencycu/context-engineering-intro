package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/evencycu/TeamsNotifyGoV2/libs/models"
	"github.com/evencycu/TeamsNotifyGoV2/services/teamsnotification/repositories"
	"github.com/google/uuid"
)

type ProvisionService interface {
	Create(ctx context.Context, req *ProvisionCreateRequest) (*ProvisionCreateResponse, error)
	Read(ctx context.Context, notifyKey string) (*ProvisionReadResponse, error)
	Update(ctx context.Context, notifyKey string, req *ProvisionUpdateRequest) (*ProvisionReadResponse, error)
	SetStatus(ctx context.Context, notifyKey string, enabled bool) error
}

type provisionService struct {
	companyRepo     repositories.CompanyRepository
	projectRepo     repositories.ProjectRepository
	destinationRepo repositories.DestinationRepository
}

func NewProvisionService(
	companyRepo repositories.CompanyRepository,
	projectRepo repositories.ProjectRepository,
	destinationRepo repositories.DestinationRepository,
) ProvisionService {
	return &provisionService{
		companyRepo:     companyRepo,
		projectRepo:     projectRepo,
		destinationRepo: destinationRepo,
	}
}

type ProvisionCreateRequest struct {
	CompanyID     uuid.UUID           `json:"companyId"`
	CreatedBy     uuid.UUID           `json:"createdBy"`
	ProjectName   string              `json:"projectName"`
	ProjectDesc   string              `json:"projectDescription"`
	TeamsTenantID string              `json:"teamsTenantId"`
	Targets       models.JSONBTargets `json:"targets"`
}

type ProvisionCreateResponse struct {
	NotifyKey   string              `json:"notifyKey"`
	Project     *models.Project     `json:"project"`
	Destination *models.Destination `json:"destination"`
}

type ProvisionReadResponse struct {
	Project      *models.Project       `json:"project"`
	Destinations []*models.Destination `json:"destinations"`
}

type ProvisionUpdateRequest struct {
	ProjectName   *string              `json:"project_name,omitempty"`
	ProjectDesc   *string              `json:"project_description,omitempty"`
	TeamsTenantID *string              `json:"teams_tenant_id,omitempty"`
	Targets       *models.JSONBTargets `json:"targets,omitempty"`
}

func (s *provisionService) Create(ctx context.Context, req *ProvisionCreateRequest) (*ProvisionCreateResponse, error) {
	if req.CompanyID == uuid.Nil || req.CreatedBy == uuid.Nil {
		return nil, errors.New("company_id and created_by are required")
	}
	if strings.TrimSpace(req.ProjectName) == "" {
		return nil, errors.New("project_name is required")
	}
	if strings.TrimSpace(req.TeamsTenantID) == "" {
		return nil, errors.New("teams_tenant_id is required")
	}
	if len(req.Targets) == 0 {
		return nil, errors.New("at least one target is required")
	}

	// 1) Check company exists
	if _, err := s.companyRepo.GetByID(ctx, req.CompanyID); err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// 2) Insert project with pre-generated ID and notify_key
	projectID := uuid.New()
	notifyKey := generateNotifyKey(projectID, req.CompanyID)

	project := &models.Project{
		BaseModel:    models.BaseModel{ID: projectID},
		CompanyID:    req.CompanyID,
		NotifyKey:    notifyKey,
		Description:  req.ProjectDesc,
		Status:       "active",
		DailyLimit:   10000,
		MonthlyLimit: 300000,
		Priority:     "normal",
		CreatedBy:    req.CreatedBy,
	}
	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// 4) Insert destination
	dest := &models.Destination{
		BaseModel:        models.BaseModel{ID: uuid.New()},
		ProjectID:        projectID,
		Name:             fmt.Sprintf("%s-default", req.ProjectName),
		Description:      "Auto-provisioned default destination",
		TeamsTenantID:    req.TeamsTenantID,
		Targets:          req.Targets,
		Status:           "active",
		ValidationStatus: "pending",
		CreatedBy:        req.CreatedBy,
	}
	if err := s.destinationRepo.Create(ctx, dest); err != nil {
		return nil, fmt.Errorf("failed to create destination: %w", err)
	}

	return &ProvisionCreateResponse{NotifyKey: notifyKey, Project: project, Destination: dest}, nil
}

func (s *provisionService) Read(ctx context.Context, notifyKey string) (*ProvisionReadResponse, error) {
	proj, err := s.projectRepo.GetByNotifyKey(ctx, notifyKey)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	dests, err := s.destinationRepo.GetByProjectID(ctx, proj.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list destinations: %w", err)
	}
	return &ProvisionReadResponse{Project: proj, Destinations: dests}, nil
}

func (s *provisionService) Update(ctx context.Context, notifyKey string, req *ProvisionUpdateRequest) (*ProvisionReadResponse, error) {
	proj, err := s.projectRepo.GetByNotifyKey(ctx, notifyKey)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	// Update project fields
	if req.ProjectDesc != nil {
		proj.Description = *req.ProjectDesc
	}
	// Name not persisted as a separate column; ignoring unless added later
	if err := s.projectRepo.Update(ctx, proj); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Update first destination if provided
	if req.TeamsTenantID != nil || req.Targets != nil {
		dests, err := s.destinationRepo.GetByProjectID(ctx, proj.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load destinations: %w", err)
		}
		if len(dests) > 0 {
			d := dests[0]
			if req.TeamsTenantID != nil {
				d.TeamsTenantID = *req.TeamsTenantID
			}
			if req.Targets != nil {
				d.Targets = *req.Targets
			}
			if err := s.destinationRepo.Update(ctx, d); err != nil {
				return nil, fmt.Errorf("failed to update destination: %w", err)
			}
		}
	}

	return s.Read(ctx, notifyKey)
}

func (s *provisionService) SetStatus(ctx context.Context, notifyKey string, enabled bool) error {
	proj, err := s.projectRepo.GetByNotifyKey(ctx, notifyKey)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}
	if enabled {
		proj.Status = "active"
	} else {
		proj.Status = "inactive"
	}
	if err := s.projectRepo.Update(ctx, proj); err != nil {
		return fmt.Errorf("failed to update project status: %w", err)
	}
	return nil
}

func generateNotifyKey(projectID, companyID uuid.UUID) string {
	sum := sha256.Sum256([]byte(projectID.String() + ":" + companyID.String()))
	// Return 64 hex chars
	return hex.EncodeToString(sum[:])
}
