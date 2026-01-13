# 檢查 Console 錯誤指南

## 🔍 檢查步驟

### 1. 打開開發者工具
- 按 `F12` 或 `Cmd+Option+I` (Mac) / `Ctrl+Shift+I` (Windows)
- 或右鍵點擊頁面 → "檢查" / "Inspect"

### 2. 檢查 Console 標籤

#### 常見錯誤類型：

**A. Chrome Extension 相關錯誤**
```
Uncaught TypeError: Cannot read property 'xxx' of undefined
chrome-extension://...
```

**B. API 調用錯誤**
```
[createTemplate] API call failed: ...
Failed to fetch
CORS error
```

**C. React 錯誤**
```
Uncaught Error: Minified React error #...
Cannot read property 'xxx' of null
```

**D. 網路錯誤**
```
Failed to load resource: net::ERR_CONNECTION_REFUSED
Failed to load resource: 404 (Not Found)
```

### 3. 檢查 Network 標籤

1. 切換到 **Network** 標籤
2. 重新載入頁面 (`Cmd+R` 或 `Ctrl+R`)
3. 查找失敗的請求（紅色標記）
4. 點擊失敗的請求，查看：
   - **Status**: 錯誤狀態碼
   - **Response**: 錯誤訊息
   - **Headers**: 請求標頭

### 4. 檢查 Chrome Extension 衝突

#### 步驟：
1. 打開 Chrome 擴展管理頁面：`chrome://extensions/`
2. 暫時禁用所有擴展
3. 重新載入頁面
4. 檢查錯誤是否消失

#### 常見衝突擴展：
- Ad blockers
- Privacy extensions
- Developer tools extensions
- React DevTools (可能導致錯誤)

### 5. 檢查 CORS 錯誤

如果看到 CORS 錯誤：
```
Access to fetch at 'http://localhost:8080/...' from origin 'http://localhost:3002' has been blocked by CORS policy
```

**解決方案**：
- 確認 Backend CORS 配置正確
- 確認 Vite proxy 配置正確

## 🐛 常見錯誤及解決方案

### 錯誤 1: Chrome Extension 注入腳本錯誤
**症狀**:
```
Uncaught TypeError: Cannot read property 'xxx' of undefined
chrome-extension://xxx/script.js:1
```

**解決方案**:
1. 禁用相關擴展
2. 或使用無痕模式測試

### 錯誤 2: API 調用失敗
**症狀**:
```
[createTemplate] API call failed: Failed to fetch
```

**解決方案**:
1. 確認 Backend server 運行在 `http://localhost:8080`
2. 檢查 Network 標籤查看具體錯誤
3. 確認 Vite proxy 配置

### 錯誤 3: React 渲染錯誤
**症狀**:
```
Uncaught Error: Minified React error #...
```

**解決方案**:
1. 檢查 React 組件是否有 null/undefined 錯誤
2. 檢查 props 是否正確傳遞
3. 查看完整錯誤堆疊

## 📋 診斷檢查清單

- [ ] 打開開發者工具 Console 標籤
- [ ] 記錄所有錯誤訊息（紅色）
- [ ] 記錄所有警告訊息（黃色）
- [ ] 檢查 Network 標籤中的失敗請求
- [ ] 檢查是否有 Chrome extension 相關錯誤
- [ ] 暫時禁用所有擴展並重新測試
- [ ] 檢查 CORS 錯誤
- [ ] 檢查 API 連接狀態

## 🔧 快速診斷腳本

在 Console 中執行以下代碼：

```javascript
// 檢查錯誤
console.log('=== Console 錯誤檢查 ===');
const errors = [];
const originalError = console.error;
console.error = function(...args) {
  errors.push(args.join(' '));
  originalError.apply(console, args);
};

// 檢查 API 連接
fetch('/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/templates')
  .then(r => r.json())
  .then(d => {
    console.log('✅ API 連接正常:', d.code === 200 ? '成功' : '失敗');
    console.log('📊 Templates 數量:', d.data?.length || 0);
  })
  .catch(e => {
    console.error('❌ API 連接失敗:', e);
  });

// 等待 2 秒後顯示錯誤
setTimeout(() => {
  console.log('\n=== 錯誤摘要 ===');
  if (errors.length === 0) {
    console.log('✅ 沒有發現錯誤');
  } else {
    console.log(`❌ 發現 ${errors.length} 個錯誤:`);
    errors.forEach((err, i) => {
      console.log(`${i + 1}. ${err}`);
    });
  }
}, 2000);
```
