package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/acronix0/Auth-Service-Go/internal/hash"
	"github.com/acronix0/Auth-Service-Go/internal/repository"
	"github.com/acronix0/Auth-Service-Go/internal/token"
)

type AuthService struct {
	log *slog.Logger
	tokenRepo repository.RefreshToken
	tokenManager token.TokenManager
	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
	hasher       hash.PasswordHasher
}

func NewAuthService(
	log *slog.Logger, 
	tokenRepo repository.RefreshToken,
	tokenManager token.TokenManager,
  accessTokenTTL time.Duration,
  refreshTokenTTL time.Duration,
  hasher hash.PasswordHasher,
) *AuthService{
	return &AuthService{
		log:               log,
    tokenRepo:         tokenRepo,
    tokenManager:       tokenManager,
    accessTokenTTL:     accessTokenTTL,
    refreshTokenTTL:    refreshTokenTTL,
    hasher:       hasher,
	}
}

func (s *AuthService) SignUp(ctx context.Context, input UserSignUpInput) (Tokens, error) {

/* 	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	user := domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
		Role: role,
		Blocked:  false,
	}

 	err = s.userRepo.Create(ctx, &user)
	if err != nil {
		return Tokens{}, err
	}  */
	//todo: user grpc
	return s.generateTokens(ctx, /* user.ID */123, input.DeviceInfo)
}
func (s *AuthService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
/* 	passwordHash, err := s.hasher.Hash(input.Password)
	if err!= nil {
    return Tokens{}, err
  }
	user, err := s.userRepo.Login(ctx, input.Email, passwordHash)
	if err != nil {
		return Tokens{}, err
	} */
	 resp,err := s.generateTokens(ctx, /* user.ID */123, input.DeviceInfo)
	return resp,err
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string, deviceInfo string) (Tokens, error) {
	userID, err := s.tokenManager.Parse(refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	valid, err := s.tokenRepo.ValidateRefreshToken(ctx, userID, refreshToken)
	if err != nil || !valid {
		return Tokens{}, err
	}

	return s.generateTokens(ctx, userID, deviceInfo)
}
func (s *AuthService)generateTokens(ctx context.Context, userID int, deviceInfo string) (Tokens, error){

	accessToken, err := s.tokenManager.NewJWT(userID, time.Hour*1)
	if err != nil {
		return Tokens{}, err
	}
	newRefreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return Tokens{}, err
	}
	err = s.tokenRepo.DeleteRefreshToken(ctx, userID, deviceInfo)
	if err != nil {
		return Tokens{}, err
	}
	err = s.tokenRepo.SaveRefreshToken(ctx, userID, newRefreshToken, time.Now().Add(s.refreshTokenTTL), deviceInfo)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
func (s *AuthService) DeleteAllRefreshTokens(ctx context.Context, userID int) error{
	return s.tokenRepo.DeleteAllRefreshTokens(ctx, userID)	
}
