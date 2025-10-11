#!/bin/bash

# Teams Notification Service Lint Script
# This script runs various linting tools on the codebase

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
LINT_DIRS=("cmd" "internal" "pkg" "api")
EXCLUDE_DIRS=("vendor" "build" "test" "deployments")

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

# Function to install golangci-lint
install_golangci_lint() {
    print_status "Installing golangci-lint..."
    
    if command_exists curl; then
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
    elif command_exists wget; then
        wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
    else
        print_error "Neither curl nor wget is available. Please install golangci-lint manually."
        exit 1
    fi
    
    print_success "golangci-lint installed"
}

# Function to run go vet
run_go_vet() {
    print_status "Running go vet..."
    
    local vet_errors=0
    
    for dir in "${LINT_DIRS[@]}"; do
        if [[ -d "$dir" ]]; then
            print_status "Checking $dir..."
            if ! go vet ./"$dir"/...; then
                vet_errors=$((vet_errors + 1))
            fi
        fi
    done
    
    if [[ $vet_errors -eq 0 ]]; then
        print_success "go vet passed"
    else
        print_error "go vet found $vet_errors issues"
        return 1
    fi
}

# Function to run go fmt
run_go_fmt() {
    print_status "Running go fmt..."
    
    local fmt_errors=0
    
    for dir in "${LINT_DIRS[@]}"; do
        if [[ -d "$dir" ]]; then
            print_status "Checking $dir..."
            local unformatted=$(go fmt ./"$dir"/...)
            if [[ -n "$unformatted" ]]; then
                print_warning "Unformatted files in $dir:"
                echo "$unformatted"
                fmt_errors=$((fmt_errors + 1))
            fi
        fi
    done
    
    if [[ $fmt_errors -eq 0 ]]; then
        print_success "go fmt passed"
    else
        print_warning "go fmt found $fmt_errors unformatted files"
        return 1
    fi
}

# Function to run goimports
run_goimports() {
    print_status "Running goimports..."
    
    if ! command_exists goimports; then
        print_status "Installing goimports..."
        go install golang.org/x/tools/cmd/goimports@latest
    fi
    
    local import_errors=0
    
    for dir in "${LINT_DIRS[@]}"; do
        if [[ -d "$dir" ]]; then
            print_status "Checking $dir..."
            local unformatted=$(goimports -l ./"$dir"/...)
            if [[ -n "$unformatted" ]]; then
                print_warning "Unformatted imports in $dir:"
                echo "$unformatted"
                import_errors=$((import_errors + 1))
            fi
        fi
    done
    
    if [[ $import_errors -eq 0 ]]; then
        print_success "goimports passed"
    else
        print_warning "goimports found $import_errors files with unformatted imports"
        return 1
    fi
}

# Function to run golangci-lint
run_golangci_lint() {
    print_status "Running golangci-lint..."
    
    if ! command_exists golangci-lint; then
        install_golangci_lint
    fi
    
    # Create .golangci.yml if it doesn't exist
    if [[ ! -f ".golangci.yml" ]]; then
        print_status "Creating .golangci.yml configuration..."
        cat > .golangci.yml << 'EOF'
linters-settings:
  gofmt:
    simplify: true
  goimports:
    local-prefixes: github.com/company/teamsnotifygov2
  gocyclo:
    min-complexity: 15
  goconst:
    min-len: 2
    min-occurrences: 2
  misspell:
    locale: US
  lll:
    line-length: 120
  gomnd:
    settings:
      mnd:
        checks: argument,case,condition,operation,return,assign
  gosec:
    severity: medium
    confidence: medium
  govet:
    check-shadowing: true
  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style
    disabled-checks:
      - dupImport
      - ifElseChain
      - octalLiteral
      - whyNoLint
      - wrapperFunc

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - exportloopref
    - funlen
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - misspell
    - nakedret
    - noctx
    - nolintlint
    - rowserrcheck
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace

  disable:
    - gochecknoglobals
    - gochecknoglobals

run:
  timeout: 5m
  issues-exit-code: 1
  tests: true
  modules-download-mode: readonly

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gomnd
        - goconst
        - funlen
    - path: cmd/
      linters:
        - gomnd
    - path: api/
      linters:
        - gomnd
        - goconst
  max-issues-per-linter: 0
  max-same-issues: 0
EOF
    fi
    
    local lint_errors=0
    
    for dir in "${LINT_DIRS[@]}"; do
        if [[ -d "$dir" ]]; then
            print_status "Linting $dir..."
            if ! golangci-lint run ./"$dir"/...; then
                lint_errors=$((lint_errors + 1))
            fi
        fi
    done
    
    if [[ $lint_errors -eq 0 ]]; then
        print_success "golangci-lint passed"
    else
        print_error "golangci-lint found $lint_errors issues"
        return 1
    fi
}

# Function to run staticcheck
run_staticcheck() {
    print_status "Running staticcheck..."
    
    if ! command_exists staticcheck; then
        print_status "Installing staticcheck..."
        go install honnef.co/go/tools/cmd/staticcheck@latest
    fi
    
    local staticcheck_errors=0
    
    for dir in "${LINT_DIRS[@]}"; do
        if [[ -d "$dir" ]]; then
            print_status "Checking $dir..."
            if ! staticcheck ./"$dir"/...; then
                staticcheck_errors=$((staticcheck_errors + 1))
            fi
        fi
    done
    
    if [[ $staticcheck_errors -eq 0 ]]; then
        print_success "staticcheck passed"
    else
        print_error "staticcheck found $staticcheck_errors issues"
        return 1
    fi
}

