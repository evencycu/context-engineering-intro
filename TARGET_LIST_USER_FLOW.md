# Target List 使用者流程規劃

## 📋 概述

本文檔規劃使用者在專案中設定和使用 Target List 的完整流程，目標是管理微軟 Teams 的人員、聊天室、頻道等資源。

---

## 🎯 使用場景

**場景**: 使用者創建了一個新專案，需要設定目標受眾列表來發送 Teams 通知。

**目標資源類型**:
- 👤 **Teams 人員** (Users) - 個人用戶
- 💬 **聊天室** (Chat Groups) - 群組聊天
- 📢 **頻道** (Channels) - Teams 頻道

---

## 🔄 完整使用流程

### 階段 1: 專案創建與初始化

#### 1.1 創建新專案
```
使用者操作:
1. 進入 Admin 頁面
2. 點擊 "Add Project"
3. 填寫專案資訊:
   - Project Name (notify_key)
   - Department/Description
   - Daily/Monthly Limits
   - Priority
4. 點擊 "Create"
```

**系統行為**:
- 創建專案記錄
- 生成 API Key
- 初始化專案設定

**數據庫狀態**:
```sql
-- projects 表新增一筆記錄
INSERT INTO projects (id, company_id, notify_key, description, ...)
VALUES ('new-project-uuid', 'company-uuid', 'project-name', '...', ...);
```

---

### 階段 2: 同步 Teams 目錄資源

#### 2.1 首次同步 Azure AD / Teams 目錄

**使用者操作**:
```
1. 進入 Audience 頁面
2. 系統提示: "需要先同步 Teams 目錄才能建立列表"
3. 點擊 "Sync Directory" 按鈕
```

**系統行為**:
```
1. 調用 POST /internal/v1/directory/sync
2. 從 Microsoft Graph API 同步:
   - Azure AD Users (人員)
   - Azure AD Groups (Teams 群組)
   - Teams Channels (頻道)
3. 儲存到數據庫:
   - azure_ad_users 表
   - azure_ad_groups 表
   - (頻道資訊從 Graph API 即時查詢)
```

**數據庫狀態**:
```sql
-- azure_ad_users 表同步用戶
INSERT INTO azure_ad_users (azure_ad_id, display_name, email, ...)
VALUES ('user-guid', 'John Doe', 'john@example.com', ...);

-- azure_ad_groups 表同步群組
INSERT INTO azure_ad_groups (azure_ad_id, display_name, mail_nickname, ...)
VALUES ('group-guid', 'Marketing Team', 'marketing', ...);
```

**同步頻率**:
- 手動觸發: 使用者點擊 "Sync Directory"
- 自動同步: 可設定每日/每週自動同步（未來功能）

---

### 階段 3: 建立 Target List

#### 3.1 建立方式選擇

使用者有三種方式建立 Target List:

**方式 A: 從 Teams 資源建立 (推薦)**
- 從已同步的 Teams 用戶、群組、頻道選擇
- 系統自動建立列表並計算數量

**方式 B: 手動上傳 CSV**
- 上傳包含 Teams 用戶 email 或 conversation ID 的 CSV
- 系統解析並建立列表

**方式 C: 動態查詢建立**
- 使用標籤或查詢條件建立動態列表
- 系統定期更新列表成員

---

#### 3.2 方式 A: 從 Teams 資源建立列表

**使用者操作流程**:

