package trace

import (
	"context"

	grpcutil "github.com/XHXHXHX/medical_marketing/util/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type traceContextKey struct{}

const IotTraceHeader = "aiforward_logid"

func ContextWithNewTraceID(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceContextKey{}, NewTraceID())
}

func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceContextKey{}, traceID)
}

func TraceIDFromContext(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(traceContextKey{}).(string)
	if len(traceID) == 0 {
		return "", false
	}
	return traceID, ok
}

func TraceIDFromIncommingGrpcContext(ctx context.Context) string {
	return grpcutil.HeaderFromContext(ctx, IotTraceHeader)
}

// 设置 grpc 的 header 和 context
func OutgoingGrpcHeaderContextWithTraceID(ctx context.Context, traceID string) (context.Context, error) {
	md := metadata.Pairs(IotTraceHeader, traceID)
	if err := grpc.SetHeader(ctx, md); err != nil {
		return nil, err
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(IotTraceHeader, traceID)), nil
}
