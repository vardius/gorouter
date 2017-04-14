package goserver

import "net/http"

type (
	MiddlewareFunc func(http.Handler) http.Handler
	middlewares    []MiddlewareFunc
)

func (m middlewares) handle(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range m {
		h = m[len(m)-1-i](h)
	}

	return h
}

func (m middlewares) handleFunc(f http.HandlerFunc) http.Handler {
	return m.handle(f)
}
