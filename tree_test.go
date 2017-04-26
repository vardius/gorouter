package goserver

// import (
// 	"testing"
// )

// func TestGetRootNode(t *testing.T) {
// 	n := newRoot("")
// 	paths := strings.Split(strings.Trim("/", "/"), "/")
// 	n.addChild(paths)

// 	equal(t, "", n.id)
// 	equal(t, "", n.regexpToString())
// 	equal(t, true, n.isRoot())
// 	equal(t, nil, n.parent)
// 	notEqual(t, nil, n.children)
// 	equal(t, true, n.isLeaf())

// 	var node *node
// 	node, _ = n.childByPath("")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("x")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/x")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/x/y")
// 	equal(t, nil, node)
// }

// func TestGetStrictNode(t *testing.T) {
// 	n := newRoot("")
// 	paths := strings.Split(strings.Trim("/x", "/"), "/")
// 	n.addChild(paths)

// 	var node *node
// 	node, _ = n.childByPath("")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/x")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("x")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/x/y")
// 	equal(t, nil, node)
// }

// func TestGetParamNode(t *testing.T) {
// 	n := newRoot("")
// 	paths := strings.Split(strings.Trim("/{x}", "/"), "/")
// 	n.addChild(paths)

// 	var node *node
// 	node, _ = n.childByPath("")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/x/y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/x")
// 	notEqual(t, nil, node)

// 	var params Params
// 	node, params = n.childByPath("x")
// 	notEqual(t, nil, node)
// 	if notEqual(t, nil, params.Value("x")) {
// 		equal(t, "x", params.Value("x"))
// 	}
// 	node, params = n.childByPath("y")
// 	notEqual(t, nil, node)
// 	if notEqual(t, nil, params.Value("x")) {
// 		equal(t, "y", params.Value("x"))
// 	}
// }

// func TestGetRegexNode(t *testing.T) {
// 	n := newRoot("")
// 	paths := strings.Split(strings.Trim("/{x:r([a-z]+)go}", "/"), "/")
// 	n.addChild(paths)

// 	var node *node
// 	node, _ = n.childByPath("")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/x")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/x/y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("x")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/rego/y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("rego")
// 	notEqual(t, nil, node)
// }

// func TestGetNestedRegexNode(t *testing.T) {
// 	n := newRoot("")
// 	paths := strings.Split(strings.Trim("/{x:r([a-z]+)go}/{y:r([a-z]+)go}", "/"), "/")
// 	n.addChild(paths)

// 	var node *node
// 	node, _ = n.childByPath("")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/x")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("//x/y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/x/y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/x/y/z")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("x")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("rego")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/rego/y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/rego/y/rego")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/rego/rego")
// 	notEqual(t, nil, node)
// }

// func TestAddRootNode(t *testing.T) {
// 	n := newRoot("")
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
// 	n := newRoot("")
// 	paths := strings.Split(strings.Trim("/example", "/"), "/")
// 	n.addChild(paths)

// 	var cn *node
// 	for _, child := range n.children {
// 		if child.id == "example" {
// 			cn = child
// 			break
// 		}
// 	}

// 	equal(t, false, n.isLeaf())
// 	if notEqual(t, nil, cn) {
// 		equal(t, true, cn.isLeaf())
// 		equal(t, "example", cn.id)
// 		equal(t, nil, cn.regexp)
// 	}
// }

// func TestAddParamNode(t *testing.T) {
// 	n := newRoot("")
// 	paths := strings.Split(strings.Trim("/{example}", "/"), "/")
// 	n.addChild(paths)

// 	var cn *node
// 	for _, child := range n.children {
// 		if child.id == "{example}" {
// 			cn = child
// 			break
// 		}
// 	}

// 	equal(t, false, n.isLeaf())
// 	if notEqual(t, nil, cn) {
// 		equal(t, true, cn.isLeaf())
// 		equal(t, "{example}", cn.id)
// 		equal(t, nil, cn.regexp)
// 	}
// }

// func TestAddRegexpNode(t *testing.T) {
// 	n := newRoot("")
// 	paths := strings.Split(strings.Trim("/{example:r([a-z]+)go}", "/"), "/")
// 	n.addChild(paths)

// 	equal(t, false, n.isLeaf())

// 	var cn *node
// 	for _, child := range n.children {
// 		if child.id == "{example:r([a-z]+)go}" {
// 			cn = child
// 			break
// 		}
// 	}

// 	if notEqual(t, nil, cn) {
// 		equal(t, true, cn.isLeaf())
// 		equal(t, "{example:r([a-z]+)go}", cn.id)
// 		if notEqual(t, nil, cn.regexp) {
// 			equal(t, true, cn.regexp.MatchString("rego"))
// 		}
// 	}
// }

// func TestWildcardConflictNodes(t *testing.T) {
// 	n := newRoot("")
// 	n.addChild(strings.Split(strings.Trim("/{x}/x", "/"), "/"))
// 	n.addChild(strings.Split(strings.Trim("/{y}/y", "/"), "/"))
// 	n.addChild(strings.Split(strings.Trim("/{z}/z", "/"), "/"))

// 	var node *node
// 	node, _ = n.childByPath("")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/x")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/t/x")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/t/y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/t/z")
// 	equal(t, nil, node)
// }

// func TestWildcardRegexpConflictNodes(t *testing.T) {
// 	n := newRoot("")
// 	n.addChild(strings.Split(strings.Trim("/{x:x([a-z]+)go}/x", "/"), "/"))
// 	n.addChild(strings.Split(strings.Trim("/{y:y([a-z]+)go}/y", "/"), "/"))
// 	n.addChild(strings.Split(strings.Trim("/{z:z([a-z]+)go}/z", "/"), "/"))

// 	var node *node
// 	node, _ = n.childByPath("")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/x")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/x/y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("x")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("y")
// 	equal(t, nil, node)
// 	node, _ = n.childByPath("/xego/x")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/yego/y")
// 	notEqual(t, nil, node)
// 	node, _ = n.childByPath("/zego/z")
// 	notEqual(t, nil, node)
// }
