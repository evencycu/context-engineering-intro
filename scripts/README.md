# 📁 Scripts 目錄整理

## 📋 概述

本目錄包含 TeamsNotifyGoV2 專案的所有腳本，按功能分類組織。

---

## 📊 腳本統計

| 分類 | 數量 | 說明 |
|------|------|------|
| **資料庫相關** | 2 個 | 資料庫管理和重置腳本 |
| **開發相關** | 3 個 | 開發工具和代碼生成腳本 |
| **測試相關** | 11 個 | 測試執行和報告生成腳本 |
| **部署相關** | 4 個 | 部署和服務管理腳本 |
| **開發部署** | 2 個 | 開發環境快速部署腳本 |
| **總計** | 22 個 | 所有腳本 |

---

## 📁 目錄結構

```
scripts/
├── database/           # 資料庫相關腳本
├── development/       # 開發相關腳本
├── testing/          # 測試相關腳本
├── deployment/       # 部署相關腳本
├── dev-deploy.sh     # 開發部署腳本 (新增)
├── quick-redeploy.sh # 快速重構腳本 (新增)
└── README.md         # 本說明文檔
```

---

## 🗄️ 資料庫相關腳本 (database/)

### 腳本列表
| 腳本名稱 | 功能說明 | 使用場景 |
|----------|----------|----------|
| `reset_database.sh` | 重置資料庫 | 開發環境重置 |
| `create_notification_center_db.sh` | 建立資料庫 | 初始化環境 |

### 使用方式
```bash
# 重置資料庫
./scripts/database/reset_database.sh

# 建立資料庫
./scripts/database/create_notification_center_db.sh
```

---

## 🔧 開發相關腳本 (development/)

### 腳本列表
| 腳本名稱 | 功能說明 | 使用場景 |
|----------|----------|----------|
| `lint.sh` | 代碼檢查 | 代碼品質檢查 |
| `build.sh` | 建置腳本 | 編譯和建置 |
| `gen.sh` | 代碼生成 | 自動生成代碼 |

### 使用方式
```bash
# 代碼檢查
./scripts/development/lint.sh

# 建置專案
./scripts/development/build.sh

# 生成代碼
./scripts/development/gen.sh
```

---

## 🧪 測試相關腳本 (testing/)

### 腳本列表
| 腳本名稱 | 功能說明 | 使用場景 |
|----------|----------|----------|
| `e2e_test.sh` | E2E 測試 | 端到端測試 |
| `load_test.js` | 負載測試 | 效能測試 |
| `setup_test_env.sh` | 測試環境準備 | 環境設定 |
| `run_e2e_tests.sh` | 整合測試 | 完整測試流程 |
| `generate_test_report.sh` | 測試報告生成 | 報告生成 |
| `test_api.sh` | API 測試 | API 功能測試 |
| `test_queue.sh` | 佇列測試 | 佇列功能測試 |
| `comprehensive_api_test.sh` | 綜合 API 測試 | 完整 API 測試 |
| `comprehensive_api_test_no_thirdparty.sh` | 無第三方 API 測試 | 隔離測試 |
| `api_test_automation.sh` | API 測試自動化 | 自動化測試 |
| `golden.sh` | 黃金測試 | 基準測試 |

### 使用方式
```bash
# E2E 測試
./scripts/testing/e2e_test.sh

# 負載測試
./scripts/testing/load_test.js

# 測試環境準備
./scripts/testing/setup_test_env.sh

# 整合測試
./scripts/testing/run_e2e_tests.sh

# 生成測試報告
./scripts/testing/generate_test_report.sh
```

---

## 🚀 部署相關腳本 (deployment/)

### 腳本列表
| 腳本名稱 | 功能說明 | 使用場景 |
|----------|----------|----------|
| `start_server.sh` | 啟動服務 | 服務啟動 |
| `deploy.sh` | 部署腳本 | 應用部署 |
| `openapi.sh` | OpenAPI 生成 | API 文檔生成 |
| `docker.sh` | Docker 管理 | 容器管理 |

### 使用方式
```bash
# 啟動服務
./scripts/deployment/start_server.sh

# 部署應用
./scripts/deployment/deploy.sh

# 生成 OpenAPI 文檔
./scripts/deployment/openapi.sh

# Docker 管理
./scripts/deployment/docker.sh
```

---

## 🎯 腳本使用指南

### 1. 開發環境設定
```bash
# 1. 建立資料庫
./scripts/database/create_notification_center_db.sh

# 2. 代碼檢查
./scripts/development/lint.sh

# 3. 建置專案
./scripts/development/build.sh
```

### 2. 測試執行
```bash
# 1. 準備測試環境
./scripts/testing/setup_test_env.sh

# 2. 執行 E2E 測試
./scripts/testing/e2e_test.sh

# 3. 生成測試報告
./scripts/testing/generate_test_report.sh
```

### 3. 部署流程
```bash
# 1. 生成 OpenAPI 文檔
./scripts/deployment/openapi.sh

# 2. 部署應用
./scripts/deployment/deploy.sh

# 3. 啟動服務
./scripts/deployment/start_server.sh
```

---

## 🚀 開發部署腳本 (新增)

### 腳本列表
| 腳本名稱 | 功能說明 | 使用場景 |
|----------|----------|----------|
| `dev-deploy.sh` | 完整開發部署 | 代碼更新後完整重新部署 |
| `quick-redeploy.sh` | 快速重構 | 代碼更新後快速重新部署 |

### 使用方式
```bash
# 完整開發部署 (推薦用於重大更新)
./scripts/dev-deploy.sh

# 快速重構 (推薦用於小修改)
./scripts/quick-redeploy.sh

# 或使用 Makefile
make dev-deploy
```

### 腳本功能對比
| 功能 | dev-deploy.sh | quick-redeploy.sh |
|------|---------------|-------------------|
| **停止服務** | ✅ 完整停止 | ✅ 快速停止 |
| **清理映像** | ✅ 清理舊映像 | ❌ 跳過清理 |
| **重新構建** | ✅ 完整構建 | ✅ 快速構建 |
| **資料庫遷移** | ✅ 執行遷移 | ❌ 跳過遷移 |
| **功能測試** | ✅ 完整測試 | ✅ 基本測試 |
| **適用場景** | 重大更新 | 小修改 |

---

## 🔧 腳本維護

### 腳本標準
1. **可執行權限**：所有腳本都應有執行權限
2. **錯誤處理**：包含適當的錯誤處理
3. **日誌記錄**：提供清晰的日誌輸出
4. **參數驗證**：驗證輸入參數

### 腳本更新
1. **版本控制**：使用 Git 追蹤腳本變更
2. **文檔更新**：更新相關文檔
3. **測試驗證**：確保腳本正常執行
4. **向後相容**：保持向後相容性

---

## 📚 相關文檔

- [開發指南](../docs/07_DEVELOPMENT/DevelopmentGuide.md)
- [測試指南](../docs/04_TEST/TestPlan.md)
- [部署指南](../docs/05_DEPLOYMENT/DeploymentGuide.md)
- [用戶手冊](../docs/06_USER_GUIDE/UserManual.md)

---

**建立時間**: 2025-10-09  
**維護者**: AI Assistant  
**狀態**: ✅ 完成