# External API 使用指南

External API 是供外部系統使用的通知發送介面，透過 `notify_key` 來識別專案並發送通知到預先配置的 Teams 目的地。

## 概述

External API 提供簡潔的 RESTful 介面，讓外部系統能夠輕鬆發送通知到 Microsoft Teams，無需管理複雜的 Bot 配置和認證。

## 基礎 URL

```
http://localhost:8080/api/v1/external
```

## 認證方式

External API 使用 `notify_key` 進行認證，無需額外的 API Key 或 JWT Token。

## API 端點

### 發送通知

**POST** `/api/v1/external/notify`

發送通知到預先配置的 Teams 目的地。

#### 請求參數

| 參數 | 類型 | 必填 | 說明 |
|------|------|------|------|
| `notify_key` | string | ✅ | 專案的 notify_key |
| `message` | string | ✅ | 通知內容 |
| `message_type` | string | ❌ | 訊息類型 (text/file/adaptive_card) |
| `priority` | string | ❌ | 優先級 (low/normal/high) |
| `targets` | array | ❌ | 目標篩選 (["all"] 或特定目標) |
| `mentions` | array | ❌ | 提及對象 |
| `metadata` | object | ❌ | 自定義元數據 |

#### 請求範例

**基本通知**
```bash
curl -X POST http://localhost:8080/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38",
    "message": "Hello World!"
  }'
```

**完整參數通知**
```bash
curl -X POST http://localhost:8080/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38",
    "message": "系統維護通知：將於今晚 10:00-11:00 進行系統維護",
    "message_type": "text",
    "priority": "high",
    "targets": ["all"],
    "mentions": ["@everyone"],
    "metadata": {
      "maintenance_id": "maint-001",
      "scheduled_time": "2025-10-02T22:00:00Z",
      "source": "monitoring-system"
    }
  }'
```

**發送到特定目標**
```bash
curl -X POST http://localhost:8080/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38",
    "message": "頻道專用訊息",
    "targets": ["19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2"]
  }'
```

#### 回應格式

**成功回應**
```json
{
  "success": true,
  "data": {
    "notification_id": "3f364043-eb50-49a6-988e-79a5be8f5ea0",
    "status": "sent",
    "message": "Notification processed successfully",
    "project_id": "f6a110f8-27f9-477e-a651-6f8cdbc6dca5",
    "project_name": "5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38",
    "destinations_count": 2,
    "results": [
      {
        "destination_id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
        "destination_name": "",
        "success": true
      },
      {
        "destination_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
        "destination_name": "",
        "success": true
      }
    ],
    "estimated_delivery": "5-10 minutes"
  }
}
```

**錯誤回應**
```json
{
  "success": false,
  "error": "Invalid notify_key",
  "details": "The provided notify_key does not exist or is inactive"
}
```

## 使用場景

### 1. 系統監控警報

```bash
# 發送 CPU 使用率警報
curl -X POST http://localhost:8080/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "monitoring-alerts",
    "message": "🚨 CPU 使用率過高：95%",
    "priority": "high",
    "metadata": {
      "alert_type": "cpu_usage",
      "threshold": 90,
      "current_value": 95,
      "server": "web-server-01"
    }
  }'
```

### 2. 部署通知

```bash
# 發送部署完成通知
curl -X POST http://localhost:8080/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "deployment-notifications",
    "message": "✅ 部署完成：v2.1.0 已成功部署到生產環境",
    "priority": "normal",
    "metadata": {
      "version": "v2.1.0",
      "environment": "production",
      "deployment_time": "2025-10-02T14:30:00Z",
      "deployed_by": "ci-cd-pipeline"
    }
  }'
```

### 3. 業務通知

```bash
# 發送訂單通知
curl -X POST http://localhost:8080/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "order-notifications",
    "message": "🛒 新訂單：訂單 #12345 已建立，金額：$299.99",
    "priority": "normal",
    "mentions": ["@sales-team"],
    "metadata": {
      "order_id": "12345",
      "amount": 299.99,
      "customer": "John Doe",
      "product": "Premium Package"
    }
  }'
```

## 參數詳解

### notify_key

專案的唯一識別碼，用於：
- 識別要發送通知的專案
- 取得預先配置的目的地
- 套用專案相關的設定（限制、優先級等）

### message

通知的主要內容，支援：
- 純文字訊息
- 表情符號和特殊字符
- 多行文字
- 基本 Markdown 格式（由 Teams 處理）

### message_type

| 類型 | 說明 | 範例 |
|------|------|------|
| `text` | 純文字訊息（預設） | "Hello World" |
| `file` | 檔案附件 | 需要額外的 `attachment` 參數 |
| `adaptive_card` | 自適應卡片 | 需要額外的 `adaptive_card` 參數 |

