# Scripts 目錄重構計畫

## 🎯 目標
重新組織腳本目錄，使其更清晰、更易維護，並符合最佳實踐。

## 📊 目前問題分析

### 現有分類問題：
1. **`root/` 目錄命名不清晰** - 不清楚用途
2. **`local-test/` 與 `testing/` 重疊** - 功能相似但分離
3. **`deployment/` 與 `root/` 功能重疊** - 都有部署相關腳本
4. **分類邏輯不一致** - 有些按功能分類，有些按環境分類

## 🏗️ 新分類結構

### 按功能分類 (Function-based)
```
scripts/
├── build/                     # 建置相關
│   ├── build.sh
│   ├── gen.sh
│   └── lint.sh
├── deploy/                    # 部署相關
│   ├── dev-deploy.sh
│   ├── smart-deploy.sh
│   ├── quick-redeploy.sh
│   ├── deploy.sh
│   └── docker.sh
├── test/                      # 測試相關
│   ├── e2e/
│   │   ├── e2e_test.sh
│   │   ├── run_e2e_tests.sh
│   │   └── setup_test_env.sh
│   ├── api/
│   │   ├── api_test_automation.sh
│   │   ├── test_api.sh
│   │   └── comprehensive_api_test.sh
│   ├── performance/
│   │   ├── load_test.js
│   │   ├── load_test_enhanced.js
│   │   └── performance_benchmark.sh
│   ├── security/
│   │   └── security_test.sh
│   └── integration/
│       ├── test-docker.sh
│       └── test_queue.sh
├── database/                  # 資料庫相關
│   ├── schema.sql
│   ├── init.sql
│   ├── create_notification_center_db.sh
│   ├── reset_database.sh
│   └── migrations/
├── docker/                    # Docker 配置
│   ├── docker-compose.yml
│   ├── docker-compose-persistent.yml
│   └── docker-compose-local.yml
└── utils/                     # 工具腳本
    ├── cleanup.sh
    ├── setup-local-env.sh
    └── generate_test_report.sh
```

### 按環境分類 (Environment-based)
```
scripts/
├── local/                     # 本地開發
│   ├── build/
│   ├── deploy/
│   ├── test/
│   └── docker/
├── dev/                       # 開發環境
│   ├── deploy/
│   └── test/
├── staging/                   # 測試環境
│   ├── deploy/
│   └── test/
└── prod/                      # 生產環境
    ├── deploy/
    └── test/
```

## 🎯 推薦方案：混合分類

### 主要按功能分類，次要按環境分類
```
scripts/
├── build/                     # 建置工具
│   ├── build.sh
│   ├── gen.sh
│   └── lint.sh
├── deploy/                    # 部署工具
│   ├── docker/                # Docker 部署
│   │   ├── deploy.sh         # 統一部署腳本 (整合所有部署模式)
│   │   └── docker.sh         # Docker 容器管理工具
│   └── k8s/                   # Kubernetes 部署
├── test/                      # 測試工具
│   ├── e2e/                   # 端到端測試
│   │   ├── e2e_test.sh
│   │   ├── run_e2e_tests.sh
│   │   └── setup_test_env.sh
│   ├── api/                   # API 測試
│   │   ├── api_test_automation.sh
│   │   ├── test_api.sh
│   │   └── comprehensive_api_test.sh
│   ├── performance/           # 效能測試
│   │   ├── load_test.js
│   │   ├── load_test_enhanced.js
│   │   └── performance_benchmark.sh
│   ├── security/              # 安全測試
│   │   └── security_test.sh
│   └── integration/           # 整合測試
│       ├── test-docker.sh
│       └── test_queue.sh
├── database/                  # 資料庫工具
│   ├── schema.sql
│   ├── init.sql
│   ├── create_notification_center_db.sh
│   ├── reset_database.sh
│   └── migrations/
├── docker/                    # Docker 配置
│   ├── docker-compose.yml
│   ├── docker-compose-persistent.yml
│   └── docker-compose-local.yml
└── utils/                     # 工具腳本
    ├── cleanup.sh
    ├── setup-local-env.sh
    └── generate_test_report.sh
```

## 🔄 遷移步驟

### 第一階段：創建新目錄結構
1. 創建新的目錄結構
2. 移動檔案到對應目錄
3. 更新所有路徑引用

### 第二階段：清理和優化
1. 移除空的目錄
2. 更新 README 檔案
3. 更新 Makefile 中的路徑

### 第三階段：測試和驗證
1. 測試所有腳本是否正常工作
2. 更新文檔
3. 驗證向後兼容性

## 📋 遷移對照表

| 現有路徑 | 新路徑 | 說明 |
|---------|--------|------|
| `scripts/root/dev-deploy.sh` | `scripts/deploy/local/dev-deploy.sh` | 本地部署腳本 |
| `scripts/root/smart-deploy.sh` | `scripts/deploy/local/smart-deploy.sh` | 智能部署腳本 |
| `scripts/root/quick-redeploy.sh` | `scripts/deploy/local/quick-redeploy.sh` | 快速重構腳本 |
| `scripts/deployment/deploy.sh` | `scripts/deploy/docker/deploy.sh` | Docker 部署腳本 |
| `scripts/deployment/docker.sh` | `scripts/deploy/docker/docker.sh` | Docker 工具腳本 |
| `scripts/local-test/test-docker.sh` | `scripts/test/integration/test-docker.sh` | 整合測試腳本 |
| `scripts/testing/e2e_test.sh` | `scripts/test/e2e/e2e_test.sh` | E2E 測試腳本 |
| `scripts/testing/api_test_automation.sh` | `scripts/test/api/api_test_automation.sh` | API 測試腳本 |
| `scripts/development/build.sh` | `scripts/build/build.sh` | 建置腳本 |

## 🎯 優點

1. **清晰的功能分類** - 每個目錄都有明確的用途
2. **易於維護** - 相關腳本集中在一起
3. **可擴展性** - 容易添加新的腳本
4. **一致性** - 統一的命名和組織方式
5. **向後兼容** - 保持現有功能不變
