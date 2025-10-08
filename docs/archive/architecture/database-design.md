# 資料庫設計

## 資料庫配置

### 當前資料庫配置

**資料庫名稱**: `notification_center`

```
postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable
```

### 為什麼是 notification_center？

這個專案使用 `notification_center` 作為主要資料庫名稱，統一管理所有通知相關的數據。

### Docker 容器中的資料庫列表

```bash
$ docker exec teamsnotify-postgres psql -U teamsnotify -l

        Name          |    Owner    
----------------------+-------------
 notification_center  | teamsnotify  ✅ 主要使用
```

## 配置方式

### 方法 1: 使用環境變數（推薦）

```bash
export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
```

### 方法 2: 使用 Docker Compose

```yaml
services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: notification_center
      POSTGRES_USER: teamsnotify
      POSTGRES_PASSWORD: teamsnotify123
    ports:
      - "5432:5432"
```

### 方法 3: 直接連線

```bash
psql -h localhost -U teamsnotify -d notification_center
```

## 實體關係圖 (ERD)

### 核心平台實體

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                COMPANIES                                        │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ name (VARCHAR(255), NOT NULL)                                                  │
│ contact_email (VARCHAR(255), NOT NULL)                                         │
│ contact_phone (VARCHAR(50))                                                    │
│ address (TEXT)                                                                  │
│ status (VARCHAR(20), CHECK: active|inactive|suspended)                         │
│ billing_enabled (BOOLEAN, DEFAULT true)                                        │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                  USERS                                          │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ company_id (UUID, FK → companies.id)                                           │
│ email (VARCHAR(255), NOT NULL)                                                 │
│ name (VARCHAR(255), NOT NULL)                                                  │
│ role (VARCHAR(50), CHECK: admin|user|viewer)                                   │
│ status (VARCHAR(20), CHECK: active|inactive|suspended)                         │
│ last_login_at (TIMESTAMP)                                                      │
│ api_key_hash (VARCHAR(255))                                                    │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                PROJECTS                                         │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ company_id (UUID, FK → companies.id)                                           │
│ notify_key (VARCHAR(255), UNIQUE, NOT NULL)                                    │
│ description (TEXT)                                                              │
│ status (VARCHAR(20), CHECK: active|inactive|suspended)                         │
│ daily_limit (INTEGER, DEFAULT 10000)                                           │
│ monthly_limit (INTEGER, DEFAULT 300000)                                        │
│ priority (VARCHAR(20), CHECK: low|normal|high)                                 │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              DESTINATIONS                                      │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ project_id (UUID, FK → projects.id)                                            │
│ name (VARCHAR(255), NOT NULL)                                                  │
│ description (TEXT)                                                              │
│ teams_tenant_id (VARCHAR(255), NOT NULL)                                       │
│ targets (JSONB, NOT NULL)                                                      │
│ bot_id (UUID, FK → teams_bots.id)                                              │
│ bot_type (VARCHAR(20), CHECK: platform|third_party)                            │
│ status (VARCHAR(20), CHECK: active|inactive|suspended)                         │
│ validation_status (VARCHAR(20), CHECK: pending|validated|failed)               │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Bot 管理實體

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                               TEAMS_BOTS                                        │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ type (VARCHAR(20), CHECK: platform|third_party)                                │
│ company_id (UUID, FK → companies.id) [僅 third_party]                          │
│ name (VARCHAR(255), NOT NULL)                                                  │
│ description (TEXT)                                                              │
│ app_id (VARCHAR(255), UNIQUE, NOT NULL)                                        │
│ app_password_hash (VARCHAR(255), NOT NULL)                                     │
│ tenant_id (VARCHAR(255))                                                       │
│ status (VARCHAR(20), CHECK: active|inactive|suspended)                         │
│ webhook_url (VARCHAR(500))                                                     │
│ capabilities (JSONB)                                                            │
│ rate_limit_per_minute (INTEGER)                                                │
│ max_concurrent_requests (INTEGER)                                              │
│ api_endpoint (VARCHAR(500)) [僅 third_party]                                   │
│ api_key_hash (VARCHAR(255)) [僅 third_party]                                   │
│ contact_email (VARCHAR(255)) [僅 third_party]                                  │
│ contact_phone (VARCHAR(50)) [僅 third_party]                                   │
│ created_by (UUID, FK → users.id) [僅 third_party]                              │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            BOT_INSTALLATIONS                                    │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ bot_id (UUID, FK → teams_bots.id)                                              │
│ bot_type (VARCHAR(20), CHECK: platform|third_party)                            │
│ teams_tenant_id (VARCHAR(255), NOT NULL)                                       │
│ conversation_type (VARCHAR(20), CHECK: personal|groupChat|channel)             │
│ conversation_id (VARCHAR(500), NOT NULL)                                       │
│ service_url (VARCHAR(500), NOT NULL)                                           │
│ recipient_id (VARCHAR(255), NOT NULL)                                          │
│ recipient_name (VARCHAR(255), NOT NULL)                                        │
│ from_id (VARCHAR(255), NOT NULL)                                               │
│ from_name (VARCHAR(255), NOT NULL)                                             │
│ from_aad_object_id (VARCHAR(255), NOT NULL)                                    │
│ installation_status (VARCHAR(20), CHECK: active|inactive|stale)                │
│ installed_at (TIMESTAMP)                                                       │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 通知管理實體

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              NOTIFICATIONS                                     │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ project_id (UUID, FK → projects.id)                                            │
│ message_type (VARCHAR(20), CHECK: text|file|adaptive_card)                     │
│ content (TEXT, NOT NULL)                                                       │
│ mentions (TEXT[])                                                               │
│ priority (VARCHAR(20), CHECK: low|normal|high)                                 │
│ status (VARCHAR(20), CHECK: pending|processing|sent|failed|cancelled)          │
│ metadata (JSONB)                                                                │
│ sent_at (TIMESTAMP)                                                             │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        NOTIFICATION_DESTINATIONS                               │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ notification_id (UUID, FK → notifications.id)                                  │
│ destination_id (UUID, FK → destinations.id)                                    │
│ conversation_id (VARCHAR(500))                                                 │
│ bot_id (UUID, FK → teams_bots.id)                                              │
│ bot_type (VARCHAR(20), CHECK: platform|third_party)                            │
│ status (VARCHAR(20), CHECK: pending|processing|sent|failed|cancelled)          │
│ error_message (TEXT)                                                            │
│ teams_message_id (VARCHAR(255))                                                │
│ sent_at (TIMESTAMP)                                                             │
│ retry_count (INTEGER, DEFAULT 0)                                               │
│ max_retries (INTEGER, DEFAULT 5)                                               │
│ next_retry_at (TIMESTAMP)                                                      │
│ first_attempt_at (TIMESTAMP)                                                   │
│ last_attempt_at (TIMESTAMP)                                                    │
│ failure_reason (VARCHAR(50))                                                   │
│ retry_after (INTEGER)                                                          │
│ actor_id (VARCHAR(100))                                                        │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## 主要業務表

