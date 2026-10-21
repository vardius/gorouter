package gorouter

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/vardius/gorouter/v5/internal/tree"
)

// stackParams is how many param values a request can carry before matching needs a heap buffer.
const stackParams = 8

// Middleware wraps a handler. It has the same shape as most net/http middleware.
type Middleware = func(http.Handler) http.Handler

// Route describes a registered route.
type Route struct {
	Method  string
	Pattern string
}

// Option configures a Router.
type Option func(*config)

type config struct {
	middleware []Middleware
	notFound   http.Handler
	notAllowed http.Handler
}

// WithMiddleware wraps the whole router, including 404, 405 and automatic OPTIONS responses.
// The first middleware is the outermost.
func WithMiddleware(mw ...Middleware) Option {
	return func(c *config) { c.middleware = append(c.middleware, mw...) }
}

// WithNotFound replaces the default 404 handler.
func WithNotFound(h http.Handler) Option {
	return func(c *config) { c.notFound = h }
}

// WithMethodNotAllowed replaces the default 405 handler. The Allow header is already set when it runs.
func WithMethodNotAllowed(h http.Handler) Option {
	return func(c *config) { c.notAllowed = h }
}

// Router is a net/http request router. Register routes and middleware first, then serve:
// registering while serving is not safe.
type Router struct {
	tree       tree.Tree[http.Handler]
	notFound   http.Handler
	notAllowed http.Handler
	handler    http.Handler
}

// New returns a Router configured by opts.
func New(opts ...Option) *Router {
	c := config{
		notFound:   http.NotFoundHandler(),
		notAllowed: http.HandlerFunc(methodNotAllowed),
	}
	for _, opt := range opts {
		opt(&c)
	}

	r := &Router{notFound: c.notFound, notAllowed: c.notAllowed}
	r.handler = http.HandlerFunc(r.dispatch)
	for _, mw := range slices.Backward(c.middleware) {
		r.handler = mw(r.handler)
	}

	return r
}

// GET registers h for GET requests to pattern.
func (r *Router) GET(pattern string, h http.HandlerFunc) {
	r.handleFunc(http.MethodGet, pattern, h)
}

// HEAD registers h for HEAD requests to pattern.
func (r *Router) HEAD(pattern string, h http.HandlerFunc) {
	r.handleFunc(http.MethodHead, pattern, h)
}

// POST registers h for POST requests to pattern.
func (r *Router) POST(pattern string, h http.HandlerFunc) {
	r.handleFunc(http.MethodPost, pattern, h)
}

// PUT registers h for PUT requests to pattern.
func (r *Router) PUT(pattern string, h http.HandlerFunc) {
	r.handleFunc(http.MethodPut, pattern, h)
}

// PATCH registers h for PATCH requests to pattern.
func (r *Router) PATCH(pattern string, h http.HandlerFunc) {
	r.handleFunc(http.MethodPatch, pattern, h)
}

// DELETE registers h for DELETE requests to pattern.
func (r *Router) DELETE(pattern string, h http.HandlerFunc) {
	r.handleFunc(http.MethodDelete, pattern, h)
}

// CONNECT registers h for CONNECT requests to pattern.
func (r *Router) CONNECT(pattern string, h http.HandlerFunc) {
	r.handleFunc(http.MethodConnect, pattern, h)
}

// OPTIONS registers h for OPTIONS requests to pattern, replacing the automatic response there.
func (r *Router) OPTIONS(pattern string, h http.HandlerFunc) {
	r.handleFunc(http.MethodOptions, pattern, h)
}

// TRACE registers h for TRACE requests to pattern.
func (r *Router) TRACE(pattern string, h http.HandlerFunc) {
	r.handleFunc(http.MethodTrace, pattern, h)
}

// Handle registers h for method and pattern. Any method works, including custom ones.
//
// A pattern is a path whose segments are literals, {name} params matching any one segment,
// or {name:regexp} params matching a segment the whole regexp matches. A trailing slash is ignored.
// Handle panics if the pattern is invalid or already registered for method.
func (r *Router) Handle(method, pattern string, h http.Handler) {
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
// from the URL path. Params in pattern stay readable with PathValue.
//
// Use and UseMethod middleware on pattern or above it wraps the mount. Middleware below the mount
// belongs in h, and registering it here panics.
func (r *Router) Mount(pattern string, h http.Handler) {
	if h == nil {
		panic("gorouter: nil handler for mount at " + pattern)
	}

	r.tree.Mount(pattern, stripSegments(pattern, h))
}

// Routes lists the registered routes in registration order. A mount is listed once, with an empty Method.
func (r *Router) Routes() []Route {
	var routes []Route
	for route := range r.tree.Routes() {
		routes = append(routes, Route{Method: route.Method, Pattern: route.Pattern})
	}

	return routes
}

// ServeHTTP dispatches req to the matching route. Param values are set on req, read them with req.PathValue,
// and req.Pattern is set to the route's pattern.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}

func (r *Router) handleFunc(method, pattern string, h http.HandlerFunc) {
	if h == nil {
		panic("gorouter: nil handler for " + method + " " + pattern)
	}

	r.tree.Handle(method, pattern, h)
}

func (r *Router) dispatch(w http.ResponseWriter, req *http.Request) {
	// Param values go to a buffer on the stack, or to a pooled one when a route has more params than it holds.
	var buf [stackParams]string
	values := buf[:]
	var pooled *[]string
	if r.tree.MaxParams() > stackParams {
		pooled = r.tree.Values()
		values = *pooled
	}

	if route := r.tree.Lookup(req.Method, req.URL.Path, values); route != nil {
		for i, name := range route.Params {
			req.SetPathValue(name, values[i])
		}
		req.Pattern = route.Pattern
		route.Handler.ServeHTTP(w, req)
	} else if allow := r.tree.Allowed(req.Method, req.URL.Path, values); allow != "" {
		w.Header().Set("Allow", allow)
		if req.Method != http.MethodOptions {
			r.notAllowed.ServeHTTP(w, req)
		}
	} else {
		r.notFound.ServeHTTP(w, req)
	}

	if pooled != nil {
		r.tree.PutValues(pooled)
	}
}

func methodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// strippedRequest lets a mount copy the request and its URL in one allocation.
type strippedRequest struct {
	req http.Request
	url url.URL
}

// stripSegments serves h with pattern's segments removed from the URL path.
// h gets a copy of the request, as with http.StripPrefix, so the caller's URL is untouched.
// The copy shares the caller's path values, so a router in h that sets a param of the same name
// as one of pattern's overwrites it: the mount's own values are put back once h returns.
func stripSegments(pattern string, h http.Handler) http.Handler {
	depth, names := tree.Depth(pattern), tree.ParamNames(pattern)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s := &strippedRequest{req: *req, url: *req.URL}
		s.url.Path = tree.StripSegments(req.URL.Path, depth)
		if req.URL.RawPath != "" {
			// URL.EscapedPath ignores RawPath once it no longer encodes Path, as when a stripped segment held %2F.
			s.url.RawPath = tree.StripSegments(req.URL.RawPath, depth)
		}
		s.req.URL = &s.url

		var buf [stackParams]string
		values := buf[:0]
		for _, name := range names {
			values = append(values, req.PathValue(name))
		}

		h.ServeHTTP(w, &s.req)

		for i, name := range names {
			req.SetPathValue(name, values[i])
		}
	})
}
