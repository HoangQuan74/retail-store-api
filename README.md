# Retail Store API

A microservices-oriented backend system for retail store management, built with Go. The system is designed around an **event-driven architecture** with three independently deployable services communicating via NATS JetStream and Redis Pub/Sub.

## Architecture Overview

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│   API Server    │       │  Socket Server  │       │    Consumer     │
│    (REST)       │       │  (WebSocket)    │       │ (Event Handler) │
│   :8080         │       │   :8081         │       │   No HTTP       │
└────────┬────────┘       └────────┬────────┘       └────────┬────────┘
         │                         │                          │
         │   Redis Pub/Sub         │                          │
         │   (real-time push)      │                          │
         ├─────────────────────────┘                          │
         │                                                    │
         │   NATS JetStream                                   │
         │   (async event processing)                         │
         ├────────────────────────────────────────────────────┘
         │
    ┌────┴────┐   ┌─────────┐   ┌───────────────┐
    │ Postgres│   │  Redis  │   │ Elasticsearch │
    └─────────┘   └─────────┘   └───────────────┘
```

**API Server** — Handles REST requests, performs CRUD operations, and publishes domain events.

**Socket Server** — Manages WebSocket connections. Subscribes to Redis Pub/Sub channels and pushes real-time updates to connected clients.

**Consumer** — Listens to NATS JetStream events and dispatches them to specialized handlers (analytics, inventory, search indexing).

### Request Processing Flow

```
HTTP Request → Handler → Service → Repository → PostgreSQL
                           │
                           ├──→ NATS Publish (domain event)
                           └──→ Redis Publish (real-time notification)
```

Each layer has a single responsibility:

| Layer          | Responsibility                          |
| -------------- | --------------------------------------- |
| **Handler**    | Input validation, HTTP response mapping |
| **Service**    | Business logic, event publishing        |
| **Repository** | Data access via sqlc-generated queries  |

### Event Processing Flow

```
API Server ── NATS Publish ──→ JetStream ──→ Consumer ──→ Handler
                                                          ├── AnalyticsHandler
                                                          ├── InventoryHandler
                                                          └── SearchIndexHandler (→ Elasticsearch)
```

## Tech Stack

| Category       | Technology                                  |
| -------------- | ------------------------------------------- |
| Language       | Go 1.25                                     |
| HTTP Framework | Gin                                         |
| Database       | PostgreSQL 16                               |
| Cache / PubSub | Redis 7                                     |
| Search Engine  | Elasticsearch 8                             |
| Event Stream   | NATS JetStream                              |
| Real-time      | WebSocket (gorilla/websocket)               |
| SQL Codegen    | [sqlc](https://sqlc.dev/)                   |
| Migrations     | [golang-migrate](https://github.com/golang-migrate/migrate) |
| API Docs       | Swagger (swaggo)                            |
| Logging        | slog + Logstash (ELK)                       |
| Deployment     | Docker, Kubernetes (kustomize)              |

## Project Structure

```
retail-store-api/
├── cmd/                             # Application entry points
│   ├── api/main.go                  # REST API server
│   ├── socket/main.go               # WebSocket server
│   └── consumer/main.go             # NATS event consumer
│
├── internal/                        # Private application code
│   ├── app/                         # Server bootstrap & dependency wiring
│   │   ├── api/                     # API server init + route registration
│   │   ├── socket/                  # Socket server init
│   │   └── consumer/                # Consumer init + handler registration
│   ├── config/                      # Environment-based configuration
│   ├── handler/                     # HTTP handlers
│   ├── service/                     # Business logic
│   ├── repository/                  # Data access layer
│   ├── model/
│   │   ├── request/                 # Request DTOs
│   │   └── response/                # Response DTOs
│   └── consumer/
│       └── handler/                 # Event handlers (analytics, inventory, search)
│
├── pkg/                             # Reusable packages
│   ├── database/                    # PostgreSQL & Redis connection helpers
│   ├── middleware/                   # HTTP middleware (logging, recovery)
│   ├── nats/                        # NATS client, publisher, stream config
│   ├── notification/                # WebSocket hub, Redis Pub/Sub bridge
│   ├── response/                    # Standardized HTTP response helpers
│   ├── logger/                      # Structured logging with Logstash support
│   └── elasticsearch/               # Elasticsearch client & indexing
│
├── db/
│   ├── migration/                   # SQL migration files (up/down)
│   ├── query/                       # SQL queries for sqlc
│   └── sqlc/                        # Auto-generated Go code (do not edit)
│
├── deploy/
│   ├── k8s/                         # Kubernetes manifests
│   │   ├── base/                    # Base resources
│   │   └── overlays/                # Environment-specific overrides (dev/prod)
│   ├── logstash/                    # Logstash pipeline config
│   └── kibana/                      # Kibana dashboards
│
├── docker-compose.yml               # Local development stack
├── Dockerfile                       # Multi-stage build
├── Makefile                         # Development commands
└── sqlc.yaml                        # sqlc configuration
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose

