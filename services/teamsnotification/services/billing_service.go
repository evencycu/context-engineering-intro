package services

import (
	"context"
	"fmt"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/libs/models"
	"github.com/evencycu/TeamsNotifyGoV2/services/teamsnotification/repositories"
	"github.com/google/uuid"
)

// BillingService defines the interface for billing operations
type BillingService interface {
	// Usage tracking
	GetUsageRecords(ctx context.Context, req *GetUsageRecordsRequest) (*UsageRecordsData, error)
	GetUsageSummary(ctx context.Context, req *GetUsageSummaryRequest) (*UsageSummaryData, error)

	// Billing plans
	GetBillingPlans(ctx context.Context) ([]*models.BillingPlan, error)
	GetBillingPlan(ctx context.Context, id uuid.UUID) (*models.BillingPlan, error)
	CreateBillingPlan(ctx context.Context, req *CreateBillingPlanRequest) (*models.BillingPlan, error)
	UpdateBillingPlan(ctx context.Context, id uuid.UUID, req *UpdateBillingPlanRequest) (*models.BillingPlan, error)

	// Project billing
	GetProjectBilling(ctx context.Context, projectID uuid.UUID) (*models.ProjectBilling, error)
	UpdateProjectBilling(ctx context.Context, projectID uuid.UUID, req *UpdateProjectBillingRequest) (*models.ProjectBilling, error)
	SetProjectBillingPlan(ctx context.Context, projectID uuid.UUID, planID uuid.UUID) (*models.ProjectBilling, error)

	// Analytics
	GetAnalyticsOverview(ctx context.Context, req *GetAnalyticsOverviewRequest) (*AnalyticsOverviewData, error)
	GetAnalyticsTrends(ctx context.Context, req *GetAnalyticsTrendsRequest) (*AnalyticsTrendsData, error)
	GetCompanyAnalytics(ctx context.Context, companyID uuid.UUID, req *GetCompanyAnalyticsRequest) (*CompanyAnalyticsData, error)

	// Usage recording
	RecordUsage(ctx context.Context, req *RecordUsageRequest) error
}

// billingService implements BillingService
type billingService struct {
	usageRepo   repositories.UsageRecordRepository
	billingRepo repositories.BillingPlanRepository
	projectRepo repositories.ProjectRepository
	userRepo    repositories.UserRepository
}

