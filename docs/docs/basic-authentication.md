---
id: basic-authentication
title: Basic Authentication
sidebar_label: Basic Authentication
---

## Basic Authentication
```go
package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/vardius/gorouter/v4"
)

func BasicAuth(next http.Handler) http.Handler {
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
	router := gorouter.New()
	router.GET("/", http.HandlerFunc(Index))	
	router.GET("/protected", http.HandlerFunc(Protected))

	router.USE("GET", "/protected", BasicAuth)

	log.Fatal(http.ListenAndServe(":8080", router))
}
```