# Retail Store API — Coding Conventions

> Convention nội bộ cho Retail Store API, viết bằng Go (Gin), PostgreSQL (sqlc + pgx), Redis, NATS, Elasticsearch.

---

## 1. Nguyên tắc chung

- Tuân thủ [Effective Go](https://go.dev/doc/effective_go) và [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).
- **Rõ ràng hơn ngắn gọn**: `productCategory` tốt hơn `pc`.
- **Ngắn gọn trong scope nhỏ**: biến trong loop 3 dòng có thể là `i`, `p`, `tx`.
- **Không stutter**: `product.ProductService` → ❌ · `product.Service` → ✅.
- **Tiếng Anh cho tất cả tên** (biến, hàm, file, comment).

---

## 2. Cấu trúc thư mục

Mỗi service có `cmd/` và `internal/` riêng, shared code nằm trong `pkg/`:

```
retail-store-api/
├── pkg/                              # Shared packages
│   ├── auth/                         # JWT
│   ├── config/                       # App configuration
│   ├── database/                     # PostgreSQL & Redis clients
│   ├── elasticsearch/                # ES client & indexing
│   ├── logger/                       # Structured logging
│   ├── middleware/                    # Gin middleware (auth, logger)
│   ├── model/
│   │   ├── request/                  # Request DTOs
│   │   └── response/                 # Response DTOs
│   ├── nats/                         # NATS connection, publisher, subjects
│   ├── notification/                 # WebSocket hub/client
│   ├── repository/                   # Data access (wrap sqlc queries)
│   └── response/                     # HTTP response helpers
│
├── services/
│   ├── api/                          # Public API :8080
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── handler/
│   │       └── service/
│   ├── admin/                        # Admin API :8082
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── handler/
│   │       └── service/
│   ├── consumer/                     # NATS consumer
│   │   ├── cmd/main.go
│   │   └── internal/
│   │       ├── consumer.go
│   │       └── handler/
│   └── socket/                       # WebSocket :8081
│       ├── cmd/main.go
│       └── internal/
│           └── handler/
│
├── db/
│   ├── query/                        # SQL cho sqlc (1 file = 1 bảng)
│   ├── migration/                    # Migrations
│   └── sqlc/                         # Generated — KHÔNG sửa tay
└── deploy/k8s/
```

Quy tắc:

- `services/<name>/cmd/main.go` chỉ làm wiring: đọc config, khởi tạo, start server.
- `services/<name>/internal/` chứa code riêng của từng service (handler, service, app setup).
- `pkg/` chứa code dùng chung: config, model, repository, middleware, infra clients.

---

## 3. Đặt tên

### Package

- 1 từ, viết thường, không `_`, không camelCase.
- Tránh: `util`, `common`, `helper`, `misc`.

### File

- snake_case: `product_handler.go`, `category_service.go`, `auth_request.go`.
- Mỗi file = 1 concern. Pattern: `<entity>_<layer>.go`.

### Biến

- camelCase (unexported), PascalCase (exported).
- Acronym giữ nguyên case: `userID`, `httpClient`, `apiKey`, `esClient`.
- Scope nhỏ → tên ngắn. Scope rộng → tên rõ ràng.

| Loại           | Convention            | Ví dụ                     |
| -------------- | --------------------- | ------------------------- |
| ID             | `<entity>ID`          | `productID`, `categoryID` |
| Số lượng       | `<entity>Count`       | `productCount`            |
| Thời điểm      | `<event>At`           | `createdAt`, `updatedAt`  |
| Boolean        | `is/has/can + <Attr>` | `isActive`, `hasStock`    |
| Context        | `ctx`                 | `ctx context.Context`     |
| Error          | `err`                 | `if err != nil`           |
| DB Transaction | `tx`                  | `tx pgx.Tx`               |

### Hằng số

- camelCase/PascalCase (Go **không** dùng `ALL_CAPS`).
- Không hard-code magic number/string trong body hàm.

### Struct & field

- Struct: PascalCase, danh từ số ít. Không suffix `Struct`, `Data`.
- Field exported: PascalCase. JSON/DB tag: snake_case.
- Request/Response: suffix rõ ràng — `CreateProductRequest`, `ProductResponse`.

```go
type Product struct {
    ID          int64     `json:"id"          db:"id"`
    Name        string    `json:"name"        db:"name"`
    CategoryID  int64     `json:"category_id" db:"category_id"`
    Price       float64   `json:"price"       db:"price"`
    CreatedAt   time.Time `json:"created_at"  db:"created_at"`
}
```

### Interface

- Single-method → suffix `-er`: `Reader`, `Publisher`.
- Multi-method → danh từ mô tả: `ProductRepository`.
- Khai báo ở nơi **dùng**, không phải nơi implement. Không prefix `I`.

### Function & method

- Động từ rõ nghĩa: `Create`, `Update`, `Delete`, `GetByID`, `List`.
- Constructor: `New<Type>` — `NewProductHandler(svc)`, `NewProductService(queries, publisher)`.
- Receiver ngắn, nhất quán: `h` cho handler, `s` cho service, `r` cho repository.

### Error

- Sentinel: prefix `Err` — `ErrProductNotFound`, `ErrInvalidInput`.
- Error type: suffix `Error` — `ValidationError`.
- Wrap: `fmt.Errorf("<layer>: <action>: %w", err)`.

---

## 4. Patterns trong project

### Handler

```go
type ProductHandler struct {
    service *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler { ... }
func (h *ProductHandler) Create(c *gin.Context)  { ... }
func (h *ProductHandler) GetByID(c *gin.Context) { ... }
func (h *ProductHandler) List(c *gin.Context)    { ... }
func (h *ProductHandler) Update(c *gin.Context)  { ... }
func (h *ProductHandler) Delete(c *gin.Context)  { ... }
```

Route registration ở `services/*/internal/router.go`.

### Service

```go
type ProductService struct {
    repo      *repository.ProductRepository
    publisher *pkgNats.Publisher
}

func NewProductService(queries *db.Queries, publisher *pkgNats.Publisher) *ProductService { ... }
func (s *ProductService) Create(ctx context.Context, req request.CreateProductRequest) (db.Product, error) { ... }
```

### Repository

Thin wrapper quanh sqlc `db.Queries` (nằm trong `pkg/repository/`):

```go
type ProductRepository struct {
    queries *db.Queries
}

func NewProductRepository(queries *db.Queries) *ProductRepository { ... }
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (db.Product, error) { ... }
```

### Consumer handler (NATS)

```go
func (h *SearchIndexHandler) HandleProductCreated(msg jetstream.Msg) { ... }
func (h *SearchIndexHandler) HandleProductUpdated(msg jetstream.Msg) { ... }
```

Dùng `msg.Ack()` khi thành công, `msg.Nak()` khi lỗi.

### Import aliases

```go
import (
    pkgNats     "retail-store-api/pkg/nats"
    pkgResponse "retail-store-api/pkg/response"
    es          "retail-store-api/pkg/elasticsearch"
    db          "retail-store-api/db/sqlc"
)
```

---

## 5. sqlc — Đặt tên query

Format: `-- name: <Verb><Entity>[By<Field>] :<return-type>`

| Operation    | Tên query                        | Return type |
| ------------ | -------------------------------- | ----------- |
| Đọc 1 record | `GetProductByID`                 | `:one`      |
| Đọc nhiều    | `ListProducts`, `ListCategories` | `:many`     |
| Đếm          | `CountProducts`                  | `:one`      |
| Tạo mới      | `CreateProduct`                  | `:one`      |
| Cập nhật     | `UpdateProduct`                  | `:one`      |
| Xoá          | `DeleteProduct`                  | `:exec`     |
| Xoá mềm      | `SoftDeleteProduct`              | `:exec`     |

Tên file `.sql`: 1 file = 1 bảng, snake_case — `product.sql`, `category.sql`, `user.sql`.

Migration format: `<sequence>_<verb>_<subject>.{up,down}.sql`

```
000001_init_schema.up.sql
000002_add_users.up.sql
```

---

## 6. Context & tham số

- `ctx context.Context` luôn là tham số đầu tiên.
- Không lưu `ctx` trong struct.
- Thứ tự: `ctx → ID → params/request → options`.

---

## 7. Test naming

- File: `<source>_test.go`.
- Function: `Test<Type>_<Method>_<Scenario>`.

```go
func TestProductService_Create_Success(t *testing.T) { ... }
func TestProductService_Create_DuplicateName(t *testing.T) { ... }
```

Table-driven test: tên case mô tả input → expectation.

---

## 8. Comment

- Exported symbol **phải** có doc comment bắt đầu bằng tên symbol.
- `TODO(name):` kèm tên người chịu trách nhiệm.
- Error logging: `slog.Error(...)` trước khi return error trong app initialization.

---

## Checklist review nhanh

- [ ] File snake_case, 1 concern / file
- [ ] Acronym giữ nguyên case (`userID`, `httpClient`)
- [ ] Không magic number/string
- [ ] Struct tag JSON/DB snake_case
- [ ] Interface ở consumer side, không prefix `I`
- [ ] Error: prefix `Err` (sentinel) hoặc suffix `Error` (type)
- [ ] sqlc query: `<Verb><Entity>[By<Field>]`
- [ ] `ctx context.Context` là tham số đầu tiên
- [ ] Exported symbol có doc comment
- [ ] Receiver nhất quán: `h`, `s`, `r`
