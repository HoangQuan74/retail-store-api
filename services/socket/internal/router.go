package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangquan/retail-store-api/pkg/notification"
	"github.com/hoangquan/retail-store-api/services/socket/internal/handler"
)

func NewRouter(hub *notification.Hub) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", handler.HealthCheck)

	notifHandler := handler.NewNotificationHandler(hub)
	router.GET("/api/v1/ws/notifications", notifHandler.HandleWebSocket)

	return router
}
