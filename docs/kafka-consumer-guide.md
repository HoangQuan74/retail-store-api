## Tổng quan: Kafka vs NATS (đang dùng)

```
NATS JetStream (hiện tại):
  Producer → Stream (RETAIL_STORE) → Consumer (durable, filter by subject)

Kafka (sẽ build):
  Producer → Topic (product-events) → Consumer Group → Partition → Consumer
```

Khác biệt cốt lõi:

| Concept           | NATS JetStream                        | Kafka                                                     |
|-------------------|---------------------------------------|-----------------------------------------------------------|
| Nơi lưu message   | Stream                                | Topic (chia thành Partitions)                             |
| Đơn vị subscribe  | Subject filter                        | Topic + Partition                                         |
| Đảm bảo thứ tự    | Per subject                           | Per partition                                             |
| Fan-out           | Nhiều durable consumer cùng subject   | Nhiều consumer group cùng topic                           |
| Load balance      | Nhiều instance cùng durable name      | Nhiều consumer cùng group, mỗi consumer nhận 1+ partition |
| Offset management | NATS tự quản lý (ack/nak)             | Consumer tự commit offset                                 |

---

## Bước 1: Thêm Kafka vào infrastructure

### 1.1. Thêm vào `docker-compose.yml`

```yaml
# Thêm vào phần Infrastructure
kafka:
  image: confluentinc/cp-kafka:7.6.0
  ports:
    - "9092:9092"
  environment:
    KAFKA_NODE_ID: 1
    KAFKA_PROCESS_ROLES: broker,controller
    KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:29092,CONTROLLER://0.0.0.0:29093,EXTERNAL://0.0.0.0:9092
    KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,EXTERNAL://localhost:9092
    KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT
    KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
    KAFKA_CONTROLLER_QUORUM_VOTERS: 1@kafka:29093
    KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    CLUSTER_ID: "MkU3OEVBNTcwNTJENDM2Qg"
  volumes:
    - kafka_data:/var/lib/kafka/data
  restart: unless-stopped
```

> **Tại sao dùng KRaft mode (không có Zookeeper)?**
> Kafka 3.5+ hỗ trợ KRaft — metadata quản lý bởi chính Kafka, không cần Zookeeper.
> Đơn giản hơn cho dev, và đây là hướng đi chính thức của Kafka.

### 1.2. Thêm config

Trong `.env`:
```
KAFKA_BROKERS=localhost:9092
```

Trong `pkg/config/config.go`, thêm:
```go
type KafkaConfig struct {
    Brokers []string
}
```

Và parse trong `Load()`:
```go
Kafka: KafkaConfig{
    Brokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
},
```

### 1.3. Chạy thử

```bash
docker compose up -d kafka
# Đợi ~10s cho Kafka start xong
docker compose exec kafka kafka-topics --bootstrap-server localhost:29092 --list
```

---

## Bước 2: Tạo Kafka package trong `pkg/`

### 2.1. Producer — `pkg/kafka/producer.go`

```go
package kafka

import (
    "encoding/json"
    "fmt"

    "github.com/IBM/sarama"
)

type Producer struct {
    producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll    // đợi tất cả replica ack
    config.Producer.Retry.Max = 3                        // retry 3 lần nếu lỗi
    config.Producer.Return.Successes = true              // cần cho SyncProducer

    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, fmt.Errorf("create kafka producer: %w", err)
    }

    return &Producer{producer: producer}, nil
}

func (p *Producer) Publish(topic string, key string, data interface{}) error {
    bytes, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("marshal message: %w", err)
    }

    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.StringEncoder(key),       // key quyết định message vào partition nào
        Value: sarama.ByteEncoder(bytes),
    }

    _, _, err = p.producer.SendMessage(msg)
    return err
}

func (p *Producer) Close() error {
    return p.producer.Close()
}
```

> **Tại sao cần Key?**
> Kafka dùng key để hash → chọn partition.
> Cùng key = cùng partition = đảm bảo thứ tự.
> VD: key = productID → tất cả event của product 42 luôn vào cùng 1 partition → đảm bảo created → updated → deleted đúng thứ tự.

