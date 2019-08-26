package mux

import (
	"regexp"
	"strings"

	"github.com/vardius/gorouter/v4/context"
)

// Node is route node
type Node struct {
	id          string
	regexp      *regexp.Regexp
	route       interface{}
	parent      *Node
	children    *Tree
	params      uint8
	isWildcard  bool
	isRegexp    bool
	isSubrouter bool
}

// NewRoot create new tree root
func NewRoot(id string) *Node {
	return NewNode(nil, id)
}

// NewNode creates new tre node
func NewNode(root *Node, id string) *Node {
	var regexp string
	isWildcard := false
	isRegexp := false

	if len(id) > 0 && id[0] == '{' {
		id = id[1 : len(id)-1]
		isWildcard = true

		if parts := strings.Split(id, ":"); len(parts) == 2 {
			id = parts[0]
			regexp = parts[1]
			regexp = regexp[:len(regexp)-1]
			isRegexp = true
		}

		if id == "" {
			panic("Empty wildcard name")
		}
	}

	n := &Node{
		id:         id,
		parent:     root,
		children:   NewTree(),
		isWildcard: isWildcard,
		isRegexp:   isRegexp,
	}

	if root != nil {
		n.params = root.params
	}

	if isWildcard {
		n.params++
	}

	if isRegexp {
		n.SetRegexp(regexp)
	}

	return n
}

// AddChild adds child by ids
func (n *Node) AddChild(ids []string) *Node {
	if len(ids) > 0 && ids[0] != "" {
		node := n.children.GetByID(ids[0])

		if node == nil {
			node = NewNode(n, ids[0])
			n.children.Insert(node)
		}

		return node.AddChild(ids[1:])
	}
	return n
}

// GetChild gets child by ids
func (n *Node) GetChild(ids []string) (*Node, context.Params) {
	if len(ids) == 0 {
		return n, make(context.Params, n.params)
	}

	child := n.children.GetByID(ids[0])
	if child != nil {
		n, params := child.GetChild(ids[1:])

		if child.isWildcard && params != nil {
			params[child.params-1].Value = ids[0]
			params[child.params-1].Key = child.id
		}

		return n, params
	}

	return nil, nil
}

// GetChildByPath accepts string path then returns:
// child node as a first arg,
// parameters built from wildcards,
// and part of path (this is used to strip request path for sub routers)
func (n *Node) GetChildByPath(path string) (*Node, context.Params, string) {
	pathLen := len(path)
	if pathLen > 0 && path[0] == '/' {
		path = path[1:]
		pathLen--
	}

	if pathLen == 0 {
		return n, make(context.Params, n.params), ""
	}

	child, part, path := n.children.GetByPath(path)

	if child != nil {
		grandChild, params, _ := child.GetChildByPath(path)

		if part != "" && params != nil {
			params[child.params-1].Value = part
			params[child.params-1].Key = child.id
		}

		if grandChild == nil && child.isSubrouter {
			return child, params, path
		}

		return grandChild, params, ""
	}

	return nil, nil, ""
}

// SetRegexp sets node regexp value
func (n *Node) SetRegexp(exp string) {
	reg, err := regexp.Compile(exp)
	if err != nil {
		panic(err)
	}

	n.regexp = reg
	n.isRegexp = true
	n.isWildcard = true
}

// SetRoute sets node route value
func (n *Node) SetRoute(r interface{}) {
	n.route = r
}

// TurnIntoSubrouter sets node as subrouter
func (n *Node) TurnIntoSubrouter() {
	n.isSubrouter = true
}

func (n *Node) regexpToString() string {
	if n.regexp == nil {
		return ""
	}
	return n.regexp.String()
}

// ID returns node's id
func (n *Node) ID() string {
	return n.id
}

// Route returns node's route
func (n *Node) Route() interface{} {
	return n.route
}

// Children returns node's children as tree
func (n *Node) Children() *Tree {
	return n.children
}

// IsSubrouter returns true if node is subrouter
func (n *Node) IsSubrouter() bool {
	return n.isSubrouter
}

// IsRoot returns true if node is root
func (n *Node) IsRoot() bool {
	return n.parent == nil
}

// IsLeaf returns true if node is root
func (n *Node) IsLeaf() bool {
	return n.children.idsLen == 0 && len(n.children.regexps) == 0 && n.children.wildcard == nil
}
