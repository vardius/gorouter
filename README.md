gorouter
================
[![Build Status](https://travis-ci.com/vardius/gorouter.svg?branch=master)](https://travis-ci.com/vardius/gorouter)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/gorouter)](https://goreportcard.com/report/github.com/vardius/gorouter)
[![codecov](https://codecov.io/gh/vardius/gorouter/branch/master/graph/badge.svg)](https://codecov.io/gh/vardius/gorouter)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvardius%2Fgorouter.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvardius%2Fgorouter?ref=badge_shield)
[![](https://godoc.org/github.com/vardius/gorouter?status.svg)](http://godoc.org/github.com/vardius/gorouter)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/gorouter/blob/master/LICENSE.md)

Go Server/API micro framework, HTTP request router, multiplexer, mux.

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
2. [Documentation](http://rafallorenz.com/gorouter)
3. [Go Server/API boilerplate using best practices DDD CQRS ES](https://github.com/vardius/go-api-boilerplate)

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

func index(w http.ResponseWriter, _ *http.Request) {
    if _, err := fmt.Fprint(w, "Welcome!\n"); err != nil {
        panic(err)
    }
}

func hello(w http.ResponseWriter, r *http.Request) {
    params, _ := context.Parameters(r.Context())
    if _, err := fmt.Fprintf(w, "hello, %s!\n", params.Value("name")); err != nil {
        panic(err)
    }
}

func main() {
    router := gorouter.New()
    router.GET("/", http.HandlerFunc(index))
    router.GET("/hello/{name}", http.HandlerFunc(hello))

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

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```

License
-------

This package is released under the MIT license. See the complete license in the package:

[LICENSE](LICENSE.md)

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvardius%2Fgorouter.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvardius%2Fgorouter?ref=badge_large)
