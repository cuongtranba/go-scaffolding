# Contributing to Go Clean Architecture Template

Thank you for your interest in contributing! This document provides guidelines for contributing to this template.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)
- [Testing Guidelines](#testing-guidelines)

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- Clear and descriptive title
- Steps to reproduce
- Expected vs actual behavior
- Go version and OS
- Relevant log output or screenshots

Use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.yml).

### Suggesting Features

Feature suggestions are welcome! Please:

- Use the [feature request template](.github/ISSUE_TEMPLATE/feature_request.yml)
- Clearly describe the problem and proposed solution
- Explain why this would be useful to most users
- Consider if it fits the template's scope and goals

### Documentation Improvements

Documentation improvements are always appreciated! This includes:

- README clarifications
- Code comments
- Architecture documentation
- ADRs (Architecture Decision Records)
- Examples and tutorials

## Development Setup

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Task (taskfile.dev)
- golang-migrate

### Setup Steps

1. **Fork and clone the repository**

```bash
git clone https://github.com/yourusername/go-scaffolding.git
cd go-scaffolding
```

2. **Install dependencies**

```bash
go mod download
go mod verify
```

3. **Start infrastructure**

```bash
task docker:up
```

4. **Run migrations**

```bash
task migrate:up
```

5. **Verify setup**

```bash
task test
task build:api
```

## Pull Request Process

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions/updates
- `chore/` - Maintenance tasks

### 2. Make Your Changes

- Write clean, readable code
- Follow the coding standards (see below)
- Add or update tests
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
task test

# Run integration tests
task test:integration

# Check coverage
task test:coverage

# Lint your code
go fmt ./...
go vet ./...
```

### 4. Commit Your Changes

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```bash
git commit -m "feat: add user authentication"
git commit -m "fix: resolve database connection pool issue"
git commit -m "docs: update API documentation"
```

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub using the PR template.

### 6. Code Review

- Address review comments promptly
- Keep the PR focused (one feature/fix per PR)
- Ensure CI checks pass
- Update documentation if needed

## Coding Standards

### Go Style Guide

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Architecture Principles

This template follows **Hexagonal Architecture**:

1. **Domain Layer** (`internal/{feature}/domain/`)
   - Pure business logic
   - Zero external dependencies
   - Domain entities and validation

2. **Ports Layer** (`internal/{feature}/ports/`)
   - Interface definitions
   - Contracts between layers

3. **Service Layer** (`internal/{feature}/service/`)
   - Use case orchestration
   - Depends only on domain and ports

4. **Adapters Layer** (`internal/{feature}/adapters/`)
   - Infrastructure implementations
   - HTTP handlers, database repositories

### Code Organization

```
internal/{feature}/
├── domain/        # Business logic - test with zero mocks
├── ports/         # Interfaces - no implementation
├── service/       # Use cases - test with mocked ports
└── adapters/      # Infrastructure - integration tests
```

### Best Practices

1. **Keep functions small** - Single responsibility
2. **Use interfaces** - Depend on abstractions
3. **Write tests first** - TDD when possible
4. **Document exported items** - Godoc comments
5. **Handle errors properly** - Don't ignore errors
6. **Use context.Context** - For cancellation and timeouts

### Example Code

```go
// Good: Clear, testable, follows architecture
package domain

import "errors"

var ErrInvalidEmail = errors.New("invalid email format")

type User struct {
    ID    string
    Email string
    Name  string
}

// NewUser creates a new user with validation
func NewUser(email, name string) (*User, error) {
    if !isValidEmail(email) {
        return nil, ErrInvalidEmail
    }
    return &User{
        Email: email,
        Name:  name,
    }, nil
}
```

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Formatting, missing semicolons, etc.
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks
- `perf:` - Performance improvements
- `ci:` - CI configuration changes

### Examples

```bash
feat: add user authentication endpoint

fix: resolve database connection pool leak

docs: update architecture decision records

test: add integration tests for user service

chore: update dependencies to latest versions
```

## Testing Guidelines

### Test Coverage Goals

- **Domain**: >95% - Critical business logic
- **Service**: >80% - Use case orchestration
- **Adapters**: >70% - Infrastructure code
- **Overall**: >80% - Comprehensive coverage

### Unit Tests

Test domain and service layers with mocks:

```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.MockUserRepository)
    service := service.NewUserService(mockRepo)

    mockRepo.On("GetByEmail", mock.Anything, "test@example.com").
        Return(nil, domain.ErrUserNotFound)
    mockRepo.On("Create", mock.Anything, mock.Anything).
        Return(nil)

    // Act
    user, err := service.CreateUser(context.Background(),
        "test@example.com", "Test User")

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
    mockRepo.AssertExpectations(t)
}
```

### Integration Tests

Use testcontainers for real infrastructure:

```go
func TestPostgresRepository_Create(t *testing.T) {
    // Setup real PostgreSQL container
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := postgres.NewUserRepository(db)
    user := &domain.User{Email: "test@example.com", Name: "Test"}

    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
}
```

### Running Tests

```bash
# Unit tests
go test ./... -v

# Integration tests (requires Docker)
go test ./cmd/api/... -run TestIntegration -v

# With coverage
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Project Structure Decisions

When adding new features:

1. **Create feature directory**: `internal/{feature}/`
2. **Follow layer structure**: domain → ports → service → adapters
3. **Write domain tests first**: Zero external dependencies
4. **Add ADR if needed**: Document architectural decisions
5. **Update documentation**: README, architecture docs

## Questions?

- 📚 Read the [Architecture Overview](./docs/architecture/overview.md)
- 💬 Ask in [GitHub Discussions](https://github.com/yourusername/go-scaffolding/discussions)
- 🐛 Check [existing issues](https://github.com/yourusername/go-scaffolding/issues)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to Go Clean Architecture Template!** 🎉
