#!/bin/bash

# 文檔同步檢查腳本
# 用於檢查文檔與實際實現的同步狀態

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
API_BASE="http://localhost:8080"
DOCS_DIR="docs"
REPORT_FILE="doc_sync_report.md"

echo -e "${BLUE}🔍 開始文檔同步檢查...${NC}"

# 檢查 API 服務是否運行
check_api_running() {
    echo -e "${YELLOW}檢查 API 服務狀態...${NC}"
    if curl -s "$API_BASE/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ API 服務運行正常${NC}"
        return 0
    else
        echo -e "${RED}❌ API 服務未運行，請先啟動服務${NC}"
        return 1
    fi
}

# 檢查監控端點
check_monitoring_endpoints() {
    echo -e "${YELLOW}檢查監控端點...${NC}"
    
    local endpoints=(
        "/health"
        "/api/v1/metrics"
        "/api/v1/config"
        "/api/v1/monitoring/health"
        "/api/v1/monitoring/performance"
        "/api/v1/monitoring/business"
        "/api/v1/monitoring/alerts"
        "/api/v1/monitoring/dashboard"
        "/api/v1/queue/stats"
    )
    
    local working_endpoints=()
    local broken_endpoints=()
    
    for endpoint in "${endpoints[@]}"; do
        if curl -s "$API_BASE$endpoint" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ $endpoint${NC}"
            working_endpoints+=("$endpoint")
        else
            echo -e "${RED}❌ $endpoint${NC}"
            broken_endpoints+=("$endpoint")
        fi
    done
    
    echo -e "\n${BLUE}監控端點檢查結果:${NC}"
    echo -e "${GREEN}正常端點: ${#working_endpoints[@]}${NC}"
    echo -e "${RED}異常端點: ${#broken_endpoints[@]}${NC}"
    
    if [ ${#broken_endpoints[@]} -gt 0 ]; then
        echo -e "${RED}異常端點列表:${NC}"
        for endpoint in "${broken_endpoints[@]}"; do
            echo -e "${RED}  - $endpoint${NC}"
        done
    fi
}

# 檢查文檔中的端點路徑
check_doc_endpoints() {
    echo -e "${YELLOW}檢查文檔中的端點路徑...${NC}"
    
    local doc_files=(
        "docs/05_DEPLOYMENT/MonitoringGuide.md"
        "docs/05_DEPLOYMENT/MonitoringConfiguration.md"
        "docs/06_USER_GUIDE/UserManual.md"
        "docs/02_ARCHITECTURE/Architecture.md"
    )
    
    local outdated_docs=()
    
    for doc_file in "${doc_files[@]}"; do
        if [ -f "$doc_file" ]; then
            # 檢查是否包含舊的端點路徑
            if grep -q "/api/v1/metrics" "$doc_file" && ! grep -q "/api/v1/monitoring/" "$doc_file"; then
                echo -e "${RED}❌ $doc_file 包含過時的端點路徑${NC}"
                outdated_docs+=("$doc_file")
            else
                echo -e "${GREEN}✅ $doc_file 端點路徑正確${NC}"
            fi
        else
            echo -e "${YELLOW}⚠️  $doc_file 不存在${NC}"
        fi
    done
    
    if [ ${#outdated_docs[@]} -gt 0 ]; then
        echo -e "\n${RED}需要更新的文檔:${NC}"
        for doc in "${outdated_docs[@]}"; do
            echo -e "${RED}  - $doc${NC}"
        done
    fi
}

# 生成同步報告
generate_sync_report() {
    echo -e "${YELLOW}生成同步報告...${NC}"
    
    cat > "$REPORT_FILE" << EOF
# 文檔同步檢查報告

**檢查時間**: $(date)
**檢查者**: 文檔同步檢查腳本

## 檢查結果摘要

### API 端點檢查
EOF

    # 檢查 API 端點並記錄結果
    local endpoints=(
        "/health"
        "/api/v1/metrics"
        "/api/v1/config"
        "/api/v1/monitoring/health"
        "/api/v1/monitoring/performance"
        "/api/v1/monitoring/business"
        "/api/v1/monitoring/alerts"
        "/api/v1/monitoring/dashboard"
        "/api/v1/queue/stats"
    )
    
    for endpoint in "${endpoints[@]}"; do
        if curl -s "$API_BASE$endpoint" > /dev/null 2>&1; then
            echo "✅ $endpoint" >> "$REPORT_FILE"
        else
            echo "❌ $endpoint" >> "$REPORT_FILE"
        fi
    done
    
    cat >> "$REPORT_FILE" << EOF

### 文檔檢查
EOF

    # 檢查文檔
    local doc_files=(
        "docs/05_DEPLOYMENT/MonitoringGuide.md"
        "docs/05_DEPLOYMENT/MonitoringConfiguration.md"
        "docs/06_USER_GUIDE/UserManual.md"
        "docs/02_ARCHITECTURE/Architecture.md"
    )
    
    for doc_file in "${doc_files[@]}"; do
        if [ -f "$doc_file" ]; then
            if grep -q "/api/v1/monitoring/" "$doc_file"; then
                echo "✅ $doc_file" >> "$REPORT_FILE"
            else
                echo "❌ $doc_file" >> "$REPORT_FILE"
            fi
        else
            echo "⚠️  $doc_file 不存在" >> "$REPORT_FILE"
        fi
    done
    
    cat >> "$REPORT_FILE" << EOF

## 建議行動

1. **檢查異常端點**: 確保所有 API 端點正常運行
2. **更新過時文檔**: 更新包含過時端點路徑的文檔
3. **定期檢查**: 建議每週運行一次此腳本

## 檢查完成

報告已生成: $REPORT_FILE
EOF

    echo -e "${GREEN}✅ 同步報告已生成: $REPORT_FILE${NC}"
}

# 主函數
main() {
    echo -e "${BLUE}📋 文檔同步檢查工具${NC}"
    echo -e "${BLUE}====================${NC}"
    
    # 檢查 API 服務
    if ! check_api_running; then
        echo -e "${RED}請先啟動 API 服務再進行檢查${NC}"
        exit 1
    fi
    
    echo ""
    
    # 檢查監控端點
    check_monitoring_endpoints
    
    echo ""
    
    # 檢查文檔
    check_doc_endpoints
    
    echo ""
    
    # 生成報告
    generate_sync_report
    
    echo -e "${GREEN}🎉 文檔同步檢查完成！${NC}"
}

# 執行主函數
main "$@"