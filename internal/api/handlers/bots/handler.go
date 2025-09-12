package bots

import (
	"net/http"
	"strconv"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/services"
	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles bot-related HTTP requests
type Handler struct {
	platformBotService   services.BotService
	thirdPartyBotService services.ThirdPartyBotService
}

// NewHandler creates a new bot handler
func NewHandler(platformBotService services.BotService, thirdPartyBotService services.ThirdPartyBotService) *Handler {
	return &Handler{
		platformBotService:   platformBotService,
		thirdPartyBotService: thirdPartyBotService,
	}
}

// RegisterRoutes registers bot routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	bots := rg.Group("/bots")
	{
		// Platform bots
		platform := bots.Group("/platform")
		{
			platform.POST("", h.CreatePlatformBot)
			platform.GET("", h.ListPlatformBots)
			platform.GET("/:id", h.GetPlatformBot)
			platform.PUT("/:id", h.UpdatePlatformBot)
			platform.DELETE("/:id", h.DeletePlatformBot)
			platform.PATCH("/:id/status", h.UpdatePlatformBotStatus)
			platform.PATCH("/:id/capabilities", h.UpdatePlatformBotCapabilities)
			platform.POST("/:id/test", h.TestPlatformBotConnection)
		}

		// Third-party bots
		thirdParty := bots.Group("/third-party")
		{
			thirdParty.POST("", h.CreateThirdPartyBot)
			thirdParty.GET("", h.ListThirdPartyBots)
			thirdParty.GET("/:id", h.GetThirdPartyBot)
			thirdParty.PUT("/:id", h.UpdateThirdPartyBot)
			thirdParty.DELETE("/:id", h.DeleteThirdPartyBot)
			thirdParty.PATCH("/:id/status", h.UpdateThirdPartyBotStatus)
			thirdParty.PATCH("/:id/api-key", h.UpdateThirdPartyBotAPIKey)
			thirdParty.POST("/:id/test", h.TestThirdPartyBotConnection)
		}

		// Common bot operations
		bots.GET("/company/:companyId", h.GetBotsByCompany)
		bots.GET("/status/:status", h.GetBotsByStatus)
	}
}

// CreatePlatformBotRequest represents a create platform bot request
type CreatePlatformBotRequest struct {
	Name                  string                 `json:"name" validate:"required,min=2,max=255"`
	Description           string                 `json:"description" validate:"required,min=10,max=500"`
	AppID                 string                 `json:"app_id" validate:"required"`
	AppPassword           string                 `json:"app_password" validate:"required"`
	TenantID              string                 `json:"tenant_id" validate:"required"`
	WebhookURL            string                 `json:"webhook_url" validate:"required,url"`
	Capabilities          map[string]interface{} `json:"capabilities"`
	RateLimitPerMinute    int                    `json:"rate_limit_per_minute" validate:"min=1,max=1000"`
	MaxConcurrentRequests int                    `json:"max_concurrent_requests" validate:"min=1,max=100"`
}

// CreateThirdPartyBotRequest represents a create third-party bot request
type CreateThirdPartyBotRequest struct {
	CompanyID             uuid.UUID              `json:"company_id" validate:"required"`
	Name                  string                 `json:"name" validate:"required,min=2,max=255"`
	Description           string                 `json:"description" validate:"required,min=10,max=500"`
	AppID                 string                 `json:"app_id" validate:"required"`
	AppPassword           string                 `json:"app_password" validate:"required"`
	TenantID              string                 `json:"tenant_id" validate:"required"`
	WebhookURL            string                 `json:"webhook_url" validate:"required,url"`
	APIEndpoint           string                 `json:"api_endpoint" validate:"required,url"`
	APIKey                string                 `json:"api_key" validate:"required"`
	Capabilities          map[string]interface{} `json:"capabilities"`
	RateLimitPerMinute    int                    `json:"rate_limit_per_minute" validate:"min=1,max=1000"`
	MaxConcurrentRequests int                    `json:"max_concurrent_requests" validate:"min=1,max=100"`
	ContactEmail          string                 `json:"contact_email" validate:"required,email"`
	ContactPhone          string                 `json:"contact_phone" validate:"required,min=10,max=20"`
	CreatedBy             uuid.UUID              `json:"created_by" validate:"required"`
}

