package app

import (
	"github.com/elastic/go-elasticsearch/v8"
	db "github.com/hoangquan/retail-store-api/db/sqlc"
	"github.com/hoangquan/retail-store-api/internal/config"
	"github.com/hoangquan/retail-store-api/pkg/auth"
	pkgNats "github.com/hoangquan/retail-store-api/pkg/nats"
	"github.com/hoangquan/retail-store-api/pkg/notification"
)

type AppContext struct {
	Config       *config.Config
	Queries      *db.Queries
	ESClient     *elasticsearch.Client
	ProductIndex string
	Publisher    *pkgNats.Publisher
	Hub          *notification.Hub
	JWTManager   *auth.JWTManager
}
