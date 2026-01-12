#!/bin/bash

# Template CRUD 功能測試腳本
# 測試前端頁面的 Template CRUD 功能

FRONTEND_URL="http://localhost:3000"
BACKEND_URL="http://localhost:8080"
TEMPLATES_PAGE="${FRONTEND_URL}/#/admin/templates"

echo "🧪 Template CRUD 功能測試"
echo "================================"
echo ""
echo "📋 測試環境："
echo "  Frontend: $FRONTEND_URL"
echo "  Backend: $BACKEND_URL"
echo "  Templates Page: $TEMPLATES_PAGE"
echo ""

# 檢查服務狀態
echo "1️⃣  檢查服務狀態..."
FRONTEND_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$FRONTEND_URL")
BACKEND_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BACKEND_URL/health")

if [ "$FRONTEND_STATUS" = "200" ]; then
  echo "  ✅ Frontend 運行中 ($FRONTEND_STATUS)"
else
  echo "  ❌ Frontend 無法訪問 ($FRONTEND_STATUS)"
  exit 1
fi

if [ "$BACKEND_STATUS" = "200" ]; then
  echo "  ✅ Backend 運行中 ($BACKEND_STATUS)"
else
  echo "  ⚠️  Backend 無法訪問 ($BACKEND_STATUS) - 將使用 Mock 數據"
fi

echo ""
echo "2️⃣  測試 API 端點..."

# 需要一個有效的 project ID，這裡使用一個測試 ID
PROJECT_ID="00000000-0000-0000-0000-000000000001"

# 測試 List Templates
echo "  📋 GET /internal/v1/projects/$PROJECT_ID/templates"
LIST_RESPONSE=$(curl -s -X GET "$BACKEND_URL/internal/v1/projects/$PROJECT_ID/templates" -w "\nHTTP_STATUS:%{http_code}")
HTTP_STATUS=$(echo "$LIST_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
LIST_BODY=$(echo "$LIST_RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" = "200" ]; then
  echo "    ✅ List Templates API 正常 ($HTTP_STATUS)"
  TEMPLATE_COUNT=$(echo "$LIST_BODY" | jq '.data | length' 2>/dev/null || echo "0")
  echo "    找到 $TEMPLATE_COUNT 個 templates"
else
  echo "    ⚠️  List Templates API 回應: $HTTP_STATUS"
  echo "    回應: $LIST_BODY"
fi

echo ""
echo "3️⃣  測試 Create Template..."

CREATE_DATA='{
  "name": "Test Template '$(date +%s)'",
  "description": "測試用 Template",
  "variables": [
    {"key": "title", "label": "Title", "type": "text"},
    {"key": "message", "label": "Message", "type": "text"}
  ],
  "defaultJsonStructure": "{\"type\":\"AdaptiveCard\",\"version\":\"1.4\",\"body\":[{\"type\":\"TextBlock\",\"text\":\"{{title}}\"},{\"type\":\"TextBlock\",\"text\":\"{{message}}\"}]}"
}'

CREATE_RESPONSE=$(curl -s -X POST "$BACKEND_URL/internal/v1/projects/$PROJECT_ID/templates" \
  -H "Content-Type: application/json" \
  -d "$CREATE_DATA" \
  -w "\nHTTP_STATUS:%{http_code}")

HTTP_STATUS=$(echo "$CREATE_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
CREATE_BODY=$(echo "$CREATE_RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" = "201" ] || [ "$HTTP_STATUS" = "200" ]; then
  echo "    ✅ Create Template API 正常 ($HTTP_STATUS)"
  TEMPLATE_ID=$(echo "$CREATE_BODY" | jq -r '.data.id // empty' 2>/dev/null)
  if [ -n "$TEMPLATE_ID" ] && [ "$TEMPLATE_ID" != "null" ]; then
    echo "    創建的 Template ID: $TEMPLATE_ID"
    
    # 測試 Update Template
    echo ""
    echo "4️⃣  測試 Update Template..."
    UPDATE_DATA='{
      "name": "Updated Test Template",
      "description": "更新後的描述",
      "variables": [
        {"key": "title", "label": "Title", "type": "text"},
        {"key": "message", "label": "Message", "type": "text"},
        {"key": "severity", "label": "Severity", "type": "select", "options": ["low", "high"]}
      ],
      "defaultJsonStructure": "{\"type\":\"AdaptiveCard\",\"version\":\"1.4\",\"body\":[{\"type\":\"TextBlock\",\"text\":\"{{title}}\",\"weight\":\"bolder\"}]}"
    }'
    
    UPDATE_RESPONSE=$(curl -s -X PUT "$BACKEND_URL/internal/v1/projects/$PROJECT_ID/templates/$TEMPLATE_ID" \
      -H "Content-Type: application/json" \
      -d "$UPDATE_DATA" \
      -w "\nHTTP_STATUS:%{http_code}")
    
    UPDATE_HTTP_STATUS=$(echo "$UPDATE_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
    UPDATE_BODY=$(echo "$UPDATE_RESPONSE" | sed '/HTTP_STATUS/d')
    
    if [ "$UPDATE_HTTP_STATUS" = "200" ]; then
      echo "    ✅ Update Template API 正常 ($UPDATE_HTTP_STATUS)"
      UPDATED_NAME=$(echo "$UPDATE_BODY" | jq -r '.data.name // empty' 2>/dev/null)
      echo "    更新後的 Template 名稱: $UPDATED_NAME"
      
      # 測試 Delete Template
      echo ""
      echo "5️⃣  測試 Delete Template..."
      DELETE_RESPONSE=$(curl -s -X DELETE "$BACKEND_URL/internal/v1/projects/$PROJECT_ID/templates/$TEMPLATE_ID" \
        -w "\nHTTP_STATUS:%{http_code}")
      
      DELETE_HTTP_STATUS=$(echo "$DELETE_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
      
      if [ "$DELETE_HTTP_STATUS" = "204" ] || [ "$DELETE_HTTP_STATUS" = "200" ]; then
        echo "    ✅ Delete Template API 正常 ($DELETE_HTTP_STATUS)"
      else
        echo "    ⚠️  Delete Template API 回應: $DELETE_HTTP_STATUS"
      fi
    else
      echo "    ⚠️  Update Template API 回應: $UPDATE_HTTP_STATUS"
      echo "    回應: $UPDATE_BODY"
    fi
  else
    echo "    ⚠️  無法從回應中取得 Template ID"
  fi
else
  echo "    ⚠️  Create Template API 回應: $HTTP_STATUS"
  echo "    回應: $CREATE_BODY"
fi

echo ""
echo "================================"
echo "✅ API 測試完成"
echo ""
echo "📝 下一步："
echo "  1. 打開瀏覽器訪問: $TEMPLATES_PAGE"
echo "  2. 檢查前端頁面是否正常顯示"
echo "  3. 嘗試在前端新增、修改、刪除 template"
echo "  4. 查看瀏覽器 Console 是否有錯誤訊息"
echo ""
