package fastrouter_test

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/vardius/gorouter/v5/fastrouter"
)

// get sends a request to r and prints the response, so examples can show what a client would see.
func get(r *fastrouter.Router, method, uri string) {
	ctx := serve(r, method, uri)
	fmt.Println(strings.TrimSpace(fmt.Sprintf("%s %s -> %d %s", method, uri, ctx.Response.StatusCode(), ctx.Response.Body())))
}

func Example() {
	r := fastrouter.New()
	r.GET("/", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprint(ctx, "welcome")
	})
	r.GET("/hello/{name}", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "hello, %s", fastrouter.Param(ctx, "name"))
	})

	get(r, "GET", "/")
	get(r, "GET", "/hello/gopher")
	// Output:
	// GET / -> 200 welcome
	// GET /hello/gopher -> 200 hello, gopher
}

// A param value points into the request's buffer. Clone it to keep it after the handler returns.
func ExampleParam() {
	var kept []string

	r := fastrouter.New()
	r.POST("/tags/{tag}", func(ctx *fasthttp.RequestCtx) {
		kept = append(kept, strings.Clone(fastrouter.Param(ctx, "tag")))
	})

	get(r, "POST", "/tags/go")
	get(r, "POST", "/tags/http")
	fmt.Println(kept)
	// Output:
	// POST /tags/go -> 200
	// POST /tags/http -> 200
	// [go http]
}

func ExampleRouter_Use() {
	logger := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			fmt.Printf("log: %s %s\n", ctx.Method(), ctx.Path())
			next(ctx)
		}
	}

	r := fastrouter.New()
	r.Use("/admin", logger)
	r.GET("/", func(ctx *fasthttp.RequestCtx) { fmt.Fprint(ctx, "home") })
	r.GET("/admin/users", func(ctx *fasthttp.RequestCtx) { fmt.Fprint(ctx, "users") })

	get(r, "GET", "/")
	get(r, "GET", "/admin/users")
	// Output:
	// GET / -> 200 home
	// log: GET /admin/users
	// GET /admin/users -> 200 users
}

func ExampleRouter_Mount() {
	v1 := fastrouter.New()
	v1.GET("/users/{id}", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "tenant %s, user %s, path %s", fastrouter.Param(ctx, "tenant"), fastrouter.Param(ctx, "id"), ctx.Path())
	})

	r := fastrouter.New()
	r.Mount("/{tenant}/v1", v1.HandleFastHTTP)

	get(r, "GET", "/acme/v1/users/7")
	// Output:
	// GET /acme/v1/users/7 -> 200 tenant acme, user 7, path /users/7
}

// A 405 response carries an Allow header listing the methods the path does accept.
func ExampleWithMethodNotAllowed() {
	r := fastrouter.New(fastrouter.WithMethodNotAllowed(func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		fmt.Fprint(ctx, "use ", string(ctx.Response.Header.Peek(fasthttp.HeaderAllow)))
	}))
	r.GET("/users", func(*fasthttp.RequestCtx) {})
	r.POST("/users", func(*fasthttp.RequestCtx) {})

	get(r, "DELETE", "/users")
	// Output:
	// DELETE /users -> 405 use GET, POST, OPTIONS
}
