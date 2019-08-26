package middleware

import (
	"fmt"
	"reflect"
)

// MiddlewareFunc is a middleware function type.
// Long story - short: it is a handler wrapper
type MiddlewareFunc interface{}

// Middleware is a slice of handler functions
type Middleware []reflect.Value

// New provides new middleware
func New(fs ...MiddlewareFunc) Middleware {
	m := make(Middleware, 0, len(fs))

	return m.Append(fs...)
}

// Append appends handlers to middlewares
func (m Middleware) Append(fs ...MiddlewareFunc) Middleware {
	for _, f := range fs {
		if f != nil {
			if err := ensureMiddlewareIsAFunc(f); err != nil {
				panic(err)
			}

			m = append(m, reflect.ValueOf(f))
		}
	}

	return m
}

// Merge merges another middlewares
func (m Middleware) Merge(n Middleware) Middleware {
	return append(m, n...)
}

func ensureMiddlewareIsAFunc(fn MiddlewareFunc) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	return nil
}
