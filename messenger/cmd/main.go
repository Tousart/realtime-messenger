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
	infraredis "github.com/tousart/messenger/internal/infrastructure/redis"
	"github.com/tousart/messenger/internal/repository/postgres"
	"github.com/tousart/messenger/internal/repository/redis"
	"github.com/tousart/messenger/internal/server"
	"github.com/tousart/messenger/internal/usecase/service"
	pkghashpassword "github.com/tousart/messenger/pkg/hashpassword"
	pkgpostgres "github.com/tousart/messenger/pkg/postgres"
	pkgredis "github.com/tousart/messenger/pkg/redis"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Load config
	cfg := config.LoadConfig()

	/*

		Подключение ко внешним инструментам

	*/

	// Connect to PSQL
	psqlDB, err := pkgpostgres.ConnectToPSQL(cfg.PostgreSQL.Addr)
	if err != nil {
		log.Fatalf("failed to connect to psql: %v\n", err)
	}
	// Create Redis-client
	redisClient := pkgredis.CreateRedisClient(cfg.Redis.Addr)
	defer redisClient.Close()
	redisPubsub := pkgredis.CreateRedisPubSubObject(context.Background(), redisClient)
	defer redisPubsub.Close()

	/*

		Создание экземпляров repository, usecase, infrastructure

	*/

	// websocket manager
	wsManager := api.NewWebSocketManager()

	// messages handler repository
	msgsHandlerRepo := redis.NewRedisMessagesHandlerRepository(redisClient, redisPubsub)

	// messages handler service
	msgsHandlerService := service.NewMessagesHandlerService(wsManager, msgsHandlerRepo)

	// go consume messages
	msgsConsumer := infraredis.NewRedisConsumer(msgsHandlerService, redisPubsub)
	go msgsConsumer.ConsumeMessages(ctx)

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

	/*

		Создание экземпляра api и запуск сервера

	*/

	// api methods router
	r := chi.NewRouter()

	// create server api
	srvApi := api.NewAPI(wsManager, msgsHandlerService, usersService)
	srvApi.WithHandlers(r)
	srvApi.WithMethods()

	// create and run server
	srv := server.NewServer(cfg.Server.Addr, r)
	srv.CreateAndRunServer(ctx)

	srv.Wg.Wait()
}
