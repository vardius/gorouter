package mux

import (
	"strings"

	"github.com/vardius/gorouter/v4/context"
	pathutils "github.com/vardius/gorouter/v4/path"
)

func NewTree() Tree {
	return make([]Node, 0)
}

type Tree []Node

func (t Tree) Match(path string) (Node, context.Params, string) {
	pathPart, subPath := pathutils.GetPart(path)

	for _, child := range t {
		if node, params, subPath := child.Match(pathPart, subPath); node != nil {
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
			node = WithSubrouter(node)
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
