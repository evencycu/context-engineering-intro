# 🗃️ ERD（資料實體關聯圖）

> 最小資料模型，支援通知任務、發送結果、安裝資訊。

```mermaid
erDiagram
    USERS {
      uuid id PK
      string aad_object_id "Entra ObjectId"
      string display_name
      string upn "email"
      timestamptz created_at
    }

    INSTALLATIONS {
      uuid id PK
      string scope "personal|groupChat|channel"
      string tenant_id
      string conversation_id
      string chat_id
      string team_id
      string channel_id
      uuid  user_id FK
      timestamptz installed_at
      timestamptz removed_at
    }

    NOTIFICATIONS {
      uuid id PK
      string request_id "idempotency key"
      string channel_type "personal|groupChat|channel"
      string destination_ref "chatId/teamId:channelId/userAadId"
      text   payload
      string status "queued|processing|sent|failed|retrying"
      timestamptz created_at
      timestamptz updated_at
      timestamptz sent_at
    }

    MESSAGE_LOGS {
      uuid id PK
      uuid notification_id FK
      string provider "graph|bot"
      int    http_status
      text   response_body
      string message_id
      string error_reason
      timestamptz created_at
    }

    USERS ||--o{ INSTALLATIONS : "has"
    NOTIFICATIONS ||--o{ MESSAGE_LOGS : "emits"
```

## 索引建議
- `INSTALLATIONS(conversation_id)`、`INSTALLATIONS(chat_id)`
- `NOTIFICATIONS(request_id)`（唯一）
- `MESSAGE_LOGS(notification_id)`

## 參考
- 安裝事件（installationUpdate）會提供 conversation/chat/team/channel 等 id，可寫入 INSTALLATIONS。
- 發送結果記到 MESSAGE_LOGS 供追蹤與稽核。

