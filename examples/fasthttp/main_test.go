package main

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestFastHTTP(t *testing.T) {
	r := newRouter()

	tests := []struct {
		method, uri string
		code        int
		body        string
	}{
		{"GET", "/", 200, "welcome\n"},
		{"GET", "/hello/gopher", 200, "hello, gopher!\n"},
		{"GET", "/api/v1/users/42", 200, "user #42\n"},
		{"GET", "/api/v1/users/gopher", 404, "Not Found"},
		{"POST", "/", 405, "Method Not Allowed"},
	}
	for _, tt := range tests {
		var ctx fasthttp.RequestCtx
		ctx.Request.Header.SetMethod(tt.method)
		ctx.Request.SetRequestURI(tt.uri)
		r.HandleFastHTTP(&ctx)

		if got := string(ctx.Response.Body()); ctx.Response.StatusCode() != tt.code || got != tt.body {
			t.Errorf("%s %s = %d %q, want %d %q", tt.method, tt.uri, ctx.Response.StatusCode(), got, tt.code, tt.body)
		}
		if got := string(ctx.Response.Header.Peek("X-Powered-By")); got != "gorouter" {
			t.Errorf("%s %s X-Powered-By = %q", tt.method, tt.uri, got)
		}
	}
}
