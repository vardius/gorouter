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
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := context.Parameters(r.Context())
		fmt.Printf("Hello, %s!\n", params.Value("name"))
	}

	router := gorouter.New()
	router.GET("/hello/{name}", http.HandlerFunc(hello))

	// for this example we will mock request
	handleNetHTTPRequest("GET", "/hello/guest", router)

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

func ExampleRouter_mount() {
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := context.Parameters(r.Context())
		fmt.Printf("Hello, %s!\n", params.Value("name"))
	}

	// gorouter as subrouter
	subrouter := gorouter.New()
	subrouter.GET("/{name}", http.HandlerFunc(hello))

	// default mux as subrouter
	// you can use eveything that implements http.Handler interface
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
