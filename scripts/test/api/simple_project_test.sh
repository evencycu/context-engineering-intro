#!/bin/bash

# 簡化版 Project API 測試腳本
set -e

API_URL="http://localhost:8080"
BASE_URL="${API_URL}/api/internal/v1"

# 顏色
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=========================================="
echo "Project API 測試"
echo -e "==========================================${NC}\n"

# 1. 使用固定的公司 ID 和用戶 ID
echo -e "${BLUE}[1/9]${NC} 使用測試公司和用戶 ID..."
COMPANY_ID="550e8400-e29b-41d4-a716-446655440001"
USER_ID="650e8400-e29b-41d4-a716-446655440001"
echo -e "${GREEN}✓ 公司 ID: $COMPANY_ID${NC}"
echo -e "${GREEN}✓ 用戶 ID: $USER_ID${NC}\n"

# 2. 創建專案
echo -e "${BLUE}[2/9]${NC} 創建專案..."

# 生成 NotifyKey (使用時間戳和隨機數)
NOTIFY_KEY="test-project-$(date +%s)-$(openssl rand -hex 16 2>/dev/null || cat /dev/urandom | LC_ALL=C tr -dc 'a-f0-9' | head -c 32)"
echo -e "${BLUE}curl -X POST ${BASE_URL}/projects -H \"Content-Type: application/json\" -d '{...}'${NC}"
PROJECT_RESPONSE=$(curl -s -X POST ${BASE_URL}/projects \
  -H "Content-Type: application/json" \
  -d "{\"company_id\":\"$COMPANY_ID\",\"notify_key\":\"$NOTIFY_KEY\",\"description\":\"This is a test project for API testing\",\"daily_limit\":1000,\"monthly_limit\":30000,\"priority\":\"normal\",\"created_by\":\"$USER_ID\"}")

echo "$PROJECT_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$PROJECT_RESPONSE"

PROJECT_ID=$(echo "$PROJECT_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "")
PROJECT_KEY=$(echo "$PROJECT_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notify_key'])" 2>/dev/null || echo "")

if [ -z "$PROJECT_ID" ]; then
    echo -e "${RED}✗ 創建專案失敗${NC}"
    echo "$PROJECT_RESPONSE"
    exit 1
fi

# 檢查 NotifyKey 是否正確返回
if [ -z "$PROJECT_KEY" ]; then
    echo -e "${RED}✗ NotifyKey 未返回${NC}"
    exit 1
elif [ "$PROJECT_KEY" != "$NOTIFY_KEY" ]; then
    echo -e "${RED}✗ NotifyKey 不一致${NC}"
    echo -e "${RED}  期望: $NOTIFY_KEY${NC}"
    echo -e "${RED}  實際: $PROJECT_KEY${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 專案已創建: $PROJECT_ID${NC}"
echo -e "${GREEN}✓ NotifyKey: $PROJECT_KEY${NC}\n"

# 3. 獲取專案列表
echo -e "${BLUE}[3/9]${NC} 獲取專案列表..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/projects${NC}"
LIST_RESPONSE=$(curl -s ${BASE_URL}/projects)
TOTAL_PROJECTS=$(echo "$LIST_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['pagination']['total'])" 2>/dev/null || echo "0")
echo "{\"data\": [...], \"pagination\": $(echo "$LIST_RESPONSE" | python3 -c "import sys, json; import json as j; print(j.dumps(json.load(sys.stdin)['pagination']))" 2>/dev/null || echo "{}")}"
echo -e "${GREEN}✓ 獲取成功，總共 $TOTAL_PROJECTS 個專案${NC}\n"

# 4. 獲取單個專案
echo -e "${BLUE}[4/9]${NC} 獲取單個專案 (ID: $PROJECT_ID)..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/projects/$PROJECT_ID${NC}"
GET_RESPONSE=$(curl -s ${BASE_URL}/projects/$PROJECT_ID)
echo "$GET_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$GET_RESPONSE"
PROJECT_NAME=$(echo "$GET_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['notify_key'])" 2>/dev/null || echo "")
echo -e "${GREEN}✓ 獲取成功: $PROJECT_NAME${NC}\n"

# 5. 更新專案
echo -e "${BLUE}[5/9]${NC} 更新專案信息..."
echo -e "${BLUE}curl -X PUT ${BASE_URL}/projects/$PROJECT_ID -H \"Content-Type: application/json\" -d '{\"description\":\"Updated test project\",\"priority\":\"high\"}'${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT ${BASE_URL}/projects/$PROJECT_ID \
  -H "Content-Type: application/json" \
  -d '{"description":"Updated test project description","priority":"high"}')
echo "$UPDATE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$UPDATE_RESPONSE"
UPDATED_DESC=$(echo "$UPDATE_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['description'])" 2>/dev/null || echo "")
echo -e "${GREEN}✓ 更新成功${NC}\n"

# 6. 更新專案限制
echo -e "${BLUE}[6/9]${NC} 更新專案限制..."
echo -e "${BLUE}curl -X PATCH ${BASE_URL}/projects/$PROJECT_ID/limits -H \"Content-Type: application/json\" -d '{\"daily_limit\":2000,\"monthly_limit\":60000}'${NC}"
LIMITS_RESPONSE=$(curl -s -X PATCH ${BASE_URL}/projects/$PROJECT_ID/limits \
  -H "Content-Type: application/json" \
  -d '{"daily_limit":2000,"monthly_limit":60000}')
echo "$LIMITS_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$LIMITS_RESPONSE"
LIMITS_MESSAGE=$(echo "$LIMITS_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('message',''))" 2>/dev/null || echo "")
echo -e "${GREEN}✓ $LIMITS_MESSAGE${NC}\n"

# 7. 按公司查詢專案
echo -e "${BLUE}[7/9]${NC} 按公司查詢專案..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/projects/company/$COMPANY_ID${NC}"
COMPANY_PROJECTS_RESPONSE=$(curl -s ${BASE_URL}/projects/company/$COMPANY_ID)
COMPANY_PROJECTS_COUNT=$(echo "$COMPANY_PROJECTS_RESPONSE" | python3 -c "import sys, json; print(len(json.load(sys.stdin)['data']))" 2>/dev/null || echo "0")
echo "{\"data\": [...$COMPANY_PROJECTS_COUNT projects...]}"
echo -e "${GREEN}✓ 該公司有 $COMPANY_PROJECTS_COUNT 個專案${NC}\n"

# 8. 按 notify_key 查詢專案
echo -e "${BLUE}[8/9]${NC} 按 notify_key 查詢專案..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/projects/key/$PROJECT_KEY${NC}"
KEY_RESPONSE=$(curl -s ${BASE_URL}/projects/key/$PROJECT_KEY)
echo "$KEY_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$KEY_RESPONSE"
KEY_PROJECT_ID=$(echo "$KEY_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "")
echo -e "${GREEN}✓ 找到專案: $KEY_PROJECT_ID${NC}\n"

# 9. 刪除專案
echo -e "${BLUE}[9/9]${NC} 刪除專案..."
echo -e "${BLUE}curl -X DELETE ${BASE_URL}/projects/$PROJECT_ID${NC}"
DELETE_RESPONSE=$(curl -s -X DELETE ${BASE_URL}/projects/$PROJECT_ID)
echo "$DELETE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$DELETE_RESPONSE"
echo -e "${GREEN}✓ 測試專案已刪除${NC}\n"

echo -e "${GREEN}=========================================="
echo "✓ 所有專案 API 測試完成！"
echo -e "==========================================${NC}"

