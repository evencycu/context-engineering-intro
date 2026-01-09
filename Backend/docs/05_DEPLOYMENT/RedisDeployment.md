# 🔴 Redis 部署文件 - TeamsNotifyGoV2

## 📋 概述

本文件詳細說明 TeamsNotifyGoV2 專案中 Redis 的部署配置，包括快取服務、佇列服務和會話管理。

---

## 🎯 Redis 架構

### 功能模組
```
┌─────────────────────────────────────────────────────────────┐
│                        Redis 服務                           │
├─────────────────┬─────────────────┬─────────────────────────┤
│   快取服務      │   佇列服務      │     會話管理            │
│                 │                 │                         │
│ • Token 快取    │ • 通知佇列      │ • 用戶會話              │
│ • 資料快取      │ • 任務佇列      │ • Bot 會話              │
│ • 會話快取      │ • 重試佇列      │ • 認證會話              │
└─────────────────┴─────────────────┴─────────────────────────┘
```

---

## 🐳 Docker 部署

### Docker Compose 配置
```yaml
# docker-compose.yml
version: '3.8'

services:
  redis:
    image: redis:7-alpine
    container_name: teamsnotify-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
      - ./configs/redis.conf:/usr/local/etc/redis/redis.conf
      - ./scripts/redis/init.sh:/docker-entrypoint-initdb.d/init.sh
    command: redis-server /usr/local/etc/redis/redis.conf
    networks:
      - teamsnotify-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
        reservations:
          memory: 256M
          cpus: '0.25'
    environment:
      - REDIS_PASSWORD=teamsnotify123
      - REDIS_MAXMEMORY=512mb
      - REDIS_MAXMEMORY_POLICY=allkeys-lru

volumes:
  redis_data:
    driver: local
```

### Redis 配置文件
```conf
# configs/redis.conf
# 基本配置
port 6379
bind 0.0.0.0
protected-mode no
timeout 0
tcp-keepalive 300

# 安全配置
requirepass teamsnotify123
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command DEBUG ""

# 持久化配置
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data

# AOF 配置
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb

# 記憶體配置
maxmemory 512mb
maxmemory-policy allkeys-lru
maxmemory-samples 5

# 日誌配置
loglevel notice
logfile ""
syslog-enabled no

# 客戶端配置
maxclients 10000
tcp-backlog 511

# 慢查詢配置
slowlog-log-slower-than 10000
slowlog-max-len 128

# 通知配置
notify-keyspace-events "Ex"
```

### Redis 初始化腳本
```bash
#!/bin/bash
# scripts/redis/init.sh

# 等待 Redis 啟動
sleep 5

# 設定 Redis 配置
redis-cli -a teamsnotify123 CONFIG SET notify-keyspace-events "Ex"
redis-cli -a teamsnotify123 CONFIG SET maxmemory-policy "allkeys-lru"

# 建立必要的鍵空間
redis-cli -a teamsnotify123 SET "teamsnotify:version" "1.0.0"
redis-cli -a teamsnotify123 SET "teamsnotify:startup" "$(date -u +%Y-%m-%dT%H:%M:%SZ)"

echo "Redis 初始化完成"
```

---

## ☸️ Kubernetes 部署

### Redis Deployment
```yaml
# deployments/redis/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: teams-notification-redis
  namespace: teams-notification
  labels:
    app: teams-notification-redis
    version: v1.0.0
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: teams-notification-redis
  template:
    metadata:
      labels:
        app: teams-notification-redis
        version: v1.0.0
    spec:
      containers:
        - name: redis
          image: redis:7-alpine
          imagePullPolicy: IfNotPresent
          ports:
            - name: redis
              containerPort: 6379
              protocol: TCP
          env:
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: teams-notification-secrets
                  key: redis-password
            - name: REDIS_MAXMEMORY
              value: "512Mi"
            - name: REDIS_MAXMEMORY_POLICY
              value: "allkeys-lru"
          command:
            - redis-server
            - /usr/local/etc/redis/redis.conf
          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "200m"
          livenessProbe:
            exec:
              command:
                - redis-cli
                - -a
                - $(REDIS_PASSWORD)
                - ping
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            exec:
              command:
                - redis-cli
                - -a
                - $(REDIS_PASSWORD)
                - ping
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3
          volumeMounts:
            - name: redis-data
              mountPath: /data
            - name: redis-config
              mountPath: /usr/local/etc/redis
              readOnly: true
            - name: redis-init
              mountPath: /docker-entrypoint-initdb.d
              readOnly: true
      volumes:
        - name: redis-data
          persistentVolumeClaim:
            claimName: redis-pvc
        - name: redis-config
          configMap:
            name: redis-config
        - name: redis-init
          configMap:
            name: redis-init
      nodeSelector:
        kubernetes.io/os: linux
```

