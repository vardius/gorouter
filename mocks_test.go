package gorouter

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/valyala/fasthttp"
)

type testLogger struct {
	t *testing.T
}

func (t testLogger) Printf(format string, args ...interface{}) {
	t.t.Logf(format, args...)
}

type mockHandler struct {
	served bool
}

func (mh *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mh.served = true
}

func (mh *mockHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	mh.served = true
}

type mockFileSystem struct {
	opened bool
}

func (mfs *mockFileSystem) Open(name string) (http.File, error) {
	mfs.opened = true
	return nil, errors.New("")
}

func mockMiddleware(body string) MiddlewareFunc {
	fn := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(body))
			h.ServeHTTP(w, r)
		})
	}

	return fn
}

func mockServeHTTP(h http.Handler, method, path string) error {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return err
	}

	h.ServeHTTP(w, req)

	return nil
}

func mockFastHTTPMiddleware(body string) FastHTTPMiddlewareFunc {
	fn := func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			fmt.Fprintf(ctx, body)

			h(ctx)
		}
	}

	return fn
}

func mockHandleFastHTTP(h fasthttp.RequestHandler, method, path string) error {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.URI().SetPath(path)

	h(ctx)

	return nil
}

func checkIfHasRootRoute(t *testing.T, r interface{}, method string) {
	switch v := r.(type) {
	case *router:
	case *fastHTTPRouter:
		if rootRoute := v.routes.GetByID(method); rootRoute == nil {
			t.Error("Route not found")
		}
	default:
		t.Error("Unsupported type")
	}
}
