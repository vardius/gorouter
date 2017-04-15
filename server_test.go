package goserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockHandler(_ http.ResponseWriter, _ *http.Request) {}

func TestInterface(t *testing.T) {
	var _ http.Handler = New()
}

func TestPOST(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.POST("/", h)

	if s.routes[post] == nil {
		t.Error("Route not found")
	}

	rmap := s.Routes()
	if rmap[post] == nil {
		t.Error("Route not found")
	}
}

func TestGET(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.GET("/", h)

	if s.routes[get] == nil {
		t.Error("Route not found")
	}

	rmap := s.Routes()
	if rmap[get] == nil {
		t.Error("Route not found")
	}
}

func TestPUT(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.PUT("/", h)

	if s.routes[put] == nil {
		t.Error("Route not found")
	}

	rmap := s.Routes()
	if rmap[put] == nil {
		t.Error("Route not found")
	}
}

func TestDELETE(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.DELETE("/", h)

	if s.routes[delete] == nil {
		t.Error("Route not found")
	}

	rmap := s.Routes()
	if rmap[delete] == nil {
		t.Error("Route not found")
	}
}

func TestPATCH(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.PATCH("/", h)

	if s.routes[patch] == nil {
		t.Error("Route not found")
	}

	rmap := s.Routes()
	if rmap[patch] == nil {
		t.Error("Route not found")
	}
}

func TestOPTIONS(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.OPTIONS("/", h)

	if s.routes[options] == nil {
		t.Error("Route not found")
	}

	rmap := s.Routes()
	if rmap[options] == nil {
		t.Error("Route not found")
	}
}

func TestNotFound(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.NotFound(h)

	if s.notFound == nil {
		t.Error("NotFound handler error")
	}
}

func TestNotAllowed(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.NotAllowed(h)

	if s.notAllowed == nil {
		t.Error("NotAllowed handler error")
	}
}

func TestServerServeFiles(t *testing.T) {
	s := New().(*server)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("File exist")
	}
}

func TestServerFiles(t *testing.T) {
	var s *server = New().(*server)

	s.ServeFiles("static", true)

	if s.fileServer == nil {
		t.Error("File serve handler error")
	}
}

func TestServerFilesError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()

	var s *server = New().(*server)
	s.ServeFiles("", true)
}

func TestServer(t *testing.T) {
	s := New().(*server)

	serverd := false
	s.GET("/:param", func(w http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := ParamsFromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params["param"] != "x" {
			t.Error("Wrong params value")
		}
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestServerNotFound(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.GET("/x", h)
	s.GET("/x", h)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/y", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Error("NotFound doesnt work")
	}
}

func TestServerNotAllowed(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.GET("/x", h)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Error("NotAllowed doesnt work")
	}
}

func TestServerOptions(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.GET("/x", h)
	s.POST("/x", h)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("OPTIONS", "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Header().Get("Allow") == "" {
		t.Error("Options doesnt work")
	}
}

func TestServerNilMiddleware(t *testing.T) {
	s := New(nil).(*server)

	s.GET("/:param", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("test"))
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Body.String() != "test" {
		t.Error("Nil middleware works")
	}
}

func TestServerPanicMiddleware(t *testing.T) {
	paniced := false
	panicMiddleware := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rcv := recover(); rcv != nil {
					paniced = true
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	s := New(panicMiddleware).(*server)

	s.GET("/:param", func(_ http.ResponseWriter, _ *http.Request) {
		panic("test panic recover")
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if paniced != true {
		t.Error("Panic has not been handled")
	}
}
