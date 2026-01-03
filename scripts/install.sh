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
        print_warning ".env.prod not found (will be created during setup)"
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

check_htpasswd() {
    if command -v htpasswd &> /dev/null; then
        return 0
    elif command -v openssl &> /dev/null; then
        return 0
    else
        return 1
    fi
}

#######################################
# Prompt for superadmin only (when using existing config)
#######################################
prompt_superadmin_only() {
    print_header "Superadmin Account Setup"

    echo -e "${CYAN}Create the superadmin account for Podoru.${NC}"
    echo ""

    read -p "Superadmin name [Administrator]: " INPUT_SUPERADMIN_NAME
    INPUT_SUPERADMIN_NAME=${INPUT_SUPERADMIN_NAME:-Administrator}

    while true; do
        read -p "Superadmin email: " INPUT_SUPERADMIN_EMAIL
        if [[ -z "$INPUT_SUPERADMIN_EMAIL" ]]; then
            print_error "Email is required"
        elif [[ ! "$INPUT_SUPERADMIN_EMAIL" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
            print_error "Invalid email format"
        else
            break
        fi
    done

    while true; do
        read -s -p "Superadmin password: " INPUT_SUPERADMIN_PASS
        echo ""
        if [[ -z "$INPUT_SUPERADMIN_PASS" ]]; then
            print_error "Password is required"
        elif [[ ${#INPUT_SUPERADMIN_PASS} -lt 8 ]]; then
            print_error "Password must be at least 8 characters"
        else
            read -s -p "Confirm password: " INPUT_SUPERADMIN_PASS_CONFIRM
            echo ""
            if [[ "$INPUT_SUPERADMIN_PASS" != "$INPUT_SUPERADMIN_PASS_CONFIRM" ]]; then
                print_error "Passwords do not match"
            else
                break
            fi
        fi
    done

    # Export for use in create_superadmin function
    export SUPERADMIN_NAME="$INPUT_SUPERADMIN_NAME"
    export SUPERADMIN_EMAIL="$INPUT_SUPERADMIN_EMAIL"
    export SUPERADMIN_PASS="$INPUT_SUPERADMIN_PASS"

    print_success "Superadmin credentials set"
}

#######################################
# Interactive Setup Wizard
#######################################
run_setup_wizard() {
    print_header "Podoru Setup Wizard"

    local env_file="$PROJECT_DIR/.env.prod"

    echo -e "${CYAN}This wizard will help you configure Podoru for production.${NC}"
    echo -e "${CYAN}Press Enter to accept default values shown in [brackets].${NC}"
    echo ""

    # Domain
    print_step "Step 1/6: Domain Configuration"
    echo "Enter the domain where Podoru will be accessible."
    echo "Example: podoru.example.com, panel.myserver.com"
    echo ""
    while true; do
        read -p "Domain: " INPUT_DOMAIN
        if [[ -z "$INPUT_DOMAIN" ]]; then
            print_error "Domain is required"
        elif [[ ! "$INPUT_DOMAIN" =~ ^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*$ ]]; then
            print_error "Invalid domain format"
        else
            break
        fi
    done

    # Email for SSL
    print_step "Step 2/6: SSL Certificate Email"
    echo "Enter your email for Let's Encrypt SSL certificates."
    echo "You'll receive expiration notices at this address."
    echo ""
    while true; do
        read -p "Email: " INPUT_EMAIL
        if [[ -z "$INPUT_EMAIL" ]]; then
            print_error "Email is required"
        elif [[ ! "$INPUT_EMAIL" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
            print_error "Invalid email format"
        else
            break
        fi
    done

    # Admin password for Traefik dashboard
    print_step "Step 3/6: Traefik Dashboard Credentials"
    echo "Set up credentials for the Traefik dashboard."
    echo "Dashboard URL: https://traefik.${INPUT_DOMAIN}/dashboard/"
    echo ""
    read -p "Admin username [admin]: " INPUT_ADMIN_USER
    INPUT_ADMIN_USER=${INPUT_ADMIN_USER:-admin}

    while true; do
        read -s -p "Admin password: " INPUT_ADMIN_PASS
        echo ""
        if [[ -z "$INPUT_ADMIN_PASS" ]]; then
            print_error "Password is required"
        elif [[ ${#INPUT_ADMIN_PASS} -lt 8 ]]; then
            print_error "Password must be at least 8 characters"
        else
            read -s -p "Confirm password: " INPUT_ADMIN_PASS_CONFIRM
            echo ""
            if [[ "$INPUT_ADMIN_PASS" != "$INPUT_ADMIN_PASS_CONFIRM" ]]; then
                print_error "Passwords do not match"
            else
                break
            fi
        fi
    done

    # Generate htpasswd
    if command -v htpasswd &> /dev/null; then
        TRAEFIK_AUTH=$(htpasswd -nb "$INPUT_ADMIN_USER" "$INPUT_ADMIN_PASS")
    else
        # Fallback using openssl
        TRAEFIK_AUTH="${INPUT_ADMIN_USER}:$(openssl passwd -apr1 "$INPUT_ADMIN_PASS")"
    fi

    # Superadmin account
    print_step "Step 4/6: Superadmin Account"
    echo "Create the superadmin account for Podoru."
    echo "This account will have full administrative access."
    echo ""

    read -p "Superadmin name [Administrator]: " INPUT_SUPERADMIN_NAME
    INPUT_SUPERADMIN_NAME=${INPUT_SUPERADMIN_NAME:-Administrator}

    while true; do
        read -p "Superadmin email: " INPUT_SUPERADMIN_EMAIL
        if [[ -z "$INPUT_SUPERADMIN_EMAIL" ]]; then
            print_error "Email is required"
        elif [[ ! "$INPUT_SUPERADMIN_EMAIL" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
            print_error "Invalid email format"
        else
            break
        fi
    done

    while true; do
        read -s -p "Superadmin password: " INPUT_SUPERADMIN_PASS
        echo ""
        if [[ -z "$INPUT_SUPERADMIN_PASS" ]]; then
            print_error "Password is required"
        elif [[ ${#INPUT_SUPERADMIN_PASS} -lt 8 ]]; then
            print_error "Password must be at least 8 characters"
        else
            read -s -p "Confirm password: " INPUT_SUPERADMIN_PASS_CONFIRM
            echo ""
            if [[ "$INPUT_SUPERADMIN_PASS" != "$INPUT_SUPERADMIN_PASS_CONFIRM" ]]; then
                print_error "Passwords do not match"
            else
                break
            fi
        fi
    done

    # Export for use in create_superadmin function
    export SUPERADMIN_NAME="$INPUT_SUPERADMIN_NAME"
    export SUPERADMIN_EMAIL="$INPUT_SUPERADMIN_EMAIL"
    export SUPERADMIN_PASS="$INPUT_SUPERADMIN_PASS"

    # Database settings
    print_step "Step 5/6: Database Configuration"
    echo "Configure PostgreSQL database settings."
    echo ""
    read -p "Database user [podoru]: " INPUT_DB_USER
    INPUT_DB_USER=${INPUT_DB_USER:-podoru}

    read -p "Database name [podoru]: " INPUT_DB_NAME
    INPUT_DB_NAME=${INPUT_DB_NAME:-podoru}

    # Application settings
    print_step "Step 6/6: Application Settings"
    echo "Configure application settings."
    echo ""
    echo "Allow public user registration after setup?"
    echo "(You can change this later in .env.prod)"
    read -p "Enable public registration [y/N]: " INPUT_REGISTRATION
    if [[ "$INPUT_REGISTRATION" =~ ^[Yy]$ ]]; then
        REGISTRATION_ENABLED="true"
    else
        REGISTRATION_ENABLED="false"
    fi

    # Generate secrets
    print_step "Generating Secure Secrets"
    local db_password=$(openssl rand -base64 24 | tr -dc 'a-zA-Z0-9' | head -c 24)
    local jwt_secret=$(openssl rand -base64 32)
    local encryption_key=$(openssl rand -base64 24 | head -c 32)

    print_success "Database password generated"
    print_success "JWT secret generated"
    print_success "Encryption key generated"

    # Write .env.prod file
    print_step "Creating Configuration File"

    cat > "$env_file" << EOF
# ===========================================
# Podoru Production Configuration
# Generated: $(date -u '+%Y-%m-%d %H:%M:%S UTC')
# ===========================================

# Domain Configuration
DOMAIN=${INPUT_DOMAIN}
ACME_EMAIL=${INPUT_EMAIL}
TRAEFIK_DOMAIN=traefik.${INPUT_DOMAIN}

# Database Configuration
DB_USER=${INPUT_DB_USER}
DB_PASSWORD=${db_password}
DB_NAME=${INPUT_DB_NAME}

# JWT Configuration
JWT_SECRET=${jwt_secret}
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Encryption Key (32 characters for AES-256)
ENCRYPTION_KEY=${encryption_key}

# Traefik Dashboard Auth
TRAEFIK_DASHBOARD_AUTH=${TRAEFIK_AUTH}

# Application Settings
REGISTRATION_ENABLED=${REGISTRATION_ENABLED}

# Docker Image (optional)
# PODORU_IMAGE=podoru:latest
EOF

    chmod 600 "$env_file"
    print_success "Configuration saved to .env.prod"

    # Summary
    print_header "Configuration Summary"
    echo -e "  ${BOLD}Domain & SSL:${NC}"
    echo -e "    Domain:              ${GREEN}${INPUT_DOMAIN}${NC}"
    echo -e "    SSL Email:           ${GREEN}${INPUT_EMAIL}${NC}"
    echo ""
    echo -e "  ${BOLD}Superadmin Account:${NC}"
    echo -e "    Name:                ${GREEN}${INPUT_SUPERADMIN_NAME}${NC}"
    echo -e "    Email:               ${GREEN}${INPUT_SUPERADMIN_EMAIL}${NC}"
    echo ""
    echo -e "  ${BOLD}Traefik Dashboard:${NC}"
    echo -e "    URL:                 ${GREEN}https://traefik.${INPUT_DOMAIN}/dashboard/${NC}"
    echo -e "    Username:            ${GREEN}${INPUT_ADMIN_USER}${NC}"
    echo ""
    echo -e "  ${BOLD}Database:${NC}"
    echo -e "    User:                ${GREEN}${INPUT_DB_USER}${NC}"
    echo -e "    Name:                ${GREEN}${INPUT_DB_NAME}${NC}"
    echo ""
    echo -e "  ${BOLD}Settings:${NC}"
    echo -e "    Public Registration: ${GREEN}${REGISTRATION_ENABLED}${NC}"
    echo ""
    echo -e "  ${YELLOW}Secrets have been auto-generated and saved to .env.prod${NC}"
    echo ""

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

# Traefik Settings
TRAEFIK_ENABLED=true
TRAEFIK_NETWORK=podoru_traefik

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
    docker compose -f docker-compose.prod.yml --env-file .env.prod build

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

create_superadmin() {
    print_info "Creating superadmin account..."

    # Wait a bit for the API to be fully ready
    sleep 3

    local api_url="http://localhost:8080/api/v1/auth/register"
    local max_attempts=10
    local attempt=1

    # Prepare JSON payload (escape quotes for docker exec)
    local payload="{\"email\":\"${SUPERADMIN_EMAIL}\",\"password\":\"${SUPERADMIN_PASS}\",\"name\":\"${SUPERADMIN_NAME}\"}"

    cd "$PROJECT_DIR"

    while [[ $attempt -le $max_attempts ]]; do
        # Make the API request via docker exec (since app isn't exposed on host)
        local response=$(docker compose -f docker-compose.prod.yml exec -T podoru \
            wget -q -O - --header="Content-Type: application/json" \
            --post-data="$payload" "$api_url" 2>&1) || true

        # Check response for success indicators
        if echo "$response" | grep -q '"success":true'; then
            print_success "Superadmin account created successfully"
            # Clear sensitive data
            unset SUPERADMIN_PASS
            return 0
        elif echo "$response" | grep -q 'already exists\|CONFLICT\|409'; then
            print_warning "User already exists (may be from previous installation)"
            return 0
        elif echo "$response" | grep -q 'Connection refused\|no route\|Unable to connect'; then
            # API not ready yet
            echo -n "."
            sleep 2
            ((attempt++))
        else
            # Some other response - might be an error or might be success
            if [[ -n "$response" ]]; then
                echo "$response"
            fi
            echo -n "."
            sleep 2
            ((attempt++))
        fi
    done

    print_error "API not responding after $max_attempts attempts"
    echo "You can create the superadmin manually after startup:"
    echo "  curl -X POST https://\$DOMAIN/api/v1/auth/register \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"email\":\"...\",\"password\":\"...\",\"name\":\"...\"}'"
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
    echo -e "${GREEN}${BOLD}========================================${NC}"
    echo -e "${GREEN}${BOLD}  SAVE THESE CREDENTIALS SECURELY!${NC}"
    echo -e "${GREEN}${BOLD}========================================${NC}"
    echo ""
    echo -e "${BOLD}Superadmin Account:${NC}"
    echo -e "  Email:    ${CYAN}${SUPERADMIN_EMAIL}${NC}"
    echo -e "  Password: ${CYAN}(the password you entered during setup)${NC}"
    echo ""
    echo -e "${BOLD}Traefik Dashboard:${NC}"
    echo -e "  URL:      ${CYAN}https://${TRAEFIK_DOMAIN:-traefik.$DOMAIN}/dashboard/${NC}"
    echo -e "  Username: ${CYAN}(the username you entered during setup)${NC}"
    echo ""
    echo -e "${GREEN}${BOLD}========================================${NC}"
    echo ""
    echo -e "${GREEN}${BOLD}Access Points:${NC}"
    echo -e "  Application:       ${CYAN}https://$DOMAIN${NC}"
    echo -e "  API Documentation: ${CYAN}https://$DOMAIN/api/v1/docs${NC}"
    echo -e "  Traefik Dashboard: ${CYAN}https://${TRAEFIK_DOMAIN:-traefik.$DOMAIN}/dashboard/${NC}"
    echo ""
    echo -e "${GREEN}${BOLD}Next Steps:${NC}"
    echo "  1. Point your domain DNS to this server's IP address"
    echo "  2. Wait for DNS propagation (may take a few minutes)"
    echo "  3. Access https://$DOMAIN and login with your superadmin account"
    echo ""
    echo -e "${GREEN}${BOLD}Useful Commands:${NC}"
    echo "  View logs:     make prod-logs"
    echo "  Stop services: make prod-down"
    echo "  Backup DB:     make prod-backup"
    echo ""
}

#######################################
# Main functions
#######################################
run_checks() {
    print_header "Pre-flight Checks"

    local failed=0
    local env_missing=0

    check_root
    check_os

    check_docker || ((failed++))
    check_docker_compose || ((failed++))
    check_docker_socket || ((failed++))
    check_ports || ((failed++))
    check_memory || ((failed++))
    check_disk || ((failed++))
    check_env_file || env_missing=1

    echo ""
    if [[ $failed -gt 0 ]]; then
        print_error "$failed check(s) failed"
        return 1
    elif [[ $env_missing -eq 1 ]]; then
        print_warning "Environment not configured (run setup wizard during install)"
        return 2
    else
        print_success "All pre-flight checks passed"
        return 0
    fi
}

run_install() {
    print_header "Podoru Production Installation"

    local env_file="$PROJECT_DIR/.env.prod"

    # Run system checks first (excluding env file)
    print_step "Running System Checks"
    local failed=0

    check_root
    check_os

    check_docker || ((failed++))
    check_docker_compose || ((failed++))
    check_docker_socket || ((failed++))
    check_ports || ((failed++))
    check_memory || ((failed++))
    check_disk || ((failed++))

    if [[ $failed -gt 0 ]]; then
        echo ""
        print_error "System checks failed. Fix the issues above and try again."
        exit 1
    fi

    echo ""
    print_success "System checks passed"

    # Check if .env.prod exists and is valid
    if [[ -f "$env_file" ]]; then
        source "$env_file"
        if [[ -n "$DOMAIN" && -n "$DB_PASSWORD" && -n "$JWT_SECRET" ]]; then
            print_info "Existing configuration found for: $DOMAIN"
            read -p "Use existing configuration? [Y/n]: " use_existing
            if [[ "$use_existing" =~ ^[Nn]$ ]]; then
                run_setup_wizard
            else
                # Still need superadmin credentials for existing config
                prompt_superadmin_only
            fi
        else
            run_setup_wizard
        fi
    else
        run_setup_wizard
    fi

    # Final confirmation
    echo ""
    read -p "Proceed with installation? [Y/n]: " confirm
    if [[ "$confirm" =~ ^[Nn]$ ]]; then
        print_info "Installation cancelled"
        exit 0
    fi

    # Run installation
    print_step "Installing Podoru"
    create_directories
    pull_images
    build_app
    start_services
    wait_for_healthy
    create_superadmin
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
    setup)
        run_setup_wizard
        ;;
    secrets)
        generate_secrets
        ;;
    env)
        create_env_template
        ;;
    *)
        echo -e "${BLUE}${BOLD}Podoru Installation Script${NC}"
        echo ""
        echo "Usage: $0 <command>"
        echo ""
        echo "Commands:"
        echo "  check    Run pre-flight checks only"
        echo "  install  Run full installation (interactive)"
        echo "  setup    Run setup wizard only (create .env.prod)"
        echo "  secrets  Generate secure secrets"
        echo "  env      Create .env.prod template (non-interactive)"
        echo ""
        exit 1
        ;;
esac
