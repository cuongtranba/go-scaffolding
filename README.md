# Go Clean Architecture Template

A production-ready Go template following hexagonal/ports & adapters architecture.

## Features

- Hexagonal Architecture with feature-sliced structure
- Multi-protocol support: REST API, gRPC, CLI, Background Workers
- Multiple databases: PostgreSQL (GORM), MongoDB, Redis
- Full observability: structured logging, metrics, tracing
- Comprehensive testing: unit, integration, E2E
- Developer experience: Taskfile, Docker, docker-compose

## Quick Start

```bash
# Clone the repository
git clone <your-repo-url>
cd go-scaffolding

# Start infrastructure
task docker:up

# Run the API server
task run:api
```

## Documentation

See [docs/](./docs/) for detailed documentation.

## License

MIT
