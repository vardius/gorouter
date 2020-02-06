---
id: sub-router
title: Mounting Sub-Router
sidebar_label: Mounting Sub-Router
---

When having multiple instance of a router you might want to mount one as a sub router of another under some route path, still keeping all middleware.

It doesn't have to be [gorouter](github.com/vardius/gorouter). You can mount other routers as well as long they implement `http.Handler` interface.

## Mount

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```go
package main

import (
   "log"
   "net/http"

   "github.com/vardius/gorouter/v4"
)

func main() {
    router := gorouter.New()
    subrouter := gorouter.New()

    router.Mount("/{param}", subrouter)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
<!--valyala/fasthttp-->
```go
package main

import (
    "log"

    "github.com/valyala/fasthttp"
    "github.com/vardius/gorouter/v4"
)

func main() {
    router := gorouter.NewFastHTTPRouter()
    subrouter := gorouter.NewFastHTTPRouter()

    router.Mount("/{param}", subrouter)

    log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleFastHTTP))
}
```
<!--END_DOCUSAURUS_CODE_TABS-->

Given example will result in all routes of a `subrouter` being available under paths prefixed with a mount path.