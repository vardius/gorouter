---
id: multidomain
title: Multidomain
sidebar_label: Multidomain
---

## HostSwitch

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```go
import (
    "fmt"
    "log"
    "net/http"

    "github.com/vardius/gorouter/v4"
    "github.com/vardius/gorouter/v4/context"
)

// We need an object that implements the http.Handler interface.
// Therefore we need a type for which we implement the ServeHTTP method.
// We just use a map here, in which we map host names (with port) to http.Handlers
type HostSwitch map[string]http.Handler

// Implement the ServerHTTP method on our new type
func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if a http.Handler is registered for the given host.
	// If yes, use it to handle the request.
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		// Handle host names for wich no handler is registered
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}

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
	// Initialize a router as usual
	router := gorouter.New()
	router.GET("/", http.HandlerFunc(Index))
	router.GET("/hello/{name}", http.HandlerFunc(Hello))

	// Make a new HostSwitch and insert the router (our http handler)
	// for example.com and port 12345
	hs := make(HostSwitch)
	hs["example.com:12345"] = router

	// Use the HostSwitch to listen and serve on port 12345
	log.Fatal(http.ListenAndServe(":12345", hs))
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

// HostSwitch is the host-handler map
// We need an object that implements the fasthttp.RequestHandler interface.
// We just use a map here, in which we map host names (with port) to fasthttp.RequestHandlers
type HostSwitch map[string]fasthttp.RequestHandler

// CheckHost Implement a CheckHost method on our new type
func (hs HostSwitch) CheckHost(ctx *fasthttp.RequestCtx) {
	// Check if a http.Handler is registered for the given host.
	// If yes, use it to handle the request.
	if handler := hs[string(ctx.Host())]; handler != nil {
		handler(ctx)
	} else {
		// Handle host names for wich no handler is registered
		ctx.Error("Forbidden", 403) // Or Redirect?
	}
}

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

	// Make a new HostSwitch and insert the router (our http handler)
	// for example.com and port 12345
	hs := make(HostSwitch)
	hs["example.com:12345"] = router.HandleFastHTTP

	// Use the HostSwitch to listen and serve on port 12345
	log.Fatal(fasthttp.ListenAndServe(":12345", hs.CheckHost))
}
```
<!--END_DOCUSAURUS_CODE_TABS-->
