package goserver

import (
	"errors"
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

	assert.NotNil(t, s.routes[POST])
	rmap := s.Routes()
	assert.NotNil(t, rmap[POST])
}

func TestGET(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.GET("/", h)

	assert.NotNil(t, s.routes[GET])
	rmap := s.Routes()
	assert.NotNil(t, rmap[GET])
}

func TestPUT(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.PUT("/", h)

	assert.NotNil(t, s.routes[PUT])
	rmap := s.Routes()
	assert.NotNil(t, rmap[PUT])
}

func TestDELETE(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.DELETE("/", h)

	assert.NotNil(t, s.routes[DELETE])
	rmap := s.Routes()
	assert.NotNil(t, rmap[DELETE])
}

func TestPATCH(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.PATCH("/", h)

	assert.NotNil(t, s.routes[PATCH])
	rmap := s.Routes()
	assert.NotNil(t, rmap[PATCH])
}

func TestOPTIONS(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.OPTIONS("/", h)

	assert.NotNil(t, s.routes[OPTIONS])
	rmap := s.Routes()
	assert.NotNil(t, rmap[OPTIONS])
}

func TestUseGlobal(t *testing.T) {
	var (
		s  *server        = New().(*server)
		mh MiddlewareFunc = mockMiddleware
	)
	s.Use("", 0, mh)

	assert.NotEqual(t, 0, len(s.middleware))
}

func TestUseRoot(t *testing.T) {
	var (
		s  *server          = New().(*server)
		h  http.HandlerFunc = mockHandler
		mh MiddlewareFunc   = mockMiddleware
	)
	s.OPTIONS("/", h)
	s.PATCH("/", h)
	s.DELETE("/", h)
	s.PUT("/", h)
	s.GET("/", h)

	s.Use("/", 0, mh)

	for _, r := range s.routes {
		assert.NotEqual(t, 0, len(r.middleware))
	}
}

func TestUseNodes(t *testing.T) {
	var (
		s  *server          = New().(*server)
		h  http.HandlerFunc = mockHandler
		mh MiddlewareFunc   = mockMiddleware
	)
	s.OPTIONS("/x", h)
	s.PATCH("/x", h)
	s.DELETE("/x", h)
	s.PUT("/x", h)
	s.GET("/x", h)

	s.Use("/x", 0, mh)

	for _, root := range s.routes {
		assert.Equal(t, 0, len(root.middleware))
		node := root.nodes["x"]
		if assert.NotNil(t, node) {
			assert.NotEqual(t, 0, len(node.middleware))
		}
	}
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

func TestOnPanic(t *testing.T) {
	var (
		s *server          = New().(*server)
		h PanicHandlerFunc = mockPanicHandler
	)
	s.OnPanic(h)
	assert.NotNil(t, s.onPanic)
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
		params, _ := ParametersFromContext(r.Context())
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
	s := New().(*server)

	paniced := false
	s.OnPanic(func(_ http.ResponseWriter, _ *http.Request, _ interface{}) {
		paniced = true
	})

	s.GET("/:param", func(_ http.ResponseWriter, _ *http.Request) {
		panic("test panic recover")
	})

	w := new(mockResponseWriter)
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, true, paniced)
}

func TestMiddlewareError(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)
	s.GET("/x", h)
	s.Use("/x", 0, func(_ http.ResponseWriter, req *http.Request) Error {
		return statusError{http.StatusBadRequest, errors.New("Bad request")}
	})

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.code)
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

func TestGlobalMiddlewareError(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHandler
	)

	s.GET("/x", h)
	s.Use("", 0, func(_ http.ResponseWriter, req *http.Request) Error {
		return statusError{http.StatusBadRequest, errors.New("Bad request")}
	})

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.code)
}
