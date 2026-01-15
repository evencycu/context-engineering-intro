package audience

import (
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles audience list-related HTTP requests
type Handler struct {
	audienceListService services.AudienceListService
}

// NewHandler creates a new audience list handler
func NewHandler(audienceListService services.AudienceListService) *Handler {
	return &Handler{
		audienceListService: audienceListService,
	}
}

// RegisterRoutes registers audience list routes
// Note: Uses :id instead of :projectId to avoid route conflict with /projects/:id
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	audience := rg.Group("/projects/:id/audience-lists")
	{
		audience.GET("", h.ListAudienceLists)
		audience.POST("", h.CreateAudienceList)
		audience.GET("/:listId", h.GetAudienceList)
		audience.PUT("/:listId", h.UpdateAudienceList)
		audience.DELETE("/:listId", h.DeleteAudienceList)
	}
}

// CreateAudienceListRequest represents a create audience list request
type CreateAudienceListRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=255"`
	Type        string                 `json:"type" validate:"omitempty,oneof=Static Dynamic"`
	Count       int                    `json:"count" validate:"min=0"`
	Description *string                `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateAudienceListRequest represents an update audience list request
type UpdateAudienceListRequest struct {
	Name        string                 `json:"name" validate:"omitempty,min=1,max=255"`
	Type        string                 `json:"type" validate:"omitempty,oneof=Static Dynamic"`
	Count       int                    `json:"count" validate:"omitempty,min=0"`
	Description *string                `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// CreateAudienceList creates a new audience list
func (h *Handler) CreateAudienceList(c *gin.Context) {
	projectIDStr := c.Param("id") // Use "id" to match route parameter
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req CreateAudienceListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert request to audience list model
	audienceList := &models.AudienceList{
		ProjectID:  projectID,
		Name:       req.Name,
		Type:       req.Type,
		Count:      req.Count,
		Description: req.Description,
		Metadata:   models.JSONBObject(req.Metadata),
	}
	if audienceList.Type == "" {
		audienceList.Type = "Static"
	}
	if audienceList.Metadata == nil {
		audienceList.Metadata = make(models.JSONBObject)
	}

	createReq := &services.CreateRequest[models.AudienceList]{
		Data: *audienceList,
	}

	createdList, err := h.audienceListService.Create(c.Request.Context(), createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create audience list",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    createdList,
		"message": "Audience list created successfully",
	})
}

// ListAudienceLists lists audience lists for a project
func (h *Handler) ListAudienceLists(c *gin.Context) {
	projectIDStr := c.Param("id") // Use "id" to match route parameter
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	lists, err := h.audienceListService.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list audience lists",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": lists,
	})
}

// GetAudienceList gets an audience list by ID
func (h *Handler) GetAudienceList(c *gin.Context) {
	idStr := c.Param("listId") // Use "listId" to distinguish from project "id"
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid audience list ID",
		})
		return
	}

	list, err := h.audienceListService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Audience list not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get audience list",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

// UpdateAudienceList updates an audience list
func (h *Handler) UpdateAudienceList(c *gin.Context) {
	idStr := c.Param("listId") // Use "listId" to distinguish from project "id"
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid audience list ID",
		})
		return
	}

	var req UpdateAudienceListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get existing list
	existingList, err := h.audienceListService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Audience list not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get audience list",
			"details": err.Error(),
		})
		return
	}

	// Update fields
	if req.Name != "" {
		existingList.Name = req.Name
	}
	if req.Type != "" {
		existingList.Type = req.Type
	}
	if req.Count > 0 {
		existingList.Count = req.Count
	}
	if req.Description != nil {
		existingList.Description = req.Description
	}
	if req.Metadata != nil {
		existingList.Metadata = models.JSONBObject(req.Metadata)
	}

	updateReq := &services.UpdateRequest[models.AudienceList]{
		Data: *existingList,
	}

	updatedList, err := h.audienceListService.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update audience list",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    updatedList,
		"message": "Audience list updated successfully",
	})
}

// DeleteAudienceList deletes an audience list
func (h *Handler) DeleteAudienceList(c *gin.Context) {
	idStr := c.Param("listId") // Use "listId" to distinguish from project "id"
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid audience list ID",
		})
		return
	}

	err = h.audienceListService.Delete(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Audience list not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete audience list",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Audience list deleted successfully",
	})
}
