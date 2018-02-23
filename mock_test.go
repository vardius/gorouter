package gorouter

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func mockServeHTTP(h http.Handler, method, path string) error {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return err
	}

	h.ServeHTTP(w, req)

	return nil
}

func mockHandler(_ http.ResponseWriter, _ *http.Request) {}

func mockMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("middleware"))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func mockMiddlewareWithBody(body string) MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(body))
			h.ServeHTTP(w, r)
		})
	}
}

func equal(t *testing.T, expected, actual interface{}) bool {
	if !areEqual(expected, actual) {
		t.Errorf("Asserts are not equal. Expected: %v, Actual: %v", expected, actual)

		return false
	}

	return true
}

func notEqual(t *testing.T, expected, actual interface{}) bool {
	if areEqual(expected, actual) {
		t.Errorf("Asserts are equal. Expected: %v, Actual: %v", expected, actual)

		return false
	}

	return true
}

func areEqual(expected, actual interface{}) bool {
	if expected == nil {
		return isNil(actual)
	}

	if actual == nil {
		return isNil(expected)
	}

	return reflect.DeepEqual(expected, actual)
}

func isNil(value interface{}) bool {
	defer func() { recover() }()
	return value == nil || reflect.ValueOf(value).IsNil()
}

func checkIfHasRootRoute(t *testing.T, router *router, method string) {
	if rootRoute := router.routes.byID(method); rootRoute == nil {
		t.Error("Route not found")
	}
}

func testBasicMethod(t *testing.T, router *router, h func(pattern string, handler http.Handler), method string) {
	serverd := false
	h("/x/y", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	checkIfHasRootRoute(t, router, method)

	err := mockServeHTTP(router, method, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}
