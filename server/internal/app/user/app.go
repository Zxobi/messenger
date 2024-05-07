package user

import (
	"context"
	"github.com/dvid-messanger/internal/app/user/grpc"
	"github.com/dvid-messanger/internal/config"
	"github.com/dvid-messanger/internal/core/service/user"
	"github.com/dvid-messanger/internal/driver/secondary/storage/user/mongo"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/dvid-messanger/pkg/database/mongodb"
	"log/slog"
	"sync"
	"time"
)

const defaultStopTimeout = 10 * time.Second

type App struct {
	log     *slog.Logger
	storage *mongo.Storage
	grpcApp *grpc.App
}

func New(log *slog.Logger, cfg *config.UserConfig) *App {
	storage := mongo.New(log, mongodb.Timeout(cfg.Storage.Timeout), mongodb.URI(cfg.Storage.ConnectUri))
	userService := user.NewUser(log, storage, storage)
	grpcApp := grpc.New(log, userService, cfg.Port)

	return &App{
		log:     log,
		storage: storage,
		grpcApp: grpcApp,
	}
}

func (app *App) MustRun() {
	if err := app.storage.Connect(context.TODO()); err != nil {
		panic(err)
	}
	app.grpcApp.MustRun()
}

func (app *App) Run() error {
	if err := app.storage.Connect(context.TODO()); err != nil {
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
		if err := app.storage.Close(ctx); err != nil {
			log.Error("storage closed with error", logger.Err(err))
		}
	}()

	wg.Wait()
}
