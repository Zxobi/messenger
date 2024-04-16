package auth

import (
	"github.com/dvid-messanger/internal/app/auth/grpc"
	"github.com/dvid-messanger/internal/lib/jwt"
	"github.com/dvid-messanger/internal/service/auth"
	"github.com/dvid-messanger/internal/storage/auth/inmem"
	"log/slog"
	"time"
)

type App struct {
	GrpcApp *grpc.App
}

func New(log *slog.Logger, grpcPort int, tokenTtl time.Duration, secret []byte) *App {
	storage := inmem.New()
	tokenMaker := jwt.NewTokenizer(secret)
	authService := auth.NewService(log, storage, storage, tokenMaker, tokenTtl)

	grpcApp := grpc.New(log, authService, grpcPort)

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
