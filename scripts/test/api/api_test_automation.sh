#!/bin/bash

# TeamsNotifyGoV2 API 自動化測試腳本
# 用於驗證所有API端點的功能和效能

set -e

# 配置
BASE_URL="http://localhost:8080"
API_VERSION="v1"
NOTIFY_KEY="cfh-alert-gogo"
TEST_RESULTS_DIR="test_results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 建立測試結果目錄
mkdir -p "$TEST_RESULTS_DIR"

echo -e "${BLUE}🚀 TeamsNotifyGoV2 API 自動化測試開始${NC}"
echo "測試時間: $(date)"
echo "測試目錄: $TEST_RESULTS_DIR"
echo "=========================================="

# 測試結果統計
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 測試函數
run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_status="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "\n${YELLOW}🧪 測試: $test_name${NC}"
    
    # 執行測試
    local start_time=$(date +%s.%N)
    local response=$(eval "$test_command" 2>&1)
    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc)
    
    # 檢查響應狀態 - 修正狀態碼提取
    local status_code=$(echo "$response" | grep -o "[0-9][0-9][0-9]$" | tail -1)
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ 通過${NC} (${duration}s)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        echo "$test_name: PASSED (${duration}s)" >> "$TEST_RESULTS_DIR/test_${TIMESTAMP}.log"
    else
        echo -e "${RED}❌ 失敗${NC} (期望: $expected_status, 實際: $status_code)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        echo "$test_name: FAILED (期望: $expected_status, 實際: $status_code)" >> "$TEST_RESULTS_DIR/test_${TIMESTAMP}.log"
        echo "響應: $response" >> "$TEST_RESULTS_DIR/test_${TIMESTAMP}.log"
    fi
}

# 1. 系統健康檢查測試
echo -e "\n${BLUE}📋 系統健康檢查測試${NC}"
run_test "系統健康檢查" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/health'" "200"
run_test "外部API健康檢查" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/$API_VERSION/external/health'" "200"

# 2. 外部API功能測試
echo -e "\n${BLUE}📋 外部API功能測試${NC}"
run_test "驗證notify key" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/$API_VERSION/external/validate/$NOTIFY_KEY'" "200"
run_test "獲取專案目的地" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/$API_VERSION/external/destinations/$NOTIFY_KEY'" "200"

# 3. 通知發送測試
echo -e "\n${BLUE}📋 通知發送測試${NC}"
run_test "基本通知發送" "curl -s -w '%{http_code}' -X POST '$BASE_URL/api/$API_VERSION/external/notify' -H 'Content-Type: application/json' -d '{\"notify_key\": \"$NOTIFY_KEY\", \"message\": \"自動化測試訊息 - $(date)\"}'" "200"

# 4. 效能測試
echo -e "\n${BLUE}📋 效能測試${NC}"
echo "執行5次連續請求測試..."

for i in {1..5}; do
    echo "效能測試 $i/5"
    start_time=$(date +%s.%N)
    curl -s "$BASE_URL/api/$API_VERSION/external/health" > /dev/null
    end_time=$(date +%s.%N)
    duration=$(echo "$end_time - $start_time" | bc)
    echo "  響應時間: ${duration}s"
    echo "效能測試 $i: ${duration}s" >> "$TEST_RESULTS_DIR/performance_${TIMESTAMP}.log"
done

# 5. 錯誤處理測試
echo -e "\n${BLUE}📋 錯誤處理測試${NC}"
run_test "無效notify key" "curl -s -w '%{http_code}' -o /dev/null '$BASE_URL/api/$API_VERSION/external/validate/invalid-key'" "200"
run_test "缺少必要參數" "curl -s -w '%{http_code}' -X POST '$BASE_URL/api/$API_VERSION/external/notify' -H 'Content-Type: application/json' -d '{}'" "400"

# 6. 負載測試
echo -e "\n${BLUE}📋 負載測試${NC}"
echo "執行10個並發請求..."

for i in {1..10}; do
    (
        curl -s "$BASE_URL/health" > /dev/null &
    ) &
done
wait
echo "✅ 並發測試完成"

# 測試結果摘要
echo -e "\n${BLUE}📊 測試結果摘要${NC}"
echo "=========================================="
echo "總測試數: $TOTAL_TESTS"
echo -e "通過: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失敗: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有測試通過！系統運行正常${NC}"
    exit 0
else
    echo -e "\n${RED}⚠️  有 $FAILED_TESTS 個測試失敗，請檢查日誌${NC}"
    exit 1
fi
