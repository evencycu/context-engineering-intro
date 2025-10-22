package billing

import (
	"net/http"
	"strconv"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles billing-related HTTP requests
type Handler struct {
	billingService services.BillingService
}

// NewHandler creates a new billing handler
func NewHandler(billingService services.BillingService) *Handler {
	return &Handler{
		billingService: billingService,
	}
}

// RegisterRoutes registers billing routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	billing := rg.Group("/billing")
	{
		// Usage tracking
		billing.GET("/usage", h.GetUsageRecords)
		billing.GET("/usage/summary", h.GetUsageSummary)
		billing.GET("/usage/company/:companyId", h.GetCompanyUsage)
		billing.GET("/usage/project/:projectId", h.GetProjectUsage)

		// Billing plans
		billing.GET("/plans", h.GetBillingPlans)
		billing.GET("/plans/:id", h.GetBillingPlan)
		billing.POST("/plans", h.CreateBillingPlan)
		billing.PUT("/plans/:id", h.UpdateBillingPlan)

		// Project billing
		billing.GET("/project/:projectId", h.GetProjectBilling)
		billing.PUT("/project/:projectId", h.UpdateProjectBilling)
		billing.POST("/project/:projectId/plan", h.SetProjectBillingPlan)

		// Analytics
		billing.GET("/analytics/overview", h.GetAnalyticsOverview)
		billing.GET("/analytics/trends", h.GetAnalyticsTrends)
		billing.GET("/analytics/company/:companyId", h.GetCompanyAnalytics)
	}
}

// GetUsageRecordsRequest represents a get usage records request
type GetUsageRecordsRequest struct {
	CompanyID  *uuid.UUID `form:"company_id"`
	ProjectID  *uuid.UUID `form:"project_id"`
	UserID     *uuid.UUID `form:"user_id"`
	RecordType string     `form:"record_type"`
	StartDate  *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate    *time.Time `form:"end_date" time_format:"2006-01-02"`
	Page       int        `form:"page,default=1"`
	PageSize   int        `form:"page_size,default=50"`
	SortBy     string     `form:"sort_by,default=created_at"`
	SortOrder  string     `form:"sort_order,default=desc"`
}

// GetUsageRecordsResponse represents the response for get usage records
type GetUsageRecordsResponse struct {
	Success bool                       `json:"success"`
	Data    *services.UsageRecordsData `json:"data,omitempty"`
	Error   string                     `json:"error,omitempty"`
}

// GetUsageRecords gets usage records with filtering and pagination
func (h *Handler) GetUsageRecords(c *gin.Context) {
	var req GetUsageRecordsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetUsageRecordsResponse{
			Success: false,
			Error:   "Invalid query parameters: " + err.Error(),
		})
		return
	}

	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 50
	}
	if req.PageSize > 1000 {
		req.PageSize = 1000
	}

	data, err := h.billingService.GetUsageRecords(c.Request.Context(), &services.GetUsageRecordsRequest{
		CompanyID:  req.CompanyID,
		ProjectID:  req.ProjectID,
		UserID:     req.UserID,
		RecordType: req.RecordType,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Page:       req.Page,
		PageSize:   req.PageSize,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetUsageRecordsResponse{
			Success: false,
			Error:   "Failed to get usage records: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetUsageRecordsResponse{
		Success: true,
		Data:    data,
	})
}

