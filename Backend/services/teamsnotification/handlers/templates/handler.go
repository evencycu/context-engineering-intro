package templates

import (
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV3/libs/response"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	templateService services.TemplateService
}

func NewHandler(templateService services.TemplateService) *Handler {
	return &Handler{templateService: templateService}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	templates := rg.Group("/projects/:id/templates")
	{
		templates.GET("", h.ListTemplates)
		templates.POST("", h.CreateTemplate)
		templates.GET("/:templateId", h.GetTemplate)
		templates.PUT("/:templateId", h.UpdateTemplate)
		templates.DELETE("/:templateId", h.DeleteTemplate)
	}

	// Global templates (available to all projects)
	global := rg.Group("/global/templates")
	{
		global.GET("", h.ListGlobalTemplates)
		global.POST("", h.CreateGlobalTemplate)
		global.GET("/:templateId", h.GetGlobalTemplate)
		global.PUT("/:templateId", h.UpdateGlobalTemplate)
		global.DELETE("/:templateId", h.DeleteGlobalTemplate)
	}
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID", err, nil)
		return
	}

	var req services.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err, nil)
		return
	}
	req.ProjectID = projectID

	template, err := h.templateService.Create(c.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(c, "Failed to create template", err, nil)
		return
	}

	response.Success(c, http.StatusCreated, "Template created successfully", template)
}

func (h *Handler) ListTemplates(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID", err, nil)
		return
	}

	templates, err := h.templateService.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		response.InternalServerError(c, "Failed to list templates", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Templates listed successfully", templates)
}

func (h *Handler) GetTemplate(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("templateId"))
	if err != nil {
		response.BadRequest(c, "Invalid template ID", err, nil)
		return
	}

	template, err := h.templateService.GetByID(c.Request.Context(), templateID)
	if err != nil {
		response.InternalServerError(c, "Failed to get template", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Template retrieved successfully", template)
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("templateId"))
	if err != nil {
		response.BadRequest(c, "Invalid template ID", err, nil)
		return
	}

	var req services.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err, nil)
		return
	}

	template, err := h.templateService.Update(c.Request.Context(), templateID, &req)
	if err != nil {
		response.InternalServerError(c, "Failed to update template", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Template updated successfully", template)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("templateId"))
	if err != nil {
		response.BadRequest(c, "Invalid template ID", err, nil)
		return
	}

	if err := h.templateService.Delete(c.Request.Context(), templateID); err != nil {
		response.InternalServerError(c, "Failed to delete template", err, nil)
		return
	}

	response.Success(c, http.StatusNoContent, "Template deleted successfully", nil)
}

// Global Template handlers (use special UUID 00000000-0000-0000-0000-000000000000 for global templates)
var globalProjectID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

func (h *Handler) CreateGlobalTemplate(c *gin.Context) {
	var req services.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err, nil)
		return
	}
	req.ProjectID = globalProjectID // Use special UUID for global templates

	template, err := h.templateService.Create(c.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(c, "Failed to create global template", err, nil)
		return
	}

	response.Success(c, http.StatusCreated, "Global template created successfully", template)
}

func (h *Handler) ListGlobalTemplates(c *gin.Context) {
	templates, err := h.templateService.GetByProjectID(c.Request.Context(), globalProjectID)
	if err != nil {
		response.InternalServerError(c, "Failed to list global templates", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Global templates listed successfully", templates)
}

func (h *Handler) GetGlobalTemplate(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("templateId"))
	if err != nil {
		response.BadRequest(c, "Invalid template ID", err, nil)
		return
	}

	template, err := h.templateService.GetByID(c.Request.Context(), templateID)
	if err != nil {
		response.InternalServerError(c, "Failed to get global template", err, nil)
		return
	}

	// Verify it's actually a global template
	if template.ProjectID != globalProjectID {
		response.BadRequest(c, "Template is not a global template", nil, nil)
		return
	}

	response.Success(c, http.StatusOK, "Global template retrieved successfully", template)
}

func (h *Handler) UpdateGlobalTemplate(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("templateId"))
	if err != nil {
		response.BadRequest(c, "Invalid template ID", err, nil)
		return
	}

	// Verify it's actually a global template
	template, err := h.templateService.GetByID(c.Request.Context(), templateID)
	if err != nil {
		response.InternalServerError(c, "Failed to get global template", err, nil)
		return
	}
	if template.ProjectID != globalProjectID {
		response.BadRequest(c, "Template is not a global template", nil, nil)
		return
	}

	var req services.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err, nil)
		return
	}

	updated, err := h.templateService.Update(c.Request.Context(), templateID, &req)
	if err != nil {
		response.InternalServerError(c, "Failed to update global template", err, nil)
		return
	}

	response.Success(c, http.StatusOK, "Global template updated successfully", updated)
}

func (h *Handler) DeleteGlobalTemplate(c *gin.Context) {
	templateID, err := uuid.Parse(c.Param("templateId"))
	if err != nil {
		response.BadRequest(c, "Invalid template ID", err, nil)
		return
	}

	// Verify it's actually a global template
	template, err := h.templateService.GetByID(c.Request.Context(), templateID)
	if err != nil {
		response.InternalServerError(c, "Failed to get global template", err, nil)
		return
	}
	if template.ProjectID != globalProjectID {
		response.BadRequest(c, "Template is not a global template", nil, nil)
		return
	}

	if err := h.templateService.Delete(c.Request.Context(), templateID); err != nil {
		response.InternalServerError(c, "Failed to delete global template", err, nil)
		return
	}

	response.Success(c, http.StatusNoContent, "Global template deleted successfully", nil)
}
