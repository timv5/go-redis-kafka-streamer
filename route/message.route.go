package route

import (
	"github.com/gin-gonic/gin"
	"go-redis-kafka-streamer/handler"
)

type MessageRouteHandler struct {
	messageHandler handler.MessageHandler
}

func NewMessageRouteHandler(messageHandler handler.MessageHandler) MessageRouteHandler {
	return MessageRouteHandler{
		messageHandler: messageHandler,
	}
}

func (h *MessageRouteHandler) MessageRoute(group *gin.RouterGroup) {
	router := group.Group("message")
	router.POST("/send", h.messageHandler.SendMessage)
	router.POST("/retrieve", h.messageHandler.RetrieveMessage)
}
