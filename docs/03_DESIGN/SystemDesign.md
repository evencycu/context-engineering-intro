# 💻 System Design

## 1. 文件資訊
- **版本**：v1.0
- **作者**：Tech Lead
- **最後更新**：2025-10-08

---

## 2. 設計目標
- 將需求（PRD、SRS）落地成可實作的技術設計
- 清楚描述模組邊界、資料流、錯誤處理與擴展點
- 支援企業級 Teams 通知服務的完整功能

---

## 3. 系統架構概覽

### 3.1 整體架構
```mermaid
graph TB
    subgraph "外部系統"
        A[應用系統 A]
        B[應用系統 B]
        C[監控系統]
    end
  
    subgraph "API Gateway"
        D[Teams Notification API]
        E[認證授權]
        F[速率限制]
    end
  
    subgraph "核心服務"
        G[通知服務]
        H[專案管理]
        I[目的地管理]
        J[Bot 管理]
        K[批次處理]
    end
  
    subgraph "佇列系統"
        L[Redis 佇列]
        M[Actor Pool]
        N[重試機制]
    end
  
    subgraph "資料層"
        O[PostgreSQL]
        P[Redis Cache]
    end
  
    subgraph "Teams 平台"
        Q[Teams Bot Framework]
        R[Teams 頻道]
        S[Teams 群組]
        T[Teams 個人]
    end
  
    A --> D
    B --> D
    C --> D
    D --> E
    D --> F
    D --> G
    G --> H
    G --> I
    G --> J
    G --> K
    G --> L
    L --> M
    M --> N
    M --> Q
    Q --> R
    Q --> S
    Q --> T
    G --> O
    G --> P
```

### 3.2 技術棧
| 類別 | 技術 | 用途 |
|------|------|------|
| Backend | Golang (Gin) | RESTful API 與佇列處理 |
| Database | PostgreSQL | 永久儲存與關聯查詢 |
| Cache/Queue | Redis | 快取與訊息佇列 |
| Infrastructure | Docker + Kubernetes | 部署與監控 |
| Authentication | Azure Entra ID | 驗證與授權 |
| Monitoring | 內建健康檢查 | 監控與日誌收集 |

---

## 4. 模組劃分與邊界

| 模組 | 職責 | 主要輸入 | 主要輸出 |
|------|------|----------|----------|
| **API Gateway** | 對外 REST 介面、驗證、參數校驗 | HTTP Request | Job Enqueue / DB Record |
| **Notification Service** | 通知管理邏輯、批次處理 | Notification Request | Notification Record |
| **Project Service** | 專案管理、權限控制 | Project Request | Project Record |
| **Destination Service** | 目的地管理、配置 | Destination Request | Destination Record |
| **Bot Service** | Bot 管理、安裝追蹤 | Bot Request | Bot Record |
| **Queue Manager** | 將訊息放入佇列，提供重試與去重 | Message DTO | Queue Item |
| **Actor Engine** | 從佇列取出並執行發送邏輯 | Queue Item | 成功/失敗結果、紀錄 |
| **Token Manager** | Token 快取與管理 | Token Request | Cached Token |
| **Integrations** | 與外部服務交互 (Teams, Graph) | DTO | HTTP Response |
| **Persistence** | 儲存 domain 資料與事件 | Domain Model | Tables / Views |

---

## 5. 主要流程（時序圖）