### 1. 公司管理
- **companies**: 公司基本資料
- **users**: 用戶帳號和權限

### 2. 專案管理
- **projects**: 通知專案設定
- **destinations**: 通知目的地配置

### 3. Bot 管理
- **teams_bots**: Teams Bot 註冊資訊
- **bot_installations**: Bot 安裝記錄

### 4. 通知處理
- **notifications**: 通知記錄
- **notification_destinations**: 通知目的地處理狀態（取代舊的 failed_notifications）

## 索引設計

### 主要索引
```sql
-- 公司查詢優化
CREATE INDEX idx_companies_status ON companies(status);
CREATE INDEX idx_companies_name ON companies(name);

-- 用戶查詢優化
CREATE INDEX idx_users_company_id ON users(company_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_api_key_hash ON users(api_key_hash);

-- 專案查詢優化
CREATE INDEX idx_projects_company_id ON projects(company_id);
CREATE INDEX idx_projects_notify_key ON projects(notify_key);
CREATE INDEX idx_projects_status ON projects(status);

-- 目的地查詢優化
CREATE INDEX idx_destinations_project_id ON destinations(project_id);
CREATE INDEX idx_destinations_teams_tenant_id ON destinations(teams_tenant_id);

-- Bot 查詢優化
CREATE INDEX idx_teams_bots_app_id ON teams_bots(app_id);
CREATE INDEX idx_teams_bots_type ON teams_bots(type);

-- Bot 安裝查詢優化
CREATE INDEX idx_bot_installations_bot_id ON bot_installations(bot_id);
CREATE INDEX idx_bot_installations_conversation_id ON bot_installations(conversation_id);
CREATE INDEX idx_bot_installations_installation_status ON bot_installations(installation_status);

-- 通知查詢優化
CREATE INDEX idx_notifications_project_id ON notifications(project_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- 通知目的地查詢優化
CREATE INDEX idx_notification_destinations_notification_id ON notification_destinations(notification_id);
CREATE INDEX idx_notification_destinations_status ON notification_destinations(status);
CREATE INDEX idx_notification_destinations_next_retry_at ON notification_destinations(next_retry_at);
CREATE INDEX idx_notification_destinations_conversation_id ON notification_destinations(conversation_id);
```

