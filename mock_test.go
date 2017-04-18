package goserver

import (
	"net/http"
	"reflect"
	"testing"
)

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
