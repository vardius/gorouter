package gorouter

import (
	"github.com/vardius/gorouter/v4/middleware"
)

type route struct {
	middleware middleware.Middleware
	handler    interface{}
	// computedHandler is an optimization to improve performance
	computedHandler interface{}
}

func newRoute(h interface{}) *route {
	return &route{
		handler:    h,
		middleware: middleware.New(),
	}
}

func (r *route) Handler() interface{} {
	// returns already cached computed handler
	return r.computedHandler
}

func (r *route) AppendMiddleware(m middleware.Middleware) {
	r.middleware = r.middleware.Merge(m)
	r.computedHandler = r.middleware.Compose(r.handler)
}

func (r *route) PrependMiddleware(m middleware.Middleware) {
	r.middleware = m.Merge(r.middleware)
	r.computedHandler = r.middleware.Compose(r.handler)
}
