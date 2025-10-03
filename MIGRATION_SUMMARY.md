# 資料庫名稱遷移總結

## 📊 變更概述

將整個專案的資料庫名稱從 `codex_teams` 統一改為 `notification_center`。

**原因**: 統一命名規範，使資料庫名稱更符合專案實際用途（通知中心）。

---

## 🔄 變更的文件清單

### 核心程式碼（1 個文件）
- ✅ `cmd/server/main.go` - 預設 DATABASE_URL

### 腳本（2 個文件）
- ✅ `start_server.sh` - 啟動腳本
- ✅ `scripts/create_notification_center_db.sh` - 新增：建立資料庫腳本

### 文檔（5 個文件）
- ✅ `docs/DATABASE_CONFIG.md` - 資料庫配置文檔
- ✅ `docs/QUICK_START_QUEUE.md` - Queue 快速開始
- ✅ `docs/QUEUE_API_EXAMPLES.md` - Queue API 範例
- ✅ `CHANGELOG_QUEUE.md` - Queue 變更日誌
- ✅ `README.md` - 主文檔

**總計**: 8 個文件已更新，1 個新文件

---

## 📝 變更內容

### 資料庫連接字串

#### 修改前
```
postgresql://teamsnotify:teamsnotify123@localhost:5432/codex_teams?sslmode=disable
```

#### 修改後
```
postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable
```

### Docker 命令

#### 修改前
```bash
docker exec teamsnotify-postgres psql -U teamsnotify -d codex_teams
```

#### 修改後
```bash
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center
```

---

## 🚀 遷移步驟

### 選項 1: 自動遷移（推薦）

使用提供的腳本自動建立新資料庫並遷移數據：

```bash
./scripts/create_notification_center_db.sh
```

這個腳本會：
1. ✅ 檢查 PostgreSQL 容器狀態
2. ✅ 檢查是否存在 codex_teams 資料庫
3. ✅ 建立 notification_center 資料庫
4. ✅ 從 codex_teams 遷移所有數據（如果存在）
5. ✅ 驗證資料庫表結構

### 選項 2: 手動遷移

#### 步驟 1: 建立新資料庫

```bash
docker exec teamsnotify-postgres psql -U teamsnotify -c "CREATE DATABASE notification_center WITH OWNER teamsnotify ENCODING 'UTF8';"
```

#### 步驟 2: 遷移數據（如果有舊資料庫）

```bash
# 從 codex_teams 導出數據
docker exec teamsnotify-postgres pg_dump -U teamsnotify codex_teams > /tmp/backup.sql

# 導入到 notification_center
docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < /tmp/backup.sql
```

#### 步驟 3: 如果是全新安裝

```bash
# 執行 schema
docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < internal/database/schema.sql

# 執行所有 migrations
for migration in scripts/migrations/*.sql; do
    docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < "$migration"
done
```

---

## ✅ 驗證遷移

### 1. 檢查資料庫是否存在

```bash
docker exec teamsnotify-postgres psql -U teamsnotify -l | grep notification_center
```

預期輸出：
```
 notification_center | teamsnotify | UTF8
```

### 2. 檢查表數量

```bash
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';"
```

預期：應該有 17+ 個表

### 3. 檢查關鍵表

```bash
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "\dt"
```

應該包含：
- companies
- users  
- projects
- teams_bots
- bot_installations
- destinations
- notifications
- **failed_notifications** (Queue 系統)

### 4. 測試連接

```bash
export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
./server
```

查看日誌確認：
```
Successfully connected to database!
Queue manager started successfully
```

---

## 🧪 測試

### 1. 編譯

```bash
go build -o server cmd/server/main.go
```

### 2. 啟動服務

```bash
./start_server.sh
```

應該看到：
```
✅ PostgreSQL connection OK
✅ Queue table exists
✅ Server started successfully
   Health: healthy
   Queue: closed
```

### 3. 測試 Queue API

```bash
./scripts/test_queue.sh
```

### 4. 手動測試

```bash
# 測試健康檢查
curl http://localhost:8080/health

# 測試 Queue 狀態
curl http://localhost:8080/api/v1/queue/status | jq '.'

# 測試熔斷器指標
curl http://localhost:8080/api/v1/queue/circuit-breaker/metrics | jq '.'
```

---

## 📊 變更統計

| 項目 | 數量 |
|-----|------|
| 修改的文件 | 8 |
| 新增的文件 | 1 |
| 更新的連接字串 | 15+ |
| 更新的 Docker 命令 | 20+ |
| 編譯測試 | ✅ 通過 |

---

## 🔧 環境變數

### 更新你的環境配置

#### .env 文件

```bash
DATABASE_URL=postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable
```

#### Shell 環境

```bash
export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
```

#### Docker Compose（如果使用）

```yaml
environment:
  - DATABASE_URL=postgresql://teamsnotify:teamsnotify123@db:5432/notification_center?sslmode=disable
```

---

## 📚 相關文檔

遷移後，請查看以下更新的文檔：

1. **[DATABASE_CONFIG.md](./docs/DATABASE_CONFIG.md)** - 完整資料庫配置說明
2. **[QUICK_START_QUEUE.md](./docs/QUICK_START_QUEUE.md)** - Queue 快速開始
3. **[README.md](./README.md)** - 主文檔

---

## ⚠️ 注意事項

### 1. 舊資料庫處理

遷移完成後，`codex_teams` 資料庫仍然存在作為備份。確認新資料庫正常運作後，可以選擇刪除：

```bash
# ⚠️ 謹慎操作！確認數據已完整遷移
docker exec teamsnotify-postgres psql -U teamsnotify -c "DROP DATABASE codex_teams;"
```

### 2. 生產環境

如果在生產環境執行遷移：

1. 📋 **備份數據** - 先完整備份舊資料庫
2. 🛑 **停止服務** - 確保沒有活動連接
3. 🔄 **執行遷移** - 使用腳本或手動步驟
4. ✅ **驗證完整性** - 檢查所有表和數據
5. 🚀 **啟動服務** - 更新環境變數並重啟

### 3. CI/CD 管道

更新你的 CI/CD 配置中的資料庫名稱：

- GitHub Actions
- GitLab CI
- Jenkins
- Docker Compose
- Kubernetes ConfigMaps/Secrets

---

## 🎯 完成檢查清單

- [x] 更新核心程式碼 (main.go)
- [x] 更新啟動腳本 (start_server.sh)
- [x] 更新所有文檔
- [x] 建立資料庫遷移腳本
- [x] 重新編譯應用程式
- [ ] 執行資料庫遷移腳本
- [ ] 驗證資料庫連接
- [ ] 測試所有 API 端點
- [ ] 測試 Queue 功能
- [ ] 更新 .env 文件
- [ ] 更新部署配置
- [ ] 通知團隊成員

---

## 📞 需要幫助？

如果遇到問題：

1. 查看 [DATABASE_CONFIG.md](./docs/DATABASE_CONFIG.md) 的故障排查章節
2. 檢查 server.log 日誌文件
3. 驗證 Docker 容器狀態
4. 確認環境變數設置正確

---

**遷移日期**: 2025-10-01  
**執行者**: TeamsNotify Team  
**狀態**: ✅ 代碼已更新，等待執行資料庫遷移

