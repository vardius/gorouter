// Command tls serves over HTTPS. net/http negotiates HTTP/2 on TLS connections by itself,
// so the router serves HTTP/2 clients with no extra setup.
//
// Behind a proxy that terminates TLS, run it with -h2c to accept HTTP/2 without TLS as well.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vardius/gorouter/v5"
)

func newRouter() *gorouter.Router {
	r := gorouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello over %s\n", r.Proto)
	})

	return r
}

func main() {
	addr := flag.String("addr", ":8443", "address to listen on")
	cert := flag.String("cert", "cert.pem", "TLS certificate file")
	key := flag.String("key", "key.pem", "TLS key file")
	h2c := flag.Bool("h2c", false, "serve HTTP/1 and unencrypted HTTP/2 without TLS")
	flag.Parse()

	srv := &http.Server{Addr: *addr, Handler: newRouter(), ReadHeaderTimeout: 5 * time.Second}

	if *h2c {
		var protocols http.Protocols
		protocols.SetHTTP1(true)
		protocols.SetUnencryptedHTTP2(true)
		srv.Protocols = &protocols
		log.Fatal(srv.ListenAndServe())
	}

	log.Fatal(srv.ListenAndServeTLS(*cert, *key))
}
