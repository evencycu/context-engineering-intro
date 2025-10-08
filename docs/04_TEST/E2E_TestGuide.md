# 🧪 E2E (End-to-End) 測試指南

## 1. 文件資訊
- **版本**：v1.0
- **作者**：QA Engineer
- **最後更新**：2025-10-08

---

## 2. E2E 測試概述

E2E 測試是驗證整個 Teams Notification API 系統從用戶請求到 Teams 通知送達的完整流程，確保所有組件協同工作正常。

### 2.1 測試範圍
- **完整業務流程**：從 API 請求到 Teams 通知送達
- **系統整合**：API Gateway → Queue → Actor → Teams API
- **資料一致性**：資料庫狀態與通知狀態同步
- **錯誤處理**：失敗重試、熔斷器機制
- **效能驗證**：高負載下的系統表現

### 2.2 測試環境
- **開發環境**：本地 Docker 容器
- **測試環境**：獨立的測試資料庫和 Redis
- **模擬環境**：使用 Mock Teams API 進行隔離測試

---

## 3. E2E 測試場景

### 3.1 基本通知流程測試

#### 場景 1：成功發送通知
```mermaid
sequenceDiagram
    participant Client as 外部系統
    participant API as API Gateway
    participant Queue as Redis Queue
    participant Actor as Actor Pool
    participant Teams as Teams API
    participant DB as PostgreSQL

    Client->>API: POST /api/v1/external/notify
    API->>DB: 儲存通知記錄
    API->>Queue: 加入佇列
    API->>Client: 回傳 202 + notification_id
    
    Queue->>Actor: 取出任務
    Actor->>Teams: 發送通知
    Teams->>Actor: 回傳成功
    Actor->>DB: 更新狀態為 sent
    Actor->>Queue: 移除任務
```

#### 場景 2：失敗重試流程
```mermaid
sequenceDiagram
    participant Client as 外部系統
    participant API as API Gateway
    participant Queue as Redis Queue
    participant Actor as Actor Pool
    participant Teams as Teams API
    participant DB as PostgreSQL

    Client->>API: POST /api/v1/external/notify
    API->>DB: 儲存通知記錄
    API->>Queue: 加入佇列
    API->>Client: 回傳 202 + notification_id
    
    Queue->>Actor: 取出任務
    Actor->>Teams: 發送通知
    Teams->>Actor: 回傳 429 Rate Limit
    Actor->>DB: 更新重試次數
    Actor->>Queue: 重新加入佇列（延遲）
    
    Note over Queue: 等待重試時間
    Queue->>Actor: 再次取出任務
    Actor->>Teams: 重新發送通知
    Teams->>Actor: 回傳成功
    Actor->>DB: 更新狀態為 sent
```

### 3.2 高負載測試場景

#### 場景 3：併發通知測試
- **目標**：驗證系統在高併發下的穩定性
- **參數**：1000 個併發請求
- **驗證點**：
  - 所有請求都能正確處理
  - 佇列不會積壓
  - 資料庫連線正常
  - 記憶體使用穩定

#### 場景 4：熔斷器測試
- **目標**：驗證熔斷器機制在 Teams API 故障時的保護作用
- **步驟**：
  1. 模擬 Teams API 返回 500 錯誤
  2. 觸發熔斷器開啟
  3. 驗證後續請求被拒絕
  4. 恢復 Teams API 正常
  5. 驗證熔斷器自動關閉

---

## 4. E2E 測試工具與框架

### 4.1 測試工具選擇

| 工具 | 用途 | 優點 |
|------|------|------|
| **k6** | 負載測試 | 高效能、易於編寫、支援 JavaScript |
| **Postman/Newman** | API 測試 | 圖形化介面、支援集合執行 |
| **curl + bash** | 簡單測試 | 輕量級、易於整合 |
| **Go test** | 單元整合 | 與專案語言一致、CI/CD 友好 |

### 4.2 測試環境準備

#### 4.2.1 Docker 環境
```bash
# 啟動完整測試環境
docker-compose -f docker-compose.test.yml up -d

# 等待服務就緒
./scripts/testing/wait_for_services.sh

# 執行資料庫遷移
make db-migrate
```

#### 4.2.2 測試資料準備
```bash
# 載入測試資料
./scripts/testing/load_test_data.sh

# 驗證測試環境
./scripts/testing/verify_environment.sh
```

---

## 5. E2E 測試執行

### 5.1 基本功能測試

#### 5.1.1 健康檢查測試
```bash
#!/bin/bash
# 測試服務健康狀態
curl -f http://localhost:8080/health || exit 1
echo "✅ 服務健康檢查通過"
```

#### 5.1.2 通知發送測試
```bash
#!/bin/bash
# 測試基本通知發送
NOTIFY_KEY="test-project-key"
MESSAGE="E2E 測試通知 - $(date)"

RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d "{
    \"notify_key\": \"$NOTIFY_KEY\",
    \"message\": \"$MESSAGE\",
    \"priority\": \"normal\"
  }")

echo "$RESPONSE" | jq '.'
```

### 5.2 負載測試

#### 5.2.1 k6 負載測試腳本
```javascript
// load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 100 },  // 漸增到 100 用戶
    { duration: '1m', target: 100 },   // 維持 100 用戶
    { duration: '30s', target: 0 },    // 漸減到 0 用戶
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'], // 95% 請求 < 300ms
    http_req_failed: ['rate<0.1'],    // 錯誤率 < 10%
  },
};

export default function() {
  const payload = JSON.stringify({
    notify_key: 'test-project-key',
    message: `Load test message ${__VU}_${__ITER}`,
    priority: 'normal'
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const response = http.post('http://localhost:8080/api/v1/external/notify', payload, params);
  
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 300ms': (r) => r.timings.duration < 300,
    'has notification_id': (r) => JSON.parse(r.body).data.notification_id !== null,
  });

  sleep(1);
}
```

