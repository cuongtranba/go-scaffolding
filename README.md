# Go Clean Architecture Template

A production-ready Go template repository implementing **Hexagonal Architecture (Ports & Adapters)** with feature-sliced structure. This template provides a solid foundation for building scalable REST APIs, gRPC services, CLI applications, and background workers.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## Features

### Architecture & Design
- ✅ **Hexagonal Architecture** - Clean separation between domain logic and infrastructure
- ✅ **Feature-Sliced Structure** - Vertical organization by domain feature
- ✅ **Dependency Injection** - Compile-time DI with Google Wire
- ✅ **SOLID Principles** - Maintainable and testable code

### Multi-Protocol Support
- ✅ **REST API** - HTTP/JSON API with Gin framework
- 🚧 **gRPC** - High-performance RPC (planned)
- 🚧 **CLI** - Command-line interface with Cobra (planned)
- 🚧 **Workers** - Background job processing (planned)

### Database Support
- ✅ **PostgreSQL** - Primary database with GORM v2
- 🚧 **MongoDB** - Document store (planned)
- 🚧 **Redis** - Caching and pub/sub (planned)

### Observability
- ✅ **Structured Logging** - JSON logging with zerolog
- ✅ **Health Checks** - Kubernetes-ready liveness/readiness endpoints
- 🚧 **Metrics** - Prometheus metrics (planned)
- 🚧 **Tracing** - OpenTelemetry distributed tracing (planned)

### Developer Experience
- ✅ **Taskfile** - Simple task runner for common operations
- ✅ **Docker Compose** - Local development environment
- ✅ **Hot Reload** - Fast development iteration
- ✅ **Database Migrations** - Version-controlled schema changes
- ✅ **Code Generation** - Mocks and Wire providers

### Testing
- ✅ **Unit Tests** - 83.8% average coverage
- ✅ **Integration Tests** - Testcontainers with real PostgreSQL
- ✅ **Mocking** - Type-safe mocks with Mockery v3
- 🚧 **E2E Tests** - Full system tests (planned)

## Architecture

This template implements **Hexagonal Architecture (Ports & Adapters)** with three main layers:

```
┌─────────────────────────────────────────────────────────┐
│                     External World                       │
│  (HTTP, gRPC, CLI, Databases, Message Queues, etc.)     │
└────────────────────┬────────────────────────────────────┘
                     │
         ┌───────────▼──────────┐
         │   Adapters (Input)   │
         │  - HTTP Handlers     │
         │  - gRPC Servers      │
         │  - CLI Commands      │
         └───────────┬──────────┘
                     │
         ┌───────────▼──────────┐
         │    Ports (Input)     │
         │  - Service Interface │
         └───────────┬──────────┘
                     │
    ┌────────────────▼────────────────┐
    │        Domain Core              │
    │  - Business Logic               │
    │  - Domain Entities              │
    │  - Domain Errors                │
    │  - Validation Rules             │
    └────────────────┬────────────────┘
                     │
         ┌───────────▼──────────┐
         │   Ports (Output)     │
         │  - Repository Iface  │
         └───────────┬──────────┘
                     │
         ┌───────────▼──────────┐
         │  Adapters (Output)   │
         │  - PostgreSQL        │
         │  - MongoDB           │
         │  - Redis             │
         └──────────────────────┘
```

### Key Principles

1. **Dependencies point inward** - Domain has zero external dependencies
2. **Interfaces at boundaries** - Ports define contracts, adapters implement them
3. **Testability** - Each layer can be tested in isolation
4. **Flexibility** - Easy to swap implementations (e.g., switch databases)

## Project Structure

```
.
├── cmd/                          # Application entry points
│   └── api/                      # REST API server
│       ├── main.go              # Server bootstrap
│       ├── wire.go              # Wire injector definition
│       └── integration_test.go  # Integration tests
├── internal/                     # Private application code
│   ├── config/                  # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── infrastructure/          # Infrastructure concerns
│   │   ├── database/           # Database connections
│   │   │   └── postgres.go
│   │   ├── health/             # Health check system
│   │   │   ├── health.go
│   │   │   └── health_test.go
│   │   └── logger/             # Logging infrastructure
│   │       ├── logger.go
│   │       └── logger_test.go
│   ├── user/                    # User feature (example domain)
│   │   ├── domain/             # Domain layer (business logic)
│   │   │   ├── user.go         # User entity with validation
│   │   │   ├── user_test.go
│   │   │   └── errors.go       # Domain-specific errors
│   │   ├── ports/              # Ports (interfaces)
│   │   │   ├── repository.go  # Repository interface
│   │   │   ├── service.go     # Service interface
│   │   │   └── mocks/         # Generated mocks
│   │   ├── service/            # Domain service implementation
│   │   │   ├── user_service.go
│   │   │   └── user_service_test.go
│   │   └── adapters/           # Infrastructure adapters
│   │       ├── postgres/       # PostgreSQL adapter
│   │       │   ├── models.go  # GORM models
│   │       │   ├── mappers.go # Domain ↔ DB mappers
│   │       │   ├── repository.go
│   │       │   └── repository_test.go
│   │       └── http/           # HTTP adapter
│   │           ├── dto.go     # Request/Response DTOs
│   │           ├── handlers.go # HTTP handlers
│   │           └── routes.go  # Route registration
│   └── wire/                    # Wire providers
│       └── providers.go
├── migrations/                   # Database migrations
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── docs/                        # Documentation
│   └── plans/                  # Design and implementation plans
├── config.yaml                  # Application configuration
├── docker-compose.yml          # Local infrastructure
├── Taskfile.yml               # Task automation
├── .gitignore
├── .mockery.yml              # Mockery configuration
├── go.mod
└── README.md
```

