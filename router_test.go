package goapi

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRootRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	r.addRoute(paths, mockHandler)

	assert.Equal(t, "/", r.Path())
	assert.Equal(t, "", r.Regexp())
	assert.Equal(t, true, r.IsRoot())
	assert.Nil(t, r.Parent())
	assert.NotNil(t, r.Nodes())
	assert.Equal(t, true, r.IsEndPoint())

	var node *route
	node, _ = r.getRoute([]string{""})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"x"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"", "x"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"x", "y"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{})
	assert.NotNil(t, node)
}

func TestGetStrictRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/x", "/"), "/")
	r.addRoute(paths, mockHandler)

	var node *route
	node, _ = r.getRoute([]string{})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{""})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"", "x"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"x", "y"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"x"})
	assert.NotNil(t, node)
}

func TestGetParamRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/:x", "/"), "/")
	r.addRoute(paths, mockHandler)

	var node *route
	var params parameters
	node, _ = r.getRoute([]string{})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{""})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"", "x"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"x", "y"})
	assert.Nil(t, node)
	node, params = r.getRoute([]string{"x"})
	assert.NotNil(t, node)
	if assert.NotNil(t, params["x"]) {
		assert.Equal(t, "x", params["x"])
	}
	node, params = r.getRoute([]string{"y"})
	if assert.NotNil(t, params["x"]) {
		assert.Equal(t, "y", params["x"])
	}
	assert.NotNil(t, node)
}

func TestGetRegexRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/:x:r([a-z]+)go", "/"), "/")
	r.addRoute(paths, mockHandler)

	var node *route
	node, _ = r.getRoute([]string{})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{""})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"", "x"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"x", "y"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"x"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"y"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"rego", "y"})
	assert.Nil(t, node)
	node, _ = r.getRoute([]string{"rego"})
	assert.NotNil(t, node)
}

func TestAddRootRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	r.addRoute(paths, mockHandler)

	assert.Equal(t, true, r.isEndPoint)
	assert.Equal(t, tree{}, r.nodes)
}

func TestAddEmptyRootRoute(t *testing.T) {
	r := newRoute(nil, "")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	r.addRoute(paths, mockHandler)

	assert.Equal(t, true, r.isEndPoint)
	assert.Equal(t, tree{}, r.nodes)
}

func TestAddEmptyRootRouteTwo(t *testing.T) {
	r := newRoute(nil, "")
	paths := strings.Split(strings.Trim("", "/"), "/")
	r.addRoute(paths, mockHandler)

	assert.Equal(t, true, r.isEndPoint)
	assert.Equal(t, tree{}, r.nodes)
}

func TestAddRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/example", "/"), "/")
	r.addRoute(paths, mockHandler)

	assert.Equal(t, false, r.isEndPoint)
	if assert.NotNil(t, r.nodes["example"]) {
		assert.Equal(t, true, r.nodes["example"].isEndPoint)
		assert.Equal(t, "example", r.nodes["example"].path)
		assert.Nil(t, r.nodes["example"].regexp)
	}
}

func TestAddParamRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/:example", "/"), "/")
	r.addRoute(paths, mockHandler)

	assert.Equal(t, false, r.isEndPoint)
	if assert.NotNil(t, r.nodes[":example"]) {
		assert.Equal(t, true, r.nodes[":example"].isEndPoint)
		assert.Equal(t, ":example", r.nodes[":example"].path)
		assert.Nil(t, r.nodes[":example"].regexp)
	}
}

func TestAddRegexpRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/:example:r([a-z]+)go", "/"), "/")
	r.addRoute(paths, mockHandler)

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
