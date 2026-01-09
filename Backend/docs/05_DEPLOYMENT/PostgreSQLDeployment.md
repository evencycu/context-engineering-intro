# 🐘 PostgreSQL 部署文件 - TeamsNotifyGoV2

## 📋 概述

本文件詳細說明 TeamsNotifyGoV2 專案中 PostgreSQL 的部署配置，包括資料庫服務、持久化儲存和備份策略。

---

## 🎯 PostgreSQL 架構

### 功能模組
```
┌─────────────────────────────────────────────────────────────┐
│                    PostgreSQL 服務                         │
├─────────────────┬─────────────────┬─────────────────────────┤
│   資料庫服務    │   持久化儲存    │     備份策略            │
│                 │                 │                         │
│ • 主要資料庫    │ • 資料持久化    │ • 自動備份              │
│ • 連線池管理    │ • 事務日誌      │ • 增量備份              │
│ • 查詢優化      │ • WAL 日誌      │ • 災難恢復              │
└─────────────────┴─────────────────┴─────────────────────────┘
```

---

## 🐳 Docker 部署

### Docker Compose 配置
```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: teamsnotify-postgres
    environment:
      POSTGRES_DB: teamsnotify
      POSTGRES_USER: teamsnotify
      POSTGRES_PASSWORD: teamsnotify123
      POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --lc-collate=C --lc-ctype=C"
      POSTGRES_HOST_AUTH_METHOD: md5
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./configs/postgresql.conf:/etc/postgresql/postgresql.conf
      - ./configs/pg_hba.conf:/etc/postgresql/pg_hba.conf
      - ./scripts/database/schema.sql:/docker-entrypoint-initdb.d/01-schema.sql
      - ./scripts/database/init.sql:/docker-entrypoint-initdb.d/02-init.sql
      - ./scripts/database/migrations:/docker-entrypoint-initdb.d/migrations
      - ./scripts/postgres/init.sh:/docker-entrypoint-initdb.d/00-init.sh
    command: postgres -c config_file=/etc/postgresql/postgresql.conf
    networks:
      - teamsnotify-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U teamsnotify -d teamsnotify"]
      interval: 10s
      timeout: 5s
      retries: 5
    deploy:
      resources:
        limits:
          memory: 1G
          cpus: '1.0'
        reservations:
          memory: 512M
          cpus: '0.5'

volumes:
  postgres_data:
    driver: local
```

### PostgreSQL 配置文件
```conf
# configs/postgresql.conf
# 基本配置
listen_addresses = '*'
port = 5432
max_connections = 100
superuser_reserved_connections = 3

# 記憶體配置
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 4MB
maintenance_work_mem = 64MB
dynamic_shared_memory_type = posix

# 檢查點配置
checkpoint_completion_target = 0.9
checkpoint_timeout = 5min
checkpoint_warning = 30s
max_wal_size = 1GB
min_wal_size = 80MB

# WAL 配置
wal_level = replica
wal_compression = on
wal_buffers = 16MB
wal_writer_delay = 200ms
commit_delay = 0
commit_siblings = 5

# 複製配置
max_wal_senders = 3
max_replication_slots = 3
hot_standby = on
hot_standby_feedback = on

# 日誌配置
log_destination = 'stderr'
logging_collector = on
log_directory = 'pg_log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_rotation_age = 1d
log_rotation_size = 100MB
log_min_duration_statement = 1000
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
log_checkpoints = on
log_connections = on
log_disconnections = on
log_lock_waits = on
log_temp_files = 0
log_autovacuum_min_duration = 0
log_error_verbosity = default

# 統計配置
track_activities = on
track_counts = on
track_io_timing = on
track_functions = all
stats_temp_directory = 'pg_stat_tmp'

# 自動清理配置
autovacuum = on
autovacuum_max_workers = 3
autovacuum_naptime = 1min
autovacuum_vacuum_threshold = 50
autovacuum_analyze_threshold = 50
autovacuum_vacuum_scale_factor = 0.2
autovacuum_analyze_scale_factor = 0.1
autovacuum_freeze_max_age = 200000000
autovacuum_multixact_freeze_max_age = 400000000
autovacuum_vacuum_cost_delay = 20ms
autovacuum_vacuum_cost_limit = 200

# 安全配置
ssl = on
ssl_cert_file = 'server.crt'
ssl_key_file = 'server.key'
ssl_ca_file = 'ca.crt'
ssl_ciphers = 'HIGH:MEDIUM:+3DES:!aNULL'
ssl_prefer_server_ciphers = on
password_encryption = scram-sha-256

# 時區配置
timezone = 'Asia/Taipei'
log_timezone = 'Asia/Taipei'

# 其他配置
default_text_search_config = 'pg_catalog.english'
escape_string_warning = on
standard_conforming_strings = on
```

