# TeamsNotifyGoV2 Makefile


# Default target

# Docker operations
docker-up:
	@echo "Starting PostgreSQL and Redis..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Services are ready!"

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	docker-compose logs -f

# Database operations
	@if [ -z "$$(docker ps -q -f name=teamsnotify-postgres)" ]; then \
		echo "PostgreSQL container is not running. Please run 'make docker-up' first."; \
		exit 1; \
	fi
	docker exec -i teamsnotify-postgres psql -U teamsnotify -d teamsnotify < internal/database/schema.sql

db-reset:
	@echo "Resetting database..."
	docker-compose down -v
	docker-compose up -d
	@sleep 10
	@echo "Database reset completed!"

# Application operations
build:
	@echo "Building application..."
	go build -o server ./cmd/server

run:
	@echo "Running application..."
	go run ./cmd/server

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -f server
	go clean

# Development setup
	@echo "Development environment is ready!"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "Adminer: http://localhost:8081"
	@echo "API Server: http://localhost:8080"

# Production build
prod-build:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Deployment commands (integrated from deploy.sh)
deploy:
	@echo "Running full deployment..."
	@./scripts/deploy/docker/deploy.sh deploy

deploy-dev:
	@echo "Running development deployment..."
	@./scripts/deploy/docker/deploy.sh dev

deploy-local:
	@echo "Running local deployment..."
	@./scripts/deploy/docker/deploy.sh local

deploy-quick:
	@echo "Running quick deployment..."
	@./scripts/deploy/docker/deploy.sh quick

deploy-docker:
	@echo "Running Docker deployment..."
	@./scripts/deploy/docker/deploy.sh docker-full

deploy-server:
	@echo "Deploying server only..."
	@./scripts/deploy/docker/deploy.sh server

deploy-openapi:
	@echo "Deploying OpenAPI only..."
	@./scripts/deploy/docker/deploy.sh openapi

deploy-stop:
	@echo "Stopping all services..."
	@./scripts/deploy/docker/deploy.sh stop

deploy-restart:
	@echo "Restarting all services..."
	@./scripts/deploy/docker/deploy.sh restart

deploy-status:
	@echo "Checking service status..."
	@./scripts/deploy/docker/deploy.sh status

deploy-logs:
	@echo "Showing server logs..."
	@./scripts/deploy/docker/deploy.sh logs

# Docker commands
docker-build:
	@echo "Building server Docker image..."
	@./scripts/deploy/docker/deploy.sh build-docker

docker-stop:
	@echo "Stopping server Docker container..."
	@./scripts/deploy/docker/deploy.sh stop-docker

docker-start:
	@echo "Starting server Docker container..."
	@./scripts/deploy/docker/deploy.sh start-docker

docker-restart:
	@echo "Restarting server Docker container..."
	@./scripts/deploy/docker/deploy.sh stop-docker
	@./scripts/deploy/docker/deploy.sh start-docker

# OpenAPI Documentation
openapi-start:
	@echo "Starting OpenAPI server..."
	@./scripts/deploy/docker/deploy.sh openapi

openapi-stop:
	@echo "Stopping OpenAPI server..."
	@./scripts/deploy/docker/deploy.sh stop

openapi-status:
	@echo "Checking OpenAPI server status..."
	@./scripts/deploy/docker/deploy.sh status

openapi-open:
	@echo "Opening Swagger UI in browser..."
	@open http://localhost:8082/ || echo "Please open http://localhost:8082/ in your browser"

# E2E Testing
test-e2e:
	@echo "Running E2E tests..."
	./scripts/test/e2e/run_e2e_tests.sh e2e

test-load:
	@echo "Running load tests..."
	./scripts/test/e2e/run_e2e_tests.sh load

test-api:
	@echo "Running API tests..."
	./scripts/test/e2e/run_e2e_tests.sh api

test-report:
	@echo "Generating test reports..."
	./scripts/utils/generate_test_report.sh generate

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Build & Run:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  dev-setup      - Setup development environment"
	@echo "  prod-build     - Build for production"
	@echo ""
	@echo "Deployment (integrated from deploy.sh):"
	@echo "  deploy         - Full deployment"
	@echo "  deploy-dev     - Development deployment"
	@echo "  deploy-local   - Local deployment"
	@echo "  deploy-quick   - Quick deployment"
	@echo "  deploy-docker  - Docker deployment"
	@echo "  deploy-server  - Deploy server only"
	@echo "  deploy-openapi - Deploy OpenAPI only"
	@echo "  deploy-stop    - Stop all services"
	@echo "  deploy-restart - Restart all services"
	@echo "  deploy-status  - Check service status"
	@echo "  deploy-logs    - Show server logs"
	@echo ""
	@echo "OpenAPI:"
	@echo "  openapi-start  - Start OpenAPI server"
	@echo "  openapi-stop   - Stop OpenAPI server"
	@echo "  openapi-status - Check OpenAPI server status"
	@echo "  openapi-open   - Open Swagger UI in browser"
	@echo ""
	@echo "Testing:"
	@echo "  test-e2e       - Run E2E tests"
	@echo "  test-load      - Run load tests"
	@echo "  test-api       - Run API tests"
	@echo "  test-report    - Generate test reports"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   - Build server Docker image"
	@echo "  docker-start   - Start server Docker container"
	@echo "  docker-stop    - Stop server Docker container"
	@echo "  docker-restart - Restart server Docker container"
	@echo ""
	@echo "SDK:"
	@echo "  sdk-generate   - Generate Go and Java SDKs"
	@echo "  sdk-clean      - Clean generated SDK files"
	@echo "  sdk-help       - Show SDK management help"
	@echo ""

# Full E2E test suite
test-full:
	@echo "Running full E2E test suite..."
	./scripts/test/e2e/run_e2e_tests.sh run

# Legacy deployment commands removed - use deploy-local instead

# Persistent Docker management
persistent-up:
	@echo "Starting persistent Docker services..."
	@docker-compose -f scripts/docker/docker-compose-persistent.yml up -d

persistent-down:
	@echo "Stopping persistent Docker services..."
	@docker-compose -f scripts/docker/docker-compose-persistent.yml down

persistent-logs:
	@echo "Showing persistent Docker logs..."
	@docker-compose -f docker-compose-persistent.yml logs -f

persistent-clean:
	@echo "Cleaning persistent Docker data..."
	@docker-compose -f scripts/docker/docker-compose-persistent.yml down -v
	@docker volume rm teamsnotify_postgres_data teamsnotify_redis_data 2>/dev/null || true

# SDK generation commands
sdk-generate:
	@echo "Generating SDKs from OpenAPI specification..."
	@./sdk/generate.sh

sdk-clean:
	@echo "Cleaning generated SDK files..."
	@rm -rf sdk/go sdk/java

sdk-help:
	@echo "SDK Management:"
	@echo "  sdk-generate - Generate Go and Java SDKs"
	@echo "  sdk-clean   - Clean generated SDK files"