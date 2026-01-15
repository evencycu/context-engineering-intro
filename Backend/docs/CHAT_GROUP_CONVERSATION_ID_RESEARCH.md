# Chat Group Conversation ID 研究報告

## 研究目的
研究是否有辦法查到所有的 chat group 的 conversation ID。

## 目前已知的限制

### ❌ Graph API `/chats` 不支援 Application-Only Context

**測試結果**（2026-01-15）:
- `GET /chats` API 返回錯誤：`"Requested API is not supported in application-only context"`
- 即使有 `Chat.Read.All` application permission 也無法使用
- 此 API 只支援 delegated permissions（需要使用者授權）

---

## 可能的方案

### 方案 1: 從 `bot_installations` 表查詢 ✅ **目前可用**

**說明**:
- Bot Framework 會在 Bot 被加入到 group chat 時觸發 `conversationUpdate` 事件
- 這些事件會被記錄到 `bot_installations` 表中
- `conversation_type = 'groupChat'` 的記錄就是 group chats

**資料庫查詢**:
```sql
SELECT 
    conversation_id,
    conversation_type,
    from_name,
    from_aad_object_id,
    email,
    installed_at,
    last_activity_at,
    installation_status
FROM bot_installations
WHERE conversation_type = 'groupChat'
  AND installation_status = 'active'
ORDER BY installed_at DESC;
```

**優點**:
- ✅ 已經有資料（從 Bot 安裝事件取得）
- ✅ 不需要額外 API 調用
- ✅ 不需要額外權限
- ✅ 可以取得 Bot 有參與的所有 group chats

**限制**:
- ⚠️ 只能取得 Bot **有參與**的 group chats
- ⚠️ 無法取得 Bot 未參與的 group chats
- ⚠️ 需要 Bot 先被加入到 group chat

**實作狀態**:
- ✅ 已實作 `ListAllChatGroups()` 方法（從 `chat_groups` 表）
- ⚠️ 需要新增從 `bot_installations` 表查詢 group chats 的方法

---

### 方案 2: 使用 `GET /users/{user-id}/chats` (Delegated Permission) ❌ **不適用**

**API**: 
```http
GET https://graph.microsoft.com/v1.0/users/{user-id}/chats?$filter=chatType eq 'groupChat'
```

**權限**: `Chat.Read` 或 `Chat.Read.All` (delegated permission)

**限制**: 
- ❌ 需要使用者授權（delegated permissions）
- ❌ 不適用於應用程式權限（app-only）
- ❌ 只能取得**特定使用者**參與的 group chats
- ❌ 無法取得所有使用者的 group chats
- ❌ 需要為每個使用者分別調用 API

**結論**: 不適用於目前的架構（application-only）

---

### 方案 3: 使用 `chat_groups` 表（手動註冊）✅ **目前使用**

**說明**:
- 目前系統有 `chat_groups` 表用於手動註冊 group chats
- 需要管理員手動註冊每個 group chat 的 conversation ID

**資料庫結構**:
```sql
CREATE TABLE chat_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    chat_id VARCHAR(500) NOT NULL,  -- This is the conversation_id
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, chat_id)
);
```

**優點**:
- ✅ 可以手動註冊任何 group chat
- ✅ 不需要 Bot 參與
- ✅ 可以按 project 組織

**限制**:
- ⚠️ 需要手動註冊
- ⚠️ 無法自動發現新的 group chats

**實作狀態**:
- ✅ 已實作 `RegisterChatGroup()` 方法
- ✅ 已實作 `ListChatGroups()` 方法
- ✅ 已實作 `ListAllChatGroups()` 方法

---

### 方案 4: 混合方案（推薦）⭐

**實作**:
1. **從 `bot_installations` 表取得 Bot 參與的 group chats**
   - 自動取得，不需要手動註冊
   - 可以定期同步

2. **從 `chat_groups` 表取得手動註冊的 group chats**
   - 可以註冊 Bot 未參與的 group chats
   - 可以按 project 組織

3. **合併兩個來源的資料**
   - 提供統一的 API 查詢所有 group chats
   - 去重複（相同的 conversation_id）

**API 設計**:
```go
// Get all group chats from both sources
func (s *directoryService) GetAllGroupChats(ctx context.Context) ([]GroupChatInfo, error) {
    // 1. Get from bot_installations
    botGroupChats, err := s.repo.GetGroupChatsFromBotInstallations(ctx)
    
    // 2. Get from chat_groups
    registeredGroupChats, err := s.repo.GetGroupChatsFromChatGroups(ctx)
    
    // 3. Merge and deduplicate
    return mergeGroupChats(botGroupChats, registeredGroupChats), nil
}
```

---

## 資料庫查詢範例

### 查詢 1: 從 `bot_installations` 表取得所有 group chats

```sql
SELECT 
    bi.conversation_id,
    bi.conversation_type,
    bi.from_name as chat_name,
    bi.from_aad_object_id,
    bi.email,
    bi.installed_at,
    bi.last_activity_at,
    bi.installation_status,
    bi.teams_tenant_id,
    bi.bot_id
FROM bot_installations bi
WHERE bi.conversation_type = 'groupChat'
  AND bi.installation_status = 'active'
ORDER BY bi.installed_at DESC;
```

