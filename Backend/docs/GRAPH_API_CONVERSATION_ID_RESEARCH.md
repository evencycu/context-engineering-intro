# Microsoft Graph API - 取得所有 Conversation ID 研究報告

## 研究目的
研究是否有 Graph API 可以取得所有的 Teams conversation IDs（包含 personal chats, group chats, channels）。

## Conversation ID 類型與取得方式

### 1. ✅ Channel Conversation ID（已實作）

**格式**: `19:xxx@thread.tacv2`

**API**: 
```http
GET https://graph.microsoft.com/v1.0/teams/{team-id}/channels
```

**權限**: `Channel.ReadBasic.All` (application permission)

**狀態**: ✅ **已實作**
- Channel 的 `id` 欄位就是 `conversation_id`
- 已在 `graph_service.go` 中實作 `GetGroupChannels()`
- 已在 Directory Sync 中同步 channels

**限制**: 
- 只有 Unified (M365) Groups 才能取得 channels
- 不是所有 Unified Groups 都是 Teams teams

---

### 2. ⚠️ Personal Chat Conversation ID（部分可用）

**格式**: 
- `19:{user-id}_{bot-id}@unq.gbl.spaces` (Teams format)
- `a:xxx` (Bot Framework format)

#### 方法 A: `/me/chats` (Delegated Permission) ❌ **不適用**

**API**: 
```http
GET https://graph.microsoft.com/v1.0/me/chats
GET https://graph.microsoft.com/v1.0/me/chats?$filter=chatType eq 'oneOnOne'
```

**權限**: `Chat.Read` (delegated permission)

**限制**: 
- ❌ 需要**使用者授權**（delegated permissions）
- ❌ 不適用於應用程式權限（app-only）
- ❌ 只能取得**當前使用者**的 chats
- ❌ 無法取得所有使用者的 chats

#### 方法 B: `/chats` (Application Permission) ⭐ **需要驗證**

**API**: 
```http
GET https://graph.microsoft.com/v1.0/chats
GET https://graph.microsoft.com/v1.0/chats?$filter=chatType eq 'oneOnOne'
```

**權限**: `Chat.Read.All` 或 `Chat.ReadWrite.All` (application permission)

**回應範例**:
```json
{
  "value": [
    {
      "id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces",
      "chatType": "oneOnOne",
      "topic": null,
      "createdDateTime": "2025-09-18T06:42:39.469Z",
      "webUrl": "https://teams.microsoft.com/l/chat/..."
    }
  ]
}
```

**重要限制**（需要驗證）:
- ⚠️ 需要**租戶管理員同意**這些權限
- ⚠️ **可能只能取得 Bot 有參與的 chats**（需要驗證）
- ⚠️ **可能無法取得所有使用者的 personal chats**（需要驗證）
- ⚠️ 需要確認實際權限範圍

#### 方法 C: 從 `bot_installations` 表（目前使用）✅

**來源**: Bot Framework installation update events

**優點**: 
- ✅ 已經有資料（從 Bot 安裝事件取得）
- ✅ 不需要額外 API 調用
- ✅ 不需要額外權限

**限制**:
- ⚠️ 只有 Bot 已安裝的使用者才有 conversation_id
- ⚠️ 需要 Bot 先被使用者安裝

---

### 3. ⚠️ Group Chat Conversation ID（需要研究）

**格式**: `19:xxx@thread.v2`

#### 方法 A: `/me/chats` (Delegated Permission) ❌ **不適用**

**API**: 
```http
GET https://graph.microsoft.com/v1.0/me/chats?$filter=chatType eq 'groupChat'
```

**權限**: `Chat.Read` (delegated permission)

**限制**: 
- ❌ 需要使用者授權
- ❌ 只能取得當前使用者的 group chats

#### 方法 B: `/chats` (Application Permission) ⭐ **需要驗證**

**API**: 
```http
GET https://graph.microsoft.com/v1.0/chats?$filter=chatType eq 'groupChat'
```

**權限**: `Chat.Read.All` 或 `Chat.ReadWrite.All` (application permission)

**回應範例**:
```json
{
  "value": [
    {
      "id": "19:f26a8d8a235f430db87a404491cd2ffc@thread.v2",
      "chatType": "groupChat",
      "topic": "Team Discussion",
      "createdDateTime": "2025-09-18T07:00:00.000Z"
    }
  ]
}
```

**重要限制**（需要驗證）:
- ⚠️ 需要租戶管理員同意
- ⚠️ **可能只能取得 Bot 有參與的 group chats**（需要驗證）

---

## Graph API Endpoints 總結

### ✅ 已實作且可用（Application Permission）

| Endpoint | 用途 | 權限 | 狀態 | Conversation ID 來源 |
|----------|------|------|------|---------------------|
| `GET /users` | 取得所有使用者 | `User.Read.All` | ✅ 已實作 | N/A |
| `GET /groups` | 取得所有群組 | `Group.Read.All` | ✅ 已實作 | N/A |
| `GET /teams/{team-id}/channels` | 取得 Team 的 Channels | `Channel.ReadBasic.All` | ✅ 已實作 | Channel `id` 欄位 |

