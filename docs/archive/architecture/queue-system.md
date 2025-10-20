# 通知隊列與熔斷器機制

## Two-loop Enqueue Design (Updated)

This system now separates write and enqueue into two independent loops to improve reliability and simplicity.

### 1) Producer Loop (API write only)
- On request, the API writes `notification_destinations` rows with `status = pending`.
- No Redis operations are performed in the request path.
- Benefits: lower latency, zero coupling to Redis availability during writes.

### 2) EnqueueWorker (background)
- Periodically scans pending `notification_destinations`.
- For each row, pushes the `notification_destination.id` into the Redis List by its `priority`:
  - `queue:notifications:high`
  - `queue:notifications:normal`
  - `queue:notifications:low`
- On successful enqueue, updates the row to `status = enqueued`.
- Enqueue has no retry loop by design; failures keep the row as `pending` to be picked up in the next scan.

### Consumer Side (unchanged)
- `QueueConsumer` performs `BLPOP` across high → normal → low and hands tasks to `ActorPool`.
- `ActorPool` spawns `NotificationActor` up to `maxActors`. When capacity is full, the consumer waits and does not dequeue.
- Redis SETNX `processing` keys ensure unique processing across multiple nodes.

### Configuration knobs
- Enqueue scan interval: `internal/actor/notification_processor.go` → `NewEnqueueWorker(..., 200*time.Millisecond)`
- Enqueue batch size: `internal/actor/enqueue_worker.go` → `GetPending(ctx, 100)`
- Actor pool size: `internal/actor/actor_pool.go` → `maxActors` parameter

### Queue visibility (ops)
- Queue lengths:
  - High: `redis-cli LLEN queue:notifications:high`
  - Normal: `redis-cli LLEN queue:notifications:normal`
  - Low: `redis-cli LLEN queue:notifications:low`
- Keys overview: `redis-cli KEYS "queue:notifications:*"`

### State transitions
```
API → write ND(status=pending)
EnqueueWorker → enqueue → ND(status=enqueued)
QueueConsumer/Actor → processing → sent/failed per delivery outcome
```

Notes:
- The previous scanner that retried notifications has been removed. Enqueue errors keep rows in `pending` for the next worker sweep.
- Circuit breaker remains in place for delivery (Teams 429/Retry-After handling).

## 📋 目錄

