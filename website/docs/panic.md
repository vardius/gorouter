---
id: panic
title: Panic Recovery
sidebar_label: Panic Recovery
---

## Recover Middleware

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```go
package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/vardius/gorouter/v4"
)

func recoverMiddleware(next http.Handler) http.Handler {
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
	router := gorouter.New(recoverMiddleware)
	router.GET("/", http.HandlerFunc(Index))
	router.GET("/panic", http.HandlerFunc(WithError))

	log.Fatal(http.ListenAndServe(":8080", router))
}
```
<!--valyala/fasthttp-->
```go
package main

import (
    "fmt"
    "log"

    "github.com/valyala/fasthttp"
    "github.com/vardius/gorouter/v4"
)

func recoverMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	fn := func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rcv := recover(); rcv != nil {
                ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
			}
		}()

		next(ctx)
	}

	return fn
}

func index(_ *fasthttp.RequestCtx) {
    fmt.Print("Welcome!\n")
}

func withError(ctx *fasthttp.RequestCtx) {
	panic("panic recover")
}

func main() {
    router := gorouter.NewFastHTTPRouter()
    router.GET("/", index)
    router.GET("/panic", withError)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```
<!--END_DOCUSAURUS_CODE_TABS-->
