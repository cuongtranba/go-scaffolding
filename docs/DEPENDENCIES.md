# Dependencies

This document lists all major dependencies used in the Go Clean Architecture Template and their versions.

Last Updated: 2025-11-22

## Go Version

- **Go**: 1.25 (Latest stable: 1.25.4)
  - Requires: Go 1.25+ for building
  - Dockerfile uses: `golang:1.25-alpine`

## Core Framework & Libraries

### Web Framework
- **Gin**: v1.11.0
  - High-performance HTTP web framework
  - Repository: https://github.com/gin-gonic/gin
  - Requires: Go 1.23+

### ORM & Database
- **GORM**: v1.31.1
  - The fantastic ORM library for Golang
  - Repository: https://github.com/go-gorm/gorm
  - PostgreSQL Driver: v1.6.0
  - SQLite Driver: v1.6.0 (for testing)

### Dependency Injection
- **Google Wire**: v0.7.0
  - Compile-time dependency injection
  - Repository: https://github.com/google/wire

### Configuration
- **Viper**: v1.21.0
  - Configuration management with environment variable support
  - Repository: https://github.com/spf13/viper

### Logging
- **Zerolog**: v1.34.0
  - Zero allocation JSON logger
  - Repository: https://github.com/rs/zerolog

## Database & Migrations

- **golang-migrate**: v4.19.0
  - Database migrations
  - Repository: https://github.com/golang-migrate/migrate

- **pgx/v5**: v5.7.6
  - PostgreSQL driver and toolkit
  - Used by GORM PostgreSQL driver

## Testing

- **Testify**: v1.11.1
  - Testing toolkit with assertions and mocks
  - Repository: https://github.com/stretchr/testify

- **Testcontainers**: v0.40.0
  - Integration testing with Docker containers
  - Repository: https://github.com/testcontainers/testcontainers-go
  - Postgres module: v0.40.0

- **SQLite**: v1.6.0 (gorm driver)
  - In-memory database for unit tests

## Utilities

- **UUID**: v1.6.0 (google/uuid)
  - UUID generation
  - Repository: https://github.com/google/uuid

## Development Tools

- **Taskfile**: External tool for task automation
  - Install: https://taskfile.dev/installation/

- **Mockery**: External tool for generating mocks
  - Version: v3 (configured in .mockery.yml)
  - Install: https://github.com/vektra/mockery

## Docker Images

- **Build**: `golang:1.25-alpine`
  - Latest Go 1.25 with Alpine Linux
  - Minimal size for building

- **Runtime**: `scratch`
  - Minimal container image (empty)
  - Binary-only deployment
  - Includes only: CA certificates, timezone data, application binary

## Infrastructure Dependencies

### OpenTelemetry (Observability)
- **otel**: v1.37.0
- **otel/metric**: v1.37.0
- **otel/trace**: v1.37.0
- **otelhttp**: v0.54.0

### Docker & Container
- **docker/docker**: v28.5.1
- **docker/go-connections**: v0.6.0

## Version Compatibility

| Component | Minimum Version | Latest Tested |
|-----------|----------------|---------------|
| Go | 1.25 | 1.25.4 |
| PostgreSQL | 15 | 16 |
| Docker | 20.10 | 28.5 |
| Docker Compose | 2.0 | Latest |

## Updating Dependencies

### Check for updates
```bash
go list -m -u all | grep -v indirect
```

### Update all dependencies
```bash
go get -u ./...
go mod tidy
```

### Update specific dependency
```bash
go get github.com/gin-gonic/gin@latest
go mod tidy
```

### Verify after updates
```bash
# Run tests
task test

# Build application
task build:api

# Run integration tests
task test:integration
```

## Security Updates

To check for security vulnerabilities:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

## Notes

- All versions are compatible with Go 1.25
- Gin requires minimum Go 1.23
- GORM and all drivers are up to date
- Testcontainers requires Docker to be running
- Wire generates code at compile-time (zero runtime dependencies)

## Change Log

### 2025-11-22
- Updated Go from 1.21 to 1.25
- Updated Dockerfile to use golang:1.25-alpine
- All dependencies verified compatible with Go 1.25
- Current versions:
  - Gin: v1.11.0
  - GORM: v1.31.1
  - Wire: v0.7.0
  - Viper: v1.21.0
  - Zerolog: v1.34.0
  - Testcontainers: v0.40.0
