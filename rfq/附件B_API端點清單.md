# 附件 B：API 端點清單
# API Endpoints Specification

**RFQ 編號**: RFQ-2025-TEAMS-NOTIFY-001  
**版本**: v1.0  
**日期**: 2025年10月23日

---

## 目錄

1. [系統端點](#1-系統端點)
2. [外部 API 端點](#2-外部-api-端點)
3. [內部 API - 公司管理](#3-內部-api---公司管理)
4. [內部 API - 用戶管理](#4-內部-api---用戶管理)
5. [內部 API - 專案管理](#5-內部-api---專案管理)
6. [內部 API - Bot 管理](#6-內部-api---bot-管理)
7. [內部 API - 目的地管理](#7-內部-api---目的地管理)
8. [內部 API - 通知管理](#8-內部-api---通知管理)
9. [內部 API - 計費管理](#9-內部-api---計費管理)
10. [內部 API - 檔案管理](#10-內部-api---檔案管理)
11. [內部 API - 監控管理](#11-內部-api---監控管理)
12. [內部 API - 佇列管理](#12-內部-api---佇列管理)

---

## 圖示說明

- ✅ **P0**: 必須實作（高優先級）
- ⚠️ **P1**: 建議實作（中優先級）
- 💡 **P2**: 可選實作（低優先級）
- 🔒 需要認證

---

## 1. 系統端點

### 1.1 健康檢查

#### GET /health ✅

**說明**: 基本健康檢查端點

**優先級**: P0

**認證**: 不需要

**請求參數**: 無

**回應範例**:
```json
{
  "status": "healthy",
  "timestamp": "2025-10-23T10:00:00Z",
  "uptime": "72h30m15s",
  "version": "v1.0.0"
}
```

---

#### GET /api/v1/metrics ⚠️

**說明**: 系統指標端點

**優先級**: P1

**認證**: 不需要（建議限制 IP 訪問）

**回應範例**:
```json
{
  "timestamp": "2025-10-23T10:00:00Z",
  "system": {
    "memory_usage": 256.5,
    "cpu_usage": 35.2,
    "goroutines": 150
  },
  "token_cache": {
    "hits": 9500,
    "misses": 500,
    "hit_rate": 0.95
  }
}
```

---

## 2. 外部 API 端點

### 2.1 通知發送

#### POST /api/v1/notify ✅

**說明**: 外部系統發送通知（核心功能）

**優先級**: P0

**認證**: 需要（API Key 或 JWT）

**請求範例**:
```json
{
  "notifyKey": "my-project-key",
  "message": "系統通知：資料庫備份完成",
  "messageType": "text",
  "priority": "normal",
  "targets": ["all"],
  "mentions": ["user@example.com"],
  "metadata": {
    "source": "backup-system",
    "job_id": "backup-001"
  }
}
```

**欄位說明**:
- `notifyKey` (必填): 專案識別碼
- `message` (必填): 訊息內容，最長 4000 字元
- `messageType` (選填): text | file | adaptive_card，預設 text
- `priority` (選填): low | normal | high，預設 normal
- `targets` (選填): 目標陣列，預設 ["all"]
  - "all": 發送到所有目的地
  - "conversation_id": 發送到特定對話
  - "email": 發送到特定用戶（personal 類型）
- `mentions` (選填): @提及的用戶 email 陣列
- `metadata` (選填): 自訂資料

**回應範例** (202 Accepted):
```json
{
  "success": true,
  "data": {
    "notificationId": "550e8400-e29b-41d4-a716-446655440000",
    "status": "processing",
    "message": "Notification queued successfully",
    "destinationsCount": 5,
    "estimated_delivery_time": "5-10 seconds"
  }
}
```

**錯誤回應** (400/404):
```json
{
  "error": "INVALID_NOTIFY_KEY",
  "message": "Notify key 'my-project-key' not found",
  "timestamp": "2025-10-23T10:00:00Z"
}
```

---

#### GET /api/v1/destinations/{notifyKey} ⚠️

**說明**: 查詢專案的所有目的地

**優先級**: P1

**認證**: 需要

**回應範例**:
```json
{
  "notifyKey": "my-project-key",
  "destinations": [
    {
      "id": "uuid",
      "type": "channel",
      "targets": [
        {
          "type": "channel",
          "conversation_id": "19:xxx@thread.tacv2",
          "display_name": "General",
          "tenant_id": "tenant-uuid"
        }
      ],
      "status": "active"
    }
  ]
}
```

---

### 2.2 專案開通

#### POST /api/v1/provision ⚠️

**說明**: 一鍵建立專案（簡化流程）

**優先級**: P1

**認證**: 需要

**請求範例**:
```json
{
  "notifyKey": "new-project",
  "companyId": "550e8400-e29b-41d4-a716-446655440000",
  "projectName": "新專案",
  "description": "專案描述"
}
```

**回應範例**:
```json
{
  "notifyKey": "new-project",
  "projectId": "uuid",
  "status": "active",
  "createdAt": "2025-10-23T10:00:00Z"
}
```

---

#### GET /api/v1/provision/{notifyKey} ⚠️

**說明**: 查詢專案開通狀態

**優先級**: P1

**認證**: 需要

**回應範例**:
```json
{
  "notifyKey": "new-project",
  "projectName": "新專案",
  "status": "active",
  "destinationsCount": 3
}
```

---

#### PUT /api/v1/provision/{notifyKey} 💡

**說明**: 更新專案資訊

**優先級**: P2

**認證**: 需要

---

#### POST /api/v1/provision/{notifyKey}/enable 💡

**說明**: 啟用專案

**優先級**: P2

**認證**: 需要

---

#### POST /api/v1/provision/{notifyKey}/disable 💡

**說明**: 停用專案

**優先級**: P2

**認證**: 需要

---

### 2.3 Bot 訊息接收

#### POST /api/v1/messages ✅

**說明**: 接收 Teams Bot Framework 訊息

**優先級**: P0

**認證**: Bot Framework Token

**請求範例** (Bot Framework Activity):
```json
{
  "type": "message",
  "id": "xxx",
  "timestamp": "2025-10-23T10:00:00Z",
  "channelId": "msteams",
  "from": {
    "id": "29:user-id",
    "name": "John Doe"
  },
  "conversation": {
    "id": "19:xxx@thread.tacv2"
  },
  "text": "Hello Bot"
}
```

**回應範例**:
```json
{
  "success": true,
  "message": "Message processed"
}
```

---

#### POST /api/v1/messages/test ⚠️

**說明**: 測試發送 Proactive Message

**優先級**: P1

**認證**: 需要

**請求範例**:
```json
{
  "conversationId": "19:xxx@thread.tacv2",
  "message": "Test message"
}
```

---

## 3. 內部 API - 公司管理

### 3.1 公司 CRUD

#### POST /internal/v1/companies ✅

**說明**: 創建公司

**優先級**: P0

**認證**: 需要 (Admin only)

**請求範例**:
```json
{
  "name": "測試公司",
  "contactEmail": "contact@company.com",
  "contactPhone": "+886-2-1234-5678",
  "address": "台北市信義區信義路五段7號",
  "status": "active",
  "billingEnabled": true
}
```

**回應範例** (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "測試公司",
  "contactEmail": "contact@company.com",
  "contactPhone": "+886-2-1234-5678",
  "address": "台北市信義區信義路五段7號",
  "status": "active",
  "billingEnabled": true,
  "createdAt": "2025-10-23T10:00:00Z",
  "updatedAt": "2025-10-23T10:00:00Z"
}
```

---

#### GET /internal/v1/companies ✅

**說明**: 列表查詢公司

**優先級**: P0

**認證**: 需要

**查詢參數**:
- `limit` (int): 每頁筆數，預設 10，最大 100
- `offset` (int): 偏移量，預設 0
- `search` (string): 搜尋關鍵字（名稱、email）

**範例**: `/internal/v1/companies?limit=20&offset=0&search=測試`

**回應範例**:
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "測試公司",
      ...
    }
  ],
  "pagination": {
    "total": 50,
    "limit": 20,
    "offset": 0,
    "pages": 3
  }
}
```

---

#### GET /internal/v1/companies/{id} ✅

**說明**: 查詢單一公司

**優先級**: P0

**認證**: 需要

**回應範例**:
```json
{
  "id": "uuid",
  "name": "測試公司",
  ...
}
```

---

#### PUT /internal/v1/companies/{id} ✅

**說明**: 更新公司資訊

**優先級**: P0

**認證**: 需要 (Admin only)

---

#### DELETE /internal/v1/companies/{id} ⚠️

**說明**: 刪除公司

**優先級**: P1

**認證**: 需要 (Admin only)

---

#### PATCH /internal/v1/companies/{id}/status ⚠️

**說明**: 更新公司狀態

**優先級**: P1

**認證**: 需要 (Admin only)

**請求範例**:
```json
{
  "status": "suspended"
}
```

---

## 4. 內部 API - 用戶管理

### 4.1 用戶 CRUD

#### POST /internal/v1/users ✅

**說明**: 創建用戶

**優先級**: P0

**認證**: 需要 (Admin/Manager)

**請求範例**:
```json
{
  "companyId": "uuid",
  "email": "user@example.com",
  "name": "張三",
  "role": "user",
  "password": "SecurePass123"
}
```

**回應範例** (201 Created):
```json
{
  "id": "uuid",
  "companyId": "uuid",
  "email": "user@example.com",
  "name": "張三",
  "role": "user",
  "status": "active",
  "createdAt": "2025-10-23T10:00:00Z"
}
```

---

#### GET /internal/v1/users ✅

**說明**: 列表查詢用戶

**優先級**: P0

**認證**: 需要

**查詢參數**: limit, offset, search

---

#### GET /internal/v1/users/{id} ✅

**說明**: 查詢單一用戶

**優先級**: P0

**認證**: 需要

---

#### PUT /internal/v1/users/{id} ✅

**說明**: 更新用戶資訊

**優先級**: P0

**認證**: 需要

---

#### DELETE /internal/v1/users/{id} ⚠️

**說明**: 刪除用戶

**優先級**: P1

**認證**: 需要 (Admin/Manager)

---

#### PATCH /internal/v1/users/{id}/password ✅

**說明**: 變更密碼

**優先級**: P0

**認證**: 需要

**請求範例**:
```json
{
  "oldPassword": "OldPass123",
  "newPassword": "NewPass456"
}
```

---

#### GET /internal/v1/users/company/{companyId} ⚠️

**說明**: 依公司查詢用戶

**優先級**: P1

**認證**: 需要

---

#### GET /internal/v1/users/role/{role} ⚠️

**說明**: 依角色查詢用戶

**優先級**: P1

**認證**: 需要

---

## 5. 內部 API - 專案管理

### 5.1 專案 CRUD

#### POST /internal/v1/projects ✅

**說明**: 創建專案

**優先級**: P0

**認證**: 需要

**請求範例**:
```json
{
  "companyId": "uuid",
  "notifyKey": "my-project-key",
  "description": "專案說明",
  "dailyLimit": 1000,
  "monthlyLimit": 30000,
  "priority": "normal",
  "createdBy": "uuid"
}
```

**欄位說明**:
- `notifyKey`: 必須唯一，建議格式: `{company}-{project}`
- `dailyLimit`: 每日通知配額，0 表示不限制
- `monthlyLimit`: 每月通知配額，0 表示不限制
- `priority`: low | normal | high，影響佇列處理順序

**回應範例** (201 Created):
```json
{
  "id": "uuid",
  "companyId": "uuid",
  "notifyKey": "my-project-key",
  "description": "專案說明",
  "dailyLimit": 1000,
  "monthlyLimit": 30000,
  "priority": "normal",
  "status": "active",
  "createdBy": "uuid",
  "createdAt": "2025-10-23T10:00:00Z"
}
```

---

#### GET /internal/v1/projects ✅

**說明**: 列表查詢專案

**優先級**: P0

**認證**: 需要

---

#### GET /internal/v1/projects/{id} ✅

**說明**: 查詢單一專案

**優先級**: P0

**認證**: 需要

---

#### PUT /internal/v1/projects/{id} ✅

**說明**: 更新專案資訊

**優先級**: P0

**認證**: 需要

---

#### DELETE /internal/v1/projects/{id} ⚠️

**說明**: 刪除專案

**優先級**: P1

**認證**: 需要

---

#### PATCH /internal/v1/projects/{id}/limits ⚠️

**說明**: 更新專案配額

**優先級**: P1

**認證**: 需要

**請求範例**:
```json
{
  "dailyLimit": 2000,
  "monthlyLimit": 60000
}
```

---

#### GET /internal/v1/projects/company/{companyId} ✅

**說明**: 依公司查詢專案

**優先級**: P0

**認證**: 需要

---

#### GET /internal/v1/projects/key/{keyName} ✅

**說明**: 依 Notify Key 查詢專案

**優先級**: P0

**認證**: 需要

---

## 6. 內部 API - Bot 管理

### 6.1 Bot CRUD

#### POST /internal/v1/bots/platform ✅

**說明**: 創建 Teams Bot

**優先級**: P0

**認證**: 需要 (Admin only)

**請求範例**:
```json
{
  "name": "My Teams Bot",
  "appId": "844146d7-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "appPassword": "secret-password",
  "capabilities": ["proactive_messaging", "bot_messaging"],
  "status": "active",
  "description": "Bot 說明"
}
```

**回應範例** (201 Created):
```json
{
  "id": "uuid",
  "name": "My Teams Bot",
  "appId": "844146d7-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "capabilities": ["proactive_messaging", "bot_messaging"],
  "status": "active",
  "description": "Bot 說明",
  "createdAt": "2025-10-23T10:00:00Z"
}
```

**注意**: `appPassword` 不會在回應中出現

---

#### GET /internal/v1/bots/platform ✅

**說明**: 列表查詢 Bot

**優先級**: P0

**認證**: 需要

---

#### GET /internal/v1/bots/platform/{id} ✅

**說明**: 查詢單一 Bot

**優先級**: P0

**認證**: 需要

---

#### PUT /internal/v1/bots/platform/{id} ✅

**說明**: 更新 Bot 資訊

**優先級**: P0

**認證**: 需要 (Admin only)

---

#### DELETE /internal/v1/bots/platform/{id} ⚠️

**說明**: 刪除 Bot

**優先級**: P1

**認證**: 需要 (Admin only)

---

#### PATCH /internal/v1/bots/platform/{id}/status ⚠️

**說明**: 更新 Bot 狀態

**優先級**: P1

**認證**: 需要 (Admin only)

**請求範例**:
```json
{
  "status": "inactive"
}
```

---

#### PATCH /internal/v1/bots/platform/{id}/capabilities 💡

**說明**: 更新 Bot 能力

**優先級**: P2

**認證**: 需要 (Admin only)

---

#### POST /internal/v1/bots/platform/{id}/test ⚠️

**說明**: 測試 Bot 連線

**優先級**: P1

**認證**: 需要

**回應範例**:
```json
{
  "success": true,
  "message": "Bot connection successful",
  "details": {
    "token_obtained": true,
    "token_expires_at": "2025-10-23T11:00:00Z"
  }
}
```

---

#### GET /internal/v1/bots/status/{status} ⚠️

**說明**: 依狀態查詢 Bot

**優先級**: P1

**認證**: 需要

---

## 7. 內部 API - 目的地管理

### 7.1 目的地 CRUD

#### POST /internal/v1/destinations ✅

**說明**: 創建目的地

**優先級**: P0

**認證**: 需要

**請求範例**:
```json
{
  "projectId": "uuid",
  "botId": "uuid",
  "type": "channel",
  "targets": [
    {
      "type": "channel",
      "conversation_id": "19:xxx@thread.tacv2",
      "display_name": "General",
      "tenant_id": "tenant-uuid"
    }
  ],
  "status": "active"
}
```

**Targets 型別說明**:
- **personal**: 個人訊息
  - 必填: `type`, `tenant_id`
  - 二擇一: `conversation_id` 或 `email`
  - 範例:
    ```json
    {
      "type": "personal",
      "email": "user@example.com",
      "display_name": "張三",
      "tenant_id": "uuid"
    }
    ```

- **groupchat**: 群組聊天
  - 必填: `type`, `conversation_id`, `tenant_id`
  - 範例:
    ```json
    {
      "type": "groupchat",
      "conversation_id": "19:xxx@thread.v2",
      "display_name": "專案討論群",
      "tenant_id": "uuid"
    }
    ```

- **channel**: 團隊頻道
  - 必填: `type`, `conversation_id`, `tenant_id`
  - 範例:
    ```json
    {
      "type": "channel",
      "conversation_id": "19:xxx@thread.tacv2",
      "display_name": "General",
      "tenant_id": "uuid"
    }
    ```

**回應範例** (201 Created):
```json
{
  "id": "uuid",
  "projectId": "uuid",
  "botId": "uuid",
  "type": "channel",
  "targets": [...],
  "status": "active",
  "createdAt": "2025-10-23T10:00:00Z"
}
```

---

#### GET /internal/v1/destinations ✅

**說明**: 列表查詢目的地

**優先級**: P0

**認證**: 需要

---

#### GET /internal/v1/destinations/{id} ✅

**說明**: 查詢單一目的地

**優先級**: P0

**認證**: 需要

---

#### PUT /internal/v1/destinations/{id} ✅

**說明**: 更新目的地

**優先級**: P0

**認證**: 需要

---

#### DELETE /internal/v1/destinations/{id} ⚠️

**說明**: 刪除目的地

**優先級**: P1

**認證**: 需要

---

#### PATCH /internal/v1/destinations/{id}/targets ✅

**說明**: 更新目的地 Targets

**優先級**: P0

**認證**: 需要

**請求範例**:
```json
{
  "targets": [
    {
      "type": "channel",
      "conversation_id": "19:new-conversation-id",
      "display_name": "New Channel",
      "tenant_id": "uuid"
    }
  ]
}
```

---

#### POST /internal/v1/destinations/{id}/validate ⚠️

**說明**: 驗證目的地 Targets（測試發送）

**優先級**: P1

**認證**: 需要

**回應範例**:
```json
{
  "valid": true,
  "errors": [],
  "warnings": ["Target 1 has not been tested recently"]
}
```

---

#### GET /internal/v1/destinations/project/{projectId} ✅

**說明**: 依專案查詢目的地

**優先級**: P0

**認證**: 需要

---

#### GET /internal/v1/destinations/bot/{botId} ⚠️

**說明**: 依 Bot 查詢目的地

**優先級**: P1

**認證**: 需要

---

#### GET /internal/v1/destinations/search ⚠️

**說明**: 搜尋目的地（進階篩選）

**優先級**: P1

**認證**: 需要

**查詢參數**:
- `search`: 關鍵字
- `type`: personal | groupchat | channel
- `status`: active | inactive

---

## 8. 內部 API - 通知管理

### 8.1 通知 CRUD

#### POST /internal/v1/notifications ✅

**說明**: 發送通知（內部 API 版本）

**優先級**: P0

**認證**: 需要

**請求範例**:
```json
{
  "projectId": "uuid",
  "senderId": "uuid",
  "messageType": "text",
  "content": "通知內容",
  "mentions": ["user@example.com"],
  "priority": "normal",
  "targets": ["all"],
  "metadata": {
    "custom_field": "value"
  }
}
```

**回應範例** (202 Accepted):
```json
{
  "notificationId": "uuid",
  "status": "processing",
  "message": "Notification queued successfully",
  "destinationsCount": 5,
  "estimated_delivery_time": "5-10 seconds"
}
```

---

#### GET /internal/v1/notifications ✅

**說明**: 列表查詢通知

**優先級**: P0

**認證**: 需要

**查詢參數**:
- `limit`, `offset`, `search`
- `status`: pending | enqueued | sending | sent | failed | cancelled

---

#### GET /internal/v1/notifications/{id} ✅

**說明**: 查詢單一通知

**優先級**: P0

**認證**: 需要

**回應範例**:
```json
{
  "id": "uuid",
  "projectId": "uuid",
  "senderId": "uuid",
  "messageType": "text",
  "content": "通知內容",
  "mentions": [],
  "priority": "normal",
  "status": "sent",
  "destinationsSent": 5,
  "destinationsFailed": 0,
  "totalDestinations": 5,
  "createdAt": "2025-10-23T10:00:00Z",
  "sentAt": "2025-10-23T10:00:15Z"
}
```

---

#### DELETE /internal/v1/notifications/{id} ⚠️

**說明**: 取消通知（僅限 pending/enqueued 狀態）

**優先級**: P1

**認證**: 需要

**回應範例**:
```json
{
  "message": "Notification cancelled successfully"
}
```

---

#### POST /internal/v1/notifications/{id}/retry ⚠️

**說明**: 手動重試失敗的通知

**優先級**: P1

**認證**: 需要

**回應範例**:
```json
{
  "message": "Notification retry initiated"
}
```

---

#### GET /internal/v1/notifications/project/{projectId} ✅

**說明**: 依專案查詢通知

**優先級**: P0

**認證**: 需要

---

#### GET /internal/v1/notifications/sender/{senderId} ⚠️

**說明**: 依發送者查詢通知

**優先級**: P1

**認證**: 需要

---

#### GET /internal/v1/notifications/status/{status} ⚠️

**說明**: 依狀態查詢通知

**優先級**: P1

**認證**: 需要

---

#### GET /internal/v1/notifications/date-range ⚠️

**說明**: 依時間範圍查詢通知

**優先級**: P1

**認證**: 需要

**查詢參數**:
- `start_date`: ISO 8601 格式
- `end_date`: ISO 8601 格式

**範例**: `/internal/v1/notifications/date-range?start_date=2025-10-01T00:00:00Z&end_date=2025-10-31T23:59:59Z`

**回應範例**:
```json
{
  "data": [...],
  "count": 150,
  "date_range": {
    "start_date": "2025-10-01T00:00:00Z",
    "end_date": "2025-10-31T23:59:59Z"
  }
}
```

---

## 9. 內部 API - 計費管理

### 9.1 使用量查詢

#### GET /internal/v1/billing/usage ⚠️

**說明**: 查詢使用量記錄

**優先級**: P1

**認證**: 需要

**查詢參數**:
- `company_id`: 公司 ID
- `project_id`: 專案 ID
- `limit`, `offset`

---

#### GET /internal/v1/billing/usage/summary ⚠️

**說明**: 使用量彙總統計

**優先級**: P1

**認證**: 需要

**查詢參數**:
- `company_id`: 公司 ID
- `project_id`: 專案 ID
- `start_date`: 開始日期
- `end_date`: 結束日期

**回應範例**:
```json
{
  "totalNotifications": 15000,
  "dailyAverage": 500,
  "monthlyTotal": 15000,
  "period": {
    "startDate": "2025-10-01",
    "endDate": "2025-10-31"
  }
}
```

---

#### GET /internal/v1/billing/usage/project/{projectId} ⚠️

**說明**: 查詢專案使用量

**優先級**: P1

**認證**: 需要

---

### 9.2 計費方案

#### GET /internal/v1/billing/plans 💡

**說明**: 列表查詢計費方案

**優先級**: P2

**認證**: 需要

---

#### POST /internal/v1/billing/plans 💡

**說明**: 創建計費方案

**優先級**: P2

**認證**: 需要 (Admin only)

---

#### GET /internal/v1/billing/plans/{id} 💡

**說明**: 查詢單一計費方案

**優先級**: P2

**認證**: 需要

---

#### PUT /internal/v1/billing/plans/{id} 💡

**說明**: 更新計費方案

**優先級**: P2

**認證**: 需要 (Admin only)

---

### 9.3 專案計費

#### GET /internal/v1/billing/project/{projectId} ⚠️

**說明**: 查詢專案計費資訊

**優先級**: P1

**認證**: 需要

**回應範例**:
```json
{
  "projectId": "uuid",
  "planId": "uuid",
  "planName": "Enterprise Plan",
  "currentUsage": 15000,
  "monthlyLimit": 50000,
  "billingEnabled": true,
  "nextBillingDate": "2025-11-01"
}
```

---

#### PUT /internal/v1/billing/project/{projectId} ⚠️

**說明**: 更新專案計費設定

**優先級**: P1

**認證**: 需要 (Admin only)

---

#### GET /internal/v1/billing/analytics/overview 💡

**說明**: 計費分析概覽

**優先級**: P2

**認證**: 需要 (Admin only)

---

## 10. 內部 API - 檔案管理

### 10.1 檔案 CRUD

#### GET /internal/v1/files 💡

**說明**: 列表查詢檔案

**優先級**: P2

**認證**: 需要

---

#### POST /internal/v1/files/upload 💡

**說明**: 上傳單一檔案

**優先級**: P2

**認證**: 需要

**請求**: multipart/form-data

---

#### POST /internal/v1/files/upload/multiple 💡

**說明**: 上傳多個檔案

**優先級**: P2

**認證**: 需要

---

#### GET /internal/v1/files/{id} 💡

**說明**: 查詢檔案資訊

**優先級**: P2

**認證**: 需要

---

#### DELETE /internal/v1/files/{id} 💡

**說明**: 刪除檔案

**優先級**: P2

**認證**: 需要

---

#### GET /internal/v1/files/{id}/download 💡

**說明**: 下載檔案

**優先級**: P2

**認證**: 需要

---

#### POST /internal/v1/files/validate 💡

**說明**: 驗證檔案（上傳前檢查）

**優先級**: P2

**認證**: 需要

---

## 11. 內部 API - 監控管理

### 11.1 系統監控

#### GET /internal/v1/monitoring/health ⚠️

**說明**: 詳細健康檢查

**優先級**: P1

**認證**: 不需要（建議限制 IP）

**回應範例**:
```json
{
  "status": "healthy",
  "timestamp": "2025-10-23T10:00:00Z",
  "uptime": "72h30m15s",
  "version": "v1.0.0",
  "services": {
    "database": {
      "status": "healthy",
      "responseTime": 5,
      "lastCheck": "2025-10-23T09:59:55Z"
    },
    "redis": {
      "status": "healthy",
      "responseTime": 2,
      "lastCheck": "2025-10-23T09:59:55Z"
    },
    "teams_api": {
      "status": "healthy",
      "lastCheck": "2025-10-23T09:55:00Z"
    }
  }
}
```

---

#### GET /internal/v1/monitoring/performance ⚠️

**說明**: 性能指標

**優先級**: P1

**認證**: 不需要（建議限制 IP）

**回應範例**:
```json
{
  "timestamp": "2025-10-23T10:00:00Z",
  "cpu": {
    "usage": 35.2,
    "cores": 4
  },
  "memory": {
    "used": 268435456,
    "total": 1073741824,
    "usage": 25.0
  },
  "database": {
    "connections": 15,
    "queryTime": 5.2
  },
  "redis": {
    "memory": 52428800,
    "keys": 1500,
    "hitRate": 0.95
  }
}
```

---

#### GET /internal/v1/monitoring/business ⚠️

**說明**: 業務指標

**優先級**: P1

**認證**: 需要

**回應範例**:
```json
{
  "timestamp": "2025-10-23T10:00:00Z",
  "notifications": {
    "totalSent": 150000,
    "successRate": 0.995,
    "averageDeliveryTime": 8.5
  },
  "users": {
    "active": 500,
    "total": 800
  },
  "companies": {
    "active": 50,
    "total": 65
  }
}
```

---

#### GET /internal/v1/monitoring/alerts ⚠️

**說明**: 告警狀態

**優先級**: P1

**認證**: 需要

---

#### GET /internal/v1/monitoring/dashboard ⚠️

**說明**: 監控儀表板（綜合資料）

**優先級**: P1

**認證**: 需要

---

## 12. 內部 API - 佇列管理

### 12.1 佇列監控

#### GET /internal/v1/queue/stats ✅

**說明**: 佇列統計資訊

**優先級**: P0

**認證**: 不需要（建議限制 IP）

**回應範例**:
```json
{
  "timestamp": "2025-10-23T10:00:00Z",
  "totalItems": 150,
  "pendingItems": 50,
  "processingItems": 10,
  "completedItems": 85,
  "failedItems": 5,
  "averageProcessingTime": 8.5,
  "throughput": 120.5
}
```

---

#### GET /internal/v1/queue/status ✅

**說明**: 佇列狀態

**優先級**: P0

**認證**: 不需要（建議限制 IP）

**回應範例**:
```json
{
  "status": "running",
  "workers": 10,
  "activeWorkers": 8,
  "queueSize": 50,
  "lastProcessed": "2025-10-23T09:59:58Z"
}
```

---

#### POST /internal/v1/queue/clear 💡

**說明**: 清空佇列（危險操作）

**優先級**: P2

**認證**: 需要 (Admin only)

---

## 附錄

### A. HTTP 狀態碼說明

| 狀態碼 | 說明 | 使用情境 |
|-------|------|---------|
| 200 OK | 成功 | GET/PUT/PATCH 成功 |
| 201 Created | 已創建 | POST 創建成功 |
| 202 Accepted | 已接受 | 異步處理接受 |
| 204 No Content | 無內容 | DELETE 成功 |
| 400 Bad Request | 錯誤請求 | 參數驗證失敗 |
| 401 Unauthorized | 未認證 | Token 無效或過期 |
| 403 Forbidden | 無權限 | 權限不足 |
| 404 Not Found | 找不到資源 | 資源不存在 |
| 409 Conflict | 衝突 | 資料衝突（如重複 Key） |
| 429 Too Many Requests | 請求過多 | 超過速率限制 |
| 500 Internal Server Error | 伺服器錯誤 | 系統內部錯誤 |
| 503 Service Unavailable | 服務不可用 | 系統維護或過載 |

### B. 錯誤碼對照

| 錯誤碼 | HTTP 狀態 | 說明 |
|-------|----------|------|
| VALIDATION_FAILED | 400 | 輸入驗證失敗 |
| INVALID_INPUT | 400 | 無效的輸入 |
| UNAUTHORIZED | 401 | 未授權 |
| INVALID_TOKEN | 401 | 無效的 Token |
| FORBIDDEN | 403 | 權限不足 |
| NOT_FOUND | 404 | 資源不存在 |
| DUPLICATE_RECORD | 409 | 資料重複 |
| INTERNAL_ERROR | 500 | 內部錯誤 |
| DATABASE_ERROR | 500 | 資料庫錯誤 |
| TEAMS_API_ERROR | 500 | Teams API 錯誤 |

### C. 優先級定義

- **P0 (必須)**: 核心功能，系統無法運作缺少此功能
- **P1 (建議)**: 重要功能，建議實作以提升系統完整性
- **P2 (可選)**: 進階功能，可在後續版本實作

### D. API 版本說明

- **v1**: 初始版本
- 版本號位於 URL 路徑中：`/api/v1/...` 或 `/internal/v1/...`
- 當有不相容的變更時，會發布新版本（v2, v3, ...）
- 舊版本會維護至少 6 個月

---

**文件結束**

**總計端點數**:
- 系統端點: 2
- 外部 API: 8
- 內部 API: 90+

**P0 必須端點**: 約 60 個  
**P1 建議端點**: 約 25 個  
**P2 可選端點**: 約 15 個

如有任何問題，請聯絡：it-rfq@company.com

