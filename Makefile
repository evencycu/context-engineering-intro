# Teams Notification API Server Makefile
# CRITICAL: Based on PRP validation requirements

.PHONY: help build test lint fmt vet clean run dev docker-build docker-run deps migrate

# Variables
BINARY_NAME=teams-notify-server
BUILD_DIR=build
GO_FILES=$(shell find . -name "*.go" | grep -v vendor)
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=${VERSION}"

# Default target
help: ## Show this help message
	@echo "Teams Notification API Server - Development Commands"
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

deps: ## Download dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

docs: ## Generate API documentation
	@echo "📚 Generating API documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/server/main.go -o docs/ --parseDependency --parseInternal; \
		echo "✅ API documentation generated in docs/"; \
	else \
		echo "⚠️  swag not installed, installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g cmd/server/main.go -o docs/ --parseDependency --parseInternal; \
		echo "✅ API documentation generated in docs/"; \
	fi

docs-serve: docs ## Generate and serve API documentation
	@echo "🌐 Serving API documentation at http://localhost:8080/docs/"
	@echo "Press Ctrl+C to stop"
	@python3 -m http.server 8080 --directory docs/ || python -m http.server 8080 --directory docs/

fmt: ## Format Go code
	@echo "🎨 Formatting Go code..."
	go fmt ./...

lint: ## Run linting (if golangci-lint is available)
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed, skipping linting"; \
		echo "   Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

vet: ## Run go vet
	@echo "🔎 Running go vet..."
	go vet ./...

test: ## Run tests
	@echo "🧪 Running tests..."
	go test -race -cover ./...

test-verbose: ## Run tests with verbose output
	@echo "🧪 Running tests (verbose)..."
	go test -v -race -cover ./...

validate: fmt vet lint test docs ## Run all validation steps (PRP Level 1 & 2)
	@echo "✅ All validations passed!"

##@ Build

build: ## Build the server binary
	@echo "🔨 Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/server

build-linux: ## Build for Linux
	@echo "🔨 Building ${BINARY_NAME} for Linux..."
	@mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ./cmd/server

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -rf ${BUILD_DIR}
	go clean

##@ Run

run: ## Run the server locally
	@echo "🚀 Running server locally..."
	go run ./cmd/server

dev: ## Run in development mode with live reload (requires air)
	@echo "🔥 Starting development server..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  'air' not installed for live reload"; \
		echo "   Install with: go install github.com/cosmtrek/air@latest"; \
		echo "   Running normally..."; \
		go run ./cmd/server; \
	fi

##@ Database


##@ Docker

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -f deployments/docker/Dockerfile -t teams-notify-server:${VERSION} .

docker-run: ## Run in Docker container
	@echo "🐳 Running Docker container..."
	docker run --rm -p 8080:8080 teams-notify-server:${VERSION}

##@ Utilities

env-check: ## Check required environment variables
	@echo "🔧 Checking environment variables..."
	@echo "Required variables:"
	@echo "  TEAMS_NOTIFY_TEAMS_APP_ID:     ${TEAMS_NOTIFY_TEAMS_APP_ID}"
	@echo "  TEAMS_NOTIFY_TEAMS_APP_SECRET: $(if ${TEAMS_NOTIFY_TEAMS_APP_SECRET},***SET***,NOT SET)"
	@echo "  TEAMS_NOTIFY_TEAMS_BOT_ID:     ${TEAMS_NOTIFY_TEAMS_BOT_ID}"
	@echo "  TEAMS_NOTIFY_AUTH_JWT_SECRET:  $(if ${TEAMS_NOTIFY_AUTH_JWT_SECRET},***SET***,NOT SET)"
	@echo "  TEAMS_NOTIFY_DATABASE_PASSWORD:$(if ${TEAMS_NOTIFY_DATABASE_PASSWORD},***SET***,NOT SET)"

config-example: ## Generate example configuration
	@echo "📝 Example configuration (config.yaml):"
	@echo "server:"
	@echo "  port: \"8080\""
	@echo "  environment: \"development\""
	@echo ""
	@echo "teams:"
	@echo "  app_id: \"your-bot-app-id\""
	@echo "  app_secret: \"your-bot-app-secret\""
	@echo "  bot_id: \"28:your-bot-app-id\""
	@echo ""
	@echo "database:"
	@echo "  host: \"localhost\""
	@echo "  port: \"5432\""
	@echo "  username: \"teams_notify\""
	@echo "  password: \"your_password\""
	@echo "  database_name: \"teams_notify\""
	@echo ""
	@echo "auth:"
	@echo "  jwt_secret: \"your-jwt-secret-here\""

logs: ## Show recent logs (if running with systemd)
	@echo "📋 Recent logs..."
	@if systemctl is-active --quiet teams-notify-server; then \
		journalctl -u teams-notify-server -n 50 -f; \
	else \
		echo "Service not running or not managed by systemd"; \
	fi

##@ CI/CD

ci: deps validate build ## Run CI pipeline locally
	@echo "🔄 CI pipeline completed successfully!"

security: ## Run security checks (if gosec is available)
	@echo "🔒 Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not installed, skipping security scan"; \
		echo "   Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

##@ Information

version: ## Show version information
	@echo "Version: ${VERSION}"
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Build time: $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')"

status: ## Show service status
	@echo "🏥 Service Status Check"
	@echo "======================"
	@if command -v curl >/dev/null 2>&1; then \
		echo "Health Check:"; \
		curl -s http://localhost:8080/health | jq . 2>/dev/null || curl -s http://localhost:8080/health; \
	else \
		echo "curl not available for health check"; \
	fi