#!/bin/bash

# Teams Notification Service Code Generation Script
# This script generates code from various sources (protobuf, OpenAPI, etc.)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROTO_DIR="api/proto"
OPENAPI_DIR="api/openapi"
GENERATED_DIR="internal/generated"
GO_OUT_DIR="internal/generated/proto"

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

# Function to install protoc
install_protoc() {
    print_status "Installing protoc..."
    
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    case $arch in
        x86_64)
            arch="x86_64"
            ;;
        arm64|aarch64)
            arch="aarch_64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    local protoc_version="21.12"
    local protoc_url="https://github.com/protocolbuffers/protobuf/releases/download/v${protoc_version}/protoc-${protoc_version}-${os}-${arch}.zip"
    
    if command_exists curl; then
        curl -LO "$protoc_url"
    elif command_exists wget; then
        wget "$protoc_url"
    else
        print_error "Neither curl nor wget is available. Please install protoc manually."
        exit 1
    fi
    
    unzip "protoc-${protoc_version}-${os}-${arch}.zip" -d /tmp/protoc
    sudo cp -r /tmp/protoc/include/* /usr/local/include/
    sudo cp /tmp/protoc/bin/protoc /usr/local/bin/
    sudo chmod +x /usr/local/bin/protoc
    
    rm -rf /tmp/protoc "protoc-${protoc_version}-${os}-${arch}.zip"
    
    print_success "protoc installed"
}

# Function to install protoc-gen-go
install_protoc_gen_go() {
    print_status "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    print_success "protoc-gen-go installed"
}

# Function to install openapi-generator
install_openapi_generator() {
    print_status "Installing openapi-generator..."
    
    if command_exists npm; then
        npm install -g @openapitools/openapi-generator-cli
    elif command_exists yarn; then
        yarn global add @openapitools/openapi-generator-cli
    else
        print_error "Neither npm nor yarn is available. Please install openapi-generator manually."
        exit 1
    fi
    
    print_success "openapi-generator installed"
}

# Function to generate protobuf files
generate_proto() {
    print_status "Generating protobuf files..."
    
    if ! command_exists protoc; then
        install_protoc
    fi
    
    if ! command_exists protoc-gen-go; then
        install_protoc_gen_go
    fi
    
    # Create output directory
    mkdir -p "$GO_OUT_DIR"
    
    # Find all proto files
    local proto_files=($(find "$PROTO_DIR" -name "*.proto"))
    
    if [[ ${#proto_files[@]} -eq 0 ]]; then
        print_warning "No proto files found in $PROTO_DIR"
        return 0
    fi
    
    # Generate Go files from proto
    for proto_file in "${proto_files[@]}"; do
        print_status "Generating from $proto_file..."
        
        protoc \
            --go_out="$GO_OUT_DIR" \
            --go_opt=paths=source_relative \
            --go-grpc_out="$GO_OUT_DIR" \
            --go-grpc_opt=paths=source_relative \
            --proto_path="$PROTO_DIR" \
            "$proto_file"
    done
    
    print_success "Protobuf files generated"
}

# Function to generate OpenAPI client
generate_openapi_client() {
    print_status "Generating OpenAPI client..."
    
    if ! command_exists openapi-generator; then
        install_openapi_generator
    fi
    
    # Find OpenAPI spec files
    local openapi_files=($(find "$OPENAPI_DIR" -name "*.yaml" -o -name "*.yml" -o -name "*.json"))
    
    if [[ ${#openapi_files[@]} -eq 0 ]]; then
        print_warning "No OpenAPI spec files found in $OPENAPI_DIR"
        return 0
    fi
    
    # Create output directory
    local client_dir="$GENERATED_DIR/openapi"
    mkdir -p "$client_dir"
    
    # Generate Go client from OpenAPI spec
    for spec_file in "${openapi_files[@]}"; do
        print_status "Generating client from $spec_file..."
        
        local spec_name=$(basename "$spec_file" | sed 's/\.[^.]*$//')
        local output_dir="$client_dir/$spec_name"
        
        openapi-generator generate \
            -i "$spec_file" \
            -g go \
            -o "$output_dir" \
            --package-name "$spec_name" \
            --git-user-id "company" \
            --git-repo-id "teamsnotifygov2" \
            --additional-properties=packageName="$spec_name",packageVersion="1.0.0"
    done
    
    print_success "OpenAPI client generated"
}

# Function to generate mocks
generate_mocks() {
    print_status "Generating mocks..."
    
    if ! command_exists mockgen; then
        print_status "Installing mockgen..."
        go install github.com/golang/mock/mockgen@latest
    fi
    
    # Create output directory
    local mocks_dir="$GENERATED_DIR/mocks"
    mkdir -p "$mocks_dir"
    
    # Find interfaces to mock
    local interface_files=($(find internal -name "*.go" -type f))
    
    for file in "${interface_files[@]}"; do
        print_status "Generating mocks for $file..."
        
        # Extract interfaces from file
        local interfaces=$(grep -n "^type.*interface" "$file" | cut -d: -f2 | sed 's/type \([A-Za-z0-9_]*\).*/\1/')
        
        for interface in $interfaces; do
            if [[ -n "$interface" ]]; then
                local mock_file="$mocks_dir/${interface,,}_mock.go"
                print_status "Generating mock for $interface..."
                
                mockgen -source="$file" -destination="$mock_file" -package=mocks
            fi
        done
    done
    
    print_success "Mocks generated"
}

