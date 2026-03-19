package consumer

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"

	"github.com/kainguyen/retail-store-api/internal/config"
	appConsumer "github.com/kainguyen/retail-store-api/internal/consumer"
	"github.com/kainguyen/retail-store-api/internal/consumer/handler"
	es "github.com/kainguyen/retail-store-api/pkg/elasticsearch"
	"github.com/kainguyen/retail-store-api/pkg/logger"
	appnats "github.com/kainguyen/retail-store-api/pkg/nats"
)

type App struct {
	nc       *nats.Conn
	consumer *appConsumer.Consumer
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger.New(cfg.Log)

	nc, err := appnats.Connect(cfg.NATS)
	if err != nil {
		return nil, err
	}

	js, err := appnats.NewJetStream(nc)
	if err != nil {
		return nil, err
	}

	esClient, err := es.NewClient(cfg.Elasticsearch)
	if err != nil {
		return nil, err
	}

	c := appConsumer.New(js)

	// Register handlers
	analyticsHandler := handler.NewAnalyticsHandler()
	inventoryHandler := handler.NewInventoryHandler()
	searchIndexHandler := handler.NewSearchIndexHandler(esClient, cfg.Elasticsearch.ProductIndex)

	c.Register(appConsumer.Subscription{
		Subject:      appnats.SubjectOrderCreated,
		ConsumerName: "analytics-order-created",
		Handler:      analyticsHandler.HandleOrderCreated,
	})
	c.Register(appConsumer.Subscription{
		Subject:      appnats.SubjectProductViewed,
		ConsumerName: "analytics-product-viewed",
		Handler:      analyticsHandler.HandleProductViewed,
	})
	c.Register(appConsumer.Subscription{
		Subject:      appnats.SubjectOrderCreated,
		ConsumerName: "inventory-order-created",
		Handler:      inventoryHandler.HandleOrderCreated,
	})
	c.Register(appConsumer.Subscription{
		Subject:      appnats.SubjectProductCreated,
		ConsumerName: "search-product-created",
		Handler:      searchIndexHandler.HandleProductCreated,
	})
	c.Register(appConsumer.Subscription{
		Subject:      appnats.SubjectProductUpdated,
		ConsumerName: "search-product-updated",
		Handler:      searchIndexHandler.HandleProductUpdated,
	})
	c.Register(appConsumer.Subscription{
		Subject:      appnats.SubjectProductDeleted,
		ConsumerName: "search-product-deleted",
		Handler:      searchIndexHandler.HandleProductDeleted,
	})

	return &App{nc: nc, consumer: c}, nil
}

func (a *App) Start() error {
	defer a.nc.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		slog.Info("Shutting down consumer...")
		cancel()
	}()

	return a.consumer.Start(ctx)
}
