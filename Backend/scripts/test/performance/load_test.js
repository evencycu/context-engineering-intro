// Teams Notification API 負載測試腳本
// 使用 k6 進行效能和壓力測試

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// 自定義指標
const errorRate = new Rate('error_rate');
const responseTime = new Trend('response_time');

// 測試配置
export let options = {
    stages: [
        { duration: '30s', target: 50 },   // 漸增到 50 用戶
        { duration: '1m', target: 50 },   // 維持 50 用戶
        { duration: '30s', target: 100 },  // 漸增到 100 用戶
        { duration: '1m', target: 100 },   // 維持 100 用戶
        { duration: '30s', target: 200 },  // 漸增到 200 用戶
        { duration: '1m', target: 200 },   // 維持 200 用戶
        { duration: '30s', target: 0 },    // 漸減到 0 用戶
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'],  // 95% 請求 < 500ms
        http_req_failed: ['rate<0.05'],    // 錯誤率 < 5%
        error_rate: ['rate<0.05'],         // 自定義錯誤率 < 5%
        response_time: ['p(95)<500'],      // 95% 回應時間 < 500ms
    },
    ext: {
        loadimpact: {
            projectID: 'teams-notification-api',
            name: 'Teams Notification API Load Test'
        }
    }
};

// 測試資料
const TEST_DATA = {
    notify_key: 'test-project-key',
    base_url: 'http://localhost:8080',
    api_base: 'http://localhost:8080/api/v1'
};

// 隨機訊息生成
function generateMessage() {
    const messages = [
        '系統狀態正常',
        '資料庫連線穩定',
        'API 回應時間正常',
        '佇列處理順暢',
        '監控指標良好',
        '效能測試進行中',
        '負載測試執行',
        '壓力測試驗證',
        '併發測試完成',
        '效能驗證通過'
    ];

    const randomMessage = messages[Math.floor(Math.random() * messages.length)];
    const timestamp = new Date().toISOString();
    const userInfo = `[User ${__VU}, Iteration ${__ITER}]`;

    return `${userInfo} ${randomMessage} - ${timestamp}`;
}

// 隨機優先級生成
function generatePriority() {
    const priorities = ['low', 'normal', 'high'];
    return priorities[Math.floor(Math.random() * priorities.length)];
}

// 隨機目標生成
function generateTargets() {
    const targets = ['all', ['specific-target-1'], ['specific-target-2']];
    return targets[Math.floor(Math.random() * targets.length)];
}

// 主要測試函數
export default function () {
    // 生成測試資料
    const message = generateMessage();
    const priority = generatePriority();
    const targets = generateTargets();

    const payload = JSON.stringify({
        notify_key: TEST_DATA.notify_key,
        message: message,
        priority: priority,
        targets: targets,
        metadata: {
            test_type: 'load_test',
            user_id: __VU,
            iteration: __ITER,
            timestamp: new Date().toISOString()
        }
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'User-Agent': 'k6-load-test'
        },
        timeout: '30s'
    };

    // 發送請求
    const startTime = Date.now();
    const response = http.post(`${TEST_DATA.api_base}/external/notify`, payload, params);
    const endTime = Date.now();
    const responseTimeMs = endTime - startTime;

    // 記錄回應時間
    responseTime.add(responseTimeMs);

    // 檢查回應
    const success = check(response, {
        'status is 200': (r) => r.status === 200,
        'response time < 500ms': (r) => r.timings.duration < 500,
        'response time < 1000ms': (r) => r.timings.duration < 1000,
        'has notification_id': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.data && body.data.notification_id !== null;
            } catch (e) {
                return false;
            }
        },
        'response body is valid JSON': (r) => {
            try {
                JSON.parse(r.body);
                return true;
            } catch (e) {
                return false;
            }
        }
    });

    // 記錄錯誤率
    errorRate.add(!success);

    // 如果請求失敗，記錄詳細資訊
    if (!success) {
        console.error(`Request failed for VU ${__VU}, Iteration ${__ITER}:`);
        console.error(`Status: ${response.status}`);
        console.error(`Response: ${response.body}`);
        console.error(`Response Time: ${response.timings.duration}ms`);
    }

    // 隨機等待時間 (1-3 秒)
    sleep(Math.random() * 2 + 1);
}

// 測試設置
export function setup() {
    console.log('🚀 開始負載測試...');
    console.log(`測試目標: ${TEST_DATA.base_url}`);
    console.log(`測試時間: ${new Date().toISOString()}`);

    // 檢查服務是否可用
    const healthCheck = http.get(`${TEST_DATA.base_url}/health`);
    if (healthCheck.status !== 200) {
        throw new Error('服務健康檢查失敗');
    }

    console.log('✅ 服務健康檢查通過');
    return { startTime: Date.now() };
}

// 測試清理
export function teardown(data) {
    const endTime = Date.now();
    const duration = (endTime - data.startTime) / 1000;

    console.log('📊 負載測試完成');
    console.log(`測試持續時間: ${duration}秒`);
    console.log(`測試時間: ${new Date().toISOString()}`);
}

// 自定義摘要報告
export function handleSummary(data) {
    const summary = {
        timestamp: new Date().toISOString(),
        test_duration: `${data.state.testRunDurationMs / 1000}s`,
        total_requests: data.metrics.http_reqs.values.count,
        successful_requests: data.metrics.http_reqs.values.count - data.metrics.http_req_failed.values.count,
        failed_requests: data.metrics.http_req_failed.values.count,
        error_rate: `${(data.metrics.http_req_failed.values.rate * 100).toFixed(2)}%`,
        avg_response_time: `${data.metrics.http_req_duration.values.avg.toFixed(2)}ms`,
        p95_response_time: `${data.metrics.http_req_duration.values.p95.toFixed(2)}ms`,
        p99_response_time: `${data.metrics.http_req_duration.values.p99.toFixed(2)}ms`,
        max_response_time: `${data.metrics.http_req_duration.values.max.toFixed(2)}ms`,
        requests_per_second: `${data.metrics.http_reqs.values.rate.toFixed(2)}`,
        data_transferred: `${(data.metrics.data_sent.values.count + data.metrics.data_received.values.count) / 1024}KB`
    };

    // 生成 JSON 報告
    const jsonReport = JSON.stringify(summary, null, 2);

    // 生成文字報告
    const textReport = `
📊 Teams Notification API 負載測試報告
=====================================
測試時間: ${summary.timestamp}
測試持續時間: ${summary.test_duration}
總請求數: ${summary.total_requests}
成功請求數: ${summary.successful_requests}
失敗請求數: ${summary.failed_requests}
錯誤率: ${summary.error_rate}
平均回應時間: ${summary.avg_response_time}
P95 回應時間: ${summary.p95_response_time}
P99 回應時間: ${summary.p99_response_time}
最大回應時間: ${summary.max_response_time}
每秒請求數: ${summary.requests_per_second}
資料傳輸量: ${summary.data_transferred}
=====================================
`;

    return {
        'test_results/load_test_summary.json': jsonReport,
        'test_results/load_test_summary.txt': textReport,
        stdout: textReport
    };
}
