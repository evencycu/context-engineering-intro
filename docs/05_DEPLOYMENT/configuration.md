# 配置說明

本文件詳細說明 Teams Notification API 的配置選項和環境變數設定。

## 最新更新 (v1.0.0)

### 🔧 **配置管理改進**
- **統一配置系統**: 所有配置現在通過 `internal/config/config.go` 統一管理
- **環境變數驗證**: 啟動時自動驗證必要的環境變數
- **配置端點**: 新增 `/config` 和 `/config/validate` API 端點
- **監控指標**: 新增 `/metrics` 端點提供系統監控數據

### 🚨 **重要變更**
- **環境變數驗證**: 服務啟動時會驗證所有必要的環境變數
- **錯誤處理**: 標準化的錯誤響應格式
- **日誌管理**: 結構化日誌輸出，支援多種格式
- **Token 快取**: Redis 支援的 Token 快取機制

## 環境變數

### 必要環境變數（支援 TN_ 前綴）

#### 資料庫配置
```bash
# PostgreSQL 連線字串
export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
# 或使用 TN_ 前綴（優先於同名無前綴變數）
export TN_DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
```

#### Teams Bot 配置
```bash
# Teams Bot 應用程式 ID
export TEAMS_BOT_APP_ID="844146d7-4ac9-4e4d-a463-d6e027714e81"

# Teams 租戶 ID
export TEAMS_TENANT_ID="051cece0-e4dc-4aed-b471-bf29824e1ee6"

# Teams Bot 密鑰
export TEAMS_BOT_APP_PASSWORD="HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy"
```

#### Redis 配置
```bash
# Redis 連線字串
export REDIS_URL="redis://localhost:6379"
# 或
export TN_REDIS_URL="redis://localhost:6379"
```

### 可選環境變數

#### 服務配置
```bash
# 服務端口 (預設: 8080)
export PORT="8080"

# 環境模式 (development/production)
export ENVIRONMENT="development"

# JWT 密鑰
export JWT_SECRET="your-jwt-secret-key"
```

#### Actor Pool 配置
```bash
# Actor 池大小 (預設: 10)
export ACTOR_POOL_SIZE="10"

# 佇列檢查間隔 (預設: 1s)
export QUEUE_CHECK_INTERVAL="1s"

# 重試延遲 (預設: 5s)
export RETRY_DELAY="5s"
```

#### 電路斷路器配置
```bash
# 失敗閾值 (預設: 5)
export CIRCUIT_BREAKER_THRESHOLD="5"

# 恢復超時 (預設: 30s)
export CIRCUIT_BREAKER_RECOVERY_TIMEOUT="30s"

# 半開狀態最大呼叫數 (預設: 3)
export CIRCUIT_BREAKER_HALF_OPEN_MAX_CALLS="3"
```

#### 速率限制配置
```bash
# 全域速率限制 (預設: 100 req/s)
export RATE_LIMIT_GLOBAL="100"

# 突發限制 (預設: 200)
export RATE_LIMIT_BURST="200"

# 讀取超時 (預設: 30s)
export READ_TIMEOUT="30s"

# 寫入超時 (預設: 30s)
export WRITE_TIMEOUT="30s"
```

## 配置驗證（Viper）

### 自動驗證
服務啟動時會自動驗證以下配置：

```bash
# 必要環境變數檢查
TEAMS_BOT_APP_ID          # Teams Bot 應用程式 ID
TEAMS_BOT_APP_PASSWORD    # Teams Bot 密鑰
TEAMS_TENANT_ID          # Teams 租戶 ID
DATABASE_URL             # 數據庫連接字串
REDIS_URL                # Redis 連接字串
```

### 手動驗證
使用 API 端點驗證配置：

```bash
# 檢查配置狀態（非敏感值）
curl http://localhost:8080/config

# 驗證配置（回傳缺失項目）
curl -X POST http://localhost:8080/config/validate
```

### 監控端點

```bash
# 系統健康檢查
curl http://localhost:8080/health

# 系統指標
curl http://localhost:8080/metrics
```

## 配置檔案

### YAML 配置檔案

#### `configs/apiserver.yaml`
```yaml
server:
  port: 8080
  environment: development
  read_timeout: 30s
  write_timeout: 30s

database:
  url: "postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

redis:
  url: "redis://localhost:6379"
  max_retries: 3
  timeout: 5s
  pool_size: 10

teams:
  bot_app_id: "844146d7-4ac9-4e4d-a463-d6e027714e81"
  tenant_id: "051cece0-e4dc-4aed-b471-bf29824e1ee6"
  bot_app_password: "HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy"

auth:
  jwt_secret: "your-jwt-secret-key"
  token_expiry: 24h

rate_limit:
  global: 100
  burst: 200

actor_pool:
  max_actors: 10
  queue_check_interval: 1s
  retry_delay: 5s

circuit_breaker:
  failure_threshold: 5
  recovery_timeout: 30s
  half_open_max_calls: 3
```

