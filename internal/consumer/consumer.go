package consumer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"

	appnats "github.com/kainguyen/retail-store-api/pkg/nats"
)

type HandlerFunc func(msg jetstream.Msg)

type Subscription struct {
	Subject      string
	ConsumerName string
	Handler      HandlerFunc
}

type Consumer struct {
	js            jetstream.JetStream
	subscriptions []Subscription
}

func New(js jetstream.JetStream) *Consumer {
	return &Consumer{js: js}
}

func (c *Consumer) Register(sub Subscription) {
	c.subscriptions = append(c.subscriptions, sub)
}

func (c *Consumer) Start(ctx context.Context) error {
	_, err := appnats.EnsureStream(ctx, c.js)
	if err != nil {
		return err
	}

	for _, sub := range c.subscriptions {
		cons, err := c.js.CreateOrUpdateConsumer(ctx, appnats.StreamRetailStore, jetstream.ConsumerConfig{
			Durable:       sub.ConsumerName,
			FilterSubject: sub.Subject,
			AckPolicy:     jetstream.AckExplicitPolicy,
		})
		if err != nil {
			return fmt.Errorf("create consumer [%s]: %w", sub.ConsumerName, err)
		}

		handler := sub.Handler
		_, err = cons.Consume(func(msg jetstream.Msg) {
			handler(msg)
		})
		if err != nil {
			return fmt.Errorf("consume [%s]: %w", sub.ConsumerName, err)
		}

		slog.Info("Consumer listening", "name", sub.ConsumerName, "subject", sub.Subject)
	}

	<-ctx.Done()
	slog.Info("Consumer stopped")
	return nil
}
