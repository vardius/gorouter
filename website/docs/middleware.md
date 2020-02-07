---
id: middleware
title: Middleware
sidebar_label: Middleware
---

Passing middleware as follow `A, B, C` will result in `A(B(C( handler )))` where handler is your handler method.

## Global Middleware

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```go
package main

import (
	"fmt"
	"log"
	"time"
	"net/http"

	"github.com/vardius/gorouter/v4"
	"github.com/vardius/gorouter/v4/context"
)

func logger(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    t1 := time.Now()
    next.ServeHTTP(w, r)
    t2 := time.Now()
    log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
  }

  return http.HandlerFunc(fn)
}

func example(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    // do smth
    next.ServeHTTP(w, r)
  }

  return http.HandlerFunc(fn)
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request) {
    params, _ := context.Parameters(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    // apply middleware to all routes
    // can pass as many as you want
    router := gorouter.New(logger, example)

    router.GET("/", http.HandlerFunc(index))
    router.GET("/hello/{name}", http.HandlerFunc(hello))

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
<!--valyala/fasthttp-->
```go
package main

import (
    "fmt"
	"log"
	"time"

    "github.com/valyala/fasthttp"
    "github.com/vardius/gorouter/v4"
)

func logger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
  fn := func(ctx *fasthttp.RequestCtx) {
    t1 := time.Now()
    next(ctx)
    t2 := time.Now()
    log.Printf("[%s] %q %v\n", ctx.Method(), ctx.Path(), t2.Sub(t1))
  }

  return fn
}

func example(next fasthttp.RequestHandler) fasthttp.RequestHandler {
  fn := func(ctx *fasthttp.RequestCtx) {
    // do smth
    next(ctx)
  }

  return fn
}

func index(_ *fasthttp.RequestCtx) {
    fmt.Print("Welcome!\n")
}

func hello(ctx *fasthttp.RequestCtx) {
    params := ctx.UserValue("params").(context.Params)
    fmt.Printf("Hello, %s!\n", params.Value("name"))
}

func main() {
    // apply middleware to all routes
    // can pass as many as you want
    router := gorouter.NewFastHTTPRouter(logger, example)

    router.GET("/", index)
    router.GET("/hello/{name}", hello)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```
<!--END_DOCUSAURUS_CODE_TABS-->

## Method Middleware

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```go
package main

import (
	"fmt"
	"log"
	"time"
	"net/http"

	"github.com/vardius/gorouter/v4"
	"github.com/vardius/gorouter/v4/context"
)

func logger(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    t1 := time.Now()
    next.ServeHTTP(w, r)
    t2 := time.Now()
    log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
  }

  return http.HandlerFunc(fn)
}

func example(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    // do smth
    next.ServeHTTP(w, r)
  }

  return http.HandlerFunc(fn)
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request) {
    params, _ := context.Parameters(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    router := gorouter.New()

    router.GET("/", http.HandlerFunc(index))
    router.GET("/hello/{name}", http.HandlerFunc(hello))

    // apply middleware to all routes with GET method
    // can pass as many as you want
    router.USE("GET", "", logger, example)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
<!--valyala/fasthttp-->
```go
package main

import (
    "fmt"
	"log"
	"time"

    "github.com/valyala/fasthttp"
    "github.com/vardius/gorouter/v4"
)

func logger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
  fn := func(ctx *fasthttp.RequestCtx) {
    t1 := time.Now()
    next(ctx)
    t2 := time.Now()
    log.Printf("[%s] %q %v\n", ctx.Method(), ctx.Path(), t2.Sub(t1))
  }

  return fn
}

func example(next fasthttp.RequestHandler) fasthttp.RequestHandler {
  fn := func(ctx *fasthttp.RequestCtx) {
    // do smth
    next(ctx)
  }

  return fn
}

func index(_ *fasthttp.RequestCtx) {
    fmt.Print("Welcome!\n")
}

func hello(ctx *fasthttp.RequestCtx) {
    params := ctx.UserValue("params").(context.Params)
    fmt.Printf("Hello, %s!\n", params.Value("name"))
}

func main() {
    router := gorouter.NewFastHTTPRouter()

    router.GET("/", index)
    router.GET("/hello/{name}", hello)

    // apply middleware to all routes with GET method
    // can pass as many as you want
    router.USE("GET", "", logger, example)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```
<!--END_DOCUSAURUS_CODE_TABS-->

## Route Middleware

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```go
package main

import (
	"fmt"
	"log"
	"time"
	"net/http"

	"github.com/vardius/gorouter/v4"
	"github.com/vardius/gorouter/v4/context"
)

func logger(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    t1 := time.Now()
    next.ServeHTTP(w, r)
    t2 := time.Now()
    log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
  }

  return http.HandlerFunc(fn)
}

func example(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    // do smth
    next.ServeHTTP(w, r)
  }

  return http.HandlerFunc(fn)
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request) {
    params, _ := context.Parameters(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    router := gorouter.New()
    router.GET("/", http.HandlerFunc(index))
    router.GET("/hello/{name}", http.HandlerFunc(hello))

    // apply middleware to route and all it children
    // can pass as many as you want
    router.USE("GET", "/hello/{name}", logger, example)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
<!--valyala/fasthttp-->
```go
package main

import (
    "fmt"
	"log"
	"time"

    "github.com/valyala/fasthttp"
    "github.com/vardius/gorouter/v4"
)

func logger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
  fn := func(ctx *fasthttp.RequestCtx) {
    t1 := time.Now()
    next(ctx)
    t2 := time.Now()
    log.Printf("[%s] %q %v\n", ctx.Method(), ctx.Path(), t2.Sub(t1))
  }

  return fn
}

func example(next fasthttp.RequestHandler) fasthttp.RequestHandler {
  fn := func(ctx *fasthttp.RequestCtx) {
    // do smth
    next(ctx)
  }

  return fn
}

func index(_ *fasthttp.RequestCtx) {
    fmt.Print("Welcome!\n")
}

func hello(ctx *fasthttp.RequestCtx) {
    params := ctx.UserValue("params").(context.Params)
    fmt.Printf("Hello, %s!\n", params.Value("name"))
}

func main() {
    router := gorouter.NewFastHTTPRouter()

    router.GET("/", index)
    router.GET("/hello/{name}", hello)

    // apply middleware to route and all it children
    // can pass as many as you want
    router.USE("GET", "/hello/{name}", logger, example)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```
<!--END_DOCUSAURUS_CODE_TABS-->
