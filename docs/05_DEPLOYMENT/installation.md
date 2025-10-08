# Teams Bot Framework Installation Update Events

本文檔記錄了 Teams Bot Framework 中不同類型的 installation update 事件和相關的 Graph API 查詢範例。

## Graph API 查詢範例

### Personal Chat
```bash
curl -X GET "https://graph.microsoft.com/v1.0/users/8d40db4b-935f-4e5a-aaae-ba86faad0e12" \
  -H "Authorization: Bearer $TOKEN"
```

**回應：**
```json
{
  "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#users/$entity",
  "businessPhones": [],
  "displayName": "even3e",
  "givenName": "even3",
  "jobTitle": null,
  "mail": "even3@cathaysinglelab.onmicrosoft.com",
  "mobilePhone": null,
  "officeLocation": null,
  "preferredLanguage": null,
  "surname": "GG",
  "userPrincipalName": "even3@cathaysinglelab.onmicrosoft.com",
  "id": "8d40db4b-935f-4e5a-aaae-ba86faad0e12"
}
```

```bash
curl -X GET "https://graph.microsoft.com/v1.0/chats/19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces" \
  -H "Authorization: Bearer $TOKEN"
```

**回應：**
```json
{
  "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#chats/$entity",
  "id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces",
  "topic": null,
  "createdDateTime": "2025-09-18T06:42:39.469Z",
  "lastUpdatedDateTime": "2025-09-18T06:42:40.79Z",
  "chatType": "oneOnOne",
  "webUrl": "https://teams.microsoft.com/l/chat/19%3A8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81%40unq.gbl.spaces/0?tenantId=051cece0-e4dc-4aed-b471-bf29824e1ee6",
  "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
  "isHiddenForAllMembers": false,
  "viewpoint": null,
  "onlineMeetingInfo": null
}
```

### Group Chat
```bash
curl -X GET "https://graph.microsoft.com/v1.0/chats/19:f26a8d8a235f430db87a404491cd2ffc@thread.v2" \
  -H "Authorization: Bearer $TOKEN"
```

**回應：**
```json
{
  "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#chats/$entity",
  "id": "19:f26a8d8a235f430db87a404491cd2ffc@thread.v2",
  "topic": "與 [Even] 的客服對話 - Support Request for User",
  "createdDateTime": "2025-07-04T08:53:16Z",
  "lastUpdatedDateTime": "2025-09-23T03:32:01.507Z",
  "chatType": "group",
  "webUrl": "https://teams.microsoft.com/l/chat/19%3Af26a8d8a235f430db87a404491cd2ffc%40thread.v2/0?tenantId=051cece0-e4dc-4aed-b471-bf29824e1ee6",
  "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
  "isHiddenForAllMembers": false,
  "viewpoint": null,
  "onlineMeetingInfo": null
}
```

### Channel
```bash
curl -X GET "https://graph.microsoft.com/v1.0/teams/2747cdba-5b58-4fdf-9cac-dc52e725e728/channels/19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2" \
  -H "Authorization: Bearer $TOKEN"
```

**回應：**
```json
{
  "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#teams('2747cdba-5b58-4fdf-9cac-dc52e725e728')/channels/$entity",
  "id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
  "createdDateTime": "2025-04-11T02:28:37.983Z",
  "displayName": "received Msg",
  "description": "Teams Send Msg Test Group",
  "isFavoriteByDefault": null,
  "email": "TeamsSendMsgTestGroup@cathay-ins.symphox.com",
  "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
  "webUrl": "https://teams.cloud.microsoft.com/l/channel/19%3Alg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1%40thread.tacv2/received%20Msg?groupId=2747cdba-5b58-4fdf-9cac-dc52e725e728&tenantId=051cece0-e4dc-4aed-b471-bf29824e1ee6&allowXTenantAccess=True",
  "membershipType": "standard",
  "isArchived": false
}
```

## Installation Update 事件範例

### 1. Personal Chat Installation Update

**事件類型：** `installationUpdate`  
**動作：** `add`  
**對話類型：** `personal`

```json
{
  "action": "add",
  "channelData": {
    "settings": {
      "selectedChannel": {
        "id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces"
      }
    },
    "source": {
      "name": "message"
    },
    "tenant": {
      "id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
    }
  },
  "channelId": "msteams",
  "conversation": {
    "conversationType": "personal",
    "id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
    "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
  },
  "entities": [
    {
      "locale": "zh-TW",
      "type": "clientInfo"
    }
  ],
  "from": {
    "aadObjectId": "8d40db4b-935f-4e5a-aaae-ba86faad0e12",
    "id": "29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw"
  },
  "id": "f:757cd071-d818-ed6a-a463-6d986aa0718d",
  "locale": "zh-TW",
  "recipient": {
    "id": "28:844146d7-4ac9-4e4d-a463-d6e027714e81",
    "name": "lab-test-teams-notify"
  },
  "serviceUrl": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/",
  "timestamp": "2025-09-18T06:42:40.192Z",
  "type": "installationUpdate"
}
```

