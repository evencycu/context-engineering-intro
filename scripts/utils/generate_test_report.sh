#!/bin/bash

# E2E 測試報告生成腳本
# 分析測試結果並生成詳細報告

set -e

# 配置
TEST_RESULTS_DIR="test_results"
REPORT_DIR="test_reports"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# 建立報告目錄
mkdir -p "$REPORT_DIR"

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

# 分析測試結果
analyze_test_results() {
    log "📊 分析測試結果..."
    
    local e2e_results=""
    local load_results=""
    
    # 尋找 E2E 測試結果
    if [ -f "$TEST_RESULTS_DIR/e2e/e2e_report_*.json" ]; then
        e2e_results=$(ls $TEST_RESULTS_DIR/e2e/e2e_report_*.json | tail -1)
    fi
    
    # 尋找負載測試結果
    if [ -f "$TEST_RESULTS_DIR/load_test_summary.json" ]; then
        load_results="$TEST_RESULTS_DIR/load_test_summary.json"
    fi
    
    # 生成綜合報告
    generate_comprehensive_report "$e2e_results" "$load_results"
}

# 生成綜合報告
generate_comprehensive_report() {
    local e2e_file="$1"
    local load_file="$2"
    
    log "📋 生成綜合測試報告..."
    
    local report_file="$REPORT_DIR/test_report_$TIMESTAMP.html"
    
    # 讀取 E2E 測試結果
    local e2e_summary=""
    if [ -n "$e2e_file" ] && [ -f "$e2e_file" ]; then
        e2e_summary=$(cat "$e2e_file")
    fi
    
    # 讀取負載測試結果
    local load_summary=""
    if [ -n "$load_file" ] && [ -f "$load_file" ]; then
        load_summary=$(cat "$load_file")
    fi
    
    # 生成 HTML 報告
    cat > "$report_file" << EOF
<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Teams Notification API 測試報告</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 2.5em;
        }
        .header p {
            margin: 10px 0 0 0;
            opacity: 0.9;
        }
        .content {
            padding: 30px;
        }
        .section {
            margin-bottom: 30px;
        }
        .section h2 {
            color: #333;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        .metrics-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin: 20px 0;
        }
        .metric-card {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 20px;
            text-align: center;
        }
        .metric-value {
            font-size: 2em;
            font-weight: bold;
            color: #667eea;
        }
        .metric-label {
            color: #666;
            margin-top: 5px;
        }
        .status-pass {
            color: #28a745;
        }
        .status-fail {
            color: #dc3545;
        }
        .status-warning {
            color: #ffc107;
        }
        .test-results {
            background: #f8f9fa;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }
        .test-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px 0;
            border-bottom: 1px solid #e9ecef;
        }
        .test-item:last-child {
            border-bottom: none;
        }
        .test-name {
            font-weight: 500;
        }
        .test-status {
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.9em;
            font-weight: 500;
        }
        .status-pass {
            background: #d4edda;
            color: #155724;
        }
        .status-fail {
            background: #f8d7da;
            color: #721c24;
        }
        .footer {
            background: #f8f9fa;
            padding: 20px;
            text-align: center;
            color: #666;
            border-top: 1px solid #e9ecef;
        }
        .chart-container {
            background: white;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }
        pre {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 4px;
            padding: 15px;
            overflow-x: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🧪 Teams Notification API 測試報告</h1>
            <p>測試執行時間: $(date -Iseconds)</p>
        </div>
        
        <div class="content">
            <div class="section">
                <h2>📊 測試摘要</h2>
                <div class="metrics-grid">
                    <div class="metric-card">
                        <div class="metric-value" id="total-tests">-</div>
                        <div class="metric-label">總測試數</div>
                    </div>
                    <div class="metric-card">
                        <div class="metric-value status-pass" id="passed-tests">-</div>
                        <div class="metric-label">通過測試</div>
                    </div>
                    <div class="metric-card">
                        <div class="metric-value status-fail" id="failed-tests">-</div>
                        <div class="metric-label">失敗測試</div>
                    </div>
                    <div class="metric-card">
                        <div class="metric-value" id="success-rate">-</div>
                        <div class="metric-label">成功率</div>
                    </div>
                </div>
            </div>
            
            <div class="section">
                <h2>🔍 E2E 測試結果</h2>
                <div class="test-results" id="e2e-results">
                    <p>正在載入 E2E 測試結果...</p>
                </div>
            </div>
            
            <div class="section">
                <h2>⚡ 效能測試結果</h2>
                <div class="test-results" id="load-results">
                    <p>正在載入效能測試結果...</p>
                </div>
            </div>
            
            <div class="section">
                <h2>📈 效能指標</h2>
                <div class="chart-container">
                    <h3>回應時間分佈</h3>
                    <div id="response-time-chart">
                        <p>回應時間統計將在此顯示</p>
                    </div>
                </div>
            </div>
            
            <div class="section">
                <h2>🔧 測試環境資訊</h2>
                <pre>
測試環境: $(hostname)
作業系統: $(uname -s) $(uname -r)
測試時間: $(date)
測試工具: k6, curl, jq
                </pre>
            </div>
        </div>
        
        <div class="footer">
            <p>Teams Notification API 測試報告 - 生成時間: $(date)</p>
        </div>
    </div>
    
    <script>
        // 載入測試結果
        function loadTestResults() {
            // 這裡可以載入實際的測試結果數據
            console.log('載入測試結果...');
        }
        
        // 初始化
        document.addEventListener('DOMContentLoaded', function() {
            loadTestResults();
        });
    </script>
</body>
</html>
EOF
    
    success "HTML 報告已生成: $report_file"
    
    # 生成 JSON 報告
    generate_json_report "$e2e_summary" "$load_summary"
    
    # 生成文字報告
    generate_text_report "$e2e_summary" "$load_summary"
}

