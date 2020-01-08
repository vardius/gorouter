package mux

import (
	"github.com/vardius/gorouter/v4/middleware"
)

// Route is an middleware aware route interface
type Route interface {
	Handler(path string) interface{}
	AppendMiddleware(m middleware.Middleware, path string)
	PrependMiddleware(m middleware.Middleware)
	ComposeMiddleware(m middleware.Middleware)
}
