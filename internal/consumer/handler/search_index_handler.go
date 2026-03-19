package handler

import (
	"encoding/json"
	"log/slog"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/nats-io/nats.go/jetstream"
	es "github.com/kainguyen/retail-store-api/pkg/elasticsearch"
)

type SearchIndexHandler struct {
	client    *elasticsearch.Client
	indexName string
}

func NewSearchIndexHandler(client *elasticsearch.Client, indexName string) *SearchIndexHandler {
	return &SearchIndexHandler{client: client, indexName: indexName}
}

func (h *SearchIndexHandler) HandleProductCreated(msg jetstream.Msg) {
	var doc es.ProductDocument
	if err := json.Unmarshal(msg.Data(), &doc); err != nil {
		slog.Error("[SearchIndex] Failed to unmarshal product created", "error", err)
		msg.Nak()
		return
	}

	if err := es.IndexProduct(h.client, h.indexName, doc); err != nil {
		slog.Error("[SearchIndex] Failed to index product", "id", doc.ID, "error", err)
		msg.Nak()
		return
	}

	slog.Info("[SearchIndex] Product indexed", "id", doc.ID, "name", doc.Name)
	msg.Ack()
}

func (h *SearchIndexHandler) HandleProductUpdated(msg jetstream.Msg) {
	var doc es.ProductDocument
	if err := json.Unmarshal(msg.Data(), &doc); err != nil {
		slog.Error("[SearchIndex] Failed to unmarshal product updated", "error", err)
		msg.Nak()
		return
	}

	if err := es.IndexProduct(h.client, h.indexName, doc); err != nil {
		slog.Error("[SearchIndex] Failed to re-index product", "id", doc.ID, "error", err)
		msg.Nak()
		return
	}

	slog.Info("[SearchIndex] Product re-indexed", "id", doc.ID, "name", doc.Name)
	msg.Ack()
}

func (h *SearchIndexHandler) HandleProductDeleted(msg jetstream.Msg) {
	var event struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		slog.Error("[SearchIndex] Failed to unmarshal product deleted", "error", err)
		msg.Nak()
		return
	}

	if err := es.DeleteProduct(h.client, h.indexName, event.ID); err != nil {
		slog.Error("[SearchIndex] Failed to delete product from index", "id", event.ID, "error", err)
		msg.Nak()
		return
	}

	slog.Info("[SearchIndex] Product removed from index", "id", event.ID)
	msg.Ack()
}
