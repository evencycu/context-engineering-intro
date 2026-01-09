# 🧪 E2E 測試指南

## 概述

本目錄包含 Teams Notification API 的完整端到端 (E2E) 測試方案，包括基本功能測試、負載測試、API 測試和測試報告生成。

## 檔案結構

```
scripts/testing/
├── README.md                    # 本檔案
├── e2e_test.sh                  # E2E 測試主腳本
├── load_test.js                 # k6 負載測試腳本
├── setup_test_env.sh            # 測試環境準備腳本
├── generate_test_report.sh      # 測試報告生成腳本
├── run_e2e_tests.sh             # 整合測試執行腳本
└── test_api.sh                  # API 測試腳本
```

## 快速開始

### 1. 執行完整測試流程

```bash
# 執行完整的 E2E 測試流程
make test-full

# 或直接執行腳本
./scripts/testing/run_e2e_tests.sh run
```

### 2. 執行特定測試

```bash
# 僅執行 E2E 測試
make test-e2e

# 僅執行負載測試
make test-load

# 僅執行 API 測試
make test-api

# 生成測試報告
make test-report
```

## 測試類型

### 1. E2E 測試 (`e2e_test.sh`)

**功能**：端到端功能測試
**範圍**：
- 基本通知發送流程
- 錯誤處理機制
- 佇列系統運作
- 資料一致性驗證

**執行方式**：
```bash
./scripts/testing/e2e_test.sh
```

### 2. 負載測試 (`load_test.js`)

**功能**：效能和壓力測試
**工具**：k6
**範圍**：
- 併發請求處理
- 回應時間測試
- 系統穩定性驗證

**執行方式**：
```bash
# 需要先安裝 k6
# macOS: brew install k6
# Ubuntu: sudo apt-get install k6

k6 run scripts/testing/load_test.js
```

### 3. API 測試 (`test_api.sh`)

**功能**：API 端點測試
**範圍**：
- 所有 API 端點
- 請求/回應驗證
- 錯誤處理測試

**執行方式**：
```bash
./scripts/testing/test_api.sh
```

## 測試環境準備

### 1. 自動環境準備

```bash
# 設定測試環境
./scripts/testing/setup_test_env.sh setup

# 驗證測試環境
./scripts/testing/setup_test_env.sh verify

# 清理測試環境
./scripts/testing/setup_test_env.sh cleanup
```

### 2. 手動環境準備

```bash
# 啟動服務
make docker-up

# 執行資料庫遷移
make db-migrate

# 啟動 API 服務
make run
```

## 測試結果

### 1. 結果目錄

```
test_results/
├── e2e/                    # E2E 測試結果
│   ├── e2e_test_*.log      # 測試日誌
│   ├── e2e_report_*.json   # 測試報告
│   └── test_data.env       # 測試資料
├── load/                   # 負載測試結果
│   └── load_test_*.json    # 負載測試結果
└── api/                   # API 測試結果
    └── api_test_*.log     # API 測試日誌
```

### 2. 報告生成

```bash
# 生成測試報告
./scripts/testing/generate_test_report.sh generate

# 查看報告摘要
./scripts/testing/generate_test_report.sh summary
```

**報告類型**：
- HTML 報告：`test_reports/test_report_*.html`
- JSON 報告：`test_reports/test_report_*.json`
- 文字報告：`test_reports/test_report_*.txt`

## 配置選項

### 1. 環境變數

```bash
# API 服務配置
export BASE_URL="http://localhost:8080"
export API_BASE="http://localhost:8080/api/v1"

# 測試資料配置
export TEST_COMPANY_ID="your-company-id"
export TEST_USER_ID="your-user-id"
export TEST_PROJECT_ID="your-project-id"
export TEST_NOTIFY_KEY="your-notify-key"
```

### 2. 負載測試配置

修改 `load_test.js` 中的 `options` 配置：

```javascript
export let options = {
  stages: [
    { duration: '30s', target: 50 },   // 漸增到 50 用戶
    { duration: '1m', target: 50 },    // 維持 50 用戶
    { duration: '30s', target: 100 },  // 漸增到 100 用戶
    { duration: '1m', target: 100 },    // 維持 100 用戶
    { duration: '30s', target: 0 },     // 漸減到 0 用戶
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% 請求 < 500ms
    http_req_failed: ['rate<0.05'],     // 錯誤率 < 5%
  },
};
```

## 故障排除

### 1. 常見問題

**問題**：服務未啟動
```bash
# 檢查服務狀態
curl http://localhost:8080/health

# 啟動服務
make docker-up
make run
```

**問題**：資料庫連線失敗
```bash
# 檢查資料庫狀態
docker ps | grep postgres

# 重啟資料庫
docker-compose restart postgres
```

**問題**：Redis 連線失敗
```bash
# 檢查 Redis 狀態
docker ps | grep redis

# 重啟 Redis
docker-compose restart redis
```

### 2. 除錯模式

```bash
# 啟用詳細日誌
export LOG_LEVEL=debug

# 執行測試
./scripts/testing/e2e_test.sh
```

### 3. 清理環境

```bash
# 清理測試資料
./scripts/testing/setup_test_env.sh cleanup

# 清理 Docker 環境
docker-compose down -v
```

## 最佳實踐

### 1. 測試執行順序

1. **環境準備**：確保所有服務正常運行
2. **基本測試**：執行 E2E 測試驗證基本功能
3. **效能測試**：執行負載測試驗證效能
4. **報告生成**：生成測試報告分析結果

### 2. 測試資料管理

- 使用獨立的測試資料
- 測試前後清理資料
- 避免硬編碼的測試資料
- 使用環境變數配置

### 3. 監控和告警

- 設定測試失敗告警
- 監控測試執行時間
- 追蹤效能回歸
- 自動化測試報告

## 參考資源

### 1. 測試工具文檔
- [k6 官方文檔](https://k6.io/docs/)
- [Postman 測試指南](https://learning.postman.com/docs/writing-scripts/test-scripts/)
- [Go 測試文檔](https://golang.org/pkg/testing/)

### 2. 專案文檔
- [E2E 測試指南](../../docs/04_TEST/E2E_TestGuide.md)
- [測試計劃](../../docs/04_TEST/TestPlan.md)
- [測試案例](../../docs/04_TEST/TestCases.md)
- [開發指南](../../docs/07_DEVELOPMENT/DevelopmentGuide.md)

### 3. 外部資源
- [Go 官方文檔](https://golang.org/doc/)
- [Gin 框架文檔](https://gin-gonic.com/docs/)
- [PostgreSQL 文檔](https://www.postgresql.org/docs/)
- [Redis 文檔](https://redis.io/documentation)

---

**建立日期**: 2025-10-08  
**版本**: v1.0  
**狀態**: ✅ 開發完成，測試通過