### 2. Group Chat Installation Update

**事件類型：** `installationUpdate`  
**動作：** `add`  
**對話類型：** `groupChat`

```json
{
  "action": "add",
  "channelData": {
    "settings": {
      "selectedChannel": {
        "id": "19:f26a8d8a235f430db87a404491cd2ffc@thread.v2"
      }
    },
    "source": {
      "name": "message"
    },
    "tenant": {
      "id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
    }
  },
  "channelId": "msteams",
  "conversation": {
    "conversationType": "groupChat",
    "id": "19:f26a8d8a235f430db87a404491cd2ffc@thread.v2",
    "isGroup": true,
    "name": "與 [Even] 的客服對話 - Support Request for User",
    "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
  },
  "entities": [
    {
      "locale": "zh-TW",
      "type": "clientInfo"
    }
  ],
  "from": {
    "aadObjectId": "bd5ab633-3f6f-48ad-be04-e144480da48f",
    "id": "29:1FrpqFNZ44LuLurxUt5tpSed4vXWXTdi428VAbwv9FSK-cI4UWRtIKdwEwFKCehU48w1hrQvIc5AXAoUtNKKBYg"
  },
  "id": "f:0c2b988b-d01a-7ad8-bdb3-896c91a06b0e",
  "locale": "zh-TW",
  "recipient": {
    "id": "28:844146d7-4ac9-4e4d-a463-d6e027714e81",
    "name": "lab-test-teams-notify"
  },
  "serviceUrl": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/",
  "timestamp": "2025-09-23T03:32:01.535Z",
  "type": "installationUpdate"
}
```

### 3. Channel Installation Update

**事件類型：** `installationUpdate`  
**動作：** `add`  
**對話類型：** `channel`

```json
{
  "action": "add",
  "channelData": {
    "channel": {
      "id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2"
    },
    "settings": {
      "selectedChannel": {
        "id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2"
      }
    },
    "source": {
      "name": "message"
    },
    "team": {
      "aadGroupId": "2747cdba-5b58-4fdf-9cac-dc52e725e728",
      "id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
      "name": "Teams Send Msg Test Group"
    },
    "tenant": {
      "id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
    }
  },
  "channelId": "msteams",
  "conversation": {
    "conversationType": "channel",
    "id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
    "isGroup": true,
    "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
  },
  "entities": [
    {
      "locale": "zh-TW",
      "type": "clientInfo"
    }
  ],
  "from": {
    "aadObjectId": "bd5ab633-3f6f-48ad-be04-e144480da48f",
    "id": "29:1FrpqFNZ44LuLurxUt5tpSed4vXWXTdi428VAbwv9FSK-cI4UWRtIKdwEwFKCehU48w1hrQvIc5AXAoUtNKKBYg"
  },
  "id": "f:f0f9e059-9224-1cca-ae84-0e10d3d003dd",
  "locale": "zh-TW",
  "recipient": {
    "id": "28:844146d7-4ac9-4e4d-a463-d6e027714e81",
    "name": "lab-test-teams-notify"
  },
  "serviceUrl": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/",
  "timestamp": "2025-09-23T03:31:05.255Z",
  "type": "installationUpdate"
}
```

## 重要欄位說明

### 通用欄位
- **`action`**: 安裝動作 (`add` 或 `remove`)
- **`conversation.conversationType`**: 對話類型 (`personal`, `groupChat`, `channel`)
- **`conversation.id`**: 對話 ID，用於識別特定的對話
- **`conversation.tenantId`**: 租戶 ID
- **`serviceUrl`**: Bot Framework 服務 URL
- **`from.aadObjectId`**: 發送者的 Azure AD 物件 ID
- **`recipient.id`**: 接收者 ID (通常是 Bot 的 ID)
- **`recipient.name`**: 接收者名稱 (通常是 Bot 名稱)

### 特殊欄位
- **Personal Chat**: 使用 `conversation.id` 作為唯一識別
- **Group Chat**: 包含 `conversation.name` 和 `isGroup: true`
- **Channel**: 包含 `team` 資訊和 `channel` 資訊

## 使用方式

這些事件會透過 Teams Bot Framework 的 webhook 端點發送到你的應用程式。你的應用程式需要：

