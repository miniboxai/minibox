package node

import (
	"context"
	"reflect"
)

type NodeStringSlice struct {
	Node
	val []string

	// tmpl string
}

func (n *NodeStringSlice) Compile(v interface{}) error {
	s, ok := v.([]string)
	if !ok {
		return &ErrInvalidItemType{Name: n.Name, Type: reflect.TypeOf([]string{})}
	}
	n.val = s
	return nil
}

func (n *NodeStringSlice) Eval(ctx context.Context) (reflect.Value, error) {
	v := reflect.ValueOf(n.val)
	n.setBound(v)
	return v, nil
}
