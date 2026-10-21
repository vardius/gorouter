# 🍃 gorouter

[![Test](https://github.com/vardius/gorouter/actions/workflows/test.yml/badge.svg)](https://github.com/vardius/gorouter/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vardius/gorouter/v5.svg)](https://pkg.go.dev/github.com/vardius/gorouter/v5)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/gorouter/v5)](https://goreportcard.com/report/github.com/vardius/gorouter/v5)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md)

A fast HTTP router for Go, for both `net/http` and `fasthttp`.

- **Zero-allocation matching.** Static routes, middleware chains and 404s allocate nothing; on fasthttp, neither do params.
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

## Mounting

```go
r.Mount("/api/v1", v1Router)                          // v1Router sees /users/7 for /api/v1/users/7
r.Mount("/orgs/{org}", orgRouter)                     // r.PathValue("org") still works inside
r.Mount("/static", http.FileServerFS(os.DirFS("public")))
```

## Benchmarks

Time and allocations to route one request, on an Apple M4 Max with Go 1.27.0, median of 5 runs.
Both versions ran on the same harness, with no-op handlers and middleware.
Run v5's own suite with `go test -bench . ./...`.

| Scenario | net/http v4 | net/http v5 | fasthttp v4 | fasthttp v5 |
| --- | ---: | ---: | ---: | ---: |
| Static route | 20.4 ns · 0 allocs | **19.2 ns · 0 allocs** | 30.4 ns · 1 alloc | **16.1 ns · 0 allocs** |
| Static, 5 segments | 28.9 ns · 0 allocs | 29.3 ns · 0 allocs | 39.1 ns · 1 alloc | **25.3 ns · 0 allocs** |
| One param | 87.5 ns · 4 allocs | 93.3 ns · **2 allocs** | 50.0 ns · 3 allocs | **34.2 ns · 0 allocs** |
| Three params + regexp | 154.2 ns · 4 allocs | **148.7 ns · 2 allocs** | 111.3 ns · 3 allocs | **83.0 ns · 0 allocs** |
| Static + 3 middleware | 158.2 ns · 6 allocs | **23.9 ns · 0 allocs** | 163.7 ns · 7 allocs | **18.1 ns · 0 allocs** |
| One param + 3 middleware | 233.6 ns · 10 allocs | **99.4 ns · 2 allocs** | 196.6 ns · 9 allocs | **35.9 ns · 0 allocs** |
| 404 | 21.2 ns · 0 allocs | **19.2 ns · 0 allocs** | 33.0 ns · 1 alloc | **14.8 ns · 0 allocs** |
| 405 | 60.6 ns · 3 allocs | **53.3 ns · 2 allocs** | 88.3 ns · 3 allocs | **65.8 ns · 1 alloc** |

On `net/http`, the two allocations for params are the ones `Request.SetPathValue` makes to store them for `r.PathValue`;
matching itself allocates nothing. A 405 allocates its `Allow` header value.

Across a 70-route, GitHub-style API (`BenchmarkAPI`), an average request takes 88.8 ns on `net/http` and 51.2 ns with zero allocations on fasthttp.

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
