# .PHONY declares targets that are not files. This prevents conflicts with files of the same name.
.PHONY: setup test test-integration test-coverage lint clean db-migrate db-migrate-primary db-migrate-logs docker-dev-up docker-dev-d docker-dev-down docker-dev-logs docker-prod-up docker-prod-down docker-prod-logs docker-clean kill-port help

# ====================================================================================
# VARIABLES
# ====================================================================================
APP_NAME=go-template

# ====================================================================================
# LOCAL UTILITY COMMANDS (For Developer Experience)
# ====================================================================================
# Setup project dependencies (For IDE and local tools)
setup:
	@echo "🔧 Setting up project dependencies for IDE..."
	@go mod tidy
	@go mod download

# Run linter on host machine
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up build artifacts..."
	@rm -rf ./bin
	@rm -f coverage.out

# Kill processes using project ports
kill-port:
	@echo "🔪 Killing processes on project ports..."
	@lsof -ti:9998 | xargs -r kill -9 2>/dev/null || echo "No process on port 9998"
	@lsof -ti:9999 | xargs -r kill -9 2>/dev/null || echo "No process on port 9999"
	@echo "✅ Port cleanup completed"

# ====================================================================================
# DOCKER DEVELOPMENT COMMANDS (Workflow หลัก)
# ====================================================================================
# Start Docker development environment (with logs)
docker-dev-up:
	@echo "🐳 Starting Docker application service [app-dev]..."
	@docker compose -f docker-compose.dev.yml up --build app-dev

# Start Docker development in background (detached)
docker-dev-d:
	@echo "🐳 Starting Docker application service [app-dev] in background..."
	@docker compose -f docker-compose.dev.yml up --build -d app-dev

# Stop Docker development environment
docker-dev-down:
	@echo "🛑 Stopping Docker development environment"
	@docker compose -f docker-compose.dev.yml down

# View Docker development logs
docker-dev-logs:
	@echo "📋 Viewing Docker development logs..."
	@docker compose -f docker-compose.dev.yml logs -f app-dev

# ====================================================================================
# DOCKER TESTING & DATABASE COMMANDS
# ====================================================================================
# Run unit tests inside a Docker container
test:
	@echo "🧪 Running unit tests inside a Docker container..."
	@docker compose -f docker-compose.dev.yml run --rm migrate go test -v ./...

# Run integration tests inside a Docker container
test-integration:
	@echo "🧪 Running integration tests inside a Docker container..."
	@docker compose -f docker-compose.dev.yml run --rm migrate go test -v ./tests/...

# ====================================================================================
# DATABASE COMMANDS
# ====================================================================================
# Run database migrations for a specific database using a one-off Docker container.
# Usage: make db-migrate db=primary
db-migrate:
ifndef db
	$(error db is not set. Usage: make db-migrate db=<primary|logs|analytics>)
endif
	@echo "🗄️  Migrating database: [$(db)] inside a Docker container..."
	# ⭐️ เพิ่ม `go run ./cmd/migrate/main.go` เข้าไปตรงนี้! ⭐️
	@docker compose -f docker-compose.dev.yml run --rm migrate go run ./cmd/migrate/main.go --db=$(db) --path=db/migrations/$(db)

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
	@echo "🚀 Development (Docker Workflow):"
	@echo "  docker-dev-up      - 🌟 Start Docker development (recommended)"
	@echo "  docker-dev-d       - Start Docker development in background"
	@echo "  docker-dev-down    - Stop Docker development"
	@echo "  docker-dev-logs    - View Docker development logs"
	@echo ""
	@echo "🧪 Testing (Docker Workflow):"
	@echo "  test               - Run unit tests inside Docker"
	@echo "  test-integration   - Run integration tests inside Docker"
	@echo ""
	@echo "🗄️  Database:"
	@echo "  db-migrate-primary - Migrate the PRIMARY database inside Docker"
	@echo "  db-migrate-logs    - Migrate the LOGS database inside Docker"
	@echo ""
	@echo "🛠️  Local Utilities:"
	@echo "  setup              - Setup Go modules for your IDE"
	@echo "  lint               - Run linter on host machine"
	@echo "  clean              - Clean build artifacts"
	@echo "  kill-port          - Kill processes on project ports"
	@echo "  docker-clean       - Clean Docker resources"
	@echo "  help               - Show this help"