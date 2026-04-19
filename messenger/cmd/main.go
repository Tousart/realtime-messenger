package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/messenger/configs"
	"github.com/tousart/messenger/internal/api/httpapi"
	"github.com/tousart/messenger/internal/api/wsapi"
	infraredis "github.com/tousart/messenger/internal/infrastructure/redis"
	"github.com/tousart/messenger/internal/repository/postgresql"
	"github.com/tousart/messenger/internal/repository/redis"
	"github.com/tousart/messenger/internal/server"
	"github.com/tousart/messenger/internal/usecase"
	"github.com/tousart/messenger/pkg/generator"
	"github.com/tousart/messenger/pkg/hashpassword"
	"github.com/tousart/messenger/pkg/logger"
	pkgpsql "github.com/tousart/messenger/pkg/postgresql"
	pkgredis "github.com/tousart/messenger/pkg/redis"
	"golang.org/x/sync/errgroup"
)

func main() {
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	ewg, ctx := errgroup.WithContext(sigCtx)

	// configs

	flags := configs.ParseFlags()
	cfg, err := configs.LoadConfig(flags.CfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

	// logger

	logger := logger.InitLogger()

	// connect to postgresql

	db, err := pkgpsql.ConnectToPSQL(
		cfg.PSQL.User, cfg.PSQL.Password, cfg.PSQL.Host, cfg.PSQL.DB, cfg.PSQL.SSLMode, cfg.PSQL.Port)
	if err != nil {
		log.Fatalf("failed connect to psql: %v\n", err)
	}

	// create redis client

	redisClient := pkgredis.NewClient(cfg.Redis.Host, cfg.Redis.Password, cfg.Redis.Port)
	defer redisClient.Close()
	redisPubsub := redisClient.CreatePubSub(ctx)

	// password hasher

	passwordHasher := hashpassword.NewBCryptPasswordHasher()

	// id generator

	idGen := generator.NewGenerator()

	// repository

	chatPub := redis.NewChatPublisher(redisClient.Client(), redisPubsub)

	sessionsRepo := redis.NewSessionsRepository(redisClient.Client())

	usersRepo := postgresql.NewUsersRepository(db)

	msgsRepo := postgresql.NewMessagesRepository(db)

	// usecase

	msgsUC := usecase.NewMessagesUsecase(msgsRepo, chatPub, idGen)

	usersUC := usecase.NewUsersService(usersRepo, sessionsRepo, passwordHasher, idGen)

	// websocket manager

	wsManager := wsapi.NewWebSocketManager(msgsUC, logger)
	wsManager.WithMethods()

	// go consume messages

	msgsConsumer := infraredis.NewRedisConsumer(wsManager, redisPubsub)
	ewg.Go(func() error {
		return msgsConsumer.ConsumeMessages(ctx)
	})

	// api

	r := chi.NewRouter()

	srvApi := httpapi.NewAPI(wsManager, msgsUC, usersUC, logger)
	srvApi.WithHandlers(r)

	// create and run server

	srv := server.NewServer(cfg.Server.Host, cfg.Server.Port, r, logger)
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
