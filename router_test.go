package gorouter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/vardius/gorouter/v5"
)

func serve(h http.Handler, method, target string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(method, target, nil))
	return w
}

// echo writes the matched pattern followed by the named params.
func echo(names ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.Pattern)
		for _, name := range names {
			fmt.Fprintf(w, " %s=%s", name, r.PathValue(name))
		}
	}
}

func text(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, body) }
}

func tag(name string) gorouter.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, name+">")
			next.ServeHTTP(w, r)
		})
	}
}

func TestRouting(t *testing.T) {
	t.Parallel()

	r := gorouter.New()
	r.GET("/", echo())
	r.GET("/users", echo())
	r.GET("/users/new", echo())
	r.GET("/users/{id:[0-9]+}", echo("id"))
	r.GET("/users/{name}", echo("name"))
	r.GET("/users/{id:[0-9]+}/posts/{post}", echo("id", "post"))
	r.POST("/users", echo())
	r.Handle("PROPFIND", "/dav/{file}", echo("file"))

	tests := []struct {
		method, target string
		code           int
		body           string
	}{
		{"GET", "/", 200, "/"},
		{"GET", "/users", 200, "/users"},
		{"GET", "/users/", 200, "/users"},
		{"POST", "/users", 200, "/users"},
		{"GET", "/users/new", 200, "/users/new"},
		{"GET", "/users/42", 200, "/users/{id:[0-9]+} id=42"},
		{"GET", "/users/bob", 200, "/users/{name} name=bob"},
		{"GET", "/users/42/posts/hello", 200, "/users/{id:[0-9]+}/posts/{post} id=42 post=hello"},
		{"GET", "/users/caf%C3%A9", 200, "/users/{name} name=café"},
		{"PROPFIND", "/dav/notes.txt", 200, "/dav/{file} file=notes.txt"},
		{"GET", "/users/42/posts", 404, "404 page not found\n"},
		{"GET", "/nope", 404, "404 page not found\n"},
	}
	for _, tt := range tests {
		w := serve(r, tt.method, tt.target)
		if w.Code != tt.code || w.Body.String() != tt.body {
			t.Errorf("%s %s = %d %q, want %d %q", tt.method, tt.target, w.Code, w.Body.String(), tt.code, tt.body)
		}
	}
}

func TestMethodHelpers(t *testing.T) {
	t.Parallel()

	r := gorouter.New()
	helpers := map[string]func(string, http.HandlerFunc){
		"GET": r.GET, "HEAD": r.HEAD, "POST": r.POST, "PUT": r.PUT, "PATCH": r.PATCH,
		"DELETE": r.DELETE, "CONNECT": r.CONNECT, "OPTIONS": r.OPTIONS, "TRACE": r.TRACE,
	}
	for method, register := range helpers {
		register("/x", text(method))
	}

	for method := range helpers {
		if got := serve(r, method, "/x").Body.String(); got != method {
			t.Errorf("%s /x = %q, want the %s handler", method, got, method)
		}
	}
}

func TestMethodNotAllowed(t *testing.T) {
	t.Parallel()

	r := gorouter.New()
	r.GET("/x", text("get"))
	r.DELETE("/x", text("delete"))
	r.POST("/", text("post"))

	tests := []struct {
		method, target, allow string
	}{
		{"PUT", "/x", "GET, DELETE, OPTIONS"},
		{"GET", "/", "POST, OPTIONS"},
	}
	for _, tt := range tests {
		w := serve(r, tt.method, tt.target)
		if w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != tt.allow {
			t.Errorf("%s %s = %d Allow %q, want 405 Allow %q", tt.method, tt.target, w.Code, w.Header().Get("Allow"), tt.allow)
		}
	}

	custom := gorouter.New(gorouter.WithMethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprint(w, "allowed: "+w.Header().Get("Allow"))
	})))
	custom.GET("/x", text("get"))

	if w := serve(custom, "POST", "/x"); w.Code != http.StatusTeapot || w.Body.String() != "allowed: GET, OPTIONS" {
		t.Errorf("custom 405 = %d %q", w.Code, w.Body.String())
	}
}

