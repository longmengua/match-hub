package grpc

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Chain multiple interceptors together
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		curr := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			curr = WrapUnaryHandler(interceptor, info, curr)
		}
		return curr(ctx, req)
	}
}

func WrapUnaryHandler(
	interceptor grpc.UnaryServerInterceptor,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) grpc.UnaryHandler {
	return func(ctx context.Context, req any) (any, error) {
		return interceptor(ctx, req, info, handler)
	}
}

type contextKey string

const TraceIDKey contextKey = "trace_id"

// Helper: get trace_id from context, traceID will be populated in TraceIDInterceptor
func GetTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(TraceIDKey).(string); ok {
		return v
	}
	return ""
}

// Interceptor: inject trace_id (generate one if missing), and log if generated
func TraceIDInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	var traceID string
	if traceIDs := md.Get("trace_id"); len(traceIDs) > 0 {
		traceID = traceIDs[0]
	} else {
		// Generate a new trace_id if not present
		// todo: can change UUID in something more meaningful
		traceID = uuid.New().String()
		log.Printf("[TRACE] Generated trace_id=%s | method=%s | req=%+v", traceID, info.FullMethod, req)
	}

	ctx = context.WithValue(ctx, TraceIDKey, traceID)
	return handler(ctx, req)
}

// Interceptor: recover from panic
func RecoveryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC] Recovered in gRPC method %s: %v", info.FullMethod, r)
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}

// Interceptor: log errors if any
func ErrorLoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		traceID := GetTraceID(ctx)
		log.Printf("[ERROR] method=%s | trace_id=%s | error=%v", info.FullMethod, traceID, err)
	}
	return resp, err
}
