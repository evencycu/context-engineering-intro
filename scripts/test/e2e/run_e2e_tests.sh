#!/bin/bash

# E2E 測試整合執行腳本
# 執行完整的端到端測試流程

set -e

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
TEST_RESULTS_DIR="test_results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# 日誌函數
log() {
    echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 檢查必要工具
check_prerequisites() {
    log "🔍 檢查必要工具..."
    
    local missing_tools=()
    
    # 檢查必要工具
    for tool in curl jq docker docker-compose; do
        if ! command -v "$tool" &> /dev/null; then
            missing_tools+=("$tool")
        fi
    done
    
    # 檢查 k6 (可選)
    if ! command -v k6 &> /dev/null; then
        warning "k6 未安裝，將跳過負載測試"
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        error "缺少必要工具: ${missing_tools[*]}"
        error "請安裝缺少的工具後重試"
        exit 1
    fi
    
    success "所有必要工具已安裝"
}

# 設定測試環境
setup_test_environment() {
    log "🚀 設定測試環境..."
    
    # 執行環境準備腳本
    if [ -f "$SCRIPT_DIR/setup_test_env.sh" ]; then
        "$SCRIPT_DIR/setup_test_env.sh" setup
    else
        error "找不到環境準備腳本: $SCRIPT_DIR/setup_test_env.sh"
        exit 1
    fi
}

# 執行 E2E 測試
run_e2e_tests() {
    log "🧪 執行 E2E 測試..."
    
    # 建立測試結果目錄
    mkdir -p "$TEST_RESULTS_DIR/e2e"
    
    # 執行 E2E 測試腳本
    if [ -f "$SCRIPT_DIR/e2e_test.sh" ]; then
        "$SCRIPT_DIR/e2e_test.sh"
    else
        error "找不到 E2E 測試腳本: $SCRIPT_DIR/e2e_test.sh"
        exit 1
    fi
}

# 執行負載測試
run_load_tests() {
    log "⚡ 執行負載測試..."
    
    # 檢查 k6 是否可用
    if ! command -v k6 &> /dev/null; then
        warning "k6 未安裝，跳過負載測試"
        return 0
    fi
    
    # 建立負載測試結果目錄
    mkdir -p "$TEST_RESULTS_DIR/load"
    
    # 執行負載測試
    if [ -f "$SCRIPT_DIR/load_test.js" ]; then
        log "執行 k6 負載測試..."
        k6 run --out json="$TEST_RESULTS_DIR/load/load_test_results.json" "$SCRIPT_DIR/load_test.js"
        success "負載測試完成"
    else
        warning "找不到負載測試腳本: $SCRIPT_DIR/load_test.js"
    fi
}

# 執行 API 測試
run_api_tests() {
    log "🔌 執行 API 測試..."
    
    # 建立 API 測試結果目錄
    mkdir -p "$TEST_RESULTS_DIR/api"
    
    # 執行 API 測試腳本
    if [ -f "$SCRIPT_DIR/test_api.sh" ]; then
        "$SCRIPT_DIR/test_api.sh" > "$TEST_RESULTS_DIR/api/api_test_results.log" 2>&1
        success "API 測試完成"
    else
        warning "找不到 API 測試腳本: $SCRIPT_DIR/test_api.sh"
    fi
}

# 生成測試報告
generate_reports() {
    log "📊 生成測試報告..."
    
    # 執行報告生成腳本
    if [ -f "$SCRIPT_DIR/generate_test_report.sh" ]; then
        "$SCRIPT_DIR/generate_test_report.sh" generate
    else
        warning "找不到報告生成腳本: $SCRIPT_DIR/generate_test_report.sh"
    fi
}

# 清理測試環境
cleanup_test_environment() {
    log "🧹 清理測試環境..."
    
    # 執行環境清理腳本
    if [ -f "$SCRIPT_DIR/setup_test_env.sh" ]; then
        "$SCRIPT_DIR/setup_test_env.sh" cleanup
    else
        warning "找不到環境清理腳本"
    fi
}

# 顯示測試結果摘要
show_test_summary() {
    log "📋 測試結果摘要"
    echo "=========================================="
    
    # 統計測試結果
    local total_tests=0
    local passed_tests=0
    local failed_tests=0
    
    # 檢查 E2E 測試結果
    if ls "$TEST_RESULTS_DIR/e2e/e2e_report_"*.json >/dev/null 2>&1; then
        local e2e_file=$(ls $TEST_RESULTS_DIR/e2e/e2e_report_*.json | tail -1)
        if [ -f "$e2e_file" ]; then
            local e2e_total=$(jq -r '.test_summary.total_tests' "$e2e_file" 2>/dev/null || echo "0")
            local e2e_passed=$(jq -r '.test_summary.passed' "$e2e_file" 2>/dev/null || echo "0")
            local e2e_failed=$(jq -r '.test_summary.failed' "$e2e_file" 2>/dev/null || echo "0")
            
            total_tests=$((total_tests + e2e_total))
            passed_tests=$((passed_tests + e2e_passed))
            failed_tests=$((failed_tests + e2e_failed))
            
            echo "E2E 測試: $e2e_passed/$e2e_total 通過"
        fi
    fi
    
    # 檢查負載測試結果
    if [ -f "$TEST_RESULTS_DIR/load/load_test_results.json" ]; then
        echo "負載測試: 已完成"
    fi
    
    # 檢查 API 測試結果
    if [ -f "$TEST_RESULTS_DIR/api/api_test_results.log" ]; then
        echo "API 測試: 已完成"
    fi
    
    echo "=========================================="
    echo "總測試數: $total_tests"
    echo "通過: $passed_tests"
    echo "失敗: $failed_tests"
    
    if [ $total_tests -gt 0 ]; then
        local success_rate=$((passed_tests * 100 / total_tests))
        echo "成功率: $success_rate%"
        
        if [ $failed_tests -eq 0 ]; then
            success "🎉 所有測試通過！"
        else
            warning "⚠️  有 $failed_tests 個測試失敗"
        fi
    else
        warning "⚠️  未找到測試結果"
    fi
    
    echo "=========================================="
    echo "測試結果目錄: $TEST_RESULTS_DIR"
    echo "報告目錄: test_reports"
}

# 顯示使用說明
show_usage() {
    echo "E2E 測試整合執行腳本"
    echo ""
    echo "使用方法:"
    echo "  $0 run        - 執行完整測試流程"
    echo "  $0 e2e        - 僅執行 E2E 測試"
    echo "  $0 load       - 僅執行負載測試"
    echo "  $0 api        - 僅執行 API 測試"
    echo "  $0 report     - 僅生成測試報告"
    echo "  $0 cleanup    - 清理測試環境"
    echo "  $0 help       - 顯示此說明"
    echo ""
    echo "範例:"
    echo "  $0 run        # 執行完整的測試流程"
    echo "  $0 e2e        # 僅執行 E2E 測試"
    echo "  $0 report     # 生成測試報告"
}

# 主函數
main() {
    local start_time=$(date +%s)
    
    echo -e "${PURPLE}🚀 Teams Notification API E2E 測試整合執行${NC}"
    echo "=========================================="
    echo "開始時間: $(date)"
    echo "專案目錄: $PROJECT_ROOT"
    echo "腳本目錄: $SCRIPT_DIR"
    echo "=========================================="
    
    case "${1:-run}" in
        "run")
            check_prerequisites
            setup_test_environment
            run_e2e_tests
            run_load_tests
            run_api_tests
            generate_reports
            show_test_summary
            ;;
        "e2e")
            check_prerequisites
            setup_test_environment
            run_e2e_tests
            generate_reports
            show_test_summary
            ;;
        "load")
            check_prerequisites
            setup_test_environment
            run_load_tests
            generate_reports
            show_test_summary
            ;;
        "api")
            check_prerequisites
            setup_test_environment
            run_api_tests
            generate_reports
            show_test_summary
            ;;
        "report")
            generate_reports
            show_test_summary
            ;;
        "cleanup")
            cleanup_test_environment
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        *)
            error "未知命令: $1"
            show_usage
            exit 1
            ;;
    esac
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo ""
    echo "=========================================="
    echo "執行時間: ${duration}秒"
    echo "完成時間: $(date)"
    echo "=========================================="
}

# 執行主函數
main "$@"
