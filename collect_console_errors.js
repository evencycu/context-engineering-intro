// Console 錯誤收集腳本
// 在瀏覽器 Console 中執行此腳本，會收集所有錯誤並顯示

(function() {
  console.log('🔍 開始收集 Console 錯誤...\n');
  
  const errors = [];
  const warnings = [];
  const apiErrors = [];
  
  // 攔截 console.error
  const originalError = console.error;
  console.error = function(...args) {
    const errorMsg = args.map(arg => {
      if (typeof arg === 'object') {
        try {
          return JSON.stringify(arg, null, 2);
        } catch {
          return String(arg);
        }
      }
      return String(arg);
    }).join(' ');
    
    errors.push({
      message: errorMsg,
      timestamp: new Date().toISOString(),
      stack: new Error().stack
    });
    
    originalError.apply(console, args);
  };
  
  // 攔截 console.warn
  const originalWarn = console.warn;
  console.warn = function(...args) {
    const warnMsg = args.map(arg => String(arg)).join(' ');
    
    // 檢查是否是 API 相關警告
    if (warnMsg.includes('API call failed') || warnMsg.includes('Falling back to mock')) {
      apiErrors.push({
        message: warnMsg,
        timestamp: new Date().toISOString()
      });
    }
    
    warnings.push({
      message: warnMsg,
      timestamp: new Date().toISOString()
    });
    
    originalWarn.apply(console, args);
  };
  
  // 攔截未捕獲的錯誤
  window.addEventListener('error', (event) => {
    errors.push({
      message: event.message,
      filename: event.filename,
      lineno: event.lineno,
      colno: event.colno,
      error: event.error?.stack,
      timestamp: new Date().toISOString()
    });
  });
  
  // 攔截 Promise rejection
  window.addEventListener('unhandledrejection', (event) => {
    errors.push({
      message: 'Unhandled Promise Rejection',
      reason: event.reason?.message || String(event.reason),
      stack: event.reason?.stack,
      timestamp: new Date().toISOString()
    });
  });
  
  // 測試 API 連接
  console.log('📡 測試 API 連接...');
  fetch('/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/templates')
    .then(r => r.json())
    .then(d => {
      console.log('✅ API 連接:', d.code === 200 ? '成功' : '失敗');
      console.log('📊 Templates 數量:', d.data?.length || 0);
    })
    .catch(e => {
      apiErrors.push({
        message: `API 連接失敗: ${e.message}`,
        timestamp: new Date().toISOString()
      });
      console.error('❌ API 連接失敗:', e);
    });
  
  // 等待 3 秒後顯示結果
  setTimeout(() => {
    console.log('\n📊 ========== 錯誤收集結果 ==========\n');
    
    if (errors.length === 0 && apiErrors.length === 0) {
      console.log('✅ 沒有發現錯誤！');
    } else {
      if (errors.length > 0) {
        console.log(`❌ 發現 ${errors.length} 個錯誤:\n`);
        errors.forEach((err, i) => {
          console.log(`${i + 1}. [${err.timestamp}]`);
          console.log(`   訊息: ${err.message}`);
          if (err.filename) {
            console.log(`   檔案: ${err.filename}:${err.lineno}:${err.colno}`);
          }
          if (err.error || err.stack) {
            console.log(`   堆疊: ${err.error || err.stack}`);
          }
          console.log('');
        });
      }
      
      if (apiErrors.length > 0) {
        console.log(`⚠️  發現 ${apiErrors.length} 個 API 相關錯誤:\n`);
        apiErrors.forEach((err, i) => {
          console.log(`${i + 1}. [${err.timestamp}] ${err.message}`);
        });
        console.log('');
      }
      
      if (warnings.length > 0) {
        console.log(`⚠️  發現 ${warnings.length} 個警告:\n`);
        warnings.slice(0, 10).forEach((warn, i) => {
          console.log(`${i + 1}. [${warn.timestamp}] ${warn.message}`);
        });
        if (warnings.length > 10) {
          console.log(`   ... 還有 ${warnings.length - 10} 個警告`);
        }
        console.log('');
      }
    }
    
    // 生成可複製的錯誤報告
    const report = {
      timestamp: new Date().toISOString(),
      url: window.location.href,
      errors: errors,
      apiErrors: apiErrors,
      warnings: warnings.slice(0, 20) // 只保留前 20 個警告
    };
    
    console.log('📋 完整錯誤報告（可複製）:');
    console.log(JSON.stringify(report, null, 2));
    
    // 將報告保存到全局變量，方便複製
    window.__consoleErrorReport = report;
    console.log('\n💡 提示: 錯誤報告已保存到 window.__consoleErrorReport，可以複製 JSON.stringify(window.__consoleErrorReport, null, 2)');
    
  }, 3000);
  
  console.log('⏳ 等待 3 秒收集錯誤...');
})();
