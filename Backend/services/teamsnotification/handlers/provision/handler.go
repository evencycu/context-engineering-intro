package provision

import (
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
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
		g.GET("/:notifyKey", h.Read)
		g.PUT("/:notifyKey", h.Update)
		g.POST("/:notifyKey/enable", h.Enable)
		g.POST("/:notifyKey/disable", h.Disable)
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
	key := c.Param("notifyKey")
	resp, err := h.svc.Read(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *Handler) Update(c *gin.Context) {
	key := c.Param("notifyKey")
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
	key := c.Param("notifyKey")
	if err := h.svc.SetStatus(c.Request.Context(), key, true); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "enabled"})
}

func (h *Handler) Disable(c *gin.Context) {
	key := c.Param("notifyKey")
	if err := h.svc.SetStatus(c.Request.Context(), key, false); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "disabled"})
}
