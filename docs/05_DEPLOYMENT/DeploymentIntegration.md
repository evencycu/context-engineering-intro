# 🔄 部署整合文件 - TeamsNotifyGoV2

## 📋 概述

本文件詳細說明 TeamsNotifyGoV2 專案中三個主要組件（App本身、Redis、PostgreSQL）的整合部署流程，包括部署順序、依賴關係和驗證步驟。

---

## 🎯 整合部署架構

### 組件關係圖
```
┌─────────────────────────────────────────────────────────────┐
│                    TeamsNotifyGoV2 部署                    │
├─────────────────┬─────────────────┬─────────────────────────┤
│   PostgreSQL    │     Redis       │      App 組件           │
│                 │                 │                         │
│ • 資料庫服務    │ • 快取服務      │ • API Server            │
│ • 持久化儲存    │ • 佇列服務      │                         │
│ • 備份策略      │ • 會話管理      │                         │
└─────────────────┴─────────────────┴─────────────────────────┘
```

---

## 📋 部署順序和依賴關係

### 1. 基礎設施層 (Infrastructure Layer)
```
PostgreSQL → Redis → App 組件
```

#### 部署順序
1. **PostgreSQL** - 資料庫服務
2. **Redis** - 快取和佇列服務
3. **App 組件** - 應用程式服務

#### 依賴關係
- **App 組件** 依賴於 **PostgreSQL** 和 **Redis**
- **Redis** 可以獨立部署
- **PostgreSQL** 必須最先部署

---

## 🐳 Docker Compose 整合部署

### 完整 Docker Compose 配置
```yaml
# docker-compose.yml
version: '3.8'

services:
  # PostgreSQL 資料庫
  postgres:
    image: postgres:15-alpine
    container_name: teamsnotify-postgres
    environment:
      POSTGRES_DB: teamsnotify
      POSTGRES_USER: teamsnotify
      POSTGRES_PASSWORD: teamsnotify123
      POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --lc-collate=C --lc-ctype=C"
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/database/schema.sql:/docker-entrypoint-initdb.d/01-schema.sql
      - ./scripts/database/init.sql:/docker-entrypoint-initdb.d/02-init.sql
    networks:
      - teamsnotify-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U teamsnotify -d teamsnotify"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis 快取和佇列
  redis:
    image: redis:7-alpine
    container_name: teamsnotify-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
      - ./configs/redis.conf:/usr/local/etc/redis/redis.conf
    command: redis-server /usr/local/etc/redis/redis.conf
    networks:
      - teamsnotify-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # API Server
  api-server:
    build: .
    container_name: teamsnotify-api-server
    ports:
      - "8080:8080"
      - "9090:9090"
    environment:
      - ENVIRONMENT=production
      - LOG_LEVEL=info
      - PORT=8080
      - DATABASE_URL=postgresql://teamsnotify:teamsnotify123@postgres:5432/teamsnotify
      - REDIS_URL=redis://redis:6379
      - TEAMS_BOT_APP_ID=${TEAMS_BOT_APP_ID}
      - TEAMS_TENANT_ID=${TEAMS_TENANT_ID}
      - TEAMS_BOT_APP_PASSWORD=${TEAMS_BOT_APP_PASSWORD}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - teamsnotify-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # Worker
  worker:
    build: .
    container_name: teamsnotify-worker
    command: ["./server", "worker"]
    environment:
      - ENVIRONMENT=production
      - LOG_LEVEL=info
      - DATABASE_URL=postgresql://teamsnotify:teamsnotify123@postgres:5432/teamsnotify
      - REDIS_URL=redis://redis:6379
      - TEAMS_BOT_APP_ID=${TEAMS_BOT_APP_ID}
      - TEAMS_TENANT_ID=${TEAMS_TENANT_ID}
      - TEAMS_BOT_APP_PASSWORD=${TEAMS_BOT_APP_PASSWORD}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - teamsnotify-network
    restart: unless-stopped

  # Admin
  admin:
    build: .
    container_name: teamsnotify-admin
    command: ["./server", "admin"]
    ports:
      - "8081:8080"
    environment:
      - ENVIRONMENT=production
      - LOG_LEVEL=info
      - PORT=8080
      - DATABASE_URL=postgresql://teamsnotify:teamsnotify123@postgres:5432/teamsnotify
      - REDIS_URL=redis://redis:6379
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - teamsnotify-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  # Adminer (資料庫管理工具)
  adminer:
    image: adminer:4.8.1
    container_name: teamsnotify-adminer
    ports:
      - "8082:8080"
    networks:
      - teamsnotify-network
    restart: unless-stopped

networks:
  teamsnotify-network:
    driver: bridge

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local
```

---

## ☸️ Kubernetes 整合部署