> **SyncProducer vs AsyncProducer:**
> - `SyncProducer`: block cho đến khi Kafka ack → chắc chắn message đã ghi, chậm hơn
> - `AsyncProducer`: fire-and-forget, nhanh hơn nhưng phải handle error qua channel
> Dùng Sync cho mutation events (product CRUD), Async cho high-throughput (analytics, logging).

### 2.2. Consumer — `pkg/kafka/consumer.go`

```go
package kafka

import (
    "context"
    "log/slog"

    "github.com/IBM/sarama"
)

// HandlerFunc xử lý 1 message
type HandlerFunc func(topic string, key, value []byte) error

// ConsumerGroup wraps sarama consumer group
type ConsumerGroup struct {
    group    sarama.ConsumerGroup
    topics   []string
    handlers map[string]HandlerFunc
}

func NewConsumerGroup(brokers []string, groupID string, topics []string) (*ConsumerGroup, error) {
    config := sarama.NewConfig()
    config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
        sarama.NewBalanceStrategyRoundRobin(),       // phân partition đều giữa consumer
    }
    config.Consumer.Offsets.Initial = sarama.OffsetNewest  // chỉ đọc message mới

    group, err := sarama.NewConsumerGroup(brokers, groupID, config)
    if err != nil {
        return nil, err
    }

    return &ConsumerGroup{
        group:    group,
        topics:   topics,
        handlers: make(map[string]HandlerFunc),
    }, nil
}

func (c *ConsumerGroup) RegisterHandler(topic string, handler HandlerFunc) {
    c.handlers[topic] = handler
}

// Start bắt đầu consume. Block cho đến khi ctx bị cancel.
func (c *ConsumerGroup) Start(ctx context.Context) error {
    handler := &groupHandler{handlers: c.handlers}

    for {
        // Consume trả về khi: rebalance xảy ra hoặc ctx cancel
        if err := c.group.Consume(ctx, c.topics, handler); err != nil {
            slog.Error("Consumer group error", "error", err)
        }

        if ctx.Err() != nil {
            return nil
        }
    }
}

func (c *ConsumerGroup) Close() error {
    return c.group.Close()
}

// groupHandler implements sarama.ConsumerGroupHandler
type groupHandler struct {
    handlers map[string]HandlerFunc
}

// Setup chạy khi consumer join group (sau rebalance)
func (h *groupHandler) Setup(sarama.ConsumerGroupSession) error {
    slog.Info("Consumer group setup — partitions assigned")
    return nil
}

// Cleanup chạy khi consumer rời group (trước rebalance)
func (h *groupHandler) Cleanup(sarama.ConsumerGroupSession) error {
    slog.Info("Consumer group cleanup — partitions revoked")
    return nil
}

// ConsumeClaim xử lý messages từ 1 partition cụ thể
func (h *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        handler, ok := h.handlers[msg.Topic]
        if !ok {
            slog.Warn("No handler for topic", "topic", msg.Topic)
            session.MarkMessage(msg, "")    // skip, commit offset
            continue
        }

        if err := handler(msg.Topic, msg.Key, msg.Value); err != nil {
            slog.Error("Handler failed",
                "topic", msg.Topic,
                "partition", msg.Partition,
                "offset", msg.Offset,
                "error", err,
            )
            // Không mark → message sẽ được xử lý lại sau rebalance
            // Hoặc mark anyway nếu muốn skip lỗi (tuỳ chiến lược)
            continue
        }

        // Commit offset — báo Kafka "tôi đã xử lý xong message này"
        session.MarkMessage(msg, "")
    }
    return nil
}
```

> **Consumer Group là gì?**
> Một "nhóm" consumer cùng đọc 1 topic.
> Kafka tự phân partition cho mỗi consumer trong group.
>
> ```
> Topic: product-events (3 partitions)
>
> Consumer Group "search-indexer":
>   consumer-1 → partition 0, 1
>   consumer-2 → partition 2
>
> Consumer Group "analytics":
>   consumer-1 → partition 0, 1, 2  (chỉ 1 instance)
> ```
>
> Mỗi partition chỉ được 1 consumer trong group đọc tại 1 thời điểm.
> → Thêm consumer = tự động rebalance partition.
> → Số consumer > số partition = consumer thừa ngồi chờ.

