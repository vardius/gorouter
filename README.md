# 🍃 gorouter

[![Test](https://github.com/vardius/gorouter/actions/workflows/test.yml/badge.svg)](https://github.com/vardius/gorouter/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vardius/gorouter/v5.svg)](https://pkg.go.dev/github.com/vardius/gorouter/v5)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/gorouter/v5)](https://goreportcard.com/report/github.com/vardius/gorouter/v5)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md)

A fast HTTP router for Go, for both `net/http` and `fasthttp`.

- **Zero-allocation routing.** Matching a route, running its middleware and finding no route allocate nothing; on fasthttp, neither do params or mounts.
- **Standard handlers.** Plain `http.Handler`s and `func(http.Handler) http.Handler` middleware; read params with the standard `r.PathValue`.
- **Params and regexps.** `/users/{id}`, `/users/{id:[0-9]+}`, with literals winning over params and backtracking when a branch fails further down.
- **Scoped middleware.** Around the whole router, the routes below a path, or one method below a path. Composed once at registration, never per request.
- **Mounting.** Hand a path to a sub-router or a file server, with the prefix stripped.
- **Correct by default.** 405 with an `Allow` header, automatic `OPTIONS`, custom methods.

📖 **[Documentation](https://vardius.github.io/gorouter/)** · [API reference](https://pkg.go.dev/github.com/vardius/gorouter/v5) · [Examples](examples) · [Migrating from v4](MIGRATION.md)

## Install

```sh
go get github.com/vardius/gorouter/v5
```

Requires Go 1.26 or newer.

## Quick start

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vardius/gorouter/v5"
)

func main() {
	r := gorouter.New()
	r.GET("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "welcome")
	})
	r.GET("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, %s!\n", r.PathValue("name"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
```

The same with fasthttp, from the `fastrouter` package, so `net/http` users never link fasthttp:

```go
r := fastrouter.New()
r.GET("/hello/{name}", func(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", fastrouter.Param(ctx, "name"))
})

log.Fatal(fasthttp.ListenAndServe(":8080", r.HandleFastHTTP))
```

## Routing

```go
r.GET("/users/new", newUser)              // literal
r.GET("/users/{id:[0-9]+}", userByID)     // regexp param: the whole segment must match
r.GET("/users/{name}", userByName)        // param: any one segment
r.GET("/users/{name}/posts", postsByUser) // /users/new/posts lands here
r.Handle("PURGE", "/cache/{key}", purge)  // any method, any http.Handler
```

Literals win over regexp params, regexp params are tried in registration order, and a plain param comes last.
If the preferred branch fails further down the path, the router backtracks to the next one.
A trailing slash is ignored and an empty segment never matches.

## Middleware

```go
r := gorouter.New(gorouter.WithMiddleware(logRequests, recoverPanics)) // whole router, 404 and 405 included
r.Use("/admin", basicAuth)                                            // every route below /admin
r.UseMethod(http.MethodPost, "/api", requireJSON)                     // POST routes below /api
```

Middleware runs in registration order, outermost first, and covers routes registered before and after the call.
`Use` picks routes by pattern, not by request path: a request runs the middleware of the route it matches,
so `/admin/users` matched by `/{section}/users` doesn't run `basicAuth`.

## Mounting

```go
r.Mount("/api/v1", v1Router)                          // v1Router sees /users/7 for /api/v1/users/7
r.Mount("/orgs/{org}", orgRouter)                     // r.PathValue("org") still works inside
r.Mount("/static", http.FileServerFS(os.DirFS("public")))
```

A mount takes every method, custom ones included. Middleware on the mount path or above wraps it;
middleware for paths below it belongs in the mounted router, and `Use` panics if asked to put it there.

## Benchmarks

Time and allocations to route one request, on an Apple M4 Max with Go 1.27.0, median of 5 runs.
Both versions ran on the same harness, with no-op handlers and middleware.
Run v5's own suite with `go test -bench . ./...`.

| Scenario | net/http v4 | net/http v5 | fasthttp v4 | fasthttp v5 |
| --- | ---: | ---: | ---: | ---: |
| Static route | 20.0 ns · 0 allocs | 19.9 ns · 0 allocs | 30.0 ns · 1 alloc | **16.9 ns · 0 allocs** |
| Static, 5 segments | 28.5 ns · 0 allocs | 30.2 ns · 0 allocs | 39.0 ns · 1 alloc | **26.2 ns · 0 allocs** |
| One param | 87.0 ns · 4 allocs | 93.6 ns · **2 allocs** | 50.1 ns · 3 allocs | **38.5 ns · 0 allocs** |
| Three params + regexp | 150.1 ns · 4 allocs | 148.5 ns · **2 allocs** | 110.2 ns · 3 allocs | **87.4 ns · 0 allocs** |
| Static + 3 middleware | 170.4 ns · 6 allocs | **23.8 ns · 0 allocs** | 173.1 ns · 7 allocs | **18.3 ns · 0 allocs** |
| One param + 3 middleware | 238.4 ns · 10 allocs | **100.8 ns · 2 allocs** | 192.8 ns · 9 allocs | **39.5 ns · 0 allocs** |
| 404 | 20.0 ns · 0 allocs | 20.3 ns · 0 allocs | 29.7 ns · 1 alloc | **15.3 ns · 0 allocs** |
| 405 | 60.8 ns · 3 allocs | **55.8 ns · 2 allocs** | 89.0 ns · 3 allocs | **69.5 ns · 1 alloc** |

On `net/http`, the two allocations for params are the ones `Request.SetPathValue` makes to store them for `r.PathValue`;
matching itself allocates nothing. Storing them that way also makes a one-param route a few nanoseconds slower than v4,
the price of reading params with the standard `r.PathValue`. A 405 allocates its `Allow` header value, and a mount one copy of the request, as `http.StripPrefix` does.

Across a 70-route, GitHub-style API (`BenchmarkAPI`), an average request takes 95.5 ns on `net/http` and 53.1 ns with zero allocations on fasthttp.

## Upgrading from v4

v5 is a new major version with a smaller API. The short version:

```go
// v4
router := gorouter.New(logger)
router.GET("/hello/{name}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params, _ := context.Parameters(r.Context())
	fmt.Fprintf(w, "hello, %s!", params.Value("name"))
}))
router.USE(http.MethodGet, "/admin", auth)

// v5
router := gorouter.New(gorouter.WithMiddleware(logger))
router.GET("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %s!", r.PathValue("name"))
})
router.UseMethod(http.MethodGet, "/admin", auth)
```

[MIGRATION.md](MIGRATION.md) maps every v4 call to its v5 form and lists the behaviour changes.

## Contributing

Bug reports and pull requests are welcome: see [CONTRIBUTING.md](.github/CONTRIBUTING.md).
For a ready-made API layout built on gorouter, see [go-api-boilerplate](https://github.com/vardius/go-api-boilerplate).

## License

MIT, see [LICENSE.md](LICENSE.md).
