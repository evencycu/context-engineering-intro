# 任務計劃：實現缺失的 Backend APIs

## 優先級排序

### 【P0 - 高優先級】立即實現

#### 1. API Key Regeneration API
- **端點**: `POST /internal/v1/projects/{projectId}/regenerate-key`
- **功能**: 重新生成項目的 notify_key
- **數據庫**: 更新 `projects.notify_key` 字段
- **影響**: Settings 頁面需要此功能

#### 2. Audience Lists API
- **端點**: 
  - `GET /internal/v1/projects/{projectId}/audience-lists`
  - `POST /internal/v1/projects/{projectId}/audience-lists`
- **功能**: 管理受眾列表（上傳、查詢）
- **數據庫**: 需要創建 `audience_lists` 表
- **影響**: Audience 頁面核心功能

### 【P1 - 中優先級】後續實現

#### 3. User Tagging API
- **端點**: 
  - `POST /internal/v1/users/{userId}/tags/{tag}`
  - `DELETE /internal/v1/users/{userId}/tags/{tag}`
- **功能**: 給用戶添加/移除標籤
- **數據庫**: 需要創建 `user_tags` 表或使用 JSONB
- **影響**: Audience 頁面標籤管理

#### 4. Channel Sync API
- **端點**: `POST /internal/v1/directory/groups/{groupId}/channels/sync`
- **功能**: 同步 Teams 頻道列表
- **數據庫**: 不需要新表，使用現有的 directory sync 機制
- **影響**: Directory 管理功能

### 【P2 - 低優先級】可選實現

#### 5. System Actions API
- **端點**: `POST /internal/v1/admin/system/actions`
- **功能**: 執行系統操作（sync_aad, restart_service, block_ip）
- **數據庫**: 記錄到 `audit_logs` 或 `system_logs`
- **影響**: Admin 系統操作

#### 6. Global Template Assignment API
- **端點**: `POST /internal/v1/global/templates/{templateId}/assign`
- **功能**: 分配全局模板到項目
- **數據庫**: 可能需要 `template_assignments` 表
- **影響**: Global Template 分配功能

## 實現步驟

### Phase 1: 數據庫 Migration
1. 創建 `audience_lists` 表
2. 創建 `user_tags` 表（如果需要）
3. 添加必要的索引

### Phase 2: Backend 實現
1. 實現 API Key Regeneration handler
2. 實現 Audience Lists handlers
3. 實現 User Tagging handlers
4. 實現 Channel Sync handler
5. 實現 System Actions handler
6. 實現 Global Template Assignment handler

### Phase 3: Frontend 整合
1. 更新 `apiService.ts` 使用新的 backend APIs
2. 移除 mock fallback
3. 更新相關頁面組件

### Phase 4: 測試
1. 單元測試
2. 整合測試
3. E2E 測試

## 預計時間

- Phase 1: 1-2 小時
- Phase 2: 4-6 小時
- Phase 3: 2-3 小時
- Phase 4: 2-3 小時

**總計**: 約 9-14 小時

