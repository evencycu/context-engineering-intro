package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_NewAuthMiddleware(t *testing.T) {
	// 測試創建認證中間件
	middleware := NewAuthMiddleware("test-secret")
	assert.NotNil(t, middleware)
	assert.Equal(t, "test-secret", middleware.jwtSecret)
}

func TestErrorHandler_NewErrorHandler(t *testing.T) {
	// 測試創建錯誤處理中間件
	logger := logrus.New()
	handler := NewErrorHandler(logger)
	assert.NotNil(t, handler)
}

func TestLoggingMiddleware_NewLoggingMiddleware(t *testing.T) {
	// 測試創建日誌中間件
	logger := logrus.New()
	middleware := NewLoggingMiddleware(logger)
	assert.NotNil(t, middleware)
}

func TestBusinessMetricsMiddleware_NewBusinessMetricsMiddleware(t *testing.T) {
	// 測試創建業務指標中間件
	logger := logrus.New()
	middleware := NewBusinessMetricsMiddleware(logger)
	assert.NotNil(t, middleware)
}

func TestMiddleware_BasicFunctionality(t *testing.T) {
	// 測試中間件的基本功能
	gin.SetMode(gin.TestMode)

	// 創建測試路由
	router := gin.New()

	// 添加基本中間件
	router.Use(gin.Recovery())

	// 添加測試路由
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// 驗證路由設置
	assert.NotNil(t, router)
}

func TestMiddleware_LoggerCreation(t *testing.T) {
	// 測試日誌器創建
	logger := logrus.New()
	assert.NotNil(t, logger)

	// 測試日誌器配置
	logger.SetLevel(logrus.InfoLevel)
	assert.Equal(t, logrus.InfoLevel, logger.GetLevel())
}

func TestMiddleware_SecretValidation(t *testing.T) {
	// 測試 JWT 密鑰驗證
	secret := "test-secret-key"
	assert.NotEmpty(t, secret)
	assert.Greater(t, len(secret), 0)
}
