package main

import (
	"context"
	"log"
	"net"

	pb "match/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHealthServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	version := "v1.0.0" // Replace with actual version retrieval logic, e.g., from build info or environment variable
	return &pb.HealthResponse{Version: version}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHealthServiceServer(grpcServer, &server{})

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
