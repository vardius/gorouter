package gorouter_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/vardius/gorouter/v5"
	"github.com/vardius/gorouter/v5/internal/benchdata"
)

// serveLoop serves the request repeatedly. Each run gets a fresh copy of it, as a server would,
// so param values set on one run don't make the next one cheaper.
func serveLoop(b *testing.B, h http.Handler, method, target string) {
	w := &discardWriter{header: http.Header{}}
	template := httptest.NewRequest(method, target, nil)
	req := new(http.Request)

	b.ReportAllocs()
	for b.Loop() {
		*req = *template
		h.ServeHTTP(w, req)
	}
}

// repeat joins n copies of segment into a path.
func repeat(segment string, n int) string {
	return strings.Repeat("/"+segment, n)
}

func BenchmarkStatic(b *testing.B) {
	for _, n := range []int{1, 5, 20} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			r := gorouter.New()
			r.GET(repeat("static", n), noop)
			serveLoop(b, r, "GET", repeat("static", n))
		})
	}
}

func BenchmarkParam(b *testing.B) {
	for _, n := range []int{1, 5, 20} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			var pattern strings.Builder
			for i := range n {
				pattern.WriteString("/{p" + strconv.Itoa(i) + "}")
			}
			r := gorouter.New()
			r.GET(pattern.String(), noop)
			serveLoop(b, r, "GET", repeat("value", n))
		})
	}
}

func BenchmarkRegexp(b *testing.B) {
	for _, n := range []int{1, 5} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			var pattern strings.Builder
			for i := range n {
				pattern.WriteString("/{p" + strconv.Itoa(i) + ":[a-z]+}")
			}
			r := gorouter.New()
			r.GET(pattern.String(), noop)
			serveLoop(b, r, "GET", repeat("value", n))
		})
	}
}

func BenchmarkMiddleware(b *testing.B) {
	r := gorouter.New()
	r.Use("/", passThrough, passThrough, passThrough)
	r.GET("/users/new", noop)
	serveLoop(b, r, "GET", "/users/new")
}

func BenchmarkNotFound(b *testing.B) {
	r := gorouter.New(gorouter.WithNotFound(http.HandlerFunc(noop)))
	r.GET("/users/{id}", noop)
	serveLoop(b, r, "GET", "/missing/path")
}

func BenchmarkMethodNotAllowed(b *testing.B) {
	r := gorouter.New(gorouter.WithMethodNotAllowed(http.HandlerFunc(noop)))
	r.GET("/users/{id}", noop)
	serveLoop(b, r, "POST", "/users/42")
}

func BenchmarkMount(b *testing.B) {
	sub := gorouter.New()
	sub.GET("/status", noop)
	r := gorouter.New()
	r.Mount("/v1", sub)
	serveLoop(b, r, "GET", "/v1/status")
}

// BenchmarkAPI routes every request of a realistic API once per op and reports the cost per request.
func BenchmarkAPI(b *testing.B) {
	r := gorouter.New()
	reqs := make([]*http.Request, 0, len(benchdata.API))
	for _, route := range benchdata.API {
		r.Handle(route.Method, route.Pattern, http.HandlerFunc(noop))
		reqs = append(reqs, httptest.NewRequest(route.Method, route.Path, nil))
	}

	w := &discardWriter{header: http.Header{}}
	req := new(http.Request)

	b.ReportAllocs()
	for b.Loop() {
		for _, template := range reqs {
			*req = *template
			r.ServeHTTP(w, req)
		}
	}
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*len(reqs)), "ns/req")
}
