package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStaticFiles(t *testing.T) {
	r, err := newRouter()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		target      string
		code        int
		contentType string
		contains    string
	}{
		{"/", 200, "text/html; charset=utf-8", "Served by gorouter"},
		{"/static/style.css", 200, "text/css; charset=utf-8", "font-family"},
		{"/static/missing.css", 404, "text/plain; charset=utf-8", "404 page not found"},
		{"/api/time", 200, "text/plain; charset=utf-8", "T"},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("GET", tt.target, nil))

		if w.Code != tt.code || w.Header().Get("Content-Type") != tt.contentType || !strings.Contains(w.Body.String(), tt.contains) {
			t.Errorf("GET %s = %d %q %q, want %d %q containing %q", tt.target,
				w.Code, w.Header().Get("Content-Type"), w.Body.String(), tt.code, tt.contentType, tt.contains)
		}
	}
}
