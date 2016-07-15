package goapi

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPOST(t *testing.T) {
	var (
		s *Server     = New()
		h HandlerFunc = mockHandler
	)
	s.POST("/", h)

	assert.NotNil(t, s.routes[POST])
}

func TestGET(t *testing.T) {
	var (
		s *Server     = New()
		h HandlerFunc = mockHandler
	)
	s.GET("/", h)

	assert.NotNil(t, s.routes[GET])
}

func TestPUT(t *testing.T) {
	var (
		s *Server     = New()
		h HandlerFunc = mockHandler
	)
	s.PUT("/", h)

	assert.NotNil(t, s.routes[PUT])
}

func TestDELETE(t *testing.T) {
	var (
		s *Server     = New()
		h HandlerFunc = mockHandler
	)
	s.DELETE("/", h)

	assert.NotNil(t, s.routes[DELETE])
}

func TestPATCH(t *testing.T) {
	var (
		s *Server     = New()
		h HandlerFunc = mockHandler
	)
	s.PATCH("/", h)

	assert.NotNil(t, s.routes[PATCH])
}

func TestOPTIONS(t *testing.T) {
	var (
		s *Server     = New()
		h HandlerFunc = mockHandler
	)
	s.OPTIONS("/", h)

	assert.NotNil(t, s.routes[OPTIONS])
}

func TestUseGlobal(t *testing.T) {
	var (
		s  *Server        = New()
		mh MiddlewareFunc = mockMiddleware
	)
	s.Use("", 0, mh)

	assert.NotEqual(t, 0, len(s.middleware))
}

func TestUseRoot(t *testing.T) {
	var (
		s  *Server        = New()
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
		s  *Server        = New()
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

func TestServer(t *testing.T) {
	s := New()

	serverd := false
	s.GET("/:param", func(w http.ResponseWriter, r *http.Request, params Params) {
		serverd = true
		if assert.NotNil(t, params["param"]) {
			assert.Equal(t, "x", params["param"])
		}
	})

	w := new(mockResponseWriter)
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, true, serverd)
}

func TestMiddlewareError(t *testing.T) {
	var (
		s *Server     = New()
		h HandlerFunc = mockHandler
	)
	s.GET("/x", h)
	s.Use("/x", 0, func(req *http.Request, _ Params) Error {
		return statusError{http.StatusBadRequest, errors.New("Bad request")}
	})

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.code)
}

func TestGlobalMiddlewareError(t *testing.T) {
	var (
		s *Server     = New()
		h HandlerFunc = mockHandler
	)

	s.GET("/x", h)
	s.Use("", 0, func(req *http.Request, _ Params) Error {
		return statusError{http.StatusBadRequest, errors.New("Bad request")}
	})

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.code)
}

func TestNotFound(t *testing.T) {
	var (
		s *Server     = New()
		h HandlerFunc = mockHandler
	)

	s.GET("/x", h)

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("GET", "/y", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.code)
}