### 部署順序腳本
```bash
#!/bin/bash
# scripts/deployment/deploy-all.sh

set -e

echo "🚀 開始部署 TeamsNotifyGoV2 到 Kubernetes..."

# 1. 建立命名空間
echo "📁 建立命名空間..."
kubectl create namespace teams-notification --dry-run=client -o yaml | kubectl apply -f -

# 2. 建立 Secrets
echo "🔐 建立 Secrets..."
kubectl create secret generic teams-notification-secrets \
  --from-literal=database-url="postgresql://teamsnotify:teamsnotify123@teams-notification-postgres:5432/teamsnotify" \
  --from-literal=redis-url="redis://teams-notification-redis:6379" \
  --from-literal=postgres-password="teamsnotify123" \
  --from-literal=redis-password="teamsnotify123" \
  --from-literal=teams-bot-app-id="${TEAMS_BOT_APP_ID}" \
  --from-literal=teams-tenant-id="${TEAMS_TENANT_ID}" \
  --from-literal=teams-bot-app-password="${TEAMS_BOT_APP_PASSWORD}" \
  --namespace=teams-notification \
  --dry-run=client -o yaml | kubectl apply -f -

# 3. 建立 ConfigMaps
echo "📋 建立 ConfigMaps..."
kubectl create configmap postgres-config \
  --from-file=configs/postgresql.conf \
  --from-file=configs/pg_hba.conf \
  --namespace=teams-notification \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create configmap redis-config \
  --from-file=configs/redis.conf \
  --namespace=teams-notification \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create configmap teams-notification-config \
  --from-file=configs/ \
  --namespace=teams-notification \
  --dry-run=client -o yaml | kubectl apply -f -

# 4. 部署 PostgreSQL
echo "🐘 部署 PostgreSQL..."
kubectl apply -f deployments/postgres/pvc.yaml
kubectl apply -f deployments/postgres/

# 等待 PostgreSQL 就緒
echo "⏳ 等待 PostgreSQL 就緒..."
kubectl wait --for=condition=ready pod -l app=teams-notification-postgres -n teams-notification --timeout=300s

# 5. 部署 Redis
echo "🔴 部署 Redis..."
kubectl apply -f deployments/redis/pvc.yaml
kubectl apply -f deployments/redis/

# 等待 Redis 就緒
echo "⏳ 等待 Redis 就緒..."
kubectl wait --for=condition=ready pod -l app=teams-notification-redis -n teams-notification --timeout=300s

# 6. 部署 App 組件
echo "🚀 部署 App 組件..."
kubectl apply -f deployments/api-server/
## (worker/admin 已移除)

# 等待 App 組件就緒
echo "⏳ 等待 App 組件就緒..."
kubectl wait --for=condition=ready pod -l app=teams-notification-api-server -n teams-notification --timeout=300s
## (worker/admin 已移除)

echo "✅ 部署完成！"
echo "📊 檢查部署狀態..."
kubectl get pods -n teams-notification
kubectl get services -n teams-notification
```

### 部署驗證腳本
```bash
#!/bin/bash
# scripts/deployment/verify-deployment.sh

echo "🔍 驗證 TeamsNotifyGoV2 部署..."

# 檢查 Pod 狀態
echo "📊 檢查 Pod 狀態..."
kubectl get pods -n teams-notification

# 檢查服務狀態
echo "🌐 檢查服務狀態..."
kubectl get services -n teams-notification

# 檢查 Ingress 狀態
echo "🔗 檢查 Ingress 狀態..."
kubectl get ingress -n teams-notification

# 檢查 PostgreSQL 連線
echo "🐘 檢查 PostgreSQL 連線..."
kubectl exec -it deployment/teams-notification-postgres -n teams-notification -- psql -U teamsnotify -d teamsnotify -c "SELECT version();"

# 檢查 Redis 連線
echo "🔴 檢查 Redis 連線..."
kubectl exec -it deployment/teams-notification-redis -n teams-notification -- redis-cli -a teamsnotify123 ping

# 檢查 API Server 健康狀態
echo "🚀 檢查 API Server 健康狀態..."
kubectl exec -it deployment/teams-notification-api-server -n teams-notification -- wget --no-verbose --tries=1 --spider http://localhost:8080/health

# 檢查 Worker 狀態
echo "⚙️ 檢查 Worker 狀態..."
kubectl logs deployment/teams-notification-worker -n teams-notification --tail=10

# 檢查 Admin 健康狀態
echo "👨‍💼 檢查 Admin 健康狀態..."
kubectl exec -it deployment/teams-notification-admin -n teams-notification -- wget --no-verbose --tries=1 --spider http://localhost:8080/health

echo "✅ 部署驗證完成！"
```

---

## 🔧 環境變數配置

