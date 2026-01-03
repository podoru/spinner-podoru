.PHONY: build run dev test clean lint migrate-up migrate-down docker-up docker-down docs help \
	dev-setup dev-check dev-quick dev-reset \
	prod-check prod-install prod-setup prod-up prod-down prod-logs prod-restart prod-backup prod-update prod-secrets

BINARY_NAME=podoru
BUILD_DIR=bin
MAIN_PATH=./cmd/podoru

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)

# Build info
VERSION?=0.1.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Linker flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(GOBIN)/$(BINARY_NAME)"

## run: Run the application
run: build
	@$(GOBIN)/$(BINARY_NAME)

## dev: Run with hot reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/air-verse/air@latest"; \
		exit 1; \
	fi

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run linter
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Done"

## tidy: Tidy go modules
tidy:
	@echo "Tidying modules..."
	@go mod tidy
	@echo "Done"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@echo "Done"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Done"

## migrate-up: Run database migrations
migrate-up:
	@echo "Running migrations..."
	@go run $(MAIN_PATH) migrate up

## migrate-down: Rollback database migrations
migrate-down:
	@echo "Rolling back migrations..."
	@go run $(MAIN_PATH) migrate down

## migrate-create: Create a new migration (usage: make migrate-create name=migration_name)
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@migrate create -ext sql -dir migrations -seq $(name)

## docker-up: Start Docker Compose services
docker-up:
	@echo "Starting Docker services..."
	@docker compose up -d
	@echo "Services started"

## docker-down: Stop Docker Compose services
docker-down:
	@echo "Stopping Docker services..."
	@docker compose down
	@echo "Services stopped"

## docker-logs: Show Docker Compose logs
docker-logs:
	@docker compose logs -f

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .
	@echo "Image built: $(BINARY_NAME):$(VERSION)"

## docker-push: Push Docker image (requires REGISTRY variable)
docker-push:
	@if [ -z "$(REGISTRY)" ]; then \
		echo "Usage: make docker-push REGISTRY=your-registry"; \
		exit 1; \
	fi
	@docker tag $(BINARY_NAME):$(VERSION) $(REGISTRY)/$(BINARY_NAME):$(VERSION)
	@docker push $(REGISTRY)/$(BINARY_NAME):$(VERSION)

## setup: Initial project setup (legacy, use dev-setup instead)
setup: deps
	@echo "Setting up project..."
	@cp -n .env.example .env 2>/dev/null || true
	@echo "Setup complete. Edit .env file with your configuration."

# =============================================================================
# Development Setup Commands
# =============================================================================

## dev-setup: Full development environment setup (recommended for first-time setup)
dev-setup:
	@chmod +x scripts/dev-setup.sh
	@./scripts/dev-setup.sh setup

## dev-check: Check development prerequisites
dev-check:
	@chmod +x scripts/dev-setup.sh
	@./scripts/dev-setup.sh check

## dev-quick: Quick start development (database + hot reload)
dev-quick:
	@chmod +x scripts/dev-setup.sh
	@./scripts/dev-setup.sh quick

## dev-reset: Reset development environment (removes data)
dev-reset:
	@chmod +x scripts/dev-setup.sh
	@./scripts/dev-setup.sh reset

## docs: Generate API documentation (Swagger/OpenAPI)
docs:
	@echo "Generating API documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/podoru/main.go -o docs --parseDependency --parseInternal --parseDepth 3; \
	else \
		echo "swag not installed. Install with: go install github.com/swaggo/swag/cmd/swag@v1.16.3"; \
		exit 1; \
	fi
	@echo "Documentation generated in docs/"

## all: Build and test
all: lint test build

# =============================================================================
# Production Commands
# =============================================================================

## prod-check: Run pre-flight checks for production deployment
prod-check:
	@chmod +x scripts/install.sh
	@./scripts/install.sh check

## prod-install: Run full production installation
prod-install:
	@chmod +x scripts/install.sh
	@./scripts/install.sh install

## prod-secrets: Generate secure secrets for production
prod-secrets:
	@chmod +x scripts/install.sh
	@./scripts/install.sh secrets