### Redis Service
```yaml
# deployments/redis/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: teams-notification-redis
  namespace: teams-notification
  labels:
    app: teams-notification-redis
spec:
  type: ClusterIP
  ports:
    - name: redis
      port: 6379
      targetPort: 6379
      protocol: TCP
  selector:
    app: teams-notification-redis
```

### Redis ConfigMap
```yaml
# deployments/redis/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: redis-config
  namespace: teams-notification
data:
  redis.conf: |
    # 基本配置
    port 6379
    bind 0.0.0.0
    protected-mode no
    timeout 0
    tcp-keepalive 300

    # 安全配置
    requirepass teamsnotify123
    rename-command FLUSHDB ""
    rename-command FLUSHALL ""
    rename-command DEBUG ""

    # 持久化配置
    save 900 1
    save 300 10
    save 60 10000
    stop-writes-on-bgsave-error yes
    rdbcompression yes
    rdbchecksum yes
    dbfilename dump.rdb
    dir /data

    # AOF 配置
    appendonly yes
    appendfilename "appendonly.aof"
    appendfsync everysec
    no-appendfsync-on-rewrite no
    auto-aof-rewrite-percentage 100
    auto-aof-rewrite-min-size 64mb

    # 記憶體配置
    maxmemory 512mb
    maxmemory-policy allkeys-lru
    maxmemory-samples 5

    # 日誌配置
    loglevel notice
    logfile ""
    syslog-enabled no

    # 客戶端配置
    maxclients 10000
    tcp-backlog 511

    # 慢查詢配置
    slowlog-log-slower-than 10000
    slowlog-max-len 128

    # 通知配置
    notify-keyspace-events "Ex"
```

### Redis PVC
```yaml
# deployments/redis/pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: redis-pvc
  namespace: teams-notification
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: standard
```

---

## 🔧 Redis 配置參數

### 快取配置
```bash
# Token 快取配置
TOKEN_CACHE_TTL=3600          # Token 快取時間 (秒)
TOKEN_CACHE_PREFIX=token:     # Token 快取前綴
TOKEN_CACHE_MAX_SIZE=1000     # 最大快取數量

# 會話快取配置
SESSION_CACHE_TTL=1800        # 會話快取時間 (秒)
SESSION_CACHE_PREFIX=session: # 會話快取前綴
SESSION_CACHE_MAX_SIZE=500    # 最大會話數量

# 資料快取配置
DATA_CACHE_TTL=300            # 資料快取時間 (秒)
DATA_CACHE_PREFIX=data:       # 資料快取前綴
DATA_CACHE_MAX_SIZE=2000     # 最大快取數量
```

### 佇列配置
```bash
# 通知佇列配置
NOTIFICATION_QUEUE_NAME=notifications
NOTIFICATION_QUEUE_PRIORITY=high
NOTIFICATION_QUEUE_RETRY=3
NOTIFICATION_QUEUE_DELAY=0

# 任務佇列配置
TASK_QUEUE_NAME=tasks
TASK_QUEUE_PRIORITY=normal
TASK_QUEUE_RETRY=5
TASK_QUEUE_DELAY=0

# 重試佇列配置
RETRY_QUEUE_NAME=retries
RETRY_QUEUE_PRIORITY=low
RETRY_QUEUE_RETRY=10
RETRY_QUEUE_DELAY=60
```

### 會話管理配置
```bash
# 用戶會話配置
USER_SESSION_TTL=3600         # 用戶會話時間 (秒)
USER_SESSION_PREFIX=user:     # 用戶會話前綴
USER_SESSION_MAX_SIZE=1000    # 最大用戶會話數

# Bot 會話配置
BOT_SESSION_TTL=7200          # Bot 會話時間 (秒)
BOT_SESSION_PREFIX=bot:       # Bot 會話前綴
BOT_SESSION_MAX_SIZE=100      # 最大 Bot 會話數

# 認證會話配置
AUTH_SESSION_TTL=1800         # 認證會話時間 (秒)
AUTH_SESSION_PREFIX=auth:     # 認證會話前綴
AUTH_SESSION_MAX_SIZE=500     # 最大認證會話數
```

---

## 📊 Redis 監控配置

