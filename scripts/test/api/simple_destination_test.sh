#!/bin/bash

# Simple Destination API Test Script
# 測試 Destination API 的基本功能

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Base URL
BASE_URL="http://localhost:8080/internal/v1"

# Fixed IDs (應該在資料庫中已存在)
COMPANY_ID="550e8400-e29b-41d4-a716-446655440001"
USER_ID=""
PROJECT_ID=""
BOT_ID=""
DESTINATION_ID=""

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to print test result
print_result() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASSED${NC}: $2"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗ FAILED${NC}: $2"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# Function to print section header
print_header() {
    echo ""
    echo -e "${BLUE}================================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================================${NC}"
}

# Function to print separator
print_separator() {
    echo -e "${YELLOW}------------------------------------------------${NC}"
}

echo -e "${BLUE}開始測試 Destinations API${NC}"
echo -e "${BLUE}Base URL: $BASE_URL${NC}"
echo ""

# ==================== 準備測試資料 ====================
print_header "準備測試資料"

# 1. 創建測試用戶
print_separator
echo -e "${YELLOW}創建測試用戶...${NC}"
echo -e "${BLUE}curl -X POST \"$BASE_URL/users\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "'"$COMPANY_ID"'",
    "username": "dest_test_user_'"$(date +%s)"'",
    "email": "dest_test_'"$(date +%s)"'@example.com",
    "display_name": "Destination Test User",
    "role": "admin"
  }')
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

USER_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
if [ -z "$USER_ID" ]; then
    echo -e "${RED}❌ 無法創建測試用戶${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 測試用戶 ID: $USER_ID${NC}"

# 2. 創建測試專案
print_separator
echo -e "${YELLOW}創建測試專案...${NC}"
echo -e "${BLUE}curl -X POST \"$BASE_URL/projects\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"
NOTIFY_KEY="dest-test-$(date +%s)"
RESPONSE=$(curl -s -X POST "$BASE_URL/projects" \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "'"$COMPANY_ID"'",
    "notify_key": "'"$NOTIFY_KEY"'",
    "description": "Test project for destination API testing",
    "daily_limit": 1000,
    "monthly_limit": 30000,
    "priority": "normal",
    "created_by": "'"$USER_ID"'"
  }')
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

PROJECT_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
if [ -z "$PROJECT_ID" ]; then
    echo -e "${RED}❌ 無法創建測試專案${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 測試專案 ID: $PROJECT_ID${NC}"

# 3. 創建測試 Bot
print_separator
echo -e "${YELLOW}創建測試 Bot...${NC}"
APP_ID="dest-bot-$(date +%s)"
echo -e "${BLUE}curl -X POST \"$BASE_URL/bots/platform\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/bots/platform" \
  -H "Content-Type: application/json" \
  -d '{
    "app_id": "'"$APP_ID"'",
    "app_password": "test-password-123",
    "tenant_id": "test-tenant-123",
    "webhook_url": "https://webhook.example.com/test",
    "name": "Destination Test Bot",
    "description": "Test bot for destination API testing",
    "capabilities": {
      "messaging": true,
      "channel": true
    },
    "rate_limit_per_minute": 60,
    "max_concurrent_requests": 10
  }')
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

BOT_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
if [ -z "$BOT_ID" ]; then
    echo -e "${RED}❌ 無法創建測試 Bot${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 測試 Bot ID: $BOT_ID${NC}"

# ==================== 測試 Destination APIs ====================
print_header "測試 1: 創建 Destination"
echo -e "${BLUE}curl -X POST \"$BASE_URL/destinations\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/destinations" \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "'"$PROJECT_ID"'",
    "name": "Test Destination",
    "description": "This is a test destination for API testing purposes",
    "teams_tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
    "targets": [
      {
        "type": "channel",
        "conversation_id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
        "display_name": "Test Channel",
        "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
      },
      {
        "type": "personal",
        "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
        "display_name": "Test User",
        "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
      }
    ],
    "bot_id": "'"$BOT_ID"'",
    "created_by": "'"$USER_ID"'"
  }')
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

DESTINATION_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
if [ -n "$DESTINATION_ID" ]; then
    print_result 0 "創建 Destination 成功 (ID: $DESTINATION_ID)"
else
    print_result 1 "創建 Destination 失敗"
    echo -e "${RED}錯誤詳情: $RESPONSE${NC}"
fi

