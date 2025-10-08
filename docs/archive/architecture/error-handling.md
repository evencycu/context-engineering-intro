# 錯誤處理架構

本文檔說明 Teams Notification API 的錯誤處理架構和標準化錯誤響應。

## 概述

系統採用統一的錯誤處理機制，提供一致的錯誤響應格式和適當的 HTTP 狀態碼。

## 錯誤分類

### 1. 認證錯誤 (4xx)
- **UNAUTHORIZED**: 未授權訪問
- **INVALID_TOKEN**: 無效或格式錯誤的 Token
- **TOKEN_EXPIRED**: Token 已過期
- **INVALID_API_KEY**: 無效的 API 金鑰

### 2. Teams API 錯誤 (4xx/5xx)
- **TEAMS_AUTH_FAILED**: Teams 認證失敗
- **TEAMS_API_ERROR**: Teams API 調用錯誤
- **TEAMS_RATE_LIMIT**: Teams API 速率限制

### 3. 數據庫錯誤 (5xx)
- **DATABASE_ERROR**: 數據庫操作失敗
- **RECORD_NOT_FOUND**: 記錄不存在
- **DUPLICATE_RECORD**: 重複記錄

### 4. 驗證錯誤 (4xx)
- **VALIDATION_FAILED**: 輸入驗證失敗
- **INVALID_INPUT**: 無效輸入
- **MISSING_FIELD**: 缺少必要欄位

### 5. 系統錯誤 (5xx)
- **INTERNAL_ERROR**: 內部服務錯誤
- **SERVICE_UNAVAILABLE**: 服務不可用
- **TIMEOUT**: 請求超時

## 錯誤響應格式

### 標準錯誤響應
```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "Input validation failed",
    "details": "Field 'email' is required",
    "cause": "Missing required field"
  },
  "timestamp": "2025-10-07T01:47:31Z",
  "request_id": "req_123456789"
}
```

### 簡化錯誤響應
```json
{
  "success": false,
  "error": "Invalid request format",
  "code": "VALIDATION_FAILED",
  "details": "Field 'email' is required",
  "timestamp": "2025-10-07T01:47:31Z",
  "request_id": "req_123456789"
}
```

## HTTP 狀態碼映射

| 錯誤代碼 | HTTP 狀態碼 | 描述 |
|---------|-------------|------|
| UNAUTHORIZED | 401 | 未授權 |
| INVALID_TOKEN | 401 | 無效 Token |
| TOKEN_EXPIRED | 401 | Token 過期 |
| INVALID_API_KEY | 401 | 無效 API 金鑰 |
| VALIDATION_FAILED | 400 | 驗證失敗 |
| INVALID_INPUT | 400 | 無效輸入 |
| MISSING_FIELD | 400 | 缺少欄位 |
| RECORD_NOT_FOUND | 404 | 記錄不存在 |
| DUPLICATE_RECORD | 409 | 重複記錄 |
| TEAMS_RATE_LIMIT | 429 | 速率限制 |
| TIMEOUT | 408 | 請求超時 |
| DATABASE_ERROR | 500 | 數據庫錯誤 |
| TEAMS_AUTH_FAILED | 500 | Teams 認證失敗 |
| TEAMS_API_ERROR | 500 | Teams API 錯誤 |
| INTERNAL_ERROR | 500 | 內部錯誤 |
| SERVICE_UNAVAILABLE | 503 | 服務不可用 |

## 錯誤處理流程

### 1. 錯誤捕獲
```go
func (h *Handler) CreateCompany(c *gin.Context) {
    var req CreateCompanyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 捕獲驗證錯誤
        appErr := errors.NewValidationError("request", err.Error())
        c.JSON(appErr.HTTPStatus, gin.H{"error": appErr})
        return
    }
    
    // 業務邏輯處理
    company, err := h.service.CreateCompany(c.Request.Context(), req)
    if err != nil {
        // 捕獲業務錯誤
        appErr := errors.NewDatabaseError("create company", err)
        c.JSON(appErr.HTTPStatus, gin.H{"error": appErr})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{"data": company})
}
```

### 2. 錯誤記錄
```go
func (h *Handler) logError(c *gin.Context, err *errors.AppError) {
    h.logger.WithFields(logrus.Fields{
        "error_code": err.Code,
        "error_message": err.Message,
        "request_id": c.GetString("request_id"),
        "user_id": c.GetString("user_id"),
        "path": c.Request.URL.Path,
        "method": c.Request.Method,
    }).Error("Request failed")
}
```

### 3. 錯誤追蹤
```go
func (h *Handler) trackError(c *gin.Context, err *errors.AppError) {
    // 發送到監控系統
    metrics.IncrementErrorCounter(err.Code)
    
    // 記錄到審計日誌
    audit.LogError(c.Request.Context(), err)
}
```

## 最佳實踐

### 1. 錯誤訊息設計
- **用戶友好**: 提供清晰的錯誤訊息
- **技術詳細**: 包含技術細節供開發者調試
- **可追蹤**: 包含請求 ID 和時間戳

### 2. 錯誤分類
- **客戶端錯誤**: 4xx 狀態碼，用戶可修正
- **服務端錯誤**: 5xx 狀態碼，需要系統修復
- **暫時性錯誤**: 可重試的錯誤
- **永久性錯誤**: 不可重試的錯誤

### 3. 錯誤監控
- **錯誤率監控**: 追蹤各類錯誤的發生率
- **錯誤趨勢**: 分析錯誤的變化趨勢
- **告警機制**: 設定錯誤閾值告警

## 配置錯誤處理

### 環境變數驗證
```go
func ValidateConfig(config *Config) error {
    var errors []string
    
    if config.Teams.AppID == "" {
        errors = append(errors, "TEAMS_BOT_APP_ID is required")
    }
    
    if config.Teams.AppPassword == "" {
        errors = append(errors, "TEAMS_BOT_APP_PASSWORD is required")
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, ", "))
    }
    
    return nil
}
```

### 配置錯誤響應
```json
{
  "valid": false,
  "errors": [
    {
      "field": "TEAMS_BOT_APP_PASSWORD",
      "message": "Environment variable is required",
      "code": "MISSING_REQUIRED_FIELD"
    }
  ],
  "warnings": [
    {
      "field": "RATE_LIMIT",
      "message": "Using default value",
      "suggestion": "Consider setting explicit rate limit"
    }
  ]
}
```

## 監控和告警

### 錯誤指標
- **錯誤率**: 錯誤請求 / 總請求數
- **錯誤分佈**: 各類錯誤的數量分佈
- **錯誤趨勢**: 錯誤率的時間趨勢
- **響應時間**: 錯誤響應的平均時間

### 告警規則
- **錯誤率告警**: 錯誤率超過 5% 時告警
- **關鍵錯誤告警**: 認證錯誤或數據庫錯誤立即告警
- **服務可用性告警**: 服務不可用時告警

## 故障排除

### 常見錯誤
1. **TEAMS_AUTH_FAILED**: 檢查 Teams Bot 配置
2. **DATABASE_ERROR**: 檢查數據庫連接
3. **VALIDATION_FAILED**: 檢查請求格式
4. **RATE_LIMIT**: 檢查請求頻率

### 調試步驟
1. 檢查錯誤日誌
2. 驗證配置設定
3. 測試外部服務連接
4. 檢查系統資源使用

## 更新日誌

### v1.0.0
- 實現統一的錯誤處理機制
- 添加標準化錯誤響應格式
- 建立錯誤分類和狀態碼映射
- 實現錯誤監控和告警
