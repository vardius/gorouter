Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Serving Files
----------------
1. [Static Files](#static-files)

## Static Files
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
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    server := goserver.New()
    server.GET("/", http.HandlerFunc(Index))
    server.GET("/hello/{name}", http.HandlerFunc(Hello))
	//If route not found and the request method equals Get
	//server will serve files from directory
	//second parameter decide if prefix should be striped
    server.ServeFiles("static", false)

    log.Fatal(http.ListenAndServe(":8080", server))
}
```
