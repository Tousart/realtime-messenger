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
	"github.com/tousart/messenger/internal/repository/postgres"
	"github.com/tousart/messenger/internal/repository/rabbitmq"
	"github.com/tousart/messenger/internal/repository/redis"
	"github.com/tousart/messenger/internal/server"
	"github.com/tousart/messenger/internal/usecase/service"
	pkghashpassword "github.com/tousart/messenger/pkg/hashpassword"
	pkgpostgres "github.com/tousart/messenger/pkg/postgres"
	pkgrabbitmq "github.com/tousart/messenger/pkg/rabbitmq"
	pkgredis "github.com/tousart/messenger/pkg/redis"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Load config
	cfg := config.LoadConfig()

	// Connections
	// Connect to PSQL
	psqlDB, err := pkgpostgres.ConnectToPSQL(cfg.PostgreSQL.Addr)
	if err != nil {
		log.Fatalf("failed to connect to psql")
	}
	// Create Redis-client
	redisClient := pkgredis.CreateRedisClient(cfg.Redis.Addr)
	// Connection to RabbitMQ
	rabbitMQConn, err := pkgrabbitmq.NewRabbitMQConnection(cfg.RabbitMQ.Addr, cfg.RabbitMQ.MessagesQueue)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq")
	}
	defer rabbitMQConn.Close()

	// messages handler repository
	msgsHandlerRepo, err := rabbitmq.NewRabbitMQMessagesHandlerRepository(rabbitMQConn.Channel(), rabbitMQConn.Queue(), cfg.RabbitMQ.MessagesQueue)
	if err != nil {
		log.Fatalf("failed to create publisher repository: %v", err)
	}

	// queues repository
	queuesRepo := redis.NewRedisQueuesRepository(redisClient, cfg.RabbitMQ.MessagesQueue) // messages queue use in rabbitmq too

	// messages handler service
	msgsHandlerService := service.NewMessagesHandlerService(msgsHandlerRepo, queuesRepo)

	// users repository
	usersRepo, err := postgres.NewPSQLUsersRepository(psqlDB)
	if err != nil {
		log.Fatalf("failed to create publisher repository: %v", err)
	}

	// Users service - create repositories and users service
	// to hashing users password
	pswrdHasher := pkghashpassword.NewBCryptPasswordHasher()
	// sessions repository
	sessionsRepo := redis.NewRedisSessionsRepository(redisClient)
	// users service
	usersService := service.NewUsersService(usersRepo, sessionsRepo, pswrdHasher)

	// api methods router
	r := chi.NewRouter()

	// create server api
	srvApi := api.NewAPI(msgsHandlerService, usersService)
	srvApi.WithHandlers(r)
	srvApi.WithMethods()

	// create and run server
	srv := server.NewServer(cfg.Server.Addr, r)
	srv.CreateAndRunServer(ctx)

	srv.Wg.Wait()
}
