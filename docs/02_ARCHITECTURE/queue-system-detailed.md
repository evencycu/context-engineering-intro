# Redis 佇列系統詳細設計

## 概述

本系統實現了一個基於 Redis 的優先序佇列系統，採用 Two-loop Enqueue Design 架構，將寫入和入隊分離為兩個獨立的迴圈，提高可靠性和簡化設計。

## 系統架構

### Two-loop Enqueue Design

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Request   │    │  EnqueueWorker  │    │ QueueConsumer   │
│                 │    │   (Background)  │    │   (Background)  │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          ▼                      ▼                      ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Write to DB     │    │ Scan Pending    │    │ BLPOP from      │
│ (status=pending)│    │ Push to Redis    │    │ Redis Queue     │
└─────────────────┘    │ (status=enqueued)│    └─────────┬───────┘
                       └─────────────────┘              │
                                                         ▼
                                                ┌─────────────────┐
                                                │   Actor Pool    │
                                                │   (Processing)  │
                                                └─────────────────┘
```

### 1. Producer Loop (API write only)
- **職責**: 接收 API 請求，寫入 `notification_destinations` 記錄
- **狀態**: `status = pending`
- **特點**: 不進行 Redis 操作，降低延遲
- **優勢**: 零耦合 Redis 可用性，提高寫入性能

### 2. EnqueueWorker (background)
- **職責**: 定期掃描 pending 狀態的 `notification_destinations`
- **操作**: 按優先級推入 Redis 隊列
- **狀態轉換**: `pending` → `enqueued`
- **特點**: 無重試機制，失敗保持 pending 狀態

### 3. Consumer Side (QueueConsumer + ActorPool)
- **QueueConsumer**: 執行 BLPOP 操作，從 Redis 取出任務
- **ActorPool**: 管理 Actor 並發，處理通知發送
- **特點**: 使用 SETNX 確保唯一處理

## 核心組件

### 1. Redis 優先序佇列

#### 佇列鍵設計
```go
// 三層優先序佇列
queue:notifications:high     // 高優先級
queue:notifications:normal   // 一般優先級  
queue:notifications:low      // 低優先級
```

#### 去重機制
```go
// 避免重複入列 (60秒 TTL)
queue:notifications:dedup:{notification_destination_id}

// 避免併發重複處理 (30分鐘 TTL)
queue:notifications:processing:{notification_destination_id}
```

#### 佇列操作
```go
type RedisQueue interface {
    Enqueue(ctx context.Context, ndID uuid.UUID, priority Priority) error
    DequeueBlocking(ctx context.Context) (uuid.UUID, error)
    Requeue(ctx context.Context, ndID uuid.UUID) error
    RequeueToFront(ctx context.Context, ndID uuid.UUID) error
}
```

### 2. EnqueueWorker

#### 核心邏輯
```go
func (w *EnqueueWorker) scanOnce(ctx context.Context) {
    // 1. 獲取 pending 狀態的記錄
    rows, err := w.repo.GetPending(ctx, 100)
    
    // 2. 按優先級推入 Redis 隊列
    for _, nd := range rows {
        prio := PriorityNormal
        switch nd.Priority {
        case "high": prio = PriorityHigh
        case "low":  prio = PriorityLow
        }
        
        // 3. 推入隊列
        if err := w.queue.Enqueue(ctx, nd.ID, prio); err != nil {
            continue // 保持 pending 狀態
        }
        
        // 4. 更新狀態為 enqueued
        w.repo.UpdateStatus(ctx, nd.ID, "enqueued")
    }
}
```

#### 配置參數
- **掃描間隔**: 200ms (可調整)
- **批次大小**: 100 筆 (可調整)
- **無重試機制**: 失敗保持 pending 狀態

### 3. QueueConsumer

#### 核心邏輯
```go
func (c *QueueConsumer) loop(ctx context.Context) {
    for {
        // 1. 檢查 Actor Pool 容量
        if !c.pool.HasCapacity() {
            time.Sleep(1 * time.Second)
            continue
        }
        
        // 2. 從 Redis 取出任務
        ndID, err := c.queue.DequeueBlocking(ctx)
        if err != nil {
            continue
        }
        
        // 3. 創建 Actor 處理
        success := c.pool.SpawnActor(ctx, ndID)
        if !success {
            // 重新入隊到高優先級
            c.queue.RequeueToFront(ctx, ndID)
        }
    }
}
```

#### 特點
- **容量檢查**: 避免取出無法處理的任務
- **BLPOP 操作**: 阻塞等待，按優先級順序
- **重入隊機制**: Actor Pool 滿時重新入隊

### 4. ActorPool

#### 核心職責
```go
type ActorPool struct {
    redis        *redis.Client
    db           ActorDB
    tokenManager TokenManager
    actors       map[uuid.UUID]*NotificationActor
    maxActors    int
}
```

#### 主要方法
```go
// 創建 Actor 處理通知
func (p *ActorPool) SpawnActor(ctx context.Context, ndID uuid.UUID) bool

