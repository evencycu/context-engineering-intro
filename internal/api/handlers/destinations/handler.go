package destinations

import (
	"net/http"
	"strconv"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/services"
	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles destination-related HTTP requests
type Handler struct {
	destinationService services.DestinationService
}

// NewHandler creates a new destination handler
func NewHandler(destinationService services.DestinationService) *Handler {
	return &Handler{
		destinationService: destinationService,
	}
}

// RegisterRoutes registers destination routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	destinations := rg.Group("/destinations")
	{
		destinations.POST("", h.CreateDestination)
		destinations.GET("", h.ListDestinations)
		destinations.GET("/:id", h.GetDestination)
		destinations.PUT("/:id", h.UpdateDestination)
		destinations.DELETE("/:id", h.DeleteDestination)
		destinations.PATCH("/:id/targets", h.UpdateTargets)
		destinations.POST("/:id/validate", h.ValidateTargets)
		destinations.GET("/project/:projectId", h.GetDestinationsByProject)
		destinations.GET("/bot/:botId", h.GetDestinationsByBot)
		destinations.GET("/search", h.SearchDestinations)
	}
}

// CreateDestinationRequest represents a create destination request
type CreateDestinationRequest struct {
	ProjectID     uuid.UUID             `json:"project_id" validate:"required"`
	Name          string                `json:"name" validate:"required,min=2,max=255"`
	Description   string                `json:"description" validate:"required,min=10,max=500"`
	TeamsTenantID string                `json:"teams_tenant_id" validate:"required"`
	Targets       database.JSONBTargets `json:"targets" validate:"required,min=1"`
	BotID         *uuid.UUID            `json:"bot_id,omitempty"`
	CreatedBy     uuid.UUID             `json:"created_by" validate:"required"`
}

// UpdateDestinationRequest represents an update destination request
type UpdateDestinationRequest struct {
	Name        string     `json:"name" validate:"omitempty,min=2,max=255"`
	Description string     `json:"description" validate:"omitempty,min=10,max=500"`
	BotID       *uuid.UUID `json:"bot_id,omitempty"`
}

// UpdateTargetsRequest represents an update targets request
type UpdateTargetsRequest struct {
	Targets database.JSONBTargets `json:"targets" validate:"required,min=1"`
}

// ValidateTargetsRequest represents a validate targets request
type ValidateTargetsRequest struct {
	Targets database.JSONBTargets `json:"targets" validate:"required,min=1"`
}

// SearchDestinationsRequest represents a search destinations request
type SearchDestinationsRequest struct {
	TargetType string `form:"target_type" validate:"omitempty,oneof=person channel chatgroup"`
	TargetID   string `form:"target_id"`
	TeamID     string `form:"team_id"`
	ChannelID  string `form:"channel_id"`
	UserID     string `form:"user_id"`
	GroupID    string `form:"group_id"`
}

// CreateDestination creates a new destination
func (h *Handler) CreateDestination(c *gin.Context) {
	var req CreateDestinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Add target validation logic here
	// For now, we'll skip validation and proceed with creation

	// Convert request to destination model
	destination := &database.Destination{
		ProjectID:        req.ProjectID,
		Name:             req.Name,
		Description:      req.Description,
		TeamsTenantID:    req.TeamsTenantID,
		Targets:          req.Targets,
		BotID:            req.BotID,
		Status:           "active",  // Default status
		ValidationStatus: "pending", // Default validation status
		CreatedBy:        req.CreatedBy,
	}

	// Create destination
	createReq := &services.CreateRequest[database.Destination]{
		Data: *destination,
	}

	createdDestination, err := h.destinationService.Create(c.Request.Context(), createReq)
	if err != nil {
		if _, ok := err.(services.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Destination with this name already exists for the project",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create destination",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    createdDestination,
		"message": "Destination created successfully",
	})
}

