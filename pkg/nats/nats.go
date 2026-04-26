package nats

import (
	"fmt"
	"log/slog"

	"github.com/hoangquan/retail-store-api/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func Connect(cfg config.NATSConfig) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("connect nats: %w", err)
	}
	slog.Info("Connected to NATS successfully")
	return nc, nil
}

func NewJetStream(nc *nats.Conn) (jetstream.JetStream, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("create jetstream: %w", err)
	}
	slog.Info("JetStream initialized successfully")
	return js, nil
}
