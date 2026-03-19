package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/kainguyen/retail-store-api/internal/handler"
	"github.com/kainguyen/retail-store-api/pkg/notification"
)

func NewRouter(hub *notification.Hub) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", handler.HealthCheck)

	handler.NewNotificationHandler(hub).RegisterRoutes(router)

	return router
}
