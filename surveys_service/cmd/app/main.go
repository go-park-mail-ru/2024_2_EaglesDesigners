package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/api"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/usecase"
	authv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/proto"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	log := logger.LoggerWithCtx(ctx, logger.Log)

	pool, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@postgres:5433/surveys")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	defer pool.Close()
	log.Println("База данных подключена")

	server := grpc.NewServer()

	repo := repository.NewRepository(pool)
	usecase := usecase.NewUsecase(repo)
	authServer := api.New(usecase)
	authv1.RegisterAuthServer(server, authServer)

	log.Println("starting server at :8084")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
	log.Println("server started at:8084")
}
