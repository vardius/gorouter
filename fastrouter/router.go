package fastrouter

import (
	"slices"
	"strings"
	"sync"
	"unsafe"

	"github.com/valyala/fasthttp"

	"github.com/vardius/gorouter/v5/internal/tree"
)

// stackParams is how many param values a request can carry before matching needs a heap buffer.
const stackParams = 8

// Middleware wraps a handler.
type Middleware = func(fasthttp.RequestHandler) fasthttp.RequestHandler

// Route describes a registered route.
type Route struct {
	Method  string
	Pattern string
}

// Option configures a Router.
type Option func(*config)

type config struct {
	middleware []Middleware
	notFound   fasthttp.RequestHandler
	notAllowed fasthttp.RequestHandler
}

// WithMiddleware wraps the whole router, including 404, 405 and automatic OPTIONS responses.
// The first middleware is the outermost.
func WithMiddleware(mw ...Middleware) Option {
	return func(c *config) { c.middleware = append(c.middleware, mw...) }
}

// WithNotFound replaces the default 404 handler.
func WithNotFound(h fasthttp.RequestHandler) Option {
	return func(c *config) { c.notFound = h }
}

// WithMethodNotAllowed replaces the default 405 handler. The Allow header is already set when it runs,
// so write the response without ctx.Error, which resets the headers.
func WithMethodNotAllowed(h fasthttp.RequestHandler) Option {
	return func(c *config) { c.notAllowed = h }
}

// Router is a fasthttp request router. Register routes and middleware first, then serve:
// registering while serving is not safe.
type Router struct {
	tree       tree.Tree[fasthttp.RequestHandler]
	notFound   fasthttp.RequestHandler
	notAllowed fasthttp.RequestHandler
	handler    fasthttp.RequestHandler
}

// New returns a Router configured by opts.
func New(opts ...Option) *Router {
	c := config{
		notFound:   errorHandler(fasthttp.StatusNotFound),
		notAllowed: errorHandler(fasthttp.StatusMethodNotAllowed),
	}
	for _, opt := range opts {
		opt(&c)
	}

	r := &Router{notFound: c.notFound, notAllowed: c.notAllowed}
	r.handler = r.dispatch
	for _, mw := range slices.Backward(c.middleware) {
		r.handler = mw(r.handler)
	}

	return r
}

// GET registers h for GET requests to pattern.
func (r *Router) GET(pattern string, h fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodGet, pattern, h)
}

// HEAD registers h for HEAD requests to pattern.
func (r *Router) HEAD(pattern string, h fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodHead, pattern, h)
}

// POST registers h for POST requests to pattern.
func (r *Router) POST(pattern string, h fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodPost, pattern, h)
}

// PUT registers h for PUT requests to pattern.
func (r *Router) PUT(pattern string, h fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodPut, pattern, h)
}

// PATCH registers h for PATCH requests to pattern.
func (r *Router) PATCH(pattern string, h fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodPatch, pattern, h)
}

// DELETE registers h for DELETE requests to pattern.
func (r *Router) DELETE(pattern string, h fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodDelete, pattern, h)
}

// CONNECT registers h for CONNECT requests to pattern.
func (r *Router) CONNECT(pattern string, h fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodConnect, pattern, h)
}

// OPTIONS registers h for OPTIONS requests to pattern, replacing the automatic response there.
func (r *Router) OPTIONS(pattern string, h fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodOptions, pattern, h)
}

// TRACE registers h for TRACE requests to pattern.
func (r *Router) TRACE(pattern string, h fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodTrace, pattern, h)
}

// Handle registers h for method and pattern. Any method works, including custom ones.
//
// A pattern is a path whose segments are literals, {name} params matching any one segment,
// or {name:regexp} params matching a segment the whole regexp matches. A trailing slash is ignored.
// Handle panics if the pattern is invalid or already registered for method.
func (r *Router) Handle(method, pattern string, h fasthttp.RequestHandler) {
	if h == nil {
		panic("gorouter: nil handler for " + method + " " + pattern)
	}

	r.tree.Handle(method, pattern, h)
}

// Use wraps every route, on every method, whose pattern starts with pattern, segment by segment.
// Routes registered before and after the call are covered. Middleware runs in registration order.
func (r *Router) Use(pattern string, mw ...Middleware) {
	r.tree.Use("", pattern, mw)
}