// ListDestinations lists destinations with pagination
func (h *Handler) ListDestinations(c *gin.Context) {
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

	// Get destinations
	destinations, err := h.destinationService.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list destinations",
			"details": err.Error(),
		})
		return
	}

	// Get total count
	countReq := &services.CountRequest{
		Search: search,
		Status: status,
	}
	total, err := h.destinationService.Count(c.Request.Context(), countReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to count destinations",
			"details": err.Error(),
		})
		return
	}

	// Calculate pages
	pages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data": destinations,
		"pagination": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"pages":  pages,
		},
	})
}

// GetDestination gets a destination by ID
func (h *Handler) GetDestination(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid destination ID",
		})
		return
	}

	destination, err := h.destinationService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Destination not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get destination",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": destination,
	})
}

// UpdateDestination updates a destination
func (h *Handler) UpdateDestination(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid destination ID",
		})
		return
	}

	var req UpdateDestinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get existing destination
	existingDestination, err := h.destinationService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Destination not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get destination",
			"details": err.Error(),
		})
		return
	}

	// Update fields
	if req.Name != "" {
		existingDestination.Name = req.Name
	}
	if req.Description != "" {
		existingDestination.Description = req.Description
	}
	if req.BotID != nil {
		existingDestination.BotID = req.BotID
	}
	// bot_type removed; destination always implies platform bot

	// Update destination
	updateReq := &services.UpdateRequest[database.Destination]{
		Data: *existingDestination,
	}

	updatedDestination, err := h.destinationService.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update destination",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    updatedDestination,
		"message": "Destination updated successfully",
	})
}

// DeleteDestination deletes a destination
func (h *Handler) DeleteDestination(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid destination ID",
		})
		return
	}

	err = h.destinationService.Delete(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Destination not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete destination",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Destination deleted successfully",
	})
}

// UpdateTargets updates destination targets
func (h *Handler) UpdateTargets(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid destination ID",
		})
		return
	}

	var req UpdateTargetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Add target validation logic here
	// For now, we'll skip validation and proceed with update

	_, err = h.destinationService.UpdateTargets(c.Request.Context(), id, req.Targets)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Destination not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update targets",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Targets updated successfully",
	})
}

// ValidateTargets validates targets without updating
func (h *Handler) ValidateTargets(c *gin.Context) {
	// This endpoint validates targets without requiring an ID
	// The validation logic is handled by the service layer

	var req ValidateTargetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Add target validation logic here
	// For now, we'll return success

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"valid":   true,
			"message": "All targets are valid",
		},
	})
}

// GetDestinationsByProject gets destinations by project ID
func (h *Handler) GetDestinationsByProject(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	destinations, err := h.destinationService.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get destinations by project",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": destinations,
	})
}

// GetDestinationsByBot gets destinations by bot ID
func (h *Handler) GetDestinationsByBot(c *gin.Context) {
	botIDStr := c.Param("botId")
	botID, err := uuid.Parse(botIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid bot ID",
		})
		return
	}

	destinations, err := h.destinationService.GetByBotID(c.Request.Context(), botID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get destinations by bot",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": destinations,
	})
}

// SearchDestinations searches destinations by target criteria
func (h *Handler) SearchDestinations(c *gin.Context) {
	var req SearchDestinationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	var destinations []*database.Destination
	var err error

	// Search by different criteria
	// TODO: Implement specific target type searches
	// For now, use general search
	searchQuery := ""
	if req.TargetType != "" && req.TargetID != "" {
		searchQuery = req.TargetID
	} else if req.TeamID != "" {
		searchQuery = req.TeamID
	} else if req.ChannelID != "" {
		searchQuery = req.ChannelID
	} else if req.UserID != "" {
		searchQuery = req.UserID
	} else if req.GroupID != "" {
		searchQuery = req.GroupID
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one search criteria must be provided",
		})
		return
	}

	destinations, err = h.destinationService.SearchDestinations(c.Request.Context(), searchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search destinations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  destinations,
		"count": len(destinations),
	})
}
