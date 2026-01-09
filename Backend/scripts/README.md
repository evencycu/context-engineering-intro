# Scripts 目錄結構

這個目錄包含所有專案相關的腳本和配置文件，按功能分類組織。

## 📁 目錄結構

### 🔨 `build/` - 建置工具
- `build.sh` - 建置腳本
- `gen.sh` - 代碼生成腳本
- `lint.sh` - 代碼檢查腳本

### 🚀 `deploy/` - 部署工具
- `docker/` - Docker 部署腳本
  - `deploy.sh` - 統一部署腳本 (支援 dev, local, quick, docker-full, server, openapi 等多種模式)
  - `docker.sh` - Docker 容器管理工具
- `k8s/` - Kubernetes 部署腳本

### 🧪 `test/` - 測試工具
- `e2e/` - 端到端測試
  - `e2e_test.sh` - E2E 測試腳本
  - `run_e2e_tests.sh` - 運行 E2E 測試
  - `setup_test_env.sh` - 設置測試環境
- `api/` - API 測試
  - `api_test_automation.sh` - API 自動化測試
  - `test_api.sh` - API 測試腳本
  - `comprehensive_api_test.sh` - 綜合 API 測試
- `performance/` - 效能測試
  - `load_test.js` - 負載測試
  - `load_test_enhanced.js` - 增強負載測試
  - `performance_benchmark.sh` - 效能基準測試
- `security/` - 安全測試
  - `security_test.sh` - 安全測試腳本
- `integration/` - 整合測試
  - `test-docker.sh` - Docker 整合測試
  - `test_queue.sh` - 佇列整合測試

### 🗄️ `database/` - 資料庫工具
- `schema.sql` - 資料庫結構
- `init.sql` - 初始化腳本
- `create_notification_center_db.sh` - 創建資料庫
- `reset_database.sh` - 重置資料庫
- `migrations/` - 資料庫遷移腳本

### 🐳 `docker/` - Docker 配置
- `docker-compose.yml` - 完整開發環境
- `docker-compose-persistent.yml` - 持久化資料服務

### 🛠️ `utils/` - 工具腳本
- `cleanup.sh` - 清理腳本
- `setup-local-env.sh` - 設置本地環境
- `generate_test_report.sh` - 生成測試報告

## 🚀 使用方法

### 建置
```bash
# 建置專案
./scripts/build/build.sh

# 代碼生成
./scripts/build/gen.sh

# 代碼檢查
./scripts/build/lint.sh
```

### 部署
```bash
# 本地部署
./scripts/deploy/local/dev-deploy.sh
./scripts/deploy/local/smart-deploy.sh local-process

# Docker 部署
./scripts/deploy/docker/deploy.sh
./scripts/deploy/docker/docker.sh
```

### 測試
```bash
# E2E 測試
./scripts/test/e2e/run_e2e_tests.sh e2e

# API 測試
./scripts/test/api/api_test_automation.sh

# 效能測試
./scripts/test/performance/performance_benchmark.sh

# 安全測試
./scripts/test/security/security_test.sh
```

### 資料庫
```bash
# 創建資料庫
./scripts/database/create_notification_center_db.sh

# 重置資料庫
./scripts/database/reset_database.sh
```

## 📋 Makefile 命令

```bash
# 建置
make build

# 部署
make dev-deploy
make deploy-local
make deploy-docker
make deploy-persistent

# 測試
make test-e2e
make test-api
make test-load
make test-report
```

## 🔄 重構歷史

這個目錄結構是從原本的混亂分類重構而來：

- **原本**: `root/`, `local-test/`, `testing/`, `deployment/` 等混亂分類
- **現在**: 按功能分類，清晰明確

重構計畫詳見 `REFACTOR_PLAN.md`。