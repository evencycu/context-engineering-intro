#!/bin/bash

# Teams Notification Service Build Script
# This script builds all services in the Teams Notification project

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="teamsnotifygov2"
VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="build"
DOCKER_REGISTRY=${DOCKER_REGISTRY:-"teams-notification"}
GO_VERSION="1.21"

# Services to build
SERVICES=("api-server" "worker" "admin")

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
        print_error "Go is not installed. Please install Go $GO_VERSION or later."
        exit 1
    fi
    
    if ! command_exists docker; then
        print_error "Docker is not installed. Please install Docker."
        exit 1
    fi
    
    # Check Go version
    GO_CURRENT_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [[ "$GO_CURRENT_VERSION" < "$GO_VERSION" ]]; then
        print_warning "Go version $GO_CURRENT_VERSION is older than recommended $GO_VERSION"
    fi
    
    print_success "Prerequisites check passed"
}

# Function to clean build directory
clean_build() {
    print_status "Cleaning build directory..."
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
    print_success "Build directory cleaned"
}

# Function to download dependencies
download_dependencies() {
    print_status "Downloading Go dependencies..."
    go mod download
    go mod tidy
    print_success "Dependencies downloaded"
}

# Function to run tests
run_tests() {
    print_status "Running tests..."
    go test -v ./...
    print_success "Tests passed"
}

# Function to run linting
run_lint() {
    print_status "Running linter..."
    
    if command_exists golangci-lint; then
        golangci-lint run
    else
        print_warning "golangci-lint not found, skipping linting"
    fi
    
    print_success "Linting completed"
}

# Function to generate protobuf files
generate_proto() {
    print_status "Generating protobuf files..."
    
    if command_exists protoc; then
        # Generate Go files from proto
        protoc --go_out=. --go_opt=paths=source_relative \
               --go-grpc_out=. --go-grpc_opt=paths=source_relative \
               api/proto/notification.proto
        print_success "Protobuf files generated"
    else
        print_warning "protoc not found, skipping protobuf generation"
    fi
}

# Function to build Go binary
build_binary() {
    local service=$1
    local os=$2
    local arch=$3
    
    print_status "Building $service for $os/$arch..."
    
    local output_name="$service"
    if [[ "$os" == "windows" ]]; then
        output_name="$service.exe"
    fi
    
    local output_path="$BUILD_DIR/$service-$os-$arch/$output_name"
    
    CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build \
        -ldflags "-X main.version=$VERSION -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -s -w" \
        -o "$output_path" \
        "./cmd/$service"
    
    print_success "Built $service for $os/$arch"
}

# Function to build all binaries
build_all_binaries() {
    print_status "Building all binaries..."
    
    local platforms=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")
    
    for service in "${SERVICES[@]}"; do
        for platform in "${platforms[@]}"; do
            IFS='/' read -r os arch <<< "$platform"
            build_binary "$service" "$os" "$arch"
        done
    done
    
    print_success "All binaries built"
}

# Function to build Docker image
build_docker_image() {
    local service=$1
    
    print_status "Building Docker image for $service..."
    
    local image_name="$DOCKER_REGISTRY/$service:$VERSION"
    local dockerfile="deployments/$service/Dockerfile"
    
    if [[ -f "$dockerfile" ]]; then
        docker build -f "$dockerfile" -t "$image_name" .
        docker tag "$image_name" "$DOCKER_REGISTRY/$service:latest"
        print_success "Docker image built: $image_name"
    else
        print_warning "Dockerfile not found for $service: $dockerfile"
    fi
}

# Function to build all Docker images
build_all_docker_images() {
    print_status "Building all Docker images..."
    
    for service in "${SERVICES[@]}"; do
        build_docker_image "$service"
    done
    
    print_success "All Docker images built"
}

# Function to create release package
create_release_package() {
    print_status "Creating release package..."
    
    local release_dir="$BUILD_DIR/release"
    mkdir -p "$release_dir"
    
    # Copy binaries
    cp -r "$BUILD_DIR"/*-linux-amd64 "$release_dir/"
    
    # Copy configuration files
    cp -r configs "$release_dir/"
    cp -r deployments "$release_dir/"
    cp -r scripts "$release_dir/"
    
    # Copy documentation
    cp README.md "$release_dir/"
    cp LICENSE "$release_dir/"
    
    # Create archive
    local archive_name="$PROJECT_NAME-$VERSION-linux-amd64.tar.gz"
    tar -czf "$BUILD_DIR/$archive_name" -C "$release_dir" .
    
    print_success "Release package created: $archive_name"
}

# Function to show help
show_help() {
    cat << EOF
Teams Notification Service Build Script

Usage: $0 [OPTIONS]

Options:
    -h, --help          Show this help message
    -v, --version       Set version (default: $VERSION)
    -c, --clean         Clean build directory before building
    -t, --test          Run tests
    -l, --lint          Run linter
    -p, --proto         Generate protobuf files
    -b, --binary        Build Go binaries only
    -d, --docker        Build Docker images only
    -a, --all           Build everything (default)
    --no-test           Skip tests
    --no-lint           Skip linting
    --no-proto          Skip protobuf generation

Examples:
    $0                          # Build everything
    $0 --clean --test           # Clean, test, and build
    $0 --binary                 # Build binaries only
    $0 --docker                 # Build Docker images only
    $0 --version 2.0.0          # Build with specific version

EOF
}

# Main function
main() {
    local clean=false
    local test=true
    local lint=true
    local proto=true
    local binary=true
    local docker=true
    local all=true
    
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
            -c|--clean)
                clean=true
                shift
                ;;
            -t|--test)
                test=true
                all=false
                shift
                ;;
            -l|--lint)
                lint=true
                all=false
                shift
                ;;
            -p|--proto)
                proto=true
                all=false
                shift
                ;;
            -b|--binary)
                binary=true
                all=false
                shift
                ;;
            -d|--docker)
                docker=true
                all=false
                shift
                ;;
            -a|--all)
                all=true
                shift
                ;;
            --no-test)
                test=false
                shift
                ;;
            --no-lint)
                lint=false
                shift
                ;;
            --no-proto)
                proto=false
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_status "Starting build process for Teams Notification Service v$VERSION"
    
    # Check prerequisites
    check_prerequisites
    
    # Clean if requested
    if [[ "$clean" == true ]]; then
        clean_build
    fi
    
    # Download dependencies
    download_dependencies
    
    # Generate protobuf files
    if [[ "$proto" == true ]]; then
        generate_proto
    fi
    
    # Run tests
    if [[ "$test" == true ]]; then
        run_tests
    fi
    
    # Run linting
    if [[ "$lint" == true ]]; then
        run_lint
    fi
    
    # Build binaries
    if [[ "$binary" == true || "$all" == true ]]; then
        build_all_binaries
    fi
    
    # Build Docker images
    if [[ "$docker" == true || "$all" == true ]]; then
        build_all_docker_images
    fi
    
    # Create release package
    if [[ "$all" == true ]]; then
        create_release_package
    fi
    
    print_success "Build process completed successfully!"
    print_status "Build artifacts are available in: $BUILD_DIR/"
}

# Run main function with all arguments
main "$@"
