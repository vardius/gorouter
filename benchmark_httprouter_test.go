package goserver

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

func benchmarkHttpRouterStaticCall(t int, b *testing.B) {
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

func benchmarkHttpRouterStaticParallel(t int, b *testing.B) {
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

func benchmarkHttpRouterWildcardCall(t int, b *testing.B) {
	var path, rpath string
	part := "/:x"
	rpart := "/x"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := httprouter.New()
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {})

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

func benchmarkHttpRouterWildcardParallel(t int, b *testing.B) {
	var path, rpath string
	part := "/:x"
	rpart := "/x"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := httprouter.New()
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {})

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

func BenchmarkHttpRouterStatic1(b *testing.B)  { benchmarkHttpRouterStaticCall(1, b) }
func BenchmarkHttpRouterStatic2(b *testing.B)  { benchmarkHttpRouterStaticCall(2, b) }
func BenchmarkHttpRouterStatic3(b *testing.B)  { benchmarkHttpRouterStaticCall(3, b) }
func BenchmarkHttpRouterStatic5(b *testing.B)  { benchmarkHttpRouterStaticCall(5, b) }
func BenchmarkHttpRouterStatic10(b *testing.B) { benchmarkHttpRouterStaticCall(10, b) }
func BenchmarkHttpRouterStatic20(b *testing.B) { benchmarkHttpRouterStaticCall(20, b) }

func BenchmarkHttpRouterWildcard1(b *testing.B)  { benchmarkHttpRouterWildcardCall(1, b) }
func BenchmarkHttpRouterWildcard2(b *testing.B)  { benchmarkHttpRouterWildcardCall(2, b) }
func BenchmarkHttpRouterWildcard3(b *testing.B)  { benchmarkHttpRouterWildcardCall(3, b) }
func BenchmarkHttpRouterWildcard5(b *testing.B)  { benchmarkHttpRouterWildcardCall(5, b) }
func BenchmarkHttpRouterWildcard10(b *testing.B) { benchmarkHttpRouterWildcardCall(10, b) }
func BenchmarkHttpRouterWildcard20(b *testing.B) { benchmarkHttpRouterWildcardCall(20, b) }

func BenchmarkHttpRouterStaticParallel1(b *testing.B)  { benchmarkHttpRouterStaticParallel(1, b) }
func BenchmarkHttpRouterStaticParallel2(b *testing.B)  { benchmarkHttpRouterStaticParallel(2, b) }
func BenchmarkHttpRouterStaticParallel3(b *testing.B)  { benchmarkHttpRouterStaticParallel(3, b) }
func BenchmarkHttpRouterStaticParallel5(b *testing.B)  { benchmarkHttpRouterStaticParallel(5, b) }
func BenchmarkHttpRouterStaticParallel10(b *testing.B) { benchmarkHttpRouterStaticParallel(10, b) }
func BenchmarkHttpRouterStaticParallel20(b *testing.B) { benchmarkHttpRouterStaticParallel(20, b) }

func BenchmarkHttpRouterWildcardParallel1(b *testing.B)  { benchmarkHttpRouterWildcardParallel(1, b) }
func BenchmarkHttpRouterWildcardParallel2(b *testing.B)  { benchmarkHttpRouterWildcardParallel(2, b) }
func BenchmarkHttpRouterWildcardParallel3(b *testing.B)  { benchmarkHttpRouterWildcardParallel(3, b) }
func BenchmarkHttpRouterWildcardParallel5(b *testing.B)  { benchmarkHttpRouterWildcardParallel(5, b) }
func BenchmarkHttpRouterWildcardParallel10(b *testing.B) { benchmarkHttpRouterWildcardParallel(10, b) }
func BenchmarkHttpRouterWildcardParallel20(b *testing.B) { benchmarkHttpRouterWildcardParallel(20, b) }
