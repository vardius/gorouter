// Command multidomain serves a different router for each host name.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/vardius/gorouter/v5"
)

// hostSwitch picks the handler for a request by its host name, ignoring the port.
type hostSwitch map[string]http.Handler

func (hs hostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host // no port
	}

	if h, ok := hs[host]; ok {
		h.ServeHTTP(w, r)
		return
	}
	http.Error(w, "unknown host", http.StatusMisdirectedRequest)
}

func newHandler() http.Handler {
	site := gorouter.New()
	site.GET("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "the website")
	})

	api := gorouter.New()
	api.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "api user %s\n", r.PathValue("id"))
	})

	return hostSwitch{
		"example.com":     site,
		"www.example.com": site,
		"api.example.com": api,
	}
}

func main() {
	srv := &http.Server{Addr: ":8080", Handler: newHandler(), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