```
步驟 1: 選擇資源類型
┌─────────────────────────────────────┐
│ Create Target List                  │
├─────────────────────────────────────┤
│ List Name: [________________]       │
│                                      │
│ Select Resources:                   │
│ ○ Users (人員)                       │
│ ○ Chat Groups (聊天室)               │
│ ○ Channels (頻道)                    │
│ ○ Mixed (混合)                       │
└─────────────────────────────────────┘

步驟 2: 選擇 Users
┌─────────────────────────────────────┐
│ Select Users                         │
├─────────────────────────────────────┤
│ Search: [________] 🔍                │
│                                      │
│ ☑ John Doe (john@example.com)       │
│ ☑ Jane Smith (jane@example.com)     │
│ ☐ Bob Johnson (bob@example.com)     │
│ ☐ Alice Brown (alice@example.com)   │
│                                      │
│ Selected: 2 users                   │
│ [Cancel] [Next]                      │
└─────────────────────────────────────┘

步驟 3: 選擇 Chat Groups (可選)
┌─────────────────────────────────────┐
│ Select Chat Groups                   │
├─────────────────────────────────────┤
│ ☑ Marketing Team Chat              │
│ ☐ Sales Team Chat                   │
│ ☐ Development Team Chat             │
│                                      │
│ Selected: 1 chat group              │
│ [Back] [Next]                        │
└─────────────────────────────────────┘

步驟 4: 選擇 Channels (可選)
┌─────────────────────────────────────┐
│ Select Channels                      │
├─────────────────────────────────────┤
│ Group: Marketing Team ▼              │
│   ☑ General                          │
│   ☑ Announcements                    │
│                                      │
│ Group: Sales Team ▼                  │
│   ☐ General                          │
│                                      │
│ Selected: 2 channels                 │
│ [Back] [Create List]                │
└─────────────────────────────────────┘
```

**系統行為**:

```typescript
// 前端調用
POST /internal/v1/projects/{projectId}/audience-lists
{
  "name": "Marketing Team List",
  "type": "Static",
  "description": "Marketing team members and channels",
  "metadata": {
    "source": "teams_resources",
    "resources": {
      "users": [
        {
          "type": "user",
          "azure_ad_id": "user-guid-1",
          "email": "john@example.com"
        },
        {
          "type": "user",
          "azure_ad_id": "user-guid-2",
          "email": "jane@example.com"
        }
      ],
      "chat_groups": [
        {
          "type": "chat_group",
          "chat_id": "19:chat-id-1",
          "name": "Marketing Team Chat"
        }
      ],
      "channels": [
        {
          "type": "channel",
          "team_id": "team-guid-1",
          "channel_id": "channel-guid-1",
          "name": "General"
        },
        {
          "type": "channel",
          "team_id": "team-guid-1",
          "channel_id": "channel-guid-2",
          "name": "Announcements"
        }
      ]
    },
    "total_count": 5  // 2 users + 1 chat + 2 channels
  }
}
```

**數據庫狀態**:
```sql
-- audience_lists 表新增列表
INSERT INTO audience_lists (
    id, project_id, name, type, count, description, metadata
) VALUES (
    'list-uuid',
    'project-uuid',
    'Marketing Team List',
    'Static',
    5,  -- 2 users + 1 chat + 2 channels
    'Marketing team members and channels',
    '{
        "source": "teams_resources",
        "resources": {
            "users": [...],
            "chat_groups": [...],
            "channels": [...]
        },
        "total_count": 5
    }'::jsonb
);
```

---

#### 3.3 方式 B: 手動上傳 CSV

**使用者操作**:
```
1. 點擊 "Upload New List"
2. 選擇 CSV 檔案
3. 系統解析檔案內容
4. 預覽解析結果
5. 確認建立列表
```

**CSV 格式範例**:
```csv
type,identifier,name
user,john@example.com,John Doe
user,jane@example.com,Jane Smith
chat_group,19:chat-id-1,Marketing Chat
channel,team-guid-1|channel-guid-1,General Channel
```

**系統行為**:
```typescript
// 前端上傳 CSV
POST /internal/v1/projects/{projectId}/audience-lists/upload
Content-Type: multipart/form-data

// 後端解析 CSV
// 驗證每個資源是否存在
// 計算總數
// 建立列表
```

---

#### 3.4 方式 C: 動態查詢建立

**使用者操作**:
```
1. 選擇 "Create Dynamic List"
2. 選擇查詢條件:
   - By Tags (標籤)
   - By Department (部門)
   - By Active Status (活躍狀態)
3. 設定更新頻率
4. 建立列表
```

**系統行為**:
```typescript
POST /internal/v1/projects/{projectId}/audience-lists
{
  "name": "Active Developers",
  "type": "Dynamic",
  "description": "Developers active in last 30 days",
  "metadata": {
    "query_type": "tag_based",
    "tags": ["developer", "backend"],
    "refresh_interval": "daily",
    "auto_update": true
  }
}
```

---

