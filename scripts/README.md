# Scripts 目錄結構

本目錄包含所有 shell 腳本和 SQL 文件，按功能分類組織。

## 📁 目錄結構

### 🗄️ database/
資料庫相關腳本和 SQL 文件
- `add_bot_installations.sql` - 添加 Teams Bot 安裝記錄
- `clear_database.sql` - 清空資料庫
- `create_notification_center_db.sh` - 創建 notification_center 資料庫
- `reset_database.sh` - 重置資料庫
- `seed_minimal.sql` - 最小化測試資料
- `test_data_enhanced.sql` - 增強測試資料

### 🚀 deployment/
部署相關腳本
- `deploy.sh` - 主要部署腳本（部署 server 和 OpenAPI）
- `docker.sh` - Docker 相關操作
- `openapi.sh` - OpenAPI 服務管理
- `start_server.sh` - 啟動 server 腳本

### 🛠️ development/
開發相關腳本
- `build.sh` - 建置腳本
- `gen.sh` - 代碼生成腳本
- `lint.sh` - 代碼檢查腳本

### 🧪 testing/
測試相關腳本
- `api_test_automation.sh` - API 測試自動化
- `comprehensive_api_test.sh` - 綜合 API 測試
- `comprehensive_api_test_no_thirdparty.sh` - 綜合 API 測試（無第三方）
- `golden.sh` - Golden 測試腳本
- `test_api.sh` - API 測試腳本
- `test_queue.sh` - 佇列測試腳本

### 📋 migrations/
資料庫遷移腳本
- `001_enhance_bot_installations.sql` - 增強 Bot 安裝表
- `002_drop_user_foreign_keys.sql` - 移除用戶外鍵
- `003_make_sender_id_nullable.sql` - 使 sender_id 可為空
- `004_update_bot_installations_schema.sql` - 更新 Bot 安裝架構
- `005_installations_unique_on_conversation.sql` - 安裝記錄唯一性約束
- `006_unify_bots_to_teams_bots.sql` - 統一 Bots 到 Teams Bots
- `007_projects_rename_key_name_to_notify_key.sql` - 重命名專案鍵名
- `009_enhance_notification_destinations_for_async_actor.sql` - 增強通知目的地表

## 🚀 快速開始

### 部署服務
```bash
# 部署所有服務
./scripts/deployment/deploy.sh deploy

# 只部署 server
./scripts/deployment/deploy.sh server

# 只部署 OpenAPI
./scripts/deployment/deploy.sh openapi
```

### 資料庫操作
```bash
# 重置資料庫
./scripts/database/reset_database.sh

# 添加測試資料
./scripts/database/seed_minimal.sql
```

### 測試
```bash
# 運行 API 測試
./scripts/testing/test_api.sh

# 運行綜合測試
./scripts/testing/comprehensive_api_test.sh
```

### 開發
```bash
# 建置專案
./scripts/development/build.sh

# 代碼檢查
./scripts/development/lint.sh
```

## 📝 注意事項

1. 所有腳本都應該在專案根目錄執行
2. 確保 Docker 容器正在運行（postgres, redis）
3. 檢查環境變數設定
4. 查看日誌文件了解執行狀態

## 🔧 環境變數

主要環境變數：
- `DATABASE_URL` - 資料庫連接字串
- `TEAMS_BOT_APP_ID` - Teams Bot 應用程式 ID
- `TEAMS_TENANT_ID` - Teams 租戶 ID
- `TEAMS_BOT_APP_PASSWORD` - Teams Bot 應用程式密碼
- `REDIS_URL` - Redis 連接字串
