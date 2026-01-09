# Queue API 使用範例

## 📌 基礎 API 測試

### 1. 健康檢查

```bash
curl -X GET http://localhost:8080/health | jq '.'
```

**預期回應**:
```json
{
  "status": "healthy",
  "timestamp": "2025-10-01T10:30:00Z",
  "version": "1.0.0"
}
```

---

## 🔍 Queue 狀態查詢

### 2. 查詢隊列狀態

```bash
curl -X GET http://localhost:8080/api/v1/queue/status | jq '.'
```

**預期回應**:
```json
{
  "total_pending": 0,
  "total_retrying": 0,
  "total_failed": 0,
  "next_retry_due": null,
  "oldest_pending": null,
  "circuit_state": "closed",
  "last_updated": "2025-10-01T10:30:00Z"
}
```

**欄位說明**:
- `total_pending`: 等待重試的通知數量
- `total_retrying`: 正在重試的通知數量
- `total_failed`: 已耗盡重試次數的通知數量
- `next_retry_due`: 下一個要重試的通知時間
- `circuit_state`: 熔斷器狀態 (closed/open/half-open)

---

## ⚡ Circuit Breaker 管理

### 3. 查詢熔斷器指標

```bash
curl -X GET http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '.'
```

**預期回應**:
```json
{
  "state": "closed",
  "total_requests": 1523,
  "success_requests": 1498,
  "failed_requests": 25,
  "failures": 0,
  "last_state_change": "2025-10-01T09:00:00Z"
}
```

**欄位說明**:
- `state`: 當前狀態
  - `closed`: 正常運作
  - `open`: 熔斷中（拒絕請求）
  - `half-open`: 測試恢復中
- `total_requests`: 總請求數
- `success_requests`: 成功請求數
- `failed_requests`: 失敗請求數
- `failures`: 連續失敗次數
- `last_state_change`: 最後狀態變更時間

### 4. 重置熔斷器

當確認服務已恢復，手動重置熔斷器：

```bash
curl -X POST http://localhost:8080/api/v1/queue/circuit-breaker/reset | jq '.'
```

**預期回應**:
```json
{
  "message": "Circuit breaker reset successfully"
}
```

**使用場景**:
- Teams API 服務已恢復正常
- 維護作業完成後
- 測試熔斷器機制

---

## 📨 發送通知測試

### 5. 發送正常通知

```bash
curl -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "your-notify-key",
    "message": "Test notification from Queue API - '$(date +%Y%m%d-%H%M%S)'",
    "targets": ["all"]
  }' | jq '.'
```

**預期回應**:
```json
{
  "notification_id": "550e8400-e29b-41d4-a716-446655440000",
  "project_name": "Test Project",
  "destinations_count": 4,
  "success_count": 4,
  "failed_count": 0,
  "results": [...]
}
```

### 6. 檢查是否有失敗進入隊列

```bash
# 發送通知後等待 2 秒
sleep 2

# 檢查隊列狀態
curl -X GET http://localhost:8080/api/v1/queue/status | jq '{
  pending: .total_pending,
  retrying: .total_retrying,
  failed: .total_failed,
  circuit: .circuit_state
}'
```

---

## 🔄 監控與持續觀察

### 7. 持續監控隊列（使用 watch）

```bash
# 每 2 秒更新一次隊列狀態
watch -n 2 'curl -s http://localhost:8080/api/v1/queue/status | jq "{pending: .total_pending, retrying: .total_retrying, failed: .total_failed, circuit: .circuit_state}"'
```

**輸出範例**:
```json
{
  "pending": 5,
  "retrying": 2,
  "failed": 0,
  "circuit": "closed"
}
```

按 `Ctrl+C` 停止監控。

### 8. 同時監控隊列和熔斷器

```bash
# 建立監控腳本
cat << 'EOF' > monitor_queue.sh
#!/bin/bash
while true; do
  clear
  echo "=== Queue Status ==="
  curl -s http://localhost:8080/api/v1/queue/status | jq '{
    pending: .total_pending,
    retrying: .total_retrying,
    failed: .total_failed,
    next_retry: .next_retry_due,
    circuit: .circuit_state
  }'
  
  echo ""
  echo "=== Circuit Breaker Metrics ==="
  curl -s http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '{
    state: .state,
    total: .total_requests,
    success: .success_requests,
    failed: .failed_requests,
    consecutive_failures: .failures
  }'
  
  sleep 3
done
EOF

chmod +x monitor_queue.sh
./monitor_queue.sh
```

---

## 🧪 完整測試流程

### 測試 1: 正常流程

```bash
#!/bin/bash

echo "=== Test 1: Normal Flow ==="

# 1. 檢查初始狀態
echo "1. Initial queue status:"
curl -s http://localhost:8080/api/v1/queue/status | jq '.total_pending'

# 2. 發送通知
echo "2. Sending notification..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "your-notify-key",
    "message": "Test notification",
    "targets": ["all"]
  }')

echo "$RESPONSE" | jq '{notification_id, success_count, failed_count}'

# 3. 檢查是否成功（隊列應該為空）
sleep 2
echo "3. Queue after success:"
curl -s http://localhost:8080/api/v1/queue/status | jq '{pending: .total_pending, circuit: .circuit_state}'
```

### 測試 2: 失敗重試流程

