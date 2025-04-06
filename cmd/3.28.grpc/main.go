package main

import (
	"log"
	"net"

	"backend-engineering/cmd/3.28.grpc/backend.engineering/gen"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	gen.RegisterTodoServer(s, &server{})

	log.Println("Server is running on port 50051...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
