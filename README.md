# Retail Store API

A microservices-oriented backend system for retail store management, built with **Go**. Designed around an **event-driven architecture** with four independently deployable services.

## Architecture

```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   API Server    │  │  Admin Server   │  │  Socket Server  │  │    Consumer     │
│   (Public)      │  │  (Admin Only)   │  │  (WebSocket)    │  │ (Event Handler) │
│   :8080         │  │   :8082         │  │   :8081         │  │   No HTTP       │
└────────┬────────┘  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘
         │                    │                     │                     │
         │   Redis Pub/Sub ───┼─────────────────────┘                     │
         │   NATS JetStream ──┼───────────────────────────────────────────┘
         │                    │
    ┌────┴────┐  ┌─────────┐  ┌────┴────┐  ┌───────────────┐
    │Postgres │  │  Redis  │  │  NATS   │  │ Elasticsearch │
    └─────────┘  └─────────┘  └─────────┘  └───────────────┘
```

| Service          | Port | Description                                                    |
| ---------------- | ---- | -------------------------------------------------------------- |
| **API Server**   | 8080 | Public REST API — read-only endpoints, authentication, search  |
| **Admin Server** | 8082 | Admin REST API — full CRUD, requires JWT with `admin` role     |
| **Socket Server**| 8081 | WebSocket — real-time notifications via Redis Pub/Sub          |
| **Consumer**     | —    | NATS JetStream listener — analytics, inventory, search indexing|

## Tech Stack

| Category        | Technology                                                   |
| --------------- | ------------------------------------------------------------ |
| Language        | Go 1.25                                                      |
| HTTP Framework  | Gin                                                          |
| Database        | PostgreSQL 16                                                |
| Cache / PubSub  | Redis 7                                                      |
| Search Engine   | Elasticsearch 8                                              |
| Event Stream    | NATS JetStream                                               |
| Real-time       | WebSocket (gorilla/websocket)                                |
| Authentication  | JWT (golang-jwt) + bcrypt                                    |
| Code Generation | [sqlc](https://sqlc.dev/)                                    |
| Migrations      | [golang-migrate](https://github.com/golang-migrate/migrate) |
| API Docs        | Swagger (swaggo)                                             |
| Logging         | slog + ELK Stack (Logstash, Kibana)                          |
| Deployment      | Docker, Kubernetes (kustomize)                               |

## Functional Requirements

### Authentication & Authorization

- User registration with email, password (bcrypt hashed), and name
- JWT-based login with configurable token expiration
- Role-based access control: `user` (read-only) and `admin` (full CRUD)
- Bearer token validation middleware with role checking

### Product Management

- **Create** / **Update** / **Delete** — Admin only, publishes domain events to NATS
- **List** (paginated) / **Get by ID** — Public, no authentication required
- Full-text search via Elasticsearch with fuzzy matching and weighted relevance (name 3x)
- Search index automatically maintained by consumer on create/update/delete events

### Category Management

- **Create** / **Update** / **Delete** — Admin only
- **List** / **Get by ID** — Public
- Products linked via foreign key (ON DELETE SET NULL)

### Real-time Notifications

- WebSocket connections managed by Socket Server
- Redis Pub/Sub bridges events to all connected clients
- Supports promotion and discount notification channels

### Event Processing (Consumer)

- **Analytics** — Tracks order creation and product views
- **Inventory** — Updates stock on order events
- **Search Indexing** — Syncs product data to Elasticsearch

## Non-Functional Requirements

### Performance & Scalability

- Horizontal scaling via Kubernetes replicas (API/Admin: 2, Socket/Consumer: 1)
- PostgreSQL connection pooling (pgx)
- Pagination on all list endpoints (limit 1–100, offset)
- Database indexes on frequently queried columns (category_id, name, email)

### Reliability

- Graceful shutdown with 5-second timeout for in-flight requests
- NATS JetStream with explicit ACK — messages redelivered on failure
- Event publishing failures logged but don't block API responses
- Health check endpoint (`GET /health`) for Kubernetes liveness/readiness probes

### Security

- Passwords hashed with bcrypt (default cost)
- JWT HS256 signing with configurable secret
- Admin endpoints isolated on separate server (port 8082)
- No sensitive data in logs or error responses

### Observability

- Structured JSON logging via Go slog
- Request logging: method, path, status, latency, client IP
- ELK stack integration: Logstash (TCP 5044), Elasticsearch, Kibana (port 5601)
- Error logging with context on all initialization failures

### Deployment

- Multi-stage Docker build (golang:1.24-alpine → alpine:3.19)
- Docker Compose for local development (all services + infrastructure)
- Kubernetes with kustomize overlays (dev/prod)
- Resource limits: CPU 100–500m, Memory 128–256Mi per service
- Liveness/readiness probes with configurable intervals

## API Endpoints

### Public API (:8080)

```
POST /api/v1/auth/register                          # Register
POST /api/v1/auth/login                             # Login → JWT token
GET  /api/v1/products?limit=20&offset=0             # List products
GET  /api/v1/products/:id                           # Get product
GET  /api/v1/categories                             # List categories
GET  /api/v1/categories/:id                         # Get category
GET  /api/v1/search/products?q=keyword              # Full-text search
GET  /health                                        # Health check
```

### Admin API (:8082) — `Authorization: Bearer <token>` required

```
POST   /api/v1/products                             # Create product
PUT    /api/v1/products/:id                         # Update product
DELETE /api/v1/products/:id                         # Delete product
POST   /api/v1/categories                           # Create category
PUT    /api/v1/categories/:id                       # Update category
DELETE /api/v1/categories/:id                       # Delete category
GET    /api/v1/products | /categories | /search     # Also available
```

### WebSocket (:8081)

```
WS ws://localhost:8081/api/v1/ws/notifications      # Real-time updates
```

Swagger UI: `http://localhost:8080/swagger/index.html`

## Quick Start

```bash
cp .env.example .env
make docker-infra                   # Start PostgreSQL, Redis, NATS, Elasticsearch
make migrate-up                     # Apply database migrations

# Run each in a separate terminal
make run-api                        # Public API   → :8080
make run-admin                      # Admin API    → :8082
make run-socket                     # WebSocket    → :8081
make run-consumer                   # Event consumer
```

## Development Commands

| Command              | Description                              |
| -------------------- | ---------------------------------------- |
| `make build`         | Build all binaries to `bin/`             |
| `make sqlc`          | Generate Go code from SQL queries        |
| `make swagger`       | Generate Swagger documentation           |
| `make migrate-up`    | Apply database migrations                |
| `make migrate-down`  | Rollback database migrations             |
| `make docker-up`     | Start all services via Docker Compose    |
| `make docker-build`  | Build Docker images                      |
| `make k8s-dev`       | Deploy to Kubernetes (dev)               |
| `make k8s-prod`      | Deploy to Kubernetes (prod)              |
