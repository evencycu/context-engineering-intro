package provision

import (
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV2/internal/api/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc services.ProvisionService
}

func NewHandler(svc services.ProvisionService) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/provision")
	{
		g.POST("", h.Create)
		g.GET("/:notify_key", h.Read)
		g.PUT("/:notify_key", h.Update)
		g.POST("/:notify_key/enable", h.Enable)
		g.POST("/:notify_key/disable", h.Disable)
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req services.ProvisionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (h *Handler) Read(c *gin.Context) {
	key := c.Param("notify_key")
	resp, err := h.svc.Read(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *Handler) Update(c *gin.Context) {
	key := c.Param("notify_key")
	var req services.ProvisionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	resp, err := h.svc.Update(c.Request.Context(), key, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *Handler) Enable(c *gin.Context) {
	key := c.Param("notify_key")
	if err := h.svc.SetStatus(c.Request.Context(), key, true); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

func (h *Handler) Disable(c *gin.Context) {
	key := c.Param("notify_key")
	if err := h.svc.SetStatus(c.Request.Context(), key, false); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}
