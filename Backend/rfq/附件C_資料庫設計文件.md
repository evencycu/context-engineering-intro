# 附件 C：資料庫設計文件
# Database Design Document

**RFQ 編號**: RFQ-2025-TEAMS-NOTIFY-001  
**版本**: v1.0  
**日期**: 2025年10月23日

---

## 目錄

1. [資料庫概述](#1-資料庫概述)
2. [實體關係圖 (ERD)](#2-實體關係圖-erd)
3. [資料表詳細設計](#3-資料表詳細設計)
4. [索引設計](#4-索引設計)
5. [檢視設計](#5-檢視設計)
6. [資料遷移策略](#6-資料遷移策略)
7. [效能優化](#7-效能優化)

---

## 1. 資料庫概述

### 1.1 資料庫配置

| 項目 | 規格 |
|------|------|
| **資料庫系統** | PostgreSQL 14+ |
| **資料庫名稱** | notification_center |
| **字符集** | UTF-8 |
| **時區** | UTC |
| **連接池** | Min: 5, Max: 25 |
| **SSL** | 生產環境：啟用，開發環境：停用 |

### 1.2 資料表清單

| 表名 | 說明 | 預估資料量 | 分區 |
|------|------|-----------|------|
| companies | 公司資料 | < 1,000 | ❌ |
| users | 用戶帳號 | < 10,000 | ❌ |
| projects | 通知專案 | < 5,000 | ❌ |
| teams_bots | Teams Bot 資訊 | < 100 | ❌ |
| destinations | 目的地配置 | < 50,000 | ❌ |
| notifications | 通知記錄 | > 1,000,000 | ✅ 按月 |
| notification_destinations | 通知分派記錄 | > 10,000,000 | ✅ 按月 |
| usage_records | 使用量記錄 | > 100,000 | ❌ |
| project_billing | 專案計費 | < 5,000 | ❌ |
| billing_plans | 計費方案 | < 50 | ❌ |
| files | 檔案資訊 | < 10,000 | ❌ |
| audit_logs | 審計日誌 | > 1,000,000 | ✅ 按月 |

---

## 2. 實體關係圖 (ERD)

### 2.1 核心資料關係

```
┌─────────────┐
│  companies  │
│  (公司)     │
└──────┬──────┘
       │ 1:N
       │
   ┌───┴───┬───────────────┬──────────┐
   │       │               │          │
   ▼       ▼               ▼          ▼
┌──────┐ ┌──────────┐  ┌──────────┐ ┌──────────────┐
│users │ │ projects  │  │ billing  │ │usage_records │
│(用戶)│ │ (專案)    │  │ (計費)   │ │(使用量)      │
└───┬──┘ └────┬──────┘  └──────────┘ └──────────────┘
    │         │
    │ N:1     │ 1:N
    │         │
    │    ┌────┴─────────┐
    │    │              │
    │    ▼              ▼
    │ ┌─────────────┐ ┌──────────────┐
    │ │destinations │ │notifications │
    │ │(目的地)     │ │(通知)        │
    │ └──────┬──────┘ └──────┬───────┘
    │        │ 1:N           │ 1:N
    │        │               │
    │        │      ┌────────┴─────────┐
    │        │      │                  │
    │        │      ▼                  │
    │        │   ┌──────────────────┐  │
    │        └──>│notification_     │<─┘
    │            │destinations      │
    │            │(通知分派)        │
    └───────────>└──────────────────┘
                 (sender_id)


┌──────────────┐
│ teams_bots   │
│ (Teams Bot)  │
└──────┬───────┘
       │ 1:N
       │
       ▼
┌──────────────┐
│ destinations │
│ (目的地)     │
└──────────────┘
```

### 2.2 計費關係

```
┌──────────────┐       ┌──────────────┐
│billing_plans │       │  projects    │
│ (計費方案)   │       │  (專案)      │
└──────┬───────┘       └──────┬───────┘
       │ 1:N                  │ 1:1
       │                      │
       └──────────┬───────────┘
                  │
                  ▼
          ┌────────────────┐
          │project_billing │
          │ (專案計費)     │
          └────────────────┘
```

---

## 3. 資料表詳細設計

### 3.1 companies (公司)

**說明**: 多租戶公司資料表

```sql
CREATE TABLE companies (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name              VARCHAR(255) NOT NULL,
    contact_email     VARCHAR(255) NOT NULL,
    contact_phone     VARCHAR(50),
    address           TEXT,
    status            VARCHAR(50) NOT NULL DEFAULT 'active',
    billing_enabled   BOOLEAN NOT NULL DEFAULT false,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- 約束
    CONSTRAINT chk_companies_status CHECK (status IN ('active', 'inactive', 'suspended')),
    CONSTRAINT chk_companies_email CHECK (contact_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- 註解
COMMENT ON TABLE companies IS '公司資料表';
COMMENT ON COLUMN companies.status IS 'active: 啟用, inactive: 停用, suspended: 暫停';
COMMENT ON COLUMN companies.billing_enabled IS '是否啟用計費功能';
```

**索引**:
```sql
CREATE INDEX idx_companies_status ON companies(status);
CREATE INDEX idx_companies_created_at ON companies(created_at DESC);
```

---

### 3.2 users (用戶)

**說明**: 系統用戶帳號表

```sql
CREATE TABLE users (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id        UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    email             VARCHAR(255) NOT NULL UNIQUE,
    name              VARCHAR(255) NOT NULL,
    password_hash     VARCHAR(255) NOT NULL,
    role              VARCHAR(50) NOT NULL DEFAULT 'user',
    status            VARCHAR(50) NOT NULL DEFAULT 'active',
    last_login_at     TIMESTAMP,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- 約束
    CONSTRAINT chk_users_role CHECK (role IN ('admin', 'manager', 'user')),
    CONSTRAINT chk_users_status CHECK (status IN ('active', 'inactive'))
);

-- 註解
COMMENT ON TABLE users IS '系統用戶表';
COMMENT ON COLUMN users.role IS 'admin: 管理員, manager: 經理, user: 一般用戶';
COMMENT ON COLUMN users.password_hash IS 'bcrypt 加密的密碼';
```

**索引**:
```sql
CREATE INDEX idx_users_company_id ON users(company_id);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE UNIQUE INDEX idx_users_email ON users(email);
```

---

### 3.3 projects (專案)

**說明**: 通知專案表（Notification Key 管理）

```sql
CREATE TABLE projects (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id        UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    notify_key        VARCHAR(255) NOT NULL UNIQUE,
    description       TEXT,
    daily_limit       INTEGER NOT NULL DEFAULT 0,
    monthly_limit     INTEGER NOT NULL DEFAULT 0,
    priority          VARCHAR(50) NOT NULL DEFAULT 'normal',
    status            VARCHAR(50) NOT NULL DEFAULT 'active',
    created_by        UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- 約束
    CONSTRAINT chk_projects_priority CHECK (priority IN ('low', 'normal', 'high')),
    CONSTRAINT chk_projects_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT chk_projects_limits CHECK (daily_limit >= 0 AND monthly_limit >= 0)
);

-- 註解
COMMENT ON TABLE projects IS '通知專案表';
COMMENT ON COLUMN projects.notify_key IS '專案唯一識別碼，用於外部 API 調用';
COMMENT ON COLUMN projects.daily_limit IS '每日通知配額，0 表示不限制';
COMMENT ON COLUMN projects.monthly_limit IS '每月通知配額，0 表示不限制';
COMMENT ON COLUMN projects.priority IS '專案優先級，影響佇列處理順序';
```

**索引**:
```sql
CREATE UNIQUE INDEX idx_projects_notify_key ON projects(notify_key);
CREATE INDEX idx_projects_company_id ON projects(company_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_priority ON projects(priority);
```

---

### 3.4 teams_bots (Teams Bot)

**說明**: Microsoft Teams Bot 資訊表

```sql
CREATE TABLE teams_bots (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name              VARCHAR(255) NOT NULL,
    app_id            VARCHAR(255) NOT NULL UNIQUE,
    app_password      TEXT NOT NULL,  -- 加密存儲
    capabilities      JSONB NOT NULL DEFAULT '[]',
    status            VARCHAR(50) NOT NULL DEFAULT 'active',
    description       TEXT,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- 約束
    CONSTRAINT chk_bots_status CHECK (status IN ('active', 'inactive', 'error'))
);

-- 註解
COMMENT ON TABLE teams_bots IS 'Microsoft Teams Bot 資訊表';
COMMENT ON COLUMN teams_bots.app_id IS 'Microsoft App ID (Bot ID)';
COMMENT ON COLUMN teams_bots.app_password IS 'Microsoft App Password (加密存儲)';
COMMENT ON COLUMN teams_bots.capabilities IS 'Bot 能力陣列，如: ["proactive_messaging", "bot_messaging"]';
```

**索引**:
```sql
CREATE UNIQUE INDEX idx_bots_app_id ON teams_bots(app_id);
CREATE INDEX idx_bots_status ON teams_bots(status);
```

---

### 3.5 destinations (目的地)

**說明**: 通知目的地配置表

```sql
CREATE TABLE destinations (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id        UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    bot_id            UUID NOT NULL REFERENCES teams_bots(id) ON DELETE RESTRICT,
    type              VARCHAR(50) NOT NULL,
    targets           JSONB NOT NULL,  -- 目標陣列
    status            VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- 約束
    CONSTRAINT chk_destinations_type CHECK (type IN ('personal', 'groupchat', 'channel')),
    CONSTRAINT chk_destinations_status CHECK (status IN ('active', 'inactive'))
);

-- 註解
COMMENT ON TABLE destinations IS '通知目的地配置表';
COMMENT ON COLUMN destinations.type IS 'personal: 個人, groupchat: 群組, channel: 頻道';
COMMENT ON COLUMN destinations.targets IS 'JSONB 陣列，包含 type, conversation_id, email, display_name, tenant_id';
```

**Targets JSONB 結構範例**:
```json
[
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
```

**索引**:
```sql
CREATE INDEX idx_destinations_project_id ON destinations(project_id);
CREATE INDEX idx_destinations_bot_id ON destinations(bot_id);
CREATE INDEX idx_destinations_type ON destinations(type);
CREATE INDEX idx_destinations_status ON destinations(status);
CREATE INDEX idx_destinations_targets_gin ON destinations USING gin(targets);  -- JSONB 索引
```

---

### 3.6 notifications (通知)

**說明**: 通知記錄表（主表）

```sql
CREATE TABLE notifications (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id              UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    sender_id               UUID REFERENCES users(id) ON DELETE SET NULL,
    message_type            VARCHAR(50) NOT NULL DEFAULT 'text',
    content                 TEXT NOT NULL,
    mentions                JSONB DEFAULT '[]',
    priority                VARCHAR(50) NOT NULL DEFAULT 'normal',
    status                  VARCHAR(50) NOT NULL DEFAULT 'pending',
    error_message           TEXT,
    destinations_sent       INTEGER NOT NULL DEFAULT 0,
    destinations_failed     INTEGER NOT NULL DEFAULT 0,
    total_destinations      INTEGER NOT NULL DEFAULT 0,
    metadata                JSONB,
    created_at              TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMP NOT NULL DEFAULT NOW(),
    sent_at                 TIMESTAMP,
    
    -- 約束
    CONSTRAINT chk_notifications_message_type CHECK (message_type IN ('text', 'file', 'adaptive_card')),
    CONSTRAINT chk_notifications_priority CHECK (priority IN ('low', 'normal', 'high')),
    CONSTRAINT chk_notifications_status CHECK (status IN ('pending', 'enqueued', 'sending', 'sent', 'failed', 'cancelled')),
    CONSTRAINT chk_notifications_content_length CHECK (LENGTH(content) <= 4000)
) PARTITION BY RANGE (created_at);  -- 按月分區

-- 註解
COMMENT ON TABLE notifications IS '通知記錄表（主表）';
COMMENT ON COLUMN notifications.message_type IS 'text: 文字, file: 檔案, adaptive_card: Adaptive Card';
COMMENT ON COLUMN notifications.mentions IS '@提及的用戶 email 陣列';
COMMENT ON COLUMN notifications.status IS '狀態流轉: pending → enqueued → sending → sent/failed';
```

**分區範例**:
```sql
-- 2025年10月分區
CREATE TABLE notifications_2025_10 PARTITION OF notifications
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

-- 2025年11月分區
CREATE TABLE notifications_2025_11 PARTITION OF notifications
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
```

**索引**:
```sql
CREATE INDEX idx_notifications_project_id ON notifications(project_id);
CREATE INDEX idx_notifications_sender_id ON notifications(sender_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_priority ON notifications(priority);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX idx_notifications_sent_at ON notifications(sent_at DESC);
```

---

### 3.7 notification_destinations (通知分派)

**說明**: 通知分派記錄表（詳細發送狀態）

```sql
CREATE TABLE notification_destinations (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id     UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    destination_id      UUID NOT NULL REFERENCES destinations(id) ON DELETE CASCADE,
    conversation_id     VARCHAR(255),
    target_email        VARCHAR(255),
    status              VARCHAR(50) NOT NULL DEFAULT 'pending',
    error_message       TEXT,
    retry_count         INTEGER NOT NULL DEFAULT 0,
    next_retry_at       TIMESTAMP,
    sent_at             TIMESTAMP,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- 約束
    CONSTRAINT chk_notif_dest_status CHECK (status IN ('pending', 'enqueued', 'sending', 'sent', 'failed')),
    CONSTRAINT chk_notif_dest_retry_count CHECK (retry_count >= 0 AND retry_count <= 10)
) PARTITION BY RANGE (created_at);  -- 按月分區

-- 註解
COMMENT ON TABLE notification_destinations IS '通知分派記錄表';
COMMENT ON COLUMN notification_destinations.conversation_id IS 'Teams 對話 ID';
COMMENT ON COLUMN notification_destinations.retry_count IS '重試次數，最多 10 次';
COMMENT ON COLUMN notification_destinations.next_retry_at IS '下次重試時間（指數退避）';
```

**分區範例**:
```sql
CREATE TABLE notification_destinations_2025_10 PARTITION OF notification_destinations
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');
```

**索引**:
```sql
CREATE INDEX idx_notif_dest_notification_id ON notification_destinations(notification_id);
CREATE INDEX idx_notif_dest_destination_id ON notification_destinations(destination_id);
CREATE INDEX idx_notif_dest_status ON notification_destinations(status);
CREATE INDEX idx_notif_dest_retry ON notification_destinations(retry_count, next_retry_at) WHERE status = 'failed';
CREATE INDEX idx_notif_dest_created_at ON notification_destinations(created_at DESC);
```

---

### 3.8 usage_records (使用量記錄)

**說明**: 專案使用量記錄表

```sql
CREATE TABLE usage_records (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id            UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    notification_count    INTEGER NOT NULL DEFAULT 0,
    date                  DATE NOT NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- 唯一約束：每個專案每天只能有一筆記錄
    CONSTRAINT uq_usage_records_project_date UNIQUE (project_id, date)
);

-- 註解
COMMENT ON TABLE usage_records IS '專案使用量記錄表';
COMMENT ON COLUMN usage_records.notification_count IS '當日發送通知數量';
```

**索引**:
```sql
CREATE INDEX idx_usage_records_project_id ON usage_records(project_id);
CREATE INDEX idx_usage_records_date ON usage_records(date DESC);
CREATE INDEX idx_usage_records_project_date ON usage_records(project_id, date);
```

---

### 3.9 billing_plans (計費方案)

**說明**: 計費方案表

```sql
CREATE TABLE billing_plans (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                  VARCHAR(255) NOT NULL,
    description           TEXT,
    price                 DECIMAL(10,2) NOT NULL,
    notification_limit    INTEGER NOT NULL,
    features              JSONB DEFAULT '[]',
    status                VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at            TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- 約束
    CONSTRAINT chk_billing_plans_price CHECK (price >= 0),
    CONSTRAINT chk_billing_plans_limit CHECK (notification_limit >= 0),
    CONSTRAINT chk_billing_plans_status CHECK (status IN ('active', 'inactive'))
);

-- 註解
COMMENT ON TABLE billing_plans IS '計費方案表';
COMMENT ON COLUMN billing_plans.price IS '月費（新台幣）';
COMMENT ON COLUMN billing_plans.notification_limit IS '每月通知配額';
COMMENT ON COLUMN billing_plans.features IS '方案功能陣列';
```

**索引**:
```sql
CREATE INDEX idx_billing_plans_status ON billing_plans(status);
```

---

### 3.10 project_billing (專案計費)

**說明**: 專案計費設定表

```sql
CREATE TABLE project_billing (
    project_id            UUID PRIMARY KEY REFERENCES projects(id) ON DELETE CASCADE,
    plan_id               UUID REFERENCES billing_plans(id) ON DELETE SET NULL,
    billing_enabled       BOOLEAN NOT NULL DEFAULT false,
    payment_method        VARCHAR(50),
    billing_cycle         VARCHAR(50) DEFAULT 'monthly',
    next_billing_date     DATE,
    billing_status        VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at            TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- 約束
    CONSTRAINT chk_project_billing_payment_method CHECK (payment_method IN ('credit_card', 'bank_transfer', 'paypal')),
    CONSTRAINT chk_project_billing_cycle CHECK (billing_cycle IN ('monthly', 'yearly')),
    CONSTRAINT chk_project_billing_status CHECK (billing_status IN ('active', 'suspended', 'cancelled'))
);

-- 註解
COMMENT ON TABLE project_billing IS '專案計費設定表';
COMMENT ON COLUMN project_billing.project_id IS '唯一鍵，一個專案只能有一筆計費記錄';
```

**索引**:
```sql
CREATE INDEX idx_project_billing_plan_id ON project_billing(plan_id);
CREATE INDEX idx_project_billing_status ON project_billing(billing_status);
CREATE INDEX idx_project_billing_next_date ON project_billing(next_billing_date);
```

---

### 3.11 files (檔案)

**說明**: 檔案資訊表

```sql
CREATE TABLE files (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename          VARCHAR(255) NOT NULL,
    original_name     VARCHAR(255) NOT NULL,
    size              BIGINT NOT NULL,
    mime_type         VARCHAR(100) NOT NULL,
    description       TEXT,
    tags              JSONB DEFAULT '[]',
    url               TEXT NOT NULL,
    uploaded_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- 約束
    CONSTRAINT chk_files_size CHECK (size > 0 AND size <= 104857600)  -- 最大 100MB
);

-- 註解
COMMENT ON TABLE files IS '檔案資訊表';
COMMENT ON COLUMN files.size IS '檔案大小（bytes）';
COMMENT ON COLUMN files.tags IS '檔案標籤陣列';
```

**索引**:
```sql
CREATE INDEX idx_files_uploaded_by ON files(uploaded_by);
CREATE INDEX idx_files_created_at ON files(created_at DESC);
CREATE INDEX idx_files_tags_gin ON files USING gin(tags);
```

---

### 3.12 audit_logs (審計日誌)

**說明**: 系統審計日誌表

```sql
CREATE TABLE audit_logs (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID REFERENCES users(id) ON DELETE SET NULL,
    action            VARCHAR(100) NOT NULL,
    resource_type     VARCHAR(100) NOT NULL,
    resource_id       UUID,
    details           JSONB,
    ip_address        VARCHAR(50),
    user_agent        TEXT,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);  -- 按月分區

-- 註解
COMMENT ON TABLE audit_logs IS '系統審計日誌表';
COMMENT ON COLUMN audit_logs.action IS '操作類型，如: CREATE, UPDATE, DELETE';
COMMENT ON COLUMN audit_logs.resource_type IS '資源類型，如: company, project, notification';
COMMENT ON COLUMN audit_logs.details IS '操作詳細資料（JSONB）';
```

**分區範例**:
```sql
CREATE TABLE audit_logs_2025_10 PARTITION OF audit_logs
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');
```

**索引**:
```sql
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
```

---

## 4. 索引設計

### 4.1 主鍵索引

所有表自動建立主鍵索引（UUID）

### 4.2 外鍵索引

所有外鍵自動建立索引以提升 JOIN 性能

### 4.3 查詢優化索引

| 表名 | 索引名稱 | 欄位 | 類型 | 用途 |
|------|---------|------|------|------|
| notifications | idx_notifications_status_priority | (status, priority) | B-tree | 佇列查詢 |
| notification_destinations | idx_notif_dest_pending | (status, created_at) | B-tree | 掃描待處理 |
| projects | idx_projects_notify_key | notify_key | Unique B-tree | 快速查找 |
| destinations | idx_destinations_targets_gin | targets | GIN | JSONB 查詢 |

### 4.4 JSONB 索引

```sql
-- destinations.targets GIN 索引（支援 JSONB 查詢）
CREATE INDEX idx_destinations_targets_gin ON destinations USING gin(targets);

-- 查詢範例
-- 查找包含特定 email 的目的地
SELECT * FROM destinations 
WHERE targets @> '[{"email": "user@example.com"}]';

-- 查找特定類型的目的地
SELECT * FROM destinations 
WHERE targets @> '[{"type": "channel"}]';
```

---

## 5. 檢視設計

### 5.1 notification_summary (通知摘要)

**用途**: 提供通知的彙總資訊，方便查詢

```sql
CREATE OR REPLACE VIEW notification_summary AS
SELECT 
    n.id AS notification_id,
    n.project_id,
    p.notify_key,
    p.company_id,
    c.name AS company_name,
    n.sender_id,
    u.name AS sender_name,
    n.message_type,
    n.priority,
    n.status,
    n.total_destinations,
    n.destinations_sent,
    n.destinations_failed,
    n.created_at,
    n.sent_at,
    EXTRACT(EPOCH FROM (n.sent_at - n.created_at)) AS processing_time_seconds
FROM notifications n
JOIN projects p ON n.project_id = p.id
JOIN companies c ON p.company_id = c.id
LEFT JOIN users u ON n.sender_id = u.id;

-- 註解
COMMENT ON VIEW notification_summary IS '通知摘要檢視，包含公司、專案、發送者資訊';
```

### 5.2 usage_summary (使用量摘要)

**用途**: 提供專案使用量彙總，方便計費查詢

```sql
CREATE OR REPLACE VIEW usage_summary AS
SELECT 
    ur.project_id,
    p.notify_key,
    p.company_id,
    c.name AS company_name,
    ur.date,
    ur.notification_count,
    SUM(ur.notification_count) OVER (
        PARTITION BY ur.project_id 
        ORDER BY ur.date 
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) AS cumulative_count
FROM usage_records ur
JOIN projects p ON ur.project_id = p.id
JOIN companies c ON p.company_id = c.id;

-- 註解
COMMENT ON VIEW usage_summary IS '使用量摘要檢視，包含累計數量';
```

### 5.3 bot_status_summary (Bot 狀態摘要)

**用途**: 提供 Bot 使用統計

```sql
CREATE OR REPLACE VIEW bot_status_summary AS
SELECT 
    tb.id AS bot_id,
    tb.name AS bot_name,
    tb.app_id,
    tb.status,
    COUNT(DISTINCT d.id) AS destination_count,
    COUNT(DISTINCT d.project_id) AS project_count,
    MAX(n.created_at) AS last_notification_at
FROM teams_bots tb
LEFT JOIN destinations d ON tb.id = d.bot_id
LEFT JOIN notification_destinations nd ON d.id = nd.destination_id
LEFT JOIN notifications n ON nd.notification_id = n.id
GROUP BY tb.id, tb.name, tb.app_id, tb.status;

-- 註解
COMMENT ON VIEW bot_status_summary IS 'Bot 狀態摘要，包含目的地數、專案數、最後通知時間';
```

---

## 6. 資料遷移策略

### 6.1 Migration Scripts

**命名規範**: `{version}_{description}.sql`

範例：
- `001_create_core_tables.sql`
- `002_create_notification_tables.sql`
- `003_create_billing_tables.sql`
- `004_create_indexes.sql`
- `005_create_views.sql`

### 6.2 初始資料

#### 預設計費方案

```sql
INSERT INTO billing_plans (name, description, price, notification_limit, features) VALUES
('Free', '免費方案', 0, 1000, '["基本通知", "最多 3 個專案"]'),
('Standard', '標準方案', 1999, 10000, '["基本通知", "最多 10 個專案", "優先支援"]'),
('Enterprise', '企業方案', 9999, 100000, '["所有功能", "無限專案", "專屬支援"]');
```

#### 預設管理員

```sql
-- 密碼需使用 bcrypt 加密
INSERT INTO companies (name, contact_email, status) VALUES
('System', 'admin@system.local', 'active');

INSERT INTO users (company_id, email, name, password_hash, role) VALUES
((SELECT id FROM companies WHERE name = 'System'), 
 'admin@system.local', 
 'System Admin', 
 '$2a$10$...', -- bcrypt hash
 'admin');
```

### 6.3 版本控制

使用 Schema Version 表追蹤遷移版本：

```sql
CREATE TABLE schema_version (
    version       INTEGER PRIMARY KEY,
    description   VARCHAR(255) NOT NULL,
    applied_at    TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## 7. 效能優化

### 7.1 分區策略

**大資料表按月分區**:
- `notifications` - 按 `created_at` 月份分區
- `notification_destinations` - 按 `created_at` 月份分區
- `audit_logs` - 按 `created_at` 月份分區

**優點**:
- 提升查詢效能（只掃描相關分區）
- 簡化資料清理（直接刪除舊分區）
- 降低索引大小

**自動分區建立**（建議使用 cron job）:
```sql
-- 建立下個月分區的函式
CREATE OR REPLACE FUNCTION create_next_month_partitions()
RETURNS void AS $$
DECLARE
    next_month DATE := date_trunc('month', CURRENT_DATE + interval '1 month');
    next_next_month DATE := next_month + interval '1 month';
BEGIN
    -- notifications
    EXECUTE format('CREATE TABLE IF NOT EXISTS notifications_%s PARTITION OF notifications FOR VALUES FROM (%L) TO (%L)',
        to_char(next_month, 'YYYY_MM'), next_month, next_next_month);
    
    -- notification_destinations
    EXECUTE format('CREATE TABLE IF NOT EXISTS notification_destinations_%s PARTITION OF notification_destinations FOR VALUES FROM (%L) TO (%L)',
        to_char(next_month, 'YYYY_MM'), next_month, next_next_month);
    
    -- audit_logs
    EXECUTE format('CREATE TABLE IF NOT EXISTS audit_logs_%s PARTITION OF audit_logs FOR VALUES FROM (%L) TO (%L)',
        to_char(next_month, 'YYYY_MM'), next_month, next_next_month);
END;
$$ LANGUAGE plpgsql;
```

### 7.2 查詢優化

#### 複合索引

針對常用查詢建立複合索引：

```sql
-- 查詢待處理通知（狀態 + 優先級 + 時間）
CREATE INDEX idx_notifications_processing ON notifications(status, priority, created_at) 
WHERE status IN ('pending', 'enqueued');

-- 查詢待重試通知
CREATE INDEX idx_notif_dest_retry_queue ON notification_destinations(status, next_retry_at) 
WHERE status = 'failed' AND retry_count < 10;
```

#### 部分索引

僅索引特定條件的資料：

```sql
-- 僅索引啟用的目的地
CREATE INDEX idx_destinations_active ON destinations(project_id) WHERE status = 'active';

-- 僅索引未完成的通知
CREATE INDEX idx_notifications_incomplete ON notifications(created_at) 
WHERE status NOT IN ('sent', 'cancelled');
```

### 7.3 連接池配置

**建議配置**:
```yaml
database:
  pool:
    min_size: 5
    max_size: 25
    max_idle_time: 300s
    max_lifetime: 1800s
```

### 7.4 Vacuum 與 Analyze

**自動 Vacuum**:
```sql
-- 啟用自動 vacuum
ALTER TABLE notifications SET (autovacuum_enabled = true);
ALTER TABLE notification_destinations SET (autovacuum_enabled = true);

-- 設定 Vacuum 參數
ALTER TABLE notifications SET (
    autovacuum_vacuum_scale_factor = 0.1,
    autovacuum_analyze_scale_factor = 0.05
);
```

---

## 附錄

### A. 資料庫備份策略

**建議備份方案**:
1. **完整備份**: 每日凌晨 2:00
2. **增量備份**: 每 6 小時
3. **WAL 歸檔**: 即時
4. **保留期限**: 30 天

**備份指令**:
```bash
# 完整備份
pg_dump -U teamsnotify -h localhost notification_center > backup_$(date +%Y%m%d).sql

# 僅資料
pg_dump -U teamsnotify -h localhost --data-only notification_center > data_$(date +%Y%m%d).sql

# 僅結構
pg_dump -U teamsnotify -h localhost --schema-only notification_center > schema.sql
```

### B. 資料保留政策

| 資料類型 | 保留期限 | 清理方式 |
|---------|---------|---------|
| notifications | 90 天 | 刪除舊分區 |
| notification_destinations | 90 天 | 刪除舊分區 |
| audit_logs | 180 天 | 刪除舊分區 |
| usage_records | 永久 | - |
| files | 依需求 | 手動清理 |

**清理範例**:
```sql
-- 刪除 3 個月前的 notifications 分區
DROP TABLE IF EXISTS notifications_2025_07;

-- 刪除 6 個月前的 audit_logs 分區
DROP TABLE IF EXISTS audit_logs_2025_04;
```

### C. 效能監控查詢

```sql
-- 查看表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 查看索引使用率
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan ASC;

-- 查看慢查詢 (需啟用 pg_stat_statements)
SELECT 
    query,
    calls,
    mean_exec_time,
    max_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

---

**文件結束**

如有任何問題，請聯絡：it-rfq@company.com

