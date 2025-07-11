package health

import (
	"context"
	"match/proto"
)

// Ensure Server implements proto.HealthServiceServer interface
var _ proto.HealthServiceServer = (*Server)(nil)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Check(ctx context.Context, req *proto.HealthRequest) (*proto.HealthResponse, error) {
	version := "v1.0.0" // Replace with actual version retrieval logic, e.g., from build info or environment variable
	if len(req.MustBeHello) == 0 {
		return &proto.HealthResponse{Error: &proto.Error{
			Code:    proto.ErrorCode_INVALID_ARGUMENT,
			Message: "MustBeHello cannot be empty",
		}}, nil
	}
	return &proto.HealthResponse{Version: version}, nil
}
