package mux

import (
	"regexp"
	"strings"
)

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

// NewNode provides new node
func NewNode(id string, parent *Node) *Node {
	if len(id) == 0 {
		return nil
	}

	if parent != nil && parent.isSubrouter {
		panic("Subrouter node can not be a parent")
	}

	var compiledRegexp *regexp.Regexp
	isWildcard := false
	isRegexp := false

	if id[0] == '{' {
		id = id[1 : len(id)-1]
		isWildcard = true

		if parts := strings.Split(id, ":"); len(parts) == 2 {
			id = parts[0]
			exp := parts[1]
			exp = exp[:len(exp)-1]

			compiledRegexp = regexp.MustCompile(exp)
			isRegexp = true
		}

		if id == "" {
			panic("Empty wildcard name")
		}
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
	}

	if node.isRegexp || node.isWildcard {
		node.maxParamsSize++
	}

	return node
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