// 檢查容量
func (p *ActorPool) HasCapacity() bool

// 獲取狀態
func (p *ActorPool) GetPoolStatus() map[string]interface{}
```

## 狀態轉換

### 完整流程
```
API Request → notification_destinations (status=pending)
     ↓
EnqueueWorker → Redis Queue (status=enqueued)
     ↓
QueueConsumer → ActorPool → NotificationActor
     ↓
Teams API → 發送結果 → 更新狀態 (sent/failed)
```

### 狀態定義
```go
const (
    StatusPending    = "pending"    // 等待處理
    StatusEnqueued   = "enqueued"   // 已入隊
    StatusProcessing = "processing" // 處理中
    StatusSent       = "sent"       // 發送成功
    StatusFailed     = "failed"     // 發送失敗
)
```

## 配置與監控

### 配置參數
```go
// EnqueueWorker 配置
enqueuer := NewEnqueueWorker(ndRepo, queue, 200*time.Millisecond)

// ActorPool 配置
pool := NewActorPool(redis, db, tokenManager, maxActors)

// QueueConsumer 配置
consumer := NewQueueConsumer(queue, pool)
```

### 監控指標
```bash
# 佇列長度
redis-cli LLEN queue:notifications:high
redis-cli LLEN queue:notifications:normal
redis-cli LLEN queue:notifications:low

# 佇列概覽
redis-cli KEYS "queue:notifications:*"

# 處理中任務
redis-cli KEYS "queue:notifications:processing:*"
```

### 資料庫監控
```sql
-- 佇列狀態視圖
CREATE OR REPLACE VIEW notification_destinations_queue_status AS
SELECT 
    COUNT(*) FILTER (WHERE status = 'pending') as pending_count,
    COUNT(*) FILTER (WHERE status = 'enqueued') as enqueued_count,
    COUNT(*) FILTER (WHERE status = 'processing') as processing_count,
    COUNT(*) FILTER (WHERE status = 'sent') as sent_count,
    COUNT(*) FILTER (WHERE status = 'failed') as failed_count
FROM notification_destinations;
```

## 錯誤處理

### 重試機制
- **EnqueueWorker**: 無重試，失敗保持 pending 狀態
- **QueueConsumer**: 使用 RequeueToFront 重新入隊
- **Actor**: 內建重試邏輯，支援指數退避

### 故障恢復
```go
// Actor Pool 滿時重新入隊
if !success {
    log.Printf("Actor pool became full, re-queuing task %s to front", ndID)
    c.queue.RequeueToFront(ctx, ndID)
}
```

## 性能優化

### 批次處理
- **EnqueueWorker**: 批次掃描 100 筆記錄
- **ActorPool**: 並發處理多個 Actor
- **QueueConsumer**: 單個任務處理

### 容量管理
- **Actor Pool**: 最大 Actor 數量限制
- **Queue Consumer**: 容量檢查避免溢出
- **Redis**: 使用 BLPOP 避免輪詢

## 部署建議

### 高流量場景
```go
// 增加 Actor 數量
maxActors := 20

// 縮短掃描間隔
interval := 100 * time.Millisecond

// 增加批次大小
batchSize := 200
```

### 低流量場景
```go
// 減少 Actor 數量
maxActors := 3

// 延長掃描間隔
interval := 1 * time.Second

// 減少批次大小
batchSize := 50
```

## 故障排查

### 常見問題
1. **佇列積壓**: 檢查 Actor Pool 容量和 Teams API 狀態
2. **重複處理**: 檢查 SETNX 機制和 processing 鍵
3. **任務丟失**: 檢查 EnqueueWorker 和 QueueConsumer 狀態

### 監控命令
```bash
# 檢查佇列狀態
curl http://localhost:8080/api/v1/queue/status

# 檢查 Actor Pool 狀態
curl http://localhost:8080/api/v1/queue/actors/status

# 重置佇列
curl -X POST http://localhost:8080/api/v1/queue/reset
```

## 最佳實踐

### 1. 容量規劃
- 根據 Teams API 限制調整 Actor 數量
- 監控佇列積壓情況
- 設置適當的告警閾值

### 2. 錯誤處理
- 使用指數退避重試
- 實現電路斷路器保護
- 記錄詳細的錯誤日誌

### 3. 監控維護
- 定期檢查佇列狀態
- 監控 Actor 處理時間
- 設置自動告警機制

---

**版本**: v1.0  
**最後更新**: 2025-10-08  
**作者**: TeamsNotify Team