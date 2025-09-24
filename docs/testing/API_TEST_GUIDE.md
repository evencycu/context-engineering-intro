## Local server start (with Teams bot env)

Use this to start the API server locally with the Teams Bot credentials and log output to a file for debugging:

```bash
export TEAMS_BOT_APP_ID=844146d7-4ac9-4e4d-a463-d6e027714e81
export TEAMS_TENANT_ID=051cece0-e4dc-4aed-b471-bf29824e1ee6
export TEAMS_BOT_APP_PASSWORD='HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy'
nohup env TEAMS_BOT_APP_ID=$TEAMS_BOT_APP_ID TEAMS_TENANT_ID=$TEAMS_TENANT_ID TEAMS_BOT_APP_PASSWORD=$TEAMS_BOT_APP_PASSWORD \
  go run ./cmd/server > /tmp/teamsnotify_server.log 2>&1 & echo $!
```

To stop the server listening on 8080 quickly:

```bash
lsof -t -i :8080 | xargs kill -9 2>/dev/null
```

# Teams Notification API 測試指南

## 概述

本文件提供了 Teams Notification Bot 平台所有 API 端點的測試方法和範例。

## 環境準備

### 啟動服務器

```bash
# 編譯並啟動服務器
go run ./cmd/server

# 或使用 Docker Compose
docker-compose up -d
```

### 健康檢查

```bash
curl -sS http://localhost:8080/health | jq .
```

## API 測試方法

### 1. Company API 測試

#### 1.1 獲取所有公司

```bash
curl -sS http://localhost:8080/api/v1/companies | jq .
```

#### 1.2 創建新公司

```bash
curl -sS -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "測試公司",
    "contact_email": "test@company.com",
    "contact_phone": "+886-2-1234-5678",
    "address": "台北市信義區信義路五段7號",
    "status": "active",
    "billing_enabled": true
  }' | jq .
```

#### 1.2.1 創建國泰投信公司 (實際範例)

```bash
# 新增國泰投信公司
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "國泰投信",
    "contact_email": "admin@cathaysite.com.tw",
    "contact_phone": "+1-555-0104",
    "address": "777 Enterprise Cathay, Taipei City, Taiwan",
    "billing_enabled": true
  }'
```

**回應範例：**

```json
{
  "data": {
    "id": "df842fd5-dfd0-44f0-b7a6-ceae9a62af3b",
    "created_at": "2025-09-22T03:42:20.425588Z",
    "updated_at": "2025-09-22T03:42:20.425589Z",
    "name": "國泰投信",
    "contact_email": "admin@cathaysite.com.tw",
    "contact_phone": "+1-555-0104",
    "address": "777 Enterprise Cathay, Taipei City, Taiwan",
    "status": "active",
    "billing_enabled": true
  },
  "message": "Company created successfully"
}
```

#### 1.3 獲取特定公司

```bash
curl -sS http://localhost:8080/api/v1/companies/{company_id} | jq .
```

#### 1.4 更新公司

```bash
curl -sS -X PUT http://localhost:8080/api/v1/companies/{company_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新後的公司名稱",
    "contact_email": "updated@company.com",
    "status": "active"
  }' | jq .
```

#### 1.5 更新公司狀態

```bash
curl -sS -X PATCH http://localhost:8080/api/v1/companies/{company_id}/status \
  -H "Content-Type: application/json" \
  -d '{"status": "inactive"}' | jq .
```

#### 1.6 刪除公司

```bash
curl -sS -X DELETE http://localhost:8080/api/v1/companies/{company_id} | jq .
```

### 2. User API 測試

#### 2.1 獲取所有用戶

```bash
curl -sS http://localhost:8080/api/v1/users | jq .
```

#### 2.2 創建新用戶

```bash
curl -sS -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "9589aad2-f5ec-4d25-8ef7-b818b1933c32",
    "email": "user@example.com",
    "name": "測試用戶",
    "password": "password123",
    "role": "user",
    "status": "active"
  }' | jq .
```

