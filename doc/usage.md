Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Usage
----------------
1. [Basic example](#basic-example)

## Basic example
```go
package main

import (
    "fmt"
    "log"
    "net/http"
	
    "github.com/vardius/goserver"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	params, _ := goserver.ParamsFromContext(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params["name"])
}

func main() {
    server := goserver.New()
    server.GET("/", Index)
    server.GET("/hello/:name", Hello)

    log.Fatal(http.ListenAndServe(":8080", server))
}
```

Advanced configuration
----------------
1. [Appling Middleware](middleware.md)
2. [Handling Errors](errors.md)
3. [Serving Files](files.md)
3. [HTTP2](http2.md)
4. [Multi-domain / Sub-domains](multi-domain.md)
4. [Authentication](auth.md)
