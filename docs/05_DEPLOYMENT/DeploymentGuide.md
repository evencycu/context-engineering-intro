# 🚀 Deployment Guide

## 1. 文件資訊

- **版本**：v1.0
- **作者**：DevOps Engineer
- **最後更新**：2025-10-08

---

## 2. 環境架構


| 環境     | URL / Host          | 用途     | 配置                 |
| -------- | ------------------- | -------- | -------------------- |
| **DEV**  | dev-api.company.com | 開發測試 | Docker Compose       |
| **UAT**  | uat-api.company.com | 驗收測試 | Kubernetes           |
| **PROD** | api.company.com     | 正式環境 | Azure Container Apps |

---

## 3. 主要元件


| 元件           | 技術           | 部署位置              | 配置                  |
| -------------- | -------------- | --------------------- | --------------------- |
| **API Server** | Golang / Gin   | Azure Container Apps  | 3 replicas            |
| **Redis**      | Redis OSS      | Azure Cache for Redis | Cluster mode          |
| **Database**   | PostgreSQL     | Azure Database        | HA with read replicas |
| **Worker**     | Go Worker      | ACA background jobs   | Auto-scaling          |
| **Monitoring** | Loki + Grafana | Azure Monitor         | Centralized logging   |

---

## 4. 部署流程

```mermaid
flowchart LR
    A["Commit to main"] --> B["GitHub Actions Build"]
    B --> C["Run Unit Tests"]
    C --> D["Docker Build & Push (ACR)"]
    D --> E["Deploy to Azure Container Apps"]
    E --> F["Run Smoke Tests"]
    F --> G["Notify Teams Channel"]
```

### 4.1 CI/CD 流程詳述

#### 4.1.1 建置階段

```yaml
# .github/workflows/deploy.yml
name: Deploy to Azure
on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.21
      - name: Run tests
        run: go test ./...
      - name: Build Docker image
        run: docker build -t teams-notification-api:${{ github.sha }} .
      - name: Push to ACR
        run: |
          az acr login --name ${{ secrets.ACR_NAME }}
          docker push ${{ secrets.ACR_NAME }}.azurecr.io/teams-notification-api:${{ github.sha }}
```

#### 4.1.2 部署階段

```yaml
  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Azure Container Apps
        run: |
          az containerapp update \
            --name teams-notification-api \
            --resource-group ${{ secrets.RG_NAME }} \
            --image ${{ secrets.ACR_NAME }}.azurecr.io/teams-notification-api:${{ github.sha }}
```

---

## 5. 環境配置

### 5.1 開發環境 (Docker Compose)

#### 5.1.1 啟動服務

```bash
# 啟動所有服務
docker-compose up -d

# 檢查服務狀態
docker-compose ps

# 查看日誌
docker-compose logs -f api-server
```

#### 5.1.2 環境變數設定

```bash
# 設定 Teams Bot 憑證
export TEAMS_BOT_APP_ID="844146d7-4ac9-4e4d-a463-d6e027714e81"
export TEAMS_TENANT_ID="051cece0-e4dc-4aed-b471-bf29824e1ee6"
export TEAMS_BOT_APP_PASSWORD="HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy"

# 設定資料庫連線
export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
```

### 5.2 測試環境 (Kubernetes)

#### 5.2.1 部署配置

```yaml
# deployments/api-server/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: teams-notification-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: teams-notification-api
  template:
    metadata:
      labels:
        app: teams-notification-api
    spec:
      containers:
      - name: api-server
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
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

#### 5.2.2 服務配置

```yaml
# deployments/api-server/service.yaml
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

### 5.3 生產環境 (Azure Container Apps)

#### 5.3.1 Azure 資源配置

```bash
# 建立資源群組
az group create --name teams-notify-rg --location eastus

# 建立 Container Registry
az acr create --resource-group teams-notify-rg --name teamsnotifyacr --sku Basic

# 建立 Container App Environment
az containerapp env create \
  --name teams-notify-env \
  --resource-group teams-notify-rg \
  --location eastus
```

#### 5.3.2 Container App 部署

```bash
# 部署 API Server
az containerapp create \
  --name teams-notification-api \
  --resource-group teams-notify-rg \
  --environment teams-notify-env \
  --image teamsnotifyacr.azurecr.io/teams-notification-api:latest \
  --target-port 8080 \
  --ingress external \
  --min-replicas 2 \
  --max-replicas 10 \
  --cpu 0.5 \
  --memory 1.0Gi
```

---

## 6. 設定檔管理

### 6.1 環境變數配置

#### 6.1.1 必要環境變數

```bash
# 資料庫配置
DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"

# Teams Bot 配置
TEAMS_BOT_APP_ID="844146d7-4ac9-4e4d-a463-d6e027714e81"
TEAMS_TENANT_ID="051cece0-e4dc-4aed-b471-bf29824e1ee6"
TEAMS_BOT_APP_PASSWORD="HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy"

# Redis 配置
REDIS_URL="redis://localhost:6379"

# 服務配置
PORT="8080"
ENVIRONMENT="production"
JWT_SECRET="your-jwt-secret-key"
```

#### 6.1.2 可選環境變數

```bash
# Actor Pool 配置
ACTOR_POOL_SIZE="10"
QUEUE_CHECK_INTERVAL="1s"
RETRY_DELAY="5s"

# 電路斷路器配置
CIRCUIT_BREAKER_THRESHOLD="5"
CIRCUIT_BREAKER_RECOVERY_TIMEOUT="30s"
CIRCUIT_BREAKER_HALF_OPEN_MAX_CALLS="3"

# 速率限制配置
RATE_LIMIT_GLOBAL="100"
RATE_LIMIT_BURST="200"
READ_TIMEOUT="30s"
WRITE_TIMEOUT="30s"
```

