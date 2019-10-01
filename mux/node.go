package mux

import (
	"regexp"

	"github.com/vardius/gorouter/v4/context"
	path_utils "github.com/vardius/gorouter/v4/path"
)

// NewNode provides new node
func NewNode(pathPart string, parent *Node) *Node {
	if len(pathPart) == 0 {
		return nil
	}

	if parent != nil && parent.isSubrouter {
		panic("Subrouter node can not be a parent")
	}

	var compiledRegexp *regexp.Regexp
	isWildcard := false
	isRegexp := false

	id, exp := path_utils.GetIDFromPart(pathPart)
	if exp != "" {
		compiledRegexp = regexp.MustCompile(exp)
		isRegexp = true
	} else if id != pathPart {
		isWildcard = true
	}

	node := &Node{
		id:         id,
		parent:     parent,
		children:   NewTree(),
		regexp:     compiledRegexp,
		isWildcard: isWildcard && !isRegexp,
		isRegexp:   isRegexp,
	}

	if parent != nil {
		node.maxParamsSize = parent.maxParamsSize
		parent.children.Insert(node)
	}

	if node.isRegexp || node.isWildcard {
		node.maxParamsSize++
	}

	return node
}

// Node is route node
type Node struct {
	id     string
	regexp *regexp.Regexp

	parent   *Node
	children *Tree

	maxParamsSize uint8

	isWildcard  bool
	isRegexp    bool
	isSubrouter bool

	route Route
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) Route() Route {
	return n.route
}

func (n *Node) Tree() *Tree {
	return n.children
}

func (n *Node) SetRoute(r Route) {
	n.route = r
}

func (n *Node) TurnIntoSubrouter() {
	n.isSubrouter = true
}

// AddChild adds child by ids
func (n *Node) AddChild(parts []string) *Node {
	if len(parts) == 0 {
		return n
	}

	id, _ := path_utils.GetIDFromPart(parts[0])

	node := n.children.GetByID(id)
	if node == nil {
		node = NewNode(parts[0], n)
	}

	return node.AddChild(parts[1:])
}

// GetByIDs gets node by IDs
// this method is used when inserting new nodes
func (n *Node) GetByIDs(ids []string) *Node {
	if len(ids) == 0 {
		return nil
	}

	node := n.children.GetByID(ids[0])

	if node != nil {
		return node.GetByIDs(ids[1:])
	}

	return n
}

func (n *Node) FindByPath(path string) (*Node, context.Params, string) {
	var params context.Params

	pathPart, path := path_utils.GetPart(path)
	node := n.children.Find(pathPart)

	if node != nil {
		if node.isSubrouter {
			return node, nil, path
		}

		node, params, path = node.FindByPath(path)
	} else {
		node = n
		params = make(context.Params, n.maxParamsSize)
	}

	if node.isRegexp || node.isWildcard {
		params[node.maxParamsSize-1].Value = pathPart
		params[node.maxParamsSize-1].Key = node.id
	}

	return node, params, path
}
