#!/bin/bash

# Development Workflow Script
# Automates the incremental development workflow defined in .cursorrules

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_step() {
    echo -e "${BLUE}▶${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to run tests
run_tests() {
    print_step "Running tests..."
    if make test; then
        print_success "All tests passed"
        return 0
    else
        print_error "Tests failed"
        return 1
    fi
}

# Function to check code quality
check_code_quality() {
    print_step "Checking code quality..."
    
    # Format code
    print_step "Formatting code..."
    if go fmt ./...; then
        print_success "Code formatted"
    else
        print_error "Code formatting failed"
        return 1
    fi
    
    # Run go vet
    print_step "Running go vet..."
    if go vet ./...; then
        print_success "go vet passed"
    else
        print_warning "go vet found issues (non-fatal)"
    fi
    
    # Run linter if available
    if command_exists golangci-lint; then
        print_step "Running golangci-lint..."
        if golangci-lint run; then
            print_success "Linter passed"
        else
            print_warning "Linter found issues (non-fatal)"
        fi
    else
        print_warning "golangci-lint not installed, skipping"
    fi
    
    return 0
}

# Function to check test coverage
check_coverage() {
    print_step "Checking test coverage..."
    coverage=$(go test -cover ./... 2>&1 | grep -oP 'coverage: \K[0-9.]+' | head -1)
    if [ -n "$coverage" ]; then
        print_success "Test coverage: ${coverage}%"
        # Check if coverage meets threshold (80%)
        if (( $(echo "$coverage >= 80" | bc -l) )); then
            print_success "Coverage meets threshold (≥80%)"
        else
            print_warning "Coverage below threshold (current: ${coverage}%, target: ≥80%)"
        fi
    else
        print_warning "Could not determine coverage"
    fi
}

# Function to check documentation sync
check_docs() {
    print_step "Checking documentation sync..."
    if [ -f "./scripts/utils/check-doc-sync.sh" ]; then
        if ./scripts/utils/check-doc-sync.sh; then
            print_success "Documentation is in sync"
        else
            print_warning "Documentation sync check found issues"
        fi
    else
        print_warning "Documentation sync script not found, skipping"
    fi
}

# Function to run full test suite
run_full_tests() {
    print_step "Running full test suite..."
    if make test-full; then
        print_success "Full test suite passed"
        return 0
    else
        print_error "Full test suite failed"
        return 1
    fi
}

# Main workflow functions

# Pre-commit workflow
pre_commit() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  Pre-Commit Workflow${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    check_code_quality || return 1
    run_tests || return 1
    check_coverage
    check_docs
    
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    print_success "Pre-commit checks completed"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Feature development workflow
feature_dev() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  Feature Development Workflow${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    print_step "1. Planning Phase"
    print_warning "Please ensure you have:"
    echo "   - Understood requirements and scope"
    echo "   - Checked existing similar implementations"
    echo "   - Planned API changes (if any)"
    echo "   - Identified affected files"
    echo ""
    read -p "Press Enter to continue to development phase..."
    
    print_step "2. Development Phase"
    print_warning "Please ensure you have:"
    echo "   - Written code following project style"
    echo "   - Added comments for complex logic"
    echo "   - Ensured comprehensive error handling"
    echo ""
    read -p "Press Enter to continue to testing phase..."
    
    print_step "3. Testing Phase"
    run_tests || return 1
    check_coverage
    
    print_step "4. Documentation Phase"
    check_docs
    print_warning "Please manually verify documentation is updated"
    
    print_step "5. Code Quality Phase"
    check_code_quality || return 1
    
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    print_success "Feature development workflow completed"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Refactoring workflow
refactor() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  Refactoring Workflow${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    print_step "1. Ensuring existing tests pass..."
    run_tests || {
        print_error "Existing tests must pass before refactoring"
        return 1
    }
    
    print_step "2. Running full test suite..."
    run_full_tests || return 1
    
    print_step "3. Checking code quality..."
    check_code_quality || return 1
    
    print_step "4. Checking documentation..."
    check_docs
    
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    print_success "Refactoring workflow completed"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Full workflow (all checks)
full() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  Full Development Workflow${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    check_code_quality || return 1
    run_full_tests || return 1
    check_coverage
    check_docs
    
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    print_success "Full workflow completed"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Help message
show_help() {
    echo "Development Workflow Script"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  pre-commit    Run pre-commit checks (format, test, lint)"
    echo "  feature       Run feature development workflow"
    echo "  refactor      Run refactoring workflow"
    echo "  full          Run full workflow (all checks)"
    echo "  test          Run tests only"
    echo "  quality       Check code quality only"
    echo "  docs          Check documentation sync only"
    echo "  help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 pre-commit    # Before committing"
    echo "  $0 feature       # When developing new feature"
    echo "  $0 refactor      # When refactoring code"
    echo "  $0 full         # Complete workflow"
}

# Main script logic
case "${1:-help}" in
    pre-commit)
        pre_commit
        ;;
    feature)
        feature_dev
        ;;
    refactor)
        refactor
        ;;
    full)
        full
        ;;
    test)
        run_tests
        ;;
    quality)
        check_code_quality
        ;;
    docs)
        check_docs
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
