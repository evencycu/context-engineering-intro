package repositories

import (
	"context"
	"time"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BillingPlanRepository defines the interface for billing plan operations
type BillingPlanRepository interface {
	Create(ctx context.Context, entity *models.BillingPlan) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.BillingPlan, error)
	GetAll(ctx context.Context) ([]*models.BillingPlan, error)
	Update(ctx context.Context, entity *models.BillingPlan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// UsageRecordRepository defines the interface for usage record operations
type UsageRecordRepository interface {
	Create(ctx context.Context, entity *models.UsageRecord) error
	CreateBatch(ctx context.Context, entities []*models.UsageRecord) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.UsageRecord, error)
	GetByCompanyID(ctx context.Context, companyID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*models.UsageRecord, int64, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*models.UsageRecord, int64, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*models.UsageRecord, int64, error)
	GetSummary(ctx context.Context, companyID *uuid.UUID, projectID *uuid.UUID, userID *uuid.UUID, startDate, endDate *time.Time) (*UsageSummary, error)
	GetTrends(ctx context.Context, companyID *uuid.UUID, projectID *uuid.UUID, userID *uuid.UUID, startDate, endDate *time.Time, groupBy string) ([]*UsageTrend, error)
}

// ProjectBillingRepository defines the interface for project billing operations
type ProjectBillingRepository interface {
	Create(ctx context.Context, entity *models.ProjectBilling) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*models.ProjectBilling, error)
	Update(ctx context.Context, entity *models.ProjectBilling) error
	Delete(ctx context.Context, projectID uuid.UUID) error
	SetBillingPlan(ctx context.Context, projectID uuid.UUID, planID uuid.UUID) error
}

// UsageSummary represents usage summary data
type UsageSummary struct {
	TotalNotifications int64   `json:"total_notifications"`
	TotalAttachments   int64   `json:"total_attachments"`
	TotalMentions      int64   `json:"total_mentions"`
	TotalAdaptiveCards int64   `json:"total_adaptive_cards"`
	TotalCost          float64 `json:"total_cost"`
	AverageCostPerDay  float64 `json:"average_cost_per_day"`
	PeakUsageDay       string  `json:"peak_usage_day"`
	PeakUsageCount     int64   `json:"peak_usage_count"`
}

// UsageTrend represents usage trend data
type UsageTrend struct {
	Date          time.Time `json:"date"`
	Notifications int64     `json:"notifications"`
	Attachments   int64     `json:"attachments"`
	Mentions      int64     `json:"mentions"`
	AdaptiveCards int64     `json:"adaptive_cards"`
	Cost          float64   `json:"cost"`
}

// billingPlanRepository implements BillingPlanRepository
type billingPlanRepository struct {
	db *sqlx.DB
}

// NewBillingPlanRepository creates a new billing plan repository
func NewBillingPlanRepository(db *sqlx.DB) BillingPlanRepository {
	return &billingPlanRepository{db: db}
}

func (r *billingPlanRepository) Create(ctx context.Context, entity *models.BillingPlan) error {
	query := `INSERT INTO billing_plans (
		id, name, description, price_per_notification, price_per_attachment,
		price_per_mention, price_per_adaptive_card, daily_limit, monthly_limit,
		status, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
	)`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.Name, entity.Description, entity.PricePerNotification,
		entity.PricePerAttachment, entity.PricePerMention, entity.PricePerAdaptiveCard,
		entity.DailyLimit, entity.MonthlyLimit, entity.Status, entity.CreatedAt, entity.UpdatedAt,
	)
	return err
}

func (r *billingPlanRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.BillingPlan, error) {
	query := `SELECT id, name, description, price_per_notification, price_per_attachment,
		price_per_mention, price_per_adaptive_card, daily_limit, monthly_limit, status,
		created_at, updated_at
		FROM billing_plans WHERE id = $1`

	var plan models.BillingPlan
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID, &plan.Name, &plan.Description, &plan.PricePerNotification,
		&plan.PricePerAttachment, &plan.PricePerMention, &plan.PricePerAdaptiveCard,
		&plan.DailyLimit, &plan.MonthlyLimit, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *billingPlanRepository) GetAll(ctx context.Context) ([]*models.BillingPlan, error) {
	query := `SELECT id, name, description, price_per_notification, price_per_attachment,
		price_per_mention, price_per_adaptive_card, daily_limit, monthly_limit, status,
		created_at, updated_at
		FROM billing_plans ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*models.BillingPlan
	for rows.Next() {
		var plan models.BillingPlan
		err := rows.Scan(
			&plan.ID, &plan.Name, &plan.Description, &plan.PricePerNotification,
			&plan.PricePerAttachment, &plan.PricePerMention, &plan.PricePerAdaptiveCard,
			&plan.DailyLimit, &plan.MonthlyLimit, &plan.Status, &plan.CreatedAt, &plan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, &plan)
	}
	return plans, nil
}

func (r *billingPlanRepository) Update(ctx context.Context, entity *models.BillingPlan) error {
	query := `UPDATE billing_plans SET
		name = $2, description = $3, price_per_notification = $4, price_per_attachment = $5,
		price_per_mention = $6, price_per_adaptive_card = $7, daily_limit = $8,
		monthly_limit = $9, status = $10, updated_at = $11
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.Name, entity.Description, entity.PricePerNotification,
		entity.PricePerAttachment, entity.PricePerMention, entity.PricePerAdaptiveCard,
		entity.DailyLimit, entity.MonthlyLimit, entity.Status, entity.UpdatedAt,
	)
	return err
}

