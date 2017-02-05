package app

import (
	"github.com/acoshift/acourse/pkg/internal"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// RecoveryUnaryInterceptor is the interceptor for panic recovery
func RecoveryUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = grpc.Errorf(codes.Internal, "%v", r)
		}
	}()
	return handler(ctx, req)
}

// AuthUnaryInterceptor is the interceptor for auththentication
func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return handler(ctx, req)
	}
	auth := md["authorization"]
	// if not provide authorization header lets it pass
	if len(auth) == 0 {
		return handler(ctx, req)
	}

	userID, err := validateHeaderToken(auth[0])
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "invalid autorization header")
	}
	rctx := internal.WithUserID(ctx, userID)
	return handler(rctx, req)
}

// UnaryInterceptors is the chain interceptors for server
func UnaryInterceptors(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return RecoveryUnaryInterceptor(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return AuthUnaryInterceptor(ctx, req, info, handler)
	})
}
