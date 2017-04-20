package goserver

import (
	"strings"
	"testing"
)

func TestGetRootNode(t *testing.T) {
	n := newRoot("/")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	n.addChild(paths)

	equal(t, "/", n.path)
	equal(t, "", n.regexpToString())
	equal(t, true, n.isRoot())
	equal(t, nil, n.parent)
	notEqual(t, nil, n.children)
	equal(t, true, n.isLeaf())

	var node *node
	node, _ = n.child([]string{""})
	equal(t, nil, node)
	node, _ = n.child([]string{"x"})
	equal(t, nil, node)
	node, _ = n.child([]string{"", "x"})
	equal(t, nil, node)
	node, _ = n.child([]string{"x", "y"})
	equal(t, nil, node)
	node, _ = n.child([]string{})
	notEqual(t, nil, node)
}

func TestGetStrictNode(t *testing.T) {
	n := newRoot("/")
	paths := strings.Split(strings.Trim("/x", "/"), "/")
	n.addChild(paths)

	var node *node
	node, _ = n.child([]string{})
	notEqual(t, nil, node)
	node, _ = n.child([]string{""})
	equal(t, nil, node)
	node, _ = n.child([]string{"", "x"})
	equal(t, nil, node)
	node, _ = n.child([]string{"x", "y"})
	equal(t, nil, node)
	node, _ = n.child([]string{"x"})
	notEqual(t, nil, node)
}

func TestGetParamNode(t *testing.T) {
	n := newRoot("/")
	paths := strings.Split(strings.Trim("/:x", "/"), "/")
	n.addChild(paths)

	var node *node
	var params Params
	node, _ = n.child([]string{})
	notEqual(t, nil, node)
	node, _ = n.child([]string{""})
	equal(t, nil, node)
	node, _ = n.child([]string{"", "x"})
	equal(t, nil, node)
	node, _ = n.child([]string{"x", "y"})
	equal(t, nil, node)
	node, params = n.child([]string{"x"})
	notEqual(t, nil, node)
	if notEqual(t, nil, params.Value("x")) {
		equal(t, "x", params.Value("x"))
	}
	node, params = n.child([]string{"y"})
	if notEqual(t, nil, params.Value("x")) {
		equal(t, "y", params.Value("x"))
	}
	notEqual(t, nil, node)
}

func TestGetRegexNode(t *testing.T) {
	n := newRoot("/")
	paths := strings.Split(strings.Trim("/:x:r([a-z]+)go", "/"), "/")
	n.addChild(paths)

	var node *node
	node, _ = n.child([]string{})
	notEqual(t, nil, node)
	node, _ = n.child([]string{""})
	equal(t, nil, node)
	node, _ = n.child([]string{"", "x"})
	equal(t, nil, node)
	node, _ = n.child([]string{"x", "y"})
	equal(t, nil, node)
	node, _ = n.child([]string{"x"})
	equal(t, nil, node)
	node, _ = n.child([]string{"y"})
	equal(t, nil, node)
	node, _ = n.child([]string{"rego", "y"})
	equal(t, nil, node)
	node, _ = n.child([]string{"rego"})
	notEqual(t, nil, node)
}

func TestGetNestedRegexNode(t *testing.T) {
	n := newRoot("/")
	paths := strings.Split(strings.Trim("/:x:r([a-z]+)go/:y:r([a-z]+)go", "/"), "/")
	n.addChild(paths)

	var node *node
	node, _ = n.child([]string{})
	notEqual(t, nil, node)
	node, _ = n.child([]string{""})
	equal(t, nil, node)
	node, _ = n.child([]string{"", "x"})
	equal(t, nil, node)
	node, _ = n.child([]string{"", "x", "y"})
	equal(t, nil, node)
	node, _ = n.child([]string{"x", "y"})
	equal(t, nil, node)
	node, _ = n.child([]string{"x", "y", "z"})
	equal(t, nil, node)
	node, _ = n.child([]string{"x"})
	equal(t, nil, node)
	node, _ = n.child([]string{"y"})
	equal(t, nil, node)
	node, _ = n.child([]string{"rego"})
	notEqual(t, nil, node)
	node, _ = n.child([]string{"rego", "y"})
	equal(t, nil, node)
	node, _ = n.child([]string{"rego", "y", "rego"})
	equal(t, nil, node)
	node, _ = n.child([]string{"rego", "rego"})
	notEqual(t, nil, node)
}

func TestAddRootNode(t *testing.T) {
	n := newRoot("/")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	n.addChild(paths)

	equal(t, true, n.isLeaf())

	if len(n.children) > 0 {
		t.Error("Rout should not contain childs")
	}
}

func TestAddEmptyRootNode(t *testing.T) {
	n := newRoot("")
	paths := strings.Split(strings.Trim("/", "/"), "/")
	n.addChild(paths)

	equal(t, true, n.isLeaf())

	if len(n.children) > 0 {
		t.Error("Rout should not contain childs")
	}
}

func TestAddEmptyRootNodeTwo(t *testing.T) {
	n := newRoot("")
	paths := strings.Split(strings.Trim("", "/"), "/")
	n.addChild(paths)

	equal(t, true, n.isLeaf())

	if len(n.children) > 0 {
		t.Error("Rout should not contain childs")
	}
}

func TestAddNode(t *testing.T) {
	n := newRoot("/")
	paths := strings.Split(strings.Trim("/example", "/"), "/")
	n.addChild(paths)

	equal(t, false, n.isLeaf())
	if notEqual(t, nil, n.children["example"]) {
		equal(t, true, n.children["example"].isLeaf())
		equal(t, "example", n.children["example"].path)
		equal(t, nil, n.children["example"].regexp)
	}
}

func TestAddParamNode(t *testing.T) {
	n := newRoot("/")
	paths := strings.Split(strings.Trim("/:example", "/"), "/")
	n.addChild(paths)

	equal(t, false, n.isLeaf())
	if notEqual(t, nil, n.children[":example"]) {
		equal(t, true, n.children[":example"].isLeaf())
		equal(t, ":example", n.children[":example"].path)
		equal(t, nil, n.children[":example"].regexp)
	}
}

func TestAddRegexpNode(t *testing.T) {
	n := newRoot("/")
	paths := strings.Split(strings.Trim("/:example:r([a-z]+)go", "/"), "/")
	n.addChild(paths)

	equal(t, false, n.isLeaf())

	route := n.children[":example:r([a-z]+)go"]
	if notEqual(t, nil, route) {
		equal(t, true, route.isLeaf())
		equal(t, ":example:r([a-z]+)go", route.path)
		if notEqual(t, nil, route.regexp) {
			equal(t, true, route.regexp.MatchString("rego"))
		}
	}
}
