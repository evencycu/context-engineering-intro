#!/bin/bash

# Teams Notification API E2E 測試腳本
# 執行完整的端到端測試流程

set -e

# 配置
BASE_URL="http://localhost:8080"
API_BASE="$BASE_URL/api/v1"
TEST_RESULTS_DIR="test_results/e2e"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
LOG_FILE="$TEST_RESULTS_DIR/e2e_test_$TIMESTAMP.log"

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# 測試統計
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
START_TIME=$(date +%s)

# 建立測試結果目錄
mkdir -p "$TEST_RESULTS_DIR"

# 日誌函數
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

# 測試函數
run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_status="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    log "🧪 執行測試: $test_name"
    
    # 執行測試
    local start_time=$(date +%s.%N)
    local response=$(eval "$test_command" 2>&1)
    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc)
    
    # 檢查結果
    if [ $? -eq 0 ]; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
        log "✅ $test_name 通過 (${duration}s)"
        echo "$response" >> "$TEST_RESULTS_DIR/success_$TIMESTAMP.log"
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
        log "❌ $test_name 失敗 (${duration}s)"
        echo "$response" >> "$TEST_RESULTS_DIR/failure_$TIMESTAMP.log"
    fi
}

# 檢查服務狀態
check_services() {
    log "🔍 檢查服務狀態..."
    
    # 檢查 API 服務
    if curl -sS "$BASE_URL/health" > /dev/null; then
        log "✅ API 服務正常"
    else
        log "❌ API 服務未運行"
        exit 1
    fi
    
    # 檢查資料庫
    if docker exec teamsnotify-postgres-local psql -U teamsnotify -d notification_center -c "SELECT 1;" > /dev/null 2>&1; then
        log "✅ 資料庫連線正常"
    else
        log "❌ 資料庫連線失敗"
        exit 1
    fi
    
    # 檢查 Redis
    if docker exec teamsnotify-redis-local redis-cli ping | grep -q "PONG"; then
        log "✅ Redis 連線正常"
    else
        log "❌ Redis 連線失敗"
        exit 1
    fi
}

# 準備測試資料
prepare_test_data() {
    log "📋 準備測試資料..."
    
    # 建立測試公司
    COMPANY_RESPONSE=$(curl -s -X POST "$API_BASE/companies" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "E2E Test Company",
            "contact_email": "test@example.com",
            "contact_phone": "+886-2-1234-5678",
            "address": "台北市信義區信義路五段7號"
        }')
    
    COMPANY_ID=$(echo "$COMPANY_RESPONSE" | jq -r '.data.id')
    log "✅ 建立測試公司: $COMPANY_ID"
    
    # 建立測試用戶
    USER_RESPONSE=$(curl -s -X POST "$API_BASE/users" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"e2e_test_user\",
            \"email\": \"e2e@example.com\",
            \"company_id\": \"$COMPANY_ID\",
            \"role\": \"user\",
            \"status\": \"active\"
        }")
    
    USER_ID=$(echo "$USER_RESPONSE" | jq -r '.data.id')
    log "✅ 建立測試用戶: $USER_ID"
    
    # 建立測試專案
    PROJECT_RESPONSE=$(curl -s -X POST "$API_BASE/projects" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\": \"E2E Test Project\",
            \"description\": \"E2E 測試專案\",
            \"company_id\": \"$COMPANY_ID\",
            \"created_by\": \"$USER_ID\"
        }")
    
    PROJECT_ID=$(echo "$PROJECT_RESPONSE" | jq -r '.data.id')
    NOTIFY_KEY=$(echo "$PROJECT_RESPONSE" | jq -r '.data.notify_key')
    log "✅ 建立測試專案: $PROJECT_ID (notify_key: $NOTIFY_KEY)"
    
    # 建立測試目的地
    DESTINATION_RESPONSE=$(curl -s -X POST "$API_BASE/destinations" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\": \"E2E Test Destination\",
            \"description\": \"E2E 測試目的地\",
            \"project_id\": \"$PROJECT_ID\",
            \"teams_tenant_id\": \"051cece0-e4dc-4aed-b471-bf29824e1ee6\",
            \"targets\": [
                {
                    \"type\": \"personal\",
                    \"conversation_id\": \"a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR\",
                    \"tenant_id\": \"051cece0-e4dc-4aed-b471-bf29824e1ee6\"
                }
            ]
        }")
    
    DESTINATION_ID=$(echo "$DESTINATION_RESPONSE" | jq -r '.data.id')
    log "✅ 建立測試目的地: $DESTINATION_ID"
    
    # 儲存測試資料供後續使用
    echo "COMPANY_ID=$COMPANY_ID" > "$TEST_RESULTS_DIR/test_data.env"
    echo "USER_ID=$USER_ID" >> "$TEST_RESULTS_DIR/test_data.env"
    echo "PROJECT_ID=$PROJECT_ID" >> "$TEST_RESULTS_DIR/test_data.env"
    echo "NOTIFY_KEY=$NOTIFY_KEY" >> "$TEST_RESULTS_DIR/test_data.env"
    echo "DESTINATION_ID=$DESTINATION_ID" >> "$TEST_RESULTS_DIR/test_data.env"
}

