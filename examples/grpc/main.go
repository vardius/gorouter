// Command grpc serves a gRPC service and a plain HTTP endpoint from one router, on one port.
// The service is gRPC health checking from connectrpc.com/grpchealth, so the example needs no generated code.
//
// gRPC runs over HTTP/2. The server accepts HTTP/2 over TLS, and HTTP/1 and unencrypted HTTP/2 without it,
// as it would behind a proxy that terminates TLS. It serves without TLS here; ListenAndServeTLS serves with it.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"github.com/vardius/gorouter/v5"
)

// maxMessageBytes caps each message after decompression. Connect has no limit by default,
// so without one a small compressed request can make the handler allocate gigabytes.
const maxMessageBytes = 4 << 20

// newRouter serves the health service, ending each of its calls after callTimeout, and /status.
func newRouter(callTimeout time.Duration) *gorouter.Router {
	r := gorouter.New()
	r.GET("/status", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	// A gRPC call is a POST to /package.Service/Method, and the service's handler picks the method from that
	// full path. Route the methods with a param, not Mount: Mount strips the service name from the path,
	// so the handler would find no method to call.
	path, health := grpchealth.NewHandler(grpchealth.NewStaticChecker(), connect.WithReadMaxBytes(maxMessageBytes))
	r.Handle(http.MethodPost, path+"{method}", health)
	// The server's timeouts don't end a stream that sends nothing over HTTP/1.1, as the Connect protocol
	// and gRPC-Web allow, so every call gets a deadline of its own.
	r.Use(path, withDeadline(callTimeout))

	return r
}

// withDeadline cancels each request's context after d. Connect ends the call when its context is done.
func withDeadline(d time.Duration) gorouter.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func newServer(addr string, h http.Handler) *http.Server {
	var protocols http.Protocols
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(true)
	protocols.SetUnencryptedHTTP2(true)

	return &http.Server{
		Addr:              addr,
		Handler:           h,
		Protocols:         &protocols,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       time.Minute,
		// Longer than a call's deadline, so the deadline ends a call cleanly before the server resets its stream.
		WriteTimeout: 2 * time.Minute,
		IdleTimeout:  2 * time.Minute,
	}
}

func main() {
	log.Fatal(newServer(":8080", newRouter(time.Minute)).ListenAndServe())
}
