# Go Clean Architecture Template - Design Document

**Date**: 2025-11-22
**Status**: Approved
**Architecture**: Hexagonal/Ports & Adapters with Feature-Sliced Structure

## Overview

This is a template repository for building Go applications following clean architecture principles with hexagonal/ports & adapters pattern. The template supports multiple application types (REST API, gRPC, CLI, background workers) and provides a production-ready foundation with observability, testing, and developer experience built-in.

## Key Decisions

### Repository Type
- **Template repository** (not a CLI tool)
- Developers clone/fork and customize for their projects
- Single Go module with clear directory structure

### Branch Strategy
- **Main branch** (`minimal`): Clean boilerplate with structure, interfaces, minimal implementation
  - Includes all adapters (HTTP, gRPC, CLI, worker) as skeleton code
  - Includes all infrastructure (Postgres, Mongo, Redis) as skeleton implementations
  - Developers delete unused components after cloning
- **Example branch** (`with-todo-app`): Fully functional TODO application
  - Demonstrates all patterns in action
  - Shows domain entities, ports, adapters, use cases, Wire integration
  - Reference implementation for learning

### Supported Application Types
- REST API services (HTTP/HTTPS)
- gRPC microservices
- CLI applications
- Background workers (job processors, scheduled tasks)

### Technology Stack

**Core Architecture**:
- Hexagonal/Ports & Adapters pattern
- Feature-sliced vertical organization
- Dependency injection via Google Wire

**Storage**:
- PostgreSQL with GORM v2
- MongoDB with official mongo-driver
- Redis with go-redis

**Configuration**:
- Viper for config management
- Environment variables (12-factor app)
- Multiple profiles (dev, staging, production)

**Observability**:
- Structured logging: zerolog
- Metrics: Prometheus
- Distributed tracing: OpenTelemetry
- Health checks for Kubernetes

**Testing**:
- Unit tests with mockery for mocking
- Integration tests with testcontainers-go
- E2E tests with full stack

**Developer Experience**:
- Taskfile.yml for common tasks
- Docker & docker-compose for local development
- Comprehensive README and architecture docs

## Architecture

### Hexagonal/Ports & Adapters Principles

**Dependency Rule**: Dependencies point inward toward the domain.

```
┌─────────────────────────────────────────┐
│         Adapters (Outer Layer)          │
│  HTTP │ gRPC │ CLI │ Postgres │ Redis   │
├─────────────────────────────────────────┤
│           Ports (Interfaces)            │
│  Services │ Repositories │ External APIs│
├─────────────────────────────────────────┤
│        Domain (Core Business)           │
│    Entities │ Value Objects │ Rules     │
└─────────────────────────────────────────┘
```

**Domain**: Pure business logic, zero external dependencies
**Ports**: Interface contracts (primary inbound, secondary outbound)
**Adapters**: Concrete implementations of ports
**Services**: Application logic orchestrating domain and ports

### Directory Structure

