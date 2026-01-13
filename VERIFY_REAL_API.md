# 確認使用真實 API 指南

## ✅ Backend 狀態確認

### 1. Backend 健康檢查
```bash
curl http://localhost:8080/health
```
**結果**: ✅ 正常
```json
{"status":"healthy","timestamp":"2026-01-13T03:38:48.950906Z","version":"1.0.0"}
```

### 2. Templates API 測試
```bash
curl http://localhost:8080/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/templates
```
**結果**: ✅ 成功，返回 2 筆 templates

## 🔍 Frontend 檢查步驟

### 步驟 1: 打開 Templates 頁面
1. 訪問: `http://localhost:3002/#/templates`
2. 打開開發者工具 (F12)

### 步驟 2: 檢查 Console 日誌

**如果使用真實 API，會看到**:
```
[API] Backend is available
```

**如果使用 Mock Data，會看到**:
```
[getTemplates] API call failed: ...
[getTemplates] Falling back to mock data
```

### 步驟 3: 檢查 Network 請求

1. 切換到 **Network** 標籤
2. 重新載入頁面 (Cmd+R 或 Ctrl+R)
3. 查找請求: `/internal/v1/projects/.../templates`

**檢查項目**:
- **Status**: 應該是 `200` (成功)
- **Response**: 應該包含 `code: 200` 和 `data` 陣列
- **Request URL**: 應該是 `http://localhost:3002/internal/v1/projects/.../templates`

### 步驟 4: 執行 API 連接檢查腳本

在 Console 中執行 `check_api_connection.js` 的內容，或直接執行：

```javascript
// 快速檢查
fetch('/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/templates')
  .then(r => r.json())
  .then(d => {
    console.log('API 狀態:', d.code === 200 ? '✅ 成功' : '❌ 失敗');
    console.log('Templates 數量:', d.data?.length || 0);
    console.log('資料:', d);
  })
  .catch(e => console.error('❌ API 錯誤:', e));
```

## 🐛 常見問題

### 問題 1: API 返回 404
**原因**: Vite proxy 配置問題或路由未註冊
**解決**:
1. 檢查 `Frontend/vite.config.ts` 的 proxy 配置
2. 確認 backend server 運行在 `http://localhost:8080`

### 問題 2: CORS 錯誤
**原因**: Backend CORS 配置問題
**解決**: Backend 已配置 CORS，應該不會有此問題

### 問題 3: 自動 Fallback 到 Mock
**原因**: API 調用失敗，`ENABLE_FALLBACK = true`
**解決**:
1. 檢查 Network 請求的錯誤訊息
2. 確認 backend server 是否運行
3. 臨時禁用 fallback（僅調試）:
   ```typescript
   // Frontend/services/index.ts
   const ENABLE_FALLBACK = false;
   ```

## 📊 預期結果

### 使用真實 API 時:
- **Console**: `[API] Backend is available`
- **Network**: Status `200`, Response 包含 `data` 陣列
- **UI**: 顯示 2 個 templates（與資料庫一致）

### 使用 Mock Data 時:
- **Console**: `[getTemplates] Falling back to mock data`
- **Network**: 沒有 `/internal/v1/...` 請求，或請求失敗
- **UI**: 可能顯示不同的 templates（來自 mock data）

## ✅ 確認清單

- [ ] Backend server 運行在 `http://localhost:8080`
- [ ] Frontend 運行在 `http://localhost:3002`
- [ ] Console 顯示 `[API] Backend is available`
- [ ] Network 顯示 `/internal/v1/projects/.../templates` 請求狀態為 `200`
- [ ] UI 顯示 2 個 templates（與資料庫一致）

## 🔧 如果仍使用 Mock Data

1. **檢查 Vite Proxy 配置**:
   ```typescript
   // Frontend/vite.config.ts
   server: {
     proxy: {
       '/internal/v1': 'http://localhost:8080',
       '/api/v1': 'http://localhost:8080',
     }
   }
   ```

2. **檢查 API Service 配置**:
   ```typescript
   // Frontend/services/apiService.ts
   const API_INTERNAL = '/internal/v1'; // 使用相對路徑，由 Vite proxy 處理
   ```

3. **重啟 Frontend Server**:
   ```bash
   cd Frontend
   npm run dev
   ```