### 5.1 通知推播（後端主動發送）
```mermaid
sequenceDiagram
    participant Client
    participant API as Notification API
    participant NS as Notification Service
    participant DB as PostgreSQL
    participant EW as EnqueueWorker
    participant RQ as Redis Queue
    participant QC as QueueConsumer
    participant AP as ActorPool
    participant NA as NotificationActor
    participant TM as TokenManager
    participant Teams as Teams API

    Client->>API: POST /external/notify
    API->>API: Validate & AuthZ
    API->>NS: SendNotification()
    NS->>DB: INSERT notification (status=pending)
    NS->>DB: INSERT notification_destinations (status=pending)
    NS-->>API: Success (async)
    API-->>Client: 200 OK (processing)
    
    Note over EW: Background Process
    EW->>DB: SELECT pending notifications
    EW->>RQ: Enqueue by priority
    EW->>DB: UPDATE status=enqueued
    
    Note over QC: Background Process
    QC->>RQ: BLPOP from queue
    QC->>AP: SpawnActor()
    AP->>NA: Create NotificationActor
    NA->>TM: GetConnectorToken()
    TM-->>NA: Cached Token
    NA->>Teams: POST /v3/conversations/{id}/activities
    Teams-->>NA: 200 OK / Error
    NA->>DB: UPDATE status=sent/failed
    NA->>AP: Actor Complete
```


### 5.2 Bot 安裝與管理
```mermaid
sequenceDiagram
    participant Teams
    participant Bot as Teams Bot
    participant API as Bot API
    participant BS as Bot Service
    participant DB as PostgreSQL
    participant Webhook as Webhook Handler

    Teams->>Bot: Installation Event
    Bot->>API: POST /webhook/installation
    API->>Webhook: Handle Installation
    Webhook->>BS: Process Installation
    BS->>DB: INSERT/UPDATE bot_installations
    BS->>DB: UPDATE installation_status
    BS-->>API: Installation Processed
    API-->>Bot: 200 OK
    Bot-->>Teams: Installation Complete
```

---

## 6. 資料模型（同步 schema.sql/init.sql，詳細見 ERD）
- 參考：`ERD.md`
- 覆蓋：companies、users、projects、teams_bots、bot_installations、notifications、notification_destinations、destinations、billing_plans、project_billing、usage_records、files、audit_logs

### 6.0 本次同步要點
- 計費改為 `project_billing`（唯一鍵：`project_id`），完全取代 `company_billing`
- `usage_records` 僅保留 `project_id`（移除 `company_id`），並在 `usage_summary` 檢視以 `project_id/notify_key/company` 彙總
- `destinations.targets` 採 JSONB 陣列，型別與欄位約束：
  - `type`: `personal` | `groupchat` | `channel`
  - `tenant_id`（必填）
  - `conversation_id`（groupchat/channel 必填；personal 與 email 二擇一）
  - `email`（personal 可用）
- `notifications` 與 `notification_destinations`：以 `targets` 展開建立分派記錄，`notification_destinations` 含重試與 `conversation_id` 欄位
- 檢視：`notification_summary`、`usage_summary`、`bot_status_summary` 已可直接用於報表

### 6.1 外部 API 相關
- `/api/v1/notify` 支援 `targets=["all"|conversation_id|email]`，依 `notify_key` 展開目的地
- `/api/v1/destinations/{notifyKey}` 回傳該專案目的地與 targets 明細（對應 `destinations` 與 JSONB targets）

### 6.1 資料庫配置
- **資料庫名稱**: `notification_center`
- **連接字串**: `postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable`
- **字符集**: UTF-8
- **時區**: UTC
- **SSL**: 開發環境禁用，生產環境啟用


---

## 7. 錯誤處理與重試

### 7.1 錯誤分類
| 錯誤類型 | HTTP 狀態碼 | 處理策略 | 重試機制 |
|----------|-------------|----------|----------|
| 參數錯誤 | 400 | 直接返回錯誤 | 不重試 |
| 認證失敗 | 401 | 返回認證錯誤 | 不重試 |
| 權限不足 | 403 | 返回權限錯誤 | 不重試 |
| 資源不存在 | 404 | 返回資源錯誤 | 不重試 |
| 速率限制 | 429 | 等待 Retry-After | 指數退避 |
| 服務錯誤 | 5xx | 重試機制 | 指數退避 |

