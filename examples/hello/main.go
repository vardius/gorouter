// Command hello is the smallest gorouter server: a few routes, a param and a regexp param.
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vardius/gorouter/v5"
)

func newRouter() *gorouter.Router {
	r := gorouter.New()
	r.GET("/", index)
	r.GET("/hello/{name}", hello)
	r.GET("/users/{id:[0-9]+}", user)

	return r
}

func index(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "welcome")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %s!\n", r.PathValue("name"))
}

func user(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "user #%s\n", r.PathValue("id"))
}

func main() {
	srv := &http.Server{Addr: ":8080", Handler: newRouter(), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
