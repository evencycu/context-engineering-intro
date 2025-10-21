package bots

import (
	"net/http"
	"strconv"

	"github.com/evencycu/TeamsNotifyGoV2/services/services"
	"github.com/evencycu/TeamsNotifyGoV2/libs/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles bot-related HTTP requests
type Handler struct {
	teamsBotService services.TeamsBotService
}

// NewHandler creates a new bot handler
func NewHandler(teamsBotService services.TeamsBotService) *Handler {
	return &Handler{
		teamsBotService: teamsBotService,
	}
}

// RegisterRoutes registers bot routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	bots := rg.Group("/bots")
	{
		// Platform (Teams) bots
		platform := bots.Group("/platform")
		{
			platform.POST("", h.CreateTeamsBotService)
			platform.GET("", h.ListteamsBotServices)
			platform.GET("/:id", h.GetteamsBotService)
			platform.PUT("/:id", h.UpdateteamsBotService)
			platform.DELETE("/:id", h.DeleteteamsBotService)
			platform.PATCH("/:id/status", h.UpdateteamsBotServiceStatus)
			platform.PATCH("/:id/capabilities", h.UpdateteamsBotServiceCapabilities)
			platform.POST("/:id/test", h.TestteamsBotServiceConnection)
		}

		// Common bot operations
		bots.GET("/status/:status", h.GetBotsByStatus)
	}
}

// CreateteamsBotServiceRequest represents a create platform bot request
type CreateteamsBotServiceRequest struct {
	Name                  string                 `json:"name" validate:"required,min=2,max=255"`
	Description           string                 `json:"description" validate:"required,min=10,max=500"`
	AppID                 string                 `json:"appId" validate:"required"`
	AppPassword           string                 `json:"appPassword" validate:"required"`
	TenantID              string                 `json:"tenantId" validate:"required"`
	WebhookURL            string                 `json:"webhookUrl" validate:"required,url"`
	Capabilities          map[string]interface{} `json:"capabilities"`
	RateLimitPerMinute    int                    `json:"rateLimitPerMinute" validate:"min=1,max=1000"`
	MaxConcurrentRequests int                    `json:"maxConcurrentRequests" validate:"min=1,max=100"`
}

// UpdateBotRequest represents an update bot request
type UpdateBotRequest struct {
	Name                  string                 `json:"name" validate:"omitempty,min=2,max=255"`
	Description           string                 `json:"description" validate:"omitempty,min=10,max=500"`
	WebhookURL            string                 `json:"webhookUrl" validate:"omitempty,url"`
	APIEndpoint           string                 `json:"apiEndpoint" validate:"omitempty,url"`
	Capabilities          map[string]interface{} `json:"capabilities"`
	RateLimitPerMinute    int                    `json:"rateLimitPerMinute" validate:"omitempty,min=1,max=1000"`
	MaxConcurrentRequests int                    `json:"maxConcurrentRequests" validate:"omitempty,min=1,max=100"`
	ContactEmail          string                 `json:"contactEmail" validate:"omitempty,email"`
	ContactPhone          string                 `json:"contactPhone" validate:"omitempty,min=10,max=20"`
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
	APIKey string `json:"apiKey" validate:"required"`
}

// TestConnectionResponse represents a test connection response
type TestConnectionResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	LatencyMs int    `json:"latency_ms,omitempty"`
}

// CreateteamsBotService creates a new platform bot
func (h *Handler) CreateTeamsBotService(c *gin.Context) {
	var req CreateteamsBotServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert request to platform bot model
	bot := &database.TeamsBot{
		Name:                  req.Name,
		Description:           func(s string) *string { return &s }(req.Description),
		AppID:                 req.AppID,
		TenantID:              func(s string) *string { return &s }(req.TenantID),
		WebhookURL:            func(s string) *string { return &s }(req.WebhookURL),
		Capabilities:          func(c database.JSONBObject) *database.JSONBObject { return &c }(req.Capabilities),
		RateLimitPerMinute:    func(i int) *int { return &i }(req.RateLimitPerMinute),
		MaxConcurrentRequests: func(i int) *int { return &i }(req.MaxConcurrentRequests),
		Status:                func(s database.BotStatus) *database.BotStatus { return &s }(database.BotStatusActive),
	}

	// Create platform bot
	createReq := &services.CreateRequest[database.TeamsBot]{
		Data: *bot,
	}

	createdBot, err := h.teamsBotService.Create(c.Request.Context(), createReq)
	if err != nil {
		if _, ok := err.(services.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{"error": "Bot with this App ID already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create platform bot", "details": err.Error()})
		return
	}

	createdBot.AppPasswordHash = ""
	c.JSON(http.StatusCreated, gin.H{"data": createdBot, "message": "Platform bot created successfully"})
}

