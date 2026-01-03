#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Configuration
REQUIRED_GO_VERSION="1.22"
REQUIRED_DOCKER_VERSION="20.10.0"
DEV_PORTS=(5432 8080 8081)

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

#######################################
# Print functions
#######################################
print_header() {
    echo -e "\n${BLUE}============================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}============================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_step() {
    echo -e "\n${CYAN}${BOLD}>> $1${NC}\n"
}

#######################################
# Version comparison
#######################################
version_gte() {
    [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" = "$2" ]
}

#######################################
# Check functions
#######################################
check_go() {
    print_info "Checking Go installation..."

    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        echo "  Install Go: https://go.dev/doc/install"
        return 1
    fi

    GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+' | head -1)

    if version_gte "$GO_VERSION" "$REQUIRED_GO_VERSION"; then
        print_success "Go version: $GO_VERSION (required: >=$REQUIRED_GO_VERSION)"
    else
        print_error "Go version $GO_VERSION is too old (required: >=$REQUIRED_GO_VERSION)"
        return 1
    fi

    return 0
}

check_docker() {
    print_info "Checking Docker installation..."

    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        echo "  Install Docker: https://docs.docker.com/engine/install/"
        return 1
    fi

    DOCKER_VERSION=$(docker version --format '{{.Server.Version}}' 2>/dev/null || echo "0.0.0")

    if version_gte "$DOCKER_VERSION" "$REQUIRED_DOCKER_VERSION"; then
        print_success "Docker version: $DOCKER_VERSION"
    else
        print_error "Docker version $DOCKER_VERSION is too old (required: >=$REQUIRED_DOCKER_VERSION)"
        return 1
    fi

    if ! docker info &> /dev/null; then
        print_error "Docker daemon is not running"
        echo "  Start Docker: sudo systemctl start docker"
        return 1
    fi
    print_success "Docker daemon is running"

    return 0
}

check_docker_compose() {
    print_info "Checking Docker Compose..."

    if docker compose version &> /dev/null; then
        COMPOSE_VERSION=$(docker compose version --short 2>/dev/null)
        COMPOSE_CMD="docker compose"
        print_success "Docker Compose version: $COMPOSE_VERSION"
        return 0
    elif command -v docker-compose &> /dev/null; then
        COMPOSE_VERSION=$(docker-compose version --short 2>/dev/null)
        COMPOSE_CMD="docker-compose"
        print_success "Docker Compose version: $COMPOSE_VERSION (legacy)"
        return 0
    else
        print_error "Docker Compose is not installed"
        return 1
    fi
}

# Detect compose command
detect_compose_cmd() {
    if docker compose version &> /dev/null; then
        echo "docker compose"
    else
        echo "docker-compose"
    fi
}

check_ports() {
    print_info "Checking port availability..."
    local ports_in_use=()

    for port in "${DEV_PORTS[@]}"; do
        if ss -tuln 2>/dev/null | grep -q ":${port} " || \
           netstat -tuln 2>/dev/null | grep -q ":${port} "; then
            ports_in_use+=($port)
            print_warning "Port $port is already in use"
        else
            print_success "Port $port is available"
        fi
    done

    if [[ ${#ports_in_use[@]} -gt 0 ]]; then
        echo ""
        print_warning "Ports in use: ${ports_in_use[*]}"
        echo "  You may need to stop existing services or adjust docker-compose.yml"
        return 1
    fi

    return 0
}

check_make() {
    print_info "Checking Make..."
    if command -v make &> /dev/null; then
        print_success "Make is installed"
        return 0
    else
        print_warning "Make is not installed (optional but recommended)"
        return 0
    fi
}

#######################################
# Setup functions
#######################################
create_env_file() {
    print_info "Setting up environment file..."

    local env_file="$PROJECT_DIR/.env"
    local env_example="$PROJECT_DIR/.env.example"

    if [[ -f "$env_file" ]]; then
        print_warning ".env file already exists"
        read -p "Overwrite with default development settings? [y/N]: " overwrite
        if [[ ! "$overwrite" =~ ^[Yy]$ ]]; then
            print_info "Keeping existing .env file"
            return 0
        fi
    fi

    if [[ -f "$env_example" ]]; then
        cp "$env_example" "$env_file"
    else
        cat > "$env_file" << 'EOF'
# Application
APP_NAME=podoru
APP_ENV=development
APP_PORT=8080
APP_DEBUG=true
REGISTRATION_ENABLED=true

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=podoru
DB_PASSWORD=podoru_secret
DB_NAME=podoru
DB_SSL_MODE=disable

# JWT
JWT_SECRET=dev-jwt-secret-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Encryption (for storing secrets like GitHub tokens)
ENCRYPTION_KEY=dev-32-byte-key-for-aes-256!!!

# Docker
DOCKER_HOST=unix:///var/run/docker.sock

# Traefik (development)
TRAEFIK_DASHBOARD_PORT=8081
TRAEFIK_HTTP_PORT=80
TRAEFIK_HTTPS_PORT=443
TRAEFIK_ACME_EMAIL=dev@localhost
EOF
    fi

    # Update with development-friendly defaults
    sed -i 's/REGISTRATION_ENABLED=false/REGISTRATION_ENABLED=true/' "$env_file" 2>/dev/null || true

    print_success "Created .env file with development settings"
}

install_dependencies() {
    print_info "Installing Go dependencies..."

    cd "$PROJECT_DIR"
    go mod download

    print_success "Dependencies installed"
}

install_dev_tools() {
    print_info "Installing development tools..."

    # Air for hot reload
    if ! command -v air &> /dev/null; then
        print_info "Installing Air (hot reload)..."
        go install github.com/air-verse/air@latest
        print_success "Air installed"
    else
        print_success "Air already installed"
    fi

    # Swag for API docs
    if ! command -v swag &> /dev/null; then
        print_info "Installing Swag (API docs)..."
        go install github.com/swaggo/swag/cmd/swag@v1.16.3
        print_success "Swag installed"
    else
        print_success "Swag already installed"
    fi

    # golang-migrate for migrations
    if ! command -v migrate &> /dev/null; then
        print_info "Installing golang-migrate..."
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
        print_success "golang-migrate installed"
    else
        print_success "golang-migrate already installed"
    fi

    # golangci-lint for linting
    if ! command -v golangci-lint &> /dev/null; then
        print_info "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        print_success "golangci-lint installed"
    else
        print_success "golangci-lint already installed"
    fi
}

start_database() {
    print_info "Starting PostgreSQL database..."

    cd "$PROJECT_DIR"

    local compose_cmd=$(detect_compose_cmd)

    # Check if postgres container is already running
    if docker ps --format '{{.Names}}' | grep -q 'podoru_postgres'; then
        print_success "PostgreSQL is already running"
        return 0
    fi

    # Start only postgres service
    $compose_cmd up -d postgres

    # Wait for postgres to be ready
    print_info "Waiting for PostgreSQL to be ready..."
    local max_attempts=30
    local attempt=1

    while [[ $attempt -le $max_attempts ]]; do
        if $compose_cmd exec -T postgres pg_isready -U podoru -d podoru &> /dev/null; then
            echo ""
            print_success "PostgreSQL is ready"
            return 0
        fi
        echo -n "."
        sleep 1
        ((attempt++))
    done

    echo ""
    print_error "PostgreSQL did not become ready in time"
    return 1
}

generate_docs() {
    print_info "Generating API documentation..."

    cd "$PROJECT_DIR"

    if command -v swag &> /dev/null; then
        swag init -g cmd/podoru/main.go -o docs --parseDependency --parseInternal --parseDepth 3 2>/dev/null
        print_success "API documentation generated"
    else
        print_warning "Swag not found, skipping docs generation"
    fi
}

build_app() {
    print_info "Building application..."

    cd "$PROJECT_DIR"
    go build -o bin/podoru ./cmd/podoru

    print_success "Application built: bin/podoru"
}

create_superadmin_dev() {
    print_info "Creating development superadmin account..."

    local api_url="http://localhost:8080/api/v1/auth/register"
    local max_attempts=10
    local attempt=1

    # Default dev credentials
    local payload=$(cat <<EOF
{
    "email": "admin@localhost",
    "password": "admin123",
    "name": "Dev Admin"
}
EOF
)

    while [[ $attempt -le $max_attempts ]]; do
        local response=$(curl -s -w "\n%{http_code}" -X POST "$api_url" \
            -H "Content-Type: application/json" \
            -d "$payload" 2>/dev/null)

        local http_code=$(echo "$response" | tail -n1)

        if [[ "$http_code" == "201" ]]; then
            print_success "Development superadmin created"
            return 0
        elif [[ "$http_code" == "409" ]]; then
            print_warning "Superadmin already exists"
            return 0
        elif [[ "$http_code" == "000" ]]; then
            echo -n "."
            sleep 1
            ((attempt++))
        else
            # API might not be running, that's okay for setup
            return 0
        fi
    done

    return 0
}

show_status() {
    echo ""
    print_header "Development Setup Complete"

    echo -e "${GREEN}${BOLD}Development Environment Ready!${NC}"
    echo ""
    echo -e "${BOLD}Services:${NC}"
    echo -e "  PostgreSQL:        ${CYAN}localhost:5432${NC}"
    echo -e "  Traefik Dashboard: ${CYAN}http://localhost:8081${NC}"
    echo ""
    echo -e "${BOLD}Default Credentials:${NC}"
    echo -e "  Email:    ${CYAN}admin@localhost${NC}"
    echo -e "  Password: ${CYAN}admin123${NC}"
    echo ""
    echo -e "${BOLD}Quick Commands:${NC}"
    echo "  make dev          # Start with hot reload"
    echo "  make run          # Build and run"
    echo "  make test         # Run tests"
    echo "  make docker-up    # Start all Docker services"
    echo "  make docker-down  # Stop Docker services"
    echo "  make docs         # Generate API docs"
    echo ""
    echo -e "${BOLD}API Endpoints:${NC}"
    echo -e "  Application:    ${CYAN}http://localhost:8080${NC}"
    echo -e "  Documentation:  ${CYAN}http://localhost:8080/api/v1/docs${NC}"
    echo -e "  Health Check:   ${CYAN}http://localhost:8080/health${NC}"
    echo ""
}

#######################################
# Main functions
#######################################
run_checks() {
    print_header "Checking Prerequisites"

    local failed=0

    check_go || ((failed++))
    check_docker || ((failed++))
    check_docker_compose || ((failed++))
    check_ports || true  # Don't fail on ports, just warn
    check_make

    echo ""
    if [[ $failed -gt 0 ]]; then
        print_error "$failed check(s) failed"
        return 1
    else
        print_success "All prerequisites met"
        return 0
    fi
}

run_setup() {
    print_header "Podoru Development Setup"

    # Run checks
    if ! run_checks; then
        print_error "Prerequisites check failed. Install missing dependencies and try again."
        exit 1
    fi

    echo ""
    read -p "Proceed with development setup? [Y/n]: " confirm
    if [[ "$confirm" =~ ^[Nn]$ ]]; then
        print_info "Setup cancelled"
        exit 0
    fi

    # Setup steps
    print_step "Setting Up Development Environment"

    create_env_file
    install_dependencies
    install_dev_tools
    start_database
    generate_docs
    build_app

    show_status

    echo -e "${YELLOW}To start the application:${NC}"
    echo "  cd $PROJECT_DIR"
    echo "  make dev   # or: make run"
    echo ""
}

run_quick() {
    print_header "Quick Development Start"

    cd "$PROJECT_DIR"

    # Ensure .env exists
    if [[ ! -f ".env" ]]; then
        create_env_file
    fi

    # Start database if not running
    if ! docker ps --format '{{.Names}}' | grep -q 'podoru_postgres'; then
        start_database
    else
        print_success "PostgreSQL is already running"
    fi

    # Start Traefik if not running
    if ! docker ps --format '{{.Names}}' | grep -q 'podoru_traefik'; then
        print_info "Starting Traefik..."
        docker-compose up -d traefik
        print_success "Traefik started (Dashboard: http://localhost:8081)"
    else
        print_success "Traefik is already running"
    fi

    # Build and run
    print_info "Starting application..."
    echo ""
    echo -e "${CYAN}Starting with hot reload...${NC}"
    echo -e "${CYAN}Press Ctrl+C to stop${NC}"
    echo ""

    if command -v air &> /dev/null; then
        air
    else
        go run ./cmd/podoru
    fi
}

run_reset() {
    print_header "Reset Development Environment"

    echo -e "${YELLOW}This will:${NC}"
    echo "  - Stop all Docker containers"
    echo "  - Remove database volume (all data)"
    echo "  - Remove .env file"
    echo ""

    read -p "Are you sure? [y/N]: " confirm
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
        print_info "Reset cancelled"
        exit 0
    fi

    cd "$PROJECT_DIR"

    local compose_cmd=$(detect_compose_cmd)

    print_info "Stopping Docker services..."
    $compose_cmd down -v 2>/dev/null || true

    print_info "Removing .env file..."
    rm -f .env

    print_info "Cleaning build artifacts..."
    rm -rf bin/
    rm -f coverage.out coverage.html

    print_success "Development environment reset"
    echo ""
    echo "Run './scripts/dev-setup.sh setup' to set up again"
}

#######################################
# Command handler
#######################################
case "${1:-}" in
    check)
        run_checks
        ;;
    setup)
        run_setup
        ;;
    quick|start)
        run_quick
        ;;
    reset)
        run_reset
        ;;
    *)
        echo -e "${BLUE}${BOLD}Podoru Development Setup${NC}"
        echo ""
        echo "Usage: $0 <command>"
        echo ""
        echo "Commands:"
        echo "  check   Check prerequisites only"
        echo "  setup   Full development setup"
        echo "  quick   Quick start (database + app)"
        echo "  reset   Reset development environment"
        echo ""
        exit 1
        ;;
esac
