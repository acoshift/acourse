package user

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Errors
var (
	ErrUserIDRequired   = grpc.Errorf(codes.InvalidArgument, "user: id required")
	ErrUserNameConflict = grpc.Errorf(codes.AlreadyExists, "user: user name already exists")
)
