package gorouter

import (
	"net/http"

	"github.com/valyala/fasthttp"
	"github.com/vardius/gorouter/v4/middleware"
	"github.com/vardius/gorouter/v4/mux"
	pathutils "github.com/vardius/gorouter/v4/path"
)

// NewFastHTTPRouter creates new Router instance, returns pointer
func NewFastHTTPRouter(fs ...FastHTTPMiddlewareFunc) FastHTTPRouter {
	return &fastHTTPRouter{
		routes:     mux.NewTree(),
		middleware: transformFastHTTPMiddlewareFunc(fs...),
	}
}

type fastHTTPRouter struct {
	routes     mux.Tree
	middleware middleware.Middleware
	fileServer fasthttp.RequestHandler
	notFound   fasthttp.RequestHandler
	notAllowed fasthttp.RequestHandler
}

func (r *fastHTTPRouter) PrettyPrint() string {
	return r.routes.PrettyPrint()
}

func (r *fastHTTPRouter) POST(p string, f fasthttp.RequestHandler) {
	r.Handle(http.MethodPost, p, f)
}

func (r *fastHTTPRouter) GET(p string, f fasthttp.RequestHandler) {
	r.Handle(http.MethodGet, p, f)
}

func (r *fastHTTPRouter) PUT(p string, f fasthttp.RequestHandler) {
	r.Handle(http.MethodPut, p, f)
}

func (r *fastHTTPRouter) DELETE(p string, f fasthttp.RequestHandler) {
	r.Handle(http.MethodDelete, p, f)
}

func (r *fastHTTPRouter) PATCH(p string, f fasthttp.RequestHandler) {
	r.Handle(http.MethodPatch, p, f)
}

func (r *fastHTTPRouter) OPTIONS(p string, f fasthttp.RequestHandler) {
	r.Handle(http.MethodOptions, p, f)
}

func (r *fastHTTPRouter) HEAD(p string, f fasthttp.RequestHandler) {
	r.Handle(http.MethodHead, p, f)
}

func (r *fastHTTPRouter) CONNECT(p string, f fasthttp.RequestHandler) {
	r.Handle(http.MethodConnect, p, f)
}

func (r *fastHTTPRouter) TRACE(p string, f fasthttp.RequestHandler) {
	r.Handle(http.MethodTrace, p, f)
}

func (r *fastHTTPRouter) USE(method, path string, fs ...FastHTTPMiddlewareFunc) {
	m := transformFastHTTPMiddlewareFunc(fs...)

	r.routes = r.routes.WithMiddleware(method+path, m, 0)
}

func (r *fastHTTPRouter) Handle(method, path string, h fasthttp.RequestHandler) {
	route := newRoute(h)

	r.routes = r.routes.WithRoute(method+path, route, 0)
}

func (r *fastHTTPRouter) Mount(path string, h fasthttp.RequestHandler) {
	for _, method := range []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	} {
		route := newRoute(h)

		r.routes = r.routes.WithSubrouter(method+path, route, 0)
	}
}

func (r *fastHTTPRouter) Compile() {
	for i, methodNode := range r.routes {
		r.routes[i].WithChildren(methodNode.Tree().Compile())
	}
}

func (r *fastHTTPRouter) NotFound(notFound fasthttp.RequestHandler) {
	r.notFound = notFound
}

func (r *fastHTTPRouter) NotAllowed(notAllowed fasthttp.RequestHandler) {
	r.notAllowed = notAllowed
}

func (r *fastHTTPRouter) ServeFiles(root string, stripSlashes int) {
	if root == "" {
		panic("gorouter.ServeFiles: empty root!")
	}

	r.fileServer = fasthttp.FSHandler(root, stripSlashes)
}

func (r *fastHTTPRouter) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	method := string(ctx.Method())
	pathAsString := string(ctx.Path())
	path := pathutils.TrimSlash(pathAsString)

	if root := r.routes.Find(method); root != nil {
		if node, treeMiddleware, params, subPath := root.Tree().Match(path); node != nil && node.Route() != nil {
			route := node.Route()
			handler := route.Handler()
			middleware := r.middleware.Merge(treeMiddleware)
			computedHandler := middleware.Compose(handler)

			h := computedHandler.(fasthttp.RequestHandler)

			if len(params) > 0 {
				ctx.SetUserValue("params", params)
			}

			if subPath != "" {
				ctx.URI().SetPathBytes(fasthttp.NewPathPrefixStripper(len("/" + subPath))(ctx))
			}

			h(ctx)
			return
		}

		if pathAsString == "/" && root.Route() != nil {
			root.Route().Handler().(fasthttp.RequestHandler)(ctx)
			return
		}
	}

	// Handle OPTIONS
	if method == http.MethodOptions {
		if allow := allowed(r.routes, method, path); len(allow) > 0 {
			ctx.Response.Header.Set("Allow", allow)
			return
		}
	} else if method == http.MethodGet && r.fileServer != nil {
		// Handle file serve
		r.fileServer(ctx)
		return
	} else {
		// Handle 405
		if allow := allowed(r.routes, method, path); len(allow) > 0 {
			ctx.Response.Header.Set("Allow", allow)
			r.serveNotAllowed(ctx)
			return
		}
	}

	// Handle 404
	r.serveNotFound(ctx)
}

func (r *fastHTTPRouter) serveNotFound(ctx *fasthttp.RequestCtx) {
	if r.notFound != nil {
		r.notFound(ctx)
	} else {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
	}
}

func (r *fastHTTPRouter) serveNotAllowed(ctx *fasthttp.RequestCtx) {
	if r.notAllowed != nil {
		r.notAllowed(ctx)
	} else {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed), fasthttp.StatusMethodNotAllowed)
	}
}

func transformFastHTTPMiddlewareFunc(fs ...FastHTTPMiddlewareFunc) middleware.Middleware {
	m := make(middleware.Middleware, len(fs))

	for i, f := range fs {
		m[i] = func(mf FastHTTPMiddlewareFunc) middleware.MiddlewareFunc {
			return func(h interface{}) interface{} {
				return mf(h.(fasthttp.RequestHandler))
			}
		}(f) // f is a reference to function so we have to wrap if with that callback
	}

	return m
}