### 6.2 YAML 配置檔案

#### 6.2.1 API Server 配置

```yaml
# configs/api-server.yaml
server:
  port: 8080
  environment: production
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

#### 6.2.2 通用配置

```yaml
# configs/common.yaml
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

## 7. 監控與日誌

### 7.1 健康檢查

```bash
# 系統健康檢查
curl http://localhost:8080/health

# 配置驗證
curl http://localhost:8080/config/validate

# 系統指標
curl http://localhost:8080/metrics
```

### 7.2 日誌管理

```bash
# 查看服務日誌
docker-compose logs -f api-server

# 查看特定日誌
tail -f server.log | grep "ERROR"

# 結構化日誌查詢
jq '.level == "ERROR"' server.log
```

### 7.3 監控指標

- **系統指標**: CPU、內存、磁盤使用率
- **應用指標**: 請求數、響應時間、錯誤率
- **業務指標**: 通知發送量、成功率、佇列長度

---

## 8. 回滾策略

### 8.1 版本管理

- 版本標籤（tag）保留三版
- 自動快照 PostgreSQL 資料
- Redis 可清除 Queue

### 8.2 回滾步驟

```bash
# 1. 停止當前版本
az containerapp stop --name teams-notification-api --resource-group teams-notify-rg

# 2. 部署前一版本
az containerapp update \
  --name teams-notification-api \
  --resource-group teams-notify-rg \
  --image teamsnotifyacr.azurecr.io/teams-notification-api:v1.0.0

# 3. 驗證回滾
curl http://api.company.com/health
```

### 8.3 資料庫回滾

```bash
# 恢復資料庫快照
az postgres flexible-server restore \
  --name teams-notify-db \
  --resource-group teams-notify-rg \
  --source-server teams-notify-db-backup \
  --restore-time 2025-10-08T10:00:00Z
```

---

## 9. 驗證步驟

### 10.1 部署後驗證

1️⃣ 呼叫 `/health` → HTTP 200
2️⃣ 測試發送一筆通知
3️⃣ 監看 Redis 任務執行情況
4️⃣ 檢查 Grafana log 指標 | 無 error spikes

### 10.2 功能驗證

```bash
# 1. 健康檢查
curl -s http://api.company.com/health | jq '.status'

# 2. 發送測試通知
curl -X POST http://api.company.com/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "test-key",
    "message": "Deployment test",
    "message_type": "text",
    "priority": "normal",
    "targets": ["all"]
  }'

# 3. 檢查佇列狀態
curl -s http://api.company.com/api/v1/queue/status | jq '.'

# 4. 檢查熔斷器狀態
curl -s http://api.company.com/api/v1/queue/circuit-breaker/metrics | jq '.state'
```

### 10.3 效能驗證

```bash
# 負載測試
k6 run --vus 100 --duration 30s load-test.js

# 監控指標
curl -s http://api.company.com/metrics | grep "notification_send_total"
```

---

## 11. 故障排除

### 11.1 常見問題

#### 11.1.1 服務無法啟動

```bash
# 檢查端口是否被佔用
lsof -i :8080

# 檢查環境變數
env | grep -E "(DATABASE_URL|REDIS_URL|TEAMS_)"

# 檢查日誌
docker-compose logs api-server
```

#### 11.1.2 資料庫連線失敗

```bash
# 檢查資料庫狀態
docker exec teamsnotify-postgres pg_isready -U teamsnotify

# 測試連線
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT 1;"
```

#### 11.1.3 Redis 連線失敗

```bash
# 檢查 Redis 狀態
docker exec teamsnotify-redis redis-cli ping

# 檢查連線
docker exec teamsnotify-redis redis-cli info
```

#### 11.1.4 Teams 認證失敗

```bash
# 檢查 Bot 憑證
echo $TEAMS_BOT_APP_ID
echo $TEAMS_TENANT_ID

# 測試 Teams API 連線
curl -X GET "https://graph.microsoft.com/v1.0/me" \
  -H "Authorization: Bearer $TOKEN"
```

### 11.2 除錯工具

- **服務日誌**: `/tmp/teamsnotify_server.log`
- **資料庫查詢**: 使用 `psql` 直接查詢
- **Redis 監控**: 使用 `redis-cli` 檢查狀態
- **系統監控**: Grafana 儀表板

---

## 12. 最佳實踐

### 12.1 安全性

- 使用環境變數儲存敏感資訊
- 定期輪換密鑰和憑證
- 限制網路存取權限
- 啟用 HTTPS 和 TLS

### 12.2 效能

- 根據負載調整 Actor Pool 大小
- 監控資料庫連線池
- 設定適當的快取策略
- 使用負載均衡器

### 12.3 監控

- 啟用健康檢查端點
- 設定日誌級別
- 監控資源使用情況
- 建立告警規則

### 12.4 備份

- 定期備份資料庫
- 備份配置檔案
- 測試恢復程序
- 建立災難恢復計劃

---

## 13. 擴展性設計

### 13.1 水平擴展

- 使用 Azure Container Apps 自動擴展
- 配置負載均衡器
- 實現無狀態設計
- 使用外部資料庫和快取

### 13.2 垂直擴展

- 調整 CPU 和記憶體限制
- 優化資料庫查詢
- 使用連接池
- 實現快取策略

---

**版本**: v1.0
**最後更新**: 2025-10-08
**作者**: TeamsNotify DevOps Team
