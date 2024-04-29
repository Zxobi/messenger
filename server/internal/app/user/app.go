package user

import (
	"context"
	"github.com/dvid-messanger/internal/app/user/grpc"
	"github.com/dvid-messanger/internal/lib/logger"
	"github.com/dvid-messanger/internal/service/user"
	"github.com/dvid-messanger/internal/storage/user/mongo"
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

func New(
	log *slog.Logger,
	grpcPort int,
	mongoTimeout time.Duration,
	mongoURI string,
	dbName string,
	usersColName string,
) *App {
	storage := mongo.New(log, dbName, usersColName, mongodb.Timeout(mongoTimeout), mongodb.URI(mongoURI))
	userService := user.NewUser(log, storage, storage)
	grpcApp := grpc.New(log, userService, grpcPort)

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
	if err := app.grpcApp.Run(); err != nil {
		return err
	}
	return nil
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
