package repository

import (
	"go-redis-kafka-streamer/dto"
	"go-redis-kafka-streamer/model"
	"gorm.io/gorm"
	"time"
)

type MessageRepositoryInterface interface {
	Insert(tx *gorm.DB, message *dto.Message) error
	Delete(tx *gorm.DB, uuid string) error
	Update(tx *gorm.DB, message *dto.Message) error
}

type MessageRepository struct{}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (repo *MessageRepository) Delete(tx *gorm.DB, uuid string) error {
	if err := tx.Where("message_id = ?", uuid).Delete(&model.Message{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MessageRepository) Update(tx *gorm.DB, message *dto.Message) error {
	savedMessage := tx.Model(model.Message{}).Where("message_id = ?", message.Id).Updates(
		model.Message{UpdatedAt: time.Now(), Header: message.Header, Body: message.Body})
	if savedMessage.Error != nil {
		return savedMessage.Error
	}
	return nil
}

func (repo *MessageRepository) Insert(tx *gorm.DB, message *dto.Message) error {
	nowTime := time.Now()
	messageEntity := model.Message{
		MessageId: message.Id,
		Body:      message.Body,
		Header:    message.Header,
		CreatedAt: nowTime,
	}

	savedMessageEntity := tx.Create(&messageEntity)
	if savedMessageEntity.Error != nil {
		return savedMessageEntity.Error
	}

	return nil
}