# ==================== 測試 2: 列出所有 Destinations ====================
print_header "測試 2: 列出所有 Destinations"
echo -e "${BLUE}curl -X GET \"$BASE_URL/destinations?limit=10&offset=0\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/destinations?limit=10&offset=0" \
  -H "Content-Type: application/json")
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(len(json.load(sys.stdin).get('data', [])))" 2>/dev/null)
if [ -n "$COUNT" ] && [ "$COUNT" -gt 0 ]; then
    print_result 0 "列出 Destinations 成功 (找到 $COUNT 筆資料)"
else
    print_result 1 "列出 Destinations 失敗"
fi

# ==================== 測試 3: 取得單一 Destination ====================
print_header "測試 3: 取得單一 Destination"
if [ -n "$DESTINATION_ID" ]; then
    echo -e "${BLUE}curl -X GET \"$BASE_URL/destinations/$DESTINATION_ID\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
    RESPONSE=$(curl -s -X GET "$BASE_URL/destinations/$DESTINATION_ID" \
      -H "Content-Type: application/json")
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    FETCHED_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
    if [ "$FETCHED_ID" = "$DESTINATION_ID" ]; then
        print_result 0 "取得 Destination 成功"
    else
        print_result 1 "取得 Destination 失敗"
    fi
else
    print_result 1 "無法測試 (沒有 Destination ID)"
fi

# ==================== 測試 4: 更新 Destination 基本資訊 ====================
print_header "測試 4: 更新 Destination 基本資訊"
if [ -n "$DESTINATION_ID" ]; then
    echo -e "${BLUE}curl -X PUT \"$BASE_URL/destinations/$DESTINATION_ID\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
    echo -e "${BLUE}  -d '{...}'${NC}"
    RESPONSE=$(curl -s -X PUT "$BASE_URL/destinations/$DESTINATION_ID" \
      -H "Content-Type: application/json" \
      -d '{
        "name": "Updated Test Destination",
        "description": "This destination has been updated for testing purposes"
      }')
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    UPDATED_NAME=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['name'])" 2>/dev/null)
    if [ "$UPDATED_NAME" = "Updated Test Destination" ]; then
        print_result 0 "更新 Destination 成功"
    else
        print_result 1 "更新 Destination 失敗"
    fi
else
    print_result 1 "無法測試 (沒有 Destination ID)"
fi

# ==================== 測試 5: 更新 Destination Targets ====================
print_header "測試 5: 更新 Destination Targets"
if [ -n "$DESTINATION_ID" ]; then
    echo -e "${BLUE}curl -X PATCH \"$BASE_URL/destinations/$DESTINATION_ID/targets\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
    echo -e "${BLUE}  -d '{...}'${NC}"
    RESPONSE=$(curl -s -X PATCH "$BASE_URL/destinations/$DESTINATION_ID/targets" \
      -H "Content-Type: application/json" \
      -d '{
        "targets": [
          {
            "type": "channel",
            "conversation_id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
            "display_name": "Updated Test Channel",
            "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
          },
          {
            "type": "groupchat",
            "conversation_id": "19:f26a8d8a235f430db87a404491cd2ffc@thread.v2",
            "display_name": "Test Group Chat",
            "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
          }
        ]
      }')
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    MESSAGE=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('message', ''))" 2>/dev/null)
    if [[ "$MESSAGE" == *"success"* ]]; then
        print_result 0 "更新 Targets 成功"
    else
        print_result 1 "更新 Targets 失敗"
    fi
else
    print_result 1 "無法測試 (沒有 Destination ID)"
fi

# ==================== 測試 6: 驗證 Targets ====================
print_header "測試 6: 驗證 Targets"
if [ -n "$DESTINATION_ID" ]; then
    echo -e "${BLUE}curl -X POST \"$BASE_URL/destinations/$DESTINATION_ID/validate\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
    echo -e "${BLUE}  -d '{...}'${NC}"
    RESPONSE=$(curl -s -X POST "$BASE_URL/destinations/$DESTINATION_ID/validate" \
      -H "Content-Type: application/json" \
      -d '{
        "targets": [
          {
            "type": "personal",
            "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
            "display_name": "Validate User",
            "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
          }
        ]
      }')
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    IS_VALID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data'].get('valid', False))" 2>/dev/null)
    if [ "$IS_VALID" = "True" ]; then
        print_result 0 "驗證 Targets 成功"
    else
        print_result 1 "驗證 Targets 失敗"
    fi
else
    print_result 1 "無法測試 (沒有 Destination ID)"
fi

