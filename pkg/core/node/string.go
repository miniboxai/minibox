package node

import (
	"context"
	"reflect"

	"minibox.ai/minibox/pkg/core/utils"
)

type NodeString struct {
	Node
	tmpl string
}

func (n *NodeString) Compile(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return &ErrInvalidString{Name: n.Name}
	}

	n.tmpl = s
	return nil
}

func (n *NodeString) Eval(ctx context.Context) (reflect.Value, error) {
	v := reflect.ValueOf(n.eval(ctx))
	n.setBound(v)
	return v, nil
}

func (n *NodeString) eval(ctx context.Context) interface{} {
	cfgCtx := utils.ContextFrom(ctx)
	return utils.FmtString(n.tmpl, cfgCtx.Map())
}
