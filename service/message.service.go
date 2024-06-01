package service

import (
	"go-redis-kafka-streamer/configs"
	"go-redis-kafka-streamer/dto"
	"go-redis-kafka-streamer/dto/response"
	"time"
)

type Message interface {
	SaveMessage(header string, body string) (response.MessageResponse, error)
	RetrieveMessage(uuid string) (response.RetrieveResponse, error)
}

type MessageService struct {
	conf         *configs.Config
	redisService *RedisService
}

func NewMessageService(config *configs.Config, redisService *RedisService) *MessageService {
	return &MessageService{conf: config, redisService: redisService}
}

func (mess *MessageService) SaveMessage(header string, body string) (response.MessageResponse, error) {
	message := &dto.Message{Header: header, Body: body}
	generatedUUID, err := mess.redisService.Save(message)
	if err != nil {
		return response.MessageResponse{}, err
	}

	return response.MessageResponse{
		ID:        generatedUUID,
		CreatedAt: time.Now(),
		Header:    header,
		Body:      body,
	}, nil
}

func (mess *MessageService) RetrieveMessage(uuid string) (response.RetrieveResponse, error) {
	message, err := mess.redisService.Fetch(uuid)
	if err != nil {
		return response.RetrieveResponse{}, err
	}

	return response.RetrieveResponse{
		ID:     uuid,
		Header: message.Header,
		Body:   message.Body,
	}, nil
}
