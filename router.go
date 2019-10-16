package gorouter

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

// MiddlewareFunc is a http middleware function type
type MiddlewareFunc func(http.Handler) http.Handler

// FastHTTPMiddlewareFunc is a fasthttp middleware function type
type FastHTTPMiddlewareFunc func(fasthttp.RequestHandler) fasthttp.RequestHandler

// Router is a micro framework, HTTP request router, multiplexer, mux
type Router interface {
	// POST adds http.Handler as router handler
	// under POST method and given patter
	POST(pattern string, handler http.Handler)

	// GET adds http.Handler as router handler
	// under GET method and given patter
	GET(pattern string, handler http.Handler)

	// PUT adds http.Handler as router handler
	// under PUT method and given patter
	PUT(pattern string, handler http.Handler)

	// DELETE adds http.Handler as router handler
	// under DELETE method and given patter
	DELETE(pattern string, handler http.Handler)

	// PATCH adds http.Handler as router handler
	// under PATCH method and given patter
	PATCH(pattern string, handler http.Handler)

	// OPTIONS adds http.Handler as router handler
	// under OPTIONS method and given patter
	OPTIONS(pattern string, handler http.Handler)

	// HEAD adds http.Handler as router handler
	// under HEAD method and given patter
	HEAD(pattern string, handler http.Handler)

	// CONNECT adds http.Handler as router handler
	// under CONNECT method and given patter
	CONNECT(pattern string, handler http.Handler)

	// TRACE adds http.Handler as router handler
	// under TRACE method and given patter
	TRACE(pattern string, handler http.Handler)

	// USE adds middleware functions ([]MiddlewareFunc)
	// to whole router branch under given method and patter
	USE(method, pattern string, fs ...MiddlewareFunc)

	// Handle adds http.Handler as router handler
	// under given method and patter
	Handle(method, pattern string, handler http.Handler)

	// Mount another handler as a subrouter
	Mount(pattern string, handler http.Handler)

	// Compile optimizes Tree nodes reducing static nodes depth when possible
	Compile()

	// ServeHTTP dispatches the request to the route handler
	// whose pattern matches the request URL
	ServeHTTP(http.ResponseWriter, *http.Request)

	// ServeFile replies to the request with the
	// contents of the named file or directory.
	ServeFiles(fs http.FileSystem, root string, strip bool)

	// NotFound replies to the request with the
	// 404 Error code
	NotFound(http.Handler)

	// NotFound replies to the request with the
	// 405 Error code
	NotAllowed(http.Handler)
}

// FastHTTPRouter is a fasthttp micro framework, HTTP request router, multiplexer, mux
type FastHTTPRouter interface {
	// POST adds fasthttp.RequestHandler as router handler
	// under POST method and given patter
	POST(pattern string, handler fasthttp.RequestHandler)

	// GET adds fasthttp.RequestHandler as router handler
	// under GET method and given patter
	GET(pattern string, handler fasthttp.RequestHandler)

	// PUT adds fasthttp.RequestHandler as router handler
	// under PUT method and given patter
	PUT(pattern string, handler fasthttp.RequestHandler)

	// DELETE adds fasthttp.RequestHandler as router handler
	// under DELETE method and given patter
	DELETE(pattern string, handler fasthttp.RequestHandler)

	// PATCH adds fasthttp.RequestHandler as router handler
	// under PATCH method and given patter
	PATCH(pattern string, handler fasthttp.RequestHandler)

	// OPTIONS adds fasthttp.RequestHandler as router handler
	// under OPTIONS method and given patter
	OPTIONS(pattern string, handler fasthttp.RequestHandler)

	// HEAD adds fasthttp.RequestHandler as router handler
	// under HEAD method and given patter
	HEAD(pattern string, handler fasthttp.RequestHandler)

	// CONNECT adds fasthttp.RequestHandler as router handler
	// under CONNECT method and given patter
	CONNECT(pattern string, handler fasthttp.RequestHandler)

	// TRACE adds fasthttp.RequestHandler as router handler
	// under TRACE method and given patter
	TRACE(pattern string, handler fasthttp.RequestHandler)

	// USE adds middleware functions ([]MiddlewareFunc)
	// to whole router branch under given method and patter
	USE(method, pattern string, fs ...FastHTTPMiddlewareFunc)

	// Handle adds fasthttp.RequestHandler as router handler
	// under given method and patter
	Handle(method, pattern string, handler fasthttp.RequestHandler)

	// Mount another handler as a subrouter
	Mount(pattern string, handler fasthttp.RequestHandler)

	// Compile optimizes Tree nodes reducing static nodes depth when possible
	Compile()

	// HandleFastHTTP dispatches the request to the route handler
	// whose pattern matches the request URL
	HandleFastHTTP(ctx *fasthttp.RequestCtx)

	// ServeFile replies to the request with the
	// contents of the named file or directory.
	ServeFiles(root string, stripSlashes int)

	// NotFound replies to the request with the
	// 404 Error code
	NotFound(fasthttp.RequestHandler)

	// NotFound replies to the request with the
	// 405 Error code
	NotAllowed(fasthttp.RequestHandler)
}