#### 2.2.1 創建國泰投信用戶 DoDoMan (實際範例)

```bash
# 新增國泰投信用戶 - 使用實際的資料庫 ID
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "df842fd5-dfd0-44f0-b7a6-ceae9a62af3b",
    "email": "admin@cathaysite.com.tw",
    "name": "DoDoMan",
    "role": "user",
    "password": "password123"
  }'
```

**回應範例：**

```json
{
  "data": {
    "id": "438bad46-a020-4ed7-a64a-7a3ce0432f7a",
    "created_at": "2025-09-22T03:44:34.834365Z",
    "updated_at": "2025-09-22T03:44:34.834365Z",
    "company_id": "df842fd5-dfd0-44f0-b7a6-ceae9a62af3b",
    "email": "admin@cathaysite.com.tw",
    "name": "DoDoMan",
    "role": "user",
    "status": "active",
    "last_login_at": null,
    "api_key_expires_at": null
  },
  "message": "User created successfully"
}
```

#### 2.3 獲取特定用戶

```bash
curl -sS http://localhost:8080/api/v1/users/{user_id} | jq .
```

#### 2.4 更新用戶

```bash
curl -sS -X PUT http://localhost:8080/api/v1/users/{user_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新後的用戶名稱",
    "email": "updated@example.com",
    "role": "admin"
  }' | jq .
```

#### 2.5 修改密碼

```bash
curl -sS -X PATCH http://localhost:8080/api/v1/users/{user_id}/password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword456"
  }' | jq .
```

#### 2.6 生成 API Key

```bash
curl -sS -X PATCH http://localhost:8080/api/v1/users/{user_id}/api-key | jq .
```

#### 2.7 按公司獲取用戶

```bash
curl -sS http://localhost:8080/api/v1/users/company/{company_id} | jq .
```

### 3. Project API 測試

#### 3.1 獲取所有項目

```bash
curl -sS http://localhost:8080/api/v1/projects | jq .
```

#### 3.2 創建新項目

```bash
curl -sS -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "9589aad2-f5ec-4d25-8ef7-b818b1933c32",
    "key_name": "test-project-2",
    "description": "第二個測試項目",
    "status": "active",
    "daily_limit": 2000,
    "monthly_limit": 60000,
    "priority": "high",
    "created_by": "643c4d7a-a18f-4caa-aff2-5a0d4439b367"
  }' | jq .
```

#### 3.3 獲取特定項目

```bash
curl -sS http://localhost:8080/api/v1/projects/{project_id} | jq .
```

#### 3.4 更新項目

```bash
curl -sS -X PUT http://localhost:8080/api/v1/projects/{project_id} \
  -H "Content-Type: application/json" \
  -d '{
    "description": "更新後的項目描述",
    "daily_limit": 3000,
    "monthly_limit": 90000
  }' | jq .
```

#### 3.5 更新項目限制

```bash
curl -sS -X PATCH http://localhost:8080/api/v1/projects/{project_id}/limits \
  -H "Content-Type: application/json" \
  -d '{
    "daily_limit": 5000,
    "monthly_limit": 150000
  }' | jq .
```

#### 3.6 按公司獲取項目

```bash
curl -sS http://localhost:8080/api/v1/projects/company/{company_id} | jq .
```

#### 3.7 按 Key Name 獲取項目

```bash
curl -sS http://localhost:8080/api/v1/projects/key/{key_name} | jq .
```

### 4. Bot API 測試

#### 4.1 Platform Bot 測試

##### 獲取所有平台機器人

```bash
curl -sS http://localhost:8080/api/v1/bots/platform | jq .
```

##### 創建平台機器人

```bash
curl -sS -X POST http://localhost:8080/api/v1/bots/platform \
  -H "Content-Type: application/json" \
  -d '{
    "name": "新平台機器人",
    "description": "新的平台機器人描述",
    "app_id": "new-app-456",
    "app_password": "password123",
    "tenant_id": "new-tenant-456",
    "status": "active",
    "webhook_url": "https://new-bot.com/webhook",
    "capabilities": {
      "receive_message": true,
      "send_message": true,
      "file_upload": true
    },
    "rate_limit_per_minute": 150,
    "max_concurrent_requests": 15
  }' | jq .
```

