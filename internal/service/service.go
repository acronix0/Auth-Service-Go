package service

import "context"
type Tokens struct {
	AccessToken  string
	RefreshToken string
}
type UserSignInInput struct {
	Email    string
	Password string
	DeviceInfo string 
}
type UserSignUpInput struct {
  Name     string
  Email    string
  Password string
  Phone    string
	DeviceInfo string
	Role      string
}
type Auth interface {
	SignUp(ctx context.Context, input UserSignUpInput) (Tokens, error)
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string, deviceInfo string) (Tokens, error)
	DeleteAllRefreshTokens(ctx context.Context, userID int) error
}