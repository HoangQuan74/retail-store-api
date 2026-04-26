package internal

import (
	"github.com/elastic/go-elasticsearch/v8"
	db "github.com/hoangquan/retail-store-api/db/sqlc"
	"github.com/hoangquan/retail-store-api/pkg/auth"
	"github.com/hoangquan/retail-store-api/pkg/config"
)

type Dependencies struct {
	Config       *config.Config
	Queries      *db.Queries
	ESClient     *elasticsearch.Client
	ProductIndex string
	JWTManager   *auth.JWTManager
}
