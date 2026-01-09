# Docker 配置檔案

這個目錄包含所有 Docker Compose 配置檔案。

## 檔案說明

- `docker-compose.yml` - 完整開發環境 (包含 API Server)
- `docker-compose-persistent.yml` - 持久化資料服務 (僅 PostgreSQL + Redis)

## 使用方法

### 完整環境
```bash
# 啟動完整開發環境
docker-compose -f scripts/docker/docker-compose.yml up -d

# 停止完整開發環境
docker-compose -f scripts/docker/docker-compose.yml down
```

### 持久化資料服務
```bash
# 啟動持久化資料服務
docker-compose -f scripts/docker/docker-compose-persistent.yml up -d

# 停止持久化資料服務
docker-compose -f scripts/docker/docker-compose-persistent.yml down
```

## 向後兼容

根目錄的 `docker-compose.yml` 和 `docker-compose-persistent.yml` 檔案會重定向到這個目錄的實際配置檔案。
