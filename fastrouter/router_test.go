package fastrouter_test

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	"github.com/vardius/gorouter/v5/fastrouter"
)

func serve(r *fastrouter.Router, method, uri string) *fasthttp.RequestCtx {
	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(method)
	ctx.Request.SetRequestURI(uri)
	r.HandleFastHTTP(ctx)
	return ctx
}

// echo writes label followed by the named params.
func echo(label string, names ...string) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Fprint(ctx, label)
		for _, name := range names {
			fmt.Fprintf(ctx, " %s=%s", name, fastrouter.Param(ctx, name))
		}
	}
}

func tag(name string) fastrouter.Middleware {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			fmt.Fprint(ctx, name+">")
			next(ctx)
		}
	}
}

func TestRouting(t *testing.T) {
	t.Parallel()

	r := fastrouter.New()
	r.GET("/", echo("root"))
	r.GET("/users/new", echo("new"))
	r.GET("/users/{id:[0-9]+}", echo("by id", "id"))
	r.GET("/users/{name}", echo("by name", "name"))
	r.GET("/users/{id:[0-9]+}/posts/{post}", echo("post", "id", "post"))
	r.Handle("PROPFIND", "/dav/{file}", echo("dav", "file"))

	tests := []struct {
		method, uri string
		code        int
		body        string
	}{
		{"GET", "/", 200, "root"},
		{"GET", "/users/new", 200, "new"},
		{"GET", "/users/42", 200, "by id id=42"},
		{"GET", "/users/bob/", 200, "by name name=bob"},
		{"GET", "/users/42/posts/hello", 200, "post id=42 post=hello"},
		{"PROPFIND", "/dav/a.txt", 200, "dav file=a.txt"},
		{"GET", "/nope", 404, "Not Found"},
	}
	for _, tt := range tests {
		ctx := serve(r, tt.method, tt.uri)
		if ctx.Response.StatusCode() != tt.code || string(ctx.Response.Body()) != tt.body {
			t.Errorf("%s %s = %d %q, want %d %q",
				tt.method, tt.uri, ctx.Response.StatusCode(), ctx.Response.Body(), tt.code, tt.body)
		}
	}
}

func TestMethodNotAllowedAndOptions(t *testing.T) {
	t.Parallel()

	r := fastrouter.New(fastrouter.WithMethodNotAllowed(func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusTeapot)
		fmt.Fprint(ctx, "allowed: ", string(ctx.Response.Header.Peek("Allow")))
	}))
	r.GET("/x", echo("get"))
	r.POST("/x", echo("post"))

	if ctx := serve(r, "PUT", "/x"); ctx.Response.StatusCode() != fasthttp.StatusTeapot ||
		string(ctx.Response.Body()) != "allowed: GET, POST, OPTIONS" {
		t.Errorf("PUT /x = %d %q", ctx.Response.StatusCode(), ctx.Response.Body())
	}
	if ctx := serve(r, "OPTIONS", "/x"); ctx.Response.StatusCode() != 200 ||
		string(ctx.Response.Header.Peek("Allow")) != "GET, POST, OPTIONS" {
		t.Errorf("OPTIONS /x = %d Allow %q", ctx.Response.StatusCode(), ctx.Response.Header.Peek("Allow"))
	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	r := fastrouter.New(fastrouter.WithNotFound(echo("custom 404")))

	if ctx := serve(r, "GET", "/y"); string(ctx.Response.Body()) != "custom 404" {
		t.Errorf("GET /y = %q, want the custom 404", ctx.Response.Body())
	}
}

func TestMiddlewareOrder(t *testing.T) {
	t.Parallel()

	r := fastrouter.New(fastrouter.WithMiddleware(tag("global")))
	r.Use("/", tag("all"))
	r.GET("/x/y", echo("handler"))
	r.POST("/x/y", echo("handler"))
	r.UseMethod(fasthttp.MethodGet, "/x", tag("get-x"))
	r.Use("/x/y", tag("first"), tag("second"))

	tests := []struct {
		method, uri, body string
	}{
		{"GET", "/x/y", "global>all>get-x>first>second>handler"},
		{"POST", "/x/y", "global>all>first>second>handler"},
	}
	for _, tt := range tests {
		if got := string(serve(r, tt.method, tt.uri).Response.Body()); got != tt.body {
			t.Errorf("%s %s = %q, want %q", tt.method, tt.uri, got, tt.body)
		}
	}
}

func TestGlobalMiddlewareWrapsNotFound(t *testing.T) {
	t.Parallel()

	// ctx.Error resets the response, so this middleware marks it after the router is done.
	r := fastrouter.New(fastrouter.WithMiddleware(func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			next(ctx)
			ctx.Response.Header.Set("X-Global", "seen")
		}
	}))

	if ctx := serve(r, "GET", "/nope"); string(ctx.Response.Header.Peek("X-Global")) != "seen" {
		t.Error("global middleware did not wrap the 404")
	}
}

