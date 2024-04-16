package chat

import (
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/app/chat/grpc"
	"github.com/dvid-messanger/internal/client/frontend"
	"github.com/dvid-messanger/internal/service/chat"
	"github.com/dvid-messanger/internal/storage/chat/inmem"
	"log/slog"
	"time"
)

type App struct {
	GrpcApp *grpc.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	feClientAddr string,
	feClientTimeout time.Duration,
	feClientRetriesCount int,
) (*App, error) {
	const op = "chat.app.New"

	chatStorage := inmem.NewChatStorage()
	messageStorage := inmem.NewMessageStorage()
	notifier, err := frontend.New(context.TODO(), log, feClientAddr, feClientTimeout, feClientRetriesCount)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chatService := chat.NewService(log, chatStorage, chatStorage, notifier, chatStorage, messageStorage, messageStorage)

	grpcApp := grpc.New(log, chatService, grpcPort)

	return &App{
		GrpcApp: grpcApp,
	}, nil
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