### ⚠️ 需要研究的新 Endpoint

| Endpoint | 用途 | 權限 | 狀態 | 備註 |
|----------|------|------|------|------|
| `GET /chats` | 取得所有 Chats | `Chat.Read.All` | ⚠️ **需要驗證** | 可能只能取得 Bot 參與的 chats |
| `GET /chats?$filter=chatType eq 'oneOnOne'` | 取得 Personal Chats | `Chat.Read.All` | ⚠️ **需要驗證** | 需要確認權限範圍 |
| `GET /chats?$filter=chatType eq 'groupChat'` | 取得 Group Chats | `Chat.Read.All` | ⚠️ **需要驗證** | 需要確認權限範圍 |

---

## 目前系統狀態

### ✅ 已實作
1. **Channels**: 
   - 透過 `/teams/{team-id}/channels` 取得
   - Channel `id` 就是 `conversation_id`
   - 已同步到 `azure_ad_channels` 表

2. **Personal Chats**: 
   - 從 `bot_installations` 表取得（Bot 安裝事件）
   - 已透過 JOIN 查詢顯示在 Users tab

### ❌ 未實作
1. **透過 Graph API 取得所有 personal chats**
2. **透過 Graph API 取得所有 group chats**

---

## 關鍵問題與驗證需求

### 問題 1: `GET /chats` 是否真的能取得所有 chats？

**需要驗證**:
- [ ] 使用 application token 測試 `GET /chats` API
- [ ] 確認是否能取得 Bot **未參與**的 chats
- [ ] 確認是否能取得**所有使用者**的 personal chats
- [ ] 確認是否能取得**所有** group chats

**測試命令**:
```bash
# 測試取得所有 chats
curl -X GET "https://graph.microsoft.com/v1.0/chats" \
  -H "Authorization: Bearer {application_token}"

# 測試取得 personal chats
curl -X GET "https://graph.microsoft.com/v1.0/chats?\$filter=chatType eq 'oneOnOne'" \
  -H "Authorization: Bearer {application_token}"

# 測試取得 group chats
curl -X GET "https://graph.microsoft.com/v1.0/chats?\$filter=chatType eq 'groupChat'" \
  -H "Authorization: Bearer {application_token}"
```

### 問題 2: `Chat.Read.All` 權限的實際範圍

**需要確認**:
- [ ] 權限是否已獲得租戶管理員同意
- [ ] 實際能取得的資料範圍
- [ ] 是否有限制（例如：只能取得 Bot 參與的 chats）

### 問題 3: 效能與分頁

**需要考慮**:
- [ ] 如果組織有大量 chats，API 調用是否會很慢？
- [ ] 是否需要處理分頁（`@odata.nextLink`）？
- [ ] 是否需要增量同步（delta query）？

---

## 建議方案

### 方案 1: 測試並實作 `/chats` API ⭐ **推薦**

**步驟**:
1. **申請權限**: 在 Azure AD App Registration 中申請 `Chat.Read.All` 權限
2. **管理員同意**: 請租戶管理員同意權限
3. **測試 API**: 使用 application token 測試 `GET /chats` API
4. **驗證範圍**: 確認實際能取得的資料範圍
5. **實作**: 如果測試成功，實作新的 Graph Service 方法

**優點**:
- ✅ 可以取得所有 Bot 有參與的 chats
- ✅ 不需要等待 Bot 安裝事件
- ✅ 可以取得 group chats

**缺點**:
- ⚠️ 需要額外權限（需要管理員同意）
- ⚠️ 可能只能取得 Bot 有參與的 chats（需要驗證）
- ⚠️ 可能無法取得所有使用者的 personal chats（需要驗證）

### 方案 2: 繼續使用 `bot_installations` 表（目前方案）

**優點**:
- ✅ 已經有資料
- ✅ 不需要額外權限
- ✅ 不需要額外 API 調用

**缺點**:
- ⚠️ 只有 Bot 已安裝的使用者才有 conversation_id
- ⚠️ 無法取得 group chats（除非 Bot 有參與）

### 方案 3: 混合方案

**實作**:
- **Personal Chats**: 繼續使用 `bot_installations` 表（已有資料）
- **Group Chats**: 使用 `/chats?$filter=chatType eq 'groupChat'` API（如果測試成功）
- **Channels**: 繼續使用 `/teams/{team-id}/channels` API（已實作）

---

## 測試計劃

### 測試 1: 驗證 `/chats` API 可用性

```bash
# 1. 確認權限已申請並獲得同意
# 2. 取得 application token
# 3. 測試 API

# 測試取得所有 chats
curl -X GET "https://graph.microsoft.com/v1.0/chats" \
  -H "Authorization: Bearer {application_token}" \
  -H "Content-Type: application/json"

# 測試過濾 personal chats
curl -X GET "https://graph.microsoft.com/v1.0/chats?\$filter=chatType eq 'oneOnOne'" \
  -H "Authorization: Bearer {application_token}"

# 測試過濾 group chats
curl -X GET "https://graph.microsoft.com/v1.0/chats?\$filter=chatType eq 'groupChat'" \
  -H "Authorization: Bearer {application_token}"
```

