# 🧪 Test Plan

## 1. 文件資訊
- **版本**：v1.0
- **作者**：QA Engineer  
- **最後更新**：2025-10-08  

---

## 2. 測試目標
確保 Teams Notification API 系統符合需求與品質標準，涵蓋功能、效能、安全與相容性。

---

## 3. 測試範圍

| 類別 | 範例項目 | 說明 |
|------|-----------|------|
| **功能測試** | API 端點、通知發送、Bot 管理 | 驗證功能正確性 |
| **整合測試** | Redis / PostgreSQL / Teams API | 系統整合與外部依賴 |
| **效能測試** | 負載、延遲、佇列深度 | 檢測高壓環境表現 |
| **安全測試** | Token 驗證、權限管控、速率限制 | 確保安全與資料保護 |
| **回歸測試** | 關鍵路徑、核心功能 | 新版不影響舊功能 |

---

## 4. 測試策略

### 4.1 測試層級
- **單元測試（Unit Test）**：80% 覆蓋率目標
- **整合測試（Integration Test）**：以 staging DB 與 mock API 執行
- **端到端測試（E2E Test）**：完整業務流程測試
- **自動化測試（CI/CD）**：GitHub Actions 於每次 PR 觸發

### 4.2 測試環境
- **開發環境**：本地 Docker 容器
- **測試環境**：獨立的測試資料庫和 Redis
- **預生產環境**：模擬生產環境配置

---

## 5. 測試工具

| 類型 | 工具 | 用途 |
|------|------|------|
| **單元測試** | Go test | 函數和方法測試 |
| **API 測試** | curl, Postman, Newman | API 端點測試 |
| **效能測試** | k6, JMeter | 負載和壓力測試 |
| **安全測試** | OWASP ZAP, Trivy | 安全漏洞掃描 |
| **CI/CD** | GitHub Actions | 自動化測試執行 |

---

## 6. 測試環境設定

### 6.1 必要條件
- Go 1.21+
- Docker 和 Docker Compose
- PostgreSQL 14+
- Redis 6+
- Teams Bot 憑證

### 6.2 環境變數
```bash
export TEAMS_BOT_APP_ID=844146d7-4ac9-4e4d-a463-d6e027714e81
export TEAMS_TENANT_ID=051cece0-e4dc-4aed-b471-bf29824e1ee6
export TEAMS_BOT_APP_PASSWORD='HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy'
export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
```

---

## 7. 測試類型詳述

### 7.1 API 測試
- **健康檢查**：服務狀態驗證
- **認證測試**：JWT 和 API Key 驗證
- **CRUD 操作**：所有資源的增刪改查
- **業務邏輯**：通知發送和處理流程

### 7.2 效能測試
- **負載測試**：並發請求處理
- **壓力測試**：系統極限測試
- **記憶體測試**：記憶體洩漏檢測

### 7.3 安全測試
- **認證繞過**：未授權存取測試
- **注入攻擊**：SQL 注入防護
- **速率限制**：防止濫用測試

---

## 8. 驗收標準

| 指標 | 目標值 | 說明 |
|------|--------|------|
| **功能測試通過率** | ≥ 95% | 所有測試案例通過率 |
| **效能測試** | p95 延遲 < 300ms | API 響應時間 |
| **錯誤率** | < 0.5% | 系統錯誤率 |
| **關鍵測試案例** | 100% 通過 | 核心功能測試 |

---

## 9. 測試執行流程

### 9.1 快速測試
```bash
# 啟動服務
bash start_server.sh

# 執行基本測試
bash scripts/test_api.sh

# 執行完整測試
bash scripts/comprehensive_api_test.sh
```

### 9.2 詳細測試
```bash
# 健康檢查
curl http://localhost:8080/health

# API 測試
curl -X POST http://localhost:8080/api/v1/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38",
    "message": "測試訊息",
    "messageType": "text",
    "priority": "normal",
    "targets": ["all"]
  }'
```

---

## 10. 測試資料管理

### 10.1 測試資料庫
- 使用 `scripts/seed_minimal.sql` 建立測試資料
- 包含公司、用戶、專案、目的地等基本資料
- 使用真實的 Teams Bot 憑證

### 10.2 測試 Bot 安裝
- 使用 `scripts/add_bot_installations.sql` 建立 Bot 安裝記錄
- 包含個人、群組、頻道等不同類型的對話
- 所有安裝狀態設為 `active`

---

## 11. 測試結果管理

### 11.1 成功指標
- ✅ 所有 API 端點正常回應
- ✅ 認證機制正常運作
- ✅ 通知可以成功發送到 Teams
- ✅ 錯誤處理機制正常
- ✅ 效能指標符合預期

### 11.2 失敗處理
- 檢查服務日誌: `tail -f server.log`
- 檢查資料庫連線: `docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center`
- 檢查 Redis 連線: `docker exec teamsnotify-redis redis-cli ping`

---

## 12. 持續整合

### 12.1 自動化測試
- 每次程式碼提交觸發測試
- 測試失敗阻止部署
- 測試結果報告和通知

### 12.2 測試覆蓋率
- 程式碼覆蓋率監控
- 未覆蓋程式碼標記
- 覆蓋率趨勢追蹤

---

## 13. 故障排除

### 13.1 常見問題
1. **服務無法啟動**: 檢查端口是否被佔用
2. **資料庫連線失敗**: 檢查 Docker 容器狀態
3. **Teams 認證失敗**: 檢查 Bot 憑證設定
4. **通知發送失敗**: 檢查 Bot 安裝狀態

### 13.2 除錯工具
- 服務日誌: `/tmp/teamsnotify_server.log`
- 資料庫查詢: 使用 `psql` 直接查詢
- Redis 監控: 使用 `redis-cli` 檢查狀態

---

## 14. 最佳實踐

### 14.1 測試設計
- 測試用例要獨立且可重複
- 使用真實的測試資料
- 包含正面和負面測試案例

### 14.2 測試維護
- 定期更新測試資料
- 保持測試腳本與 API 同步
- 記錄測試結果和問題

### 14.3 效能優化
- 並行執行測試用例
- 使用測試資料庫
- 避免外部依賴

---

**版本**: v1.0  
**最後更新**: 2025-10-08  
**作者**: TeamsNotify QA Team
