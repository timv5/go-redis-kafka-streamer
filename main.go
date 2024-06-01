package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-redis-kafka-streamer/configs"
	"go-redis-kafka-streamer/handler"
	"go-redis-kafka-streamer/route"
	"go-redis-kafka-streamer/service"
	"log"
	"strconv"
)

var (
	server                 *gin.Engine
	MessageController      handler.MessageHandler
	MessageRouteController route.MessageRouteHandler
)

func main() {
	// set configs
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic("Could not initialize app")
	}

	// connect to database
	configs.ConnectToDB(&config)

	// Initialize a new Redis client
	redisDb, err := strconv.Atoi(config.RedisDb)
	if err != nil {
		panic("Could not initialize app, error converting redis db configs")
	}
	redisDatabase := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: "",
		DB:       redisDb,
	})

	// Ping the Redis server to check the connection
	ctx := context.Background()
	pong, err := redisDatabase.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}
	log.Println("Redis ping response:", pong)

	// initialize service
	redisService := service.NewRedisService(redisDatabase)
	messageService := service.NewMessageService(&config, redisService)

	// initialize controllers and routes
	MessageController = handler.NewMessageHandler(configs.DB, messageService, &config)
	MessageRouteController = route.NewMessageRouteHandler(MessageController)

	server = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	MessageRouteController.MessageRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
