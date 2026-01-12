#!/bin/bash

# Template CRUD API Test Script
# Usage: ./api_test.sh [base_url]
# Example: ./api_test.sh http://localhost:8080/api/internal/v1

BASE_URL="${1:-http://localhost:8080/api/internal/v1}"
PROJECT_ID="${2:-00000000-0000-0000-0000-000000000001}"

echo "🧪 Testing Template CRUD APIs"
echo "Base URL: $BASE_URL"
echo "Project ID: $PROJECT_ID"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Create Template
echo "📝 Test 1: Create Template"
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/projects/$PROJECT_ID/templates" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Template",
    "description": "Test Description",
    "variables": [
      {"key": "title", "label": "Title", "type": "text"},
      {"key": "message", "label": "Message", "type": "text"}
    ],
    "defaultJsonStructure": "{\"type\":\"AdaptiveCard\",\"version\":\"1.4\",\"body\":[{\"type\":\"TextBlock\",\"text\":\"{{title}}\"},{\"type\":\"TextBlock\",\"text\":\"{{message}}\"}]}"
  }')

echo "$CREATE_RESPONSE" | jq '.' 2>/dev/null || echo "$CREATE_RESPONSE"

TEMPLATE_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.id // empty' 2>/dev/null)

if [ -z "$TEMPLATE_ID" ] || [ "$TEMPLATE_ID" = "null" ]; then
  echo -e "${RED}❌ Failed to create template${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Template created with ID: $TEMPLATE_ID${NC}"
echo ""

# Test 2: List Templates
echo "📋 Test 2: List Templates"
LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/projects/$PROJECT_ID/templates")
echo "$LIST_RESPONSE" | jq '.' 2>/dev/null || echo "$LIST_RESPONSE"

if echo "$LIST_RESPONSE" | jq -e '.data' > /dev/null 2>&1; then
  echo -e "${GREEN}✅ Templates listed successfully${NC}"
else
  echo -e "${RED}❌ Failed to list templates${NC}"
fi
echo ""

# Test 3: Get Template
echo "🔍 Test 3: Get Template"
GET_RESPONSE=$(curl -s -X GET "$BASE_URL/projects/$PROJECT_ID/templates/$TEMPLATE_ID")
echo "$GET_RESPONSE" | jq '.' 2>/dev/null || echo "$GET_RESPONSE"

if echo "$GET_RESPONSE" | jq -e '.data.id' > /dev/null 2>&1; then
  echo -e "${GREEN}✅ Template retrieved successfully${NC}"
else
  echo -e "${RED}❌ Failed to get template${NC}"
fi
echo ""

# Test 4: Update Template
echo "✏️  Test 4: Update Template"
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/projects/$PROJECT_ID/templates/$TEMPLATE_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Template",
    "description": "Updated Description",
    "variables": [
      {"key": "title", "label": "Title", "type": "text"},
      {"key": "message", "label": "Message", "type": "text"},
      {"key": "severity", "label": "Severity", "type": "select", "options": ["low", "medium", "high"]}
    ],
    "defaultJsonStructure": "{\"type\":\"AdaptiveCard\",\"version\":\"1.4\",\"body\":[{\"type\":\"TextBlock\",\"text\":\"{{title}}\",\"weight\":\"bolder\"},{\"type\":\"TextBlock\",\"text\":\"{{message}}\"}]}"
  }')

echo "$UPDATE_RESPONSE" | jq '.' 2>/dev/null || echo "$UPDATE_RESPONSE"

if echo "$UPDATE_RESPONSE" | jq -e '.data.name == "Updated Template"' > /dev/null 2>&1; then
  echo -e "${GREEN}✅ Template updated successfully${NC}"
else
  echo -e "${RED}❌ Failed to update template${NC}"
fi
echo ""

# Test 5: Delete Template
echo "🗑️  Test 5: Delete Template"
DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/projects/$PROJECT_ID/templates/$TEMPLATE_ID" -w "\nHTTP_STATUS:%{http_code}")
HTTP_STATUS=$(echo "$DELETE_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
DELETE_BODY=$(echo "$DELETE_RESPONSE" | sed '/HTTP_STATUS/d')

echo "$DELETE_BODY" | jq '.' 2>/dev/null || echo "$DELETE_BODY"

if [ "$HTTP_STATUS" = "204" ] || [ "$HTTP_STATUS" = "200" ]; then
  echo -e "${GREEN}✅ Template deleted successfully${NC}"
else
  echo -e "${RED}❌ Failed to delete template (HTTP $HTTP_STATUS)${NC}"
fi
echo ""

echo -e "${GREEN}🎉 All tests completed!${NC}"
