package messages

import (
    "encoding/json"
    "net/http"

    "github.com/evencycu/TeamsNotifyGoV2/internal/api/services"
    "github.com/gin-gonic/gin"
)

// Handler handles /api/messages webhook from Bot Framework
type Handler struct {
    service services.MessagesService
}

func NewHandler(service services.MessagesService) *Handler {
    return &Handler{service: service}
}

// Handle receives Bot Framework Activity and persists installation info
func (h *Handler) Handle(c *gin.Context) {
    var act services.Activity
    var raw map[string]any

    if err := c.ShouldBindJSON(&raw); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
        return
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

// mapToStruct marshals then unmarshals to map into struct
func mapToStruct(m map[string]any, out any) error {
    b, err := json.Marshal(m)
    if err != nil { return err }
    return json.Unmarshal(b, out)
}


