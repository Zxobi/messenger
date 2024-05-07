package main

import (
	"github.com/dvid-messanger/internal/app/chat"
	"github.com/dvid-messanger/internal/config"
	"github.com/dvid-messanger/internal/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustSetupLogger(cfg.Env, cfg.LogLevel)

	log.Info("starting", slog.Any("cfg", cfg))

	application, err := chat.New(
		log,
		&cfg.Services.Chat,
		&cfg.Clients.Frontend,
	)
	if err != nil {
		panic(err)
	}
	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping", slog.String("signal", sign.String()))
	application.Stop()

	log.Info("stopped")
}
