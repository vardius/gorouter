Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/goserver)](https://goreportcard.com/report/github.com/vardius/goserver)
[![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver)
[![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/goserver/blob/master/LICENSE.md)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Appling Middleware
----------------
1. [Global Middlewares](#global-middlewares)
2. [Method Middlewares](#method-middlewares)
3. [Route Middlewares](#route-middlewares)

## Global Middlewares
```go
package main

import (
	"fmt"
	"log"
	"time"
	"net/http"

    "github.com/vardius/goserver"
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
	//do smth
    next.ServeHTTP(w, r)
  }

  return http.HandlerFunc(fn)
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	params, _ := goserver.FromContext(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
	//apply middlewares to all routes
	//can pass as many as you want
    server := goserver.New(logger, example)

    server.GET("/", Index)
    server.GET("/hello/{name}", Hello)

    log.Fatal(http.ListenAndServe(":8080", server))
}
```
## Method Middlewares
```go
package main

import (
	"fmt"
	"log"
	"time"
	"net/http"

    "github.com/vardius/goserver"
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
	//do smth
    next.ServeHTTP(w, r)
  }

  return http.HandlerFunc(fn)
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	params, _ := goserver.FromContext(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    server := goserver.New()

    server.GET("/", http.HandlerFunc(Index))
    server.GET("/hello/{name}", http.HandlerFunc(Hello))

	//apply middlewares to all routes with GET method
	//can pass as many as you want
    server.USE("GET", "", logger, example)

    log.Fatal(http.ListenAndServe(":8080", server))
}
```
## Route Middlewares
```go
package main

import (
	"fmt"
	"log"
	"time"
	"net/http"

    "github.com/vardius/goserver"
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
	//do smth
    next.ServeHTTP(w, r)
  }

  return http.HandlerFunc(fn)
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	params, _ := goserver.FromContext(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    server := goserver.New()
    server.GET("/", http.HandlerFunc(Index))
    server.GET("/hello/{name}", http.HandlerFunc(Hello))

	//apply midlewares to route and all it childs
	//can pass as many as you want
    server.USE("GET", "/hello/{name}", logger, example)

    log.Fatal(http.ListenAndServe(":8080", server))
}
```
