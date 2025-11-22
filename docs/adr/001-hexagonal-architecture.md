# ADR-001: Use Hexagonal Architecture

## Status

**Accepted** - 2025-11-22

## Context

We need an architecture that:
- Keeps business logic independent of infrastructure
- Supports multiple delivery mechanisms (REST, gRPC, CLI)
- Makes the system highly testable
- Allows easy swapping of infrastructure components
- Scales well for large teams

Traditional layered architecture (Controller → Service → Repository) has several drawbacks:
- Business logic leaks into controllers and repositories
- Tight coupling to frameworks and databases
- Difficult to test without the full stack
- Hard to adapt to new delivery mechanisms

## Decision

We will use **Hexagonal Architecture** (Ports and Adapters pattern) for our application structure.

### Key Components

**Domain Core** (`internal/{feature}/domain/`)
- Contains business entities and logic
- Has zero external dependencies
- Defines domain errors and validation rules

**Ports** (`internal/{feature}/ports/`)
- Define interfaces for input (service) and output (repository)
- Act as contracts between layers
- Allow dependency inversion

**Service** (`internal/{feature}/service/`)
- Implements use case orchestration
- Depends only on domain and port interfaces
- Contains business workflow logic

**Adapters** (`internal/{feature}/adapters/`)
- Implement port interfaces
- Handle infrastructure concerns (HTTP, database, etc.)
- Translate between domain and external systems

### Dependency Rule

Dependencies point inward:
```
External World → Adapters → Ports → Domain
```

The domain never depends on infrastructure.

## Consequences

### Positive

✅ **Business Logic Independence**
- Domain code has no framework dependencies
- Easy to understand and reason about
- Can be used across different applications

✅ **Testability**
- Domain: Pure unit tests, no mocks
- Service: Unit tests with interface mocks
- Adapters: Integration tests with real infrastructure

✅ **Flexibility**
- Swap databases without changing business logic
- Support multiple protocols (REST, gRPC, CLI) easily
- Migrate frameworks with minimal impact

✅ **Team Scalability**
- Clear boundaries enable parallel development
- Developers can work on different adapters independently
- New team members understand architecture quickly

✅ **Maintainability**
- Changes are localized to specific layers
- Adapter changes don't affect domain logic
- Clear separation of concerns

### Negative

⚠️ **Learning Curve**
- Team needs to understand hexagonal architecture
- More upfront design effort required
- Requires discipline to maintain boundaries

⚠️ **More Files**
- More directories and files than simple layered architecture
- Need to navigate between domain, ports, and adapters
- Can feel like over-engineering for simple CRUD

⚠️ **Indirection**
- Extra layer of interfaces (ports)
- Need to map between domain entities and database models
- More code for simple operations

### Mitigations

- Provide comprehensive documentation and examples
- Use feature-sliced structure to keep related code together
- Create generators/templates for new features
- Start with good examples (User feature)

## Alternatives Considered

### 1. Traditional Layered Architecture

**Pros**: Simple, familiar, less code
**Cons**: Tight coupling, hard to test, inflexible

**Why Rejected**: Doesn't meet our testability and flexibility requirements

### 2. Clean Architecture (Uncle Bob)

**Pros**: Similar benefits to hexagonal
**Cons**: More layers (entities, use cases, interface adapters, frameworks), more complex

**Why Rejected**: Hexagonal architecture is simpler while providing the same benefits

### 3. Domain-Driven Design (DDD)

**Pros**: Rich domain modeling, ubiquitous language
**Cons**: Higher complexity, requires domain experts

**Why Rejected**: We're using DDD principles within hexagonal architecture, not as a replacement

## References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) - Alistair Cockburn
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Robert C. Martin
- [Ports and Adapters Pattern](https://herbertograca.com/2017/09/14/ports-adapters-architecture/) - Herberto Graça
