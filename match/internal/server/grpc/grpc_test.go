package grpc_test

import (
	"context"
	"errors"
	server "match/internal/server/grpc"
	"match/proto"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func waitForServerReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := grpc.Dial("localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithTimeout(200*time.Millisecond),
		)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("gRPC server did not become ready in time")
}

func TestGRPCServerLifecycle(t *testing.T) {
	grpcServer := server.NewGRPCServer(nil)
	// 啟動 gRPC server
	go func() {
		if err := grpcServer.Start(); err != nil {
			t.Errorf("failed to start gRPC server: %v", err)
		}
	}()

	// 等待 server readiness
	if err := waitForServerReady(3 * time.Second); err != nil {
		t.Fatalf("server not ready: %v", err)
	}

	// 建立 gRPC client
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	client := proto.NewHealthServiceClient(conn)

	// 呼叫 Health Check
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
	defer cancel1()

	// _, err = client.Check(ctx1, &proto.HealthRequest{})
	// if err != nil {
	// 	t.Errorf("expected successful health check, got error: %v", err)
	// }
	res, err := client.Check(ctx1, &proto.HealthRequest{MustBeHello: ""})
	if res.Error.TraceId == "" {
		t.Error("Expected trace_id to be set in response, got empty")
	}
	if err != nil {
		t.Errorf("Unexpected error, got: %v", err)
	}
	if res.Error.Code != proto.ErrorCode_INVALID_ARGUMENT {
		t.Errorf("Unexpected error, got: %v", res.Error)
	}

	// 關閉 server
	grpcServer.Stop()

	// 等待一點時間以確保 server 停下
	time.Sleep(200 * time.Millisecond)

	// 嘗試使用新 context 再次呼叫 RPC（應該失敗）
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()

	_, err = client.Check(ctx2, &proto.HealthRequest{})
	if err == nil {
		t.Error("expected error after server shutdown, got nil")
	}
}
