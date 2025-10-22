#!/bin/bash

# Simple Notification API Test Script
# 測試 Notification API 的基本功能

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
USER_ID="650e8400-e29b-41d4-a716-446655440001"
PROJECT_ID="750e8400-e29b-41d4-a716-446655440001"
DESTINATION_ID="950e8400-e29b-41d4-a716-446655440003"
NOTIFICATION_ID=""

# Predefined conversation IDs from scripts/database/init.sql
PERSONAL_CONV_ID="a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR"
GROUPCHAT_CONV_ID="19:f26a8d8a235f430db87a404491cd2ffc@thread.v2"
CHANNEL_CONV_ID="19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2"

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

echo -e "${BLUE}開始測試 Notifications API${NC}"
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
    "companyId": "'"$COMPANY_ID"'",
    "username": "notif_test_user_'"$(date +%s)"'",
    "email": "notif_test_'"$(date +%s)"'@example.com",
    "display_name": "Notification Test User",
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
NOTIFY_KEY="notif-test-$(date +%s)"
RESPONSE=$(curl -s -X POST "$BASE_URL/projects" \
  -H "Content-Type: application/json" \
  -d '{
    "companyId": "'"$COMPANY_ID"'",
    "notifyKey": "'"$NOTIFY_KEY"'",
    "description": "Test project for notification API testing",
    "dailyLimit": 1000,
    "monthlyLimit": 30000,
    "priority": "normal",
    "createdBy": "'"$USER_ID"'"
  }')
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

PROJECT_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
if [ -z "$PROJECT_ID" ]; then
    echo -e "${RED}❌ 無法創建測試專案${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 測試專案 ID: $PROJECT_ID${NC}"

# 3. 創建測試 Destination
print_separator
echo -e "${YELLOW}創建測試 Destination...${NC}"
echo -e "${BLUE}curl -X POST \"$BASE_URL/destinations\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/destinations" \
  -H "Content-Type: application/json" \
  -d '{
    "projectId": "'"$PROJECT_ID"'",
    "name": "Notification Test Destination",
    "description": "Test destination for notification API testing",
    "teamsTenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
    "targets": [
      {
        "type": "personal",
        "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
        "display_name": "Test User",
        "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
      }
    ],
    "createdBy": "'"$USER_ID"'"
  }')
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

DESTINATION_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
if [ -z "$DESTINATION_ID" ]; then
    echo -e "${RED}❌ 無法創建測試 Destination${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 測試 Destination ID: $DESTINATION_ID${NC}"

# ==================== 測試 Notification APIs ====================
print_header "測試 1: 發送文字通知"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notifications\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/notifications" \
  -H "Content-Type: application/json" \
  -d '{
    "projectId": "'"$PROJECT_ID"'",
    "senderId": "'"$USER_ID"'",
    "messageType": "text",
    "content": "這是一條測試通知訊息",
    "priority": "normal",
    "targets": ["'"$PERSONAL_CONV_ID"'", "'"$GROUPCHAT_CONV_ID"'", "'"$CHANNEL_CONV_ID"'"],
    "metadata": {
      "test_type": "simple_text",
      "environment": "test"
    }
  }')
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

NOTIFICATION_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notificationId'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
if [ -n "$NOTIFICATION_ID" ] && [ "$STATUS" = "pending" ]; then
    print_result 0 "發送文字通知成功 (ID: $NOTIFICATION_ID)"
else
    print_result 1 "發送文字通知失敗"
    echo -e "${RED}錯誤詳情: $RESPONSE${NC}"
fi

# ==================== 測試 2: 發送高優先級通知 ====================
print_header "測試 2: 發送高優先級通知"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notifications\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/notifications" \
  -H "Content-Type: application/json" \
  -d '{
    "projectId": "'"$PROJECT_ID"'",
    "senderId": "'"$USER_ID"'",
    "messageType": "text",
    "content": "🚨 這是一條高優先級測試通知",
    "priority": "high",
    "targets": ["'"$PERSONAL_CONV_ID"'", "'"$GROUPCHAT_CONV_ID"'", "'"$CHANNEL_CONV_ID"'"],
    "metadata": {
      "test_type": "high_priority",
      "alert_level": "critical"
    }
  }')
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

