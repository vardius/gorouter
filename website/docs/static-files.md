---
id: static-files
title: Static Files
sidebar_label: Static Files
---

## Static Files

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```go
package main

import (
    "fmt"
    "log"
    "net/http"
	
    "github.com/vardius/gorouter/v4"
    "github.com/vardius/gorouter/v4/context"
)

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request) {
    params, _ := context.Parameters(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    router := gorouter.New()
    router.GET("/", http.HandlerFunc(index))
    router.GET("/hello/{name}", http.HandlerFunc(hello))

    // If route not found and the request method equals Get
    // router will serve files from directory
    // third parameter decide if prefix should be striped
    router.ServeFiles(http.Dir("static"), "static", false)

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

    // If route not found and the request method equals Get
    // router will serve files from directory
	// Will serve files from /var/www/static with path /static/*
	// because strip 1 slash (/static/favicon.ico -> /favicon.ico)
	router.ServeFiles("/var/www/static", 1)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```
<!--END_DOCUSAURUS_CODE_TABS-->
