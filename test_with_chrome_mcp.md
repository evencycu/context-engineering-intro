# 使用 Chrome MCP 自動化測試 Template CRUD

## 🎯 測試目標
使用 Chrome MCP 自動化測試 Template CRUD 功能。

## 📋 測試步驟

### Step 1: 打開 Templates 頁面

請在 Chrome MCP 中執行：
```
打開 http://localhost:3000/#/admin/templates
```

### Step 2: 檢查頁面載入

檢查以下元素是否存在：
- 標題 "Message Templates"
- "Create Template" 按鈕
- Project Switcher（右上角）

### Step 3: 執行自動化測試腳本

在瀏覽器 Console 中執行以下代碼（或使用 Chrome MCP 執行）：

```javascript
// 複製 test_template_crud_automated.js 的內容到 Console
// 然後執行：
runTemplateCRUDTests()
```

### Step 4: 手動測試流程

如果自動化腳本無法執行，請按照以下步驟手動測試：

#### 4.1 測試新增 Template

1. **點擊 "Create Template" 按鈕**
   - 使用 Chrome MCP: `點擊包含文字 "Create Template" 的按鈕`

2. **填寫表單**
   - Template Name: `Test Template`
   - Description: `Automated test template`
   - Variables: 添加 `title` 和 `message`
   - JSON Structure:
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

3. **點擊 "Save Template" 按鈕**
   - 使用 Chrome MCP: `點擊包含文字 "Save Template" 的按鈕`

4. **檢查結果**
   - Modal 應該關閉
   - 新 Template 應該出現在列表中
   - Network 標籤中應該有 `POST /internal/v1/projects/{projectId}/templates` 請求

#### 4.2 測試編輯 Template

1. **點擊 Template 卡片的編輯圖標**
   - 使用 Chrome MCP: `點擊包含 Edit2 圖標的按鈕`

2. **修改表單內容**
   - 修改 Template Name 為 `Updated Test Template`
   - 添加新變數 `severity`

3. **點擊 "Save Template" 按鈕**

4. **檢查結果**
   - Template 名稱應該更新
   - Network 標籤中應該有 `PUT /internal/v1/projects/{projectId}/templates/{templateId}` 請求

#### 4.3 測試刪除 Template

1. **點擊 Template 卡片的刪除圖標**
   - 使用 Chrome MCP: `點擊包含 Trash2 圖標的按鈕`

2. **確認刪除**
   - 確認對話框應該出現
   - 點擊確認

3. **檢查結果**
   - Template 應該從列表中移除
   - Network 標籤中應該有 `DELETE /internal/v1/projects/{projectId}/templates/{templateId}` 請求

## 🔍 Chrome MCP 命令範例

### 打開頁面
```
打開 http://localhost:3000/#/admin/templates
```

### 點擊按鈕
```
點擊包含文字 "Create Template" 的按鈕
```

### 填寫表單
```
在名為 "Template Name" 的輸入框中輸入 "Test Template"
在名為 "Description" 的文字區域中輸入 "Test description"
```

### 檢查元素
```
檢查頁面中是否包含文字 "Message Templates"
檢查頁面中是否包含按鈕 "Create Template"
```

### 截圖
```
截圖當前頁面
```

### 檢查 Console
```
檢查瀏覽器 Console 是否有錯誤訊息
```

### 檢查 Network
```
檢查 Network 標籤中是否有對 /templates 的 API 請求
```

## 📊 預期結果

### 成功指標：
- ✅ 頁面正常載入，沒有錯誤
- ✅ 可以成功新增 Template
- ✅ 可以成功編輯 Template
- ✅ 可以成功刪除 Template
- ✅ 所有 API 請求返回 200/201/204 狀態碼
- ✅ Console 中沒有錯誤訊息

### 失敗指標：
- ❌ 頁面無法載入
- ❌ 按鈕無法點擊
- ❌ API 請求返回錯誤狀態碼（400/500）
- ❌ Console 中有錯誤訊息
- ❌ Modal 無法打開或關閉

## 🐛 除錯

如果測試失敗，請檢查：

1. **服務狀態**
   ```bash
   curl http://localhost:3000
   curl http://localhost:8080/health
   ```

2. **Console 錯誤**
   - 打開開發者工具
   - 查看 Console 標籤
   - 複製錯誤訊息

3. **Network 請求**
   - 打開開發者工具
   - 查看 Network 標籤
   - 檢查 API 請求的狀態碼和回應內容

4. **Project Context**
   - 確保頁面右上角有選擇 Project
   - 如果沒有，先創建一個 Project

## 📝 測試報告模板

```
測試日期: _______________
測試人員: _______________

測試結果:
- [ ] 全部通過
- [ ] 部分通過
- [ ] 全部失敗

通過的測試:
1. 
2. 

失敗的測試:
1. 
2. 

錯誤訊息:
1. 
2. 

截圖: (附上截圖)
```
