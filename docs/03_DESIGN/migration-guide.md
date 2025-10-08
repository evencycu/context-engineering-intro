# 遷移指南

本指南幫助您從舊版本升級到最新版本的 Teams Notification API。

## 版本變更

### v1.0.0 → v1.1.0 (2025-10-02)

#### 主要變更
- **移除 `failed_notifications` 表**: 功能整合到 `notification_destinations`
- **新增 Actor 系統**: 非同步通知處理
- **Redis 整合**: 佇列管理和電路斷路器
- **新增 Provision API**: 一鍵建立專案和目的地
- **新增 External API**: 外部系統通知發送

#### 資料庫變更
- 新增 `notification_destinations` 欄位
- 移除 `failed_notifications` 相關表
- 更新索引和約束

## 遷移步驟

### 1. 備份現有資料
```bash
# 備份資料庫
pg_dump -h localhost -U teamsnotify -d notification_center > backup_$(date +%Y%m%d_%H%M%S).sql

# 備份配置檔案
cp config.example.env config.backup.env
```

### 2. 停止服務
```bash
# 停止現有服務
pkill -f "./server"

# 停止 Docker 容器
docker-compose down
```

### 3. 更新程式碼
```bash
# 拉取最新程式碼
git pull origin main

# 更新依賴
go mod tidy
```

### 4. 執行資料庫遷移
```bash
# 啟動資料庫
docker-compose up -d postgres redis

# 執行遷移腳本
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -f /scripts/migrations/009_enhance_notification_destinations_for_async_actor.sql
```

### 5. 更新環境變數
```bash
# 新增 Redis 相關環境變數
export REDIS_URL="redis://localhost:6379"

# 確保 Teams Bot 憑證正確
export TEAMS_BOT_APP_ID=844146d7-4ac9-4e4d-a463-d6e027714e81
export TEAMS_TENANT_ID=051cece0-e4dc-4aed-b471-bf29824e1ee6
export TEAMS_BOT_APP_PASSWORD='HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy'
```

### 6. 重新編譯和啟動
```bash
# 編譯新版本
go build -o server cmd/server/main.go

# 啟動服務
bash start_server.sh
```

### 7. 驗證遷移
```bash
# 檢查服務健康狀態
curl http://localhost:8080/health

# 檢查資料庫結構
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "\d notification_destinations"

# 測試新功能
curl -X POST http://localhost:8080/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "test-key",
    "message": "Migration test",
    "message_type": "text",
    "priority": "normal",
    "targets": ["all"]
  }'
```

## 資料遷移

### 從 failed_notifications 遷移
如果您的舊版本有 `failed_notifications` 資料，需要手動遷移：

```sql
-- 遷移失敗通知到 notification_destinations
INSERT INTO notification_destinations (
    id, notification_id, destination_id, conversation_id,
    bot_id, bot_type, status, error_message, retry_count,
    max_retries, next_retry_at, first_attempt_at, last_attempt_at,
    failure_reason, created_at, updated_at
)
SELECT 
    id, notification_id, destination_id, conversation_id,
    bot_id, bot_type, 'failed', error_message, retry_count,
    max_retries, next_retry_at, first_attempt_at, last_attempt_at,
    failure_reason, created_at, updated_at
FROM failed_notifications
WHERE status = 'failed';
```

### 清理舊資料
```sql
-- 刪除舊的 failed_notifications 表
DROP TABLE IF EXISTS failed_notifications;
DROP VIEW IF EXISTS failed_notifications_queue_status;
```

## 配置變更

### 新增配置項目
```yaml
# configs/api-server.yaml
redis:
  url: "redis://localhost:6379"
  max_retries: 3
  timeout: 5s

actor_pool:
  max_actors: 10
  queue_check_interval: 1s
  retry_delay: 5s

circuit_breaker:
  failure_threshold: 5
  recovery_timeout: 30s
  half_open_max_calls: 3
```

### 環境變數更新
```bash
# 新增環境變數
export REDIS_URL="redis://localhost:6379"
export ACTOR_POOL_SIZE="10"
export CIRCUIT_BREAKER_THRESHOLD="5"
```

## API 變更

### 新增端點
- `POST /api/v1/provision` - 建立專案和目的地
- `GET /api/v1/provision/{notify_key}` - 取得專案資訊
- `POST /api/v1/external/notify` - 外部系統通知發送

### 變更的端點
- `POST /api/v1/notifications` - 現在為非同步處理
- 所有通知相關端點都支援新的狀態管理

### 移除的端點
- 所有 `failed_notifications` 相關端點
- 舊的佇列管理端點

## 相容性

### 向後相容
- 現有的 API 端點保持不變
- 資料庫結構向後相容
- 配置檔案格式相容

### 不向後相容
- `failed_notifications` 表已移除
- 舊的佇列管理方式已棄用
- 某些內部 API 結構已變更

## 回滾計畫

如果遷移失敗，可以回滾到舊版本：

### 1. 停止新服務
```bash
pkill -f "./server"
```

### 2. 恢復資料庫
```bash
# 從備份恢復資料庫
psql -h localhost -U teamsnotify -d notification_center < backup_20251002_120000.sql
```

### 3. 恢復舊版本
```bash
# 切換到舊版本
git checkout v1.0.0

# 重新編譯
go build -o server cmd/server/main.go

# 啟動舊服務
./server
```

## 測試遷移

### 1. 功能測試
```bash
# 執行完整測試套件
bash scripts/comprehensive_api_test.sh

# 測試新功能
bash scripts/test_provision_api.sh
bash scripts/test_external_api.sh
```

### 2. 效能測試
```bash
# 測試並發處理
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/external/notify \
    -H "Content-Type: application/json" \
    -d '{"notify_key": "test", "message": "Test '$i'", "message_type": "text"}' &
done
wait
```

### 3. 資料完整性測試
```sql
-- 檢查資料遷移是否正確
SELECT 
    COUNT(*) as total_notifications,
    COUNT(CASE WHEN status = 'sent' THEN 1 END) as sent,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending
FROM notification_destinations;
```

## 常見問題

### Q: 遷移後服務無法啟動
A: 檢查 Redis 連線和環境變數設定

### Q: 通知發送失敗
A: 檢查 Bot 安裝狀態和 Teams 憑證

### Q: 效能變慢
A: 調整 Actor Pool 大小和 Redis 配置

### Q: 資料遺失
A: 從備份恢復資料庫

## 支援

如果遷移過程中遇到問題，請：

1. 檢查服務日誌
2. 查看資料庫狀態
3. 確認環境變數設定
4. 聯絡技術支援

### 診斷資訊收集
```bash
# 收集診斷資訊
bash scripts/collect_diagnostics.sh > migration_diagnostics.log
```
