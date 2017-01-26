package assignment

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Errors
var (
	ErrAssignmentNotFound     = grpc.Errorf(codes.NotFound, "assignment: not found")
	ErrAssignmentNotOpen      = grpc.Errorf(codes.FailedPrecondition, "assignment: not open")
	ErrUserAssignmentNotFound = grpc.Errorf(codes.NotFound, "assignment: user assignment not found")
	ErrPermissionDenied       = grpc.Errorf(codes.PermissionDenied, "assignment: user do not have permission")
)