// NewBillingService creates a new billing service
func NewBillingService(
	usageRepo repositories.UsageRecordRepository,
	billingRepo repositories.BillingPlanRepository,
	projectRepo repositories.ProjectRepository,
	userRepo repositories.UserRepository,
) BillingService {
	return &billingService{
		usageRepo:   usageRepo,
		billingRepo: billingRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// Request/Response types for billing operations

type GetUsageRecordsRequest struct {
	CompanyID  *uuid.UUID
	ProjectID  *uuid.UUID
	UserID     *uuid.UUID
	RecordType string
	StartDate  *time.Time
	EndDate    *time.Time
	Page       int
	PageSize   int
	SortBy     string
	SortOrder  string
}

type UsageRecordsData struct {
	Records    []*models.UsageRecord `json:"records"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

type GetUsageSummaryRequest struct {
	CompanyID *uuid.UUID
	ProjectID *uuid.UUID
	UserID    *uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
	GroupBy   string
}

type UsageSummaryData struct {
	TotalNotifications int64            `json:"total_notifications"`
	TotalAttachments   int64            `json:"total_attachments"`
	TotalMentions      int64            `json:"total_mentions"`
	TotalAdaptiveCards int64            `json:"total_adaptive_cards"`
	TotalCost          float64          `json:"total_cost"`
	AverageCostPerDay  float64          `json:"average_cost_per_day"`
	PeakUsageDay       string           `json:"peak_usage_day"`
	PeakUsageCount     int64            `json:"peak_usage_count"`
	GroupedData        []UsageGroupData `json:"grouped_data"`
}

type UsageGroupData struct {
	Date          string  `json:"date"`
	Notifications int64   `json:"notifications"`
	Attachments   int64   `json:"attachments"`
	Mentions      int64   `json:"mentions"`
	AdaptiveCards int64   `json:"adaptive_cards"`
	Cost          float64 `json:"cost"`
}

type CreateBillingPlanRequest struct {
	Name                 string
	Description          string
	PricePerNotification float64
	PricePerAttachment   float64
	PricePerMention      float64
	PricePerAdaptiveCard float64
	DailyLimit           *int
	MonthlyLimit         *int
	Status               string
}

type UpdateBillingPlanRequest struct {
	Name                 *string
	Description          *string
	PricePerNotification *float64
	PricePerAttachment   *float64
	PricePerMention      *float64
	PricePerAdaptiveCard *float64
	DailyLimit           *int
	MonthlyLimit         *int
	Status               *string
}

type UpdateProjectBillingRequest struct {
	PaymentMethod   *string
	BillingCycle    *string
	NextBillingDate *string
	BillingStatus   *string
}

type GetAnalyticsOverviewRequest struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type AnalyticsOverviewData struct {
	TotalCompanies     int64              `json:"total_companies"`
	TotalProjects      int64              `json:"total_projects"`
	TotalNotifications int64              `json:"total_notifications"`
	TotalRevenue       float64            `json:"total_revenue"`
	ActiveUsers        int64              `json:"active_users"`
	TopCompanies       []CompanyUsageData `json:"top_companies"`
	NotificationTrends []TrendData        `json:"notification_trends"`
	RevenueTrends      []TrendData        `json:"revenue_trends"`
}

type CompanyUsageData struct {
	CompanyID     uuid.UUID `json:"company_id"`
	CompanyName   string    `json:"company_name"`
	Notifications int64     `json:"notifications"`
	Cost          float64   `json:"cost"`
}

type TrendData struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type GetAnalyticsTrendsRequest struct {
	StartDate *time.Time
	EndDate   *time.Time
	GroupBy   string
}

type AnalyticsTrendsData struct {
	NotificationTrends []TrendData `json:"notification_trends"`
	RevenueTrends      []TrendData `json:"revenue_trends"`
	UserTrends         []TrendData `json:"user_trends"`
	CostTrends         []TrendData `json:"cost_trends"`
}

type GetCompanyAnalyticsRequest struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type CompanyAnalyticsData struct {
	CompanyID          uuid.UUID          `json:"company_id"`
	CompanyName        string             `json:"company_name"`
	TotalNotifications int64              `json:"total_notifications"`
	TotalCost          float64            `json:"total_cost"`
	AverageDailyUsage  float64            `json:"average_daily_usage"`
	PeakUsageDay       string             `json:"peak_usage_day"`
	ProjectBreakdown   []ProjectUsageData `json:"project_breakdown"`
	MonthlyTrends      []TrendData        `json:"monthly_trends"`
}

type ProjectUsageData struct {
	ProjectID     uuid.UUID `json:"project_id"`
	ProjectName   string    `json:"project_name"`
	Notifications int64     `json:"notifications"`
	Cost          float64   `json:"cost"`
	Percentage    float64   `json:"percentage"`
}

type RecordUsageRequest struct {
	ProjectID  uuid.UUID
	UserID     *uuid.UUID
	RecordType string
	Quantity   int
	UnitCost   float64
	TotalCost  float64
	Metadata   map[string]any
}

// Implementation methods

func (s *billingService) GetUsageRecords(ctx context.Context, req *GetUsageRecordsRequest) (*UsageRecordsData, error) {
	// Implementation will be added when repository methods are available
	return &UsageRecordsData{
		Records:    []*models.UsageRecord{},
		Total:      0,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: 0,
	}, nil
}

func (s *billingService) GetUsageSummary(ctx context.Context, req *GetUsageSummaryRequest) (*UsageSummaryData, error) {
	// Implementation will be added when repository methods are available
	return &UsageSummaryData{
		TotalNotifications: 0,
		TotalAttachments:   0,
		TotalMentions:      0,
		TotalAdaptiveCards: 0,
		TotalCost:          0.0,
		AverageCostPerDay:  0.0,
		PeakUsageDay:       "",
		PeakUsageCount:     0,
		GroupedData:        []UsageGroupData{},
	}, nil
}

func (s *billingService) GetBillingPlans(ctx context.Context) ([]*models.BillingPlan, error) {
	// Implementation will be added when repository methods are available
	return []*models.BillingPlan{}, nil
}

func (s *billingService) GetBillingPlan(ctx context.Context, id uuid.UUID) (*models.BillingPlan, error) {
	// Implementation will be added when repository methods are available
	return nil, fmt.Errorf("not implemented")
}

func (s *billingService) CreateBillingPlan(ctx context.Context, req *CreateBillingPlanRequest) (*models.BillingPlan, error) {
	// Implementation will be added when repository methods are available
	return nil, fmt.Errorf("not implemented")
}

func (s *billingService) UpdateBillingPlan(ctx context.Context, id uuid.UUID, req *UpdateBillingPlanRequest) (*models.BillingPlan, error) {
	// Implementation will be added when repository methods are available
	return nil, fmt.Errorf("not implemented")
}

func (s *billingService) GetProjectBilling(ctx context.Context, projectID uuid.UUID) (*models.ProjectBilling, error) {
	// Implementation will be added when repository methods are available
	return nil, fmt.Errorf("not implemented")
}

func (s *billingService) UpdateProjectBilling(ctx context.Context, projectID uuid.UUID, req *UpdateProjectBillingRequest) (*models.ProjectBilling, error) {
	// Implementation will be added when repository methods are available
	return nil, fmt.Errorf("not implemented")
}

func (s *billingService) SetProjectBillingPlan(ctx context.Context, projectID uuid.UUID, planID uuid.UUID) (*models.ProjectBilling, error) {
	// Implementation will be added when repository methods are available
	return nil, fmt.Errorf("not implemented")
}

func (s *billingService) GetAnalyticsOverview(ctx context.Context, req *GetAnalyticsOverviewRequest) (*AnalyticsOverviewData, error) {
	// Implementation will be added when repository methods are available
	return &AnalyticsOverviewData{
		TotalCompanies:     0,
		TotalProjects:      0,
		TotalNotifications: 0,
		TotalRevenue:       0.0,
		ActiveUsers:        0,
		TopCompanies:       []CompanyUsageData{},
		NotificationTrends: []TrendData{},
		RevenueTrends:      []TrendData{},
	}, nil
}

func (s *billingService) GetAnalyticsTrends(ctx context.Context, req *GetAnalyticsTrendsRequest) (*AnalyticsTrendsData, error) {
	// Implementation will be added when repository methods are available
	return &AnalyticsTrendsData{
		NotificationTrends: []TrendData{},
		RevenueTrends:      []TrendData{},
		UserTrends:         []TrendData{},
		CostTrends:         []TrendData{},
	}, nil
}

func (s *billingService) GetCompanyAnalytics(ctx context.Context, companyID uuid.UUID, req *GetCompanyAnalyticsRequest) (*CompanyAnalyticsData, error) {
	// Implementation will be added when repository methods are available
	return &CompanyAnalyticsData{
		CompanyID:          companyID,
		CompanyName:        "",
		TotalNotifications: 0,
		TotalCost:          0.0,
		AverageDailyUsage:  0.0,
		PeakUsageDay:       "",
		ProjectBreakdown:   []ProjectUsageData{},
		MonthlyTrends:      []TrendData{},
	}, nil
}

func (s *billingService) RecordUsage(ctx context.Context, req *RecordUsageRequest) error {
	// Validate required fields
	if req.ProjectID == uuid.Nil {
		return fmt.Errorf("project_id is required")
	}
	if req.RecordType == "" {
		return fmt.Errorf("record_type is required")
	}

	// Create usage record
	usageRecord := &models.UsageRecord{
		BaseModel: models.BaseModel{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ProjectID:  req.ProjectID,
		UserID:     req.UserID,
		RecordType: req.RecordType,
		Quantity:   req.Quantity,
		UnitCost:   req.UnitCost,
		TotalCost:  req.TotalCost,
		Metadata:   models.JSONBObject(req.Metadata),
	}

	// Save to models
	if err := s.usageRepo.Create(ctx, usageRecord); err != nil {
		return fmt.Errorf("failed to create usage record: %w", err)
	}

	return nil
}
