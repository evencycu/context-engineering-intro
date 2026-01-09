# Teams Bot Platform - 使用範例

## Destinations 的 Targets 陣列設計

### 設計概念

新的 `destinations` 表使用 `targets` JSONB 陣列來支援多個 Teams 目標的組合，讓一個目的地可以包含：

- 多個 channel
- 多個 person  
- 多個 chatgroup
- 或以上任意組合

### 資料結構

```sql
-- destinations 表結構
CREATE TABLE destinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    teams_tenant_id VARCHAR(255) NOT NULL,
    targets JSONB NOT NULL DEFAULT '[]'::jsonb, -- Array of Teams targets
    bot_id UUID, -- 主要使用的 Bot
    -- bot_type removed; destination always implies platform bot
    status VARCHAR(20) DEFAULT 'active',
    validation_status VARCHAR(20) DEFAULT 'pending',
    last_validated_at TIMESTAMP WITH TIME ZONE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_targets_array CHECK (jsonb_array_length(targets) > 0)
);
```

### Go 模型

```go
// TeamsTarget represents a single Teams target
type TeamsTarget struct {
    Type         string `json:"type"` // person, channel, chatgroup
    TeamID       string `json:"team_id,omitempty"`
    ChannelID    string `json:"channel_id,omitempty"`
    UserID       string `json:"user_id,omitempty"`
    GroupID      string `json:"group_id,omitempty"`
    DisplayName  string `json:"display_name,omitempty"`
    Description  string `json:"description,omitempty"`
    Metadata     map[string]any `json:"metadata,omitempty"`
}

// Destination represents a Teams target group (can contain multiple targets)
type Destination struct {
    BaseModel
    ProjectID        uuid.UUID      `json:"project_id" db:"project_id"`
    Name             string         `json:"name" db:"name"`
    Description      string         `json:"description" db:"description"`
    TeamsTenantID    string         `json:"teams_tenant_id" db:"teams_tenant_id"`
    Targets          []TeamsTarget  `json:"targets" db:"targets"`
    BotID            *uuid.UUID     `json:"bot_id" db:"bot_id"`
    BotType          *BotType       `json:"bot_type" db:"bot_type"`
    Status           string         `json:"status" db:"status"`
    ValidationStatus string         `json:"validation_status" db:"validation_status"`
    LastValidatedAt  *time.Time     `json:"last_validated_at" db:"last_validated_at"`
    CreatedBy        uuid.UUID      `json:"created_by" db:"created_by"`
}
```

## 使用範例

### 1. 多個 Channel 的目的地

```json
{
  "name": "開發團隊通知群組",
  "description": "包含所有開發相關的 Teams channels",
  "teams_tenant_id": "tenant-123",
  "targets": [
    {
      "type": "channel",
      "team_id": "team-1",
      "channel_id": "channel-dev-general",
      "display_name": "開發討論區",
      "description": "一般開發討論"
    },
    {
      "type": "channel", 
      "team_id": "team-1",
      "channel_id": "channel-dev-alerts",
      "display_name": "系統警報",
      "description": "系統監控警報"
    },
    {
      "type": "channel",
      "team_id": "team-2", 
      "channel_id": "channel-qa-testing",
      "display_name": "QA 測試",
      "description": "測試相關討論"
    }
  ],
  "bot_id": "bot-platform-1",
  "bot_type": "platform"
}
```

### 2. 混合類型的目的地

```json
{
  "name": "緊急通知群組",
  "description": "包含重要人員和頻道的緊急通知",
  "teams_tenant_id": "tenant-123",
  "targets": [
    {
      "type": "person",
      "user_id": "user-manager-1",
      "display_name": "張經理",
      "description": "技術經理"
    },
    {
      "type": "person",
      "user_id": "user-lead-1", 
      "display_name": "李組長",
      "description": "技術組長"
    },
    {
      "type": "channel",
      "team_id": "team-1",
      "channel_id": "channel-emergency",
      "display_name": "緊急事件頻道",
      "description": "緊急事件處理"
    },
    {
      "type": "chatgroup",
      "group_id": "group-oncall",
      "display_name": "On-call 群組",
      "description": "值班人員群組"
    }
  ],
  "bot_id": "bot-platform-1",
  "bot_type": "platform"
}
```

### 3. 多個 Person 的目的地

```json
{
  "name": "管理層通知",
  "description": "所有管理層人員",
  "teams_tenant_id": "tenant-123", 
  "targets": [
    {
      "type": "person",
      "user_id": "user-ceo",
      "display_name": "王執行長",
      "description": "執行長"
    },
    {
      "type": "person",
      "user_id": "user-cto",
      "display_name": "陳技術長", 
      "description": "技術長"
    },
    {
      "type": "person",
      "user_id": "user-cfo",
      "display_name": "林財務長",
      "description": "財務長"
    }
  ],
  "bot_id": "bot-platform-1",
  "bot_type": "platform"
}
```

### 4. 多個 Chatgroup 的目的地

```json
{
  "name": "專案群組通知",
  "description": "所有專案相關的群組",
  "teams_tenant_id": "tenant-123",
  "targets": [
    {
      "type": "chatgroup",
      "group_id": "group-project-alpha",
      "display_name": "Alpha 專案群組",
      "description": "Alpha 專案開發團隊"
    },
    {
      "type": "chatgroup", 
      "group_id": "group-project-beta",
      "display_name": "Beta 專案群組",
      "description": "Beta 專案開發團隊"
    },
    {
      "type": "chatgroup",
      "group_id": "group-project-gamma", 
      "display_name": "Gamma 專案群組",
      "description": "Gamma 專案開發團隊"
    }
  ],
  "bot_id": "bot-platform-1",
  "bot_type": "platform"
}
```

## 查詢範例

### 1. 查詢包含特定 channel 的所有目的地

```sql
SELECT * FROM destinations 
WHERE targets @> '[{"type": "channel", "channel_id": "channel-dev-general"}]';
```

### 2. 查詢包含特定 user 的所有目的地

```sql
SELECT * FROM destinations 
WHERE targets @> '[{"type": "person", "user_id": "user-manager-1"}]';
```

### 3. 查詢包含多個 targets 的目的地

```sql
SELECT * FROM destinations 
WHERE targets @> '[{"type": "channel"}, {"type": "person"}]';
```

### 4. 統計每個目的地的 target 數量

```sql
SELECT 
    id,
    name,
    jsonb_array_length(targets) as target_count,
    targets
FROM destinations;
```

## 優點

1. **靈活性**: 一個目的地可以包含任意數量和類型的 Teams 目標
2. **效率**: 減少目的地數量，簡化管理
3. **擴展性**: 容易添加新的目標類型
4. **查詢能力**: 使用 PostgreSQL JSONB 的強大查詢功能
5. **一致性**: 所有相關目標在一個記錄中，確保一致性

## 注意事項

1. **驗證**: 需要確保每個 target 的必填欄位都有值
2. **索引**: 已建立適當的 GIN 索引來支援 JSONB 查詢
3. **約束**: 使用 CHECK 約束確保 targets 陣列不為空
4. **效能**: 大量 targets 可能影響查詢效能，需要適當的索引策略