### 階段 4: 管理 Target List

#### 4.1 查看列表

**使用者操作**:
```
進入 Audience 頁面 → Target Lists tab
```

**顯示內容**:
```
┌─────────────────────────────────────────────────────────┐
│ Target Lists                                             │
├─────────────────────────────────────────────────────────┤
│ List Name        │ Recipients │ Last Updated │ Type    │
├──────────────────┼────────────┼──────────────┼─────────┤
│ Marketing Team   │ 5          │ 2024-01-15   │ Static  │
│ Active Users     │ 0          │ 2024-01-15   │ Dynamic │
│ Sales Team       │ 3          │ 2024-01-14   │ Static  │
└──────────────────┴────────────┴──────────────┴─────────┘
```

#### 4.2 編輯列表

**使用者操作**:
```
1. 點擊列表名稱或 "Edit" 按鈕
2. 修改列表內容:
   - 新增/移除資源
   - 修改名稱/描述
3. 儲存變更
```

**系統行為**:
```typescript
PUT /internal/v1/projects/{projectId}/audience-lists/{listId}
{
  "name": "Updated Marketing Team",
  "count": 6,
  "metadata": {
    "resources": {
      "users": [...],  // 新增或移除用戶
      "channels": [...] // 新增或移除頻道
    }
  }
}
```

#### 4.3 刪除列表

**使用者操作**:
```
1. 點擊列表的 "Delete" 按鈕
2. 確認刪除
```

**系統行為**:
```sql
DELETE FROM audience_lists WHERE id = 'list-uuid';
```

---

### 階段 5: 使用 Target List 發送通知

#### 5.1 在 Compose 頁面選擇列表

**使用者操作**:
```
1. 進入 Compose 頁面
2. 選擇 Recipient Type: "List"
3. 選擇 Target List: "Marketing Team List"
4. 選擇 Template 或輸入自訂內容
5. 發送通知
```

**系統行為**:
```typescript
POST /internal/v1/notifications
{
  "project_id": "project-uuid",
  "recipient_type": "List",
  "recipient": "Marketing Team List",
  "list_id": "list-uuid",
  "template_id": "template-uuid",
  "content": {...},
  "priority": "Normal"
}
```

**後端處理**:
```
1. 查詢 audience_lists 表獲取列表 metadata
2. 解析 metadata.resources 獲取所有目標:
   - Users: 轉換為 TeamsTarget (type: "person", email: "...")
   - Chat Groups: 轉換為 TeamsTarget (type: "chatgroup", conversation_id: "...")
   - Channels: 轉換為 TeamsTarget (type: "channel", conversation_id: "...")
3. 為每個目標建立 notification_destinations
4. 發送通知到每個目標
```

---

## 📊 數據模型關係

```
projects (專案)
  └── audience_lists (目標列表)
        └── metadata.resources (JSONB)
              ├── users[] (引用 azure_ad_users)
              ├── chat_groups[] (引用 chat_groups)
              └── channels[] (引用 azure_ad_groups + Graph API)

azure_ad_users (Teams 用戶)
  └── (同步自 Microsoft Graph API)

azure_ad_groups (Teams 群組)
  └── (同步自 Microsoft Graph API)
        └── channels (從 Graph API 即時查詢)

chat_groups (專案聊天群組)
  └── (手動註冊或從列表建立)
```

---

## 🔧 技術實現細節

### 1. Teams 資源類型映射

| Teams 資源 | 數據庫表 | Graph API 端點 | Target Type |
|-----------|---------|---------------|-------------|
| User | azure_ad_users | /users | person |
| Chat Group | chat_groups | N/A (手動註冊) | chatgroup |
| Channel | N/A (即時查詢) | /teams/{id}/channels | channel |

### 2. Metadata 結構設計

```json
{
  "source": "teams_resources" | "csv_upload" | "dynamic_query",
  "resources": {
    "users": [
      {
        "type": "user",
        "azure_ad_id": "guid",
        "email": "user@example.com",
        "display_name": "John Doe"
      }
    ],
    "chat_groups": [
      {
        "type": "chat_group",
        "chat_id": "19:chat-id",
        "name": "Team Chat"
      }
    ],
    "channels": [
      {
        "type": "channel",
        "team_id": "team-guid",
        "channel_id": "channel-guid",
        "name": "General",
        "team_name": "Marketing Team"
      }
    ]
  },
  "total_count": 5,
  "last_sync": "2024-01-15T10:30:00Z"
}
```

