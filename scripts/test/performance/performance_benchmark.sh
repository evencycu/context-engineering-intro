#!/bin/bash

# Performance Benchmark Test Script for Teams Notification API
# Tests various performance aspects and establishes benchmarks

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
API_BASE_URL="http://localhost:8080"
TEST_RESULTS_DIR="test_results/performance"
REPORT_FILE="$TEST_RESULTS_DIR/performance_benchmark_$(date +%Y%m%d_%H%M%S).json"

# Create test results directory
mkdir -p "$TEST_RESULTS_DIR"

# Performance thresholds (in milliseconds)
HEALTH_CHECK_THRESHOLD=100
METRICS_THRESHOLD=200
CONFIG_THRESHOLD=150
QUEUE_STATS_THRESHOLD=100
COMPANIES_THRESHOLD=300
USERS_THRESHOLD=300
EXTERNAL_NOTIFY_THRESHOLD=500

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Helper functions
log() {
    echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
    ((PASSED_TESTS++))
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    ((FAILED_TESTS++))
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Performance test functions
test_endpoint_performance() {
    local endpoint="$1"
    local threshold="$2"
    local test_name="$3"
    local method="${4:-GET}"
    local data="${5:-}"
    
    log "測試 $test_name 性能..."
    ((TOTAL_TESTS++))
    
    local start_time=$(date +%s)
    local response
    local end_time
    
    if [[ "$method" == "POST" ]]; then
        response=$(curl -s -w "%{http_code}" -o /dev/null \
            -X POST "$API_BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -w "%{http_code}" -o /dev/null \
            -X GET "$API_BASE_URL$endpoint")
    fi
    
    end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    if [[ "$response" == "200" || "$response" == "404" ]]; then
        if [[ $duration -le $threshold ]]; then
            success "$test_name 性能良好: ${duration}ms (閾值: ${threshold}ms)"
            return 0
        else
            warning "$test_name 性能較慢: ${duration}ms (閾值: ${threshold}ms)"
            return 1
        fi
    else
        error "$test_name 請求失敗: HTTP $response"
        return 1
    fi
}

test_concurrent_performance() {
    local endpoint="$1"
    local concurrent_users="$2"
    local test_name="$3"
    
    log "測試 $test_name 併發性能 ($concurrent_users 個併發用戶)..."
    ((TOTAL_TESTS++))
    
    local start_time=$(date +%s)
    local pids=()
    
    # Start concurrent requests
    for i in $(seq 1 $concurrent_users); do
        (
            curl -s -w "%{http_code}" -o /dev/null \
                -X GET "$API_BASE_URL$endpoint" > /tmp/curl_result_$i 2>&1
        ) &
        pids+=($!)
    done
    
    # Wait for all requests to complete
    for pid in "${pids[@]}"; do
        wait $pid
    done
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    # Check results
    local success_count=0
    local total_count=$concurrent_users
    
    for i in $(seq 1 $concurrent_users); do
        if [[ -f "/tmp/curl_result_$i" ]]; then
            local result=$(cat /tmp/curl_result_$i)
            if [[ "$result" == "200" ]]; then
                ((success_count++))
            fi
            rm -f /tmp/curl_result_$i
        fi
    done
    
    local success_rate=$((success_count * 100 / total_count))
    
    if [[ $success_rate -ge 95 ]]; then
        success "$test_name 併發性能良好: ${success_rate}% 成功率, ${duration}ms 總時間"
        return 0
    else
        error "$test_name 併發性能不佳: ${success_rate}% 成功率, ${duration}ms 總時間"
        return 1
    fi
}

test_memory_usage() {
    log "測試記憶體使用情況..."
    ((TOTAL_TESTS++))
    
    # Get memory usage from metrics endpoint
    local response=$(curl -s "$API_BASE_URL/api/v1/metrics")
    local memory_usage=$(echo "$response" | jq -r '.system.memory_usage // 0')
    
    if [[ "$memory_usage" -le 100 ]]; then
        success "記憶體使用正常: ${memory_usage}MB"
        return 0
    else
        warning "記憶體使用較高: ${memory_usage}MB"
        return 1
    fi
}

test_response_consistency() {
    local endpoint="$1"
    local test_name="$2"
    local iterations="$3"
    
    log "測試 $test_name 回應一致性 ($iterations 次請求)..."
    ((TOTAL_TESTS++))
    
    local response_times=()
    local status_codes=()
    
    for i in $(seq 1 $iterations); do
        local start_time=$(date +%s)
        local response=$(curl -s -w "%{http_code}" -o /dev/null \
            -X GET "$API_BASE_URL$endpoint")
        local end_time=$(date +%s)
        
        response_times+=($((end_time - start_time)))
        status_codes+=($response)
        
        sleep 0.1
    done
    
    # Calculate statistics
    local total_time=0
    local min_time=${response_times[0]}
    local max_time=${response_times[0]}
    
    for time in "${response_times[@]}"; do
        total_time=$((total_time + time))
        if [[ $time -lt $min_time ]]; then
            min_time=$time
        fi
        if [[ $time -gt $max_time ]]; then
            max_time=$time
        fi
    done
    
    local avg_time=$((total_time / iterations))
    local time_variance=$((max_time - min_time))
    
    # Check if all responses are successful
    local success_count=0
    for status in "${status_codes[@]}"; do
        if [[ "$status" == "200" ]]; then
            ((success_count++))
        fi
    done
    
    if [[ $success_count -eq $iterations && $time_variance -le 100 ]]; then
        success "$test_name 回應一致: 平均 ${avg_time}ms, 變異 ${time_variance}ms"
        return 0
    else
        error "$test_name 回應不一致: ${success_count}/${iterations} 成功, 變異 ${time_variance}ms"
        return 1
    fi
}

test_throughput() {
    local endpoint="$1"
    local test_name="$2"
    local duration_seconds="$3"
    
    log "測試 $test_name 吞吐量 ($duration_seconds 秒)..."
    ((TOTAL_TESTS++))
    
    local start_time=$(date +%s)
    local request_count=0
    local success_count=0
    
    while [[ $(($(date +%s) - start_time)) -lt $duration_seconds ]]; do
        local response=$(curl -s -w "%{http_code}" -o /dev/null \
            -X GET "$API_BASE_URL$endpoint")
        
        ((request_count++))
        if [[ "$response" == "200" ]]; then
            ((success_count++))
        fi
        
        sleep 0.01
    done
    
    local actual_duration=$(($(date +%s) - start_time))
    local requests_per_second=$((request_count / actual_duration))
    local success_rate=$((success_count * 100 / request_count))
    
    if [[ $success_rate -ge 95 && $requests_per_second -ge 10 ]]; then
        success "$test_name 吞吐量良好: ${requests_per_second} req/s, ${success_rate}% 成功率"
        return 0
    else
        error "$test_name 吞吐量不佳: ${requests_per_second} req/s, ${success_rate}% 成功率"
        return 1
    fi
}

# Main test execution
main() {
    log "🚀 開始性能基準測試..."
    echo "=========================================="
    
    # Check if API is running
    if ! curl -s "$API_BASE_URL/health" > /dev/null; then
        error "API 服務未運行，請先啟動服務"
        exit 1
    fi
    
    # Run performance tests
    test_endpoint_performance "/health" $HEALTH_CHECK_THRESHOLD "健康檢查"
    test_endpoint_performance "/api/v1/metrics" $METRICS_THRESHOLD "系統指標"
    test_endpoint_performance "/api/v1/config" $CONFIG_THRESHOLD "配置信息"
    test_endpoint_performance "/api/v1/queue/stats" $QUEUE_STATS_THRESHOLD "佇列狀態"
    test_endpoint_performance "/api/v1/companies" $COMPANIES_THRESHOLD "公司 API"
    test_endpoint_performance "/api/v1/users" $USERS_THRESHOLD "用戶 API"
    
    # Test external notification with error handling
    test_endpoint_performance "/api/v1/notify" $EXTERNAL_NOTIFY_THRESHOLD \
        "外部通知" "POST" '{"notify_key":"test","message":"test","targets":["test"]}'
    
    # Concurrent performance tests
    test_concurrent_performance "/health" 10 "健康檢查併發"
    test_concurrent_performance "/api/v1/metrics" 5 "系統指標併發"
    test_concurrent_performance "/api/v1/companies" 3 "公司 API 併發"
    
    # Memory usage test
    test_memory_usage
    
    # Response consistency tests
    test_response_consistency "/health" "健康檢查一致性" 10
    test_response_consistency "/api/v1/metrics" "系統指標一致性" 5
    
    # Throughput tests
    test_throughput "/health" "健康檢查吞吐量" 5
    test_throughput "/api/v1/metrics" "系統指標吞吐量" 3
    
    # Generate report
    log "📊 生成性能基準測試報告..."
    
    local report_data=$(cat << EOF
{
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "test_summary": {
        "total_tests": $TOTAL_TESTS,
        "passed_tests": $PASSED_TESTS,
        "failed_tests": $FAILED_TESTS,
        "success_rate": "$(echo "scale=2; $PASSED_TESTS * 100 / $TOTAL_TESTS" | bc)%"
    },
    "performance_benchmarks": {
        "health_check_threshold": ${HEALTH_CHECK_THRESHOLD},
        "metrics_threshold": ${METRICS_THRESHOLD},
        "config_threshold": ${CONFIG_THRESHOLD},
        "queue_stats_threshold": ${QUEUE_STATS_THRESHOLD},
        "companies_threshold": ${COMPANIES_THRESHOLD},
        "users_threshold": ${USERS_THRESHOLD},
        "external_notify_threshold": ${EXTERNAL_NOTIFY_THRESHOLD}
    }
}
EOF
)
    
    echo "$report_data" > "$REPORT_FILE"
    
    # Print summary
    echo "=========================================="
    log "📊 性能基準測試結果摘要"
    echo "總測試數: $TOTAL_TESTS"
    echo "通過: $PASSED_TESTS"
    echo "失敗: $FAILED_TESTS"
    echo "成功率: $(echo "scale=2; $PASSED_TESTS * 100 / $TOTAL_TESTS" | bc)%"
    echo "報告文件: $REPORT_FILE"
    echo "=========================================="
    
    if [ $FAILED_TESTS -eq 0 ]; then
        success "🎉 所有性能基準測試通過！"
        exit 0
    else
        warning "⚠️  發現 $FAILED_TESTS 個性能問題，請檢查報告"
        exit 1
    fi
}

# Run main function
main "$@"
