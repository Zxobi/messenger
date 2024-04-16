package http

import (
	"context"
	"fmt"
	"github.com/dvid-messanger/internal/server/frontend/ws"
	"github.com/gorilla/mux"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type App struct {
	log        *slog.Logger
	port       int
	wsPath     string
	wsServer   *ws.WsServer
	httpServer *http.Server
}

func New(
	log *slog.Logger,
	wsServer *ws.WsServer,
	port int,
	wsPath string,
) *App {
	return &App{
		log:      log,
		port:     port,
		wsPath:   wsPath,
		wsServer: wsServer,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "http.Run"
	log := a.log.With(slog.String("op", op))

	router := mux.NewRouter()
	router.HandleFunc(a.wsPath, a.wsServer.Handle)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.httpServer = &http.Server{Addr: l.Addr().String(), Handler: router}

	log.Info("http server running", slog.String("addr", l.Addr().String()))

	if err = a.httpServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "http.Stop"
	a.log.With(slog.String("op", op)).Info("stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() { cancel() }()

	_ = a.httpServer.Shutdown(ctx)
}
