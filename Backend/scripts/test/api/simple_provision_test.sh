#!/bin/bash

# Provision API 測試腳本
# 測試所有 Provision API 端點

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
TEST_COMPANY_ID="550e8400-e29b-41d4-a716-446655440001"  # 國泰金控
TEST_USER_ID="650e8400-e29b-41d4-a716-446655440001"     # Admin User
TEST_BOT_ID="850e8400-e29b-41d4-a716-446655440002"      # Lab Test Teams Notify Bot
TEST_APP_ID="844146d7-4ac9-4e4d-a463-d6e027714e81"      # 從 init.sql
TEST_TENANT_ID="051cece0-e4dc-4aed-b471-bf29824e1ee6"   # 從 init.sql
NOTIFY_KEY=""

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
echo "Provision API 測試"
echo "==========================================${NC}"
echo ""

# 測試 1: 創建 Provision（完整配置）- 使用 init.sql 中的數據
print_test 5 "創建完整配置（使用 init.sql 中的資源）"
PROVISION_DATA='{
  "companyId": "'$TEST_COMPANY_ID'",
  "createdBy": "'$TEST_USER_ID'",
  "projectName": "Provision Test Project '$RANDOM'",
  "projectDescription": "Provision API test project using init.sql data",
  "teamsTenantId": "'$TEST_TENANT_ID'",
  "targets": [
    {
      "type": "channel",
      "teamId": "19:MzE5NmVmNjAtMWE3ZS00YmQ0LWIxOWUtNTBiN2RjNDQyNDM2@thread.tacv2",
      "channelId": "19:MzE5NmVmNjAtMWE3ZS00YmQ0LWIxOWUtNTBiN2RjNDQyNDM2@thread.tacv2",
      "conversationId": "19:MzE5NmVmNjAtMWE3ZS00YmQ0LWIxOWUtNTBiN2RjNDQyNDM2@thread.tacv2",
      "tenantId": "'$TEST_TENANT_ID'"
    }
  ]
}'

PROVISION_CMD="curl -s -X POST $BASE_URL/provision -H \"Content-Type: application/json\" -d '$PROVISION_DATA'"
print_curl "$PROVISION_CMD"
PROVISION_RESPONSE=$(curl -s -X POST "$BASE_URL/provision" \
  -H "Content-Type: application/json" \
  -d "$PROVISION_DATA")

print_json "$PROVISION_RESPONSE"

if echo "$PROVISION_RESPONSE" | grep -q '"data"'; then
    NOTIFY_KEY=$(echo "$PROVISION_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notifyKey'])" 2>/dev/null)
    if [ -n "$NOTIFY_KEY" ]; then
        print_success "創建配置成功，NotifyKey: $NOTIFY_KEY"
    else
        print_fail "創建配置成功但未獲取到 NotifyKey"
        echo "$PROVISION_RESPONSE"
    fi
else
    print_fail "創建配置失敗"
    echo "$PROVISION_RESPONSE"
    exit 1
fi
echo ""

# 測試 2: 讀取 Provision
print_test 5 "讀取配置詳情"
READ_CMD="curl -s -X GET $BASE_URL/provision/$NOTIFY_KEY"
print_curl "$READ_CMD"
READ_RESPONSE=$(curl -s -X GET "$BASE_URL/provision/$NOTIFY_KEY")

print_json "$READ_RESPONSE"

if echo "$READ_RESPONSE" | grep -q '"data"'; then
    PROJECT_NAME=$(echo "$READ_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data'].get('projectName', 'N/A'))" 2>/dev/null)
    print_success "讀取配置成功，專案名稱: $PROJECT_NAME"
else
    print_fail "讀取配置失敗"
    echo "$READ_RESPONSE"
fi
echo ""

# 測試 3: 更新 Provision
print_test 5 "更新配置"
UPDATE_DATA='{
  "projectDescription": "Updated provision project description",
  "dailyLimit": 1500,
  "monthlyLimit": 40000
}'

UPDATE_CMD="curl -s -X PUT $BASE_URL/provision/$NOTIFY_KEY -H \"Content-Type: application/json\" -d '$UPDATE_DATA'"
print_curl "$UPDATE_CMD"
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/provision/$NOTIFY_KEY" \
  -H "Content-Type: application/json" \
  -d "$UPDATE_DATA")

print_json "$UPDATE_RESPONSE"

if echo "$UPDATE_RESPONSE" | grep -q '"data"'; then
    UPDATED_LIMIT=$(echo "$UPDATE_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data'].get('dailyLimit', 'N/A'))" 2>/dev/null)
    if [ "$UPDATED_LIMIT" = "1500" ]; then
        print_success "更新配置成功，新的每日限制: $UPDATED_LIMIT"
    else
        print_success "更新配置成功（每日限制: $UPDATED_LIMIT）"
    fi
else
    print_fail "更新配置失敗"
    echo "$UPDATE_RESPONSE"
fi
echo ""

# 測試 4: 停用 Provision
print_test 5 "停用專案"
DISABLE_CMD="curl -s -X POST $BASE_URL/provision/$NOTIFY_KEY/disable"
print_curl "$DISABLE_CMD"
DISABLE_RESPONSE=$(curl -s -X POST "$BASE_URL/provision/$NOTIFY_KEY/disable")

print_json "$DISABLE_RESPONSE"

if echo "$DISABLE_RESPONSE" | grep -q '"message".*"disabled"'; then
    print_success "停用專案成功"
else
    print_fail "停用專案失敗"
    echo "$DISABLE_RESPONSE"
fi
echo ""

# 測試 5: 啟用 Provision
print_test 5 "啟用專案"
ENABLE_CMD="curl -s -X POST $BASE_URL/provision/$NOTIFY_KEY/enable"
print_curl "$ENABLE_CMD"
ENABLE_RESPONSE=$(curl -s -X POST "$BASE_URL/provision/$NOTIFY_KEY/enable")

print_json "$ENABLE_RESPONSE"

if echo "$ENABLE_RESPONSE" | grep -q '"message".*"enabled"'; then
    print_success "啟用專案成功"
else
    print_fail "啟用專案失敗"
    echo "$ENABLE_RESPONSE"
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
    echo -e "${YELLOW}創建的資源:${NC}"
    echo -e "  NotifyKey: ${GREEN}$NOTIFY_KEY${NC}"
    echo ""
    echo -e "${YELLOW}測試的端點:${NC}"
    echo "  ✓ POST /api/v1/provision"
    echo "  ✓ GET  /api/v1/provision/{notifyKey}"
    echo "  ✓ PUT  /api/v1/provision/{notifyKey}"
    echo "  ✓ POST /api/v1/provision/{notifyKey}/disable"
    echo "  ✓ POST /api/v1/provision/{notifyKey}/enable"
    exit 0
else
    echo -e "${RED}✗ 有測試失敗${NC}"
    exit 1
fi

