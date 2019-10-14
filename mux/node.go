package mux

import (
	"fmt"
	"regexp"

	"github.com/vardius/gorouter/v4/context"
	pathutils "github.com/vardius/gorouter/v4/path"
)

const wildcardRegexp = "\\w+"

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
		node = withRegexp(static, regexp.MustCompile(exp))
	} else if id != pathPart {
		static.maxParamsSize++
		node = withWildcard(static)
	} else {
		node = static
	}

	return node
}

// Node is route node
type Node interface {
	Match(path string) (Node, context.Params, string)
	Merge(node Node) Node

	Name() string
	MaxParamsSize() uint8
	Tree() Tree
	Route() Route

	Rename(string)
	WithRoute(r Route)
	WithChildren(t Tree)
}

type staticNode struct {
	name          string
	route         Route
	children      Tree
	maxParamsSize uint8
}

func (n *staticNode) Merge(node Node) Node {
	var newNode Node
	name := n.name

	n.WithChildren(node.Tree())

	switch castedNode := node.(type) {
	case *subrouterNode:
		newNode = n.Merge(castedNode.Node)
	case *staticNode:
		n.name = fmt.Sprintf("%s/%s", name, node.Name())
		newNode = n
	case *wildcardNode:
		n.name = fmt.Sprintf("%s/{%s}", name, node.Name())
		newNode = withRegexp(n, regexp.MustCompile(fmt.Sprintf("^%s\\/(?P<%s>%s)", name, castedNode.Name(), wildcardRegexp)))
	case *regexpNode:
		n.name = fmt.Sprintf("%s/{%s}", name, node.Name())
		newNode = withRegexp(n, regexp.MustCompile(fmt.Sprintf("^%s\\/(?P<%s>%s)", name, castedNode.Name(), castedNode.regexp.String())))
	}

	return newNode
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

func (n *staticNode) Rename(name string) {
	n.name = name
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

func withWildcard(parent Node) *wildcardNode {
	return &wildcardNode{Node: parent}
}

type wildcardNode struct {
	Node
}

func (n *wildcardNode) Merge(node Node) Node {
	var newNode Node
	name := n.Name()

	n.WithChildren(node.Tree())

	switch castedNode := node.(type) {
	case *subrouterNode:
		newNode = n.Merge(castedNode.Node)
	case *staticNode:
		n.Rename(fmt.Sprintf("{%s}/%s", name, castedNode.name))
		newNode = withRegexp(n, regexp.MustCompile(fmt.Sprintf("^(?P<%s>%s)\\/(?P<%s>%s)", name, wildcardRegexp, castedNode.name, castedNode.name)))
	case *wildcardNode:
		n.Rename(fmt.Sprintf("{%s}/{%s}", name, castedNode.Name()))
		newNode = withRegexp(n, regexp.MustCompile(fmt.Sprintf("^(?P<%s>%s)\\/(?P<%s>%s)", name, wildcardRegexp, castedNode.Name(), wildcardRegexp)))
	case *regexpNode:
		n.Rename(fmt.Sprintf("{%s}/{%s}", name, castedNode.Name()))
		newNode = withRegexp(n, regexp.MustCompile(fmt.Sprintf("^(?P<%s>%s)\\/(?P<%s>%s)", name, wildcardRegexp, castedNode.Name(), castedNode.regexp.String())))
	}

	return newNode
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
		node, params, subPath = n.Node.Tree().Match(subPath)

		if node == nil {
			return nil, nil, ""
		}
	}

	params.Set(maxParamsSize-1, n.Node.Name(), pathPart)

	return node, params, subPath
}

func withRegexp(parent Node, regexp *regexp.Regexp) *regexpNode {
	return &regexpNode{
		Node:   parent,
		regexp: regexp,
	}
}

type regexpNode struct {
	Node

	regexp *regexp.Regexp
}

func (n *regexpNode) Merge(node Node) Node {
	var newNode Node

	n.WithChildren(node.Tree())

	switch castedNode := node.(type) {
	case *subrouterNode:
		newNode = n.Merge(castedNode.Node)
	case *staticNode:
		n.regexp = regexp.MustCompile(fmt.Sprintf("^(?P<%s>%s)\\/(?P<%s>%s)", n.Name(), n.regexp.String(), castedNode.name, castedNode.name))
		n.Rename(fmt.Sprintf("{%s}/%s", n.Name(), castedNode.name))
		newNode = n
	case *wildcardNode:
		n.regexp = regexp.MustCompile(fmt.Sprintf("^(?P<%s>%s)\\/(?P<%s>%s)", n.Name(), n.regexp.String(), castedNode.Name(), wildcardRegexp))
		n.Rename(fmt.Sprintf("{%s}/{%s}", n.Name(), castedNode.Name()))
		newNode = n
	case *regexpNode:
		n.regexp = regexp.MustCompile(fmt.Sprintf("^(?P<%s>%s)\\/(?P<%s>%s)", n.Name(), n.regexp.String(), castedNode.Name(), castedNode.regexp.String()))
		n.Rename(fmt.Sprintf("{%s}/{%s}", n.Name(), castedNode.Name()))
		newNode = n
	}

	return newNode
}

func (n *regexpNode) Match(path string) (Node, context.Params, string) {
	match := n.regexp.FindStringSubmatch(path)

	if len(match) == 0 {
		return nil, nil, ""
	}

	var node Node
	var params context.Params

	subPath := path[len(match[0])+1:]
	maxParamsSize := n.MaxParamsSize()

	if subPath == "" {
		node = n
		params = make(context.Params, maxParamsSize)
	} else {
		node, params, subPath = n.Node.Tree().Match(subPath)

		if node == nil {
			return nil, nil, ""
		}
	}

	names := n.regexp.SubexpNames()
	numSubexp := uint8(n.regexp.NumSubexp())

	for i, value := range match[1:] {
		params.Set(maxParamsSize-numSubexp+uint8(i), names[i+1], value)
	}

	return node, params, subPath
}

func withSubrouter(parent Node) *subrouterNode {
	return &subrouterNode{Node: parent}
}

type subrouterNode struct {
	Node
}

func (n *subrouterNode) Match(path string) (Node, context.Params, string) {
	pathPart, subPath := pathutils.GetPart(path)

	node, params, _ := n.Node.Match(pathPart)

	return node, params, subPath
}

func (n *subrouterNode) Merge(node Node) Node {
	return n
}

func (n *subrouterNode) WithChildren(t Tree) {
	panic("Subrouter node can not have children.")
}
