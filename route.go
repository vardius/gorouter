package gorouter

import (
	"fmt"
	"github.com/vardius/gorouter/v4/middleware"
)

type route struct {
	path       string
	middleware middleware.Middleware
	handler    interface{}
	// computedHandler is an optimization to improve performance
	computedHandler interface{}
}

func newRoute(h interface{}) *route {
	if h == nil {
		panic("Handler can not be nil.")
	}

	return &route{
		handler:    h,
		middleware: middleware.New(),
	}
}

func (r *route) Handler(path string) interface{} {
	fmt.Printf("Handler path: %v\n", path)
	fmt.Printf("Handler mapped: %v\n", r.path)
	if r.path == path {
		r.ComposeMiddleware(r.middleware)
	}
	// returns already cached computed handler
	return r.computedHandler
}

func (r *route) AppendMiddleware(m middleware.Middleware, path string) {
	r.path = path
	r.middleware = r.middleware.Merge(m)
	//r.computedHandler = r.middleware.Compose(r.handler)
}

func (r *route) PrependMiddleware(m middleware.Middleware) {
	r.middleware = m.Merge(r.middleware)
	r.computedHandler = r.middleware.Compose(r.handler)
}

func (r *route) ComposeMiddleware(m middleware.Middleware) {
	r.computedHandler = r.middleware.Compose(r.handler)
}
