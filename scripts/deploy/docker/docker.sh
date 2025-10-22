#!/bin/bash

# Teams Notification Service Docker Script
# This script manages Docker containers and images for the Teams Notification project

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="teams-notification"
VERSION=${VERSION:-"1.0.0"}
DOCKER_REGISTRY=${DOCKER_REGISTRY:-"teams-notification"}
COMPOSE_FILE="scripts/docker/docker-compose.yml"
COMPOSE_DEV_FILE="docker-compose.dev.yml"
COMPOSE_PROD_FILE="docker-compose.prod.yml"

# Services
SERVICES=("apiserver" "postgres" "redis")

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

# Function to check if Docker is running
check_docker() {
    if ! command_exists docker; then
        print_error "Docker is not installed. Please install Docker."
        exit 1
    fi
    
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker."
        exit 1
    fi
    
    print_success "Docker is available"
}

# Function to build Docker image
build_image() {
    local service=$1
    local tag=$2
    
    print_status "Building Docker image for $service:$tag..."
    
    local dockerfile="build/teamsnotification/Dockerfile"
    if [[ ! -f "$dockerfile" ]]; then
        print_error "Dockerfile not found: $dockerfile"
        return 1
    fi
    
    local image_name="$DOCKER_REGISTRY/$service:$tag"
    
    docker build \
        -f "$dockerfile" \
        -t "$image_name" \
        -t "$DOCKER_REGISTRY/$service:latest" \
        .
    
    print_success "Built image: $image_name"
}

# Function to build all images
build_all_images() {
    print_status "Building all Docker images..."
    
    for service in "${SERVICES[@]}"; do
        if [[ "$service" != "postgres" && "$service" != "redis" ]]; then
            build_image "$service" "$VERSION"
        fi
    done
    
    print_success "All images built"
}

# Function to push image to registry
push_image() {
    local service=$1
    local tag=$2
    
    print_status "Pushing image $service:$tag to registry..."
    
    local image_name="$DOCKER_REGISTRY/$service:$tag"
    
    docker push "$image_name"
    docker push "$DOCKER_REGISTRY/$service:latest"
    
    print_success "Pushed image: $image_name"
}

# Function to push all images
push_all_images() {
    print_status "Pushing all Docker images..."
    
    for service in "${SERVICES[@]}"; do
        if [[ "$service" != "postgres" && "$service" != "redis" ]]; then
            push_image "$service" "$VERSION"
        fi
    done
    
    print_success "All images pushed"
}

# Function to start services
start_services() {
    local env=${1:-"dev"}
    
    print_status "Starting services in $env environment..."
    
    local compose_file=""
    case $env in
        "dev")
            compose_file="$COMPOSE_DEV_FILE"
            ;;
        "prod")
            compose_file="$COMPOSE_PROD_FILE"
            ;;
        *)
            compose_file="$COMPOSE_FILE"
            ;;
    esac
    
    if [[ ! -f "$compose_file" ]]; then
        print_error "Compose file not found: $compose_file"
        return 1
    fi
    
    docker-compose -f "$compose_file" up -d
    
    print_success "Services started"
}

# Function to stop services
stop_services() {
    local env=${1:-"dev"}
    
    print_status "Stopping services in $env environment..."
    
    local compose_file=""
    case $env in
        "dev")
            compose_file="$COMPOSE_DEV_FILE"
            ;;
        "prod")
            compose_file="$COMPOSE_PROD_FILE"
            ;;
        *)
            compose_file="$COMPOSE_FILE"
            ;;
    esac
    
    if [[ ! -f "$compose_file" ]]; then
        print_error "Compose file not found: $compose_file"
        return 1
    fi
    
    docker-compose -f "$compose_file" down
    
    print_success "Services stopped"
}

# Function to restart services
restart_services() {
    local env=${1:-"dev"}
    
    print_status "Restarting services in $env environment..."
    
    stop_services "$env"
    start_services "$env"
    
    print_success "Services restarted"
}

# Function to show service status
show_status() {
    local env=${1:-"dev"}
    
    print_status "Service status in $env environment..."
    
    local compose_file=""
    case $env in
        "dev")
            compose_file="$COMPOSE_DEV_FILE"
            ;;
        "prod")
            compose_file="$COMPOSE_PROD_FILE"
            ;;
        *)
            compose_file="$COMPOSE_FILE"
            ;;
    esac
    
    if [[ ! -f "$compose_file" ]]; then
        print_error "Compose file not found: $compose_file"
        return 1
    fi
    
    docker-compose -f "$compose_file" ps
}

# Function to show service logs
show_logs() {
    local service=$1
    local env=${2:-"dev"}
    local follow=${3:-false}
    
    print_status "Showing logs for $service in $env environment..."
    
    local compose_file=""
    case $env in
        "dev")
            compose_file="$COMPOSE_DEV_FILE"
            ;;
        "prod")
            compose_file="$COMPOSE_PROD_FILE"
            ;;
        *)
            compose_file="$COMPOSE_FILE"
            ;;
    esac
    
    if [[ ! -f "$compose_file" ]]; then
        print_error "Compose file not found: $compose_file"
        return 1
    fi
    
    local cmd="docker-compose -f $compose_file logs"
    
    if [[ "$follow" == true ]]; then
        cmd="$cmd -f"
    fi
    
    if [[ -n "$service" ]]; then
        cmd="$cmd $service"
    fi
    
    eval "$cmd"
}