### 監控指標
```yaml
# deployments/redis/monitoring.yaml
apiVersion: v1
kind: ServiceMonitor
metadata:
  name: teams-notification-redis
  namespace: teams-notification
  labels:
    app: teams-notification-redis
spec:
  selector:
    matchLabels:
      app: teams-notification-redis
  endpoints:
    - port: redis
      interval: 30s
      path: /metrics
```

### 監控查詢
```bash
# 記憶體使用率
redis-cli -a teamsnotify123 INFO memory | grep used_memory_human

# 連線數
redis-cli -a teamsnotify123 INFO clients | grep connected_clients

# 鍵空間大小
redis-cli -a teamsnotify123 DBSIZE

# 慢查詢
redis-cli -a teamsnotify123 SLOWLOG GET 10

# 快取命中率
redis-cli -a teamsnotify123 INFO stats | grep keyspace
```

---

## 🚀 部署流程

### 1. Docker 部署
```bash
# 啟動 Redis 服務
docker-compose up -d redis

# 檢查 Redis 狀態
docker-compose ps redis

# 檢查 Redis 日誌
docker-compose logs redis

# 連接到 Redis
docker-compose exec redis redis-cli -a teamsnotify123
```

### 2. Kubernetes 部署
```bash
# 建立命名空間
kubectl create namespace teams-notification

# 建立 Secrets
kubectl create secret generic teams-notification-secrets \
  --from-literal=redis-password="teamsnotify123"

# 建立 ConfigMap
kubectl create configmap redis-config \
  --from-file=configs/redis.conf

# 建立 PVC
kubectl apply -f deployments/redis/pvc.yaml

# 部署 Redis
kubectl apply -f deployments/redis/

# 檢查部署狀態
kubectl get pods -n teams-notification -l app=teams-notification-redis
```

### 3. 驗證部署
```bash
# 檢查 Pod 狀態
kubectl get pods -n teams-notification

# 檢查服務狀態
kubectl get services -n teams-notification

# 檢查 Redis 連線
kubectl exec -it deployment/teams-notification-redis -n teams-notification -- redis-cli -a teamsnotify123 ping

# 檢查 Redis 配置
kubectl exec -it deployment/teams-notification-redis -n teams-notification -- redis-cli -a teamsnotify123 CONFIG GET "*"
```

---

## 🔧 維護和監控

### 備份策略
```bash
# 建立備份腳本
#!/bin/bash
# scripts/redis/backup.sh

BACKUP_DIR="/backups/redis"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="redis_backup_${DATE}.rdb"

# 建立備份目錄
mkdir -p $BACKUP_DIR

# 執行備份
redis-cli -a teamsnotify123 BGSAVE

# 等待備份完成
while [ $(redis-cli -a teamsnotify123 LASTSAVE) -eq $(redis-cli -a teamsnotify123 LASTSAVE) ]; do
  sleep 1
done

# 複製備份文件
cp /data/dump.rdb $BACKUP_DIR/$BACKUP_FILE

# 清理舊備份 (保留7天)
find $BACKUP_DIR -name "redis_backup_*.rdb" -mtime +7 -delete

echo "Redis 備份完成: $BACKUP_FILE"
```

### 監控腳本
```bash
#!/bin/bash
# scripts/redis/monitor.sh

# 檢查 Redis 狀態
redis-cli -a teamsnotify123 ping

# 檢查記憶體使用率
MEMORY_USAGE=$(redis-cli -a teamsnotify123 INFO memory | grep used_memory_human | cut -d: -f2 | tr -d '\r')
echo "記憶體使用率: $MEMORY_USAGE"

# 檢查連線數
CONNECTIONS=$(redis-cli -a teamsnotify123 INFO clients | grep connected_clients | cut -d: -f2 | tr -d '\r')
echo "連線數: $CONNECTIONS"

# 檢查鍵空間大小
KEY_COUNT=$(redis-cli -a teamsnotify123 DBSIZE)
echo "鍵空間大小: $KEY_COUNT"
```

---

## 📚 相關文檔

- [部署組件文件](./DeploymentComponents.md)
- [App本身部署文件](./AppDeployment.md)
- [PostgreSQL 部署文件](./PostgreSQLDeployment.md)
- [Redis 配置文件](../configs/redis.conf)
- [Docker Compose 配置](../docker-compose.yml)

---

**建立時間**: 2025-10-09  
**維護者**: AI Assistant  
**狀態**: ✅ 完成  
**組件**: Redis 快取和佇列服務
