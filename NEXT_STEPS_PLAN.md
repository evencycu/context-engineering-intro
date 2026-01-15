# 接下來的開發計劃

## 📊 當前狀態總結

### ✅ 已完成 (P0)
1. **API Key Regeneration API** ✓
   - 後端實現完成
   - 前端集成完成
   - 測試通過

2. **Audience Lists API** ✓
   - 數據庫 Migration 完成
   - 後端 CRUD 實現完成
   - 前端集成完成
   - 測試通過

---

## 🎯 下一步計劃

### 【P1 - 中優先級】接下來實現

#### 1. User Tagging API
**優先級**: P1  
**預計時間**: 2-3 小時

**端點**:
- `POST /internal/v1/users/{userId}/tags/{tag}` - 添加標籤
- `DELETE /internal/v1/users/{userId}/tags/{tag}` - 移除標籤
- `GET /internal/v1/users/{userId}/tags` - 獲取用戶標籤列表

**需要完成**:
- [ ] 數據庫 Migration: 創建 `user_tags` 表（或使用 JSONB）
- [ ] Backend: 實現 UserTagRepository
- [ ] Backend: 實現 UserTagService
- [ ] Backend: 實現 UserTagHandler
- [ ] Backend: 註冊 handler 到應用
- [ ] Frontend: 更新 `apiService.ts` 添加標籤 API
- [ ] Frontend: 更新 `Audience.tsx` 頁面使用新 API
- [ ] 測試: 後端 API 測試
- [ ] 測試: 前端集成測試

**影響頁面**: `Frontend/pages/Audience.tsx` (User Tags tab)

---

#### 2. Channel Sync API
**優先級**: P1  
**預計時間**: 1-2 小時

**端點**:
- `POST /internal/v1/directory/groups/{groupId}/channels/sync` - 同步 Teams 頻道

**需要完成**:
- [ ] Backend: 檢查現有的 directory sync 機制
- [ ] Backend: 實現 Channel Sync handler
- [ ] Backend: 整合現有的 GraphService
- [ ] Backend: 註冊 handler 到應用
- [ ] Frontend: 更新 `apiService.ts` 添加 sync API
- [ ] Frontend: 更新相關頁面使用新 API
- [ ] 測試: 後端 API 測試

**影響頁面**: Directory 管理相關頁面

---

### 【P2 - 低優先級】後續實現

#### 3. System Actions API
**優先級**: P2  
**預計時間**: 2-3 小時

**端點**:
- `POST /internal/v1/admin/system/actions` - 執行系統操作

**功能**:
- `sync_aad` - 同步 Azure AD
- `restart_service` - 重啟服務
- `block_ip` - 封鎖 IP

**需要完成**:
- [ ] Backend: 實現 SystemActions handler
- [ ] Backend: 實現各類系統操作邏輯
- [ ] Backend: 添加 audit_logs 記錄
- [ ] Backend: 註冊 handler 到應用
- [ ] Frontend: 更新 `apiService.ts` 添加 system actions API
- [ ] Frontend: 更新 Admin 頁面使用新 API
- [ ] 測試: 後端 API 測試

**影響頁面**: `Frontend/pages/Admin.tsx`

---

#### 4. Global Template Assignment API
**優先級**: P2  
**預計時間**: 2-3 小時

**端點**:
- `POST /internal/v1/global/templates/{templateId}/assign` - 分配全局模板到項目

**需要完成**:
- [ ] 數據庫: 決定是否需要 `template_assignments` 表
- [ ] Backend: 實現 TemplateAssignment handler
- [ ] Backend: 註冊 handler 到應用
- [ ] Frontend: 更新 `apiService.ts` 添加 assignment API
- [ ] Frontend: 更新相關頁面使用新 API
- [ ] 測試: 後端 API 測試

**影響頁面**: `Frontend/pages/Admin.tsx` (Global Templates)

---

### 【測試與文檔】持續進行

#### 5. 單元測試
**優先級**: 中  
**預計時間**: 3-4 小時

**需要完成**:
- [ ] 為 User Tagging API 添加單元測試
- [ ] 為 Channel Sync API 添加單元測試
- [ ] 為 System Actions API 添加單元測試
- [ ] 為 Global Template Assignment API 添加單元測試
- [ ] 確保測試覆蓋率 ≥ 80%

**文件位置**: `Backend/services/teamsnotification/handlers/*/basic_test.go`

---

#### 6. OpenAPI 文檔更新
**優先級**: 中  
**預計時間**: 1-2 小時

**需要完成**:
- [ ] 更新 OpenAPI 規範添加 User Tagging API
- [ ] 更新 OpenAPI 規範添加 Channel Sync API
- [ ] 更新 OpenAPI 規範添加 System Actions API
- [ ] 更新 OpenAPI 規範添加 Global Template Assignment API
- [ ] 驗證 OpenAPI 規範與實現同步

