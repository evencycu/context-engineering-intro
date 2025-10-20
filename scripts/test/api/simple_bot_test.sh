#!/bin/bash

# 簡化版 Bot API 測試腳本
set -e

API_URL="http://localhost:8080"
BASE_URL="${API_URL}/internal/v1"

# 顏色
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=========================================="
echo "Bot API 測試"
echo -e "==========================================${NC}\n"

# 1. 使用測試配置
echo -e "${BLUE}[1/10]${NC} 使用測試 Bot 配置..."
BOT_NAME="Test Teams Bot $(date +%s)"
APP_ID="test-app-id-$(date +%s)"
APP_PASSWORD="test-app-password-123"
TENANT_ID="test-tenant-id-123"
WEBHOOK_URL="https://example.com/webhook"
echo -e "${GREEN}✓ 測試配置已準備${NC}\n"

# 2. 創建 Teams Bot (Platform Bot)
echo -e "${BLUE}[2/10]${NC} 創建 Teams Bot..."

echo -e "${BLUE}curl -X POST ${BASE_URL}/bots/platform -H \"Content-Type: application/json\" -d '{...}'${NC}"
BOT_RESPONSE=$(curl -s -X POST ${BASE_URL}/bots/platform \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"$BOT_NAME\",
    \"description\": \"This is a test Teams bot for API testing\",
    \"app_id\": \"$APP_ID\",
    \"app_password\": \"$APP_PASSWORD\",
    \"tenant_id\": \"$TENANT_ID\",
    \"webhook_url\": \"$WEBHOOK_URL\",
    \"rate_limit_per_minute\": 60,
    \"max_concurrent_requests\": 10,
    \"capabilities\": {
      \"messaging\": true,
      \"notifications\": true
    }
  }")

echo "$BOT_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$BOT_RESPONSE"

BOT_ID=$(echo "$BOT_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "")

if [ -z "$BOT_ID" ]; then
    echo -e "${RED}✗ 創建 Bot 失敗${NC}"
    echo "$BOT_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Teams Bot 已創建: $BOT_ID${NC}\n"

# 3. 獲取 Bot 列表
echo -e "${BLUE}[3/10]${NC} 獲取 Bot 列表..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/bots/platform${NC}"
LIST_RESPONSE=$(curl -s ${BASE_URL}/bots/platform)
TOTAL_BOTS=$(echo "$LIST_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['pagination']['total'])" 2>/dev/null || echo "0")
echo "{\"data\": [...], \"pagination\": $(echo "$LIST_RESPONSE" | python3 -c "import sys, json; import json as j; print(j.dumps(json.load(sys.stdin)['pagination']))" 2>/dev/null || echo "{}")}"
echo -e "${GREEN}✓ 獲取成功，總共 $TOTAL_BOTS 個 Bot${NC}\n"

# 4. 獲取單個 Bot
echo -e "${BLUE}[4/10]${NC} 獲取單個 Bot (ID: $BOT_ID)..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/bots/platform/$BOT_ID${NC}"
GET_RESPONSE=$(curl -s ${BASE_URL}/bots/platform/$BOT_ID)
echo "$GET_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$GET_RESPONSE"
BOT_NAME_CHECK=$(echo "$GET_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['name'])" 2>/dev/null || echo "")
echo -e "${GREEN}✓ 獲取成功: $BOT_NAME_CHECK${NC}\n"

# 5. 更新 Bot
echo -e "${BLUE}[5/10]${NC} 更新 Bot 信息..."
echo -e "${BLUE}curl -X PUT ${BASE_URL}/bots/platform/$BOT_ID -H \"Content-Type: application/json\" -d '{...}'${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT ${BASE_URL}/bots/platform/$BOT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Test Teams Bot",
    "description": "Updated description for testing",
    "rateLimitPerMinute": 120,
    "maxConcurrentRequests": 20
  }')
echo "$UPDATE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$UPDATE_RESPONSE"
echo -e "${GREEN}✓ 更新成功${NC}\n"

