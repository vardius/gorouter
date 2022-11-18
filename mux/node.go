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
		middleware:    middleware.NewCollection(),
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

// RouteAware represents route aware Node
type RouteAware interface {
	// MatchRoute matches given path to Route within Node and its Tree
	MatchRoute(path string) (Route, context.Params)

	// Route provides Node's Route if assigned
	Route() Route
	// WithRoute assigns Route to given Node
	WithRoute(r Route)

	// MaxParamsSize provides maximum number of parameters Route can have for given Node
	MaxParamsSize() uint8
	// SkipSubPath sets skipSubPath node property to true
	// will skip children match search and return current node directly
	// this value is used when matching subrouter
	SkipSubPath()
}

// MiddlewareAware represents middleware aware node
type MiddlewareAware interface {
	// MatchMiddleware collects middleware from all nodes within tree matching given path
	// middleware is merged in order nodes where created, collecting from top to bottom
	MatchMiddleware(path string) middleware.Collection

	// Middleware provides Node's middleware collection
	Middleware() middleware.Collection
	// AppendMiddleware appends middleware collection to Node
	AppendMiddleware(m middleware.Collection)
	// PrependMiddleware prepends middleware collection to Node
	PrependMiddleware(m middleware.Collection)
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
	middleware middleware.Collection

	maxParamsSize uint8
	skipSubPath   bool
}

func (n *staticNode) MatchRoute(path string) (Route, context.Params) {
	nameLength := len(n.name)
	pathLength := len(path)

	if pathLength >= nameLength && n.name == path[:nameLength] {
		if nameLength == pathLength || n.skipSubPath {
			return n.route, make(context.Params, n.maxParamsSize)
		}
		if path[nameLength:nameLength+1] == "/" { // skip slashes only
			return n.children.MatchRoute(path[nameLength+1:]) // +1 because we wan to skip slash as well
		}
	}

	return nil, nil
}

func (n *staticNode) MatchMiddleware(path string) middleware.Collection {
	nameLength := len(n.name)
	pathLength := len(path)

	if pathLength >= nameLength && n.name == path[:nameLength] {
		if nameLength == pathLength || n.skipSubPath {
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

func (n *staticNode) Middleware() middleware.Collection {
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

func (n *staticNode) AppendMiddleware(m middleware.Collection) {
	n.middleware = n.middleware.Merge(m)
}

func (n *staticNode) PrependMiddleware(m middleware.Collection) {
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

func (n *wildcardNode) MatchRoute(path string) (Route, context.Params) {
	pathPart, subPath := pathutils.GetPart(path)
	maxParamsSize := n.MaxParamsSize()

	var route Route
	var params context.Params

	if subPath == "" || n.staticNode.skipSubPath {
		route = n.route
		params = make(context.Params, maxParamsSize)
	} else {
		route, params = n.children.MatchRoute(subPath)
		if route == nil {
			return nil, nil
		}
	}

	params.Set(maxParamsSize-1, n.name, pathPart)

	return route, params
}

func (n *wildcardNode) MatchMiddleware(path string) middleware.Collection {
	_, subPath := pathutils.GetPart(path)

	if subPath == "" || n.staticNode.skipSubPath {
		return n.middleware
	}

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

func (n *regexpNode) MatchRoute(path string) (Route, context.Params) {
	pathPart, subPath := pathutils.GetPart(path)
	if !n.regexp.MatchString(pathPart) {
		return nil, nil
	}

	maxParamsSize := n.MaxParamsSize()

	var route Route
	var params context.Params

	if subPath == "" || n.staticNode.skipSubPath {
		route = n.route
		params = make(context.Params, maxParamsSize)
	} else {
		route, params = n.children.MatchRoute(subPath)
		if route == nil {
			return nil, nil
		}
	}

	params.Set(maxParamsSize-1, n.name, pathPart)

	return route, params
}

func (n *regexpNode) MatchMiddleware(path string) middleware.Collection {
	pathPart, subPath := pathutils.GetPart(path)
	if !n.regexp.MatchString(pathPart) {
		return nil
	}

	if subPath == "" || n.staticNode.skipSubPath {
		return n.middleware
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
