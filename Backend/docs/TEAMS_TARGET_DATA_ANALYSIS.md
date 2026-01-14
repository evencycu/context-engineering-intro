# Teams Notification 目標資料分析

## 執行摘要

目前資料庫中的 **Azure AD Users** 和 **Azure AD Groups** 資料**可以部分用於 Teams notify 目標**，但需要額外處理才能完整支援所有 Teams 目標類型。

## Teams Notification 目標類型

根據 `TeamsTarget` 模型定義，Teams 通知支援三種目標類型：

```go
type TeamsTarget struct {
    Type           string `json:"type"` // person, channel, chatgroup
    Email          string `json:"email,omitempty"`
    ConversationID string `json:"conversation_id,omitempty"`
    DisplayName    string `json:"display_name,omitempty"`
    TenantID       string `json:"tenant_id,omitempty"`
}
```

### 1. `personal` - 個人聊天
- **需求**: `conversation_id` 或 `email`
- **現有資料可用性**: ✅ **完全可用**
- **說明**: 
  - `azure_ad_users` 表已有 `email` 欄位
  - 236 筆使用者資料都包含 email
  - 可直接用於 `personal` 類型的通知目標

### 2. `channel` - Teams 頻道
- **需求**: `conversation_id`（頻道的 conversation ID）
- **現有資料可用性**: ⚠️ **部分可用，需要額外同步**
- **說明**:
  - 程式碼中已有 `GetGroupChannels()` 方法（`graph_service.go`）
  - 可以透過 `/teams/{team-id}/channels` API 取得頻道列表
  - **重要發現**: Graph API 回傳的 channel `id` **就是 `conversation_id`**！
    - 格式：`19:xxx@thread.tacv2`（標準 Teams conversation ID 格式）
    - 測試結果：`"id": "19:CpbIxO-wpsRWncD-kN2jMx598dZi-FOKjOCqtKMqxRI1@thread.tacv2"`
    - **可直接使用，無需額外轉換**
  - **限制**: 只有 **Unified (M365) Groups** 才能取得 Teams channels
  - 目前資料庫中有 **至少 10+ 個 Unified Groups**（例如：`LAB-SASE PA POC`, `LAB-DW戰情室` 等）
  - **需要**: 
    1. 識別哪些 AD Groups 是 Unified Groups
    2. 對每個 Unified Group 呼叫 `GetGroupChannels()` 取得 channels
    3. 將 channels 存入資料庫（目前沒有 `azure_ad_channels` 表）
    4. Channel 的 `id` 欄位可直接作為 `conversation_id` 使用

### 3. `groupChat` - Teams 群組聊天
- **需求**: `conversation_id`（chat 的 conversation ID）
- **現有資料可用性**: ❌ **無法直接取得**
- **說明**:
  - Graph API 有 `/me/chats` 端點可以取得使用者加入的 chats
  - 但這需要**使用者授權**（delegated permissions），不適用於應用程式權限（app-only）
  - 目前系統使用 **Client Credentials Flow**（應用程式權限），無法取得 chats
  - **替代方案**: 
    - 目前有 `ChatGroup` 模型和 `RegisterChatGroup()` 方法
    - 需要**手動註冊** chat IDs（可能是透過 Bot Framework 或其他方式取得）

## 目前資料庫狀態

### Azure AD Users
- **數量**: 236 筆
- **欄位**: `azure_ad_id`, `display_name`, `email`, `job_title`, `department`
- **可用於**: ✅ `personal` target（透過 email）

### Azure AD Groups
- **數量**: 154 筆
- **欄位**: `azure_ad_id`, `display_name`, `description`, `group_types`
- **Unified Groups**: 至少 10+ 個（`group_types` 包含 `"Unified"`）
- **可用於**: ⚠️ 需要先取得 channels 才能用於 `channel` target

### Azure AD Channels
- **狀態**: ❌ **尚未同步**
- **需要**: 新增 `azure_ad_channels` 表並實作同步邏輯

### Chat Groups
- **狀態**: ⚠️ **僅支援手動註冊**
- **表**: `chat_groups`（已存在於 schema）
- **欄位**: `project_id`, `name`, `chat_id`

