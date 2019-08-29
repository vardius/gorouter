package gorouter

import (
	"bytes"

	"github.com/valyala/fasthttp"
	"github.com/vardius/gorouter/v4/middleware"
	"github.com/vardius/gorouter/v4/mux"
)

// NewFastHTTPRouter creates new Router instance, returns pointer
func NewFastHTTPRouter(fs ...FastHTTPMiddlewareFunc) FastHTTPRouter {
	m := transformFastHTTPMiddlewareFunc(fs...)

	return &fastHTTPRouter{
		routes:     mux.NewTree(),
		middleware: m,
	}
}

type fastHTTPRouter struct {
	routes     *mux.Tree
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

func (r *fastHTTPRouter) Handle(method, p string, h fasthttp.RequestHandler) {
	route := newRoute(h)
	route.PrependMiddleware(r.middleware)

	node := addNode(r.routes, method, p)
	node.SetRoute(route)
}

func (r *fastHTTPRouter) Mount(path string, h fasthttp.RequestHandler) {
	for _, method := range []string{GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD, CONNECT, TRACE} {
		route := newRoute(h)
		route.PrependMiddleware(r.middleware)

		node := addNode(r.routes, method, path)
		node.SetRoute(route)
		node.TurnIntoSubrouter()
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
	path := string(ctx.Path())
	method := string(ctx.Method())

	root := r.routes.GetByID(method)
	if root != nil {
		node, params, subPath := root.GetChildByPath(path)

		if node != nil && node.Route() != nil {
			if h := node.Route().Handler().(fasthttp.RequestHandler); h != nil {

				for _, param := range params {
					ctx.SetUserValue(param.Key, param.Value)
				}

				if subPath != "" {
					ctx.URI().SetPathBytes(fasthttp.NewPathPrefixStripper(len(subPath))(ctx))
				}

				h(ctx)
				return
			}
		}
	}

	// Handle OPTIONS
	if bytes.Equal(ctx.Method(), []byte(OPTIONS)) {
		if allow := allowed(r.routes, method, path); len(allow) > 0 {
			ctx.Response.Header.Set("Allow", allow)
			return
		}
	} else if bytes.Equal(ctx.Method(), []byte(GET)) && r.fileServer != nil {
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