# Function to generate wire files
generate_wire() {
    print_status "Generating wire files..."
    
    if ! command_exists wire; then
        print_status "Installing wire..."
        go install github.com/google/wire/cmd/wire@latest
    fi
    
    # Find wire.go files
    local wire_files=($(find . -name "wire.go" -type f))
    
    for wire_file in "${wire_files[@]}"; do
        print_status "Generating wire for $wire_file..."
        
        local dir=$(dirname "$wire_file")
        (cd "$dir" && wire)
    done
    
    print_success "Wire files generated"
}

# Function to generate swagger docs
generate_swagger() {
    print_status "Generating Swagger documentation..."
    
    if ! command_exists swag; then
        print_status "Installing swag..."
        go install github.com/swaggo/swag/cmd/swag@latest
    fi
    
    # Find main.go files
    local main_files=($(find cmd -name "main.go" -type f))
    
    for main_file in "${main_files[@]}"; do
        print_status "Generating swagger for $main_file..."
        
        local dir=$(dirname "$main_file")
        local service_name=$(basename "$dir")
        
        # Generate swagger docs
        swag init -g "$main_file" -o "docs/$service_name" --parseDependency --parseInternal
    done
    
    print_success "Swagger documentation generated"
}

# Function to clean generated files
clean_generated() {
    print_status "Cleaning generated files..."
    
    rm -rf "$GENERATED_DIR"
    rm -rf "docs"
    
    # Remove generated files in wire directories
    find . -name "wire_gen.go" -type f -delete
    
    print_success "Generated files cleaned"
}

# Function to show help
show_help() {
    cat << EOF
Teams Notification Service Code Generation Script

Usage: $0 [OPTIONS]

Options:
    -h, --help          Show this help message
    -a, --all           Generate all code (default)
    --proto             Generate protobuf files only
    --openapi           Generate OpenAPI client only
    --mocks             Generate mocks only
    --wire              Generate wire files only
    --swagger           Generate Swagger docs only
    --clean             Clean generated files
    --install           Install required tools

Examples:
    $0                          # Generate all code
    $0 --proto                  # Generate protobuf files only
    $0 --openapi                # Generate OpenAPI client only
    $0 --clean                  # Clean generated files
    $0 --install                # Install required tools

EOF
}

# Main function
main() {
    local run_all=true
    local run_proto=false
    local run_openapi=false
    local run_mocks=false
    local run_wire=false
    local run_swagger=false
    local clean=false
    local install_tools=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -a|--all)
                run_all=true
                shift
                ;;
            --proto)
                run_proto=true
                run_all=false
                shift
                ;;
            --openapi)
                run_openapi=true
                run_all=false
                shift
                ;;
            --mocks)
                run_mocks=true
                run_all=false
                shift
                ;;
            --wire)
                run_wire=true
                run_all=false
                shift
                ;;
            --swagger)
                run_swagger=true
                run_all=false
                shift
                ;;
            --clean)
                clean=true
                shift
                ;;
            --install)
                install_tools=true
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_status "Starting code generation for Teams Notification Service"
    
    # Install tools if requested
    if [[ "$install_tools" == true ]]; then
        install_protoc
        install_protoc_gen_go
        install_openapi_generator
        print_success "Tools installed"
        exit 0
    fi
    
    # Clean if requested
    if [[ "$clean" == true ]]; then
        clean_generated
        exit 0
    fi
    
    # Run specific generators or all
    if [[ "$run_all" == true ]]; then
        generate_proto
        generate_openapi_client
        generate_mocks
        generate_wire
        generate_swagger
    else
        [[ "$run_proto" == true ]] && generate_proto
        [[ "$run_openapi" == true ]] && generate_openapi_client
        [[ "$run_mocks" == true ]] && generate_mocks
        [[ "$run_wire" == true ]] && generate_wire
        [[ "$run_swagger" == true ]] && generate_swagger
    fi
    
    print_success "Code generation completed successfully!"
}

# Run main function with all arguments
main "$@"
