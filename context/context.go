package context

import (
	"context"
)

type key struct{}

// WithParams stores params in context
func WithParams(ctx context.Context, params Params) context.Context {
	return context.WithValue(ctx, key{}, params)
}

// Parameters extracts the request Params ctx, if present.
func Parameters(ctx context.Context) (Params, bool) {
	params, ok := ctx.Value(key{}).(Params)
	return params, ok
}
