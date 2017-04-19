Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Handling Errors
----------------
1. [Recover Middleware](#recover-middleware)

## Recover Middleware
```go
package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/vardius/goserver"
)

func recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok!\n")
}

func WithError(w http.ResponseWriter, r *http.Request) {
	panic("panic recover")
}

func main() {
	server := goserver.New(recoverfunc)
	server.GET("/", Index)	
	server.GET("/panic", WithError)

	log.Fatal(http.ListenAndServe(":8080", server))
}
```
