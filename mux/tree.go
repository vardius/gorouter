package mux

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/vardius/gorouter/v4/context"
	"github.com/vardius/gorouter/v4/middleware"
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
		case *staticNode:
			fmt.Fprintf(buff, "\t%s\n", node.Name())
		case *wildcardNode:
			fmt.Fprintf(buff, "\t{%s}\n", node.Name())
		case *regexpNode:
			fmt.Fprintf(buff, "\t{%s:%s}\n", node.Name(), node.regexp.String())
		case *subrouterNode:
			fmt.Fprintf(buff, "\t_%s\n", node.Name())
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

		if len(child.Tree()) == 1 {
			switch node := child.(type) {
			case *staticNode:
				if staticNode, ok := node.Tree()[0].(*staticNode); ok {
					node.WithChildren(staticNode.Tree())
					node.AppendMiddleware(staticNode.Middleware())
					node.name = fmt.Sprintf("%s/%s", node.name, staticNode.name)

					t[i] = node
				}
				// skip
				// case *wildcardNode:
				// case *regexpNode:
				// case *subrouterNode:
			}
		}
	}

	return t
}

// Match path to Node
func (t Tree) Match(path string) (Node, middleware.Middleware, context.Params, string) {
	for _, child := range t {
		if node, m, params, subPath := child.Match(path); node != nil {
			return node, m, params, subPath
		}
	}

	return nil, nil, nil, ""
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
// Route is set to Node under the give path, if Node does not exist it is created
func (t Tree) WithRoute(path string, route Route, maxParamsSize uint8) Tree {
	path = pathutils.TrimSlash(path)
	if path == "" {
		return t
	}

	parts := strings.Split(path, "/")
	name, _ := pathutils.GetNameFromPart(parts[0])
	node := t.Find(name)
	newTree := t

	if node == nil {
		node = NewNode(parts[0], maxParamsSize)
		newTree = t.withNode(node)
	}

	if len(parts) == 1 {
		node.WithRoute(route)
	} else {
		node.WithChildren(node.Tree().WithRoute(strings.Join(parts[1:], "/"), route, node.MaxParamsSize()))
	}

	return newTree
}

// WithMiddleware returns new Tree with Middleware appended to given Node
// Middleware is appended to Node under the give path, if Node does not exist it is created
func (t Tree) WithMiddleware(path string, m middleware.Middleware, maxParamsSize uint8) Tree {
	path = pathutils.TrimSlash(path)
	if path == "" {
		return t
	}

	parts := strings.Split(path, "/")
	name, _ := pathutils.GetNameFromPart(parts[0])
	node := t.Find(name)
	newTree := t

	if node == nil {
		node = NewNode(parts[0], maxParamsSize)
		newTree = t.withNode(node)
	}

	if len(parts) == 1 {
		node.AppendMiddleware(m)
	} else {
		node.WithChildren(node.Tree().WithMiddleware(strings.Join(parts[1:], "/"), m, node.MaxParamsSize()))
	}

	return newTree
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
	newTree := t

	if node == nil {
		node = NewNode(parts[0], maxParamsSize)
		if len(parts) == 1 {
			node = withSubrouter(node)
		}
		newTree = t.withNode(node)
	}

	if len(parts) == 1 {
		node.WithRoute(route)
	} else {
		node.WithChildren(node.Tree().WithSubrouter(strings.Join(parts[1:], "/"), route, node.MaxParamsSize()))
	}

	return newTree
}

// withNode inserts node to Tree
// Nodes are sorted static, regexp, wildcard
func (t Tree) withNode(node Node) Tree {
	if node == nil {
		return t
	}

	newTree := append(t, node)

	// Sort Nodes in order [statics, regexps, wildcards]
	sort.Slice(newTree, func(i, j int) bool {
		return isMoreImportant(newTree[i], newTree[j])
	})

	return newTree
}

func isMoreImportant(left Node, right Node) bool {
	if leftNode, ok := left.(*subrouterNode); ok {
		return isMoreImportant(leftNode.Node, right)
	}

	if rightNode, ok := right.(*subrouterNode); ok {
		return isMoreImportant(left, rightNode.Node)
	}

	switch leftNode := left.(type) {
	case *staticNode:
		if rightNode, ok := right.(*staticNode); ok {
			return len(leftNode.name) < len(rightNode.name)
		}
		return true
	case *regexpNode:
		if _, ok := right.(*wildcardNode); ok {
			return true
		}
		if rightNode, ok := right.(*regexpNode); ok {
			return len(leftNode.regexp.String()) < len(rightNode.regexp.String())
		}
		return false
		// case *wildcardNode:
	}

	return false
}
