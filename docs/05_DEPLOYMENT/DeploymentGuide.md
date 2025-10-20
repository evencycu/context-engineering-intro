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
curl -X POST http://api.company.com/api/v1/notify \
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

## 14. 部署模式比較

### 14.1 本地進程部署 (Local Process)

**優點**:
- 最快啟動速度，無需 Docker 映像構建
- 最佳調試體驗，可直接使用 IDE 斷點
- 最低資源使用，應用程式直接運行在主機

**缺點**:
- 環境隔離性差，可能受主機環境影響
- 需要手動管理依賴
- 數據可能不持久

**使用場景**: 快速開發和深度調試

### 14.2 本地 Docker 部署 (Local Docker)

**優點**:
- 環境隔離性好，避免環境衝突
- 確保開發環境與生產環境一致性
- 快速重置，每次啟動都是乾淨環境

**缺點**:
- 啟動速度較慢，需要構建映像
- 數據不持久，容器停止後數據丟失
- 調試相對複雜

**使用場景**: 功能測試和集成測試

### 14.3 持久化 Docker 部署 (Persistent Docker)

**優點**:
- 數據持久化，適合長期開發
- 保持 Docker 環境隔離優勢
- 更接近生產環境的部署方式

**缺點**:
- 啟動速度中等
- 相對本地進程部署佔用更多資源
- 調試相對複雜

**使用場景**: 長期功能開發和數據依賴測試

### 14.4 部署模式選擇建議

| 開發階段 | 推薦模式 | 原因 |
|---------|---------|------|
| 快速迭代 | 本地進程 | 最快反饋和最佳調試體驗 |
| 功能測試 | 本地 Docker | 環境隔離和快速重置 |
| 長期開發 | 持久化 Docker | 數據保留和真實環境模擬 |

## 15. 常見部署問題與解決方案

### 15.1 端口衝突問題

**問題**: 多個部署模式同時運行造成端口衝突

**解決方案**:
```bash
# 檢查端口使用情況
lsof -i :8080
lsof -i :5432
lsof -i :6379

# 停止所有相關容器
docker stop $(docker ps -q --filter "name=teamsnotify")

# 清理容器
docker rm $(docker ps -aq --filter "name=teamsnotify")
```

### 15.2 數據庫連接問題

**問題**: 本地進程模式無法連接到數據庫

**解決方案**:
```bash
# 確保數據庫容器運行
docker-compose -f scripts/local-test/docker-compose-local.yml up -d postgres redis

# 檢查數據庫連接
docker exec teamsnotify-postgres-local pg_isready -U teamsnotify

# 測試連接
docker exec teamsnotify-postgres-local psql -U teamsnotify -d notification_center -c "SELECT 1;"
```

### 15.3 密碼認證失敗

**問題**: 持久化模式密碼認證失敗

**解決方案**:
```bash
# 檢查環境變數
echo $DATABASE_URL
echo $REDIS_URL

# 重新設定正確的連接字符串
export DATABASE_URL="postgresql://teamsnotify:teamsnotify_password@localhost:5432/notification_center?sslmode=disable"
export REDIS_URL="redis://localhost:6379/0"
```

### 15.4 容器重複問題

**問題**: 同時存在多組容器造成資源浪費

**解決方案**:
```bash
# 清理所有相關容器
docker stop teamsnotify-postgres-persistent teamsnotify-redis-persistent teamsnotify-postgres-local teamsnotify-redis-local
docker rm teamsnotify-postgres-persistent teamsnotify-redis-persistent teamsnotify-postgres-local teamsnotify-redis-local

# 確認清理完成
docker ps -a
```

## 16. 智能部署腳本使用指南

### 16.1 部署腳本概覽

| 腳本 | 用途 | 模式 |
|------|------|------|
| `scripts/smart-deploy.sh` | 智能部署主腳本 | 支持三種模式 |
| `scripts/dev-deploy.sh` | 開發部署腳本 | 本地 Docker |
| `scripts/quick-redeploy.sh` | 快速重構腳本 | 本地 Docker |

### 16.2 使用示例

```bash
# 本地進程部署
./scripts/smart-deploy.sh local-process

# 本地 Docker 部署 (清理模式)
./scripts/smart-deploy.sh local-docker --clean --migrate --test

# 持久化 Docker 部署
./scripts/smart-deploy.sh persistent-docker --clean --migrate

# 快速重構 (僅 API Server)
./scripts/quick-redeploy.sh
```

### 16.3 部署後驗證

```bash
# 健康檢查
curl http://localhost:8080/health

# 系統指標
curl http://localhost:8080/api/v1/metrics

# 配置驗證
curl -X POST http://localhost:8080/api/v1/config/validate

# 隊列狀態
curl http://localhost:8080/api/v1/queue/stats
```

---

**版本**: v1.1
**最後更新**: 2025-10-11
**作者**: TeamsNotify DevOps Team

---

## 🔧 故障排除指南

### 常見問題與解決方案

#### 1. Actor Pool 問題
**症狀**: "Actor pool full" 錯誤，通知處理緩慢
**原因**: Actor 沒有正確釋放，導致池子飽和
**解決方案**:
```bash
# 檢查 Actor Pool 狀態
curl http://localhost:8080/api/v1/metrics | jq '.actor_pool'

# 重啟服務
docker-compose restart api-server
```

#### 2. Teams API 認證失敗
**症狀**: "401 Unauthorized" 錯誤
**原因**: Teams API 憑證配置錯誤
**解決方案**:
```bash
# 檢查環境變數
docker exec teamsnotify-api-server env | grep TEAMS

# 更新憑證
docker exec teamsnotify-api-server sh -c 'export TEAMS_BOT_APP_ID="your-app-id"'
```

#### 3. 數據庫連接問題
**症狀**: "Failed to connect to database" 錯誤
**原因**: 數據庫連接字符串或密碼錯誤
**解決方案**:
```bash
# 檢查數據庫狀態
docker exec teamsnotify-postgres pg_isready

# 檢查連接字符串
echo $DATABASE_URL
```

#### 4. Redis 連接問題
**症狀**: "Redis connection failed" 錯誤
**原因**: Redis 服務未啟動或連接配置錯誤
**解決方案**:
```bash
# 檢查 Redis 狀態
docker exec teamsnotify-redis redis-cli ping

# 重啟 Redis
docker-compose restart redis
```

### 監控和診斷

#### 系統健康檢查
```bash
# 檢查整體健康狀態
curl http://localhost:8080/api/v1/monitoring/health

# 檢查業務指標
curl http://localhost:8080/api/v1/monitoring/business

# 檢查性能指標
curl http://localhost:8080/api/v1/monitoring/performance
```

#### 日誌分析
```bash
# 查看 API 服務日誌
docker logs teamsnotify-api-server --tail 100

# 查看錯誤日誌
docker logs teamsnotify-api-server 2>&1 | grep ERROR

# 查看 Actor Pool 日誌
docker logs teamsnotify-api-server 2>&1 | grep "Actor"
```

### 性能優化建議

#### 1. Actor Pool 調優
- 根據負載調整 `max_actors` 參數
- 監控 Actor 完成率和處理時間
- 避免 Actor 積壓

#### 2. 數據庫優化
- 定期清理舊的 notification_destinations
- 優化數據庫索引
- 監控連接池使用情況

#### 3. Redis 優化
- 設置適當的 TTL
- 監控記憶體使用
- 配置持久化策略

