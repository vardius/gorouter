---
id: apphandler
title: App Handler
sidebar_label: App Handler
---

## Use custom handler type in your application

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/vardius/gorouter/v4"
)

type AppHandlerFunc func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r) and handles error
func (f AppHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err.Error())
	}
}

func Index(w http.ResponseWriter, r *http.Request) error {
	return errors.New("I am app handler which can return error")
}

func main() {
	router := gorouter.New()
	router.GET("/", AppHandlerFunc(Index))

	log.Fatal(http.ListenAndServe(":8080", router))
}
```
<!--valyala/fasthttp-->
```go
package main

import (
	"errors"
	"log"

    "github.com/valyala/fasthttp"
    "github.com/vardius/gorouter/v4"
)

type AppHandlerFunc func(ctx *fasthttp.RequestCtx) error

// HandleFastHTTP calls f(ctx) and handles error
func (f AppHandlerFunc) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	if err := f(ctx); err != nil {
		ctx.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}

func index(_ *fasthttp.RequestCtx) error {
	return errors.New("I am app handler which can return error")
}

func main() {
	router := gorouter.NewFastHTTPRouter()
	router.GET("/", AppHandlerFunc(index))

	log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```
<!--END_DOCUSAURUS_CODE_TABS-->
