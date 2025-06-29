package service

import (
	"context"
	"match/internal/errors"
	"match/proto"

	"google.golang.org/grpc/codes"
)

type Server struct {
	proto.UnimplementedHealthServiceServer
}

func (s *Server) Check(ctx context.Context, req *proto.HealthRequest) (*proto.HealthResponse, error) {
	version := "v1.0.0" // Replace with actual version retrieval logic, e.g., from build info or environment variable
	if len(req.MustBeHello) == 0 {
		return nil, errors.NewGRPCError(codes.InvalidArgument, "invalid input")
	}
	return &proto.HealthResponse{Version: version}, nil
}
