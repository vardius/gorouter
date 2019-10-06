package mux

import (
	"regexp"

	"github.com/vardius/gorouter/v4/context"
	pathutils "github.com/vardius/gorouter/v4/path"
)

// NewNode provides new node
func NewNode(pathPart string, parent Node) Node {
	if len(pathPart) == 0 {
		return nil
	}

	switch parent.(type) {
	case *subrouterNode:
		panic("Subrouter node can not be a parent")
	}

	id, exp := pathutils.GetNameFromPart(pathPart)
	static := &staticNode{
		name:     id,
		children: NewTree(),
	}

	if parent != nil {
		static.maxParamsSize = parent.MaxParamsSize()
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

	if parent != nil {
		parent.WithChildren(parent.Tree().WithNode(node))
	}

	return node
}

// Node is route node
type Node interface {
	Find(names []string) Node
	Match(pathPart string, subPath string) (Node, context.Params, string)

	Name() string
	MaxParamsSize() uint8
	Tree() Tree
	Route() Route

	WithRoute(r Route)
	WithChildren(t Tree)
	WithChild(parts []string) Node
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

func (n *staticNode) WithChild(parts []string) Node {
	if len(parts) == 0 {
		return n
	}

	name, _ := pathutils.GetNameFromPart(parts[0])
	node := n.Find([]string{name})

	if node == nil {
		node = NewNode(parts[0], n)
	}

	return node.WithChild(parts[1:])
}

func (n *staticNode) Find(names []string) Node {
	if len(names) == 0 {
		return n
	}

	if node := n.children.Find(names[0]); node != nil {
		return node.Find(names[1:])
	}

	return n
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

func (n *subrouterNode) Match(pathPart string, subPath string) (Node, context.Params, string) {
	if node, params, _ := n.Node.Match(pathPart, ""); node != nil {

		return node, params, subPath
	}

	return nil, nil, ""
}
