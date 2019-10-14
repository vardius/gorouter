package mux

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vardius/gorouter/v4/context"
	pathutils "github.com/vardius/gorouter/v4/path"
)

func NewTree() Tree {
	return make([]Node, 0)
}

type Tree []Node

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

func (t Tree) Compile() Tree {
	for i, child := range t {
		child.WithChildren(child.Tree().Compile())

		if len(child.Tree()) == 1 {
			t[i] = child.Merge(child.Tree()[0])
		}
	}

	return t
}

func (t Tree) Match(path string) (Node, context.Params, string) {
	for _, child := range t {
		if node, params, subPath := child.Match(path); node != nil {
			return node, params, subPath
		}
	}

	return nil, nil, ""
}

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

func (t Tree) WithNode(node Node) Tree {
	if node == nil {
		return t
	}

	if len(t) > 0 {
		switch (t)[0].(type) {
		case *wildcardNode:
			panic("Tree already contains wildcard node")
		}
	}

	switch n := node.(type) {
	case *wildcardNode:
		return append(append(NewTree(), n), t...)
	case *staticNode:
		return append(append(NewTree(), n), t...)
	case *regexpNode:
		return append(t, n)
	case *subrouterNode:
		return append(t, n)
	default:
		return append(t, n)
	}
}
