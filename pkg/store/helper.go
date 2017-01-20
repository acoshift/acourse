package store

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Errors
var (
	ErrNotFound = grpc.Errorf(codes.NotFound, "not found")
	ErrInternal = grpc.Errorf(codes.Internal, "internal")
)

func errInternalWith(err error) error {
	return grpc.Errorf(codes.Internal, err.Error())
}
