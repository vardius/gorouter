// Command middleware shows the three places middleware can go: around the whole router,
// around the routes below a path, and around one method below a path.
package main

import (
	"io"
	"log"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/vardius/gorouter/v5"
)

func newRouter(logger *slog.Logger) *gorouter.Router {
	// Around the whole router: these also see 404 and 405 responses.
	r := gorouter.New(gorouter.WithMiddleware(logRequests(logger), recoverPanics(logger)))

	// Around every route below /api, on every method.
	r.Use("/api", noStore)

	// Around POST routes below /api only.
	r.UseMethod(http.MethodPost, "/api", requireJSON)

	r.GET("/", index)
	r.GET("/api/status", status)
	r.POST("/api/echo", echo)
	r.GET("/panic", func(http.ResponseWriter, *http.Request) { panic("something broke") })

	return r
}

func index(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "welcome\n")
}

func status(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
}

func echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, http.MaxBytesReader(w, r.Body, 1<<20))
}

// logRequests logs each request's method, path and duration.
func logRequests(logger *slog.Logger) gorouter.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("request", "method", r.Method, "path", r.URL.Path, "took", time.Since(start))
		})
	}
}

// recoverPanics turns a panicking handler into a 500 response.
func recoverPanics(logger *slog.Logger) gorouter.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic", "path", r.URL.Path, "error", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// noStore keeps API responses out of caches.
func noStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// requireJSON rejects request bodies that aren't JSON.
func requireJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type")); mediaType != "application/json" {
			http.Error(w, "send JSON", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	srv := &http.Server{Addr: ":8080", Handler: newRouter(logger), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
