package server

import (
	"context"
	"fmt"
	user_service "github.com/levelord1311/grpc-microservices/grpc-user-service/internal/app/user-service"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/service/user"
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
	userService *user.Service
}

func NewGrpcServer(userService *user.Service) *GrpcServer {
	return &GrpcServer{
		userService: userService,
	}
}

func (s *GrpcServer) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcAddr := "127.0.0.1:6002"

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

	pb.RegisterUserServiceServer(grpcServer, user_service.NewUserService(*s.userService))

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
