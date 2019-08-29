package mux

import (
	"github.com/vardius/gorouter/v4/middleware"
)

// Route is an middleware aware route interface
type Route interface {
	Handler() interface{}
	AppendMiddleware(m middleware.Middleware)
	PrependMiddleware(m middleware.Middleware)
}
