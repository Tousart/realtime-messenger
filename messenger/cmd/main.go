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

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// config
	cfg := config.LoadConfig()

	// messages handler repository
	msgsHandlerRepo, err := rabbitmq.NewRabbitMQMessagesHandlerRepository(cfg.RabbitMQ.Addr, cfg.RabbitMQ.MessagesQueue)
	if err != nil {
		log.Fatalf("failed to create publisher repository: %s", err.Error())
	}
	defer msgsHandlerRepo.Close()

	// queues repository
	queuesRepo := redis.NewRedisQueuesRepository(cfg.Redis.Addr, "messages")

	// messages handler service
	msgsHandlerService := service.NewMessagesHandlerService(msgsHandlerRepo, queuesRepo)

	// api methods router
	r := chi.NewRouter()

	// create server api
	srvApi := api.NewAPI(msgsHandlerService)
	srvApi.WithHandlers(r)
	srvApi.WithMethods()

	// create and run server
	srv := server.NewServer(cfg.Server.Addr, r)
	srv.CreateAndRunServer(ctx)

	srv.Wg.Wait()
}