**文件位置**: `Backend/api/openapi/teams-notification-api.yaml`

---

## 📅 建議執行順序

### 第一階段：P1 APIs (本週)
1. **User Tagging API** (2-3 小時)
   - 影響 Audience 頁面，用戶體驗優先
   - 相對簡單，可以先完成

2. **Channel Sync API** (1-2 小時)
   - 依賴現有 directory sync 機制
   - 可以快速實現

**預計完成時間**: 3-5 小時

---

### 第二階段：測試與文檔 (本週末)
1. **單元測試** (3-4 小時)
   - 為 P1 APIs 添加測試
   - 確保代碼質量

2. **OpenAPI 文檔更新** (1-2 小時)
   - 更新文檔保持同步

**預計完成時間**: 4-6 小時

---

### 第三階段：P2 APIs (下週)
1. **System Actions API** (2-3 小時)
   - Admin 功能，優先級較低

2. **Global Template Assignment API** (2-3 小時)
   - Template 管理增強功能

**預計完成時間**: 4-6 小時

---

## 🎯 本週目標

### 主要目標
- ✅ 完成 P0 APIs (已完成)
- 🎯 完成 P1 APIs (User Tagging + Channel Sync)
- 🎯 添加單元測試
- 🎯 更新 OpenAPI 文檔

### 次要目標
- 開始 P2 APIs 實現
- 優化現有代碼
- 改進錯誤處理

---

## 📝 開發工作流程

### 每個 API 實現步驟
1. **數據庫 Migration** (如果需要)
   - 創建 migration 文件
   - 執行 migration
   - 驗證表結構

2. **Backend 實現**
   - 創建 Repository (如果需要)
   - 創建 Service
   - 創建 Handler
   - 註冊到應用

3. **Frontend 集成**
   - 更新 `apiService.ts`
   - 更新相關頁面組件
   - 測試 UI 功能

4. **測試**
   - 後端 API 測試 (curl/Postman)
   - 前端集成測試 (瀏覽器)
   - 單元測試 (Go test)

5. **文檔**
   - 更新 OpenAPI 規範
   - 更新 README (如果需要)

---

## 🔍 技術要點

### User Tagging API
- **數據庫設計**: 考慮使用 JSONB (`users.tags`) 或獨立表 (`user_tags`)
- **性能**: 如果使用 JSONB，需要添加 GIN 索引
- **查詢**: 支持按標籤查詢用戶

### Channel Sync API
- **依賴**: 現有的 `GraphService` 和 `DirectoryService`
- **同步邏輯**: 調用 Microsoft Graph API 獲取頻道列表
- **錯誤處理**: 處理 Graph API 限流和錯誤

### System Actions API
- **安全性**: 需要管理員權限驗證
- **審計**: 所有操作記錄到 audit_logs
- **異步**: 某些操作可能需要異步執行

### Global Template Assignment API
- **數據模型**: 決定是否需要中間表
- **權限**: 只有全局管理員可以分配模板
- **驗證**: 確保模板和項目存在

---

## 📊 進度追蹤

### 已完成
- [x] P0: API Key Regeneration API
- [x] P0: Audience Lists API
- [x] 數據庫 Migration (audience_lists)
- [x] 前端集成 (Settings, Audience)

### 進行中
- [ ] P1: User Tagging API
- [ ] P1: Channel Sync API

### 待開始
- [ ] P2: System Actions API
- [ ] P2: Global Template Assignment API
- [ ] 單元測試
- [ ] OpenAPI 文檔更新

---

## 🚀 快速開始

### 開始實現 User Tagging API

```bash
# 1. 創建數據庫 migration
cd Backend/scripts/database/migrations
# 創建 014_add_user_tags_table.sql

# 2. 執行 migration
docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < 014_add_user_tags_table.sql

# 3. 創建 Repository
# Backend/services/teamsnotification/repositories/user_tag_repository.go

# 4. 創建 Service
# Backend/services/teamsnotification/services/service.go (添加 UserTagService)

# 5. 創建 Handler
# Backend/services/teamsnotification/handlers/users/handler.go (添加標籤方法)

# 6. 註冊到應用
# Backend/pkg/app/app.go (註冊 handler)

# 7. 更新前端
# Frontend/services/apiService.ts (添加標籤 API)
# Frontend/pages/Audience.tsx (使用新 API)

# 8. 測試
# 後端: curl 測試
# 前端: 瀏覽器測試
```

---

## 📞 需要協助？

如果遇到問題：
1. 檢查現有類似實現（如 Audience Lists API）
2. 參考 OpenAPI 規範
3. 查看測試文件了解預期行為
4. 檢查數據庫 schema

---

**最後更新**: 2026-01-13  
**狀態**: P0 完成，準備開始 P1
