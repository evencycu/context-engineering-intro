#!/bin/bash

# Teams Notification Service - Deployment Script
# This script deploys the server and OpenAPI service with proper environment configuration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVER_PORT="8080"
OPENAPI_PORT="8082"
CONTAINER_NAME="teams-swagger-ui"
DATABASE_NAME="notification_center"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.21 or later."
        exit 1
    fi
    
    if ! command_exists docker; then
        print_error "Docker is not installed. Please install Docker."
        exit 1
    fi
    
    if ! command_exists curl; then
        print_error "curl is not installed. Please install curl."
        exit 1
    fi
    
    print_success "Prerequisites check passed"
}

# Function to check database connection
check_database() {
    print_status "Checking database connection..."
    
    # Check if postgres container is running
    if ! docker ps --format "table {{.Names}}" | grep -q "teamsnotify-postgres"; then
        print_error "PostgreSQL container 'teamsnotify-postgres' is not running"
        print_status "Please start the database first:"
        print_status "  docker-compose up -d postgres redis"
        exit 1
    fi
    
    # Test database connection
    if ! docker exec teamsnotify-postgres psql -U teamsnotify -d $DATABASE_NAME -c "SELECT 1" > /dev/null 2>&1; then
        print_error "Cannot connect to PostgreSQL database '$DATABASE_NAME'"
        exit 1
    fi
    
    print_success "Database connection OK"
}

# Function to build server
build_server() {
    print_status "Building server..."
    
    # Clean previous build
    rm -f ./server
    
    # Build server
    go build -o server cmd/server/main.go
    if [ $? -ne 0 ]; then
        print_error "Server build failed"
        exit 1
    fi
    
    print_success "Server built successfully"
}

# Function to stop existing services
stop_services() {
    print_status "Stopping existing services..."
    
    # Stop server
    if [ -f "./server.pid" ]; then
        SERVER_PID=$(cat server.pid)
        if ps -p $SERVER_PID > /dev/null 2>&1; then
            print_status "Stopping server (PID: $SERVER_PID)..."
            kill $SERVER_PID 2>/dev/null || true
            sleep 2
        fi
        rm -f server.pid
    fi
    
    # Kill any remaining server processes
    pkill -f './server' 2>/dev/null || true
    
    # Stop OpenAPI container
    if docker ps -a --format "table {{.Names}}" | grep -q "^$CONTAINER_NAME$"; then
        print_status "Stopping OpenAPI container..."
        docker stop $CONTAINER_NAME >/dev/null 2>&1 || true
        docker rm $CONTAINER_NAME >/dev/null 2>&1 || true
    fi
    
    print_success "Services stopped"
}

# Function to start server
start_server() {
    print_status "Starting server..."
    
    # Set environment variables
    export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/$DATABASE_NAME?sslmode=disable"
    export TEAMS_BOT_APP_ID="${TEAMS_BOT_APP_ID:-844146d7-4ac9-4e4d-a463-d6e027714e81}"
    export TEAMS_TENANT_ID="${TEAMS_TENANT_ID:-051cece0-e4dc-4aed-b471-bf29824e1ee6}"
    export TEAMS_BOT_APP_PASSWORD="${TEAMS_BOT_APP_PASSWORD:-HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy}"
    export REDIS_URL="${REDIS_URL:-redis://localhost:6379}"
    
    print_status "Configuration:"
    print_status "  Database: $DATABASE_NAME"
    print_status "  App ID: $TEAMS_BOT_APP_ID"
    print_status "  Tenant ID: $TEAMS_TENANT_ID"
    print_status "  Redis: $REDIS_URL"
    
    # Start server in background
    nohup ./server > server.log 2>&1 &
    SERVER_PID=$!
    echo $SERVER_PID > server.pid
    
    # Wait for server to start
    print_status "Waiting for server to start..."
    sleep 3
    
    # Check if server is running
    if ps -p $SERVER_PID > /dev/null 2>&1; then
        print_success "Server started (PID: $SERVER_PID)"
    else
        print_error "Server failed to start. Check server.log for details:"
        tail -20 server.log
        exit 1
    fi
    
    # Test server health
    print_status "Testing server health..."
    for i in {1..10}; do
        if curl -s http://localhost:$SERVER_PORT/health >/dev/null 2>&1; then
            print_success "Server health check passed"
            break
        fi
        if [ $i -eq 10 ]; then
            print_error "Server health check failed after 10 attempts"
            print_error "Check server.log for details:"
            tail -20 server.log
            exit 1
        fi
        sleep 1
    done
}

