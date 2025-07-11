package match

import (
	"context"
	"match/pkg/engin"
	"match/proto"
)

// Ensure Server implements proto.HealthServiceServer interface
var _ proto.MatchServiceServer = (*Server)(nil)

type Server struct {
	EnginManager *engin.EngineManager
}

func NewServer(enginManager *engin.EngineManager) *Server {
	return &Server{
		EnginManager: enginManager,
	}
}

// GetOrderBook implements proto.MatchServiceServer.
func (s *Server) GetOrderBook(context.Context, *proto.MatchRequest) (*proto.MatchResponse, error) {
	panic("unimplemented")
}
