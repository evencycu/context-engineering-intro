# Teams Notify API Documentation

## 概述

Teams Notify API 是一個基於 Go 和 Gin 框架的 RESTful API，提供 Teams Bot 中台功能，支援多個 Teams 目標的組合通知。

## 功能特點

- **多目標支援**: 一個目的地可以包含多個 Teams 目標（channel、person、chatgroup）
- **Bot 管理**: 支援平台 Bot（第三方 Bot 已暫時停用）
- **智能路由**: 根據條件自動選擇合適的 Bot
- **認證授權**: JWT 和 API Key 雙重認證
- **完整 CRUD**: 所有基礎模型的完整 CRUD 操作
- **日誌監控**: 完整的請求日誌和錯誤處理

## API 架構

### 基礎 URL
```
http://localhost:8080/api/v1
```

### 認證方式

#### 1. JWT Token 認證
```bash
Authorization: Bearer <jwt_token>
```

#### 2. API Key 認證
```bash
X-API-Key: <api_key>
```

## API 端點

### 1. 公司管理 (Companies)

#### 創建公司
```http
POST /api/v1/companies
Content-Type: application/json

{
  "name": "示例公司",
  "contactEmail": "contact@example.com",
  "contactPhone": "+886-2-1234-5678",
  "address": "台北市信義區信義路五段7號",
  "billing_enabled": true
}
```

#### 獲取公司列表
```http
GET /api/v1/companies?limit=10&offset=0&search=示例&status=active
```

#### 獲取單一公司
```http
GET /api/v1/companies/{id}
```

#### 更新公司
```http
PUT /api/v1/companies/{id}
Content-Type: application/json

{
  "name": "更新後的公司名稱",
  "contactEmail": "new@example.com"
}
```

#### 刪除公司
```http
DELETE /api/v1/companies/{id}
```

### 2. 用戶管理 (Users)

#### 創建用戶
```http
POST /api/v1/users
Content-Type: application/json

{
  "companyId": "uuid",
  "email": "user@example.com",
  "name": "張三",
  "role": "admin",
  "password": "password123"
}
```

#### 獲取用戶列表
```http
GET /api/v1/users?limit=10&offset=0&role=admin&status=active
```

#### 生成 API Key
```http
PATCH /api/v1/users/{id}/api-key
```

#### 撤銷 API Key
```http
DELETE /api/v1/users/{id}/api-key
```

### 3. 專案管理 (Projects)

#### 創建專案
```http
POST /api/v1/projects
Content-Type: application/json

{
  "companyId": "uuid",
  "notifyKey": "project_alpha",
  "description": "Alpha 專案通知",
  "dailyLimit": 1000,
  "monthlyLimit": 30000,
  "priority": "normal"
}
```

#### 更新專案限制
```http
PATCH /api/v1/projects/{id}/limits
Content-Type: application/json

{
  "dailyLimit": 2000,
  "monthlyLimit": 60000
}
```

### 4. Bot 管理 (Bots)

#### 創建平台 Bot
```http
POST /api/v1/bots/platform
Content-Type: application/json

{
  "name": "主要通知 Bot",
  "description": "負責所有通知的主要 Bot",
  "appId": "bot-app-id",
  "appPassword": "bot-password",
  "tenantId": "tenant-id",
  "webhookUrl": "https://webhook.example.com",
  "capabilities": {
    "send_message": true,
    "send_file": true,
    "send_adaptive_card": true
  },
  "rate_limit_per_minute": 100,
  "max_concurrent_requests": 10
}
```

> 注意：第三方 Bot 相關端點目前已停用（已自伺服器路由移除）。如需恢復，請先完成相依的服務/儲存層並重新開啟路由。

#### 測試 Bot 連接
```http
POST /api/v1/bots/platform/{id}/test
```

### 5. 目的地管理 (Destinations)

