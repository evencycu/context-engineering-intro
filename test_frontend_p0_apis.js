/**
 * Frontend P0 API 測試腳本
 * 在瀏覽器控制台中執行此腳本來測試前端 API 調用
 * 
 * 使用方法：
 * 1. 打開 http://localhost:3002
 * 2. 打開瀏覽器開發者工具 (F12)
 * 3. 切換到 Console 標籤
 * 4. 複製並執行此腳本
 */

(async function testFrontendP0APIs() {
  console.log('=== 前端 P0 API 測試開始 ===\n');

  // 假設我們有一個 project ID（從當前頁面獲取或使用測試 ID）
  const testProjectId = '750e8400-e29b-41d4-a716-446655440001';

  // 測試 1: API Key Regeneration
  console.log('1. 測試 API Key Regeneration...');
  try {
    const response = await fetch(`http://localhost:8080/internal/v1/projects/${testProjectId}/regenerate-key`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    const data = await response.json();
    if (response.ok && data.data && data.data.apiKey) {
      console.log('✅ API Key Regeneration 成功:', data.data.apiKey);
    } else {
      console.error('❌ API Key Regeneration 失敗:', data);
    }
  } catch (error) {
    console.error('❌ API Key Regeneration 錯誤:', error);
  }

  console.log('\n');

  // 測試 2: Get Audience Lists
  console.log('2. 測試 Get Audience Lists...');
  try {
    const response = await fetch(`http://localhost:8080/internal/v1/projects/${testProjectId}/audience-lists`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });
    const data = await response.json();
    if (response.ok) {
      console.log('✅ Get Audience Lists 成功:', data.data || []);
    } else {
      console.error('❌ Get Audience Lists 失敗:', data);
    }
  } catch (error) {
    console.error('❌ Get Audience Lists 錯誤:', error);
  }

  console.log('\n');

  // 測試 3: Create Audience List
  console.log('3. 測試 Create Audience List...');
  try {
    const response = await fetch(`http://localhost:8080/internal/v1/projects/${testProjectId}/audience-lists`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: `Test List ${Date.now()}`,
        type: 'Static',
        count: 100,
      }),
    });
    const data = await response.json();
    if (response.ok && data.data && data.data.id) {
      console.log('✅ Create Audience List 成功:', data.data);
      window.testAudienceListId = data.data.id; // 保存 ID 供後續測試使用
    } else {
      console.error('❌ Create Audience List 失敗:', data);
    }
  } catch (error) {
    console.error('❌ Create Audience List 錯誤:', error);
  }

  console.log('\n=== 前端 P0 API 測試完成 ===');
  console.log('\n提示: 如果創建成功，可以使用 window.testAudienceListId 進行更新和刪除測試');
})();
