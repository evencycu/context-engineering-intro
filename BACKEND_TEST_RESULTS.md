# P0 後端 API 測試結果

## ✅ 測試時間
2026-01-13

## ✅ 測試結果總結

### 1. API Key Regeneration API
**端點**: `POST /internal/v1/projects/{projectId}/regenerate-key`

**測試結果**: ✅ 通過
- 成功重新生成 API key
- 返回格式: `{"data": {"apiKey": "proj_bBFeTE7EM6"}, "message": "API key regenerated successfully"}`
- HTTP Status: 200

**測試命令**:
```bash
curl -X POST "http://localhost:8080/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/regenerate-key" \
  -H "Content-Type: application/json"
```

### 2. Audience Lists API

#### 2.1 GET /internal/v1/projects/{projectId}/audience-lists
**測試結果**: ✅ 通過
- 成功獲取空列表（初始狀態）
- 返回格式: `{"data": null}` 或 `{"data": [...]}`
- HTTP Status: 200

#### 2.2 POST /internal/v1/projects/{projectId}/audience-lists
**測試結果**: ✅ 通過
- 成功創建 audience list
- 返回創建的列表數據，包含 id, name, type, count 等字段
- HTTP Status: 201

**測試命令**:
```bash
curl -X POST "http://localhost:8080/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/audience-lists" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Audience List", "type": "Static", "count": 100}'
```

#### 2.3 GET /internal/v1/projects/{projectId}/audience-lists/{listId}
**測試結果**: ✅ 通過
- 成功獲取單個 audience list
- 返回完整的列表數據
- HTTP Status: 200

#### 2.4 PUT /internal/v1/projects/{projectId}/audience-lists/{listId}
**測試結果**: ✅ 通過
- 成功更新 audience list
- 更新了 name 和 count 字段
- HTTP Status: 200

**測試命令**:
```bash
curl -X PUT "http://localhost:8080/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/audience-lists/{listId}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Test Audience List", "count": 200}'
```

#### 2.5 DELETE /internal/v1/projects/{projectId}/audience-lists/{listId}
**測試結果**: ✅ 通過
- 成功刪除 audience list
- 返回成功消息
- HTTP Status: 200

**測試命令**:
```bash
curl -X DELETE "http://localhost:8080/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/audience-lists/{listId}" \
  -H "Content-Type: application/json"
```

## 📊 測試統計

- **總測試數**: 6
- **通過**: 6 ✅
- **失敗**: 0 ❌
- **通過率**: 100%

## 🔧 修復的問題

1. **路由衝突**: 修復了 `/projects/:projectId` 與 `/projects/:id` 的路由衝突
   - 解決方案: 使用 `:id` 作為路由參數，使用 `:listId` 區分 audience list ID

2. **錯誤處理**: 修復了 `RegenerateKey` 方法中的錯誤處理
   - 解決方案: 正確處理 `sql.ErrNoRows` 錯誤（表示 key 可用）

3. **編譯錯誤**: 修復了未使用的變量
   - 解決方案: 移除未使用的 `existing` 變量

## ✅ 數據庫 Migration

- Migration 文件: `013_add_audience_lists_table.sql`
- 狀態: ✅ 已執行
- 表結構: ✅ 驗證通過

