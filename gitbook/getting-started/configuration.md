# Configuration

Podoru is configured through environment variables and a YAML configuration file.

## Configuration Priority

Configuration is loaded in this order (later overrides earlier):

1. Default values in code
2. `configs/config.yaml` file
3. Environment variables

## Environment Variables

### Application

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | `podoru` |
| `APP_ENV` | Environment (`development` / `production`) | `development` |
| `APP_PORT` | HTTP server port | `8080` |
| `APP_DEBUG` | Enable debug mode | `false` |
| `REGISTRATION_ENABLED` | Allow new user registration | `false` |

### Database

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database username | `podoru` |
| `DB_PASSWORD` | Database password | (required) |
| `DB_NAME` | Database name | `podoru` |
| `DB_SSL_MODE` | SSL mode (`disable` / `require`) | `disable` |

### Security

| Variable | Description | Default |
|----------|-------------|---------|
| `JWT_SECRET` | Secret key for JWT tokens (32+ chars) | (required) |
| `JWT_ACCESS_EXPIRY` | Access token lifetime | `15m` |
| `JWT_REFRESH_EXPIRY` | Refresh token lifetime | `168h` (7 days) |
| `ENCRYPTION_KEY` | AES-256 encryption key (32 chars) | (required) |

### Docker

| Variable | Description | Default |
|----------|-------------|---------|
| `DOCKER_HOST` | Docker daemon socket | `unix:///var/run/docker.sock` |

### Traefik

| Variable | Description | Default |
|----------|-------------|---------|
| `TRAEFIK_ENABLED` | Enable Traefik integration | `true` |
| `TRAEFIK_NETWORK` | Docker network for Traefik | `podoru_traefik` |
| `TRAEFIK_DASHBOARD_PORT` | Traefik dashboard port | `8081` |
| `TRAEFIK_HTTP_PORT` | HTTP entrypoint port | `80` |
| `TRAEFIK_HTTPS_PORT` | HTTPS entrypoint port | `443` |
| `TRAEFIK_ACME_EMAIL` | Email for Let's Encrypt | (required for SSL) |

### Logging

| Variable | Description | Default |
|----------|-------------|---------|
| `LOGGER_LEVEL` | Log level (`debug` / `info` / `warn` / `error`) | `info` |
| `LOGGER_FORMAT` | Output format (`json` / `text`) | `json` |
| `LOGGER_OUTPUT` | Output destination (`stdout` / `file`) | `stdout` |

## Configuration File

The default configuration file is at `configs/config.yaml`:

```yaml
app:
  name: podoru
  env: production
  port: 8080
  debug: false
  registration_enabled: false

database:
  host: localhost
  port: 5432
  user: podoru
  password: ""  # Use environment variable
  name: podoru
  ssl_mode: disable
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

jwt:
  secret: ""  # Use environment variable
  access_expiry: 15m
  refresh_expiry: 168h

encryption:
  key: ""  # Use environment variable

docker:
  host: unix:///var/run/docker.sock

traefik:
  enabled: true
  dashboard_port: 8081
  http_port: 80
  https_port: 443
  acme_email: ""
  network: podoru_traefik

logger:
  level: info
  format: json
  output: stdout
```

## Production Configuration

For production, create a `.env.prod` file:

```bash
# Domain
DOMAIN=podoru.example.com
ACME_EMAIL=admin@example.com

# Database (use strong passwords)
DB_USER=podoru
DB_PASSWORD=super-secure-random-password-here
DB_NAME=podoru

# Security (generate with: openssl rand -base64 32)
JWT_SECRET=your-32-character-jwt-secret-key
ENCRYPTION_KEY=your-32-character-encrypt-key!

# Traefik
TRAEFIK_ENABLED=true
TRAEFIK_NETWORK=podoru_traefik

# Traefik Dashboard Auth (generate with: htpasswd -nb admin password)
TRAEFIK_DASHBOARD_AUTH=admin:$apr1$xyz...

# Disable registration after first user
REGISTRATION_ENABLED=false
```

## Generating Secrets

### JWT Secret

```bash
openssl rand -base64 32
```

### Encryption Key

Must be exactly 32 characters:

```bash
openssl rand -base64 24 | head -c 32
```

### Traefik Dashboard Password

```bash
# Install htpasswd
apt install apache2-utils

# Generate password hash
htpasswd -nb admin your-password
```

## Development Configuration

For development, use `.env`:

```bash
APP_ENV=development
APP_DEBUG=true
REGISTRATION_ENABLED=true

DB_HOST=localhost
DB_PORT=5432
DB_USER=podoru
DB_PASSWORD=podoru_secret
DB_NAME=podoru

JWT_SECRET=dev-jwt-secret-change-in-production
ENCRYPTION_KEY=dev-32-byte-key-for-aes-256!!!

TRAEFIK_ENABLED=true
TRAEFIK_NETWORK=podoru_traefik
```

## Validating Configuration

Check your configuration:

```bash
# Test database connection
docker compose exec postgres pg_isready -U podoru -d podoru

# Test API health
curl http://localhost:8080/health

# View current config (development only)
curl http://localhost:8080/api/v1/debug/config  # if debug enabled
```