### pg_hba.conf 配置
```conf
# configs/pg_hba.conf
# PostgreSQL Client Authentication Configuration File

# TYPE  DATABASE        USER            ADDRESS                 METHOD

# "local" is for Unix domain socket connections only
local   all             all                                     trust

# IPv4 local connections:
host    all             all             127.0.0.1/32            md5
host    all             all             0.0.0.0/0               md5

# IPv6 local connections:
host    all             all             ::1/128                 md5

# Allow replication connections from localhost, by a user with the
# replication privilege.
local   replication     all                                     trust
host    replication     all             127.0.0.1/32            md5
host    replication     all             ::1/128                 md5
```

### PostgreSQL 初始化腳本
```bash
#!/bin/bash
# scripts/postgres/init.sh

# 等待 PostgreSQL 啟動
sleep 10

# 建立必要的資料庫
psql -U teamsnotify -d teamsnotify -c "CREATE DATABASE IF NOT EXISTS teamsnotify_test;"
psql -U teamsnotify -d teamsnotify -c "CREATE DATABASE IF NOT EXISTS teamsnotify_backup;"

# 建立必要的使用者
psql -U teamsnotify -d teamsnotify -c "CREATE USER IF NOT EXISTS teamsnotify_readonly WITH PASSWORD 'readonly123';"
psql -U teamsnotify -d teamsnotify -c "CREATE USER IF NOT EXISTS teamsnotify_backup WITH PASSWORD 'backup123';"

# 設定權限
psql -U teamsnotify -d teamsnotify -c "GRANT SELECT ON ALL TABLES IN SCHEMA public TO teamsnotify_readonly;"
psql -U teamsnotify -d teamsnotify -c "GRANT SELECT ON ALL TABLES IN SCHEMA public TO teamsnotify_backup;"

# 建立索引
psql -U teamsnotify -d teamsnotify -c "CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);"
psql -U teamsnotify -d teamsnotify -c "CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);"
psql -U teamsnotify -d teamsnotify -c "CREATE INDEX IF NOT EXISTS idx_notifications_priority ON notifications(priority);"

echo "PostgreSQL 初始化完成"
```

---

## ☸️ Kubernetes 部署

### PostgreSQL Deployment
```yaml
# deployments/postgres/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: teams-notification-postgres
  namespace: teams-notification
  labels:
    app: teams-notification-postgres
    version: v1.0.0
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: teams-notification-postgres
  template:
    metadata:
      labels:
        app: teams-notification-postgres
        version: v1.0.0
    spec:
      containers:
        - name: postgres
          image: postgres:15-alpine
          imagePullPolicy: IfNotPresent
          ports:
            - name: postgres
              containerPort: 5432
              protocol: TCP
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
            - name: POSTGRES_INITDB_ARGS
              value: "--encoding=UTF-8 --lc-collate=C --lc-ctype=C"
            - name: POSTGRES_HOST_AUTH_METHOD
              value: "md5"
          command:
            - postgres
            - -c
            - config_file=/etc/postgresql/postgresql.conf
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "1Gi"
              cpu: "500m"
          livenessProbe:
            exec:
              command:
                - pg_isready
                - -U
                - teamsnotify
                - -d
                - teamsnotify
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            exec:
              command:
                - pg_isready
                - -U
                - teamsnotify
                - -d
                - teamsnotify
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 3
          volumeMounts:
            - name: postgres-data
              mountPath: /var/lib/postgresql/data
            - name: postgres-config
              mountPath: /etc/postgresql
              readOnly: true
            - name: postgres-init
              mountPath: /docker-entrypoint-initdb.d
              readOnly: true
      volumes:
        - name: postgres-data
          persistentVolumeClaim:
            claimName: postgres-pvc
        - name: postgres-config
          configMap:
            name: postgres-config
        - name: postgres-init
          configMap:
            name: postgres-init
      nodeSelector:
        kubernetes.io/os: linux
```

### PostgreSQL Service
```yaml
# deployments/postgres/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: teams-notification-postgres
  namespace: teams-notification
  labels:
    app: teams-notification-postgres
spec:
  type: ClusterIP
  ports:
    - name: postgres
      port: 5432
      targetPort: 5432
      protocol: TCP
  selector:
    app: teams-notification-postgres
```

