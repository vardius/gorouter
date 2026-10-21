package main

import (
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	r := newRouter()

	tests := []struct {
		method, target string
		code           int
		body           string
	}{
		{"GET", "/", 200, "welcome\n"},
		{"GET", "/hello/gopher", 200, "hello, gopher!\n"},
		{"GET", "/users/42", 200, "user #42\n"},
		{"GET", "/users/gopher", 404, "404 page not found\n"},
		{"POST", "/", 405, "Method Not Allowed\n"},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(tt.method, tt.target, nil))
		if w.Code != tt.code || w.Body.String() != tt.body {
			t.Errorf("%s %s = %d %q, want %d %q", tt.method, tt.target, w.Code, w.Body.String(), tt.code, tt.body)
		}
	}
}
