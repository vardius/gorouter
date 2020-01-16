package mux

import (
	"regexp"

	"github.com/vardius/gorouter/v4/context"
	"github.com/vardius/gorouter/v4/middleware"
	pathutils "github.com/vardius/gorouter/v4/path"
)

// NewNode provides new mux Node
func NewNode(pathPart string, maxParamsSize uint8) Node {
	if len(pathPart) == 0 {
		return nil
	}

	name, exp := pathutils.GetNameFromPart(pathPart)
	static := &staticNode{
		name:          name,
		children:      NewTree(),
		middleware:    middleware.New(),
		maxParamsSize: maxParamsSize,
	}

	var node Node

	if exp != "" {
		static.maxParamsSize++
		node = withRegexp(static, regexp.MustCompile(exp))
	} else if name != pathPart {
		static.maxParamsSize++
		node = withWildcard(static)
	} else {
		node = static
	}

	return node
}

// Node represents mux Node
// Can match path and provide routes
type Node interface {
	// Match matches given path to Node within Node and its Tree
	Match(path string) (Node, middleware.Middleware, context.Params, string)

	// Name provides Node name
	Name() string
	// Tree provides next level Node Tree
	Tree() Tree
	// Route provides Node's Route if assigned
	Route() Route
	// Middleware provides Node's middleware
	Middleware() middleware.Middleware

	// Name provides maximum number of parameters Route can have for given Node
	MaxParamsSize() uint8

	// WithRoute assigns Route to given Node
	WithRoute(r Route)
	// WithChildren sets Node's Tree
	WithChildren(t Tree)
	// AppendMiddleware appends middleware to Node
	AppendMiddleware(m middleware.Middleware)
	// PrependMiddleware prepends middleware to Node
	PrependMiddleware(m middleware.Middleware)

	// SkipSubPath sets skipSubPath node property to true
	// will skip children match search and return current node directly
	// this value is used when matching subrouter
	SkipSubPath()
}

type staticNode struct {
	name     string
	children Tree

	route      Route
	middleware middleware.Middleware

	maxParamsSize uint8
	skipSubPath   bool
}

func (n *staticNode) Match(path string) (Node, middleware.Middleware, context.Params, string) {
	nameLength := len(n.name)
	pathLength := len(path)
	if pathLength >= nameLength && n.name == path[:nameLength] || regexp.MustCompile(`{|}`).MatchString(n.name) {
		if nameLength+1 >= pathLength {
			// is there a better solution here ? it stopped working once the braces were included in node.name
			return n, n.middleware, context.Params{{Key: "param", Value: path}}, ""
		}

		if n.skipSubPath {
			return n, n.middleware, context.Params{{Key: "param", Value: path}}, path[nameLength+1:]
		}
		node, treeMiddleware, params, p := n.children.Match(path[nameLength+1:]) // +1 because we wan to skip slash as well
		return node, n.middleware.Merge(treeMiddleware), params, p
	}

	return nil, nil, nil, ""
}

func (n *staticNode) Name() string {
	return n.name
}

func (n *staticNode) Tree() Tree {
	return n.children
}

func (n *staticNode) Route() Route {
	return n.route
}

func (n *staticNode) Middleware() middleware.Middleware {
	return n.middleware
}

func (n *staticNode) MaxParamsSize() uint8 {
	return n.maxParamsSize
}

func (n *staticNode) WithChildren(t Tree) {
	n.children = t
}

func (n *staticNode) WithRoute(r Route) {
	n.route = r
}

func (n *staticNode) AppendMiddleware(m middleware.Middleware) {
	n.middleware = n.middleware.Merge(m)
}

func (n *staticNode) PrependMiddleware(m middleware.Middleware) {
	n.middleware = m.Merge(n.middleware)
}

func (n *staticNode) SkipSubPath() {
	n.skipSubPath = true
}

func withWildcard(parent *staticNode) *wildcardNode {
	return &wildcardNode{staticNode: parent}
}

type wildcardNode struct {
	*staticNode
}

func (n *wildcardNode) Match(path string) (Node, middleware.Middleware, context.Params, string) {
	pathPart, subPath := pathutils.GetPart(path)
	maxParamsSize := n.MaxParamsSize()

	var node Node
	var treeMiddleware middleware.Middleware
	var params context.Params

	if subPath == "" || n.staticNode.skipSubPath {
		node = n
		treeMiddleware = n.Middleware()
		params = make(context.Params, maxParamsSize)
	} else {
		node, treeMiddleware, params, subPath = n.children.Match(subPath)

		if node == nil {
			return nil, nil, nil, ""
		}

		treeMiddleware = n.middleware.Merge(treeMiddleware)
	}

	params.Set(maxParamsSize-1, n.name, pathPart)

	return node, treeMiddleware, params, subPath
}

func withRegexp(parent *staticNode, regexp *regexp.Regexp) *regexpNode {
	return &regexpNode{
		staticNode: parent,
		regexp:     regexp,
	}
}

type regexpNode struct {
	*staticNode

	regexp *regexp.Regexp
}

func (n *regexpNode) Match(path string) (Node, middleware.Middleware, context.Params, string) {
	pathPart, subPath := pathutils.GetPart(path)
	if !n.regexp.MatchString(pathPart) {
		return nil, nil, nil, ""
	}

	maxParamsSize := n.MaxParamsSize()

	var node Node
	var treeMiddleware middleware.Middleware
	var params context.Params

	if subPath == "" || n.staticNode.skipSubPath {
		node = n
		treeMiddleware = n.Middleware()
		params = make(context.Params, maxParamsSize)
	} else {
		node, treeMiddleware, params, subPath = n.children.Match(subPath)

		if node == nil {
			return nil, nil, nil, ""
		}

		treeMiddleware = n.middleware.Merge(treeMiddleware)
	}

	params.Set(maxParamsSize-1, n.name, pathPart)

	return node, treeMiddleware, params, subPath
}

func withSubrouter(parent Node) *subrouterNode {
	parent.SkipSubPath()

	return &subrouterNode{Node: parent}
}

type subrouterNode struct {
	Node
}

func (n *subrouterNode) WithChildren(_ Tree) {
	panic("Subrouter node can not have children.")
}
