#!/bin/bash

# Teams Notification API 快速測試腳本
# 使用方法: ./test_api.sh

set -e

BASE_URL="http://localhost:8080"
API_BASE="$BASE_URL/api/v1"

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 測試函數
test_api() {
    local name="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    
    echo -e "${BLUE}測試 $name...${NC}"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -sS "$url")
    else
        response=$(curl -sS -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
    fi
    
    # 檢查響應是否包含錯誤
    if echo "$response" | jq -e '.error' > /dev/null 2>&1; then
        echo -e "${RED}❌ $name 失敗${NC}"
        echo "$response" | jq '.error'
        return 1
    else
        echo -e "${GREEN}✅ $name 成功${NC}"
        return 0
    fi
}

# 檢查服務器是否運行
check_server() {
    echo -e "${YELLOW}檢查服務器狀態...${NC}"
    if curl -sS "$BASE_URL/health" > /dev/null; then
        echo -e "${GREEN}✅ 服務器運行正常${NC}"
    else
        echo -e "${RED}❌ 服務器未運行，請先啟動服務器${NC}"
        echo "運行: go run ./cmd/server"
        exit 1
    fi
}

# 獲取統計信息
get_stats() {
    echo -e "${YELLOW}獲取 API 統計信息...${NC}"
    
    companies=$(curl -sS "$API_BASE/companies" | jq '.data | length')
    users=$(curl -sS "$API_BASE/users" | jq '.data | length')
    projects=$(curl -sS "$API_BASE/projects" | jq '.data | length')
    platform_bots=$(curl -sS "$API_BASE/bots/platform" | jq '.data | length')
    third_party_bots=$(curl -sS "$API_BASE/bots/third-party" | jq '.data | length')
    destinations=$(curl -sS "$API_BASE/destinations" | jq '.data | length')
    notifications=$(curl -sS "$API_BASE/notifications" | jq '.data | length')
    
    echo -e "${GREEN}📊 數據統計:${NC}"
    echo "   公司: $companies"
    echo "   用戶: $users"
    echo "   項目: $projects"
    echo "   平台機器人: $platform_bots"
    echo "   第三方機器人: $third_party_bots"
    echo "   目的地: $destinations"
    echo "   通知: $notifications"
}

# 測試基本 API
test_basic_apis() {
    echo -e "${YELLOW}測試基本 API...${NC}"
    
    # 健康檢查
    test_api "健康檢查" "GET" "$BASE_URL/health"
    
    # 獲取所有資源
    test_api "獲取公司列表" "GET" "$API_BASE/companies"
    test_api "獲取用戶列表" "GET" "$API_BASE/users"
    test_api "獲取項目列表" "GET" "$API_BASE/projects"
    test_api "獲取平台機器人列表" "GET" "$API_BASE/bots/platform"
    test_api "獲取第三方機器人列表" "GET" "$API_BASE/bots/third-party"
    test_api "獲取目的地列表" "GET" "$API_BASE/destinations"
    test_api "獲取通知列表" "GET" "$API_BASE/notifications"
}

# 測試創建操作
test_create_operations() {
    echo -e "${YELLOW}測試創建操作...${NC}"
    
    # 創建公司
    test_api "創建公司" "POST" "$API_BASE/companies" '{
        "name": "API測試公司",
        "contact_email": "test@api-test.com",
        "contact_phone": "+886-2-1111-2222",
        "address": "台北市測試區測試路123號",
        "status": "active",
        "billing_enabled": true
    }'
    
    # 創建目的地
    test_api "創建目的地" "POST" "$API_BASE/destinations" '{
        "project_id": "198f1130-20a9-4c7c-a504-6a055d27e8db",
        "name": "API測試目的地",
        "description": "API測試用目的地",
        "teams_tenant_id": "api-test-tenant",
        "targets": [
            {
                "type": "channel",
                "team_id": "api-test-team",
                "channel_id": "api-test-channel",
                "display_name": "API測試頻道"
            }
        ],
        "status": "active",
        "validation_status": "pending",
        "created_by": "643c4d7a-a18f-4caa-aff2-5a0d4439b367"
    }'
    
    # 發送通知
    test_api "發送通知" "POST" "$API_BASE/notifications" '{
        "project_id": "198f1130-20a9-4c7c-a504-6a055d27e8db",
        "sender_id": "643c4d7a-a18f-4caa-aff2-5a0d4439b367",
        "message_type": "text",
        "content": "API測試通知消息",
        "mentions": ["@test"],
        "priority": "normal",
        "destinations": ["c5133b8b-6e5c-4359-be17-221483569b97"]
    }'
}

