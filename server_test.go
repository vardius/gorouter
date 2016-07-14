package goapi

import (
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
