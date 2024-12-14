package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/api"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/repository"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/auth/usecase"
	authv1 "github.com/go-park-mail-ru/2024_2_EaglesDesigner/auth_service/internal/proto"
	dbConfig "github.com/go-park-mail-ru/2024_2_EaglesDesigner/db/config"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/metric"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	pool := connectToPSQL()
	defer pool.Close()
	log.Println("База данных подключена")

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	log := logger.LoggerWithCtx(ctx, logger.Log)

	server := grpc.NewServer()

	repo := repository.NewRepository(pool)
	usecase := usecase.NewUsecase(repo)
	authServer := api.New(usecase)
	authv1.RegisterAuthServer(server, authServer)

	go func() {
		log.Println("starting server at :8081")
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
		log.Println("server started at:8081")
	}()
	metric.CollectMetrics()
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8087", nil); err != nil {
		log.Fatalf("failed to start HTTP server %v", err)
	}
}

func connectToPSQL() *pgxpool.Pool {
	config, err := dbConfig.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Формируем строку подключения
	config.Database.MaxPoolSize = 6
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?pool_max_conns=%d",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
		config.Database.MaxPoolSize,
	)

	// Настройка пула соединений
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	return pool
}
