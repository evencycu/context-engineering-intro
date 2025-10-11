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

# 3. Third-Party Bot 測試
echo -e "\n${BLUE}📋 Third-Party Bot 測試${NC}"

# 創建Third-Party Bot
echo "  📝 創建Third-Party Bot..."
THIRD_PARTY_BOT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/bots/third-party" \
    -H "Content-Type: application/json" \
    -d '{
        "company_id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Test Third-Party Bot",
        "description": "A test third-party bot for automated testing",
        "app_id": "test-third-party-bot-'$(date +%s)'",
        "app_password": "test-password-123",
        "tenant_id": "test-tenant-123",
        "webhook_url": "https://example.com/webhook",
        "api_endpoint": "https://api.example.com/bot",
        "api_key": "test-api-key-123",
        "contact_email": "test@example.com",
        "contact_phone": "+1234567890",
        "capabilities": {"messaging": true, "notifications": true},
        "rate_limit_per_minute": 60,
        "max_concurrent_requests": 10,
        "created_by": "550e8400-e29b-41d4-a716-446655440000"
    }')

if echo "$THIRD_PARTY_BOT_RESPONSE" | jq -e '.data.id' > /dev/null; then
    THIRD_PARTY_BOT_ID=$(echo "$THIRD_PARTY_BOT_RESPONSE" | jq -r '.data.id')
    echo "    ✅ Third-Party Bot 創建成功 (ID: $THIRD_PARTY_BOT_ID)"
    ((PASSED_TESTS++))
else
    echo "    ❌ Third-Party Bot 創建失敗"
    echo "    Response: $THIRD_PARTY_BOT_RESPONSE"
    ((FAILED_TESTS++))
fi

# Third-Party Bot 測試
if [ -n "$THIRD_PARTY_BOT_ID" ]; then
    run_test "Third-Party Bot 列表" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/v1/bots/third-party?limit=5'" "200"
    run_test "Third-Party Bot 獲取" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/v1/bots/third-party/$THIRD_PARTY_BOT_ID'" "200"
    run_test "Third-Party Bot 狀態更新" "curl -s -w '%{http_code}' -X PATCH '$BASE_URL/api/v1/bots/third-party/$THIRD_PARTY_BOT_ID/status' -H 'Content-Type: application/json' -d '{\"status\": \"inactive\"}'" "200"
    run_test "Third-Party Bot 連接測試" "curl -s -w '%{http_code}' -X POST '$BASE_URL/api/v1/bots/third-party/$THIRD_PARTY_BOT_ID/test'" "200"
    
    # API Key 測試
    echo "  🔑 測試API Key功能..."
    API_KEY_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/bots/third-party/$THIRD_PARTY_BOT_ID/api-key" \
        -H "Content-Type: application/json" \
        -d '{
            "notify_key": "test-notify-key",
            "permissions": {"send_message": true, "read_status": true},
            "rate_limit_per_minute": 60,
            "created_by": "550e8400-e29b-41d4-a716-446655440000"
        }')
    
    if echo "$API_KEY_RESPONSE" | jq -e '.data.api_key' > /dev/null; then
        echo "    ✅ API Key 生成成功"
        ((PASSED_TESTS++))
    else
        echo "    ❌ API Key 生成失敗"
        echo "    Response: $API_KEY_RESPONSE"
        ((FAILED_TESTS++))
    fi
    
    run_test "API Key 撤銷" "curl -s -w '%{http_code}' -X DELETE '$BASE_URL/api/v1/bots/third-party/$THIRD_PARTY_BOT_ID/api-key'" "200"
    
    # 清理Third-Party Bot
    echo "  🧹 清理Third-Party Bot..."
    THIRD_PARTY_BOT_DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/v1/bots/third-party/$THIRD_PARTY_BOT_ID")
    if echo "$THIRD_PARTY_BOT_DELETE_RESPONSE" | jq -e '.message' > /dev/null; then
        echo "    ✅ Third-Party Bot 刪除成功"
        ((PASSED_TESTS++))
    else
        echo "    ❌ Third-Party Bot 刪除失敗"
        ((FAILED_TESTS++))
    fi
fi

# 4. 外部API功能測試
echo -e "\n${BLUE}📋 外部API功能測試${NC}"
run_test "驗證notify key" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/$API_VERSION/external/validate/$NOTIFY_KEY'" "200"
run_test "獲取專案目的地" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/$API_VERSION/external/destinations/$NOTIFY_KEY'" "200"

# 5. 通知發送測試
echo -e "\n${BLUE}📋 通知發送測試${NC}"
run_test "基本通知發送" "curl -s -w '%{http_code}' -X POST '$BASE_URL/api/$API_VERSION/external/notify' -H 'Content-Type: application/json' -d '{\"notify_key\": \"$NOTIFY_KEY\", \"message\": \"自動化測試訊息 - $(date)\"}'" "200"

# 6. 效能測試
echo -e "\n${BLUE}📋 效能測試${NC}"
echo "執行5次連續請求測試..."

for i in {1..5}; do
    echo "效能測試 $i/5"
    start_time=$(date +%s.%N)
    curl -s "$BASE_URL/api/$API_VERSION/external/health" > /dev/null
    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)
    echo "  響應時間: ${duration}s"
    echo "效能測試 $i: ${duration}s" >> "$TEST_RESULTS_DIR/performance_${TIMESTAMP}.log"
done

# 7. 錯誤處理測試
echo -e "\n${BLUE}📋 錯誤處理測試${NC}"
run_test "無效notify key" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/$API_VERSION/external/validate/invalid-key'" "200"
run_test "缺少必要參數" "curl -s -w '%{http_code}' -X POST '$BASE_URL/api/$API_VERSION/external/notify' -H 'Content-Type: application/json' -d '{}'" "400"

# 8. 負載測試
echo -e "\n${BLUE}📋 負載測試${NC}"
echo "執行10個並發請求..."

for i in {1..10}; do
    (
        curl -s "$BASE_URL/health" > /dev/null &
    ) &
done
wait
echo "✅ 並發測試完成"

# 測試結果摘要
echo -e "\n${BLUE}📊 測試結果摘要${NC}"
echo "=========================================="
echo "總測試數: $TOTAL_TESTS"
echo -e "通過: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失敗: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有測試通過！系統運行正常${NC}"
    echo -e "${GREEN}✅ Platform Bot 和 Third-Party Bot 功能完整${NC}"
    exit 0
else
    echo -e "\n${RED}⚠️  有 $FAILED_TESTS 個測試失敗，請檢查日誌${NC}"
    exit 1
fi