// GetUsageSummaryRequest represents a get usage summary request
type GetUsageSummaryRequest struct {
	CompanyID *uuid.UUID `form:"company_id"`
	ProjectID *uuid.UUID `form:"project_id"`
	UserID    *uuid.UUID `form:"user_id"`
	StartDate *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"end_date" time_format:"2006-01-02"`
	GroupBy   string     `form:"group_by,default=day"` // day, week, month
}

// GetUsageSummaryResponse represents the response for get usage summary
type GetUsageSummaryResponse struct {
	Success bool                       `json:"success"`
	Data    *services.UsageSummaryData `json:"data,omitempty"`
	Error   string                     `json:"error,omitempty"`
}

// GetUsageSummary gets usage summary statistics
func (h *Handler) GetUsageSummary(c *gin.Context) {
	var req GetUsageSummaryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetUsageSummaryResponse{
			Success: false,
			Error:   "Invalid query parameters: " + err.Error(),
		})
		return
	}

	// Set default date range if not provided
	if req.StartDate == nil {
		now := time.Now()
		startDate := now.AddDate(0, 0, -30) // Last 30 days
		req.StartDate = &startDate
	}
	if req.EndDate == nil {
		now := time.Now()
		req.EndDate = &now
	}

	data, err := h.billingService.GetUsageSummary(c.Request.Context(), &services.GetUsageSummaryRequest{
		CompanyID: req.CompanyID,
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		GroupBy:   req.GroupBy,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetUsageSummaryResponse{
			Success: false,
			Error:   "Failed to get usage summary: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetUsageSummaryResponse{
		Success: true,
		Data:    data,
	})
}

// GetCompanyUsage gets usage records for a specific company
func (h *Handler) GetCompanyUsage(c *gin.Context) {
	companyIDStr := c.Param("companyId")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetUsageRecordsResponse{
			Success: false,
			Error:   "Invalid company ID format",
		})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	var startDate, endDate *time.Time
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		}
	}

	data, err := h.billingService.GetUsageRecords(c.Request.Context(), &services.GetUsageRecordsRequest{
		CompanyID: &companyID,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetUsageRecordsResponse{
			Success: false,
			Error:   "Failed to get company usage: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetUsageRecordsResponse{
		Success: true,
		Data:    data,
	})
}

// GetProjectUsage gets usage records for a specific project
func (h *Handler) GetProjectUsage(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetUsageRecordsResponse{
			Success: false,
			Error:   "Invalid project ID format",
		})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	var startDate, endDate *time.Time
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		}
	}

	data, err := h.billingService.GetUsageRecords(c.Request.Context(), &services.GetUsageRecordsRequest{
		ProjectID: &projectID,
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetUsageRecordsResponse{
			Success: false,
			Error:   "Failed to get project usage: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetUsageRecordsResponse{
		Success: true,
		Data:    data,
	})
}

// GetBillingPlans gets all billing plans
func (h *Handler) GetBillingPlans(c *gin.Context) {
	plans, err := h.billingService.GetBillingPlans(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get billing plans: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plans,
	})
}

// GetBillingPlan gets a specific billing plan
func (h *Handler) GetBillingPlan(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid plan ID format",
		})
		return
	}

	plan, err := h.billingService.GetBillingPlan(c.Request.Context(), planID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Billing plan not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plan,
	})
}

// CreateBillingPlanRequest represents a create billing plan request
type CreateBillingPlanRequest struct {
	Name                 string  `json:"name" validate:"required,min=3,max=100"`
	Description          string  `json:"description" validate:"required,min=10,max=500"`
	PricePerNotification float64 `json:"pricePerNotification" validate:"required,min=0"`
	PricePerAttachment   float64 `json:"pricePerAttachment" validate:"required,min=0"`
	PricePerMention      float64 `json:"pricePerMention" validate:"required,min=0"`
	PricePerAdaptiveCard float64 `json:"pricePerAdaptiveCard" validate:"required,min=0"`
	DailyLimit           *int    `json:"dailyLimit" validate:"omitempty,min=1,max=100000"`
	MonthlyLimit         *int    `json:"monthlyLimit" validate:"omitempty,min=1,max=1000000"`
	Status               string  `json:"status" validate:"required,oneof=active inactive"`
}

// CreateBillingPlan creates a new billing plan
func (h *Handler) CreateBillingPlan(c *gin.Context) {
	var req CreateBillingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	plan, err := h.billingService.CreateBillingPlan(c.Request.Context(), &services.CreateBillingPlanRequest{
		Name:                 req.Name,
		Description:          req.Description,
		PricePerNotification: req.PricePerNotification,
		PricePerAttachment:   req.PricePerAttachment,
		PricePerMention:      req.PricePerMention,
		PricePerAdaptiveCard: req.PricePerAdaptiveCard,
		DailyLimit:           req.DailyLimit,
		MonthlyLimit:         req.MonthlyLimit,
		Status:               req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create billing plan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    plan,
	})
}

// UpdateBillingPlanRequest represents an update billing plan request
type UpdateBillingPlanRequest struct {
	Name                 *string  `json:"name" validate:"omitempty,min=3,max=100"`
	Description          *string  `json:"description" validate:"omitempty,min=10,max=500"`
	PricePerNotification *float64 `json:"pricePerNotification" validate:"omitempty,min=0"`
	PricePerAttachment   *float64 `json:"pricePerAttachment" validate:"omitempty,min=0"`
	PricePerMention      *float64 `json:"pricePerMention" validate:"omitempty,min=0"`
	PricePerAdaptiveCard *float64 `json:"pricePerAdaptiveCard" validate:"omitempty,min=0"`
	DailyLimit           *int     `json:"dailyLimit" validate:"omitempty,min=1,max=100000"`
	MonthlyLimit         *int     `json:"monthlyLimit" validate:"omitempty,min=1,max=1000000"`
	Status               *string  `json:"status" validate:"omitempty,oneof=active inactive"`
}

// UpdateBillingPlan updates a billing plan
func (h *Handler) UpdateBillingPlan(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid plan ID format",
		})
		return
	}

	var req UpdateBillingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	plan, err := h.billingService.UpdateBillingPlan(c.Request.Context(), planID, &services.UpdateBillingPlanRequest{
		Name:                 req.Name,
		Description:          req.Description,
		PricePerNotification: req.PricePerNotification,
		PricePerAttachment:   req.PricePerAttachment,
		PricePerMention:      req.PricePerMention,
		PricePerAdaptiveCard: req.PricePerAdaptiveCard,
		DailyLimit:           req.DailyLimit,
		MonthlyLimit:         req.MonthlyLimit,
		Status:               req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update billing plan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plan,
	})
}

// GetProjectBilling gets billing information for a project
func (h *Handler) GetProjectBilling(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid project ID format",
		})
		return
	}

	billing, err := h.billingService.GetProjectBilling(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Project billing not found: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    billing,
	})
}

// UpdateProjectBillingRequest represents an update project billing request
type UpdateProjectBillingRequest struct {
	PaymentMethod    *string `json:"paymentMethod" validate:"omitempty,oneof=credit_card bank_transfer invoice"`
	BillingCycle     *string `json:"billingCycle" validate:"omitempty,oneof=monthly yearly"`
	NextBillingDate *string `json:"nextBillingDate" validate:"omitempty,datetime=2006-01-02"`
	BillingStatus    *string `json:"billingStatus" validate:"omitempty,oneof=active suspended cancelled"`
}

// UpdateProjectBilling updates project billing information
func (h *Handler) UpdateProjectBilling(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid project ID format",
		})
		return
	}

	var req UpdateProjectBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	billing, err := h.billingService.UpdateProjectBilling(c.Request.Context(), projectID, &services.UpdateProjectBillingRequest{
		PaymentMethod:    req.PaymentMethod,
		BillingCycle:     req.BillingCycle,
		NextBillingDate: req.NextBillingDate,
		BillingStatus:    req.BillingStatus,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update project billing: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    billing,
	})
}

// SetProjectBillingPlanRequest represents a set project billing plan request
type SetProjectBillingPlanRequest struct {
	BillingPlanID uuid.UUID `json:"billingPlanId" validate:"required"`
}

// SetProjectBillingPlan sets the billing plan for a project
func (h *Handler) SetProjectBillingPlan(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid project ID format",
		})
		return
	}

	var req SetProjectBillingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	billing, err := h.billingService.SetProjectBillingPlan(c.Request.Context(), projectID, req.BillingPlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to set project billing plan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    billing,
	})
}

// GetAnalyticsOverview gets analytics overview
func (h *Handler) GetAnalyticsOverview(c *gin.Context) {
	var startDate, endDate *time.Time
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		}
	}

	overview, err := h.billingService.GetAnalyticsOverview(c.Request.Context(), &services.GetAnalyticsOverviewRequest{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get analytics overview: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    overview,
	})
}

// GetAnalyticsTrends gets analytics trends
func (h *Handler) GetAnalyticsTrends(c *gin.Context) {
	var startDate, endDate *time.Time
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		}
	}

	groupBy := c.DefaultQuery("group_by", "day")
	trends, err := h.billingService.GetAnalyticsTrends(c.Request.Context(), &services.GetAnalyticsTrendsRequest{
		StartDate: startDate,
		EndDate:   endDate,
		GroupBy:   groupBy,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get analytics trends: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    trends,
	})
}

// GetCompanyAnalytics gets analytics for a specific company
func (h *Handler) GetCompanyAnalytics(c *gin.Context) {
	companyIDStr := c.Param("companyId")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid company ID format",
		})
		return
	}

	var startDate, endDate *time.Time
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &t
		}
	}

	analytics, err := h.billingService.GetCompanyAnalytics(c.Request.Context(), companyID, &services.GetCompanyAnalyticsRequest{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get company analytics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analytics,
	})
}