## 檢查資料庫表

### 查看所有表
```sql
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
ORDER BY table_name;
```

### 查看表結構
```sql
\d+ companies
\d+ users
\d+ projects
\d+ destinations
\d+ teams_bots
\d+ bot_installations
\d+ notifications
\d+ notification_destinations
```

### 常用查詢

#### 1. 檢查通知佇列狀態
```sql
SELECT 
    status,
    COUNT(*) as count,
    MIN(created_at) as oldest,
    MAX(created_at) as newest
FROM notification_destinations 
GROUP BY status
ORDER BY count DESC;
```

#### 2. 檢查重試佇列
```sql
SELECT 
    COUNT(*) as retry_count,
    MIN(next_retry_at) as next_retry,
    MAX(retry_count) as max_retries
FROM notification_destinations 
WHERE status = 'failed' 
AND next_retry_at IS NOT NULL
AND next_retry_at <= NOW();
```

#### 3. 檢查 Bot 安裝狀態
```sql
SELECT 
    installation_status,
    COUNT(*) as count
FROM bot_installations 
GROUP BY installation_status;
```

#### 4. 檢查專案使用情況
```sql
SELECT 
    p.notify_key,
    p.description,
    COUNT(n.id) as notification_count,
    COUNT(nd.id) as destination_count
FROM projects p
LEFT JOIN notifications n ON p.id = n.project_id
LEFT JOIN notification_destinations nd ON n.id = nd.notification_id
GROUP BY p.id, p.notify_key, p.description
ORDER BY notification_count DESC;
```

## 資料庫維護

### 清理舊資料
```sql
-- 清理 30 天前的已發送通知
DELETE FROM notification_destinations 
WHERE status = 'sent' 
AND sent_at < NOW() - INTERVAL '30 days';

-- 清理 90 天前的失敗通知
DELETE FROM notification_destinations 
WHERE status = 'failed' 
AND last_attempt_at < NOW() - INTERVAL '90 days';
```

### 統計資訊更新
```sql
-- 更新表統計資訊
ANALYZE companies;
ANALYZE users;
ANALYZE projects;
ANALYZE destinations;
ANALYZE teams_bots;
ANALYZE bot_installations;
ANALYZE notifications;
ANALYZE notification_destinations;
```

### 備份建議
```bash
# 完整備份
pg_dump -h localhost -U teamsnotify -d notification_center > backup_$(date +%Y%m%d_%H%M%S).sql

# 僅結構備份
pg_dump -h localhost -U teamsnotify -d notification_center --schema-only > schema_backup.sql

# 僅資料備份
pg_dump -h localhost -U teamsnotify -d notification_center --data-only > data_backup.sql
```
