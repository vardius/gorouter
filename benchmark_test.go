package gorouter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/valyala/fasthttp"
)

func benchmarkStatic(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := New()
	s.GET(path, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.ServeHTTP(w, req)
		}
	})
}

func benchmarkWildcard(t int, b *testing.B) {
	var path, rpath string
	part := "/{x}"
	rpart := "/x"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New()
	s.GET(path, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, rpath, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.ServeHTTP(w, req)
		}
	})
}

func benchmarkRegexp(t int, b *testing.B) {
	var path, rpath string
	part := "/{x:r([a-z]+)go}"
	rpart := "/rxgo"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New()
	s.GET(path, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, rpath, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.ServeHTTP(w, req)
		}
	})
}

func BenchmarkStatic1(b *testing.B)  { benchmarkStatic(1, b) }
func BenchmarkStatic2(b *testing.B)  { benchmarkStatic(2, b) }
func BenchmarkStatic3(b *testing.B)  { benchmarkStatic(3, b) }
func BenchmarkStatic5(b *testing.B)  { benchmarkStatic(5, b) }
func BenchmarkStatic10(b *testing.B) { benchmarkStatic(10, b) }
func BenchmarkStatic20(b *testing.B) { benchmarkStatic(20, b) }

func BenchmarkWildcard1(b *testing.B)  { benchmarkWildcard(1, b) }
func BenchmarkWildcard2(b *testing.B)  { benchmarkWildcard(2, b) }
func BenchmarkWildcard3(b *testing.B)  { benchmarkWildcard(3, b) }
func BenchmarkWildcard5(b *testing.B)  { benchmarkWildcard(5, b) }
func BenchmarkWildcard10(b *testing.B) { benchmarkWildcard(10, b) }
func BenchmarkWildcard20(b *testing.B) { benchmarkWildcard(20, b) }

func BenchmarkRegexp1(b *testing.B)  { benchmarkRegexp(1, b) }
func BenchmarkRegexp2(b *testing.B)  { benchmarkRegexp(2, b) }
func BenchmarkRegexp3(b *testing.B)  { benchmarkRegexp(3, b) }
func BenchmarkRegexp5(b *testing.B)  { benchmarkRegexp(5, b) }
func BenchmarkRegexp10(b *testing.B) { benchmarkRegexp(10, b) }
func BenchmarkRegexp20(b *testing.B) { benchmarkRegexp(20, b) }

func benchmarkFastHTTPStatic(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := NewFastHTTPRouter()
	s.GET(path, func(_ *fasthttp.RequestCtx) {})

	ctx := buildFastHTTPRequestContext(http.MethodGet, path)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.HandleFastHTTP(ctx)
		}
	})
}

func benchmarkFastHTTPWildcard(t int, b *testing.B) {
	var path, rpath string
	part := "/{x}"
	rpart := "/x"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := NewFastHTTPRouter()
	s.GET(path, func(_ *fasthttp.RequestCtx) {})

	ctx := buildFastHTTPRequestContext(http.MethodGet, rpath)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.HandleFastHTTP(ctx)
		}
	})
}

func benchmarkFastHTTPRegexp(t int, b *testing.B) {
	var path, rpath string
	part := "/{x:r([a-z]+)go}"
	rpart := "/rxgo"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := NewFastHTTPRouter()
	s.GET(path, func(_ *fasthttp.RequestCtx) {})

	ctx := buildFastHTTPRequestContext(http.MethodGet, rpath)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.HandleFastHTTP(ctx)
		}
	})
}

func BenchmarkFastHTTPStatic1(b *testing.B)  { benchmarkFastHTTPStatic(1, b) }
func BenchmarkFastHTTPStatic2(b *testing.B)  { benchmarkFastHTTPStatic(2, b) }
func BenchmarkFastHTTPStatic3(b *testing.B)  { benchmarkFastHTTPStatic(3, b) }
func BenchmarkFastHTTPStatic5(b *testing.B)  { benchmarkFastHTTPStatic(5, b) }
func BenchmarkFastHTTPStatic10(b *testing.B) { benchmarkFastHTTPStatic(10, b) }
func BenchmarkFastHTTPStatic20(b *testing.B) { benchmarkFastHTTPStatic(20, b) }

func BenchmarkFastHTTPWildcard1(b *testing.B)  { benchmarkFastHTTPWildcard(1, b) }
func BenchmarkFastHTTPWildcard2(b *testing.B)  { benchmarkFastHTTPWildcard(2, b) }
func BenchmarkFastHTTPWildcard3(b *testing.B)  { benchmarkFastHTTPWildcard(3, b) }
func BenchmarkFastHTTPWildcard5(b *testing.B)  { benchmarkFastHTTPWildcard(5, b) }
func BenchmarkFastHTTPWildcard10(b *testing.B) { benchmarkFastHTTPWildcard(10, b) }
func BenchmarkFastHTTPWildcard20(b *testing.B) { benchmarkFastHTTPWildcard(20, b) }

func BenchmarkFastHTTPRegexp1(b *testing.B)  { benchmarkFastHTTPRegexp(1, b) }
func BenchmarkFastHTTPRegexp2(b *testing.B)  { benchmarkFastHTTPRegexp(2, b) }
func BenchmarkFastHTTPRegexp3(b *testing.B)  { benchmarkFastHTTPRegexp(3, b) }
func BenchmarkFastHTTPRegexp5(b *testing.B)  { benchmarkFastHTTPRegexp(5, b) }
func BenchmarkFastHTTPRegexp10(b *testing.B) { benchmarkFastHTTPRegexp(10, b) }
func BenchmarkFastHTTPRegexp20(b *testing.B) { benchmarkFastHTTPRegexp(20, b) }
