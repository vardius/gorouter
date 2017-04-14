package goserver

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	var _ http.Handler = New()
}

func TestPOST(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.POST("/", h)

	assert.NotNil(t, s.routes[post])
	rmap := s.Routes()
	assert.NotNil(t, rmap[post])
}

func TestGET(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.GET("/", h)

	assert.NotNil(t, s.routes[get])
	rmap := s.Routes()
	assert.NotNil(t, rmap[get])
}

func TestPUT(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.PUT("/", h)

	assert.NotNil(t, s.routes[put])
	rmap := s.Routes()
	assert.NotNil(t, rmap[put])
}

func TestDELETE(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.DELETE("/", h)

	assert.NotNil(t, s.routes[delete])
	rmap := s.Routes()
	assert.NotNil(t, rmap[delete])
}

func TestPATCH(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.PATCH("/", h)

	assert.NotNil(t, s.routes[patch])
	rmap := s.Routes()
	assert.NotNil(t, rmap[patch])
}

func TestOPTIONS(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.OPTIONS("/", h)

	assert.NotNil(t, s.routes[options])
	rmap := s.Routes()
	assert.NotNil(t, rmap[options])
}

func TestNotFound(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.NotFound(h)
	assert.NotNil(t, s.notFound)
}

func TestNotAllowed(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.NotAllowed(h)
	assert.NotNil(t, s.notAllowed)
}

func TestServerFiles(t *testing.T) {
	var s *server = New().(*server)
	s.ServeFiles("static", true)
	assert.NotNil(t, s.fileServer)
}

func TestServerFilesError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
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
		params, _ := ParamsFromContext(r.Context())
		if assert.NotNil(t, params["param"]) {
			assert.Equal(t, "x", params["param"])
		}
	})

	w := new(mockResponseWriter)
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, true, serverd)
}

func TestServerPanic(t *testing.T) {
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

	w := new(mockResponseWriter)
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, true, paniced)
}

func TestServerNotFound(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.GET("/x", h)

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("GET", "/y", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.code)
}

func TestServerNotAllowed(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.GET("/x", h)

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("POST", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.code)
}

func TestServerOptions(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.GET("/x", h)
	s.POST("/x", h)

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("OPTIONS", "/x", nil)
	s.ServeHTTP(w, req)
	assert.NotEmpty(t, w.Header().Get("Allow"))
}