##### 獲取特定平台機器人

```bash
curl -sS http://localhost:8080/api/v1/bots/platform/{bot_id} | jq .
```

##### 更新平台機器人

```bash
curl -sS -X PUT http://localhost:8080/api/v1/bots/platform/{bot_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新後的平台機器人",
    "description": "更新後的描述",
    "capabilities": {
      "receive_message": true,
      "send_message": true,
      "file_upload": true,
      "adaptive_cards": true
    }
  }' | jq .
```

##### 更新機器人狀態

```bash
curl -sS -X PATCH http://localhost:8080/api/v1/bots/platform/{bot_id}/status \
  -H "Content-Type: application/json" \
  -d '{"status": "inactive"}' | jq .
```

##### 測試機器人連接

```bash
curl -sS -X POST http://localhost:8080/api/v1/bots/platform/{bot_id}/test | jq .
```

#### 4.2 Third Party Bot 測試

##### 獲取所有第三方機器人

```bash
curl -sS http://localhost:8080/api/v1/bots/third-party | jq .
```

##### 創建第三方機器人

```bash
curl -sS -X POST http://localhost:8080/api/v1/bots/third-party \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "9589aad2-f5ec-4d25-8ef7-b818b1933c32",
    "name": "新第三方機器人",
    "description": "新的第三方機器人描述",
    "app_id": "third-party-app-456",
    "app_password": "password123",
    "tenant_id": "third-party-tenant-456",
    "status": "active",
    "webhook_url": "https://third-party-bot.com/webhook",
    "api_endpoint": "https://third-party-bot.com/api",
    "api_key": "api-key-123",
    "capabilities": {
      "receive_message": true,
      "send_message": true,
      "file_upload": true,
      "adaptive_cards": true
    },
    "rate_limit_per_minute": 300,
    "max_concurrent_requests": 30,
    "contact_email": "contact@third-party-bot.com",
    "contact_phone": "+886-2-9876-5432",
    "created_by": "643c4d7a-a18f-4caa-aff2-5a0d4439b367"
  }' | jq .
```

##### 獲取特定第三方機器人

```bash
curl -sS http://localhost:8080/api/v1/bots/third-party/{bot_id} | jq .
```

##### 更新第三方機器人

```bash
curl -sS -X PUT http://localhost:8080/api/v1/bots/third-party/{bot_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新後的第三方機器人",
    "description": "更新後的描述",
    "contact_email": "updated@third-party-bot.com"
  }' | jq .
```

##### 更新 API Key

```bash
curl -sS -X PATCH http://localhost:8080/api/v1/bots/third-party/{bot_id}/api-key \
  -H "Content-Type: application/json" \
  -d '{"api_key": "new-api-key-456"}' | jq .
```

##### 測試第三方機器人連接

```bash
curl -sS -X POST http://localhost:8080/api/v1/bots/third-party/{bot_id}/test | jq .
```

#### 4.3 通用 Bot 查詢

##### 按公司獲取機器人

```bash
curl -sS http://localhost:8080/api/v1/bots/company/{company_id} | jq .
```

##### 按狀態獲取機器人

```bash
curl -sS http://localhost:8080/api/v1/bots/status/{status} | jq .
```

### 5. Destination API 測試

#### 5.1 獲取所有目的地

```bash
curl -sS http://localhost:8080/api/v1/destinations | jq .
```

#### 5.2 創建新目的地

