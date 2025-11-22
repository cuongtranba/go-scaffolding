# Architecture Overview

This document provides a high-level overview of the Go Clean Architecture Template's architecture, design principles, and structure.

## Table of Contents

- [Architecture Style](#architecture-style)
- [Core Principles](#core-principles)
- [Project Organization](#project-organization)
- [Layering Strategy](#layering-strategy)
- [Dependency Management](#dependency-management)
- [Testing Strategy](#testing-strategy)

## Architecture Style

This template implements **Hexagonal Architecture** (also known as **Ports and Adapters**), combined with **Feature-Sliced Design** for vertical organization.

### Why Hexagonal Architecture?

1. **Business Logic Independence** - Domain logic has zero dependencies on infrastructure
2. **Testability** - Each layer can be tested in complete isolation
3. **Flexibility** - Easy to swap implementations (databases, frameworks, protocols)
4. **Maintainability** - Clear boundaries and responsibilities
5. **Evolvability** - Architecture supports growth and change

## Core Principles

### 1. Dependency Inversion

Dependencies always point inward toward the domain core:

```
External World → Adapters → Ports → Domain
```

The domain never depends on infrastructure. Infrastructure adapts to the domain's interfaces.

### 2. Separation of Concerns

Each layer has a specific responsibility:

- **Domain**: Business logic, entities, validation rules
- **Ports**: Interface definitions (contracts)
- **Service**: Use case orchestration
- **Adapters**: Infrastructure implementations (HTTP, DB, etc.)

### 3. Explicit Boundaries

Interfaces (ports) define clear boundaries between layers. No layer can directly depend on another layer's implementation.

### 4. Testability

Every layer can be tested independently:

- Domain: Pure unit tests, no mocks needed
- Service: Unit tests with repository mocks
- Adapters: Integration tests with real infrastructure

## Project Organization

### Feature-Sliced Structure

Features are organized vertically by domain concern:

```
internal/
├── user/       # User management feature
├── post/       # Blog post feature
└── product/    # Product catalog feature
```

Each feature is self-contained with its own:
- Domain logic
- Port definitions
- Service implementation
- Infrastructure adapters

### Standard Feature Structure

```
internal/{feature}/
├── domain/        # Business entities and logic
│   ├── {entity}.go
│   ├── {entity}_test.go
│   └── errors.go
├── ports/         # Interface definitions
│   ├── repository.go
│   └── service.go
├── service/       # Use case implementations
│   ├── {feature}_service.go
│   └── {feature}_service_test.go
└── adapters/      # Infrastructure implementations
    ├── postgres/  # Database adapter
    ├── http/      # HTTP adapter
    ├── grpc/      # gRPC adapter (optional)
    └── mongodb/   # Alternative storage (optional)
```

## Layering Strategy

### Layer 1: Domain Core

**Location**: `internal/{feature}/domain/`

**Responsibilities**:
- Define business entities
- Implement domain validation
- Define domain-specific errors
- Contain business rules

**Dependencies**: None (zero external dependencies)

**Example**:
```go
// internal/user/domain/user.go
type User struct {
    ID    string
    Email string
    Name  string
}

func NewUser(email, name string) (*User, error) {
    // Domain validation
    if !isValidEmail(email) {
        return nil, ErrInvalidEmail
    }
    // Business logic
    return &User{...}, nil
}
```

### Layer 2: Ports (Interfaces)

**Location**: `internal/{feature}/ports/`

**Responsibilities**:
- Define service interfaces (input ports)
- Define repository interfaces (output ports)
- Act as contracts between layers

**Dependencies**: Domain only

**Example**:
```go
// internal/user/ports/repository.go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    GetByID(ctx context.Context, id string) (*domain.User, error)
}
```

### Layer 3: Service (Use Cases)

**Location**: `internal/{feature}/service/`

**Responsibilities**:
- Orchestrate use cases
- Coordinate between domain and repositories
- Handle business workflows

**Dependencies**: Domain and ports (interfaces only)

**Example**:
```go
// internal/user/service/user_service.go
type UserService struct {
    repo ports.UserRepository
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
    // Orchestration logic
    user, err := domain.NewUser(email, name)
    return user, s.repo.Create(ctx, user)
}
```

### Layer 4: Adapters (Infrastructure)

**Location**: `internal/{feature}/adapters/{adapter_type}/`

**Responsibilities**:
- Implement port interfaces
- Handle infrastructure concerns
- Translate between domain and external systems

**Dependencies**: Domain, ports, and infrastructure libraries

**Example**:
```go
// internal/user/adapters/postgres/repository.go
type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    // Infrastructure implementation
    model := ToUserModel(user)
    return r.db.Create(model).Error
}
```

## Dependency Management

### Compile-Time Dependency Injection

We use **Google Wire** for compile-time dependency injection:

**Benefits**:
- No runtime reflection overhead
- Compile-time errors for missing dependencies
- Clear dependency graphs
- Easy to understand and debug

**Example**:
```go
// internal/wire/providers.go
func ProvideUserService(repo ports.UserRepository) ports.UserService {
    return service.NewUserService(repo)
}

// cmd/api/wire.go
func initializeApp(configPath string) (*gin.Engine, func(), error) {
    wire.Build(wire.ProviderSet)
    return nil, nil, nil
}
```

### Interface-Based Design

All dependencies are injected through interfaces:

```go
// Service depends on interface, not implementation
type UserService struct {
    repo ports.UserRepository  // Interface
}

// Adapter implements interface
type postgresRepository struct {
    db *gorm.DB
}

func (r *postgresRepository) Create(...) error { ... }
```

## Testing Strategy

### Unit Tests

Test domain and service layers with mocks:

```go
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(mocks.MockUserRepository)
    service := service.NewUserService(mockRepo)

    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

    user, err := service.CreateUser(ctx, "test@example.com", "Test")
    assert.NoError(t, err)
}
```

### Integration Tests

Test adapters with real infrastructure using testcontainers:

```go
func TestPostgresRepository_Create(t *testing.T) {
    // Start real PostgreSQL container
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := postgres.NewUserRepository(db)
    user := &domain.User{...}

    err := repo.Create(ctx, user)
    assert.NoError(t, err)
}
```

### Test Coverage Goals

- **Domain**: >95% - Critical business logic
- **Service**: >80% - Use case orchestration
- **Adapters**: >70% - Infrastructure code
- **Overall**: >80% - Comprehensive coverage

## Best Practices

### 1. Domain-Driven Design

- Use ubiquitous language from the business domain
- Model entities based on real-world concepts
- Keep business logic in the domain layer

### 2. Thin Adapters

- Adapters should only translate between domain and infrastructure
- Keep adapter code simple and focused
- Complex logic belongs in the domain or service layer

### 3. Immutable Entities

- Prefer immutable domain entities where possible
- Use factory functions (e.g., `NewUser`) for creation
- Validation happens at entity creation

### 4. Context Propagation

- Always pass `context.Context` as the first parameter
- Use context for cancellation and timeouts
- Propagate context through all layers

### 5. Error Handling

- Define domain-specific errors in the domain layer
- Map infrastructure errors to domain errors in adapters
- Use sentinel errors for common cases

## Further Reading

- [Hexagonal Architecture Details](./hexagonal-architecture.md)
- [Dependency Flow](./dependency-flow.md)
- [ADR: Why Hexagonal Architecture](../adr/001-hexagonal-architecture.md)
- [ADR: Feature-Sliced Structure](../adr/002-feature-sliced-structure.md)