#### 5.2.2 執行負載測試
```bash
# 執行 k6 負載測試
k6 run scripts/testing/load_test.js

# 生成詳細報告
k6 run --out json=results.json scripts/testing/load_test.js
```

### 5.3 錯誤處理測試

#### 5.3.1 無效請求測試
```bash
#!/bin/bash
# 測試無效的 notify_key
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "invalid-key",
    "message": "Test message"
  }')

# 驗證錯誤回應
echo "$RESPONSE" | jq '.success' | grep -q "false" || exit 1
echo "✅ 無效請求處理正確"
```

#### 5.3.2 熔斷器測試
```bash
#!/bin/bash
# 模擬 Teams API 故障
# 1. 修改配置使 Teams API 返回錯誤
# 2. 發送多個請求觸發熔斷器
# 3. 驗證熔斷器狀態
# 4. 恢復正常配置
# 5. 驗證熔斷器關閉
```

---

## 6. 測試結果分析

### 6.1 效能指標

| 指標 | 目標值 | 監控方法 |
|------|--------|----------|
| **回應時間** | P95 < 300ms | k6 統計 |
| **吞吐量** | > 1000 req/s | 負載測試 |
| **錯誤率** | < 1% | 監控日誌 |
| **佇列深度** | < 100 | Redis 監控 |
| **資料庫連線** | < 80% | PostgreSQL 監控 |

### 6.2 測試報告

#### 6.2.1 自動化報告生成
```bash
#!/bin/bash
# 生成測試報告
./scripts/testing/generate_report.sh

# 報告包含：
# - 測試執行摘要
# - 效能指標分析
# - 錯誤統計
# - 建議改進項目
```

#### 6.2.2 測試結果範例
```json
{
  "test_summary": {
    "total_tests": 150,
    "passed": 147,
    "failed": 3,
    "success_rate": "98%"
  },
  "performance": {
    "avg_response_time": "120ms",
    "p95_response_time": "280ms",
    "throughput": "1200 req/s"
  },
  "errors": [
    {
      "type": "timeout",
      "count": 2,
      "description": "Redis connection timeout"
    },
    {
      "type": "validation",
      "count": 1,
      "description": "Invalid notify_key format"
    }
  ]
}
```

---

## 7. 持續整合

### 7.1 CI/CD 整合

#### 7.1.1 GitHub Actions 配置
```yaml
# .github/workflows/e2e-test.yml
name: E2E Tests
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  e2e-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.24
      
      - name: Start test environment
        run: |
          docker-compose -f docker-compose.test.yml up -d
          ./scripts/testing/wait_for_services.sh
      
      - name: Run E2E tests
        run: |
          ./scripts/testing/run_e2e_tests.sh
      
      - name: Generate report
        run: |
          ./scripts/testing/generate_report.sh
      
      - name: Upload test results
        uses: actions/upload-artifact@v3
        with:
          name: e2e-test-results
          path: test_results/
```

### 7.2 測試環境管理

#### 7.2.1 環境隔離
```bash
# 使用獨立的測試環境
export TEST_ENV=true
export DATABASE_URL="postgresql://test:test@localhost:5433/testdb"
export REDIS_URL="redis://localhost:6380"
```

#### 7.2.2 資料清理
```bash
# 測試前清理
./scripts/testing/cleanup_test_data.sh

# 測試後清理
./scripts/testing/cleanup_after_test.sh
```

---

## 8. 最佳實踐

### 8.1 測試設計原則
- **獨立性**：每個測試案例獨立執行
- **可重複性**：測試結果一致
- **隔離性**：不影響其他測試
- **可維護性**：易於更新和擴展

### 8.2 測試資料管理
- 使用固定的測試資料
- 測試前後清理資料
- 避免硬編碼的測試資料
- 使用環境變數配置

### 8.3 監控與告警
- 設定測試失敗告警
- 監控測試執行時間
- 追蹤效能回歸
- 自動化測試報告

---

## 9. 故障排除

### 9.1 常見問題

#### 9.1.1 測試環境問題
```bash
# 檢查服務狀態
docker ps | grep teamsnotify

# 檢查日誌
docker logs teamsnotify-postgres
docker logs teamsnotify-redis

# 重啟服務
docker-compose restart
```

#### 9.1.2 測試失敗分析
```bash
# 查看詳細錯誤
./scripts/testing/debug_failed_tests.sh

# 檢查系統資源
./scripts/testing/check_system_resources.sh

# 分析效能瓶頸
./scripts/testing/analyze_performance.sh
```

### 9.2 除錯工具

#### 9.2.1 日誌分析
```bash
# 收集測試日誌
./scripts/testing/collect_logs.sh

# 分析錯誤模式
./scripts/testing/analyze_error_patterns.sh
```

#### 9.2.2 效能分析
```bash
# 生成效能報告
./scripts/testing/generate_performance_report.sh

# 識別瓶頸
./scripts/testing/identify_bottlenecks.sh
```

---

## 10. 參考資源

### 10.1 測試工具文檔
- [k6 官方文檔](https://k6.io/docs/)
- [Postman 測試指南](https://learning.postman.com/docs/writing-scripts/test-scripts/)
- [Go 測試文檔](https://golang.org/pkg/testing/)

### 10.2 專案相關文檔
- [測試計劃](./TestPlan.md)
- [測試案例](./TestCases.md)
- [開發指南](../07_DEVELOPMENT/DevelopmentGuide.md)
- [部署指南](../05_DEPLOYMENT/DeploymentGuide.md)

---

**建立日期**: 2025-10-08  
**版本**: v1.0  
**狀態**: ✅ 開發完成，測試通過
