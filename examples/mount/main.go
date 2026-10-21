// Command mount splits an API into sub-routers, each mounted under a path of the main router.
// A sub-router sees paths relative to its mount point, and params in the mount path stay readable.
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
	r.GET("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "see /api/v1 and /orgs/{org}")
	})
	r.Mount("/api/v1", apiV1())
	r.Mount("/orgs/{org}", orgs())

	return r
}

// apiV1 knows nothing about where it is mounted: its routes start at /.
func apiV1() *gorouter.Router {
	r := gorouter.New()
	r.GET("/status", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "v1 ok")
	})
	r.GET("/users/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "v1 user %s\n", r.PathValue("id"))
	})

	return r
}

// orgs reads org from the mount path and repo from its own route.
func orgs() *gorouter.Router {
	r := gorouter.New()
	r.GET("/repos/{repo}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s/%s (path %s)\n", r.PathValue("org"), r.PathValue("repo"), r.URL.Path)
	})

	return r
}

func main() {
	srv := &http.Server{Addr: ":8080", Handler: newRouter(), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
