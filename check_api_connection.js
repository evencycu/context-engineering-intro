// API 連接檢查腳本
// 在瀏覽器 Console 中執行此腳本

console.log('🔍 開始檢查 API 連接...\n');

// 1. 檢查 Backend 健康狀態
console.log('1️⃣ 檢查 Backend 健康狀態...');
fetch('http://localhost:8080/health')
  .then(r => r.json())
  .then(d => {
    console.log('✅ Backend 健康:', d);
    console.log('');
  })
  .catch(e => {
    console.error('❌ Backend 無法連接:', e);
    console.error('   請確認 Backend server 是否運行在 http://localhost:8080');
    console.log('');
  });

// 2. 檢查 Project Templates API
console.log('2️⃣ 檢查 Project Templates API...');
fetch('http://localhost:8080/internal/v1/projects/750e8400-e29b-41d4-a716-446655440001/templates')
  .then(r => r.json())
  .then(d => {
    if (d.code === 200) {
      console.log('✅ Project Templates API 成功');
      console.log('📊 Templates 數量:', d.data?.length || 0);
      if (d.data && d.data.length > 0) {
        console.log('📋 Templates 列表:');
        d.data.forEach((t, i) => {
          console.log(`   ${i + 1}. ${t.name} (ID: ${t.id})`);
        });
      }
    } else {
      console.error('❌ API 返回錯誤:', d);
    }
    console.log('');
  })
  .catch(e => {
    console.error('❌ Project Templates API 失敗:', e);
    console.log('');
  });

// 3. 檢查 Global Templates API
setTimeout(() => {
  console.log('3️⃣ 檢查 Global Templates API...');
  fetch('http://localhost:8080/internal/v1/global/templates')
    .then(r => {
      if (r.status === 404) {
        console.warn('⚠️  Global Templates API 返回 404');
        console.warn('   可能原因: Backend server 未重啟，新路由未載入');
        console.warn('   解決方案: 重啟 Backend server (cd Backend && make run)');
        return r.json();
      }
      return r.json();
    })
    .then(d => {
      if (d.code === 200) {
        console.log('✅ Global Templates API 成功');
        console.log('📊 Global Templates 數量:', d.data?.length || 0);
      } else {
        console.error('❌ Global Templates API 錯誤:', d);
      }
      console.log('');
    })
    .catch(e => {
      console.error('❌ Global Templates API 失敗:', e);
      console.log('');
    });
}, 500);

// 4. 檢查 Frontend Service 配置
setTimeout(() => {
  console.log('4️⃣ 檢查 Frontend Service 配置...');
  console.log('   請檢查 Frontend/services/index.ts:');
  console.log('   - USE_MOCK_ONLY:', false);
  console.log('   - ENABLE_FALLBACK:', true);
  console.log('   如果 API 失敗，會自動 fallback 到 mock data');
  console.log('');
}, 1000);