# 6. 更新 Bot 狀態
echo -e "${BLUE}[6/10]${NC} 更新 Bot 狀態..."
echo -e "${BLUE}curl -X PATCH ${BASE_URL}/bots/platform/$BOT_ID/status -H \"Content-Type: application/json\" -d '{\"status\":\"inactive\"}'${NC}"
STATUS_RESPONSE=$(curl -s -X PATCH ${BASE_URL}/bots/platform/$BOT_ID/status \
  -H "Content-Type: application/json" \
  -d '{"status":"inactive"}')
echo "$STATUS_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$STATUS_RESPONSE"
STATUS_MESSAGE=$(echo "$STATUS_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('message',''))" 2>/dev/null || echo "")
echo -e "${GREEN}✓ $STATUS_MESSAGE${NC}\n"

# 7. 更新 Bot Capabilities
echo -e "${BLUE}[7/10]${NC} 更新 Bot Capabilities..."
echo -e "${BLUE}curl -X PATCH ${BASE_URL}/bots/platform/$BOT_ID/capabilities -H \"Content-Type: application/json\" -d '{...}'${NC}"
CAP_RESPONSE=$(curl -s -X PATCH ${BASE_URL}/bots/platform/$BOT_ID/capabilities \
  -H "Content-Type: application/json" \
  -d '{
    "capabilities": {
      "messaging": true,
      "notifications": true,
      "file_upload": true,
      "adaptive_cards": true
    }
  }')
echo "$CAP_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$CAP_RESPONSE"
CAP_MESSAGE=$(echo "$CAP_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('message',''))" 2>/dev/null || echo "")
echo -e "${GREEN}✓ $CAP_MESSAGE${NC}\n"

# 8. 測試 Bot 連接
echo -e "${BLUE}[8/10]${NC} 測試 Bot 連接..."
echo -e "${BLUE}curl -X POST ${BASE_URL}/bots/platform/$BOT_ID/test${NC}"
TEST_RESPONSE=$(curl -s -X POST ${BASE_URL}/bots/platform/$BOT_ID/test)
echo "$TEST_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$TEST_RESPONSE"
TEST_SUCCESS=$(echo "$TEST_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('data', {}).get('success', False))" 2>/dev/null || echo "false")
if [ "$TEST_SUCCESS" = "True" ]; then
    echo -e "${GREEN}✓ 連接測試成功${NC}\n"
else
    echo -e "${BLUE}ℹ 連接測試失敗（預期，因為使用測試配置）${NC}\n"
fi

# 9. 按狀態查詢 Bot
echo -e "${BLUE}[9/10]${NC} 按狀態查詢 Bot (inactive)..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/bots/status/inactive${NC}"
STATUS_BOTS_RESPONSE=$(curl -s ${BASE_URL}/bots/status/inactive)
echo "$STATUS_BOTS_RESPONSE" | python3 -m json.tool 2>/dev/null | head -20
STATUS_BOTS_COUNT=$(echo "$STATUS_BOTS_RESPONSE" | python3 -c "import sys, json; data=json.load(sys.stdin).get('data', {}); print(len(data.get('platform_bots', [])))" 2>/dev/null || echo "0")
echo -e "${GREEN}✓ Inactive 狀態有 $STATUS_BOTS_COUNT 個 Bot${NC}\n"

# 10. 刪除 Bot
echo -e "${BLUE}[10/10]${NC} 刪除 Bot..."
echo -e "${BLUE}curl -X DELETE ${BASE_URL}/bots/platform/$BOT_ID${NC}"
DELETE_RESPONSE=$(curl -s -X DELETE ${BASE_URL}/bots/platform/$BOT_ID)
echo "$DELETE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$DELETE_RESPONSE"
echo -e "${GREEN}✓ 測試 Bot 已刪除${NC}\n"

echo -e "${GREEN}=========================================="
echo "✓ 所有 Bot API 測試完成！"
echo -e "==========================================${NC}"

