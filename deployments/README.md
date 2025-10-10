# 🚀 TeamsNotifyGoV2 部署文件

## 📋 概述

本目錄包含 TeamsNotifyGoV2 專案的三個環境部署文件：UT (測試環境)、UAT (驗收測試環境)、PROD (生產環境)。

---

## 📁 目錄結構

```
deployments/
├── ut/                    # 測試環境
│   ├── namespace.yaml     # 命名空間
│   ├── secrets.yaml       # 機密配置
│   ├── configmap.yaml     # 配置映射
│   ├── postgres.yaml      # PostgreSQL 部署
│   ├── redis.yaml         # Redis 部署
│   ├── api-server.yaml    # API 服務部署
│   └── kustomization.yaml # Kustomize 配置
├── uat/                   # 驗收測試環境
│   ├── namespace.yaml
│   ├── secrets.yaml
│   ├── configmap.yaml
│   ├── postgres.yaml
│   ├── redis.yaml
│   ├── api-server.yaml
│   └── kustomization.yaml
├── prod/                  # 生產環境
│   ├── namespace.yaml
│   ├── secrets.yaml
│   ├── configmap.yaml
│   ├── postgres.yaml
│   ├── redis.yaml
│   ├── api-server.yaml
│   └── kustomization.yaml
└── README.md              # 本說明文件
```

---

## 🎯 環境配置

### UT (測試環境)
- **用途**：單元測試和開發測試
- **資源**：最小資源配置
- **持久化**：使用 emptyDir
- **映像標籤**：ut-latest
- **特點**：啟用調試模式、模擬外部 API

### UAT (驗收測試環境)
- **用途**：用戶驗收測試
- **資源**：中等資源配置
- **持久化**：使用 PVC
- **映像標籤**：uat-latest
- **特點**：接近生產環境配置

### PROD (生產環境)
- **用途**：生產環境
- **資源**：高資源配置
- **持久化**：使用高性能儲存
- **映像標籤**：prod-latest
- **特點**：高可用性、高安全性

---

## 🚀 部署方式

### 1. 使用 kubectl 部署
```bash
# 部署 UT 環境
kubectl apply -f deployments/ut/

# 部署 UAT 環境
kubectl apply -f deployments/uat/

# 部署 PROD 環境
kubectl apply -f deployments/prod/
```

### 2. 使用 Kustomize 部署
```bash
# 部署 UT 環境
kubectl apply -k deployments/ut/

# 部署 UAT 環境
kubectl apply -k deployments/uat/

# 部署 PROD 環境
kubectl apply -k deployments/prod/
```

### 3. 使用 Helm 部署 (未來)
```bash
# 部署 UT 環境
helm install teams-notification-ut ./helm-chart --set environment=ut

# 部署 UAT 環境
helm install teams-notification-uat ./helm-chart --set environment=uat

# 部署 PROD 環境
helm install teams-notification-prod ./helm-chart --set environment=prod
```

---

## 🔧 環境變數配置

### 通用環境變數
```bash
# 基本配置
ENVIRONMENT=ut|uat|prod
LOG_LEVEL=debug|info|warn
PORT=8080

# 資料庫配置
DATABASE_URL=postgresql://username:password@host:port/database
DB_HOST=teams-notification-postgres
DB_PORT=5432
DB_NAME=teamsnotify
DB_USER=teamsnotify
DB_SSL_MODE=disable|prefer|require

# Redis 配置
REDIS_URL=redis://host:port
REDIS_HOST=teams-notification-redis
REDIS_PORT=6379
REDIS_PASSWORD=password

# Teams 配置
TEAMS_BOT_APP_ID=your-app-id
TEAMS_TENANT_ID=your-tenant-id
TEAMS_BOT_APP_PASSWORD=your-app-password
```

### 環境特定配置
```bash
# UT 環境
DEBUG_MODE=true
MOCK_EXTERNAL_APIS=true
TEST_DATA_ENABLED=true

# UAT 環境
DEBUG_MODE=false
MOCK_EXTERNAL_APIS=false
TEST_DATA_ENABLED=false

# PROD 環境
DEBUG_MODE=false
MOCK_EXTERNAL_APIS=false
TEST_DATA_ENABLED=false
```

