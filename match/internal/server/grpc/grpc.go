package grpc

import (
	"log"
	"match/internal/server/grpc/health"
	"match/internal/server/grpc/match"
	"match/pkg/engin"
	"match/proto"
	"net"
	"time"

	"google.golang.org/grpc"
)

var (
	grpcSrv  *grpc.Server
	listener net.Listener
)

type GRPCServer struct {
	engineManager *engin.EngineManager
}

func NewGRPCServer(enginManager *engin.EngineManager) *GRPCServer {
	return &GRPCServer{
		engineManager: enginManager,
	}
}

func (g *GRPCServer) Start() error {
	listener, _ = net.Listen("tcp", ":50051")

	// Create gRPC server with middleware interceptors
	grpcSrv = grpc.NewServer(
		grpc.UnaryInterceptor(ChainUnaryInterceptors(
			TraceIDInterceptor,
			ErrorLoggingInterceptor,
			RecoveryInterceptor,
		)),
	)

	// 在這裡註冊你的 gRPC service，例如：
	proto.RegisterHealthServiceServer(grpcSrv, health.NewServer())
	proto.RegisterMatchServiceServer(grpcSrv, match.NewServer(g.engineManager))

	log.Println("gRPC server listening on :50051")
	return grpcSrv.Serve(listener)
}

func (g *GRPCServer) Stop() {
	log.Println("Shutting down gRPC server...")

	// 支援優雅關閉
	stopped := make(chan struct{})
	go func() {
		grpcSrv.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("gRPC server shutdown completed")
	case <-time.After(5 * time.Second):
		log.Println("gRPC server shutdown timed out, forcing stop")
		grpcSrv.Stop()
	}
}
