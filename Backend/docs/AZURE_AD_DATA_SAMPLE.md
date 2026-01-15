# Azure AD 資料庫結構與 Sample 資料

## 資料表結構

### 1. `azure_ad_groups` - Azure AD Groups

**用途**: 儲存從 Microsoft Graph API 同步的 Azure AD Groups

**表結構**:
```sql
CREATE TABLE azure_ad_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    azure_ad_id VARCHAR(255) NOT NULL UNIQUE,  -- Azure AD Group ID
    display_name VARCHAR(255) NOT NULL,        -- Group 顯示名稱
    description TEXT,                           -- Group 描述
    group_types JSONB DEFAULT '[]'::jsonb,      -- Group 類型陣列
    synced_at TIMESTAMP WITH TIME ZONE,         -- 最後同步時間
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**Sample 資料**:
```json
{
  "id": "68ca278a-10ef-4af7-ae48-0363ec68ec7d",
  "azure_ad_id": "8303127b-ae6e-4797-bcde-7bbb280609b9",
  "display_name": "CFH Teams Bot",
  "description": "CFH Teams Bot",
  "group_types": ["Unified"],
  "synced_at": "2026-01-14T07:17:56.813589Z",
  "created_at": "2026-01-13T10:15:26.211503Z",
  "updated_at": "2026-01-14T07:17:56.813589Z"
}
```

**重要欄位說明**:
- `group_types`: JSONB 陣列，可能的值：
  - `["Unified"]`: M365/Teams Group（**會出現在 Teams**）
  - `["DynamicMembership"]`: 動態成員群組
  - `["Unified", "DynamicMembership"]`: 同時是 M365 Group 和動態成員群組
  - `[]`: 一般 Security Group（**不會出現在 Teams**）

### 2. `azure_ad_channels` - Teams Channels

**用途**: 儲存從 Teams Unified Groups 同步的 Channels

**表結構**:
```sql
CREATE TABLE azure_ad_channels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    azure_ad_id VARCHAR(500) NOT NULL UNIQUE,   -- Channel ID (也是 conversation_id)
    team_id VARCHAR(255) NOT NULL,               -- 父 Team/Group ID (對應 azure_ad_groups.azure_ad_id)
    display_name VARCHAR(255) NOT NULL,          -- Channel 顯示名稱
    description TEXT,                            -- Channel 描述
    membership_type VARCHAR(50),                 -- standard, private, shared
    conversation_id VARCHAR(500) NOT NULL,       -- 與 azure_ad_id 相同（用於 Teams 通知）
    synced_at TIMESTAMP WITH TIME ZONE,          -- 最後同步時間
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**Sample 資料**:
```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "azure_ad_id": "19:abc123def456@thread.tacv2",
  "team_id": "8303127b-ae6e-4797-bcde-7bbb280609b9",
  "display_name": "General",
  "description": "General discussion channel",
  "membership_type": "standard",
  "conversation_id": "19:abc123def456@thread.tacv2",
  "synced_at": "2026-01-14T07:17:56.813589Z",
  "created_at": "2026-01-14T07:17:56.813589Z",
  "updated_at": "2026-01-14T07:17:56.813589Z"
}
```

**重要欄位說明**:
- `team_id`: 對應到 `azure_ad_groups.azure_ad_id`（必須是 `Unified` 類型的 Group）
- `azure_ad_id` 和 `conversation_id`: 相同值，用於 Teams Bot 發送通知
- `membership_type`: 
  - `standard`: 公開頻道
  - `private`: 私人頻道
  - `shared`: 共享頻道

## 資料關聯

```
azure_ad_groups (Unified 類型)
    ↓ (team_id 對應 azure_ad_id)
azure_ad_channels
```

**只有 `group_types` 包含 `"Unified"` 的 Group 才會有 Channels**

## 查詢範例

### 1. 只列出 Teams Groups（會出現在 Teams 的 Groups）

```sql
SELECT * FROM azure_ad_groups 
WHERE group_types @> '["Unified"]'::jsonb 
ORDER BY display_name;
```

### 2. 查詢特定 Team 的所有 Channels

```sql
SELECT * FROM azure_ad_channels 
WHERE team_id = '8303127b-ae6e-4797-bcde-7bbb280609b9'
ORDER BY display_name;
```

### 3. 查詢所有 Teams Groups 及其 Channels

```sql
SELECT 
    g.azure_ad_id as group_id,
    g.display_name as group_name,
    c.azure_ad_id as channel_id,
    c.display_name as channel_name,
    c.membership_type
FROM azure_ad_groups g
LEFT JOIN azure_ad_channels c ON g.azure_ad_id = c.team_id
WHERE g.group_types @> '["Unified"]'::jsonb
ORDER BY g.display_name, c.display_name;
```

## API 端點

### 獲取所有 Groups（包含非 Teams Groups）
```
GET /internal/v1/directory/groups
```

### 獲取只會出現在 Teams 的 Groups
```
GET /internal/v1/directory/groups/teams
```

### 獲取所有 Channels
```
GET /internal/v1/directory/channels
```

### 獲取特定 Team 的 Channels
```
GET /internal/v1/directory/channels?team_id={team_id}
```
