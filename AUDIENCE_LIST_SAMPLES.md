# Target List (Audience List) 樣本數據與數據庫格式

## 📊 數據庫表結構

### 表名: `audience_lists`

```sql
CREATE TABLE audience_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) DEFAULT 'Static' CHECK (type IN ('Static', 'Dynamic')),
    count INTEGER DEFAULT 0,
    description TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, name)
);
```

### 索引
- `idx_audience_lists_project_id` - 按項目查詢
- `idx_audience_lists_type` - 按類型查詢
- `idx_audience_lists_last_updated` - 按更新時間查詢

### 約束
- `UNIQUE(project_id, name)` - 同一項目內名稱唯一
- `CHECK (type IN ('Static', 'Dynamic'))` - 類型限制
- `FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE` - 外鍵約束

---

## 📝 樣本數據

### 樣本 1: Static List (靜態列表)

#### SQL INSERT 語句
```sql
INSERT INTO audience_lists (
    id,
    project_id,
    name,
    type,
    count,
    description,
    metadata,
    last_updated,
    created_at,
    updated_at
) VALUES (
    'a1b2c3d4-e5f6-7890-abcd-ef1234567890',
    '750e8400-e29b-41d4-a716-446655440001',
    'Marketing Team',
    'Static',
    150,
    'Marketing department team members for promotional notifications',
    '{
        "source": "manual_upload",
        "uploaded_by": "user@example.com",
        "file_name": "marketing_team_2024.csv",
        "tags": ["marketing", "promotions"],
        "last_sync": "2024-01-15T10:30:00Z"
    }'::jsonb,
    '2024-01-15T10:30:00Z',
    '2024-01-15T10:30:00Z',
    '2024-01-15T10:30:00Z'
);
```

#### 數據庫中的實際數據
```
id: a1b2c3d4-e5f6-7890-abcd-ef1234567890
project_id: 750e8400-e29b-41d4-a716-446655440001
name: Marketing Team
type: Static
count: 150
description: Marketing department team members for promotional notifications
metadata: {"source": "manual_upload", "uploaded_by": "user@example.com", "file_name": "marketing_team_2024.csv", "tags": ["marketing", "promotions"], "last_sync": "2024-01-15T10:30:00Z"}
last_updated: 2024-01-15 10:30:00+00
created_at: 2024-01-15 10:30:00+00
updated_at: 2024-01-15 10:30:00+00
```

