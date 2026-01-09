# 測試概覽

Teams Notification API 提供完整的測試策略和工具，確保系統的穩定性和可靠性。

## 測試策略

### 1. 單元測試
- **範圍**: 個別函數和方法的測試
- **工具**: Go 內建 testing 包
- **覆蓋率**: 目標 80% 以上

### 2. 整合測試
- **範圍**: API 端點和服務整合
- **工具**: curl, 測試腳本
- **環境**: 本地開發環境

### 3. 端到端測試
- **範圍**: 完整業務流程測試
- **工具**: 自動化測試腳本
- **環境**: 測試環境

## 測試環境設定

### 必要條件
- Go 1.21+
- Docker 和 Docker Compose
- PostgreSQL 14+
- Redis 6+
- Teams Bot 憑證

### 環境變數
```bash
export TEAMS_BOT_APP_ID=844146d7-4ac9-4e4d-a463-d6e027714e81
export TEAMS_TENANT_ID=051cece0-e4dc-4aed-b471-bf29824e1ee6
export TEAMS_BOT_APP_PASSWORD='HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy'
export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
```

## 測試類型

### API 測試
- **健康檢查**: 服務狀態驗證
- **認證測試**: JWT 和 API Key 驗證
- **CRUD 操作**: 所有資源的增刪改查
- **業務邏輯**: 通知發送和處理流程

### 效能測試
- **負載測試**: 並發請求處理
- **壓力測試**: 系統極限測試
- **記憶體測試**: 記憶體洩漏檢測

### 安全測試
- **認證繞過**: 未授權存取測試
- **注入攻擊**: SQL 注入防護
- **速率限制**: 防止濫用測試

## 測試工具

### 1. 自動化測試腳本
- `scripts/comprehensive_api_test.sh` - 完整 API 測試
- `scripts/test_api.sh` - 基本 API 測試
- `scripts/test_queue.sh` - 佇列功能測試

### 2. 手動測試工具
- **curl**: HTTP 請求測試
- **Postman**: API 集合測試
- **Swagger UI**: 互動式 API 測試

### 3. 監控工具
- **日誌監控**: 即時日誌查看
- **健康檢查**: 服務狀態監控
- **效能指標**: 回應時間和吞吐量

## 測試資料

### 測試資料庫
- 使用 `scripts/seed_minimal.sql` 建立測試資料
- 包含公司、用戶、專案、目的地等基本資料
- 使用真實的 Teams Bot 憑證

### 測試 Bot 安裝
- 使用 `scripts/add_bot_installations.sql` 建立 Bot 安裝記錄
- 包含個人、群組、頻道等不同類型的對話
- 所有安裝狀態設為 `active`

## 測試執行

### 快速測試
```bash
# 啟動服務
bash start_server.sh

# 執行基本測試
bash scripts/test_api.sh

# 執行完整測試
bash scripts/comprehensive_api_test.sh
```

### 詳細測試
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

## 測試結果

### 成功指標
- ✅ 所有 API 端點正常回應
- ✅ 認證機制正常運作
- ✅ 通知可以成功發送到 Teams
- ✅ 錯誤處理機制正常
- ✅ 效能指標符合預期

### 失敗處理
- 檢查服務日誌: `tail -f server.log`
- 檢查資料庫連線: `docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center`
- 檢查 Redis 連線: `docker exec teamsnotify-redis redis-cli ping`

## 持續整合

### 自動化測試
- 每次程式碼提交觸發測試
- 測試失敗阻止部署
- 測試結果報告和通知

### 測試覆蓋率
- 程式碼覆蓋率監控
- 未覆蓋程式碼標記
- 覆蓋率趨勢追蹤

## 故障排除

### 常見問題
1. **服務無法啟動**: 檢查端口是否被佔用
2. **資料庫連線失敗**: 檢查 Docker 容器狀態
3. **Teams 認證失敗**: 檢查 Bot 憑證設定
4. **通知發送失敗**: 檢查 Bot 安裝狀態

### 除錯工具
- 服務日誌: `/tmp/teamsnotify_server.log`
- 資料庫查詢: 使用 `psql` 直接查詢
- Redis 監控: 使用 `redis-cli` 檢查狀態

## 最佳實踐

### 測試設計
- 測試用例要獨立且可重複
- 使用真實的測試資料
- 包含正面和負面測試案例

### 測試維護
- 定期更新測試資料
- 保持測試腳本與 API 同步
- 記錄測試結果和問題

### 效能優化
- 並行執行測試用例
- 使用測試資料庫
- 避免外部依賴
