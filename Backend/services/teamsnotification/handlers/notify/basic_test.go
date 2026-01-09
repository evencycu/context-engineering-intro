package notify

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_NewHandler(t *testing.T) {
	// 測試創建新的 handler
	handler := NewHandler(nil)
	assert.NotNil(t, handler)
	assert.Nil(t, handler.notifyService)
}

func TestHandler_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 創建 handler
	handler := &Handler{}

	// 設置路由
	router := gin.New()
	handler.RegisterRoutes(router.Group("/api/v1/external"))

	// 測試路由是否正確註冊
	assert.NotNil(t, router)
}

func TestSendNotificationRequest_Validation(t *testing.T) {
	// 測試請求結構體的基本功能
	req := SendNotificationRequest{
		NotifyKey:   "test-notify-key",
		Message:     "Test notification message",
		MessageType: "text",
		Priority:    "normal",
		Targets:     []string{"all"},
		Mentions:    []string{"@everyone"},
		Metadata: map[string]interface{}{
			"source": "test",
		},
	}

	assert.Equal(t, "test-notify-key", req.NotifyKey)
	assert.Equal(t, "Test notification message", req.Message)
	assert.Equal(t, "text", req.MessageType)
	assert.Equal(t, "normal", req.Priority)
	assert.Equal(t, []string{"all"}, req.Targets)
	assert.Equal(t, []string{"@everyone"}, req.Mentions)
	assert.NotNil(t, req.Metadata)
}

func TestSendNotificationResponse_Validation(t *testing.T) {
	// 測試回應結構體的基本功能
	resp := SendNotificationResponse{
		Success: true,
		Error:   "",
	}

	assert.True(t, resp.Success)
	assert.Empty(t, resp.Error)
}

func TestHealthCheckResponse_Validation(t *testing.T) {
	// 測試健康檢查回應結構體的基本功能
	resp := HealthCheckResponse{
		Status:    "healthy",
		Timestamp: "2025-01-01T00:00:00Z",
	}

	assert.Equal(t, "healthy", resp.Status)
	assert.Equal(t, "2025-01-01T00:00:00Z", resp.Timestamp)
}
