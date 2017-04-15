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

func newMiddleware(fs ...MiddlewareFunc) middlewares {
	ms := make(middlewares, 0, len(fs))
	for _, f := range fs {
		if f != nil {
			ms = append(ms, f)
		}
	}

	return ms
}
