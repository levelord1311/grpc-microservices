package main

import (
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/server"
	"log"
)

func main() {
	s := server.NewGrpcServer()
	if err := s.Start(); err != nil {
		log.Fatal("Failed to create gRPC server", err)
	}

}