# Function to run ineffassign
run_ineffassign() {
    print_status "Running ineffassign..."
    
    if ! command_exists ineffassign; then
        print_status "Installing ineffassign..."
        go install github.com/gordonklaus/ineffassign@latest
    fi
    
    local ineffassign_errors=0
    
    for dir in "${LINT_DIRS[@]}"; do
        if [[ -d "$dir" ]]; then
            print_status "Checking $dir..."
            if ! ineffassign ./"$dir"/...; then
                ineffassign_errors=$((ineffassign_errors + 1))
            fi
        fi
    done
    
    if [[ $ineffassign_errors -eq 0 ]]; then
        print_success "ineffassign passed"
    else
        print_error "ineffassign found $ineffassign_errors issues"
        return 1
    fi
}

# Function to run misspell
run_misspell() {
    print_status "Running misspell..."
    
    if ! command_exists misspell; then
        print_status "Installing misspell..."
        go install github.com/client9/misspell/cmd/misspell@latest
    fi
    
    local misspell_errors=0
    
    # Check Go files
    for dir in "${LINT_DIRS[@]}"; do
        if [[ -d "$dir" ]]; then
            print_status "Checking $dir..."
            if ! misspell -error ./"$dir"/...; then
                misspell_errors=$((misspell_errors + 1))
            fi
        fi
    done
    
    # Check documentation files
    local doc_files=("README.md" "ARCHITECTURE.md" "INITIAL.md")
    for file in "${doc_files[@]}"; do
        if [[ -f "$file" ]]; then
            print_status "Checking $file..."
            if ! misspell -error "$file"; then
                misspell_errors=$((misspell_errors + 1))
            fi
        fi
    done
    
    if [[ $misspell_errors -eq 0 ]]; then
        print_success "misspell passed"
    else
        print_error "misspell found $misspell_errors issues"
        return 1
    fi
}

# Function to run all linters
run_all_linters() {
    print_status "Running all linters..."
    
    local total_errors=0
    
    # Run individual linters
    run_go_vet || total_errors=$((total_errors + 1))
    run_go_fmt || total_errors=$((total_errors + 1))
    run_goimports || total_errors=$((total_errors + 1))
    run_golangci_lint || total_errors=$((total_errors + 1))
    run_staticcheck || total_errors=$((total_errors + 1))
    run_ineffassign || total_errors=$((total_errors + 1))
    run_misspell || total_errors=$((total_errors + 1))
    
    if [[ $total_errors -eq 0 ]]; then
        print_success "All linters passed!"
    else
        print_error "Linting failed with $total_errors errors"
        return 1
    fi
}

# Function to show help
show_help() {
    cat << EOF
Teams Notification Service Lint Script

Usage: $0 [OPTIONS]

Options:
    -h, --help          Show this help message
    -a, --all           Run all linters (default)
    --go-vet            Run go vet only
    --go-fmt            Run go fmt only
    --goimports         Run goimports only
    --golangci-lint     Run golangci-lint only
    --staticcheck       Run staticcheck only
    --ineffassign       Run ineffassign only
    --misspell          Run misspell only
    --install           Install required tools

Examples:
    $0                          # Run all linters
    $0 --go-vet                 # Run go vet only
    $0 --golangci-lint          # Run golangci-lint only
    $0 --install                # Install required tools

EOF
}

# Main function
main() {
    local run_all=true
    local run_go_vet=false
    local run_go_fmt=false
    local run_goimports=false
    local run_golangci_lint=false
    local run_staticcheck=false
    local run_ineffassign=false
    local run_misspell=false
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
            --go-vet)
                run_go_vet=true
                run_all=false
                shift
                ;;
            --go-fmt)
                run_go_fmt=true
                run_all=false
                shift
                ;;
            --goimports)
                run_goimports=true
                run_all=false
                shift
                ;;
            --golangci-lint)
                run_golangci_lint=true
                run_all=false
                shift
                ;;
            --staticcheck)
                run_staticcheck=true
                run_all=false
                shift
                ;;
            --ineffassign)
                run_ineffassign=true
                run_all=false
                shift
                ;;
            --misspell)
                run_misspell=true
                run_all=false
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
    
    print_status "Starting linting process for Teams Notification Service"
    
    # Install tools if requested
    if [[ "$install_tools" == true ]]; then
        install_golangci_lint
        print_success "Tools installed"
        exit 0
    fi
    
    # Run specific linters or all
    if [[ "$run_all" == true ]]; then
        run_all_linters
    else
        local total_errors=0
        
        [[ "$run_go_vet" == true ]] && run_go_vet || total_errors=$((total_errors + 1))
        [[ "$run_go_fmt" == true ]] && run_go_fmt || total_errors=$((total_errors + 1))
        [[ "$run_goimports" == true ]] && run_goimports || total_errors=$((total_errors + 1))
        [[ "$run_golangci_lint" == true ]] && run_golangci_lint || total_errors=$((total_errors + 1))
        [[ "$run_staticcheck" == true ]] && run_staticcheck || total_errors=$((total_errors + 1))
        [[ "$run_ineffassign" == true ]] && run_ineffassign || total_errors=$((total_errors + 1))
        [[ "$run_misspell" == true ]] && run_misspell || total_errors=$((total_errors + 1))
        
        if [[ $total_errors -eq 0 ]]; then
            print_success "Linting completed successfully!"
        else
            print_error "Linting failed with $total_errors errors"
            exit 1
        fi
    fi
}

# Run main function with all arguments
main "$@"