```
go-scaffolding/
├── cmd/                          # Application entry points
│   ├── api/                      # HTTP REST API server
│   │   ├── main.go
│   │   └── wire.go              # Wire dependency graph
│   ├── grpc-server/             # gRPC service server
│   │   ├── main.go
│   │   └── wire.go
│   ├── cli/                     # CLI application
│   │   ├── main.go
│   │   └── wire.go
│   └── worker/                  # Background worker
│       ├── main.go
│       └── wire.go
│
├── internal/                     # Private application code
│   ├── {feature}/               # Feature slice (e.g., user, todo)
│   │   ├── domain/              # Domain entities and logic
│   │   │   ├── entity.go
│   │   │   ├── errors.go
│   │   │   └── value_objects.go
│   │   ├── ports/               # Interface definitions
│   │   │   ├── service.go       # Primary ports (inbound)
│   │   │   └── repository.go    # Secondary ports (outbound)
│   │   ├── service/             # Use case implementations
│   │   │   └── service.go
│   │   └── adapters/            # Port implementations
│   │       ├── http/            # HTTP handlers
│   │       │   ├── handlers.go
│   │       │   ├── dto.go
│   │       │   └── routes.go
│   │       ├── grpc/            # gRPC service impl
│   │       │   ├── server.go
│   │       │   └── pb/          # Generated protobuf
│   │       ├── cli/             # CLI commands
│   │       │   └── commands.go
│   │       ├── postgres/        # PostgreSQL repository
│   │       │   ├── repository.go
│   │       │   ├── models.go    # GORM models
│   │       │   └── mappers.go   # Model <-> Entity
│   │       ├── mongo/           # MongoDB repository
│   │       │   └── repository.go
│   │       └── redis/           # Redis cache
│   │           └── cache.go
│   │
│   ├── infrastructure/          # Shared infrastructure
│   │   ├── database/
│   │   │   ├── postgres.go     # GORM connection setup
│   │   │   ├── mongo.go        # MongoDB client setup
│   │   │   └── redis.go        # Redis client setup
│   │   ├── logger/             # Structured logging wrapper
│   │   │   └── logger.go
│   │   ├── metrics/            # Prometheus metrics
│   │   │   └── metrics.go
│   │   ├── tracing/            # OpenTelemetry setup
│   │   │   └── tracing.go
│   │   └── health/             # Health check handlers
│   │       └── health.go
│   │
│   ├── config/                  # Configuration management
│   │   ├── config.go
│   │   └── validator.go
│   │
│   └── wire/                    # Wire providers
│       └── providers.go
│
├── pkg/                         # Public reusable packages
│   └── ...
│
├── api/                         # API definitions
│   └── proto/                   # Protocol buffer definitions
│       └── {feature}/
│           └── service.proto
│
├── migrations/                  # Database migrations
│   └── {feature}/
│       ├── 001_initial.up.sql
│       └── 001_initial.down.sql
│
├── test/                        # Test utilities and E2E tests
│   ├── e2e/                    # End-to-end tests
│   └── helpers/                # Test fixtures and utilities
│
├── deployments/                 # Deployment configs
│   ├── docker/
│   ├── kubernetes/
│   └── grafana/                # Dashboard templates
│
├── docs/                        # Documentation
│   ├── architecture/           # ADRs and diagrams
│   └── plans/                  # Design documents
│
├── docker-compose.yml          # Local development stack
├── Dockerfile                  # Production container
├── Taskfile.yml               # Task runner commands
├── go.mod
├── go.sum
└── README.md
```

### Feature Slice Structure

Each feature (domain bounded context) is self-contained:

**Example: `/internal/user/`**
```
user/
├── domain/
│   ├── user.go              # User entity (ID, Email, Name, etc.)
│   ├── errors.go            # ErrUserNotFound, ErrInvalidEmail
│   └── email.go             # Email value object
├── ports/
│   ├── service.go           # UserService interface
│   └── repository.go        # UserRepository interface
├── service/
│   └── user_service.go      # Implements UserService
└── adapters/
    ├── http/
    │   ├── handlers.go      # CreateUser, GetUser handlers
    │   ├── dto.go           # CreateUserRequest, UserResponse
    │   └── routes.go        # Register routes
    ├── grpc/
    │   └── server.go        # gRPC UserService implementation
    ├── postgres/
    │   ├── repository.go    # Implements UserRepository
    │   ├── models.go        # GORM UserModel with tags
    │   └── mappers.go       # toUserEntity, fromUserEntity
    └── redis/
        └── cache.go         # User caching layer
```

**Rationale**: Co-locates related code, scales well for large domains, clear feature boundaries.

## Data Flow

### Inbound Request Flow (HTTP Example)

```
1. HTTP Request → Gin Router
2. Router → Middleware (logging, tracing, metrics)
3. Middleware → HTTP Handler (adapter)
4. Handler validates request → creates DTO
5. Handler calls Service interface (port)
6. Service executes business logic using domain entities
7. Service calls Repository interface (port)
8. Repository adapter (Postgres/Mongo) performs data operations
9. Repository returns domain entities
10. Service returns result to handler
11. Handler maps domain entity → response DTO
12. Handler sends HTTP response
```

