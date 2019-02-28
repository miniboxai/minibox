package v1

import (
	"context"

	"minibox.ai/pkg/core/node"
	"minibox.ai/pkg/core/utils"
)

type Nodes map[string]node.Noder

func (ns Nodes) Eval(ctx context.Context) error {
	if err := ns.evalGroup(prioriEvals, ctx); err != nil {
		return err
	}

	var nkeys []string
	for k, _ := range ns {
		nkeys = append(nkeys, k)
	}
	restEvels := minusSlice(nkeys, prioriEvals)
	if err := ns.evalGroup(restEvels, ctx); err != nil {
		return err
	}
	return nil
}

func (ns Nodes) evalGroup(evals []string, ctx context.Context) error {
	for _, pe := range evals {
		if n, ok := ns[pe]; ok {
			if v, err := n.Eval(ctx); err != nil {
				return err
			} else {
				cfgCtx := utils.ContextFrom(ctx)
				cfgCtx.Set(pe, v)
			}
		}
	}
	return nil
}

func compileString(v interface{}) *node.NodeString {
	var n node.NodeString
	n.Compile(v)
	return &n
}

func compileFunc(v interface{}, handle func(ctnt string, ctx context.Context) interface{}) *node.NodeFunc {
	var n node.NodeFunc
	n.Compile(v)
	n.Func = handle
	return &n
}

func compileStruct(v interface{}) *node.NodeStruct {
	var n node.NodeStruct
	n.Compile(v)
	return &n

}

func compileStringSlice(v []interface{}) *node.NodeStringSlice {
	var n node.NodeStringSlice
	var ret = make([]string, 0, len(v))
	for _, m := range v {
		if s, ok := m.(string); ok {
			ret = append(ret, s)
		}
	}

	n.Compile(ret)
	return &n
}

func compileEnvs(val []interface{}) *node.NodeFunc {
	var zero = []Env{}

	return compileFunc(zero, func(ctnt string, ctx context.Context) interface{} {

		ess := toStringSlice(val)
		var envs = make([]Env, 0, len(ess))
		for _, cenv := range ess {
			val := compileTmpl(cenv, ctx)
			if env, err := ParseEnv(val); err != nil {
				panic(err)
			} else {
				envs = append(envs, *env)
			}
		}
		return envs
	})
}
