// Template CRUD 自動化測試腳本
// 可以在瀏覽器 Console 中執行，或使用 Puppeteer/Playwright

const TEST_CONFIG = {
  baseURL: 'http://localhost:3000',
  templatesPage: '/#/admin/templates',
  apiBase: 'http://localhost:8080/internal/v1',
  projectId: null, // 將從頁面獲取
};

// 測試結果記錄
const testResults = {
  passed: [],
  failed: [],
  errors: []
};

function logTest(name, passed, error = null) {
  if (passed) {
    testResults.passed.push(name);
    console.log(`✅ ${name}`);
  } else {
    testResults.failed.push(name);
    testResults.errors.push({ test: name, error });
    console.error(`❌ ${name}`, error);
  }
}

// 等待函數
function wait(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// 等待元素出現
async function waitForElement(selector, timeout = 5000) {
  const startTime = Date.now();
  while (Date.now() - startTime < timeout) {
    const element = document.querySelector(selector);
    if (element) return element;
    await wait(100);
  }
  throw new Error(`Element not found: ${selector}`);
}

// 測試 1: 檢查頁面載入
async function testPageLoad() {
  try {
    const title = await waitForElement('h2');
    const titleText = title.textContent;
    if (titleText.includes('Message Templates') || titleText.includes('Templates')) {
      logTest('Page Load', true);
      return true;
    } else {
      logTest('Page Load', false, 'Title not found');
      return false;
    }
  } catch (error) {
    logTest('Page Load', false, error.message);
    return false;
  }
}

// 測試 2: 檢查 Create Template 按鈕
async function testCreateButton() {
  try {
    const buttons = Array.from(document.querySelectorAll('button'));
    const createButton = buttons.find(btn => 
      btn.textContent.includes('Create Template') || 
      btn.textContent.includes('Create')
    );
    
    if (createButton) {
      logTest('Create Button Exists', true);
      return createButton;
    } else {
      logTest('Create Button Exists', false, 'Create button not found');
      return null;
    }
  } catch (error) {
    logTest('Create Button Exists', false, error.message);
    return null;
  }
}

// 測試 3: 測試新增 Template
async function testCreateTemplate() {
  try {
    // 點擊 Create Template 按鈕
    const createButton = await testCreateButton();
    if (!createButton) {
      logTest('Create Template - Open Modal', false, 'Create button not found');
      return false;
    }
    
    createButton.click();
    await wait(500);
    
    // 檢查 Modal 是否打開
    const modal = document.querySelector('[role="dialog"]') || 
                  document.querySelector('.fixed.inset-0');
    
    if (!modal) {
      logTest('Create Template - Open Modal', false, 'Modal not opened');
      return false;
    }
    
    logTest('Create Template - Open Modal', true);
    
    // 填寫表單
    const nameInput = modal.querySelector('input[type="text"]') || 
                      Array.from(modal.querySelectorAll('input')).find(inp => 
                        inp.placeholder?.includes('name') || 
                        inp.value === ''
                      );
    
    if (!nameInput) {
      logTest('Create Template - Fill Form', false, 'Name input not found');
      return false;
    }
    
    // 填寫名稱
    nameInput.value = `Test Template ${Date.now()}`;
    nameInput.dispatchEvent(new Event('input', { bubbles: true }));
    nameInput.dispatchEvent(new Event('change', { bubbles: true }));
    
    // 填寫描述
    const descriptionInput = modal.querySelector('textarea');
    if (descriptionInput) {
      descriptionInput.value = 'Automated test template';
      descriptionInput.dispatchEvent(new Event('input', { bubbles: true }));
      descriptionInput.dispatchEvent(new Event('change', { bubbles: true }));
    }
    
    // 填寫 JSON
    const jsonInputs = Array.from(modal.querySelectorAll('textarea'));
    const jsonInput = jsonInputs.find(ta => 
      ta.className.includes('font-mono') || 
      ta.className.includes('bg-slate')
    );
    
    if (jsonInput) {
      const testJSON = JSON.stringify({
        type: "AdaptiveCard",
        version: "1.4",
        body: [
          { type: "TextBlock", text: "{{title}}" },
          { type: "TextBlock", text: "{{message}}" }
        ]
      }, null, 2);
      
      jsonInput.value = testJSON;
      jsonInput.dispatchEvent(new Event('input', { bubbles: true }));
      jsonInput.dispatchEvent(new Event('change', { bubbles: true }));
    }
    
    logTest('Create Template - Fill Form', true);
    
    // 點擊 Save 按鈕
    const saveButton = Array.from(modal.querySelectorAll('button')).find(btn =>
      btn.textContent.includes('Save') && !btn.disabled
    );
    
    if (!saveButton) {
      logTest('Create Template - Save', false, 'Save button not found');
      return false;
    }
    
    // 監聽 API 請求
    const originalFetch = window.fetch;
    let apiCalled = false;
    let apiSuccess = false;
    
    window.fetch = function(...args) {
      const url = args[0];
      if (typeof url === 'string' && url.includes('/templates')) {
        apiCalled = true;
        return originalFetch.apply(this, args).then(response => {
          if (response.ok) {
            apiSuccess = true;
          }
          return response;
        });
      }
      return originalFetch.apply(this, args);
    };
    
    saveButton.click();
    await wait(2000);
    
    // 恢復 fetch
    window.fetch = originalFetch;
    
    if (apiCalled) {
      logTest('Create Template - API Called', true);
      if (apiSuccess) {
        logTest('Create Template - API Success', true);
      } else {
        logTest('Create Template - API Success', false, 'API returned error');
      }
    } else {
      logTest('Create Template - API Called', false, 'API not called');
    }
    
    return apiSuccess;
  } catch (error) {
    logTest('Create Template', false, error.message);
    return false;
  }
}

// 測試 4: 檢查 Console 錯誤
function testConsoleErrors() {
  const errors = window.console._errors || [];
  if (errors.length === 0) {
    logTest('Console Errors', true);
    return true;
  } else {
    logTest('Console Errors', false, `Found ${errors.length} errors: ${errors.join(', ')}`);
    return false;
  }
}

// 主測試函數
async function runAllTests() {
  console.log('🧪 Starting Template CRUD Automated Tests...\n');
  
  // 攔截 console.error
  const originalError = console.error;
  const errors = [];
  console.error = function(...args) {
    errors.push(args.join(' '));
    originalError.apply(console, args);
  };
  window.console._errors = errors;
  
  await testPageLoad();
  await wait(500);
  
  await testCreateButton();
  await wait(500);
  
  await testCreateTemplate();
  await wait(1000);
  
  testConsoleErrors();
  
  // 顯示測試結果
  console.log('\n📊 Test Results Summary:');
  console.log(`✅ Passed: ${testResults.passed.length}`);
  console.log(`❌ Failed: ${testResults.failed.length}`);
  
  if (testResults.failed.length > 0) {
    console.log('\n❌ Failed Tests:');
    testResults.failed.forEach(test => {
      console.log(`  - ${test}`);
    });
    console.log('\n🔍 Errors:');
    testResults.errors.forEach(({ test, error }) => {
      console.log(`  - ${test}: ${error}`);
    });
  }
  
  return {
    passed: testResults.passed.length,
    failed: testResults.failed.length,
    results: testResults
  };
}

// 如果在瀏覽器中執行
if (typeof window !== 'undefined') {
  window.runTemplateCRUDTests = runAllTests;
  console.log('📝 Test functions loaded. Run: runTemplateCRUDTests()');
}

// 導出給 Node.js 使用
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { runAllTests, testPageLoad, testCreateButton, testCreateTemplate };
}