func TestMount(t *testing.T) {
	t.Parallel()

	sub := fastrouter.New()
	sub.GET("/", echo("sub root", "tenant"))
	sub.GET("/users/{id}", echo("sub user", "tenant", "id"))
	sub.GET("/path", func(ctx *fasthttp.RequestCtx) { fmt.Fprint(ctx, string(ctx.Path())) })

	r := fastrouter.New()
	r.Mount("/t/{tenant}", sub.HandleFastHTTP)
	r.Use("/t/{tenant}", tag("tenant"))

	tests := []struct {
		uri, body string
	}{
		{"/t/acme", "tenant>sub root tenant=acme"},
		// The tenant value has to survive the path rewrite that reuses its buffer.
		{"/t/acme/users/7", "tenant>sub user tenant=acme id=7"},
		{"/t/acme/path", "tenant>/path"},
	}
	for _, tt := range tests {
		if got := string(serve(r, "GET", tt.uri).Response.Body()); got != tt.body {
			t.Errorf("GET %s = %q, want %q", tt.uri, got, tt.body)
		}
	}
}

func TestParamsReleasedWithRequest(t *testing.T) {
	t.Parallel()

	r := fastrouter.New()
	r.GET("/users/{id}", echo("user", "id"))

	ctx := serve(r, "GET", "/users/7")
	if fastrouter.Param(ctx, "id") != "7" {
		t.Fatalf("Param during the request = %q, want 7", fastrouter.Param(ctx, "id"))
	}

	// fasthttp resets user values once a request is done, which hands the params back to the pool.
	ctx.ResetUserValues()
	if got := fastrouter.Param(ctx, "id"); got != "" {
		t.Errorf("Param after the request = %q, want empty", got)
	}
}

func TestServer(t *testing.T) {
	t.Parallel()

	r := fastrouter.New()
	r.GET("/users/{id}", echo("user", "id"))
	r.POST("/users/{id}", echo("update", "id"))

	ln := fasthttputil.NewInmemoryListener()
	t.Cleanup(func() { _ = ln.Close() })
	go func() { _ = fasthttp.Serve(ln, r.HandleFastHTTP) }()

	conn, err := ln.Dial()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	br := bufio.NewReader(conn)

	// Pooled params are reused across these keep-alive requests.
	for _, tt := range []struct{ request, want string }{
		{"GET /users/1 HTTP/1.1\r\nHost: x\r\n\r\n", "user id=1"},
		{"GET /users/22 HTTP/1.1\r\nHost: x\r\n\r\n", "user id=22"},
		{"OPTIONS * HTTP/1.1\r\nHost: x\r\n\r\n", "Allow: GET, POST, OPTIONS"},
	} {
		if got := roundTrip(t, conn, br, tt.request); !strings.Contains(got, tt.want) {
			t.Errorf("%q: response %q, want it to contain %q", tt.request, got, tt.want)
		}
	}
}

func roundTrip(t *testing.T, conn net.Conn, br *bufio.Reader, request string) string {
	t.Helper()

	if _, err := conn.Write([]byte(request)); err != nil {
		t.Fatal(err)
	}

	var resp fasthttp.Response
	if err := resp.Read(br); err != nil {
		t.Fatal(err)
	}

	return resp.Header.String() + string(resp.Body())
}

func TestConcurrentRequests(t *testing.T) {
	t.Parallel()

	r := fastrouter.New()
	r.Use("/a", tag("a"))
	r.Use("/b", tag("b"))
	r.GET("/a/{id}", echo("A", "id"))
	r.GET("/b/{id}", echo("B", "id"))

	var wg sync.WaitGroup
	for g := range 8 {
		wg.Go(func() {
			for i := range 200 {
				uri, want := fmt.Sprintf("/a/%d", i), fmt.Sprintf("a>A id=%d", i)
				if (g+i)%2 == 0 {
					uri, want = fmt.Sprintf("/b/%d", i), fmt.Sprintf("b>B id=%d", i)
				}
				ctx := serve(r, "GET", uri)
				if got := string(ctx.Response.Body()); got != want {
					t.Errorf("GET %s = %q, want %q", uri, got, want)
					return
				}
				ctx.ResetUserValues()
			}
		})
	}
	wg.Wait()
}

func TestRegistrationPanics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		register func(*fastrouter.Router)
		want     string
	}{
		{"nil handler", func(r *fastrouter.Router) { r.GET("/", nil) }, "nil handler for GET /"},
		{"nil mount", func(r *fastrouter.Router) { r.Mount("/v1", nil) }, "nil handler for mount at /v1"},
		{"invalid pattern", func(r *fastrouter.Router) { r.GET("x", echo("")) }, "must start with /"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				msg, _ := recover().(string)
				if !strings.Contains(msg, tt.want) {
					t.Errorf("panic = %q, want it to contain %q", msg, tt.want)
				}
			}()
			tt.register(fastrouter.New())
		})
	}
}