# ==================== 測試 7: 按專案查詢 Destinations ====================
print_header "測試 7: 按專案查詢 Destinations"
if [ -n "$PROJECT_ID" ]; then
    echo -e "${BLUE}curl -X GET \"$BASE_URL/destinations/project/$PROJECT_ID\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
    RESPONSE=$(curl -s -X GET "$BASE_URL/destinations/project/$PROJECT_ID" \
      -H "Content-Type: application/json")
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; data = json.load(sys.stdin).get('data'); print(len(data) if data is not None else 0)" 2>/dev/null)
    if [ -n "$COUNT" ]; then
        print_result 0 "按專案查詢 Destinations 成功 (找到 $COUNT 筆資料)"
    else
        print_result 1 "按專案查詢 Destinations 失敗"
    fi
else
    print_result 1 "無法測試 (沒有 Project ID)"
fi

# ==================== 測試 8: 按 Bot 查詢 Destinations ====================
print_header "測試 8: 按 Bot 查詢 Destinations"
if [ -n "$BOT_ID" ]; then
    echo -e "${BLUE}curl -X GET \"$BASE_URL/destinations/bot/$BOT_ID\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
    RESPONSE=$(curl -s -X GET "$BASE_URL/destinations/bot/$BOT_ID" \
      -H "Content-Type: application/json")
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; data = json.load(sys.stdin).get('data'); print(len(data) if data is not None else 0)" 2>/dev/null)
    if [ -n "$COUNT" ]; then
        print_result 0 "按 Bot 查詢 Destinations 成功 (找到 $COUNT 筆資料)"
    else
        print_result 1 "按 Bot 查詢 Destinations 失敗"
    fi
else
    print_result 1 "無法測試 (沒有 Bot ID)"
fi

# ==================== 測試 9: 搜尋 Destinations ====================
print_header "測試 9: 搜尋 Destinations"
echo -e "${BLUE}curl -X GET \"$BASE_URL/destinations/search?target_type=channel&target_id=updated-channel-456\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/destinations/search?target_type=channel&target_id=updated-channel-456" \
  -H "Content-Type: application/json")
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('count', 0))" 2>/dev/null)
if [ -n "$COUNT" ] && [ "$COUNT" -ge 0 ]; then
    print_result 0 "搜尋 Destinations 成功 (找到 $COUNT 筆資料)"
else
    print_result 1 "搜尋 Destinations 失敗"
fi

# ==================== 測試 10: 刪除 Destination ====================
print_header "測試 10: 刪除 Destination"
if [ -n "$DESTINATION_ID" ]; then
    echo -e "${BLUE}curl -X DELETE \"$BASE_URL/destinations/$DESTINATION_ID\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
    RESPONSE=$(curl -s -X DELETE "$BASE_URL/destinations/$DESTINATION_ID" \
      -H "Content-Type: application/json")
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    MESSAGE=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('message', ''))" 2>/dev/null)
    if [[ "$MESSAGE" == *"success"* ]]; then
        print_result 0 "刪除 Destination 成功"
    else
        print_result 1 "刪除 Destination 失敗"
    fi
else
    print_result 1 "無法測試 (沒有 Destination ID)"
fi

# ==================== 清理測試資料 ====================
print_header "清理測試資料"

# 刪除測試 Bot
if [ -n "$BOT_ID" ]; then
    echo -e "${YELLOW}刪除測試 Bot: $BOT_ID${NC}"
    curl -s -X DELETE "$BASE_URL/bots/platform/$BOT_ID" > /dev/null
    echo -e "${GREEN}✓ Bot 已刪除${NC}"
fi

# 刪除測試專案
if [ -n "$PROJECT_ID" ]; then
    echo -e "${YELLOW}刪除測試專案: $PROJECT_ID${NC}"
    curl -s -X DELETE "$BASE_URL/projects/$PROJECT_ID" > /dev/null
    echo -e "${GREEN}✓ 專案已刪除${NC}"
fi

# 刪除測試用戶
if [ -n "$USER_ID" ]; then
    echo -e "${YELLOW}刪除測試用戶: $USER_ID${NC}"
    curl -s -X DELETE "$BASE_URL/users/$USER_ID" > /dev/null
    echo -e "${GREEN}✓ 用戶已刪除${NC}"
fi

# ==================== 測試總結 ====================
print_header "測試總結"
echo -e "總測試數: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "通過: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失敗: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有測試通過！${NC}"
    exit 0
else
    echo -e "\n${RED}❌ 有測試失敗${NC}"
    exit 1
fi

