# Podoru

A Docker container management platform similar to EasyPanel. Manage containers, services, networks, and Docker Swarm clusters through a REST API.

## Features

- **Authentication**: Email/Password with JWT tokens (access + refresh)
- **Multi-tenancy**: Teams/Organizations with role-based access (owner, admin, member)
- **Projects**: Organize services within teams
- **Services**: Deploy via Docker image, Dockerfile, or docker-compose
- **GitHub Integration**: Auto-deploy on push with webhook support
- **Traefik Gateway**: Automatic SSL and routing management
- **Docker Swarm**: Cluster management and service scaling

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin
- **Database**: PostgreSQL 16
- **Gateway**: Traefik v3
- **Container Runtime**: Docker Engine + Docker Swarm

## Prerequisites

- Go 1.22+
- PostgreSQL 16+
- Docker Engine
- Make

## Quick Start

### 1. Clone the repository

```bash
git clone git@github.com:podoru/spinner-podoru.git
cd spinner-podoru
```

### 2. Setup environment

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your configuration
vim .env
```

### 3. Start PostgreSQL (using Docker)

```bash
docker compose up -d postgres
```

### 4. Run the application

```bash
# Install dependencies
make deps

# Run database migrations and start server
make run
```

The API will be available at `http://localhost:8080`

## Development

### Available Commands

```bash
make help           # Show all available commands
make build          # Build the application
make run            # Build and run
make dev            # Run with hot reload (requires air)
make test           # Run tests
make test-coverage  # Run tests with coverage report
make lint           # Run linter
make fmt            # Format code
make docs           # Generate API documentation
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
```

### Running with Hot Reload

```bash
# Install air
go install github.com/air-verse/air@latest

# Run with hot reload
make dev
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

## API Documentation

Interactive API documentation is available via Scalar UI:

```
http://localhost:8080/api/v1/docs
```

OpenAPI JSON spec:
```
http://localhost:8080/api/v1/docs/openapi.json
```

### Generate/Update Documentation

```bash
# Install swag (if not installed)
go install github.com/swaggo/swag/cmd/swag@v1.16.3

# Generate docs
make docs
```

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user (first user becomes superadmin) |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/auth/logout` | Logout |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/me` | Get current user |
| PUT | `/api/v1/users/me` | Update current user |
| PUT | `/api/v1/users/me/password` | Change password |

### Teams
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/teams` | List user's teams |
| POST | `/api/v1/teams` | Create team |
| GET | `/api/v1/teams/:teamId` | Get team |
| PUT | `/api/v1/teams/:teamId` | Update team |
| DELETE | `/api/v1/teams/:teamId` | Delete team |
| GET | `/api/v1/teams/:teamId/members` | List members |
| POST | `/api/v1/teams/:teamId/members` | Add member |
| PUT | `/api/v1/teams/:teamId/members/:userId` | Update member role |
| DELETE | `/api/v1/teams/:teamId/members/:userId` | Remove member |

### Projects
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/teams/:teamId/projects` | List projects |
| POST | `/api/v1/teams/:teamId/projects` | Create project |
| GET | `/api/v1/projects/:projectId` | Get project |
| PUT | `/api/v1/projects/:projectId` | Update project |
| DELETE | `/api/v1/projects/:projectId` | Delete project |

### Services
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/projects/:projectId/services` | List services |
| POST | `/api/v1/projects/:projectId/services` | Create service |
| GET | `/api/v1/services/:serviceId` | Get service |
| PUT | `/api/v1/services/:serviceId` | Update service |
| DELETE | `/api/v1/services/:serviceId` | Delete service |
| POST | `/api/v1/services/:serviceId/deploy` | Deploy service |
| POST | `/api/v1/services/:serviceId/start` | Start service |
| POST | `/api/v1/services/:serviceId/stop` | Stop service |
| POST | `/api/v1/services/:serviceId/restart` | Restart service |
| POST | `/api/v1/services/:serviceId/scale` | Scale replicas |
| GET | `/api/v1/services/:serviceId/logs` | Get logs |

## Configuration

Configuration is loaded from `configs/config.yaml` or environment variables.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CONFIG_PATH` | Path to config file | `configs/config.yaml` |
| `APP_PORT` | Server port | `8080` |
| `APP_ENV` | Environment (development/production) | `development` |
| `REGISTRATION_ENABLED` | Allow new registrations | `false` |
| `DATABASE_HOST` | PostgreSQL host | `localhost` |
| `DATABASE_PORT` | PostgreSQL port | `5432` |
| `DATABASE_USER` | PostgreSQL user | `podoru` |
| `DATABASE_PASSWORD` | PostgreSQL password | - |
| `DATABASE_NAME` | PostgreSQL database | `podoru` |
| `JWT_SECRET` | JWT signing secret | - |
| `ENCRYPTION_KEY` | AES encryption key (32 bytes) | - |

## Project Structure

```
podoru/
├── cmd/podoru/          # Application entrypoint
├── configs/             # Configuration files
├── docs/                # Generated API documentation
├── internal/
│   ├── adapter/
│   │   ├── http/        # HTTP handlers, middleware, DTOs
│   │   ├── repository/  # Database implementations
│   │   └── gateway/     # External services (Docker, GitHub, Traefik)
│   ├── domain/
│   │   ├── entity/      # Business entities
│   │   └── repository/  # Repository interfaces
│   ├── infrastructure/  # Config, database, logger
│   ├── mocks/           # Test mocks
│   └── usecase/         # Business logic
├── migrations/          # Database migrations
└── pkg/                 # Shared utilities
```

## First User Setup

The first registered user automatically becomes a **superadmin**. After the first user registers:

1. Registration can be disabled via `REGISTRATION_ENABLED=false`
2. New users can only be added by existing admins through team invitations

## License

MIT License - see [LICENSE](LICENSE) for details.
