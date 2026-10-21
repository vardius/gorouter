package gorouter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vardius/gorouter/v5"
)

// discardWriter is a ResponseWriter that allocates nothing, so only the router's allocations are counted.
type discardWriter struct{ header http.Header }

func (w *discardWriter) Header() http.Header         { return w.header }
func (w *discardWriter) Write(b []byte) (int, error) { return len(b), nil }
func (w *discardWriter) WriteHeader(int)             {}

func noop(http.ResponseWriter, *http.Request) {}

// passThrough is middleware that allocates nothing per request.
func passThrough(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
}

// TestAllocations guards the router's allocations per request. Lower a number when an optimisation lands;
// never raise one without a reason in the commit.
func TestAllocations(t *testing.T) {
	r := gorouter.New(
		gorouter.WithNotFound(http.HandlerFunc(noop)),
		gorouter.WithMethodNotAllowed(http.HandlerFunc(noop)),
	)
	r.GET("/", noop)
	r.GET("/a/b/c/d/e", noop)
	r.GET("/users/{id}", noop)
	r.GET("/repos/{owner}/{repo}/issues/{number:[0-9]+}", noop)
	r.GET("/admin/stats", noop)
	r.Use("/admin", passThrough, passThrough, passThrough)
	sub := gorouter.New()
	sub.GET("/status", noop)
	r.Mount("/v1", sub)

	tests := []struct {
		name, method, target string
		want                 float64
	}{
		{"root", "GET", "/", 0},
		{"static, 5 segments", "GET", "/a/b/c/d/e", 0},
		{"static with 3 middlewares", "GET", "/admin/stats", 0},
		{"404", "GET", "/missing", 0},
		// net/http's own SetPathValue allocates its map: the floor for params on a standard *http.Request.
		{"1 param", "GET", "/users/42", 2},
		{"3 params with a regexp", "GET", "/repos/vardius/gorouter/issues/7", 2},
		// The Allow value and http.Header.Set's slice.
		{"405", "POST", "/users/42", 2},
		// One copy of the request and its URL, as http.StripPrefix makes.
		{"mount", "GET", "/v1/status", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &discardWriter{header: http.Header{}}
			template := httptest.NewRequest(tt.method, tt.target, nil)
			req := new(http.Request)

			got := testing.AllocsPerRun(1000, func() {
				*req = *template // a fresh request each run, without allocating one
				r.ServeHTTP(w, req)
			})
			if got != tt.want {
				t.Errorf("%s %s allocates %v times per request, want %v", tt.method, tt.target, got, tt.want)
			}
		})
	}
}