---

## 📊 資源配置

### UT 環境資源
| 組件 | CPU 請求 | CPU 限制 | 記憶體請求 | 記憶體限制 |
|------|----------|----------|------------|------------|
| API Server | 100m | 200m | 256Mi | 512Mi |
| PostgreSQL | 100m | 200m | 256Mi | 512Mi |
| Redis | 50m | 100m | 128Mi | 256Mi |

### UAT 環境資源
| 組件 | CPU 請求 | CPU 限制 | 記憶體請求 | 記憶體限制 |
|------|----------|----------|------------|------------|
| API Server | 200m | 500m | 512Mi | 1Gi |
| PostgreSQL | 200m | 500m | 512Mi | 1Gi |
| Redis | 100m | 200m | 256Mi | 512Mi |

### PROD 環境資源
| 組件 | CPU 請求 | CPU 限制 | 記憶體請求 | 記憶體限制 |
|------|----------|----------|------------|------------|
| API Server | 500m | 1000m | 1Gi | 2Gi |
| PostgreSQL | 500m | 1000m | 1Gi | 2Gi |
| Redis | 200m | 500m | 512Mi | 1Gi |

---

## 🔐 安全配置

### Secrets 管理
- 使用 Kubernetes Secrets 儲存敏感資訊
- 所有密碼和 API 金鑰都經過 Base64 編碼
- 生產環境應使用外部密鑰管理系統

### 網路安全
- 使用 NetworkPolicy 限制網路存取
- 啟用 TLS 加密
- 使用 Service Mesh 進行流量管理

### 存取控制
- 使用 RBAC 控制存取權限
- 限制管理員存取
- 啟用審計日誌

---

## 📈 監控和日誌

### 監控配置
- 啟用 Prometheus 監控
- 配置 Grafana 儀表板
- 設定告警規則

### 日誌管理
- 使用 Fluentd 收集日誌
- 集中化日誌儲存
- 日誌輪轉和清理

### 健康檢查
- 配置 Liveness 和 Readiness 探針
- 設定健康檢查端點
- 自動故障恢復

---

## 🚀 部署腳本

### 自動化部署腳本
```bash
#!/bin/bash
# scripts/deployment/deploy.sh

ENVIRONMENT=$1
if [ -z "$ENVIRONMENT" ]; then
  echo "Usage: $0 <ut|uat|prod>"
  exit 1
fi

echo "🚀 部署 TeamsNotifyGoV2 到 $ENVIRONMENT 環境..."

# 檢查環境
if [ "$ENVIRONMENT" != "ut" ] && [ "$ENVIRONMENT" != "uat" ] && [ "$ENVIRONMENT" != "prod" ]; then
  echo "❌ 無效的環境: $ENVIRONMENT"
  exit 1
fi

# 部署到指定環境
kubectl apply -k deployments/$ENVIRONMENT/

# 等待部署完成
kubectl wait --for=condition=ready pod -l app=teams-notification-api-server -n teams-notification-$ENVIRONMENT --timeout=300s

echo "✅ 部署完成！"
```

### 驗證部署腳本
```bash
#!/bin/bash
# scripts/deployment/verify.sh

ENVIRONMENT=$1
if [ -z "$ENVIRONMENT" ]; then
  echo "Usage: $0 <ut|uat|prod>"
  exit 1
fi

echo "🔍 驗證 $ENVIRONMENT 環境部署..."

# 檢查 Pod 狀態
kubectl get pods -n teams-notification-$ENVIRONMENT

# 檢查服務狀態
kubectl get services -n teams-notification-$ENVIRONMENT

# 檢查健康狀態
kubectl exec -it deployment/teams-notification-api-server -n teams-notification-$ENVIRONMENT -- curl http://localhost:8080/health

echo "✅ 驗證完成！"
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
**環境**: UT, UAT, PROD
