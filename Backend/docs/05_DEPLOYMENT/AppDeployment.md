# 🚀 App本身部署文件 - TeamsNotifyGoV2

## 📋 概述

本文件詳細說明 TeamsNotifyGoV2 應用程式本身的部署配置，現僅包含 API Server。

---

## 🎯 應用程式架構

### 組件關係圖
```
┌─────────────────────────────────────────────────────────────┐
│                    TeamsNotifyGoV2 App                      │
├─────────────────┬─────────────────┬─────────────────────────┤
│   API Server    │                 │                         │
│                 │                 │                         │
│ • REST API     │ • 背景任務      │ • 管理介面              │
│ • 認證授權     │ • 佇列處理      │ • 監控面板              │
│ • 業務邏輯     │ • 通知發送      │ • 配置管理              │
│ • 資料存取     │ • 錯誤重試      │ • 系統狀態              │
└─────────────────┴─────────────────┴─────────────────────────┘
```

---

## 1️⃣ API Server 部署

### 📋 組件說明
- **主要功能**：提供 REST API 服務
- **端口**：8080 (HTTP), 9090 (Metrics)
- **依賴**：PostgreSQL, Redis

### 🐳 Docker 部署

#### Dockerfile
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 複製 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 複製源碼
COPY . .

# 建置應用程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# 生產環境映像
FROM alpine:latest

# 安裝必要工具
RUN apk --no-cache add ca-certificates tzdata

# 設定工作目錄
WORKDIR /root/

# 複製建置好的二進制文件
COPY --from=builder /app/server .

# 複製配置檔案
COPY --from=builder /app/configs ./configs

# 複製資料庫腳本
COPY --from=builder /app/scripts/database ./scripts/database

# 設定時區
ENV TZ=Asia/Taipei

# 暴露端口
EXPOSE 8080 9090

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 啟動命令
CMD ["./server", "api"]
```

#### Docker Compose 配置
```yaml
# docker-compose.yml
version: '3.8'

services:
  apiserver:
    build: .
    container_name: teamsnotify-apiserver
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
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
        reservations:
          memory: 256M
          cpus: '0.25'
```

### ☸️ Kubernetes 部署

#### Deployment 配置
```yaml
# deployments/apiserver/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: teams-notification-apiserver
  namespace: teams-notification
  labels:
    app: teams-notification-apiserver
    version: v1.0.0
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  selector:
    matchLabels:
      app: teams-notification-apiserver
  template:
    metadata:
      labels:
        app: teams-notification-apiserver
        version: v1.0.0
    spec:
      containers:
        - name: apiserver
          image: teams-notification/apiserver:v1.0.0
          imagePullPolicy: IfNotPresent
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
            - name: metrics
              containerPort: 9090
              protocol: TCP
          env:
            - name: ENVIRONMENT
              value: "production"
            - name: LOG_LEVEL
              value: "info"
            - name: PORT
              value: "8080"
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
            - name: TEAMS_BOT_APP_ID
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: teams-bot-app-id
            - name: TEAMS_TENANT_ID
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: teams-tenant-id
            - name: TEAMS_BOT_APP_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: teams-bot-app-password
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
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3
          volumeMounts:
            - name: config-volume
              mountPath: /root/configs
              readOnly: true
      volumes:
        - name: config-volume
          configMap:
            name: teams-notification-config
      nodeSelector:
        kubernetes.io/os: linux
      tolerations:
        - key: "node-role.kubernetes.io/master"
          operator: "Exists"
          effect: "NoSchedule"
```

#### Service 配置
```yaml
# deployments/apiserver/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: teams-notification-apiserver
  namespace: teams-notification
  labels:
    app: teams-notification-apiserver
spec:
  type: LoadBalancer
  ports:
    - name: http
      port: 80
      targetPort: 8080
      protocol: TCP
    - name: metrics
      port: 9090
      targetPort: 9090
      protocol: TCP
  selector:
    app: teams-notification-apiserver
```

#### Ingress 配置
```yaml
# deployments/apiserver/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: teams-notification-apiserver
  namespace: teams-notification
  annotations:
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/rate-limit-window: "1m"
spec:
  tls:
    - hosts:
        - api.teamsnotify.com
      secretName: teams-notification-tls
  rules:
    - host: api.teamsnotify.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: teams-notification-apiserver
                port:
                  number: 80
```

---

## 2️⃣ Worker 部署（已移除）

### 📋 組件說明
- **主要功能**：背景任務處理
- **端口**：無對外端口
- **依賴**：PostgreSQL, Redis

### 🐳 Docker 部署

#### Docker Compose 配置
```yaml
# docker-compose.yml
services:
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
    deploy:
      resources:
        limits:
          memory: 256M
          cpus: '0.25'
        reservations:
          memory: 128M
          cpus: '0.1'
```

### ☸️ Kubernetes 部署

#### Deployment 配置
```yaml
# deployments/worker/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: teams-notification-worker
  namespace: teams-notification
  labels:
    app: teams-notification-worker
    version: v1.0.0
