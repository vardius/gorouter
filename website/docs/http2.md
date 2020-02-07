---
id: http2
title: HTTP2
sidebar_label: HTTP2
---

## [The Go Blog - HTTP/2 Server Push](https://blog.golang.org/h2push)

## Pusher
<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```go
package main

import (
    "log"
    "net/http"
	
    "golang.org/x/net/http2"
    "github.com/vardius/gorouter/v4"
)

func Pusher(w http.ResponseWriter, r *http.Request) {
    if pusher, ok := w.(http.Pusher); ok {
        // Push is supported.
        options := &http.PushOptions{
            Header: http.Header{
                "Accept-Encoding": r.Header["Accept-Encoding"],
            },
        }
        if err := pusher.Push("/script.js", options); err != nil {
            log.Printf("Failed to push: %v", err)
        }
    }
    // ...
}

func main() {
    router := gorouter.New()
    router.GET("/", http.HandlerFunc(Pusher))

    http2.ConfigureServer(router, &http2.Server{})
    log.Fatal(router.ListenAndServeTLS("router.crt", "router.key"))
}
```
<!--valyala/fasthttp-->
HTTP/2 implementation for fasthttp is [under construction...](https://github.com/fasthttp/http2)
<!--END_DOCUSAURUS_CODE_TABS-->