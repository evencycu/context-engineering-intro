# AD Sync 實現計劃

## 📊 當前狀態分析

### ✅ 已實現（後端）
1. **DirectoryService.SyncDirectory()** - 完整實現
   - 同步 Users (Azure AD)
   - 同步 Groups (Azure AD)
   - 支持 Delta Query 增量同步
   - 錯誤處理和日誌記錄

2. **API 端點** - `POST /internal/v1/directory/sync`
   - Handler 已實現
   - 異步執行（後台運行）
   - 返回 202 Accepted

3. **GraphService** - Microsoft Graph API 整合
   - GetUsers() - 獲取用戶列表
   - GetGroups() - 獲取群組列表
   - GetGroupChannels() - 獲取頻道列表（即時查詢）

### ⚠️ 部分實現（前端）
1. **apiService.syncDirectory()** - ✅ 已實現
2. **service.syncDirectory()** - ❌ 未導出
3. **UI 按鈕** - ⚠️ 只在 Admin 頁面，且可能調用錯誤的 API

### ❌ 缺失功能
1. **同步狀態查詢** - 無法查詢同步進度
2. **同步歷史記錄** - 沒有記錄同步歷史
3. **Audience 頁面整合** - 沒有 sync 按鈕和狀態顯示
4. **錯誤處理 UI** - 沒有用戶友好的錯誤提示
5. **同步進度顯示** - 沒有進度條或狀態指示

---

## 🎯 改進計劃

### Phase 1: 完善基礎功能（優先級：P0）

#### 1.1 前端 Service 層整合
**目標**: 確保 syncDirectory 可以在前端使用

**任務**:
- [ ] 在 `Frontend/services/index.ts` 導出 `syncDirectory`
- [ ] 添加錯誤處理和重試邏輯
- [ ] 添加同步狀態查詢方法

**代碼變更**:
```typescript
// Frontend/services/index.ts
syncDirectory: (): Promise<void> =>
  apiService.syncDirectory(),

// 新增：查詢同步狀態
getSyncStatus: (): Promise<SyncStatus> =>
  apiService.getSyncStatus(),
```

#### 1.2 Audience 頁面整合
**目標**: 在 Audience 頁面添加 sync 功能

**任務**:
- [ ] 添加 "Sync Directory" 按鈕
- [ ] 顯示同步狀態（未同步/同步中/已同步）
- [ ] 顯示最後同步時間
- [ ] 顯示同步統計（用戶數、群組數）

**UI 設計**:
```
┌─────────────────────────────────────────┐
│ Audience & Tags                         │
├─────────────────────────────────────────┤
│ Directory Status:                       │
│   Last Sync: 2024-01-15 10:30          │
│   Users: 150 | Groups: 25              │
│   [Sync Directory] 🔄                   │
│                                         │
│ [Target Lists] [User Tags] [Chat Groups]│
└─────────────────────────────────────────┘
```

#### 1.3 同步狀態 API（後端）
**目標**: 提供同步狀態查詢端點

**任務**:
- [ ] 添加 `GET /internal/v1/directory/sync/status` 端點
- [ ] 返回同步狀態、進度、統計信息
- [ ] 記錄最後同步時間和結果

**API 設計**:
```go
type SyncStatusResponse struct {
    IsSyncing    bool      `json:"is_syncing"`
    LastSyncAt   *time.Time `json:"last_sync_at"`
    LastSyncStatus string   `json:"last_sync_status"` // success, failed, in_progress
    UserCount     int       `json:"user_count"`
    GroupCount    int       `json:"group_count"`
    ErrorMessage  string    `json:"error_message,omitempty"`
}
```

---

### Phase 2: 增強功能（優先級：P1）

#### 2.1 同步進度顯示
**目標**: 顯示同步進度和實時統計

**任務**:
- [ ] WebSocket 或 Polling 獲取同步進度
- [ ] 顯示進度條
- [ ] 顯示當前同步的資源類型（Users/Groups）

#### 2.2 同步歷史記錄
**目標**: 記錄和顯示同步歷史

**任務**:
- [ ] 創建 `directory_sync_logs` 表
- [ ] 記錄每次同步的時間、狀態、統計
- [ ] 在 UI 顯示同步歷史

**數據庫設計**:
```sql
CREATE TABLE directory_sync_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sync_type VARCHAR(50) NOT NULL, -- 'full', 'incremental'
    status VARCHAR(50) NOT NULL, -- 'success', 'failed', 'partial'
    user_count INTEGER DEFAULT 0,
    group_count INTEGER DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### 2.3 錯誤處理改進
**目標**: 提供更好的錯誤提示和恢復機制

**任務**:
- [ ] 區分不同類型的錯誤（權限、網路、限流等）
- [ ] 提供錯誤恢復建議
- [ ] 部分失敗時顯示詳細信息

---

### Phase 3: 自動化（優先級：P2）

#### 3.1 自動同步排程
**目標**: 自動定期同步目錄

**任務**:
- [ ] 實現定時任務（cron job）
- [ ] 可配置同步頻率
- [ ] 支援增量同步（Delta Query）

#### 3.2 同步通知
**目標**: 同步完成後通知使用者

**任務**:
- [ ] 同步完成後顯示通知
- [ ] 同步失敗時發送警報
- [ ] 可選的郵件通知

---

## 🔧 技術實現細節

### 1. 同步狀態查詢實現

**後端**:
```go
// Backend/services/teamsnotification/handlers/directory/handler.go
func (h *Handler) GetSyncStatus(c *gin.Context) {
    status := h.directoryService.GetSyncStatus(c.Request.Context())
    response.Success(c, http.StatusOK, "Sync status retrieved", status)
}

