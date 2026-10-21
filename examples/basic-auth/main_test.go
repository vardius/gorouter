package main

import (
	"net/http/httptest"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	r := newRouter("admin", "s3cret")

	tests := []struct {
		name, target, user, password string
		code                         int
		body                         string
	}{
		{"public page needs nothing", "/", "", "", 200, "public page\n"},
		{"no credentials", "/admin", "", "", 401, "Unauthorized\n"},
		{"wrong password", "/admin", "admin", "guess", 401, "Unauthorized\n"},
		{"right credentials", "/admin", "admin", "s3cret", 200, "admin dashboard\n"},
		{"routes below are covered", "/admin/users/gopher", "", "", 401, "Unauthorized\n"},
		{"and let through", "/admin/users/gopher", "admin", "s3cret", 200, "managing gopher\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.target, nil)
			if tt.user != "" {
				req.SetBasicAuth(tt.user, tt.password)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.code || w.Body.String() != tt.body {
				t.Errorf("GET %s = %d %q, want %d %q", tt.target, w.Code, w.Body.String(), tt.code, tt.body)
			}
			if tt.code == 401 && w.Header().Get("WWW-Authenticate") == "" {
				t.Error("401 without a WWW-Authenticate header")
			}
		})
	}
}
