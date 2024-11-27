package api

import (
	"context"
	"log"
	"time"

	authv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/proto"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"
	"github.com/prometheus/client_golang/prometheus"
)

//go:generate mockgen -source=api.go -destination=mocks/mocks.go

type Auth interface {
	Authenticate(ctx context.Context, in *authv1.AuthRequest) (*authv1.AuthResponse, error)
	Registration(ctx context.Context, in *authv1.RegistrationRequest) (*authv1.Nothing, error)
	GetUserDataByUsername(ctx context.Context, in *authv1.GetUserDataByUsernameRequest) (*authv1.GetUserDataByUsernameResponse, error)
	CreateJWT(ctx context.Context, in *authv1.CreateJWTRequest) (*authv1.Token, error)
	GetUserByJWT(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error)
	IsAuthorized(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error)
}

func init() {
	prometheus.MustRegister(requestAuthDuration)
	log := logger.LoggerWithCtx(context.Background(), logger.Log)
	log.Info("Метрики для авторизации зарегистрированы")
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

var requestAuthDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "grpc_authenticate_request_duration_seconds",
		Help: "grpcRequest",
	},
	[]string{"method"},
)

func (s Server) Authenticate(ctx context.Context, in *authv1.AuthRequest) (*authv1.AuthResponse, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthDuration, "Authenticate")
	}()
	return s.auth.Authenticate(ctx, in)
}

func (s Server) Registration(ctx context.Context, in *authv1.RegistrationRequest) (*authv1.Nothing, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthDuration, "Registration")
	}()
	return s.auth.Registration(ctx, in)
}

func (s Server) GetUserDataByUsername(ctx context.Context, in *authv1.GetUserDataByUsernameRequest) (*authv1.GetUserDataByUsernameResponse, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthDuration, "GetUserDataByUsername")
	}()
	return s.auth.GetUserDataByUsername(ctx, in)
}

func (s Server) CreateJWT(ctx context.Context, in *authv1.CreateJWTRequest) (*authv1.Token, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthDuration, "CreateJWT")
	}()
	return s.auth.CreateJWT(ctx, in)
}

func (s Server) GetUserByJWT(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthDuration, "GetUserByJWT")
	}()
	return s.auth.GetUserByJWT(ctx, in)
}

func (s Server) IsAuthorized(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthDuration, "IsAuthorized")
	}()

	log.Printf("Пришел запрос на авторизацию %v", in.Token)
	return s.auth.IsAuthorized(ctx, in)
}
