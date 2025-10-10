#!/bin/bash

# TeamsNotifyGoV2 開發部署腳本
# 用於快速部署最新代碼到 Docker 環境

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日誌函數
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
}

# 檢查必要工具
check_requirements() {
    log "檢查必要工具..."
    
    if ! command -v docker &> /dev/null; then
        error "Docker 未安裝或未啟動"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose 未安裝"
        exit 1
    fi
    
    success "必要工具檢查完成"
}

# 停止現有服務
stop_services() {
    log "停止現有服務..."
    
    # 停止本地 Docker 測試環境
    if [ -f "scripts/local-test/docker-compose-local.yml" ]; then
        docker-compose -f scripts/local-test/docker-compose-local.yml down 2>/dev/null || true
    fi
    
    # 停止主 Docker 環境
    make docker-down 2>/dev/null || true
    
    # 停止本地運行的服務
    pkill -f "go run" 2>/dev/null || true
    pkill -f "./bin/server" 2>/dev/null || true
    
    success "服務已停止"
}

# 清理舊的 Docker 映像
cleanup_images() {
    log "清理舊的 Docker 映像..."
    
    # 清理未使用的映像
    docker image prune -f 2>/dev/null || true
    
    # 清理特定的應用映像
    docker rmi teams-notification/api-server:local 2>/dev/null || true
    docker rmi local-test-api-server 2>/dev/null || true
    
    success "Docker 映像清理完成"
}

# 構建應用
build_application() {
    log "構建應用程式..."
    
    # 清理舊的構建文件
    rm -rf bin/
    
    # 構建應用
    make build
    
    if [ $? -eq 0 ]; then
        success "應用程式構建完成"
    else
        error "應用程式構建失敗"
        exit 1
    fi
}

# 啟動 Docker 服務
start_docker_services() {
    log "啟動 Docker 服務..."
    
    # 啟動基礎服務 (PostgreSQL, Redis)
    make docker-up
    
    # 等待服務就緒
    log "等待服務就緒..."
    sleep 10
    
    # 檢查服務狀態
    if docker ps | grep -q "teamsnotify-postgres" && docker ps | grep -q "teamsnotify-redis"; then
        success "Docker 服務啟動完成"
    else
        error "Docker 服務啟動失敗"
        exit 1
    fi
}

# 運行資料庫遷移
run_migrations() {
    log "運行資料庫遷移..."
    
    # 等待資料庫完全就緒
    sleep 5
    
    # 運行遷移
    make db-migrate
    
    if [ $? -eq 0 ]; then
        success "資料庫遷移完成"
    else
        warning "資料庫遷移可能失敗，但繼續執行"
    fi
}

# 啟動應用服務
start_application() {
    log "啟動應用服務..."
    
    # 在背景啟動應用
    make run &
    
    # 等待應用啟動
    sleep 5
    
    # 檢查應用是否正常運行
    if curl -s http://localhost:8080/health > /dev/null; then
        success "應用服務啟動完成"
    else
        error "應用服務啟動失敗"
        exit 1
    fi
}

# 測試新功能
test_new_features() {
    log "測試新功能..."
    
    # 測試健康檢查
    if curl -s http://localhost:8080/health | grep -q "healthy"; then
        success "健康檢查通過"
    else
        error "健康檢查失敗"
        return 1
    fi
    
    # 測試系統監控端點
    if curl -s http://localhost:8080/api/v1/metrics | grep -q "timestamp"; then
        success "系統監控端點正常"
    else
        warning "系統監控端點可能未完全實現"
    fi
    
    # 測試配置端點
    if curl -s http://localhost:8080/api/v1/config | grep -q "server"; then
        success "配置端點正常"
    else
        warning "配置端點可能未完全實現"
    fi
    
    success "功能測試完成"
}

# 顯示服務信息
show_services_info() {
    log "服務信息："
    echo ""
    echo "🌐 API Server: http://localhost:8080"
    echo "📊 Health Check: http://localhost:8080/health"
    echo "📈 Metrics: http://localhost:8080/api/v1/metrics"
    echo "⚙️  Config: http://localhost:8080/api/v1/config"
    echo "🔧 Config Validate: http://localhost:8080/api/v1/config/validate"
    echo ""
    echo "📋 可用命令："
    echo "  make docker-logs    # 查看 Docker 日誌"
    echo "  make docker-down   # 停止 Docker 服務"
    echo "  pkill -f server    # 停止應用服務"
    echo ""
}

# 主函數
main() {
    echo "🚀 TeamsNotifyGoV2 開發部署腳本"
    echo "=================================="
    echo ""
    
    check_requirements
    stop_services
    cleanup_images
    build_application
    start_docker_services
    run_migrations
    start_application
    test_new_features
    show_services_info
    
    success "部署完成！"
}

# 執行主函數
main "$@"
