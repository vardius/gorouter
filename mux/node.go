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
	// Name provides maximum number of parameters Route can have for given Node
	MaxParamsSize() uint8
	// Tree provides next level Node Tree
	Tree() Tree
	// Route provides Node's Route if assigned
	Route() Route

	// WithRoute assigns Route to given Node
	WithRoute(r Route)
	// WithChildren sets Node's Tree
	WithChildren(t Tree)
}

type staticNode struct {
	name          string
	route         Route
	children      Tree
	maxParamsSize uint8
}

func (n *staticNode) Match(path string) (Node, context.Params, string) {
	nameLength := len(n.name)
	pathLength := len(path)

	if pathLength >= nameLength && n.name == path[:nameLength] {
		if nameLength+1 >= pathLength {
			return n, make(context.Params, n.maxParamsSize), ""
		}

		return n.children.Match(path[nameLength+1:]) // +1 because we wan to skip slash as well
	}

	return nil, nil, ""
}

func (n *staticNode) WithChildren(t Tree) {
	n.children = t
}

func (n *staticNode) Name() string {
	return n.name
}

func (n *staticNode) MaxParamsSize() uint8 {
	return n.maxParamsSize
}

func (n *staticNode) Tree() Tree {
	return n.children
}

func (n *staticNode) Route() Route {
	return n.route
}

func (n *staticNode) WithRoute(r Route) {
	n.route = r
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

	if subPath == "" {
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

	if subPath == "" {
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
	return &subrouterNode{Node: parent}
}

type subrouterNode struct {
	Node
}

func (n *subrouterNode) Match(path string) (Node, context.Params, string) {
	switch node := n.Node.(type) {
	case *staticNode:
		nameLength := len(node.name)
		n, params, _ := node.Match(path[:nameLength])

		if nameLength < len(path) {
			return n, params, path[nameLength+1:]
		}

		return n, params, ""
	case *wildcardNode:
		pathPart, subPath := pathutils.GetPart(path)
		n, params, _ := node.Match(pathPart)

		return n, params, subPath
	case *regexpNode:
		pathPart, subPath := pathutils.GetPart(path)
		n, params, _ := node.Match(pathPart)

		return n, params, subPath
	case *subrouterNode:
		return node.Match(path)
	}

	return nil, nil, ""
}

func (n *subrouterNode) WithChildren(t Tree) {
	panic("Subrouter node can not have children.")
}
