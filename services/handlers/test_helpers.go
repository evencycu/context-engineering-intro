package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService 是服務層的模擬接口
type MockService interface {
	// 這裡可以定義需要模擬的方法
}

// TestHelper 提供測試輔助功能
type TestHelper struct {
	router *gin.Engine
}

// NewTestHelper 創建新的測試輔助器
func NewTestHelper() *TestHelper {
	gin.SetMode(gin.TestMode)
	return &TestHelper{
		router: gin.New(),
	}
}

// SetupRouter 設置測試路由
func (th *TestHelper) SetupRouter(handler interface{}) {
	// 這裡可以根據不同的 handler 設置路由
}

// MakeRequest 發送 HTTP 請求
func (th *TestHelper) MakeRequest(method, url string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	th.router.ServeHTTP(w, req)
	return w
}

// AssertJSONResponse 斷言 JSON 回應
func (th *TestHelper) AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	assert.Equal(t, expectedStatus, w.Code)

	if expectedBody != nil {
		var actualBody, expectedJSON interface{}
		json.Unmarshal(w.Body.Bytes(), &actualBody)
		expectedBytes, _ := json.Marshal(expectedBody)
		json.Unmarshal(expectedBytes, &expectedJSON)
		assert.Equal(t, expectedJSON, actualBody)
	}
}

// MockRepository 模擬資料庫操作
type MockRepository struct {
	mock.Mock
}

// MockService 模擬服務層
type MockServiceLayer struct {
	mock.Mock
}

// 通用測試數據
var (
	TestCompanyID = "test-company-id"
	TestUserID    = "test-user-id"
	TestProjectID = "test-project-id"
	TestBotID     = "test-bot-id"
)

// 測試用的 JSON 數據
var (
	TestCompanyJSON = map[string]interface{}{
		"name":        "Test Company",
		"description": "Test Company Description",
		"status":      "active",
	}

	TestUserJSON = map[string]interface{}{
		"email":      "test@example.com",
		"name":       "Test User",
		"role":       "admin",
		"status":     "active",
		"company_id": TestCompanyID,
	}

	TestProjectJSON = map[string]interface{}{
		"name":        "Test Project",
		"description": "Test Project Description",
		"company_id":  TestCompanyID,
		"status":      "active",
	}
)
