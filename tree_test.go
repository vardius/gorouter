package goserver

// import (
// 	"strings"
// 	"testing"
// )

// func TestGetRootNode(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/", "/"), "/")
// 	n.addChild(paths)

// 	equal(t, "/", n.pattern)
// 	equal(t, "", n.regexpToString())
// 	equal(t, true, n.isRoot())
// 	equal(t, nil, n.parent)
// 	notEqual(t, nil, n.children)
// 	equal(t, true, n.isLeaf())

// 	var node *node
// 	node, _ = n.child([]string{""})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"", "x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{})
// 	notEqual(t, nil, node)
// }

// func TestGetRootNodeNotRecursive(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/", "/"), "/")
// 	n.addChild(paths)

// 	equal(t, "/", n.pattern)
// 	equal(t, "", n.regexpToString())
// 	equal(t, true, n.isRoot())
// 	equal(t, nil, n.parent)
// 	notEqual(t, nil, n.children)
// 	equal(t, true, n.isLeaf())

// 	var node *node
// 	node, _ = n.childNotRecursive([]string{""})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"x"})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"", "x"})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"x", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{})
// 	notEqual(t, nil, node)
// }

// func TestGetStrictNode(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/x", "/"), "/")
// 	n.addChild(paths)

// 	var node *node
// 	node, _ = n.child([]string{})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{""})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"", "x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x"})
// 	notEqual(t, nil, node)
// }

// func TestGetParamNode(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/{x}", "/"), "/")
// 	n.addChild(paths)

// 	var node *node
// 	node, _ = n.child([]string{})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{""})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"", "x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x", "y"})
// 	equal(t, nil, node)

// 	var params Params
// 	node, params = n.child([]string{"x"})
// 	notEqual(t, nil, node)
// 	if notEqual(t, nil, params.Value("x")) {
// 		equal(t, "x", params.Value("x"))
// 	}
// 	node, params = n.child([]string{"y"})
// 	notEqual(t, nil, node)
// 	if notEqual(t, nil, params.Value("x")) {
// 		equal(t, "y", params.Value("x"))
// 	}
// }

// func TestGetRegexNode(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/{x:r([a-z]+)go}", "/"), "/")
// 	n.addChild(paths)

// 	var node *node
// 	node, _ = n.child([]string{})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{""})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"", "x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"rego", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"rego"})
// 	notEqual(t, nil, node)
// }

// func TestGetParamNodeNotRecursive(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/{x}", "/"), "/")
// 	n.addChild(paths)

// 	var node *node
// 	node, _ = n.childNotRecursive([]string{})
// 	notEqual(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{""})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"", "x"})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"x", "y"})
// 	equal(t, nil, node)

// 	var params Params
// 	node, params = n.childNotRecursive([]string{"x"})
// 	notEqual(t, nil, node)
// 	if notEqual(t, nil, params.Value("x")) {
// 		equal(t, "x", params.Value("x"))
// 	}
// 	node, params = n.childNotRecursive([]string{"y"})
// 	notEqual(t, nil, node)
// 	if notEqual(t, nil, params.Value("x")) {
// 		equal(t, "y", params.Value("x"))
// 	}
// }

// func TestGetRegexNodeNotRecursive(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/{x:r([a-z]+)go}", "/"), "/")
// 	n.addChild(paths)

// 	var node *node
// 	node, _ = n.childNotRecursive([]string{})
// 	notEqual(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{""})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"", "x"})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"x", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"x"})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"y"})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"rego", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.childNotRecursive([]string{"rego"})
// 	notEqual(t, nil, node)
// }

// func TestGetNestedRegexNode(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/{x:r([a-z]+)go/{y:r([a-z]+)go", "/"), "/")
// 	n.addChild(paths)

// 	var node *node
// 	node, _ = n.child([]string{})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{""})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"", "x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"", "x", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x", "y", "z"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"rego"})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{"rego", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"rego", "y", "rego"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"rego", "rego"})
// 	notEqual(t, nil, node)
// }

// func TestAddRootNode(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/", "/"), "/")
// 	n.addChild(paths)

