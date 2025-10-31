#!/bin/bash

# Simple External Notify API Test Script
# 測試外部使用者透過 /api/v1/notify 發送通知

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Base URL for external API
BASE_URL="http://localhost:8080/api/v1"

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

echo -e "${BLUE}開始測試 External Notify API${NC}"
echo -e "${BLUE}Base URL: $BASE_URL${NC}"
echo ""

# ==================== 測試 1: 基本文字通知 ====================
print_header "測試 1: 基本文字通知"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notify\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "cfh-alert-gogo",
    "message": "這是一條外部測試通知訊息",
    "messageType": "text",
    "priority": "normal",
    "targets": ["all"],
    "metadata": {
      "source": "external_test",
      "environment": "test"
    }
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

NOTIFICATION_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notification_id'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
DESTINATIONS_COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['destinations_count'])" 2>/dev/null)

if [ -n "$NOTIFICATION_ID" ] && [ "$STATUS" = "pending" ] && [ "$DESTINATIONS_COUNT" -gt 0 ]; then
    print_result 0 "基本文字通知發送成功 (ID: $NOTIFICATION_ID, Destinations: $DESTINATIONS_COUNT)"
else
    print_result 1 "基本文字通知發送失敗"
    echo -e "${RED}錯誤詳情: $RESPONSE${NC}"
fi

# ==================== 測試 2: 高優先級通知 ====================
print_header "測試 2: 高優先級通知"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notify\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "cfh-alert-gogo",
    "message": "🚨 這是一條高優先級外部測試通知",
    "messageType": "text",
    "priority": "high",
    "targets": ["all"],
    "metadata": {
      "source": "external_test",
      "alert_level": "critical",
      "environment": "test"
    }
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

HIGH_PRIORITY_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notification_id'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
DESTINATIONS_COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['destinations_count'])" 2>/dev/null)

if [ -n "$HIGH_PRIORITY_ID" ] && [ "$STATUS" = "pending" ] && [ "$DESTINATIONS_COUNT" -gt 0 ]; then
    print_result 0 "高優先級通知發送成功 (ID: $HIGH_PRIORITY_ID, Destinations: $DESTINATIONS_COUNT)"
else
    print_result 1 "高優先級通知發送失敗"
fi

# ==================== 測試 3: 包含 Mentions 的通知 ====================
print_header "測試 3: 包含 Mentions 的通知"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notify\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "cfh-alert-gogo",
    "message": "Hi @TestUser, 請查看這條重要訊息",
    "messageType": "text",
    "priority": "normal",
    "targets": ["all"],
    "mentions": ["TestUser"],
    "metadata": {
      "source": "external_test",
      "environment": "test"
    }
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

MENTION_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notification_id'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
DESTINATIONS_COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['destinations_count'])" 2>/dev/null)

if [ -n "$MENTION_ID" ] && [ "$STATUS" = "pending" ] && [ "$DESTINATIONS_COUNT" -gt 0 ]; then
    print_result 0 "包含 Mentions 的通知發送成功 (ID: $MENTION_ID, Destinations: $DESTINATIONS_COUNT)"
else
    print_result 1 "包含 Mentions 的通知發送失敗"
fi

# ==================== 測試 4: 指定特定目標的通知 ====================
print_header "測試 4: 指定特定目標的通知"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notify\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"

# 使用 init.sql 中的 conversation IDs
PERSONAL_CONV_ID="a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR"
GROUPCHAT_CONV_ID="19:f26a8d8a235f430db87a404491cd2ffc@thread.v2"
CHANNEL_CONV_ID="19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2"
EXTRA_GROUPCHAT_CONV_ID="19:3a6943ad620946199a061ce2b87ea8b9@thread.v2"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "cfh-alert-gogo",
    "message": "這是一條指定目標的測試通知",
    "messageType": "text",
    "priority": "normal",
    "targets": ["'"$PERSONAL_CONV_ID"'", "'"$GROUPCHAT_CONV_ID"'", "'"$CHANNEL_CONV_ID"'", "'"$EXTRA_GROUPCHAT_CONV_ID"'"],
    "metadata": {
      "source": "external_test",
      "target_type": "specific",
      "environment": "test"
    }
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

TARGETED_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notification_id'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
DESTINATIONS_COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['destinations_count'])" 2>/dev/null)

if [ -n "$TARGETED_ID" ] && [ "$STATUS" = "pending" ] && [ "$DESTINATIONS_COUNT" -gt 0 ]; then
    print_result 0 "指定目標通知發送成功 (ID: $TARGETED_ID, Destinations: $DESTINATIONS_COUNT)"
else
    print_result 1 "指定目標通知發送失敗"
