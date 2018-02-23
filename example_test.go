package gorouter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/vardius/gorouter"
)

func Example() {
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := gorouter.FromContext(r.Context())
		fmt.Printf("Hello, %s!\n", params.Value("name"))
	}

	router := gorouter.New()
	router.GET("/hello/{name}", http.HandlerFunc(hello))

	// Normally you would call ListenAndServe starting an HTTP server
	// with a given address and router as a handler
	// log.Fatal(http.ListenAndServe(":8080", router))
	// but for this example we will mock request

	w := httptest.NewRecorder()
	req, err := http.NewRequest(gorouter.GET, "/hello/guest", nil)
	if err != nil {
		return
	}

	router.ServeHTTP(w, req)

	// Output:
	// Hello, guest!
}

func ExampleMiddlewareFunc_global() {
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := gorouter.FromContext(r.Context())
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

	// Normally you would call ListenAndServe starting an HTTP server
	// with a given address and router as a handler
	// log.Fatal(http.ListenAndServe(":8080", router))
	// but for this example we will mock request

	w := httptest.NewRecorder()
	req, err := http.NewRequest(gorouter.GET, "/hello/guest", nil)
	if err != nil {
		return
	}

	router.ServeHTTP(w, req)

	// Output:
	// [GET] "/hello/guest"
	// Hello, guest!
}

func ExampleMiddlewareFunc_route() {
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := gorouter.FromContext(r.Context())
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

	// Normally you would call ListenAndServe starting an HTTP server
	// with a given address and router as a handler
	// log.Fatal(http.ListenAndServe(":8080", router))
	// but for this example we will mock request

	w := httptest.NewRecorder()
	req, err := http.NewRequest(gorouter.GET, "/hello/guest", nil)
	if err != nil {
		return
	}

	router.ServeHTTP(w, req)

	// Output:
	// [GET] "/hello/guest"
	// Hello, guest!
}

func ExampleMiddlewareFunc_method() {
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := gorouter.FromContext(r.Context())
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

	// Normally you would call ListenAndServe starting an HTTP server
	// with a given address and router as a handler
	// log.Fatal(http.ListenAndServe(":8080", router))
	// but for this example we will mock request

	w := httptest.NewRecorder()
	req, err := http.NewRequest(gorouter.GET, "/hello/guest", nil)
	if err != nil {
		return
	}

	router.ServeHTTP(w, req)

	// Output:
	// [GET] "/hello/guest"
	// Hello, guest!
}

func ExampleRouter_mount() {
	hello := func(w http.ResponseWriter, r *http.Request) {
		params, _ := gorouter.FromContext(r.Context())
		fmt.Printf("Hello, %s!\n", params.Value("name"))
	}

	subrouter := gorouter.New()
	subrouter.GET("/{name}", http.HandlerFunc(hello))

	router := gorouter.New()
	router.Mount("/hello", subrouter)

	// Normally you would call ListenAndServe starting an HTTP server
	// with a given address and router as a handler
	// log.Fatal(http.ListenAndServe(":8080", router))
	// but for this example we will mock request

	w := httptest.NewRecorder()
	req, err := http.NewRequest(gorouter.GET, "/hello/guest", nil)
	if err != nil {
		return
	}

	router.ServeHTTP(w, req)

	// Output:
	// Hello, guest!
}
