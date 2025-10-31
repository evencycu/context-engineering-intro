#!/bin/bash

# Simple Adaptive Card Test Script
# 測試 Teams Adaptive Card 功能

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

echo -e "${BLUE}開始測試 Teams Adaptive Card 功能${NC}"
echo -e "${BLUE}Base URL: $BASE_URL${NC}"

# ==================== 測試 1: 基本 Adaptive Card ====================
print_header "測試 1: 基本 Adaptive Card"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "cfh-alert-gogo",
    "message": "這是一條包含 Adaptive Card 的通知",
    "messageType": "adaptive_card",
    "priority": "normal",
    "targets": ["all"],
    "attachments": [
      {
        "contentType": "application/vnd.microsoft.card.adaptive",
        "content": {
          "$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
          "type": "AdaptiveCard",
          "version": "1.5",
          "body": [
            { "type": "TextBlock", "text": "Teams Notification API 測試", "weight": "Bolder", "size": "Medium" },
            { "type": "TextBlock", "text": "這是一條測試通知，包含 Adaptive Card 內容。", "wrap": true },
            { "type": "FactSet", "facts": [
              { "title": "測試時間", "value": "2024-01-01 12:00:00" },
              { "title": "測試類型", "value": "Adaptive Card" },
              { "title": "優先級", "value": "Normal" }
            ]}
          ],
          "actions": [
            { "type": "Action.OpenUrl", "title": "查看詳情", "url": "https://docs.microsoft.com/en-us/adaptive-cards/" }
          ]
        }
      }
    ]
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

NOTIFICATION_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notification_id'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
DESTINATIONS_COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['destinations_count'])" 2>/dev/null)

if [ -n "$NOTIFICATION_ID" ] && [ "$STATUS" = "pending" ] && [ "$DESTINATIONS_COUNT" -gt 0 ]; then
    print_result 0 "基本 Adaptive Card 發送成功 (ID: $NOTIFICATION_ID, Destinations: $DESTINATIONS_COUNT)"
else
    print_result 1 "基本 Adaptive Card 發送失敗"
fi

# ==================== 測試 2: 包含圖片的 Adaptive Card ====================
print_header "測試 2: 包含圖片的 Adaptive Card"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "cfh-alert-gogo",
    "message": "包含圖片的 Adaptive Card 測試",
    "messageType": "adaptive_card",
    "priority": "normal",
    "targets": ["all"],
    "attachments": [
      {
        "contentType": "application/vnd.microsoft.card.adaptive",
        "content": {
          "$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
          "type": "AdaptiveCard",
          "version": "1.5",
          "body": [
            { "type": "TextBlock", "text": "📊 系統狀態報告", "weight": "Bolder", "size": "Large" },
            { "type": "Image", "url": "https://via.placeholder.com/400x200/0078D4/FFFFFF?text=Teams+Notification+API", "altText": "Teams Notification API Logo", "size": "Medium" },
            { "type": "TextBlock", "text": "系統運行正常，所有服務都在正常運行中。", "wrap": true },
            { "type": "ColumnSet", "columns": [
              { "type": "Column", "width": "stretch", "items": [
                { "type": "TextBlock", "text": "✅ 服務狀態", "weight": "Bolder" },
                { "type": "TextBlock", "text": "正常運行" }
              ]},
              { "type": "Column", "width": "stretch", "items": [
                { "type": "TextBlock", "text": "📈 使用率", "weight": "Bolder" },
                { "type": "TextBlock", "text": "75%" }
              ]}
            ]}
          ],
          "actions": [
            { "type": "Action.Submit", "title": "確認收到", "data": { "action": "acknowledge" } }
          ]
        }
      }
    ]
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

IMAGE_CARD_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notification_id'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
DESTINATIONS_COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['destinations_count'])" 2>/dev/null)

if [ -n "$IMAGE_CARD_ID" ] && [ "$STATUS" = "pending" ] && [ "$DESTINATIONS_COUNT" -gt 0 ]; then
    print_result 0 "包含圖片的 Adaptive Card 發送成功 (ID: $IMAGE_CARD_ID, Destinations: $DESTINATIONS_COUNT)"
else
    print_result 1 "包含圖片的 Adaptive Card 發送失敗"
fi

# ==================== 測試 3: 包含輸入欄位的 Adaptive Card ====================
print_header "測試 3: 包含輸入欄位的 Adaptive Card"

RESPONSE=$(curl -s -X POST "$BASE_URL/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "cfh-alert-gogo",
    "message": "包含輸入欄位的 Adaptive Card 測試",
    "messageType": "adaptive_card",
    "priority": "normal",
    "targets": ["all"],
    "attachments": [
      {
        "contentType": "application/vnd.microsoft.card.adaptive",
        "content": {
          "$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
          "type": "AdaptiveCard",
          "version": "1.5",
          "body": [
            { "type": "TextBlock", "text": "📝 意見回饋表單", "weight": "Bolder", "size": "Large" },
            { "type": "TextBlock", "text": "請填寫以下表單來提供您的意見回饋：", "wrap": true },
            { "type": "Input.Text", "id": "feedback_title", "label": "標題", "placeholder": "請輸入標題", "maxLength": 100 },
            { "type": "Input.Text", "id": "feedback_content", "label": "內容", "placeholder": "請輸入您的意見回饋", "isMultiline": true, "maxLength": 500 },
            { "type": "Input.ChoiceSet", "id": "feedback_rating", "label": "評分", "choices": [
              { "title": "非常滿意", "value": "5" },
              { "title": "滿意", "value": "4" },
              { "title": "普通", "value": "3" },
              { "title": "不滿意", "value": "2" },
              { "title": "非常不滿意", "value": "1" }
            ], "style": "compact" }
          ],
          "actions": [
            { "type": "Action.Submit", "title": "提交回饋", "data": { "action": "submit_feedback", "formType": "feedback" } },
            { "type": "Action.Submit", "title": "取消", "data": { "action": "cancel", "formType": "feedback" } }
          ]
        }
      }
    ]
  }')

echo "$RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE"

FORM_CARD_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notification_id'])" 2>/dev/null)
STATUS=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['status'])" 2>/dev/null)
DESTINATIONS_COUNT=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['destinations_count'])" 2>/dev/null)

if [ -n "$FORM_CARD_ID" ] && [ "$STATUS" = "pending" ] && [ "$DESTINATIONS_COUNT" -gt 0 ]; then
    print_result 0 "包含輸入欄位的 Adaptive Card 發送成功 (ID: $FORM_CARD_ID, Destinations: $DESTINATIONS_COUNT)"
else
    print_result 1 "包含輸入欄位的 Adaptive Card 發送失敗"
fi

# ==================== 測試總結 ====================
print_header "Adaptive Card 測試總結"
echo -e "總測試數: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "通過: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失敗: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有 Adaptive Card 測試通過！${NC}"
    exit 0
else
    echo -e "\n${RED}❌ 有測試失敗${NC}"
    exit 1
fi
