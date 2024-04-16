package main

import (
	"github.com/dvid-messanger/internal/app/frontend"
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

	application, err := frontend.New(
		log,
		cfg.Services.Frontend.Port,
		[]byte(cfg.Services.Auth.Secret),
		cfg.Services.Frontend.WsPort,
		cfg.Services.Frontend.WsBasePath,
		cfg.Services.Frontend.SendBuffSize,
		cfg.Services.Frontend.RBuffSize,
		cfg.Services.Frontend.WBuffSize,
		cfg.Services.Frontend.HsTimeout,
		cfg.Services.Frontend.MsgLimit,
		cfg.Services.Frontend.WriteWait,
		cfg.Services.Frontend.PongWait,
		cfg.Clients.User.Address,
		cfg.Clients.User.Timeout,
		cfg.Clients.User.RetriesCount,
		cfg.Clients.Auth.Address,
		cfg.Clients.Auth.Timeout,
		cfg.Clients.Auth.RetriesCount,
		cfg.Clients.Chat.Address,
		cfg.Clients.Chat.Timeout,
		cfg.Clients.Chat.RetriesCount,
	)
	if err != nil {
		log.Error("unable to setup application", logger.Err(err))
		return
	}

	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping", slog.String("signal", sign.String()))
	application.Stop()

	log.Info("stopped")
}
