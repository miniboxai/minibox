package errors

import (
	"errors"
)

var ErrNonImplement = errors.New("this feature is non implement")
var ErrAlreadyRecord = errors.New("already have a some record")
