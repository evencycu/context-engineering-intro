#!/bin/bash

# Security Test Script for Teams Notification API
# Tests various security aspects of the API

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
API_BASE_URL="http://localhost:8080"
TEST_RESULTS_DIR="test_results/security"
REPORT_FILE="$TEST_RESULTS_DIR/security_test_$(date +%Y%m%d_%H%M%S).json"

# Create test results directory
mkdir -p "$TEST_RESULTS_DIR"

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

test_security() {
    local test_name="$1"
    local test_function="$2"
    
    log "🧪 執行安全測試: $test_name"
    ((TOTAL_TESTS++))
    
    if $test_function; then
        success "$test_name 通過"
        return 0
    else
        error "$test_name 失敗"
        return 1
    fi
}

# Security test functions
test_sql_injection() {
    log "測試 SQL 注入防護..."
    
    # Test SQL injection in various parameters
    local sql_payloads=(
        "'; DROP TABLE users; --"
        "1' OR '1'='1"
        "admin'--"
        "1' UNION SELECT * FROM users--"
    )
    
    local failed=0
    for payload in "${sql_payloads[@]}"; do
        local response=$(curl -s -w "%{http_code}" -o /dev/null \
            -X GET "$API_BASE_URL/api/v1/companies?search=$payload" \
            -H "Content-Type: application/json")
        
        if [[ "$response" == "400" || "$response" == "422" ]]; then
            success "SQL 注入防護有效: $payload"
        else
            error "SQL 注入防護可能無效: $payload (HTTP $response)"
            ((failed++))
        fi
    done
    
    return $failed
}

test_xss_protection() {
    log "測試 XSS 防護..."
    
    local xss_payloads=(
        "<script>alert('XSS')</script>"
        "javascript:alert('XSS')"
        "<img src=x onerror=alert('XSS')>"
        "';alert('XSS');//"
    )
    
    local failed=0
    for payload in "${xss_payloads[@]}"; do
        local response=$(curl -s -X POST "$API_BASE_URL/api/v1/notify" \
            -H "Content-Type: application/json" \
            -d "{\"notify_key\":\"$payload\",\"message\":\"test\",\"targets\":[\"test\"]}")
        
        # Check if XSS payload is reflected in response
        if echo "$response" | grep -q "$payload"; then
            error "XSS 防護可能無效: $payload"
            ((failed++))
        else
            success "XSS 防護有效: $payload"
        fi
    done
    
    return $failed
}

test_rate_limiting() {
    log "測試速率限制..."
    
    local rate_limit_requests=0
    local rate_limit_hit=false
    
    # Send rapid requests to trigger rate limiting
    for i in {1..150}; do
        local response=$(curl -s -w "%{http_code}" -o /dev/null \
            -X GET "$API_BASE_URL/health")
        
        if [[ "$response" == "429" ]]; then
            rate_limit_hit=true
            break
        fi
        
        ((rate_limit_requests++))
        sleep 0.01
    done
    
    if $rate_limit_hit; then
        success "速率限制有效 (在 $rate_limit_requests 請求後觸發)"
        return 0
    else
        warning "速率限制可能未生效 (發送了 $rate_limit_requests 請求)"
        return 1
    fi
}

test_authentication_bypass() {
    log "測試認證繞過..."
    
    # Test accessing protected endpoints without authentication
    local protected_endpoints=(
        "/api/v1/companies"
        "/api/v1/users"
        "/api/v1/projects"
    )
    
    local failed=0
    for endpoint in "${protected_endpoints[@]}"; do
        local response=$(curl -s -w "%{http_code}" -o /dev/null \
            -X GET "$API_BASE_URL$endpoint")
        
        if [[ "$response" == "401" || "$response" == "403" ]]; then
            success "認證保護有效: $endpoint"
        else
            error "認證保護可能無效: $endpoint (HTTP $response)"
            ((failed++))
        fi
    done
    
    return $failed
}

test_input_validation() {
    log "測試輸入驗證..."
    
    local invalid_inputs=(
        '{"notify_key":"","message":"","targets":[]}'
        '{"notify_key":null,"message":null,"targets":null}'
        '{"notify_key":123,"message":456,"targets":"invalid"}'
        '{"invalid_field":"test"}'
    )
    
    local failed=0
    for input in "${invalid_inputs[@]}"; do
        local response=$(curl -s -X POST "$API_BASE_URL/api/v1/notify" \
            -H "Content-Type: application/json" \
            -d "$input")
        
        if echo "$response" | grep -q "error\|validation\|invalid"; then
            success "輸入驗證有效: $input"
        else
            error "輸入驗證可能無效: $input"
            ((failed++))
        fi
    done
    
    return $failed
}

