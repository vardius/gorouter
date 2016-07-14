package goapi

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRootRoute(t *testing.T) {
	r := newRoute("/")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	r.addRoute(paths, func(w http.ResponseWriter, req *http.Request) {})

	r.getRoute(paths)

	assert.Nil(t, r.getRoute([]string{""}))
	assert.Nil(t, r.getRoute([]string{"x"}))
	assert.Nil(t, r.getRoute([]string{"", "x"}))
	assert.Nil(t, r.getRoute([]string{"x", "y"}))
	assert.NotNil(t, r.getRoute([]string{}))
}

func TestGetStrictRoute(t *testing.T) {
	r := newRoute("/")
	paths := strings.Split(strings.Trim("/x", "/"), "/")
	r.addRoute(paths, func(w http.ResponseWriter, req *http.Request) {})

	r.getRoute(paths)

	assert.Nil(t, r.getRoute([]string{}))
	assert.Nil(t, r.getRoute([]string{""}))
	assert.Nil(t, r.getRoute([]string{"", "x"}))
	assert.Nil(t, r.getRoute([]string{"x", "y"}))
	assert.NotNil(t, r.getRoute([]string{"x"}))
}

func TestGetParamRoute(t *testing.T) {
	r := newRoute("/")
	paths := strings.Split(strings.Trim("/:x", "/"), "/")
	r.addRoute(paths, func(w http.ResponseWriter, req *http.Request) {})

	r.getRoute(paths)

	assert.Nil(t, r.getRoute([]string{}))
	assert.Nil(t, r.getRoute([]string{""}))
	assert.Nil(t, r.getRoute([]string{"", "x"}))
	assert.Nil(t, r.getRoute([]string{"x", "y"}))
	assert.NotNil(t, r.getRoute([]string{"x"}))
	assert.NotNil(t, r.getRoute([]string{"y"}))
}

func TestGetRegexRoute(t *testing.T) {
	r := newRoute("/")
	paths := strings.Split(strings.Trim("/:x:r([a-z]+)go", "/"), "/")
	r.addRoute(paths, func(w http.ResponseWriter, req *http.Request) {})

	r.getRoute(paths)

	assert.Nil(t, r.getRoute([]string{}))
	assert.Nil(t, r.getRoute([]string{""}))
	assert.Nil(t, r.getRoute([]string{"", "x"}))
	assert.Nil(t, r.getRoute([]string{"x", "y"}))
	assert.Nil(t, r.getRoute([]string{"x"}))
	assert.Nil(t, r.getRoute([]string{"y"}))
	assert.Nil(t, r.getRoute([]string{"rego", "y"}))
	assert.NotNil(t, r.getRoute([]string{"rego"}))
}

func TestAddRootRoute(t *testing.T) {
	r := newRoute("/")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	r.addRoute(paths, func(w http.ResponseWriter, req *http.Request) {})

	assert.Equal(t, true, r.isEndPoint)
	assert.Equal(t, tree{}, r.nodes)
}

func TestAddEmptyRootRoute(t *testing.T) {
	r := newRoute("")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	r.addRoute(paths, func(w http.ResponseWriter, req *http.Request) {})

	assert.Equal(t, true, r.isEndPoint)
	assert.Equal(t, tree{}, r.nodes)
}

func TestAddEmptyRootRouteTwo(t *testing.T) {
	r := newRoute("")
	paths := strings.Split(strings.Trim("", "/"), "/")
	r.addRoute(paths, func(w http.ResponseWriter, req *http.Request) {})

	assert.Equal(t, true, r.isEndPoint)
	assert.Equal(t, tree{}, r.nodes)
}

func TestAddRoute(t *testing.T) {
	r := newRoute("/")
	paths := strings.Split(strings.Trim("/example", "/"), "/")
	r.addRoute(paths, func(w http.ResponseWriter, req *http.Request) {})

	assert.Equal(t, false, r.isEndPoint)
	if assert.NotNil(t, r.nodes["example"]) {
		assert.Equal(t, true, r.nodes["example"].isEndPoint)
		assert.Equal(t, "example", r.nodes["example"].path)
		assert.Nil(t, r.nodes["example"].regexp)
	}
}

func TestAddParamRoute(t *testing.T) {
	r := newRoute("/")
	paths := strings.Split(strings.Trim("/:example", "/"), "/")
	r.addRoute(paths, func(w http.ResponseWriter, req *http.Request) {})

	assert.Equal(t, false, r.isEndPoint)
	if assert.NotNil(t, r.nodes[":example"]) {
		assert.Equal(t, true, r.nodes[":example"].isEndPoint)
		assert.Equal(t, ":example", r.nodes[":example"].path)
		assert.Nil(t, r.nodes[":example"].regexp)
	}
}

func TestAddRegexpRoute(t *testing.T) {
	r := newRoute("/")
	paths := strings.Split(strings.Trim("/:example:r([a-z]+)go", "/"), "/")
	r.addRoute(paths, func(w http.ResponseWriter, req *http.Request) {})

	assert.Equal(t, false, r.isEndPoint)

	route := r.nodes[":example:r([a-z]+)go"]
	if assert.NotNil(t, route) {
		assert.Equal(t, true, route.isEndPoint)
		assert.Equal(t, ":example:r([a-z]+)go", route.path)
		if assert.NotNil(t, route.regexp) {
			assert.Equal(t, true, route.regexp.MatchString("rego"))
		}
	}
}
