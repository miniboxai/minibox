package acl

import "errors"

var (
	ErrInvalidRoleString  = errors.New("invalid role string format")
	ErrAlwaysRegisterRole = errors.New("always register role name")
)
