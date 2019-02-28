package utils

import (
	"context"

	"github.com/delicb/gstring"
	pcontext "minibox.ai/minibox/pkg/core/context"
	"minibox.ai/minibox/pkg/core/errors"
)

var ContextKey = struct{}{}

func ContextFrom(ctx context.Context) *pcontext.Context {
	kv := ctx.Value(ContextKey)
	if kv == nil {
		panic(errors.ErrMissingContext)
	}

	if cfgCtx, ok := kv.(*pcontext.Context); !ok {
		panic(errors.ErrMissingContext)
	} else {
		return cfgCtx
	}
}

func FmtString(tmpl string, ctx map[string]interface{}) string {
	return gstring.Sprintm(tmpl, ctx)
}
