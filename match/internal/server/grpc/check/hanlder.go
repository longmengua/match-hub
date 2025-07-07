package check

import (
	"context"
	"fmt"
	"match/internal/enum"
	"match/internal/errors"
	"match/proto"
)

// Ensure Server implements proto.HealthServiceServer interface
var _ proto.HealthServiceServer = (*Server)(nil)

type Server struct{}

func (s *Server) Check(ctx context.Context, req *proto.HealthRequest) (*proto.HealthResponse, error) {
	version := "v1.0.0" // Replace with actual version retrieval logic, e.g., from build info or environment variable
	if len(req.MustBeHello) == 0 {
		return nil, errors.New(enum.InvalidParams, fmt.Errorf("Invalid request params"))
	}
	return &proto.HealthResponse{Version: version}, nil
}