fi

# ==================== 測試 5: Google Drive 連結分享 ====================
print_header "測試 5: Google Drive 連結分享"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notify\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "cfh-alert-gogo",
    "message": "📎 重要文件分享\n\n請查看以下 Google Drive 連結：\nhttps://drive.google.com/file/d/1cpLlDyS_tKHFHoMXcRdLHGkSd6nID7nR/view?usp=drive_link\n\n請及時查看相關內容。",
    "messageType": "text",
    "priority": "normal",
    "targets": ["all"],
    "metadata": {
      "source": "manual_notification",
      "file_type": "google_drive",
      "file_id": "1cpLlDyS_tKHFHoMXcRdLHGkSd6nID7nR",
      "environment": "test"
    }
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

GOOGLE_DRIVE_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notification_id'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
DESTINATIONS_COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['destinations_count'])" 2>/dev/null)

if [ -n "$GOOGLE_DRIVE_ID" ] && [ "$STATUS" = "pending" ] && [ "$DESTINATIONS_COUNT" -gt 0 ]; then
    print_result 0 "Google Drive 連結分享成功 (ID: $GOOGLE_DRIVE_ID, Destinations: $DESTINATIONS_COUNT)"
else
    print_result 1 "Google Drive 連結分享失敗"
fi

# ==================== 測試 6: 低優先級通知 ====================
print_header "測試 6: 低優先級通知"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notify\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "cfh-alert-gogo",
    "message": "這是一條低優先級測試通知",
    "messageType": "text",
    "priority": "low",
    "targets": ["all"],
    "metadata": {
      "source": "external_test",
      "priority_level": "low",
      "environment": "test"
    }
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

LOW_PRIORITY_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notification_id'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
DESTINATIONS_COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['destinations_count'])" 2>/dev/null)

if [ -n "$LOW_PRIORITY_ID" ] && [ "$STATUS" = "pending" ] && [ "$DESTINATIONS_COUNT" -gt 0 ]; then
    print_result 0 "低優先級通知發送成功 (ID: $LOW_PRIORITY_ID, Destinations: $DESTINATIONS_COUNT)"
else
    print_result 1 "低優先級通知發送失敗"
fi

# ==================== 測試 7: 錯誤的 notifyKey ====================
print_header "測試 7: 錯誤的 notifyKey"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notify\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "invalid-key-12345",
    "message": "這是一條測試通知",
    "messageType": "text",
    "priority": "normal",
    "targets": ["all"]
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

ERROR=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('error', ''))" 2>/dev/null)

if [ -n "$ERROR" ] && [[ "$ERROR" == *"not found"* ]]; then
    print_result 0 "錯誤的 notifyKey 正確返回錯誤"
else
    print_result 1 "錯誤的 notifyKey 沒有正確處理"
fi

# ==================== 測試 8: 缺少必要欄位 ====================
print_header "測試 8: 缺少必要欄位"
echo -e "${BLUE}curl -X POST \"$BASE_URL/notify\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${BLUE}  -d '{...}'${NC}"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "cfh-alert-gogo",
    "messageType": "text",
    "priority": "normal",
    "targets": ["all"]
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

ERROR=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('error', ''))" 2>/dev/null)

if [ -n "$ERROR" ] && [[ "$ERROR" == *"required"* ]]; then
    print_result 0 "缺少必要欄位正確返回錯誤"
else
    print_result 1 "缺少必要欄位沒有正確處理"
fi

# ==================== 測試 9: 獲取專案目標列表 ====================
print_header "測試 9: 獲取專案目標列表"
echo -e "${BLUE}curl -X GET \"$BASE_URL/destinations/cfh-alert-gogo\" \\${NC}"
echo -e "${BLUE}  -H \"Content-Type: application/json\"${NC}"

RESPONSE=$(curl -s -X GET "$BASE_URL/destinations/cfh-alert-gogo" \
  -H "Content-Type: application/json")

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

DESTINATIONS=$(echo "$RESPONSE" | python3 -c "import sys, json; data = json.load(sys.stdin).get('destinations', []); print(len(data))" 2>/dev/null)

if [ -n "$DESTINATIONS" ] && [ "$DESTINATIONS" -gt 0 ]; then
    print_result 0 "獲取專案目標列表成功 (找到 $DESTINATIONS 個目標)"
else
    print_result 1 "獲取專案目標列表失敗"
fi

# ==================== 測試總結 ====================
print_header "測試總結"
echo -e "總測試數: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "通過: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失敗: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有外部通知測試通過！${NC}"
    exit 0
else
    echo -e "\n${RED}❌ 有測試失敗${NC}"
    exit 1
fi
