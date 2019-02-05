package gorouter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func benchmarkStaticCall(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := New().(*router)
	s.GET(path, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

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

func benchmarkStaticParallel(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := New().(*router)
	s.GET(path, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	req, err := http.NewRequest(GET, path, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		w := httptest.NewRecorder()
		for pb.Next() {
			s.ServeHTTP(w, req)
		}
	})
}

func benchmarkWildcardCall(t int, b *testing.B) {
	var path, rpath string
	part := "/{x}"
	rpart := "/x"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*router)
	s.GET(path, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

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

func benchmarkWildcardParallel(t int, b *testing.B) {
	var path, rpath string
	part := "/{x}"
	rpart := "/x"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*router)
	s.GET(path, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	req, err := http.NewRequest(GET, rpath, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		w := httptest.NewRecorder()
		for pb.Next() {
			s.ServeHTTP(w, req)
		}
	})
}

func benchmarkRegexpCall(t int, b *testing.B) {
	var path, rpath string
	part := "/{x:r([a-z]+)go}"
	rpart := "/rxgo"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*router)
	s.GET(path, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

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

func benchmarkRegexpParallel(t int, b *testing.B) {
	var path, rpath string
	part := "/{x:r([a-z]+)go}"
	rpart := "/rxgo"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*router)
	s.GET(path, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	req, err := http.NewRequest(GET, rpath, nil)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		w := httptest.NewRecorder()
		for pb.Next() {
			s.ServeHTTP(w, req)
		}
	})
}

func BenchmarkStatic1(b *testing.B)  { benchmarkStaticCall(1, b) }
func BenchmarkStatic2(b *testing.B)  { benchmarkStaticCall(2, b) }
func BenchmarkStatic3(b *testing.B)  { benchmarkStaticCall(3, b) }
func BenchmarkStatic5(b *testing.B)  { benchmarkStaticCall(5, b) }
func BenchmarkStatic10(b *testing.B) { benchmarkStaticCall(10, b) }
func BenchmarkStatic20(b *testing.B) { benchmarkStaticCall(20, b) }

func BenchmarkWildcard1(b *testing.B)  { benchmarkWildcardCall(1, b) }
func BenchmarkWildcard2(b *testing.B)  { benchmarkWildcardCall(2, b) }
func BenchmarkWildcard3(b *testing.B)  { benchmarkWildcardCall(3, b) }
func BenchmarkWildcard5(b *testing.B)  { benchmarkWildcardCall(5, b) }
func BenchmarkWildcard10(b *testing.B) { benchmarkWildcardCall(10, b) }
func BenchmarkWildcard20(b *testing.B) { benchmarkWildcardCall(20, b) }

func BenchmarkRegexp1(b *testing.B)  { benchmarkRegexpCall(1, b) }
func BenchmarkRegexp2(b *testing.B)  { benchmarkRegexpCall(2, b) }
func BenchmarkRegexp3(b *testing.B)  { benchmarkRegexpCall(3, b) }
func BenchmarkRegexp5(b *testing.B)  { benchmarkRegexpCall(5, b) }
func BenchmarkRegexp10(b *testing.B) { benchmarkRegexpCall(10, b) }
func BenchmarkRegexp20(b *testing.B) { benchmarkRegexpCall(20, b) }

func BenchmarkStaticParallel1(b *testing.B)  { benchmarkStaticParallel(1, b) }
func BenchmarkStaticParallel2(b *testing.B)  { benchmarkStaticParallel(2, b) }
func BenchmarkStaticParallel3(b *testing.B)  { benchmarkStaticParallel(3, b) }
func BenchmarkStaticParallel5(b *testing.B)  { benchmarkStaticParallel(5, b) }
func BenchmarkStaticParallel10(b *testing.B) { benchmarkStaticParallel(10, b) }
func BenchmarkStaticParallel20(b *testing.B) { benchmarkStaticParallel(20, b) }

func BenchmarkWildcardParallel1(b *testing.B)  { benchmarkWildcardParallel(1, b) }
func BenchmarkWildcardParallel2(b *testing.B)  { benchmarkWildcardParallel(2, b) }
func BenchmarkWildcardParallel3(b *testing.B)  { benchmarkWildcardParallel(3, b) }
func BenchmarkWildcardParallel5(b *testing.B)  { benchmarkWildcardParallel(5, b) }
func BenchmarkWildcardParallel10(b *testing.B) { benchmarkWildcardParallel(10, b) }
func BenchmarkWildcardParallel20(b *testing.B) { benchmarkWildcardParallel(20, b) }

func BenchmarkRegexpParallel1(b *testing.B)  { benchmarkRegexpParallel(1, b) }
func BenchmarkRegexpParallel2(b *testing.B)  { benchmarkRegexpParallel(2, b) }
func BenchmarkRegexpParallel3(b *testing.B)  { benchmarkRegexpParallel(3, b) }
func BenchmarkRegexpParallel5(b *testing.B)  { benchmarkRegexpParallel(5, b) }
func BenchmarkRegexpParallel10(b *testing.B) { benchmarkRegexpParallel(10, b) }
func BenchmarkRegexpParallel20(b *testing.B) { benchmarkRegexpParallel(20, b) }
