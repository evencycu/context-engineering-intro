# Changelog - Queue & Circuit Breaker Implementation

## 版本 v1.1.0 - 2025-10-01

### 🎉 新增功能

#### 1. 失敗通知重試隊列系統
- ✅ 實現持久化隊列（PostgreSQL）
- ✅ 自動入隊失敗的通知
- ✅ 後台 Worker 處理重試
- ✅ 支援最多 5 次重試（可配置）

#### 2. Circuit Breaker（熔斷器）機制
- ✅ 三種狀態：Closed（正常）、Open（熔斷）、Half-Open（測試）
- ✅ 自動熔斷保護
- ✅ 防止服務雪崩
- ✅ 手動重置功能

#### 3. 智能重試策略
- ✅ 指數退避算法（Exponential Backoff）
- ✅ 隨機抖動（Jitter）避免驚群效應
- ✅ 支援 Teams API 的 Retry-After header
- ✅ 可配置的重試參數

#### 4. 錯誤分類與處理
- ✅ 智能識別可重試錯誤
  - `rate_limit` (429)
  - `timeout` (408)
  - `server_error` (5xx)
  - `network_error`
- ✅ 自動跳過不可重試錯誤（4xx）

#### 5. 監控與管理 API
- ✅ `/api/v1/queue/status` - 查詢隊列狀態
- ✅ `/api/v1/queue/circuit-breaker/metrics` - 熔斷器指標
- ✅ `/api/v1/queue/circuit-breaker/reset` - 重置熔斷器

### 📁 新增文件

#### 核心程式碼
```
internal/queue/
├── circuit_breaker.go      # 熔斷器實現
├── retry_policy.go          # 重試策略
├── failed_notification.go   # 失敗通知模型
├── repository.go            # 資料庫操作
├── manager.go               # 隊列管理器
└── error_handler.go         # 錯誤處理工具

internal/api/handlers/queue/
└── handler.go               # Queue API handlers
```

#### 資料庫遷移
```
scripts/migrations/
└── 008_add_failed_notifications_queue.sql  # 建立 failed_notifications 表
```

#### 測試與文檔
```
scripts/
└── test_queue.sh            # Queue 功能測試腳本

docs/
├── QUEUE_CIRCUIT_BREAKER.md  # 完整技術文檔
├── QUICK_START_QUEUE.md       # 快速開始指南
└── QUEUE_API_EXAMPLES.md      # API 使用範例
```

### 🔄 修改文件

#### 主程式
- `cmd/server/main.go`
  - 初始化 Queue Manager
  - 啟動後台 Workers
  - 註冊 Queue Handler

#### Server 路由
- `internal/api/server.go`
  - 新增 Queue API 路由註冊

#### 文檔更新
- `README.md`
  - 新增 Queue 功能說明
  - 更新功能列表

### 📊 資料庫變更

#### 新增表: `failed_notifications`
```sql
CREATE TABLE failed_notifications (
    id UUID PRIMARY KEY,
    notification_id UUID NOT NULL,
    project_id UUID NOT NULL,
    target_id VARCHAR(500) NOT NULL,
    message TEXT NOT NULL,
    reason VARCHAR(50) NOT NULL,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 5,
    next_retry_at TIMESTAMP WITH TIME ZONE NOT NULL,
    retry_after INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_attempt_at TIMESTAMP WITH TIME ZONE,
    installation_id UUID,
    bot_id UUID,
    metadata JSONB DEFAULT '{}'
);
```

#### 新增索引
- `idx_failed_notifications_next_retry` - 高效查詢待重試通知
- `idx_failed_notifications_notification_id` - 按通知 ID 查詢
- `idx_failed_notifications_project_id` - 按專案查詢
- `idx_failed_notifications_reason` - 按失敗原因統計

#### 新增視圖: `failed_notifications_queue_status`
提供即時隊列狀態統計。

### ⚙️ 配置參數

#### Queue Manager 預設配置
```go
WorkerPool:    3                    // 並發 worker 數量
PollInterval:  10 * time.Second     // 輪詢間隔
BatchSize:     10                   // 批次處理大小
CleanupPeriod: 24 * time.Hour       // 清理週期
```

#### Circuit Breaker 預設配置
```go
MaxFailures:   5                    // 失敗閾值
CBTimeout:     1 * time.Minute      // 熔斷超時
```

#### Retry Policy 預設配置
```go
MaxRetries:     5                   // 最大重試次數
InitialBackoff: 1 * time.Second     // 初始退避
MaxBackoff:     5 * time.Minute     // 最大退避
Multiplier:     2.0                 // 退避倍數
JitterFraction: 0.1                 // 10% 抖動
```

### 🧪 測試方式

#### 1. 執行 Migration
```bash
docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < scripts/migrations/008_add_failed_notifications_queue.sql
```

#### 2. 編譯並啟動
```bash
go build -o server cmd/server/main.go
DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable" ./server
```

#### 3. 測試 Queue API
```bash
./scripts/test_queue.sh
```

#### 4. 手動測試
```bash
# 查詢隊列狀態
curl http://localhost:8080/api/v1/queue/status | jq '.'

# 查詢熔斷器指標
curl http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '.'
```

### 📈 效能指標

#### 重試時間表（預設配置）
| 重試次數 | 基礎退避 | 實際範圍 (±10%) |
|---------|---------|----------------|
| 1       | 1秒     | 0.9s - 1.1s    |
| 2       | 2秒     | 1.8s - 2.2s    |
| 3       | 4秒     | 3.6s - 4.4s    |
| 4       | 8秒     | 7.2s - 8.8s    |
| 5       | 16秒    | 14.4s - 17.6s  |

#### 熔斷器行為
- **觸發條件**: 連續 5 次失敗
- **熔斷時長**: 1 分鐘
- **恢復測試**: 3 個請求
- **狀態轉換**: Closed → Open → Half-Open → Closed

### 🔍 監控重點

#### 需要關注的指標
1. **隊列積壓**: `total_pending` > 100 需要關注
2. **熔斷器狀態**: `circuit_state` = "open" 超過 5 分鐘需要告警
3. **失敗率**: `failed_requests / total_requests` > 10% 需要調查
4. **重試成功率**: 追蹤最終成功送達的比例

#### 日誌關鍵字
```
Queue manager started successfully
Worker X: Processing N notifications
Notification X retry succeeded
Notification X retry failed
Circuit breaker is open
```

### 🚀 後續改進計劃

- [ ] 實現通知優先級隊列
- [ ] 整合 Prometheus metrics
- [ ] 實現 Dead Letter Queue
- [ ] 支援自定義重試策略
- [ ] 實現 Grafana 監控面板
- [ ] 新增告警規則配置
- [ ] 優化批次處理效能
- [ ] 實現跨節點協調（分散式場景）

### 📚 相關文檔

- [完整技術文檔](./docs/QUEUE_CIRCUIT_BREAKER.md)
- [快速開始指南](./docs/QUICK_START_QUEUE.md)
- [API 使用範例](./docs/QUEUE_API_EXAMPLES.md)
- [主 README](./README.md)

### 🙏 致謝

本功能參考了以下最佳實踐：
- Microsoft Teams Bot Framework Rate Limiting Guidelines
- Martin Fowler's Circuit Breaker Pattern
- AWS Exponential Backoff and Jitter Strategy

---

**實現日期**: 2025-10-01  
**版本**: v1.1.0  
**狀態**: ✅ 已完成，準備測試

