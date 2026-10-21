package main

import (
	"net/http/httptest"
	"testing"
)

func TestMount(t *testing.T) {
	r := newRouter()

	tests := []struct {
		target string
		code   int
		body   string
	}{
		{"/api/v1/status", 200, "v1 ok\n"},
		{"/api/v1/users/7", 200, "v1 user 7\n"},
		{"/api/v1/users/gopher", 404, "404 page not found\n"},
		{"/orgs/vardius/repos/gorouter", 200, "vardius/gorouter (path /repos/gorouter)\n"},
		{"/api/v2/status", 404, "404 page not found\n"},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", tt.target, nil))
		if w.Code != tt.code || w.Body.String() != tt.body {
			t.Errorf("GET %s = %d %q, want %d %q", tt.target, w.Code, w.Body.String(), tt.code, tt.body)
		}
	}
}