#### 詳細錯誤類型
- **認證錯誤**: UNAUTHORIZED, INVALID_TOKEN, TOKEN_EXPIRED, INVALID_API_KEY
- **Teams API 錯誤**: TEAMS_AUTH_FAILED, TEAMS_API_ERROR, TEAMS_RATE_LIMIT
- **數據庫錯誤**: DATABASE_ERROR, RECORD_NOT_FOUND, DUPLICATE_RECORD
- **驗證錯誤**: VALIDATION_FAILED, INVALID_INPUT, MISSING_FIELD
- **系統錯誤**: INTERNAL_ERROR, SERVICE_UNAVAILABLE, TIMEOUT

### 7.2 重試策略
```go
type RetryPolicy struct {
    MaxRetries     int           // 最多重試次數
    InitialBackoff time.Duration // 初始等待時間
    MaxBackoff     time.Duration // 最長等待時間
    Multiplier     float64       // 退避倍數
    JitterFraction float64       // 隨機抖動
}
```

### 7.3 去重策略
- **SETNX 鎖定**: `queue:notifications:processing:{id}` (TTL: 30分鐘)
- **資料庫唯一約束**: 防止重複處理
- **冪等性設計**: 支援重複請求

---

## 8. 安全性

### 8.1 認證與授權
- **認證**: Azure Entra ID (OAuth2/OIDC)
- **授權**: RBAC 角色 (admin, operator, viewer)
- **API Key**: 服務間認證
- **JWT Token**: 用戶會話管理

### 8.2 傳輸安全
- **HTTPS**: 強制加密傳輸
- **HSTS**: HTTP Strict Transport Security
- **TLS 1.2+**: 最低加密標準

### 8.3 資料安全
- **敏感資料加密**: 密碼、API Key 加密存儲
- **Key Vault**: 生產環境密鑰管理
- **日誌遮罩**: PII/Token 遮罩處理

---

## 9. 觀測性（Observability）

### 9.1 監控指標
| 指標類型 | 指標名稱 | 描述 | 告警閾值 |
|----------|----------|------|----------|
| 業務指標 | notification_send_total | 通知發送總數 | - |
| 業務指標 | notification_success_rate | 發送成功率 | < 95% |
| 性能指標 | api_response_time_p95 | API 響應時間 P95 | > 2s |
| 性能指標 | queue_length | 佇列長度 | > 1000 |
| 錯誤指標 | error_rate_4xx | 4xx 錯誤率 | > 5% |
| 錯誤指標 | error_rate_5xx | 5xx 錯誤率 | > 1% |

### 9.2 日誌管理
- **格式**: 結構化 JSON
- **欄位**: trace_id, span_id, request_id
- **級別**: Debug, Info, Warn, Error
- **聚合**: 集中式日誌收集

### 9.3 追蹤
- **OpenTelemetry**: 分散式追蹤
- **Span**: API → Worker → 外部服務
- **Context**: 請求上下文傳遞

---

## 10. 配置與開關

| Key | 範例值 | 說明 |
|-----|--------|------|
| **佇列配置** | | |
| QUEUE_MAX_RETRY | 5 | 最大重試次數 |
| QUEUE_BATCH_SIZE | 100 | 批次處理大小 |
| QUEUE_SCAN_INTERVAL | 200ms | 掃描間隔 |
| **外部服務** | | |
| GRAPH_BASE_URL | https://graph.microsoft.com/ | Graph 端點 |
| BOT_SVC_URL | https://smba.trafficmanager.net/... | Bot Framework 端點 |
| **速率限制** | | |
| RATE_LIMIT_QPS | 50 | 對外呼叫限速 |
| RATE_LIMIT_BURST | 100 | 突發請求限制 |
| **Token 管理** | | |
| TOKEN_CACHE_TTL | 3000s | Token 快取時間 |
| TOKEN_REFRESH_AHEAD | 300s | 提前刷新時間 |
| **Actor 配置** | | |
| ACTOR_MAX_COUNT | 10 | 最大 Actor 數量 |
| ACTOR_TIMEOUT | 30s | Actor 超時時間 |