```bash
curl -sS -X POST http://localhost:8080/api/v1/destinations \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "198f1130-20a9-4c7c-a504-6a055d27e8db",
    "name": "測試目的地",
    "description": "測試目的地描述",
    "teams_tenant_id": "test-tenant-789",
    "targets": [
      {
        "type": "channel",
        "team_id": "team-456",
        "channel_id": "channel-789",
        "display_name": "測試頻道"
      },
      {
        "type": "person",
        "user_id": "user-456",
        "display_name": "測試用戶"
      }
    ],
    "status": "active",
    "validation_status": "pending",
    "created_by": "643c4d7a-a18f-4caa-aff2-5a0d4439b367"
  }' | jq .
```

#### 5.2.1 創建個人通知目的地 (實際範例)

```bash
# 新增個人通知目的地 - 使用實際的資料庫 ID 和 conversation ID
curl -X POST http://localhost:8080/api/v1/destinations \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "750e8400-e29b-41d4-a716-446655440001",
    "name": "Personal Conversation Target",
    "description": "Personal conversation destination for proactive messaging",
    "teams_tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
    "targets": [
      {
        "type": "person",
        "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
        "user_id": "29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw",
        "aad_object_id": "8d40db4b-935f-4e5a-aaae-ba86faad0e12",
        "display_name": "Test User",
        "description": "Personal conversation target for proactive messaging"
      }
    ],
    "bot_id": "850e8400-e29b-41d4-a716-446655440001",
    "bot_type": "platform",
    "created_by": "650e8400-e29b-41d4-a716-446655440001"
  }'
```

**回應範例：**

```json
{
  "data": {
    "id": "922c78f0-382c-4893-bae3-dc4bd38d7a0b",
    "created_at": "2025-09-22T05:17:24.270033Z",
    "updated_at": "2025-09-22T05:17:24.270033Z",
    "project_id": "750e8400-e29b-41d4-a716-446655440001",
    "name": "Personal Conversation Target",
    "description": "Personal conversation destination for proactive messaging",
    "teams_tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
    "targets": [
      {
        "type": "person",
        "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
        "user_id": "29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw",
        "aad_object_id": "8d40db4b-935f-4e5a-aaae-ba86faad0e12",
        "display_name": "Test User",
        "description": "Personal conversation target for proactive messaging"
      }
    ],
    "bot_id": "850e8400-e29b-41d4-a716-446655440001",
    "bot_type": "platform",
    "status": "active",
    "validation_status": "pending",
    "last_validated_at": null,
    "created_by": "650e8400-e29b-41d4-a716-446655440001"
  },
  "message": "Destination created successfully"
}
```

**重要說明：**

- `conversation_id` 是發送 proactive message 的關鍵欄位
- 這個 ID 來自 Teams 的 installationUpdate 事件中的 `conversation.id`
- 用於 Bot Framework 的 proactive messaging API 端點

#### 5.3 獲取特定目的地

```bash
curl -sS http://localhost:8080/api/v1/destinations/{destination_id} | jq .
```

#### 5.4 更新目的地

```bash
curl -sS -X PUT http://localhost:8080/api/v1/destinations/{destination_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "更新後的目的地名稱",
    "description": "更新後的描述",
    "status": "active"
  }' | jq .
```

#### 5.5 更新目的地目標

```bash
curl -sS -X PATCH http://localhost:8080/api/v1/destinations/{destination_id}/targets \
  -H "Content-Type: application/json" \
  -d '{
    "targets": [
      {
        "type": "channel",
        "team_id": "team-789",
        "channel_id": "channel-123",
        "display_name": "新頻道"
      }
    ]
  }' | jq .
```

#### 5.6 驗證目的地目標

```bash
curl -sS -X POST http://localhost:8080/api/v1/destinations/{destination_id}/validate | jq .
```

#### 5.7 按項目獲取目的地

```bash
curl -sS http://localhost:8080/api/v1/destinations/project/{project_id} | jq .
```

#### 5.8 按機器人獲取目的地

```bash
curl -sS http://localhost:8080/api/v1/destinations/bot/{bot_id} | jq .
```

#### 5.9 搜索目的地

```bash
curl -sS "http://localhost:8080/api/v1/destinations/search?q=測試" | jq .
```

### 6. Notification API 測試

