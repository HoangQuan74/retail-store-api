package app

import (
	"github.com/elastic/go-elasticsearch/v8"
	db "github.com/kainguyen/retail-store-api/db/sqlc"
	"github.com/kainguyen/retail-store-api/internal/config"
	appnats "github.com/kainguyen/retail-store-api/pkg/nats"
	"github.com/kainguyen/retail-store-api/pkg/notification"
)

type AppContext struct {
	Config       *config.Config
	Queries      *db.Queries
	ESClient     *elasticsearch.Client
	ProductIndex string
	Publisher    *appnats.Publisher
	Hub          *notification.Hub
}
