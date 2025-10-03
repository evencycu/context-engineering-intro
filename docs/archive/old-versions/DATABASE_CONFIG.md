# 資料庫配置說明

## 📊 當前資料庫配置

### 實際使用的資料庫

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

## 🔧 配置方式

### 方法 1: 使用環境變數（推薦）

```bash
export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
./server
```

### 方法 2: 使用啟動腳本（最簡單）

```bash
./start_server.sh
```

腳本會自動：
- ✅ 檢查資料庫連接
- ✅ 檢查並執行必要的 migration
- ✅ 設置環境變數
- ✅ 啟動服務

### 方法 3: 修改預設值（已完成）

在 `cmd/server/main.go` 中：

```go
dbURL := getEnv("DATABASE_URL", "postgres://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable")
```

## 📋 資料庫結構

### 主要業務表

```
notification_center 資料庫包含：

核心業務表：
- companies
- users
- projects
- teams_bots
- bot_installations
- destinations
- notifications
- notification_destinations

🆕 Queue 相關（Redis + Actor）：
- notification_destinations（含 retry_count / next_retry_at / failure_reason / retry_after / actor_id）

其他功能表：
- audit_logs
- company_billing
- billing_plans
- usage_records
- system_logs
- system_settings
- feature_flags
- bot_health_status
- bot_routing_rules
```

### 檢查資料庫表

```bash
# 列出所有表
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "\dt"

# 檢查異步欄位是否存在（notification_destinations）
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "
SELECT 
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='notification_destinations' AND column_name='next_retry_at') AS has_next_retry_at,
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='notification_destinations' AND column_name='failure_reason') AS has_failure_reason,
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='notification_destinations' AND column_name='retry_after') AS has_retry_after;"
```

## 🗄️ Migration 管理

### 執行新的 Migration

```bash
# 異步欄位 migration（Redis + Actor）
docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < scripts/migrations/009_enhance_notification_destinations_for_async_actor.sql
```

### 檢查 Migration 狀態

```bash
# 檢查 notification_destinations 異步欄位
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "
SELECT 
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='notification_destinations' AND column_name='next_retry_at') AS has_next_retry_at,
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='notification_destinations' AND column_name='failure_reason') AS has_failure_reason,
  EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='notification_destinations' AND column_name='retry_after') AS has_retry_after;"
```

### 所有 Migrations

```
scripts/migrations/
├── 001_enhance_bot_installations.sql
├── 002_drop_user_foreign_keys.sql
├── 003_make_sender_id_nullable.sql
├── 004_update_bot_installations_schema.sql
├── 005_installations_unique_on_conversation.sql
├── 006_unify_bots_to_teams_bots.sql
├── 007_projects_rename_key_name_to_notify_key.sql
└── 009_enhance_notification_destinations_for_async_actor.sql  🆕
```

## 🔍 常用查詢

### 查看隊列統計（以 notification_destinations 為準）

```sql
SELECT 
    COALESCE(failure_reason, 'unknown') AS reason,
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE status IN ('pending','processing') AND retry_count < max_retries) as pending,
    COUNT(*) FILTER (WHERE status='failed' OR retry_count >= max_retries) as exhausted
FROM notification_destinations
GROUP BY COALESCE(failure_reason, 'unknown');
```

### 查看最近的失敗通知（DLQ）

```sql
SELECT 
    id, destination_id, conversation_id,
    failure_reason as reason,
    retry_count,
    next_retry_at,
    error_message,
    updated_at as created_at
FROM notification_destinations
WHERE status IN ('pending','failed')
ORDER BY COALESCE(next_retry_at, updated_at) DESC
LIMIT 10;
```

### 清理測試數據（僅開發環境）

```sql
-- 重置已達上限的重試計數
UPDATE notification_destinations SET retry_count=0 WHERE retry_count >= max_retries;

-- 刪除 30 天前的最終失敗紀錄
DELETE FROM notification_destinations WHERE status='failed' AND last_attempt_at < NOW() - INTERVAL '30 days';
```

## 🔄 資料庫連接配置

### PostgreSQL 連接參數

| 參數 | 值 | 說明 |
|------|-----|------|
| Host | localhost | 本地開發 |
| Port | 5432 | PostgreSQL 預設 port |
| Database | notification_center | 資料庫名稱 |
| User | teamsnotify | 資料庫用戶 |
| Password | teamsnotify123 | 密碼 |
| SSL Mode | disable | 本地開發不使用 SSL |

### 完整連接字串格式

```
postgresql://[user]:[password]@[host]:[port]/[database]?sslmode=[mode]
```

實際範例：
```
postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable
```

## 🐛 故障排查

### 問題：無法連接資料庫

```bash
# 1. 檢查 Docker 容器狀態
docker ps | grep postgres

# 2. 測試連接
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT 1"

# 3. 查看容器日誌
docker logs teamsnotify-postgres

# 4. 重啟容器
docker restart teamsnotify-postgres
```

### 問題：表不存在

```bash
# 檢查表是否存在
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "\dt"

# 執行缺失的 migration（009）
docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < scripts/migrations/009_enhance_notification_destinations_for_async_actor.sql
```

### 問題：Migration 已執行但表結構不對

```bash
# 刪除舊制式表 / 視圖（如仍存在）
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "DROP VIEW IF EXISTS failed_notifications_queue_status CASCADE;"
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "DROP TABLE IF EXISTS failed_notifications CASCADE;"

# 重新執行 migration（009）
docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < scripts/migrations/009_enhance_notification_destinations_for_async_actor.sql
```

## 📝 最佳實踐

### 開發環境

1. **使用環境變數**: 不要硬編碼資料庫連接
2. **定期備份**: 在執行 migration 前備份
3. **測試 Migration**: 先在測試資料庫執行

### 生產環境

1. **啟用 SSL**: `sslmode=require` 或 `sslmode=verify-full`
2. **使用連接池**: 配置適當的連接池大小
3. **監控連接數**: 避免連接耗盡
4. **定期備份**: 自動化備份策略

## 🔐 安全注意事項

1. ⚠️ **不要提交密碼**: 使用環境變數
2. ⚠️ **使用強密碼**: 生產環境改用強密碼
3. ⚠️ **限制連接**: 設置適當的防火牆規則
4. ⚠️ **啟用 SSL**: 生產環境必須使用 SSL

## 📚 相關文檔

- [Queue 快速開始](./QUICK_START_QUEUE.md)
- [Queue 完整文檔](./QUEUE_CIRCUIT_BREAKER.md)
- [主 README](../README.md)

---

**最後更新**: 2025-10-01  
**資料庫版本**: PostgreSQL 15

