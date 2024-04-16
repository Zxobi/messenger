package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateId(id []byte, name string) error {
	if len(id) == 0 {
		return status.Error(codes.InvalidArgument, name+" is required")
	}
	if len(id) != 16 {
		return status.Error(codes.InvalidArgument, name+" is bad")
	}

	return nil
}
