package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/tousart/messenger/internal/api"
	"github.com/tousart/messenger/internal/repository/rabbitmq"
	"github.com/tousart/messenger/internal/repository/redis"
	"github.com/tousart/messenger/internal/server"
	"github.com/tousart/messenger/internal/usecase/service"
)

type Config struct {
	serverAddr    string
	rabbitMQAddr  string
	messagesQueue string
	redisAddr     string
	nodeAddr      string
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// environment variables

	cfg := Config{}

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("no env params, using os environments")
		cfg.serverAddr = os.Getenv("SERVER_ADDR")
		cfg.rabbitMQAddr = os.Getenv("RABBITMQ_ADDR_DOCKER")
		cfg.messagesQueue = os.Getenv("MESSAGES_QUEUE")
		cfg.redisAddr = os.Getenv("REDIS_ADDR_DOCKER")
		cfg.nodeAddr = os.Getenv("NODE_ADDR")
	} else {
		cfg.serverAddr = os.Getenv("SERVER_ADDR")
		cfg.rabbitMQAddr = os.Getenv("RABBITMQ_ADDR_LOCALHOST")
		cfg.messagesQueue = os.Getenv("MESSAGES_QUEUE")
		cfg.redisAddr = os.Getenv("REDIS_ADDR_LOCALHOST")
		cfg.nodeAddr = os.Getenv("NODE_ADDR")
	}

	// publisher repository

	publisherRepo, err := rabbitmq.NewRabbitMQPublisherRepository(cfg.rabbitMQAddr, cfg.messagesQueue)
	if err != nil {
		log.Fatalf("failed to create publisher repository: %s", err.Error())
	}
	defer publisherRepo.Close()

	// publisher service

	publisherService := service.NewRabbitMQPublisherService(publisherRepo)

	// set chat-nodes repository

	setRepo := redis.NewRedisNodesSetRepository(cfg.redisAddr, cfg.nodeAddr)

	// set chat-nodes repository

	setService := service.NewRedisChatNodesSetService(setRepo)

	// api methods router

	r := chi.NewRouter()

	// create server api

	srvApi := api.NewAPI(publisherService, setService)
	srvApi.WithHandlers(r)
	srvApi.WithMethods()

	// create and run server

	srv := server.NewServer(cfg.serverAddr, r)
	srv.CreateAndRunServer(ctx)

	srv.Wg.Wait()
}
