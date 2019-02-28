package node

import (
	"context"
	"reflect"
)

type NodeStruct struct {
	Node
	val interface{}
}

func (n *NodeStruct) Compile(v interface{}) error {
	// s, ok := v.(string)
	// if !ok {
	// 	return &ErrInvalidString{Name: n.Name}
	// }

	n.val = v
	return nil
}

func (n *NodeStruct) Eval(ctx context.Context) (reflect.Value, error) {
	v := reflect.ValueOf(n.val)
	n.setBound(v)
	return v, nil
}
