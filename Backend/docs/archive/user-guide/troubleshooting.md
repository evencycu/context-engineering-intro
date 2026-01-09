# 故障排除指南

本指南幫助您解決 Teams Notification API 使用過程中遇到的常見問題。

## 常見問題

### 1. 服務無法啟動

#### 問題描述
服務啟動時出現錯誤或無法訪問。

#### 可能原因
- 端口被佔用
- 環境變數未設定
- 資料庫連線失敗

#### 解決方案
```bash
# 檢查端口是否被佔用
lsof -i :8080

# 停止佔用端口的程序
pkill -f "./server"

# 檢查環境變數
echo $TEAMS_BOT_APP_ID
echo $TEAMS_TENANT_ID
echo $TEAMS_BOT_APP_PASSWORD

# 重新啟動服務
bash start_server.sh
```

### 2. 資料庫連線失敗

#### 問題描述
出現資料庫連線錯誤。

#### 可能原因
- PostgreSQL 容器未啟動
- 資料庫憑證錯誤
- 網路連線問題

#### 解決方案
```bash
# 檢查 Docker 容器狀態
docker ps | grep postgres

# 啟動 PostgreSQL 容器
docker-compose up -d postgres

# 測試資料庫連線
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT 1;"

# 檢查資料庫是否存在
docker exec teamsnotify-postgres psql -U teamsnotify -l
```

### 3. Teams 認證失敗

#### 問題描述
發送通知時出現認證錯誤。

#### 可能原因
- Bot 憑證錯誤
- Teams 應用程式未註冊
- 權限設定問題

#### 解決方案
```bash
# 檢查環境變數
export TEAMS_BOT_APP_ID=844146d7-4ac9-4e4d-a463-d6e027714e81
export TEAMS_TENANT_ID=051cece0-e4dc-4aed-b471-bf29824e1ee6
export TEAMS_BOT_APP_PASSWORD='HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy'

# 重新啟動服務
bash start_server.sh

# 檢查 Bot 註冊狀態
curl http://localhost:8080/api/v1/bots/platform
```

### 4. 通知發送失敗

#### 問題描述
通知無法成功發送到 Teams。

#### 可能原因
- Bot 未安裝到目標對話
- 對話 ID 錯誤
- 權限不足

#### 解決方案
```bash
# 檢查 Bot 安裝狀態
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "
SELECT conversation_id, installation_status 
FROM bot_installations 
WHERE installation_status = 'active';"

# 更新 Bot 安裝狀態
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "
UPDATE bot_installations 
SET installation_status = 'active' 
WHERE installation_status = 'stale';"

# 檢查目的地配置
curl http://localhost:8080/api/v1/destinations
```

### 5. API 回應錯誤

#### 問題描述
API 請求返回錯誤狀態碼。

#### 常見錯誤碼
- **400 Bad Request**: 請求參數錯誤
- **401 Unauthorized**: 認證失敗
- **403 Forbidden**: 權限不足
- **404 Not Found**: 資源不存在
- **429 Too Many Requests**: 請求過於頻繁
- **500 Internal Server Error**: 內部錯誤

#### 解決方案
```bash
# 檢查服務健康狀態
curl http://localhost:8080/health

# 查看服務日誌
tail -f server.log

# 檢查資料庫狀態
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "
SELECT COUNT(*) FROM notifications;"
```

## 除錯工具

### 1. 日誌查看
```bash
# 查看服務日誌
tail -f server.log

# 查看 Docker 容器日誌
docker logs teamsnotify-postgres
docker logs teamsnotify-redis

# 查看特定時間的日誌
grep "2025-10-02" server.log
```

### 2. 資料庫查詢
```bash
# 連接到資料庫
docker exec -it teamsnotify-postgres psql -U teamsnotify -d notification_center

# 檢查表結構
\d+ notifications
\d+ notification_destinations

# 查詢資料
SELECT * FROM notifications ORDER BY created_at DESC LIMIT 10;
```

