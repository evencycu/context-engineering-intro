# Teams Notification API - 最小需求 RFQ
# Minimal Requirements RFQ

**RFQ 編號**: RFQ-2025-TEAMS-NOTIFY-MINIMAL-001  
**版本**: v1.0  
**日期**: 2025年10月22日  
**類型**: 最小可行產品 (MVP)

---

## 概述

本 RFQ 定義 Teams Notification API 的最小需求，專注於核心功能，確保系統能夠正常運作並滿足基本業務需求。

## 核心原則

- **最小化**: 僅包含絕對必要的功能
- **穩定性**: 確保核心功能穩定可靠
- **可擴展**: 為未來功能擴展預留空間
- **易用性**: 簡化外部系統整合流程

---

## 1. 系統端點 (P0 - 必須)

### 1.1 健康檢查

#### GET /health ✅

**說明**: 基本健康檢查端點

**優先級**: P0

**認證**: 不需要

**回應範例**:
```json
{
  "status": "healthy",
  "timestamp": "2025-10-22T10:00:00Z",
  "uptime": "72h30m15s",
  "version": "v1.0.0"
}
```

---

## 2. 外部 API 端點 (P0 - 必須)

### 2.1 通知發送

#### POST /api/v1/notify ✅

**說明**: 外部系統發送通知（核心功能）

**優先級**: P0