// UpdateBotRequest represents an update bot request
type UpdateBotRequest struct {
	Name                  string                 `json:"name" validate:"omitempty,min=2,max=255"`
	Description           string                 `json:"description" validate:"omitempty,min=10,max=500"`
	WebhookURL            string                 `json:"webhook_url" validate:"omitempty,url"`
	APIEndpoint           string                 `json:"api_endpoint" validate:"omitempty,url"`
	Capabilities          map[string]interface{} `json:"capabilities"`
	RateLimitPerMinute    int                    `json:"rate_limit_per_minute" validate:"omitempty,min=1,max=1000"`
	MaxConcurrentRequests int                    `json:"max_concurrent_requests" validate:"omitempty,min=1,max=100"`
	ContactEmail          string                 `json:"contact_email" validate:"omitempty,email"`
	ContactPhone          string                 `json:"contact_phone" validate:"omitempty,min=10,max=20"`
}

// UpdateStatusRequest represents an update status request
type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended maintenance"`
}

// UpdateCapabilitiesRequest represents an update capabilities request
type UpdateCapabilitiesRequest struct {
	Capabilities map[string]interface{} `json:"capabilities" validate:"required"`
}

// UpdateAPIKeyRequest represents an update API key request
type UpdateAPIKeyRequest struct {
	APIKey string `json:"api_key" validate:"required"`
}

// TestConnectionResponse represents a test connection response
type TestConnectionResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	LatencyMs int    `json:"latency_ms,omitempty"`
}

// CreatePlatformBot creates a new platform bot
func (h *Handler) CreatePlatformBot(c *gin.Context) {
	var req CreatePlatformBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert request to platform bot model
	bot := &database.PlatformBot{
		Name:                  req.Name,
		Description:           req.Description,
		AppID:                 req.AppID,
		TenantID:              req.TenantID,
		WebhookURL:            req.WebhookURL,
		Capabilities:          req.Capabilities,
		RateLimitPerMinute:    req.RateLimitPerMinute,
		MaxConcurrentRequests: req.MaxConcurrentRequests,
		Status:                database.BotStatusActive, // Default status
		// AppPassword will be hashed in service layer
	}

	// Create platform bot
	createReq := &services.CreateRequest[database.PlatformBot]{
		Data: *bot,
	}

	createdBot, err := h.platformBotService.Create(c.Request.Context(), createReq)
	if err != nil {
		if _, ok := err.(services.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Bot with this App ID already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create platform bot",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	createdBot.AppPasswordHash = ""

	c.JSON(http.StatusCreated, gin.H{
		"data":    createdBot,
		"message": "Platform bot created successfully",
	})
}

// CreateThirdPartyBot creates a new third-party bot
func (h *Handler) CreateThirdPartyBot(c *gin.Context) {
	var req CreateThirdPartyBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert request to third-party bot model
	bot := &database.ThirdPartyBot{
		CompanyID:             req.CompanyID,
		Name:                  req.Name,
		Description:           req.Description,
		AppID:                 req.AppID,
		TenantID:              req.TenantID,
		WebhookURL:            req.WebhookURL,
		APIEndpoint:           req.APIEndpoint,
		Capabilities:          req.Capabilities,
		RateLimitPerMinute:    req.RateLimitPerMinute,
		MaxConcurrentRequests: req.MaxConcurrentRequests,
		ContactEmail:          req.ContactEmail,
		ContactPhone:          req.ContactPhone,
		CreatedBy:             req.CreatedBy,
		Status:                database.BotStatusActive, // Default status
		// AppPassword and APIKey will be hashed in service layer
	}

	// Create third-party bot
	createReq := &services.CreateRequest[database.ThirdPartyBot]{
		Data: *bot,
	}

	createdBot, err := h.thirdPartyBotService.Create(c.Request.Context(), createReq)
	if err != nil {
		if _, ok := err.(services.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Bot with this App ID already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create third-party bot",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	createdBot.AppPasswordHash = ""
	createdBot.APIKeyHash = ""

	c.JSON(http.StatusCreated, gin.H{
		"data":    createdBot,
		"message": "Third-party bot created successfully",
	})
}

// ListPlatformBots lists platform bots with pagination
func (h *Handler) ListPlatformBots(c *gin.Context) {
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

	// Get platform bots
	bots, err := h.platformBotService.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list platform bots",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	for _, bot := range bots {
		bot.AppPasswordHash = ""
	}

	// Get total count
	countReq := &services.CountRequest{
		Search: search,
		Status: status,
	}
	total, err := h.platformBotService.Count(c.Request.Context(), countReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to count platform bots",
			"details": err.Error(),
		})
		return
	}

	// Calculate pages
	pages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data": bots,
		"pagination": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"pages":  pages,
		},
	})
}

