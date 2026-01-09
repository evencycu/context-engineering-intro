# 🧪 Test Plan

## 1. 文件資訊
- **版本**：v1.0
- **作者**：QA Engineer  
- **最後更新**：2025-10-08  

---

## 2. 測試目標
確保系統符合需求與品質標準，涵蓋功能、效能、安全與相容性。

---

## 3. 測試範圍
| 類別 | 範例項目 | 說明 |
|------|-----------|------|
| 功能測試 | API、UI、Bot回覆 | 驗證功能正確性 |
| 整合測試 | Redis / DB / Graph API | 系統整合與外部依賴 |
| 效能測試 | 負載、延遲、佇列深度 | 檢測高壓環境表現 |
| 安全測試 | Token驗證、權限管控 | 確保安全與資料保護 |
| 回歸測試 | 關鍵路徑 | 新版不影響舊功能 |

---

## 4. 測試策略
- 單元測試（Unit Test）：80% 覆蓋率目標。
- 整合測試（Integration Test）：以 staging DB 與 mock API 執行。
- 自動化測試（CI/CD）：GitHub Actions 於每次 PR 觸發。
- 壓力測試（Load Test）：使用 k6 / Locust 模擬 1000 併發。

---

## 5. 測試工具
| 類型 | 工具 |
|------|------|
| 單元測試 | Go test / Pytest |
| API 測試 | Postman / Newman |
| 效能測試 | k6 / JMeter |
| 安全測試 | OWASP ZAP / Trivy |
| CI/CD | GitHub Actions / Azure Pipeline |

---

## 6. 驗收標準
- 功能測試通過率 ≥ 95%
- 效能測試：p95 延遲 < 300ms
- 錯誤率 < 0.5%
- 關鍵測試案例全通過
