package projects

import (
	"net/http"
	"strconv"

	"github.com/evencycu/TeamsNotifyGoV2/services/services"
	"github.com/evencycu/TeamsNotifyGoV2/libs/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles project-related HTTP requests
type Handler struct {
	projectService services.ProjectService
}

// NewHandler creates a new project handler
func NewHandler(projectService services.ProjectService) *Handler {
	return &Handler{
		projectService: projectService,
	}
}

// RegisterRoutes registers project routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	projects := rg.Group("/projects")
	{
		projects.POST("", h.CreateProject)
		projects.GET("", h.ListProjects)
		projects.GET("/:id", h.GetProject)
		projects.PUT("/:id", h.UpdateProject)
		projects.DELETE("/:id", h.DeleteProject)
		projects.PATCH("/:id/limits", h.UpdateLimits)
		projects.GET("/company/:companyId", h.GetProjectsByCompany)
		projects.GET("/key/:keyName", h.GetProjectByKeyName)
	}
}

// CreateProjectRequest represents a create project request
type CreateProjectRequest struct {
	CompanyID    uuid.UUID `json:"companyId" validate:"required"`
	NotifyKey    string    `json:"notifyKey" validate:"required,min=3,max=255"`
	Description  string    `json:"description" validate:"required,min=10,max=500"`
	DailyLimit   int       `json:"dailyLimit" validate:"min=1,max=10000"`
	MonthlyLimit int       `json:"monthlyLimit" validate:"min=1,max=300000"`
	Priority     string    `json:"priority" validate:"oneof=low normal high"`
	CreatedBy    uuid.UUID `json:"createdBy" validate:"required"`
}

// UpdateProjectRequest represents an update project request
type UpdateProjectRequest struct {
	NotifyKey   string `json:"notifyKey" validate:"omitempty,min=3,max=50,alphanum"`
	Description string `json:"description" validate:"omitempty,min=10,max=500"`
	Priority    string `json:"priority" validate:"omitempty,oneof=low normal high"`
}

// UpdateLimitsRequest represents an update limits request
type UpdateLimitsRequest struct {
	DailyLimit   int `json:"dailyLimit" validate:"required,min=1,max=10000"`
	MonthlyLimit int `json:"monthlyLimit" validate:"required,min=1,max=300000"`
}

// CreateProject creates a new project
func (h *Handler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Convert request to project model
	project := &database.Project{
		CompanyID:    req.CompanyID,
		NotifyKey:    req.NotifyKey,
		Description:  req.Description,
		DailyLimit:   req.DailyLimit,
		MonthlyLimit: req.MonthlyLimit,
		Priority:     req.Priority,
		CreatedBy:    req.CreatedBy,
		Status:       "active", // Default status
	}

	// Create project
	createReq := &services.CreateRequest[database.Project]{
		Data: *project,
	}

	createdProject, err := h.projectService.Create(c.Request.Context(), createReq)
	if err != nil {
		if _, ok := err.(services.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Project with this key name already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create project",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    createdProject,
		"message": "Project created successfully",
	})
}

// ListProjects lists projects with pagination
func (h *Handler) ListProjects(c *gin.Context) {
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

	// Get projects
	projects, err := h.projectService.List(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list projects",
			"details": err.Error(),
		})
		return
	}

	// Get total count
	countReq := &services.CountRequest{
		Search: search,
		Status: status,
	}
	total, err := h.projectService.Count(c.Request.Context(), countReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to count projects",
			"details": err.Error(),
		})
		return
	}

	// Calculate pages
	pages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data": projects,
		"pagination": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
			"pages":  pages,
		},
	})
}

// GetProject gets a project by ID
func (h *Handler) GetProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	project, err := h.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get project",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": project,
	})
}

// UpdateProject updates a project
func (h *Handler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get existing project
	existingProject, err := h.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get project",
			"details": err.Error(),
		})
		return
	}

	// Update fields
	if req.NotifyKey != "" {
		existingProject.NotifyKey = req.NotifyKey
	}
	if req.Description != "" {
		existingProject.Description = req.Description
	}
	if req.Priority != "" {
		existingProject.Priority = req.Priority
	}

	// Update project
	updateReq := &services.UpdateRequest[database.Project]{
		Data: *existingProject,
	}

	updatedProject, err := h.projectService.Update(c.Request.Context(), id, updateReq)
	if err != nil {
		if _, ok := err.(services.ConflictError); ok {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Project with this key name already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update project",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    updatedProject,
		"message": "Project updated successfully",
	})
}

// DeleteProject deletes a project
func (h *Handler) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	err = h.projectService.Delete(c.Request.Context(), id)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete project",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project deleted successfully",
	})
}

// UpdateLimits updates project limits
func (h *Handler) UpdateLimits(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid project ID",
		})
		return
	}

	var req UpdateLimitsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	err = h.projectService.UpdateLimits(c.Request.Context(), id, req.DailyLimit, req.MonthlyLimit)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update project limits",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project limits updated successfully",
	})
}

// GetProjectsByCompany gets projects by company ID
func (h *Handler) GetProjectsByCompany(c *gin.Context) {
	companyIDStr := c.Param("companyId")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid company ID",
		})
		return
	}

	projects, err := h.projectService.GetByCompanyID(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get projects by company",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": projects,
	})
}

// GetProjectByKeyName gets a project by key name
func (h *Handler) GetProjectByKeyName(c *gin.Context) {
	keyName := c.Param("keyName")

	project, err := h.projectService.GetByKeyName(c.Request.Context(), keyName)
	if err != nil {
		if _, ok := err.(services.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get project by key name",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": project,
	})
}