# 生成 JSON 報告
generate_json_report() {
    local e2e_summary="$1"
    local load_summary="$2"
    
    local json_file="$REPORT_DIR/test_report_$TIMESTAMP.json"
    
    cat > "$json_file" << EOF
{
  "report_metadata": {
    "timestamp": "$(date -Iseconds)",
    "report_version": "1.0",
    "test_environment": "$(hostname)",
    "generated_by": "generate_test_report.sh"
  },
  "e2e_test_results": $e2e_summary,
  "load_test_results": $load_summary,
  "summary": {
    "total_tests": 0,
    "passed_tests": 0,
    "failed_tests": 0,
    "success_rate": "0%",
    "overall_status": "unknown"
  }
}
EOF
    
    success "JSON 報告已生成: $json_file"
}

# 生成文字報告
generate_text_report() {
    local e2e_summary="$1"
    local load_summary="$2"
    
    local text_file="$REPORT_DIR/test_report_$TIMESTAMP.txt"
    
    cat > "$text_file" << EOF
========================================
Teams Notification API 測試報告
========================================
生成時間: $(date)
測試環境: $(hostname)
報告版本: 1.0

📊 測試摘要
========================================
總測試數: 0
通過測試: 0
失敗測試: 0
成功率: 0%

🔍 E2E 測試結果
========================================
$e2e_summary

⚡ 效能測試結果
========================================
$load_summary

🔧 測試環境資訊
========================================
作業系統: $(uname -s) $(uname -r)
測試工具: k6, curl, jq
測試目錄: $TEST_RESULTS_DIR
報告目錄: $REPORT_DIR

📋 建議事項
========================================
1. 檢查失敗的測試案例
2. 分析效能瓶頸
3. 優化系統配置
4. 更新測試案例

========================================
EOF
    
    success "文字報告已生成: $text_file"
}

# 顯示報告摘要
show_report_summary() {
    log "📋 測試報告摘要"
    echo "=========================================="
    echo "報告目錄: $REPORT_DIR"
    echo "生成時間: $(date)"
    echo ""
    
    # 列出生成的報告檔案
    if [ -d "$REPORT_DIR" ]; then
        echo "生成的報告檔案:"
        ls -la "$REPORT_DIR"/*.html 2>/dev/null || echo "  - 無 HTML 報告"
        ls -la "$REPORT_DIR"/*.json 2>/dev/null || echo "  - 無 JSON 報告"
        ls -la "$REPORT_DIR"/*.txt 2>/dev/null || echo "  - 無文字報告"
    fi
    
    echo ""
    echo "查看報告:"
    echo "  HTML: open $REPORT_DIR/test_report_$TIMESTAMP.html"
    echo "  JSON: cat $REPORT_DIR/test_report_$TIMESTAMP.json"
    echo "  文字: cat $REPORT_DIR/test_report_$TIMESTAMP.txt"
}

# 顯示使用說明
show_usage() {
    echo "測試報告生成腳本"
    echo ""
    echo "使用方法:"
    echo "  $0 generate  - 生成測試報告"
    echo "  $0 summary   - 顯示報告摘要"
    echo "  $0 help      - 顯示此說明"
    echo ""
    echo "範例:"
    echo "  $0 generate  # 生成完整的測試報告"
    echo "  $0 summary   # 顯示報告摘要"
}

# 主函數
main() {
    case "${1:-generate}" in
        "generate")
            echo -e "${BLUE}📊 生成測試報告${NC}"
            echo "=========================================="
            analyze_test_results
            show_report_summary
            ;;
        "summary")
            show_report_summary
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        *)
            echo -e "${RED}未知命令: $1${NC}"
            show_usage
            exit 1
            ;;
    esac
}

# 執行主函數
main "$@"
