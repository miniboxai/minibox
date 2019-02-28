package v1

import (
	"errors"
	"fmt"
)

var ErrMissingFramework = errors.New("missing `framework` setting")

type ErrNotSupportFramework struct {
	Name string
}

type ErrInvalidStruct struct {
	Section string
	val     interface{}
	Suggest string
}

type ErrInvalidEnvString struct {
	Str string
}

func (err *ErrNotSupportFramework) Error() string {
	return fmt.Sprintf("dont' have this `%s` Framework support", err.Name)
}

func (err *ErrInvalidStruct) Error() string {
	var s string
	if empty(err.Suggest) {
		s = fmt.Sprintf("invalid section `%s`, type :%T\n", err.Section, err.val)
	} else {
		s = fmt.Sprintf("invalid section `%s`, type :%T, only %s",
			err.Section, err.val, err.Suggest)
	}
	return s
}

func (err *ErrInvalidEnvString) Error() string {
	return fmt.Sprintf("invalid Env string format `%s`, must has `Name=Value`", err.Str)
}
