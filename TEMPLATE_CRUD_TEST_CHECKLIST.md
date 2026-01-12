# Template CRUD 功能測試檢查清單

## 🎯 測試目標
驗證前端頁面 `http://localhost:3000/#/admin/templates` 的 Template CRUD 功能是否正常運作。

## 📋 前置條件檢查

### 1. 服務狀態
- [x] Frontend 運行中 (`http://localhost:3000`)
- [x] Backend 運行中 (`http://localhost:8080`)
- [ ] 資料庫中有至少一個 Project（用於關聯 Template）

### 2. 瀏覽器準備
- [ ] 打開 Chrome 開發者工具（F12）
- [ ] 切換到 **Console** 標籤（查看錯誤訊息）
- [ ] 切換到 **Network** 標籤（查看 API 請求）

## 🧪 測試步驟

### Test 1: 載入 Templates 列表
1. 訪問 `http://localhost:3000/#/admin/templates`
2. 檢查頁面是否正常載入
3. 檢查是否顯示 "No templates found" 或現有的 templates 列表
4. **預期結果**: 頁面正常顯示，沒有錯誤訊息

**檢查點**:
- [ ] 頁面正常載入
- [ ] 沒有 Console 錯誤
- [ ] Network 標籤中有 `GET /internal/v1/projects/{projectId}/templates` 請求
- [ ] API 請求返回 200 狀態碼

---

### Test 2: 新增 Template
1. 點擊 **"Create Template"** 按鈕
2. 填寫表單：
   - **Template Name**: `Test Template`
   - **Description**: `This is a test template`
   - **Variables**: 點擊 "Add" 添加變數 `title` 和 `message`
   - **JSON Structure**: 
     ```json
     {
       "type": "AdaptiveCard",
       "version": "1.4",
       "body": [
         {
           "type": "TextBlock",
           "text": "{{title}}"
         },
         {
           "type": "TextBlock",
           "text": "{{message}}"
         }
       ]
     }
     ```
3. 點擊 **"Save Template"** 按鈕
4. **預期結果**: Template 成功創建，modal 關閉，列表更新

**檢查點**:
- [ ] Modal 正常打開
- [ ] 表單可以正常輸入
- [ ] JSON 驗證通過（如果輸入無效 JSON 會顯示錯誤）
- [ ] 點擊 Save 後沒有錯誤訊息
- [ ] Network 標籤中有 `POST /internal/v1/projects/{projectId}/templates` 請求
- [ ] API 請求返回 201 或 200 狀態碼
- [ ] 新創建的 Template 出現在列表中

**可能的錯誤**:
- ❌ `Project ID is required` - 檢查 `currentProject.id` 是否存在
- ❌ `Failed to create template: foreign key constraint` - Project ID 不存在於資料庫
- ❌ `Invalid JSON structure` - JSON 格式錯誤

---

### Test 3: 編輯 Template
1. 在 Template 列表中，點擊某個 Template 的 **編輯圖標**（鉛筆圖標）
2. 修改表單內容：
   - 修改 **Template Name** 為 `Updated Test Template`
   - 添加新的變數 `severity`
   - 修改 JSON Structure
3. 點擊 **"Save Template"** 按鈕
4. **預期結果**: Template 成功更新，modal 關閉，列表顯示更新後的內容

**檢查點**:
- [ ] Modal 正常打開並顯示現有資料
- [ ] 所有欄位都正確預填（name, description, variables, JSON）
- [ ] 修改後可以正常儲存
- [ ] Network 標籤中有 `PUT /internal/v1/projects/{projectId}/templates/{templateId}` 請求
- [ ] API 請求返回 200 狀態碼
- [ ] 列表中的 Template 顯示更新後的內容

**可能的錯誤**:
- ❌ `Project ID is required` - 檢查傳遞的 `projectId` 參數
- ❌ `Template not found` - Template ID 不存在

---

