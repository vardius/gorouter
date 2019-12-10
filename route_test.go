package gorouter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vardius/gorouter/v4/context"
	"github.com/vardius/gorouter/v4/middleware"
)

func TestRouter(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("4"))
	})

	buildMiddlewareFunc := func(body string) middleware.MiddlewareFunc {
		fn := func(h interface{}) interface{} {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(body))
				h.(http.Handler).ServeHTTP(w, r)
			})
		}

		return fn
	}

	m1 := buildMiddlewareFunc("1")
	m2 := buildMiddlewareFunc("2")
	m3 := buildMiddlewareFunc("3")

	r := newRoute(handler)
	r.AppendMiddleware(middleware.New(m1, m2, m3))

	h := r.Handler()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	h.(http.Handler).ServeHTTP(w, req)

	if w.Body.String() != "1234" {
		t.Errorf("The router doesn't work correctly. Expected 1234, Actual: %s", w.Body.String())
	}
}

func TestParams(t *testing.T) {
	param := context.Param{Key: "key", Value: "value"}
	params := context.Params{param}

	if params.Value("key") != "value" {
		t.Error("Invalid params value")
	}
}

func TestInvalidParams(t *testing.T) {
	param := context.Param{Key: "key", Value: "value"}
	params := context.Params{param}

	if params.Value("invalid_key") != "" {
		t.Error("Invalid params value")
	}
}

func TestNilHandler(t *testing.T) {
	panicked := false
	defer func() {
		if rcv := recover(); rcv != nil {
			panicked = true
		}
	}()

	r := newRoute(nil)
	if h := r.Handler(); h != nil {
		t.Error("Handler should be equal nil")
	}

	if panicked != true {
		t.Error("Router should panic if handler is nil")
	}
}
