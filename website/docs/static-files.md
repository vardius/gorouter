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

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
    params, _ := context.Parameters(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    router := gorouter.New()
    router.GET("/", http.HandlerFunc(Index))
    router.GET("/hello/{name}", http.HandlerFunc(Hello))
    // If route not found and the request method equals Get
    // router will serve files from directory
    // third parameter decide if prefix should be striped
    router.ServeFiles(http.Dir("static"), "static", false)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
<!--valyala/fasthttp-->
```go
// Example coming soon...
```
<!--END_DOCUSAURUS_CODE_TABS-->