#### 6.1 獲取所有通知

```bash
curl -sS http://localhost:8080/api/v1/notifications | jq .
```

#### 6.2 發送通知

```bash
curl -sS -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "198f1130-20a9-4c7c-a504-6a055d27e8db",
    "sender_id": "643c4d7a-a18f-4caa-aff2-5a0d4439b367",
    "message_type": "text",
    "content": "這是一個測試通知消息 - 包含中文內容",
    "mentions": ["@admin", "@user", "@manager"],
    "priority": "high",
    "destinations": ["c5133b8b-6e5c-4359-be17-221483569b97"]
  }' | jq .
```

#### 6.3 發送帶附件的通知

```bash
curl -sS -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "198f1130-20a9-4c7c-a504-6a055d27e8db",
    "sender_id": "643c4d7a-a18f-4caa-aff2-5a0d4439b367",
    "message_type": "file",
    "content": "請查看附件",
    "attachment": {
      "type": "image",
      "url": "https://example.com/image.jpg",
      "name": "測試圖片.jpg",
      "size": 1024000
    },
    "priority": "normal",
    "destinations": ["c5133b8b-6e5c-4359-be17-221483569b97"]
  }' | jq .
```

#### 6.4 發送 Adaptive Card 通知

```bash
curl -sS -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "198f1130-20a9-4c7c-a504-6a055d27e8db",
    "sender_id": "643c4d7a-a18f-4caa-aff2-5a0d4439b367",
    "message_type": "adaptive_card",
    "content": "自適應卡片通知",
    "adaptive_card": {
      "type": "AdaptiveCard",
      "version": "1.0",
      "body": [
        {
          "type": "TextBlock",
          "text": "測試自適應卡片",
          "weight": "bolder",
          "size": "medium"
        }
      ]
    },
    "priority": "high",
    "destinations": ["c5133b8b-6e5c-4359-be17-221483569b97"]
  }' | jq .
```

#### 6.5 獲取特定通知

```bash
curl -sS http://localhost:8080/api/v1/notifications/{notification_id} | jq .
```

#### 6.6 重試通知

```bash
curl -sS -X POST http://localhost:8080/api/v1/notifications/{notification_id}/retry | jq .
```

#### 6.7 取消通知

```bash
curl -sS -X DELETE http://localhost:8080/api/v1/notifications/{notification_id} | jq .
```

#### 6.8 按項目獲取通知

```bash
curl -sS http://localhost:8080/api/v1/notifications/project/{project_id} | jq .
```

#### 6.9 按發送者獲取通知

```bash
curl -sS http://localhost:8080/api/v1/notifications/sender/{sender_id} | jq .
```

#### 6.10 按狀態獲取通知

```bash
curl -sS http://localhost:8080/api/v1/notifications/status/{status} | jq .
```

#### 6.11 按日期範圍獲取通知

```bash
curl -sS "http://localhost:8080/api/v1/notifications/date-range?start_date=2025-09-01&end_date=2025-09-30" | jq .
```

### 7. Messages / 主動訊息（Proactive）測試

本節提供臨時測試端點，用於不經過資料庫，直接以 Bot Framework 主動發送簡訊給指定的 Teams 會話。

#### 7.1 必填環境變數

- `TEAMS_BOT_APP_ID`: 機器人 App ID
- `TEAMS_BOT_APP_PASSWORD`: 機器人密鑰

請在啟動服務前匯出環境變數：

```bash
export TEAMS_BOT_APP_ID="<your-bot-app-id>"
export TEAMS_BOT_APP_PASSWORD="<your-bot-app-password>"
go run ./cmd/server
```

#### 7.2 最小可用請求 Payload（需自行帶入實際值）

