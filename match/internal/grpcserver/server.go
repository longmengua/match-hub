package grpcserver

import (
	"log"
	"match/internal/grpcserver/handler"
	"match/proto"
	"net"
	"time"

	"google.golang.org/grpc"
)

var (
	grpcSrv  *grpc.Server
	listener net.Listener
)

func Start() error {
	listener, _ = net.Listen("tcp", ":50051")

	grpcSrv = grpc.NewServer()
	// 在這裡註冊你的 gRPC service，例如：
	proto.RegisterHealthServiceServer(grpcSrv, &handler.Server{})

	log.Println("gRPC server listening on :50051")
	return grpcSrv.Serve(listener)
}

func Stop() {
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
