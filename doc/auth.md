Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Authentication
----------------
1. [Basic Authentication](#basic-authentication)

## Basic Authentication
```go
package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/vardius/goserver"
)

func BasicAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
        requiredUser := "gordon"
        requiredPassword := "secret!"
        
        // Get the Basic Authentication credentials
        user, password, hasAuth := r.BasicAuth()
        
        if hasAuth && user == requiredUser && password == requiredPassword {
            return nil;
        } else {
            w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
            http.Error(w,
                http.StatusText(http.StatusUnauthorized),
                http.StatusUnauthorized,
            )
        }
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Not protected!\n")
}

func Protected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Protected!\n")
}

func main() {
	server := goserver.New()
	server.GET("/", http.HandlerFunc(Index))	
	server.GET("/protected", http.HandlerFunc(Protected))

	server.USE("GET", "/protected", BasicAuth)

	log.Fatal(http.ListenAndServe(":8080", server))
}
```