func TestOptions(t *testing.T) {
	t.Parallel()

	r := gorouter.New()
	r.GET("/x", text("get"))
	r.POST("/x", text("post"))
	r.PUT("/y", text("put"))
	r.OPTIONS("/y", text("custom options"))

	if w := serve(r, "OPTIONS", "/x"); w.Code != 200 || w.Header().Get("Allow") != "GET, POST, OPTIONS" || w.Body.Len() != 0 {
		t.Errorf("OPTIONS /x = %d Allow %q body %q", w.Code, w.Header().Get("Allow"), w.Body.String())
	}
	if w := serve(r, "OPTIONS", "/y"); w.Body.String() != "custom options" {
		t.Errorf("OPTIONS /y = %q, want the registered handler", w.Body.String())
	}

	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.URL.Path = "*"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Header().Get("Allow") != "GET, POST, PUT, OPTIONS" {
		t.Errorf("OPTIONS * Allow = %q", w.Header().Get("Allow"))
	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	r := gorouter.New(gorouter.WithNotFound(text("custom 404")))
	r.GET("/x", text("x"))

	if w := serve(r, "GET", "/y"); w.Body.String() != "custom 404" {
		t.Errorf("GET /y = %q, want the custom 404", w.Body.String())
	}
}

func TestMiddlewareOrder(t *testing.T) {
	t.Parallel()

	r := gorouter.New(gorouter.WithMiddleware(tag("global1"), tag("global2")))
	r.Use("/", tag("all"))
	r.GET("/x/y", text("handler"))
	r.POST("/x/y", text("handler"))
	r.GET("/xy", text("handler"))
	r.UseMethod(http.MethodGet, "/x", tag("get-x"))
	r.Use("/x/y", tag("first"), tag("second"))

	tests := []struct {
		method, target, body string
	}{
		{"GET", "/x/y", "global1>global2>all>get-x>first>second>handler"},
		{"POST", "/x/y", "global1>global2>all>first>second>handler"},
		{"GET", "/xy", "global1>global2>all>handler"},
		// Global middleware also wraps 404 and 405 responses; route middleware does not.
		{"GET", "/nope", "global1>global2>404 page not found\n"},
		{"PUT", "/xy", "global1>global2>Method Not Allowed\n"},
	}
	for _, tt := range tests {
		if got := serve(r, tt.method, tt.target).Body.String(); got != tt.body {
			t.Errorf("%s %s = %q, want %q", tt.method, tt.target, got, tt.body)
		}
	}
}

func TestMiddlewareComposedOnce(t *testing.T) {
	t.Parallel()

	built := 0
	r := gorouter.New()
	r.Use("/", func(next http.Handler) http.Handler {
		built++
		return next
	})
	r.GET("/x", text("x"))

	before := built
	for range 3 {
		serve(r, "GET", "/x")
	}
	if built != before {
		t.Errorf("middleware was built %d more times while serving, want 0", built-before)
	}
}

func TestMount(t *testing.T) {
	t.Parallel()

	sub := gorouter.New(gorouter.WithMiddleware(tag("sub-global")))
	sub.GET("/", echo("tenant"))
	sub.GET("/users/{id}", echo("tenant", "id"))
	sub.Use("/users", tag("sub-users"))

	r := gorouter.New(gorouter.WithMiddleware(tag("global")))
	r.Mount("/t/{tenant}", sub)
	r.Use("/t/{tenant}", tag("tenant"))

	// The tags write to the body, so every response here is a 200: the bodies tell them apart.
	tests := []struct {
		method, target, body string
	}{
		{"GET", "/t/acme", "global>tenant>sub-global>/ tenant=acme"},
		{"GET", "/t/acme/users/7", "global>tenant>sub-global>sub-users>/users/{id} tenant=acme id=7"},
		{"DELETE", "/t/acme/users/7", "global>tenant>sub-global>Method Not Allowed\n"},
		{"GET", "/t/acme/nope", "global>tenant>sub-global>404 page not found\n"},
		{"GET", "/t", "global>404 page not found\n"},
	}
	for _, tt := range tests {
		if got := serve(r, tt.method, tt.target).Body.String(); got != tt.body {
			t.Errorf("%s %s = %q, want %q", tt.method, tt.target, got, tt.body)
		}
	}
}

func TestMountKeepsCallerRequest(t *testing.T) {
	t.Parallel()

	var seen string
	r := gorouter.New(gorouter.WithMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req)
			seen = req.URL.Path
		})
	}))
	r.Mount("/v1", gorouter.New())

	serve(r, "GET", "/v1/users")
	if seen != "/v1/users" {
		t.Errorf("caller's path after Mount = %q, want it unchanged", seen)
	}
}

func TestMountCustomMethod(t *testing.T) {
	t.Parallel()

	sub := gorouter.New()
	sub.Handle("PROPFIND", "/{file}", echo("file"))

	r := gorouter.New()
	r.Mount("/dav", sub)

	if w := serve(r, "PROPFIND", "/dav/notes.txt"); w.Code != http.StatusOK || w.Body.String() != "/{file} file=notes.txt" {
		t.Errorf("PROPFIND /dav/notes.txt = %d %q, want the mounted route", w.Code, w.Body.String())
	}
}

func TestMountKeepsEscapedPath(t *testing.T) {
	t.Parallel()

	var got string
	r := gorouter.New()
	r.Mount("/files", http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) { got = req.URL.EscapedPath() }))

	serve(r, "GET", "/files/a%2Fb")
	if got != "/a%2Fb" {
		t.Errorf("mounted handler's escaped path = %q, want /a%%2Fb", got)
	}
}

