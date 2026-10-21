package main

import (
	"net/http/httptest"
	"testing"
)

func TestMultidomain(t *testing.T) {
	h := newHandler()

	tests := []struct {
		host, target string
		code         int
		body         string
	}{
		{"example.com", "/", 200, "the website\n"},
		{"www.example.com:8080", "/", 200, "the website\n"},
		{"api.example.com", "/users/7", 200, "api user 7\n"},
		{"api.example.com", "/", 404, "404 page not found\n"},
		{"other.example.com", "/", 421, "unknown host\n"},
	}
	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.target, nil)
		req.Host = tt.host
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != tt.code || w.Body.String() != tt.body {
			t.Errorf("GET %s%s = %d %q, want %d %q", tt.host, tt.target, w.Code, w.Body.String(), tt.code, tt.body)
		}
	}
}