spec:
  replicas: 2
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  selector:
    matchLabels:
      app: teams-notification-worker
  template:
    metadata:
      labels:
        app: teams-notification-worker
        version: v1.0.0
    spec:
      containers:
        - name: worker
          image: teams-notification/worker:v1.0.0
          imagePullPolicy: IfNotPresent
          command: ["./server", "worker"]
          env:
            - name: ENVIRONMENT
              value: "production"
            - name: LOG_LEVEL
              value: "info"
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
            - name: TEAMS_BOT_APP_ID
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: teams-bot-app-id
            - name: TEAMS_TENANT_ID
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: teams-tenant-id
            - name: TEAMS_BOT_APP_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: teams-bot-app-password
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "250m"
          volumeMounts:
            - name: config-volume
              mountPath: /root/configs
              readOnly: true
      volumes:
        - name: config-volume
          configMap:
            name: teams-notification-config
      nodeSelector:
        kubernetes.io/os: linux
```

---

## 3️⃣ Admin 部署（已移除）

### 📋 組件說明
- **主要功能**：管理介面
- **端口**：8080
- **依賴**：PostgreSQL, Redis

### 🐳 Docker 部署

#### Docker Compose 配置
```yaml
# docker-compose.yml
services:
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
    deploy:
      resources:
        limits:
          memory: 256M
          cpus: '0.25'
        reservations:
          memory: 128M
          cpus: '0.1'
```

### ☸️ Kubernetes 部署

#### Deployment 配置
```yaml
# deployments/admin/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: teams-notification-admin
  namespace: teams-notification
  labels:
    app: teams-notification-admin
    version: v1.0.0
spec:
  replicas: 1
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
  selector:
    matchLabels:
      app: teams-notification-admin
  template:
    metadata:
      labels:
        app: teams-notification-admin
        version: v1.0.0
    spec:
      containers:
        - name: admin
          image: teams-notification/admin:v1.0.0
          imagePullPolicy: IfNotPresent
          command: ["./server", "admin"]
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          env:
            - name: ENVIRONMENT
              value: "production"
            - name: LOG_LEVEL
              value: "info"
            - name: PORT
              value: "8080"
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
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "250m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3
          volumeMounts:
            - name: config-volume
              mountPath: /root/configs
              readOnly: true
      volumes:
        - name: config-volume
          configMap:
            name: teams-notification-config
      nodeSelector:
        kubernetes.io/os: linux
```

#### Service 配置
```yaml
# deployments/admin/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: teams-notification-admin
  namespace: teams-notification
  labels:
    app: teams-notification-admin
spec:
  type: LoadBalancer
  ports:
    - name: http
      port: 80
      targetPort: 8080
      protocol: TCP
  selector:
    app: teams-notification-admin
```

---

## 🔧 環境變數配置

### 必要環境變數
```bash
# 基本配置
ENVIRONMENT=production
LOG_LEVEL=info
PORT=8080

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

# Teams 配置
TEAMS_BOT_APP_ID=your-app-id
TEAMS_TENANT_ID=your-tenant-id
TEAMS_BOT_APP_PASSWORD=your-app-password

# 安全配置
JWT_SECRET=your-jwt-secret
ENCRYPTION_KEY=your-encryption-key
```

### 可選環境變數
```bash
# 監控配置
METRICS_ENABLED=true
METRICS_PORT=9090
PROMETHEUS_ENABLED=true

# 快取配置
CACHE_TTL=3600
CACHE_MAX_SIZE=1000

# 佇列配置
QUEUE_WORKERS=5
QUEUE_RETRY_ATTEMPTS=3
QUEUE_RETRY_DELAY=5s
```

---

## 🚀 部署流程

### 1. 建置映像
```bash
# 建置 API Server
docker build -t teams-notification/apiserver:v1.0.0 .

## (Worker/Admin 已移除)
```

### 2. 推送映像
```bash
# 推送到 Registry
docker push teams-notification/apiserver:v1.0.0
## (Worker/Admin 已移除)
```

### 3. 部署到 Kubernetes
```bash
# 建立命名空間
kubectl create namespace teams-notification

# 建立 Secrets
kubectl create secret generic teams-notification-secrets \
  --from-literal=database-url="postgresql://teamsnotify:password@postgres:5432/teamsnotify" \
  --from-literal=redis-url="redis://redis:6379" \
  --from-literal=teams-bot-app-id="your-app-id" \
  --from-literal=teams-tenant-id="your-tenant-id" \
  --from-literal=teams-bot-app-password="your-app-password"

# 建立 ConfigMap
kubectl create configmap teams-notification-config \
  --from-file=configs/

# 部署組件
kubectl apply -f deployments/apiserver/
## (Worker/Admin 已移除)
```

### 4. 驗證部署
```bash
# 檢查 Pod 狀態
kubectl get pods -n teams-notification

# 檢查服務狀態
kubectl get services -n teams-notification

# 檢查日誌
kubectl logs -f deployment/teams-notification-apiserver -n teams-notification
```

---

## 📚 相關文檔

- [部署組件文件](./DeploymentComponents.md)
- [Redis 部署文件](./RedisDeployment.md)
- [PostgreSQL 部署文件](./PostgreSQLDeployment.md)
- [Docker Compose 配置](../docker-compose.yml)
- [Kubernetes 部署文件](../deployments/)

---

**建立時間**: 2025-10-09  
**維護者**: AI Assistant  
**狀態**: ✅ 完成  
**組件**: API Server, Worker, Admin
