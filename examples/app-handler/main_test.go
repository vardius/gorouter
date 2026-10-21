package main

import (
	"net/http/httptest"
	"testing"
)

func TestAppHandler(t *testing.T) {
	r := newRouter()

	tests := []struct {
		target string
		code   int
		body   string
	}{
		{"/users/1", 200, `{"id":"1","name":"Gopher"}` + "\n"},
		{"/users/2", 404, "no such user\n"},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", tt.target, nil))
		if w.Code != tt.code || w.Body.String() != tt.body {
			t.Errorf("GET %s = %d %q, want %d %q", tt.target, w.Code, w.Body.String(), tt.code, tt.body)
		}
	}
}
