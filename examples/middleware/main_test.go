package main

import (
	"bytes"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddleware(t *testing.T) {
	var logs bytes.Buffer
	r := newRouter(slog.New(slog.NewTextHandler(&logs, nil)))

	tests := []struct {
		method, target, contentType, body string
		code                              int
		cacheControl, response            string
	}{
		{"GET", "/", "", "", 200, "", "welcome\n"},
		{"GET", "/api/status", "", "", 200, "no-store", `{"status":"ok"}`},
		{"POST", "/api/echo", "application/json; charset=utf-8", `{"a":1}`, 200, "no-store", `{"a":1}`},
		{"POST", "/api/echo", "text/plain", "a", 415, "no-store", "send JSON\n"},
		{"GET", "/panic", "", "", 500, "", "Internal Server Error\n"},
	}
	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body))
		req.Header.Set("Content-Type", tt.contentType)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != tt.code || w.Body.String() != tt.response || w.Header().Get("Cache-Control") != tt.cacheControl {
			t.Errorf("%s %s = %d %q Cache-Control %q, want %d %q %q", tt.method, tt.target,
				w.Code, w.Body.String(), w.Header().Get("Cache-Control"), tt.code, tt.response, tt.cacheControl)
		}
	}

	for _, want := range []string{"path=/api/status", "path=/panic", "level=ERROR msg=panic"} {
		if !strings.Contains(logs.String(), want) {
			t.Errorf("logs don't contain %q:\n%s", want, logs.String())
		}
	}
}
