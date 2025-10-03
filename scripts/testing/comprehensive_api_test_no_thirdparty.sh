#!/bin/bash

# TeamsNotifyGoV2 完整API測試腳本
# 包含Platform Bot和Third-Party Bot的完整測試

set -e

# 配置
BASE_URL="http://localhost:8080"
API_VERSION="v1"
NOTIFY_KEY="cfh-alert-gogo"
TEST_RESULTS_DIR="test_results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 建立測試結果目錄
mkdir -p "$TEST_RESULTS_DIR"

echo -e "${BLUE}🚀 TeamsNotifyGoV2 完整API測試開始${NC}"
echo "測試時間: $(date)"
echo "測試目錄: $TEST_RESULTS_DIR"
echo "=========================================="

# 測試結果統計
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 測試函數
run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_status="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "\n${YELLOW}🧪 測試: $test_name${NC}"
    
    # 執行測試
    local start_time=$(date +%s.%N)
    local response=$(eval "$test_command" 2>&1)
    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc)
    
    # 檢查響應狀態
    local status_code=$(echo "$response" | grep -o "[0-9][0-9][0-9]$" | tail -1)
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ 通過${NC} (${duration}s)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        echo "$test_name: PASSED (${duration}s)" >> "$TEST_RESULTS_DIR/test_${TIMESTAMP}.log"
    else
        echo -e "${RED}❌ 失敗${NC} (期望: $expected_status, 實際: $status_code)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        echo "$test_name: FAILED (期望: $expected_status, 實際: $status_code)" >> "$TEST_RESULTS_DIR/test_${TIMESTAMP}.log"
        echo "響應: $response" >> "$TEST_RESULTS_DIR/test_${TIMESTAMP}.log"
    fi
}

# 1. 系統健康檢查測試
echo -e "\n${BLUE}📋 系統健康檢查測試${NC}"
run_test "系統健康檢查" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/health'" "200"
run_test "外部API健康檢查" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/$API_VERSION/external/health'" "200"

# 2. Platform Bot 測試
echo -e "\n${BLUE}📋 Platform Bot 測試${NC}"

# 創建Platform Bot
echo "  📝 創建Platform Bot..."
PLATFORM_BOT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/bots/platform" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Test Platform Bot",
        "description": "A test platform bot for automated testing",
        "app_id": "test-platform-bot-'$(date +%s)'",
        "app_password": "test-password-123",
        "tenant_id": "test-tenant-123",
        "webhook_url": "https://example.com/webhook",
        "capabilities": {"messaging": true, "notifications": true},
        "rate_limit_per_minute": 60,
        "max_concurrent_requests": 10
    }')

if echo "$PLATFORM_BOT_RESPONSE" | jq -e '.data.id' > /dev/null; then
    PLATFORM_BOT_ID=$(echo "$PLATFORM_BOT_RESPONSE" | jq -r '.data.id')
    echo "    ✅ Platform Bot 創建成功 (ID: $PLATFORM_BOT_ID)"
    ((PASSED_TESTS++))
else
    echo "    ❌ Platform Bot 創建失敗"
    echo "    Response: $PLATFORM_BOT_RESPONSE"
    ((FAILED_TESTS++))
fi

# Platform Bot 測試
if [ -n "$PLATFORM_BOT_ID" ]; then
    run_test "Platform Bot 列表" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/v1/bots/platform?limit=5'" "200"
    run_test "Platform Bot 獲取" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/v1/bots/platform/$PLATFORM_BOT_ID'" "200"
    run_test "Platform Bot 狀態更新" "curl -s -w '%{http_code}' -X PATCH '$BASE_URL/api/v1/bots/platform/$PLATFORM_BOT_ID/status' -H 'Content-Type: application/json' -d '{\"status\": \"inactive\"}'" "200"
    run_test "Platform Bot 連接測試" "curl -s -w '%{http_code}' -X POST '$BASE_URL/api/v1/bots/platform/$PLATFORM_BOT_ID/test'" "200"
    
    # 清理Platform Bot
    echo "  🧹 清理Platform Bot..."
    PLATFORM_BOT_DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/v1/bots/platform/$PLATFORM_BOT_ID")
    if echo "$PLATFORM_BOT_DELETE_RESPONSE" | jq -e '.message' > /dev/null; then
        echo "    ✅ Platform Bot 刪除成功"
        ((PASSED_TESTS++))
    else
        echo "    ❌ Platform Bot 刪除失敗"
        ((FAILED_TESTS++))
    fi
fi

