package handler

import (
	"go-redis-kafka-streamer/configs"
	"go-redis-kafka-streamer/dto/request"
	"go-redis-kafka-streamer/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessageHandler struct {
	postgresDB     *gorm.DB
	messageService *service.MessageService
	config         *configs.Config
}

func NewMessageHandler(postgresDB *gorm.DB, messageService *service.MessageService, config *configs.Config) MessageHandler {
	return MessageHandler{
		postgresDB:     postgresDB,
		messageService: messageService,
		config:         config,
	}
}

func (messageHandler MessageHandler) RetrieveMessage(ctx *gin.Context) {
	var retrievePayload *request.RetrieveRequest
	if err := ctx.ShouldBindJSON(&retrievePayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if retrievePayload.UUID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "wrong request params"})
		return
	}

	retrieveResponse, err := messageHandler.messageService.RetrieveMessage(retrievePayload.UUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, retrieveResponse)
	}
}

func (messageHandler MessageHandler) SendMessage(ctx *gin.Context) {
	var messagePayload *request.MessageRequest
	if err := ctx.ShouldBindJSON(&messagePayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if messagePayload.Body == "" || messagePayload.Header == "" || messagePayload.UUID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "wrong request params"})
		return
	}

	messageResponse, err := messageHandler.messageService.SaveMessage(messagePayload.UUID, messagePayload.Header, messagePayload.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status:": "error", "message": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, messageResponse)
	}
}
