// Command static-files serves an embedded directory of files next to regular routes.
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/vardius/gorouter/v5"
)

//go:embed public
var public embed.FS

func newRouter() (*gorouter.Router, error) {
	files, err := fs.Sub(public, "public")
	if err != nil {
		return nil, err
	}

	r := gorouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, files, "index.html")
	})
	r.GET("/api/time", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, time.Now().UTC().Format(time.RFC3339))
	})

	// Every request below /static goes to the file server, with /static removed from the path.
	r.Mount("/static", http.FileServerFS(files))

	return r, nil
}

func main() {
	r, err := newRouter()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{Addr: ":8080", Handler: r, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
