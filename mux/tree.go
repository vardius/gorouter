package mux

import (
	"github.com/vardius/gorouter/v4/context"
	pathutils "github.com/vardius/gorouter/v4/path"
)

func NewTree() Tree {
	return make([]Node, 0)
}

type Tree []Node

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

func (t Tree) Match(path string) (Node, context.Params, string) {
	pathPart, subPath := pathutils.GetPart(path)

	for _, child := range t {
		if node, params, subPath := child.Match(pathPart, subPath); node != nil {
			return node, params, subPath
		}
	}

	return nil, nil, path
}
