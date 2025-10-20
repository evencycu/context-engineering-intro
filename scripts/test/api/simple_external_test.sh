#!/bin/bash

# External API 測試腳本
# 測試所有 External API 端點

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# API 基礎 URL
BASE_URL="http://localhost:8080/api/v1"

# 測試計數器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 測試數據 - 從 init.sql 中獲取
TEST_NOTIFY_KEY="cfh-alert-gogo"                              # Project notify_key
TEST_DESTINATION_ID="950e8400-e29b-41d4-a716-446655440003"    # Golden Sample Destination
TEST_PROJECT_ID="750e8400-e29b-41d4-a716-446655440001"        # Holdings notification alerts

# 輔助函數：打印測試標題
print_test() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "${BLUE}[$TOTAL_TESTS/$1]${NC} $2..."
}

# 輔助函數：打印成功
print_success() {
    PASSED_TESTS=$((PASSED_TESTS + 1))
    echo -e "${GREEN}✓ $1${NC}"
}

# 輔助函數：打印失敗
print_fail() {
    FAILED_TESTS=$((FAILED_TESTS + 1))
    echo -e "${RED}✗ $1${NC}"
}

# 輔助函數：打印 curl 命令
print_curl() {
    echo -e "${BLUE}curl $1${NC}"
}

# 輔助函數：打印 JSON 響應
print_json() {
    echo "$1" | python3 -m json.tool 2>/dev/null || echo "$1"
}

echo -e "${BLUE}=========================================="
echo "External API 測試"
echo "==========================================${NC}"
echo ""
echo -e "${YELLOW}測試資源 (from init.sql):${NC}"
echo "  NotifyKey: $TEST_NOTIFY_KEY"
echo "  Project: 750e8400-e29b-41d4-a716-446655440001 (Holdings notification alerts)"
echo "  Destination: 950e8400-e29b-41d4-a716-446655440003 (All Bot Installations - Golden Sample)"
echo "  - Personal Chat"
echo "  - Group Chat"
echo "  - Channel"
echo ""

# 測試 1: 發送基本文字通知
print_test 5 "發送基本文字通知"
NOTIFY_DATA='{
  "notifyKey": "'$TEST_NOTIFY_KEY'",
  "message": "External API Test - Basic Text Message",
  "messageType": "text",
  "priority": "normal"
}'

NOTIFY_CMD="curl -s -X POST $BASE_URL/notify -H \"Content-Type: application/json\" -d '$NOTIFY_DATA'"
print_curl "$NOTIFY_CMD"
NOTIFY_RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d "$NOTIFY_DATA")

print_json "$NOTIFY_RESPONSE"

if echo "$NOTIFY_RESPONSE" | grep -q '"success".*true'; then
    print_success "發送基本通知成功"
else
    print_fail "發送基本通知失敗"
    echo "$NOTIFY_RESPONSE"
fi
echo ""

# 測試 2: 發送帶優先級的通知
print_test 5 "發送高優先級通知"
HIGH_PRIORITY_DATA='{
  "notifyKey": "'$TEST_NOTIFY_KEY'",
  "message": "External API Test - High Priority Alert 🚨",
  "messageType": "text",
  "priority": "high"
}'

HIGH_CMD="curl -s -X POST $BASE_URL/notify -H \"Content-Type: application/json\" -d '$HIGH_PRIORITY_DATA'"
print_curl "$HIGH_CMD"
HIGH_RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d "$HIGH_PRIORITY_DATA")

print_json "$HIGH_RESPONSE"

if echo "$HIGH_RESPONSE" | grep -q '"success".*true'; then
    print_success "發送高優先級通知成功"
else
    print_fail "發送高優先級通知失敗"
    echo "$HIGH_RESPONSE"
fi
echo ""

