package mux

import (
	"github.com/vardius/gorouter/v4/context"
	path_utils "github.com/vardius/gorouter/v4/path"
)

// Tree of routes
type Tree struct {
	wildcard *Node
	statics  *NodeMap
	regexps  []*Node
}

// NewTree provides new node tree
func NewTree() *Tree {
	return &Tree{
		statics: NewNodeMap(),
		regexps: make([]*Node, 0),
	}
}

func (t *Tree) Insert(node *Node) {
	if node == nil {
		return
	}

	if t.wildcard != nil {
		panic("Tree already contains a wildcard node!")
	}

	if node.isRegexp {
		t.regexps = append(t.regexps, node)
	} else if node.isWildcard {
		t.wildcard = node

		// wildcard node will match every case, reset other collections
		t.statics = NewNodeMap()
		t.regexps = make([]*Node, 0)
	} else {
		t.statics.Append(node)
	}
}

// Find finds node by ID
func (t *Tree) Find(id string) *Node {
	if id == "" {
		return nil
	}

	if t.wildcard != nil {
		return t.wildcard
	}

	staticNode := t.statics.Find(id)
	if staticNode != nil {
		return staticNode
	}

	for _, regexpNode := range t.regexps {
		if regexpNode.regexp.MatchString(id) {
			return regexpNode
		}
	}

	return nil
}

func (t *Tree) FindByPath(path string) (*Node, context.Params, string) {
	if path == "" || path == "/" {
		return nil, nil, ""
	}

	path = path_utils.TrimSlash(path)
	id, path := path_utils.GetPart(path)

	node := t.Find(id)

	if node == nil {
		return nil, nil, path
	}

	if node.isSubrouter {
		return node, nil, path
	}

	nextNode, params, subPath := node.children.FindByPath(path)

	if params == nil {
		params = make(context.Params, node.maxParamsSize)
	}

	if node.isRegexp || node.isWildcard {
		params[node.maxParamsSize-1].Value = id
		params[node.maxParamsSize-1].Key = node.id
	}

	if nextNode != nil {
		node = nextNode
	}

	return node, params, subPath
}

func (t *Tree) StaticNodes() []*Node {
	return t.statics.nodes
}

func (t *Tree) RegexpNodes() []*Node {
	return t.regexps
}

func (t *Tree) WildcardNode() *Node {
	return t.wildcard
}
