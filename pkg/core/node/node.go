package node

import (
	"context"
	"errors"
	"reflect"
)

type Noder interface {
	Eval(context.Context) (reflect.Value, error)
	// BoundTo(interface{}) Noder
}

type Node struct {
	Name    string
	Type    reflect.Type
	bounder *reflect.Value
	// children []Noder
	// Func    func(ctx context.Context, args ...interface{}) interface{}
}

func (n *Node) setBound(val reflect.Value) error {
	if n.bounder == nil {
		return ErrNotHaveBound
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	bounder := n.bounder.Elem()

	if !bounder.CanSet() {
		return ErrCanNotBeSet
	}

	if bounder.Type() != val.Type() {
		return ErrCanNotSetDifferentType
	}

	bounder.Set(val)
	return nil
}

func (n *Node) Named(name string) *Node {
	n.Name = name
	return n
}

func (n *Node) BoundTo(val interface{}) *Node {
	v := reflect.ValueOf(val)

	if v.Kind() != reflect.Ptr {
		panic(ErrMustPointer)
	}

	n.bounder = &v
	return n
}

func (n *Node) Eval(ctx context.Context) (reflect.Value, error) {
	return reflect.Value{}, errors.New("Node Base can't eval")
}
