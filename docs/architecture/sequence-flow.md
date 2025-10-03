# Queue & Circuit Breaker Sequence Flow

## 🏗️ 系統架構特性

- ✅ **多節點部署**: API Server 可水平擴展，多個 Node 並行處理
- ✅ **Redis Queue**: 使用 Redis 實現分散式隊列
- ✅ **Redis Circuit Breaker**: 使用 Redis 共享熔斷器狀態
- ✅ **Dead Letter Queue**: 達到最大重試後以 `notification_destinations.status=failed` 表示
- ✅ **智能退避**: 優先使用 429 Retry-After，否則使用指數退避

## 📋 目錄

1. [通知發送完整流程](#1-通知發送完整流程)
2. [失敗通知入隊流程](#2-失敗通知入隊流程-redis-queue)
3. [Worker 重試流程](#3-worker-重試流程-多節點)
4. [Circuit Breaker 狀態管理](#4-circuit-breaker-狀態管理-redis)
5. [Dead Letter Queue 處理](#5-dead-letter-queue-處理)
6. [API 查詢流程](#6-api-查詢流程)

---

## 1. 通知發送完整流程

### 成功場景

```mermaid
sequenceDiagram
    participant Client as 客戶端
    participant API as External Handler
    participant Service as External Service
    participant Broadcast as Broadcast Service
    participant Teams as Teams API
    participant DB as Database

    Client->>API: POST /api/v1/external/notify
    Note over Client,API: {notify_key, message, targets}
    
    API->>API: 驗證請求
    API->>Service: SendNotification()
    
    Service->>DB: 查詢 Project (by notify_key)
    DB-->>Service: Project 資料
    
    Service->>DB: 查詢 Destinations
    DB-->>Service: Destinations 列表
    
    Service->>DB: 建立 Notification 記錄
    DB-->>Service: Notification ID
    
    Service->>Broadcast: SendToTargets()
    Note over Service,Broadcast: targets, message
    
    loop 每個 Target
        Broadcast->>DB: 查詢 Bot Installation
        DB-->>Broadcast: Installation 資料
        
        Broadcast->>Broadcast: 取得 Bot Credentials
        Broadcast->>Teams: POST /v3/conversations/{id}/activities
        Note over Broadcast,Teams: Authorization: Bearer {token}
        
        Teams-->>Broadcast: 200 OK
        Note over Teams: 發送成功
    end
    
    Broadcast-->>Service: BroadcastResult
    Note over Broadcast,Service: {successful: 4, failed: 0}
    
    Service->>DB: 更新 Notification 狀態
    
    Service-->>API: Response
    API-->>Client: 200 OK
    Note over API,Client: {notification_id, success: true}
```

### 失敗場景（觸發 Queue）

```mermaid
sequenceDiagram
    participant Client as 客戶端
    participant API as External Handler
    participant Service as External Service
    participant Broadcast as Broadcast Service
    participant Teams as Teams API
    participant Queue as Queue Manager
    participant DB as Database

    Client->>API: POST /api/v1/external/notify
    
    API->>Service: SendNotification()
    Service->>Broadcast: SendToTargets()
    
    Broadcast->>Teams: POST /v3/conversations/{id}/activities
    
    alt Rate Limit (429)
        Teams-->>Broadcast: 429 Too Many Requests
        Note over Teams: Retry-After: 60
        
        Broadcast->>Broadcast: 檢查錯誤類型
        Note over Broadcast: IsRetryableError(429) = true
        
        Broadcast->>DB: INSERT INTO notification_destinations (status=pending, retry_count=0, max_retries=5, next_retry_at, failure_reason, retry_after)
        Note over Broadcast,Queue: {<br/>notification_id,<br/>target_id,<br/>message,<br/>reason: "rate_limit",<br/>retry_after: 60<br/>}
        
        Queue->>DB: UPDATE notification_destinations SET status='pending', retry_count=retry_count+1, next_retry_at=... WHERE id=...
        Note over Queue,DB: retry_count=0<br/>next_retry_at=now+60s
        DB-->>Queue: Success
        
        Queue-->>Broadcast: Enqueued
        
    else Timeout (408)
        Teams-->>Broadcast: 408 Request Timeout
        Broadcast->>DB: INSERT INTO notification_destinations (...)
        Note over Broadcast,Queue: reason: "timeout"
        
    else Server Error (5xx)
        Teams-->>Broadcast: 500 Internal Server Error
        Broadcast->>DB: INSERT INTO notification_destinations (...)
        Note over Broadcast,Queue: reason: "server_error"
        
    else Unauthorized (401/403)
        Teams-->>Broadcast: 401 Unauthorized
        Note over Broadcast: 不可重試，記錄錯誤
        Broadcast->>DB: 標記 Installation 為 stale
    end
    
    Broadcast-->>Service: BroadcastResult
    Note over Broadcast,Service: {successful: 3, failed: 1}
    
    Service-->>API: Partial Success
    API-->>Client: 200 OK
    Note over API,Client: {status: "partial",<br/>failed_count: 1}
```

---

## 2. 失敗通知入隊流程 (Redis Queue)

```mermaid
sequenceDiagram
    participant Broadcast as Broadcast Service
    participant ErrorHandler as Error Handler
    participant Queue as Queue Manager
    participant Policy as Retry Policy
    participant Redis as Redis Service
    participant DB as PostgreSQL

    Broadcast->>ErrorHandler: ClassifyFailureReason(statusCode, headers)
    
    alt 429 Rate Limit
        ErrorHandler->>ErrorHandler: 提取 Retry-After header
        
        alt 有 Retry-After header
            ErrorHandler-->>Broadcast: retry_after: 60 seconds (from header)
            Note over ErrorHandler: 🎯 優先使用 Teams API 的建議
        else 無 Retry-After header
            ErrorHandler-->>Broadcast: retry_after: null
            Note over ErrorHandler: 將使用預設指數退避
        end
        
    else 408 Timeout
        ErrorHandler-->>Broadcast: FailureReasonTimeout
        
    else 5xx Server Error
        ErrorHandler-->>Broadcast: FailureReasonServerError
    end
    
    Broadcast->>DB: 建立 NotificationDestination (status=pending, retry_count=0, max_retries=5)
    Note over Broadcast: {<br/>notification_id,<br/>project_id,<br/>target_id,<br/>message,<br/>reason,<br/>error_message,<br/>retry_count: 0,<br/>retry_after: 60 (from header)<br/>}
    
    Broadcast->>DB: INSERT INTO notification_destinations (...)
    
    Queue->>Policy: CalculateBackoff(retry_count=0, retry_after)
    
    alt 有 429 Retry-After (優先)
        Policy->>Policy: 使用 API 回傳的 Retry-After
        Note over Policy: retry_after = 60s (from header)<br/>+ 1s buffer = 61s
        Policy-->>Queue: next_retry_at = now + 61s
        
    else 無 Retry-After (使用指數退避)
        Policy->>Policy: InitialBackoff = 1 second
        Policy->>Policy: 加入隨機抖動 (±10%)
        Policy-->>Queue: next_retry_at = now + 0.9~1.1s
    end
    
    Queue->>Queue: 設置 max_retries = 5
    Queue->>Queue: 設置 next_retry_at
    Queue->>Queue: 計算 score = next_retry_at.unix()
    
    par 並行寫入
        Queue->>Redis: ZADD failed_queue {score} {notification_json}
        Note over Redis: Redis Sorted Set<br/>score = 重試時間戳<br/>member = 通知資料
        Redis-->>Queue: OK
        
    and 
        Queue->>DB: UPDATE notification_destinations SET status='pending', retry_count=retry_count+1, next_retry_at=... WHERE id=...
        Note over DB: 持久化備份<br/>防止 Redis 資料遺失
        DB-->>Queue: OK
    end
    
    Queue-->>Broadcast: Enqueued Successfully
    
    Note over Redis: 通知已加入 Redis Queue<br/>等待任一 Worker 處理
```

---

## 3. Worker 重試流程 (多節點)

### 多節點 Worker 並行處理

```mermaid
sequenceDiagram
    participant Node1 as Worker Node 1
    participant Node2 as Worker Node 2
    participant Node3 as Worker Node 3
    participant Redis as Redis Service
    participant RedisCB as Redis Circuit Breaker
    participant DB as PostgreSQL
    participant Teams as Teams API
    participant DLQ as Dead Letter Queue

    Note over Node1,Node3: 多個 API Server Node<br/>同時運行 Worker

    par Node 1 處理
        loop 每 10 秒
            Node1->>RedisCB: GET circuit_breaker:state
            
            alt Circuit Breaker OPEN
                RedisCB-->>Node1: "open"
                Note over Node1: ⏸️ 跳過處理<br/>等待熔斷器恢復
                
            else Circuit Breaker CLOSED/HALF-OPEN
                RedisCB-->>Node1: "closed"
                
                Node1->>Redis: ZPOPMIN failed_queue 5
                Note over Redis: 原子操作<br/>取出 5 個最早到期的通知<br/>避免多節點重複處理
                Redis-->>Node1: [notification1, notification2]
                
                loop 處理每個通知
                    Node1->>Node1: 解析通知資料
                    Note over Node1: retry_count = 2<br/>max_retries = 5
                    
                    Node1->>Teams: 發送通知
                    
                    alt 發送成功
                        Teams-->>Node1: 200 OK
                        Node1->>RedisCB: INCR circuit_breaker:success
                        Node1->>DB: UPDATE notification_destinations SET status='sent', sent_at=NOW() WHERE id=...
                        Note over Node1: ✅ 成功，移除備份
                        
                    else 429 Rate Limit (有 Retry-After)
                        Teams-->>Node1: 429, Retry-After: 90
                        Note over Node1: 🎯 使用 Retry-After = 90s
                        Node1->>Node1: retry_count++
                        
                        alt retry_count < max_retries
                            Node1->>Node1: next_retry = now + 90s
                            Node1->>Redis: ZADD failed_queue {score} {notification}
                            Node1->>DB: UPDATE retry_count, next_retry_at
                            Note over Node1: 📅 重新入隊，等待 90 秒
                        else
                            Node1->>DLQ: 移至 Dead Letter Queue
                            Note over Node1: ❌ 達到 max_retries
                        end
                        
                        Node1->>RedisCB: INCR circuit_breaker:failures
                        
                    else 408/5xx (無 Retry-After)
                        Teams-->>Node1: 408 Timeout
                        Node1->>Node1: retry_count++
                        
                        alt retry_count < max_retries
                            Node1->>Node1: 計算指數退避<br/>backoff = 1s * 2^(retry_count) ± 10%
                            Node1->>Redis: ZADD failed_queue {score} {notification}
                            Node1->>DB: UPDATE retry_count, next_retry_at
                            Note over Node1: 📅 使用指數退避
                        else
                            Node1->>DLQ: 移至 Dead Letter Queue
                        end
                    end
                end
            end
        end
        
    and Node 2 處理
        loop 每 10 秒 (錯開時間)
            Note over Node2: 延遲 3 秒啟動
            Node2->>RedisCB: GET circuit_breaker:state
            RedisCB-->>Node2: "closed"
            
            Node2->>Redis: ZPOPMIN failed_queue 5
            Note over Redis: 原子操作<br/>Node2 取得不同的通知
            Redis-->>Node2: [notification3, notification4, notification5]
            
            Note over Node2: 處理邏輯同 Node 1
        end
        
    and Node 3 處理
        loop 每 10 秒 (錯開時間)
            Note over Node3: 延遲 6 秒啟動
            Node3->>RedisCB: GET circuit_breaker:state
            RedisCB-->>Node3: "closed"
            
            Node3->>Redis: ZPOPMIN failed_queue 5
            Redis-->>Node3: [] (隊列已空)
            
            Note over Node3: 無待處理項目
        end
    end
    
    Note over Redis: ✅ Redis Sorted Set 確保:<br/>1. 原子操作，無重複處理<br/>2. 按時間排序<br/>3. 多節點負載均衡
```

### 重試退避策略範例

```mermaid
sequenceDiagram
    participant N as Notification
    participant Teams as Teams API
    participant Policy as Retry Policy
    
    Note over N,Policy: 初始失敗 (429 with Retry-After)
    N->>Teams: 發送通知
    Teams-->>N: 429, Retry-After: 60
    N->>Policy: CalculateBackoff(0, retry_after=60)
    Note over Policy: 🎯 優先使用 Retry-After header
    Policy-->>N: 60 seconds (+ 1s buffer = 61s)
    Note over N: retry_count=0, 等待 61s
    
    Note over N,Policy: 第 1 次重試失敗 (無 Retry-After)
    N->>Teams: 發送通知
    Teams-->>N: 408 Timeout (no header)
    N->>Policy: CalculateBackoff(1, retry_after=null)
    Note over Policy: 使用指數退避 + 抖動
    Policy-->>N: 2 seconds (1.8~2.2s)
    Note over N: retry_count=1, 等待 2s
    
    Note over N,Policy: 第 2 次重試失敗
    N->>Policy: CalculateBackoff(2)
    Policy-->>N: 4 seconds (3.6~4.4s)
    Note over N: retry_count=2, 等待 4s
    
    Note over N,Policy: 第 3 次重試失敗
    N->>Policy: CalculateBackoff(3)
    Policy-->>N: 8 seconds (7.2~8.8s)
    Note over N: retry_count=3, 等待 8s
    
    Note over N,Policy: 第 4 次重試失敗
    N->>Policy: CalculateBackoff(4)
    Policy-->>N: 16 seconds (14.4~17.6s)
    Note over N: retry_count=4, 等待 16s
    
    Note over N,Policy: 第 5 次重試失敗
    N->>Policy: ShouldRetry(5)
    Policy-->>N: false (retry_count >= max_retries)
    Note over N: ❌ 移至 Dead Letter Queue
```

---

## 4. Circuit Breaker 狀態管理 (Redis)

### 使用 Redis 的分散式熔斷器

```mermaid
sequenceDiagram
    participant Node1 as Worker Node 1
    participant Node2 as Worker Node 2
    participant Redis as Redis Circuit Breaker
    participant Teams as Teams API

    Note over Redis: Redis Keys:<br/>cb:state = "closed"<br/>cb:failures = 0<br/>cb:last_failure = null<br/>cb:last_429_retry_after = null

    rect rgb(200, 255, 200)
        Note over Redis: 狀態: CLOSED (正常)
        
        par Node 1 發送
            Node1->>Redis: GET cb:state
            Redis-->>Node1: "closed"
            Node1->>Redis: GET cb:failures
            Redis-->>Node1: 0
            
            Node1->>Teams: 發送通知
            
            alt 成功
                Teams-->>Node1: 200 OK
                Node1->>Redis: SET cb:failures 0
                Node1->>Redis: INCR cb:success_count
                
            else 429 Rate Limit
                Teams-->>Node1: 429, Retry-After: 120
                Note over Node1: 🎯 記錄 Retry-After
                Node1->>Redis: INCR cb:failures
                Node1->>Redis: SET cb:last_429_retry_after 120
                Node1->>Redis: SET cb:last_failure {timestamp}
                Redis-->>Node1: failures = 1
                
            else 其他錯誤
                Teams-->>Node1: 408/5xx
                Node1->>Redis: INCR cb:failures
                Node1->>Redis: SET cb:last_failure {timestamp}
                Redis-->>Node1: failures = 1
            end
            
        and Node 2 發送 (並行)
            Node2->>Redis: GET cb:state
            Redis-->>Node2: "closed"
            
            Node2->>Teams: 發送通知
            Teams-->>Node2: 429
            
            Node2->>Redis: INCR cb:failures
            Redis-->>Node2: failures = 2
        end
        
        Note over Redis: ⚠️ 失敗計數累積中
        
        loop 持續失敗
            Node1->>Redis: INCR cb:failures
            Redis-->>Node1: failures = 5
            
            alt failures >= 5
                Node1->>Redis: Lua Script: 檢查並設置 OPEN
                Note over Redis: 原子操作:<br/>if failures >= 5 then<br/>  SET cb:state "open"<br/>  SET cb:open_until (now + retry_after or 60s)<br/>end
                
                Redis->>Redis: GET cb:last_429_retry_after
                Note over Redis: retry_after = 120s (from last 429)
                
                Redis->>Redis: SET cb:state "open"
                Redis->>Redis: SET cb:open_until {now + 120s}
                Note over Redis: 🔴 觸發熔斷！<br/>使用 Retry-After = 120s<br/>如無則使用 60s 預設值
                
                Redis-->>Node1: State changed to OPEN
            end
        end
    end

    rect rgb(255, 200, 200)
        Note over Redis: 狀態: OPEN (熔斷)
        
        par 多節點檢查
            Node1->>Redis: GET cb:state
            Redis-->>Node1: "open"
            Node1->>Redis: GET cb:open_until
            Redis-->>Node1: {timestamp}
            
            alt 未到 open_until 時間
                Node1->>Node1: 計算剩餘時間
                Note over Node1: ⏸️ 跳過處理<br/>還需等待 90 秒
                
            else 已過 open_until 時間
                Node1->>Redis: Lua Script: OPEN → HALF-OPEN
                Note over Redis: if now >= open_until then<br/>  SET cb:state "half-open"<br/>  SET cb:half_open_tests 0<br/>end
                
                Redis->>Redis: SET cb:state "half-open"
                Redis->>Redis: SET cb:half_open_tests 0
                Note over Redis: 🟡 進入測試狀態
                Redis-->>Node1: State changed to HALF-OPEN
            end
            
        and
            Node2->>Redis: GET cb:state
            Redis-->>Node2: "open"
            Note over Node2: ⏸️ 同樣跳過處理<br/>所有節點遵循相同狀態
        end
    end

    rect rgb(255, 255, 200)
        Note over Redis: 狀態: HALF-OPEN (測試)
        
        Note over Redis: 允許最多 3 個測試請求
        
        Node1->>Redis: Lua Script: 獲取測試配額
        Note over Redis: if state == "half-open" and<br/>   half_open_tests < 3 then<br/>  INCR cb:half_open_tests<br/>  return true<br/>end
        
        Redis->>Redis: INCR cb:half_open_tests
        Redis-->>Node1: Allowed (1/3)
        
        Node1->>Teams: 測試發送
        
        alt 測試成功
            Teams-->>Node1: 200 OK
            Node1->>Redis: INCR cb:half_open_success
            
            Node1->>Redis: GET cb:half_open_success
            Redis-->>Node1: 3
            
            alt 3 個測試都成功
                Node1->>Redis: Lua Script: HALF-OPEN → CLOSED
                Note over Redis: if half_open_success >= 3 then<br/>  SET cb:state "closed"<br/>  SET cb:failures 0<br/>  DEL cb:open_until<br/>end
                
                Redis->>Redis: SET cb:state "closed"
                Redis->>Redis: SET cb:failures 0
                Redis->>Redis: DEL cb:last_429_retry_after
                Note over Redis: 🟢 恢復正常
                Redis-->>Node1: State changed to CLOSED
            end
            
        else 測試失敗
            Teams-->>Node1: 429/Error
            Node1->>Redis: Lua Script: HALF-OPEN → OPEN
            
            Redis->>Redis: 檢查是否有新的 Retry-After
            alt 有 Retry-After
                Note over Redis: 使用新的 Retry-After = 180s
                Redis->>Redis: SET cb:state "open"
                Redis->>Redis: SET cb:open_until {now + 180s}
            else 無 Retry-After
                Redis->>Redis: SET cb:state "open"
                Redis->>Redis: SET cb:open_until {now + 120s (last retry_after)}
            end
            
            Note over Redis: 🔴 重新熔斷<br/>優先使用最新 Retry-After
            Redis-->>Node1: State changed to OPEN
        end
    end
```

### 狀態轉換圖

```
                     ┌─────────────────┐
                     │     CLOSED      │
                     │   (正常運行)     │
                     └────────┬────────┘
                              │
                     連續失敗 ≥ 5 次
                              │
                              ▼
                     ┌─────────────────┐
            ┌────────│      OPEN       │
            │        │   (熔斷拒絕)     │◄────┐
            │        └────────┬────────┘     │
            │                 │               │
            │        等待 1 分鐘後             │
            │                 │               │
            │                 ▼               │
            │        ┌─────────────────┐     │
            │        │   HALF-OPEN     │     │
            │        │  (測試恢復)      │     │
            │        └────────┬────────┘     │
            │                 │               │
            │        ┌────────┴────────┐     │
            │        │                 │     │
            │  3個測試請求成功   3個測試請求中  │
            │        │           任一失敗     │
            │        ▼                 │     │
            └────成功────────────────失敗───┘
```

---

## 5. Dead Letter Queue 處理

### 達到最大重試後的處理流程

```mermaid
sequenceDiagram
    participant Worker as Queue Worker
    participant Redis as Redis Queue
    participant DB as PostgreSQL
    participant DLQ as Notification Destinations (DLQ)
    participant Alert as Alert System

    Worker->>Redis: ZPOPMIN failed_queue 5
    Redis-->>Worker: [notification]
    
    Worker->>Worker: 解析通知
    Note over Worker: retry_count = 4<br/>max_retries = 5
    
    Worker->>Worker: 嘗試重試
    Note over Worker: 第 5 次重試
    
    alt 第 5 次重試失敗
        Worker->>Worker: retry_count++ (now = 5)
        Worker->>Worker: retry_count >= max_retries?
        Note over Worker: ✅ True, 達到上限
        
        Worker->>Worker: 準備 Dead Notification
        Note over Worker: {<br/>  notification_id,<br/>  project_id,<br/>  target_id,<br/>  message,<br/>  original_reason: "rate_limit",<br/>  total_retry_attempts: 5,<br/>  first_failed_at: "2025-10-01 10:00:00",<br/>  last_failed_at: "2025-10-01 10:05:00",<br/>  last_error: "429 Too Many Requests",<br/>  status: "exhausted"<br/>}
        
        par 並行處理
        Worker->>DLQ: UPDATE notification_destinations SET status='failed', last_attempt_at=NOW() WHERE id=...
        Note over DLQ: 將目的地標記為 failed（視為 Dead Letter）
            DLQ-->>Worker: OK
            
        and
            Worker->>DB: UPDATE notification_destinations SET status='failed' WHERE id=...
            Note over DB: 從重試隊列移除
            DB-->>Worker: OK
            
        and
            Worker->>Alert: 發送告警
            Note over Alert: 通知運維團隊<br/>有通知最終失敗
            Alert-->>Worker: Alert Sent
        end
        
        Worker->>Worker: 記錄 Metrics
        Note over Worker: dead_letter_count++<br/>failure_reason: "rate_limit"
        
    else 第 5 次重試成功
        Note over Worker: ✅ 成功發送
        Worker->>Redis: (不重新入隊)
        Worker->>DB: UPDATE notification_destinations SET status='failed' WHERE id=...
        Note over Worker: 🎉 最後一次機會成功！
    end
```

### Dead Notifications 表結構

```sql
CREATE TABLE dead_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL,
    project_id UUID NOT NULL,
    target_id TEXT NOT NULL,
    target_type TEXT NOT NULL,
    message JSONB NOT NULL,
    
    -- 失敗原因
    original_reason TEXT NOT NULL, -- rate_limit, timeout, server_error
    last_error TEXT,
    
    -- 重試歷史
    total_retry_attempts INT NOT NULL DEFAULT 5,
    first_failed_at TIMESTAMP NOT NULL,
    last_failed_at TIMESTAMP NOT NULL,
    
    -- 狀態
    status TEXT NOT NULL DEFAULT 'exhausted', -- exhausted, manual_retry, resolved
    
    -- 元數據
    metadata JSONB, -- 額外的診斷資訊
    
    -- 時間戳
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP,
    
    -- 索引
    INDEX idx_dead_notifications_project (project_id),
    INDEX idx_dead_notifications_status (status),
    INDEX idx_dead_notifications_created (created_at DESC)
);
```

### Dead Letter Queue 查詢 API

```mermaid
sequenceDiagram
    participant Client as 客戶端
    participant API as API Handler
    participant DB as PostgreSQL

    Client->>API: GET /api/v1/queue/dead-letter
    Note over Client: 查詢參數:<br/>?status=exhausted<br/>&limit=50<br/>&offset=0
    
    API->>DB: SELECT ... FROM notification_destinations WHERE status='failed'
    Note over DB: ORDER BY last_attempt_at DESC<br/>LIMIT 50 OFFSET 0
    
    DB-->>API: destinations (failed)
    
    API-->>Client: 200 OK
    Note over Client: {<br/>  "total": 123,<br/>  "items": [<br/>    {<br/>      "id": "uuid",<br/>      "notification_id": "uuid",<br/>      "target_id": "channel-id",<br/>      "reason": "rate_limit",<br/>      "total_attempts": 5,<br/>      "first_failed_at": "...",<br/>      "last_failed_at": "..."<br/>    }<br/>  ]<br/>}
```

### 手動重試 Dead Letter

```mermaid
sequenceDiagram
    participant Admin as 管理員
    participant API as API Handler
    participant DB as PostgreSQL
    participant Redis as Redis Queue
    participant DLQ as Notification Destinations (DLQ)

    Admin->>API: POST /api/v1/queue/dead-letter/{id}/retry
    Note over Admin: 手動重試失敗通知
    
    API->>DLQ: SELECT FROM notification_destinations WHERE id = {id} AND status='failed'
    DLQ-->>API: notification_destination (failed)
    
    API->>API: 驗證可重試
    Note over API: status = "exhausted"<br/>resolved_at IS NULL
    
    API->>API: 重建 Notification
    Note over API: {<br/>  notification_id,<br/>  target_id,<br/>  message,<br/>  retry_count: 0,<br/>  max_retries: 5<br/>}
    
    par 並行處理
        API->>Redis: ZADD failed_queue {score} {notification}
        Note over Redis: 重新加入重試隊列<br/>score = now (立即重試)
        Redis-->>API: OK
        
    and
        API->>DB: INSERT INTO notification_destinations
        DB-->>API: OK
        
    and
        API->>DLQ: UPDATE notification_destinations<br/>SET status = 'pending', retry_count = 0, next_retry_at = NOW()
        DLQ-->>API: OK
    end
    
    API-->>Admin: 200 OK
    Note over Admin: {<br/>  "message": "Notification requeued",<br/>  "queue_position": 1<br/>}
    
    Note over Redis: Worker 將在下一輪處理
```

---

## 6. API 查詢流程

### 6.1 查詢隊列狀態

```mermaid
sequenceDiagram
    participant Client as 客戶端
    participant Handler as Queue Handler
    participant Redis as Redis Service
    participant RedisCB as Redis Circuit Breaker
    participant DB as PostgreSQL

    Client->>Handler: GET /api/v1/queue/status
    
    par 並行查詢
        Handler->>Redis: ZCARD failed_queue
        Redis-->>Handler: total_pending: 15
        
    and
        Handler->>Redis: ZCOUNT failed_queue -inf {now}
        Note over Redis: score <= now (可重試)
        Redis-->>Handler: ready_for_retry: 5
        
    and
        Handler->>Redis: ZRANGE failed_queue 0 0 WITHSCORES
        Note over Redis: 取得最早的通知
        Redis-->>Handler: next_retry_due: timestamp
        
    and
        Handler->>DB: SELECT COUNT(*) FROM dead_notifications
        DB-->>Handler: total_dead: 3
        
    and
        Handler->>RedisCB: GET cb:state
        RedisCB-->>Handler: circuit_state: "closed"
        
    and
        Handler->>RedisCB: GET cb:open_until
        RedisCB-->>Handler: open_until: null
    end
    
    Handler->>Handler: 組合完整狀態
    Note over Handler: {<br/>  total_pending: 15,<br/>  ready_for_retry: 5,<br/>  total_dead: 3,<br/>  circuit_state: "closed",<br/>  next_retry_due: "...",<br/>  open_until: null<br/>}
    
    Handler-->>Client: 200 OK
    Note over Client: {<br/>  "total_pending": 15,<br/>  "ready_for_retry": 5,<br/>  "total_dead": 3,<br/>  "circuit_state": "closed",<br/>  "next_retry_due": "2025-10-01T14:05:00Z",<br/>  "open_until": null,<br/>  "last_updated": "2025-10-01T14:00:00Z"<br/>}
```

### 6.2 查詢熔斷器指標

```mermaid
sequenceDiagram
    participant Client as 客戶端
    participant Handler as Queue Handler
    participant Redis as Redis Circuit Breaker

    Client->>Handler: GET /api/v1/queue/circuit-breaker/metrics
    
    par 並行查詢 Redis
        Handler->>Redis: GET cb:state
        Redis-->>Handler: "closed"
        
    and
        Handler->>Redis: GET cb:failures
        Redis-->>Handler: 0
        
    and
        Handler->>Redis: GET cb:success_count
        Redis-->>Handler: 1498
        
    and
        Handler->>Redis: GET cb:failed_count
        Redis-->>Handler: 25
        
    and
        Handler->>Redis: GET cb:last_state_change
        Redis-->>Handler: "2025-10-01T13:50:08Z"
        
    and
        Handler->>Redis: GET cb:last_429_retry_after
        Redis-->>Handler: null
        
    and
        Handler->>Redis: GET cb:open_until
        Redis-->>Handler: null
    end
    
    Handler->>Handler: 組合 Metrics
    Note over Handler: {<br/>  state: "closed",<br/>  failures: 0,<br/>  success_count: 1498,<br/>  failed_count: 25,<br/>  last_429_retry_after: null,<br/>  open_until: null<br/>}
    
    Handler-->>Client: 200 OK
    Note over Client: {<br/>  "state": "closed",<br/>  "consecutive_failures": 0,<br/>  "total_success": 1498,<br/>  "total_failed": 25,<br/>  "last_429_retry_after": null,<br/>  "open_until": null,<br/>  "last_state_change": "2025-10-01T13:50:08Z"<br/>}
```

### 6.3 重置熔斷器

```mermaid
sequenceDiagram
    participant Admin as 管理員
    participant Handler as Queue Handler
    participant Redis as Redis Circuit Breaker

    Admin->>Handler: POST /api/v1/queue/circuit-breaker/reset
    Note over Admin: 手動重置熔斷器
    
    Handler->>Redis: Multi/Exec Transaction
    Note over Redis: 原子操作批次執行
    
    par 批次重置 Redis Keys
        Handler->>Redis: SET cb:state "closed"
        Handler->>Redis: SET cb:failures 0
        Handler->>Redis: DEL cb:open_until
        Handler->>Redis: DEL cb:last_429_retry_after
        Handler->>Redis: DEL cb:half_open_tests
        Handler->>Redis: DEL cb:half_open_success
        Handler->>Redis: SET cb:last_state_change {now}
    end
    
    Redis->>Redis: EXEC (執行事務)
    Redis-->>Handler: OK
    
    Handler->>Handler: 記錄日誌
    Note over Handler: "Circuit breaker manually reset by admin"
    
    Handler-->>Admin: 200 OK
    Note over Admin: {<br/>  "message": "Circuit breaker reset successfully",<br/>  "new_state": "closed",<br/>  "timestamp": "2025-10-01T14:00:00Z"<br/>}
    
    Note over Redis: 🟢 熔斷器已重置<br/>所有節點立即恢復正常運行
```

---

## 📊 時序圖總覽

### 完整系統互動 (多節點 + Redis)

```mermaid
sequenceDiagram
    participant C as Client
    participant Node1 as API Node 1
    participant Node2 as API Node 2
    participant Redis as Redis
    participant RedisCB as Redis CB
    participant T as Teams API
    participant DB as PostgreSQL
    participant DLQ as Dead Letter

    Note over C,DLQ: 階段 1: 初始請求 (任一節點)
    C->>Node1: POST /api/v1/external/notify
    Node1->>Node1: Broadcast Service
    Node1->>T: 嘗試發送
    
    alt 成功
        T-->>Node1: 200 OK
        Node1->>RedisCB: INCR cb:success_count
        Node1-->>C: 成功
        
    else 429 Rate Limit
        T-->>Node1: 429, Retry-After: 120
        
        Note over Node1,Redis: 階段 2: 入隊處理
        Node1->>Node1: 提取 Retry-After = 120s
        Node1->>Node1: score = now + 120s
        
        par 並行寫入
            Node1->>Redis: ZADD failed_queue {score} {notification}
            Redis-->>Node1: OK
        and
            Node1->>DB: UPDATE notification_destinations SET retry_count=retry_count+1, next_retry_at=...
            DB-->>Node1: OK
        and
            Node1->>RedisCB: INCR cb:failures
            Node1->>RedisCB: SET cb:last_429_retry_after 120
        end
        
        Node1-->>C: 202 Accepted (已入隊重試)
    end
    
    Note over Node2,DLQ: 階段 3: 背景重試 (多節點並行)
    
    par Node 1 Worker
        loop 每 10 秒
            Node1->>RedisCB: GET cb:state
            RedisCB-->>Node1: "closed"
            
            Node1->>Redis: ZPOPMIN failed_queue 5
            Redis-->>Node1: [notification1, notification2]
            
            loop 處理每個通知
                Node1->>T: 重試發送
                
                alt 成功
                    T-->>Node1: 200 OK
                    Node1->>RedisCB: INCR cb:success_count
                    Node1->>DB: UPDATE notification_destinations SET status='failed' WHERE id=...
                    
                else 再次失敗 (有 Retry-After)
                    T-->>Node1: 429, Retry-After: 180
                    Node1->>Node1: retry_count++
                    
                    alt retry_count < 5
                        Node1->>Redis: ZADD failed_queue {now+180s} {notification}
                        Node1->>DB: UPDATE retry_count
                    else retry_count >= 5
                        Node1->>DLQ: INSERT dead_notifications
                        Node1->>DB: UPDATE notification_destinations SET status='failed' WHERE id=...
                        Note over DLQ: ❌ 移至 Dead Letter
                    end
                    
                    Node1->>RedisCB: INCR cb:failures
                    
                    alt failures >= 5
                        Node1->>RedisCB: SET cb:state "open"
                        Node1->>RedisCB: SET cb:open_until {now+180s}
                        Note over RedisCB: 🔴 觸發熔斷 (使用 Retry-After)
                    end
                    
                else 再次失敗 (無 Retry-After)
                    T-->>Node1: 408 Timeout
                    Node1->>Node1: backoff = 2^(retry_count) ± 10%
                    Node1->>Redis: ZADD failed_queue {now+backoff} {notification}
                end
            end
        end
        
    and Node 2 Worker (並行)
        loop 每 10 秒
            Node2->>RedisCB: GET cb:state
            RedisCB-->>Node2: "open"
            Note over Node2: ⏸️ 熔斷中，跳過處理
        end
    end
    
    Note over C,DLQ: 階段 4: 監控查詢
    C->>Node1: GET /api/v1/queue/status
    
    par 並行查詢
        Node1->>Redis: ZCARD failed_queue
        Redis-->>Node1: 10
    and
        Node1->>RedisCB: GET cb:state
        RedisCB-->>Node1: "open"
    and
        Node1->>RedisCB: GET cb:open_until
        RedisCB-->>Node1: {timestamp}
    and
        Node1->>DB: SELECT COUNT(*) FROM dead_notifications
        DB-->>Node1: 3
    end
    
    Node1-->>C: {pending: 10, state: "open", dead: 3}
```

---

## 🔢 重要時間參數

| 參數 | 預設值 | 說明 | 備註 |
|------|--------|------|------|
| Worker 輪詢間隔 | 10 秒 | Worker 檢查隊列的頻率 | 可配置 |
| 批次處理大小 | 5 個/節點 | 每次處理的通知數量 | 避免雪崩 |
| 初始退避時間 | 1 秒 | 第一次重試的等待時間 | 無 Retry-After 時 |
| 最大退避時間 | 5 分鐘 | 重試等待的上限 | 指數退避上限 |
| 退避倍數 | 2.0 | 每次重試時間翻倍 | 1s → 2s → 4s → 8s → 16s |
| 隨機抖動 | ±10% | 避免雪崩效應 | 分散重試時間 |
| 最大重試次數 | 5 次 | 超過後移至 DLQ | 可配置 |
| 熔斷失敗閾值 | 5 次 | 連續失敗觸發熔斷 | 可配置 |
| 熔斷超時時間 | **優先使用 Retry-After** | Open 狀態持續時間 | **429 header 優先** |
| 預設熔斷時間 | 60 秒 | 無 Retry-After 時使用 | 後備機制 |
| 半開測試數量 | 3 個 | Half-Open 狀態的測試請求 | 成功後恢復 |
| Redis Key TTL | 7 天 | Redis 鍵過期時間 | 自動清理 |
| DLQ 保留時間 | 30 天 | Dead Letter 記錄保留期 | 可手動清理 |

---

## 📝 狀態說明

### 通知狀態
- **pending**: 等待重試（Redis score > now）
- **ready**: 可以重試（Redis score <= now）
- **retrying**: 正在重試中（已被 Worker ZPOPMIN）
- **exhausted**: 已耗盡重試次數（移至 dead_notifications）
- **success**: 重試成功（已從 Redis 和 DB 移除）

### 失敗原因
- **rate_limit**: 429 Too Many Requests（可能有 Retry-After header）
- **timeout**: 408 Request Timeout
- **server_error**: 5xx Server Error
- **network_error**: 網路連接問題
- **unauthorized**: 401/403 認證失敗（**不重試**）
- **bad_request**: 4xx Client Error（**不重試**）
- **unknown**: 其他未知錯誤

### Circuit Breaker 狀態 (Redis 共享)
- **closed**: 正常運行，所有節點允許請求
- **open**: 熔斷狀態，所有節點拒絕請求
- **half-open**: 測試狀態，允許有限數量的測試請求（3 個）

### Dead Letter 狀態
- **exhausted**: 達到最大重試次數
- **manual_retry**: 管理員手動重試
- **resolved**: 問題已解決（不再重試）

---

## 🎯 關鍵決策點

### 1. 是否入隊？
```go
func ShouldEnqueue(statusCode int, headers http.Header) bool {
    if statusCode == 429 || statusCode == 408 || (statusCode >= 500 && statusCode < 600) {
        // 可重試錯誤 → 入隊
        retryAfter := ExtractRetryAfter(headers) // 提取 Retry-After
        EnqueueToRedis(notification, retryAfter)
        return true
    }
    // 4xx (除 408) → 不重試，記錄錯誤
    LogError("Non-retryable error")
    return false
}
```

### 2. 計算重試時間？（優先使用 Retry-After）
```go
func CalculateNextRetry(retryCount int, retryAfter *int) time.Time {
    if retryAfter != nil {
        // 🎯 優先使用 429 的 Retry-After header
        return time.Now().Add(time.Duration(*retryAfter) * time.Second).Add(1 * time.Second) // +1s buffer
    }
    
    // 使用指數退避
    backoff := math.Pow(2, float64(retryCount)) * time.Second
    if backoff > 5*time.Minute {
        backoff = 5 * time.Minute // 最大 5 分鐘
    }
    
    // 加入隨機抖動 (±10%)
    jitter := backoff * 0.1 * (rand.Float64()*2 - 1)
    return time.Now().Add(backoff + jitter)
}
```

### 3. 何時移至 Dead Letter？
```go
func ProcessRetry(notification *Notification) {
    result := SendToTeams(notification)
    
    if result.Success {
        RemoveFromRedis(notification.ID)
        RemoveFromDB(notification.ID)
        return
    }
    
    notification.RetryCount++
    
    if notification.RetryCount >= MaxRetries { // 5
        // ❌ 達到最大重試次數
        MoveToDLQ(notification)
        RemoveFromRedis(notification.ID)
        RemoveFromDB(notification.ID)
        AlertAdmin(notification)
    } else {
        // 📅 重新入隊
        nextRetry := CalculateNextRetry(notification.RetryCount, result.RetryAfter)
        RequeueToRedis(notification, nextRetry)
        UpdateDB(notification)
    }
}
```

### 4. 何時觸發熔斷？（使用 Redis Retry-After）
```go
func OnFailure(err error, retryAfter *int) {
    redis.Incr("cb:failures")
    failures := redis.Get("cb:failures")
    
    if failures >= CircuitBreakerThreshold { // 5
        // 🔴 觸發熔斷
        openDuration := 60 * time.Second // 預設 60 秒
        
        if retryAfter != nil {
            // 🎯 優先使用最後一次 429 的 Retry-After
            redis.Set("cb:last_429_retry_after", *retryAfter)
            openDuration = time.Duration(*retryAfter) * time.Second
        } else if lastRetryAfter := redis.Get("cb:last_429_retry_after"); lastRetryAfter != nil {
            // 使用先前記錄的 Retry-After
            openDuration = time.Duration(*lastRetryAfter) * time.Second
        }
        
        redis.Set("cb:state", "open")
        redis.Set("cb:open_until", time.Now().Add(openDuration).Unix())
        
        LogInfo("Circuit breaker opened for %v (from Retry-After)", openDuration)
    }
}
```

### 5. Redis 資料結構
```redis
# Queue: Sorted Set (按時間排序)
ZADD failed_queue {next_retry_timestamp} {notification_json}
ZPOPMIN failed_queue 5  # 原子操作，多節點安全

# Circuit Breaker: String Keys
SET cb:state "closed|open|half-open"
SET cb:failures 0
SET cb:last_429_retry_after 120  # 最後一次 429 的 Retry-After
SET cb:open_until {unix_timestamp}
SET cb:half_open_tests 0
SET cb:success_count 1000
SET cb:failed_count 25
```

---

---

## 🏗️ 技術特性總結

### ✅ 已實現的關鍵特性

1. **多節點水平擴展**
   - 使用 Redis Sorted Set 實現分散式隊列
   - ZPOPMIN 原子操作避免重複處理
   - 每個節點獨立運行 Worker，負載自動平衡

2. **智能退避策略**
   - 🎯 **優先使用 429 Retry-After header**
   - 無 header 時使用指數退避 (1s, 2s, 4s, 8s, 16s)
   - 加入隨機抖動 (±10%) 避免雪崩

3. **Redis 共享熔斷器**
   - 所有節點共享熔斷器狀態
   - 觸發熔斷時優先使用 Retry-After 作為 open_until
   - 支援手動重置熔斷器

4. **Dead Letter Queue**
   - 達到 max_retries (5次) 自動移至 dead_notifications
   - 支援手動重試失敗通知
   - 保留完整的失敗歷史供分析

5. **雙重持久化**
   - Redis: 高效能隊列，按時間排序
   - PostgreSQL: 資料備份，防止 Redis 資料遺失

### 📊 系統容量

| 指標 | 值 | 說明 |
|------|-----|------|
| 節點數量 | 無限制 | 可水平擴展 |
| 每節點吞吐 | ~30 通知/分鐘 | 5 個/批次 × 6 次/分鐘 |
| 隊列容量 | > 100萬 | Redis Sorted Set |
| 重試延遲 | 0-300 秒 | 根據 Retry-After 動態調整 |
| 熔斷恢復 | 60-300 秒 | 根據 Retry-After 動態調整 |

---

**建立日期**: 2025-10-01  
**版本**: v2.0 (Redis + Multi-Node)  
**作者**: TeamsNotifyGoV2 Team  
**相關文檔**:
- [QUEUE_CIRCUIT_BREAKER.md](./QUEUE_CIRCUIT_BREAKER.md)
- [CHANGES_SUMMARY.md](../CHANGES_SUMMARY.md)
- [README.md](../README.md)

