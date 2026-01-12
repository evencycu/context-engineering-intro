# 🧪 Global Template CRUD 快速測試指南

## ✅ 準備狀態檢查

- ✅ Frontend Server: `http://localhost:3000` (運行中)
- ✅ Backend Server: `http://localhost:8080` (運行中)
- ✅ Global Template 功能: 已修復並可用

## 🚀 快速測試步驟

### 1. 打開頁面
訪問：`http://localhost:3000/#/admin/templates`

### 2. 測試創建 (30秒)
1. 點擊 **"New Global Template"** 按鈕
2. 填寫：
   - Name: `Test Template`
   - Description: `Quick test`
   - Variables: 輸入 `title` 點擊 Add
   - JSON: 
     ```json
     {
       "type": "AdaptiveCard",
       "version": "1.4",
       "body": [{"type": "TextBlock", "text": "{{title}}"}]
     }
     ```
3. 點擊 **"Save Template"**
4. ✅ 應該看到 "Test Template" 出現在列表中

### 3. 測試編輯 (20秒)
1. 點擊 "Test Template" 行的 **Edit 圖標** (鉛筆)
2. 修改 Name 為 `Updated Test Template`
3. 點擊 **"Save Template"**
4. ✅ 應該看到名稱更新

### 4. 測試刪除 (10秒)
1. 點擊 **Delete 圖標** (垃圾桶)
2. 確認刪除
3. ✅ Template 應該從列表中移除

## ✅ 測試通過標準

- [x] 可以打開頁面
- [x] 可以點擊 "New Global Template"
- [x] Modal 正常打開
- [x] 可以填寫表單
- [x] 可以成功儲存
- [x] Template 出現在列表
- [x] 可以編輯
- [x] 可以刪除

## 🐛 如果遇到問題

1. **按鈕沒反應**
   - 檢查 Console (F12) 是否有錯誤
   - 確認頁面已完全載入

2. **Modal 沒打開**
   - 檢查 `handleOpenCreateGlobalTemplate` 是否被調用
   - 檢查 `showGlobalTemplateModal` 狀態

3. **儲存失敗**
   - 檢查 Console 錯誤訊息
   - 檢查 Network 標籤中的請求
   - 確認 JSON 格式正確

4. **列表沒更新**
   - 檢查 `getGlobalTemplates` 是否被調用
   - 檢查 mock API 是否正常運作

## 📝 測試結果記錄

測試時間: _______________

結果: [ ] 全部通過  [ ] 部分通過  [ ] 失敗

備註: 
_________________________________