### Quick Start

```bash
# Clone the repository
git clone <repo-url>
cd retail-store-api

# Copy and configure environment variables
cp .env.example .env

# Start infrastructure services (PostgreSQL, Redis, NATS, Elasticsearch)
make docker-up

# Run database migrations
make migrate-up

# Start all services (run each in a separate terminal)
make run-api          # REST API         → http://localhost:8080
make run-socket       # WebSocket server → ws://localhost:8081
make run-consumer     # Event consumer   → listens on NATS JetStream
```

### Using Docker

```bash
# Build and run the entire stack
make docker-build
docker compose up
```

## API Endpoints

### Health Check

```
GET  /health
```

### Categories

```
POST   /api/v1/categories
GET    /api/v1/categories
GET    /api/v1/categories/:id
PUT    /api/v1/categories/:id
DELETE /api/v1/categories/:id
```

### Products

```
POST   /api/v1/products
GET    /api/v1/products?limit=20&offset=0
GET    /api/v1/products/:id
PUT    /api/v1/products/:id
DELETE /api/v1/products/:id
```

### Search

```
GET  /api/v1/search/products?q=keyword&limit=20&offset=0
```

### WebSocket

```
WS   ws://localhost:8081/api/v1/ws/notifications
```

### API Documentation

Swagger UI is available at `http://localhost:8080/swagger/index.html` when the API server is running.

```bash
# Regenerate Swagger docs after modifying annotations
make swagger
```

## Development

### Available Commands

| Command              | Description                              |
| -------------------- | ---------------------------------------- |
| `make run-api`       | Start the API server                     |
| `make run-socket`    | Start the WebSocket server               |
| `make run-consumer`  | Start the NATS consumer                  |
| `make build`         | Build all binaries to `bin/`             |
| `make sqlc`          | Generate Go code from SQL queries        |
| `make swagger`       | Generate Swagger documentation           |
| `make migrate-up`    | Apply database migrations                |
| `make migrate-down`  | Rollback database migrations             |
| `make docker-up`     | Start infrastructure via Docker Compose  |
| `make docker-build`  | Build Docker images                      |
| `make k8s-dev`       | Deploy to Kubernetes (dev overlay)       |

### Adding a New API Endpoint

1. Define request/response structs in `internal/model/`
2. Add SQL queries in `db/query/` and run `make sqlc`
3. Implement repository in `internal/repository/`
4. Implement service in `internal/service/`
5. Implement handler in `internal/handler/`
6. Register routes in `internal/app/api/router.go`
7. Run `make swagger` to update API docs

### Adding a New Event Handler

1. Define the subject in `pkg/nats/subjects.go`
2. Implement the handler in `internal/consumer/handler/`
3. Register it in `internal/app/consumer/consumer.go`

### Database Changes

```bash
# Create migration files
# db/migration/000002_<name>.up.sql
# db/migration/000002_<name>.down.sql

# Apply
make migrate-up
```

## Deployment

The project supports deployment via **Docker Compose** for local/staging environments and **Kubernetes** with kustomize overlays for production.

```bash
# Kubernetes (dev)
make k8s-dev

# Kubernetes (prod) — apply production overlay
kubectl apply -k deploy/k8s/overlays/prod
```

Logging is handled through the **ELK stack** (Elasticsearch, Logstash, Kibana) with structured log output via Go's `slog` package.