### 3. 列表建立 API 設計

```typescript
// 從 Teams 資源建立
POST /internal/v1/projects/{projectId}/audience-lists
{
  "name": "Marketing Team",
  "type": "Static",
  "description": "Marketing team members",
  "resources": {
    "user_ids": ["user-guid-1", "user-guid-2"],
    "chat_group_ids": ["chat-group-uuid-1"],
    "channel_ids": [
      {"team_id": "team-guid-1", "channel_id": "channel-guid-1"}
    ]
  }
}

// 上傳 CSV
POST /internal/v1/projects/{projectId}/audience-lists/upload
Content-Type: multipart/form-data
file: <CSV file>

// 動態查詢
POST /internal/v1/projects/{projectId}/audience-lists
{
  "name": "Active Developers",
  "type": "Dynamic",
  "query": {
    "type": "tag_based",
    "tags": ["developer"],
    "refresh_interval": "daily"
  }
}
```

---

## 🎨 UI/UX 設計建議

### 1. 首次使用引導

```
┌─────────────────────────────────────────┐
│ Welcome to Target Lists!                 │
├─────────────────────────────────────────┤
│ Before creating lists, you need to:     │
│                                          │
│ 1. Sync Teams Directory                  │
│    [Sync Directory]                      │
│                                          │
│ 2. Create your first list                │
│    [Create List]                         │
└─────────────────────────────────────────┘
```

### 2. 列表建立精靈 (Wizard)

```
Step 1: Basic Info
  └─ List Name, Description

Step 2: Select Resources
  └─ Users / Chat Groups / Channels

Step 3: Review & Create
  └─ Preview, Confirm
```

### 3. 列表詳情頁面

```
┌─────────────────────────────────────────┐
│ Marketing Team List                      │
├─────────────────────────────────────────┤
│ Type: Static                            │
│ Count: 5                                 │
│ Last Updated: 2024-01-15 10:30          │
│                                          │
│ Resources:                               │
│ • Users (2)                              │
│   - John Doe (john@example.com)         │
│   - Jane Smith (jane@example.com)       │
│                                          │
│ • Chat Groups (1)                        │
│   - Marketing Team Chat                  │
│                                          │
│ • Channels (2)                           │
│   - Marketing Team > General             │
│   - Marketing Team > Announcements      │
│                                          │
│ [Edit] [Delete] [Use in Compose]         │
└─────────────────────────────────────────┘
```

---

## 📝 待實現功能清單

### Phase 1: 基礎功能 (P0)
- [x] Audience Lists CRUD API
- [ ] 從 Teams 資源建立列表 UI
- [ ] CSV 上傳功能
- [ ] 列表詳情頁面

### Phase 2: 進階功能 (P1)
- [ ] 動態列表查詢
- [ ] 列表預覽功能
- [ ] 批量操作（批量新增/移除資源）
- [ ] 列表複製功能

### Phase 3: 優化功能 (P2)
- [ ] 列表模板
- [ ] 自動同步 Teams 資源
- [ ] 列表使用統計
- [ ] 列表分享功能

---

## 🔍 技術考量

### 1. 性能優化
- **大量資源**: 如果列表包含數百個資源，考慮分頁顯示
- **即時查詢**: Channels 從 Graph API 即時查詢，考慮快取
- **批量操作**: 批量新增/移除資源時，使用事務處理

### 2. 數據一致性
- **資源刪除**: 如果 Teams 資源被刪除，列表中的引用需要處理
- **同步狀態**: 顯示資源最後同步時間，提示使用者是否需要重新同步

### 3. 錯誤處理
- **Graph API 限流**: 處理 Microsoft Graph API 限流錯誤
- **無效資源**: 驗證資源是否存在，提示使用者無效資源
- **權限問題**: 檢查使用者是否有權限訪問特定 Teams 資源

---

**最後更新**: 2026-01-13  
**狀態**: 規劃階段，待實現
