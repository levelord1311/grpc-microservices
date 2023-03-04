package server

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/repo"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/service"
	pb "github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GrpcServer struct {
	db        *sqlx.DB
	batchSize uint64
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func (s *GrpcServer) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcAddr := "127.0.0.1:8080"

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer l.Close()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 1 * time.Minute,
			MaxConnectionAge:  1 * time.Minute,
			Time:              1 * time.Minute,
			Timeout:           15 * time.Second,
		}),
	)

	r := repo.NewRepo(s.db, s.batchSize)

	pb.RegisterUserServiceServer(grpcServer, service.NewUserService(r))

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Printf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		log.Printf("ctx.Done: %v", done)
	}

	grpcServer.GracefulStop()
	log.Println("grpcServer is shut down correctly")
	return nil
}
