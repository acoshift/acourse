package app

import (
	"github.com/acoshift/acourse/pkg/acourse"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// AuthUnaryInterceptor is the interceptor for auththentication
func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
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
	rctx := context.WithValue(ctx, acourse.KeyUserID, userID)
	return handler(rctx, req)
}
