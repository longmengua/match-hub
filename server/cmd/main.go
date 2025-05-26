package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"server/internal/httpserver"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 啟動 HTTP server
	go func() {
		if err := httpserver.Start(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 啟動 gRPC server
	// go func() {
	// 	if err := grpcserver.Start(); err != nil {
	// 		log.Fatalf("gRPC server error: %v", err)
	// 	}
	// }()

	// 等待中斷訊號
	<-ctx.Done()
	log.Println("Main: shutdown signal received")

	// 呼叫個別 shutdown 函式
	httpserver.Stop()
	// grpcserver.Stop()

	log.Println("Main: all servers shutdown cleanly")
}