## prod-env: Create .env.prod template file (non-interactive)
prod-env:
	@chmod +x scripts/install.sh
	@./scripts/install.sh env

## prod-setup: Run interactive setup wizard
prod-setup:
	@chmod +x scripts/install.sh
	@./scripts/install.sh setup

## prod-up: Start production services
prod-up:
	@echo "Starting production services..."
	@docker compose -f docker-compose.prod.yml --env-file .env.prod up -d
	@echo "Services started. Run 'make prod-logs' to view logs."

## prod-down: Stop production services
prod-down:
	@echo "Stopping production services..."
	@docker compose -f docker-compose.prod.yml --env-file .env.prod down
	@echo "Services stopped."

## prod-restart: Restart production services
prod-restart:
	@echo "Restarting production services..."
	@docker compose -f docker-compose.prod.yml --env-file .env.prod restart
	@echo "Services restarted."

## prod-logs: View production logs (follow mode)
prod-logs:
	@docker compose -f docker-compose.prod.yml logs -f

## prod-status: Show production services status
prod-status:
	@echo "Production Services Status:"
	@echo "============================"
	@docker compose -f docker-compose.prod.yml ps
	@echo ""
	@echo "Resource Usage:"
	@docker stats --no-stream $$(docker compose -f docker-compose.prod.yml ps -q) 2>/dev/null || true

## prod-backup: Backup PostgreSQL database
prod-backup:
	@echo "Creating database backup..."
	@mkdir -p backups
	@docker compose -f docker-compose.prod.yml exec -T postgres pg_dump -U $${DB_USER:-podoru} $${DB_NAME:-podoru} | gzip > backups/podoru_$$(date +%Y%m%d_%H%M%S).sql.gz
	@echo "Backup created: backups/podoru_$$(date +%Y%m%d_%H%M%S).sql.gz"
	@ls -lh backups/*.sql.gz | tail -5

## prod-restore: Restore PostgreSQL database (usage: make prod-restore file=backups/backup.sql.gz)
prod-restore:
	@if [ -z "$(file)" ]; then \
		echo "Usage: make prod-restore file=backups/backup.sql.gz"; \
		echo "Available backups:"; \
		ls -lh backups/*.sql.gz 2>/dev/null || echo "  No backups found"; \
		exit 1; \
	fi
	@echo "Restoring database from $(file)..."
	@read -p "This will overwrite the current database. Continue? [y/N] " confirm && [ "$$confirm" = "y" ] || exit 1
	@gunzip -c $(file) | docker compose -f docker-compose.prod.yml exec -T postgres psql -U $${DB_USER:-podoru} $${DB_NAME:-podoru}
	@echo "Database restored."

## prod-update: Update production to latest version
prod-update:
	@echo "Updating Podoru to latest version..."
	@echo ""
	@echo "Step 1: Creating backup..."
	@$(MAKE) prod-backup
	@echo ""
	@echo "Step 2: Pulling latest code..."
	@git pull origin master
	@echo ""
	@echo "Step 3: Rebuilding application..."
	@docker compose -f docker-compose.prod.yml --env-file .env.prod build --no-cache podoru
	@echo ""
	@echo "Step 4: Restarting services..."
	@docker compose -f docker-compose.prod.yml --env-file .env.prod up -d podoru
	@echo ""
	@echo "Update complete. Run 'make prod-logs' to check for errors."

## prod-shell: Open shell in production app container
prod-shell:
	@docker compose -f docker-compose.prod.yml exec podoru sh

## prod-db-shell: Open PostgreSQL shell
prod-db-shell:
	@docker compose -f docker-compose.prod.yml exec postgres psql -U $${DB_USER:-podoru} $${DB_NAME:-podoru}

## prod-clean: Remove all production data (DANGEROUS)
prod-clean:
	@echo "WARNING: This will remove ALL production data including the database!"
	@read -p "Type 'DELETE' to confirm: " confirm && [ "$$confirm" = "DELETE" ] || exit 1
	@docker compose -f docker-compose.prod.yml --env-file .env.prod down -v
	@echo "All production data removed."
