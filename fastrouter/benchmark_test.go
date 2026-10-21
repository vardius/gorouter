package fastrouter_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/vardius/gorouter/v5/fastrouter"
	"github.com/vardius/gorouter/v5/internal/benchdata"
)

// serveLoop serves the request repeatedly on one reused ctx, resetting its user values as fasthttp's server does.
// The URI is parsed once, so the loop measures routing rather than fasthttp's parser.
func serveLoop(b *testing.B, r *fastrouter.Router, method, uri string) {
	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod(method)
	ctx.Request.SetRequestURI(uri)

	b.ReportAllocs()
	for b.Loop() {
		r.HandleFastHTTP(ctx)
		ctx.ResetUserValues()
	}
}

// repeat joins n copies of segment into a path.
func repeat(segment string, n int) string {
	return strings.Repeat("/"+segment, n)
}

func BenchmarkStatic(b *testing.B) {
	for _, n := range []int{1, 5, 20} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			r := fastrouter.New()
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
			r := fastrouter.New()
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
			r := fastrouter.New()
			r.GET(pattern.String(), noop)
			serveLoop(b, r, "GET", repeat("value", n))
		})
	}
}

func BenchmarkMiddleware(b *testing.B) {
	r := fastrouter.New()
	r.Use("/", passThrough, passThrough, passThrough)
	r.GET("/users/new", noop)
	serveLoop(b, r, "GET", "/users/new")
}

func BenchmarkNotFound(b *testing.B) {
	r := fastrouter.New(fastrouter.WithNotFound(noop))
	r.GET("/users/{id}", noop)
	serveLoop(b, r, "GET", "/missing/path")
}

func BenchmarkMethodNotAllowed(b *testing.B) {
	r := fastrouter.New(fastrouter.WithMethodNotAllowed(noop))
	r.GET("/users/{id}", noop)
	serveLoop(b, r, "POST", "/users/42")
}

// BenchmarkMount includes parsing the URI on every run, because the mount rewrites the path.
func BenchmarkMount(b *testing.B) {
	sub := fastrouter.New()
	sub.GET("/status", noop)
	r := fastrouter.New()
	r.Mount("/v1", sub.HandleFastHTTP)

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod("GET")

	b.ReportAllocs()
	for b.Loop() {
		ctx.Request.SetRequestURI("/v1/status")
		r.HandleFastHTTP(ctx)
		ctx.ResetUserValues()
	}
}

// BenchmarkAPI routes every request of a realistic API once per op and reports the cost per request.
func BenchmarkAPI(b *testing.B) {
	r := fastrouter.New()
	ctxs := make([]*fasthttp.RequestCtx, 0, len(benchdata.API))
	for _, route := range benchdata.API {
		r.Handle(route.Method, route.Pattern, noop)

		ctx := new(fasthttp.RequestCtx)
		ctx.Request.Header.SetMethod(route.Method)
		ctx.Request.SetRequestURI(route.Path)
		ctxs = append(ctxs, ctx)
	}

	b.ReportAllocs()
	for b.Loop() {
		for _, ctx := range ctxs {
			r.HandleFastHTTP(ctx)
			ctx.ResetUserValues()
		}
	}
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*len(ctxs)), "ns/req")
}
