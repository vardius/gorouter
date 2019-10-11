package gorouter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/valyala/fasthttp"
	"github.com/vardius/gorouter/v4"
	"github.com/vardius/gorouter/v4/context"
)

func handleNetHTTPRequest(method, path string, handler http.Handler) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, nil)

	handler.ServeHTTP(w, req)
}

func handleFastHTTPRequest(method, path string, handler fasthttp.RequestHandler) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.URI().SetPath(path)

	handler(ctx)
}

func Example() {
	index := func(_ http.ResponseWriter, _ *http.Request) {
		fmt.Printf("Hello\n")
	}

	hello := func(_ http.ResponseWriter, r *http.Request) {
		params, _ := context.Parameters(r.Context())
		fmt.Printf("Hello, %s!\n", params.Value("name"))
	}

	router := gorouter.New()
	router.GET("/", http.HandlerFunc(index))
	router.GET("/hello/{name}", http.HandlerFunc(hello))

	// for this example we will mock request
	handleNetHTTPRequest("GET", "/", router)
	handleNetHTTPRequest("GET", "/hello/guest", router)

	// Output:
	// Hello
	// Hello, guest!
}

func Example_second() {
	hello := func(ctx *fasthttp.RequestCtx) {
		fmt.Printf("Hello, %s!\n", ctx.UserValue("name"))
	}

	router := gorouter.NewFastHTTPRouter()
	router.GET("/hello/{name}", hello)

	// for this example we will mock request
	handleFastHTTPRequest("GET", "/hello/guest", router.HandleFastHTTP)

	// Output:
	// Hello, guest!
}

func ExampleMiddlewareFunc() {
	// Global middleware example
	// applies to all routes
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := context.Parameters(r.Context())
		fmt.Printf("Hello, %s!\n", params.Value("name"))
	}

	logger := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[%s] %q\n", r.Method, r.URL.String())
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	// apply middlewares to all routes
	// can pass as many as you want
	router := gorouter.New(logger)
	router.GET("/hello/{name}", http.HandlerFunc(hello))

	// for this example we will mock request
	handleNetHTTPRequest("GET", "/hello/guest", router)

	// Output:
	// [GET] "/hello/guest"
	// Hello, guest!
}

func ExampleMiddlewareFunc_second() {
	// Route level middleware example
	// applies to route and its lower tree
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := context.Parameters(r.Context())
		fmt.Printf("Hello, %s!\n", params.Value("name"))
	}

	logger := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[%s] %q\n", r.Method, r.URL.String())
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	router := gorouter.New()
	router.GET("/hello/{name}", http.HandlerFunc(hello))

	// apply middlewares to route and all it children
	// can pass as many as you want
	router.USE("GET", "/hello/{name}", logger)

	// for this example we will mock request
	handleNetHTTPRequest("GET", "/hello/guest", router)

	// Output:
	// [GET] "/hello/guest"
	// Hello, guest!
}

func ExampleMiddlewareFunc_third() {
	// Http method middleware example
	// applies to all routes under this method
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := context.Parameters(r.Context())
		fmt.Printf("Hello, %s!\n", params.Value("name"))
	}

	logger := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("[%s] %q\n", r.Method, r.URL.String())
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	router := gorouter.New()
	router.GET("/hello/{name}", http.HandlerFunc(hello))

	// apply middlewares to all routes with GET method
	// can pass as many as you want
	router.USE("GET", "", logger)

	// for this example we will mock request
	handleNetHTTPRequest("GET", "/hello/guest", router)

	// Output:
	// [GET] "/hello/guest"
	// Hello, guest!
}

func ExampleFastHTTPMiddlewareFunc() {
	// Global middleware example
	// applies to all routes
	hello := func(ctx *fasthttp.RequestCtx) {
		fmt.Printf("Hello, %s!\n", ctx.UserValue("name"))
	}

	logger := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		fn := func(ctx *fasthttp.RequestCtx) {
			fmt.Printf("[%s] %q\n", ctx.Method(), ctx.Path())
			next(ctx)
		}

		return fn
	}

	router := gorouter.NewFastHTTPRouter(logger)
	router.GET("/hello/{name}", hello)

	// for this example we will mock request
	handleFastHTTPRequest("GET", "/hello/guest", router.HandleFastHTTP)

	// Output:
	// [GET] "/hello/guest"
	// Hello, guest!
}

func ExampleFastHTTPMiddlewareFunc_second() {
	// Route level middleware example
	// applies to route and its lower tree
	hello := func(ctx *fasthttp.RequestCtx) {
		fmt.Printf("Hello, %s!\n", ctx.UserValue("name"))
	}

	logger := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		fn := func(ctx *fasthttp.RequestCtx) {
			fmt.Printf("[%s] %q\n", ctx.Method(), ctx.Path())
			next(ctx)
		}

		return fn
	}

	router := gorouter.NewFastHTTPRouter()
	router.GET("/hello/{name}", hello)

	// apply middlewares to route and all it children
	// can pass as many as you want
	router.USE("GET", "/hello/{name}", logger)

	// for this example we will mock request
	handleFastHTTPRequest("GET", "/hello/guest", router.HandleFastHTTP)

	// Output:
	// [GET] "/hello/guest"
	// Hello, guest!
}

func ExampleFastHTTPMiddlewareFunc_third() {
	// Http method middleware example
	// applies to all routes under this method
	hello := func(ctx *fasthttp.RequestCtx) {
		fmt.Printf("Hello, %s!\n", ctx.UserValue("name"))
	}

	logger := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		fn := func(ctx *fasthttp.RequestCtx) {
			fmt.Printf("[%s] %q\n", ctx.Method(), ctx.Path())
			next(ctx)
		}

		return fn
	}

	router := gorouter.NewFastHTTPRouter()
	router.GET("/hello/{name}", hello)

	// apply middlewares to all routes with GET method
	// can pass as many as you want
	router.USE("GET", "", logger)

	// for this example we will mock request
	handleFastHTTPRequest("GET", "/hello/guest", router.HandleFastHTTP)

	// Output:
	// [GET] "/hello/guest"
	// Hello, guest!
}

func ExampleRouter_mount() {
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := context.Parameters(r.Context())
		fmt.Printf("Hello, %s!\n", params.Value("name"))
	}

	// gorouter as subrouter
	subrouter := gorouter.New()
	subrouter.GET("/{name}", http.HandlerFunc(hello))

	// default mux as subrouter
	// you can use everything that implements http.Handler interface
	unknownSubrouter := http.NewServeMux()
	unknownSubrouter.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Hi, guest!")
	})

	router := gorouter.New()
	router.Mount("/hello", subrouter)
	router.Mount("/hi", unknownSubrouter)

	// for this example we will mock request
	handleNetHTTPRequest("GET", "/hello/guest", router)
	handleNetHTTPRequest("GET", "/hi/guest", router)

	// Output:
	// Hello, guest!
	// Hi, guest!
}

func ExampleFastHTTPRouter_mount() {
	hello := func(ctx *fasthttp.RequestCtx) {
		fmt.Printf("Hello, %s!\n", ctx.UserValue("name"))
	}

	subrouter := gorouter.NewFastHTTPRouter()
	subrouter.GET("/{name}", hello)

	router := gorouter.NewFastHTTPRouter()
	router.Mount("/hello", subrouter.HandleFastHTTP)

	// for this example we will mock request
	handleFastHTTPRequest("GET", "/hello/guest", router.HandleFastHTTP)

	// Output:
	// Hello, guest!
}
