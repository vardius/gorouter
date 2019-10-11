package gorouter

import (
	"net/http"
	"strings"

	"github.com/vardius/gorouter/v4/context"
	"github.com/vardius/gorouter/v4/middleware"
	"github.com/vardius/gorouter/v4/mux"
	pathutils "github.com/vardius/gorouter/v4/path"
)

// New creates new net/http Router instance, returns pointer
func New(fs ...MiddlewareFunc) Router {
	return &router{
		routes:     mux.NewTree(),
		middleware: transformMiddlewareFunc(fs...),
	}
}

type router struct {
	routes     mux.Tree
	middleware middleware.Middleware
	fileServer http.Handler
	notFound   http.Handler
	notAllowed http.Handler
}

func (r *router) POST(p string, f http.Handler) {
	r.Handle(POST, p, f)
}

func (r *router) GET(p string, f http.Handler) {
	r.Handle(GET, p, f)
}

func (r *router) PUT(p string, f http.Handler) {
	r.Handle(PUT, p, f)
}

func (r *router) DELETE(p string, f http.Handler) {
	r.Handle(DELETE, p, f)
}

func (r *router) PATCH(p string, f http.Handler) {
	r.Handle(PATCH, p, f)
}

func (r *router) OPTIONS(p string, f http.Handler) {
	r.Handle(OPTIONS, p, f)
}

func (r *router) HEAD(p string, f http.Handler) {
	r.Handle(HEAD, p, f)
}

func (r *router) CONNECT(p string, f http.Handler) {
	r.Handle(CONNECT, p, f)
}

func (r *router) TRACE(p string, f http.Handler) {
	r.Handle(TRACE, p, f)
}

func (r *router) USE(method, p string, fs ...MiddlewareFunc) {
	m := transformMiddlewareFunc(fs...)

	addMiddleware(r.routes, method, p, m)
}

func (r *router) Handle(method, path string, h http.Handler) {
	route := newRoute(h)
	route.PrependMiddleware(r.middleware)

	r.routes = r.routes.WithRoute(method+path, route, 0)
}

func (r *router) Mount(path string, h http.Handler) {
	for _, method := range []string{TRACE, CONNECT, HEAD, OPTIONS, PATCH, DELETE, PUT, POST, GET} {
		route := newRoute(h)
		route.PrependMiddleware(r.middleware)

		r.routes = r.routes.WithSubrouter(method+path, route, 0)
	}
}

func (r *router) NotFound(notFound http.Handler) {
	r.notFound = notFound
}

func (r *router) NotAllowed(notAllowed http.Handler) {
	r.notAllowed = notAllowed
}

func (r *router) ServeFiles(fs http.FileSystem, root string, strip bool) {
	if root == "" {
		panic("gorouter.ServeFiles: empty root!")
	}
	handler := http.FileServer(fs)
	if strip {
		handler = http.StripPrefix("/"+root+"/", handler)
	}
	r.fileServer = handler
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := pathutils.TrimSlash(req.URL.Path)

	if root := r.routes.Find(req.Method); root != nil {
		if node, params, subPath := root.Tree().Match(path); node != nil && node.Route() != nil {
			h := node.Route().Handler().(http.Handler)
			req = req.WithContext(context.WithParams(req.Context(), params))

			if subPath != "" {
				h = http.StripPrefix(strings.TrimSuffix(req.URL.Path, "/"+subPath), h)
			}

			h.ServeHTTP(w, req)
			return
		}

		if req.URL.Path == "/" && root.Route() != nil {
			root.Route().Handler().(http.Handler).ServeHTTP(w, req)
			return
		}
	}

	// Handle OPTIONS
	if req.Method == OPTIONS {
		if allow := allowed(r.routes, req.Method, path); len(allow) > 0 {
			w.Header().Set("Allow", allow)
			return
		}
	} else if req.Method == GET && r.fileServer != nil {
		// Handle file serve
		r.fileServer.ServeHTTP(w, req)
		return
	} else {
		// Handle 405
		if allow := allowed(r.routes, req.Method, path); len(allow) > 0 {
			w.Header().Set("Allow", allow)
			r.serveNotAllowed(w, req)
			return
		}
	}

	// Handle 404
	r.serveNotFound(w, req)
}

func (r *router) serveNotFound(w http.ResponseWriter, req *http.Request) {
	if r.notFound != nil {
		r.notFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (r *router) serveNotAllowed(w http.ResponseWriter, req *http.Request) {
	if r.notAllowed != nil {
		r.notAllowed.ServeHTTP(w, req)
	} else {
		http.Error(w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}

func transformMiddlewareFunc(fs ...MiddlewareFunc) middleware.Middleware {
	m := make(middleware.Middleware, len(fs))

	for i, f := range fs {
		m[i] = func(mf MiddlewareFunc) middleware.MiddlewareFunc {
			return func(h interface{}) interface{} {
				return mf(h.(http.Handler))
			}
		}(f) // f is a reference to function so we have to wrap if with that callback
	}

	return m
}
