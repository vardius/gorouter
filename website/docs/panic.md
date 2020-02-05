---
id: panic
title: Panic Recovery
sidebar_label: Panic Recovery
---

## Recover Middleware
```go
package main

import (
	"fmt"
	"log"
	"net/http"

        "github.com/vardius/gorouter/v4"
)

func recoverMiddleware(next http.Handler) http.Handler {
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
	router := gorouter.New(recoverMiddleware)
	router.GET("/", http.HandlerFunc(Index))
	router.GET("/panic", http.HandlerFunc(WithError))

	log.Fatal(http.ListenAndServe(":8080", router))
}
```
