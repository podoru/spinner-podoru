# Installation

This guide covers installing Podoru in a production environment.

## Prerequisites

Ensure your server meets these requirements:

| Requirement | Minimum | Recommended |
|-------------|---------|-------------|
| OS | Ubuntu 20.04 / Debian 11 | Ubuntu 22.04 |
| Docker | 20.10.0 | Latest |
| Docker Compose | 2.0.0 | Latest |
| RAM | 2 GB | 4 GB |
| Storage | 20 GB | 50 GB |

### Required Ports

| Port | Service | Description |
|------|---------|-------------|
| 80 | Traefik | HTTP traffic (redirects to HTTPS) |
| 443 | Traefik | HTTPS traffic |
| 8080 | Podoru API | Application API |

## Automated Installation

### Step 1: Download the Install Script

```bash
curl -fsSL https://raw.githubusercontent.com/podoru/spinner-podoru/master/scripts/install.sh -o install.sh
chmod +x install.sh
```

### Step 2: Run Pre-flight Checks

Verify your system is ready:

```bash
./install.sh check
```

This checks:
- Operating system compatibility
- Docker installation and version
- Docker Compose installation
- Port availability
- Docker socket access

### Step 3: Run Installation

```bash
./install.sh install
```

You will be prompted for:

| Setting | Description | Example |
|---------|-------------|---------|
| Domain | Your Podoru domain | `podoru.example.com` |
| ACME Email | Email for Let's Encrypt | `admin@example.com` |
| Admin Email | Superadmin email | `admin@example.com` |
| Admin Password | Superadmin password | (secure password) |

### Step 4: Verify Installation

After installation completes:

```bash
# Check services are running
docker compose -f docker-compose.prod.yml ps

# Test API health
curl https://your-domain.com/health
```

## Manual Installation

If you prefer manual installation:

### 1. Clone Repository

```bash
git clone git@github.com:podoru/spinner-podoru.git
cd spinner-podoru
```

### 2. Create Environment File

```bash
cp .env.example .env.prod
```

Edit `.env.prod` with your settings:

```bash
# Domain
DOMAIN=podoru.example.com
ACME_EMAIL=admin@example.com

# Database
DB_USER=podoru
DB_PASSWORD=your-secure-db-password
DB_NAME=podoru

# Security
JWT_SECRET=your-32-byte-jwt-secret-key-here
ENCRYPTION_KEY=your-32-byte-encryption-key!!

# Traefik
TRAEFIK_ENABLED=true
TRAEFIK_NETWORK=podoru_traefik
TRAEFIK_DASHBOARD_AUTH=admin:$apr1$...  # htpasswd generated

# App
REGISTRATION_ENABLED=false
```

### 3. Start Services

```bash
docker compose -f docker-compose.prod.yml up -d
```

### 4. Create Superadmin

```bash
# Enable registration temporarily
docker compose -f docker-compose.prod.yml exec podoru \
  wget -q -O - --header="Content-Type: application/json" \
  --post-data='{"email":"admin@example.com","password":"secure-password","name":"Admin"}' \
  http://localhost:8080/api/v1/auth/register
```

## DNS Configuration

Point your domain to your server's IP address:

```
Type: A
Name: podoru (or @ for root)
Value: YOUR_SERVER_IP
TTL: 300
```

If using subdomains for services, add a wildcard:

```
Type: A
Name: *.podoru
Value: YOUR_SERVER_IP
TTL: 300
```

## SSL Certificates

Podoru uses Traefik with Let's Encrypt for automatic SSL. Certificates are:

- Automatically requested on first access
- Stored in a Docker volume (`traefik_data`)
- Renewed automatically before expiration

## Updating

To update Podoru:

```bash
cd /path/to/spinner-podoru
git pull
docker compose -f docker-compose.prod.yml up -d --build
```

## Uninstallation

To completely remove Podoru:

```bash
# Stop and remove containers
docker compose -f docker-compose.prod.yml down

# Remove volumes (WARNING: deletes all data)
docker compose -f docker-compose.prod.yml down -v

# Remove installation directory
cd .. && rm -rf spinner-podoru
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port
sudo lsof -i :80
sudo lsof -i :443

# Stop conflicting service
sudo systemctl stop nginx  # or apache2
```

### Docker Permission Denied

```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Re-login or run
newgrp docker
```

### SSL Certificate Issues

```bash
# Check Traefik logs
docker compose -f docker-compose.prod.yml logs traefik

# Verify domain DNS
dig +short your-domain.com
```

### Database Connection Failed

```bash
# Check PostgreSQL is running
docker compose -f docker-compose.prod.yml ps postgres

# Check logs
docker compose -f docker-compose.prod.yml logs postgres
```