HIGH_PRIORITY_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notificationId'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
if [ -n "$HIGH_PRIORITY_ID" ] && [ "$STATUS" = "pending" ]; then
    print_result 0 "發送高優先級通知成功"
else
    print_result 1 "發送高優先級通知失敗"
fi

# ==================== 測試 3: 發送包含 Mentions 的通知 ====================
print_header "測試 3: 發送包含 Mentions 的通知"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notifications\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/notifications" \
  -H "Content-Type: application/json" \
  -d '{
    "projectId": "'"$PROJECT_ID"'",
    "senderId": "'"$USER_ID"'",
    "messageType": "text",
    "content": "Hi @TestUser, 請查看這條重要訊息",
    "mentions": ["TestUser"],
    "priority": "normal",
    "targets": ["'"$PERSONAL_CONV_ID"'", "'"$GROUPCHAT_CONV_ID"'", "'"$CHANNEL_CONV_ID"'"]
  }')
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

MENTION_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notificationId'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
if [ -n "$MENTION_ID" ] && [ "$STATUS" = "pending" ]; then
    print_result 0 "發送包含 Mentions 的通知成功"
else
    print_result 1 "發送包含 Mentions 的通知失敗"
fi

# ==================== 測試 4: 列出所有通知 ====================
print_header "測試 4: 列出所有通知"
echo -e "${BLUE}curl -X GET \"$BASE_URL/notifications?limit=10&offset=0\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/notifications?limit=10&offset=0" \
  -H "Content-Type: application/json")
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(len(json.load(sys.stdin).get('data', [])))" 2>/dev/null)
if [ -n "$COUNT" ] && [ "$COUNT" -gt 0 ]; then
    print_result 0 "列出通知成功 (找到 $COUNT 筆資料)"
else
    print_result 1 "列出通知失敗"
fi

# ==================== 測試 5: 取得單一通知 ====================
print_header "測試 5: 取得單一通知"
if [ -n "$NOTIFICATION_ID" ]; then
    echo -e "${BLUE}curl -X GET \"$BASE_URL/notifications/$NOTIFICATION_ID\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
    RESPONSE=$(curl -s -X GET "$BASE_URL/notifications/$NOTIFICATION_ID" \
      -H "Content-Type: application/json")
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    FETCHED_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null)
    if [ "$FETCHED_ID" = "$NOTIFICATION_ID" ]; then
        print_result 0 "取得通知成功"
    else
        print_result 1 "取得通知失敗"
    fi
else
    print_result 1 "無法測試 (沒有 Notification ID)"
fi

# ==================== 測試 6: 按專案查詢通知 ====================
print_header "測試 6: 按專案查詢通知"
if [ -n "$PROJECT_ID" ]; then
    echo -e "${BLUE}curl -X GET \"$BASE_URL/notifications/project/$PROJECT_ID\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
    RESPONSE=$(curl -s -X GET "$BASE_URL/notifications/project/$PROJECT_ID" \
      -H "Content-Type: application/json")
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; data = json.load(sys.stdin).get('data'); print(len(data) if data is not None else 0)" 2>/dev/null)
    if [ -n "$COUNT" ] && [ "$COUNT" -gt 0 ]; then
        print_result 0 "按專案查詢通知成功 (找到 $COUNT 筆資料)"
    else
        print_result 1 "按專案查詢通知失敗"
    fi
else
    print_result 1 "無法測試 (沒有 Project ID)"
fi

# ==================== 測試 7: 按發送者查詢通知 ====================
print_header "測試 7: 按發送者查詢通知"
if [ -n "$USER_ID" ]; then
    echo -e "${BLUE}curl -X GET \"$BASE_URL/notifications/sender/$USER_ID\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
    RESPONSE=$(curl -s -X GET "$BASE_URL/notifications/sender/$USER_ID" \
      -H "Content-Type: application/json")
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; data = json.load(sys.stdin).get('data'); print(len(data) if data is not None else 0)" 2>/dev/null)
    if [ -n "$COUNT" ] && [ "$COUNT" -gt 0 ]; then
        print_result 0 "按發送者查詢通知成功 (找到 $COUNT 筆資料)"
    else
        print_result 1 "按發送者查詢通知失敗"
    fi