### Feature Organization

Each feature (e.g., `user`, `post`, `product`) follows this structure:

```
internal/
└── {feature}/
    ├── domain/        # Business logic, entities, validation
    ├── ports/         # Interface definitions
    ├── service/       # Service implementation (uses repository port)
    └── adapters/      # Infrastructure implementations
        ├── postgres/  # Database adapter
        ├── http/      # HTTP adapter
        ├── grpc/      # gRPC adapter (if needed)
        └── mongodb/   # Alternative database (if needed)
```

## Getting Started

### Prerequisites

- **Go 1.25+** - [Download](https://go.dev/dl/) (Latest: 1.25.4)
- **Docker & Docker Compose** - [Download](https://docs.docker.com/get-docker/)
- **Task** (optional) - [Install](https://taskfile.dev/installation/)
- **golang-migrate** - [Install](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Installation

1. **Clone the template**

```bash
git clone https://github.com/yourusername/go-scaffolding.git my-project
cd my-project
```

2. **Update module name**

```bash
# Replace 'github.com/yourusername/go-scaffolding' with your module name
find . -type f -name '*.go' -exec sed -i 's|github.com/yourusername/go-scaffolding|github.com/yourorg/yourproject|g' {} +
go mod edit -module github.com/yourorg/yourproject
go mod tidy
```

3. **Start infrastructure**

```bash
# Using Task (recommended)
task docker:up

# Or using docker compose directly
docker compose up -d postgres
```

4. **Run database migrations**

```bash
task migrate:up

# Or manually
migrate -path ./migrations -database "postgresql://postgres:postgres@localhost:5432/app?sslmode=disable" up
```

5. **Run the API server**

```bash
# Using Task
task run:api

# Or build and run
task build:api
./bin/api

# Or run directly
go run ./cmd/api
```

6. **Verify it's running**

```bash
# Check health
curl http://localhost:8080/health/live

# Create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test User"}'

# List users
curl http://localhost:8080/users
```

## API Documentation

### Health Endpoints

#### GET /health/live
**Liveness check** - Always returns 200 if the server is running

```bash
curl http://localhost:8080/health/live
```

Response:
```json
{
  "status": "healthy",
  "checks": {
    "liveness": {"status": "healthy"}
  }
}
```

#### GET /health/ready
**Readiness check** - Returns 200 if ready to serve traffic, 503 if not ready

```bash
curl http://localhost:8080/health/ready
```

Response:
```json
{
  "status": "healthy",
  "checks": {
    "database": {"status": "healthy"}
  }
}
```

### User Endpoints

#### POST /users
Create a new user

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe"
  }'
```

Response (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "name": "John Doe",
  "created_at": "2025-11-22T10:00:00Z",
  "updated_at": "2025-11-22T10:00:00Z"
}
```

Errors:
- `400 Bad Request` - Invalid email format or missing required fields
- `409 Conflict` - Email already exists

#### GET /users/:id
Get a user by ID

```bash
curl http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000
```

Response (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "name": "John Doe",
  "created_at": "2025-11-22T10:00:00Z",
  "updated_at": "2025-11-22T10:00:00Z"
}
```

Errors:
- `404 Not Found` - User not found

#### GET /users/email/:email
Get a user by email address

```bash
curl http://localhost:8080/users/email/john@example.com
```

Response (200 OK): Same as GET /users/:id

#### GET /users?limit=10&offset=0
List users with pagination

```bash
curl "http://localhost:8080/users?limit=10&offset=0"
```

Response (200 OK):
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john@example.com",
      "name": "John Doe",
      "created_at": "2025-11-22T10:00:00Z",
      "updated_at": "2025-11-22T10:00:00Z"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

Query Parameters:
- `limit` - Number of users to return (default: 10, max: 100)
- `offset` - Number of users to skip (default: 0)

Results are ordered by `created_at DESC` (newest first).

Errors:
- `400 Bad Request` - Limit exceeds 100

#### PUT /users/:id
Update a user's name

```bash
curl -X PUT http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated"}'
```

Response (200 OK):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "name": "John Updated",
  "created_at": "2025-11-22T10:00:00Z",
  "updated_at": "2025-11-22T10:05:00Z"
}
```

Note: Email cannot be updated for data integrity reasons.

Errors:
- `400 Bad Request` - Invalid name
- `404 Not Found` - User not found

#### DELETE /users/:id
Delete a user (soft delete)

```bash
curl -X DELETE http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000
```

Response (204 No Content): Empty body

Errors:
- `404 Not Found` - User not found

## Development

### Available Tasks

View all available tasks:

```bash
task --list
```

Common tasks:

```bash
# Development
task run:api              # Run API server
task docker:up           # Start infrastructure
task docker:down         # Stop infrastructure

# Testing
task test                # Run unit tests
task test:integration    # Run integration tests (requires Docker)
task test:coverage       # Generate coverage report

# Database
task migrate:up          # Run migrations
task migrate:down        # Rollback last migration
task migrate:create      # Create new migration

# Code Generation
task mock:generate       # Generate mocks with Mockery
task wire:generate       # Generate Wire dependency injection code

# Build
task build:api           # Build API binary
```

### Adding a New Feature

Example: Adding a `Post` feature

1. **Create domain layer**

```bash
mkdir -p internal/post/domain
```

Create `internal/post/domain/post.go`:
```go
package domain

import (
    "errors"
    "time"
)

type Post struct {
    ID        string
    Title     string
    Content   string
    AuthorID  string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewPost(title, content, authorID string) (*Post, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    // Add validation...

    return &Post{
        Title:     title,
        Content:   content,
        AuthorID:  authorID,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}
```

2. **Define ports**

Create `internal/post/ports/repository.go`:
```go
package ports

import (
    "context"
    "github.com/yourorg/yourproject/internal/post/domain"
)

type PostRepository interface {
    Create(ctx context.Context, post *domain.Post) error
    GetByID(ctx context.Context, id string) (*domain.Post, error)
    // Add other methods...
}
```

3. **Implement service**

Create `internal/post/service/post_service.go`:
```go
package service

import (
    "context"
    "github.com/yourorg/yourproject/internal/post/domain"
    "github.com/yourorg/yourproject/internal/post/ports"
)

type PostService struct {
    repo ports.PostRepository
}

func NewPostService(repo ports.PostRepository) *PostService {
    return &PostService{repo: repo}
}

func (s *PostService) CreatePost(ctx context.Context, title, content, authorID string) (*domain.Post, error) {
    post, err := domain.NewPost(title, content, authorID)
    if err != nil {
        return nil, err
    }

    if err := s.repo.Create(ctx, post); err != nil {
        return nil, err
    }

    return post, nil
}
```

4. **Implement adapters**
   - Create PostgreSQL adapter in `internal/post/adapters/postgres/`
   - Create HTTP adapter in `internal/post/adapters/http/`

5. **Wire it up**

Update `internal/wire/providers.go` to include Post providers.

6. **Write tests**
   - Unit tests for domain and service
   - Integration tests for adapters

### Configuration

Configuration is managed via `config.yaml` and environment variables.

Environment variables override config file values using the pattern: `SECTION_KEY`

Examples:
```bash
# Override database host
export POSTGRES_HOST=localhost

# Override log level
export APP_LOG_LEVEL=debug

# Override HTTP port
export APP_HTTP_PORT=3000
```

## Testing

### Unit Tests

Run all unit tests:

```bash
task test
```

Run tests for a specific package:

```bash
go test ./internal/user/domain/... -v
```

Run tests with coverage:

```bash
task test:coverage
# Open coverage.html in browser
```

### Integration Tests

Integration tests use testcontainers to spin up real PostgreSQL containers.

```bash
# Requires Docker
task test:integration

# Or directly
go test ./cmd/api/... -run TestIntegration -v
```

### Test Coverage

Current coverage: **83.8%**

| Package | Coverage |
|---------|----------|
| config | 90% |
| logger | 69% |
| health | 100% |
| user/domain | 97% |
| user/service | 62% |
| user/adapters/postgres | 81% |

### Writing Tests

Example unit test with mocks:

```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.MockUserRepository)
    service := service.NewUserService(mockRepo)

    mockRepo.On("GetByEmail", mock.Anything, "test@example.com").
        Return(nil, domain.ErrUserNotFound)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
        Return(nil)

    // Act
    user, err := service.CreateUser(context.Background(), "test@example.com", "Test User")

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
    mockRepo.AssertExpectations(t)
}
```

## Deployment

### Docker

Build the Docker image:

```bash
docker build -t myapp:latest .
```

Run the container:

```bash
docker run -p 8080:8080 \
  -e POSTGRES_HOST=postgres \
  -e POSTGRES_PORT=5432 \
  myapp:latest
```

### Kubernetes

Example Kubernetes manifests are in `deployments/k8s/` (to be added).

Health check configuration:
```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

## Contributing

Contributions are welcome! Please read the [contributing guidelines](CONTRIBUTING.md) first.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`task test:all`)
5. Commit your changes (`git commit -m 'feat: add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test additions or changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) by Alistair Cockburn
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin
- [Feature-Sliced Design](https://feature-sliced.design/)

---

**Made with ❤️ by developers, for developers**
