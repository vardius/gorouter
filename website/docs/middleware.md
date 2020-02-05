---
id: middleware
title: Middleware
sidebar_label: Middleware
---

Passing middleware as follow `A, B, C` will result in `A(B(C( handler )))` where handler is your handler method.

## Global Middleware
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

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
    params, _ := context.Parameters(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    // apply middleware to all routes
    // can pass as many as you want
    router := gorouter.New(logger, example)

    router.GET("/", http.HandlerFunc(Index))
    router.GET("/hello/{name}", http.HandlerFunc(Hello))

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
## Method Middleware
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

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
    params, _ := context.Parameters(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    router := gorouter.New()

    router.GET("/", http.HandlerFunc(Index))
    router.GET("/hello/{name}", http.HandlerFunc(Hello))

    // apply middleware to all routes with GET method
    // can pass as many as you want
    router.USE("GET", "", logger, example)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
## Route Middleware
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

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
    params, _ := context.Parameters(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    router := gorouter.New()
    router.GET("/", http.HandlerFunc(Index))
    router.GET("/hello/{name}", http.HandlerFunc(Hello))

    // apply middleware to route and all it children
    // can pass as many as you want
    router.USE("GET", "/hello/{name}", logger, example)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```