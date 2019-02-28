package node

import (
	"context"
	"errors"
	"reflect"
)

type NodeFunc struct {
	NodeString
	Func func(string, context.Context) interface{}
}

func (n *NodeFunc) Eval(ctx context.Context) (reflect.Value, error) {
	val := n.NodeString.eval(ctx)
	ctnt, ok := val.(string)

	if !ok {
		return reflect.Value{}, errors.New("can't convert to string")
	}

	v := reflect.ValueOf(n.Func(ctnt, ctx))
	if err := n.setBound(v); err != nil {
		return reflect.Value{}, err
	}
	return v, nil
}