#### 創建目的地（多目標）
```http
POST /api/v1/destinations
Content-Type: application/json

{
  "projectId": "uuid",
  "name": "開發團隊通知群組",
  "description": "包含所有開發相關的 Teams 目標",
  "teamsTenantId": "tenant-123",
  "targets": [
    {
      "type": "channel",
      "team_id": "team-1",
      "channel_id": "channel-dev-general",
      "display_name": "開發討論區",
      "description": "一般開發討論"
    },
    {
      "type": "person",
      "user_id": "user-manager-1",
      "display_name": "張經理",
      "description": "技術經理"
    },
    {
      "type": "chatgroup",
      "group_id": "group-oncall",
      "display_name": "On-call 群組",
      "description": "值班人員群組"
    }
  ],
  "botId": "bot-uuid",
  "bot_type": "platform"
}
```

#### 更新目的地目標
```http
PATCH /api/v1/destinations/{id}/targets
Content-Type: application/json

{
  "targets": [
    {
      "type": "channel",
      "team_id": "team-2",
      "channel_id": "channel-qa-testing",
      "display_name": "QA 測試",
      "description": "測試相關討論"
    }
  ]
}
```

#### 驗證目標
```http
POST /api/v1/destinations/{id}/validate
Content-Type: application/json

{
  "targets": [
    {
      "type": "channel",
      "team_id": "team-1",
      "channel_id": "channel-dev"
    }
  ]
}
```

#### 搜尋目的地
```http
GET /api/v1/destinations/search?target_type=channel&target_id=channel-dev
GET /api/v1/destinations/search?team_id=team-1
GET /api/v1/destinations/search?user_id=user-123
```

### 6. 通知管理 (Notifications)

#### 發送通知
```http
POST /api/v1/notifications
Content-Type: application/json

{
  "projectId": "uuid",
  "messageType": "text",
  "content": "系統維護通知：將於今晚 10:00-11:00 進行系統維護",
  "mentions": ["@張經理", "@李組長"],
  "priority": "high",
  "metadata": {
    "maintenance_id": "maint-001",
    "scheduled_time": "2024-01-15T22:00:00Z"
  },
  "destinations": ["dest-uuid-1", "dest-uuid-2"]
}
```

#### 發送檔案通知
```http
POST /api/v1/notifications
Content-Type: application/json

{
  "projectId": "uuid",
  "messageType": "file",
  "content": "請查看附件中的報告",
  "attachment": {
    "file_name": "monthly_report.pdf",
    "file_url": "https://files.example.com/reports/monthly_report.pdf",
    "file_size": 1024000,
    "mime_type": "application/pdf"
  },
  "priority": "normal",
  "destinations": ["dest-uuid-1"]
}
```

#### 發送自適應卡片
```http
POST /api/v1/notifications
Content-Type: application/json

{
  "projectId": "uuid",
  "messageType": "adaptive_card",
  "content": "系統警報",
  "adaptive_card": {
    "type": "AdaptiveCard",
    "version": "1.3",
    "body": [
      {
        "type": "TextBlock",
        "text": "系統警報",
        "size": "Large",
        "weight": "Bolder",
        "color": "Attention"
      },
      {
        "type": "TextBlock",
        "text": "CPU 使用率過高：95%",
        "wrap": true
      }
    ],
    "actions": [
      {
        "type": "Action.OpenUrl",
        "title": "查看詳情",
        "url": "https://monitoring.example.com/alerts/123"
      }
    ]
  },
  "priority": "urgent",
  "destinations": ["dest-uuid-1"]
}
```

#### 重試失敗的通知
```http
POST /api/v1/notifications/{id}/retry
```

#### 取消待發送的通知
```http
DELETE /api/v1/notifications/{id}
```

#### 按日期範圍查詢通知
```http
GET /api/v1/notifications/date-range?start_date=2024-01-01T00:00:00Z&end_date=2024-01-31T23:59:59Z
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

## 環境變數

```bash
PORT=8080                    # 服務端口
ENVIRONMENT=development      # 環境（development/production）
JWT_SECRET=your-secret-key   # JWT 密鑰
```

## 啟動服務

```bash
# 安裝依賴
go mod tidy

# 啟動服務
go run cmd/server/main.go
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

## 注意事項

1. **認證**: 大部分 API 需要 JWT 或 API Key 認證
2. **目標驗證**: 創建目的地時會驗證 Teams 目標的有效性
3. **速率限制**: 預設限制 100 請求/秒，可配置
4. **日誌**: 所有請求都會被記錄，包含敏感數據的請求會特別標記
5. **錯誤處理**: 統一的錯誤處理機制，提供詳細的錯誤信息