1. 解析 `installationUpdate` 事件
2. 根據 `action` 欄位決定是新增還是移除安裝
3. 儲存或更新 `bot_installations` 表格中的記錄
4. 使用 `conversation.id` 和 `conversationType` 進行後續的訊息發送

## 測試區（curl 範例）

以下三個範例可直接打到本地端的 `/api/v1/messages` 端點，模擬 Teams 的 installationUpdate 事件。

### Personal Chat 測試
```bash
curl -X POST "http://localhost:8080/api/v1/messages" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "add",
    "channelData": {
      "settings": {
        "selectedChannel": {
          "id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces"
        }
      },
      "source": { "name": "message" },
      "tenant": { "id": "051cece0-e4dc-4aed-b471-bf29824e1ee6" }
    },
    "channelId": "msteams",
    "conversation": {
      "conversationType": "personal",
      "id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
      "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
    },
    "entities": [ { "locale": "zh-TW", "type": "clientInfo" } ],
    "from": {
      "aadObjectId": "8d40db4b-935f-4e5a-aaae-ba86faad0e12",
      "id": "29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw"
    },
    "id": "f:757cd071-d818-ed6a-a463-6d986aa0718d",
    "locale": "zh-TW",
    "recipient": { "id": "28:844146d7-4ac9-4e4d-a463-d6e027714e81", "name": "lab-test-teams-notify" },
    "serviceUrl": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/",
    "timestamp": "2025-09-18T06:42:40.192Z",
    "type": "installationUpdate"
  }'
```

### Group Chat 測試
```bash
curl -X POST "http://localhost:8080/api/v1/messages" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "add",
    "channelData": {
      "settings": { "selectedChannel": { "id": "19:f26a8d8a235f430db87a404491cd2ffc@thread.v2" } },
      "source": { "name": "message" },
      "tenant": { "id": "051cece0-e4dc-4aed-b471-bf29824e1ee6" }
    },
    "channelId": "msteams",
    "conversation": {
      "conversationType": "groupChat",
      "id": "19:f26a8d8a235f430db87a404491cd2ffc@thread.v2",
      "isGroup": true,
      "name": "與 [Even] 的客服對話 - Support Request for User",
      "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
    },
    "entities": [ { "locale": "zh-TW", "type": "clientInfo" } ],
    "from": {
      "aadObjectId": "bd5ab633-3f6f-48ad-be04-e144480da48f",
      "id": "29:1FrpqFNZ44LuLurxUt5tpSed4vXWXTdi428VAbwv9FSK-cI4UWRtIKdwEwFKCehU48w1hrQvIc5AXAoUtNKKBYg"
    },
    "id": "f:0c2b988b-d01a-7ad8-bdb3-896c91a06b0e",
    "locale": "zh-TW",
    "recipient": { "id": "28:844146d7-4ac9-4e4d-a463-d6e027714e81", "name": "lab-test-teams-notify" },
    "serviceUrl": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/",
    "timestamp": "2025-09-23T03:32:01.535Z",
    "type": "installationUpdate"
  }'
```

### Channel 測試
```bash
curl -X POST "http://localhost:8080/api/v1/messages" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "add",
    "channelData": {
      "channel": { "id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2" },
      "settings": { "selectedChannel": { "id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2" } },
      "source": { "name": "message" },
      "team": {
        "aadGroupId": "2747cdba-5b58-4fdf-9cac-dc52e725e728",
        "id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
        "name": "Teams Send Msg Test Group"
      },
      "tenant": { "id": "051cece0-e4dc-4aed-b471-bf29824e1ee6" }
    },
    "channelId": "msteams",
    "conversation": {
      "conversationType": "channel",
      "id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
      "isGroup": true,
      "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
    },
    "entities": [ { "locale": "zh-TW", "type": "clientInfo" } ],
    "from": {
      "aadObjectId": "bd5ab633-3f6f-48ad-be04-e144480da48f",
      "id": "29:1FrpqFNZ44LuLurxUt5tpSed4vXWXTdi428VAbwv9FSK-cI4UWRtIKdwEwFKCehU48w1hrQvIc5AXAoUtNKKBYg"
    },
    "id": "f:f0f9e059-9224-1cca-ae84-0e10d3d003dd",
    "locale": "zh-TW",
    "recipient": { "id": "28:844146d7-4ac9-4e4d-a463-d6e027714e81", "name": "lab-test-teams-notify" },
    "serviceUrl": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/",
    "timestamp": "2025-09-23T03:31:05.255Z",
    "type": "installationUpdate"
  }'
```