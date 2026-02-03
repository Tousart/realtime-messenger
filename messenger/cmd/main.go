package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/messenger/internal/api"
	"github.com/tousart/messenger/internal/repository/rabbitmq"
	"github.com/tousart/messenger/internal/server"
	"github.com/tousart/messenger/internal/usecase/service"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// environment variables

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("failed to load env params")
	}

	// publisher repository

	publisherRepo, err := rabbitmq.NewRabbitMQPublisherRepository(os.Getenv("RABBITMQ_ADDR"), os.Getenv("MESSAGES_QUEUE"))
	if err != nil {
		log.Fatal("failed to create publisher repository")
	}
	defer publisherRepo.Close()

	// publisher service

	publisherService := service.NewRabbitMQPublisherService(publisherRepo)

	// api methods router

	r := chi.NewRouter()

	// create server api

	srvApi := api.NewAPI(publisherService)
	srvApi.WithHandlers(r)
	srvApi.WithMethods()

	// create and run server

	srv := server.NewServer(os.Getenv("SERVER_ADDR"), r)
	srv.CreateAndRunServer(ctx)

	srv.Wg.Wait()
}
