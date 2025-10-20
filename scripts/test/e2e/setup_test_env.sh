#!/bin/bash

# E2E 測試環境準備腳本
# 設定測試環境、載入測試資料、驗證服務狀態

set -e

# 配置
BASE_URL="http://localhost:8080"
API_BASE="$BASE_URL/api/v1"
TEST_DATA_DIR="test_data"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日誌函數
log() {
    echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 檢查必要工具
check_prerequisites() {
    log "🔍 檢查必要工具..."
    
    local missing_tools=()
    
    # 檢查 curl
    if ! command -v curl &> /dev/null; then
        missing_tools+=("curl")
    fi
    
    # 檢查 jq
    if ! command -v jq &> /dev/null; then
        missing_tools+=("jq")
    fi
    
    # 檢查 docker
    if ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    fi
    
    # 檢查 docker-compose
    if ! command -v docker-compose &> /dev/null; then
        missing_tools+=("docker-compose")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        error "缺少必要工具: ${missing_tools[*]}"
        error "請安裝缺少的工具後重試"
        exit 1
    fi
    
    success "所有必要工具已安裝"
}

# 啟動測試環境
start_test_environment() {
    log "🚀 啟動測試環境..."
    
    # 檢查 Docker 是否運行
    if ! docker info &> /dev/null; then
        error "Docker 未運行，請啟動 Docker 後重試"
        exit 1
    fi
    
    # 檢查現有服務
    log "檢查現有服務..."
    local postgres_running=false
    local redis_running=false
    
    if docker ps --format "table {{.Names}}" | grep -q "teamsnotify-postgres"; then
        postgres_running=true
        log "PostgreSQL 容器已運行，重用現有服務"
    fi
    
    if docker ps --format "table {{.Names}}" | grep -q "teamsnotify-redis"; then
        redis_running=true
        log "Redis 容器已運行，重用現有服務"
    fi
    
    # 只在需要時啟動服務
    if [ "$postgres_running" = false ] || [ "$redis_running" = false ]; then
        log "啟動 PostgreSQL 和 Redis..."
        docker-compose up -d postgres redis
    fi
    
    # 等待服務就緒
    log "等待服務就緒..."
    sleep 10
    
    # 檢查 PostgreSQL
    local max_attempts=30
    local attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if docker exec teamsnotify-postgres pg_isready -U teamsnotify -d notification_center &> /dev/null; then
            success "PostgreSQL 已就緒"
            break
        fi
        attempt=$((attempt + 1))
        log "等待 PostgreSQL... ($attempt/$max_attempts)"
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        error "PostgreSQL 啟動超時"
        exit 1
    fi
    
    # 檢查 Redis
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if docker exec teamsnotify-redis redis-cli ping | grep -q "PONG"; then
            success "Redis 已就緒"
            break
        fi
        attempt=$((attempt + 1))
        log "等待 Redis... ($attempt/$max_attempts)"
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        error "Redis 啟動超時"
        exit 1
    fi
}

# 執行資料庫遷移
    # 清理現有資料庫結構
    log "清理現有資料庫結構..."
    docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" 2>/dev/null || true
run_database_migrations() {
    log "📊 執行資料庫遷移..."
    # 清理現有資料庫結構
    log "清理現有資料庫結構..."
    docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" 2>/dev/null || true
    
    # 檢查遷移檔案是否存在（更新為現有路徑）
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
    
    if [ ! -f "$PROJECT_ROOT/scripts/database/schema.sql" ]; then
        error "找不到資料庫結構檔案: $PROJECT_ROOT/scripts/database/schema.sql"
        exit 1
    fi
    
    # 確保目標資料庫存在（舊資料卷可能未包含 notification_center 資料庫）
    if ! docker exec teamsnotify-postgres psql -U teamsnotify -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='notification_center'" | grep -q 1; then
        log "建立資料庫 notification_center..."
        docker exec teamsnotify-postgres psql -U teamsnotify -d postgres -c "CREATE DATABASE notification_center;"
    fi

    # 執行遷移
    docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < "$PROJECT_ROOT/scripts/database/schema.sql"
    
    if [ $? -eq 0 ]; then
        success "資料庫遷移完成"
    else
        error "資料庫遷移失敗"
        exit 1
    fi
    
    # 載入初始數據
    if [ -f "$PROJECT_ROOT/scripts/database/init.sql" ]; then
        log "載入初始數據..."
        docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < "$PROJECT_ROOT/scripts/database/init.sql"
        if [ $? -eq 0 ]; then
            success "初始數據載入完成"
        else
            warning "初始數據載入失敗，繼續執行測試"
        fi
    else
        warning "找不到初始數據檔案: $PROJECT_ROOT/scripts/database/init.sql"
    fi
}

# 載入測試資料
load_test_data() {
    log "📋 載入測試資料..."
    
    # 建立測試資料目錄
    mkdir -p "$TEST_DATA_DIR"
    
    # 建立測試公司
    log "建立測試公司..."
    local company_response=$(curl -s -X POST "$API_BASE/companies" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "E2E Test Company",
            "contact_email": "e2e-test@example.com",
            "contact_phone": "+886-2-1234-5678",
            "address": "台北市信義區信義路五段7號"
        }')
    
    local company_id=$(echo "$company_response" | jq -r '.data.id')
    if [ "$company_id" = "null" ] || [ -z "$company_id" ]; then
        error "建立測試公司失敗"
        echo "$company_response" | jq '.'
        exit 1
    fi
    
    success "測試公司已建立: $company_id"
    
    # 建立測試用戶
    log "建立測試用戶..."
    local user_response=$(curl -s -X POST "$API_BASE/users" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"e2e_test_user\",
            \"email\": \"e2e-test@example.com\",
            \"company_id\": \"$company_id\",
            \"role\": \"user\",
            \"status\": \"active\"
        }")
    
    local user_id=$(echo "$user_response" | jq -r '.data.id')
    if [ "$user_id" = "null" ] || [ -z "$user_id" ]; then
        error "建立測試用戶失敗"
        echo "$user_response" | jq '.'
        exit 1
    fi
    
    success "測試用戶已建立: $user_id"
    
    # 建立測試專案
    log "建立測試專案..."
    local project_response=$(curl -s -X POST "$API_BASE/projects" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\": \"E2E Test Project\",
            \"description\": \"E2E 測試專案\",
            \"company_id\": \"$company_id\",
            \"created_by\": \"$user_id\"
        }")
    
    local project_id=$(echo "$project_response" | jq -r '.data.id')
    local notify_key=$(echo "$project_response" | jq -r '.data.notify_key')
    
    if [ "$project_id" = "null" ] || [ -z "$project_id" ]; then
        error "建立測試專案失敗"
        echo "$project_response" | jq '.'
        exit 1
    fi
    
    success "測試專案已建立: $project_id (notify_key: $notify_key)"
    
    # 建立測試目的地
    log "建立測試目的地..."
    # 從資料庫取得最新且 active 的 conversation_id 以確保正確
    local conv_id=$(docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -tAc "SELECT conversation_id FROM bot_installations WHERE installation_status='active' ORDER BY installed_at DESC LIMIT 1;")
    conv_id=$(echo "$conv_id" | tr -d '[:space:]')
    local destination_response=$(curl -s -X POST "$API_BASE/destinations" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\": \"E2E Test Destination\",
            \"description\": \"E2E 測試目的地\",
            \"project_id\": \"$project_id\",
            \"teams_tenant_id\": \"051cece0-e4dc-4aed-b471-bf29824e1ee6\",
            \"targets\": [
                {
                    \"type\": \"personal\",
                    \"conversation_id\": \"$conv_id\",
                    \"tenant_id\": \"051cece0-e4dc-4aed-b471-bf29824e1ee6\"
                }
            ]
        }")
    
    local destination_id=$(echo "$destination_response" | jq -r '.data.id')
    if [ "$destination_id" = "null" ] || [ -z "$destination_id" ]; then
        error "建立測試目的地失敗"
        echo "$destination_response" | jq '.'
        exit 1
    fi
    
    success "測試目的地已建立: $destination_id"
    
    # 儲存測試資料
    cat > "$TEST_DATA_DIR/test_data.env" << EOF
COMPANY_ID=$company_id
USER_ID=$user_id
PROJECT_ID=$project_id
NOTIFY_KEY=$notify_key
DESTINATION_ID=$destination_id
EOF
    
    success "測試資料已儲存到 $TEST_DATA_DIR/test_data.env"
}

