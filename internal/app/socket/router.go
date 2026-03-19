package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/kainguyen/retail-store-api/internal/app"
	"github.com/kainguyen/retail-store-api/internal/handler"
)

func NewRouter(ctx *app.AppContext) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", handler.HealthCheck)

	handler.NewNotificationHandler(ctx, router)

	return router
}
