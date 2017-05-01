package goserver

import (
	"bytes"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Goserver benchmark tests functions

func benchmarkGoserverStaticCall(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := New().(*server)
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

func benchmarkGoserverStaticParallel(t int, b *testing.B) {
	var path string
	part := "/x"
	for i := 0; i < t; i++ {
		path += part
	}

	s := New().(*server)
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

func benchmarkGoserverWildcardCall(t int, b *testing.B) {
	var path, rpath string
	part := "/{x}"
	rpart := "/x"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*server)
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

func benchmarkGoserverWildcardParallel(t int, b *testing.B) {
	var path, rpath string
	part := "/{x}"
	rpart := "/x"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*server)
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

func benchmarkGoserverRegexpCall(t int, b *testing.B) {
	var path, rpath string
	part := "/{x:r([a-z]+)go}"
	rpart := "/rxgo"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*server)
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

func benchmarkGoserverRegexpParallel(t int, b *testing.B) {
	var path, rpath string
	part := "/{x:r([a-z]+)go}"
	rpart := "/rxgo"
	for i := 0; i < t; i++ {
		path += part
		rpath += rpart
	}

	s := New().(*server)
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

//Goserver benchmark tests
func BenchmarkGoserverStatic1(b *testing.B)  { benchmarkGoserverStaticCall(1, b) }
func BenchmarkGoserverStatic2(b *testing.B)  { benchmarkGoserverStaticCall(2, b) }
func BenchmarkGoserverStatic3(b *testing.B)  { benchmarkGoserverStaticCall(3, b) }
func BenchmarkGoserverStatic5(b *testing.B)  { benchmarkGoserverStaticCall(5, b) }
func BenchmarkGoserverStatic10(b *testing.B) { benchmarkGoserverStaticCall(10, b) }
func BenchmarkGoserverStatic20(b *testing.B) { benchmarkGoserverStaticCall(20, b) }

func BenchmarkGoserverWildcard1(b *testing.B)  { benchmarkGoserverWildcardCall(1, b) }
func BenchmarkGoserverWildcard2(b *testing.B)  { benchmarkGoserverWildcardCall(2, b) }
func BenchmarkGoserverWildcard3(b *testing.B)  { benchmarkGoserverWildcardCall(3, b) }
func BenchmarkGoserverWildcard5(b *testing.B)  { benchmarkGoserverWildcardCall(5, b) }
func BenchmarkGoserverWildcard10(b *testing.B) { benchmarkGoserverWildcardCall(10, b) }
func BenchmarkGoserverWildcard20(b *testing.B) { benchmarkGoserverWildcardCall(20, b) }

func BenchmarkGoserverRegexp1(b *testing.B)  { benchmarkGoserverRegexpCall(1, b) }
func BenchmarkGoserverRegexp2(b *testing.B)  { benchmarkGoserverRegexpCall(2, b) }
func BenchmarkGoserverRegexp3(b *testing.B)  { benchmarkGoserverRegexpCall(3, b) }
func BenchmarkGoserverRegexp5(b *testing.B)  { benchmarkGoserverRegexpCall(5, b) }
func BenchmarkGoserverRegexp10(b *testing.B) { benchmarkGoserverRegexpCall(10, b) }
func BenchmarkGoserverRegexp20(b *testing.B) { benchmarkGoserverRegexpCall(20, b) }

func BenchmarkGoserverStaticParallel1(b *testing.B)  { benchmarkGoserverStaticParallel(1, b) }
func BenchmarkGoserverStaticParallel2(b *testing.B)  { benchmarkGoserverStaticParallel(2, b) }
func BenchmarkGoserverStaticParallel3(b *testing.B)  { benchmarkGoserverStaticParallel(3, b) }
func BenchmarkGoserverStaticParallel5(b *testing.B)  { benchmarkGoserverStaticParallel(5, b) }
func BenchmarkGoserverStaticParallel10(b *testing.B) { benchmarkGoserverStaticParallel(10, b) }
func BenchmarkGoserverStaticParallel20(b *testing.B) { benchmarkGoserverStaticParallel(20, b) }

func BenchmarkGoserverWildcardParallel1(b *testing.B)  { benchmarkGoserverWildcardParallel(1, b) }
func BenchmarkGoserverWildcardParallel2(b *testing.B)  { benchmarkGoserverWildcardParallel(2, b) }
func BenchmarkGoserverWildcardParallel3(b *testing.B)  { benchmarkGoserverWildcardParallel(3, b) }
func BenchmarkGoserverWildcardParallel5(b *testing.B)  { benchmarkGoserverWildcardParallel(5, b) }
func BenchmarkGoserverWildcardParallel10(b *testing.B) { benchmarkGoserverWildcardParallel(10, b) }
func BenchmarkGoserverWildcardParallel20(b *testing.B) { benchmarkGoserverWildcardParallel(20, b) }

func BenchmarkGoserverRegexpParallel1(b *testing.B)  { benchmarkGoserverRegexpParallel(1, b) }
func BenchmarkGoserverRegexpParallel2(b *testing.B)  { benchmarkGoserverRegexpParallel(2, b) }
func BenchmarkGoserverRegexpParallel3(b *testing.B)  { benchmarkGoserverRegexpParallel(3, b) }
func BenchmarkGoserverRegexpParallel5(b *testing.B)  { benchmarkGoserverRegexpParallel(5, b) }
func BenchmarkGoserverRegexpParallel10(b *testing.B) { benchmarkGoserverRegexpParallel(10, b) }
func BenchmarkGoserverRegexpParallel20(b *testing.B) { benchmarkGoserverRegexpParallel(20, b) }

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