### Test 4: 刪除 Template
1. 在 Template 列表中，點擊某個 Template 的 **刪除圖標**（垃圾桶圖標）
2. 確認刪除對話框
3. 點擊確認
4. **預期結果**: Template 成功刪除，從列表中移除

**檢查點**:
- [ ] 確認對話框正常顯示
- [ ] 點擊確認後沒有錯誤訊息
- [ ] Network 標籤中有 `DELETE /internal/v1/projects/{projectId}/templates/{templateId}` 請求
- [ ] API 請求返回 204 或 200 狀態碼
- [ ] Template 從列表中移除

**可能的錯誤**:
- ❌ `Project ID is required` - 檢查傳遞的 `projectId` 參數
- ❌ `Template not found` - Template ID 不存在

---

### Test 5: JSON 驗證
1. 嘗試創建一個 Template，但輸入無效的 JSON：
   ```json
   {
     "type": "AdaptiveCard",
     "version": "1.4",
     "body": [
       {
         "type": "TextBlock",
         "text": "{{title}}"
       }
     ]
     // 缺少閉合括號
   ```
2. 點擊 **"Save Template"**
3. **預期結果**: 顯示錯誤訊息 "Invalid JSON structure. Please check your JSON syntax."

**檢查點**:
- [ ] 無效 JSON 被正確攔截
- [ ] 顯示明確的錯誤訊息
- [ ] 不會發送 API 請求

---

## 🔍 除錯指南

### 如果所有操作都失敗：

1. **檢查 Console 錯誤**:
   - 打開開發者工具 Console
   - 查看是否有紅色錯誤訊息
   - 複製錯誤訊息內容

2. **檢查 Network 請求**:
   - 打開開發者工具 Network 標籤
   - 嘗試操作（新增/修改/刪除）
   - 查看對應的 API 請求：
     - 請求 URL 是否正確
     - 請求方法（GET/POST/PUT/DELETE）是否正確
     - 請求狀態碼（200/201/204/400/500）
     - 請求 Body 格式是否正確（camelCase）
     - 回應內容是什麼

3. **檢查 Project Context**:
   - 確認頁面右上角的 Project Switcher 有選擇 Project
   - 如果沒有 Project，需要先到 Admin > Organization 創建 Project

4. **檢查 Backend 日誌**:
   ```bash
   # 查看 backend 日誌
   tail -f /tmp/teamsnotify_server.log
   # 或如果使用 docker
   docker logs -f teamsnotify-backend
   ```

### 常見問題解決：

**問題**: `foreign key constraint` 錯誤
- **原因**: Project ID 不存在於資料庫
- **解決**: 先創建一個 Project，或使用現有的 Project ID

**問題**: `Project ID is required` 錯誤
- **原因**: `currentProject.id` 為空
- **解決**: 確保在 Project Switcher 中選擇了 Project

**問題**: API 返回 404
- **原因**: API 路徑不正確
- **解決**: 檢查 API 路徑是否為 `/internal/v1/projects/{projectId}/templates`

**問題**: API 返回 500
- **原因**: Backend 內部錯誤
- **解決**: 查看 Backend 日誌找出具體錯誤

---

## ✅ 測試完成標準

所有測試通過的標準：
- [x] 可以正常載入 Templates 列表
- [x] 可以成功新增 Template
- [x] 可以成功編輯 Template
- [x] 可以成功刪除 Template
- [x] JSON 驗證正常工作
- [x] 錯誤訊息清晰明確
- [x] 沒有 Console 錯誤
- [x] 所有 API 請求返回正確狀態碼

---

## 📝 測試結果記錄

**測試日期**: _______________

**測試人員**: _______________

**測試結果**:
- [ ] 全部通過
- [ ] 部分通過（請記錄失敗的測試項目）
- [ ] 全部失敗（請記錄錯誤訊息）

**發現的問題**:
1. 
2. 
3. 

**建議改進**:
1. 
2. 
3. 
