#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REQUIRED_DOCKER_VERSION="20.10.0"
REQUIRED_COMPOSE_VERSION="2.0.0"
MIN_MEMORY_MB=2048
MIN_DISK_GB=10
REQUIRED_PORTS=(80 443 5432 8080)

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

#######################################
# Version comparison
#######################################
version_gte() {
    [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" = "$2" ]
}

#######################################
# Pre-flight checks
#######################################
check_root() {
    if [[ $EUID -eq 0 ]]; then
        print_warning "Running as root. Consider using a non-root user with docker group."
    fi
}

check_os() {
    print_info "Checking operating system..."
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        print_success "OS: $PRETTY_NAME"
    else
        print_warning "Could not determine OS version"
    fi
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
        print_success "Docker version: $DOCKER_VERSION (required: >=$REQUIRED_DOCKER_VERSION)"
    else
        print_error "Docker version $DOCKER_VERSION is too old (required: >=$REQUIRED_DOCKER_VERSION)"
        return 1
    fi

    # Check if docker daemon is running
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

    # Check for docker compose (v2) or docker-compose (v1)
    if docker compose version &> /dev/null; then
        COMPOSE_VERSION=$(docker compose version --short 2>/dev/null || echo "0.0.0")
        COMPOSE_CMD="docker compose"
    elif command -v docker-compose &> /dev/null; then
        COMPOSE_VERSION=$(docker-compose version --short 2>/dev/null || echo "0.0.0")
        COMPOSE_CMD="docker-compose"
    else
        print_error "Docker Compose is not installed"
        return 1
    fi

    if version_gte "$COMPOSE_VERSION" "$REQUIRED_COMPOSE_VERSION"; then
        print_success "Docker Compose version: $COMPOSE_VERSION (required: >=$REQUIRED_COMPOSE_VERSION)"
    else
        print_error "Docker Compose version $COMPOSE_VERSION is too old (required: >=$REQUIRED_COMPOSE_VERSION)"
        return 1
    fi

    return 0
}

check_docker_socket() {
    print_info "Checking Docker socket access..."

    if [[ ! -S /var/run/docker.sock ]]; then
        print_error "Docker socket not found at /var/run/docker.sock"
        return 1
    fi

    if ! docker ps &> /dev/null; then
        print_error "Cannot access Docker socket. Add user to docker group:"
        echo "  sudo usermod -aG docker \$USER"
        echo "  Then log out and back in"
        return 1
    fi

    print_success "Docker socket is accessible"
    return 0
}

check_ports() {
    print_info "Checking port availability..."
    local ports_in_use=()

    for port in "${REQUIRED_PORTS[@]}"; do
        if ss -tuln 2>/dev/null | grep -q ":${port} " || \
           netstat -tuln 2>/dev/null | grep -q ":${port} "; then
            ports_in_use+=($port)
            print_error "Port $port is already in use"

            # Try to identify what's using the port
            local process=$(ss -tlnp 2>/dev/null | grep ":${port} " | awk '{print $7}' | head -1)
            if [[ -n "$process" ]]; then
                echo "  Process: $process"
            fi
        else
            print_success "Port $port is available"
        fi
    done

    if [[ ${#ports_in_use[@]} -gt 0 ]]; then
        echo ""
        print_warning "The following ports are in use: ${ports_in_use[*]}"
        echo "  Stop the services using these ports before continuing"
        return 1
    fi

    return 0
}

check_memory() {
    print_info "Checking available memory..."

    local total_mem_kb=$(grep MemTotal /proc/meminfo | awk '{print $2}')
    local total_mem_mb=$((total_mem_kb / 1024))

    if [[ $total_mem_mb -ge $MIN_MEMORY_MB ]]; then
        print_success "Memory: ${total_mem_mb}MB available (required: >=${MIN_MEMORY_MB}MB)"
    else
        print_error "Insufficient memory: ${total_mem_mb}MB (required: >=${MIN_MEMORY_MB}MB)"
        return 1
    fi

    return 0
}

check_disk() {
    print_info "Checking available disk space..."

    local available_gb=$(df -BG "$PROJECT_DIR" | tail -1 | awk '{print $4}' | sed 's/G//')

    if [[ $available_gb -ge $MIN_DISK_GB ]]; then
        print_success "Disk space: ${available_gb}GB available (required: >=${MIN_DISK_GB}GB)"
    else
        print_error "Insufficient disk space: ${available_gb}GB (required: >=${MIN_DISK_GB}GB)"
        return 1
    fi

    return 0
}

check_env_file() {
    print_info "Checking environment configuration..."

    local env_file="$PROJECT_DIR/.env.prod"

    if [[ ! -f "$env_file" ]]; then
        print_warning ".env.prod not found"
        echo "  Creating from template..."
        create_env_template
        print_info "Please edit $env_file with your configuration"
        return 1
    fi

    # Check required variables
    source "$env_file"
    local missing_vars=()

    [[ -z "$DOMAIN" ]] && missing_vars+=("DOMAIN")
    [[ -z "$DB_PASSWORD" ]] && missing_vars+=("DB_PASSWORD")
    [[ -z "$JWT_SECRET" ]] && missing_vars+=("JWT_SECRET")
    [[ -z "$ENCRYPTION_KEY" ]] && missing_vars+=("ENCRYPTION_KEY")
    [[ -z "$ACME_EMAIL" ]] && missing_vars+=("ACME_EMAIL")
    [[ -z "$TRAEFIK_DASHBOARD_AUTH" ]] && missing_vars+=("TRAEFIK_DASHBOARD_AUTH")

    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        print_error "Missing required environment variables:"
        for var in "${missing_vars[@]}"; do
            echo "  - $var"
        done
        return 1
    fi

    print_success "Environment configuration is valid"
    return 0
}

#######################################
# Setup functions
#######################################
create_env_template() {
    local env_file="$PROJECT_DIR/.env.prod"

    cat > "$env_file" << 'EOF'
# ===========================================
# Podoru Production Configuration
# ===========================================

# Domain Configuration (REQUIRED)
# Your domain name (e.g., podoru.example.com)
DOMAIN=

# ACME/Let's Encrypt email for SSL certificates (REQUIRED)
ACME_EMAIL=

# Traefik dashboard domain (optional, defaults to traefik.${DOMAIN})
# TRAEFIK_DOMAIN=

# Database Configuration (REQUIRED)
DB_USER=podoru
DB_PASSWORD=
DB_NAME=podoru

# JWT Configuration (REQUIRED)
# Generate with: openssl rand -base64 32
JWT_SECRET=
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Encryption Key (REQUIRED)
# Must be exactly 32 characters for AES-256
# Generate with: openssl rand -base64 24 | head -c 32
ENCRYPTION_KEY=

# Traefik Dashboard Auth (REQUIRED)
# Generate with: htpasswd -nb admin your_password
# Format: username:password_hash
TRAEFIK_DASHBOARD_AUTH=

# Application Settings
REGISTRATION_ENABLED=false

# Docker Image (optional)
# PODORU_IMAGE=podoru:latest
EOF

    chmod 600 "$env_file"
    print_success "Created $env_file"
}

generate_secrets() {
    print_info "Generating secure secrets..."

    local db_password=$(openssl rand -base64 24 | tr -dc 'a-zA-Z0-9' | head -c 24)
    local jwt_secret=$(openssl rand -base64 32)
    local encryption_key=$(openssl rand -base64 24 | head -c 32)

    echo ""
    echo "Generated secrets (save these securely):"
    echo "========================================="
    echo "DB_PASSWORD=$db_password"
    echo "JWT_SECRET=$jwt_secret"
    echo "ENCRYPTION_KEY=$encryption_key"
    echo ""
    echo "For TRAEFIK_DASHBOARD_AUTH, run:"
    echo "  htpasswd -nb admin your_password"
    echo ""
}

create_directories() {
    print_info "Creating required directories..."

    mkdir -p "$PROJECT_DIR/backups"
    chmod 750 "$PROJECT_DIR/backups"

    print_success "Created backup directory"
}

pull_images() {
    print_info "Pulling Docker images..."

    cd "$PROJECT_DIR"
    docker compose -f docker-compose.prod.yml pull

    print_success "Docker images pulled"
}

build_app() {
    print_info "Building Podoru application..."

    cd "$PROJECT_DIR"
    docker compose -f docker-compose.prod.yml build

    print_success "Application built"
}

start_services() {
    print_info "Starting services..."

    cd "$PROJECT_DIR"
    docker compose -f docker-compose.prod.yml --env-file .env.prod up -d

    print_success "Services started"
}

wait_for_healthy() {
    print_info "Waiting for services to be healthy..."

    local max_attempts=30
    local attempt=1

    while [[ $attempt -le $max_attempts ]]; do
        if docker compose -f docker-compose.prod.yml ps 2>/dev/null | grep -q "unhealthy\|starting"; then
            echo -n "."
            sleep 2
            ((attempt++))
        else
            echo ""
            print_success "All services are healthy"
            return 0
        fi
    done

    echo ""
    print_error "Services did not become healthy in time"
    docker compose -f docker-compose.prod.yml ps
    return 1
}

show_status() {
    echo ""
    print_header "Installation Complete"

    cd "$PROJECT_DIR"
    source .env.prod

    echo "Services Status:"
    docker compose -f docker-compose.prod.yml ps

    echo ""
    echo "Access Points:"
    echo "  - Application:       https://$DOMAIN"
    echo "  - API Documentation: https://$DOMAIN/api/v1/docs"
    echo "  - Traefik Dashboard: https://${TRAEFIK_DOMAIN:-traefik.$DOMAIN}/dashboard/"
    echo ""
    echo "Useful Commands:"
    echo "  - View logs:     make prod-logs"
    echo "  - Stop services: make prod-down"
    echo "  - Backup DB:     make prod-backup"
    echo ""
}

#######################################
# Main functions
#######################################
run_checks() {
    print_header "Pre-flight Checks"

    local failed=0

    check_root
    check_os

    check_docker || ((failed++))
    check_docker_compose || ((failed++))
    check_docker_socket || ((failed++))
    check_ports || ((failed++))
    check_memory || ((failed++))
    check_disk || ((failed++))
    check_env_file || ((failed++))

    echo ""
    if [[ $failed -gt 0 ]]; then
        print_error "$failed check(s) failed"
        return 1
    else
        print_success "All pre-flight checks passed"
        return 0
    fi
}

run_install() {
    print_header "Podoru Production Installation"

    # Run checks first
    if ! run_checks; then
        print_error "Pre-flight checks failed. Fix the issues and try again."
        exit 1
    fi

    echo ""
    read -p "Proceed with installation? [y/N] " -n 1 -r
    echo ""

    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Installation cancelled"
        exit 0
    fi

    create_directories
    pull_images
    build_app
    start_services
    wait_for_healthy
    show_status
}

#######################################
# Command handler
#######################################
case "${1:-}" in
    check)
        run_checks
        ;;
    install)
        run_install
        ;;
    secrets)
        generate_secrets
        ;;
    env)
        create_env_template
        ;;
    *)
        echo "Podoru Installation Script"
        echo ""
        echo "Usage: $0 <command>"
        echo ""
        echo "Commands:"
        echo "  check    Run pre-flight checks only"
        echo "  install  Run full installation"
        echo "  secrets  Generate secure secrets"
        echo "  env      Create .env.prod template"
        echo ""
        exit 1
        ;;
esac
