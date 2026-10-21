# Migrating from v4 to v5

v5 keeps the pattern syntax and the idea of the router, and replaces the API around it:
concrete types instead of interfaces, functional options instead of setters, the standard `r.PathValue`
instead of a custom params type, and a separate `fastrouter` package for fasthttp.

## 1. Update the import path

```sh
go get github.com/vardius/gorouter/v5
```

```go
import "github.com/vardius/gorouter/v5"            // net/http
import "github.com/vardius/gorouter/v5/fastrouter" // fasthttp
```

The `context`, `middleware`, `mux` and `path` packages are gone. v5 needs Go 1.26 or newer.

## 2. Before and after

### net/http

```go
// v4
import (
	"github.com/vardius/gorouter/v4"
	"github.com/vardius/gorouter/v4/context"
)

func hello(w http.ResponseWriter, r *http.Request) {
	params, _ := context.Parameters(r.Context())
	fmt.Fprintf(w, "hello, %s!", params.Value("name"))
}

func main() {
	router := gorouter.New(logger)
	router.NotFound(http.HandlerFunc(notFound))
	router.GET("/hello/{name}", http.HandlerFunc(hello))
	router.USE(http.MethodGet, "/admin", auth)
	router.USEANY("/api", cors)
	router.Mount("/v1", v1)
	router.ServeFiles(http.Dir("static"), "static", true)

	log.Fatal(http.ListenAndServe(":8080", router))
}
```

```go
// v5
import "github.com/vardius/gorouter/v5"

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %s!", r.PathValue("name"))
}

func main() {
	router := gorouter.New(
		gorouter.WithMiddleware(logger),
		gorouter.WithNotFound(http.HandlerFunc(notFound)),
	)
	router.GET("/hello/{name}", hello)
	router.UseMethod(http.MethodGet, "/admin", auth)
	router.Use("/api", cors)
	router.Mount("/v1", v1)
	router.Mount("/static", http.FileServerFS(os.DirFS("static")))

	log.Fatal(http.ListenAndServe(":8080", router))
}
```

### fasthttp

```go
// v4
import (
	"github.com/valyala/fasthttp"
	"github.com/vardius/gorouter/v4"
	"github.com/vardius/gorouter/v4/context"
)

func hello(ctx *fasthttp.RequestCtx) {
	params := ctx.UserValue("params").(context.Params)
	fmt.Fprintf(ctx, "hello, %s!", params.Value("name"))
}

func main() {
	router := gorouter.NewFastHTTPRouter(logger)
	router.GET("/hello/{name}", hello)
	router.ServeFiles("static", 1)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```

```go
// v5
import (
	"github.com/valyala/fasthttp"
	"github.com/vardius/gorouter/v5/fastrouter"
)

func hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!", fastrouter.Param(ctx, "name"))
}

func main() {
	router := fastrouter.New(fastrouter.WithMiddleware(logger))
	router.GET("/hello/{name}", hello)
	router.Mount("/static", fasthttp.FSHandler("static", 0))

	log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```

## 3. Every v4 call and its v5 form

The fasthttp column is the same with `fastrouter` in place of `gorouter` and `fasthttp.RequestHandler` in place of `http.Handler`.

