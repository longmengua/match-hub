package main

import (
	"context"
	"log"
	"net"

	pb "grpc-hello-world/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreetServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	message := "Hello, " + req.Name
	return &pb.HelloResponse{Message: message}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &server{})

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
