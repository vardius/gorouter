package goserver

import (
	"net/http"
)

type (
	route struct {
		middleware  middleware
		handlerFunc http.HandlerFunc
	}
	// Param object to hold request parameter
	Param struct {
		Key   string
		Value string
	}
	//Params slice returned from request context
	Params []Param
)

//Get Params value by key
func (p Params) Value(key string) string {
	for i := range p {
		if p[i].Key == key {
			return p[i].Value
		}
	}
	return ""
}

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
