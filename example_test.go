package gorouter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/vardius/gorouter"
)

func logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %q\n", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	params, _ := gorouter.FromContext(r.Context())
	fmt.Printf("Hello, %s!\n", params.Value("name"))
}

func Example() {
	router := gorouter.New()
	router.GET("/hello/{name}", http.HandlerFunc(Hello))

	// Normally you would call
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
	router := gorouter.New(logger)
	router.GET("/hello/{name}", http.HandlerFunc(Hello))

	// Normally you would call
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
	router := gorouter.New()
	router.GET("/hello/{name}", http.HandlerFunc(Hello))

	router.USE("GET", "/hello/{name}", logger)

	// Normally you would call
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
