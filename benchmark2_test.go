package goserver

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

func jbenchmarkStrictCall(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := httprouter.New()
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, path, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ServeHTTP(w, req)
	}
}

func BenchmarkJStrict1(b *testing.B)   { jbenchmarkStrictCall(1, b) }
func BenchmarkJStrict2(b *testing.B)   { jbenchmarkStrictCall(2, b) }
func BenchmarkJStrict3(b *testing.B)   { jbenchmarkStrictCall(3, b) }
func BenchmarkJStrict5(b *testing.B)   { jbenchmarkStrictCall(5, b) }
func BenchmarkJStrict10(b *testing.B)  { jbenchmarkStrictCall(10, b) }
func BenchmarkJStrict100(b *testing.B) { jbenchmarkStrictCall(100, b) }

func jbenchmarkStrictParallel(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := httprouter.New()
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, path, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			s.ServeHTTP(w, req)
		}
	})
}

func BenchmarkJStrictParallel1(b *testing.B)   { jbenchmarkStrictParallel(1, b) }
func BenchmarkJStrictParallel2(b *testing.B)   { jbenchmarkStrictParallel(2, b) }
func BenchmarkJStrictParallel3(b *testing.B)   { jbenchmarkStrictParallel(3, b) }
func BenchmarkJStrictParallel5(b *testing.B)   { jbenchmarkStrictParallel(5, b) }
func BenchmarkJStrictParallel10(b *testing.B)  { jbenchmarkStrictParallel(10, b) }
func BenchmarkJStrictParallel100(b *testing.B) { jbenchmarkStrictParallel(100, b) }
