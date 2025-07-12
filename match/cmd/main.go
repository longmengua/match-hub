package main

import (
	"context"
	"log"
	"match/config"
	"match/internal/server/grpc"
	"match/pkg/engin"
	"match/pkg/redisclient"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 初始化撮合引擎管理器等
	engineMap := make(map[string]*engin.Engine)
	engineMap["BTCUSDT"] = engin.NewEngine()
	engineMap["ETHUSDT"] = engin.NewEngine()
	engineMap["DOGEUSDT"] = engin.NewEngine()
	enginManager := engin.NewEngineManager(engineMap)

	// 初始化 Redis 連線
	_, err := redisclient.New("single", config.RedisAddrs, config.RedisPass, config.RedisDB)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// 啟動 gRPC server
	grpcServer := grpc.NewGRPCServer(enginManager)
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// 等待中斷訊號
	<-ctx.Done()
	log.Println("Main: shutdown signal received")

	// 呼叫個別 shutdown 函式
	grpcServer.Stop()

	log.Println("Main: all servers shutdown cleanly")
}