- [概述](#概述)
- [架構設計](#架構設計)
- [核心組件](#核心組件)
- [工作流程](#工作流程)
- [資料庫設計](#資料庫設計)
- [API 端點](#api-端點)
- [使用範例](#使用範例)
- [配置說明](#配置說明)
- [監控與維護](#監控與維護)

---

## 概述

### 背景

當系統向 Microsoft Teams 發送通知時，可能會遇到以下問題：

1. **速率限制 (Rate Limit)**: Teams API 返回 `429 Too Many Requests` 錯誤
2. **網路問題**: 暫時性的網路故障或超時
3. **服務錯誤**: Teams 服務暫時不可用 (5xx 錯誤)

為了解決這些問題，我們實現了一個**失敗重試隊列**和**熔斷器 (Circuit Breaker)** 機制。

### 目標

- ✅ 自動重試失敗的通知
- ✅ 使用指數退避策略避免過度重試
- ✅ 遵守 Teams API 的 `Retry-After` header
- ✅ 熔斷器保護，防止服務雪崩
- ✅ 持久化失敗記錄，防止數據丟失
- ✅ 提供監控和管理 API

---

## 架構設計

```
┌─────────────────────────────────────────────────────────────┐
│                     Notification Request                     │
└───────────────────────────────┬─────────────────────────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │  Notification Service │
                    │  (Create & Store)     │
                    └───────────┬───────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │  notification_destinations │
                    │  (status: pending)    │
                    └───────────┬───────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │    Actor Pool V2      │
                    │  (Background Workers) │
                    └───────────┬───────────┘
                                │
                    ┌───────────┴───────────┐
                    │                       │
                    ▼                       ▼
            ┌───────────────┐   ┌──────────────────┐
            │ Circuit       │   │ Teams Sender     │
            │ Breaker       │   │ (Send to Teams)  │
            │ Check         │   │                  │
            └───────┬───────┘   └────────┬─────────┘
                    │                    │
                    ▼                    ▼
            ┌─────────────────────────────────┐
            │     Send to Teams API          │
            └─────────────┬───────────────────┘
                          │
            ┌─────────────┴───────────┐
            │                         │
        ✅ Success              ❌ Failure
            │                         │
            ▼                         ▼
    ┌───────────────┐       ┌──────────────────┐
    │ Update Status │       │ Update Retry     │
    │ (status: sent)│       │ (retry_count++)  │
    └───────────────┘       └──────────────────┘
```

### 核心設計原則

1. **完全非同步**: 所有通知發送都是非同步的
2. **Actor 模式**: 使用 Actor Pool 處理通知發送
3. **Redis 佇列**: 使用 Redis Sorted Set 管理重試佇列
4. **電路斷路器**: 所有 Actor 共享電路斷路器狀態
5. **統一資料表**: 使用 `notification_destinations` 管理所有狀態

---

## 核心組件

### 1. Redis 優先序佇列（Producer/Consumer）

本系統以 Redis List + SETNX 實作三層優先序佇列，並採 Producer/Consumer 架構實作：

- 佇列鍵名（每個優先序一個 List）：
  - `queue:notifications:high`
  - `queue:notifications:normal`
  - `queue:notifications:low`
- 去重鍵（避免重複入列，短期有效）：
  - `queue:notifications:dedup:{notification_destination_id}`（SETNX，TTL 約 60s）
- 處理中鍵（避免併發重複處理）：
  - `queue:notifications:processing:{notification_destination_id}`（SETNX，TTL 約 30m）

Producer（在 `NotificationService.SendNotification` 持久化 `notification_destinations` 後）會依 `priority=high|normal|low` 將每筆 `notification_destination.id` 推入對應的 List；
Consumer（`QueueConsumer`）會以 `BLPOP` 順序檢查 `high -> normal -> low`，取出一筆後，以 SETNX 設置 processing key，成功後交由 `ActorPool.SpawnActor` 建立 `NotificationActor` 進行發送流程。

此設計具備：
- 優先序保證（高優先佇列先被消化）
- 去重與併發防重（dedup + processing）
- 橫向擴展（多個 Consumer 節點可安全競爭 `BLPOP`）

### 2. Circuit Breaker (熔斷器)

**文件**: `internal/actor/redis_circuit_breaker.go`

#### 三種狀態

```
┌─────────┐  Max Failures   ┌─────────┐  Timeout     ┌──────────┐
│ Closed  │ ───────────────►│  Open   │ ──────────►  │ Half-Open │
│(正常)   │                 │(熔斷)   │              │(測試)     │
└────▲────┘                 └─────────┘              └─────┬─────┘
     │                                                      │
     └──────────────── Success Requests ───────────────────┘
```

1. **Closed (關閉狀態)**: 正常運作，允許所有請求
   - 當失敗次數達到閾值 → 切換到 **Open**

2. **Open (開啟狀態)**: 熔斷狀態，拒絕所有請求
   - 避免持續向不可用的服務發送請求
   - 等待一段時間後 → 切換到 **Half-Open**

3. **Half-Open (半開狀態)**: 測試狀態，允許少量請求測試服務恢復
   - 如果測試請求成功 → 切換回 **Closed**
   - 如果測試請求失敗 → 切換回 **Open**

#### 配置參數

```go
type CircuitBreaker struct {
    maxFailures     uint32        // 觸發熔斷的失敗次數閾值 (預設: 5)
    timeout         time.Duration // 熔斷後等待時間 (預設: 1 分鐘)
    halfOpenMaxReqs uint32        // 半開狀態允許的測試請求數 (預設: 3)
}
```

#### 使用方式

```go
// 執行受保護的操作
err := circuitBreaker.Execute(ctx, func() error {
    return sendToTeams(message)
})

if err == ErrCircuitOpen {
    // 熔斷器開啟，服務不可用
    log.Println("Circuit breaker is open, service unavailable")
}
```

---

### 3. Retry Policy (重試策略)

**文件**: `internal/actor/notification_actor.go` (內建重試邏輯)

#### 指數退避 (Exponential Backoff) + 隨機抖動 (Jitter)

重試時間計算公式：

```
backoff = initialBackoff × (multiplier ^ retryCount)
backoff = min(backoff, maxBackoff)
backoff += random(-jitter, +jitter)
```

#### 預設配置

```go
DefaultRetryPolicy() *RetryPolicy {
    return &RetryPolicy{
        MaxRetries:     5,                    // 最多重試 5 次
        InitialBackoff: 1 * time.Second,      // 初始等待 1 秒
        MaxBackoff:     5 * time.Minute,      // 最長等待 5 分鐘
        Multiplier:     2.0,                  // 每次翻倍
        JitterFraction: 0.1,                  // 10% 隨機抖動
    }
}
```

#### 重試時間範例

| 重試次數 | 基礎退避時間 | 實際範圍 (±10% jitter) |
|---------|------------|---------------------|
| 1       | 1 秒       | 0.9s - 1.1s         |
| 2       | 2 秒       | 1.8s - 2.2s         |
| 3       | 4 秒       | 3.6s - 4.4s         |
| 4       | 8 秒       | 7.2s - 8.8s         |
| 5       | 16 秒      | 14.4s - 17.6s       |

#### 處理 Retry-After Header

Teams API 在返回 429 錯誤時會包含 `Retry-After` header，指示何時可以重試：

```go
func RateLimitBackoff(retryAfterHeader string, defaultBackoff time.Duration) time.Duration {
    if retryAfter := GetRetryAfterDuration(retryAfterHeader); retryAfter > 0 {
        return retryAfter + (1 * time.Second) // 額外加 1 秒緩衝
    }
    return defaultBackoff
}
```

---

### 3. Notification Destinations (通知目的地管理)

**文件**: `internal/database/models.go`

#### 狀態管理

```go
type NotificationDestination struct {
    BaseModel
    NotificationID uuid.UUID  `json:"notification_id" db:"notification_id"`
    DestinationID  uuid.UUID  `json:"destination_id" db:"destination_id"`
    ConversationID *string    `json:"conversation_id" db:"conversation_id"`
    BotID          *uuid.UUID `json:"botId" db:"botId"`
    BotType        *BotType   `json:"bot_type" db:"bot_type"`
    Status         string     `json:"status" db:"status"` // pending, processing, sent, failed, cancelled
    ErrorMessage   *string    `json:"error_message" db:"error_message"`
    TeamsMessageID *string    `json:"teams_message_id" db:"teams_message_id"`
    SentAt         *time.Time `json:"sent_at" db:"sent_at"`
    RetryCount     int        `json:"retry_count" db:"retry_count"`
    MaxRetries     int        `json:"max_retries" db:"max_retries"`
    NextRetryAt    *time.Time `json:"next_retry_at" db:"next_retry_at"`
    FirstAttemptAt *time.Time `json:"first_attempt_at" db:"first_attempt_at"`
    LastAttemptAt  *time.Time `json:"last_attempt_at" db:"last_attempt_at"`
    FailureReason  *string    `json:"failure_reason" db:"failure_reason"`
    RetryAfter     *int       `json:"retry_after" db:"retry_after"`
    ActorID        *string    `json:"actor_id" db:"actor_id"`
}
```

#### 狀態轉換

```go
const (
    StatusPending    = "pending"    // 等待處理
    StatusProcessing = "processing" // 處理中
    StatusSent       = "sent"       // 發送成功
    StatusFailed     = "failed"     // 發送失敗
    StatusCancelled  = "cancelled"  // 已取消
)
```

---

### 4. Actor Pool V2 (Actor 池管理)

**文件**: `internal/actor/actor_pool_v2.go`

#### 核心職責

1. **Actor 管理**: 管理多個並發的 Notification Actor
2. **Redis 佇列**: 使用 Redis Sorted Set 管理重試佇列
3. **Circuit Breaker Integration**: 所有 Actor 共享電路斷路器狀態
4. **自動重試**: 根據 `next_retry_at` 自動重試失敗的通知

#### 配置參數

```go
type ActorPool struct {
    redis        *redis.Client
    db           ActorDB
    actors       map[uuid.UUID]*NotificationActor
    maxActors    int           // 最大 Actor 數量 (預設: 10)
    workerTicker *time.Ticker  // 輪詢間隔 (預設: 10 秒)
}
```

#### 工作流程

```go
// 1. 啟動 Actor Pool
actorPool.Start(ctx)

// 2. 定期輪詢重試佇列
func (p *ActorPool) pollAndSpawnActors(ctx context.Context) {
    // 檢查電路斷路器狀態
    if circuitBreaker.IsOpen() {
        return // 所有 Actor 暫停
    }
    
    // 從資料庫取出待重試的通知
    notifications := p.db.GetRetryReadyNotificationDestinations(ctx, 10)
    
    // 為每個通知創建 Actor
    for _, nd := range notifications {
        if len(p.actors) < p.maxActors {
            p.SpawnActor(ctx, nd.ID)
        }
    }
}

// 3. Actor 處理通知
func (actor *NotificationActor) Start(ctx context.Context) {
    // 檢查電路斷路器
    if circuitBreaker.IsOpen() {
        actor.Stop() // Actor 停止
        return
    }
    
    // 發送到 Teams
    result := actor.SendToTeams(ctx)
    
    // 更新狀態
    actor.UpdateStatus(result)
}
```

---

## 工作流程

### 完整流程圖

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant NotificationService
    participant ActorPool
    participant Actor
    participant Teams
    participant Redis
    participant DB

    Client->>API: POST /external/notify
    API->>NotificationService: SendNotification()
    NotificationService->>DB: INSERT notification (status: pending)
    NotificationService->>DB: INSERT notification_destinations (status: pending)
    NotificationService-->>API: Success (async)
    API-->>Client: 200 OK (processing)
    
    Note over ActorPool: Background Process
    ActorPool->>DB: SELECT retry-ready notifications
    ActorPool->>ActorPool: Check Circuit Breaker
    
    alt Circuit Breaker Closed
        ActorPool->>Actor: SpawnActor()
        Actor->>DB: UPDATE status = processing
        Actor->>Teams: Send Notification
        
        alt Success
            Teams-->>Actor: 200 OK
            Actor->>DB: UPDATE status = sent, sent_at = NOW()
            Actor->>ActorPool: Actor Complete
        else Failure
            Teams-->>Actor: Error (429/5xx)
            Actor->>DB: UPDATE retry_count++, next_retry_at, failure_reason
            Actor->>Redis: Add to retry queue
            Actor->>ActorPool: Actor Complete
        end
    else Circuit Breaker Open
        ActorPool->>ActorPool: All Actors Paused
        Note over ActorPool: Wait for Circuit Breaker to close
    end
```

### 步驟說明

#### 1. 通知建立 (完全非同步)

```go
// Notification Service 建立通知
func (s *notificationService) SendNotification(ctx context.Context, req *SendNotificationRequest) (*database.Notification, error) {
    // 1. 建立通知記錄
    notification := &database.Notification{
        ProjectID:   req.ProjectID,
        MessageType: req.MessageType,
        Content:     req.Message,
        Status:      "pending", // 開始為待處理狀態
    }
    err := s.repo.Create(ctx, notification)
    
    // 2. 為每個目的地建立 notification_destinations 記錄
    for _, dest := range destinations {
        for _, target := range dest.Targets {
            nd := &database.NotificationDestination{
                NotificationID: notification.ID,
                DestinationID:  dest.ID,
                ConversationID: &target.ConversationID,
                Status:         "pending",
                RetryCount:     0,
                MaxRetries:     5,
            }
            s.notificationDestRepo.CreateBatch(ctx, []*database.NotificationDestination{nd})
        }
    }
    
    return notification, nil
}
```

#### 2. Actor Pool 輪詢處理

```go
// Actor Pool 定期檢查待處理通知 (每 10 秒)
func (p *ActorPool) pollAndSpawnActors(ctx context.Context) {
    // 檢查電路斷路器狀態
    if circuitBreaker.IsOpen() {
        log.Println("Circuit breaker open, all actors paused")
        return
    }
    
    // 從資料庫取出待重試的通知
    notifications := p.db.GetRetryReadyNotificationDestinations(ctx, 10)
    
    // 為每個通知創建 Actor
    for _, nd := range notifications {
        if len(p.actors) < p.maxActors {
            p.SpawnActor(ctx, nd.ID)
        }
    }
}
```

#### 3. Actor 處理通知

```go
func (actor *NotificationActor) Start(ctx context.Context) {
    // 1. 更新狀態為處理中
    actor.db.UpdateNotificationDestination(ctx, actor.notificationDestID, &NotificationDestinationUpdate{
        Status: stringPtr("processing"),
        ActorID: stringPtr(actor.actorID),
    })
    
    // 2. 檢查電路斷路器
    if circuitBreaker.IsOpen() {
        actor.Stop()
        return
    }
    
    // 3. 發送到 Teams
    result := actor.teamsSender.Send(ctx, actor.notificationDest)
    
    // 4. 根據結果更新狀態
    if result.Success {
        actor.db.UpdateNotificationDestination(ctx, actor.notificationDestID, &NotificationDestinationUpdate{
            Status: stringPtr("sent"),
            TeamsMessageID: &result.MessageID,
            SentAt: &result.SentAt,
        })
    } else {
        // 失敗：更新重試資訊
        nextRetry := calculateNextRetry(actor.retryCount + 1, result.RetryAfter)
        actor.db.UpdateNotificationDestination(ctx, actor.notificationDestID, &NotificationDestinationUpdate{
            Status: stringPtr("pending"),
            RetryCount: intPtr(actor.retryCount + 1),
            NextRetryAt: &nextRetry,
            FailureReason: &result.FailureReason,
            ErrorMessage: &result.ErrorMessage,
        })
    }
}
```

---

## 資料庫設計

### Notification Destinations 表結構

**Migration 文件**: `scripts/migrations/009_enhance_notification_destinations_for_async_actor.sql`

```sql
-- 增強 notification_destinations 表以支援非同步處理
ALTER TABLE notification_destinations ADD COLUMN IF NOT EXISTS conversation_id VARCHAR(500);
ALTER TABLE notification_destinations ADD COLUMN IF NOT EXISTS next_retry_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE notification_destinations ADD COLUMN IF NOT EXISTS first_attempt_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE notification_destinations ADD COLUMN IF NOT EXISTS last_attempt_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE notification_destinations ADD COLUMN IF NOT EXISTS failure_reason VARCHAR(50);
ALTER TABLE notification_destinations ADD COLUMN IF NOT EXISTS retry_after INTEGER;
ALTER TABLE notification_destinations ADD COLUMN IF NOT EXISTS actor_id VARCHAR(100);
```

### 索引設計

```sql
-- 高效查詢待重試通知
CREATE INDEX IF NOT EXISTS idx_notification_destinations_next_retry 
    ON notification_destinations(next_retry_at) 
    WHERE retry_count < max_retries AND status = 'pending';

-- 按通知 ID 查詢
CREATE INDEX IF NOT EXISTS idx_notification_destinations_notification_id 
    ON notification_destinations(notification_id);

-- 按狀態查詢
CREATE INDEX IF NOT EXISTS idx_notification_destinations_status 
    ON notification_destinations(status);

-- 按對話 ID 查詢
CREATE INDEX IF NOT EXISTS idx_notification_destinations_conversation_id 
    ON notification_destinations(conversation_id);
```

### 監控視圖

```sql
CREATE OR REPLACE VIEW notification_destinations_queue_status AS
SELECT 
    COUNT(*) FILTER (WHERE status = 'pending' AND next_retry_at > NOW()) 
        as pending_count,
    COUNT(*) FILTER (WHERE status = 'pending' AND next_retry_at <= NOW()) 
        as ready_for_retry_count,
    COUNT(*) FILTER (WHERE status = 'processing') 
        as processing_count,
    COUNT(*) FILTER (WHERE status = 'sent') 
        as sent_count,
    COUNT(*) FILTER (WHERE status = 'failed') 
        as failed_count,
    MIN(next_retry_at) FILTER (WHERE status = 'pending' AND next_retry_at IS NOT NULL) 
        as next_retry_due,
    MIN(created_at) FILTER (WHERE status = 'pending') 
        as oldest_pending,
    MAX(updated_at) as last_activity
FROM notification_destinations;
```

---

## API 端點

### 1. 查詢佇列狀態

```bash
GET /api/v1/queue/status
```

**回應範例**:
```json
{
  "pending_count": 15,
  "ready_for_retry_count": 3,
  "processing_count": 2,
  "sent_count": 1250,
  "failed_count": 5,
  "next_retry_due": "2025-10-01T10:30:00Z",
  "oldest_pending": "2025-10-01T10:00:00Z",
  "circuit_state": "closed",
  "last_updated": "2025-10-01T10:25:30Z"
}
```

### 2. 查詢電路斷路器指標

```bash
GET /api/v1/queue/circuit-breaker/metrics
```

**回應範例**:
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

### 3. 重置電路斷路器

```bash
POST /api/v1/queue/circuit-breaker/reset
```

**使用場景**: 當確認 Teams 服務已恢復，手動重置電路斷路器

**回應範例**:
```json
{
  "message": "Circuit breaker reset successfully"
}
```

### 4. 查詢通知目的地狀態

```bash
GET /api/v1/queue/notifications/{notification_id}/destinations
```

**回應範例**:
```json
{
  "notification_id": "3f364043-eb50-49a6-988e-79a5be8f5ea0",
  "destinations": [
    {
      "id": "dest-1",
      "conversation_id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
      "status": "sent",
      "sent_at": "2025-10-01T10:25:30Z",
      "teams_message_id": "msg-123"
    },
    {
      "id": "dest-2",
      "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
      "status": "pending",
      "retry_count": 2,
      "next_retry_at": "2025-10-01T10:30:00Z",
      "failure_reason": "rate_limit"
    }
  ]
}
```

---

## 使用範例

### 測試腳本

**文件**: `scripts/test_queue.sh`

```bash
#!/bin/bash

# 1. 檢查佇列狀態
curl -X GET http://localhost:8080/api/v1/queue/status | jq '.'

# 2. 查看電路斷路器指標
curl -X GET http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '.'

# 3. 監控佇列變化
watch -n 2 'curl -s http://localhost:8080/api/v1/queue/status | jq "."'

# 4. 重置電路斷路器（如需要）
curl -X POST http://localhost:8080/api/v1/queue/circuit-breaker/reset

# 5. 查詢特定通知的目的地狀態
NOTIFICATION_ID="3f364043-eb50-49a6-988e-79a5be8f5ea0"
curl -X GET "http://localhost:8080/api/v1/queue/notifications/$NOTIFICATION_ID/destinations" | jq '.'
```

執行測試：
```bash
chmod +x scripts/test_queue.sh
./scripts/test_queue.sh
```

---

## 配置說明

### 環境變數

在 `cmd/server/main.go` 中初始化：

```go
// Actor Pool 配置
actorDB := actor.NewActorDB(notificationDestRepo, installationRepo)
actorPool := actor.NewActorPool(redisClient, actorDB, 10) // 最大 10 個並發 Actor

// 電路斷路器配置
circuitBreaker := &redisCircuitBreaker{
    redis: redisClient,
    key:   "circuit_breaker:teams_api",
    config: CircuitBreakerConfig{
        MaxFailures:     5,                    // 5 次失敗後熔斷
        Timeout:         1 * time.Minute,      // 熔斷 1 分鐘後測試恢復
        HalfOpenMaxReqs: 3,                    // 半開狀態最多 3 個測試請求
    },
}
```

### 調整建議

#### 高流量場景
```go
// 增加 Actor 數量
actorPool := actor.NewActorPool(redisClient, actorDB, 20)

// 更頻繁輪詢
actorPool.workerTicker = time.NewTicker(5 * time.Second)

// 調整電路斷路器參數
config: CircuitBreakerConfig{
    MaxFailures:     10,                   // 提高失敗閾值
    Timeout:         30 * time.Second,     // 縮短熔斷時間
    HalfOpenMaxReqs: 5,                    // 增加測試請求數
}
```

#### 低流量場景
```go
// 減少 Actor 數量
actorPool := actor.NewActorPool(redisClient, actorDB, 3)

// 降低輪詢頻率
actorPool.workerTicker = time.NewTicker(30 * time.Second)

// 調整電路斷路器參數
config: CircuitBreakerConfig{
    MaxFailures:     3,                    // 降低失敗閾值
    Timeout:         2 * time.Minute,      // 延長熔斷時間
    HalfOpenMaxReqs: 1,                    // 減少測試請求數
}
```

---

## 監控與維護

### 監控指標

1. **佇列積壓**
   - `pending_count`: 待處理數量
   - `ready_for_retry_count`: 準備重試數量
   - `processing_count`: 處理中數量
   - 告警閾值: pending_count > 100

2. **電路斷路器狀態**
   - `circuit_state`: closed/open/half-open
   - 告警: state = "open" 超過 5 分鐘

3. **失敗率**
   - `failed_count / (sent_count + failed_count)`
   - 告警閾值: > 5%

4. **Actor 狀態**
   - 活躍 Actor 數量
   - Actor 處理時間
   - Actor 錯誤率

### 維護操作

#### 1. 查看佇列積壓

```sql
SELECT 
    status,
    COUNT(*) as count,
    AVG(retry_count) as avg_retries,
    MIN(created_at) as oldest,
    MAX(updated_at) as newest
FROM notification_destinations
GROUP BY status
ORDER BY count DESC;
```

#### 2. 查看失敗原因統計

```sql
SELECT 
    failure_reason,
    COUNT(*) as count,
    AVG(retry_count) as avg_retries
FROM notification_destinations
WHERE status = 'pending' AND failure_reason IS NOT NULL
GROUP BY failure_reason
ORDER BY count DESC;
```

#### 3. 清理過期記錄

```sql
-- 手動清理 7 天前已耗盡的通知
UPDATE notification_destinations 
SET status = 'failed'
WHERE retry_count >= max_retries
  AND status = 'pending'
  AND updated_at < NOW() - INTERVAL '7 days';
```

#### 4. 重試指定通知

```sql
-- 重置特定通知的重試次數
UPDATE notification_destinations
SET retry_count = 0,
    next_retry_at = NOW(),
    status = 'pending',
    failure_reason = NULL,
    error_message = NULL
WHERE notification_id = 'xxx-xxx-xxx';
```

#### 5. 查看 Actor 狀態

```sql
-- 查看正在處理的通知
SELECT 
    actor_id,
    COUNT(*) as processing_count,
    MIN(updated_at) as oldest_processing
FROM notification_destinations
WHERE status = 'processing' AND actor_id IS NOT NULL
GROUP BY actor_id;
```

### 日誌監控

關鍵日誌：
```
Actor Pool V2 started successfully           // Actor Pool 啟動成功
Actor X started for notification Y           // Actor 啟動
Actor X: Processing notification Y           // 處理通知
Notification Y sent successfully             // 發送成功
Notification Y failed: error message         // 發送失敗
Circuit breaker is open, all actors paused   // 電路斷路器開啟，所有 Actor 暫停
Actor X completed for notification Y         // Actor 完成
```

---

## 故障排查

### 問題：佇列積壓過多

**可能原因**:
1. Teams API 持續返回錯誤
2. Actor 數量不足
3. 電路斷路器長期開啟

**解決方法**:
```bash
# 1. 檢查電路斷路器狀態
curl http://localhost:8080/api/v1/queue/circuit-breaker/metrics

# 2. 如果電路斷路器開啟，確認 Teams 服務正常後重置
curl -X POST http://localhost:8080/api/v1/queue/circuit-breaker/reset

# 3. 增加 Actor 數量（需重啟服務）
# 修改 actorPool := actor.NewActorPool(redisClient, actorDB, 20)
```

### 問題：通知重複發送

**可能原因**:
- Actor 處理邏輯問題
- 資料庫事務未正確提交
- Actor 狀態管理問題

**檢查**:
```sql
-- 查看是否有重複處理
SELECT conversation_id, COUNT(*)
FROM notification_destinations
WHERE status = 'sent'
GROUP BY conversation_id
HAVING COUNT(*) > 1;

-- 查看 Actor 狀態
SELECT actor_id, COUNT(*) as processing_count
FROM notification_destinations
WHERE status = 'processing'
GROUP BY actor_id;
```

### 問題：Actor 停止工作

**可能原因**:
- 電路斷路器開啟
- Redis 連線問題
- 資料庫連線問題

**檢查**:
```bash
# 1. 檢查電路斷路器狀態
curl http://localhost:8080/api/v1/queue/circuit-breaker/metrics

# 2. 檢查 Redis 連線
docker exec teamsnotify-redis redis-cli ping

# 3. 檢查資料庫連線
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT 1;"

# 4. 查看服務日誌
tail -f server.log | grep "Actor"
```

---

## 最佳實踐

### 1. 合理設置重試次數

```go
// 根據業務重要性調整
CriticalNotification.MaxRetries = 10  // 重要通知多重試
NormalNotification.MaxRetries = 5     // 一般通知
LowPriorityNotification.MaxRetries = 3 // 低優先級少重試
```

### 2. 監控告警設置

- 佇列積壓 > 100 → 警告
- 佇列積壓 > 500 → 嚴重警告
- 電路斷路器開啟 > 5 分鐘 → 告警
- 失敗率 > 10% → 告警
- Actor 處理時間 > 30 秒 → 告警

### 3. 定期清理

```sql
-- 自動清理 7 天前的記錄
UPDATE notification_destinations 
SET status = 'failed'
WHERE retry_count >= max_retries
  AND status = 'pending'
  AND updated_at < NOW() - INTERVAL '7 days';
```

### 4. Actor 管理最佳實踐

- 監控 Actor 數量，避免過多或過少
- 定期檢查 Actor 狀態，確保正常運作
- 在電路斷路器開啟時，所有 Actor 會自動暫停
- 使用 Redis 共享電路斷路器狀態，支援多節點部署

---

## 參考資料

- [Microsoft Teams Bot Rate Limits](https://learn.microsoft.com/en-us/microsoftteams/platform/bots/how-to/rate-limit)
- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html)
- [Exponential Backoff and Jitter](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)

---

## 版本歷史

- **v1.0.0** (2025-10-01)
  - 初始版本
  - 實現基礎隊列和熔斷器機制
  - 支援指數退避重試
  - 提供監控 API

---

**作者**: TeamsNotify Team  
**最後更新**: 2025-10-01