**Key points**:
- HTTP handler never touches database directly
- Service layer is protocol-agnostic
- Domain entities flow through all layers
- Adapters handle mapping (DTOs ↔ Entities ↔ Models)

### Database Layer with GORM

**Separation of concerns**:
- **Domain entities**: Pure Go structs, no tags, business-focused
- **GORM models**: Database-specific structs with GORM tags
- **Mappers**: Convert between models and entities

**Example**:
```go
// domain/user.go
type User struct {
    ID        string
    Email     Email  // Value object
    Name      string
    CreatedAt time.Time
}

// adapters/postgres/models.go
type UserModel struct {
    ID        string    `gorm:"primaryKey"`
    Email     string    `gorm:"unique;not null"`
    Name      string    `gorm:"not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

// adapters/postgres/mappers.go
func toUserEntity(m *UserModel) (*domain.User, error) {
    email, err := domain.NewEmail(m.Email)
    if err != nil {
        return nil, err
    }
    return &domain.User{
        ID:        m.ID,
        Email:     email,
        Name:      m.Name,
        CreatedAt: m.CreatedAt,
    }, nil
}
```

**Migration strategy**:
- Development: GORM AutoMigrate for rapid iteration
- Production: Manual migrations via golang-migrate for control
- Migration files per feature in `/migrations/{feature}/`

## Dependency Injection with Wire

**Wire philosophy**: Compile-time dependency injection, no runtime reflection.

### Wire Setup per Entry Point

Each `cmd/` has its own `wire.go`:

```go
// cmd/api/wire.go
//go:build wireinject

func InitializeApp(cfg *config.Config) (*app.Application, error) {
    wire.Build(
        // Infrastructure
        infrastructure.NewPostgresDB,
        infrastructure.NewMongoDB,
        infrastructure.NewRedis,
        infrastructure.NewLogger,
        infrastructure.NewTracer,

        // Repositories
        postgres.NewUserRepository,
        mongo.NewProductRepository,
        redis.NewCache,

        // Services
        service.NewUserService,
        service.NewProductService,

        // HTTP handlers
        httpHandlers.NewUserHandler,
        httpHandlers.NewProductHandler,

        // Application
        app.NewHTTPServer,
    )
    return nil, nil
}
```

**Wire generates** `wire_gen.go` with all initialization logic, constructor calls, and error handling.

### Shared Providers

Common providers in `/internal/wire/providers.go`:

```go
func ProvideUserService(repo ports.UserRepository, logger *logger.Logger) ports.UserService {
    return service.NewUserService(repo, logger)
}

func ProvidePostgresUserRepo(db *gorm.DB) ports.UserRepository {
    return postgres.NewUserRepository(db)
}
```

## Configuration Management

**Viper setup** in `/internal/config/`:

```go
type Config struct {
    App         AppConfig
    Postgres    PostgresConfig
    MongoDB     MongoConfig
    Redis       RedisConfig
    Observability ObservabilityConfig
}

type AppConfig struct {
    Name        string `mapstructure:"name"`
    Environment string `mapstructure:"environment"` // dev, staging, prod
    HTTPPort    int    `mapstructure:"http_port"`
    GRPCPort    int    `mapstructure:"grpc_port"`
}

type PostgresConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Database string `mapstructure:"database"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
}
```

**Configuration sources** (priority order):
1. Environment variables (highest priority)
2. Config file (`config.{environment}.yaml`)
3. Default values

**Validation**: Required fields checked on startup, app fails fast with clear errors.

## Testing Strategy

### Unit Tests

**Location**: Alongside code (`*_test.go`)
**Focus**: Domain logic and service layer
**Dependencies**: Mocked using mockery

**Mock generation**:
```go
//go:generate mockery --name=UserRepository --dir=ports --output=ports/mocks
```

**Example**:
```go
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
    service := service.NewUserService(mockRepo, logger)

    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

    err := service.CreateUser(ctx, "test@example.com", "Test User")
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

**Task command**: `task test`

### Integration Tests

**Location**: Adapter directories (`*_integration_test.go`)
**Build tag**: `//go:build integration`
**Dependencies**: Real databases via testcontainers-go