test_headers_security() {
    log "測試安全標頭..."
    
    local response=$(curl -s -I "$API_BASE_URL/health")
    
    local security_headers=(
        "X-Content-Type-Options"
        "X-Frame-Options"
        "X-XSS-Protection"
        "Strict-Transport-Security"
    )
    
    local missing_headers=0
    for header in "${security_headers[@]}"; do
        if echo "$response" | grep -qi "$header"; then
            success "安全標頭存在: $header"
        else
            warning "安全標頭缺失: $header"
            ((missing_headers++))
        fi
    done
    
    return $missing_headers
}

test_cors_configuration() {
    log "測試 CORS 配置..."
    
    local response=$(curl -s -H "Origin: https://malicious.com" \
        -H "Access-Control-Request-Method: POST" \
        -X OPTIONS "$API_BASE_URL/api/v1/notify")
    
    if echo "$response" | grep -qi "access-control-allow-origin"; then
        success "CORS 配置存在"
        return 0
    else
        warning "CORS 配置可能缺失"
        return 1
    fi
}

test_error_information_disclosure() {
    log "測試錯誤信息洩露..."
    
    # Test various error conditions
    local error_tests=(
        "GET /api/v1/nonexistent"
        "POST /api/v1/notify -d '{}'"
        "GET /api/v1/companies/999999"
    )
    
    local information_disclosed=0
    for test in "${error_tests[@]}"; do
        local response=$(curl -s $test 2>/dev/null)
        
        # Check for sensitive information in error responses
        if echo "$response" | grep -qi "database\|sql\|stack\|trace\|internal"; then
            error "可能洩露敏感信息: $test"
            ((information_disclosed++))
        else
            success "錯誤信息安全: $test"
        fi
    done
    
    return $information_disclosed
}

# Main test execution
main() {
    log "🔒 開始安全測試..."
    echo "=========================================="
    
    # Check if API is running
    if ! curl -s "$API_BASE_URL/health" > /dev/null; then
        error "API 服務未運行，請先啟動服務"
        exit 1
    fi
    
    # Run security tests
    test_security "SQL 注入防護" test_sql_injection
    test_security "XSS 防護" test_xss_protection
    test_security "速率限制" test_rate_limiting
    test_security "認證繞過防護" test_authentication_bypass
    test_security "輸入驗證" test_input_validation
    test_security "安全標頭" test_headers_security
    test_security "CORS 配置" test_cors_configuration
    test_security "錯誤信息洩露防護" test_error_information_disclosure
    
    # Generate report
    log "📊 生成安全測試報告..."
    
    local report_data=$(cat << EOF
{
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "test_summary": {
        "total_tests": $TOTAL_TESTS,
        "passed_tests": $PASSED_TESTS,
        "failed_tests": $FAILED_TESTS,
        "success_rate": "$(echo "scale=2; $PASSED_TESTS * 100 / $TOTAL_TESTS" | bc)%"
    },
    "security_tests": {
        "sql_injection_protection": $([ $? -eq 0 ] && echo "passed" || echo "failed"),
        "xss_protection": $([ $? -eq 0 ] && echo "passed" || echo "failed"),
        "rate_limiting": $([ $? -eq 0 ] && echo "passed" || echo "failed"),
        "authentication_bypass_protection": $([ $? -eq 0 ] && echo "passed" || echo "failed"),
        "input_validation": $([ $? -eq 0 ] && echo "passed" || echo "failed"),
        "security_headers": $([ $? -eq 0 ] && echo "passed" || echo "failed"),
        "cors_configuration": $([ $? -eq 0 ] && echo "passed" || echo "failed"),
        "error_information_disclosure": $([ $? -eq 0 ] && echo "passed" || echo "failed")
    }
}
EOF
)
    
    echo "$report_data" > "$REPORT_FILE"
    
    # Print summary
    echo "=========================================="
    log "📊 安全測試結果摘要"
    echo "總測試數: $TOTAL_TESTS"
    echo "通過: $PASSED_TESTS"
    echo "失敗: $FAILED_TESTS"
    echo "成功率: $(echo "scale=2; $PASSED_TESTS * 100 / $TOTAL_TESTS" | bc)%"
    echo "報告文件: $REPORT_FILE"
    echo "=========================================="
    
    if [ $FAILED_TESTS -eq 0 ]; then
        success "🎉 所有安全測試通過！"
        exit 0
    else
        error "⚠️  發現 $FAILED_TESTS 個安全問題，請檢查報告"
        exit 1
    fi
}

# Run main function
main "$@"
