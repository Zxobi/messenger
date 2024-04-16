package frontend

import (
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/app/frontend/grpc"
	"github.com/dvid-messanger/internal/app/frontend/http"
	"github.com/dvid-messanger/internal/client/auth"
	"github.com/dvid-messanger/internal/client/chat"
	"github.com/dvid-messanger/internal/client/user"
	"github.com/dvid-messanger/internal/lib/jwt"
	"github.com/dvid-messanger/internal/server/frontend/ws"
	"github.com/dvid-messanger/internal/server/frontend/ws/handler"
	fe "github.com/dvid-messanger/internal/service/frontend"
	"log/slog"
	"sync"
	"time"
)

type App struct {
	HttpApp *http.App
	GrpcApp *grpc.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	secret []byte,

	wsPort int,
	wsPath string,
	sendBuffSize int,
	rBuffSize int,
	wBuffSize int,
	hsTimeout time.Duration,
	msgLimit int64,
	writeWait time.Duration,
	pongWait time.Duration,
	userClientAddr string,
	userClientTimeout time.Duration,
	userClientRetriesCount int,
	authClientAddr string,
	authClientTimeout time.Duration,
	authClientRetriesCount int,
	chatClientAddr string,
	chatClientTimeout time.Duration,
	chatClientRetriesCount int,
) (*App, error) {
	const op = "frontend.New"

	userClient, err := user.New(
		context.TODO(),
		log,
		userClientAddr,
		userClientTimeout,
		userClientRetriesCount,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	authClient, err := auth.New(
		context.TODO(),
		log,
		authClientAddr,
		authClientTimeout,
		authClientRetriesCount,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chatClient, err := chat.New(
		context.TODO(),
		log,
		chatClientAddr,
		chatClientTimeout,
		chatClientRetriesCount,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	verifier := jwt.NewTokenizer(secret)
	registry := fe.NewClientRegistry(log)

	wsHandler := handler.NewHandler(log, userClient, authClient, chatClient, verifier, registry)
	wsServer := ws.NewWsServer(
		log,
		wsHandler,
		wsHandler,
		sendBuffSize,
		rBuffSize,
		wBuffSize,
		hsTimeout,
		msgLimit,
		writeWait,
		pongWait,
	)

	httpApp := http.New(log, wsServer, wsPort, wsPath)

	notifier := fe.NewNotifier(log, registry, registry)
	grpcApp := grpc.New(log, notifier, grpcPort)

	return &App{
		HttpApp: httpApp,
		GrpcApp: grpcApp,
	}, nil
}

func (app *App) MustRun() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		app.HttpApp.MustRun()
		wg.Done()
	}()
	go func() {
		app.GrpcApp.MustRun()
		wg.Done()
	}()

	wg.Wait()
}

func (app *App) Run() error {
	if err := app.HttpApp.Run(); err != nil {
		return err
	}

	if err := app.GrpcApp.Run(); err != nil {
		return err
	}

	return nil
}

func (app *App) Stop() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		app.HttpApp.Stop()
		wg.Done()
	}()
	go func() {
		app.GrpcApp.Stop()
		wg.Done()
	}()

	wg.Wait()
}
