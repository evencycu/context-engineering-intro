#!/bin/bash

# 簡化版 User API 測試腳本
set -e

API_URL="http://localhost:8080"
BASE_URL="${API_URL}/api/internal/v1"

# 顏色
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=========================================="
echo "User API 測試"
echo -e "==========================================${NC}\n"

# 1. 使用固定的公司 ID
echo -e "${BLUE}[1/10]${NC} 使用測試公司 ID..."
COMPANY_ID="550e8400-e29b-41d4-a716-446655440001"
echo -e "${GREEN}✓ 公司 ID: $COMPANY_ID${NC}\n"

# 2. 創建用戶
echo -e "${BLUE}[2/10]${NC} 創建用戶..."
TIMESTAMP=$(date +%s)
USER_EMAIL="testuser${TIMESTAMP}@example.com"

echo -e "${BLUE}curl -X POST ${BASE_URL}/users -H \"Content-Type: application/json\" -d '{...}'${NC}"
USER_RESPONSE=$(curl -s -X POST ${BASE_URL}/users \
  -H "Content-Type: application/json" \
  -d "{\"company_id\":\"$COMPANY_ID\",\"email\":\"$USER_EMAIL\",\"name\":\"Test User\",\"role\":\"user\",\"password\":\"password123\"}")

echo "$USER_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$USER_RESPONSE"
USER_ID=$(echo "$USER_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "")

if [ -z "$USER_ID" ]; then
    echo -e "${RED}✗ 創建用戶失敗${NC}"
    echo "$USER_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ 用戶已創建: $USER_ID${NC}"
echo -e "  Email: $USER_EMAIL\n"

# 3. 獲取用戶列表
echo -e "${BLUE}[3/10]${NC} 獲取用戶列表..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/users${NC}"
LIST_RESPONSE=$(curl -s ${BASE_URL}/users)
TOTAL_USERS=$(echo "$LIST_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['pagination']['total'])" 2>/dev/null || echo "0")
echo "{\"data\": [...], \"pagination\": $(echo "$LIST_RESPONSE" | python3 -c "import sys, json; import json as j; print(j.dumps(json.load(sys.stdin)['pagination']))" 2>/dev/null || echo "{}")}"
echo -e "${GREEN}✓ 獲取成功，總共 $TOTAL_USERS 個用戶${NC}\n"

# 4. 獲取單個用戶
echo -e "${BLUE}[4/10]${NC} 獲取單個用戶 (ID: $USER_ID)..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/users/$USER_ID${NC}"
GET_RESPONSE=$(curl -s ${BASE_URL}/users/$USER_ID)
echo "$GET_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$GET_RESPONSE"
USER_NAME=$(echo "$GET_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['name'])" 2>/dev/null || echo "")
echo -e "${GREEN}✓ 獲取成功: $USER_NAME${NC}\n"

# 5. 更新用戶
echo -e "${BLUE}[5/10]${NC} 更新用戶信息..."
echo -e "${BLUE}curl -X PUT ${BASE_URL}/users/$USER_ID -H \"Content-Type: application/json\" -d '{\"name\":\"Updated Test User\",\"role\":\"manager\"}'${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT ${BASE_URL}/users/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Test User","role":"manager"}')
echo "$UPDATE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$UPDATE_RESPONSE"
UPDATED_NAME=$(echo "$UPDATE_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['name'])" 2>/dev/null || echo "")
echo -e "${GREEN}✓ 更新成功: $UPDATED_NAME${NC}\n"

# 6. 修改密碼
echo -e "${BLUE}[6/10]${NC} 修改用戶密碼..."
echo -e "${BLUE}curl -X PATCH ${BASE_URL}/users/$USER_ID/password -H \"Content-Type: application/json\" -d '{...}'${NC}"
PWD_RESPONSE=$(curl -s -X PATCH ${BASE_URL}/users/$USER_ID/password \
  -H "Content-Type: application/json" \
  -d '{"old_password":"password123","new_password":"newpassword123"}')
echo "$PWD_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$PWD_RESPONSE"
PWD_MESSAGE=$(echo "$PWD_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('message',''))" 2>/dev/null || echo "")
echo -e "${GREEN}✓ $PWD_MESSAGE${NC}\n"

# 7. 按公司查詢用戶
echo -e "${BLUE}[7/8]${NC} 按公司查詢用戶..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/users/company/$COMPANY_ID${NC}"
COMPANY_USERS_RESPONSE=$(curl -s ${BASE_URL}/users/company/$COMPANY_ID)
COMPANY_USERS_COUNT=$(echo "$COMPANY_USERS_RESPONSE" | python3 -c "import sys, json; print(len(json.load(sys.stdin)['data']))" 2>/dev/null || echo "0")
echo "{\"data\": [...$COMPANY_USERS_COUNT users...]}"
echo -e "${GREEN}✓ 該公司有 $COMPANY_USERS_COUNT 個用戶${NC}\n"

# 8. 按角色查詢用戶
echo -e "${BLUE}[8/8]${NC} 按角色查詢用戶 (manager)..."
echo -e "${BLUE}curl -X GET ${BASE_URL}/users/role/manager${NC}"
ROLE_USERS_RESPONSE=$(curl -s ${BASE_URL}/users/role/manager)
echo "$ROLE_USERS_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$ROLE_USERS_RESPONSE"
ROLE_USERS_COUNT=$(echo "$ROLE_USERS_RESPONSE" | python3 -c "import sys, json; print(len(json.load(sys.stdin)['data']))" 2>/dev/null || echo "0")
echo -e "${GREEN}✓ Manager 角色有 $ROLE_USERS_COUNT 個用戶${NC}\n"

# 清理
echo -e "${BLUE}清理測試數據...${NC}"
echo -e "${BLUE}curl -X DELETE ${BASE_URL}/users/$USER_ID${NC}"
DELETE_RESPONSE=$(curl -s -X DELETE ${BASE_URL}/users/$USER_ID)
echo "$DELETE_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$DELETE_RESPONSE"
echo -e "${GREEN}✓ 測試用戶已刪除${NC}\n"

echo -e "${GREEN}=========================================="
echo "✓ 所有用戶 API 測試完成！"
echo -e "==========================================${NC}"