**認證**: 不需要（使用 notifyKey 認證）

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
    "notification_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "pending",
    "message": "Notification queued for processing",
    "project_id": "750e8400-e29b-41d4-a716-446655440001",
    "project_name": "my-project-key",
    "destinations_count": 3,
    "estimated_delivery": "5-10 minutes"
  }
}
```

**錯誤回應** (400/404):
```json
{
  "success": false,
  "error": "project not found for notify_key: invalid-key-12345"
}
```

---

### 2.2 目的地查詢

#### GET /api/v1/destinations/{notifyKey} ✅

**說明**: 查詢專案的所有目的地

**優先級**: P0

**認證**: 不需要

**回應範例**:
```json
{
  "success": true,
  "destinations": [
    {
      "id": "950e8400-e29b-41d4-a716-446655440001",
      "name": "Development Team",
      "description": "Development team notifications",
      "status": "active",
      "targets": [
        {
          "type": "channel",
          "display_name": "Development Channel"
        }
      ]
    },
    {
      "id": "950e8400-e29b-41d4-a716-446655440003",
      "name": "All Bot Installations",
      "description": "Golden sample - sends to all active bot installations",
      "status": "active",
      "targets": [
        {
          "type": "personal",
          "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
          "display_name": "Test User Personal",
          "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
        },
        {
          "type": "groupchat",
          "conversation_id": "19:f26a8d8a235f430db87a404491cd2ffc@thread.v2",
          "display_name": "Test Group Chat",
          "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
        },
        {
          "type": "channel",
          "conversation_id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
          "display_name": "Test Channel",
          "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
        }
      ]
    }
  ]
}
```

---

### 2.3 專案開通 (Provision)

#### POST /api/v1/provision ✅

**說明**: 一鍵建立專案和目的地（簡化流程）

**優先級**: P0

**認證**: 需要

**請求範例**:
```json
{
  "companyId": "550e8400-e29b-41d4-a716-446655440001",
  "createdBy": "650e8400-e29b-41d4-a716-446655440001",
  "project_name": "新專案",
  "project_description": "專案描述",
  "teamsTenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
  "targets": [
    {
      "type": "channel",
      "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
      "conversation_id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2"
    },
    {
      "type": "personal",
      "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
      "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR"
    }
  ]
}
```

**回應範例** (201 Created):
```json
{
  "success": true,
  "data": {
    "notifyKey": "5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38",
    "projectId": "750e8400-e29b-41d4-a716-446655440001",
    "projectName": "新專案",
    "status": "active",
    "destinationsCount": 1,
    "createdAt": "2025-10-22T10:00:00Z"
  }
}
```

---

#### GET /api/v1/provision/{notifyKey} ✅

**說明**: 查詢專案開通狀態

**優先級**: P0

**認證**: 需要

**回應範例**:
```json
{
  "success": true,
  "data": {
    "notifyKey": "5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38",
    "projectName": "新專案",
    "projectDescription": "專案描述",
    "status": "active",
    "destinationsCount": 1,
    "teamsTenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
    "createdAt": "2025-10-22T10:00:00Z"
  }
}
```

---

### 2.4 Bot 訊息接收

#### POST /api/v1/messages ✅

**說明**: 接收 Teams Bot Framework 訊息

**優先級**: P0

**認證**: Bot Framework Token

**請求範例** (Bot Framework Activity):
```json
{
  "type": "message",
  "id": "xxx",
  "timestamp": "2025-10-22T10:00:00Z",
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

## 3. 內部管理 API (P0 - 必須)

### 3.1 專案管理

#### POST /internal/v1/projects ✅

**說明**: 創建專案

**優先級**: P0

**認證**: 需要

**請求範例**:
```json
{
  "companyId": "550e8400-e29b-41d4-a716-446655440001",
  "notifyKey": "my-project-key",
  "description": "專案說明",
  "dailyLimit": 1000,
  "monthlyLimit": 30000,
  "priority": "normal",
  "createdBy": "650e8400-e29b-41d4-a716-446655440001"
}
```

**回應範例** (201 Created):
```json
{
  "id": "750e8400-e29b-41d4-a716-446655440001",
  "companyId": "550e8400-e29b-41d4-a716-446655440001",
  "notifyKey": "my-project-key",
  "description": "專案說明",
  "dailyLimit": 1000,
  "monthlyLimit": 30000,
  "priority": "normal",
  "status": "active",
  "createdBy": "650e8400-e29b-41d4-a716-446655440001",
  "createdAt": "2025-10-22T10:00:00Z"
}
```

---

#### GET /internal/v1/projects/key/{keyName} ✅

**說明**: 依 Notify Key 查詢專案

**優先級**: P0

**認證**: 需要

---

### 3.2 Bot 管理

#### POST /internal/v1/bots/platform ✅

**說明**: 創建 Teams Bot

**優先級**: P0

**認證**: 需要 (Admin only)

**請求範例**:
```json
{
  "name": "My Teams Bot",
  "appId": "844146d7-4ac9-4e4d-a463-d6e027714e81",
  "appPassword": "secret-password",
  "capabilities": ["proactive_messaging", "bot_messaging"],
  "status": "active",
  "description": "Bot 說明"
}
```

---

#### GET /internal/v1/bots/platform ✅

**說明**: 列表查詢 Bot

**優先級**: P0

**認證**: 需要

---

### 3.3 目的地管理

#### POST /internal/v1/destinations ✅

**說明**: 創建目的地

**優先級**: P0

**認證**: 需要

**請求範例**:
```json
{
  "projectId": "750e8400-e29b-41d4-a716-446655440001",
  "name": "Development Team",
  "description": "Development team notifications",
  "teamsTenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
  "targets": [
    {
      "type": "personal",
      "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
      "display_name": "Test User Personal",
      "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
    }
  ],
  "botId": "850e8400-e29b-41d4-a716-446655440001",
  "status": "active",
  "createdBy": "650e8400-e29b-41d4-a716-446655440001"
}
```

---

#### GET /internal/v1/destinations/project/{projectId} ✅

**說明**: 依專案查詢目的地

**優先級**: P0

**認證**: 需要

---

### 3.4 通知管理

#### GET /internal/v1/notifications/{id} ✅

**說明**: 查詢單一通知

**優先級**: P0

**認證**: 需要

**回應範例**:
```json
{
  "id": "uuid",
  "projectId": "750e8400-e29b-41d4-a716-446655440001",
  "senderId": "650e8400-e29b-41d4-a716-446655440001",
  "messageType": "text",
  "content": "通知內容",
  "mentions": [],
  "priority": "normal",
  "status": "sent",
  "destinationsSent": 3,
  "destinationsFailed": 0,
  "totalDestinations": 3,
  "createdAt": "2025-10-22T10:00:00Z",
  "sentAt": "2025-10-22T10:00:15Z"
}
```

---

### 3.5 佇列監控

#### GET /internal/v1/queue/stats ✅

**說明**: 佇列統計資訊

**優先級**: P0

**認證**: 不需要（建議限制 IP）

**回應範例**:
```json
{
  "timestamp": "2025-10-22T10:00:00Z",
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

## 4. 資料模型 (最小需求)

### 4.1 核心表格

- **companies**: 公司基本資料
- **users**: 用戶帳號
- **projects**: 專案設定（包含 notifyKey）
- **teams_bots**: Teams Bot 配置
- **destinations**: 通知目的地
- **notifications**: 通知記錄
- **notification_destinations**: 通知分派記錄
- **usage_records**: 使用量記錄（專案級計費）

### 4.2 關鍵欄位

- **projects.notify_key**: 唯一識別碼，用於外部 API 認證
- **destinations.targets**: JSONB 陣列，支援 personal/groupchat/channel
- **notification_destinations**: 每個 target 的發送狀態追蹤

---

## 5. 技術需求

### 5.1 基礎設施

- **語言**: Go 1.21+
- **框架**: Gin Web Framework
- **資料庫**: PostgreSQL 15+
- **快取/佇列**: Redis
- **部署**: Docker + Docker Compose

### 5.2 外部整合

- **Microsoft Teams Bot Framework**: 發送通知
- **Microsoft Graph API**: Token 管理
- **Teams API**: 速率限制合規（50 RPS global, 7 RPS per conversation）

### 5.3 安全需求

- **認證**: JWT Token（內部 API）
- **授權**: RBAC 角色（admin, user, viewer）
- **傳輸**: HTTPS 強制
- **資料**: 敏感資料加密存儲

---

## 6. 非功能性需求

### 6.1 性能

- **API 響應時間**: P95 < 2 秒
- **通知發送**: 平均 5-10 秒送達
- **並發處理**: 支援 100+ 並發請求
- **佇列處理**: 支援 1000+ 待處理項目

### 6.2 可靠性

- **可用性**: 99.5% 以上
- **重試機制**: 自動重試失敗通知
- **錯誤處理**: 完整的錯誤碼和訊息
- **監控**: 基本健康檢查和指標

### 6.3 可擴展性

- **水平擴展**: 支援多實例部署
- **資料庫**: 支援讀寫分離
- **佇列**: Redis 分散式佇列
- **快取**: Token 和配置快取

---

## 7. 測試需求

### 7.1 功能測試

- **外部 API**: 完整的功能測試腳本
- **內部 API**: CRUD 操作測試
- **整合測試**: Teams API 整合測試
- **錯誤處理**: 各種錯誤情境測試

### 7.2 性能測試

- **負載測試**: 模擬高並發場景
- **壓力測試**: 系統極限測試
- **穩定性測試**: 長時間運行測試

---

## 8. 交付物

### 8.1 代碼交付

- **源代碼**: 完整的 Go 源代碼
- **Docker 配置**: Dockerfile 和 docker-compose.yml
- **部署腳本**: 自動化部署腳本
- **測試腳本**: 完整的測試套件

### 8.2 文檔交付

- **API 文檔**: OpenAPI 3.0 規格
- **部署指南**: 詳細的部署說明
- **用戶手冊**: 外部 API 使用指南
- **運維手冊**: 監控和維護指南

### 8.3 環境交付

- **開發環境**: 本地開發環境
- **測試環境**: 功能測試環境
- **生產環境**: 生產部署環境

---

## 9. 時程規劃

### 9.1 開發階段

- **第 1 週**: 基礎架構和核心 API
- **第 2 週**: Bot 整合和通知發送
- **第 3 週**: 佇列系統和重試機制
- **第 4 週**: 測試和文檔

### 9.2 測試階段

- **第 5 週**: 功能測試和修復
- **第 6 週**: 性能測試和優化
- **第 7 週**: 整合測試和部署

### 9.3 交付階段

- **第 8 週**: 最終測試和文檔完善
- **第 9 週**: 生產部署和驗收

---

## 10. 驗收標準

### 10.1 功能驗收

- ✅ 所有 P0 端點正常運作
- ✅ 外部 API 完整功能測試通過
- ✅ Teams 通知成功發送
- ✅ 錯誤處理正確

### 10.2 性能驗收

- ✅ API 響應時間符合要求
- ✅ 通知發送時間符合要求
- ✅ 系統穩定性測試通過
- ✅ 負載測試通過

### 10.3 文檔驗收

- ✅ API 文檔完整準確
- ✅ 部署指南可執行
- ✅ 用戶手冊清晰易懂
- ✅ 運維手冊實用

---

## 附錄

### A. 最小需求端點清單

| 端點 | 方法 | 優先級 | 說明 |
|------|------|--------|------|
| `/health` | GET | P0 | 健康檢查 |
| `/api/v1/notify` | POST | P0 | 發送通知 |
| `/api/v1/destinations/{notifyKey}` | GET | P0 | 查詢目的地 |
| `/api/v1/provision` | POST | P0 | 一鍵建立專案 |
| `/api/v1/provision/{notifyKey}` | GET | P0 | 查詢專案開通狀態 |
| `/api/v1/messages` | POST | P0 | Bot 訊息接收 |
| `/internal/v1/projects` | POST | P0 | 創建專案 |
| `/internal/v1/projects/key/{keyName}` | GET | P0 | 查詢專案 |
| `/internal/v1/bots/platform` | POST/GET | P0 | Bot 管理 |
| `/internal/v1/destinations` | POST | P0 | 創建目的地 |
| `/internal/v1/destinations/project/{projectId}` | GET | P0 | 查詢目的地 |
| `/internal/v1/notifications/{id}` | GET | P0 | 查詢通知 |
| `/internal/v1/queue/stats` | GET | P0 | 佇列統計 |

**總計**: 13 個核心端點

### B. 資料庫表格 (最小需求)

- companies
- users  
- projects
- teams_bots
- destinations
- notifications
- notification_destinations
- usage_records

### C. 外部依賴

- Microsoft Teams Bot Framework
- Microsoft Graph API
- PostgreSQL
- Redis

---

**文件結束**

**最小需求總結**:
- **核心端點**: 13 個
- **資料表格**: 8 個
- **開發時程**: 9 週
- **預估工作量**: 2-3 人月

如有任何問題，請聯絡：it-rfq@company.com
