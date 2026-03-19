package socket

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/kainguyen/retail-store-api/internal/config"
	"github.com/kainguyen/retail-store-api/pkg/database"
	"github.com/kainguyen/retail-store-api/pkg/logger"
	"github.com/kainguyen/retail-store-api/pkg/notification"
)

type App struct {
	cfg       *config.Config
	rdb       *redis.Client
	server    *http.Server
	hub       *notification.Hub
	subCancel context.CancelFunc
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger.New(cfg.Log)

	rdb, err := database.NewRedis(cfg.Redis)
	if err != nil {
		return nil, err
	}

	hub := notification.NewHub()

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &App{
		cfg: cfg,
		rdb: rdb,
		hub: hub,
		server: &http.Server{
			Addr:    ":" + cfg.Socket.Port,
			Handler: NewRouter(hub),
		},
	}, nil
}

func (a *App) Start() error {
	defer a.rdb.Close()

	go a.hub.Run()

	subCtx, cancel := context.WithCancel(context.Background())
	a.subCancel = cancel
	sub := notification.NewSubscriber(a.rdb, a.hub)
	go sub.Start(subCtx)

	errCh := make(chan error, 1)
	go func() {
		slog.Info("Socket server starting", "port", a.cfg.Socket.Port)
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
	slog.Info("Shutting down Socket server...")
	a.subCancel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return a.server.Shutdown(ctx)
}
