package goserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func benchmarkGoserverStaticCall(t int, b *testing.B) {
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

func benchmarkGoserverStaticParallel(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := New().(*server)
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request) {})

	req, err := http.NewRequest(GET, path, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		w := httptest.NewRecorder()
		for pb.Next() {
			buf.Reset()
			s.ServeHTTP(w, req)
		}
	})
}

func benchmarkGoserverWildcardCall(t int, b *testing.B) {
	var path, rpath string
	part := "/:x"
	rpart := "/x"
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

func benchmarkGoserverWildcardParallel(t int, b *testing.B) {
	var path, rpath string
	part := "/:x"
	rpart := "/x"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*server)
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request) {})

	req, err := http.NewRequest(GET, rpath, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		w := httptest.NewRecorder()
		for pb.Next() {
			buf.Reset()
			s.ServeHTTP(w, req)
		}
	})
}

func benchmarkGoserverRegexpCall(t int, b *testing.B) {
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

func benchmarkGoserverRegexpParallel(t int, b *testing.B) {
	var path, rpath string
	part := "/:x:r([a-z]+)go"
	rpart := "/rxgo"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*server)
	s.GET(path, func(_ http.ResponseWriter, _ *http.Request) {})

	req, err := http.NewRequest(GET, rpath, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		w := httptest.NewRecorder()
		for pb.Next() {
			buf.Reset()
			s.ServeHTTP(w, req)
		}
	})
}

func BenchmarkGoserverStatic1(b *testing.B)  { benchmarkGoserverStaticCall(1, b) }
func BenchmarkGoserverStatic2(b *testing.B)  { benchmarkGoserverStaticCall(2, b) }
func BenchmarkGoserverStatic3(b *testing.B)  { benchmarkGoserverStaticCall(3, b) }
func BenchmarkGoserverStatic5(b *testing.B)  { benchmarkGoserverStaticCall(5, b) }
func BenchmarkGoserverStatic10(b *testing.B) { benchmarkGoserverStaticCall(10, b) }
func BenchmarkGoserverStatic20(b *testing.B) { benchmarkGoserverStaticCall(20, b) }

// func BenchmarkGoserverWildcard1(b *testing.B)  { benchmarkGoserverWildcardCall(1, b) }
// func BenchmarkGoserverWildcard2(b *testing.B)  { benchmarkGoserverWildcardCall(2, b) }
// func BenchmarkGoserverWildcard3(b *testing.B)  { benchmarkGoserverWildcardCall(3, b) }
// func BenchmarkGoserverWildcard5(b *testing.B)  { benchmarkGoserverWildcardCall(5, b) }
// func BenchmarkGoserverWildcard10(b *testing.B) { benchmarkGoserverWildcardCall(10, b) }
// func BenchmarkGoserverWildcard20(b *testing.B) { benchmarkGoserverWildcardCall(20, b) }

// func BenchmarkGoserverRegexp1(b *testing.B)  { benchmarkGoserverRegexpCall(1, b) }
// func BenchmarkGoserverRegexp2(b *testing.B)  { benchmarkGoserverRegexpCall(2, b) }
// func BenchmarkGoserverRegexp3(b *testing.B)  { benchmarkGoserverRegexpCall(3, b) }
// func BenchmarkGoserverRegexp5(b *testing.B)  { benchmarkGoserverRegexpCall(5, b) }
// func BenchmarkGoserverRegexp10(b *testing.B) { benchmarkGoserverRegexpCall(10, b) }
// func BenchmarkGoserverRegexp20(b *testing.B) { benchmarkGoserverRegexpCall(20, b) }

func BenchmarkGoserverStaticParallel1(b *testing.B)  { benchmarkGoserverStaticParallel(1, b) }
func BenchmarkGoserverStaticParallel2(b *testing.B)  { benchmarkGoserverStaticParallel(2, b) }
func BenchmarkGoserverStaticParallel3(b *testing.B)  { benchmarkGoserverStaticParallel(3, b) }
func BenchmarkGoserverStaticParallel5(b *testing.B)  { benchmarkGoserverStaticParallel(5, b) }
func BenchmarkGoserverStaticParallel10(b *testing.B) { benchmarkGoserverStaticParallel(10, b) }
func BenchmarkGoserverStaticParallel20(b *testing.B) { benchmarkGoserverStaticParallel(20, b) }

// func BenchmarkGoserverWildcardParallel1(b *testing.B)  { benchmarkGoserverWildcardParallel(1, b) }
// func BenchmarkGoserverWildcardParallel2(b *testing.B)  { benchmarkGoserverWildcardParallel(2, b) }
// func BenchmarkGoserverWildcardParallel3(b *testing.B)  { benchmarkGoserverWildcardParallel(3, b) }
// func BenchmarkGoserverWildcardParallel5(b *testing.B)  { benchmarkGoserverWildcardParallel(5, b) }
// func BenchmarkGoserverWildcardParallel10(b *testing.B) { benchmarkGoserverWildcardParallel(10, b) }
// func BenchmarkGoserverWildcardParallel20(b *testing.B) { benchmarkGoserverWildcardParallel(20, b) }

// func BenchmarkGoserverRegexpParallel1(b *testing.B)  { benchmarkGoserverRegexpParallel(1, b) }
// func BenchmarkGoserverRegexpParallel2(b *testing.B)  { benchmarkGoserverRegexpParallel(2, b) }
// func BenchmarkGoserverRegexpParallel3(b *testing.B)  { benchmarkGoserverRegexpParallel(3, b) }
// func BenchmarkGoserverRegexpParallel5(b *testing.B)  { benchmarkGoserverRegexpParallel(5, b) }
// func BenchmarkGoserverRegexpParallel10(b *testing.B) { benchmarkGoserverRegexpParallel(10, b) }
// func BenchmarkGoserverRegexpParallel20(b *testing.B) { benchmarkGoserverRegexpParallel(20, b) }
