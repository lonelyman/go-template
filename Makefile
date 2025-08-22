# .PHONY declares targets that are not files
.PHONY: setup run dev test test-integration test-coverage lint build build-dev build-prod clean db-migrate docker-dev docker-dev-bg docker-dev-stop docker-dev-logs docker-prod-up docker-prod-down docker-prod-logs docker-clean kill-port help

# ====================================================================================
# VARIABLES
# ====================================================================================
APP_NAME=go-template
DOCKER_IMAGE_PROD=$(APP_NAME):latest

# ====================================================================================
# DEVELOPMENT COMMANDS
# ====================================================================================
# Setup project dependencies
setup:
	go mod tidy
	go mod download

# Run the application locally (without .env)
run:
	go run ./cmd/api/main.go

# Run the application locally with .env loaded
dev:
	@if [ ! -f .env ]; then echo "❌ .env file not found"; exit 1; fi
	@echo "🔧 Loading environment variables from .env"
	@set -a && source .env && set +a && echo "🚀 Starting Go application on port $${PORT}"
	@set -a && source .env && set +a && echo "🌐 Access URL: http://localhost:$${PORT}"
	@set -a && source .env && set +a && go run ./cmd/api/main.go

# ====================================================================================
# DOCKER DEVELOPMENT COMMANDS (แนะนำใช้)
# ====================================================================================
# Docker development (เสถียร, ดู logs realtime)
docker-dev:
	@echo "🐳 Starting Docker development environment"
	@echo "🔧 Using existing PostgreSQL container"
	docker compose -f docker-compose.dev.yml up --build

# Docker development in background
docker-dev-bg:
	@echo "🐳 Starting Docker development environment in background"
	@echo "🔧 Using existing PostgreSQL container"
	docker compose -f docker-compose.dev.yml up --build -d

# Stop Docker development
docker-dev-stop:
	@echo "🛑 Stopping Docker development environment"
	docker compose -f docker-compose.dev.yml down

# View Docker logs
docker-dev-logs:
	@echo "📋 Viewing Docker development logs"
	docker compose -f docker-compose.dev.yml logs -f app-dev

# ====================================================================================
# TESTING COMMANDS
# ====================================================================================
# Run unit tests
test:
	go test -v ./...

# Run integration tests
test-integration:
	go test -v ./tests/...

# Run tests with coverage report
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run linter
lint:
	golangci-lint run

# ====================================================================================
# BUILD COMMANDS
# ====================================================================================
# Build for local development (fast compilation)
build-dev:
	go build -gcflags="all=-N -l" -o ./bin/$(APP_NAME)-dev ./cmd/api/main.go

# Build for production (optimized)
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./bin/$(APP_NAME) ./cmd/api/main.go

# Default build command
build: build-prod

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	@rm -rf ./bin
	@rm -f coverage.out
	@rm -f app.log app.pid

# ====================================================================================
# DATABASE COMMANDS
# ====================================================================================
# Run database migrations
db-migrate:
	go run ./cmd/migrate/main.go

# ====================================================================================
# DOCKER PRODUCTION COMMANDS
# ====================================================================================
# Start production environment
docker-prod-up:
	docker compose --env-file .env.production -f docker-compose.yml up --build -d

# Stop production environment
docker-prod-down:
	docker compose -f docker-compose.yml down

# View production logs
docker-prod-logs:
	docker compose -f docker-compose.yml logs -f app

# ====================================================================================
# UTILITY COMMANDS
# ====================================================================================
# Kill processes using project ports
kill-port:
	@echo "🔍 Checking for processes on project ports..."
	@lsof -ti:9998 | xargs -r kill -9 2>/dev/null || echo "No process on port 9998 (local dev)"
	@lsof -ti:9999 | xargs -r kill -9 2>/dev/null || echo "No process on port 9999 (docker dev host)"
	@lsof -ti:9090 | xargs -r kill -9 2>/dev/null || echo "No process on port 9090 (docker prod host)"
	@lsof -ti:9088 | xargs -r kill -9 2>/dev/null || echo "No process on port 9088 (docker prod container)"
	@echo "✅ Port cleanup completed"

# Clean up unused docker resources
docker-clean:
	@echo "🧹 Pruning Docker system..."
	@docker system prune -af

# ====================================================================================
# HELP
# ====================================================================================
# Show available commands
help:
	@echo "📚 Available commands:"
	@echo ""
	@echo "🚀 Development:"
	@echo "  setup            - Setup project dependencies"
	@echo "  run              - Run locally (without .env)"
	@echo "  dev              - Run locally with .env loaded"
	@echo "  docker-dev       - 🌟 Run in Docker (recommended, stable)"
	@echo "  docker-dev-bg    - Run in Docker background"
	@echo "  docker-dev-stop  - Stop Docker development"
	@echo "  docker-dev-logs  - View Docker logs"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test             - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  lint             - Run linter"
	@echo ""
	@echo "🏗️  Building:"
	@echo "  build-dev        - Build for development"
	@echo "  build-prod       - Build for production"
	@echo "  build            - Default build (production)"
	@echo "  clean            - Clean build artifacts"
	@echo ""
	@echo "🗄️  Database:"
	@echo "  db-migrate       - Run database migrations"
	@echo ""
	@echo "🐳 Production:"
	@echo "  docker-prod-up   - Start production environment"
	@echo "  docker-prod-down - Stop production environment"
	@echo "  docker-prod-logs - View production logs"
	@echo ""
	@echo "🛠️  Utilities:"
	@echo "  kill-port        - Kill processes on project ports"
	@echo "  docker-clean     - Clean Docker resources"
	@echo "  help             - Show this help"