```bash
#!/bin/bash

echo "=== Test 2: Failure & Retry Flow ==="

# 1. 發送會失敗的通知（使用無效 target）
echo "1. Sending notification (expected to fail)..."
curl -s -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "your-notify-key",
    "message": "Test failed notification",
    "targets": ["nonexistent@example.com"]
  }' | jq '{notification_id, failed_count}'

# 2. 等待失敗通知進入隊列
sleep 3

# 3. 檢查隊列（應該有待重試的通知）
echo "2. Queue after failure:"
curl -s http://localhost:8080/api/v1/queue/status | jq '.'

# 4. 監控重試過程（持續 30 秒）
echo "3. Monitoring retry process for 30 seconds..."
for i in {1..10}; do
  echo "[$i] $(date +%H:%M:%S)"
  curl -s http://localhost:8080/api/v1/queue/status | jq '{pending: .total_pending, retrying: .total_retrying, circuit: .circuit_state}'
  sleep 3
done
```

### 測試 3: 熔斷器測試

```bash
#!/bin/bash

echo "=== Test 3: Circuit Breaker Test ==="

# 1. 檢查初始熔斷器狀態
echo "1. Initial circuit breaker state:"
curl -s http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '{state, failures}'

# 2. 模擬多次失敗（根據實際情況調整）
echo "2. Simulating failures..."
for i in {1..6}; do
  curl -s -X POST http://localhost:8080/api/v1/notify \
    -H "Content-Type: application/json" \
    -d "{
      \"notify_key\": \"your-notify-key\",
      \"message\": \"Fail test $i\",
      \"targets\": [\"invalid-target-$i@example.com\"]
    }" > /dev/null
  sleep 1
done

# 3. 檢查熔斷器是否開啟
sleep 5
echo "3. Circuit breaker after failures:"
curl -s http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '{state, failures, failed_requests}'

# 4. 重置熔斷器
echo "4. Resetting circuit breaker..."
curl -s -X POST http://localhost:8080/api/v1/queue/circuit-breaker/reset | jq '.'

# 5. 確認已重置
echo "5. Circuit breaker after reset:"
curl -s http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '{state, failures}'
```

---

## 📊 資料庫查詢

### 查看失敗通知詳情

```sql
-- 連接到資料庫
docker exec -it teamsnotify-postgres psql -U teamsnotify -d notification_center

-- 查看所有待重試通知
SELECT 
    id,
    target_id,
    reason,
    retry_count,
    max_retries,
    next_retry_at,
    created_at
FROM notification_destinations
WHERE retry_count < max_retries
ORDER BY next_retry_at;

-- 統計失敗原因
SELECT 
    reason,
    COUNT(*) as count,
    AVG(retry_count) as avg_retries
FROM notification_destinations
GROUP BY reason;

-- 查看隊列狀態視圖
-- 取代視圖的等價查詢：
SELECT 
  COUNT(*) FILTER (WHERE status IN ('pending','processing') AND retry_count < max_retries) AS pending,
  COUNT(*) FILTER (WHERE status='failed' OR retry_count >= max_retries) AS exhausted,
  MIN(next_retry_at) AS next_ready_at
FROM notification_destinations;
```

---

## 🎯 常用組合命令

### 快速診斷

```bash
#!/bin/bash
# 一鍵查看系統狀態

echo "=== Server Health ==="
curl -s http://localhost:8080/health | jq '{status, timestamp}'

echo ""
echo "=== Queue Status ==="
curl -s http://localhost:8080/api/v1/queue/status | jq '{pending: .total_pending, retrying: .total_retrying, failed: .total_failed, circuit: .circuit_state}'

echo ""
echo "=== Circuit Breaker ==="
curl -s http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '{state, success_rate: ((.success_requests / .total_requests * 100) | floor)}'

echo ""
echo "=== Database Queue Count ==="
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -t -c "SELECT COUNT(*) FROM notification_destinations WHERE status IN ('pending','processing') AND retry_count < max_retries;"
```

保存為 `quick_status.sh` 並執行：

```bash
chmod +x quick_status.sh
./quick_status.sh
```

---

## 🔧 故障排查命令

### 問題：隊列積壓

```bash
# 1. 檢查積壓數量
curl -s http://localhost:8080/api/v1/queue/status | jq '.total_pending'

# 2. 檢查熔斷器是否開啟
curl -s http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '.state'

# 3. 如果熔斷器開啟，重置它
curl -X POST http://localhost:8080/api/v1/queue/circuit-breaker/reset

# 4. 查看資料庫中最舊的待重試通知
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "
SELECT target_id, reason, retry_count, next_retry_at 
FROM notification_destinations 
WHERE retry_count < max_retries 
ORDER BY created_at 
LIMIT 5;"
```

### 問題：重試不生效

```bash
# 1. 檢查 server 日誌
tail -f server.log | grep -i "worker\|retry\|queue"

# 2. 手動觸發重試（更新 next_retry_at）
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "
UPDATE notification_destinations 
SET next_retry_at = NOW() 
WHERE retry_count < max_retries;"

# 3. 等待並檢查隊列變化
watch -n 2 'curl -s http://localhost:8080/api/v1/queue/status | jq .total_pending'
```

---

## 📝 注意事項

1. **Replace Placeholders**: 將 `your-notify-key` 替換為實際的 notify_key
2. **Test Environment**: 建議先在測試環境執行
3. **Monitor Resources**: 大量測試前確認系統資源充足
4. **Database Backup**: 重要操作前備份資料庫

---

**相關文檔**:
- [完整 Queue 文檔](./QUEUE_CIRCUIT_BREAKER.md)
- [快速開始指南](./QUICK_START_QUEUE.md)