### PostgreSQL ConfigMap
```yaml
# deployments/postgres/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgres-config
  namespace: teams-notification
data:
  postgresql.conf: |
    # 基本配置
    listen_addresses = '*'
    port = 5432
    max_connections = 100
    superuser_reserved_connections = 3

    # 記憶體配置
    shared_buffers = 256MB
    effective_cache_size = 1GB
    work_mem = 4MB
    maintenance_work_mem = 64MB
    dynamic_shared_memory_type = posix

    # 檢查點配置
    checkpoint_completion_target = 0.9
    checkpoint_timeout = 5min
    checkpoint_warning = 30s
    max_wal_size = 1GB
    min_wal_size = 80MB

    # WAL 配置
    wal_level = replica
    wal_compression = on
    wal_buffers = 16MB
    wal_writer_delay = 200ms
    commit_delay = 0
    commit_siblings = 5

    # 複製配置
    max_wal_senders = 3
    max_replication_slots = 3
    hot_standby = on
    hot_standby_feedback = on

    # 日誌配置
    log_destination = 'stderr'
    logging_collector = on
    log_directory = 'pg_log'
    log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
    log_rotation_age = 1d
    log_rotation_size = 100MB
    log_min_duration_statement = 1000
    log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
    log_checkpoints = on
    log_connections = on
    log_disconnections = on
    log_lock_waits = on
    log_temp_files = 0
    log_autovacuum_min_duration = 0
    log_error_verbosity = default

    # 統計配置
    track_activities = on
    track_counts = on
    track_io_timing = on
    track_functions = all
    stats_temp_directory = 'pg_stat_tmp'

    # 自動清理配置
    autovacuum = on
    autovacuum_max_workers = 3
    autovacuum_naptime = 1min
    autovacuum_vacuum_threshold = 50
    autovacuum_analyze_threshold = 50
    autovacuum_vacuum_scale_factor = 0.2
    autovacuum_analyze_scale_factor = 0.1
    autovacuum_freeze_max_age = 200000000
    autovacuum_multixact_freeze_max_age = 400000000
    autovacuum_vacuum_cost_delay = 20ms
    autovacuum_vacuum_cost_limit = 200

    # 安全配置
    ssl = on
    ssl_cert_file = 'server.crt'
    ssl_key_file = 'server.key'
    ssl_ca_file = 'ca.crt'
    ssl_ciphers = 'HIGH:MEDIUM:+3DES:!aNULL'
    ssl_prefer_server_ciphers = on
    password_encryption = scram-sha-256

    # 時區配置
    timezone = 'Asia/Taipei'
    log_timezone = 'Asia/Taipei'

    # 其他配置
    default_text_search_config = 'pg_catalog.english'
    escape_string_warning = on
    standard_conforming_strings = on

  pg_hba.conf: |
    # PostgreSQL Client Authentication Configuration File
    # TYPE  DATABASE        USER            ADDRESS                 METHOD
    local   all             all                                     trust
    host    all             all             127.0.0.1/32            md5
    host    all             all             0.0.0.0/0               md5
    host    all             all             ::1/128                 md5
    local   replication     all                                     trust
    host    replication     all             127.0.0.1/32            md5
    host    replication     all             ::1/128                 md5
```

### PostgreSQL PVC
```yaml
# deployments/postgres/pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: teams-notification
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
  storageClassName: standard
```

---

## 🔧 PostgreSQL 配置參數

### 資料庫連線配置
```bash
# 基本連線配置
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

# 連線超時配置
DB_CONNECT_TIMEOUT=30s
DB_QUERY_TIMEOUT=60s
DB_STATEMENT_TIMEOUT=300s
```

### 備份配置
```bash
# 備份配置
BACKUP_ENABLED=true
BACKUP_SCHEDULE="0 2 * * *"  # 每天凌晨2點
BACKUP_RETENTION_DAYS=30
BACKUP_STORAGE_PATH=/backups
BACKUP_COMPRESSION=true
BACKUP_ENCRYPTION=true

# 增量備份配置
INCREMENTAL_BACKUP=true
INCREMENTAL_BACKUP_SCHEDULE="0 */6 * * *"  # 每6小時
INCREMENTAL_BACKUP_RETENTION_DAYS=7
```

### 監控配置
```bash
# 監控配置
MONITORING_ENABLED=true
MONITORING_INTERVAL=30s
MONITORING_METRICS=true
MONITORING_QUERIES=true
MONITORING_CONNECTIONS=true
MONITORING_LOCKS=true
```

---

## 📊 PostgreSQL 監控配置

### 監控指標
```yaml
# deployments/postgres/monitoring.yaml
apiVersion: v1
kind: ServiceMonitor
metadata:
  name: teams-notification-postgres
  namespace: teams-notification
  labels:
    app: teams-notification-postgres
spec:
  selector:
    matchLabels:
      app: teams-notification-postgres
  endpoints:
    - port: postgres
      interval: 30s
      path: /metrics
```