# Function to start OpenAPI service
start_openapi() {
    print_status "Starting OpenAPI service..."
    
    # Check if API file exists
    if [ ! -f "api/openapi/teams-notification-api.yaml" ]; then
        print_error "OpenAPI file not found: api/openapi/teams-notification-api.yaml"
        exit 1
    fi
    
    # Start OpenAPI container
    docker run -d \
        --name $CONTAINER_NAME \
        -p $OPENAPI_PORT:8080 \
        -e SWAGGER_JSON=/api/teams-notification-api.yaml \
        -v "$(pwd)/api/openapi:/api" \
        swaggerapi/swagger-ui >/dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        print_success "OpenAPI service started"
    else
        print_error "Failed to start OpenAPI service"
        exit 1
    fi
    
    # Wait for OpenAPI to be ready
    print_status "Waiting for OpenAPI service to be ready..."
    sleep 3
    
    # Test OpenAPI health
    for i in {1..10}; do
        if curl -s http://localhost:$OPENAPI_PORT >/dev/null 2>&1; then
            print_success "OpenAPI service is ready"
            break
        fi
        if [ $i -eq 10 ]; then
            print_warning "OpenAPI service may not be ready yet"
        fi
        sleep 1
    done
}

# Function to show service status
show_status() {
    print_status "Service Status:"
    echo ""
    
    # Server status
    if [ -f "./server.pid" ]; then
        SERVER_PID=$(cat server.pid)
        if ps -p $SERVER_PID > /dev/null 2>&1; then
            print_success "Server: Running (PID: $SERVER_PID)"
        else
            print_error "Server: Not running"
        fi
    else
        print_error "Server: Not running"
    fi
    
    # OpenAPI status
    if docker ps --format "table {{.Names}}" | grep -q "^$CONTAINER_NAME$"; then
        print_success "OpenAPI: Running"
    else
        print_error "OpenAPI: Not running"
    fi
    
    echo ""
    print_status "Service URLs:"
    print_status "  Server API: http://localhost:$SERVER_PORT"
    print_status "  Server Health: http://localhost:$SERVER_PORT/health"
    print_status "  OpenAPI UI: http://localhost:$OPENAPI_PORT"
    print_status "  OpenAPI Spec: http://localhost:$OPENAPI_PORT/api/teams-notification-api.yaml"
    echo ""
}

# Function to show help
show_help() {
    cat << EOF
Teams Notification Service - Deployment Script

Usage: $0 [COMMAND]

Commands:
    deploy      - Deploy both server and OpenAPI service (default)
    server      - Deploy server only
    openapi     - Deploy OpenAPI service only
    stop        - Stop all services
    restart     - Restart all services
    status      - Show service status
    logs        - Show server logs
    help        - Show this help message

Environment Variables:
    TEAMS_BOT_APP_ID        - Teams Bot Application ID
    TEAMS_TENANT_ID         - Teams Tenant ID  
    TEAMS_BOT_APP_PASSWORD  - Teams Bot Application Password
    REDIS_URL               - Redis connection URL
    DATABASE_URL            - PostgreSQL connection URL

Examples:
    $0                      # Deploy everything
    $0 server               # Deploy server only
    $0 openapi              # Deploy OpenAPI only
    $0 restart              # Restart all services
    $0 status               # Check status

EOF
}

# Function to show logs
show_logs() {
    if [ -f "server.log" ]; then
        print_status "Server logs (last 50 lines):"
        tail -50 server.log
    else
        print_warning "No server log file found"
    fi
}

# Main function
main() {
    case "${1:-deploy}" in
        deploy)
            check_prerequisites
            check_database
            stop_services
            build_server
            start_server
            start_openapi
            show_status
            ;;
        server)
            check_prerequisites
            check_database
            stop_services
            build_server
            start_server
            show_status
            ;;
        openapi)
            check_prerequisites
            start_openapi
            show_status
            ;;
        stop)
            stop_services
            print_success "All services stopped"
            ;;
        restart)
            check_prerequisites
            check_database
            stop_services
            build_server
            start_server
            start_openapi
            show_status
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
