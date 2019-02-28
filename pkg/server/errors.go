package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"minibox.ai/pkg/server/internal/acl"
)

var (
	ErrInvalidUser        = status.Error(codes.Internal, "invalid user")
	ErrInvalidObjectStore = status.Error(codes.Internal, "invalid object store")
	ErrNotFoundResource   = status.Error(codes.NotFound, "not found this resource")
	ErrDenyPrivateProject = status.Error(codes.PermissionDenied, "can't access private Project")
)

func ErrInvalidPermission(act acl.Action) error {
	return status.Errorf(codes.PermissionDenied, "Dont has '%s' permssion", act)
}

func ErrInvalidNamespace(name string) error {
	return status.Errorf(codes.InvalidArgument, "Invalid namespace '%s'", name)
}
