Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

The fastest Go Server/API micro framwework, HTTP request router, multiplexer, mux.

ABOUT
==================================================
Contributors:

* [Rafał Lorenz](http://rafallorenz.com)

Want to contribute ? Feel free to send pull requests!

Have problems, bugs, feature ideas?
We are using the github [issue tracker](https://github.com/vardius/goserver/issues) to manage them.

HOW TO USE
==================================================

[GoDoc](http://godoc.org/github.com/vardius/goserver)
-------
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
	params, _ := goserver.ParamsFromContext(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params["name"])
}

func main() {
    server := goserver.New()
    server.GET("/", Index)
    server.GET("/hello/:name", Hello)

    log.Fatal(http.ListenAndServe(":8080", server))
}
```
## HTTP2 example
```go
package main

import (
    "fmt"
    "log"
    "net/http"
	
    "golang.org/x/net/http2"
    "github.com/vardius/goserver"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
    params, _ := goserver.ParamsFromContext(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params["name"])
}

func main() {
    server := goserver.New()
    server.GET("/", Index)
    server.GET("/hello/:name", Hello)

    http2.ConfigureServer(server, &http2.Server{})
    log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
```
## Serve files
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
    params, _ := goserver.ParamsFromContext(r.Context())
    fmt.Fprintf(w, "hello, %s!\n", params["name"])
}

func main() {
    server := goserver.New()
    server.GET("/", Index)
    server.GET("/hello/:name", Hello)
	//If route not found and the request method equals Get
	//server will serve files from directory
	//second parameter decide if prefix should be striped
    server.ServeFiles("static", false)

    log.Fatal(http.ListenAndServe(":8080", server))
}
```
## Multi-domain / Sub-domains
```go
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

func main() {
	// Initialize a server as usual
	server := goserver.New()
	server.GET("/", Index)
	server.GET("/hello/:name", Hello)

	// Make a new HostSwitch and insert the server (our http handler)
	// for example.com and port 12345
	hs := make(HostSwitch)
	hs["example.com:12345"] = server

	// Use the HostSwitch to listen and serve on port 12345
	log.Fatal(http.ListenAndServe(":12345", hs))
}
```
## Basic Authentication
### Useing middleware
```go
package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/vardius/goserver"
)

type (
	statusError struct {
		code int
		err  error
	}
)

func BasicAuth(h http.Handler) Error {
		fn := func(w http.ResponseWriter, r *http.Request) {
	requiredUser := "gordon"
	requiredPassword := "secret!"
	
	// Get the Basic Authentication credentials
	user, password, hasAuth := r.BasicAuth()
	
	if hasAuth && user == requiredUser && password == requiredPassword {
		return nil;
	} else {
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w,
			http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized,
		)
	}
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Not protected!\n")
}

func Protected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Protected!\n")
}

func main() {
	server := goserver.New(BasicAuth)
	server.GET("/", Index)	
	server.GET("/protected", Protected)

	log.Fatal(http.ListenAndServe(":8080", server))
}
```
### Not useing middleware
```go
package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/vardius/goserver"
)

func BasicAuth(h http.Handler, requiredUser, requiredPassword string) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()
		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, r)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Not protected!\n")
}

func Protected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Protected!\n")
}

func main() {
	user := "gordon"
	pass := "secret!"

	server := goserver.New()
	server.GET("/", Index)
	server.GET("/protected", BasicAuth(Protected, user, pass))

	log.Fatal(http.ListenAndServe(":8080", server))
}
```
## Handle errors
```go
package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/vardius/goserver"
)

type (
	statusError struct {
		code int
		err  error
	}
)

func recoverfunc(next http.Handler) http.Handler {
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
	server := goserver.New(recoverfunc)
	server.GET("/", Index)	
	server.GET("/panic", WithError)

	log.Fatal(http.ListenAndServe(":8080", server))
}
```

License
-------

This package is released under the MIT license. See the complete license in the package:

[LICENSE](LICENSE.md)