func (r *billingPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM billing_plans WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// usageRecordRepository implements UsageRecordRepository
type usageRecordRepository struct {
	db *sqlx.DB
}

// NewUsageRecordRepository creates a new usage record repository
func NewUsageRecordRepository(db *sqlx.DB) UsageRecordRepository {
	return &usageRecordRepository{db: db}
}

func (r *usageRecordRepository) Create(ctx context.Context, entity *models.UsageRecord) error {
	query := `INSERT INTO usage_records (
		id, project_id, user_id, record_type, quantity, unit_cost, total_cost, metadata, created_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	)`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.UserID, entity.RecordType, entity.Quantity,
		entity.UnitCost, entity.TotalCost, entity.Metadata, entity.CreatedAt,
	)
	return err
}

func (r *usageRecordRepository) CreateBatch(ctx context.Context, entities []*models.UsageRecord) error {
	if len(entities) == 0 {
		return nil
	}

	query := `INSERT INTO usage_records (
		id, project_id, user_id, record_type, quantity, unit_cost, total_cost, metadata, created_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	)`

	for _, entity := range entities {
		_, err := r.db.ExecContext(ctx, query,
			entity.ID, entity.ProjectID, entity.UserID, entity.RecordType, entity.Quantity,
			entity.UnitCost, entity.TotalCost, entity.Metadata, entity.CreatedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *usageRecordRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UsageRecord, error) {
	query := `SELECT id, project_id, user_id, record_type, quantity, unit_cost, total_cost, metadata, created_at
		FROM usage_records WHERE id = $1`

	var record models.UsageRecord
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID, &record.ProjectID, &record.UserID, &record.RecordType, &record.Quantity,
		&record.UnitCost, &record.TotalCost, &record.Metadata, &record.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *usageRecordRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*models.UsageRecord, int64, error) {
	// Implementation will be added
	return []*models.UsageRecord{}, 0, nil
}

func (r *usageRecordRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*models.UsageRecord, int64, error) {
	// Implementation will be added
	return []*models.UsageRecord{}, 0, nil
}

func (r *usageRecordRepository) GetByUserID(ctx context.Context, userID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*models.UsageRecord, int64, error) {
	// Implementation will be added
	return []*models.UsageRecord{}, 0, nil
}

func (r *usageRecordRepository) GetSummary(ctx context.Context, companyID *uuid.UUID, projectID *uuid.UUID, userID *uuid.UUID, startDate, endDate *time.Time) (*UsageSummary, error) {
	// Implementation will be added
	return &UsageSummary{}, nil
}

func (r *usageRecordRepository) GetTrends(ctx context.Context, companyID *uuid.UUID, projectID *uuid.UUID, userID *uuid.UUID, startDate, endDate *time.Time, groupBy string) ([]*UsageTrend, error) {
	// Implementation will be added
	return []*UsageTrend{}, nil
}

// projectBillingRepository implements ProjectBillingRepository
type projectBillingRepository struct {
	db *sqlx.DB
}

// NewProjectBillingRepository creates a new project billing repository
func NewProjectBillingRepository(db *sqlx.DB) ProjectBillingRepository {
	return &projectBillingRepository{db: db}
}

func (r *projectBillingRepository) Create(ctx context.Context, entity *models.ProjectBilling) error {
	query := `INSERT INTO project_billing (
		id, project_id, billing_plan_id, billing_status, payment_method,
		billing_cycle, next_billing_date, total_usage_cost, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	)`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.ProjectID, entity.BillingPlanID, entity.BillingStatus,
		entity.PaymentMethod, entity.BillingCycle, entity.NextBillingDate, entity.TotalUsageCost,
		entity.CreatedAt, entity.UpdatedAt,
	)
	return err
}

func (r *projectBillingRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*models.ProjectBilling, error) {
	query := `SELECT id, project_id, billing_plan_id, billing_status, payment_method,
		billing_cycle, next_billing_date, total_usage_cost, created_at, updated_at
		FROM project_billing WHERE project_id = $1`

	var billing models.ProjectBilling
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&billing.ID, &billing.ProjectID, &billing.BillingPlanID, &billing.BillingStatus,
		&billing.PaymentMethod, &billing.BillingCycle, &billing.NextBillingDate, &billing.TotalUsageCost,
		&billing.CreatedAt, &billing.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *projectBillingRepository) Update(ctx context.Context, entity *models.ProjectBilling) error {
	query := `UPDATE project_billing SET
		billing_plan_id = $2, billing_status = $3, payment_method = $4,
		billing_cycle = $5, next_billing_date = $6, total_usage_cost = $7, updated_at = $8
		WHERE project_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		entity.ProjectID, entity.BillingPlanID, entity.BillingStatus, entity.PaymentMethod,
		entity.BillingCycle, entity.NextBillingDate, entity.TotalUsageCost, entity.UpdatedAt,
	)
	return err
}

func (r *projectBillingRepository) Delete(ctx context.Context, projectID uuid.UUID) error {
	query := `DELETE FROM project_billing WHERE project_id = $1`
	_, err := r.db.ExecContext(ctx, query, projectID)
	return err
}

func (r *projectBillingRepository) SetBillingPlan(ctx context.Context, projectID uuid.UUID, planID uuid.UUID) error {
	query := `UPDATE project_billing SET billing_plan_id = $2, updated_at = $3 WHERE project_id = $1`
	_, err := r.db.ExecContext(ctx, query, projectID, planID, time.Now())
	return err
}