# 驗證測試環境
verify_test_environment() {
    log "🔍 驗證測試環境..."
    
    # 檢查 API 服務
    if curl -sS "$BASE_URL/health" > /dev/null; then
        success "API 服務正常"
    else
        error "API 服務未運行"
        exit 1
    fi
    
    # 檢查資料庫連線
    if docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT 1;" > /dev/null 2>&1; then
        success "資料庫連線正常"
    else
        error "資料庫連線失敗"
        exit 1
    fi
    
    # 檢查 Redis 連線
    if docker exec teamsnotify-redis redis-cli ping | grep -q "PONG"; then
        success "Redis 連線正常"
    else
        error "Redis 連線失敗"
        exit 1
    fi
    
    # 檢查佇列狀態
    local queue_status=$(curl -sS "$API_BASE/queue/stats")
    if echo "$queue_status" | jq -r '.circuit_state' | grep -q "closed"; then
        success "佇列系統正常"
    else
        warning "佇列系統狀態異常"
        echo "$queue_status" | jq '.'
    fi
}

# 清理測試環境
cleanup_test_environment() {
    log "🧹 清理測試環境..."
    
    # 停止容器
    docker-compose down
    
    success "測試環境已清理"
}

# 顯示使用說明
show_usage() {
    echo "E2E 測試環境準備腳本"
    echo ""
    echo "使用方法:"
    echo "  $0 setup     - 設定測試環境"
    echo "  $0 verify    - 驗證測試環境"
    echo "  $0 cleanup   - 清理測試環境"
    echo "  $0 help      - 顯示此說明"
    echo ""
    echo "範例:"
    echo "  $0 setup     # 設定完整的測試環境"
    echo "  $0 verify    # 驗證環境是否正常"
    echo "  $0 cleanup   # 清理測試環境"
}

# 主函數
main() {
    case "${1:-setup}" in
        "setup")
            echo -e "${BLUE}🚀 設定 E2E 測試環境${NC}"
            echo "=========================================="
            check_prerequisites
            start_test_environment
            run_database_migrations
            load_test_data
            verify_test_environment
            success "✅ 測試環境設定完成！"
            ;;
        "verify")
            echo -e "${BLUE}🔍 驗證測試環境${NC}"
            echo "=========================================="
            verify_test_environment
            success "✅ 測試環境驗證通過！"
            ;;
        "cleanup")
            echo -e "${BLUE}🧹 清理測試環境${NC}"
            echo "=========================================="
            cleanup_test_environment
            success "✅ 測試環境已清理！"
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
}

# 執行主函數
main "$@"
