package node

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotHaveBound           = errors.New("this Node not have bound")
	ErrCanNotBeSet            = errors.New("this Node can't be set new variable")
	ErrCanNotSetDifferentType = errors.New("can't assign between two Nodes different type")
	ErrMustPointer            = errors.New("must a pointer variable ")
	ErrMissingContext         = errors.New("missing contextKey key of Context")
)

type ErrInvalidConfigItem struct {
	Name string
}

type ErrInvalidStringSlice struct {
	Name string
}

type ErrInvalidString struct {
	Name string
}

type ErrInvalidItemType struct {
	Name string
	Type reflect.Type
}

func (e *ErrInvalidConfigItem) Error() string {
	return fmt.Sprintf("does not exist '%s' config item", e.Name)
}

func (e *ErrInvalidStringSlice) Error() string {
	return fmt.Sprintf("does not exist '%s' []string type slice item", e.Name)
}

func (e *ErrInvalidString) Error() string {
	return fmt.Sprintf("does not exist '%s' string type item", e.Name)
}

func (e *ErrInvalidItemType) Error() string {
	return fmt.Sprintf("can't get item '%s' %s type item", e.Name, e.Type)
}
