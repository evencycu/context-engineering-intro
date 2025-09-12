# TeamsNotifyGoV2 Makefile

.PHONY: help build run test clean docker-up docker-down docker-logs db-migrate db-reset

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