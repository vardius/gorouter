package gorouter

import (
	"context"
	"net/http"
)

type key int

const paramsKey key = 0

func newContext(req *http.Request, params Params) context.Context {
	return context.WithValue(req.Context(), paramsKey, params)
}

//FromContext extracts the request Params ctx, if present.
func FromContext(ctx context.Context) (Params, bool) {
	params, ok := ctx.Value(paramsKey).(Params)
	return params, ok
}