### 查詢 2: 從 `chat_groups` 表取得所有註冊的 group chats

```sql
SELECT 
    cg.id,
    cg.project_id,
    cg.name,
    cg.chat_id as conversation_id,
    cg.created_at,
    cg.updated_at
FROM chat_groups cg
ORDER BY cg.name;
```

### 查詢 3: 合併兩個來源（去重複）

```sql
-- Get all unique group chat conversation IDs
SELECT DISTINCT conversation_id
FROM (
    -- From bot_installations
    SELECT conversation_id
    FROM bot_installations
    WHERE conversation_type = 'groupChat'
      AND installation_status = 'active'
    
    UNION
    
    -- From chat_groups
    SELECT chat_id as conversation_id
    FROM chat_groups
) AS all_group_chats
ORDER BY conversation_id;
```

---

## 建議實作方案

### 階段 1: 查詢 `bot_installations` 表的 group chats

**新增 Repository 方法**:
```go
func (r *directoryRepository) GetGroupChatsFromBotInstallations(ctx context.Context) ([]GroupChatInfo, error) {
    query := `
        SELECT 
            bi.conversation_id,
            bi.conversation_type,
            bi.from_name as chat_name,
            bi.from_aad_object_id,
            bi.email,
            bi.installed_at,
            bi.last_activity_at,
            bi.installation_status,
            bi.teams_tenant_id
        FROM bot_installations bi
        WHERE bi.conversation_type = 'groupChat'
          AND bi.installation_status = 'active'
        ORDER BY bi.installed_at DESC
    `
    var groupChats []GroupChatInfo
    err := r.db.SelectContext(ctx, &groupChats, query)
    return groupChats, err
}
```

### 階段 2: 合併兩個來源

**新增 Service 方法**:
```go
func (s *directoryService) GetAllGroupChats(ctx context.Context) ([]GroupChatInfo, error) {
    // Get from bot_installations
    botGroupChats, err := s.repo.GetGroupChatsFromBotInstallations(ctx)
    if err != nil {
        return nil, err
    }
    
    // Get from chat_groups
    registeredGroupChats, err := s.repo.ListAllChatGroups(ctx)
    if err != nil {
        return nil, err
    }
    
    // Merge and deduplicate
    return mergeGroupChats(botGroupChats, registeredGroupChats), nil
}
```

### 階段 3: 新增 API Endpoint

**新增 Handler 方法**:
```go
func (h *Handler) GetAllGroupChats(c *gin.Context) {
    groupChats, err := h.directoryService.GetAllGroupChats(c.Request.Context())
    if err != nil {
        response.InternalServerError(c, "Failed to get group chats", err, nil)
        return
    }
    
    response.Success(c, http.StatusOK, "Group chats retrieved successfully", groupChats)
}
```

---

## 結論

### ✅ 可以取得的方式

1. **從 `bot_installations` 表**:
   - ✅ 可以取得 Bot 有參與的所有 group chats
   - ✅ 不需要額外 API 調用
   - ✅ 不需要額外權限
   - ⚠️ 只能取得 Bot 參與的 group chats

2. **從 `chat_groups` 表**:
   - ✅ 可以取得手動註冊的 group chats
   - ✅ 可以註冊 Bot 未參與的 group chats
   - ⚠️ 需要手動註冊

### ❌ 無法取得的方式

1. **Graph API `/chats`**:
   - ❌ 不支援 application-only context
   - ❌ 即使有權限也無法使用

2. **Graph API `/users/{user-id}/chats`**:
   - ❌ 需要 delegated permissions
   - ❌ 只能取得特定使用者的 group chats
   - ❌ 不適用於 application-only 架構

### 📋 建議

**推薦方案**: 混合方案
- 從 `bot_installations` 表自動取得 Bot 參與的 group chats
- 從 `chat_groups` 表取得手動註冊的 group chats
- 提供統一的 API 查詢所有 group chats

這樣可以：
- ✅ 自動取得 Bot 參與的 group chats
- ✅ 手動註冊 Bot 未參與的 group chats
- ✅ 提供完整的 group chats 列表

---

## 測試建議

### 測試 1: 查詢 `bot_installations` 表的 group chats

```sql
-- Check how many group chats are in bot_installations
SELECT COUNT(*) as total_group_chats
FROM bot_installations
WHERE conversation_type = 'groupChat'
  AND installation_status = 'active';
```

### 測試 2: 查詢 `chat_groups` 表的 group chats

```sql
-- Check how many group chats are registered
SELECT COUNT(*) as registered_group_chats
FROM chat_groups;
```

### 測試 3: 檢查是否有重複

```sql
-- Check for duplicates between bot_installations and chat_groups
SELECT 
    bi.conversation_id,
    'bot_installations' as source
FROM bot_installations bi
WHERE bi.conversation_type = 'groupChat'
  AND bi.installation_status = 'active'
  AND EXISTS (
      SELECT 1 
      FROM chat_groups cg 
      WHERE cg.chat_id = bi.conversation_id
  );
```
