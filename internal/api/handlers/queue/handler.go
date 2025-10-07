package queue

import (
	"net/http"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Handler provides queue-related endpoints (observability)
type Handler struct {
	redis *redis.Client
}

func NewHandler(redis *redis.Client) *Handler {
	return &Handler{redis: redis}
}

// RegisterRoutes implements api.Handler
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/queue/stats", h.GetQueueStats)
}

type queueStatsResponse struct {
	High   int64 `json:"high"`
	Normal int64 `json:"normal"`
	Low    int64 `json:"low"`
}

// GetQueueStats returns lengths of the three priority queues
func (h *Handler) GetQueueStats(c *gin.Context) {
	ctx := context.Background()
	high := h.redis.LLen(ctx, "queue:notifications:high").Val()
	normal := h.redis.LLen(ctx, "queue:notifications:normal").Val()
	low := h.redis.LLen(ctx, "queue:notifications:low").Val()

	c.JSON(http.StatusOK, gin.H{
		"data": queueStatsResponse{High: high, Normal: normal, Low: low},
	})
}
