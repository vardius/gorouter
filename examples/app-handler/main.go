// Command app-handler uses handlers that return an error, turned into a response in one place.
// Handle takes any http.Handler, so a custom handler type registers like any other.
package main

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/vardius/gorouter/v5"
)

// appHandler is a handler that can fail. Its ServeHTTP writes the error response.
type appHandler func(http.ResponseWriter, *http.Request) error

func (h appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err == nil {
		return
	}

	if he, ok := errors.AsType[*httpError](err); ok {
		http.Error(w, he.message, he.status)
		return
	}
	slog.Error("handler failed", "path", r.URL.Path, "error", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// httpError is an error that carries the response the client should get.
type httpError struct {
	status  int
	message string
}

func (e *httpError) Error() string { return e.message }

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = map[string]user{"1": {ID: "1", Name: "Gopher"}}

func newRouter() *gorouter.Router {
	r := gorouter.New()
	r.Handle(http.MethodGet, "/users/{id}", appHandler(getUser))

	return r
}

func getUser(w http.ResponseWriter, r *http.Request) error {
	u, ok := users[r.PathValue("id")]
	if !ok {
		return &httpError{status: http.StatusNotFound, message: "no such user"}
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(u)
}

func main() {
	srv := &http.Server{Addr: ":8080", Handler: newRouter(), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
