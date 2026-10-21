package main

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestHTTP2OverTLS(t *testing.T) {
	srv := httptest.NewUnstartedServer(newRouter())
	srv.EnableHTTP2 = true
	srv.StartTLS()
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Proto != "HTTP/2.0" || string(body) != "hello over HTTP/2.0\n" {
		t.Errorf("GET / = %s %q, want HTTP/2.0", resp.Proto, body)
	}
}
