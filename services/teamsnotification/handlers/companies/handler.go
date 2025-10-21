package companies

import (
	"net/http"
	"strconv"

	database "github.com/evencycu/TeamsNotifyGoV2/libs/models"
	"github.com/evencycu/TeamsNotifyGoV2/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles company-related HTTP requests
type Handler struct {
	companyService services.CompanyService
}

// NewHandler creates a new company handler
func NewHandler(companyService services.CompanyService) *Handler {
	return &Handler{
		companyService: companyService,
	}
}

// RegisterRoutes registers company routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	companies := rg.Group("/companies")
	{
		companies.POST("", h.CreateCompany)
		companies.GET("", h.ListCompanies)
		companies.GET("/:id", h.GetCompany)
		companies.PUT("/:id", h.UpdateCompany)
		companies.DELETE("/:id", h.DeleteCompany)
		companies.PATCH("/:id/status", h.UpdateCompanyStatus)
		companies.PATCH("/:id/billing", h.UpdateBillingStatus)
	}
}

// CreateCompanyRequest represents a create company request
type CreateCompanyRequest struct {
	Name           string `json:"name" validate:"required,min=2,max=255"`
	ContactEmail   string `json:"contactEmail" validate:"required,email"`
	ContactPhone   string `json:"contactPhone" validate:"required,min=10,max=20"`
	Address        string `json:"address" validate:"required,min=10,max=500"`
	BillingEnabled bool   `json:"billingEnabled"`
}

// UpdateCompanyRequest represents an update company request
type UpdateCompanyRequest struct {
	Name           string `json:"name" validate:"omitempty,min=2,max=255"`
	ContactEmail   string `json:"contactEmail" validate:"omitempty,email"`
	ContactPhone   string `json:"contactPhone" validate:"omitempty,min=10,max=20"`
	Address        string `json:"address" validate:"omitempty,min=10,max=500"`
	BillingEnabled *bool  `json:"billingEnabled"`
}

// UpdateStatusRequest represents an update status request
type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended"`
}

// UpdateBillingRequest represents an update billing request
type UpdateBillingRequest struct {
	Enabled bool `json:"enabled"`
}

// CreateCompany creates a new company
func (h *Handler) CreateCompany(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert request to company model
	company := &database.Company{
		Name:           req.Name,
		ContactEmail:   req.ContactEmail,
		ContactPhone:   req.ContactPhone,
		Address:        req.Address,
		BillingEnabled: req.BillingEnabled,
		Status:         "active", // Default status
	}

	// Create company
	createReq := &services.CreateRequest[database.Company]{
		Data: *company,
	}

	createdCompany, err := h.companyService.Create(c.Request.Context(), createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create company",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    createdCompany,
		"message": "Company created successfully",
	})
}

// ListCompanies lists companies with pagination
func (h *Handler) ListCompanies(c *gin.Context) {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	search := c.Query("search")
	status := c.Query("status")

	// Create list request
	req := &services.ListRequest{
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
		Search: search,
		Status: status,
	}

	// Get companies
	companies, err := h.companyService.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list companies",
			"details": err.Error(),
		})
		return
	}

	// Get total count
	countReq := &services.CountRequest{
		Search: search,
		Status: status,
	}
	total, err := h.companyService.Count(c.Request.Context(), countReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to count companies",
			"details": err.Error(),
		})
		return
	}

	// Calculate pages
	pages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data": companies,
		"pagination": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"pages":  pages,
		},
	})
}

// GetCompany gets a company by ID
func (h *Handler) GetCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	company, err := h.companyService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Company not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get company",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": company,
	})
}

// UpdateCompany updates a company
func (h *Handler) UpdateCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get existing company
	existingCompany, err := h.companyService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Company not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get company",
			"details": err.Error(),
		})
		return
	}

	// Update fields
	if req.Name != "" {
		existingCompany.Name = req.Name
	}
	if req.ContactEmail != "" {
		existingCompany.ContactEmail = req.ContactEmail
	}
	if req.ContactPhone != "" {
		existingCompany.ContactPhone = req.ContactPhone
	}
	if req.Address != "" {
		existingCompany.Address = req.Address
	}
	if req.BillingEnabled != nil {
		existingCompany.BillingEnabled = *req.BillingEnabled
	}

	// Update company
	updateReq := &services.UpdateRequest[database.Company]{
		Data: *existingCompany,
	}

	updatedCompany, err := h.companyService.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update company",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    updatedCompany,
		"message": "Company updated successfully",
	})
}

// DeleteCompany deletes a company
func (h *Handler) DeleteCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	err = h.companyService.Delete(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Company not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete company",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company deleted successfully",
	})
}

// UpdateCompanyStatus updates company status
func (h *Handler) UpdateCompanyStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	err = h.companyService.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Company not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update company status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company status updated successfully",
	})
}

// UpdateBillingStatus updates company billing status
func (h *Handler) UpdateBillingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	var req UpdateBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Enabled {
		err = h.companyService.EnableBilling(c.Request.Context(), id)
	} else {
		err = h.companyService.DisableBilling(c.Request.Context(), id)
	}

	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Company not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update billing status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Billing status updated successfully",
	})
}
