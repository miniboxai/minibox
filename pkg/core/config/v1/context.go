package v1

import (
	"context"

	pcontext "minibox.ai/minibox/pkg/core/context"
)

func NewContext(prjName string, initials map[string]interface{}) context.Context {
	var initialParams = map[string]interface{}{
		"projectName": prjName,
	}

	initialParams = mergeMap(initialParams, initials)

	ctx := context.WithValue(
		context.Background(),
		pcontext.ContextKey,
		pcontext.NewContext(initialParams))

	return ctx
}

func mergeMap(dst, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		dst[k] = v
	}

	return dst
}