### 測試 2: 驗證權限範圍

**測試項目**:
- [ ] 是否能取得 Bot 未參與的 personal chats？
- [ ] 是否能取得所有使用者的 personal chats？
- [ ] 是否能取得 Bot 未參與的 group chats？
- [ ] 回應中是否包含所有必要的欄位（id, chatType, topic 等）？

### 測試 3: 效能測試

**測試項目**:
- [ ] API 回應時間
- [ ] 分頁處理（`@odata.nextLink`）
- [ ] 大量資料的處理能力

---

## 參考資料

### Microsoft 官方文檔
- [List chats - Microsoft Graph API](https://learn.microsoft.com/en-us/graph/api/chat-list)
- [Chat permissions - Microsoft Graph API](https://learn.microsoft.com/en-us/graph/permissions-reference#chat-permissions)
- [Get channel - Microsoft Graph API](https://learn.microsoft.com/en-us/graph/api/channel-get)

### 社群討論
- [How to fetch chat list and their members with MS Graph API](https://techcommunity.microsoft.com/discussions/teamsdeveloper/how-to-fetch-chat-list-and-their-members-with-ms-graph-api/2272191)
- [Microsoft Teams - List all chats regardless of users](https://stackoverflow.com/questions/62791536/microsoft-teams-list-all-chats-regardless-of-users)

---

## 測試結果（2026-01-15）

### 測試環境
- **Client ID**: `844146d7-4ac9-4e4d-a463-d6e027714e81`
- **Tenant ID**: `051cece0-e4dc-4aed-b471-bf29824e1ee6`
- **權限類型**: Application Permission (Client Credentials Flow)

### 測試結果

#### ✅ Access Token 取得成功
- Token 成功取得，可以使用 Graph API

#### ❌ `/chats` API 測試失敗

**錯誤訊息**:
```json
{
  "error": {
    "code": "BadRequest",
    "message": "Requested API is not supported in application-only context",
    "innerError": {
      "date": "2026-01-15T06:00:40",
      "request-id": "5b37a7a3-a712-428e-9573-8945deea955d",
      "client-request-id": "5b37a7a3-a712-428e-9573-8945deea955d"
    }
  }
}
```

**結論**:
- ❌ `GET /chats` API **不支援 application-only context**
- ❌ 即使有 `Chat.Read.All` application permission，也無法使用此 API
- ❌ 此 API 只支援 **delegated permissions**（需要使用者授權）

### 測試的 API Endpoints

1. **GET /chats** - ❌ 失敗（不支援 application-only context）
2. **GET /chats?$filter=chatType eq 'oneOnOne'** - ❌ 失敗（不支援 application-only context）
3. **GET /chats?$filter=chatType eq 'groupChat'** - ❌ 失敗（不支援 application-only context）

---

## 結論

### ✅ 目前可用
1. **Channels**: 
   - ✅ 已實作，channel `id` 就是 `conversation_id`
   - ✅ 已同步到資料庫
   - ✅ 支援 application permissions

2. **Personal Chats**: 
   - ✅ 從 `bot_installations` 表取得（Bot 安裝事件）
   - ✅ 已透過 JOIN 查詢顯示在 Users tab
   - ✅ 不需要額外 API 調用

### ❌ 無法使用（已驗證）
1. **`GET /chats` API**: 
   - ❌ **不支援 application-only context**
   - ❌ 即使有 `Chat.Read.All` application permission 也無法使用
   - ❌ 只支援 delegated permissions（需要使用者授權）
   - ❌ **無法透過此 API 取得所有 chats**

2. **Group Chats**: 
   - ❌ 無法透過 `/chats` API 取得（不支援 application-only context）
   - ⚠️ 目前只能透過 `bot_installations` 表取得（Bot 有參與的 group chats）

### 📋 最終建議

#### 方案 1: 繼續使用現有方案（推薦）✅

**Personal Chats**:
- ✅ 繼續使用 `bot_installations` 表
- ✅ 從 Bot 安裝事件中取得 conversation_id
- ✅ 不需要額外 API 調用

**Group Chats**:
- ⚠️ 繼續使用 `chat_groups` 表（手動註冊）
- ⚠️ 或從 `bot_installations` 表取得（Bot 有參與的）

**Channels**:
- ✅ 繼續使用 `/teams/{team-id}/channels` API
- ✅ 已實作並同步

#### 方案 2: 使用 Delegated Permissions（不推薦）

如果要使用 `/chats` API，需要：
- ❌ 改用 delegated permissions（需要使用者授權）
- ❌ 無法在後端自動同步
- ❌ 需要使用者互動授權
- ❌ 不符合目前的架構設計（application-only）

---

## 重要結論

⚠️ **`GET /chats` API 無法在 application-only context 下使用**

即使：
- ✅ 有 `Chat.Read.All` application permission
- ✅ 有租戶管理員同意
- ✅ Access token 有效

**仍然無法使用**，因為此 API 只支援 delegated permissions。

**建議**：
- ✅ 繼續使用現有方案（`bot_installations` 表）
- ✅ 不需要實作 `/chats` API 調用
- ✅ 目前的方案已經足夠