#### JSON 格式 (API 響應)
```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "project_id": "750e8400-e29b-41d4-a716-446655440001",
  "name": "Marketing Team",
  "type": "Static",
  "count": 150,
  "description": "Marketing department team members for promotional notifications",
  "metadata": {
    "source": "manual_upload",
    "uploaded_by": "user@example.com",
    "file_name": "marketing_team_2024.csv",
    "tags": ["marketing", "promotions"],
    "last_sync": "2024-01-15T10:30:00Z"
  },
  "last_updated": "2024-01-15T10:30:00Z",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

---

### 樣本 2: Dynamic List (動態列表)

#### SQL INSERT 語句
```sql
INSERT INTO audience_lists (
    id,
    project_id,
    name,
    type,
    count,
    description,
    metadata,
    last_updated,
    created_at,
    updated_at
) VALUES (
    'b2c3d4e5-f6a7-8901-bcde-f23456789012',
    '750e8400-e29b-41d4-a716-446655440001',
    'Active Users (Last 30 Days)',
    'Dynamic',
    0,
    'Users who have been active in the last 30 days, automatically updated',
    '{
        "query_type": "time_based",
        "query": "SELECT user_id FROM user_activity WHERE last_active > NOW() - INTERVAL ''30 days''",
        "refresh_interval": "daily",
        "last_refresh": "2024-01-15T00:00:00Z",
        "auto_update": true
    }'::jsonb,
    '2024-01-15T00:00:00Z',
    '2024-01-14T15:20:00Z',
    '2024-01-15T00:00:00Z'
);
```

#### JSON 格式 (API 響應)
```json
{
  "id": "b2c3d4e5-f6a7-8901-bcde-f23456789012",
  "project_id": "750e8400-e29b-41d4-a716-446655440001",
  "name": "Active Users (Last 30 Days)",
  "type": "Dynamic",
  "count": 0,
  "description": "Users who have been active in the last 30 days, automatically updated",
  "metadata": {
    "query_type": "time_based",
    "query": "SELECT user_id FROM user_activity WHERE last_active > NOW() - INTERVAL '30 days'",
    "refresh_interval": "daily",
    "last_refresh": "2024-01-15T00:00:00Z",
    "auto_update": true
  },
  "last_updated": "2024-01-15T00:00:00Z",
  "created_at": "2024-01-14T15:20:00Z",
  "updated_at": "2024-01-15T00:00:00Z"
}
```

---

### 樣本 3: Simple Static List (簡單靜態列表)

#### SQL INSERT 語句
```sql
INSERT INTO audience_lists (
    project_id,
    name,
    type,
    count
) VALUES (
    '750e8400-e29b-41d4-a716-446655440001',
    'Test List',
    'Static',
    100
);
```

#### JSON 格式 (API 響應)
```json
{
  "id": "c3d4e5f6-a7b8-9012-cdef-345678901234",
  "project_id": "750e8400-e29b-41d4-a716-446655440001",
  "name": "Test List",
  "type": "Static",
  "count": 100,
  "description": null,
  "metadata": {},
  "last_updated": "2024-01-15T12:00:00Z",
  "created_at": "2024-01-15T12:00:00Z",
  "updated_at": "2024-01-15T12:00:00Z"
}
```

---

## 🔧 Go 模型結構

### 定義位置: `Backend/libs/models/models.go`

```go
type AudienceList struct {
    BaseModel
    ProjectID   uuid.UUID       `json:"project_id" db:"project_id"`
    Name        string           `json:"name" db:"name"`
    Type        string           `json:"type" db:"type"` // 'Static' or 'Dynamic'
    Count       int              `json:"count" db:"count"`
    Description *string          `json:"description" db:"description"`
    Metadata    JSONBObject      `json:"metadata" db:"metadata"`
    LastUpdated time.Time        `json:"last_updated" db:"last_updated"`
}
```

### BaseModel 結構
```go
type BaseModel struct {
    ID        uuid.UUID `json:"id" db:"id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

---

## 📡 API 請求/響應格式

### GET /internal/v1/projects/{projectId}/audience-lists

#### 請求
```http
GET /internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/audience-lists
Content-Type: application/json
```

#### 響應 (成功)
```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "project_id": "750e8400-e29b-41d4-a716-446655440001",
      "name": "Marketing Team",
      "type": "Static",
      "count": 150,
      "description": "Marketing department team members",
      "metadata": {
        "source": "manual_upload",
        "tags": ["marketing"]
      },
      "last_updated": "2024-01-15T10:30:00Z",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": "b2c3d4e5-f6a7-8901-bcde-f23456789012",
      "project_id": "750e8400-e29b-41d4-a716-446655440001",
      "name": "Active Users",
      "type": "Dynamic",
      "count": 0,
      "description": null,
      "metadata": {
        "query_type": "time_based"
      },
      "last_updated": "2024-01-15T00:00:00Z",
      "created_at": "2024-01-14T15:20:00Z",
      "updated_at": "2024-01-15T00:00:00Z"
    }
  ]
}
```

#### 響應 (空列表)
```json
{
  "data": null
}
```

---

### POST /internal/v1/projects/{projectId}/audience-lists

#### 請求
```http
POST /internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/audience-lists
Content-Type: application/json

{
  "name": "Sales Team",
  "type": "Static",
  "count": 200,
  "description": "Sales department team members"
}
```

#### 響應 (成功 - 201 Created)
```json
{
  "data": {
    "id": "d4e5f6a7-b8c9-0123-def4-567890123456",
    "project_id": "750e8400-e29b-41d4-a716-446655440001",
    "name": "Sales Team",
    "type": "Static",
    "count": 200,
    "description": "Sales department team members",
    "metadata": {},
    "last_updated": "2024-01-15T14:00:00Z",
    "created_at": "2024-01-15T14:00:00Z",
    "updated_at": "2024-01-15T14:00:00Z"
  },
  "message": "Audience list created successfully"
}
```

---

### PUT /internal/v1/projects/{projectId}/audience-lists/{listId}

#### 請求
```http
PUT /internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/audience-lists/a1b2c3d4-e5f6-7890-abcd-ef1234567890
Content-Type: application/json