// Backend/services/teamsnotification/services/directory_service.go
func (s *directoryService) GetSyncStatus(ctx context.Context) (*SyncStatus, error) {
    // 查詢數據庫獲取統計
    userCount, _ := s.repo.CountUsers(ctx)
    groupCount, _ := s.repo.CountGroups(ctx)
    
    // 獲取最後同步時間
    lastSync, _ := s.repo.GetLastSyncTime(ctx)
    
    return &SyncStatus{
        IsSyncing: s.isSyncing,
        LastSyncAt: lastSync,
        UserCount: userCount,
        GroupCount: groupCount,
    }, nil
}
```

**前端**:
```typescript
// Frontend/services/apiService.ts
getSyncStatus: async (): Promise<SyncStatus> => {
  const response = await apiCall<SyncStatus>('/directory/sync/status');
  return extractData(response)!;
},

// Frontend/pages/Audience.tsx
const [syncStatus, setSyncStatus] = useState<SyncStatus | null>(null);

useEffect(() => {
  loadSyncStatus();
}, []);

const loadSyncStatus = async () => {
  try {
    const status = await service.getSyncStatus();
    setSyncStatus(status);
  } catch (error) {
    console.error('Failed to load sync status:', error);
  }
};

const handleSyncDirectory = async () => {
  setIsSyncing(true);
  try {
    await service.syncDirectory();
    // 輪詢狀態直到完成
    pollSyncStatus();
  } catch (error) {
    // 錯誤處理
  } finally {
    setIsSyncing(false);
  }
};
```

### 2. 同步進度輪詢

```typescript
const pollSyncStatus = async () => {
  const interval = setInterval(async () => {
    const status = await service.getSyncStatus();
    setSyncStatus(status);
    
    if (!status.isSyncing) {
      clearInterval(interval);
      // 重新載入數據
      loadData();
    }
  }, 2000); // 每 2 秒輪詢一次
};
```

---

## 📋 實現檢查清單

### Phase 1: 基礎功能
- [ ] 後端：添加 `GET /directory/sync/status` 端點
- [ ] 後端：實現 `GetSyncStatus()` 方法
- [ ] 後端：添加 `CountUsers()` 和 `CountGroups()` repository 方法
- [ ] 前端：在 `apiService.ts` 添加 `getSyncStatus()`
- [ ] 前端：在 `index.ts` 導出 `syncDirectory` 和 `getSyncStatus`
- [ ] 前端：在 Audience 頁面添加 sync 按鈕和狀態顯示
- [ ] 前端：實現同步狀態輪詢
- [ ] 測試：測試同步功能
- [ ] 測試：測試狀態查詢

### Phase 2: 增強功能
- [ ] 數據庫：創建 `directory_sync_logs` 表
- [ ] 後端：記錄同步日誌
- [ ] 後端：實現同步歷史查詢 API
- [ ] 前端：顯示同步歷史
- [ ] 前端：改進錯誤處理 UI

### Phase 3: 自動化
- [ ] 後端：實現定時任務
- [ ] 後端：支援增量同步
- [ ] 配置：添加同步頻率配置

---

## 🚀 快速開始

### 步驟 1: 後端實現同步狀態 API

```bash
# 1. 添加 repository 方法
# Backend/services/teamsnotification/repositories/directory_repository.go
func (r *directoryRepository) CountUsers(ctx context.Context) (int, error)
func (r *directoryRepository) CountGroups(ctx context.Context) (int, error)
func (r *directoryRepository) GetLastSyncTime(ctx context.Context) (*time.Time, error)

# 2. 添加 service 方法
# Backend/services/teamsnotification/services/directory_service.go
func (s *directoryService) GetSyncStatus(ctx context.Context) (*SyncStatus, error)

# 3. 添加 handler
# Backend/services/teamsnotification/handlers/directory/handler.go
func (h *Handler) GetSyncStatus(c *gin.Context)
```

### 步驟 2: 前端整合

```bash
# 1. 更新 apiService.ts
# 2. 更新 index.ts
# 3. 更新 Audience.tsx
```

---

## 📊 預期效果

完成後，使用者可以：
1. ✅ 在 Audience 頁面看到目錄同步狀態
2. ✅ 一鍵觸發目錄同步
3. ✅ 查看同步進度和結果
4. ✅ 了解有多少用戶和群組已同步
5. ✅ 查看同步歷史記錄

---

**最後更新**: 2026-01-13  
**優先級**: P0 - 立即實現  
**預計時間**: 2-3 小時
