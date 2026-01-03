.PHONY: build run dev test clean lint migrate-up migrate-down docker-up docker-down docs help

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

## setup: Initial project setup
setup: deps
	@echo "Setting up project..."
	@cp -n .env.example .env 2>/dev/null || true
	@echo "Setup complete. Edit .env file with your configuration."

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
