package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kainguyen/retail-store-api/internal/app"
	"github.com/kainguyen/retail-store-api/internal/handler"
	"github.com/kainguyen/retail-store-api/pkg/middleware"
)

func NewRouter(ctx *app.AppContext) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", handler.HealthCheck)

	handler.NewCategoryHandler(ctx, router)
	handler.NewProductHandler(ctx, router)
	handler.NewSearchHandler(ctx, router)

	return router
}
