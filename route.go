package gorouter

import (
	"net/http"

	"github.com/vardius/gorouter/v4/middleware"
)

type route struct {
	middleware middleware.Middleware
	handler    http.Handler
}

func (r *route) getHandler() http.Handler {
	if r.handler != nil {
		return r.middleware.Handle(r.handler)
	}
	return nil
}

func (r *route) appendMiddleware(m middleware.Middleware) {
	r.middleware = r.middleware.Merge(m)
}

func (r *route) prependMiddleware(m middleware.Middleware) {
	r.middleware = m.Merge(r.middleware)
}

func newRoute(h http.Handler) *route {
	return &route{
		handler:    h,
		middleware: middleware.New(),
	}
}
