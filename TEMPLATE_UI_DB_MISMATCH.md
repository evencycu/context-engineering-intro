# Template UI 與 DB 不一致問題分析

## 問題描述
- **資料庫**: 只有 2 筆 templates（正確）
- **UI 顯示**: 與資料庫不同（可能顯示更多或不同的 templates）

## 根本原因

### 1. Frontend 使用 Mock Data
Frontend 有 fallback 機制，當 API 調用失敗時會使用 mock data：

```typescript
// Frontend/services/index.ts
const USE_MOCK_ONLY = false; // 設為 false，但會 fallback
const ENABLE_FALLBACK = true; // 如果 API 失敗，會使用 mock

getTemplates: (projectId: string): Promise<Template[]> =>
  withFallback(
    () => apiService.getTemplates(projectId),  // 真實 API
    () => mockApiService.getTemplates(projectId), // Mock data
    'getTemplates'
  ),
```

### 2. Mock Data 與 DB 資料不匹配

**Mock Data** (`Frontend/constants.ts`):
- Project IDs: `proj_fin_01`, `proj_hr_02`, `proj_sec_03`, `proj_nw_01`
- Templates: 2 個 templates（屬於 `proj_fin_01` 和 `proj_hr_02`）

**資料庫**:
- Project IDs: `750e8400-e29b-41d4-a716-446655440001`, `750e8400-e29b-41d4-a716-446655440002`
- Templates: 2 個 templates（屬於 `750e8400-e29b-41d4-a716-446655440001`）

### 3. Global Templates API 404 錯誤

Global templates API 返回 404，可能是因為：
- 路由沒有正確註冊
- 或者 backend server 需要重啟

## 解決方案

### 方案 1: 確保使用真實 API（推薦）

1. **檢查 Backend 是否運行**:
   ```bash
   curl http://localhost:8080/health
   ```

2. **檢查 API 是否可訪問**:
   ```bash
   curl http://localhost:8080/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/templates
   ```

3. **檢查 Frontend Console**:
   - 打開瀏覽器開發者工具
   - 查看 Console 是否有 API 錯誤
   - 查看 Network 標籤，確認 API 調用是否成功

### 方案 2: 禁用 Mock Fallback（僅用於調試）

在 `Frontend/services/index.ts` 中：
```typescript
const ENABLE_FALLBACK = false; // 禁用 fallback，強制使用真實 API
```

### 方案 3: 修復 Global Templates API

Global templates API 返回 404，需要檢查：
1. Backend server 是否已重啟（載入新的路由）
2. 路由是否正確註冊

## 驗證步驟

### 1. 檢查 Backend API
```bash
# 檢查 project templates
curl http://localhost:8080/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/templates

# 應該返回 2 個 templates
```

### 2. 檢查 Frontend Console
- 打開瀏覽器開發者工具
- 查看 Console 日誌：
  - `[getTemplates] API call failed` → 使用 mock data
  - `[API] Backend is available` → 使用真實 API

### 3. 檢查 Network 請求
- 打開 Network 標籤
- 查找 `/internal/v1/projects/.../templates` 請求
- 檢查狀態碼：
  - `200` → 成功，使用真實 API
  - `404/500` → 失敗，fallback 到 mock

## 預期行為

### 如果使用真實 API:
- UI 應該顯示 **2 個 templates**（與資料庫一致）
- Templates 屬於 project `750e8400-e29b-41d4-a716-446655440001`

### 如果使用 Mock Data:
- UI 可能顯示 **2 個 templates**（但內容不同）
- Templates 屬於 mock projects (`proj_fin_01`, `proj_hr_02`)

## 建議行動

1. **立即**: 檢查 Frontend Console，確認是否使用 mock data
2. **修復**: 如果 API 失敗，修復 backend 連接問題
3. **驗證**: 確認 UI 顯示的 templates 與資料庫一致
