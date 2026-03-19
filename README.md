# Retail Store API

Hệ thống quản lý cửa hàng tạp hoá, viết bằng Go.

## Yêu cầu

- Go 1.24+
- PostgreSQL
- Redis
- NATS Server
- [sqlc](https://sqlc.dev/) (generate code từ SQL)
- [golang-migrate](https://github.com/golang-migrate/migrate) (chạy migration)

## Cài đặt

```bash
# Clone repo
git clone <repo-url>
cd retail-store-api

# Copy file env
cp .env.example .env
# Sửa .env cho phù hợp với môi trường local

# Cài dependencies
go mod tidy

# Tạo database
createdb retail_store

# Chạy migration
make migrate-up

# Generate sqlc (nếu sửa file SQL)
make sqlc
```

## Chạy

Dự án có 3 server, chạy độc lập trên 3 terminal:

```bash
make run-api        # HTTP API         → localhost:8080
make run-socket     # WebSocket server → localhost:8081
make run-consumer   # NATS consumer    → không có HTTP, chỉ lắng nghe event
```

## Cấu trúc dự án

```
retail-store-api/
│
├── cmd/                          # Entry point cho từng server
│   ├── api/main.go               # HTTP API server
│   ├── socket/main.go            # WebSocket server
│   └── consumer/main.go          # NATS consumer
│
├── db/                           # Database
│   ├── migration/                # SQL migration files (up/down)
│   ├── query/                    # SQL queries cho sqlc
│   └── sqlc/                     # Code được generate tự động (KHÔNG sửa tay)
│
├── internal/                     # Code nội bộ (không export ra ngoài)
│   ├── app/                      # Khởi tạo từng server
│   │   ├── api/                  # API server: kết nối DB, Redis, setup router
│   │   ├── socket/               # Socket server: kết nối Redis, setup WebSocket
│   │   └── consumer/             # Consumer: kết nối NATS, đăng ký handler
│   │
│   ├── config/                   # Đọc biến môi trường từ .env
│   │
│   ├── handler/                  # Xử lý HTTP request (validate input, trả response)
│   ├── service/                  # Business logic (xử lý nghiệp vụ)
│   ├── repository/               # Truy vấn database (gọi sqlc)
│   │
│   ├── model/                    # Định nghĩa struct
│   │   ├── request/              # Struct cho request body (CreateProductRequest, ...)
│   │   └── response/             # Struct cho response data (ProductResponse, ...)
│   │
│   └── consumer/                 # Consumer engine + handler xử lý event
│       └── handler/              # Từng handler: analytics, inventory, ...
│
├── pkg/                          # Code dùng chung (copy sang project khác được)
│   ├── database/                 # Helper kết nối DB (PostgreSQL, Redis)
│   ├── middleware/               # Gin middleware (logger, ...)
│   ├── nats/                     # NATS connection, publisher, subjects, stream
│   ├── notification/             # WebSocket hub, Redis Pub/Sub subscriber
│   └── response/                 # Helper: Success(), Error()
│
├── .env.example                  # Mẫu biến môi trường
├── Makefile                      # Lệnh tắt
└── sqlc.yaml                     # Cấu hình sqlc
```

## Kiến trúc

### Tổng quan

Dự án gồm **3 server chạy độc lập**, giao tiếp với nhau qua **Redis** và **NATS**:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  API Server │     │Socket Server│     │  Consumer   │
│   :8080     │     │   :8081     │     │  (no HTTP)  │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       │    ┌──────────────┘                   │
       │    │  Redis Pub/Sub                   │
       │    │  (realtime notification)         │
       │    │                                  │
       ▼    ▼                                  │
   ┌──────────┐                                │
   │  Redis   │                                │
   └──────────┘                                │
       │                                       │
       │         NATS JetStream                │
       │         (event processing)            │
       │                                       │
       ▼                                       ▼
   ┌──────────┐                         ┌──────────────┐
   │   NATS   │ ──────────────────────▶ │   Handlers   │
   └──────────┘                         │ - analytics  │
       │                                │ - inventory  │
       ▼                                └──────────────┘
   ┌──────────┐
   │ Postgres │
   └──────────┘
```

### Flow xử lý request (API Server)

```
Request → Handler → Service → Repository → DB (sqlc)
                      │
                      └──→ NATS Publish (nếu cần gửi event)
```

Mỗi layer có nhiệm vụ riêng:

| Layer | Nhiệm vụ | Ví dụ |
|---|---|---|
| **Handler** | Validate input, trả response | Kiểm tra `name` không rỗng |
| **Service** | Xử lý business logic | Tính giá sau giảm, check tồn kho |
| **Repository** | Gọi database | Gọi sqlc query |

### Flow WebSocket (Socket Server)

```
Client ──WebSocket──▶ Socket Server
                          │
Redis PUBLISH ──────▶ Subscriber ──▶ Hub ──▶ Broadcast tới tất cả clients
```

### Flow Event (Consumer)

```
API Server ──NATS Publish──▶ NATS JetStream ──▶ Consumer ──▶ Handler
                                                              │
                                                    ┌─────────┼──────────┐
                                                    ▼         ▼          ▼
                                               Analytics  Inventory    ...
```

## API Endpoints

### Health
```
GET /health
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

### WebSocket
```
ws://localhost:8081/api/v1/ws/notifications
```

## Swagger

Swagger UI có sẵn tại `http://localhost:8080/swagger/index.html` khi chạy API server.

Sau khi thêm/sửa swagger annotations, chạy lại:

```bash
make swagger
```

### Snippets

Trong Cursor/VSCode, gõ prefix rồi nhấn **Tab** để sinh swagger annotation block:

| Prefix | Sinh ra |
|---|---|
| `swag-get` | GET endpoint (với path param) |
| `swag-list` | GET list (có sẵn limit/offset) |
| `swag-post` | POST endpoint (có body) |
| `swag-put` | PUT endpoint (có id + body) |
| `swag-del` | DELETE endpoint (có id) |

## Thêm tính năng mới

### Thêm API endpoint mới

1. Tạo request/response struct trong `internal/model/request/` và `internal/model/response/`
2. Tạo repository trong `internal/repository/`
3. Tạo service trong `internal/service/`
4. Tạo handler trong `internal/handler/`, dùng snippet `swag-*` để thêm swagger annotations
5. Đăng ký handler trong `internal/app/api/router.go`
6. Chạy `make swagger` để cập nhật docs

### Thêm consumer handler mới

1. Thêm subject mới trong `pkg/nats/subjects.go`
2. Tạo handler trong `internal/consumer/handler/`
3. Đăng ký trong `internal/app/consumer/consumer.go`

### Thêm SQL query mới

1. Viết query trong `db/query/`
2. Chạy `make sqlc` để generate code
3. Dùng trong repository

### Thêm migration mới

1. Tạo file `db/migration/000002_<tên>.up.sql` và `.down.sql`
2. Chạy `make migrate-up`

## Makefile

| Lệnh | Mô tả |
|---|---|
| `make run-api` | Chạy API server |
| `make run-socket` | Chạy WebSocket server |
| `make run-consumer` | Chạy NATS consumer |
| `make build` | Build tất cả ra thư mục `bin/` |
| `make sqlc` | Generate Go code từ SQL |
| `make swagger` | Generate swagger docs |
| `make migrate-up` | Chạy migration |
| `make migrate-down` | Rollback migration |
