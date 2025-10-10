# 🧪 本地測試指南 - TeamsNotifyGoV2

## 📋 概述

本目錄包含 TeamsNotifyGoV2 專案的本地測試腳本和配置，支援兩種測試方式：
1. **Docker Compose** - 簡單的容器化測試
2. **Kubernetes (minikube)** - 完整的 Kubernetes 環境測試

---

## 🎯 測試方式

### 1️⃣ Docker Compose 測試 (推薦)

#### 快速開始
```bash
# 啟動本地 Docker 環境
./scripts/local-test/start-docker.sh

# 測試 API
./scripts/local-test/test-docker.sh

# 停止環境
docker-compose -f scripts/local-test/docker-compose-local.yml down
```

#### 服務 URL
- **API Server**: http://localhost:8080

#### 管理命令
```bash
# 查看日誌
docker-compose -f scripts/local-test/docker-compose-local.yml logs -f

# 重啟服務
docker-compose -f scripts/local-test/docker-compose-local.yml restart

# 停止服務
docker-compose -f scripts/local-test/docker-compose-local.yml down

# 清理資料
docker-compose -f scripts/local-test/docker-compose-local.yml down -v
```

### 2️⃣ Kubernetes (minikube) 測試

#### 環境準備
```bash
# 建立本地 Kubernetes 環境
./scripts/local-test/setup-local-env.sh

# 部署 UT 環境
./scripts/local-test/deploy-ut.sh

# 測試 API
./scripts/local-test/test-api.sh

# 清理環境
./scripts/local-test/cleanup.sh
```

#### 取得服務 URL
```bash
# 取得 API Server URL
minikube service teams-notification-api-server -n teams-notification-ut --url

# 開啟 Kubernetes 儀表板
minikube dashboard
```

---

## 🔧 環境配置

### 必要工具
- **Docker**: 容器化運行環境
- **kubectl**: Kubernetes 命令行工具
- **minikube**: 本地 Kubernetes 集群
- **jq**: JSON 處理工具

### 環境變數
```bash
# Teams 配置 (可選)
export TEAMS_BOT_APP_ID="your-app-id"
export TEAMS_TENANT_ID="your-tenant-id"
export TEAMS_BOT_APP_PASSWORD="your-app-password"
```

---

## 📊 測試內容

### API 測試
- **健康檢查**: `/health`
- **指標端點**: `/metrics`
- **外部 API**: `/api/v1/external/notify`
- **配置端點**: `/config`

### 資料庫測試
- **連線測試**: PostgreSQL 連線
- **資料驗證**: 資料庫查詢
- **遷移測試**: 資料庫結構

### Redis 測試
- **連線測試**: Redis 連線
- **快取測試**: 快取功能
- **佇列測試**: 佇列功能

---

## 🚀 快速測試流程

### Docker Compose 流程
```bash
# 1. 啟動環境
./scripts/local-test/start-docker.sh

# 2. 測試 API
./scripts/local-test/test-docker.sh

# 3. 查看日誌
docker-compose -f scripts/local-test/docker-compose-local.yml logs -f

# 4. 停止環境
docker-compose -f scripts/local-test/docker-compose-local.yml down
```

### Kubernetes 流程
```bash
# 1. 建立環境
./scripts/local-test/setup-local-env.sh

# 2. 部署應用
./scripts/local-test/deploy-ut.sh

# 3. 測試 API
./scripts/local-test/test-api.sh

# 4. 清理環境
./scripts/local-test/cleanup.sh
```

---

## 🔧 故障排除

### 常見問題

#### 1. Docker 相關
```bash
# 檢查 Docker 狀態
docker info

# 檢查容器狀態
docker-compose -f scripts/local-test/docker-compose-local.yml ps

# 查看容器日誌
docker-compose -f scripts/local-test/docker-compose-local.yml logs api-server
```

#### 2. Kubernetes 相關
```bash
# 檢查 minikube 狀態
minikube status

# 檢查 Pod 狀態
kubectl get pods -n teams-notification-ut

# 查看 Pod 日誌
kubectl logs -f deployment/teams-notification-api-server -n teams-notification-ut
```

#### 3. 網路問題
```bash
# 檢查服務狀態
kubectl get services -n teams-notification-ut

# 檢查 Ingress
kubectl get ingress -n teams-notification-ut

# 測試連線
curl -v http://localhost:8080/health
```

---

## 📚 相關文檔

- [部署組件文件](../docs/05_DEPLOYMENT/DeploymentComponents.md)
- [App本身部署文件](../docs/05_DEPLOYMENT/AppDeployment.md)
- [Redis 部署文件](../docs/05_DEPLOYMENT/RedisDeployment.md)
- [PostgreSQL 部署文件](../docs/05_DEPLOYMENT/PostgreSQLDeployment.md)
- [部署整合文件](../docs/05_DEPLOYMENT/DeploymentIntegration.md)

---

**建立時間**: 2025-10-09  
**維護者**: AI Assistant  
**狀態**: ✅ 完成  
**測試方式**: Docker Compose, Kubernetes (minikube)
