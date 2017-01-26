package payment

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Errors
var (
	ErrPaymentNotFound = grpc.Errorf(codes.NotFound, "payment: not found")
)
