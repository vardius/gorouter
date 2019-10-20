package mux

import (
	"regexp"

	"github.com/vardius/gorouter/v4/context"
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
	Match(path string) (Node, context.Params, string)

	// Name provides Node name
	Name() string
	// Tree provides next level Node Tree
	Tree() Tree
	// Route provides Node's Route if assigned
	Route() Route

	// Name provides maximum number of parameters Route can have for given Node
	MaxParamsSize() uint8

	// WithRoute assigns Route to given Node
	WithRoute(r Route)
	// WithChildren sets Node's Tree
	WithChildren(t Tree)

	// SkipSubPath sets skipSubPath node property to true
	// will skip children match search and return current node directly
	// this value is used when matching subrouter
	SkipSubPath()
}

type staticNode struct {
	name     string
	children Tree

	route Route

	maxParamsSize uint8
	skipSubPath   bool
}

func (n *staticNode) Match(path string) (Node, context.Params, string) {
	nameLength := len(n.name)
	pathLength := len(path)

	if pathLength >= nameLength && n.name == path[:nameLength] {
		if nameLength+1 >= pathLength {
			return n, make(context.Params, n.maxParamsSize), ""
		}

		if n.skipSubPath {
			return n, make(context.Params, n.maxParamsSize), path[nameLength+1:]
		}

		return n.children.Match(path[nameLength+1:]) // +1 because we wan to skip slash as well
	}

	return nil, nil, ""
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

func (n *staticNode) MaxParamsSize() uint8 {
	return n.maxParamsSize
}

func (n *staticNode) WithChildren(t Tree) {
	n.children = t
}

func (n *staticNode) WithRoute(r Route) {
	n.route = r
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

func (n *wildcardNode) Match(path string) (Node, context.Params, string) {
	pathPart, subPath := pathutils.GetPart(path)
	maxParamsSize := n.MaxParamsSize()

	var node Node
	var params context.Params

	if n.staticNode.skipSubPath || subPath == "" {
		node = n
		params = make(context.Params, maxParamsSize)
	} else {
		node, params, subPath = n.children.Match(subPath)

		if node == nil {
			return nil, nil, ""
		}
	}

	params.Set(maxParamsSize-1, n.name, pathPart)

	return node, params, subPath
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

func (n *regexpNode) Match(path string) (Node, context.Params, string) {
	pathPart, subPath := pathutils.GetPart(path)
	if !n.regexp.MatchString(pathPart) {
		return nil, nil, ""
	}

	maxParamsSize := n.MaxParamsSize()

	var node Node
	var params context.Params

	if n.staticNode.skipSubPath || subPath == "" {
		node = n
		params = make(context.Params, maxParamsSize)
	} else {
		node, params, subPath = n.children.Match(subPath)

		if node == nil {
			return nil, nil, ""
		}
	}

	params.Set(maxParamsSize-1, n.name, pathPart)

	return node, params, subPath
}

func withSubrouter(parent Node) *subrouterNode {
	parent.SkipSubPath()

	return &subrouterNode{Node: parent}
}

type subrouterNode struct {
	Node
}

func (n *subrouterNode) Match(path string) (Node, context.Params, string) {
	node, params, subPath := n.Node.Match(path)

	return node, params, subPath
}

func (n *subrouterNode) WithChildren(t Tree) {
	panic("Subrouter node can not have children.")
}
