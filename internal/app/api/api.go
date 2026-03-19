package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"

	db "github.com/kainguyen/retail-store-api/db/sqlc"
	"github.com/kainguyen/retail-store-api/internal/app"
	"github.com/kainguyen/retail-store-api/internal/config"
	"github.com/kainguyen/retail-store-api/pkg/database"
	es "github.com/kainguyen/retail-store-api/pkg/elasticsearch"
	"github.com/kainguyen/retail-store-api/pkg/logger"
	appnats "github.com/kainguyen/retail-store-api/pkg/nats"
)

type App struct {
	cfg    *config.Config
	pool   *pgxpool.Pool
	rdb    *redis.Client
	nc     *nats.Conn
	server *http.Server
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger.New(cfg.Log)

	pool, err := database.NewPostgres(cfg.DB)
	if err != nil {
		return nil, err
	}

	rdb, err := database.NewRedis(cfg.Redis)
	if err != nil {
		return nil, err
	}

	nc, err := appnats.Connect(cfg.NATS)
	if err != nil {
		return nil, err
	}

	js, err := appnats.NewJetStream(nc)
	if err != nil {
		return nil, err
	}

	if _, err := appnats.EnsureStream(context.Background(), js); err != nil {
		return nil, err
	}

	publisher := appnats.NewPublisher(js)

	esClient, err := es.NewClient(cfg.Elasticsearch)
	if err != nil {
		return nil, err
	}

	if err := es.EnsureProductIndex(esClient, cfg.Elasticsearch.ProductIndex); err != nil {
		return nil, err
	}

	queries := db.New(pool)
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	appCtx := &app.AppContext{
		Config:       cfg,
		Queries:      queries,
		ESClient:     esClient,
		ProductIndex: cfg.Elasticsearch.ProductIndex,
		Publisher:    publisher,
	}

	return &App{
		cfg:  cfg,
		pool: pool,
		rdb:  rdb,
		nc:   nc,
		server: &http.Server{
			Addr:    ":" + cfg.App.Port,
			Handler: NewRouter(appCtx),
		},
	}, nil
}

func (a *App) Start() error {
	defer a.pool.Close()
	defer a.rdb.Close()
	defer a.nc.Close()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("API server starting", "port", a.cfg.App.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-quit:
		return a.Shutdown()
	}
}

func (a *App) Shutdown() error {
	slog.Info("Shutting down API server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return a.server.Shutdown(ctx)
}
