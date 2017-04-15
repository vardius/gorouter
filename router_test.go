package goserver

import (
	"reflect"
	"strings"
	"testing"
)

func equal(t *testing.T, expected, actual interface{}) bool {
	if !areEqual(expected, actual) {
		t.Errorf("Asserts are not equal. Expected: %v, Actual: %v", expected, actual)

		return false
	}

	return true
}

func notEqual(t *testing.T, expected, actual interface{}) bool {
	if areEqual(expected, actual) {
		t.Errorf("Asserts are equal. Expected: %v, Actual: %v", expected, actual)

		return false
	}

	return true
}

func areEqual(expected, actual interface{}) bool {
	if expected == nil {
		return isNil(actual)
	}

	if actual == nil {
		return isNil(expected)
	}

	return reflect.DeepEqual(expected, actual)
}

func isNil(value interface{}) bool {
	defer func() { recover() }()
	return value == nil || reflect.ValueOf(value).IsNil()
}

func TestGetRootRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	r.addRoute(paths, mockHandler)

	equal(t, "/", r.Path())
	equal(t, "", r.Regexp())
	equal(t, true, r.IsRoot())
	equal(t, nil, r.Parent())
	notEqual(t, nil, r.Nodes())
	equal(t, true, r.IsEndPoint())

	var node *route
	node, _ = r.getRoute([]string{""})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"x"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"", "x"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"x", "y"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{})
	notEqual(t, nil, node)
}

func TestGetStrictRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/x", "/"), "/")
	r.addRoute(paths, mockHandler)

	var node *route
	node, _ = r.getRoute([]string{})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{""})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"", "x"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"x", "y"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"x"})
	notEqual(t, nil, node)
}

func TestGetParamRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/:x", "/"), "/")
	r.addRoute(paths, mockHandler)

	var node *route
	var params Params
	node, _ = r.getRoute([]string{})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{""})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"", "x"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"x", "y"})
	equal(t, nil, node)
	node, params = r.getRoute([]string{"x"})
	notEqual(t, nil, node)
	if notEqual(t, nil, params["x"]) {
		equal(t, "x", params["x"])
	}
	node, params = r.getRoute([]string{"y"})
	if notEqual(t, nil, params["x"]) {
		equal(t, "y", params["x"])
	}
	notEqual(t, nil, node)
}

func TestGetRegexRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/:x:r([a-z]+)go", "/"), "/")
	r.addRoute(paths, mockHandler)

	var node *route
	node, _ = r.getRoute([]string{})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{""})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"", "x"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"x", "y"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"x"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"y"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"rego", "y"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"rego"})
	notEqual(t, nil, node)
}

func TestGetNestedRegexRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/:x:r([a-z]+)go/:y:r([a-z]+)go", "/"), "/")
	r.addRoute(paths, mockHandler)

	var node *route
	node, _ = r.getRoute([]string{})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{""})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"", "x"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"", "x", "y"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"x", "y"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"x", "y", "z"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"x"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"y"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"rego"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"rego", "y"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"rego", "y", "rego"})
	equal(t, nil, node)
	node, _ = r.getRoute([]string{"rego", "rego"})
	notEqual(t, nil, node)
}

func TestAddRootRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	r.addRoute(paths, mockHandler)

	equal(t, true, r.isEndPoint)

	if len(r.nodes) > 0 {
		t.Error("Rout should not contain childs")
	}
}

func TestAddEmptyRootRoute(t *testing.T) {
	r := newRoute(nil, "")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	r.addRoute(paths, mockHandler)

	equal(t, true, r.isEndPoint)

	if len(r.nodes) > 0 {
		t.Error("Rout should not contain childs")
	}
}

func TestAddEmptyRootRouteTwo(t *testing.T) {
	r := newRoute(nil, "")
	paths := strings.Split(strings.Trim("", "/"), "/")
	r.addRoute(paths, mockHandler)

	equal(t, true, r.isEndPoint)

	if len(r.nodes) > 0 {
		t.Error("Rout should not contain childs")
	}
}

func TestAddRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/example", "/"), "/")
	r.addRoute(paths, mockHandler)

	equal(t, false, r.isEndPoint)
	if notEqual(t, nil, r.nodes["example"]) {
		equal(t, true, r.nodes["example"].isEndPoint)
		equal(t, "example", r.nodes["example"].path)
		equal(t, nil, r.nodes["example"].regexp)
	}
}

func TestAddParamRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/:example", "/"), "/")
	r.addRoute(paths, mockHandler)

	equal(t, false, r.isEndPoint)
	if notEqual(t, nil, r.nodes[":example"]) {
		equal(t, true, r.nodes[":example"].isEndPoint)
		equal(t, ":example", r.nodes[":example"].path)
		equal(t, nil, r.nodes[":example"].regexp)
	}
}

func TestAddRegexpRoute(t *testing.T) {
	r := newRoute(nil, "/")
	paths := strings.Split(strings.Trim("/:example:r([a-z]+)go", "/"), "/")
	r.addRoute(paths, mockHandler)

	equal(t, false, r.isEndPoint)

	route := r.nodes[":example:r([a-z]+)go"]
	if notEqual(t, nil, route) {
		equal(t, true, route.isEndPoint)
		equal(t, ":example:r([a-z]+)go", route.path)
		if notEqual(t, nil, route.regexp) {
			equal(t, true, route.regexp.MatchString("rego"))
		}
	}
}
