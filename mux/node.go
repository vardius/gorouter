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
		node = WithRegexp(static, regexp.MustCompile(exp))
	} else if id != pathPart {
		static.maxParamsSize++
		node = WithWildcard(static)
	} else {
		node = static
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
	if n.name != pathPart {
		return nil, nil, ""
	}

	if node, params, _ := n.children.Match(subPath); node != nil {
		return node, params, ""
	}

	return n, make(context.Params, n.maxParamsSize), ""
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

// WithWildcard returns a copy of parent with a wildcard.
func WithWildcard(parent Node) Node {
	return &wildcardNode{Node: parent}
}

type wildcardNode struct {
	Node
}

func (n *wildcardNode) Match(pathPart string, subPath string) (Node, context.Params, string) {
	if node, params, _ := n.Tree().Match(subPath); node != nil {
		params.Set(n.MaxParamsSize()-1, n.Name(), pathPart)

		return node, params, ""
	}

	params := make(context.Params, n.MaxParamsSize())

	params.Set(n.MaxParamsSize()-1, n.Name(), pathPart)

	return n, params, ""
}

// WithRegexp returns a copy of parent with a regexp wildcard.
func WithRegexp(parent Node, regexp *regexp.Regexp) Node {
	return &regexpNode{Node: parent, regexp: regexp}
}

type regexpNode struct {
	Node
	regexp *regexp.Regexp
}

func (n *regexpNode) Match(pathPart string, subPath string) (Node, context.Params, string) {
	if !n.regexp.MatchString(pathPart) {
		return nil, nil, ""
	}

	if node, params, _ := n.Tree().Match(subPath); node != nil {
		params.Set(n.MaxParamsSize()-1, n.Name(), pathPart)

		return node, params, ""
	}

	params := make(context.Params, n.MaxParamsSize())

	params.Set(n.MaxParamsSize()-1, n.Name(), pathPart)

	return n, params, ""
}

// WithSubrouter returns a copy of parent as a subrouter.
func WithSubrouter(parent Node) Node {
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
