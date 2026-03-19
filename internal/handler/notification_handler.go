package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kainguyen/retail-store-api/pkg/notification"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type NotificationHandler struct {
	hub *notification.Hub
}

func NewNotificationHandler(hub *notification.Hub) *NotificationHandler {
	return &NotificationHandler{hub: hub}
}

func (h *NotificationHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/ws/notifications", h.HandleWebSocket)
}

func (h *NotificationHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", "error", err)
		return
	}

	client := notification.NewClient(h.hub, conn)
	client.Register()

	go client.WritePump()
	go client.ReadPump()
}
