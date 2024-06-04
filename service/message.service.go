package service

import (
	"go-redis-kafka-streamer/configs"
	"go-redis-kafka-streamer/dto"
	"go-redis-kafka-streamer/dto/response"
	"go-redis-kafka-streamer/repository"
	"gorm.io/gorm"
	"time"
)

type Message interface {
	SaveMessage(header string, body string) (response.MessageResponse, error)
	RetrieveMessage(uuid string) (response.RetrieveResponse, error)
}

type MessageService struct {
	conf              *configs.Config
	redisService      *RedisService
	messageRepository *repository.MessageRepository
	postgresDB        *gorm.DB
}

func NewMessageService(config *configs.Config, redisService *RedisService, messageRepository *repository.MessageRepository, postgresDB *gorm.DB) *MessageService {
	return &MessageService{conf: config, redisService: redisService, messageRepository: messageRepository, postgresDB: postgresDB}
}

func (mess *MessageService) SaveMessage(header string, body string) (response.MessageResponse, error) {
	message := &dto.Message{Header: header, Body: body}

	tx := mess.getDbConnection()

	generatedUUID, err := mess.redisService.Save(message)
	message.Id = generatedUUID
	if err != nil {
		return response.MessageResponse{}, err
	}

	err = mess.messageRepository.Insert(tx, message)
	if err != nil {
		tx.Rollback()
		return response.MessageResponse{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
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

func (mess *MessageService) getDbConnection() *gorm.DB {
	tx := mess.postgresDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	return tx
}
