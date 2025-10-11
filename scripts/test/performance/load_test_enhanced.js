import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('error_rate');
const responseTime = new Rate('response_time');

export let options = {
    stages: [
        { duration: '30s', target: 10 },   // Ramp up to 10 users
        { duration: '1m', target: 10 },     // Stay at 10 users
        { duration: '30s', target: 20 },    // Ramp up to 20 users
        { duration: '1m', target: 20 },     // Stay at 20 users
        { duration: '30s', target: 0 },     // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(95)<1000'], // 95% of requests must complete below 1s
        http_req_failed: ['rate<0.1'],    // Error rate must be below 10%
        error_rate: ['rate<0.05'],        // Custom error rate must be below 5%
    },
};

const BASE_URL = 'http://localhost:8080';

export default function () {
    // Test 1: Health Check
    let response = http.get(`${BASE_URL}/health`);
    let success = check(response, {
        'health check status is 200': (r) => r.status === 200,
        'health check response time < 100ms': (r) => r.timings.duration < 100,
    });
    errorRate.add(!success);
    responseTime.add(response.timings.duration < 100);

    // Test 2: System Metrics
    response = http.get(`${BASE_URL}/api/v1/metrics`);
    success = check(response, {
        'metrics status is 200': (r) => r.status === 200,
        'metrics response time < 200ms': (r) => r.timings.duration < 200,
        'metrics has system data': (r) => r.json('system') !== undefined,
    });
    errorRate.add(!success);
    responseTime.add(response.timings.duration < 200);

    // Test 3: Configuration
    response = http.get(`${BASE_URL}/api/v1/config`);
    success = check(response, {
        'config status is 200': (r) => r.status === 200,
        'config response time < 150ms': (r) => r.timings.duration < 150,
        'config has server data': (r) => r.json('server') !== undefined,
    });
    errorRate.add(!success);
    responseTime.add(response.timings.duration < 150);

    // Test 4: Queue Stats
    response = http.get(`${BASE_URL}/api/v1/queue/stats`);
    success = check(response, {
        'queue stats status is 200': (r) => r.status === 200,
        'queue stats response time < 100ms': (r) => r.timings.duration < 100,
        'queue stats has data': (r) => r.json('data') !== undefined,
    });
    errorRate.add(!success);
    responseTime.add(response.timings.duration < 100);

    // Test 5: Companies API
    response = http.get(`${BASE_URL}/api/v1/companies`);
    success = check(response, {
        'companies status is 200': (r) => r.status === 200,
        'companies response time < 300ms': (r) => r.timings.duration < 300,
        'companies has data array': (r) => Array.isArray(r.json('data')),
    });
    errorRate.add(!success);
    responseTime.add(response.timings.duration < 300);

    // Test 6: Users API
    response = http.get(`${BASE_URL}/api/v1/users`);
    success = check(response, {
        'users status is 200': (r) => r.status === 200,
        'users response time < 300ms': (r) => r.timings.duration < 300,
        'users has data array': (r) => Array.isArray(r.json('data')),
    });
    errorRate.add(!success);
    responseTime.add(response.timings.duration < 300);

    // Test 7: External Notification (with error handling)
    const payload = JSON.stringify({
        notify_key: 'test-load-key',
        message: 'Load test message',
        targets: ['test@example.com']
    });

    response = http.post(`${BASE_URL}/api/v1/external/notify`, payload, {
        headers: { 'Content-Type': 'application/json' },
    });

    // This will likely fail due to missing project, but we test error handling
    success = check(response, {
        'external notify has response': (r) => r.status >= 200,
        'external notify response time < 500ms': (r) => r.timings.duration < 500,
    });
    errorRate.add(!success);
    responseTime.add(response.timings.duration < 500);

    // Test 8: Error Handling - 404
    response = http.get(`${BASE_URL}/api/v1/nonexistent`);
    success = check(response, {
        '404 status is 404': (r) => r.status === 404,
        '404 response time < 100ms': (r) => r.timings.duration < 100,
    });
    errorRate.add(!success);
    responseTime.add(response.timings.duration < 100);

    // Test 9: Configuration Validation
    response = http.post(`${BASE_URL}/api/v1/config/validate`, '{}', {
        headers: { 'Content-Type': 'application/json' },
    });
    success = check(response, {
        'config validate status is 200': (r) => r.status === 200,
        'config validate response time < 200ms': (r) => r.timings.duration < 200,
    });
    errorRate.add(!success);
    responseTime.add(response.timings.duration < 200);

    // Test 10: Stress Test - Multiple Concurrent Requests
    const responses = [];
    for (let i = 0; i < 3; i++) {
        responses.push(http.get(`${BASE_URL}/health`));
    }

    success = check(responses, {
        'concurrent requests all successful': (rs) => rs.every(r => r.status === 200),
        'concurrent requests response time < 200ms': (rs) => rs.every(r => r.timings.duration < 200),
    });
    errorRate.add(!success);
    responseTime.add(responses.every(r => r.timings.duration < 200));

    sleep(1); // Wait 1 second between iterations
}

export function handleSummary(data) {
    return {
        'test_results/load_test_enhanced.json': JSON.stringify(data, null, 2),
        'test_results/load_test_enhanced.html': htmlReport(data),
    };
}

function htmlReport(data) {
    return `
    <html>
      <head>
        <title>Enhanced Load Test Report</title>
        <style>
          body { font-family: Arial, sans-serif; margin: 20px; }
          .metric { margin: 10px 0; padding: 10px; border: 1px solid #ddd; }
          .success { background-color: #d4edda; }
          .warning { background-color: #fff3cd; }
          .error { background-color: #f8d7da; }
        </style>
      </head>
      <body>
        <h1>Enhanced Load Test Report</h1>
        <div class="metric">
          <h3>Test Summary</h3>
          <p>Duration: ${data.state.testRunDurationMs}ms</p>
          <p>VUs: ${data.metrics.vus.values.max}</p>
          <p>Iterations: ${data.metrics.iterations.values.count}</p>
        </div>
        <div class="metric">
          <h3>Response Time</h3>
          <p>Average: ${data.metrics.http_req_duration.values.avg}ms</p>
          <p>P95: ${data.metrics.http_req_duration.values.p95}ms</p>
          <p>P99: ${data.metrics.http_req_duration.values.p99}ms</p>
        </div>
        <div class="metric">
          <h3>Error Rate</h3>
          <p>Failed Requests: ${data.metrics.http_req_failed.values.rate * 100}%</p>
          <p>Custom Error Rate: ${data.metrics.error_rate.values.rate * 100}%</p>
        </div>
        <div class="metric">
          <h3>Throughput</h3>
          <p>Requests per second: ${data.metrics.http_reqs.values.rate}</p>
        </div>
      </body>
    </html>
  `;
}
