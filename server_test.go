package goapi

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPOST(t *testing.T) {
	s := New()
	s.POST("/", func(w http.ResponseWriter, req *http.Request) {})

	assert.NotNil(t, s.routes[POST])
}

func TestGET(t *testing.T) {
	s := New()
	s.GET("/", func(w http.ResponseWriter, req *http.Request) {})

	assert.NotNil(t, s.routes[GET])
}

func TestPUT(t *testing.T) {
	s := New()
	s.PUT("/", func(w http.ResponseWriter, req *http.Request) {})

	assert.NotNil(t, s.routes[PUT])
}

func TestDELETE(t *testing.T) {
	s := New()
	s.DELETE("/", func(w http.ResponseWriter, req *http.Request) {})

	assert.NotNil(t, s.routes[DELETE])
}

func TestPATCH(t *testing.T) {
	s := New()
	s.PATCH("/", func(w http.ResponseWriter, req *http.Request) {})

	assert.NotNil(t, s.routes[PATCH])
}

func TestOPTIONS(t *testing.T) {
	s := New()
	s.OPTIONS("/", func(w http.ResponseWriter, req *http.Request) {})

	assert.NotNil(t, s.routes[OPTIONS])
}

func TestUseGlobal(t *testing.T) {
	s := New()
	s.Use("", 0, func(req *http.Request) Error { return nil })

	assert.NotEqual(t, 0, len(s.middleware))
}

func TestUseRoot(t *testing.T) {
	s := New()
	s.OPTIONS("/", func(w http.ResponseWriter, req *http.Request) {})
	s.PATCH("/", func(w http.ResponseWriter, req *http.Request) {})
	s.DELETE("/", func(w http.ResponseWriter, req *http.Request) {})
	s.PUT("/", func(w http.ResponseWriter, req *http.Request) {})
	s.GET("/", func(w http.ResponseWriter, req *http.Request) {})

	s.Use("/", 0, func(req *http.Request) Error { return nil })

	for _, r := range s.routes {
		assert.NotEqual(t, 0, len(r.middleware))
	}
}

func TestUseNodes(t *testing.T) {
	s := New()
	s.OPTIONS("/x", func(w http.ResponseWriter, req *http.Request) {})
	s.PATCH("/x", func(w http.ResponseWriter, req *http.Request) {})
	s.DELETE("/x", func(w http.ResponseWriter, req *http.Request) {})
	s.PUT("/x", func(w http.ResponseWriter, req *http.Request) {})
	s.GET("/x", func(w http.ResponseWriter, req *http.Request) {})

	s.Use("/x", 0, func(req *http.Request) Error { return nil })

	for _, root := range s.routes {
		assert.Equal(t, 0, len(root.middleware))
		node := root.nodes["x"]
		if assert.NotNil(t, node) {
			assert.NotEqual(t, 0, len(node.middleware))
		}
	}
}

type mockResponseWriter struct {
	header http.Header
	code   int
}

func (m *mockResponseWriter) Header() (h http.Header) {
	return m.header
}
func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}
func (m *mockResponseWriter) WriteHeader(i int) {
	m.code = i
}

func TestServer(t *testing.T) {
	s := New()

	serverd := false
	s.GET("/:param", func(w http.ResponseWriter, r *http.Request) {
		serverd = true
	})

	w := new(mockResponseWriter)
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, true, serverd)
}

type (
	statusError struct {
		Code int
		Err  error
	}
)

func (se statusError) Error() string {
	return se.Err.Error()
}

func (se statusError) Status() int {
	return se.Code
}

func TestMiddlewareError(t *testing.T) {
	s := New()

	s.GET("/x", func(w http.ResponseWriter, req *http.Request) {})
	s.Use("/x", 0, func(req *http.Request) Error {
		return statusError{http.StatusBadRequest, errors.New("Bad request")}
	})

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.code)
}

func TestGlobalMiddlewareError(t *testing.T) {
	s := New()

	s.GET("/x", func(w http.ResponseWriter, req *http.Request) {})
	s.Use("", 0, func(req *http.Request) Error {
		return statusError{http.StatusBadRequest, errors.New("Bad request")}
	})

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("GET", "/x", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.code)
}

func TestNotFound(t *testing.T) {
	s := New()

	s.GET("/x", func(w http.ResponseWriter, req *http.Request) {})

	w := new(mockResponseWriter)
	w.header = http.Header{}
	req, _ := http.NewRequest("GET", "/y", nil)
	s.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.code)
}
