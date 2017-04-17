package goserver

import "net/http"

type (
	middlewareFunc func(http.Handler) http.Handler
	middlewares    []middlewareFunc
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

func (m middlewares) append(mf ...middlewareFunc) {
	m = append(m, mf...)
}

func (m middlewares) merge(n middlewares) {
	m = append(m, n...)
}

func newMiddleware(fs ...middlewareFunc) middlewares {
	ms := make(middlewares, 0, len(fs))
	for _, f := range fs {
		if f != nil {
			ms = append(ms, f)
		}
	}

	return ms
}