> **Offset là gì?**
> Vị trí đọc của consumer trong partition. Mỗi message có 1 offset (số tăng dần).
> `MarkMessage()` = commit offset = báo "đã xử lý đến đây".
> Consumer restart → đọc từ offset đã commit → không mất, không trùng (at-least-once).

> **Rebalance là gì?**
> Khi consumer join/leave group, Kafka phân lại partition.
> `Setup()` → nhận partition mới.
> `Cleanup()` → trả partition cũ.
> Trong lúc rebalance, consume bị pause ~vài giây.

---

## Bước 3: Tạo Kafka consumer service

### 3.1. Cấu trúc

```
services/
└── kafka-consumer/
    ├── cmd/main.go
    └── internal/
        ├── app.go
        └── handler/
            └── product_handler.go
```

### 3.2. `services/kafka-consumer/internal/app.go`

```go
package internal

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "your-module/pkg/config"
    "your-module/pkg/kafka"
    "your-module/pkg/logger"
    "your-module/services/kafka-consumer/internal/handler"
)

// Định nghĩa topic names
const (
    TopicProductEvents = "product-events"
)

type App struct {
    consumer *kafka.ConsumerGroup
}

func New() (*App, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, err
    }

    logger.New(cfg.Log)

    // Tạo consumer group
    consumer, err := kafka.NewConsumerGroup(
        cfg.Kafka.Brokers,
        "retail-kafka-consumer",     // group ID
        []string{TopicProductEvents},
    )
    if err != nil {
        slog.Error("Failed to create Kafka consumer group", "error", err)
        return nil, err
    }

    // Register handlers
    productHandler := handler.NewProductHandler()
    consumer.RegisterHandler(TopicProductEvents, productHandler.Handle)

    return &App{consumer: consumer}, nil
}

func (a *App) Start() error {
    defer a.consumer.Close()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        <-quit
        slog.Info("Shutting down kafka consumer...")
        cancel()
    }()

    slog.Info("Kafka consumer starting")
    return a.consumer.Start(ctx)
}
```

### 3.3. `services/kafka-consumer/internal/handler/product_handler.go`

```go
package handler

import (
    "encoding/json"
    "fmt"
    "log/slog"
)

type ProductEvent struct {
    Action string      `json:"action"`     // "created", "updated", "deleted"
    Data   ProductData `json:"data"`
}

type ProductData struct {
    ID          int64   `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Quantity    int32   `json:"quantity"`
    CategoryID  int64   `json:"category_id"`
}

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
    return &ProductHandler{}
}

func (h *ProductHandler) Handle(topic string, key, value []byte) error {
    var event ProductEvent
    if err := json.Unmarshal(value, &event); err != nil {
        return fmt.Errorf("unmarshal product event: %w", err)
    }

    switch event.Action {
    case "created":
        slog.Info("Product created", "id", event.Data.ID, "name", event.Data.Name)
        // TODO: index vào Elasticsearch
    case "updated":
        slog.Info("Product updated", "id", event.Data.ID, "name", event.Data.Name)
        // TODO: re-index
    case "deleted":
        slog.Info("Product deleted", "id", event.Data.ID)
        // TODO: xoá khỏi index
    default:
        slog.Warn("Unknown action", "action", event.Action)
    }

    return nil
}
```

---

## Bước 4: Producer side — publish event từ admin service

Trong `services/admin/internal/service/product_service.go`, thêm Kafka producer bên cạnh NATS:

```go
type ProductService struct {
    repo      *repository.ProductRepository
    publisher *pkgNats.Publisher
    kafka     *kafka.Producer         // thêm
}

