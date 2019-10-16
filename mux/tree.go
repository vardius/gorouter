package mux

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vardius/gorouter/v4/context"
	pathutils "github.com/vardius/gorouter/v4/path"
)

// NewTree provides new empty Tree
func NewTree() Tree {
	return make([]Node, 0)
}

// Tree slice of mux Nodes
type Tree []Node

// PrettyPrint prints the tree text representation to console
func (t Tree) PrettyPrint() string {
	buff := &bytes.Buffer{}

	for _, child := range t {
		switch node := child.(type) {
		case *subrouterNode:
			fmt.Fprintf(buff, "\t_%s\n", node.Name())
		case *staticNode:
			fmt.Fprintf(buff, "\t%s\n", node.Name())
		case *wildcardNode:
			fmt.Fprintf(buff, "\t{%s}\n", node.Name())
		case *regexpNode:
			fmt.Fprintf(buff, "\t{%s:%s}\n", node.Name(), node.regexp.String())
		}

		if len(child.Tree()) > 0 {
			fmt.Fprintf(buff, "\t%s", child.Tree().PrettyPrint())
		}
	}

	return buff.String()
}

// Compile optimizes Tree nodes reducing static nodes depth when possible
func (t Tree) Compile() Tree {
	for i, child := range t {
		child.WithChildren(child.Tree().Compile())

		if staticNode, ok := child.(*staticNode); ok && len(child.Tree()) == 1 {
			node := child.Tree()[0]

			staticNode.WithChildren(node.Tree())
			staticNode.name = fmt.Sprintf("%s/%s", staticNode.name, node.Name())

			t[i] = staticNode
		}
	}

	return t
}

// Match path to Node
func (t Tree) Match(path string) (Node, context.Params, string) {
	for _, child := range t {
		if node, params, subPath := child.Match(path); node != nil {
			return node, params, subPath
		}
	}

	return nil, nil, ""
}

// Find Node inside a tree by name
func (t Tree) Find(name string) Node {
	if name == "" {
		return nil
	}

	for _, child := range t {
		if child.Name() == name {
			return child
		}
	}

	return nil
}

// WithRoute returns new Tree with Route set to Node
// Route is set to Node under the give path, ff Node does not exist it is created
func (t Tree) WithRoute(path string, route Route, maxParamsSize uint8) Tree {
	path = pathutils.TrimSlash(path)
	if path == "" {
		return t
	}

	parts := strings.Split(path, "/")
	name, _ := pathutils.GetNameFromPart(parts[0])
	node := t.Find(name)

	if node == nil {
		node = NewNode(parts[0], maxParamsSize)
		t = t.WithNode(node)
	}

	if len(parts) == 1 {
		node.WithRoute(route)
	} else {
		node.WithChildren(node.Tree().WithRoute(strings.Join(parts[1:], "/"), route, node.MaxParamsSize()))
	}

	return t
}

// WithSubrouter returns new Tree with new Route set to Subrouter Node
// Route is set to Node under the give path, ff Node does not exist it is created
func (t Tree) WithSubrouter(path string, route Route, maxParamsSize uint8) Tree {
	path = pathutils.TrimSlash(path)
	if path == "" {
		return t
	}

	parts := strings.Split(path, "/")
	name, _ := pathutils.GetNameFromPart(parts[0])
	node := t.Find(name)

	if node == nil {
		node = NewNode(parts[0], maxParamsSize)
		if len(parts) == 1 {
			node = withSubrouter(node)
		}
		t = t.WithNode(node)
	}

	if len(parts) == 1 {
		node.WithRoute(route)
	} else {
		node.WithChildren(node.Tree().WithSubrouter(strings.Join(parts[1:], "/"), route, node.MaxParamsSize()))
	}

	return t
}

// WithNode inserts node to Tree
// Nodes are sorted static, regexp, wildcard
func (t Tree) WithNode(node Node) Tree {
	if node == nil {
		return t
	}

	t = append(t, node)

	// @TODO: sort tree

	return t
}
