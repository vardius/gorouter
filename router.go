package gorouter

import (
	"net/http"
	"strings"

	"github.com/vardius/gorouter/v4/context"
	"github.com/vardius/gorouter/v4/middleware"
	"github.com/vardius/gorouter/v4/mux"
	path_utils "github.com/vardius/gorouter/v4/path"
)

// HTTP methods constants
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	HEAD    = "HEAD"
	CONNECT = "CONNECT"
	TRACE   = "TRACE"
)

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
	USE(method, pattern string, fs ...middleware.MiddlewareFunc)

	// Handle adds http.Handler as router handler
	// under given method and patter
	Handle(method, pattern string, handler http.Handler)

	// HandleFunc adds http.HandlerFunc as router handler
	// under given method and patter
	HandleFunc(method, pattern string, handler http.HandlerFunc)

	// Mount another handler as a subrouter
	Mount(pattern string, handler http.Handler)

	// ServeHTTP dispatches the request to the route handler
	// whose pattern matches the request URL
	ServeHTTP(http.ResponseWriter, *http.Request)

	// ServeFile replies to the request with the
	// contents of the named file or directory.
	ServeFiles(root http.FileSystem, path string, strip bool)

	// NotFound replies to the request with the
	// 404 Error code
	NotFound(http.Handler)

	// NotFound replies to the request with the
	// 405 Error code
	NotAllowed(http.Handler)
}

type router struct {
	routes     *mux.Tree
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

func (r *router) USE(method, p string, fs ...middleware.MiddlewareFunc) {
	r.addMiddleware(method, p, fs...)
}

func (r *router) Handle(m, p string, h http.Handler) {
	r.addRoute(m, p, h)
}

func (r *router) HandleFunc(m, p string, f http.HandlerFunc) {
	r.addRoute(m, p, f)
}

func (r *router) Mount(p string, h http.Handler) {
	r.addRoute(GET, p, h).TurnIntoSubrouter()
	r.addRoute(POST, p, h).TurnIntoSubrouter()
	r.addRoute(PUT, p, h).TurnIntoSubrouter()
	r.addRoute(DELETE, p, h).TurnIntoSubrouter()
	r.addRoute(PATCH, p, h).TurnIntoSubrouter()
	r.addRoute(OPTIONS, p, h).TurnIntoSubrouter()
	r.addRoute(HEAD, p, h).TurnIntoSubrouter()
	r.addRoute(CONNECT, p, h).TurnIntoSubrouter()
	r.addRoute(TRACE, p, h).TurnIntoSubrouter()
}

func (r *router) NotFound(notFound http.Handler) {
	r.notFound = notFound
}

func (r *router) NotAllowed(notAllowed http.Handler) {
	r.notAllowed = notAllowed
}

func (r *router) ServeFiles(root http.FileSystem, path string, strip bool) {
	if path == "" {
		panic("gorouter.ServeFiles: empty path!")
	}
	handler := http.FileServer(root)
	if strip {
		handler = http.StripPrefix("/"+path+"/", handler)
	}
	r.fileServer = handler
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	root := r.routes.GetByID(req.Method)
	if root != nil {
		node, params, subPath := root.GetChildByPath(req.URL.Path)

		if node != nil && node.Route() != nil {
			if h := node.Route().(*route).getHandler(); h != nil {
				req = req.WithContext(context.WithParams(req.Context(), params))

				if subPath != "" {
					h = http.StripPrefix(strings.TrimSuffix(req.URL.Path, subPath), h)
				}

				h.ServeHTTP(w, req)
				return
			}
		}
	}

	// Handle OPTIONS
	if req.Method == OPTIONS {
		if allow := r.allowed(req.Method, req.URL.Path); len(allow) > 0 {
			w.Header().Set("Allow", allow)
			return
		}
	} else if req.Method == GET && r.fileServer != nil {
		// Handle file serve
		r.fileServer.ServeHTTP(w, req)
		return
	} else {
		// Handle 405
		if allow := r.allowed(req.Method, req.URL.Path); len(allow) > 0 {
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

func (r *router) addRoute(method, path string, f http.Handler) *mux.Node {
	root := r.routes.GetByID(method)
	if root == nil {
		root = mux.NewRoot(method)
		r.routes.Insert(root)
	}

	paths := path_utils.Split(path)
	route := newRoute(f)
	route.prependMiddleware(r.middleware)

	n := root.AddChild(paths)
	n.SetRoute(route)

	return n
}

func (r *router) addMiddleware(method, path string, fs ...middleware.MiddlewareFunc) {
	type recFunc func(recFunc, *mux.Node, middleware.Middleware)
	c := func(c recFunc, n *mux.Node, m middleware.Middleware) {
		if n.Route() != nil {
			n.Route().(*route).appendMiddleware(m)
		}
		for _, child := range n.Children().StaticNodes() {
			c(c, child, m)
		}
		for _, child := range n.Children().RegexpNodes() {
			c(c, child, m)
		}
		if n.Children().WildcardNode() != nil {
			c(c, n.Children().WildcardNode(), m)
		}
	}

	paths := path_utils.Split(path)

	// routes tree roots should be http method nodes only
	for _, root := range r.routes.StaticNodes() {
		if method == "" || method == root.ID() {
			node, _ := root.GetChild(paths)
			if node != nil {
				c(c, node, middleware.New(fs...))
			}
		}
	}
}

func (r *router) allowed(method, path string) (allow string) {
	if path == "*" {
		// routes tree roots should be http method nodes only
		for _, root := range r.routes.StaticNodes() {
			if root.ID() == OPTIONS {
				continue
			}
			if len(allow) == 0 {
				allow = root.ID()
			} else {
				allow += ", " + root.ID()
			}
		}
	} else {
		// routes tree roots should be http method nodes only
		for _, root := range r.routes.StaticNodes() {
			if root.ID() == method || root.ID() == OPTIONS {
				continue
			}

			n, _, _ := root.GetChildByPath(path)
			if n != nil && n.Route() != nil {
				if len(allow) == 0 {
					allow = root.ID()
				} else {
					allow += ", " + root.ID()
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return allow
}

// New creates new Router instance, return pointer
func New(fs ...middleware.MiddlewareFunc) Router {
	return &router{
		routes:     mux.NewTree(),
		middleware: middleware.New(fs...),
	}
}
