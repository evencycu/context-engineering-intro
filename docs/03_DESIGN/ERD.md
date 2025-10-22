# 🗃️ ERD（資料實體關聯圖）

> Teams Notification Bot Platform - 完整的資料模型設計

## 核心實體關係圖

```mermaid
erDiagram
    COMPANIES {
        uuid id PK
        string name
        string contact_email
        string contact_phone
        text address
        string status "active|inactive|suspended"
        boolean billing_enabled
        timestamptz created_at
        timestamptz updated_at
    }

    USERS {
        uuid id PK
        uuid company_id FK
        string email
        string name
        string role "admin|user|viewer"
        string status "active|inactive|suspended"
        timestamptz last_login_at
        string password_hash
        string api_key_hash
        timestamptz api_key_expires_at
        timestamptz created_at
        timestamptz updated_at
    }

    PROJECTS {
        uuid id PK
        uuid company_id FK
        string notify_key
        text description
        string status "active|inactive|suspended"
        integer daily_limit
        integer monthly_limit
        string priority "low|normal|high|urgent"
        uuid created_by FK
        timestamptz created_at
        timestamptz updated_at
    }

    TEAMS_BOTS {
        uuid id PK
        string type "platform|third_party"
        uuid company_id FK
        string name
        text description
        string app_id
        string app_password_hash
        string tenant_id
        string status "active|inactive|suspended"
        string webhook_url
        jsonb capabilities
        integer rate_limit_per_minute
        integer max_concurrent_requests
        string api_endpoint
        string api_key_hash
        string contact_email
        string contact_phone
        uuid created_by FK
        timestamptz created_at
        timestamptz updated_at
    }

    BOT_INSTALLATIONS {
        uuid id PK
        uuid bot_id FK
        string bot_type "platform|third_party"
        string teams_tenant_id
        string conversation_id
        string teams_chat_id
        string teams_team_id
        string teams_channel_id
        string teams_user_id
        string scope "personal|team|groupChat"
        string service_url
        string recipient_id
        string from_id
        string from_aad_object_id
        string installation_status "active|inactive|uninstalled|stale"
        timestamptz last_activity_at
        timestamptz installed_at
        timestamptz uninstalled_at
        timestamptz created_at
        timestamptz updated_at
    }

    DESTINATIONS {
        uuid id PK
        uuid project_id FK
        string name
        text description
        string destination_type "teams_channel|teams_group|teams_user|email|webhook"
        jsonb destination_config
        string status "active|inactive|suspended"
        uuid created_by FK
        timestamptz created_at
        timestamptz updated_at
    }

    NOTIFICATIONS {
        uuid id PK
        uuid project_id FK
        string message_type "text|card|adaptive"
        text content
        jsonb metadata
        string status "pending|processing|sent|failed|cancelled"
        integer priority "1-5"
        timestamptz scheduled_at
        timestamptz sent_at
        timestamptz expires_at
        uuid created_by FK
        timestamptz created_at
        timestamptz updated_at
    }

    NOTIFICATION_DESTINATIONS {
        uuid id PK
        uuid notification_id FK
        uuid destination_id FK
        string conversation_id
        uuid bot_id FK
        string bot_type "platform|third_party"
        string status "pending|enqueued|processing|sent|failed|cancelled"
        text error_message
        string teams_message_id
        timestamptz sent_at
        integer retry_count
        integer max_retries
        timestamptz next_retry_at
        timestamptz first_attempt_at
        timestamptz last_attempt_at
        string failure_reason
        integer retry_after
        string actor_id
        string priority "low|normal|high"
        timestamptz created_at
        timestamptz updated_at
    }


    BILLING_PLANS {
        uuid id PK
        string name
        text description
        decimal price_per_notification
        decimal price_per_month
        integer max_notifications_per_month
        integer max_projects
        integer max_users
        jsonb features
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    PROJECT_BILLING {
        uuid id PK
        uuid project_id FK
        uuid billing_plan_id FK
        string billing_status "active|suspended|cancelled"
        string payment_method "credit_card|bank_transfer|invoice"
        string billing_cycle "monthly|yearly"
        date next_billing_date
        decimal total_usage_cost
        timestamptz created_at
        timestamptz updated_at
    }

    USAGE_RECORDS {
        uuid id PK
        uuid company_id FK
        uuid project_id FK
        uuid user_id FK
        string record_type "notification|api_call|storage"
        integer quantity
        decimal unit_cost
        decimal total_cost
        jsonb metadata
        timestamptz created_at
    }

    FILES {
        uuid id PK
        uuid project_id FK
        string file_name
        bigint file_size
        string content_type
        string file_key
        text file_url
        text description
        text[] tags
        boolean is_public
        timestamptz expires_at
        timestamptz created_at
        timestamptz updated_at
    }


    AUDIT_LOGS {
        uuid id PK
        uuid user_id FK
        string action
        string resource_type
        uuid resource_id
        jsonb old_values
        jsonb new_values
        string ip_address
        string user_agent
        timestamptz created_at
    }

    %% 關聯關係
    COMPANIES ||--o{ USERS : "has"
    COMPANIES ||--o{ PROJECTS : "owns"
    COMPANIES ||--o{ TEAMS_BOTS : "manages"
    PROJECTS ||--o{ PROJECT_BILLING : "billed"
    COMPANIES ||--o{ USAGE_RECORDS : "tracks"

    USERS ||--o{ PROJECTS : "creates"
    USERS ||--o{ TEAMS_BOTS : "creates"
    USERS ||--o{ DESTINATIONS : "creates"
    USERS ||--o{ NOTIFICATIONS : "sends"
    USERS ||--o{ USAGE_RECORDS : "generates"
    USERS ||--o{ AUDIT_LOGS : "performs"

    PROJECTS ||--o{ DESTINATIONS : "contains"
    PROJECTS ||--o{ NOTIFICATIONS : "sends"
    PROJECTS ||--o{ FILES : "stores"
    PROJECTS ||--o{ USAGE_RECORDS : "tracks"

    TEAMS_BOTS ||--o{ BOT_INSTALLATIONS : "installed_as"
    TEAMS_BOTS ||--o{ NOTIFICATION_DESTINATIONS : "sends_via"

    BOT_INSTALLATIONS ||--o{ NOTIFICATION_DESTINATIONS : "receives"

    DESTINATIONS ||--o{ NOTIFICATION_DESTINATIONS : "targets"

    NOTIFICATIONS ||--o{ NOTIFICATION_DESTINATIONS : "sends_to"
    NOTIFICATIONS ||--o{ USAGE_RECORDS : "generates"



    BILLING_PLANS ||--o{ PROJECT_BILLING : "plans"
```

