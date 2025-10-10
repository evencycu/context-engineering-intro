# 🚀 部署組件文件 - TeamsNotifyGoV2

## 📋 概述

本文件詳細說明 TeamsNotifyGoV2 專案的三個主要部署組件：App本身、Redis、PostgreSQL。

---

## 🎯 部署組件架構

### 組件關係圖
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   App本身       │    │     Redis       │    │   PostgreSQL    │
│                 │    │                 │    │                 │
│ • API Server    │◄──►│ • 快取服務      │    │ • 資料庫服務    │
│                 │    │ • 佇列服務      │    │ • 持久化儲存    │
│                 │    │ • 會話管理      │    │ • 資料備份      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

---

## 1️⃣ App本身 - 應用程式部署（僅 API Server）

### 📋 組件說明
- **API Server**：主要的 REST API 服務
- **Worker**：背景任務處理服務
- **Admin**：管理介面服務

### 🐳 Docker 部署

#### Docker Compose 配置
```yaml
# docker-compose.yml
version: '3.8'

services:
  api-server:
    build: .
    container_name: teamsnotify-api-server
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://teamsnotify:teamsnotify123@postgres:5432/teamsnotify
      - REDIS_URL=redis://redis:6379
      - LOG_LEVEL=info
    depends_on:
      - postgres
      - redis
    networks:
      - teamsnotify-network
    restart: unless-stopped

  ## (worker/admin 已移除)
```

#### Dockerfile 配置
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./server"]
```

### ☸️ Kubernetes 部署（僅 API Server）

#### API Server Deployment
```yaml
# deployments/api-server/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: teams-notification-api-server
  namespace: teams-notification
spec:
  replicas: 3
  selector:
    matchLabels:
      app: teams-notification-api-server
  template:
    metadata:
      labels:
        app: teams-notification-api-server
    spec:
      containers:
        - name: api-server
          image: teams-notification/api-server:v1.0.0
          ports:
            - containerPort: 8080
              name: http
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: database-url
            - name: REDIS_URL
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: redis-url
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
```

#### Service 配置
```yaml
# deployments/api-server/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: teams-notification-api-server
  namespace: teams-notification
spec:
  selector:
    app: teams-notification-api-server
  ports:
    - name: http
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

### 🔧 環境變數配置

#### 必要環境變數
```bash
# 資料庫配置
DATABASE_URL=postgresql://username:password@host:port/database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=teamsnotify
DB_USER=teamsnotify
DB_PASSWORD=password

# Redis 配置
REDIS_URL=redis://host:port
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# 應用配置
LOG_LEVEL=info
ENVIRONMENT=production
PORT=8080

# Teams 配置
TEAMS_BOT_APP_ID=your-app-id
TEAMS_TENANT_ID=your-tenant-id
TEAMS_BOT_APP_PASSWORD=your-app-password
```

---

## 2️⃣ Redis - 快取和佇列服務

### 📋 組件說明
- **快取服務**：Token 快取、會話快取
- **佇列服務**：通知佇列、任務佇列
- **會話管理**：用戶會話、Bot 會話

### 🐳 Docker 部署

#### Redis 配置
```yaml
# docker-compose.yml
services:
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
```

#### Redis 配置文件
```conf
# configs/redis.conf
# 基本配置
port 6379
bind 0.0.0.0
protected-mode no

# 持久化配置
save 900 1
save 300 10
save 60 10000

# 記憶體配置
maxmemory 512mb
maxmemory-policy allkeys-lru

# 日誌配置
loglevel notice
logfile ""

# 安全配置
requirepass your-redis-password
```

### ☸️ Kubernetes 部署

#### Redis Deployment
```yaml
# deployments/redis/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: teams-notification-redis
  namespace: teams-notification
spec:
  replicas: 1
  selector:
    matchLabels:
      app: teams-notification-redis
  template:
    metadata:
      labels:
        app: teams-notification-redis
    spec:
      containers:
        - name: redis
          image: redis:7-alpine
          ports:
            - containerPort: 6379
              name: redis
          env:
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: redis-password
          command:
            - redis-server
            - --requirepass
            - $(REDIS_PASSWORD)
          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "200m"
          volumeMounts:
            - name: redis-data
              mountPath: /data
      volumes:
        - name: redis-data
          persistentVolumeClaim:
            claimName: redis-pvc
```

#### Redis Service
```yaml
# deployments/redis/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: teams-notification-redis
  namespace: teams-notification
spec:
  selector:
    app: teams-notification-redis
  ports:
    - name: redis
      port: 6379
      targetPort: 6379
  type: ClusterIP
```

### 🔧 Redis 配置參數

#### 快取配置
```bash
# Token 快取配置
TOKEN_CACHE_TTL=3600          # Token 快取時間 (秒)
TOKEN_CACHE_PREFIX=token:     # Token 快取前綴
TOKEN_CACHE_MAX_SIZE=1000     # 最大快取數量

# 會話快取配置
SESSION_CACHE_TTL=1800        # 會話快取時間 (秒)
SESSION_CACHE_PREFIX=session: # 會話快取前綴
```

