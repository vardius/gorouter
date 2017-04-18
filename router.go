package goserver

import (
	"net/http"
)

type (
	route struct {
		middleware  middleware
		handlerFunc http.HandlerFunc
	}
)

func (r *route) handler() http.Handler {
	if r.handlerFunc != nil {
		return r.middleware.handleFunc(r.handlerFunc)
	}
	return nil
}

func (r *route) addMiddleware(m middleware) {
	r.middleware = r.middleware.merge(m)
}

func newRoute(h http.HandlerFunc) *route {
	return &route{
		handlerFunc: h,
		middleware:  newMiddleware(),
	}
}