**Example**:
```go
//go:build integration

func TestPostgresUserRepository_Create(t *testing.T) {
    ctx := context.Background()

    // Start PostgreSQL container
    pgContainer, err := postgres.RunContainer(ctx)
    require.NoError(t, err)
    defer pgContainer.Terminate(ctx)

    // Get connection string
    connStr, err := pgContainer.ConnectionString(ctx)
    require.NoError(t, err)

    // Initialize GORM
    db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
    require.NoError(t, err)

    // Auto migrate
    db.AutoMigrate(&models.UserModel{})

    // Test repository
    repo := postgres.NewUserRepository(db)
    user := &domain.User{...}

    err = repo.Create(ctx, user)
    assert.NoError(t, err)
}
```

**Task command**: `task test:integration`

### E2E Tests

**Location**: `/test/e2e/`
**Build tag**: `//go:build e2e`
**Setup**: Docker Compose spins up entire stack
**Scope**: Full workflows through HTTP/gRPC APIs

**Example**:
```go
//go:build e2e

func TestUserWorkflow(t *testing.T) {
    // Assumes docker-compose is running
    client := &http.Client{}

    // Create user
    resp, err := client.Post("http://localhost:8080/api/v1/users", ...)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    // Get user
    resp, err = client.Get("http://localhost:8080/api/v1/users/123")
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

**Task command**: `task test:e2e`

### Test Utilities

**Location**: `/test/helpers/`

- Fixture builders for domain entities
- Database seeders
- Assertion helpers
- Test context setup

## Observability

### Structured Logging (zerolog)

**Setup**: `/internal/infrastructure/logger/`

```go
logger := logger.New(cfg.Observability.LogLevel)
logger.Info().Str("service", "user-api").Msg("Starting application")
```

**Features**:
- Context-aware: trace IDs attached to logger
- JSON output in production
- Pretty console output in development
- HTTP middleware auto-logs requests

### Metrics (Prometheus)

**Setup**: `/internal/infrastructure/metrics/`

**Exposed metrics**:
- HTTP: request count, duration, status codes
- Database: query duration, connection pool stats
- Business: feature-specific counters/gauges

**Endpoint**: `GET /metrics`

### Distributed Tracing (OpenTelemetry)

**Setup**: `/internal/infrastructure/tracing/`

**Instrumentation**:
- Auto: HTTP handlers, gRPC servers
- Manual: Service methods, repository calls
- Context propagation through all layers

**Exporters**: Jaeger, Zipkin, OTLP (configurable)

### Health Checks

**Endpoints**:
- `GET /health/live` - Liveness (app running)
- `GET /health/ready` - Readiness (dependencies healthy)

**Checks**:
- Postgres connection
- MongoDB connection
- Redis connection
- Custom per-adapter checks

## Multi-Protocol Support

### HTTP/REST API

**Framework**: Gin
**Entry point**: `/cmd/api/`
**Handlers**: `/internal/{feature}/adapters/http/`

**Features**:
- RESTful conventions
- Middleware: logging, tracing, metrics, CORS, recovery
- Request validation
- Swagger/OpenAPI via swaggo

**Example route**:
```go
router.POST("/api/v1/users", userHandler.Create)
router.GET("/api/v1/users/:id", userHandler.GetByID)
```

### gRPC Service

**Entry point**: `/cmd/grpc-server/`
**Protos**: `/api/proto/{feature}/`
**Implementation**: `/internal/{feature}/adapters/grpc/`

**Features**:
- Interceptors: logging, tracing, auth
- gRPC-gateway (optional HTTP/JSON proxy)
- Reflection for development tools

### CLI Application

**Framework**: Cobra + Viper
**Entry point**: `/cmd/cli/`
**Commands**: `/internal/{feature}/adapters/cli/`

**Features**:
- Subcommand structure
- Flag parsing
- Output formatting (table, JSON, YAML)
- Interactive prompts

### Background Worker

**Entry point**: `/cmd/worker/`

**Features**:
- Message queue (RabbitMQ/NATS/Kafka)
- Worker pool with graceful shutdown
- Retry with exponential backoff
- Dead letter queue
- Cron scheduler option

## Developer Experience

### Taskfile.yml Commands

```yaml
tasks:
  run:api:      # Start HTTP API server
  run:grpc:     # Start gRPC server
  run:cli:      # Run CLI application
  run:worker:   # Start background worker

  test:         # Run unit tests
  test:integration:  # Run integration tests
  test:e2e:     # Run E2E tests
  test:all:     # Run all tests

  mock:generate:  # Regenerate all mocks

  docker:up:    # Start docker-compose stack
  docker:down:  # Stop docker-compose stack

  migrate:up:   # Run database migrations
  migrate:down: # Rollback migrations

  lint:         # Run linters
  format:       # Format code

  build:api:    # Build API binary
  build:grpc:   # Build gRPC binary
  build:all:    # Build all binaries
