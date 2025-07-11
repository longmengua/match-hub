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

func StartGRPCServer(enginManager *engin.EngineManager) error {
	listener, _ = net.Listen("tcp", ":50051")

	// Create gRPC server with middleware interceptors
	grpcSrv = grpc.NewServer(
		grpc.UnaryInterceptor(ChainUnaryInterceptors(
			TraceIDInterceptor,
			RecoveryInterceptor,
		)),
	)

	// 在這裡註冊你的 gRPC service，例如：
	proto.RegisterHealthServiceServer(grpcSrv, health.NewServer())
	proto.RegisterMatchServiceServer(grpcSrv, match.NewServer(enginManager))

	log.Println("gRPC server listening on :50051")
	return grpcSrv.Serve(listener)
}

func StopGRPCServer() {
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