```json
{
  "activity": {
    "type": "installationUpdate",
    "channelId": "msteams",
    "serviceUrl": "https://smba.trafficmanager.net/apac/<tenantId>/",
    "conversation": {
      "conversationType": "personal",
      "id": "<conversationId>",
      "tenantId": "<tenantId>"
    },
    "from": {
      "id": "<fromId>",
      "aadObjectId": "<aadObjectId>"
    },
    "recipient": {
      "id": "28:<TEAMS_BOT_APP_ID>"
    },
    "channelData": {
      "tenant": {
        "id": "<tenantId>"
      }
    }
  },
  "text": "這是一個主動訊息測試"
}
```

必填欄位與說明：

- `activity.serviceUrl`: Bot Framework endpoint（依租戶/區域）。
- `activity.channelId`: 固定 `msteams`。
- `activity.type`: 建議 `installationUpdate`，便於服務端解析欄位。
- `activity.conversation.id`: 目標會話 ID（個人/群組/頻道）。
- `activity.conversation.tenantId`: AAD 租戶 ID。
- `activity.channelData.tenant.id`: AAD 租戶 ID（與上方一致）。
- `activity.recipient.id`: Bot 成員 ID，通常為 `28:<TEAMS_BOT_APP_ID>`。
- `activity.from.id`: 來源成員 ID（常由 Teams 事件提供）。
- `text`: 要發送的訊息文字內容。

#### 7.3 範例：使用現有測試資料發送

```bash
curl -s -X POST http://localhost:8080/api/v1/messages/proactive/test \
  -H 'Content-Type: application/json' \
  -d @- <<'JSON'
{
  "activity": {
    "type": "installationUpdate",
    "channelId": "msteams",
    "serviceUrl": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/",
    "conversation": {
      "conversationType": "personal",
      "id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
      "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
    },
    "from": {
      "id": "29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw",
      "aadObjectId": "8d40db4b-935f-4e5a-aaae-ba86faad0e12"
    },
    "recipient": {
      "id": "28:844146d7-4ac9-4e4d-a463-d6e027714e81"
    },
    "channelData": {
      "tenant": {
        "id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
      }
    }
  },
  "text": "這是一個主動訊息測試"
}
JSON
```

成功條件速查：

- Token 取得：以 `tenantId` 組合 `https://login.microsoftonline.com/{tenantId}/oauth2/v2.0/token`。
- 發送目的地：使用 `serviceUrl` 與 `conversation.id`。
- 常見錯誤：
  - 未設定或錯誤的 `TEAMS_BOT_APP_ID/TEAMS_BOT_APP_PASSWORD`
  - `serviceUrl` 與租戶/區域不一致
  - `conversation.id` 過期或非該租戶

## 批量測試腳本

### 完整 API 測試腳本

```bash
#!/bin/bash

echo "=== Teams Notification API 完整測試 ==="
echo ""

# 健康檢查
echo "1. 健康檢查:"
curl -sS http://localhost:8080/health | jq -r '.status'
echo ""

# 統計所有 API
echo "2. API 統計:"
echo "   公司總數: $(curl -sS http://localhost:8080/api/v1/companies | jq '.data | length')"
echo "   用戶總數: $(curl -sS http://localhost:8080/api/v1/users | jq '.data | length')"
echo "   項目總數: $(curl -sS http://localhost:8080/api/v1/projects | jq '.data | length')"
echo "   平台機器人總數: $(curl -sS http://localhost:8080/api/v1/bots/platform | jq '.data | length')"
echo "   第三方機器人總數: $(curl -sS http://localhost:8080/api/v1/bots/third-party | jq '.data | length')"
echo "   目的地總數: $(curl -sS http://localhost:8080/api/v1/destinations | jq '.data | length')"
echo "   通知總數: $(curl -sS http://localhost:8080/api/v1/notifications | jq '.data | length')"
echo ""

echo "✅ 所有 API 測試完成！"
```

### 性能測試腳本