```

### Docker Setup

**docker-compose.yml**: Local development stack
- PostgreSQL with initialization
- MongoDB with replica set
- Redis
- Jaeger (optional, for tracing)

**Dockerfile**: Multi-stage production build
- Builder stage with Go compilation
- Runtime stage with minimal image
- Non-root user
- Health check

### Documentation

**README.md**: Getting started, architecture overview, task commands
**docs/architecture/**: ADRs (Architecture Decision Records)
**docs/plans/**: Design documents like this one

## Implementation Plan

### Phase 1: Foundation
1. Initialize Go module
2. Set up directory structure
3. Create Taskfile.yml
4. Set up docker-compose with databases
5. Implement config package with Viper
6. Set up infrastructure: logger, metrics, tracing, health

### Phase 2: Domain & Ports
7. Create example feature (user or todo) domain layer
8. Define ports (service and repository interfaces)
9. Implement service layer

### Phase 3: Adapters
10. PostgreSQL adapter with GORM (models, mappers, repository)
11. MongoDB adapter
12. Redis adapter
13. HTTP adapter (Gin setup, handlers, DTOs)
14. gRPC adapter (protos, generated code, server)
15. CLI adapter (Cobra commands)
16. Worker adapter (message queue setup)

### Phase 4: Dependency Injection
17. Set up Wire providers
18. Create wire.go for each cmd/
19. Generate wire_gen.go
20. Test bootstrapping each entry point

### Phase 5: Testing
21. Set up mockery code generation
22. Write example unit tests
23. Set up testcontainers for integration tests
24. Write example integration tests
25. Set up E2E test structure
26. Write example E2E tests
27. Create test helpers and fixtures

### Phase 6: Documentation & Polish
28. Write comprehensive README
29. Create architecture diagrams
30. Write ADRs for key decisions
31. Add code comments and examples
32. Create Swagger/OpenAPI docs
33. Grafana dashboard templates

### Phase 7: Example Branch
34. Create example branch from main
35. Implement TODO application
36. Add full CRUD operations
37. Demonstrate all patterns
38. Add comprehensive tests
39. Document the example

## Success Criteria

1. Developers can clone the template and have a running application in < 5 minutes
2. All four application types (HTTP, gRPC, CLI, worker) work out of the box
3. All three databases (Postgres, Mongo, Redis) are integrated and tested
4. Tests pass: unit, integration, and E2E
5. Observability works: logs, metrics, traces, health checks
6. Clear documentation explains architecture and how to extend
7. Example branch demonstrates all patterns with working code
8. Docker setup allows local development without installing dependencies
9. Task commands provide smooth developer workflow
10. Code follows Go best practices and clean architecture principles

## Non-Goals

- Not a CLI generator tool (just a template to clone)
- Not opinionated about business domain (generic structure)
- Not a framework (developers maintain full control)
- Not supporting every database (just Postgres, Mongo, Redis)
- Not including auth/authz (developers add their own)
- Not production Kubernetes manifests (just local development)

## Future Enhancements (Out of Scope for V1)

- Authentication/authorization middleware examples
- Multi-tenancy support
- Event sourcing patterns
- CQRS implementation
- Message queue adapters (RabbitMQ, Kafka specifics)
- More example applications (e-commerce, blog, etc.)
- Kubernetes production manifests
- CI/CD pipeline templates
- Performance benchmarking suite