# 基本功能測試
test_basic_functionality() {
    log "🧪 開始基本功能測試..."
    
    # 載入測試資料
    source "$TEST_RESULTS_DIR/test_data.env"
    
    # 測試 1: 健康檢查
    run_test "健康檢查" "curl -sS '$BASE_URL/health' | jq -r '.status' | grep -q 'healthy'"
    
    # 測試 2: 佇列狀態檢查
    run_test "佇列狀態檢查" "curl -sS '$API_BASE/queue/status' | jq -r '.circuit_state' | grep -q 'closed'"
    
    # 測試 3: 外部通知發送
    run_test "外部通知發送" "curl -sS -X POST '$API_BASE/external/notify' \
        -H 'Content-Type: application/json' \
        -d '{\"notify_key\": \"$NOTIFY_KEY\", \"message\": \"E2E 測試通知 - $(date)\"}' \
        | jq -r '.success' | grep -q 'true'"
    
    # 測試 4: 通知狀態查詢
    sleep 2  # 等待處理
    run_test "通知狀態查詢" "curl -sS '$API_BASE/notifications' | jq -r '.data | length' | grep -q '^[0-9]'"
}

# 錯誤處理測試
test_error_handling() {
    log "🧪 開始錯誤處理測試..."
    
    # 測試 1: 無效的 notify_key
    run_test "無效 notify_key 處理" "curl -sS -X POST '$API_BASE/external/notify' \
        -H 'Content-Type: application/json' \
        -d '{\"notify_key\": \"invalid-key\", \"message\": \"Test\"}' \
        | jq -r '.success' | grep -q 'false'"
    
    # 測試 2: 缺少必要參數
    run_test "缺少必要參數處理" "curl -sS -X POST '$API_BASE/external/notify' \
        -H 'Content-Type: application/json' \
        -d '{\"message\": \"Test\"}' \
        | jq -r '.success' | grep -q 'false'"
    
    # 測試 3: 無效的 JSON 格式
    run_test "無效 JSON 格式處理" "curl -sS -X POST '$API_BASE/external/notify' \
        -H 'Content-Type: application/json' \
        -d 'invalid-json' \
        | jq -r '.error' | grep -q '.'"
}

# 效能測試
test_performance() {
    log "🧪 開始效能測試..."
    
    source "$TEST_RESULTS_DIR/test_data.env"
    
    # 測試 1: 單一請求回應時間
    run_test "單一請求回應時間" "curl -w '%{time_total}' -sS -X POST '$API_BASE/external/notify' \
        -H 'Content-Type: application/json' \
        -d '{\"notify_key\": \"$NOTIFY_KEY\", \"message\": \"Performance test\"}' \
        -o /dev/null | awk '{if (\$1 < 1.0) exit 0; else exit 1}'"
    
    # 測試 2: 併發請求處理
    log "📊 執行併發測試 (10 個併發請求)..."
    for i in {1..10}; do
        (
            curl -sS -X POST "$API_BASE/external/notify" \
                -H "Content-Type: application/json" \
                -d "{\"notify_key\": \"$NOTIFY_KEY\", \"message\": \"Concurrent test $i\"}" \
                > /dev/null 2>&1
        ) &
    done
    wait
    
    run_test "併發請求處理" "echo '併發測試完成'"
}

# 佇列系統測試
test_queue_system() {
    log "🧪 開始佇列系統測試..."
    
    source "$TEST_RESULTS_DIR/test_data.env"
    
    # 測試 1: 佇列狀態監控
    run_test "佇列狀態監控" "curl -sS '$API_BASE/queue/status' | jq -r '.total_pending' | grep -q '^[0-9]'"
    
    # 測試 2: 熔斷器狀態
    run_test "熔斷器狀態檢查" "curl -sS '$API_BASE/queue/circuit-breaker/metrics' | jq -r '.state' | grep -q 'closed'"
    
    # 測試 3: 佇列處理能力
    log "📊 測試佇列處理能力..."
    for i in {1..5}; do
        curl -sS -X POST "$API_BASE/external/notify" \
            -H "Content-Type: application/json" \
            -d "{\"notify_key\": \"$NOTIFY_KEY\", \"message\": \"Queue test $i\"}" \
            > /dev/null 2>&1
    done
    
    sleep 3  # 等待處理
    
    run_test "佇列處理能力" "curl -sS '$API_BASE/queue/status' | jq -r '.total_pending' | awk '{if (\$1 <= 5) exit 0; else exit 1}'"
}

