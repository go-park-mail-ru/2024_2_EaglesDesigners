package api

import (
	"context"
	"log"

	authv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/proto"
)

type Auth interface {
	Authenticate(ctx context.Context, in *authv1.AuthRequest) (*authv1.AuthResponse, error)
	Registration(ctx context.Context, in *authv1.RegistrationRequest) (*authv1.Nothing, error)
	GetUserDataByUsername(ctx context.Context, in *authv1.GetUserDataByUsernameRequest) (*authv1.GetUserDataByUsernameResponse, error)
	CreateJWT(ctx context.Context, in *authv1.CreateJWTRequest) (*authv1.Token, error)
	GetUserByJWT(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error)
	IsAuthorized(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error)
}

type Server struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func New(auth Auth) authv1.AuthServer {
	return Server{
		auth: auth,
	}
}

func (s Server) Authenticate(ctx context.Context, in *authv1.AuthRequest) (*authv1.AuthResponse, error) {
	return s.auth.Authenticate(ctx, in)
}

func (s Server) Registration(ctx context.Context, in *authv1.RegistrationRequest) (*authv1.Nothing, error) {
	return s.auth.Registration(ctx, in)
}

func (s Server) GetUserDataByUsername(ctx context.Context, in *authv1.GetUserDataByUsernameRequest) (*authv1.GetUserDataByUsernameResponse, error) {
	return s.auth.GetUserDataByUsername(ctx, in)
}

func (s Server) CreateJWT(ctx context.Context, in *authv1.CreateJWTRequest) (*authv1.Token, error) {
	return s.auth.CreateJWT(ctx, in)
}

func (s Server) GetUserByJWT(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error) {
	return s.auth.GetUserByJWT(ctx, in)
}

func (s Server) IsAuthorized(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error) {
	log.Printf("Пришел запрос на авторизацию %v", in.Token)
	return s.auth.IsAuthorized(ctx, in)
}
