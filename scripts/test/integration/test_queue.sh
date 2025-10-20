#!/bin/bash

# Test script for notification queue and circuit breaker functionality
# Tests queue status, circuit breaker metrics, and queue management

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
API_URL="${BASE_URL}/api/v1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Function to print colored output
print_test() {
    echo -e "${BLUE}[TEST $1]${NC} $2"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
    ((TESTS_FAILED++))
}

print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

# Function to make API calls
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4

    ((TESTS_RUN++))
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "${API_URL}${endpoint}")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            "${API_URL}${endpoint}")
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -eq "$expected_status" ]; then
        print_success "HTTP $http_code - $method $endpoint"
        echo "$body"
        return 0
    else
        print_error "HTTP $http_code (expected $expected_status) - $method $endpoint"
        echo "$body"
        return 1
    fi
}

# Main test execution
echo "========================================"
echo "    Notification Queue Test Suite"
echo "========================================"
echo "Base URL: $BASE_URL"
echo ""

# Test 1: Health check
print_test 1 "Health check"
if api_call "GET" "/health" "" 200 > /dev/null; then
    print_info "Server is healthy"
else
    print_error "Server is not responding"
    exit 1
fi
echo ""

# Test 2: Get queue status
print_test 2 "Get queue status"
if queue_status=$(api_call "GET" "/queue/stats" "" 200); then
    print_info "Queue Status:"
    echo "$queue_status" | jq '.' 2>/dev/null || echo "$queue_status"
else
    print_error "Failed to get queue status"
fi
echo ""

# Test 3: Get circuit breaker metrics
print_test 3 "Get circuit breaker metrics"
if cb_metrics=$(api_call "GET" "/queue/circuit-breaker/metrics" "" 200); then
    print_info "Circuit Breaker Metrics:"
    echo "$cb_metrics" | jq '.' 2>/dev/null || echo "$cb_metrics"
    
    # Check circuit breaker state
    cb_state=$(echo "$cb_metrics" | jq -r '.state' 2>/dev/null || echo "unknown")
    print_info "Circuit Breaker State: $cb_state"
else
    print_error "Failed to get circuit breaker metrics"
fi
echo ""

# Test 4: Reset circuit breaker (if it's open)
if [ "$cb_state" = "open" ]; then
    print_test 4 "Reset circuit breaker"
    if api_call "POST" "/queue/circuit-breaker/reset" "" 200 > /dev/null; then
        print_info "Circuit breaker reset successfully"
    else
        print_error "Failed to reset circuit breaker"
    fi
    echo ""
fi

# Test 5: Trigger a notification to test queue (optional)
print_test 5 "Send test notification (to potentially trigger queue)"
print_info "This test sends a notification that might fail and enter the queue"

# You would need a valid notify_key and targets for this
# For now, we'll skip this test if no test data is available
TEST_NOTIFY_KEY="${TEST_NOTIFY_KEY:-}"
if [ -n "$TEST_NOTIFY_KEY" ]; then
    test_notification='{
        "notifyKey": "'"$TEST_NOTIFY_KEY"'",
        "message": "Queue test notification - '"$(date +%s)"'",
        "targets": ["test@example.com"]
    }'
    
    if api_call "POST" "/external/notify" "$test_notification" 200 > /dev/null; then
        print_info "Test notification sent"
        sleep 2  # Wait a bit for processing
        
        # Check queue status again
        print_info "Checking queue status after notification..."
        api_call "GET" "/queue/stats" "" 200 | jq '.' 2>/dev/null
    else
        print_info "Test notification failed (expected - for testing queue)"
    fi
else
    print_info "Skipping test notification (set TEST_NOTIFY_KEY to enable)"
fi
echo ""

# Test 6: Monitor queue over time (optional)
print_test 6 "Monitor queue status changes"
print_info "Monitoring queue for 10 seconds..."

for i in {1..5}; do
    sleep 2
    status=$(api_call "GET" "/queue/stats" "" 200 2>/dev/null)
    pending=$(echo "$status" | jq -r '.total_pending' 2>/dev/null || echo "0")
    retrying=$(echo "$status" | jq -r '.total_retrying' 2>/dev/null || echo "0")
    failed=$(echo "$status" | jq -r '.total_failed' 2>/dev/null || echo "0")
    circuit=$(echo "$status" | jq -r '.circuit_state' 2>/dev/null || echo "unknown")
    
    print_info "[$i] Pending: $pending, Retrying: $retrying, Failed: $failed, Circuit: $circuit"
done
echo ""

# Summary
echo "========================================"
echo "           Test Summary"
echo "========================================"
echo "Tests Run:    $TESTS_RUN"
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi

