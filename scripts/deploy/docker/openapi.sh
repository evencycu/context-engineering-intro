#!/bin/bash

# Teams Notification API - OpenAPI Server Management Script

# 顏色定義
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

CONTAINER_NAME="teams-swagger-ui"
PORT="8082"
API_FILE="api/openapi/teams-notification-api.yaml"

# 函數：顯示幫助信息
show_help() {
    echo -e "${BLUE}Teams Notification API - OpenAPI Server Management${NC}"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start     - 啟動 OpenAPI 服務器 (Swagger UI)"
    echo "  stop      - 停止 OpenAPI 服務器"
    echo "  restart   - 重啟 OpenAPI 服務器"
    echo "  status    - 檢查服務器狀態"
    echo "  logs      - 查看服務器日誌"
    echo "  open      - 在瀏覽器中打開 Swagger UI"
    echo "  help      - 顯示此幫助信息"
    echo ""
    echo "服務器信息:"
    echo "  URL: http://localhost:$PORT"
    echo "  API 文件: $API_FILE"
    echo "  容器名稱: $CONTAINER_NAME"
}

# 函數：檢查容器是否存在
container_exists() {
    docker ps -a --format "table {{.Names}}" | grep -q "^$CONTAINER_NAME$"
}

# 函數：檢查容器是否運行
container_running() {
    docker ps --format "table {{.Names}}" | grep -q "^$CONTAINER_NAME$"
}

# 函數：啟動服務器
start_server() {
    echo -e "${YELLOW}啟動 OpenAPI 服務器...${NC}"
    
    # 檢查 API 文件是否存在
    if [ ! -f "$API_FILE" ]; then
        echo -e "${RED}錯誤: API 文件不存在: $API_FILE${NC}"
        exit 1
    fi
    
    # 如果容器已存在，先停止
    if container_exists; then
        echo -e "${YELLOW}停止現有容器...${NC}"
        docker stop $CONTAINER_NAME >/dev/null 2>&1
        docker rm $CONTAINER_NAME >/dev/null 2>&1
    fi
    
    # 啟動新容器
    echo -e "${YELLOW}啟動 Swagger UI 容器...${NC}"
    docker run -d \
        --name $CONTAINER_NAME \
        -p $PORT:8080 \
        -e SWAGGER_JSON=/api/teams-notification-api.yaml \
        -v "$(pwd)/api/openapi:/api" \
        swaggerapi/swagger-ui >/dev/null
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ OpenAPI 服務器啟動成功！${NC}"
        echo -e "${BLUE}🌐 Swagger UI: http://localhost:$PORT${NC}"
        echo -e "${BLUE}📄 API 規範: http://localhost:$PORT/api/teams-notification-api.yaml${NC}"
    else
        echo -e "${RED}❌ 啟動失敗${NC}"
        exit 1
    fi
}

# 函數：停止服務器
stop_server() {
    echo -e "${YELLOW}停止 OpenAPI 服務器...${NC}"
    
    if container_exists; then
        docker stop $CONTAINER_NAME >/dev/null 2>&1
        docker rm $CONTAINER_NAME >/dev/null 2>&1
        echo -e "${GREEN}✅ OpenAPI 服務器已停止${NC}"
    else
        echo -e "${YELLOW}⚠️  容器不存在${NC}"
    fi
}

# 函數：重啟服務器
restart_server() {
    echo -e "${YELLOW}重啟 OpenAPI 服務器...${NC}"
    stop_server
    sleep 2
    start_server
}

# 函數：檢查狀態
check_status() {
    echo -e "${BLUE}OpenAPI 服務器狀態:${NC}"
    echo ""
    
    if container_exists; then
        if container_running; then
            echo -e "${GREEN}✅ 狀態: 運行中${NC}"
            echo -e "${BLUE}🌐 URL: http://localhost:$PORT${NC}"
            echo -e "${BLUE}📄 API 文件: $API_FILE${NC}"
            
            # 檢查健康狀態
            if curl -s http://localhost:$PORT >/dev/null 2>&1; then
                echo -e "${GREEN}✅ 健康檢查: 通過${NC}"
            else
                echo -e "${RED}❌ 健康檢查: 失敗${NC}"
            fi
        else
            echo -e "${YELLOW}⚠️  狀態: 已停止${NC}"
        fi
    else
        echo -e "${RED}❌ 狀態: 未創建${NC}"
    fi
}

# 函數：查看日誌
show_logs() {
    if container_exists; then
        echo -e "${BLUE}OpenAPI 服務器日誌:${NC}"
        docker logs $CONTAINER_NAME
    else
        echo -e "${RED}❌ 容器不存在${NC}"
    fi
}

# 函數：在瀏覽器中打開
open_browser() {
    if container_running; then
        echo -e "${YELLOW}在瀏覽器中打開 Swagger UI...${NC}"
        
        # 檢測操作系統
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            open "http://localhost:$PORT"
        elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
            # Linux
            xdg-open "http://localhost:$PORT"
        elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
            # Windows
            start "http://localhost:$PORT"
        else
            echo -e "${YELLOW}請手動打開瀏覽器訪問: http://localhost:$PORT${NC}"
        fi
    else
        echo -e "${RED}❌ 服務器未運行，請先啟動服務器${NC}"
    fi
}

# 主函數
main() {
    case "${1:-help}" in
        start)
            start_server
            ;;
        stop)
            stop_server
            ;;
        restart)
            restart_server
            ;;
        status)
            check_status
            ;;
        logs)
            show_logs
            ;;
        open)
            open_browser
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            echo -e "${RED}未知命令: $1${NC}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 執行主函數
main "$@"
