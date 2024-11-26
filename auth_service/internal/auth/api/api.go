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

type Auth interface {
	Authenticate(ctx context.Context, in *authv1.AuthRequest) (*authv1.AuthResponse, error)
	Registration(ctx context.Context, in *authv1.RegistrationRequest) (*authv1.Nothing, error)
	GetUserDataByUsername(ctx context.Context, in *authv1.GetUserDataByUsernameRequest) (*authv1.GetUserDataByUsernameResponse, error)
	CreateJWT(ctx context.Context, in *authv1.CreateJWTRequest) (*authv1.Token, error)
	GetUserByJWT(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error)
	IsAuthorized(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error)
}

func init() {
	prometheus.MustRegister(requestAuthenticateDuration, requestCreateJWTDuration, requestGetUserByJWTDuration, requestGetUserDataByUsernameDuration,
		requestIsAuthorizedDuration, requestRegistrationDuration)
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

var requestAuthenticateDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "Authenticate_grpc_request_duration_seconds",
		Help: "grpcRequest",
	},
	[]string{"method"},
)

func (s Server) Authenticate(ctx context.Context, in *authv1.AuthRequest) (*authv1.AuthResponse, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestAuthenticateDuration, "grpc")
	}()
	return s.auth.Authenticate(ctx, in)
}

var requestRegistrationDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "Registration_grpc_request_duration_seconds",
		Help: "grpcRequest",
	},
	[]string{"method"},
)

func (s Server) Registration(ctx context.Context, in *authv1.RegistrationRequest) (*authv1.Nothing, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestRegistrationDuration, "grpc")
	}()
	return s.auth.Registration(ctx, in)
}

var requestGetUserDataByUsernameDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "GetUserDataByUsername_grpc_request_duration_seconds",
		Help: "grpcRequest",
	},
	[]string{"method"},
)

func (s Server) GetUserDataByUsername(ctx context.Context, in *authv1.GetUserDataByUsernameRequest) (*authv1.GetUserDataByUsernameResponse, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestGetUserDataByUsernameDuration, "grpc")
	}()
	return s.auth.GetUserDataByUsername(ctx, in)
}

var requestCreateJWTDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "CreateJWT_grpc_request_duration_seconds",
		Help: "grpcRequest",
	},
	[]string{"method"},
)

func (s Server) CreateJWT(ctx context.Context, in *authv1.CreateJWTRequest) (*authv1.Token, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestCreateJWTDuration, "grpc")
	}()
	return s.auth.CreateJWT(ctx, in)
}

var requestGetUserByJWTDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "GetUserByJWT_grpc_request_duration_seconds",
		Help: "grpcRequest",
	},
	[]string{"method"},
)

func (s Server) GetUserByJWT(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestGetUserByJWTDuration, "grpc")
	}()
	return s.auth.GetUserByJWT(ctx, in)
}

var requestIsAuthorizedDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "IsAuthorized_grpc_request_duration_seconds",
		Help: "grpcRequest",
	},
	[]string{"method"},
)

func (s Server) IsAuthorized(ctx context.Context, in *authv1.Token) (*authv1.UserJWT, error) {
	start := time.Now()
	defer func() {
		metric.WriteRequestDuration(start, requestIsAuthorizedDuration, "grpc")
	}()

	log.Printf("Пришел запрос на авторизацию %v", in.Token)
	return s.auth.IsAuthorized(ctx, in)
}
