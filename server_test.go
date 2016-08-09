package goapi

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPOST(t *testing.T) {
	var (
		s *server     = New().(*server)
		h HandlerFunc = mockHandler
	)
	s.POST("/", h)

	assert.NotNil(t, s.routes[POST])
	rmap := s.Routes()
	assert.NotNil(t, rmap[POST])
}

func TestGET(t *testing.T) {
	var (
		s *server     = New().(*server)
		h HandlerFunc = mockHandler
	)
	s.GET("/", h)

	assert.NotNil(t, s.routes[GET])
	rmap := s.Routes()
	assert.NotNil(t, rmap[GET])
}

func TestPUT(t *testing.T) {
	var (
		s *server     = New().(*server)
		h HandlerFunc = mockHandler
	)
	s.PUT("/", h)

	assert.NotNil(t, s.routes[PUT])
	rmap := s.Routes()
	assert.NotNil(t, rmap[PUT])
}

func TestDELETE(t *testing.T) {
	var (
		s *server     = New().(*server)
		h HandlerFunc = mockHandler
	)
	s.DELETE("/", h)

	assert.NotNil(t, s.routes[DELETE])
	rmap := s.Routes()
	assert.NotNil(t, rmap[DELETE])
}

func TestPATCH(t *testing.T) {
	var (
		s *server     = New().(*server)
		h HandlerFunc = mockHandler
	)
	s.PATCH("/", h)

	assert.NotNil(t, s.routes[PATCH])
	rmap := s.Routes()
	assert.NotNil(t, rmap[PATCH])
}

func TestOPTIONS(t *testing.T) {
	var (
		s *server     = New().(*server)
		h HandlerFunc = mockHandler
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
		s  *server        = New().(*server)
		h  HandlerFunc    = mockHandler
		mh MiddlewareFunc = mockMiddleware
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
		s  *server        = New().(*server)
		h  HandlerFunc    = mockHandler
		mh MiddlewareFunc = mockMiddleware
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
		h http.HandlerFunc = mockHttpHandler
	)
	s.NotFound(h)
	assert.NotNil(t, s.notFound)
}

func TestNotAllowed(t *testing.T) {
	var (
		s *server          = New().(*server)
		h http.HandlerFunc = mockHttpHandler
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

func TestServer(t *testing.T) {
	s := New().(*server)

	serverd := false
	s.GET("/:param", func(w http.ResponseWriter, r *http.Request, ctx *Context) {
		serverd = true
		if assert.NotNil(t, ctx.Params["param"]) {
			assert.Equal(t, "x", ctx.Params["param"])
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

	s.GET("/:param", func(_ http.ResponseWriter, _ *http.Request, _ *Context) {
		panic("test panic recover")
	})

	w := new(mockResponseWriter)
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, true, paniced)
}

func TestMiddlewareError(t *testing.T) {
	var (
		s *server     = New().(*server)
		h HandlerFunc = mockHandler
	)
	s.GET("/x", h)
	s.Use("/x", 0, func(req *http.Request, _ *Context) Error {
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
		s *server     = New().(*server)
		h HandlerFunc = mockHandler
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
		s *server     = New().(*server)
		h HandlerFunc = mockHandler
	)

	s.GET("/x", h)

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("POST", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.code)
}

func TestGlobalMiddlewareError(t *testing.T) {
	var (
		s *server     = New().(*server)
		h HandlerFunc = mockHandler
	)

	s.GET("/x", h)
	s.Use("", 0, func(req *http.Request, _ *Context) Error {
		return statusError{http.StatusBadRequest, errors.New("Bad request")}
	})

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.code)
}
