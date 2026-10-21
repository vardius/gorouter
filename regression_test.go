package gorouter

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/vardius/gorouter/v4/context"
)

func writeHandler(body string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(body))
	})
}

func serve(h http.Handler, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	return w
}

func serveFastHTTP(r FastHTTPRouter, method, path string) *fasthttp.RequestCtx {
	ctx := buildFastHTTPRequestContext(method, path)
	r.HandleFastHTTP(ctx)
	return ctx
}

// Concurrent requests to sibling routes must each run their own middleware chain.
func TestConcurrentMiddlewareChains(t *testing.T) {
	t.Parallel()

	router := New()
	router.GET("/a", writeHandler("A"))
	router.GET("/b", writeHandler("B"))
	// Four root middlewares leave spare capacity in the node's backing array.
	router.USE(http.MethodGet, "/", mockMiddleware("1"), mockMiddleware("2"), mockMiddleware("3"))
	router.USE(http.MethodGet, "/", mockMiddleware("4"))
	router.USE(http.MethodGet, "/a", mockMiddleware("a"))
	router.USE(http.MethodGet, "/b", mockMiddleware("b"))

	var wg sync.WaitGroup
	for g := range 8 {
		wg.Go(func() {
			for i := range 500 {
				path, want := "/a", "1234aA"
				if (i+g)%2 == 0 {
					path, want = "/b", "1234bB"
				}
				if got := serve(router, http.MethodGet, path).Body.String(); got != want {
					t.Errorf("GET %s = %q, want %q", path, got, want)
					return
				}
			}
		})
	}
	wg.Wait()
}

func TestCompileKeepsRoutes(t *testing.T) {
	t.Parallel()

	router := New()
	router.GET("/a", writeHandler("A"))
	router.GET("/a/b", writeHandler("AB"))
	router.GET("/x/y", writeHandler("XY"))
	router.Compile()

	for path, want := range map[string]string{"/a": "A", "/a/b": "AB", "/x/y": "XY"} {
		if got := serve(router, http.MethodGet, path).Body.String(); got != want {
			t.Errorf("GET %s after Compile = %q, want %q", path, got, want)
		}
	}
}

func TestMiddlewareRequiresSegmentBoundary(t *testing.T) {
	t.Parallel()

	router := New()
	router.GET("/abc", writeHandler("H"))
	router.USE(http.MethodGet, "/a", mockMiddleware("m"))

	if got := serve(router, http.MethodGet, "/abc").Body.String(); got != "H" {
		t.Errorf("GET /abc = %q, want %q: middleware of /a leaked", got, "H")
	}
}

func TestRegexpMatchesWholeSegment(t *testing.T) {
	t.Parallel()

	router := New()
	router.GET("/{id:[0-9]+}", writeHandler("H"))

	tests := []struct {
		path string
		code int
	}{
		{"/123", http.StatusOK},
		{"/abc123", http.StatusNotFound},
		{"/123abc", http.StatusNotFound},
	}
	for _, tt := range tests {
		if got := serve(router, http.MethodGet, tt.path).Code; got != tt.code {
			t.Errorf("GET %s = %d, want %d", tt.path, got, tt.code)
		}
	}
}

func TestRegexpMayContainColon(t *testing.T) {
	t.Parallel()

	router := New()
	router.GET("/{t:(?:a|b)}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, _ := context.Parameters(r.Context())
		_, _ = w.Write([]byte(params.Value("t")))
	}))

	if w := serve(router, http.MethodGet, "/a"); w.Body.String() != "a" {
		t.Errorf("GET /a = %d %q, want 200 %q", w.Code, w.Body.String(), "a")
	}
	if w := serve(router, http.MethodGet, "/zzz"); w.Code != http.StatusNotFound {
		t.Errorf("GET /zzz = %d, want 404", w.Code)
	}
}

func TestUnclosedWildcardPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Error("registering /{id should panic")
		}
	}()

	New().GET("/{id", writeHandler(""))
}

func TestRootPathMethodNotAllowed(t *testing.T) {
	t.Parallel()

	router := New()
	router.POST("/", writeHandler("P"))

	w := serve(router, http.MethodGet, "/")
	if w.Code != http.StatusMethodNotAllowed || w.Header().Get("Allow") != "POST, OPTIONS" {
		t.Errorf("GET / = %d Allow %q, want 405 Allow %q", w.Code, w.Header().Get("Allow"), "POST, OPTIONS")
	}

	fastRouter := NewFastHTTPRouter()
	fastRouter.POST("/", func(_ *fasthttp.RequestCtx) {})

	ctx := serveFastHTTP(fastRouter, http.MethodGet, "/")
	if ctx.Response.StatusCode() != fasthttp.StatusMethodNotAllowed {
		t.Errorf("fasthttp GET / = %d, want 405", ctx.Response.StatusCode())
	}
}

func TestMountOverMiddlewareNode(t *testing.T) {
	t.Parallel()

	sub := New()
	sub.GET("/x", writeHandler("SUB"))

	router := New()
	router.USE(http.MethodGet, "/a", mockMiddleware("m"))
	router.Mount("/a", sub)

	if got := serve(router, http.MethodGet, "/a/x").Body.String(); got != "mSUB" {
		t.Errorf("GET /a/x = %q, want %q", got, "mSUB")
	}
}

func TestMountOverRoutesPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Error("mounting over /a/b should panic")
		}
	}()

	router := New()
	router.GET("/a/b", writeHandler("AB"))
	router.Mount("/a", New())
}
