# Architecture

Podoru follows Clean Architecture principles with clear separation of concerns.

## Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP Layer                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │    Handlers     │  │   Middleware    │  │    DTOs     │ │
│  └────────┬────────┘  └────────┬────────┘  └─────────────┘ │
└───────────┼─────────────────────┼───────────────────────────┘
            │                     │
            ▼                     ▼
┌─────────────────────────────────────────────────────────────┐
│                      Use Case Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Auth UseCase  │  │ Service UseCase │  │  Deployment │ │
│  └────────┬────────┘  └────────┬────────┘  └──────┬──────┘ │
└───────────┼─────────────────────┼─────────────────┼─────────┘
            │                     │                 │
            ▼                     ▼                 ▼
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │    Entities     │  │  Repositories   │  │  Interfaces │ │
│  └─────────────────┘  └────────┬────────┘  └─────────────┘ │
└────────────────────────────────┼────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────┐
│                   Infrastructure Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   PostgreSQL    │  │  Docker Client  │  │   Config    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
podoru/
├── cmd/
│   └── podoru/
│       └── main.go              # Application entrypoint
│
├── configs/
│   └── config.yaml              # Default configuration
│
├── gitbook/                     # Documentation
│
├── internal/
│   ├── adapter/
│   │   ├── http/
│   │   │   ├── handler/         # HTTP request handlers
│   │   │   ├── middleware/      # Auth, logging, etc.
│   │   │   ├── dto/             # Request/Response structs
│   │   │   └── router/          # Route definitions
│   │   └── repository/          # Database implementations
│   │
│   ├── domain/
│   │   ├── entity/              # Business entities
│   │   ├── repository/          # Repository interfaces
│   │   └── docker/              # Docker interfaces
│   │
│   ├── infrastructure/
│   │   ├── config/              # Config loading
│   │   ├── database/            # DB connection, migrations
│   │   ├── docker/              # Docker client
│   │   └── logger/              # Logging setup
│   │
│   └── usecase/
│       ├── auth/                # Authentication logic
│       ├── user/                # User management
│       ├── team/                # Team management
│       ├── project/             # Project management
│       ├── service/             # Service management
│       └── deployment/          # Container deployment
│
├── migrations/                  # SQL migrations
│
└── scripts/
    ├── install.sh              # Production installer
    └── dev-setup.sh            # Development setup
```

## Layers

### HTTP Layer (`internal/adapter/http/`)

Handles HTTP concerns:

- **Handlers**: Parse requests, call use cases, format responses
- **Middleware**: Authentication, logging, error recovery
- **DTOs**: Request/response serialization
- **Router**: Route registration with Gin

### Use Case Layer (`internal/usecase/`)

Contains business logic:

- Orchestrates operations
- Enforces business rules
- Coordinates between repositories
- No knowledge of HTTP or database specifics

### Domain Layer (`internal/domain/`)

Core business definitions:

- **Entities**: Business objects (User, Team, Service, etc.)
- **Repository Interfaces**: Data access contracts
- **Docker Interfaces**: Container management contracts

### Infrastructure Layer (`internal/infrastructure/`)

External implementations:

- **Database**: PostgreSQL connection, migrations
- **Docker**: Docker SDK client
- **Config**: Viper-based configuration

## Key Components

### Entities

```go
// internal/domain/entity/service.go
type Service struct {
    ID            uuid.UUID
    ProjectID     uuid.UUID
    Name          string
    Slug          string
    DeployType    DeployType
    Image         *string
    Status        ServiceStatus
    ContainerID   *string
    // ...
}
```

### Repository Interface

```go
// internal/domain/repository/service_repository.go
type ServiceRepository interface {
    Create(ctx context.Context, service *entity.Service) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error)
    Update(ctx context.Context, service *entity.Service) error
    Delete(ctx context.Context, id uuid.UUID) error
    // ...
}
```

### Use Case

```go
// internal/usecase/deployment/deployment_usecase.go
type UseCase struct {
    serviceRepo      repository.ServiceRepository
    containerManager docker.ContainerManager
    // ...
}

func (uc *UseCase) Deploy(ctx context.Context, userID, serviceID uuid.UUID) (*entity.Deployment, error) {
    // Business logic here
}
```

### Handler

```go
// internal/adapter/http/handler/service_handler.go
func (h *ServiceHandler) Deploy(c *gin.Context) {
    // Parse request
    serviceID := c.Param("serviceId")

    // Call use case
    deployment, err := h.deploymentUC.Deploy(ctx, userID, serviceID)

    // Format response
    c.JSON(http.StatusOK, dto.SuccessResponse(deployment))
}
```

## Dependency Injection

Dependencies are wired in `main.go`:

```go
// Infrastructure
db := database.NewConnection(cfg.Database)
dockerClient := docker.NewClient(cfg.Docker)

// Repositories
serviceRepo := repository.NewServiceRepository(db)
domainRepo := repository.NewDomainRepository(db)

// Use Cases
deploymentUC := deployment.NewUseCase(
    serviceRepo,
    domainRepo,
    dockerClient,
    cfg.Traefik,
)

// Handlers
serviceHandler := handler.NewServiceHandler(serviceUC, deploymentUC)

// Router
router := gin.New()
router.POST("/services/:id/deploy", serviceHandler.Deploy)
```

## Data Flow

1. HTTP request arrives at handler
2. Handler validates and parses request
3. Handler calls use case with domain objects
4. Use case executes business logic
5. Use case calls repositories/services as needed
6. Response flows back up through layers

## Docker Integration

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Deployment UC  │────▶│ ContainerManager│────▶│  Docker Engine  │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                               │
                               ▼
                        ┌─────────────────┐
                        │ Traefik Network │
                        └─────────────────┘
```

The ContainerManager interface abstracts Docker operations:

```go
type ContainerManager interface {
    PullImage(ctx context.Context, image string) error
    CreateContainer(ctx context.Context, config *ContainerConfig) (string, error)
    StartContainer(ctx context.Context, containerID string) error
    StopContainer(ctx context.Context, containerID string, timeout *int) error
    ValidateNetwork(ctx context.Context, networkName string) error
    // ...
}
```
