# .PHONY declares targets that are not files. This prevents conflicts with files of the same name.
.PHONY: setup run dev test test-integration test-coverage lint build build-dev build-prod clean db-migrate db-migrate-primary db-migrate-logs docker-dev-up docker-dev-d docker-dev-down docker-dev-logs docker-prod-up docker-prod-down docker-prod-logs docker-clean kill-port help

# ====================================================================================
# VARIABLES
# ====================================================================================
APP_NAME=go-template
DOCKER_IMAGE_PROD=$(APP_NAME):latest
DOCKER_IMAGE_DEV=$(APP_NAME):dev

# ====================================================================================
# DEVELOPMENT COMMANDS
# ====================================================================================
# Setup project dependencies
setup:
	@echo "🔧 Setting up project dependencies..."
	@go mod tidy
	@go mod download

# Run the application locally (without .env)
run:
	@echo "🚀 Starting Go application (without .env)..."
	@go run ./cmd/api/main.go

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
# Start Docker development environment (with logs)
docker-dev-up:
	@echo "🐳 Starting Docker development environment..."
	@docker compose -f docker-compose.dev.yml up --build

# Start Docker development in background (detached)
docker-dev-d:
	@echo "🐳 Starting Docker development environment in background..."
	@docker compose -f docker-compose.dev.yml up --build -d

# Stop Docker development environment
docker-dev-down:
	@echo "🛑 Stopping Docker development environment"
	@docker compose -f docker-compose.dev.yml down

# View Docker development logs
docker-dev-logs:
	@echo "📋 Viewing Docker development logs..."
	@docker compose -f docker-compose.dev.yml logs -f app-dev

# ====================================================================================
# TESTING COMMANDS
# ====================================================================================
# Run unit tests
test:
	@echo "🧪 Running unit tests..."
	@go test -v ./...

# Run integration tests
test-integration:
	@echo "🧪 Running integration tests..."
	@go test -v ./tests/...

# Run tests with coverage report
test-coverage:
	@echo "📊 Generating test coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Run linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run

# ====================================================================================
# BUILD COMMANDS
# ====================================================================================
# Build for local development (fast compilation)
build-dev:
	@echo "🏗️  Building for development (fast compile)..."
	@go build -gcflags="all=-N -l" -o ./bin/$(APP_NAME)-dev ./cmd/api/main.go

# Build for production (optimized for size and performance)
build-prod:
	@echo "🏗️  Building for production (optimized)..."
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./bin/$(APP_NAME) ./cmd/api/main.go

# Default build command
build: build-prod

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up build artifacts..."
	@rm -rf ./bin
	@rm -f coverage.out

# ====================================================================================
# DATABASE COMMANDS (✨ ปรับปรุงใหม่!)
# ====================================================================================
# Run database migrations for a specific database.
# Usage: make db-migrate db=<primary|logs|analytics>
db-migrate:
ifndef db
	$(error db is not set. Usage: make db-migrate db=<primary|logs|analytics>)
endif
	@echo "🗄️  Migrating database: [$(db)]..."
	@go run ./cmd/migrate/main.go --db=$(db) --path=db/migrations/$(db)

# --- Shortcuts for convenience ---
db-migrate-primary:
	@make db-migrate db=primary

db-migrate-logs:
	@make db-migrate db=logs


# ====================================================================================
# DOCKER PRODUCTION COMMANDS
# ====================================================================================
# Start production environment
docker-prod-up:
	@echo "🚀 Starting production environment..."
	@docker compose --env-file .env.production -f docker-compose.prod.yml up --build -d

# Stop production environment
docker-prod-down:
	@echo "🛑 Stopping production environment"
	@docker compose -f docker-compose.prod.yml down

# View production logs
docker-prod-logs:
	@echo "📋 Viewing production logs..."
	@docker compose -f docker-compose.prod.yml logs -f app

# ====================================================================================
# UTILITY COMMANDS
# ====================================================================================
# Kill processes using project ports
kill-port:
	@echo "🔪 Killing processes on project ports..."
	@lsof -ti:8080 | xargs -r kill -9 2>/dev/null || echo "No process on port 8080"
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
	@echo "  setup              - Setup project dependencies"
	@echo "  run                - Run locally (without .env)"
	@echo "  dev                - Run locally with .env loaded"
	@echo "  docker-dev-up      - 🌟 Start Docker development (recommended)"
	@echo "  docker-dev-d       - Start Docker development in background"
	@echo "  docker-dev-down    - Stop Docker development"
	@echo "  docker-dev-logs    - View Docker development logs"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test               - Run unit tests"
	@echo "  test-integration   - Run integration tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  lint               - Run linter"
	@echo ""
	@echo "🏗️  Building:"
	@echo "  build-dev          - Build for development"
	@echo "  build-prod         - Build for production"
	@echo "  build              - Default build (production)"
	@echo "  clean              - Clean build artifacts"
	@echo ""
	@echo "🗄️  Database:"
	@echo "  db-migrate-primary - Migrate the PRIMARY database"
	@echo "  db-migrate-logs    - Migrate the LOGS database"
	@echo ""
	@echo "🐳 Production:"
	@echo "  docker-prod-up     - Start production environment"
	@echo "  docker-prod-down   - Stop production environment"
	@echo "  docker-prod-logs   - View production logs"
	@echo ""
	@echo "🛠️  Utilities:"
	@echo "  kill-port          - Kill processes on project ports"
	@echo "  docker-clean       - Clean Docker resources"
	@echo "  help               - Show this help"