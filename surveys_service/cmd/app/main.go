package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	surveysv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/surveys_service/internal/proto"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/surveys_service/internal/surveys/api"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/surveys_service/internal/surveys/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/surveys_service/internal/surveys/usecase"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	lis, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	log := logger.LoggerWithCtx(ctx, logger.Log)

	pool, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@postgres:5432/patefon")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	defer pool.Close()
	log.Println("База данных подключена")

	server := grpc.NewServer()

	repo := repository.NewServeyRepository(pool)
	usecase := usecase.NewUsecase(repo)
	surveyServer := api.New(usecase)
	surveysv1.RegisterSurveysServer(server, surveyServer)

	log.Println("starting server at :8084")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
	log.Println("server started at:8084")
}
