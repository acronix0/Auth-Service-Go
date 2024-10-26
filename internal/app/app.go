package app

import (
	"log/slog"

	grpcapp "github.com/acronix0/Auth-Service-Go/internal/app/grpc"
	"github.com/acronix0/Auth-Service-Go/internal/config"
	"github.com/acronix0/Auth-Service-Go/internal/database"
	"github.com/acronix0/Auth-Service-Go/internal/hash"
	"github.com/acronix0/Auth-Service-Go/internal/repository"
	srv "github.com/acronix0/Auth-Service-Go/internal/service"
	"github.com/acronix0/Auth-Service-Go/internal/token"
)


type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {

	postgreClient, err := database.New(&cfg.SqlConnection)
	if err != nil {
		log.Error(err.Error())
	}
	tokenManager, err := token.NewManager(cfg.AuthConfig.JWT.Secret)
	if err!=nil {
		log.Error(err.Error())
	}
	repositories := repository.NewRepositories(postgreClient.GetDB())
	hasher := hash.NewSHA1Hasher(cfg.AuthConfig.PasswordSalt)
	authService := srv.NewAuthService(log, repositories.RefreshToken, tokenManager,cfg.AuthConfig.JWT.AccessTokenTTL, cfg.AuthConfig.JWT.RefreshTokenTTL,hasher)

	grpcApp := grpcapp.New(log, cfg.GRPC.Port, authService)

	return &App{
		GRPCServer: grpcApp,
	}
}