# 資料一致性測試
test_data_consistency() {
    log "🧪 開始資料一致性測試..."
    
    source "$TEST_RESULTS_DIR/test_data.env"
    
    # 測試 1: 通知記錄一致性
    run_test "通知記錄一致性" "curl -sS '$API_BASE/notifications' | jq -r '.data | length' | awk '{if (\$1 > 0) exit 0; else exit 1}'"
    
    # 測試 2: 專案資料一致性
    run_test "專案資料一致性" "curl -sS '$API_BASE/projects/$PROJECT_ID' | jq -r '.data.id' | grep -q '$PROJECT_ID'"
    
    # 測試 3: 目的地資料一致性
    run_test "目的地資料一致性" "curl -sS '$API_BASE/destinations/$DESTINATION_ID' | jq -r '.data.id' | grep -q '$DESTINATION_ID'"
}

# 清理測試資料
cleanup_test_data() {
    log "🧹 清理測試資料..."
    
    source "$TEST_RESULTS_DIR/test_data.env"
    
    # 刪除測試目的地
    curl -sS -X DELETE "$API_BASE/destinations/$DESTINATION_ID" > /dev/null 2>&1 || true
    
    # 刪除測試專案
    curl -sS -X DELETE "$API_BASE/projects/$PROJECT_ID" > /dev/null 2>&1 || true
    
    # 刪除測試用戶
    curl -sS -X DELETE "$API_BASE/users/$USER_ID" > /dev/null 2>&1 || true
    
    # 刪除測試公司
    curl -sS -X DELETE "$API_BASE/companies/$COMPANY_ID" > /dev/null 2>&1 || true
    
    log "✅ 測試資料清理完成"
}

# 生成測試報告
generate_report() {
    local end_time=$(date +%s)
    local duration=$((end_time - START_TIME))
    
    log "📊 生成測試報告..."
    
    cat > "$TEST_RESULTS_DIR/e2e_report_$TIMESTAMP.json" << EOF
{
  "test_summary": {
    "timestamp": "$(date -Iseconds)",
    "duration": "${duration}s",
    "total_tests": $TOTAL_TESTS,
    "passed": $PASSED_TESTS,
    "failed": $FAILED_TESTS,
    "success_rate": "$(echo "scale=2; $PASSED_TESTS * 100 / $TOTAL_TESTS" | bc)%"
  },
  "environment": {
    "base_url": "$BASE_URL",
    "api_version": "v1",
    "test_timestamp": "$TIMESTAMP"
  },
  "test_categories": {
    "basic_functionality": "基本功能測試",
    "error_handling": "錯誤處理測試", 
    "performance": "效能測試",
    "queue_system": "佇列系統測試",
    "data_consistency": "資料一致性測試"
  }
}
EOF
    
    log "📋 測試報告已生成: $TEST_RESULTS_DIR/e2e_report_$TIMESTAMP.json"
}

# 主執行流程
main() {
    echo -e "${BLUE}🚀 Teams Notification API E2E 測試開始${NC}"
    echo "測試時間: $(date)"
    echo "測試目錄: $TEST_RESULTS_DIR"
    echo "=========================================="
    
    # 檢查服務狀態
    check_services
    
    # 準備測試資料
    prepare_test_data
    
    # 執行測試類別
    test_basic_functionality
    test_error_handling
    test_performance
    test_queue_system
    test_data_consistency
    
    # 清理測試資料
    cleanup_test_data
    
    # 生成測試報告
    generate_report
    
    # 顯示測試結果
    echo -e "\n${PURPLE}📊 測試結果摘要${NC}"
    echo "=========================================="
    echo -e "總測試數: ${BLUE}$TOTAL_TESTS${NC}"
    echo -e "通過: ${GREEN}$PASSED_TESTS${NC}"
    echo -e "失敗: ${RED}$FAILED_TESTS${NC}"
    echo -e "成功率: ${YELLOW}$(echo "scale=2; $PASSED_TESTS * 100 / $TOTAL_TESTS" | bc)%${NC}"
    echo -e "執行時間: ${BLUE}$(echo "scale=2; $(date +%s) - $START_TIME" | bc)秒${NC}"
    echo "=========================================="
    
    if [ $FAILED_TESTS -eq 0 ]; then
        echo -e "${GREEN}🎉 所有測試通過！${NC}"
        exit 0
    else
        echo -e "${RED}❌ 有 $FAILED_TESTS 個測試失敗${NC}"
        exit 1
    fi
}

# 執行主函數
main "$@"
