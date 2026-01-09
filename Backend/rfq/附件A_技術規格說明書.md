# 附件 A：技術規格說明書
# Technical Specification Document

**RFQ 編號**: RFQ-2025-TEAMS-NOTIFY-001  
**版本**: v1.0  
**日期**: 2025年10月23日

---

## 目錄

1. [系統概述](#1-系統概述)
2. [技術架構](#2-技術架構)
3. [功能模組規格](#3-功能模組規格)
4. [資料庫設計](#4-資料庫設計)
5. [API 規格](#5-api-規格)
6. [非功能需求](#6-非功能需求)
7. [整合規格](#7-整合規格)
8. [安全性規格](#8-安全性規格)
9. [部署規格](#9-部署規格)
10. [測試規格](#10-測試規格)

---

## 1. 系統概述

### 1.1 系統定位

Teams Notification API Server 是一個企業級的 Microsoft Teams 通知服務中台，提供統一的 RESTful API 供內部各應用系統調用，實現 Teams 訊息的集中管理與發送。

### 1.2 核心特性

- ✅ **預先註冊機制**: 採用 Notification Key (Project) 映射目的地，確保安全可控
- ✅ **多目標支援**: 支援 Channel、GroupChat、Personal 三種目的地類型
- ✅ **異步處理**: 採用 Actor Pattern + Redis Queue 實現高效異步處理
- ✅ **重試保護**: 指數退避 + 熔斷器 + Teams Retry-After 智能處理
- ✅ **多租戶架構**: 支援公司、專案、用戶三級權限管理
- ✅ **完整計費**: 專案級使用量統計與計費報表

### 1.3 技術亮點

| 技術特性 | 實作方式 | 效益 |
|---------|---------|------|
| **Two-loop Enqueue** | Producer 快速寫 DB + Background Worker 批次入隊 | 降低 API 延遲 |
| **Actor Pool** | Goroutine Pool 管理並發 Actor | 可控並發、防止資源耗盡 |
| **Token Cache** | Redis 快取 Teams Token（50min TTL） | 減少 API 調用、提升效能 |
| **Circuit Breaker** | 熔斷器保護外部 API | 快速失敗、防止級聯故障 |
| **Rate Limiting** | 符合 Teams API 限制（50/7 RPS） | 避免被 Microsoft 限流 |

---

## 2. 技術架構

### 2.1 技術棧

#### 後端技術

| 類別 | 技術選型 | 版本要求 | 用途 |
|------|---------|---------|------|
| **程式語言** | Golang | 1.21+ | 後端開發 |
| **Web 框架** | Gin | 1.9+ | HTTP Server |
| **ORM** | GORM 或 sqlx | - | 資料庫存取 |
| **主要資料庫** | PostgreSQL | 14+ | 永久儲存 |
| **快取/佇列** | Redis | 7+ | 快取與訊息佇列 |
| **JWT** | golang-jwt/jwt | 5.0+ | 認證 Token |
| **日誌** | logrus | 1.9+ | 結構化日誌 |
| **HTTP Client** | net/http | stdlib | API 調用 |

#### 開發工具

- **版本控制**: Git
- **程式碼檢查**: golint, gofmt, go vet
- **測試框架**: testing (stdlib), testify
- **文件產生**: godoc, swagger
- **CI/CD**: 需支援 GitHub Actions 或 GitLab CI

#### 部署技術

- **容器化**: Docker 20+
- **容器編排**: Kubernetes 1.24+ 或 Docker Compose
- **配置管理**: YAML 配置檔 + 環境變數
- **監控**: 內建健康檢查端點

### 2.2 系統架構圖

```
┌─────────────────────────────────────────────────────────────────┐
│                        外部系統 (Clients)                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐         │
│  │ 應用 A    │  │ 應用 B    │  │ 應用 C    │  │ 監控系統  │         │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘         │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTPS/REST API
┌────────────────────────┼────────────────────────────────────────┐
│                   API Gateway Layer                              │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │ Gin HTTP Server                                           │   │
│  │  - CORS Middleware                                        │   │
│  │  - JWT Auth Middleware                                    │   │
│  │  - Request ID Middleware                                  │   │
│  │  - Logging Middleware                                     │   │
│  │  - Rate Limiting Middleware                               │   │
│  └──────────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────┼────────────────────────────────────────┐
│                 Application Layer                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  Handlers   │  │  Services   │  │ Repositories │             │
│  │  (HTTP)     │──│  (Business) │──│  (Data)     │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│                           │                                       │
│  ┌─────────────────────────────────────────────────────┐         │
│  │         Actor System (Async Processing)              │         │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐          │         │
│  │  │ Actor 1  │  │ Actor 2  │  │ Actor N  │          │         │
│  │  └──────────┘  └──────────┘  └──────────┘          │         │
│  │         ActorPool (Goroutine Pool)                   │         │
│  └─────────────────────────────────────────────────────┘         │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────┼────────────────────────────────────────┐
│                    Data Layer                                    │
│  ┌──────────────────┐              ┌──────────────────┐         │
│  │   PostgreSQL     │              │      Redis       │         │
│  │  - Companies     │              │  - Token Cache   │         │
│  │  - Users         │              │  - Queue         │         │
│  │  - Projects      │              │  - Circuit State │         │
│  │  - Notifications │              │  - Locks         │         │
│  │  - Destinations  │              └──────────────────┘         │
│  │  - Billing       │                                            │
│  └──────────────────┘                                            │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────┼────────────────────────────────────────┐
│                Integration Layer                                 │
│  ┌──────────────────┐  ┌──────────────────┐                     │
│  │ Teams Bot API    │  │  Graph API       │                     │
│  │ (Bot Framework)  │  │  (Optional)      │                     │
│  └──────────────────┘  └──────────────────┘                     │
└──────────────────────────────────────────────────────────────────┘
```

### 2.3 資料流程

#### 通知發送流程（Two-loop Design）

```
1. Producer Loop (API Handler):
   ┌─────────────────────────────────────────────────────┐
   │ Client Request → API Validation → Create Record    │
   │ → Insert notification (status=pending)              │
   │ → Insert notification_destinations (status=pending) │
   │ → Return 202 Accepted (Async)                       │
   └─────────────────────────────────────────────────────┘
                         ↓
2. Enqueue Worker (Background):
   ┌─────────────────────────────────────────────────────┐
   │ Scan pending records → Sort by priority             │
   │ → Push to Redis Queue (high/normal/low)             │
   │ → Update status=enqueued                            │
   └─────────────────────────────────────────────────────┘
                         ↓
3. Consumer Loop (QueueConsumer + ActorPool):
   ┌─────────────────────────────────────────────────────┐
   │ BLPOP from Redis → SETNX Lock (dedup)               │
   │ → Spawn Actor → Get Token from Cache                │
   │ → Call Teams API → Update status (sent/failed)      │
   │ → Retry on failure (exponential backoff)            │
   │ → Circuit Breaker protection                        │
   └─────────────────────────────────────────────────────┘
```

---

## 3. 功能模組規格

### 3.1 模組清單

| 模組代碼 | 模組名稱 | 優先級 | 核心功能 |
|---------|---------|-------|---------|
| M-001 | 公司管理 | P0 | 公司 CRUD、狀態管理 |
| M-002 | 用戶管理 | P0 | 用戶 CRUD、權限管理、密碼變更 |
| M-003 | 專案管理 | P0 | 專案 CRUD、Notify Key 管理、配額設定 |
| M-004 | Bot 管理 | P0 | Teams Bot CRUD、狀態追蹤 |
| M-005 | 目的地管理 | P0 | 目的地 CRUD、Targets 配置 (JSONB) |
| M-006 | 通知管理 | P0 | 發送通知、狀態查詢、取消通知 |
| M-007 | 廣播服務 | P0 | 多目標廣播、Targets 解析 |
| M-008 | 佇列系統 | P0 | Redis Queue、Actor Pool、重試機制 |
| M-009 | Token 管理 | P1 | Token 快取、自動刷新 |
| M-010 | 計費系統 | P1 | 使用量記錄、計費報表 |
| M-011 | 外部 API | P1 | Provision API、Notify API |
| M-012 | 監控系統 | P1 | 健康檢查、系統指標、告警 |
| M-013 | 檔案管理 | P2 | 檔案上傳、下載 (未來擴充) |

### 3.2 模組詳細規格

#### M-001: 公司管理模組

**功能描述**: 管理多租戶公司資料

**主要 API**:
- `POST /internal/v1/companies` - 創建公司
- `GET /internal/v1/companies` - 列表查詢（支援分頁）
- `GET /internal/v1/companies/{id}` - 查詢單一公司
- `PUT /internal/v1/companies/{id}` - 更新公司資訊
- `DELETE /internal/v1/companies/{id}` - 刪除公司
- `PATCH /internal/v1/companies/{id}/status` - 更新狀態

**資料欄位**:
```go
type Company struct {
    ID            uuid.UUID  `json:"id"`
    Name          string     `json:"name"`           // 公司名稱
    ContactEmail  string     `json:"contact_email"`  // 聯絡信箱
    ContactPhone  string     `json:"contact_phone"`  // 聯絡電話
    Address       string     `json:"address"`        // 地址
    Status        string     `json:"status"`         // active/inactive/suspended
    BillingEnabled bool      `json:"billing_enabled"` // 是否啟用計費
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
}
```

**驗證規則**:
- `name`: 必填，2-255 字元
- `contact_email`: 必填，有效 email 格式
- `status`: 必填，限定值 active/inactive/suspended

---

#### M-002: 用戶管理模組

**功能描述**: 管理系統用戶與權限

**主要 API**:
- `POST /internal/v1/users` - 創建用戶
- `GET /internal/v1/users` - 列表查詢
- `GET /internal/v1/users/{id}` - 查詢單一用戶
- `PUT /internal/v1/users/{id}` - 更新用戶資訊
- `DELETE /internal/v1/users/{id}` - 刪除用戶
- `PATCH /internal/v1/users/{id}/password` - 變更密碼
- `GET /internal/v1/users/company/{companyId}` - 依公司查詢
- `GET /internal/v1/users/role/{role}` - 依角色查詢

**資料欄位**:
```go
type User struct {
    ID          uuid.UUID  `json:"id"`
    CompanyID   uuid.UUID  `json:"company_id"`
    Email       string     `json:"email"`         // 登入帳號
    Name        string     `json:"name"`          // 顯示名稱
    PasswordHash string    `json:"-"`             // 密碼雜湊（不回傳）
    Role        string     `json:"role"`          // admin/manager/user
    Status      string     `json:"status"`        // active/inactive
    LastLoginAt *time.Time `json:"last_login_at"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
```

**角色權限**:
- `admin`: 完整系統管理權限
- `manager`: 公司內資源管理權限
- `user`: 基本使用權限

**密碼要求**:
- 長度 ≥ 8 字元
- 使用 bcrypt 加密存儲
- 支援密碼變更功能

---

#### M-003: 專案管理模組

**功能描述**: 管理通知專案（Notification Key）

**主要 API**:
- `POST /internal/v1/projects` - 創建專案
- `GET /internal/v1/projects` - 列表查詢
- `GET /internal/v1/projects/{id}` - 查詢單一專案
- `PUT /internal/v1/projects/{id}` - 更新專案
- `DELETE /internal/v1/projects/{id}` - 刪除專案
- `PATCH /internal/v1/projects/{id}/limits` - 更新配額
- `GET /internal/v1/projects/company/{companyId}` - 依公司查詢
- `GET /internal/v1/projects/key/{keyName}` - 依 Key 查詢

**資料欄位**:
```go
type Project struct {
    ID           uuid.UUID  `json:"id"`
    CompanyID    uuid.UUID  `json:"company_id"`
    NotifyKey    string     `json:"notify_key"`     // 唯一識別碼
    Description  string     `json:"description"`    // 專案說明
    DailyLimit   int        `json:"daily_limit"`    // 每日配額
    MonthlyLimit int        `json:"monthly_limit"`  // 每月配額
    Priority     string     `json:"priority"`       // low/normal/high
    Status       string     `json:"status"`         // active/inactive
    CreatedBy    uuid.UUID  `json:"created_by"`     // 建立者
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}
```

**唯一約束**:
- `notify_key` 必須唯一

---

#### M-004: Teams Bot 管理模組

**功能描述**: 管理 Microsoft Teams Bot 資訊

**主要 API**:
- `POST /internal/v1/bots/platform` - 創建 Bot
- `GET /internal/v1/bots/platform` - 列表查詢
- `GET /internal/v1/bots/platform/{id}` - 查詢單一 Bot
- `PUT /internal/v1/bots/platform/{id}` - 更新 Bot
- `DELETE /internal/v1/bots/platform/{id}` - 刪除 Bot
- `PATCH /internal/v1/bots/platform/{id}/status` - 更新狀態
- `POST /internal/v1/bots/platform/{id}/test` - 測試連線

**資料欄位**:
```go
type TeamsBot struct {
    ID           uuid.UUID    `json:"id"`
    Name         string       `json:"name"`          // Bot 名稱
    AppID        string       `json:"app_id"`        // Microsoft App ID
    AppPassword  string       `json:"-"`             // App Secret（加密）
    Capabilities []string     `json:"capabilities"`  // Bot 能力
    Status       string       `json:"status"`        // active/inactive/error
    Description  string       `json:"description"`   // 說明
    CreatedAt    time.Time    `json:"created_at"`
    UpdatedAt    time.Time    `json:"updated_at"`
}
```

**安全要求**:
- `app_password` 必須加密存儲
- API 回應不可回傳明文密碼

---

#### M-005: 目的地管理模組

**功能描述**: 管理通知目的地（Channel/GroupChat/Personal）

**主要 API**:
- `POST /internal/v1/destinations` - 創建目的地
- `GET /internal/v1/destinations` - 列表查詢
- `GET /internal/v1/destinations/{id}` - 查詢單一目的地
- `PUT /internal/v1/destinations/{id}` - 更新目的地
- `DELETE /internal/v1/destinations/{id}` - 刪除目的地
- `PATCH /internal/v1/destinations/{id}/targets` - 更新 Targets
- `POST /internal/v1/destinations/{id}/validate` - 驗證 Targets
- `GET /internal/v1/destinations/project/{projectId}` - 依專案查詢
- `GET /internal/v1/destinations/bot/{botId}` - 依 Bot 查詢

**資料欄位**:
```go
type Destination struct {
    ID          uuid.UUID    `json:"id"`
    ProjectID   uuid.UUID    `json:"project_id"`
    BotID       uuid.UUID    `json:"bot_id"`
    Type        string       `json:"type"`         // personal/groupchat/channel
    Targets     []Target     `json:"targets"`      // JSONB 陣列
    Status      string       `json:"status"`       // active/inactive
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
}

type Target struct {
    Type            string `json:"type"`              // personal/groupchat/channel
    ConversationID  string `json:"conversation_id"`   // Teams 對話 ID
    Email           string `json:"email,omitempty"`   // 僅 personal 類型
    DisplayName     string `json:"display_name"`      // 顯示名稱
    TenantID        string `json:"tenant_id"`         // Tenant ID
}
```

**Targets 欄位約束**:
- `type`: 必填，限定值 personal/groupchat/channel
- `tenant_id`: 必填
- `conversation_id`: groupchat/channel 必填
- `email`: personal 類型可填（與 conversation_id 二擇一）

---

#### M-006: 通知管理模組

**功能描述**: 管理通知發送與狀態追蹤

**主要 API**:
- `POST /internal/v1/notifications` - 發送通知
- `GET /internal/v1/notifications` - 列表查詢
- `GET /internal/v1/notifications/{id}` - 查詢單一通知
- `DELETE /internal/v1/notifications/{id}` - 取消通知
- `POST /internal/v1/notifications/{id}/retry` - 手動重試
- `GET /internal/v1/notifications/project/{projectId}` - 依專案查詢
- `GET /internal/v1/notifications/sender/{senderId}` - 依發送者查詢
- `GET /internal/v1/notifications/status/{status}` - 依狀態查詢
- `GET /internal/v1/notifications/date-range` - 依時間範圍查詢

**資料欄位**:
```go
type Notification struct {
    ID                 uuid.UUID   `json:"id"`
    ProjectID          uuid.UUID   `json:"project_id"`
    SenderID           *uuid.UUID  `json:"sender_id,omitempty"`
    MessageType        string      `json:"message_type"`      // text/file/adaptive_card
    Content            string      `json:"content"`           // 訊息內容
    Mentions           []string    `json:"mentions"`          // @提及
    Priority           string      `json:"priority"`          // low/normal/high
    Status             string      `json:"status"`            // pending/enqueued/sending/sent/failed
    ErrorMessage       string      `json:"error_message"`     // 錯誤訊息
    DestinationsSent   int         `json:"destinations_sent"` // 成功數
    DestinationsFailed int         `json:"destinations_failed"` // 失敗數
    TotalDestinations  int         `json:"total_destinations"` // 總數
    CreatedAt          time.Time   `json:"created_at"`
    UpdatedAt          time.Time   `json:"updated_at"`
    SentAt             *time.Time  `json:"sent_at"`           // 發送完成時間
}
```

**狀態流轉**:
```
pending → enqueued → sending → sent (成功)
                            → failed (失敗)
```

---

#### M-007: 廣播服務模組

**功能描述**: 支援一次發送到多個目的地

**主要功能**:
1. 接收廣播請求（targets 參數）
2. 解析 targets（支援 "all"、conversation_id、email）
3. 展開為多個 notification_destinations 記錄
4. 批次處理發送

**Targets 解析邏輯**:
```go
// targets=["all"] → 發送到專案所有目的地
// targets=["conversation_id_1", "conversation_id_2"] → 發送到指定對話
// targets=["user@example.com"] → 發送到指定 email（personal 類型）
// targets=["all", "exclude:conversation_id_1"] → 排除特定目的地
```

---

#### M-008: 佇列系統模組

**功能描述**: 異步處理通知發送

**核心組件**:
1. **EnqueueWorker**: 背景掃描 pending 記錄，推入 Redis
2. **QueueConsumer**: 從 Redis 取出任務
3. **ActorPool**: 管理並發 Actor
4. **NotificationActor**: 實際執行發送邏輯

**Redis 佇列設計**:
```
Queue Keys:
- queue:notifications:high    (高優先級)
- queue:notifications:normal  (一般優先級)
- queue:notifications:low     (低優先級)

Lock Keys:
- queue:notifications:processing:{id}  (處理鎖，TTL 30min)

Circuit Breaker Keys:
- circuit:teams_api:state      (熔斷器狀態)
- circuit:teams_api:failures   (失敗計數)
```

**重試策略**:
```go
type RetryConfig struct {
    MaxRetries     int           = 5       // 最大重試次數
    InitialBackoff time.Duration = 1s     // 初始等待
    MaxBackoff     time.Duration = 5min   // 最長等待
    Multiplier     float64       = 2.0    // 退避倍數
    JitterFraction float64       = 0.1    // 隨機抖動 10%
}

// 重試時間計算
backoff = min(InitialBackoff * (Multiplier ^ retry_count), MaxBackoff)
actual_backoff = backoff * (1 + random(-JitterFraction, JitterFraction))
```

**熔斷器規則**:
- 連續失敗 5 次 → 進入 Open 狀態（不調用 API）
- Open 狀態持續 30 秒 → 進入 Half-Open 狀態（允許 1 次嘗試）
- Half-Open 成功 → 回到 Closed 狀態
- Half-Open 失敗 → 回到 Open 狀態

---

#### M-009: Token 管理模組

**功能描述**: 管理 Microsoft Teams Access Token

**核心功能**:
1. Token 快取（Redis）
2. 自動刷新（TTL 前 10 分鐘）
3. 多 Bot 支援

**Token 快取設計**:
```
Redis Key: token:teams:{bot_id}
TTL: 3000 seconds (50 minutes)
Value: {
    "access_token": "eyJ0...",
    "expires_at": "2025-10-23T12:00:00Z"
}
```

**Token 取得流程**:
```
1. Check Redis Cache
   ├─ Hit → Return cached token
   └─ Miss → Request new token from Microsoft
              └─ Cache in Redis (TTL 50min)
              └─ Return token
```

---

#### M-010: 計費系統模組

**功能描述**: 記錄使用量與產生計費報表

**主要 API**:
- `GET /internal/v1/billing/usage` - 使用量查詢
- `GET /internal/v1/billing/usage/summary` - 使用量彙總
- `GET /internal/v1/billing/usage/project/{projectId}` - 專案使用量
- `GET /internal/v1/billing/project/{projectId}` - 專案計費資訊
- `PUT /internal/v1/billing/project/{projectId}` - 更新計費設定

**資料欄位**:
```go
type UsageRecord struct {
    ID                uuid.UUID  `json:"id"`
    ProjectID         uuid.UUID  `json:"project_id"`
    NotificationCount int        `json:"notification_count"` // 通知數量
    Date              time.Time  `json:"date"`              // 統計日期
    CreatedAt         time.Time  `json:"created_at"`
}

type ProjectBilling struct {
    ProjectID      uuid.UUID  `json:"project_id"`      // 唯一鍵
    BillingEnabled bool       `json:"billing_enabled"` // 是否啟用計費
    PlanID         *uuid.UUID `json:"plan_id"`         // 計費方案 ID
    CurrentUsage   int        `json:"current_usage"`   // 當月使用量
    MonthlyLimit   int        `json:"monthly_limit"`   // 每月配額
    // ... 其他欄位
}
```

**使用量統計邏輯**:
- 每次成功發送通知 → `notification_count++`
- 每日彙總到 `usage_records`
- 提供 `usage_summary` View 供查詢

---

#### M-011: 外部 API 模組

**功能描述**: 提供簡化的外部調用 API

**主要 API**:

1. **Provision API**: 一鍵建立專案
   ```
   POST /api/v1/provision
   Body: {
       "notifyKey": "my-project",
       "companyId": "uuid",
       "projectName": "My Project",
       "description": "Project description"
   }
   ```

2. **Notify API**: 外部系統發送通知
   ```
   POST /api/v1/notify
   Body: {
       "notifyKey": "my-project",
       "message": "Hello Teams",
       "messageType": "text",
       "priority": "normal",
       "targets": ["all"]  // or ["conversation_id"] or ["email"]
   }
   ```

3. **Destinations API**: 查詢專案目的地
   ```
   GET /api/v1/destinations/{notifyKey}
   Response: {
       "notifyKey": "my-project",
       "destinations": [...]
   }
   ```

---

#### M-012: 監控系統模組

**功能描述**: 系統健康檢查與指標監控

**主要 API**:
- `GET /health` - 基本健康檢查
- `GET /api/v1/metrics` - 系統指標
- `GET /internal/v1/monitoring/health` - 詳細健康狀態
- `GET /internal/v1/monitoring/performance` - 性能指標
- `GET /internal/v1/monitoring/business` - 業務指標
- `GET /internal/v1/queue/stats` - 佇列統計
- `GET /internal/v1/queue/status` - 佇列狀態

**健康檢查項目**:
```json
{
  "status": "healthy",
  "timestamp": "2025-10-23T10:00:00Z",
  "uptime": "72h30m",
  "version": "v1.0.0",
  "checks": {
    "database": {
      "status": "healthy",
      "responseTime": "5ms"
    },
    "redis": {
      "status": "healthy",
      "responseTime": "2ms"
    },
    "teams_api": {
      "status": "healthy",
      "lastCheck": "2025-10-23T09:55:00Z"
    }
  }
}
```

**系統指標**:
- CPU 使用率
- 記憶體使用率
- Goroutine 數量
- Database 連接數
- Redis 連接數
- API 請求數/秒
- API 響應時間 (P50/P95/P99)
- 錯誤率 (4xx/5xx)
- 佇列長度
- Token 快取命中率

---

## 4. 資料庫設計

### 4.1 資料庫配置

- **資料庫**: PostgreSQL 14+
- **字符集**: UTF-8
- **時區**: UTC
- **連接池**: Min 5, Max 25
- **SSL**: 生產環境啟用

### 4.2 主要資料表

詳細的 ERD 圖請參閱**附件 C**。

#### 核心資料表清單

| 表名 | 說明 | 預估資料量 |
|------|------|-----------|
| `companies` | 公司資料 | < 1000 |
| `users` | 用戶帳號 | < 10000 |
| `projects` | 通知專案 | < 5000 |
| `teams_bots` | Teams Bot 資訊 | < 100 |
| `destinations` | 目的地配置 | < 50000 |
| `notifications` | 通知記錄 | > 1M (需分區) |
| `notification_destinations` | 通知分派記錄 | > 10M (需分區) |
| `usage_records` | 使用量記錄 | > 100K |
| `project_billing` | 專案計費 | < 5000 |
| `billing_plans` | 計費方案 | < 50 |
| `files` | 檔案資訊 | < 10000 |
| `audit_logs` | 審計日誌 | > 1M (需分區) |

#### 索引策略

**高頻查詢索引**:
```sql
-- Projects
CREATE INDEX idx_projects_notify_key ON projects(notify_key);
CREATE INDEX idx_projects_company_id ON projects(company_id);

-- Destinations
CREATE INDEX idx_destinations_project_id ON destinations(project_id);
CREATE INDEX idx_destinations_bot_id ON destinations(bot_id);

-- Notifications
CREATE INDEX idx_notifications_project_id ON notifications(project_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);

-- Notification Destinations
CREATE INDEX idx_notif_dest_notification_id ON notification_destinations(notification_id);
CREATE INDEX idx_notif_dest_status ON notification_destinations(status);
CREATE INDEX idx_notif_dest_retry ON notification_destinations(retry_count, next_retry_at);
```

#### 分區策略

**大資料表分區**（建議）:
```sql
-- notifications 按月分區
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    ...
) PARTITION BY RANGE (created_at);

CREATE TABLE notifications_2025_10 PARTITION OF notifications
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

-- 每月自動建立新分區
```

### 4.3 JSONB 欄位設計

#### destinations.targets

```json
{
  "targets": [
    {
      "type": "channel",
      "conversation_id": "19:xxx@thread.tacv2",
      "display_name": "General",
      "tenant_id": "tenant-uuid"
    },
    {
      "type": "personal",
      "email": "user@example.com",
      "display_name": "John Doe",
      "tenant_id": "tenant-uuid"
    }
  ]
}
```

**JSONB 查詢範例**:
```sql
-- 查詢包含特定 email 的目的地
SELECT * FROM destinations 
WHERE targets @> '[{"email": "user@example.com"}]';

-- 查詢特定類型的目的地
SELECT * FROM destinations 
WHERE targets @> '[{"type": "channel"}]';
```

---

## 5. API 規格

### 5.1 API 端點總覽

完整的 API 端點清單請參閱**附件 B**。

### 5.2 API 設計原則

1. **RESTful 風格**: 使用標準 HTTP 動詞 (GET/POST/PUT/DELETE/PATCH)
2. **統一路徑格式**: 
   - Internal API: `/internal/v1/{resource}`
   - External API: `/api/v1/{resource}`
3. **統一錯誤格式**: 
   ```json
   {
     "error": "ERROR_CODE",
     "message": "Human readable message",
     "details": "Additional information",
     "timestamp": "2025-10-23T10:00:00Z"
   }
   ```
4. **分頁標準**: 
   - Query: `?limit=10&offset=0`
   - Response: `{"data": [...], "pagination": {"total": 100, "limit": 10, "offset": 0}}`

### 5.3 認證機制

#### JWT Token 格式

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user-id",
    "company_id": "company-id",
    "role": "admin",
    "exp": 1729756800,
    "iat": 1729670400
  }
}
```

#### API Key 格式

- Header: `X-API-Key: {api_key}`
- 用於服務間認證

### 5.4 OpenAPI 規範

完整的 OpenAPI 3.0 規範已提供在專案中：
- 檔案路徑: `/api/openapi/teams-notification-api.yaml`
- 包含所有端點、Schema、範例

---

## 6. 非功能需求

### 6.1 性能需求

| 指標 | 目標值 | 測量方法 |
|------|-------|---------|
| API 響應時間 (P50) | < 200ms | Load testing |
| API 響應時間 (P95) | < 500ms | Load testing |
| API 響應時間 (P99) | < 1000ms | Load testing |
| 通知發送成功率 | ≥ 99% | Production monitoring |
| 併發請求支援 | ≥ 100 req/s | Load testing |
| 佇列處理能力 | ≥ 1000 notifications/min | Performance testing |
| Token 快取命中率 | ≥ 95% | Redis monitoring |

### 6.2 可用性需求

| 指標 | 目標值 |
|------|-------|
| 系統可用性 (Uptime) | ≥ 99.5% |
| 平均故障恢復時間 (MTTR) | < 30 minutes |
| 資料備份頻率 | 每日 |
| 資料保留期限 | ≥ 90 天 |

### 6.3 擴展性需求

- 支援水平擴展（多 Pod 部署）
- 無狀態設計（Session 存 Redis）
- Database Read Replica 支援
- Redis Cluster 支援

### 6.4 安全性需求

| 需求項目 | 實作方式 |
|---------|---------|
| 傳輸加密 | HTTPS/TLS 1.2+ |
| 資料加密 | 敏感欄位加密 (AES-256) |
| 認證機制 | JWT + API Key |
| 密碼強度 | ≥ 8 字元，bcrypt 雜湊 |
| SQL 防注入 | 使用 Prepared Statements |
| XSS 防護 | 輸入驗證與輸出編碼 |
| Rate Limiting | 每 IP 100 req/min |
| 審計日誌 | 所有關鍵操作記錄 |

---

## 7. 整合規格

### 7.1 Microsoft Teams Bot Framework

**整合端點**:
- Token Endpoint: `https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token`
- Bot Service: `https://smba.trafficmanager.net/apis/`

**發送訊息流程**:
```
1. Get Access Token (cached in Redis)
   POST https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token
   Body: client_id={app_id}&client_secret={secret}&grant_type=client_credentials

2. Send Proactive Message
   POST https://smba.trafficmanager.net/apis/v3/conversations/{conversationId}/activities
   Headers: Authorization: Bearer {token}
   Body: {
       "type": "message",
       "text": "Hello",
       "from": {"id": "{bot_id}"}
   }
```

**錯誤處理**:
- 401 Unauthorized → Token 失效，重新取得
- 429 Too Many Requests → 讀取 Retry-After header，延遲重試
- 5xx Server Error → 指數退避重試

### 7.2 Microsoft Graph API (選用)

如需使用 Graph API 取得用戶資訊或群組資訊：

**端點**: `https://graph.microsoft.com/v1.0/`

**常用 API**:
- 取得用戶: `GET /users/{id}`
- 取得團隊: `GET /teams/{id}`
- 取得頻道: `GET /teams/{id}/channels`

---

## 8. 安全性規格

### 8.1 認證流程

#### JWT 認證流程

```
1. Login
   POST /api/v1/auth/login
   Body: {"email": "user@example.com", "password": "xxx"}
   Response: {"token": "eyJ0...", "expires_at": "..."}

2. Use Token
   GET /internal/v1/projects
   Headers: Authorization: Bearer eyJ0...

3. Token Validation
   - Verify signature
   - Check expiration
   - Extract user_id & role
```

### 8.2 授權設計

#### RBAC 角色權限表

| 資源 | Admin | Manager | User |
|------|-------|---------|------|
| 公司管理 | CRUD | R | - |
| 用戶管理 | CRUD | CRUD (同公司) | R (self) |
| 專案管理 | CRUD | CRUD (同公司) | R |
| 通知發送 | ✓ | ✓ | ✓ |
| 計費查詢 | ✓ | ✓ (同公司) | - |
| 系統設定 | ✓ | - | - |

### 8.3 資料加密

**加密欄位**:
- `users.password_hash` - bcrypt
- `teams_bots.app_password` - AES-256-GCM
- API Keys - SHA-256 hash

**加密實作**:
```go
// 密碼加密
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// API Key 加密
block, _ := aes.NewCipher(key)
gcm, _ := cipher.NewGCM(block)
ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
```

---

## 9. 部署規格

### 9.1 Docker 映像要求

**Dockerfile 規範**:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

**映像大小**: < 50 MB

### 9.2 環境變數配置

**必要環境變數**:
```bash
# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USERNAME=teamsnotify
DATABASE_PASSWORD=secret
DATABASE_DATABASE_NAME=notification_center

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=secret

# Teams Bot
TEAMS_BOT_APP_ID=xxx
TEAMS_BOT_APP_PASSWORD=xxx

# JWT
JWT_SECRET=your-secret-key

# Server
SERVER_PORT=8080
SERVER_ENVIRONMENT=production
LOGGING_LEVEL=info
```

### 9.3 Kubernetes 部署

**最低資源需求**:
```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

**建議副本數**: 2-3 個 Pod

**Health Check**:
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

---

## 10. 測試規格

### 10.1 單元測試

**覆蓋率要求**: ≥ 70%

**測試範圍**:
- 所有 Service 層業務邏輯
- 所有 Repository 層資料存取
- 關鍵的 Handler 邏輯
- Utility 函式

**測試工具**:
- `testing` (標準庫)
- `testify/assert`
- `testify/mock`

**範例**:
```go
func TestProjectService_CreateProject(t *testing.T) {
    // Setup
    mockRepo := new(MockProjectRepository)
    service := NewProjectService(mockRepo)
    
    // Test data
    input := &CreateProjectInput{
        CompanyID: uuid.New(),
        NotifyKey: "test-project",
    }
    
    // Mock
    mockRepo.On("Create", mock.Anything).Return(nil)
    
    // Execute
    result, err := service.CreateProject(input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

### 10.2 整合測試

**測試範圍**:
- API 端對端測試
- 資料庫整合測試
- Redis 整合測試
- Teams API Mock 測試

**測試工具**:
- Shell scripts (cURL)
- Go integration tests
- Docker Compose (測試環境)

**範例**:
```bash
#!/bin/bash
# Test: Create and send notification

# 1. Create project
response=$(curl -X POST http://localhost:8080/internal/v1/projects \
  -H "Content-Type: application/json" \
  -d '{"company_id":"xxx","notify_key":"test"}')

# 2. Send notification
curl -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{"notifyKey":"test","message":"Hello"}'
  
# 3. Verify status
curl http://localhost:8080/internal/v1/notifications?limit=1
```

### 10.3 性能測試

**測試工具**: 
- Apache Bench
- wrk
- Go benchmark

**測試場景**:
1. API 響應時間測試
2. 併發請求測試
3. 佇列處理能力測試
4. 資料庫連接池測試

**範例**:
```bash
# 併發測試
ab -n 10000 -c 100 -H "Authorization: Bearer xxx" \
   http://localhost:8080/internal/v1/projects
```

### 10.4 驗收測試清單

**功能測試** (100% 通過):
- [ ] 公司 CRUD 操作
- [ ] 用戶 CRUD 操作與權限驗證
- [ ] 專案 CRUD 操作與 Notify Key 唯一性
- [ ] Bot 管理與狀態更新
- [ ] 目的地管理與 Targets 驗證
- [ ] 通知發送（文字訊息）
- [ ] 通知發送（Adaptive Card）
- [ ] 廣播功能（targets=all）
- [ ] 外部 API (/api/v1/notify)
- [ ] 使用量統計與計費查詢

**非功能測試**:
- [ ] API 響應時間 P95 < 500ms
- [ ] 通知發送成功率 ≥ 99%
- [ ] 併發 100 req/s 正常運行
- [ ] 重試機制正確運作
- [ ] 熔斷器保護機制有效
- [ ] 健康檢查端點回應正常
- [ ] 日誌格式符合規範

---

## 附錄

### A. 技術術語對照

| 英文 | 中文 | 說明 |
|------|------|------|
| Notification Key | 通知金鑰 | 專案唯一識別碼 |
| Destination | 目的地 | Teams 發送目標 |
| Actor | 執行器 | 異步任務處理單元 |
| Circuit Breaker | 熔斷器 | 故障保護機制 |
| Exponential Backoff | 指數退避 | 重試延遲策略 |
| JSONB | JSON Binary | PostgreSQL JSON 資料型態 |

### B. 參考文件

- Microsoft Teams Bot Framework: https://docs.microsoft.com/en-us/microsoftteams/platform/bots/
- Microsoft Graph API: https://docs.microsoft.com/en-us/graph/
- Go Best Practices: https://go.dev/doc/effective_go
- PostgreSQL JSONB: https://www.postgresql.org/docs/current/datatype-json.html

---

**文件結束**

如有任何技術問題，請聯絡：it-rfq@company.com