```bash
#!/bin/bash

echo "=== API 性能測試 ==="
echo ""

# 測試響應時間
echo "測試各 API 響應時間:"
echo "健康檢查: $(curl -sS -w "%{time_total}s" -o /dev/null http://localhost:8080/health)"
echo "公司 API: $(curl -sS -w "%{time_total}s" -o /dev/null http://localhost:8080/api/v1/companies)"
echo "用戶 API: $(curl -sS -w "%{time_total}s" -o /dev/null http://localhost:8080/api/v1/users)"
echo "項目 API: $(curl -sS -w "%{time_total}s" -o /dev/null http://localhost:8080/api/v1/projects)"
echo "平台機器人 API: $(curl -sS -w "%{time_total}s" -o /dev/null http://localhost:8080/api/v1/bots/platform)"
echo "第三方機器人 API: $(curl -sS -w "%{time_total}s" -o /dev/null http://localhost:8080/api/v1/bots/third-party)"
echo "目的地 API: $(curl -sS -w "%{time_total}s" -o /dev/null http://localhost:8080/api/v1/destinations)"
echo "通知 API: $(curl -sS -w "%{time_total}s" -o /dev/null http://localhost:8080/api/v1/notifications)"
```

## 錯誤處理測試

### 測試無效的 ID

```bash
# 測試不存在的公司 ID
curl -sS http://localhost:8080/api/v1/companies/00000000-0000-0000-0000-000000000000 | jq .

# 測試無效的 UUID 格式
curl -sS http://localhost:8080/api/v1/companies/invalid-uuid | jq .
```

### 測試無效的請求數據

```bash
# 測試缺少必填字段
curl -sS -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{"name": "測試公司"}' | jq .

# 測試無效的 JSON
curl -sS -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{"name": "測試公司", "status": "invalid"}' | jq .
```

## 數據驗證測試

### 測試 JSONB 字段

```bash
# 測試複雜的 targets 數組
curl -sS -X POST http://localhost:8080/api/v1/destinations \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "198f1130-20a9-4c7c-a504-6a055d27e8db",
    "name": "複雜目標測試",
    "description": "測試複雜的目標配置",
    "teams_tenant_id": "complex-tenant",
    "targets": [
      {
        "type": "channel",
        "team_id": "team-1",
        "channel_id": "channel-1",
        "display_name": "頻道1"
      },
      {
        "type": "person",
        "user_id": "user-1",
        "display_name": "用戶1"
      },
      {
        "type": "chatgroup",
        "group_id": "group-1",
        "display_name": "聊天群組1"
      }
    ],
    "status": "active",
    "validation_status": "pending",
    "created_by": "643c4d7a-a18f-4caa-aff2-5a0d4439b367"
  }' | jq .
```

### 測試中文內容

```bash
# 測試中文公司名稱
curl -sS -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "台灣科技股份有限公司",
    "contact_email": "info@taiwan-tech.com",
    "contact_phone": "+886-2-2345-6789",
    "address": "台北市信義區信義路五段7號101大樓",
    "status": "active",
    "billing_enabled": true
  }' | jq .
```

## 注意事項

1. **端口配置**: 默認使用 8080 端口，確保端口未被占用
2. **數據庫連接**: 確保 PostgreSQL 數據庫正在運行
3. **UUID 格式**: 所有 ID 必須使用有效的 UUID 格式
4. **JSON 格式**: 請求數據必須是有效的 JSON 格式
5. **必填字段**: 創建操作需要提供所有必填字段
6. **外鍵關聯**: 確保引用的 ID 在相關表中存在

## 故障排除

### 常見錯誤

1. **500 Internal Server Error**: 檢查服務器日誌
2. **404 Not Found**: 檢查 API 路徑和 ID
3. **400 Bad Request**: 檢查請求數據格式
4. **422 Unprocessable Entity**: 檢查數據驗證規則

### 日誌查看

```bash
# 查看服務器日誌
tail -f server.log

# 或直接查看終端輸出
go run ./cmd/server
```

## 更新日誌

- **v1.0.0** (2025-09-15): 初始版本，包含所有核心 API 測試方法
- 支持完整的 CRUD 操作
- 支持 JSONB 字段測試
- 支持中文內容測試
- 包含性能測試和錯誤處理測試
