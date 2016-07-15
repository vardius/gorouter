package goapi

import "net/http"

type (
	statusError struct {
		code int
		err  error
	}
	mockResponseWriter struct {
		header http.Header
		code   int
	}
)

func mockHandler(_ http.ResponseWriter, _ *http.Request, _ Params) {}
func mockMiddleware(_ *http.Request, _ Params) Error {
	return nil
}
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
func (se statusError) Error() string {
	return se.err.Error()
}
func (se statusError) Status() int {
	return se.code
}
