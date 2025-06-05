package main

import (
	"context"
	"log"
	"match/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
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
