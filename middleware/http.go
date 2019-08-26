package middleware

import (
	"net/http"
	"reflect"
)

// Handle provides http handler
func (m Middleware) Handle(handler http.Handler) http.Handler {
	if handler == nil {
		handler = http.DefaultServeMux
	}

	for i := range m {
		args := []reflect.Value{reflect.ValueOf(handler)}
		wrapper := m[len(m)-1-i]
		handler = wrapper.Call(args)[0].Interface().(http.Handler)
	}

	return handler
}

// HandleFunc provides http handler from http.HandlerFunc
func (m Middleware) HandleFunc(f http.HandlerFunc) http.Handler {
	return m.Handle(f)
}
