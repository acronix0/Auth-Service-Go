package app

import (
	"log"
	"log/slog"
	"os"

	"github.com/acronix0/Auth-Service-Go/internal/config"
	"github.com/acronix0/Auth-Service-Go/internal/database"
	"github.com/acronix0/Auth-Service-Go/internal/repository"
	"github.com/acronix0/Auth-Service-Go/internal/service"
	"github.com/acronix0/Auth-Service-Go/internal/token"
	"github.com/acronix0/Auth-Service-Go/internal/hash"
)

type serviceProvider struct {
	config *config.Config
	dataBase *database.Database
	authRepository repository.RefreshToken
	logger *slog.Logger
	authService service.Auth
	tokenManager *token.Manager
	hash *hash.SHA1Hasher
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{config: cfg}
}
func (s *serviceProvider) Config() *config.Config{
	if s.config == nil {
    s.config = config.MustLoad() 
	}
	return s.config
}
func (s *serviceProvider) Logger() *slog.Logger{
	if s.logger == nil {
    s.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
  }
  return s.logger
}
func (s *serviceProvider) AuthRepository() repository.RefreshToken {
	if s.authRepository == nil {
    s.authRepository = repository.NewRefreshTokenRepository(s.Database().GetDB())
  }

  return s.authRepository
}
func (s *serviceProvider) TokenManager() *token.Manager{
	if s.tokenManager == nil {
    tokenManager, err := token.NewManager(s.Config().AuthConfig.JWT.Secret)
		if err!= nil {
			log.Fatalf("failed to load token manager: %s", err.Error())
		}else{
			s.tokenManager = tokenManager
		}
  }

  return s.tokenManager
}
func (s *serviceProvider) Hasher() *hash.SHA1Hasher{
	if s.hash == nil {
    s.hash = hash.NewSHA1Hasher(s.Config().AuthConfig.PasswordSalt)
  }

  return s.hash
}
func (s *serviceProvider) AuthService() service.Auth {
	if s.authService == nil {
    s.authService = service.NewAuthService(
		s.Logger(), 
		s.AuthRepository(), 
		s.TokenManager(), 
		s.Config().AuthConfig.JWT.AccessTokenTTL, 
		s.Config().AuthConfig.JWT.RefreshTokenTTL, 
		s.Hasher())
  }

  return s.authService

}
func (s *serviceProvider) Database() *database.Database{
	if s.dataBase == nil {
		db, err := database.New(&s.Config().SqlConnection)
  	if err!= nil {
    	log.Fatalf("failed to connect to database: %s", err.Error())
  	}
		s.dataBase = db
	}

	return s.dataBase
}