# 測試 3: 發送帶 Metadata 的通知
print_test 5 "發送帶 Metadata 的通知"
METADATA_DATA='{
  "notifyKey": "'$TEST_NOTIFY_KEY'",
  "message": "External API Test - Message with Metadata",
  "messageType": "text",
  "priority": "normal",
  "metadata": {
    "source": "external_api_test",
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
    "environment": "test",
    "version": "1.0"
  }
}'

METADATA_CMD="curl -s -X POST $BASE_URL/notify -H \"Content-Type: application/json\" -d '$METADATA_DATA'"
print_curl "$METADATA_CMD"
METADATA_RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d "$METADATA_DATA")

print_json "$METADATA_RESPONSE"

if echo "$METADATA_RESPONSE" | grep -q '"success".*true'; then
    print_success "發送帶 Metadata 通知成功"
else
    print_fail "發送帶 Metadata 通知失敗"
    echo "$METADATA_RESPONSE"
fi
echo ""

# 測試 4: 獲取 Project Destinations
print_test 5 "獲取 Project 的 Destinations"
DEST_CMD="curl -s -X GET $BASE_URL/destinations/$TEST_NOTIFY_KEY"
print_curl "$DEST_CMD"
DEST_RESPONSE=$(curl -s -X GET "$BASE_URL/destinations/$TEST_NOTIFY_KEY")

print_json "$DEST_RESPONSE"

if echo "$DEST_RESPONSE" | grep -q '"success".*true'; then
    DEST_COUNT=$(echo "$DEST_RESPONSE" | python3 -c "import sys, json; print(len(json.load(sys.stdin).get('destinations', [])))" 2>/dev/null || echo "0")
    print_success "獲取 Destinations 成功，共 $DEST_COUNT 個目的地"
else
    print_fail "獲取 Destinations 失敗"
    echo "$DEST_RESPONSE"
fi
echo ""

# 測試 5: 測試無效的 NotifyKey（應該失敗）
print_test 5 "測試無效的 NotifyKey（預期失敗）"
INVALID_DATA='{
  "notifyKey": "invalid-notify-key-12345",
  "message": "This should fail",
  "messageType": "text"
}'

INVALID_CMD="curl -s -X POST $BASE_URL/notify -H \"Content-Type: application/json\" -d '$INVALID_DATA'"
print_curl "$INVALID_CMD"
INVALID_RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d "$INVALID_DATA")

print_json "$INVALID_RESPONSE"

if echo "$INVALID_RESPONSE" | grep -q '"success".*false'; then
    ERROR_MSG=$(echo "$INVALID_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('error', 'N/A')[:50])" 2>/dev/null)
    print_success "正確返回錯誤: $ERROR_MSG"
else
    print_fail "應該返回錯誤但沒有"
    echo "$INVALID_RESPONSE"
fi
echo ""

# 測試總結
echo -e "${BLUE}=========================================="
echo "測試總結"
echo "==========================================${NC}"
echo -e "總測試數: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "通過: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失敗: ${RED}$FAILED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✓ 所有測試通過！${NC}"
    echo ""
    echo -e "${YELLOW}測試的端點:${NC}"
    echo "  ✓ POST /api/v1/notify - 發送通知"
    echo "  ✓ GET  /api/v1/destinations/{notifyKey} - 獲取目的地"
    echo ""
    echo -e "${YELLOW}測試場景:${NC}"
    echo "  ✓ 基本文字通知"
    echo "  ✓ 高優先級通知"
    echo "  ✓ 帶 Metadata 的通知"
    echo "  ✓ 獲取 Project Destinations"
    echo "  ✓ 無效 NotifyKey 錯誤處理"
    echo ""
    echo -e "${YELLOW}通知發送到:${NC}"
    echo "  📱 Personal Chat (1 target)"
    echo "  👥 Group Chat (1 target)"
    echo "  📢 Channel (1 target)"
    echo "  Total: 3 個 Teams 目標"
    exit 0
else
    echo -e "${RED}✗ 有測試失敗${NC}"
    exit 1
fi

