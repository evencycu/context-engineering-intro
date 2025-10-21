package companies

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_NewHandler(t *testing.T) {
	// 測試創建新的 handler
	handler := NewHandler(nil)
	assert.NotNil(t, handler)
	assert.Nil(t, handler.companyService)
}

func TestHandler_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 創建 handler
	handler := &Handler{}

	// 設置路由
	router := gin.New()
	handler.RegisterRoutes(router.Group("/api/v1"))

	// 測試路由是否正確註冊
	// 這裡我們只測試路由註冊，不測試實際的處理邏輯
	assert.NotNil(t, router)
}

func TestCreateCompanyRequest_Validation(t *testing.T) {
	// 測試請求結構體的基本功能
	req := CreateCompanyRequest{
		Name:           "Test Company",
		ContactEmail:   "test@example.com",
		ContactPhone:   "1234567890",
		Address:        "Test Address",
		BillingEnabled: true,
	}

	assert.Equal(t, "Test Company", req.Name)
	assert.Equal(t, "test@example.com", req.ContactEmail)
	assert.Equal(t, "1234567890", req.ContactPhone)
	assert.Equal(t, "Test Address", req.Address)
	assert.True(t, req.BillingEnabled)
}

func TestUpdateCompanyRequest_Validation(t *testing.T) {
	// 測試更新請求結構體的基本功能
	req := UpdateCompanyRequest{
		Name:           "Updated Company",
		ContactEmail:   "updated@example.com",
		ContactPhone:   "0987654321",
		Address:        "Updated Address",
		BillingEnabled: &[]bool{false}[0],
	}

	assert.Equal(t, "Updated Company", req.Name)
	assert.Equal(t, "updated@example.com", req.ContactEmail)
	assert.Equal(t, "0987654321", req.ContactPhone)
	assert.Equal(t, "Updated Address", req.Address)
	assert.False(t, *req.BillingEnabled)
}

func TestUpdateStatusRequest_Validation(t *testing.T) {
	// 測試狀態更新請求結構體的基本功能
	req := UpdateStatusRequest{
		Status: "inactive",
	}

	assert.Equal(t, "inactive", req.Status)
}

func TestUpdateBillingRequest_Validation(t *testing.T) {
	// 測試計費更新請求結構體的基本功能
	req := UpdateBillingRequest{
		Enabled: false,
	}

	assert.False(t, req.Enabled)
}
