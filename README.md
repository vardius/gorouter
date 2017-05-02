Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/goserver)](https://goreportcard.com/report/github.com/vardius/goserver)
[![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)
[![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/goserver/blob/master/LICENSE.md)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

ABOUT
==================================================
Contributors:

* [Rafał Lorenz](http://rafallorenz.com)

Want to contribute ? Feel free to send pull requests!

Have problems, bugs, feature ideas?
We are using the github [issue tracker](https://github.com/vardius/goserver/issues) to manage them.

HOW TO USE
==================================================

1. [GoDoc](http://godoc.org/github.com/vardius/goserver)
2. [Documentation](https://github.com/vardius/goserver/wiki)
3. [Benchmarks](doc/benchmark.md)

## Basic example
```go
package main

import (
    "fmt"
    "log"
    "net/http"
	
    "github.com/vardius/goserver"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	params, _ := goserver.FromContext(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params.Value("name"))
}

func main() {
    server := goserver.New()
    server.GET("/", http.HandlerFunc(Index))
    server.GET("/hello/{name}", http.HandlerFunc(Hello))

    log.Fatal(http.ListenAndServe(":8080", server))
}
```

License
-------

This package is released under the MIT license. See the complete license in the package:

[LICENSE](LICENSE.md)
