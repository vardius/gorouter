package mux

import (
	"regexp"

	"github.com/vardius/gorouter/v4/context"
	pathutils "github.com/vardius/gorouter/v4/path"
)

// NewNode provides new node
func NewNode(pathPart string, maxParamsSize uint8) Node {
	if len(pathPart) == 0 {
		return nil
	}

	id, exp := pathutils.GetNameFromPart(pathPart)
	static := &staticNode{
		name:          id,
		children:      NewTree(),
		maxParamsSize: maxParamsSize,
	}

	var node Node

	if exp != "" {
		static.maxParamsSize++
		node = withRegexp(withWildcard(static), regexp.MustCompile(exp))
	} else if id != pathPart {
		static.maxParamsSize++
		node = withWildcard(static)
	} else {
		node = static
	}

	if id != pathPart {
		static.maxParamsSize++
		node = withWildcard(static)

		if exp != "" {
			node = withRegexp(withWildcard(static), regexp.MustCompile(exp))
		}
	}

	return node
}

// Node is route node
type Node interface {
	Match(pathPart string, subPath string) (Node, context.Params, string)

	Name() string
	MaxParamsSize() uint8
	Tree() Tree
	Route() Route

	WithRoute(r Route)
	WithChildren(t Tree)
}

type staticNode struct {
	name          string
	route         Route
	children      Tree
	maxParamsSize uint8
}

func (n *staticNode) WithChildren(t Tree) {
	n.children = t
}

func (n *staticNode) Match(pathPart string, subPath string) (Node, context.Params, string) {
	if n.name == pathPart {
		if subPath == "" {
			return n, make(context.Params, n.maxParamsSize), ""
		}

		if node, params, subPath := n.children.Match(subPath); node != nil {
			return node, params, subPath
		}
	}

	return nil, nil, ""
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

func (n *wildcardNode) Match(pathPart string, subPath string) (Node, context.Params, string) {
	if subPath == "" {
		params := make(context.Params, n.maxParamsSize)

		params.Set(n.maxParamsSize-1, n.staticNode.name, pathPart)

		return n, params, ""
	}

	if node, params, subPath := n.Tree().Match(subPath); node != nil {
		params.Set(n.maxParamsSize-1, n.staticNode.name, pathPart)

		return node, params, subPath
	}

	return nil, nil, ""
}

func withRegexp(parent *wildcardNode, regexp *regexp.Regexp) *regexpNode {
	return &regexpNode{wildcardNode: parent, regexp: regexp}
}

type regexpNode struct {
	*wildcardNode
	regexp *regexp.Regexp
}

func (n *regexpNode) Match(pathPart string, subPath string) (Node, context.Params, string) {
	if n.regexp.MatchString(pathPart) {
		return n.wildcardNode.Match(pathPart, subPath)
	}

	return nil, nil, ""
}

func withSubrouter(parent Node) *subrouterNode {
	return &subrouterNode{Node: parent}
}

type subrouterNode struct {
	Node
}

func (n *subrouterNode) WithChildren(t Tree) {
	panic("Subrouter node can not have children.")
}

func (n *subrouterNode) Match(pathPart string, subPath string) (Node, context.Params, string) {
	if node, params, _ := n.Node.Match(pathPart, ""); node != nil {
		return node, params, subPath
	}

	return nil, nil, ""
}
