package handler

import (
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

type AnalyticsHandler struct{}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{}
}

type OrderCreatedEvent struct {
	OrderID    int64   `json:"order_id"`
	CustomerID int64   `json:"customer_id"`
	Total      float64 `json:"total"`
}

type ProductViewedEvent struct {
	ProductID  int64 `json:"product_id"`
	CustomerID int64 `json:"customer_id"`
}

func (h *AnalyticsHandler) HandleOrderCreated(msg jetstream.Msg) {
	var event OrderCreatedEvent
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		slog.Error("[Analytics] Failed to unmarshal order event", "error", err)
		msg.Nak()
		return
	}

	slog.Info("[Analytics] Order created", "order_id", event.OrderID, "customer_id", event.CustomerID, "total", event.Total)
	msg.Ack()
}

func (h *AnalyticsHandler) HandleProductViewed(msg jetstream.Msg) {
	var event ProductViewedEvent
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		slog.Error("[Analytics] Failed to unmarshal product viewed event", "error", err)
		msg.Nak()
		return
	}

	slog.Info("[Analytics] Product viewed", "product_id", event.ProductID, "customer_id", event.CustomerID)
	msg.Ack()
}
