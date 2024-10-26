package auth

import (
	"context"

	authv1 "github.com/acronix0/REST-API-Go-protos/gen/go/auth"
	srv "github.com/acronix0/Auth-Service-Go/internal/service"
	"google.golang.org/grpc"
)
type serverApi struct{
	authv1.UnimplementedAuthServer
	auth srv.Auth
}

func RegisterServer(grpc *grpc.Server, auth srv.Auth){
  authv1.RegisterAuthServer(grpc, &serverApi{auth: auth})
}

func (s *serverApi) SignIn(
	ctx context.Context,
	req *authv1.SignInRequest,
)(*authv1.SignInResponse, error){
	res,err := s.auth.SignIn(ctx, srv.UserSignInInput{
		Email: req.GetEmail(),
    Password: req.GetPassword(),
		DeviceInfo: req.GetDeviceInfo(),
	})
	if err != nil {
		return nil, nil
	}
	return &authv1.SignInResponse{JwtToken: res.AccessToken+"123", RefreshToken: res.RefreshToken}, nil
}

func (s *serverApi) SignUp(
  ctx context.Context,
  req *authv1.SignUpRequest,
) (*authv1.SignUpResponse, error){
	panic("impls me")
}

func (s *serverApi) SignOut(
	ctx context.Context,
  req *authv1.SignOutRequest,
)(*authv1.SignOutResponse, error){
	panic("impls me")
}

func (s *serverApi) RefreshToken(
  ctx context.Context,
  req *authv1.RefreshTokenRequest,
)(*authv1.RefreshTokenResponse, error){
	panic("impls me")
}