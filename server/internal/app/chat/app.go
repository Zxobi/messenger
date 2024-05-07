package chat

import (
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/app/chat/grpc"
	"github.com/dvid-messanger/internal/client/frontend"
	"github.com/dvid-messanger/internal/config"
	"github.com/dvid-messanger/internal/lib/logger"
	"github.com/dvid-messanger/internal/service/chat"
	"github.com/dvid-messanger/internal/storage/chat/mongo"
	"github.com/dvid-messanger/internal/storage/chat/scylla"
	"github.com/dvid-messanger/pkg/database/mongodb"
	"github.com/dvid-messanger/pkg/database/scylladb"
	"github.com/gocql/gocql"
	"log/slog"
	"sync"
	"time"
)

const defaultStopTimeout = 10 * time.Second

type App struct {
	log         *slog.Logger
	grpcApp     *grpc.App
	chatStorage *mongo.Storage
}

func New(
	log *slog.Logger,
	cfg *config.ChatConfig,
	feClientCfg *config.ClientConfig,
) (*App, error) {
	const op = "chat.app.New"

	chatStorage := mongo.New(log, mongodb.Timeout(cfg.ChatStorage.Timeout), mongodb.URI(cfg.ChatStorage.ConnectUri))
	messageStorage := scylla.New(log)
	if err := messageStorage.Connect(
		scylladb.CreateCluster(
			gocql.Quorum,
			cfg.MessageStorage.Keyspace,
			cfg.MessageStorage.Hosts...,
		),
	); err != nil {
		log.Error("failed to connect to message storage", slog.String("op", op))
		return nil, err
	}
	notifier, err := frontend.New(context.TODO(), log, feClientCfg.Address, feClientCfg.Timeout, feClientCfg.RetriesCount)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chatService := chat.NewService(log, chatStorage, chatStorage, notifier, chatStorage, messageStorage, messageStorage)

	grpcApp := grpc.New(log, chatService, cfg.Port)

	return &App{
		log:         log,
		grpcApp:     grpcApp,
		chatStorage: chatStorage,
	}, nil
}

func (app *App) MustRun() {
	if err := app.chatStorage.Connect(context.TODO()); err != nil {
		panic(err)
	}
	app.grpcApp.MustRun()
}

func (app *App) Run() error {
	if err := app.chatStorage.Connect(context.TODO()); err != nil {
		return err
	}
	return app.grpcApp.Run()
}

func (app *App) Stop() {
	const op = "app.Stop"
	log := app.log.With(slog.String("op", op))

	var wg sync.WaitGroup
	wg.Add(2)

	ctx, cancel := context.WithTimeout(context.Background(), defaultStopTimeout)
	defer cancel()

	go func() {
		defer wg.Done()
		app.grpcApp.Stop()
	}()
	go func() {
		defer wg.Done()
		if err := app.chatStorage.Close(ctx); err != nil {
			log.Error("storage closed with error", logger.Err(err))
		}
	}()

	wg.Wait()
}
