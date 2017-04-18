package goserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func benchmarkStrictCall(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := New().(*server)
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request) {})

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

func BenchmarkStrict1(b *testing.B)   { benchmarkStrictCall(1, b) }
func BenchmarkStrict2(b *testing.B)   { benchmarkStrictCall(2, b) }
func BenchmarkStrict3(b *testing.B)   { benchmarkStrictCall(3, b) }
func BenchmarkStrict5(b *testing.B)   { benchmarkStrictCall(5, b) }
func BenchmarkStrict10(b *testing.B)  { benchmarkStrictCall(10, b) }
func BenchmarkStrict100(b *testing.B) { benchmarkStrictCall(100, b) }

func benchmarkStrictParallel(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := New().(*server)
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request) {})

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

func BenchmarkStrictParallel1(b *testing.B)   { benchmarkStrictParallel(1, b) }
func BenchmarkStrictParallel2(b *testing.B)   { benchmarkStrictParallel(2, b) }
func BenchmarkStrictParallel3(b *testing.B)   { benchmarkStrictParallel(3, b) }
func BenchmarkStrictParallel5(b *testing.B)   { benchmarkStrictParallel(5, b) }
func BenchmarkStrictParallel10(b *testing.B)  { benchmarkStrictParallel(10, b) }
func BenchmarkStrictParallel100(b *testing.B) { benchmarkStrictParallel(100, b) }

func benchmarkRegexpCall(t int, b *testing.B) {
	var path, rpath string
	part := "/:x:r([a-z]+)go"
	rpart := "/rxgo"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*server)
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request) {})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, rpath, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.ServeHTTP(w, req)
	}
}

func BenchmarkRegexp1(b *testing.B)   { benchmarkRegexpCall(1, b) }
func BenchmarkRegexp2(b *testing.B)   { benchmarkRegexpCall(2, b) }
func BenchmarkRegexp3(b *testing.B)   { benchmarkRegexpCall(3, b) }
func BenchmarkRegexp5(b *testing.B)   { benchmarkRegexpCall(5, b) }
func BenchmarkRegexp10(b *testing.B)  { benchmarkRegexpCall(10, b) }
func BenchmarkRegexp100(b *testing.B) { benchmarkRegexpCall(100, b) }

func benchmarkRegexpParallel(t int, b *testing.B) {
	var path, rpath string
	part := "/:x:r([a-z]+)go"
	rpart := "/rxgo"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*server)
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request) {})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, rpath, nil)
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

func BenchmarkRegexpParallel1(b *testing.B)   { benchmarkRegexpParallel(1, b) }
func BenchmarkRegexpParallel2(b *testing.B)   { benchmarkRegexpParallel(2, b) }
func BenchmarkRegexpParallel3(b *testing.B)   { benchmarkRegexpParallel(3, b) }
func BenchmarkRegexpParallel5(b *testing.B)   { benchmarkRegexpParallel(5, b) }
func BenchmarkRegexpParallel10(b *testing.B)  { benchmarkRegexpParallel(10, b) }
func BenchmarkRegexpParallel100(b *testing.B) { benchmarkRegexpParallel(100, b) }
