package fastrouter_test

import (
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/vardius/gorouter/v5/fastrouter"
	"github.com/vardius/gorouter/v5/internal/race"
)

func noop(*fasthttp.RequestCtx) {}

// passThrough is middleware that allocates nothing per request.
func passThrough(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) { next(ctx) }
}

// TestAllocations guards the router's allocations per request. Lower a number when an optimisation lands;
// never raise one without a reason in the commit.
func TestAllocations(t *testing.T) {
	skipUnderRace(t)

	r := fastrouter.New(fastrouter.WithNotFound(noop), fastrouter.WithMethodNotAllowed(noop))
	r.GET("/", noop)
	r.GET("/a/b/c/d/e", noop)
	r.GET("/users/{id}", func(ctx *fasthttp.RequestCtx) { _ = fastrouter.Param(ctx, "id") })
	r.GET("/repos/{owner}/{repo}/issues/{number:[0-9]+}", noop)
	r.GET("/admin/stats", noop)
	r.Use("/admin", passThrough, passThrough, passThrough)
	sub := fastrouter.New()
	sub.GET("/status", noop)
	r.Mount("/v1", sub.HandleFastHTTP)
	r.Mount("/t/{tenant}", sub.HandleFastHTTP)

	tests := []struct {
		name, method, uri string
		want              float64
	}{
		{"root", "GET", "/", 0},
		{"static, 5 segments", "GET", "/a/b/c/d/e", 0},
		{"static with 3 middlewares", "GET", "/admin/stats", 0},
		{"404", "GET", "/missing", 0},
		{"1 param", "GET", "/users/42", 0},
		{"3 params with a regexp", "GET", "/repos/vardius/gorouter/issues/7", 0},
		{"mount", "GET", "/v1/status", 0},
		{"mount with a param", "GET", "/t/acme/status", 0},
		// The Allow value.
		{"405", "POST", "/users/42", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := new(fasthttp.RequestCtx)
			ctx.Request.Header.SetMethod(tt.method)

			got := testing.AllocsPerRun(1000, func() {
				ctx.Request.SetRequestURI(tt.uri) // the mount rewrites the path, so start over each run
				r.HandleFastHTTP(ctx)
				ctx.ResetUserValues() // what fasthttp does after each request: params go back to the pool
			})
			if got != tt.want {
				t.Errorf("%s %s allocates %v times per request, want %v", tt.method, tt.uri, got, tt.want)
			}
		})
	}
}

// TestAllocationsWithManyParams checks that a route with more params than fit on the stack
// costs nothing, on its own requests or anyone else's.
func TestAllocationsWithManyParams(t *testing.T) {
	skipUnderRace(t)

	r := fastrouter.New(fastrouter.WithNotFound(noop))
	r.GET("/", noop)
	r.GET("/{a}/{b}/{c}/{d}/{e}/{f}/{g}/{h}/{i}", noop)

	for _, uri := range []string{"/", "/missing", "/1/2/3/4/5/6/7/8/9"} {
		ctx := new(fasthttp.RequestCtx)
		ctx.Request.Header.SetMethod("GET")
		ctx.Request.SetRequestURI(uri)

		got := testing.AllocsPerRun(1000, func() {
			r.HandleFastHTTP(ctx)
			ctx.ResetUserValues()
		})
		if got != 0 {
			t.Errorf("GET %s allocates %v times per request, want 0", uri, got)
		}
	}
}

func skipUnderRace(t *testing.T) {
	t.Helper()
	if race.Enabled {
		t.Skip("sync.Pool drops items at random under -race, so allocation counts mean nothing")
	}
}
