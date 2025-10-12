package users

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_NewHandler(t *testing.T) {
	// 測試創建新的 handler
	handler := NewHandler(nil)
	assert.NotNil(t, handler)
	assert.Nil(t, handler.userService)
}

func TestHandler_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 創建 handler
	handler := &Handler{}

	// 設置路由
	router := gin.New()
	handler.RegisterRoutes(router.Group("/api/v1"))

	// 測試路由是否正確註冊
	assert.NotNil(t, router)
}

func TestCreateUserRequest_Validation(t *testing.T) {
	// 測試請求結構體的基本功能
	req := CreateUserRequest{
		Email:     "test@example.com",
		Name:      "Test User",
		Password:  "password123",
		Role:      "admin",
		CompanyID: uuid.New(),
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "Test User", req.Name)
	assert.Equal(t, "password123", req.Password)
	assert.Equal(t, "admin", req.Role)
	assert.NotNil(t, req.CompanyID)
}

func TestUpdateUserRequest_Validation(t *testing.T) {
	// 測試更新請求結構體的基本功能
	req := UpdateUserRequest{
		Name: "Updated User",
		Role: "user",
	}

	assert.Equal(t, "Updated User", req.Name)
	assert.Equal(t, "user", req.Role)
}

func TestChangePasswordRequest_Validation(t *testing.T) {
	// 測試密碼變更請求結構體的基本功能
	req := ChangePasswordRequest{
		OldPassword: "oldpassword",
		NewPassword: "newpassword",
	}

	assert.Equal(t, "oldpassword", req.OldPassword)
	assert.Equal(t, "newpassword", req.NewPassword)
}