### 統一環境變數
```bash
# 基本配置
ENVIRONMENT=production
LOG_LEVEL=info
PORT=8080

# 資料庫配置
DATABASE_URL=postgresql://teamsnotify:teamsnotify123@postgres:5432/teamsnotify
DB_HOST=postgres
DB_PORT=5432
DB_NAME=teamsnotify
DB_USER=teamsnotify
DB_PASSWORD=teamsnotify123

# Redis 配置
REDIS_URL=redis://redis:6379
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=teamsnotify123

# Teams 配置
TEAMS_BOT_APP_ID=your-app-id
TEAMS_TENANT_ID=your-tenant-id
TEAMS_BOT_APP_PASSWORD=your-app-password

# 安全配置
JWT_SECRET=your-jwt-secret
ENCRYPTION_KEY=your-encryption-key
```

### 環境變數驗證
```bash
#!/bin/bash
# scripts/deployment/validate-env.sh

echo "🔍 驗證環境變數..."

# 檢查必要環境變數
REQUIRED_VARS=(
  "TEAMS_BOT_APP_ID"
  "TEAMS_TENANT_ID"
  "TEAMS_BOT_APP_PASSWORD"
)

for var in "${REQUIRED_VARS[@]}"; do
  if [ -z "${!var}" ]; then
    echo "❌ 缺少必要環境變數: $var"
    exit 1
  else
    echo "✅ $var 已設定"
  fi
done

echo "✅ 環境變數驗證完成！"
```

---

## 📊 監控和日誌

### 統一監控配置
```yaml
# deployments/monitoring/service-monitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: teams-notification-monitor
  namespace: teams-notification
  labels:
    app: teams-notification
spec:
  selector:
    matchLabels:
      app: teams-notification
  endpoints:
    - port: http
      interval: 30s
      path: /metrics
    - port: metrics
      interval: 30s
      path: /metrics
```

### 日誌聚合配置
```yaml
# deployments/logging/fluentd-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
  namespace: teams-notification
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/*teams-notification*.log
      pos_file /var/log/fluentd-containers.log.pos
      tag kubernetes.*
      format json
    </source>
    
    <match kubernetes.**>
      @type elasticsearch
      host elasticsearch.logging.svc.cluster.local
      port 9200
      index_name teams-notification
    </match>
```

---

## 🚀 部署流程

### 1. 準備階段
```bash
# 1. 檢查環境
./scripts/deployment/validate-env.sh

# 2. 建置映像
docker build -t teams-notification/api-server:v1.0.0 .
docker build -t teams-notification/worker:v1.0.0 .
docker build -t teams-notification/admin:v1.0.0 .

# 3. 推送映像
docker push teams-notification/api-server:v1.0.0
docker push teams-notification/worker:v1.0.0
docker push teams-notification/admin:v1.0.0
```

### 2. Docker 部署
```bash
# 啟動所有服務
docker-compose up -d

# 檢查服務狀態
docker-compose ps

# 檢查日誌
docker-compose logs -f
```

### 3. Kubernetes 部署
```bash
# 執行部署腳本
./scripts/deployment/deploy-all.sh

# 驗證部署
./scripts/deployment/verify-deployment.sh
```

### 4. 清理部署
```bash
# Docker 清理
docker-compose down -v

# Kubernetes 清理
kubectl delete namespace teams-notification
```

---

## 🔧 故障排除

### 常見問題
1. **PostgreSQL 連線失敗**
   - 檢查資料庫服務狀態
   - 驗證連線字串
   - 檢查網路連線

2. **Redis 連線失敗**
   - 檢查 Redis 服務狀態
   - 驗證密碼設定
   - 檢查記憶體使用

3. **App 組件啟動失敗**
   - 檢查依賴服務
   - 驗證環境變數
   - 檢查日誌輸出

### 除錯命令
```bash
# 檢查 Pod 日誌
kubectl logs -f deployment/teams-notification-api-server -n teams-notification

# 檢查服務連線
kubectl exec -it deployment/teams-notification-api-server -n teams-notification -- curl http://teams-notification-postgres:5432

# 檢查環境變數
kubectl exec -it deployment/teams-notification-api-server -n teams-notification -- env | grep -E "(DATABASE|REDIS|TEAMS)"
```

---

## 📚 相關文檔

- [部署組件文件](./DeploymentComponents.md)
- [App本身部署文件](./AppDeployment.md)
- [Redis 部署文件](./RedisDeployment.md)
- [PostgreSQL 部署文件](./PostgreSQLDeployment.md)
- [Docker Compose 配置](../docker-compose.yml)
- [Kubernetes 部署文件](../deployments/)

---

**建立時間**: 2025-10-09  
**維護者**: AI Assistant  
**狀態**: ✅ 完成  
**組件**: 整合部署流程
