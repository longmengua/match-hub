package grpc_test

import (
	"context"
	"match/internal/server/grpc"
	"testing"

	googleGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestTraceIDInterceptor_WithTraceID(t *testing.T) {
	traceID := "test-trace-123"
	md := metadata.Pairs("trace_id", traceID)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	called := false
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		val := ctx.Value("trace_id")
		if val != traceID {
			t.Errorf("expected trace_id %s, got %v", traceID, val)
		}
		called = true
		return "ok", nil
	}

	resp, err := grpc.TraceIDInterceptor(ctx, nil, &googleGRPC.UnaryServerInfo{
		FullMethod: "/test/trace",
	}, handler)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Errorf("expected response 'ok', got %v", resp)
	}
	if !called {
		t.Error("handler not called")
	}
}

func TestTraceIDInterceptor_WithoutTraceID(t *testing.T) {
	ctx := context.Background()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		val := ctx.Value("trace_id")
		if val != nil {
			t.Errorf("expected no trace_id, got %v", val)
		}
		return "ok", nil
	}

	resp, err := grpc.TraceIDInterceptor(ctx, nil, &googleGRPC.UnaryServerInfo{
		FullMethod: "/test/trace",
	}, handler)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Errorf("expected response 'ok', got %v", resp)
	}
}

func TestRecoveryInterceptor_NoPanic(t *testing.T) {
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "safe", nil
	}

	resp, err := grpc.RecoveryInterceptor(context.Background(), nil, &googleGRPC.UnaryServerInfo{
		FullMethod: "/test/recovery",
	}, handler)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != "safe" {
		t.Errorf("unexpected response: %v", resp)
	}
}

func TestRecoveryInterceptor_WithPanic(t *testing.T) {
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("boom")
	}

	resp, err := grpc.RecoveryInterceptor(context.Background(), nil, &googleGRPC.UnaryServerInfo{
		FullMethod: "/test/recovery",
	}, handler)

	if resp != nil {
		t.Errorf("expected nil resp on panic, got: %v", resp)
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.Internal {
		t.Errorf("expected gRPC internal error, got: %v", err)
	}
}

func TestChainUnaryInterceptors_OrderAndChaining(t *testing.T) {
	order := []string{}

	interceptor1 := func(ctx context.Context, req interface{}, info *googleGRPC.UnaryServerInfo, handler googleGRPC.UnaryHandler) (interface{}, error) {
		order = append(order, "1-before")
		resp, err := handler(ctx, req)
		order = append(order, "1-after")
		return resp, err
	}

	interceptor2 := func(ctx context.Context, req interface{}, info *googleGRPC.UnaryServerInfo, handler googleGRPC.UnaryHandler) (interface{}, error) {
		order = append(order, "2-before")
		resp, err := handler(ctx, req)
		order = append(order, "2-after")
		return resp, err
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		order = append(order, "handler")
		return "done", nil
	}

	chain := grpc.ChainUnaryInterceptors(interceptor1, interceptor2)
	resp, err := chain(context.Background(), nil, &googleGRPC.UnaryServerInfo{
		FullMethod: "/test/chain",
	}, handler)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != "done" {
		t.Errorf("expected response 'done', got: %v", resp)
	}

	expectedOrder := []string{"1-before", "2-before", "handler", "2-after", "1-after"}
	for i, step := range expectedOrder {
		if order[i] != step {
			t.Errorf("order[%d]: expected %s, got %s", i, step, order[i])
		}
	}
}
