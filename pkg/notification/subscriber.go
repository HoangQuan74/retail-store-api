package notification

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

var Channels = []string{
	"notifications:promotion",
	"notifications:discount",
}

type Subscriber struct {
	rdb *redis.Client
	hub *Hub
}

func NewSubscriber(rdb *redis.Client, hub *Hub) *Subscriber {
	return &Subscriber{rdb: rdb, hub: hub}
}

func (s *Subscriber) Start(ctx context.Context) {
	pubsub := s.rdb.Subscribe(ctx, Channels...)
	defer pubsub.Close()

	slog.Info("Subscribed to Redis channels", "channels", Channels)

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			slog.Info("Subscriber stopped")
			return
		case msg := <-ch:
			if msg == nil {
				return
			}
			slog.Info("Received event", "channel", msg.Channel)
			s.hub.Broadcast <- []byte(msg.Payload)
		}
	}
}
