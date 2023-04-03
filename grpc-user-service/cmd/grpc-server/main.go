package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/config"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/server"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/service/user"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/service/user/repo"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/database"
	"github.com/pressly/goose/v3"
	"log"
	"time"
)

var (
	batchSize uint64 = 2
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal("failed to read config:", err)
	}
	cfg := config.GetConfigInstance()

	migration := flag.Bool("migration", false, "start with migration up")
	flag.Parse()

	initCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, err := database.ConnectDB(initCtx, cfg.Database)
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			fmt.Println("migration failed:", err)
		}
	}

	r := repo.NewRepo(db, batchSize)
	userService := user.NewService(r)

	s := server.NewGrpcServer(userService)
	if err := s.Start(); err != nil {
		log.Fatal("Failed to create gRPC server", err)
	}

}
