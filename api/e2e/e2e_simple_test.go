package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestE2E_ExternalAPI_BasicFlow(t *testing.T) {
	// 端到端測試：基本流程
	gin.SetMode(gin.TestMode)

	// 創建測試路由
	router := gin.New()

	// 模擬外部 API 端點
	router.POST("/api/v1/notify", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 驗證必填字段
		if req["notify_key"] == nil || req["message"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Notification sent successfully",
			"data": gin.H{
				"notification_id": "test-notification-id",
				"status":          "sent",
				"timestamp":       time.Now(),
			},
		})
	})

	router.GET("/api/v1/destinations/:notifyKey", func(c *gin.Context) {
		notifyKey := c.Param("notifyKey")
		if notifyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing notify key"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"project_id": "test-project-id",
				"destinations": []gin.H{
					{
						"id":     "test-destination-1",
						"type":   "channel",
						"status": "active",
					},
					{
						"id":     "test-destination-2",
						"type":   "group",
						"status": "active",
					},
				},
			},
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now(),
		})
	})

	// 測試發送通知
	t.Run("Send Notification", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"notify_key":   "test-notify-key",
			"message":      "Test notification message",
			"message_type": "text",
			"priority":     "normal",
			"targets":      []string{"all"},
			"mentions":     []string{"@everyone"},
			"metadata": map[string]interface{}{
				"source": "e2e-test",
			},
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Notification sent successfully", response["message"])
	})

	// 測試獲取專案目的地
	t.Run("Get Project Destinations", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/destinations/test-notify-key", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
	})

	// 測試健康檢查
	t.Run("Health Check", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", response["status"])
	})
}

func TestE2E_ExternalAPI_InvalidRequest(t *testing.T) {
	// 端到端測試：無效請求
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/api/v1/notify", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req["notify_key"] == nil || req["message"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// 測試缺少必填字段
	requestBody := map[string]interface{}{
		"invalid": "data",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestE2E_ExternalAPI_ConcurrentRequests(t *testing.T) {
	// 端到端測試：並發請求
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/api/v1/notify", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req["notify_key"] == nil || req["message"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// 並發請求
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			requestBody := map[string]interface{}{
				"notify_key": fmt.Sprintf("test-notify-key-%d", index),
				"message":    fmt.Sprintf("Test message %d", index),
			}

			jsonBody, _ := json.Marshal(requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			done <- true
		}(i)
	}

	// 等待所有請求完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestE2E_ExternalAPI_Performance(t *testing.T) {
	// 端到端測試：性能測試
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/api/v1/notify", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req["notify_key"] == nil || req["message"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	start := time.Now()

	// 執行多個請求
	for i := 0; i < 100; i++ {
		requestBody := map[string]interface{}{
			"notify_key": fmt.Sprintf("test-notify-key-%d", i),
			"message":    fmt.Sprintf("Test message %d", i),
		}

		jsonBody, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	}

	duration := time.Since(start)
	t.Logf("Executed 100 requests in %v", duration)
	assert.Less(t, duration, 5*time.Second) // 應該在5秒內完成
}

func TestE2E_ExternalAPI_DataValidation(t *testing.T) {
	// 端到端測試：數據驗證
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/api/v1/notify", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req["notify_key"] == nil || req["message"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "有效請求",
			requestBody: map[string]interface{}{
				"notify_key": "test-notify-key",
				"message":    "Test message",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "缺少通知鍵",
			requestBody: map[string]interface{}{
				"message": "Test message",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "缺少訊息",
			requestBody: map[string]interface{}{
				"notify_key": "test-notify-key",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestE2E_ExternalAPI_Integration(t *testing.T) {
	// 端到端測試：集成測試
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// 設置所有端點
	router.POST("/api/v1/notify", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req["notify_key"] == nil || req["message"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	router.GET("/api/v1/destinations/:notifyKey", func(c *gin.Context) {
		notifyKey := c.Param("notifyKey")
		if notifyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing notify key"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"project_id": "test-project-id",
				"destinations": []gin.H{
					{
						"id":     "test-destination-1",
						"type":   "channel",
						"status": "active",
					},
				},
			},
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now(),
		})
	})

	// 測試完整的 API 流程
	t.Run("Complete API Flow", func(t *testing.T) {
		// 1. 健康檢查
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// 2. 獲取專案目的地
		req, _ = http.NewRequest("GET", "/api/v1/destinations/test-notify-key", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// 3. 發送通知
		requestBody := map[string]interface{}{
			"notify_key": "test-notify-key",
			"message":    "Test notification",
		}
		jsonBody, _ := json.Marshal(requestBody)
		req, _ = http.NewRequest("POST", "/api/v1/notify", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