### 監控查詢
```bash
# 資料庫大小
psql -U teamsnotify -d teamsnotify -c "SELECT pg_size_pretty(pg_database_size('teamsnotify'));"

# 連線數
psql -U teamsnotify -d teamsnotify -c "SELECT count(*) FROM pg_stat_activity;"

# 慢查詢
psql -U teamsnotify -d teamsnotify -c "SELECT query, mean_time, calls FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"

# 鎖等待
psql -U teamsnotify -d teamsnotify -c "SELECT * FROM pg_locks WHERE NOT granted;"

# 資料庫統計
psql -U teamsnotify -d teamsnotify -c "SELECT * FROM pg_stat_database WHERE datname = 'teamsnotify';"
```

---

## 🚀 部署流程

### 1. Docker 部署
```bash
# 啟動 PostgreSQL 服務
docker-compose up -d postgres

# 檢查 PostgreSQL 狀態
docker-compose ps postgres

# 檢查 PostgreSQL 日誌
docker-compose logs postgres

# 連接到 PostgreSQL
docker-compose exec postgres psql -U teamsnotify -d teamsnotify
```

### 2. Kubernetes 部署
```bash
# 建立命名空間
kubectl create namespace teams-notification

# 建立 Secrets
kubectl create secret generic teams-notification-secrets \
  --from-literal=postgres-password="teamsnotify123"

# 建立 ConfigMap
kubectl create configmap postgres-config \
  --from-file=configs/postgresql.conf \
  --from-file=configs/pg_hba.conf

# 建立 PVC
kubectl apply -f deployments/postgres/pvc.yaml

# 部署 PostgreSQL
kubectl apply -f deployments/postgres/

# 檢查部署狀態
kubectl get pods -n teams-notification -l app=teams-notification-postgres
```

### 3. 驗證部署
```bash
# 檢查 Pod 狀態
kubectl get pods -n teams-notification

# 檢查服務狀態
kubectl get services -n teams-notification

# 檢查 PostgreSQL 連線
kubectl exec -it deployment/teams-notification-postgres -n teams-notification -- psql -U teamsnotify -d teamsnotify -c "SELECT version();"

# 檢查資料庫大小
kubectl exec -it deployment/teams-notification-postgres -n teams-notification -- psql -U teamsnotify -d teamsnotify -c "SELECT pg_size_pretty(pg_database_size('teamsnotify'));"
```

---

## 🔧 維護和監控

### 備份策略
```bash
# 建立備份腳本
#!/bin/bash
# scripts/postgres/backup.sh

BACKUP_DIR="/backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="postgres_backup_${DATE}.sql"

# 建立備份目錄
mkdir -p $BACKUP_DIR

# 執行備份
pg_dump -U teamsnotify -h localhost -d teamsnotify > $BACKUP_DIR/$BACKUP_FILE

# 壓縮備份文件
gzip $BACKUP_DIR/$BACKUP_FILE

# 清理舊備份 (保留30天)
find $BACKUP_DIR -name "postgres_backup_*.sql.gz" -mtime +30 -delete

echo "PostgreSQL 備份完成: $BACKUP_FILE.gz"
```

### 監控腳本
```bash
#!/bin/bash
# scripts/postgres/monitor.sh

# 檢查 PostgreSQL 狀態
pg_isready -U teamsnotify -d teamsnotify

# 檢查資料庫大小
DB_SIZE=$(psql -U teamsnotify -d teamsnotify -t -c "SELECT pg_size_pretty(pg_database_size('teamsnotify'));" | tr -d ' ')
echo "資料庫大小: $DB_SIZE"

# 檢查連線數
CONNECTIONS=$(psql -U teamsnotify -d teamsnotify -t -c "SELECT count(*) FROM pg_stat_activity;" | tr -d ' ')
echo "連線數: $CONNECTIONS"

# 檢查慢查詢
SLOW_QUERIES=$(psql -U teamsnotify -d teamsnotify -t -c "SELECT count(*) FROM pg_stat_statements WHERE mean_time > 1000;" | tr -d ' ')
echo "慢查詢數: $SLOW_QUERIES"
```

---

## 📚 相關文檔

- [部署組件文件](./DeploymentComponents.md)
- [App本身部署文件](./AppDeployment.md)
- [Redis 部署文件](./RedisDeployment.md)
- [PostgreSQL 配置文件](../configs/postgresql.conf)
- [資料庫 Schema](../scripts/database/schema.sql)

---

**建立時間**: 2025-10-09  
**維護者**: AI Assistant  
**狀態**: ✅ 完成  
**組件**: PostgreSQL 資料庫服務
