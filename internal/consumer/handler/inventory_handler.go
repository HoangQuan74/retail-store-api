package handler

import (
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

type InventoryHandler struct{}

func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{}
}

func (h *InventoryHandler) HandleOrderCreated(msg jetstream.Msg) {
	var event OrderCreatedEvent
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		slog.Error("[Inventory] Failed to unmarshal order event", "error", err)
		msg.Nak()
		return
	}

	slog.Info("[Inventory] Deducting stock for order", "order_id", event.OrderID)

	// TODO: trừ tồn kho trong DB

	msg.Ack()
}
