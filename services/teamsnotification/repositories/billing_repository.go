package repositories

import (
	"context"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/libs/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BillingPlanRepository defines the interface for billing plan operations
type BillingPlanRepository interface {
	Create(ctx context.Context, entity *database.BillingPlan) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.BillingPlan, error)
	GetAll(ctx context.Context) ([]*database.BillingPlan, error)
	Update(ctx context.Context, entity *database.BillingPlan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// UsageRecordRepository defines the interface for usage record operations
type UsageRecordRepository interface {
	Create(ctx context.Context, entity *database.UsageRecord) error
	CreateBatch(ctx context.Context, entities []*database.UsageRecord) error
	GetByID(ctx context.Context, id uuid.UUID) (*database.UsageRecord, error)
	GetByCompanyID(ctx context.Context, companyID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*database.UsageRecord, int64, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*database.UsageRecord, int64, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*database.UsageRecord, int64, error)
	GetSummary(ctx context.Context, companyID *uuid.UUID, projectID *uuid.UUID, userID *uuid.UUID, startDate, endDate *time.Time) (*UsageSummary, error)
	GetTrends(ctx context.Context, companyID *uuid.UUID, projectID *uuid.UUID, userID *uuid.UUID, startDate, endDate *time.Time, groupBy string) ([]*UsageTrend, error)
}

// CompanyBillingRepository defines the interface for company billing operations
type CompanyBillingRepository interface {
	Create(ctx context.Context, entity *database.CompanyBilling) error
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*database.CompanyBilling, error)
	Update(ctx context.Context, entity *database.CompanyBilling) error
	Delete(ctx context.Context, companyID uuid.UUID) error
	SetBillingPlan(ctx context.Context, companyID uuid.UUID, planID uuid.UUID) error
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

func (r *billingPlanRepository) Create(ctx context.Context, entity *database.BillingPlan) error {
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

func (r *billingPlanRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.BillingPlan, error) {
	query := `SELECT id, name, description, price_per_notification, price_per_attachment,
		price_per_mention, price_per_adaptive_card, daily_limit, monthly_limit, status,
		created_at, updated_at
		FROM billing_plans WHERE id = $1`

	var plan database.BillingPlan
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

func (r *billingPlanRepository) GetAll(ctx context.Context) ([]*database.BillingPlan, error) {
	query := `SELECT id, name, description, price_per_notification, price_per_attachment,
		price_per_mention, price_per_adaptive_card, daily_limit, monthly_limit, status,
		created_at, updated_at
		FROM billing_plans ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*database.BillingPlan
	for rows.Next() {
		var plan database.BillingPlan
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

func (r *billingPlanRepository) Update(ctx context.Context, entity *database.BillingPlan) error {
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

func (r *usageRecordRepository) Create(ctx context.Context, entity *database.UsageRecord) error {
	query := `INSERT INTO usage_records (
		id, company_id, project_id, user_id, notification_id, bot_id, bot_type,
		record_type, quantity, unit_price, total_cost, billing_period, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
	)`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.ProjectID, entity.UserID, entity.NotificationID,
		entity.BotID, entity.BotType, entity.RecordType, entity.Quantity,
		entity.UnitPrice, entity.TotalCost, entity.BillingPeriod, entity.CreatedAt, entity.UpdatedAt,
	)
	return err
}

func (r *usageRecordRepository) CreateBatch(ctx context.Context, entities []*database.UsageRecord) error {
	if len(entities) == 0 {
		return nil
	}

	query := `INSERT INTO usage_records (
		id, company_id, project_id, user_id, notification_id, bot_id, bot_type,
		record_type, quantity, unit_price, total_cost, billing_period, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
	)`

	for _, entity := range entities {
		_, err := r.db.ExecContext(ctx, query,
			entity.ID, entity.CompanyID, entity.ProjectID, entity.UserID, entity.NotificationID,
			entity.BotID, entity.BotType, entity.RecordType, entity.Quantity,
			entity.UnitPrice, entity.TotalCost, entity.BillingPeriod, entity.CreatedAt, entity.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *usageRecordRepository) GetByID(ctx context.Context, id uuid.UUID) (*database.UsageRecord, error) {
	query := `SELECT id, company_id, project_id, user_id, notification_id, bot_id, bot_type,
		record_type, quantity, unit_price, total_cost, billing_period, created_at, updated_at
		FROM usage_records WHERE id = $1`

	var record database.UsageRecord
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&record.ID, &record.CompanyID, &record.ProjectID, &record.UserID, &record.NotificationID,
		&record.BotID, &record.BotType, &record.RecordType, &record.Quantity,
		&record.UnitPrice, &record.TotalCost, &record.BillingPeriod, &record.CreatedAt, &record.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *usageRecordRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*database.UsageRecord, int64, error) {
	// Implementation will be added
	return []*database.UsageRecord{}, 0, nil
}

func (r *usageRecordRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*database.UsageRecord, int64, error) {
	// Implementation will be added
	return []*database.UsageRecord{}, 0, nil
}

func (r *usageRecordRepository) GetByUserID(ctx context.Context, userID uuid.UUID, startDate, endDate *time.Time, page, pageSize int) ([]*database.UsageRecord, int64, error) {
	// Implementation will be added
	return []*database.UsageRecord{}, 0, nil
}

func (r *usageRecordRepository) GetSummary(ctx context.Context, companyID *uuid.UUID, projectID *uuid.UUID, userID *uuid.UUID, startDate, endDate *time.Time) (*UsageSummary, error) {
	// Implementation will be added
	return &UsageSummary{}, nil
}

func (r *usageRecordRepository) GetTrends(ctx context.Context, companyID *uuid.UUID, projectID *uuid.UUID, userID *uuid.UUID, startDate, endDate *time.Time, groupBy string) ([]*UsageTrend, error) {
	// Implementation will be added
	return []*UsageTrend{}, nil
}

// companyBillingRepository implements CompanyBillingRepository
type companyBillingRepository struct {
	db *sqlx.DB
}

// NewCompanyBillingRepository creates a new company billing repository
func NewCompanyBillingRepository(db *sqlx.DB) CompanyBillingRepository {
	return &companyBillingRepository{db: db}
}

func (r *companyBillingRepository) Create(ctx context.Context, entity *database.CompanyBilling) error {
	query := `INSERT INTO company_billing (
		id, company_id, billing_plan_id, status, billing_email, payment_method,
		currency, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	)`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID, entity.CompanyID, entity.BillingPlanID, entity.Status,
		entity.BillingEmail, entity.PaymentMethod, entity.Currency, entity.CreatedAt, entity.UpdatedAt,
	)
	return err
}

func (r *companyBillingRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*database.CompanyBilling, error) {
	query := `SELECT id, company_id, billing_plan_id, status, billing_email, payment_method,
		currency, created_at, updated_at
		FROM company_billing WHERE company_id = $1`

	var billing database.CompanyBilling
	err := r.db.QueryRowContext(ctx, query, companyID).Scan(
		&billing.ID, &billing.CompanyID, &billing.BillingPlanID, &billing.Status,
		&billing.BillingEmail, &billing.PaymentMethod, &billing.Currency, &billing.CreatedAt, &billing.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *companyBillingRepository) Update(ctx context.Context, entity *database.CompanyBilling) error {
	query := `UPDATE company_billing SET
		billing_plan_id = $2, status = $3, billing_email = $4, payment_method = $5,
		currency = $6, updated_at = $7
		WHERE company_id = $1`

	_, err := r.db.ExecContext(ctx, query,
		entity.CompanyID, entity.BillingPlanID, entity.Status, entity.BillingEmail,
		entity.PaymentMethod, entity.Currency, entity.UpdatedAt,
	)
	return err
}

func (r *companyBillingRepository) Delete(ctx context.Context, companyID uuid.UUID) error {
	query := `DELETE FROM company_billing WHERE company_id = $1`
	_, err := r.db.ExecContext(ctx, query, companyID)
	return err
}

func (r *companyBillingRepository) SetBillingPlan(ctx context.Context, companyID uuid.UUID, planID uuid.UUID) error {
	query := `UPDATE company_billing SET billing_plan_id = $2, updated_at = $3 WHERE company_id = $1`
	_, err := r.db.ExecContext(ctx, query, companyID, planID, time.Now())
	return err
}
