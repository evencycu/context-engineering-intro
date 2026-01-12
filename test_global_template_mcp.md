# Chrome MCP 自動化測試 Global Template CRUD

## 🎯 測試目標
驗證 `http://localhost:3000/#/admin/templates` 頁面的 Global Template CRUD 功能。

## 📋 Chrome MCP 測試命令

### Step 1: 打開頁面
```
打開 http://localhost:3000/#/admin/templates
```

### Step 2: 等待頁面載入
```
等待 2 秒
```

### Step 3: 檢查頁面元素
```
檢查頁面中是否包含文字 "Global Templates"
檢查頁面中是否包含按鈕 "New Global Template"
截圖當前頁面，保存為 global_template_page_loaded.png
```

### Step 4: 測試創建 Global Template
```
點擊包含文字 "New Global Template" 的按鈕
等待 1 秒
檢查頁面中是否包含文字 "Create New Global Template" 或 "Global Template"
截圖當前頁面，保存為 global_template_modal_opened.png
```

### Step 5: 填寫表單
```
在名為 "Template Name" 的輸入框中輸入 "Test Global Template"
在名為 "Description" 的文字區域中輸入 "Automated test template"
在變數輸入框中輸入 "title" 並點擊 "Add" 按鈕
在變數輸入框中輸入 "message" 並點擊 "Add" 按鈕
在包含 class "font-mono" 或 "bg-slate" 的文字區域中輸入以下 JSON:
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
截圖當前頁面，保存為 global_template_form_filled.png
```

### Step 6: 儲存 Template
```
點擊包含文字 "Save Template" 的按鈕
等待 2 秒
檢查頁面中是否包含文字 "Test Global Template"
截圖當前頁面，保存為 global_template_created.png
```

### Step 7: 測試編輯
```
點擊包含文字 "Test Global Template" 的行中的 Edit 圖標（鉛筆圖標）
等待 1 秒
檢查頁面中是否包含文字 "Edit Global Template"
修改 Template Name 為 "Updated Test Global Template"
點擊包含文字 "Save Template" 的按鈕
等待 2 秒
檢查頁面中是否包含文字 "Updated Test Global Template"
截圖當前頁面，保存為 global_template_updated.png
```

### Step 8: 測試刪除
```
點擊包含文字 "Updated Test Global Template" 的行中的 Delete 圖標（垃圾桶圖標）
確認刪除對話框
等待 2 秒
檢查頁面中是否不再包含文字 "Updated Test Global Template"
截圖當前頁面，保存為 global_template_deleted.png
```

### Step 9: 檢查 Console
```
檢查瀏覽器 Console 是否有錯誤訊息
```

### Step 10: 檢查 Network
```
檢查 Network 標籤中是否有對應的 API 請求
```

## ✅ 預期結果

### 成功指標：
- ✅ 頁面正常載入
- ✅ "New Global Template" 按鈕可以點擊
- ✅ Modal 正常打開
- ✅ 表單可以正常填寫
- ✅ 可以成功儲存 Template
- ✅ 新 Template 出現在列表中
- ✅ 可以成功編輯 Template
- ✅ 可以成功刪除 Template
- ✅ Console 中沒有錯誤訊息

### 失敗指標：
- ❌ 頁面無法載入
- ❌ 按鈕無法點擊
- ❌ Modal 無法打開
- ❌ 表單無法填寫
- ❌ 儲存失敗
- ❌ Console 中有錯誤訊息

## 🐛 除錯

如果測試失敗，請檢查：

1. **Frontend Server 狀態**
   ```bash
   curl http://localhost:3000
   ```

2. **Console 錯誤**
   - 打開開發者工具
   - 查看 Console 標籤
   - 複製錯誤訊息

3. **Network 請求**
   - 打開開發者工具
   - 查看 Network 標籤
   - 檢查 API 請求的狀態碼和回應內容

4. **代碼檢查**
   - 確認 `Admin.tsx` 中的 handler 函數都已正確實現
   - 確認 `services/index.ts` 中的 API 方法都已添加
