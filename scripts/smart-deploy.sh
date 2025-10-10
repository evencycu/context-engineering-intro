#!/bin/bash

# TeamsNotifyGoV2 智能部署腳本
# 支持多種部署模式：local process, local docker, persistent docker

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# 日誌函數
log() {
    echo -e "${BLUE}[$(date +'%H:%M:%S')]${NC} $1"
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

info() {
    echo -e "${PURPLE}ℹ️  $1${NC}"
}

# 顯示使用說明
show_usage() {
    echo "🚀 TeamsNotifyGoV2 智能部署腳本"
    echo "=================================="
    echo ""
    echo "用法: $0 [模式] [選項]"
    echo ""
    echo "部署模式:"
    echo "  local-process    - 本地進程部署 (最快，適合開發)"
    echo "  local-docker     - 本地 Docker 部署 (隔離環境)"
    echo "  persistent-docker - 持久化 Docker 部署 (保持數據庫數據)"
    echo ""
    echo "選項:"
    echo "  --clean          - 清理舊的映像和容器"
    echo "  --migrate        - 執行資料庫遷移"
    echo "  --test           - 執行功能測試"
    echo "  --help           - 顯示此說明"
    echo ""
    echo "範例:"
    echo "  $0 local-process                    # 本地進程部署"
    echo "  $0 local-docker --clean --test      # Docker 部署並清理和測試"
    echo "  $0 persistent-docker --migrate      # 持久化部署並遷移"
    echo ""
}

# 解析參數
DEPLOY_MODE=""
CLEAN_MODE=false
MIGRATE_MODE=false
TEST_MODE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        local-process|local-docker|persistent-docker)
            DEPLOY_MODE="$1"
            shift
            ;;
        --clean)
            CLEAN_MODE=true
            shift
            ;;
        --migrate)
            MIGRATE_MODE=true
            shift
            ;;
        --test)
            TEST_MODE=true
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            error "未知參數: $1"
            show_usage
            exit 1
            ;;
    esac
done

# 檢查部署模式
if [ -z "$DEPLOY_MODE" ]; then
    error "請指定部署模式"
    show_usage
    exit 1
fi

# 檢查必要工具
check_requirements() {
    log "檢查必要工具..."
    
    if ! command -v go &> /dev/null; then
        error "Go 未安裝"
        exit 1
    fi
    
    if [ "$DEPLOY_MODE" != "local-process" ]; then
        if ! command -v docker &> /dev/null; then
            error "Docker 未安裝或未啟動"
            exit 1
        fi
        
        if ! command -v docker-compose &> /dev/null; then
            error "Docker Compose 未安裝"
            exit 1
        fi
    fi
    
    success "必要工具檢查完成"
}

# 停止現有服務
stop_services() {
    log "停止現有服務..."
    
    # 停止本地進程
    pkill -f "go run" 2>/dev/null || true
    pkill -f "./bin/server" 2>/dev/null || true
    
    # 停止 Docker 服務
    if [ "$DEPLOY_MODE" != "local-process" ]; then
        docker-compose down 2>/dev/null || true
        docker-compose -f scripts/local-test/docker-compose-local.yml down 2>/dev/null || true
    fi
    
    success "服務已停止"
}

