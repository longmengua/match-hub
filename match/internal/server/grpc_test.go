package server_test

import (
	"context"
	"match/internal/server"
	"match/proto"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestGRPCServerLifecycle(t *testing.T) {
	// 啟動 gRPC 伺服器
	go func() {
		if err := server.StartGRPCServer(); err != nil {
			t.Errorf("failed to start gRPC server: %v", err)
		}
	}()

	// 等待 server 啟動（實際應更健壯，可改用 health check）
	time.Sleep(500 * time.Millisecond)

	// 建立 gRPC client 連線，驗證 server 有啟動
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		t.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewHealthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 執行 Health Check 或其他測試 RPC
	_, err = client.Check(ctx, &proto.HealthRequest{})
	if err != nil {
		t.Errorf("Ping RPC failed: %v", err)
	}

	// 停止 gRPC server
	server.StopGRPCServer()

	// 再次請求應失敗（可選）
	_, err = client.Check(ctx, &proto.HealthRequest{})
	if err == nil {
		t.Error("expected error after server shutdown, got nil")
	}
}
