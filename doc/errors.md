Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/goserver)](https://goreportcard.com/report/github.com/vardius/goserver)
[![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)
[![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/goserver/blob/master/LICENSE.md)

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
	server.GET("/", http.HandlerFunc(Index))
	server.GET("/panic", http.HandlerFunc(WithError))

	log.Fatal(http.ListenAndServe(":8080", server))
}
```
