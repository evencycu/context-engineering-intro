// Templates 頁面診斷腳本
// 在瀏覽器 Console 中執行此腳本

(function() {
  console.log('🔍 開始診斷 Templates 頁面...\n');
  
  const diagnostics = {
    errors: [],
    warnings: [],
    apiStatus: null,
    pageElements: {},
    timestamp: new Date().toISOString()
  };
  
  // 1. 收集現有錯誤
  console.log('1️⃣ 收集 Console 錯誤...');
  
  // 攔截新的錯誤
  const originalError = console.error;
  const originalWarn = console.warn;
  
  console.error = function(...args) {
    const msg = args.map(a => String(a)).join(' ');
    diagnostics.errors.push({
      message: msg,
      timestamp: new Date().toISOString()
    });
    originalError.apply(console, args);
  };
  
  console.warn = function(...args) {
    const msg = args.map(a => String(a)).join(' ');
    if (msg.includes('API') || msg.includes('mock') || msg.includes('fallback')) {
      diagnostics.warnings.push({
        message: msg,
        timestamp: new Date().toISOString()
      });
    }
    originalWarn.apply(console, args);
  };
  
  // 2. 檢查頁面元素
  console.log('2️⃣ 檢查頁面元素...');
  
  const title = document.querySelector('h2');
  diagnostics.pageElements.title = title ? title.textContent : '未找到';
  
  const createButton = Array.from(document.querySelectorAll('button')).find(btn => 
    btn.textContent.includes('Create') || btn.textContent.includes('New')
  );
  diagnostics.pageElements.createButton = createButton ? '找到' : '未找到';
  
  const templates = document.querySelectorAll('[class*="template"], [data-template]');
  diagnostics.pageElements.templateCount = templates.length;
  
  // 3. 檢查 API 連接
  console.log('3️⃣ 測試 API 連接...');
  
  // 獲取當前 project ID（從 React state 或 URL）
  let projectId = null;
  
  // 嘗試從 window 對象獲取（如果 React DevTools 可用）
  if (window.__REACT_DEVTOOLS_GLOBAL_HOOK__) {
    console.log('   React DevTools 可用，嘗試獲取 project ID...');
  }
  
  // 嘗試從 localStorage 或 sessionStorage
  try {
    const projectData = localStorage.getItem('currentProject') || sessionStorage.getItem('currentProject');
    if (projectData) {
      const project = JSON.parse(projectData);
      projectId = project.id;
    }
  } catch (e) {
    // 忽略
  }
  
  // 使用默認 project ID
  if (!projectId) {
    projectId = '750e8400-e29b-41d4-a716-446655440001';
    console.log(`   使用默認 project ID: ${projectId}`);
  }
  
  // 測試 Templates API
  fetch(`/internal/v1/projects/${projectId}/templates`)
    .then(async (r) => {
      const data = await r.json();
      diagnostics.apiStatus = {
        url: `/internal/v1/projects/${projectId}/templates`,
        status: r.status,
        statusText: r.statusText,
        success: r.ok,
        data: data
      };
      
      if (r.ok && data.code === 200) {
        console.log(`   ✅ API 連接成功`);
        console.log(`   📊 Templates 數量: ${data.data?.length || 0}`);
        if (data.data && data.data.length > 0) {
          console.log(`   📋 Templates 列表:`);
          data.data.forEach((t, i) => {
            console.log(`      ${i + 1}. ${t.name} (ID: ${t.id})`);
          });
        }
      } else {
        console.error(`   ❌ API 返回錯誤: ${data.message || data.error || r.statusText}`);
        diagnostics.errors.push({
          message: `API 錯誤: ${data.message || data.error || r.statusText}`,
          status: r.status
        });
      }
    })
    .catch((e) => {
      console.error(`   ❌ API 連接失敗: ${e.message}`);
      diagnostics.errors.push({
        message: `API 連接失敗: ${e.message}`,
        type: 'network'
      });
      diagnostics.apiStatus = {
        error: e.message,
        success: false
      };
    })
    .finally(() => {
      // 4. 檢查 React 組件狀態
      console.log('4️⃣ 檢查 React 組件...');
      
      // 嘗試檢查 React 組件
      const root = document.getElementById('root');
      if (root && root._reactInternalInstance) {
        console.log('   React 組件已掛載');
      }
      
      // 5. 檢查 Network 請求
      console.log('5️⃣ 檢查 Network 請求...');
      console.log('   請在 Network 標籤中檢查：');
      console.log('   - /internal/v1/projects/.../templates 請求');
      console.log('   - 狀態碼應該是 200');
      
      // 6. 顯示診斷結果
      setTimeout(() => {
        console.log('\n📊 ========== 診斷結果 ==========\n');
        
        console.log('📄 頁面元素:');
        console.log(`   - 標題: ${diagnostics.pageElements.title}`);
        console.log(`   - Create 按鈕: ${diagnostics.pageElements.createButton}`);
        console.log(`   - Template 元素數量: ${diagnostics.pageElements.templateCount}`);
        console.log('');
        
        if (diagnostics.apiStatus) {
          console.log('🌐 API 狀態:');
          if (diagnostics.apiStatus.success) {
            console.log(`   ✅ 連接成功`);
            console.log(`   📊 Templates: ${diagnostics.apiStatus.data?.data?.length || 0} 個`);
          } else {
            console.log(`   ❌ 連接失敗`);
            console.log(`   錯誤: ${diagnostics.apiStatus.error || diagnostics.apiStatus.data?.message}`);
          }
          console.log('');
        }
        
        if (diagnostics.errors.length > 0) {
          console.log(`❌ 發現 ${diagnostics.errors.length} 個錯誤:\n`);
          diagnostics.errors.forEach((err, i) => {
            console.log(`${i + 1}. ${err.message}`);
            if (err.status) console.log(`   狀態碼: ${err.status}`);
            if (err.type) console.log(`   類型: ${err.type}`);
          });
          console.log('');
        } else {
          console.log('✅ 沒有發現錯誤');
          console.log('');
        }
        
        if (diagnostics.warnings.length > 0) {
          console.log(`⚠️  發現 ${diagnostics.warnings.length} 個警告:\n`);
          diagnostics.warnings.forEach((warn, i) => {
            console.log(`${i + 1}. ${warn.message}`);
          });
          console.log('');
        }
        
        // 生成完整報告
        const report = {
          ...diagnostics,
          url: window.location.href,
          userAgent: navigator.userAgent
        };
        
        console.log('📋 完整診斷報告:');
        console.log(JSON.stringify(report, null, 2));
        
        // 保存到全局變量
        window.__templatesPageDiagnostics = report;
        console.log('\n💡 提示: 診斷報告已保存到 window.__templatesPageDiagnostics');
        console.log('   可以複製: JSON.stringify(window.__templatesPageDiagnostics, null, 2)');
        
      }, 2000);
    });
  
  console.log('⏳ 等待診斷完成...');
})();
