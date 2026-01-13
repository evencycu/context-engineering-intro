# Template 資料庫分析報告

## 資料庫結構

### Templates 表結構
```sql
CREATE TABLE templates (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    variables JSONB DEFAULT '[]'::jsonb,
    default_json_structure TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(project_id, name)
);
```

## 當前資料庫狀態

### 1. Project Level Templates
- **位置**: `templates` 表
- **識別方式**: `project_id` 指向實際的 project UUID
- **當前資料**: 2 個 templates
  - `Welcome Message` (project_id: `750e8400-e29b-41d4-a716-446655440001`)
  - `System Alert` (project_id: `750e8400-e29b-41d4-a716-446655440001`)

### 2. Global Templates
- **位置**: `templates` 表（同一個表）
- **識別方式**: `project_id = '00000000-0000-0000-0000-000000000000'`
- **當前資料**: 0 個 templates
- **問題**: Global project 不存在於資料庫中

## 問題分析

### 問題 1: Global Project 不存在
- **原因**: `init.sql` 中的 global project 插入語句可能還沒執行
- **影響**: 無法創建 global templates（外鍵約束會失敗）

### 問題 2: 用戶創建的 Templates 沒有保存
- **可能原因**:
  1. Backend API 調用失敗
  2. 資料庫連接問題
  3. 外鍵約束失敗（如果 project 不存在）

## 解決方案

### 步驟 1: 創建 Global Project
```sql
INSERT INTO projects (id, company_id, notify_key, description, status, daily_limit, monthly_limit, priority, created_by, created_at, updated_at) 
VALUES 
('00000000-0000-0000-0000-000000000000', 
 '550e8400-e29b-41d4-a716-446655440001', 
 '__global__', 
 'Global templates project (available to all projects)', 
 'active', 0, 0, 'normal', 
 '650e8400-e29b-41d4-a716-446655440001', 
 NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
```

### 步驟 2: 檢查 Templates 創建流程
1. 檢查 Backend API 日誌
2. 檢查資料庫連接
3. 驗證外鍵約束

## 查詢 Templates 的 SQL

### 查詢所有 Templates（包含類型）
```sql
SELECT 
    t.id,
    t.name,
    t.project_id,
    p.notify_key as project_name,
    CASE 
        WHEN t.project_id = '00000000-0000-0000-0000-000000000000' 
        THEN 'Global' 
        ELSE 'Project' 
    END as template_type
FROM templates t
LEFT JOIN projects p ON t.project_id = p.id
ORDER BY t.created_at DESC;
```

### 查詢 Project Level Templates
```sql
SELECT * FROM templates 
WHERE project_id != '00000000-0000-0000-0000-000000000000'
ORDER BY created_at DESC;
```

### 查詢 Global Templates
```sql
SELECT * FROM templates 
WHERE project_id = '00000000-0000-0000-0000-000000000000'
ORDER BY created_at DESC;
```

## 資料庫位置總結

| Template 類型 | 資料表 | 識別欄位 | 當前數量 |
|--------------|--------|----------|---------|
| Project Level | `templates` | `project_id` = 實際 project UUID | 2 |
| Global | `templates` | `project_id` = `00000000-0000-0000-0000-000000000000` | 0 |

**重要**: Project level 和 Global templates 都存儲在同一個 `templates` 表中，通過 `project_id` 欄位區分。
