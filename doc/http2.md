Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

The fastest Go Server/API micro framwework, HTTP request router, multiplexer, mux.

HTTP2
----------------
1. [Pusher](#pusher)
2. [The Go Blog - HTTP/2 Server Push](https://blog.golang.org/h2push)

## HTTP2
```go
package main

import (
    "log"
    "net/http"
	
    "golang.org/x/net/http2"
    "github.com/vardius/goserver"
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
    server := goserver.New()
    server.GET("/", Pusher)

    http2.ConfigureServer(server, &http2.Server{})
    log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
```