# 測試 JSONB 字段
test_jsonb_fields() {
    echo -e "${YELLOW}測試 JSONB 字段...${NC}"
    
    # 測試複雜的 targets 數組
    test_api "複雜目標配置" "POST" "$API_BASE/destinations" '{
        "project_id": "198f1130-20a9-4c7c-a504-6a055d27e8db",
        "name": "JSONB測試目的地",
        "description": "測試複雜的JSONB字段",
        "teams_tenant_id": "jsonb-test-tenant",
        "targets": [
            {
                "type": "channel",
                "team_id": "team-jsonb-1",
                "channel_id": "channel-jsonb-1",
                "display_name": "JSONB頻道1"
            },
            {
                "type": "person",
                "user_id": "user-jsonb-1",
                "display_name": "JSONB用戶1"
            },
            {
                "type": "chatgroup",
                "group_id": "group-jsonb-1",
                "display_name": "JSONB群組1"
            }
        ],
        "status": "active",
        "validation_status": "pending",
        "created_by": "643c4d7a-a18f-4caa-aff2-5a0d4439b367"
    }'
}

# 測試中文內容
test_chinese_content() {
    echo -e "${YELLOW}測試中文內容...${NC}"
    
    # 測試中文公司名稱
    test_api "中文公司名稱" "POST" "$API_BASE/companies" '{
        "name": "台灣科技股份有限公司",
        "contact_email": "info@taiwan-tech.com",
        "contact_phone": "+886-2-2345-6789",
        "address": "台北市信義區信義路五段7號101大樓",
        "status": "active",
        "billing_enabled": true
    }'
    
    # 測試中文通知內容
    test_api "中文通知內容" "POST" "$API_BASE/notifications" '{
        "project_id": "198f1130-20a9-4c7c-a504-6a055d27e8db",
        "sender_id": "643c4d7a-a18f-4caa-aff2-5a0d4439b367",
        "message_type": "text",
        "content": "這是一個包含中文內容的測試通知消息，用於驗證系統對中文字符的處理能力。",
        "mentions": ["@管理員", "@用戶"],
        "priority": "high",
        "destinations": ["c5133b8b-6e5c-4359-be17-221483569b97"]
    }'
}

# 性能測試
test_performance() {
    echo -e "${YELLOW}性能測試...${NC}"
    
    echo "測試各 API 響應時間:"
    
    # 測試響應時間
    health_time=$(curl -sS -w "%{time_total}" -o /dev/null "$BASE_URL/health")
    companies_time=$(curl -sS -w "%{time_total}" -o /dev/null "$API_BASE/companies")
    users_time=$(curl -sS -w "%{time_total}" -o /dev/null "$API_BASE/users")
    projects_time=$(curl -sS -w "%{time_total}" -o /dev/null "$API_BASE/projects")
    platform_bots_time=$(curl -sS -w "%{time_total}" -o /dev/null "$API_BASE/bots/platform")
    third_party_bots_time=$(curl -sS -w "%{time_total}" -o /dev/null "$API_BASE/bots/third-party")
    destinations_time=$(curl -sS -w "%{time_total}" -o /dev/null "$API_BASE/destinations")
    notifications_time=$(curl -sS -w "%{time_total}" -o /dev/null "$API_BASE/notifications")
    
    echo "   健康檢查: ${health_time}s"
    echo "   公司 API: ${companies_time}s"
    echo "   用戶 API: ${users_time}s"
    echo "   項目 API: ${projects_time}s"
    echo "   平台機器人 API: ${platform_bots_time}s"
    echo "   第三方機器人 API: ${third_party_bots_time}s"
    echo "   目的地 API: ${destinations_time}s"
    echo "   通知 API: ${notifications_time}s"
}

# 主函數
main() {
    echo -e "${GREEN}=== Teams Notification API 快速測試 ===${NC}"
    echo ""
    
    # 檢查服務器
    check_server
    echo ""
    
    # 獲取統計信息
    get_stats
    echo ""
    
    # 測試基本 API
    test_basic_apis
    echo ""
    
    # 測試創建操作
    test_create_operations
    echo ""
    
    # 測試 JSONB 字段
    test_jsonb_fields
    echo ""
    
    # 測試中文內容
    test_chinese_content
    echo ""
    
    # 性能測試
    test_performance
    echo ""
    
    echo -e "${GREEN}🎉 所有測試完成！${NC}"
}

# 運行主函數
main "$@"