// 	equal(t, true, n.isLeaf())

// 	if len(n.children) > 0 {
// 		t.Error("Rout should not contain childs")
// 	}
// }

// func TestAddEmptyRootNode(t *testing.T) {
// 	n := newRoot("")
// 	paths := strings.Split(strings.Trim("/", "/"), "/")
// 	n.addChild(paths)

// 	equal(t, true, n.isLeaf())

// 	if len(n.children) > 0 {
// 		t.Error("Rout should not contain childs")
// 	}
// }

// func TestAddEmptyRootNodeTwo(t *testing.T) {
// 	n := newRoot("")
// 	paths := strings.Split(strings.Trim("", "/"), "/")
// 	n.addChild(paths)

// 	equal(t, true, n.isLeaf())

// 	if len(n.children) > 0 {
// 		t.Error("Rout should not contain childs")
// 	}
// }

// func TestAddNode(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/example", "/"), "/")
// 	n.addChild(paths)

// 	var cn *node
// 	for _, child := range n.children {
// 		if child.pattern == "example" {
// 			cn = child
// 			break
// 		}
// 	}

// 	equal(t, false, n.isLeaf())
// 	if notEqual(t, nil, cn) {
// 		equal(t, true, cn.isLeaf())
// 		equal(t, "example", cn.pattern)
// 		equal(t, nil, cn.regexp)
// 	}
// }

// func TestAddParamNode(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/{example}", "/"), "/")
// 	n.addChild(paths)

// 	var cn *node
// 	for _, child := range n.children {
// 		if child.pattern == "{example}" {
// 			cn = child
// 			break
// 		}
// 	}

// 	equal(t, false, n.isLeaf())
// 	if notEqual(t, nil, cn) {
// 		equal(t, true, cn.isLeaf())
// 		equal(t, "{example}", cn.pattern)
// 		equal(t, nil, cn.regexp)
// 	}
// }

// func TestAddRegexpNode(t *testing.T) {
// 	n := newRoot("/")
// 	paths := strings.Split(strings.Trim("/{example:r([a-z]+)go}", "/"), "/")
// 	n.addChild(paths)

// 	equal(t, false, n.isLeaf())

// 	var cn *node
// 	for _, child := range n.children {
// 		if child.pattern == "{example}:r([a-z]+)go}" {
// 			cn = child
// 			break
// 		}
// 	}

// 	if notEqual(t, nil, cn) {
// 		equal(t, true, cn.isLeaf())
// 		equal(t, "{example}:r([a-z]+)go}", cn.pattern)
// 		if notEqual(t, nil, cn.regexp) {
// 			equal(t, true, cn.regexp.MatchString("rego"))
// 		}
// 	}
// }

// func TestWildcardConflictNodes(t *testing.T) {
// 	n := newRoot("/")
// 	n.addChild(strings.Split(strings.Trim("/{x}/x", "/"), "/"))
// 	n.addChild(strings.Split(strings.Trim("/{y}/y", "/"), "/"))
// 	n.addChild(strings.Split(strings.Trim("/{z}/z", "/"), "/"))

// 	var node *node
// 	node, _ = n.child([]string{})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{""})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"", "x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"t", "x"})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{"t", "y"})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{"t", "z"})
// 	notEqual(t, nil, node)
// }

// func TestWildcardRegexpConflictNodes(t *testing.T) {
// 	n := newRoot("/")
// 	n.addChild(strings.Split(strings.Trim("/{x:x([a-z]+)go}/x", "/"), "/"))
// 	n.addChild(strings.Split(strings.Trim("/{y:y([a-z]+)go}/y", "/"), "/"))
// 	n.addChild(strings.Split(strings.Trim("/{z:z([a-z]+)go}/z", "/"), "/"))

// 	var node *node
// 	node, _ = n.child([]string{})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{""})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"", "x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x", "y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"x"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"y"})
// 	equal(t, nil, node)
// 	node, _ = n.child([]string{"xego", "x"})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{"yego", "y"})
// 	notEqual(t, nil, node)
// 	node, _ = n.child([]string{"zego", "z"})
// 	notEqual(t, nil, node)
// }