// ListteamsBotServices lists platform bots with pagination
func (h *Handler) ListteamsBotServices(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	search := c.Query("search")
	status := c.Query("status")

	req := &services.ListRequest{Limit: limit, Offset: offset, SortBy: sortBy, Order: order, Search: search, Status: status}
	bots, err := h.teamsBotService.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list platform bots", "details": err.Error()})
		return
	}
	for _, bot := range bots {
		bot.AppPasswordHash = ""
	}
	total, err := h.teamsBotService.Count(c.Request.Context(), &services.CountRequest{Search: search, Status: status})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count platform bots", "details": err.Error()})
		return
	}
	pages := int((total + int64(limit) - 1) / int64(limit))
	c.JSON(http.StatusOK, gin.H{"data": bots, "pagination": gin.H{"total": total, "limit": limit, "offset": offset, "pages": pages}})
}

// GetteamsBotService gets a platform bot by ID
func (h *Handler) GetteamsBotService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bot ID"})
		return
	}
	bot, err := h.teamsBotService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get platform bot", "details": err.Error()})
		return
	}
	bot.AppPasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"data": bot})
}

// UpdateteamsBotService updates a platform bot
func (h *Handler) UpdateteamsBotService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bot ID"})
		return
	}
	var req UpdateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	existingBot, err := h.teamsBotService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get platform bot", "details": err.Error()})
		return
	}
	if req.Name != "" {
		existingBot.Name = req.Name
	}
	if req.Description != "" {
		d := req.Description
		existingBot.Description = &d
	}
	if req.WebhookURL != "" {
		w := req.WebhookURL
		existingBot.WebhookURL = &w
	}
	if req.Capabilities != nil {
		c := database.JSONBObject(req.Capabilities)
		existingBot.Capabilities = &c
	}
	if req.RateLimitPerMinute > 0 {
		r := req.RateLimitPerMinute
		existingBot.RateLimitPerMinute = &r
	}
	if req.MaxConcurrentRequests > 0 {
		m := req.MaxConcurrentRequests
		existingBot.MaxConcurrentRequests = &m
	}
	updateReq := &services.UpdateRequest[database.TeamsBot]{Data: *existingBot}
	updatedBot, err := h.teamsBotService.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update platform bot", "details": err.Error()})
		return
	}
	updatedBot.AppPasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"data": updatedBot, "message": "Platform bot updated successfully"})
}

// DeleteteamsBotService deletes a platform bot
func (h *Handler) DeleteteamsBotService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bot ID"})
		return
	}
	if err := h.teamsBotService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete platform bot", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Platform bot deleted successfully"})
}

// UpdateteamsBotServiceStatus updates platform bot status
func (h *Handler) UpdateteamsBotServiceStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bot ID"})
		return
	}
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	if err := h.teamsBotService.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update teams bot status", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Platform bot status updated successfully"})
}

// UpdateteamsBotServiceCapabilities updates platform bot capabilities
func (h *Handler) UpdateteamsBotServiceCapabilities(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bot ID"})
		return
	}
	var req UpdateCapabilitiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	if err := h.teamsBotService.UpdateCapabilities(c.Request.Context(), id, req.Capabilities); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update platform bot capabilities", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Platform bot capabilities updated successfully"})
}

// TestteamsBotServiceConnection tests platform bot connection
func (h *Handler) TestteamsBotServiceConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bot ID"})
		return
	}
	if err := h.teamsBotService.TestConnection(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusOK, gin.H{"data": TestConnectionResponse{Success: false, Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": TestConnectionResponse{Success: true, Message: "Connection test successful"}})
}

// GetBotsByStatus gets bots by status
func (h *Handler) GetBotsByStatus(c *gin.Context) {
	status := c.Param("status")
	teamsBotServices, err := h.teamsBotService.GetByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get platform bots by status", "details": err.Error()})
		return
	}
	for _, bot := range teamsBotServices {
		bot.AppPasswordHash = ""
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"platform_bots": teamsBotServices}})
}
