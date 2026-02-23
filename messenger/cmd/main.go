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
	"golang.org/x/sync/errgroup"
)

func main() {
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	ewg, ctx := errgroup.WithContext(sigCtx)
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
	redisClient := pkgredis.NewClient(cfg.Redis.Addr)
	defer redisClient.Close()
	redisPubsub := redisClient.CreatePubSub(context.Background())

	/*

		Создание экземпляров repository, usecase, infrastructure

	*/

	// websocket manager
	wsManager := api.NewWebSocketManager()
	// messages handler repository
	msgsHandlerRepo := redis.NewRedisMessagesHandlerRepository(redisClient.Client(), redisPubsub)
	// messages handler service
	msgsHandlerService := service.NewMessagesHandlerService(wsManager, msgsHandlerRepo)
	// go consume messages
	msgsConsumer := infraredis.NewRedisConsumer(msgsHandlerService, redisPubsub)
	ewg.Go(func() error {
		return msgsConsumer.ConsumeMessages(ctx)
	})
	// users repository
	usersRepo, err := postgres.NewPSQLUsersRepository(psqlDB)
	if err != nil {
		log.Fatalf("failed to create publisher repository: %v", err)
	}
	// to hashing users password
	pswrdHasher := pkghashpassword.NewBCryptPasswordHasher()
	// sessions repository
	sessionsRepo := redis.NewRedisSessionsRepository(redisClient.Client())
	// users service
	usersService := service.NewUsersService(usersRepo, sessionsRepo, pswrdHasher)

	/*

		Создание экземпляра api и запуск сервера

	*/

	// api methods router
	r := chi.NewRouter()
	// create server api
	srvApi := api.NewAPI(wsManager, msgsHandlerService, usersService)
	isProd := true // isProd - boolean flag to local development (false if local else true)
	srvApi.WithHandlers(r, isProd)
	srvApi.WithMethods()
	// create and run server
	srv := server.NewServer(cfg.Server.Addr, r)
	ewg.Go(func() error {
		return srv.CreateAndRunServer(ctx)
	})

	ewg.Go(func() error {
		return srv.ShutdownServer(ctx)
	})

	if err := ewg.Wait(); err != nil {
		log.Printf("main error: %v\n", err)
		return
	}
}