// ListThirdPartyBots lists third-party bots with pagination
func (h *Handler) ListThirdPartyBots(c *gin.Context) {
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

	// Get third-party bots
	bots, err := h.thirdPartyBotService.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list third-party bots",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	for _, bot := range bots {
		bot.AppPasswordHash = ""
		bot.APIKeyHash = ""
	}

	// Get total count
	countReq := &services.CountRequest{
		Search: search,
		Status: status,
	}
	total, err := h.thirdPartyBotService.Count(c.Request.Context(), countReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to count third-party bots",
			"details": err.Error(),
		})
		return
	}

	// Calculate pages
	pages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data": bots,
		"pagination": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"pages":  pages,
		},
	})
}

// GetPlatformBot gets a platform bot by ID
func (h *Handler) GetPlatformBot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	bot, err := h.platformBotService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Platform bot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get platform bot",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	bot.AppPasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"data": bot,
	})
}

// GetThirdPartyBot gets a third-party bot by ID
func (h *Handler) GetThirdPartyBot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	bot, err := h.thirdPartyBotService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Third-party bot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get third-party bot",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	bot.AppPasswordHash = ""
	bot.APIKeyHash = ""

	c.JSON(http.StatusOK, gin.H{
		"data": bot,
	})
}

// UpdatePlatformBot updates a platform bot
func (h *Handler) UpdatePlatformBot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	var req UpdateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get existing bot
	existingBot, err := h.platformBotService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Platform bot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get platform bot",
			"details": err.Error(),
		})
		return
	}

	// Update fields
	if req.Name != "" {
		existingBot.Name = req.Name
	}
	if req.Description != "" {
		existingBot.Description = req.Description
	}
	if req.WebhookURL != "" {
		existingBot.WebhookURL = req.WebhookURL
	}
	if req.Capabilities != nil {
		existingBot.Capabilities = req.Capabilities
	}
	if req.RateLimitPerMinute > 0 {
		existingBot.RateLimitPerMinute = req.RateLimitPerMinute
	}
	if req.MaxConcurrentRequests > 0 {
		existingBot.MaxConcurrentRequests = req.MaxConcurrentRequests
	}

	// Update bot
	updateReq := &services.UpdateRequest[database.PlatformBot]{
		Data: *existingBot,
	}

	updatedBot, err := h.platformBotService.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update platform bot",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	updatedBot.AppPasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"data":    updatedBot,
		"message": "Platform bot updated successfully",
	})
}

// UpdateThirdPartyBot updates a third-party bot
func (h *Handler) UpdateThirdPartyBot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	var req UpdateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get existing bot
	existingBot, err := h.thirdPartyBotService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Third-party bot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get third-party bot",
			"details": err.Error(),
		})
		return
	}

	// Update fields
	if req.Name != "" {
		existingBot.Name = req.Name
	}
	if req.Description != "" {
		existingBot.Description = req.Description
	}
	if req.WebhookURL != "" {
		existingBot.WebhookURL = req.WebhookURL
	}
	if req.APIEndpoint != "" {
		existingBot.APIEndpoint = req.APIEndpoint
	}
	if req.Capabilities != nil {
		existingBot.Capabilities = req.Capabilities
	}
	if req.RateLimitPerMinute > 0 {
		existingBot.RateLimitPerMinute = req.RateLimitPerMinute
	}
	if req.MaxConcurrentRequests > 0 {
		existingBot.MaxConcurrentRequests = req.MaxConcurrentRequests
	}
	if req.ContactEmail != "" {
		existingBot.ContactEmail = req.ContactEmail
	}
	if req.ContactPhone != "" {
		existingBot.ContactPhone = req.ContactPhone
	}

	// Update bot
	updateReq := &services.UpdateRequest[database.ThirdPartyBot]{
		Data: *existingBot,
	}

	updatedBot, err := h.thirdPartyBotService.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update third-party bot",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	updatedBot.AppPasswordHash = ""
	updatedBot.APIKeyHash = ""

	c.JSON(http.StatusOK, gin.H{
		"data":    updatedBot,
		"message": "Third-party bot updated successfully",
	})
}