### 3. Redis 監控
```bash
# 連接到 Redis
docker exec -it teamsnotify-redis redis-cli

# 查看所有鍵
KEYS *

# 查看佇列狀態
ZRANGE notification_queue 0 -1 WITHSCORES
```

### 4. 網路診斷
```bash
# 檢查端口監聽
netstat -tlnp | grep :8080

# 測試本地連線
curl -v http://localhost:8080/health

# 檢查防火牆設定
sudo ufw status
```

## 效能問題

### 1. 回應時間過慢

#### 可能原因
- 資料庫查詢慢
- 網路延遲
- 資源不足

#### 解決方案
```bash
# 檢查資料庫效能
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "
EXPLAIN ANALYZE SELECT * FROM notifications WHERE project_id = 'xxx';"

# 檢查系統資源
top
htop

# 檢查記憶體使用
free -h
```

### 2. 記憶體洩漏

#### 症狀
- 記憶體使用持續增長
- 系統變慢
- 服務崩潰

#### 解決方案
```bash
# 監控記憶體使用
watch -n 1 'ps aux | grep server'

# 重啟服務
pkill -f "./server"
bash start_server.sh

# 檢查 Go 記憶體設定
go run -race cmd/server/main.go
```

## 安全問題

### 1. 認證繞過

#### 檢查項目
- API Key 是否正確設定
- JWT Token 是否有效
- 權限檢查是否正常

#### 解決方案
```bash
# 重新生成 API Key
curl -X PATCH http://localhost:8080/api/v1/users/{user_id}/api-key

# 檢查認證中間件
grep -r "auth" internal/api/middleware/
```

### 2. 資料洩漏

#### 檢查項目
- 敏感資料是否加密
- 日誌是否包含敏感資訊
- 資料庫連線是否安全

#### 解決方案
```bash
# 檢查環境變數
env | grep -i password
env | grep -i secret

# 檢查日誌內容
grep -i "password\|secret" server.log
```

## 監控和警報

### 1. 健康檢查
```bash
# 自動健康檢查腳本
#!/bin/bash
while true; do
    if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "Service is down at $(date)"
        # 發送警報
    fi
    sleep 30
done
```

### 2. 效能監控
```bash
# 監控 API 回應時間
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/health

# curl-format.txt 內容:
#      time_namelookup:  %{time_namelookup}\n
#         time_connect:  %{time_connect}\n
#      time_appconnect:  %{time_appconnect}\n
#     time_pretransfer:  %{time_pretransfer}\n
#        time_redirect:  %{time_redirect}\n
#   time_starttransfer:  %{time_starttransfer}\n
#                      ----------\n
#           time_total:  %{time_total}\n
```

## 預防措施

### 1. 定期維護
- 定期清理舊日誌
- 監控磁碟空間
- 更新依賴套件

### 2. 備份策略
- 定期備份資料庫
- 備份配置檔案
- 測試恢復程序

### 3. 監控設定
- 設定資源使用警報
- 監控錯誤率
- 追蹤效能指標

## 聯絡支援

如果問題無法解決，請提供以下資訊：

1. **錯誤訊息**: 完整的錯誤日誌
2. **環境資訊**: 作業系統、Go 版本、Docker 版本
3. **重現步驟**: 詳細的重現步驟
4. **系統狀態**: 服務狀態、資源使用情況
5. **相關日誌**: 服務日誌、資料庫日誌

### 日誌收集腳本
```bash
#!/bin/bash
# 收集診斷資訊
echo "=== System Info ===" > diagnostic.log
uname -a >> diagnostic.log
docker version >> diagnostic.log
go version >> diagnostic.log

echo "=== Service Status ===" >> diagnostic.log
ps aux | grep server >> diagnostic.log
curl http://localhost:8080/health >> diagnostic.log

echo "=== Database Status ===" >> diagnostic.log
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT COUNT(*) FROM notifications;" >> diagnostic.log

echo "=== Recent Logs ===" >> diagnostic.log
tail -100 server.log >> diagnostic.log
```
