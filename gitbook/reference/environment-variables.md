# Environment Variables

Complete reference of all environment variables.

## Application

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `APP_NAME` | Application name | `podoru` | No |
| `APP_ENV` | Environment mode | `development` | No |
| `APP_PORT` | HTTP server port | `8080` | No |
| `APP_DEBUG` | Enable debug mode | `false` | No |
| `REGISTRATION_ENABLED` | Allow user registration | `false` | No |

### APP_ENV Values

- `development` - Debug logging, detailed errors
- `production` - Optimized logging, sanitized errors

## Database

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | PostgreSQL host | `localhost` | No |
| `DB_PORT` | PostgreSQL port | `5432` | No |
| `DB_USER` | Database user | `podoru` | No |
| `DB_PASSWORD` | Database password | - | **Yes** |
| `DB_NAME` | Database name | `podoru` | No |
| `DB_SSL_MODE` | SSL mode | `disable` | No |
| `DB_MAX_OPEN_CONNS` | Max open connections | `25` | No |
| `DB_MAX_IDLE_CONNS` | Max idle connections | `5` | No |
| `DB_CONN_MAX_LIFETIME` | Connection lifetime | `5m` | No |

### DB_SSL_MODE Values

- `disable` - No SSL
- `require` - SSL required
- `verify-ca` - Verify CA certificate
- `verify-full` - Verify CA and hostname

## Security

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `JWT_SECRET` | JWT signing key (32+ chars) | - | **Yes** |
| `JWT_ACCESS_EXPIRY` | Access token lifetime | `15m` | No |
| `JWT_REFRESH_EXPIRY` | Refresh token lifetime | `168h` | No |
| `ENCRYPTION_KEY` | AES-256 key (32 chars) | - | **Yes** |

### Generating Secrets

```bash
# JWT Secret (32+ characters)
openssl rand -base64 32

# Encryption Key (exactly 32 characters)
openssl rand -base64 24 | head -c 32
```

## Docker

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DOCKER_HOST` | Docker socket | `unix:///var/run/docker.sock` | No |

## Traefik

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TRAEFIK_ENABLED` | Enable Traefik integration | `true` | No |
| `TRAEFIK_NETWORK` | Docker network for routing | `podoru_traefik` | No |
| `TRAEFIK_DASHBOARD_PORT` | Dashboard port | `8081` | No |
| `TRAEFIK_HTTP_PORT` | HTTP entrypoint | `80` | No |
| `TRAEFIK_HTTPS_PORT` | HTTPS entrypoint | `443` | No |
| `TRAEFIK_ACME_EMAIL` | Let's Encrypt email | - | For SSL |

## Logging

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `LOGGER_LEVEL` | Log level | `info` | No |
| `LOGGER_FORMAT` | Output format | `json` | No |
| `LOGGER_OUTPUT` | Output destination | `stdout` | No |

### LOGGER_LEVEL Values

- `debug` - All messages
- `info` - Info and above
- `warn` - Warnings and errors
- `error` - Errors only

### LOGGER_FORMAT Values

- `json` - JSON structured logs
- `text` - Human-readable text

## Production Environment

Example `.env.prod`:

```bash
# Application
APP_ENV=production
APP_PORT=8080
APP_DEBUG=false
REGISTRATION_ENABLED=false

# Domain
DOMAIN=podoru.example.com
ACME_EMAIL=admin@example.com

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=podoru
DB_PASSWORD=super-secure-password-here
DB_NAME=podoru
DB_SSL_MODE=disable

# Security
JWT_SECRET=your-32-character-jwt-secret-key
ENCRYPTION_KEY=your-32-character-encrypt-key!

# Docker
DOCKER_HOST=unix:///var/run/docker.sock

# Traefik
TRAEFIK_ENABLED=true
TRAEFIK_NETWORK=podoru_traefik
TRAEFIK_DASHBOARD_AUTH=admin:$apr1$...

# Logging
LOGGER_LEVEL=info
LOGGER_FORMAT=json
```

## Development Environment

Example `.env`:

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

LOGGER_LEVEL=debug
LOGGER_FORMAT=text
```
