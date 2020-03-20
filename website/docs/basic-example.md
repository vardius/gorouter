---
id: basic-example
title: Basic example
sidebar_label: Basic example
---

[gorouter](github.com/vardius/gorouter) supports following http implementations:

- [net/http](https://golang.org/pkg/net/http/)
- [fasthttp](https://github.com/valyala/fasthttp)

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

func index(w http.ResponseWriter, _ *http.Request) {
    if _, err := fmt.Fprint(w, "Welcome!\n"); err != nil {
        panic(err)
    }
}

func hello(w http.ResponseWriter, r *http.Request) {
    params, _ := context.Parameters(r.Context())
    if _, err := fmt.Fprintf(w, "hello, %s!\n", params.Value("name")); err != nil {
        panic(err)
    }
}

func main() {
    router := gorouter.New()
    router.GET("/", http.HandlerFunc(index))
    router.GET("/hello/{name}", http.HandlerFunc(hello))

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
    "github.com/vardius/gorouter/v4/context"
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

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```
<!--END_DOCUSAURUS_CODE_TABS-->
