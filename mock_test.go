package goserver

import "net/http"

type (
	mockResponseWriter struct {
		header http.Header
		code   int
	}
)

func mockHandler(_ http.ResponseWriter, _ *http.Request) {}

func (m *mockResponseWriter) Header() (h http.Header) {
	return m.header
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(i int) {
	m.code = i
}
