# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) documenting significant architectural decisions made in this project.

## What is an ADR?

An Architecture Decision Record (ADR) is a document that captures an important architectural decision made along with its context and consequences.

## Format

Each ADR follows this structure:

1. **Title** - Short noun phrase
2. **Status** - Proposed, Accepted, Deprecated, Superseded
3. **Context** - What is the issue we're seeing?
4. **Decision** - What have we decided to do?
5. **Consequences** - What becomes easier or harder?

## Index

| ADR | Title | Status |
|-----|-------|--------|
| [001](./001-hexagonal-architecture.md) | Use Hexagonal Architecture | Accepted |
| [002](./002-feature-sliced-structure.md) | Feature-Sliced Vertical Structure | Accepted |
| [003](./003-google-wire-di.md) | Google Wire for Dependency Injection | Accepted |
| [004](./004-gorm-for-postgresql.md) | GORM for PostgreSQL | Accepted |
| [005](./005-testcontainers-integration-tests.md) | Testcontainers for Integration Tests | Accepted |

## Creating a New ADR

1. Copy the template from `000-template.md`
2. Number it sequentially
3. Fill in all sections
4. Update this index
5. Commit with message: `docs: add ADR-NNN {title}`
