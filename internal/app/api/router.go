package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kainguyen/retail-store-api/docs"
	"github.com/kainguyen/retail-store-api/internal/app"
	"github.com/kainguyen/retail-store-api/internal/handler"
	"github.com/kainguyen/retail-store-api/pkg/middleware"
)

func NewRouter(ctx *app.AppContext) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", handler.HealthCheck)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	handler.NewCategoryHandler(ctx, router)
	handler.NewProductHandler(ctx, router)
	handler.NewSearchHandler(ctx, router)

	return router
}
