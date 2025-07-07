package main

import (
	"context"
	"log"
	server "match/internal/server/grpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 啟動 gRPC server
	go func() {
		if err := server.StartGRPCServer(); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// 等待中斷訊號
	<-ctx.Done()
	log.Println("Main: shutdown signal received")

	// 呼叫個別 shutdown 函式
	server.StopGRPCServer()

	log.Println("Main: all servers shutdown cleanly")
}
