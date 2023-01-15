package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Unauthenticated ...
func Unauthenticated() error {
	return status.Error(codes.Unauthenticated, "unauthenticated")
}

// Internal ...
func Internal() error {
	return status.Error(codes.Internal, "internal")
}

// NotFound ...
func NotFound() error {
	return status.Error(codes.NotFound, "not found")
}

// BadRequest ...
func BadRequest() error {
	return status.Errorf(codes.InvalidArgument, "bad request")
}
