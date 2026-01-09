package messages

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/services"
	"github.com/gin-gonic/gin"
)

// Handler handles /api/messages webhook from Bot Framework
type Handler struct {
	service services.MessagesService
}

func NewHandler(service services.MessagesService) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers message routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	messages := rg.Group("/messages")
	{
		messages.POST("", h.Handle)
		messages.POST("/test", h.ProactiveTest)
	}
}

// Handle receives Bot Framework Activity and persists installation info
func (h *Handler) Handle(c *gin.Context) {
	var act services.Activity
	var raw map[string]any

	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}

	// Log full payload
	if b, err := json.MarshalIndent(raw, "", "  "); err == nil {
		log.Printf("/api/v1/messages payload: %s", string(b))
	}

	// Best-effort map to typed Activity
	if err := mapToStruct(raw, &act); err != nil {
		// proceed with minimal fields if mapping fails
	}

	if err := h.service.HandleActivity(c.Request.Context(), &act, raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ProactiveTestRequest represents the payload to send a proactive message without DB
type ProactiveTestRequest struct {
	Activity map[string]any `json:"activity"` // the installationUpdate payload
	Text     string         `json:"text"`     // message to send
}

// ProactiveTest sends a proactive message to the provided conversation using env creds
func (h *Handler) ProactiveTest(c *gin.Context) {
	var req ProactiveTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}
	if req.Activity == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "activity is required"})
		return
	}
	if req.Text == "" {
		req.Text = "Hello from proactive test"
	}

	// Use service to send proactive without DB
	if err := h.service.SendProactiveTest(c.Request.Context(), req.Activity, req.Text); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": string(models.NotificationStatusSent)})
}

// mapToStruct marshals then unmarshals to map into struct
func mapToStruct(m map[string]any, out any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
