Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/goserver)](https://goreportcard.com/report/github.com/vardius/goserver)
[![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver)
[![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/goserver/blob/master/LICENSE.md)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Multi-domain / Sub-domains
----------------
1. [HostSwitch](#hostswitch)

## HostSwitch
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
	server.GET("/", http.HandlerFunc(Index))
	server.GET("/hello/{name}", http.HandlerFunc(Hello))

	// Make a new HostSwitch and insert the server (our http handler)
	// for example.com and port 12345
	hs := make(HostSwitch)
	hs["example.com:12345"] = server

	// Use the HostSwitch to listen and serve on port 12345
	log.Fatal(http.ListenAndServe(":12345", hs))
}
```
