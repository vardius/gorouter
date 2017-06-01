package gorouter

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

//GoRouter benchmark tests functions

func benchmarkGoRouterStaticCall(t int, b *testing.B) {
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

func benchmarkGoRouterStaticParallel(t int, b *testing.B) {
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
		var buf bytes.Buffer
		w := httptest.NewRecorder()
		for pb.Next() {
			buf.Reset()
			s.ServeHTTP(w, req)
		}
	})
}

func benchmarkGoRouterWildcardCall(t int, b *testing.B) {
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

func benchmarkGoRouterWildcardParallel(t int, b *testing.B) {
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
		var buf bytes.Buffer
		w := httptest.NewRecorder()
		for pb.Next() {
			buf.Reset()
			s.ServeHTTP(w, req)
		}
	})
}

func benchmarkGoRouterRegexpCall(t int, b *testing.B) {
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

func benchmarkGoRouterRegexpParallel(t int, b *testing.B) {
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
		var buf bytes.Buffer
		w := httptest.NewRecorder()
		for pb.Next() {
			buf.Reset()
			s.ServeHTTP(w, req)
		}
	})
}

//HttpRouter benchmark tests functions

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

//GoRouter benchmark tests
func BenchmarkGoRouterStatic1(b *testing.B)  { benchmarkGoRouterStaticCall(1, b) }
func BenchmarkGoRouterStatic2(b *testing.B)  { benchmarkGoRouterStaticCall(2, b) }
func BenchmarkGoRouterStatic3(b *testing.B)  { benchmarkGoRouterStaticCall(3, b) }
func BenchmarkGoRouterStatic5(b *testing.B)  { benchmarkGoRouterStaticCall(5, b) }
func BenchmarkGoRouterStatic10(b *testing.B) { benchmarkGoRouterStaticCall(10, b) }
func BenchmarkGoRouterStatic20(b *testing.B) { benchmarkGoRouterStaticCall(20, b) }

func BenchmarkGoRouterWildcard1(b *testing.B)  { benchmarkGoRouterWildcardCall(1, b) }
func BenchmarkGoRouterWildcard2(b *testing.B)  { benchmarkGoRouterWildcardCall(2, b) }
func BenchmarkGoRouterWildcard3(b *testing.B)  { benchmarkGoRouterWildcardCall(3, b) }
func BenchmarkGoRouterWildcard5(b *testing.B)  { benchmarkGoRouterWildcardCall(5, b) }
func BenchmarkGoRouterWildcard10(b *testing.B) { benchmarkGoRouterWildcardCall(10, b) }
func BenchmarkGoRouterWildcard20(b *testing.B) { benchmarkGoRouterWildcardCall(20, b) }

func BenchmarkGoRouterRegexp1(b *testing.B)  { benchmarkGoRouterRegexpCall(1, b) }
func BenchmarkGoRouterRegexp2(b *testing.B)  { benchmarkGoRouterRegexpCall(2, b) }
func BenchmarkGoRouterRegexp3(b *testing.B)  { benchmarkGoRouterRegexpCall(3, b) }
func BenchmarkGoRouterRegexp5(b *testing.B)  { benchmarkGoRouterRegexpCall(5, b) }
func BenchmarkGoRouterRegexp10(b *testing.B) { benchmarkGoRouterRegexpCall(10, b) }
func BenchmarkGoRouterRegexp20(b *testing.B) { benchmarkGoRouterRegexpCall(20, b) }

func BenchmarkGoRouterStaticParallel1(b *testing.B)  { benchmarkGoRouterStaticParallel(1, b) }
func BenchmarkGoRouterStaticParallel2(b *testing.B)  { benchmarkGoRouterStaticParallel(2, b) }
func BenchmarkGoRouterStaticParallel3(b *testing.B)  { benchmarkGoRouterStaticParallel(3, b) }
func BenchmarkGoRouterStaticParallel5(b *testing.B)  { benchmarkGoRouterStaticParallel(5, b) }
func BenchmarkGoRouterStaticParallel10(b *testing.B) { benchmarkGoRouterStaticParallel(10, b) }
func BenchmarkGoRouterStaticParallel20(b *testing.B) { benchmarkGoRouterStaticParallel(20, b) }

func BenchmarkGoRouterWildcardParallel1(b *testing.B)  { benchmarkGoRouterWildcardParallel(1, b) }
func BenchmarkGoRouterWildcardParallel2(b *testing.B)  { benchmarkGoRouterWildcardParallel(2, b) }
func BenchmarkGoRouterWildcardParallel3(b *testing.B)  { benchmarkGoRouterWildcardParallel(3, b) }
func BenchmarkGoRouterWildcardParallel5(b *testing.B)  { benchmarkGoRouterWildcardParallel(5, b) }
func BenchmarkGoRouterWildcardParallel10(b *testing.B) { benchmarkGoRouterWildcardParallel(10, b) }
func BenchmarkGoRouterWildcardParallel20(b *testing.B) { benchmarkGoRouterWildcardParallel(20, b) }

func BenchmarkGoRouterRegexpParallel1(b *testing.B)  { benchmarkGoRouterRegexpParallel(1, b) }
func BenchmarkGoRouterRegexpParallel2(b *testing.B)  { benchmarkGoRouterRegexpParallel(2, b) }
func BenchmarkGoRouterRegexpParallel3(b *testing.B)  { benchmarkGoRouterRegexpParallel(3, b) }
func BenchmarkGoRouterRegexpParallel5(b *testing.B)  { benchmarkGoRouterRegexpParallel(5, b) }
func BenchmarkGoRouterRegexpParallel10(b *testing.B) { benchmarkGoRouterRegexpParallel(10, b) }
func BenchmarkGoRouterRegexpParallel20(b *testing.B) { benchmarkGoRouterRegexpParallel(20, b) }

//HttpRouter benchmark tests for comparison
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