#### 佇列配置
```bash
# 通知佇列配置
NOTIFICATION_QUEUE_NAME=notifications
NOTIFICATION_QUEUE_PRIORITY=high
NOTIFICATION_QUEUE_RETRY=3

# 任務佇列配置
TASK_QUEUE_NAME=tasks
TASK_QUEUE_PRIORITY=normal
TASK_QUEUE_RETRY=5
```

---

## 3️⃣ PostgreSQL - 資料庫服務

### 📋 組件說明
- **資料庫服務**：主要資料儲存
- **持久化儲存**：資料備份和恢復
- **連線池管理**：資料庫連線優化

### 🐳 Docker 部署

#### PostgreSQL 配置
```yaml
# docker-compose.yml
services:
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
```

#### PostgreSQL 配置文件
```conf
# configs/postgresql.conf
# 基本配置
listen_addresses = '*'
port = 5432
max_connections = 100

# 記憶體配置
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 4MB
maintenance_work_mem = 64MB

# 日誌配置
log_destination = 'stderr'
logging_collector = on
log_directory = 'pg_log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_rotation_age = 1d
log_rotation_size = 100MB

# 安全配置
ssl = on
ssl_cert_file = 'server.crt'
ssl_key_file = 'server.key'
```

### ☸️ Kubernetes 部署

#### PostgreSQL Deployment
```yaml
# deployments/postgres/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: teams-notification-postgres
  namespace: teams-notification
spec:
  replicas: 1
  selector:
    matchLabels:
      app: teams-notification-postgres
  template:
    metadata:
      labels:
        app: teams-notification-postgres
    spec:
      containers:
        - name: postgres
          image: postgres:15-alpine
          ports:
            - containerPort: 5432
              name: postgres
          env:
            - name: POSTGRES_DB
              value: "teamsnotify"
            - name: POSTGRES_USER
              value: "teamsnotify"
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: postgres-password
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "1Gi"
              cpu: "500m"
          volumeMounts:
            - name: postgres-data
              mountPath: /var/lib/postgresql/data
            - name: postgres-config
              mountPath: /etc/postgresql
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: postgres-pvc
        - name: postgres-config
          configMap:
            name: postgres-config
```

#### PostgreSQL Service
```yaml
# deployments/postgres/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: teams-notification-postgres
  namespace: teams-notification
spec:
  selector:
    app: teams-notification-postgres
  ports:
    - name: postgres
      port: 5432
      targetPort: 5432
  type: ClusterIP
```

### 🔧 PostgreSQL 配置參數

#### 資料庫配置
```bash
# 資料庫連線配置
DB_HOST=localhost
DB_PORT=5432
DB_NAME=teamsnotify
DB_USER=teamsnotify
DB_PASSWORD=password
DB_SSL_MODE=require

# 連線池配置
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300s
DB_CONN_MAX_IDLE_TIME=60s
```

#### 備份配置
```bash
# 備份配置
BACKUP_ENABLED=true
BACKUP_SCHEDULE="0 2 * * *"  # 每天凌晨2點
BACKUP_RETENTION_DAYS=30
BACKUP_STORAGE_PATH=/backups
```

---

## 🔄 整合部署流程

### 1. 環境準備
```bash
# 1. 建立命名空間
kubectl create namespace teams-notification

# 2. 建立 Secrets
kubectl create secret generic teams-notification-secrets \
  --from-literal=database-url="postgresql://teamsnotify:password@postgres:5432/teamsnotify" \
  --from-literal=redis-url="redis://redis:6379" \
  --from-literal=postgres-password="password" \
  --from-literal=redis-password="password"
```

### 2. 部署順序
```bash
# 1. 部署 PostgreSQL
kubectl apply -f deployments/postgres/

# 2. 部署 Redis
kubectl apply -f deployments/redis/

# 3. 部署 App 組件
kubectl apply -f deployments/api-server/
kubectl apply -f deployments/worker/
kubectl apply -f deployments/admin/
```

### 3. 驗證部署
```bash
# 檢查 Pod 狀態
kubectl get pods -n teams-notification

# 檢查服務狀態
kubectl get services -n teams-notification

# 檢查日誌
kubectl logs -f deployment/teams-notification-api-server -n teams-notification
```

---

## 📚 相關文檔

- [部署指南](./DeploymentGuide.md)
- [Docker Compose 配置](../docker-compose.yml)
- [Kubernetes 部署文件](../deployments/)
- [資料庫 Schema](../scripts/database/schema.sql)
- [Redis 配置](../configs/redis.conf)

---

**建立時間**: 2025-10-09  
**維護者**: AI Assistant  
**狀態**: ✅ 完成  
**組件**: App本身、Redis、PostgreSQL