## 核心實體說明

### 1. 公司與用戶管理
- **COMPANIES**: 公司/組織實體
- **USERS**: 系統用戶，支援多角色
- **PROJECTS**: 專案管理，每個公司可有多個專案

### 2. Bot 管理
- **TEAMS_BOTS**: 統一的 Bot 管理（平台 Bot + 第三方 Bot）
- **BOT_INSTALLATIONS**: Bot 安裝記錄，支援多種安裝範圍

### 3. 通知系統
- **NOTIFICATIONS**: 通知主體
- **NOTIFICATION_DESTINATIONS**: 通知目的地，支援異步處理
- **DESTINATIONS**: 目的地配置


### 4. 計費系統
- **BILLING_PLANS**: 計費方案
- **PROJECT_BILLING**: 專案計費記錄
- **USAGE_RECORDS**: 使用記錄

### 5. 檔案管理
- **FILES**: 檔案存儲

### 6. 審計
- **AUDIT_LOGS**: 操作審計日誌

## 索引設計

### 主要查詢索引
```sql
-- 通知查詢
CREATE INDEX idx_notifications_project_status ON notifications(project_id, status);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- 通知目的地查詢
CREATE INDEX idx_notification_destinations_status ON notification_destinations(status);
CREATE INDEX idx_notification_destinations_next_retry ON notification_destinations(next_retry_at) 
    WHERE retry_count < max_retries AND status = 'pending';

-- Bot 安裝查詢
CREATE INDEX idx_bot_installations_bot_tenant ON bot_installations(bot_id, bot_type, teams_tenant_id);
CREATE INDEX idx_bot_installations_conversation ON bot_installations(conversation_id);

-- 使用記錄查詢
CREATE INDEX idx_usage_records_company_created ON usage_records(company_id, created_at);
CREATE INDEX idx_usage_records_project_created ON usage_records(project_id, created_at);
```

### 唯一約束
```sql
-- 專案通知鍵唯一性
ALTER TABLE projects ADD CONSTRAINT projects_company_notify_key_unique 
    UNIQUE (company_id, notify_key);

-- Bot 安裝唯一性
ALTER TABLE bot_installations ADD CONSTRAINT bot_installations_bot_tenant_conversation_unique 
    UNIQUE (bot_id, bot_type, teams_tenant_id, conversation_id);

-- 專案計費唯一性
ALTER TABLE project_billing ADD CONSTRAINT project_billing_project_unique 
    UNIQUE (project_id);
```

## 資料庫配置

### 連接字串
```
postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable
```

### 主要配置
- **資料庫名稱**: `notification_center`
- **字符集**: UTF-8
- **時區**: UTC
- **SSL**: 開發環境禁用，生產環境啟用

## 遷移策略

### 版本控制
- 使用 `scripts/migrations/` 目錄管理資料庫變更
- 每個遷移檔案包含版本號和描述
- 支援前向和後向遷移

### 主要遷移
1. **001_enhance_bot_installations.sql**: 增強 Bot 安裝表
2. **006_unify_bots_to_teams_bots.sql**: 統一 Bot 管理
3. **009_enhance_notification_destinations_for_async_actor.sql**: 支援異步處理
4. **011_add_billing_tables.sql**: 添加計費系統
5. **012_add_file_tables.sql**: 添加檔案管理

## 性能優化

### 分區策略
- 按公司分區 `usage_records` 表

### 查詢優化
- 使用適當的索引
- 避免全表掃描
- 使用 EXPLAIN 分析查詢計劃

### 維護策略
- 定期清理過期資料
- 監控索引使用情況
- 優化慢查詢

---

**版本**: v1.0  
**最後更新**: 2025-10-08  
**作者**: TeamsNotify Team
