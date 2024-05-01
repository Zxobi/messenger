package main

import (
	"github.com/dvid-messanger/internal/app/user"
	"github.com/dvid-messanger/internal/config"
	"github.com/dvid-messanger/internal/lib/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustSetupLogger(cfg.Env, cfg.LogLevel)

	log.Info("starting", slog.Any("cfg", cfg))

	application := user.New(log, &cfg.Services.User)
	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping", slog.String("signal", sign.String()))
	application.Stop()

	log.Info("stopped")
}