# 清理模式
cleanup_if_needed() {
    if [ "$CLEAN_MODE" = true ]; then
        log "清理舊的映像和容器..."
        
        # 清理未使用的映像
        docker image prune -f 2>/dev/null || true
        
        # 清理特定的應用映像
        docker rmi teams-notification/api-server:local 2>/dev/null || true
        docker rmi local-test-api-server 2>/dev/null || true
        
        success "清理完成"
    fi
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

# 本地進程部署
deploy_local_process() {
    log "部署到本地進程..."
    
    # 啟動資料庫服務 (如果需要的話)
    if [ "$MIGRATE_MODE" = true ]; then
        log "啟動資料庫服務..."
        make docker-up
        sleep 10
        
        log "執行資料庫遷移..."
        make db-migrate
    fi
    
    # 啟動應用
    log "啟動應用服務..."
    make run &
    
    # 等待應用啟動
    sleep 3
    
    success "本地進程部署完成"
}

# 本地 Docker 部署
deploy_local_docker() {
    log "部署到本地 Docker..."
    
    # 啟動 Docker 環境
    ./scripts/local-test/start-docker.sh
    
    success "本地 Docker 部署完成"
}

# 持久化 Docker 部署
deploy_persistent_docker() {
    log "部署到持久化 Docker..."
    
    # 檢查持久化容器是否存在
    if ! docker ps -a | grep -q "teamsnotify-postgres-persistent"; then
        log "創建持久化 PostgreSQL 容器..."
        docker run -d \
            --name teamsnotify-postgres-persistent \
            -e POSTGRES_DB=notification_center \
            -e POSTGRES_USER=teamsnotify \
            -e POSTGRES_PASSWORD=teamsnotify \
            -p 5432:5432 \
            -v teamsnotify_postgres_data:/var/lib/postgresql/data \
            postgres:15-alpine
    else
        log "啟動現有持久化 PostgreSQL 容器..."
        docker start teamsnotify-postgres-persistent
    fi
    
    if ! docker ps -a | grep -q "teamsnotify-redis-persistent"; then
        log "創建持久化 Redis 容器..."
        docker run -d \
            --name teamsnotify-redis-persistent \
            -p 6379:6379 \
            -v teamsnotify_redis_data:/data \
            redis:7-alpine
    else
        log "啟動現有持久化 Redis 容器..."
        docker start teamsnotify-redis-persistent
    fi
    
    # 等待服務就緒
    sleep 5
    
    # 執行遷移
    if [ "$MIGRATE_MODE" = true ]; then
        log "執行資料庫遷移..."
        make db-migrate
    fi
    
    # 啟動應用
    log "啟動應用服務..."
    make run &
    
    # 等待應用啟動
    sleep 3
    
    success "持久化 Docker 部署完成"
}

# 執行測試
run_tests() {
    if [ "$TEST_MODE" = true ]; then
        log "執行功能測試..."
        
        # 等待服務就緒
        sleep 2
        
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
    fi
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
    
    if [ "$DEPLOY_MODE" = "local-process" ]; then
        echo "💡 部署模式: 本地進程 (最快)"
        echo "📋 停止命令: pkill -f 'go run'"
    elif [ "$DEPLOY_MODE" = "local-docker" ]; then
        echo "💡 部署模式: 本地 Docker (隔離)"
        echo "📋 停止命令: docker-compose -f scripts/local-test/docker-compose-local.yml down"
    elif [ "$DEPLOY_MODE" = "persistent-docker" ]; then
        echo "💡 部署模式: 持久化 Docker (保持數據)"
        echo "📋 停止命令: docker stop teamsnotify-postgres-persistent teamsnotify-redis-persistent"
        echo "🗑️  清理命令: docker rm teamsnotify-postgres-persistent teamsnotify-redis-persistent"
    fi
    echo ""
}

# 主函數
main() {
    echo "🚀 TeamsNotifyGoV2 智能部署腳本"
    echo "=================================="
    echo ""
    info "部署模式: $DEPLOY_MODE"
    info "清理模式: $CLEAN_MODE"
    info "遷移模式: $MIGRATE_MODE"
    info "測試模式: $TEST_MODE"
    echo ""
    
    check_requirements
    stop_services
    cleanup_if_needed
    build_application
    
    case $DEPLOY_MODE in
        "local-process")
            deploy_local_process
            ;;
        "local-docker")
            deploy_local_docker
            ;;
        "persistent-docker")
            deploy_persistent_docker
            ;;
    esac
    
    run_tests
    show_services_info
    
    success "部署完成！"
}

# 執行主函數
main "$@"
