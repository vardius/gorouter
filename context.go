package goserver

import (
	"context"
	"net/http"
)

type (
	key int
	//Params type returned from request context
	Params map[string]string
)

const paramsKey key = 0

func newContextFromRequest(req *http.Request, params Params) context.Context {
	return context.WithValue(req.Context(), paramsKey, params)
}

func ParamsFromContext(ctx context.Context) (Params, bool) {
	params, ok := ctx.Value(paramsKey).(Params)
	return params, ok
}
