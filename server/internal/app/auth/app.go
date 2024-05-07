package auth

import (
	"context"
	"github.com/dvid-messanger/internal/app/auth/grpc"
	"github.com/dvid-messanger/internal/config"
	"github.com/dvid-messanger/internal/core/service/auth"
	"github.com/dvid-messanger/internal/driver/secondary/storage/auth/mongo"
	"github.com/dvid-messanger/internal/pkg/jwt"
	"github.com/dvid-messanger/internal/pkg/logger"
	"github.com/dvid-messanger/pkg/database/mongodb"
	"log/slog"
	"sync"
	"time"
)

const defaultStopTimeout = 10 * time.Second

type App struct {
	log     *slog.Logger
	grpcApp *grpc.App
	storage *mongo.Storage
}

func New(log *slog.Logger, cfg *config.AuthConfig) *App {
	storage := mongo.New(log, mongodb.Timeout(cfg.Storage.Timeout), mongodb.URI(cfg.Storage.ConnectUri))
	tokenMaker := jwt.NewTokenizer([]byte(cfg.Secret))
	authService := auth.NewService(log, storage, storage, tokenMaker, cfg.TokenTTL)

	grpcApp := grpc.New(log, authService, cfg.Port)

	return &App{
		log:     log,
		grpcApp: grpcApp,
		storage: storage,
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
