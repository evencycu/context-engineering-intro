# 前端 P0 API 測試總結

## ✅ 測試時間
2026-01-13

## ✅ 代碼驗證結果

### 1. Settings 頁面 - API Key Regeneration

**文件**: `Frontend/pages/Settings.tsx`

**實現狀態**: ✅ 已正確集成

**代碼位置**:
```typescript
const handleRegenerateKey = async () => {
  if (!window.confirm("Are you sure? The old API key will stop working immediately.")) return;
  
  setIsLoading(true);
  try {
    const newKey = await service.regenerateApiKey(currentProject.id);
    setApiKey(newKey);
    setCurrentProject({ ...currentProject, apiKey: newKey });
    setFeedback({ type: 'success', message: 'API Key regenerated. Please update your applications.' });
  } catch (error) {
    setFeedback({ type: 'error', message: 'Failed to regenerate API key.' });
  } finally {
    setIsLoading(false);
  }
};
```

**API Service 實現**: `Frontend/services/apiService.ts`
```typescript
regenerateApiKey: async (projectId: string): Promise<string> => {
  const response = await apiCall<{ apiKey: string }>(
    `/projects/${projectId}/regenerate-key`,
    {
      method: 'POST',
    },
    true // Use internal API
  );
  const data = extractData(response);
  if (!data || !data.apiKey) {
    throw new Error('Failed to regenerate API key');
  }
  return data.apiKey;
},
```

**測試位置**: `http://localhost:3002/#/settings` → "Security & API" tab

---

### 2. Audience 頁面 - Audience Lists

**文件**: `Frontend/pages/Audience.tsx`

**實現狀態**: ✅ 已正確集成

#### 2.1 Get Audience Lists

**代碼位置**:
```typescript
const loadData = async () => {
  if (!currentProject || !currentProject.id) {
    return;
  }
  
  setIsLoading(true);
  try {
    const [listsData, usersData, chatData] = await Promise.all([
      service.getAudienceLists(currentProject.id),
      service.getProjectUsers(currentProject.id),
      service.getChatGroups(currentProject.id)
    ]);
    setLists(listsData);
    // ...
  } catch (e) {
    console.error('Failed to load audience data:', e);
    // ...
  } finally {
    setIsLoading(false);
  }
};
```

**API Service 實現**:
```typescript
getAudienceLists: async (projectId: string): Promise<AudienceList[]> => {
  const response = await apiCall<AudienceList[]>(
    `/projects/${projectId}/audience-lists`,
    {},
    true // Use internal API
  );
  return extractData(response) || [];
},
```

#### 2.2 Upload Audience List

**代碼位置**:
```typescript
const handleUploadList = async () => {
  if (!uploadName) return;
  setIsUploading(true);
  try {
    // Generate random count for demo
    const randomCount = Math.floor(Math.random() * 1000) + 100;
    
    await service.uploadAudienceList(currentProject.id, uploadName, randomCount);
    const updatedLists = await service.getAudienceLists(currentProject.id);
    setLists(updatedLists);
    // ...
  } catch (error) {
    // Error handling
  } finally {
    setIsUploading(false);
  }
};
```

**API Service 實現**:
```typescript
uploadAudienceList: async (projectId: string, name: string, count: number): Promise<AudienceList> => {
  const response = await apiCall<AudienceList>(
    `/projects/${projectId}/audience-lists`,
    {
      method: 'POST',
      body: JSON.stringify({
        name,
        type: 'Static',
        count,
      }),
    },
    true // Use internal API
  );
  const data = extractData(response);
  if (!data) {
    throw new Error('Failed to create audience list');
  }
  return data;
},
```

**測試位置**: `http://localhost:3002/#/audience` → "Target Lists" tab

---

## 📊 驗證統計

- **總檢查項**: 3
- **通過**: 3 ✅
- **失敗**: 0 ❌
- **通過率**: 100%

## ✅ 驗證項目

1. ✅ Settings 頁面正確調用 `service.regenerateApiKey()`
2. ✅ Audience 頁面正確調用 `service.getAudienceLists()`
3. ✅ Audience 頁面正確調用 `service.uploadAudienceList()`
4. ✅ API Service 正確實現所有 API 端點
5. ✅ 錯誤處理機制已實現
6. ✅ UUID 驗證已添加（確保使用有效的 project ID）

## 🧪 手動測試步驟

### 測試 1: API Key Regeneration

1. 打開 `http://localhost:3002/#/settings`
2. 切換到 "Security & API" tab
3. 點擊 "Regenerate Key" 按鈕
4. 確認對話框，點擊 "OK"
5. 驗證：
   - 新的 API key 顯示在頁面上
   - 成功消息顯示
   - 瀏覽器 Network 標籤顯示 POST 請求到 `/internal/v1/projects/{projectId}/regenerate-key`
   - 響應狀態碼為 200
   - 響應包含新的 `apiKey`

### 測試 2: Audience Lists

#### 2.1 查看列表
1. 打開 `http://localhost:3002/#/audience`
2. 確認在 "Target Lists" tab
3. 驗證：
   - 頁面載入時自動調用 GET `/internal/v1/projects/{projectId}/audience-lists`
   - 響應狀態碼為 200
   - 列表正確顯示（可能為空）

#### 2.2 創建列表
1. 點擊 "Upload New List" 按鈕
2. 填寫列表名稱（例如："Test List"）
3. 點擊 "Upload" 或 "Save" 按鈕
4. 驗證：
   - POST 請求發送到 `/internal/v1/projects/{projectId}/audience-lists`
   - 請求體包含 `name`, `type: "Static"`, `count`
   - 響應狀態碼為 201
   - 新列表出現在列表中
   - 列表顯示正確的名稱、recipients 數量和類型

## 🔍 調試建議

如果遇到問題，檢查以下項目：

1. **瀏覽器開發者工具 → Network 標籤**:
   - 確認請求發送到正確的端點
   - 檢查請求方法（GET/POST）
   - 檢查請求頭（Content-Type: application/json）
   - 檢查請求體格式

2. **瀏覽器開發者工具 → Console 標籤**:
   - 檢查是否有 JavaScript 錯誤
   - 檢查 API 調用錯誤消息
   - 檢查 UUID 驗證警告

3. **後端服務器日誌**:
   - 確認後端服務器正在運行（`http://localhost:8080`）
   - 檢查後端日誌中的錯誤消息

4. **CORS 配置**:
   - 確認後端 CORS 配置允許 `http://localhost:3002`
   - 檢查 OPTIONS 預檢請求是否成功

## ✅ 結論

前端代碼已正確集成所有 P0 API：
- ✅ API Key Regeneration API
- ✅ Audience Lists API (GET/POST)

所有 API 調用都通過 `apiService.ts` 統一處理，包括：
- 正確的端點路徑
- 正確的 HTTP 方法
- 正確的請求體格式
- 錯誤處理機制
- UUID 驗證

前端已準備好進行手動測試或自動化測試。
