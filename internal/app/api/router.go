package api

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	db "github.com/kainguyen/retail-store-api/db/sqlc"
	"github.com/kainguyen/retail-store-api/internal/handler"
	"github.com/kainguyen/retail-store-api/pkg/middleware"
	appnats "github.com/kainguyen/retail-store-api/pkg/nats"
)

func NewRouter(queries *db.Queries, esClient *elasticsearch.Client, productIndex string, publisher *appnats.Publisher) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", handler.HealthCheck)

	handler.NewCategoryHandler(queries).RegisterRoutes(router)
	handler.NewProductHandler(queries, publisher).RegisterRoutes(router)
	handler.NewSearchHandler(esClient, productIndex).RegisterRoutes(router)

	return router
}
