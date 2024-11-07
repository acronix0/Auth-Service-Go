package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	grpcapp "github.com/acronix0/Auth-Service-Go/internal/app/grpc"
	"github.com/acronix0/Auth-Service-Go/internal/config"

)


type App struct {
	serviceProvider *serviceProvider
	GRPCServer *grpcapp.App
	log *slog.Logger
	config *config.Config
}
func NewApp(ctx context.Context) (*App, error){
	a := &App{config: config.MustLoad()}
	err := a.initDeps(ctx)
	if err!= nil {
    return nil, err
  }
	return a, nil
}

func (a *App)Run() error{
	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error{
	inits := []func(context.Context) error{
		a.initLogger,
    a.initServiceProvider,
		a.initGRPCServer,
	}	
	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
func (a *App) initServiceProvider(_ context.Context) error{
	a.serviceProvider = newServiceProvider(a.config)
	return nil
}


func (a *App) initGRPCServer(_ context.Context) error{
	a.GRPCServer = grpcapp.New(a.log, a.config.GRPC.Port, a.serviceProvider.AuthService())
	return nil
}
func (a *App) runGRPCServer() error{
	const op = "app.runGRPCServer"
  log := a.log.With(slog.String("op", op))
  log.Info("Starting gRPC server")
	go func(){
		a.GRPCServer.MustRun()
	}()
  stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	a.GRPCServer.Stop()
	log.Info("Server stopped")
  return nil
}
func (a *App)initLogger(_ context.Context) error {

	switch a.config.Env {
	case config.EnvLocal:
		a.log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.EnvProd:
		a.log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return nil
}