package health

import (
	"context"
	"match/proto"
)

// Ensure Server implements proto.HealthServiceServer interface
var _ proto.MatchServiceServer = (*Server)(nil)

type Server struct{}

// GetOrderBook implements proto.MatchServiceServer.
func (s *Server) GetOrderBook(context.Context, *proto.MatchRequest) (*proto.MatchResponse, error) {
	panic("unimplemented")
}
