package goserver

import (
	"context"
	"net/http"
)

type key int

const paramsKey key = 0

func newContextFromRequest(req *http.Request, params parameters) context.Context {
	return context.WithValue(req.Context(), paramsKey, params)
}

func ParametersFromContext(ctx context.Context) (parameters, bool) {
	params, ok := ctx.Value(paramsKey).(parameters)
	return params, ok
}
