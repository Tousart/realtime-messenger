package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/messenger/config"
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

	// config

	cfg := config.LoadConfig()

	// publisher repository

	publisherRepo, err := rabbitmq.NewRabbitMQPublisherRepository(cfg.RabbitMQ.Addr, cfg.RabbitMQ.MessagesQueue)
	if err != nil {
		log.Fatalf("failed to create publisher repository: %s", err.Error())
	}
	defer publisherRepo.Close()

	// publisher service

	publisherService := service.NewMessagesPublisherService(publisherRepo)

	// set chat-nodes repository

	receiverRepo := redis.NewRedisNodesReceiverRepository(cfg.Redis.Addr, cfg.Server.NodeAddr)

	// set chat-nodes repository

	receiverService := service.NewNodesReceiverService(receiverRepo)

	// api methods router

	r := chi.NewRouter()

	// create server api

	srvApi := api.NewAPI(publisherService, receiverService)
	srvApi.WithHandlers(r)
	srvApi.WithMethods()

	// create and run server

	srv := server.NewServer(cfg.Server.Addr, r)
	srv.CreateAndRunServer(ctx)

	srv.Wg.Wait()
}
