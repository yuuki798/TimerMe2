package main

import (
	"fmt"
	"log"
	"net"
	"task_service/internal/service"

	"google.golang.org/grpc"
	pb "proto/task"
)

func main() {
	service.InitDB()
	fmt.Println("50051 listening..")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &service.Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
