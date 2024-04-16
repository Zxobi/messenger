package grpc

import (
	"fmt"
	grpcchat "github.com/dvid-messanger/internal/server/chat"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRpcServer *grpc.Server
	port       int
}

func New(log *slog.Logger, chat grpcchat.Chat, port int) *App {
	gRpcServer := grpc.NewServer()
	grpcchat.Register(gRpcServer, chat)

	return &App{
		log:        log,
		gRpcServer: gRpcServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server running", slog.String("addr", l.Addr().String()))

	if err = a.gRpcServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping grpc server")
	a.gRpcServer.GracefulStop()
}