{
  "name": "Marketing Team Updated",
  "count": 175,
  "description": "Updated marketing team list"
}
```

#### 響應 (成功)
```json
{
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "project_id": "750e8400-e29b-41d4-a716-446655440001",
    "name": "Marketing Team Updated",
    "type": "Static",
    "count": 175,
    "description": "Updated marketing team list",
    "metadata": {
      "source": "manual_upload",
      "tags": ["marketing"]
    },
    "last_updated": "2024-01-15T10:30:00Z",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T15:00:00Z"
  },
  "message": "Audience list updated successfully"
}
```

---

## 📋 字段說明

| 字段 | 類型 | 必填 | 說明 |
|------|------|------|------|
| `id` | UUID | 自動生成 | 唯一標識符 |
| `project_id` | UUID | 是 | 所屬項目 ID |
| `name` | String(255) | 是 | 列表名稱（項目內唯一） |
| `type` | String(50) | 否 | 類型：'Static' 或 'Dynamic'，默認 'Static' |
| `count` | Integer | 否 | 受眾數量，默認 0 |
| `description` | Text | 否 | 列表描述 |
| `metadata` | JSONB | 否 | 元數據（JSON 對象），默認 {} |
| `last_updated` | Timestamp | 自動 | 最後更新時間 |
| `created_at` | Timestamp | 自動 | 創建時間 |
| `updated_at` | Timestamp | 自動 | 更新時間（觸發器自動更新） |

---

## 💡 Metadata 字段使用示例

### Static List Metadata
```json
{
  "source": "manual_upload",
  "uploaded_by": "user@example.com",
  "file_name": "team_list.csv",
  "file_size": 10240,
  "tags": ["marketing", "promotions"],
  "last_sync": "2024-01-15T10:30:00Z",
  "upload_date": "2024-01-15T10:30:00Z"
}
```

### Dynamic List Metadata
```json
{
  "query_type": "time_based",
  "query": "SELECT user_id FROM user_activity WHERE last_active > NOW() - INTERVAL '30 days'",
  "refresh_interval": "daily",
  "last_refresh": "2024-01-15T00:00:00Z",
  "auto_update": true,
  "refresh_time": "00:00:00"
}
```

### Tag-based List Metadata
```json
{
  "query_type": "tag_based",
  "tags": ["developer", "backend"],
  "include_all": false,
  "last_refresh": "2024-01-15T12:00:00Z"
}
```

---

## 🔍 查詢示例

### 查詢特定項目的所有列表
```sql
SELECT * FROM audience_lists 
WHERE project_id = '750e8400-e29b-41d4-a716-446655440001'
ORDER BY created_at DESC;
```

### 查詢 Static 類型的列表
```sql
SELECT * FROM audience_lists 
WHERE project_id = '750e8400-e29b-41d4-a716-446655440001'
  AND type = 'Static'
ORDER BY name;
```

### 查詢包含特定 metadata 的列表
```sql
SELECT * FROM audience_lists 
WHERE project_id = '750e8400-e29b-41d4-a716-446655440001'
  AND metadata @> '{"tags": ["marketing"]}'::jsonb;
```

### 查詢最近更新的列表
```sql
SELECT * FROM audience_lists 
WHERE project_id = '750e8400-e29b-41d4-a716-446655440001'
ORDER BY last_updated DESC
LIMIT 10;
```

---

## 🧪 測試數據插入腳本

```sql
-- 插入測試數據
INSERT INTO audience_lists (project_id, name, type, count, description, metadata) VALUES
('750e8400-e29b-41d4-a716-446655440001', 'Marketing Team', 'Static', 150, 'Marketing department', '{"source": "manual"}'::jsonb),
('750e8400-e29b-41d4-a716-446655440001', 'Sales Team', 'Static', 200, 'Sales department', '{"source": "manual"}'::jsonb),
('750e8400-e29b-41d4-a716-446655440001', 'Active Users', 'Dynamic', 0, 'Active users in last 30 days', '{"query_type": "time_based"}'::jsonb),
('750e8400-e29b-41d4-a716-446655440001', 'Developers', 'Static', 50, 'Development team', '{"tags": ["developer"]}'::jsonb);
```

---

**最後更新**: 2026-01-13
