# API 概覽

Teams Notify Center 是一個基於 Go 和 Gin 框架的 RESTful API Server，提供 Teams Bot 送訊息功能到 Teams 安裝 Teams Bot 的對象。

## 功能特點

- **多目標支援**: 一個目的地可以包含多個 Teams 目標（channel、person、chatgroup）
- **Bot 管理**: 支援平台 Bot（第三方 Bot 已暫時停用）
- **智能路由**: 根據 Tag 選擇符合的 Bot
- **認證授權**: API key委外APIM服務，這個服務不處理
- **日誌監控**: 完整的請求日誌和錯誤處理

## API 架構

### 基礎 URL

```
http://localhost:8080/api/v1
```


## 主要 API 端點

### 核心管理 API

- **公司管理** (`/companies`) - 公司資料管理
- **用戶管理** (`/users`) - 用戶帳號和權限管理
- **專案管理** (`/projects`) - 通知專案管理
- **Bot 管理** (`/bots`) - Teams Bot 管理
- **目的地管理** (`/destinations`) - 通知目的地管理
- **通知管理** (`/notifications`) - 通知發送和管理

### 特殊功能 API

- **Provision API** (`/provision`) - 一鍵建立專案和目的地
- **External API** (`/external/notify`) - 外部系統通知發送
- **Queue Management API** (`/queue`) - 內部佇列管理和監控
- **訊息管理** (`/messages`) - 訊息歷史和狀態查詢

## 快速開始

### 1. 建立專案和目的地

使用 [Provision API](provision-api.md) 快速建立：

```bash
curl -X POST http://localhost:8080/api/v1/provision \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "e4f160f4-4917-443d-86e9-29071967b768",
    "created_by": "11111111-1111-1111-1111-111111111111",
    "project_name": "我的通知專案",
    "project_description": "用於發送系統通知",
    "teams_tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
    "targets": [
      {
        "type": "channel",
        "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
        "conversation_id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2"
      }
    ]
  }'
```

### 2. 發送通知

使用 [External API](external-api.md) 發送通知：

```bash
curl -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38",
    "message": "Hello World!",
    "message_type": "text",
    "priority": "normal",
    "targets": ["all"]
  }'
```

## 錯誤處理

### 標準錯誤格式

```json
{
  "error": "錯誤描述",
  "details": "詳細錯誤信息（可選）"
}
```

### 常見錯誤碼

- `400 Bad Request`: 請求參數錯誤
- `401 Unauthorized`: 未認證
- `403 Forbidden`: 權限不足
- `404 Not Found`: 資源不存在
- `409 Conflict`: 資源衝突
- `429 Too Many Requests`: 請求過於頻繁
- `500 Internal Server Error`: 內部服務器錯誤

## 分頁

### 分頁參數

- `limit`: 每頁數量（1-100，預設 10）
- `offset`: 偏移量（預設 0）
- `sort_by`: 排序欄位（預設 created_at）
- `order`: 排序方向（asc/desc，預設 desc）

### 分頁響應格式

```json
{
  "data": [...],
  "pagination": {
    "total": 100,
    "limit": 10,
    "offset": 0,
    "pages": 10
  }
}
```

## 健康檢查

```http
GET /health
```

響應：

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

## 詳細文檔

- [Provision API](provision-api.md) - 專案和目的地管理
- [External API](external-api.md) - 外部系統通知發送
- [Queue Management API](queue-management-api.md) - 內部佇列管理
- [OpenAPI 規範](openapi/teams-notification-api.yaml) - 完整 API 規範

## 注意事項

1. **認證**: 大部分 API 需要 JWT 或 API Key 認證
2. **目標驗證**: 創建目的地時會驗證 Teams 目標的有效性
3. **速率限制**: 預設限制 100 請求/秒，可配置
4. **日誌**: 所有請求都會被記錄，包含敏感數據的請求會特別標記
5. **錯誤處理**: 統一的錯誤處理機制，提供詳細的錯誤信息