func TestMountPutsParamsBack(t *testing.T) {
	t.Parallel()

	sub := gorouter.New()
	sub.GET("/{id}", echo("id"))

	var after string
	r := gorouter.New()
	r.Mount("/t/{id}", sub)
	r.Use("/t/{id}", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req)
			after = req.PathValue("id")
		})
	})

	if got := serve(r, "GET", "/t/acme/42").Body.String(); got != "/{id} id=42" {
		t.Errorf("GET /t/acme/42 = %q, want the mounted route's id", got)
	}
	if after != "acme" {
		t.Errorf("mount's id once the mounted router returned = %q, want acme", after)
	}
}

func TestOptionsAsteriskSkipsMounts(t *testing.T) {
	t.Parallel()

	sub := gorouter.New()
	sub.GET("/", text("sub"))

	r := gorouter.New()
	r.POST("/login", text(""))
	r.Mount("/{tenant}", sub)

	// The mount takes every method, so the server as a whole allows them all.
	want := "GET, HEAD, POST, PUT, PATCH, DELETE, CONNECT, TRACE, OPTIONS"
	if w := serve(r, "OPTIONS", "*"); w.Header().Get("Allow") != want {
		t.Errorf("OPTIONS * Allow = %q, want %q", w.Header().Get("Allow"), want)
	}
}

func TestMountFileServer(t *testing.T) {
	t.Parallel()

	files := fstest.MapFS{"hello.txt": {Data: []byte("hi")}}
	r := gorouter.New()
	r.Mount("/static", http.FileServerFS(files))

	if w := serve(r, "GET", "/static/hello.txt"); w.Code != 200 || w.Body.String() != "hi" {
		t.Errorf("GET /static/hello.txt = %d %q", w.Code, w.Body.String())
	}
}

func TestManyParams(t *testing.T) {
	t.Parallel()

	// More params than the router keeps on the stack.
	var pattern, target strings.Builder
	var names, want []string
	for i := range 20 {
		names = append(names, fmt.Sprintf("p%d", i))
		want = append(want, fmt.Sprint(i))
		pattern.WriteString("/{" + names[i] + "}")
		target.WriteString("/" + want[i])
	}

	r := gorouter.New()
	r.GET(pattern.String(), func(w http.ResponseWriter, req *http.Request) {
		var got []string
		for _, name := range names {
			got = append(got, req.PathValue(name))
		}
		fmt.Fprint(w, strings.Join(got, " "))
	})

	if got := serve(r, "GET", target.String()).Body.String(); got != strings.Join(want, " ") {
		t.Errorf("20 params = %q, want %q", got, strings.Join(want, " "))
	}
}

func TestRoutes(t *testing.T) {
	t.Parallel()

	r := gorouter.New()
	r.GET("/a", text(""))
	r.POST("/b/{id}", text(""))
	r.Mount("/v1", gorouter.New())

	want := []gorouter.Route{{Method: "GET", Pattern: "/a"}, {Method: "POST", Pattern: "/b/{id}"}, {Method: "", Pattern: "/v1"}}
	if got := r.Routes(); !slices.Equal(got, want) {
		t.Errorf("Routes() = %v, want %v", got, want)
	}
}

func TestRegistrationPanics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		register func(*gorouter.Router)
		want     string
	}{
		{"nil handler func", func(r *gorouter.Router) { r.GET("/", nil) }, "nil handler for GET /"},
		{"nil handler", func(r *gorouter.Router) { r.Handle("GET", "/", nil) }, "nil handler for GET /"},
		{"nil mount", func(r *gorouter.Router) { r.Mount("/v1", nil) }, "nil handler for mount at /v1"},
		{"invalid pattern", func(r *gorouter.Router) { r.GET("/{id", text("")) }, "not closed"},
		{"duplicate", func(r *gorouter.Router) {
			r.GET("/a", text(""))
			r.GET("/a/", text(""))
		}, "GET /a/ conflicts with GET /a"},
		{"middleware below a mount", func(r *gorouter.Router) {
			r.Mount("/v1", gorouter.New())
			r.Use("/v1/admin", tag("admin"))
		}, "below the mount at /v1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				msg, _ := recover().(string)
				if !strings.Contains(msg, tt.want) {
					t.Errorf("panic = %q, want it to contain %q", msg, tt.want)
				}
			}()
			tt.register(gorouter.New())
		})
	}
}

func TestConcurrentRequests(t *testing.T) {
	t.Parallel()

	r := gorouter.New()
	r.Use("/", tag("1"), tag("2"))
	r.Use("/a", tag("a"))
	r.Use("/b", tag("b"))
	r.GET("/a/{id}", echo("id"))
	r.GET("/b/{id}", echo("id"))

	var wg sync.WaitGroup
	for g := range 8 {
		wg.Go(func() {
			for i := range 200 {
				target, want := fmt.Sprintf("/a/%d", i), fmt.Sprintf("1>2>a>/a/{id} id=%d", i)
				if (g+i)%2 == 0 {
					target, want = fmt.Sprintf("/b/%d", i), fmt.Sprintf("1>2>b>/b/{id} id=%d", i)
				}
				if got := serve(r, "GET", target).Body.String(); got != want {
					t.Errorf("GET %s = %q, want %q", target, got, want)
					return
				}
			}
		})
	}
	wg.Wait()
}