### priority

| 優先級 | 說明 | 處理順序 |
|--------|------|----------|
| `low` | 低優先級 | 最後處理 |
| `normal` | 一般優先級（預設） | 正常處理 |
| `high` | 高優先級 | 優先處理 |

### targets

目標篩選選項：
- `["all"]` - 發送到所有配置的目的地（預設）
- `["conversation_id1", "conversation_id2"]` - 發送到特定對話
- `["email1@example.com", "email2@example.com"] - 發送到特定電子郵件

### mentions

提及對象，會在 Teams 中高亮顯示：
- `["@everyone"]` - 提及所有人
- `["@張三", "@李四"]` - 提及特定人員
- `["@channel"]` - 提及整個頻道

### metadata

自定義元數據，用於：
- 追蹤通知來源
- 儲存業務相關資訊
- 支援後續查詢和分析

## 錯誤處理

### 常見錯誤碼

| 錯誤碼 | 說明 | 解決方案 |
|--------|------|----------|
| `400` | 請求參數錯誤 | 檢查必填參數和格式 |
| `404` | notify_key 不存在 | 確認 notify_key 正確 |
| `429` | 請求過於頻繁 | 降低發送頻率 |
| `500` | 內部服務器錯誤 | 稍後重試或聯絡支援 |

### 錯誤回應範例

```json
{
  "success": false,
  "error": "Invalid notify_key",
  "details": "The notify_key 'invalid-key' does not exist",
  "code": "INVALID_NOTIFY_KEY"
}
```

## 最佳實踐

### 1. 錯誤處理

```bash
#!/bin/bash
# 發送通知並處理錯誤

send_notification() {
  local notify_key="$1"
  local message="$2"
  
  response=$(curl -s -X POST http://localhost:8080/api/v1/external/notify \
    -H "Content-Type: application/json" \
    -d "{
      \"notify_key\": \"$notify_key\",
      \"message\": \"$message\"
    }")
  
  success=$(echo "$response" | jq -r '.success')
  
  if [ "$success" = "true" ]; then
    echo "✅ 通知發送成功"
    echo "$response" | jq '.data.notification_id'
  else
    echo "❌ 通知發送失敗"
    echo "$response" | jq -r '.error'
    exit 1
  fi
}

# 使用範例
send_notification "my-project" "測試訊息"
```

### 2. 批次發送

```bash
#!/bin/bash
# 批次發送多個通知

messages=(
  "系統狀態正常"
  "資料庫連線穩定"
  "API 回應時間正常"
)

for message in "${messages[@]}"; do
  curl -X POST http://localhost:8080/api/v1/external/notify \
    -H "Content-Type: application/json" \
    -d "{
      \"notify_key\": \"system-status\",
      \"message\": \"$message\",
      \"priority\": \"normal\"
    }" > /dev/null
  
  echo "已發送：$message"
  sleep 1  # 避免請求過於頻繁
done
```

### 3. 監控和重試

```bash
#!/bin/bash
# 發送通知並監控狀態

notify_key="monitoring-alerts"
message="系統警報測試"

# 發送通知
response=$(curl -s -X POST http://localhost:8080/api/v1/external/notify \
  -H "Content-Type: application/json" \
  -d "{
    \"notify_key\": \"$notify_key\",
    \"message\": \"$message\",
    \"priority\": \"high\"
  }")

notification_id=$(echo "$response" | jq -r '.data.notification_id')

if [ "$notification_id" != "null" ]; then
  echo "通知已發送，ID: $notification_id"
  
  # 等待處理完成
  sleep 5
  
  # 檢查狀態（需要額外的 API 端點）
  echo "檢查通知狀態..."
else
  echo "發送失敗：$(echo "$response" | jq -r '.error')"
fi
```

## 限制和配額

### 專案限制

每個專案都有以下限制：
- **每日限制**: 預設 10,000 則通知
- **每月限制**: 預設 300,000 則通知
- **速率限制**: 100 請求/秒

### 訊息限制

- **訊息長度**: 最大 8,000 字符
- **提及數量**: 最多 50 個提及
- **元數據大小**: 最大 1KB

## 監控和日誌

### 健康檢查

```bash
# 檢查服務狀態
curl http://localhost:8080/health
```

### 日誌查詢

```bash
# 查看服務日誌
tail -f server.log | grep "external/notify"
```

## 相關文檔

- [API 概覽](overview.md) - 完整的 API 說明
- [Provision API](provision-api.md) - 專案和目的地管理
- [快速開始](../getting-started/quick-start.md) - 快速上手指南
- [故障排除](../user-guide/troubleshooting.md) - 常見問題解決