| v4 | v5 |
| --- | --- |
| `gorouter.New(mw...)` | `gorouter.New(gorouter.WithMiddleware(mw...))` |
| `gorouter.NewFastHTTPRouter(mw...)` | `fastrouter.New(fastrouter.WithMiddleware(mw...))` |
| `gorouter.Router`, `gorouter.FastHTTPRouter` (interfaces) | `*gorouter.Router`, `*fastrouter.Router` |
| `gorouter.MiddlewareFunc` | `gorouter.Middleware`, an alias of `func(http.Handler) http.Handler` |
| `gorouter.FastHTTPMiddlewareFunc` | `fastrouter.Middleware` |
| `router.GET(p, http.HandlerFunc(f))`, and the other methods | `router.GET(p, f)`: the method helpers take a handler func |
| `router.GET(p, h)` with any other `http.Handler` | `router.Handle(http.MethodGet, p, h)` |
| `router.Handle(method, p, h)` | unchanged, and any method works, custom ones included |
| `context.Parameters(r.Context())`, then `params.Value("id")` | `r.PathValue("id")` |
| `ctx.UserValue("params").(context.Params).Value("id")` | `fastrouter.Param(ctx, "id")` |
| `router.USE(method, p, mw...)` | `router.UseMethod(method, p, mw...)` |
| `router.USE(method, "", mw...)` | `router.UseMethod(method, "/", mw...)` |
| `router.USEANY(p, mw...)` | `router.Use(p, mw...)` |
| `router.NotFound(h)` | `gorouter.WithNotFound(h)`, passed to `New` |
| `router.NotAllowed(h)` | `gorouter.WithMethodNotAllowed(h)`, passed to `New` |
| `router.Mount(p, h)` | unchanged |
| `router.ServeFiles(fs, "static", true)` | `router.Mount("/static", http.FileServer(fs))` |
| `router.ServeFiles(fs, root, false)` | `gorouter.WithNotFound(http.FileServer(fs))`: serves files for paths no route matches |
| fasthttp `router.ServeFiles(root, n)` | `router.Mount("/static", fasthttp.FSHandler(root, 0))`, or `fastrouter.WithNotFound(fasthttp.FSHandler(root, n))` |
| `router.Compile()` | removed: the tree is always compact |
| `router.PrettyPrint()` | `router.Routes()` lists method and pattern for each route |

## 4. Behaviour changes

- **`Use` picks routes by pattern, not by request path. Check this one if you protect routes with it.**
  v4 chose middleware by walking the request path; v5 attaches it to the routes whose pattern starts with the `Use` pattern,
  and a request runs the middleware of the route it matches. With `Use("/admin", auth)` and a route `/{section}/users`,
  a request for `/admin/users` ran `auth` in v4 and doesn't in v5. Register protected routes under the protected pattern,
  or use the middleware on the wildcard route too.
- **`Use` below a mount panics.** `Use("/api/v1/admin", mw)` next to `Mount("/api/v1", sub)` can never run, because `sub` routes
  those paths: use it on `sub` instead.
- **Regexps match the whole segment.** Before v4.6.0, `{id:[0-9]+}` matched `12abc`; v5 anchors every regexp as `^(?:regexp)$`.
  If you relied on a partial match, write it out, for example `{id:[0-9]+.*}`.
- **A regexp can't contain a slash.** v4 accepted `{p:a/b}` but split it at the slash, so `/x/{p:a/b}` matched paths such as `/x/zzz/b}`; v5 rejects the pattern.
- **Invalid patterns panic at registration with a message naming the problem**: an unclosed `{`, an empty param name,
  a segment that mixes text and a param, or a repeated param name.
- **Registering a route twice panics.** v4 accepted the second registration without a word.
- **`Use` scopes by segment.** Before v4.6.0, `USE("GET", "/api", mw)` also wrapped `/apiary`; v5's `Use("/api", mw)` covers `/api` and the routes below it only.
- **Middleware is composed at registration**, once per route, instead of on every request. It still runs in registration order, outermost first, and still covers routes registered later.
- **Params on net/http are set with `Request.SetPathValue`**, so `r.PathValue` works in handlers and in middleware registered with `Use`. `Request.Pattern` holds the matched pattern, which logging and metrics middleware can read.
- **fasthttp param values are valid until the handler returns.** They live in a buffer reused once the request is done. Use `strings.Clone` to keep one longer.
- **`Mount` takes every method**, custom ones included. v4 answered 405 to a custom method such as `PROPFIND` below a mount.
- **On fasthttp, `Mount` decodes the path once and puts it back afterwards.** v4 decoded the stripped path a second time, so `/api/users/a%2525b` reached the mounted router as `a%b`, and middleware around the mount saw the stripped path once the mounted handler returned.
- **fasthttp 405 responses carry the `Allow` header.** v4's default 405 handler called `ctx.Error`, which resets the headers, so the header never reached the client.
- **`Routes()` replaces `PrettyPrint()`** and lists a mount once, with an empty method, since a mount takes every method.

Unchanged: 404 and 405 responses, the `Allow` header on `net/http`, automatic `OPTIONS` responses, backtracking between literal and param routes,
trailing slashes being ignored, middleware passed to `New` wrapping 404 and 405 responses, and `Mount` stripping the mount path.

## 5. Staying on v4

v4 keeps working, and its fixes ship on the `v4` branch. Moving to v5 is only needed for the new API and the speed-up.
