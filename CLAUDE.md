# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Run services (each in a separate terminal)
make run-api          # REST API server on :8080
make run-socket       # WebSocket server on :8081
make run-consumer     # NATS event consumer (no HTTP)

# Build all binaries to bin/
make build

# Infrastructure
make docker-up        # Start all Docker services (app + infra)
make docker-infra     # Start infra only (Postgres, Redis, NATS, ES)
make docker-build     # Build Docker images
make docker-down      # Stop all

# Database
make migrate-up       # Apply migrations
make migrate-down     # Rollback migrations
make sqlc             # Regenerate Go code from SQL queries (db/query/ → db/sqlc/)

# API docs
make swagger          # Regenerate Swagger from handler annotations

# Kubernetes
make k8s-dev          # Deploy dev overlay
make k8s-prod         # Deploy prod overlay
```

## Architecture

Three independent services share a single repo but run as separate processes:

- **API Server** (`cmd/api`) — REST API, CRUD, publishes events to NATS and notifications to Redis Pub/Sub.
- **Socket Server** (`cmd/socket`) — WebSocket server. Subscribes to Redis Pub/Sub and pushes real-time updates to connected clients via a Hub/Client pattern.
- **Consumer** (`cmd/consumer`) — Subscribes to NATS JetStream. Dispatches events to handlers: analytics, inventory, and Elasticsearch search indexing.

### Request flow (API)

```
Handler → Service → Repository → sqlc-generated queries → PostgreSQL
                ↓
          NATS Publish (domain event) + Redis Publish (real-time notification)
```

### Dependency wiring

All shared dependencies are passed via `AppContext` (`internal/app/context.go`), which holds Config, sqlc Queries, Elasticsearch client, NATS Publisher, and WebSocket Hub. Handlers receive `AppContext` + Gin router in their constructors and register their own routes.

### Event system

- Stream name: `RETAIL_STORE` with WorkQueuePolicy
- Subject wildcards: `orders.>`, `products.>`
- Specific subjects defined in `pkg/nats/subjects.go` (e.g., `products.created`, `orders.created`)
- Consumer registration in `internal/app/consumer/consumer.go` — each subscription has a durable consumer name and explicit ack

### Key patterns

- **Handlers** (`internal/handler/`): receive `AppContext`, set up their route group in the constructor, use `pkg/response.Success()` / `pkg/response.Error()` for standardized responses.
- **Services** (`internal/service/`): contain business logic, call repository methods, publish NATS events after mutations.
- **Repositories** (`internal/repository/`): thin wrappers around sqlc-generated `db.Queries`.
- **DTOs**: request structs in `internal/model/request/`, response structs in `internal/model/response/`.

## Code Generation

- **sqlc**: SQL queries in `db/query/` generate type-safe Go code into `db/sqlc/`. Never edit `db/sqlc/` manually. Config in `sqlc.yaml` (engine: PostgreSQL, driver: pgx/v5).
- **Swagger**: annotations live in handler files. Run `make swagger` after changes.

## Infrastructure Dependencies

PostgreSQL 16, Redis 7, NATS 2 (JetStream enabled), Elasticsearch 8.12, Logstash + Kibana for logging. All defined in `docker-compose.yml`. Config loaded from `.env` via godotenv with defaults in `internal/config/config.go`.
