package service

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go-redis-kafka-streamer/dto"
)

type RedisServiceInterface interface {
	Save(message *dto.Message) (string, error)
	Fetch(key string) (*dto.Message, error)
}

type RedisService struct {
	redisClient *redis.Client
}

func NewRedisService(redisClient *redis.Client) *RedisService {
	return &RedisService{redisClient: redisClient}
}

func (red *RedisService) Save(message *dto.Message) (string, error) {
	jsonMessage, err := json.Marshal(message)

	ctx := context.Background()
	generatedUUID := generateUUid(message.Header, message.Body)
	err = red.redisClient.Set(ctx, generatedUUID, jsonMessage, 0).Err()
	if err != nil {
		return generatedUUID, err
	} else {
		return generatedUUID, nil
	}
}

func (red *RedisService) Fetch(key string) (*dto.Message, error) {
	ctx := context.Background()
	val, err := red.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var message dto.Message
	err = json.Unmarshal([]byte(val), &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func generateUUid(header string, body string) string {
	str := header + body
	calculatedUUid := uuid.NewSHA1(uuid.Nil, []byte(str))
	return calculatedUUid.String()
}
