package gorouter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing/fstest"

	"github.com/vardius/gorouter/v5"
)

// get sends a request to h and prints the response, so examples can show what a client would see.
func get(h http.Handler, method, target string) {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(method, target, nil))
	fmt.Printf("%s %s -> %d %s\n", method, target, w.Code, strings.TrimSpace(w.Body.String()))
}

func Example() {
	r := gorouter.New()
	r.GET("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "welcome")
	})
	r.GET("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, %s", r.PathValue("name"))
	})

	get(r, "GET", "/")
	get(r, "GET", "/hello/gopher")
	// Output:
	// GET / -> 200 welcome
	// GET /hello/gopher -> 200 hello, gopher
}

// Literals win over regexp params, which win over plain params. A branch that fails further down
// falls back to the next one.
func Example_priority() {
	r := gorouter.New()
	show := func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, r.Pattern) }
	r.GET("/users/new", show)
	r.GET("/users/{id:[0-9]+}", show)
	r.GET("/users/{name}", show)
	r.GET("/users/{name}/posts", show)

	get(r, "GET", "/users/new")
	get(r, "GET", "/users/42")
	get(r, "GET", "/users/gopher")
	get(r, "GET", "/users/new/posts")
	// Output:
	// GET /users/new -> 200 /users/new
	// GET /users/42 -> 200 /users/{id:[0-9]+}
	// GET /users/gopher -> 200 /users/{name}
	// GET /users/new/posts -> 200 /users/{name}/posts
}

func ExampleRouter_Use() {
	logger := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("log:", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}

	r := gorouter.New()
	r.Use("/admin", logger)
	r.GET("/", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "home") })
	r.GET("/admin/users", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "users") })

	get(r, "GET", "/")
	get(r, "GET", "/admin/users")
	// Output:
	// GET / -> 200 home
	// log: GET /admin/users
	// GET /admin/users -> 200 users
}

func ExampleRouter_UseMethod() {
	readOnly := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "read-only", http.StatusForbidden)
		})
	}

	r := gorouter.New()
	r.UseMethod("DELETE", "/", readOnly)
	r.GET("/files/{name}", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, r.PathValue("name")) })
	r.DELETE("/files/{name}", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "deleted") })

	get(r, "GET", "/files/notes.txt")
	get(r, "DELETE", "/files/notes.txt")
	// Output:
	// GET /files/notes.txt -> 200 notes.txt
	// DELETE /files/notes.txt -> 403 read-only
}

func ExampleRouter_Mount() {
	v1 := gorouter.New()
	v1.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "tenant %s, user %s, path %s", r.PathValue("tenant"), r.PathValue("id"), r.URL.Path)
	})

	r := gorouter.New()
	r.Mount("/{tenant}/v1", v1)

	get(r, "GET", "/acme/v1/users/7")
	// Output:
	// GET /acme/v1/users/7 -> 200 tenant acme, user 7, path /users/7
}

// Mount a file server to serve files below a path.
func ExampleRouter_Mount_files() {
	files := fstest.MapFS{"hello.txt": {Data: []byte("hello from a file")}}

	r := gorouter.New()
	r.Mount("/static", http.FileServerFS(files))

	get(r, "GET", "/static/hello.txt")
	// Output:
	// GET /static/hello.txt -> 200 hello from a file
}

func ExampleRouter_Routes() {
	r := gorouter.New()
	r.GET("/users", func(http.ResponseWriter, *http.Request) {})
	r.POST("/users", func(http.ResponseWriter, *http.Request) {})
	r.GET("/users/{id}", func(http.ResponseWriter, *http.Request) {})

	for _, route := range r.Routes() {
		fmt.Println(route.Method, route.Pattern)
	}
	// Output:
	// GET /users
	// POST /users
	// GET /users/{id}
}

// Middleware passed to New wraps the whole router, so it also sees 404 and 405 responses.
func ExampleWithMiddleware() {
	recoverer := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, "something went wrong", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}

	r := gorouter.New(gorouter.WithMiddleware(recoverer))
	r.GET("/panic", func(http.ResponseWriter, *http.Request) { panic("boom") })

	get(r, "GET", "/panic")
	// Output:
	// GET /panic -> 500 something went wrong
}

func ExampleWithNotFound() {
	r := gorouter.New(gorouter.WithNotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "no page at "+r.URL.Path, http.StatusNotFound)
	})))

	get(r, "GET", "/missing")
	// Output:
	// GET /missing -> 404 no page at /missing
}

// A 405 response carries an Allow header listing the methods the path does accept.
func ExampleWithMethodNotAllowed() {
	r := gorouter.New(gorouter.WithMethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "use "+w.Header().Get("Allow"), http.StatusMethodNotAllowed)
	})))
	r.GET("/users", func(http.ResponseWriter, *http.Request) {})
	r.POST("/users", func(http.ResponseWriter, *http.Request) {})

	get(r, "DELETE", "/users")
	// Output:
	// DELETE /users -> 405 use GET, POST, OPTIONS
}

// Handle takes any http.Handler and any method, including custom ones.
func ExampleRouter_Handle() {
	r := gorouter.New()
	r.Handle("PURGE", "/cache/{key}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "purged ", r.PathValue("key"))
	}))

	get(r, "PURGE", "/cache/home")
	// Output:
	// PURGE /cache/home -> 200 purged home
}
