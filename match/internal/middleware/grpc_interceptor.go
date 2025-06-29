package middleware

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Chain multiple interceptors together
func ChainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Chain interceptors recursively
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
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return interceptor(ctx, req, info, handler)
	}
}

// traceIDInterceptor injects or logs a traceID from metadata
func TraceIDInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if traceIDs := md.Get("trace_id"); len(traceIDs) > 0 {
			ctx = context.WithValue(ctx, "trace_id", traceIDs[0])
		}
	}
	// You can log or propagate the trace_id here if needed
	return handler(ctx, req)
}

// recoveryInterceptor recovers from panics and returns a gRPC error
func RecoveryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered in gRPC method %s: %v", info.FullMethod, r)
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}
