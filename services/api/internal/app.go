package internal

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

	db "github.com/hoangquan/retail-store-api/db/sqlc"
	"github.com/hoangquan/retail-store-api/pkg/auth"
	"github.com/hoangquan/retail-store-api/pkg/config"
	"github.com/hoangquan/retail-store-api/pkg/database"
	es "github.com/hoangquan/retail-store-api/pkg/elasticsearch"
	"github.com/hoangquan/retail-store-api/pkg/logger"
	pkgNats "github.com/hoangquan/retail-store-api/pkg/nats"
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
		slog.Error("Failed to load config", "error", err)
		return nil, err
	}

	logger.New(cfg.Log)

	pool, err := database.NewPostgres(cfg.DB)
	if err != nil {
		slog.Error("Failed to connect to PostgreSQL", "error", err)
		return nil, err
	}

	rdb, err := database.NewRedis(cfg.Redis)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		return nil, err
	}

	nc, err := pkgNats.Connect(cfg.NATS)
	if err != nil {
		slog.Error("Failed to connect to NATS", "error", err)
		return nil, err
	}

	js, err := pkgNats.NewJetStream(nc)
	if err != nil {
		slog.Error("Failed to create JetStream", "error", err)
		return nil, err
	}

	if _, err := pkgNats.EnsureStream(context.Background(), js); err != nil {
		slog.Error("Failed to ensure NATS stream", "error", err)
		return nil, err
	}

	esClient, err := es.NewClient(cfg.Elasticsearch)
	if err != nil {
		slog.Error("Failed to connect to Elasticsearch", "error", err)
		return nil, err
	}

	if err := es.EnsureProductIndex(esClient, cfg.Elasticsearch.ProductIndex); err != nil {
		slog.Error("Failed to ensure Elasticsearch product index", "error", err)
		return nil, err
	}

	queries := db.New(pool)
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)

	deps := &Dependencies{
		Config:       cfg,
		Queries:      queries,
		ESClient:     esClient,
		ProductIndex: cfg.Elasticsearch.ProductIndex,
		JWTManager:   jwtManager,
	}

	return &App{
		cfg:  cfg,
		pool: pool,
		rdb:  rdb,
		nc:   nc,
		server: &http.Server{
			Addr:    ":" + cfg.App.Port,
			Handler: NewRouter(deps),
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
			slog.Error("API server failed", "error", err)
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