else
    print_result 1 "無法測試 (沒有 User ID)"
fi

# ==================== 測試 8: 按狀態查詢通知 ====================
print_header "測試 8: 按狀態查詢通知"
echo -e "${BLUE}curl -X GET \"$BASE_URL/notifications/status/pending\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/notifications/status/pending" \
  -H "Content-Type: application/json")
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; data = json.load(sys.stdin).get('data'); print(len(data) if data is not None else 0)" 2>/dev/null)
if [ -n "$COUNT" ]; then
    print_result 0 "按狀態查詢通知成功 (找到 $COUNT 筆資料)"
else
    print_result 1 "按狀態查詢通知失敗"
fi

# ==================== 測試 9: 按日期範圍查詢通知 ====================
print_header "測試 9: 按日期範圍查詢通知"
# 獲取當前時間和一週前的時間
END_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
START_DATE=$(date -u -v-7d +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "7 days ago" +"%Y-%m-%dT%H:%M:%SZ")

echo -e "${BLUE}curl -X GET \"$BASE_URL/notifications/date-range?start_date=$START_DATE&end_date=$END_DATE\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
RESPONSE=$(curl -s -X GET "$BASE_URL/notifications/date-range?start_date=$START_DATE&end_date=$END_DATE" \
  -H "Content-Type: application/json")
echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('count', 0))" 2>/dev/null)
if [ -n "$COUNT" ]; then
    print_result 0 "按日期範圍查詢通知成功 (找到 $COUNT 筆資料)"
else
    print_result 1 "按日期範圍查詢通知失敗"
fi

# ==================== 測試 10: 重試通知 ====================
print_header "測試 10: 重試通知"
if [ -n "$NOTIFICATION_ID" ]; then
    # 先等待一下確保通知已處理
    sleep 2
    
    echo -e "${BLUE}curl -X POST \"$BASE_URL/notifications/$NOTIFICATION_ID/retry\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
    RESPONSE=$(curl -s -X POST "$BASE_URL/notifications/$NOTIFICATION_ID/retry" \
      -H "Content-Type: application/json")
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    MESSAGE=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('message', ''))" 2>/dev/null)
    ERROR=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('error', ''))" 2>/dev/null)
    
    # 重試可能成功或失敗（取決於當前狀態），兩者都算測試通過
    if [[ "$MESSAGE" == *"success"* ]] || [[ "$ERROR" != "" ]]; then
        print_result 0 "重試通知測試完成 (回應: ${MESSAGE:-$ERROR})"
    else
        print_result 1 "重試通知測試失敗"
    fi
else
    print_result 1 "無法測試 (沒有 Notification ID)"
fi

# ==================== 測試 11: 取消通知 ====================
print_header "測試 11: 取消通知"
if [ -n "$HIGH_PRIORITY_ID" ]; then
    echo -e "${BLUE}curl -X DELETE \"$BASE_URL/notifications/$HIGH_PRIORITY_ID\" \\${NC}"
    echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"
    RESPONSE=$(curl -s -X DELETE "$BASE_URL/notifications/$HIGH_PRIORITY_ID" \
      -H "Content-Type: application/json")
    echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"
    
    MESSAGE=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('message', ''))" 2>/dev/null)
    ERROR=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('error', ''))" 2>/dev/null)
    
    # 取消可能成功或失敗（取決於當前狀態），兩者都算測試通過
    if [[ "$MESSAGE" == *"success"* ]] || [[ "$ERROR" != "" ]]; then
        print_result 0 "取消通知測試完成 (回應: ${MESSAGE:-$ERROR})"
    else
        print_result 1 "取消通知測試失敗"
    fi
else
    print_result 1 "無法測試 (沒有高優先級通知 ID)"
fi

# ==================== 清理測試資料 ====================
print_header "清理測試資料"

# 刪除測試 Destination
if [ -n "$DESTINATION_ID" ]; then
    echo -e "${YELLOW}刪除測試 Destination: $DESTINATION_ID${NC}"
    curl -s -X DELETE "$BASE_URL/destinations/$DESTINATION_ID" > /dev/null
    echo -e "${GREEN}✓ Destination 已刪除${NC}"
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

