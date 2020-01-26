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

type RouteAware interface {
	// MatchRoute matches given path to Route within Node and its Tree
	MatchRoute(path string) (Route, context.Params, string)

	// Route provides Node's Route if assigned
	Route() Route
	// WithRoute assigns Route to given Node
	WithRoute(r Route)

	// Name provides maximum number of parameters Route can have for given Node
	MaxParamsSize() uint8
	// SkipSubPath sets skipSubPath node property to true
	// will skip children match search and return current node directly
	// this value is used when matching subrouter
	SkipSubPath()
}

type MiddlewareAware interface {
	// MatchMiddleware collects middleware from all nodes within tree matching given path
	// middleware is merged in order nodes where created, collecting from top to bottom
	MatchMiddleware(path string) middleware.Middleware

	// Middleware provides Node's middleware
	Middleware() middleware.Middleware
	// AppendMiddleware appends middleware to Node
	AppendMiddleware(m middleware.Middleware)
	// PrependMiddleware prepends middleware to Node
	PrependMiddleware(m middleware.Middleware)
}

// Node represents mux Node
// Can match path and provide routes
type Node interface {
	RouteAware
	MiddlewareAware

	// Name provides Node name
	Name() string
	// Tree provides next level Node Tree
	Tree() Tree
	// WithChildren sets Node's Tree
	WithChildren(t Tree)
}

type staticNode struct {
	name     string
	children Tree

	route      Route
	middleware middleware.Middleware

	maxParamsSize uint8
	skipSubPath   bool
}

func (n *staticNode) MatchRoute(path string) (Route, context.Params, string) {
	nameLength := len(n.name)
	pathLength := len(path)

	if pathLength >= nameLength && n.name == path[:nameLength] {
		if nameLength+1 >= pathLength {
			return n.route, make(context.Params, n.maxParamsSize), ""
		}

		if n.skipSubPath {
			return n.route, make(context.Params, n.maxParamsSize), path[nameLength+1:]
		}

		return n.children.MatchRoute(path[nameLength+1:]) // +1 because we wan to skip slash as well
	}

	return nil, nil, ""
}

func (n *staticNode) MatchMiddleware(path string) middleware.Middleware {
	nameLength := len(n.name)
	pathLength := len(path)

	if pathLength >= nameLength && n.name == path[:nameLength] {
		if nameLength+1 >= pathLength {
			return n.middleware
		}

		if treeMiddleware := n.children.MatchMiddleware(path[nameLength+1:]); treeMiddleware != nil { // +1 because we wan to skip slash as well

			return n.middleware.Merge(treeMiddleware)
		}

		return n.middleware
	}

	return nil
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

func (n *wildcardNode) MatchRoute(path string) (Route, context.Params, string) {
	pathPart, subPath := pathutils.GetPart(path)
	maxParamsSize := n.MaxParamsSize()

	var route Route
	var params context.Params

	if subPath == "" || n.staticNode.skipSubPath {
		route = n.route
		params = make(context.Params, maxParamsSize)
	} else {
		route, params, subPath = n.children.MatchRoute(subPath)
		if route == nil {
			return nil, nil, ""
		}
	}

	params.Set(maxParamsSize-1, n.name, pathPart)

	return route, params, subPath
}

func (n *wildcardNode) MatchMiddleware(path string) middleware.Middleware {
	_, subPath := pathutils.GetPart(path)

	if treeMiddleware := n.children.MatchMiddleware(subPath); treeMiddleware != nil {

		return n.middleware.Merge(treeMiddleware)
	}

	return n.middleware
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

func (n *regexpNode) MatchRoute(path string) (Route, context.Params, string) {
	pathPart, subPath := pathutils.GetPart(path)
	if !n.regexp.MatchString(pathPart) {
		return nil, nil, ""
	}

	maxParamsSize := n.MaxParamsSize()

	var route Route
	var params context.Params

	if subPath == "" || n.staticNode.skipSubPath {
		route = n.route
		params = make(context.Params, maxParamsSize)
	} else {
		route, params, subPath = n.children.MatchRoute(subPath)
		if route == nil {
			return nil, nil, ""
		}
	}

	params.Set(maxParamsSize-1, n.name, pathPart)

	return route, params, subPath
}

func (n *regexpNode) MatchMiddleware(path string) middleware.Middleware {
	pathPart, subPath := pathutils.GetPart(path)
	if !n.regexp.MatchString(pathPart) {
		return nil
	}

	if treeMiddleware := n.children.MatchMiddleware(subPath); treeMiddleware != nil {

		return n.middleware.Merge(treeMiddleware)
	}

	return n.middleware
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
