# Podoru

A self-hosted Docker container management platform with automatic SSL and domain routing. Deploy and manage containers through a simple REST API with Traefik integration.

## Features

- **Container Management**: Deploy containers from Docker images with full lifecycle control
- **Automatic SSL**: Let's Encrypt integration via Traefik for automatic HTTPS
- **Domain Routing**: Bind custom domains to services with automatic Traefik configuration
- **Multi-tenancy**: Teams and projects with role-based access control
- **Authentication**: JWT-based auth with access and refresh tokens
- **GitHub Integration**: Webhook support for auto-deploy on push
- **Docker Swarm**: Optional cluster management and service scaling

## Quick Start

### Development Setup

```bash
# Clone and setup
git clone git@github.com:podoru/spinner-podoru.git
cd spinner-podoru

# Run development setup (installs deps, starts database)
./scripts/dev-setup.sh setup

# Start with hot reload
make dev
```

### Production Installation

```bash
# Download and run install script
curl -fsSL https://raw.githubusercontent.com/podoru/spinner-podoru/master/scripts/install.sh -o install.sh
chmod +x install.sh
./install.sh install
```

See the [Installation Guide](gitbook/getting-started/installation.md) for detailed instructions.

## Documentation

Full documentation is available in the [GitBook](gitbook/README.md):

- [Getting Started](gitbook/getting-started/README.md)
- [Configuration](gitbook/getting-started/configuration.md)
- [Deployment Guide](gitbook/guides/deployment.md)
- [Domain Management](gitbook/guides/domains.md)
- [API Reference](gitbook/api/README.md)

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.22+ |
| Framework | Gin |
| Database | PostgreSQL 16 |
| Reverse Proxy | Traefik v3 |
| Container Runtime | Docker Engine |

## Development Commands

```bash
make help           # Show all available commands
make dev            # Run with hot reload
make dev-full       # Run with hot reload + Traefik
make run            # Build and run
make test           # Run tests
make docs           # Generate API documentation
make docker-up      # Start all Docker services
make docker-down    # Stop Docker services
```

## API Overview

| Endpoint | Description |
|----------|-------------|
| `POST /api/v1/auth/register` | Register new user |
| `POST /api/v1/auth/login` | Login and get tokens |
| `GET /api/v1/teams` | List user's teams |
| `POST /api/v1/teams/:id/projects` | Create project |
| `POST /api/v1/projects/:id/services` | Create service |
| `POST /api/v1/services/:id/deploy` | Deploy service |
| `POST /api/v1/services/:id/domains` | Add domain |

Interactive API docs: `http://localhost:8080/api/v1/docs`

## Architecture

```
podoru/
├── cmd/podoru/              # Application entrypoint
├── configs/                 # Configuration files
├── gitbook/                 # Documentation
├── internal/
│   ├── adapter/http/        # HTTP handlers and middleware
│   ├── domain/              # Business entities and interfaces
│   ├── infrastructure/      # Config, database, Docker client
│   └── usecase/             # Business logic
├── migrations/              # Database migrations
└── scripts/                 # Installation and setup scripts
```

## License

MIT License - see [LICENSE](LICENSE) for details.
