package goserver

import (
	"net/http"
)

type (
	route struct {
		middleware middleware
		handler    http.Handler
	}
	// Param object to hold request parameter
	Param struct {
		Key   string
		Value string
	}
	//Params slice returned from request context
	Params []Param
)

//Value of the request parameter by name
func (p Params) Value(key string) string {
	for i := range p {
		if p[i].Key == key {
			return p[i].Value
		}
	}
	return ""
}

func (r *route) chain() http.Handler {
	if r.handler != nil {
		return r.middleware.handle(r.handler)
	}
	return nil
}

func (r *route) addMiddleware(m middleware) {
	r.middleware = r.middleware.merge(m)
}

func newRoute(h http.Handler) *route {
	return &route{
		handler:    h,
		middleware: newMiddleware(),
	}
}