func (s *ProductService) Create(ctx context.Context, req request.CreateProductRequest) (db.Product, error) {
    // ... tạo product ...

    // Publish qua NATS (giữ nguyên)
    s.publishEvent(ctx, pkgNats.SubjectProductCreated, ...)

    // Publish qua Kafka (thêm mới)
    s.kafka.Publish("product-events", fmt.Sprintf("%d", product.ID), map[string]interface{}{
        "action": "created",
        "data":   product,
    })
    //          ↑ key = productID
    //          → cùng product luôn vào cùng partition
    //          → đảm bảo thứ tự created → updated → deleted

    return product, nil
}
```

---

## Bước 5: Chạy & test

### 5.1. Start infrastructure

```bash
docker compose up -d kafka
```

### 5.2. Tạo topic (optional — Kafka auto-create, nhưng nên tạo tay để kiểm soát số partition)

```bash
docker compose exec kafka kafka-topics \
  --bootstrap-server localhost:29092 \
  --create \
  --topic product-events \
  --partitions 3 \
  --replication-factor 1
```

> **Chọn số partition:**
> - Partition = đơn vị song song. 3 partition → tối đa 3 consumer cùng group đọc song song.
> - Không thể giảm partition sau khi tạo (chỉ tăng được).
> - Rule of thumb: bắt đầu với 3-6, scale lên khi cần.

### 5.3. Start consumer

```bash
go run services/kafka-consumer/cmd/main.go
```

### 5.4. Test publish message thủ công

```bash
docker compose exec kafka kafka-console-producer \
  --bootstrap-server localhost:29092 \
  --topic product-events \
  --property "key.separator=:" \
  --property "parse.key=true"
```

Gõ vào:
```
42:{"action":"created","data":{"id":42,"name":"Test Product","price":99.99}}
```

Consumer sẽ log ra: `Product created id=42 name="Test Product"`

### 5.5. Xem consumer group status

```bash
docker compose exec kafka kafka-consumer-groups \
  --bootstrap-server localhost:29092 \
  --group retail-kafka-consumer \
  --describe
```

Output cho thấy: partition nào đang được consumer nào đọc, current offset, lag.

---

## Bước 6: Hiểu sâu — Những điều cần biết khi production

### 6.1. At-least-once vs Exactly-once

```
At-least-once (default, phổ biến nhất):
  Process message → Commit offset
  Nếu crash giữa process và commit → message được xử lý lại

Exactly-once (phức tạp):
  Cần Kafka transaction + idempotent producer
  Hoặc consumer tự đảm bảo idempotent (dùng message ID để dedup)
```

**Thực tế:** Hầu hết hệ thống dùng at-least-once + idempotent handler. Đơn giản và đủ tốt.

### 6.2. Dead Letter Topic

Message lỗi vĩnh viễn → publish sang dead letter topic thay vì retry mãi:

```go
func (h *groupHandler) ConsumeClaim(session, claim) error {
    for msg := range claim.Messages() {
        err := handler(msg)
        if err != nil {
            retryCount := getRetryCount(msg.Headers)
            if retryCount >= 3 {
                // Gửi sang dead letter topic
                producer.Publish("product-events.dlq", msg.Key, msg.Value)
                session.MarkMessage(msg, "")
                continue
            }
            // Publish lại với retry count + 1
            republishWithRetry(msg, retryCount+1)
        }
        session.MarkMessage(msg, "")
    }
}
```

### 6.3. Monitoring

Các metric quan trọng:
- **Consumer lag**: offset mới nhất - offset đã commit. Lag tăng = consumer không kịp xử lý.
- **Rebalance frequency**: rebalance thường xuyên = consumer không ổn định.
- **Processing time per message**: giúp dự đoán throughput.

---

## Tổng kết flow

```
[Admin Service]
    │
    ├── POST /products → ProductService.Create()
    │       │
    │       ├── Save to PostgreSQL
    │       ├── Publish to NATS (products.created)      ← consumer NATS hiện tại
    │       └── Publish to Kafka (product-events)        ← consumer Kafka mới
    │
    ▼
[Kafka Broker]
    │
    ├── Topic: product-events
    │     ├── Partition 0 ──→ [kafka-consumer instance 1]
    │     ├── Partition 1 ──→ [kafka-consumer instance 2]
    │     └── Partition 2 ──→ [kafka-consumer instance 1]  (2 partition vì chỉ có 2 instance)
    │
    ▼
[Kafka Consumer Service]
    │
    └── ProductHandler.Handle()
            ├── action: "created" → index ES
            ├── action: "updated" → re-index ES
            └── action: "deleted" → delete from ES
```
