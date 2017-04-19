Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Usage
----------------
1. [Basic example](#basic-example)
2. [Routing](#routing)

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
## Routing
The router determines how to handle that request. Goserver uses a routing tree. Once one branch of the tree matches, only routes inside that branch are considered, not any routes after that branch. When instantiating server, the root node of router tree is created.
### Defining Routes
A full route definition contain up to three parts:
1. HTTP method under which route will be available
2. The URL path route. This is matched against the URL passed to the server, and can contain named wildcard placeholders *(e.g. :placeholders)* to match dynamic parts in the URL.
3. `http.HandleFunc`, which tells the server to handle matched requests to the router with handler.
Take the following example:
```
server.GET("/hello/:name:r([a-z]+)go", func Hello(w http.ResponseWriter, r *http.Request) {
    params, _ := goserver.ParamsFromContext(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params["name"])
})
```
In this case, the route is matched by `/hello/rxxxxxgo` for example, because the `:name` wildcard matches the regular expression wildcard given (`r([a-z]+)go`). However, `/hello/foo` does not match, because "foo" fails the *name* wildcard. When using wildcards, these are returned in the map from request context. The part of the path that the wildcard matched (e.g. *rxxxxxgo*) is used as value.
### Route types
- Static `/hello`
will match requests matching given route
- Named `/:name`
will match requests matching given route scheme
- Regexp `/:name:[a-z]+`
will match requests matching given route scheme and its regexp

Advanced configuration
----------------
1. [Appling Middleware](middleware.md)
2. [Handling Errors](errors.md)
3. [Serving Files](files.md)
3. [HTTP2](http2.md)
4. [Multi-domain / Sub-domains](multi-domain.md)
4. [Authentication](auth.md)
