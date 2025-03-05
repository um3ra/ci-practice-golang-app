package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate() error
}

func ValidateInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if val, ok := req.(validator); ok {
		if err := val.Validate(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
		}
	}

	return handler(ctx, req)
}
