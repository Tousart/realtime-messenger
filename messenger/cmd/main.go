package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/messenger/config"
	httpapi "github.com/tousart/messenger/internal/api/http"
	wsapi "github.com/tousart/messenger/internal/api/websocket"
	infraredis "github.com/tousart/messenger/internal/infrastructure/redis"
	"github.com/tousart/messenger/internal/repository/postgres"
	"github.com/tousart/messenger/internal/repository/redis"
	"github.com/tousart/messenger/internal/server"
	"github.com/tousart/messenger/internal/usecase"
	pkghashpassword "github.com/tousart/messenger/pkg/hashpassword"
	pkglogger "github.com/tousart/messenger/pkg/logger"
	pkgpostgres "github.com/tousart/messenger/pkg/postgres"
	pkgredis "github.com/tousart/messenger/pkg/redis"
	"golang.org/x/sync/errgroup"
)

func main() {
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	ewg, ctx := errgroup.WithContext(sigCtx)

	// config
	cfg := config.LoadConfig()

	// logger
	logger := pkglogger.InitLogger()

	/*

		Подключение ко внешним инструментам

	*/

	// Connect to PSQL
	db, err := pkgpostgres.ConnectToPSQL(cfg.PostgreSQL.Addr)
	if err != nil {
		log.Fatalf("failed to connect to psql: %v\n", err)
	}
	// Create Redis-client
	redisClient := pkgredis.NewClient(cfg.Redis.Addr)
	defer redisClient.Close()
	redisPubsub := redisClient.CreatePubSub(context.Background())

	/*

		Создание экземпляров repository, usecase

	*/
	// to hashing users password
	pswrdHasher := pkghashpassword.NewBCryptPasswordHasher()
	// repository
	messagesRepo := redis.NewRedisMessagesHandlerRepository(redisClient.Client(), redisPubsub)
	sessionsRepo := redis.NewRedisSessionsRepository(redisClient.Client())
	usersRepo, err := postgres.NewPSQLUsersRepository(db)
	if err != nil {
		log.Fatalf("failed to create publisher repository: %v", err)
	}
	chatsRepo := postgres.NewChatsRepository(db)

	// messages handler service
	messagesUC := usecase.NewMessagesUsecase(messagesRepo, chatsRepo)

	// users service
	usersService := usecase.NewUsersService(usersRepo, sessionsRepo, pswrdHasher)

	/*

		Создание экземпляров api, infrastructure и запуск сервера

	*/

	// websocket manager
	wsManager := wsapi.NewWebSocketManager(messagesUC, logger)

	// go consume messages
	msgsConsumer := infraredis.NewRedisConsumer(wsManager, redisPubsub)
	ewg.Go(func() error {
		return msgsConsumer.ConsumeMessages(ctx)
	})

	// api methods router
	r := chi.NewRouter()
	// create server api
	srvApi := httpapi.NewAPI(wsManager, messagesUC, usersService, logger)
	isProd := true // isProd - boolean flag to local development (false if local else true)
	srvApi.WithHandlers(r, isProd)
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
