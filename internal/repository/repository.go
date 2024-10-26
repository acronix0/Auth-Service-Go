package repository

import (
	"context"
	"database/sql"
	"time"
)

type RefreshToken interface {
	SaveRefreshToken(ctx context.Context, userID int, refreshToken string, expiresAt time.Time, deviceInfo string) error
	ValidateRefreshToken(ctx context.Context, userID int, refreshToken string) (bool, error)
	DeleteRefreshToken(ctx context.Context, userID int, deviceInfo string) error
	DeleteAllRefreshTokens(ctx context.Context, userID int) error
}

type Repositories struct{
	RefreshToken RefreshToken
}

func NewRepositories(db *sql.DB) *Repositories{
	return &Repositories{RefreshToken: NewRefreshTokenRepository(db)}
}