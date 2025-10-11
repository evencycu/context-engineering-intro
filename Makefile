# TeamsNotifyGoV2 Makefile

.PHONY: help build run test clean docker-up docker-down docker-logs db-migrate db-reset openapi-start openapi-stop openapi-status openapi-open test-e2e test-load test-api test-report test-full

# Default target
help:
	@echo "Available commands:"
	@echo "  docker-up      - Start PostgreSQL and Redis with Docker Compose"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  docker-logs    - Show Docker container logs"
	@echo "  db-migrate     - Run database migrations"
	@echo "  db-reset       - Reset database (drop and recreate)"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo ""
	@echo "E2E Testing:"
	@echo "  test-e2e       - Run E2E tests"
	@echo "  test-load      - Run load tests"
	@echo "  test-api       - Run API tests"
	@echo "  test-report    - Generate test reports"
	@echo "  test-full      - Run full E2E test suite"
	@echo ""
	@echo "Development:"
	@echo "  dev-deploy     - Deploy latest code to Docker environment"
	@echo "  deploy-local   - Deploy to local process (fastest)"
	@echo "  deploy-docker  - Deploy to local Docker (isolated)"
	@echo "  deploy-persistent - Deploy to persistent Docker (keeps data)"
	@echo ""
	@echo "Persistent Docker:"
	@echo "  persistent-up    - Start persistent Docker services"
	@echo "  persistent-down  - Stop persistent Docker services"
	@echo "  persistent-logs  - Show persistent Docker logs"
	@echo "  persistent-clean - Clean persistent Docker data"
	@echo ""
	@echo "OpenAPI Documentation:"
	@echo "  openapi-start  - Start OpenAPI server (Swagger UI)"
	@echo "  openapi-stop   - Stop OpenAPI server"
	@echo "  openapi-status - Check OpenAPI server status"
	@echo "  openapi-open   - Open Swagger UI in browser"

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
db-migrate:
	@echo "Running database migrations..."
	@if [ -z "$$(docker ps -q -f name=teamsnotify-postgres)" ]; then \
		echo "PostgreSQL container is not running. Please run 'make docker-up' first."; \
		exit 1; \
	fi
	docker exec -i teamsnotify-postgres psql -U teamsnotify -d teamsnotify < internal/database/schema.sql
	@echo "Database migration completed!"

db-reset:
	@echo "Resetting database..."
	docker-compose down -v
	docker-compose up -d
	@sleep 10
	@echo "Database reset completed!"

# Application operations
build:
	@echo "Building application..."
	go build -o bin/server ./cmd/server

run:
	@echo "Running application..."
	go run ./cmd/server

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Development setup
dev-setup: docker-up db-migrate
	@echo "Development environment is ready!"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "Adminer: http://localhost:8081"
	@echo "API Server: http://localhost:8080"

# Production build
prod-build:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server ./cmd/server

# OpenAPI Documentation
openapi-start:
	@echo "Starting OpenAPI server..."
	./scripts/openapi.sh start

openapi-stop:
	@echo "Stopping OpenAPI server..."
	./scripts/openapi.sh stop

openapi-status:
	@echo "Checking OpenAPI server status..."
	./scripts/openapi.sh status

openapi-open:
	@echo "Opening Swagger UI in browser..."
	./scripts/openapi.sh open

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

# Full E2E test suite
test-full:
	@echo "Running full E2E test suite..."
	./scripts/test/e2e/run_e2e_tests.sh run

# Development deployment
dev-deploy:
	@echo "Running development deployment..."
	./scripts/deploy/local/dev-deploy.sh

# Smart deployment with different modes
deploy-local:
	@echo "Deploying to local process..."
	./scripts/deploy/local/smart-deploy.sh local-process

deploy-docker:
	@echo "Deploying to local Docker..."
	./scripts/deploy/local/smart-deploy.sh local-docker

deploy-persistent:
	@echo "Deploying to persistent Docker..."
	./scripts/deploy/local/smart-deploy.sh persistent-docker

# Persistent Docker management
persistent-up:
	@echo "Starting persistent Docker services..."
	docker-compose -f scripts/docker/docker-compose-persistent.yml up -d

persistent-down:
	@echo "Stopping persistent Docker services..."
	docker-compose -f scripts/docker/docker-compose-persistent.yml down

persistent-logs:
	@echo "Showing persistent Docker logs..."
	docker-compose -f docker-compose-persistent.yml logs -f

persistent-clean:
	@echo "Cleaning persistent Docker data..."
	docker-compose -f scripts/docker/docker-compose-persistent.yml down -v
	docker volume rm teamsnotify_postgres_data teamsnotify_redis_data 2>/dev/null || true