## Graph API 支援情況

### ✅ 已實作
1. **取得 Users**: `/users` - ✅ 已同步
2. **取得 Groups**: `/groups` - ✅ 已同步
3. **取得 Channels**: `/teams/{team-id}/channels` - ✅ 已實作方法，但未同步

### ❌ 無法取得（使用應用程式權限）
1. **取得 Chats**: `/me/chats` - 需要使用者授權（delegated permissions）
2. **取得所有 Chats**: 沒有應用程式權限的端點可以取得所有 chats

### ⚠️ 需要研究
1. **取得 Channel Conversation ID**: 
   - Channels 的 `id` 不等於 `conversation_id`
   - 需要透過 Bot Framework 或其他方式取得 `conversation_id`
   - 或使用 `/teams/{team-id}/channels/{channel-id}` 的 `webUrl` 來推導

## 建議實作方案

### Phase 1: 支援 Channels（優先）
1. **新增 `azure_ad_channels` 表**
   ```sql
   CREATE TABLE azure_ad_channels (
       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
       azure_ad_id VARCHAR(255) NOT NULL UNIQUE,
       team_id VARCHAR(255) NOT NULL, -- 對應 azure_ad_groups.azure_ad_id
       display_name VARCHAR(255) NOT NULL,
       description TEXT,
       membership_type VARCHAR(50), -- standard, private, shared
       conversation_id VARCHAR(255), -- 需要額外取得
       synced_at TIMESTAMP WITH TIME ZONE,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
   );
   ```

2. **擴充 Directory Sync**
   - 在 `SyncDirectory()` 中，對每個 Unified Group 呼叫 `GetGroupChannels()`
   - 將 channels 存入 `azure_ad_channels` 表

3. **取得 Conversation ID** ✅ **已解決**
   - Graph API 回傳的 channel `id` **就是 `conversation_id`**
   - 格式：`19:xxx@thread.tacv2`（標準 Teams conversation ID）
   - **無需額外處理，可直接使用**

### Phase 2: 支援 Chat Groups（較低優先）
1. **保持現有手動註冊機制**
   - 使用者透過 UI 手動輸入 chat ID
   - 或透過 Bot Framework 取得後手動註冊

2. **研究替代方案**
   - 是否可以使用 Bot Framework 的 API 來取得 chats
   - 或提供工具讓使用者更容易取得 chat IDs

## 結論

### ✅ 目前可用
- **Personal targets**: 236 個 AD Users 的 email 可直接使用

### ⚠️ 需要額外實作
- **Channel targets**: 
  - 需要同步 Unified Groups 的 channels
  - ✅ Channel 的 `id` 就是 `conversation_id`，可直接使用
  - 預估可取得 **數十到數百個 channels**（取決於 Unified Groups 數量）

### ❌ 無法自動取得
- **Chat Group targets**: 
  - 需要手動註冊
  - 或需要改用使用者授權流程（不建議，會增加複雜度）

## 下一步行動

1. **立即**: 實作 Channels 同步功能
   - 新增 `azure_ad_channels` 表
   - 擴充 `SyncDirectory()` 同步 Unified Groups 的 channels
   - Channel `id` 可直接作為 `conversation_id` 使用
2. **短期**: 測試 Channels 同步功能並驗證可用性
3. **長期**: 評估是否需要支援 Chat Groups 的自動同步（可能需要改變授權流程）

## 測試結果

### Channel API 測試
```bash
# 測試 Unified Group: LAB-SASE PA POC (005d564d-f571-4faf-a969-3ed1fd60d389)
GET /internal/v1/directory/groups/{groupId}/channels

# 回應：
{
  "id": "19:CpbIxO-wpsRWncD-kN2jMx598dZi-FOKjOCqtKMqxRI1@thread.tacv2",
  "display_name": "General",
  "description": "LAB-SASE PA POC",
  "membership_type": "standard"
}
```

**結論**: Channel 的 `id` 欄位就是 `conversation_id`，格式為 `19:xxx@thread.tacv2`，可直接用於 Teams notification 的 `channel` target。
