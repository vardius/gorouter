---
id: basic-authentication
title: Basic Authentication
sidebar_label: Basic Authentication
---

## Basic Authentication

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"crypto/subtle"

	"github.com/vardius/gorouter/v4"
)

const (
	requiredUser     = []byte("gordon")
	requiredPassword = []byte("secret!")
)

func BasicAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if !hasAuth || subtle.ConstantTimeCompare(requiredUser, []byte(user)) != 1 || subtle.ConstantTimeCompare(requiredPassword, []byte(pass)) != 1 {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Not protected!\n")
}

func protected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Protected!\n")
}

func main() {
	router := gorouter.New()
	router.GET("/", http.HandlerFunc(index))	
	router.GET("/protected", http.HandlerFunc(protected))

	router.USE("GET", "/protected", BasicAuth)

	log.Fatal(http.ListenAndServe(":8080", router))
}
```
<!--valyala/fasthttp-->
```go
package main

import (
	"crypto/subtle"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/vardius/gorouter/v4"
)

const (
	basicAuthPrefix = []byte("Basic ")
	requiredUser     = []byte("gordon")
	requiredPassword = []byte("secret!")
)

func BasicAuth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		// Get the Basic Authentication credentials
		auth := ctx.Request.Header.Peek("Authorization")
		if bytes.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(string(auth[len(basicAuthPrefix):]))
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 && subtle.ConstantTimeCompare(requiredUser, pair[0]) == 1 && subtle.ConstantTimeCompare(requiredPassword, pair[1]) == 1 {
					// Delegate request to the given handle
					next(ctx)
					return
				}
			}
		}

		// Request Basic Authentication otherwise
		ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
	}

	return fn
}

func index(_ *fasthttp.RequestCtx) {
    fmt.Print("Not Protected!\n")
}

func protected(_ *fasthttp.RequestCtx) {
    fmt.Print("Protected!\n")
}

func main() {
    router := gorouter.NewFastHTTPRouter()
    router.GET("/", index)
    router.GET("/protected", protected)

    router.USE("GET", "/protected", BasicAuth)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```
<!--END_DOCUSAURUS_CODE_TABS-->
