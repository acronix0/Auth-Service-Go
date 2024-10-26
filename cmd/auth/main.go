package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	app "github.com/acronix0/Auth-Service-Go/internal/app"
	"github.com/acronix0/Auth-Service-Go/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log := setUpLogger(cfg.Env)
	application := app.New(log, cfg)
	go func(){
		application.GRPCServer.MustRun()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	application.GRPCServer.Stop()
	log.Info("Server stopped")
}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case config.EnvLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}