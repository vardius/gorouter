package gorouter

import (
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

func (r *fastHTTPRouter) POST(p string, f fasthttp.RequestHandler) {
	r.Handle(POST, p, f)
}

func (r *fastHTTPRouter) GET(p string, f fasthttp.RequestHandler) {
	r.Handle(GET, p, f)
}

func (r *fastHTTPRouter) PUT(p string, f fasthttp.RequestHandler) {
	r.Handle(PUT, p, f)
}

func (r *fastHTTPRouter) DELETE(p string, f fasthttp.RequestHandler) {
	r.Handle(DELETE, p, f)
}

func (r *fastHTTPRouter) PATCH(p string, f fasthttp.RequestHandler) {
	r.Handle(PATCH, p, f)
}

func (r *fastHTTPRouter) OPTIONS(p string, f fasthttp.RequestHandler) {
	r.Handle(OPTIONS, p, f)
}

func (r *fastHTTPRouter) HEAD(p string, f fasthttp.RequestHandler) {
	r.Handle(HEAD, p, f)
}

func (r *fastHTTPRouter) CONNECT(p string, f fasthttp.RequestHandler) {
	r.Handle(CONNECT, p, f)
}

func (r *fastHTTPRouter) TRACE(p string, f fasthttp.RequestHandler) {
	r.Handle(TRACE, p, f)
}

func (r *fastHTTPRouter) USE(method, p string, fs ...FastHTTPMiddlewareFunc) {
	m := transformFastHTTPMiddlewareFunc(fs...)

	addMiddleware(r.routes, method, p, m)
}

func (r *fastHTTPRouter) Handle(method, path string, h fasthttp.RequestHandler) {
	route := newRoute(h)
	route.PrependMiddleware(r.middleware)

	r.routes = r.routes.WithRoute(method+path, route, 0)
}

func (r *fastHTTPRouter) Mount(path string, h fasthttp.RequestHandler) {
	for _, method := range []string{TRACE, CONNECT, HEAD, OPTIONS, PATCH, DELETE, PUT, POST, GET} {
		route := newRoute(h)
		route.PrependMiddleware(r.middleware)

		r.routes = r.routes.WithSubrouter(method+path, route, 0)
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
		if node, params, subPath := root.Tree().Match(path); node != nil && node.Route() != nil {
			h := node.Route().Handler().(fasthttp.RequestHandler)

			for _, param := range params {
				ctx.SetUserValue(param.Key, param.Value)
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
	if method == OPTIONS {
		if allow := allowed(r.routes, method, path); len(allow) > 0 {
			ctx.Response.Header.Set("Allow", allow)
			return
		}
	} else if method == GET && r.fileServer != nil {
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
