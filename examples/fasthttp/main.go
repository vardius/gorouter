// Command fasthttp is the hello example on fasthttp, with a middleware and a mounted sub-router.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/vardius/gorouter/v5/fastrouter"
)

func newRouter() *fastrouter.Router {
	api := fastrouter.New()
	api.GET("/users/{id:[0-9]+}", user)

	r := fastrouter.New(fastrouter.WithMiddleware(serverHeader))
	r.GET("/", index)
	r.GET("/hello/{name}", hello)
	r.Mount("/api/v1", api.HandleFastHTTP)

	return r
}

func index(ctx *fasthttp.RequestCtx) {
	fmt.Fprintln(ctx, "welcome")
}

func hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", fastrouter.Param(ctx, "name"))
}

func user(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "user #%s\n", fastrouter.Param(ctx, "id"))
}

// serverHeader runs for every request, including 404 and 405 responses.
func serverHeader(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("X-Powered-By", "gorouter")
		next(ctx)
	}
}

func main() {
	srv := &fasthttp.Server{Handler: newRouter().HandleFastHTTP, ReadTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe(":8080"))
}
