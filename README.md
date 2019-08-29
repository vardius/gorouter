gorouter
================
[![Build Status](https://travis-ci.org/vardius/gorouter.svg?branch=master)](https://travis-ci.org/vardius/gorouter)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/gorouter)](https://goreportcard.com/report/github.com/vardius/gorouter)
[![codecov](https://codecov.io/gh/vardius/gorouter/branch/master/graph/badge.svg)](https://codecov.io/gh/vardius/gorouter)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvardius%2Fgorouter.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvardius%2Fgorouter?ref=badge_shield)
[![](https://godoc.org/github.com/vardius/gorouter?status.svg)](http://godoc.org/github.com/vardius/gorouter)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/gorouter/blob/master/LICENSE.md)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

ABOUT
==================================================
Contributors:

* [Rafał Lorenz](http://rafallorenz.com)

Want to contribute ? Feel free to send pull requests!

Have problems, bugs, feature ideas?
We are using the github [issue tracker](https://github.com/vardius/gorouter/issues) to manage them.

HOW TO USE
==================================================

1. [GoDoc](http://godoc.org/github.com/vardius/gorouter)
2. [Documentation](https://github.com/vardius/gorouter/wiki)
3. [Benchmarks](https://github.com/vardius/gorouter/wiki/Benchmarks)
4. [Go Server/API boilerplate using best practices DDD CQRS ES](https://github.com/vardius/go-api-boilerplate)

## Basic example
### [net/http](https://golang.org/pkg/net/http/)
```go
package main

import (
    "fmt"
    "log"
    "net/http"
	
    "github.com/vardius/gorouter/v4"
    "github.com/vardius/gorouter/v4/context"
)

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

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
### [fasthttp](https://github.com/valyala/fasthttp)
```go
package main

import (
    "fmt"
    "log"

	"github.com/valyala/fasthttp"
    "github.com/vardius/gorouter/v4"
)

func Index(ctx *fasthttp.RequestCtx) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(ctx *fasthttp.RequestCtx) {
    fmt.Fprintf(w, "hello, %s!\n", ctx.UserValue("name"))
}

func main() {
    router := gorouter.NewFastHTTPRouter()
    router.GET("/", Index)
    router.GET("/hello/{name}", Hello)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```

## Advanced examples
- [Routing](https://github.com/vardius/gorouter/wiki/Routing)
- [Middleware](https://github.com/vardius/gorouter/wiki/Middleware)
- [Mounting Sub Router](https://github.com/vardius/gorouter/wiki/Mounting-Sub-Router)
- [Serving Files](https://github.com/vardius/gorouter/wiki/Serving-Files)
- [Authentication](https://github.com/vardius/gorouter/wiki/Authentication)
- [Handling Panic](https://github.com/vardius/gorouter/wiki/Handling-Panic)
- [HTTP2](https://github.com/vardius/gorouter/wiki/HTTP2)
- [Multidomain](https://github.com/vardius/gorouter/wiki/Multidomain)

License
-------

This package is released under the MIT license. See the complete license in the package:

[LICENSE](LICENSE.md)

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvardius%2Fgorouter.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvardius%2Fgorouter?ref=badge_large)
