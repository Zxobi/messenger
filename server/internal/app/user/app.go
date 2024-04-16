package user

import (
	"github.com/dvid-messanger/internal/app/user/grpc"
	"github.com/dvid-messanger/internal/service/user"
	"github.com/dvid-messanger/internal/storage/user/inmem"
	"log/slog"
)

type App struct {
	GrpcApp *grpc.App
}

func New(log *slog.Logger, grpcPort int) *App {
	storage := inmem.New()
	userService := user.NewUser(log, storage, storage)
	grpcApp := grpc.New(log, userService, grpcPort)

	return &App{
		GrpcApp: grpcApp,
	}
}

func (app *App) MustRun() {
	app.GrpcApp.MustRun()
}

func (app *App) Run() error {
	return app.GrpcApp.Run()
}

func (app *App) Stop() {
	app.GrpcApp.Stop()
}
