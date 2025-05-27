package service

import (
	"context"
	"match/proto"
)

type Server struct {
	proto.UnimplementedHealthServiceServer
}

func (s *Server) Check(ctx context.Context, req *proto.HealthRequest) (*proto.HealthResponse, error) {
	version := "v1.0.0" // Replace with actual version retrieval logic, e.g., from build info or environment variable
	return &proto.HealthResponse{Version: version}, nil
}