// UseMethod is Use restricted to one method.
func (r *Router) UseMethod(method, pattern string, mw ...Middleware) {
	r.tree.Use(method, pattern, mw)
}

// Mount sends every request at or below pattern, on any method, to h with the pattern's segments removed
// from the URI path. Params in pattern stay readable with Param. Once h returns, the path and params
// are put back, so middleware around the mount sees the request it was given.
//
// Use and UseMethod middleware on pattern or above it wraps the mount. Middleware below the mount
// belongs in h, and registering it here panics.
func (r *Router) Mount(pattern string, h fasthttp.RequestHandler) {
	if h == nil {
		panic("gorouter: nil handler for mount at " + pattern)
	}

	r.tree.Mount(pattern, stripSegments(tree.Depth(pattern), h))
}

// Routes lists the registered routes in registration order. A mount is listed once, with an empty Method.
func (r *Router) Routes() []Route {
	var routes []Route
	for route := range r.tree.Routes() {
		routes = append(routes, Route{Method: route.Method, Pattern: route.Pattern})
	}

	return routes
}

// HandleFastHTTP dispatches ctx to the matching route. Read param values with Param.
// Pass it to fasthttp as r.HandleFastHTTP.
func (r *Router) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	r.handler(ctx)
}

func (r *Router) dispatch(ctx *fasthttp.RequestCtx) {
	// Views of fasthttp's buffers: matching only reads them, so neither is copied.
	method := bytesToString(ctx.Method())
	path := bytesToString(ctx.Path())
	if method == fasthttp.MethodOptions && string(ctx.RequestURI()) == "*" {
		path = "*"
	}

	// Param values go to a buffer on the stack, or to a pooled one when a route has more params than it holds.
	var buf [stackParams]string
	values := buf[:]
	var pooled *[]string
	if r.tree.MaxParams() > stackParams {
		pooled = r.tree.Values()
		values = *pooled
	}

	if route := r.tree.Lookup(method, path, values); route != nil {
		if len(route.Params) > 0 {
			setParams(ctx, path, route.Params, values)
		}
		route.Handler(ctx)
	} else if allow := r.tree.Allowed(method, path, values); allow != "" {
		ctx.Response.Header.Set(fasthttp.HeaderAllow, allow)
		if method != fasthttp.MethodOptions {
			r.notAllowed(ctx)
		}
	} else {
		r.notFound(ctx)
	}

	if pooled != nil {
		r.tree.PutValues(pooled)
	}
}

// errorHandler writes status and its message. Unlike ctx.Error it keeps the headers already set,
// such as Allow or ones set by middleware.
func errorHandler(status int) fasthttp.RequestHandler {
	msg := fasthttp.StatusMessage(status)
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(status)
		ctx.SetContentType("text/plain; charset=utf-8")
		ctx.SetBodyString(msg)
	}
}

// mountPaths holds a mount's copy of the caller's path and the stripped path it hands on.
type mountPaths struct {
	caller, stripped []byte
}

var mountPathsPool = sync.Pool{New: func() any { return new(mountPaths) }}

// stripSegments serves h with the first depth segments removed from the URI path,
// then puts the caller's path and params back.
func stripSegments(depth int, h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		uri := ctx.URI()
		paths := mountPathsPool.Get().(*mountPaths)
		paths.caller = append(paths.caller[:0], uri.PathOriginal()...)
		paths.stripped = appendPath(paths.stripped[:0], tree.StripSegments(bytesToString(ctx.Path()), depth))
		mark := markParams(ctx)

		uri.SetPathBytes(paths.stripped)
		h(ctx)
		uri.SetPathBytes(paths.caller)

		resetParams(ctx, mark)
		mountPathsPool.Put(paths)
	}
}

// appendPath appends path, already decoded, escaped so that SetPath decodes it back to path:
// it percent-decodes its argument, and path's own % signs must survive that.
func appendPath(dst []byte, path string) []byte {
	for {
		i := strings.IndexByte(path, '%')
		if i < 0 {
			return append(dst, path...)
		}
		dst = append(dst, path[:i+1]...)
		dst = append(dst, "25"...)
		path = path[i+1:]
	}
}

// bytesToString views b as a string without copying. The string is only valid while b is unchanged.
func bytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