#### `configs/common.yaml`
```yaml
logging:
  level: info
  format: json
  output: stdout

monitoring:
  health_check_interval: 30s
  metrics_enabled: true

security:
  cors_enabled: true
  cors_origins:
    - "http://localhost:3000"
    - "http://localhost:8080"
```

### 環境變數檔案

#### `.env` 檔案
```bash
# 資料庫配置
DATABASE_URL=postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable

# Teams Bot 配置
TEAMS_BOT_APP_ID=844146d7-4ac9-4e4d-a463-d6e027714e81
TEAMS_TENANT_ID=051cece0-e4dc-4aed-b471-bf29824e1ee6
TEAMS_BOT_APP_PASSWORD=HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy

# Redis 配置
REDIS_URL=redis://localhost:6379

# 服務配置
PORT=8080
ENVIRONMENT=development
JWT_SECRET=your-jwt-secret-key

# Actor Pool 配置
ACTOR_POOL_SIZE=10
QUEUE_CHECK_INTERVAL=1s
RETRY_DELAY=5s

# 電路斷路器配置
CIRCUIT_BREAKER_THRESHOLD=5
CIRCUIT_BREAKER_RECOVERY_TIMEOUT=30s
CIRCUIT_BREAKER_HALF_OPEN_MAX_CALLS=3

# 速率限制配置
RATE_LIMIT_GLOBAL=100
RATE_LIMIT_BURST=200
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
```

## Docker 配置

### Docker Compose 配置

#### `docker-compose.yml`
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: notification_center
      POSTGRES_USER: teamsnotify
      POSTGRES_PASSWORD: teamsnotify123
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/migrations:/docker-entrypoint-initdb.d

  redis:
    image: redis:6-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  apiserver:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://teamsnotify:teamsnotify123@postgres:5432/notification_center?sslmode=disable
      - REDIS_URL=redis://redis:6379
      - TEAMS_BOT_APP_ID=844146d7-4ac9-4e4d-a463-d6e027714e81
      - TEAMS_TENANT_ID=051cece0-e4dc-4aed-b471-bf29824e1ee6
      - TEAMS_BOT_APP_PASSWORD=HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy
    depends_on:
      - postgres
      - redis

volumes:
  postgres_data:
  redis_data:
```

### Dockerfile 配置

#### `Dockerfile`
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./server"]
```

## Kubernetes 配置

### 部署配置

#### `deployments/apiserver/deployment.yaml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: teams-notification-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: teams-notification-api
  template:
    metadata:
      labels:
        app: teams-notification-api
    spec:
      containers:
      - name: apiserver
        image: teams-notification-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: redis-secret
              key: url
        - name: TEAMS_BOT_APP_ID
          valueFrom:
            secretKeyRef:
              name: teams-secret
              key: app-id
        - name: TEAMS_TENANT_ID
          valueFrom:
            secretKeyRef:
              name: teams-secret
              key: tenant-id
        - name: TEAMS_BOT_APP_PASSWORD
          valueFrom:
            secretKeyRef:
              name: teams-secret
              key: app-password
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

### 服務配置

#### `deployments/apiserver/service.yaml`
```yaml
apiVersion: v1
kind: Service
metadata:
  name: teams-notification-api-service
spec:
  selector:
    app: teams-notification-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

## 配置驗證

### 檢查配置
```bash
# 檢查環境變數
env | grep -E "(DATABASE_URL|REDIS_URL|TEAMS_)"

# 檢查服務配置
curl http://localhost:8080/health

# 檢查資料庫連線
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT 1;"

# 檢查 Redis 連線
docker exec teamsnotify-redis redis-cli ping
```

### 配置測試
```bash
# 測試 API 端點
curl -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "test-key",
    "message": "Configuration test",
    "messageType": "text",
    "priority": "normal",
    "targets": ["all"]
  }'
```

## 最佳實踐

### 1. 安全性
- 使用環境變數儲存敏感資訊
- 定期輪換密鑰和憑證
- 限制網路存取權限

### 2. 效能
- 根據負載調整 Actor Pool 大小
- 監控資料庫連線池
- 設定適當的快取策略

### 3. 監控
- 啟用健康檢查端點
- 設定日誌級別
- 監控資源使用情況

### 4. 備份
- 定期備份資料庫
- 備份配置檔案
- 測試恢復程序