# Function to clean up
cleanup() {
    print_status "Cleaning up Docker resources..."
    
    # Stop all services
    docker-compose -f "$COMPOSE_FILE" down 2>/dev/null || true
    docker-compose -f "$COMPOSE_DEV_FILE" down 2>/dev/null || true
    docker-compose -f "$COMPOSE_PROD_FILE" down 2>/dev/null || true
    
    # Remove unused images
    docker image prune -f
    
    # Remove unused volumes
    docker volume prune -f
    
    # Remove unused networks
    docker network prune -f
    
    print_success "Cleanup completed"
}

# Function to create Docker Compose files
create_compose_files() {
    print_status "Creating Docker Compose files..."
    
    # Development compose file
    cat > "$COMPOSE_DEV_FILE" << 'EOF'
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: teams_notification
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  apiserver:
    build:
      context: .
      dockerfile: build/teamsnotification/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - ENVIRONMENT=development
      - DATABASE_URL=postgres://postgres:password@postgres:5432/teams_notification?sslmode=disable
      - REDIS_URL=redis://redis:6379/0
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    volumes:
      - ./configs:/etc/config:ro
      - ./logs:/var/log

volumes:
  postgres_data:
  redis_data:
EOF

    # Production compose file
    cat > "$COMPOSE_PROD_FILE" << 'EOF'
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${POSTGRES_DB:-teams_notification}
      POSTGRES_USER: ${POSTGRES_USER:-postgres}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  apiserver:
    image: teams-notification/apiserver:latest
    ports:
      - "8080:8080"
    environment:
      - ENVIRONMENT=production
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
      - TEAMS_APP_ID=${TEAMS_APP_ID}
      - TEAMS_APP_PASSWORD=${TEAMS_APP_PASSWORD}
      - TEAMS_TENANT_ID=${TEAMS_TENANT_ID}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    volumes:
      - ./configs:/etc/config:ro
      - ./logs:/var/log
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
EOF

    print_success "Docker Compose files created"
}

# Function to show help
show_help() {
    cat << EOF
Teams Notification Service Docker Script

Usage: $0 [COMMAND] [OPTIONS]

Commands:
    build [SERVICE]     Build Docker image(s)
    push [SERVICE]      Push Docker image(s) to registry
    start [ENV]         Start services (dev|prod)
    stop [ENV]          Stop services (dev|prod)
    restart [ENV]       Restart services (dev|prod)
    status [ENV]        Show service status
    logs [SERVICE] [ENV] Show service logs
    cleanup             Clean up Docker resources
    create-compose      Create Docker Compose files

Options:
    -h, --help          Show this help message
    -v, --version       Set version (default: $VERSION)
    -f, --follow        Follow logs (for logs command)
    -r, --registry      Set Docker registry (default: $DOCKER_REGISTRY)

Examples:
    $0 build                    # Build all images
    $0 build apiserver         # Build apiserver image only
    $0 push                     # Push all images
    $0 start dev                # Start development environment
    $0 stop prod                # Stop production environment
    $0 logs apiserver dev      # Show apiserver logs
    $0 logs apiserver dev -f   # Follow apiserver logs
    $0 status prod              # Show production status
    $0 cleanup                  # Clean up resources

EOF
}

# Main function
main() {
    local command=""
    local service=""
    local env="dev"
    local follow=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -r|--registry)
                DOCKER_REGISTRY="$2"
                shift 2
                ;;
            -f|--follow)
                follow=true
                shift
                ;;
            build|push|start|stop|restart|status|logs|cleanup|create-compose)
                command="$1"
                shift
                ;;
            apiserver|postgres|redis)
                service="$1"
                shift
                ;;
            dev|prod)
                env="$1"
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    if [[ -z "$command" ]]; then
        print_error "No command specified"
        show_help
        exit 1
    fi
    
    print_status "Docker management for Teams Notification Service v$VERSION"
    
    # Check Docker
    check_docker
    
    # Execute command
    case $command in
        build)
            if [[ -n "$service" ]]; then
                build_image "$service" "$VERSION"
            else
                build_all_images
            fi
            ;;
        push)
            if [[ -n "$service" ]]; then
                push_image "$service" "$VERSION"
            else
                push_all_images
            fi
            ;;
        start)
            start_services "$env"
            ;;
        stop)
            stop_services "$env"
            ;;
        restart)
            restart_services "$env"
            ;;
        status)
            show_status "$env"
            ;;
        logs)
            show_logs "$service" "$env" "$follow"
            ;;
        cleanup)
            cleanup
            ;;
        create-compose)
            create_compose_files
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
