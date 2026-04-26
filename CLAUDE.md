# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Run services (each in a separate terminal)
make run-api          # Public API server on :8080 (read-only + auth)
make run-admin        # Admin API server on :8082 (full CRUD, JWT required)
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

Four independent services share a single repo but run as separate processes. Each service has its own `cmd/` entrypoint and `internal/` package:

```
services/
├── api/          # Public REST API (:8080) — read-only + auth
├── admin/        # Admin REST API (:8082) — full CRUD, JWT required
├── consumer/     # NATS JetStream consumer — no HTTP
└── socket/       # WebSocket server (:8081) — real-time notifications
```

### Project layout

```
retail-store-api/
├── pkg/                          # Shared packages (all services import from here)
│   ├── auth/                     # JWT generation & validation
│   ├── config/                   # App configuration (env vars)
│   ├── database/                 # PostgreSQL & Redis clients
│   ├── elasticsearch/            # ES client & product indexing
│   ├── logger/                   # Structured logging (slog)
│   ├── middleware/               # Gin middleware (auth, logger)
│   ├── model/                    # Shared DTOs
│   │   ├── request/              # Request structs
│   │   └── response/             # Response structs
│   ├── nats/                     # NATS connection, publisher, subjects
│   ├── notification/             # WebSocket hub/client/subscriber
│   ├── repository/               # Data access (wraps sqlc queries)
│   └── response/                 # HTTP response helpers
├── services/
│   ├── api/
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── handler/          # HTTP handlers (read-only)
│   │       └── service/          # Business logic (no writes)
│   ├── admin/
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── handler/          # HTTP handlers (full CRUD)
│   │       └── service/          # Business logic + NATS publish
│   ├── consumer/
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── consumer.go       # NATS subscription manager
│   │       └── handler/          # Event handlers (analytics, inventory, search index)
│   └── socket/
│       ├── cmd/main.go
│       └── internal/
│           └── handler/          # WebSocket handler
├── db/
│   ├── query/                    # SQL for sqlc (1 file = 1 table)
│   ├── migration/                # Database migrations
│   └── sqlc/                     # Generated Go code — DO NOT EDIT
└── deploy/k8s/                   # Kubernetes manifests (base + overlays)
```

### Request flow

```
Handler → Service → Repository → sqlc-generated queries → PostgreSQL
                ↓
          NATS Publish (domain event) + Redis Publish (real-time notification)
```

### Authentication

- JWT-based auth using `pkg/auth/jwt.go` (generation/validation) and `pkg/middleware/auth.go` (Gin middleware).
- `middleware.Auth(jwtManager)` validates the Bearer token and sets `user_id`, `email`, `role` in Gin context.
- `middleware.RequireRole("admin")` checks the role from context.
- Auth endpoints in each service's `internal/handler/auth_handler.go` with business logic in `internal/service/auth_service.go`.
- Config: `JWT_SECRET` and `JWT_EXPIRATION_HOURS` env vars, loaded into `config.JWTConfig`.

### Dependency wiring

Each service manages its own dependencies in `internal/app.go` via a `Dependencies` struct. There is no shared `AppContext` — each service initializes only the infrastructure it needs. Route registration is done in each service's `internal/router.go`.

### Event system

- Stream name: `RETAIL_STORE` with WorkQueuePolicy
- Subject wildcards: `orders.>`, `products.>`
- Specific subjects defined in `pkg/nats/subjects.go` (e.g., `products.created`, `orders.created`)
- Consumer registration in `services/consumer/internal/consumer.go` — each subscription has a durable consumer name and explicit ack

### Key patterns

- **Handlers** (`services/*/internal/handler/`): constructed with service dependency, expose methods like `Create`, `List`, `GetByID`. Routes wired in each service's `internal/router.go`.
- **Services** (`services/*/internal/service/`): contain business logic, call repository methods, publish NATS events after mutations (admin only).
- **Repositories** (`pkg/repository/`): thin wrappers around sqlc-generated `db.Queries`. Shared across services.
- **DTOs**: request structs in `pkg/model/request/`, response structs in `pkg/model/response/`.
- **Import aliases**: use `pkgNats` for `pkg/nats`, `pkgResponse` for `pkg/response`, `es` for `pkg/elasticsearch`.
- **Error logging**: always `slog.Error(...)` before returning errors in app initialization (`New()` functions).

## Code Generation

- **sqlc**: SQL queries in `db/query/` generate type-safe Go code into `db/sqlc/`. Never edit `db/sqlc/` manually. Config in `sqlc.yaml` (engine: PostgreSQL, driver: pgx/v5).
- **Swagger**: annotations live in handler files. Run `make swagger` after changes.

## Infrastructure Dependencies

PostgreSQL 16, Redis 7, NATS 2 (JetStream enabled), Elasticsearch 8.12, Logstash + Kibana for logging. All defined in `docker-compose.yml`. Config loaded from `.env` with defaults in `pkg/config/config.go`.
