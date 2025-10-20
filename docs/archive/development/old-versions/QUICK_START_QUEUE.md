# Queue 和 Circuit Breaker 快速開始指南

## 🚀 快速開始

### 1. 執行資料庫 Migration

```bash
# 套用 009 異步欄位 migration（notification_destinations）
docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < scripts/migrations/009_enhance_notification_destinations_for_async_actor.sql
```

### 2. 編譯並啟動服務

```bash
# 編譯
go build -o server cmd/server/main.go

# 啟動 (記得帶上環境變數)
DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable" \
TEAMS_BOT_APP_ID="your-app-id" \
TEAMS_TENANT_ID="your-tenant-id" \
TEAMS_BOT_APP_PASSWORD="your-password" \
./server
```

### 3. 測試 Queue API

```bash
# 執行測試腳本
./scripts/test_queue.sh
```

## 📊 監控隊列

### 查看隊列狀態

```bash
curl -X GET http://localhost:8080/api/v1/queue/status | jq '.'
```

**回應範例**:
```json
{
  "total_pending": 0,
  "total_retrying": 0,
  "total_failed": 0,
  "circuit_state": "closed",
  "last_updated": "2025-10-01T10:30:00Z"
}
```

### 查看熔斷器指標

```bash
curl -X GET http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '.'
```

**回應範例**:
```json
{
  "state": "closed",
  "total_requests": 100,
  "success_requests": 98,
  "failed_requests": 2,
  "failures": 0,
  "last_state_change": "2025-10-01T09:00:00Z"
}
```

## 🧪 模擬測試場景

### 場景 1: 正常通知（應該成功）

```bash
# 發送通知（使用有效的 notify_key）
curl -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "your-notify-key",
    "message": "Test notification",
    "targets": ["user@example.com"]
  }'

# 檢查隊列（應該為空）
curl http://localhost:8080/api/v1/queue/status | jq '.total_pending'
```

### 場景 2: 模擬失敗（會進入隊列）

```bash
# 使用不存在的 target 發送（預期失敗）
curl -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "your-notify-key",
    "message": "Test failed notification",
    "targets": ["nonexistent@example.com"]
  }'

# 等待幾秒後檢查隊列
sleep 3
curl http://localhost:8080/api/v1/queue/status | jq '.'
```

### 場景 3: 監控重試過程

```bash
# 持續監控隊列狀態
watch -n 2 'curl -s http://localhost:8080/api/v1/queue/status | jq "{pending: .total_pending, retrying: .total_retrying, failed: .total_failed, circuit: .circuit_state}"'
```

## 🔧 常見操作

### 重置熔斷器

當確認 Teams 服務已恢復時：

```bash
curl -X POST http://localhost:8080/api/v1/queue/circuit-breaker/reset
```

### 查詢資料庫中的失敗通知

```sql
-- 查看所有待重試的通知
SELECT 
    id,
    target_id,
    reason,
    retry_count,
    next_retry_at,
    error_message
FROM notification_destinations
WHERE retry_count < max_retries
ORDER BY next_retry_at;

-- 統計失敗原因
SELECT 
    reason,
    COUNT(*) as count
FROM notification_destinations
GROUP BY reason;
```

### 手動重試特定通知

```sql
-- 重置重試計數，立即重試
UPDATE notification_destinations
SET retry_count = 0,
    next_retry_at = NOW()
WHERE id = 'notification-uuid';
```

## 📈 效能調優

### 調整 Worker 數量

在 `cmd/server/main.go` 中修改：

```go
queueConfig := &queue.Config{
    WorkerPool:   5,  // 從 3 增加到 5
    // ...
}
```

### 調整輪詢頻率

```go
queueConfig := &queue.Config{
    PollInterval: 5 * time.Second,  // 從 10 秒縮短到 5 秒
    // ...
}
```

### 調整批次大小

```go
queueConfig := &queue.Config{
    BatchSize: 20,  // 從 10 增加到 20
    // ...
}
```

## ⚠️ 故障排查

### 問題：隊列一直積壓

**檢查**:
1. Circuit Breaker 是否開啟
2. Worker 是否正常運行
3. Teams API 是否正常

```bash
# 1. 檢查熔斷器
curl http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '.state'

# 2. 檢查日誌
tail -f server.log | grep "Worker"

# 3. 如果 CB 開啟，重置它
curl -X POST http://localhost:8080/api/v1/queue/circuit-breaker/reset
```

### 問題：通知沒有進入隊列

**可能原因**:
- 錯誤類型不可重試（如 4xx 錯誤）
- 資料庫連接問題

**檢查**:
```bash
# 查看 server 日誌中的錯誤
tail -f server.log | grep -i "enqueue\|failed"

# 確認資料庫連接
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT COUNT(*) FROM notification_destinations WHERE status IN ('pending','failed');"
```

## 📚 更多資訊

詳細文檔請參考: [QUEUE_CIRCUIT_BREAKER.md](./QUEUE_CIRCUIT_BREAKER.md)

## 🎯 下一步

- [ ] 整合到 Broadcast Service（自動入隊）
- [ ] 實現通知優先級
- [ ] 添加 Prometheus metrics
- [ ] 實現 Dead Letter Queue
- [ ] 建立告警規則