// DeletePlatformBot deletes a platform bot
func (h *Handler) DeletePlatformBot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	err = h.platformBotService.Delete(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Platform bot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete platform bot",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Platform bot deleted successfully",
	})
}

// DeleteThirdPartyBot deletes a third-party bot
func (h *Handler) DeleteThirdPartyBot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	err = h.thirdPartyBotService.Delete(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Third-party bot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete third-party bot",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Third-party bot deleted successfully",
	})
}

// UpdatePlatformBotStatus updates platform bot status
func (h *Handler) UpdatePlatformBotStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
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

	err = h.platformBotService.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Platform bot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update platform bot status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Platform bot status updated successfully",
	})
}

// UpdateThirdPartyBotStatus updates third-party bot status
func (h *Handler) UpdateThirdPartyBotStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
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

	err = h.thirdPartyBotService.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Third-party bot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update third-party bot status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Third-party bot status updated successfully",
	})
}

// UpdatePlatformBotCapabilities updates platform bot capabilities
func (h *Handler) UpdatePlatformBotCapabilities(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	var req UpdateCapabilitiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	err = h.platformBotService.UpdateCapabilities(c.Request.Context(), id, req.Capabilities)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Platform bot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update platform bot capabilities",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Platform bot capabilities updated successfully",
	})
}

// UpdateThirdPartyBotAPIKey updates third-party bot API key
func (h *Handler) UpdateThirdPartyBotAPIKey(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	var req UpdateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	err = h.thirdPartyBotService.UpdateAPIKey(c.Request.Context(), id, req.APIKey)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Third-party bot not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update third-party bot API key",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Third-party bot API key updated successfully",
	})
}

// TestPlatformBotConnection tests platform bot connection
func (h *Handler) TestPlatformBotConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	err = h.platformBotService.TestConnection(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Platform bot not found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": TestConnectionResponse{
				Success: false,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": TestConnectionResponse{
			Success: true,
			Message: "Connection test successful",
		},
	})
}

// TestThirdPartyBotConnection tests third-party bot connection
func (h *Handler) TestThirdPartyBotConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	err = h.thirdPartyBotService.TestConnection(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Third-party bot not found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": TestConnectionResponse{
				Success: false,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": TestConnectionResponse{
			Success: true,
			Message: "Connection test successful",
		},
	})
}

// GetBotsByCompany gets bots by company ID
func (h *Handler) GetBotsByCompany(c *gin.Context) {
	companyIDStr := c.Param("companyId")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	// Get third-party bots for the company
	thirdPartyBots, err := h.thirdPartyBotService.GetByCompanyID(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get third-party bots by company",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	for _, bot := range thirdPartyBots {
		bot.AppPasswordHash = ""
		bot.APIKeyHash = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"third_party_bots": thirdPartyBots,
		},
	})
}

// GetBotsByStatus gets bots by status
func (h *Handler) GetBotsByStatus(c *gin.Context) {
	status := c.Param("status")

	// Get platform bots by status
	platformBots, err := h.platformBotService.GetByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get platform bots by status",
			"details": err.Error(),
		})
		return
	}

	// Get third-party bots by status
	thirdPartyBots, err := h.thirdPartyBotService.GetByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get third-party bots by status",
			"details": err.Error(),
		})
		return
	}

	// Remove sensitive fields from response
	for _, bot := range platformBots {
		bot.AppPasswordHash = ""
	}
	for _, bot := range thirdPartyBots {
		bot.AppPasswordHash = ""
		bot.APIKeyHash = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"platform_bots":    platformBots,
			"third_party_bots": thirdPartyBots,
		},
	})
}