---


## 11. 性能優化

### 11.1 快取策略
- **Token 快取**: Redis 存儲，50 分鐘 TTL
- **配置快取**: 內存存儲，啟動時載入
- **資料快取**: 查詢結果快取

### 11.2 並發處理
- **Goroutine 池**: 控制並發數量
- **Actor Pool**: 異步任務處理
- **連接池**: 資料庫和 Redis 連接池

### 11.3 資料庫優化
- **索引設計**: 針對查詢模式優化
- **分區策略**: 按時間分區大表
- **查詢優化**: 避免 N+1 查詢

---

## 12. 延伸點

### 12.1 多租戶支援
- **租戶隔離**: TenantId 列、索引、隔離
- **資源限制**: 按租戶限制資源使用
- **計費管理**: 多租戶計費系統

### 12.2 多管道支援
- **Email**: SMTP 發送
- **SMS**: 簡訊服務
- **LINE**: LINE Bot 整合
- **Slack**: Slack Bot 整合

### 12.3 佇列系統擴展
- **Kafka**: 替換 Redis Stream
- **抽象層**: QueueProvider 介面
- **多佇列**: 支援不同優先級佇列

---

## 13. 故障恢復

### 13.1 故障檢測
- **健康檢查**: 定期健康檢查
- **熔斷器**: 故障隔離
- **重試機制**: 自動重試

### 13.2 故障恢復
- **自動重啟**: 服務自動重啟
- **故障轉移**: 流量切換
- **資料恢復**: 資料一致性保證

---

## 14. 監控與維護

### 14.1 監控儀表板
- **系統指標**: CPU、內存、磁盤
- **應用指標**: 請求數、響應時間、錯誤率
- **業務指標**: 通知發送量、成功率

### 14.2 告警機制
- **即時告警**: 關鍵錯誤即時通知
- **趨勢告警**: 性能趨勢異常告警
- **容量告警**: 資源使用率告警

### 14.3 維護操作
- **定期備份**: 資料庫備份策略
- **日誌清理**: 定期清理過期日誌
- **性能調優**: 定期性能分析與優化

## 15. API 端點概覽

### 15.1 核心管理 API
- **公司管理** (`/companies`) - 公司資料管理
- **用戶管理** (`/users`) - 用戶帳號和權限管理
- **專案管理** (`/projects`) - 通知專案管理
- **Bot 管理** (`/bots`) - Teams Bot 管理
- **目的地管理** (`/destinations`) - 通知目的地管理
- **通知管理** (`/notifications`) - 通知發送和管理

### 15.2 特殊功能 API（外部）
- **Provision API** (`/provision`) - 一鍵建立專案和目的地
- **Notify API** (`/api/v1/notify`) - 外部系統通知發送
- **Destinations API** (`/api/v1/destinations/{notifyKey}`) - 查詢專案目的地
- **Queue Management API** (`/queue`) - 內部佇列管理和監控
- **訊息管理** (`/messages`) - 訊息歷史和狀態查詢

### 15.3 基礎 URL
```
http://localhost:8080/api/v1
```

## 16. 序列流程設計

### 16.1 通知發送完整流程
- **成功場景**: 客戶端 → External Handler → External Service → Broadcast Service → Teams API → Database
- **失敗場景**: 包含重試機制和電路斷路器保護
- **多節點部署**: 支援水平擴展和並行處理

### 16.2 佇列管理流程
- **Redis Queue**: 使用 Redis 實現分散式隊列
- **Redis Circuit Breaker**: 使用 Redis 共享熔斷器狀態
- **Dead Letter Queue**: 達到最大重試後標記為 failed
- **智能退避**: 優先使用 429 Retry-After，否則使用指數退避

---

**版本**: v1.0  
**最後更新**: 2025-10-08  
**作者**: TeamsNotify Team
