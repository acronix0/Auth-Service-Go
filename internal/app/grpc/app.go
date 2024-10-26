package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/acronix0/Auth-Service-Go/internal/grpc/auth"
	srv "github.com/acronix0/Auth-Service-Go/internal/service"
	"google.golang.org/grpc"
)

type App struct{
	log *slog.Logger
	gRPCServer *grpc.Server
	port int
}

func New(logger *slog.Logger, port int, auth srv.Auth) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.RegisterServer(gRPCServer,auth)
	return &App{
    log:       logger,
    gRPCServer: gRPCServer,
    port:      port,
  }
}
func (a *App) MustRun(){
	if err:= a.Run(); err != nil {
		panic(err)
	}
}
func (a *App) Run() error{
	const op = "grpcapp.Run"
	log := a.log.With(
		slog.String("op",op),
		slog.Int("port", a.port),
	)

	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("grpc server is running", slog.String("addr",listner.Addr().String()))
  if err:= a.gRPCServer.Serve(listner); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

  log := a.log.With(slog.String("op", op))
  log.Info("stopping gRPC server")
  
	a.gRPCServer.GracefulStop()
}