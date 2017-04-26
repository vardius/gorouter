package goserver

import (
	"regexp"
	"strings"
)

type (
	node struct {
		id         string
		ids        []string
		parent     *node
		children   map[int]*node
		length     int
		regexp     *regexp.Regexp
		route      *route
		params     uint8
		isWildcard bool
	}
)

func (n *node) isRoot() bool {
	return n.parent == nil
}

func (n *node) isLeaf() bool {
	return len(n.children) == 0
}

func (n *node) regexpToString() string {
	if n.regexp == nil {
		return ""
	}
	return n.regexp.String()
}

func (n *node) setRegexp(exp string) {
	reg, err := regexp.Compile(exp)
	if err == nil {
		n.regexp = reg
	}
}

func (n *node) setRoute(r *route) {
	n.route = r
}

func (n *node) addChild(ids []string) *node {
	if len(ids) > 0 && ids[0] != "" {
		var node *node
		for i, id := range n.ids {
			if id == ids[0] {
				node = n.children[i]
				break
			}
		}

		if node == nil {
			node = newNode(n, ids[0])
			n.ids = append(n.ids, ids[0])
			n.children[n.length] = node
			n.length++
		}

		return node.addChild(ids[1:])
	}
	return n
}

func (n *node) childById(id string) *node {
	if id != "" && n.length > 0 {
		for i, cId := range n.ids {
			child := n.children[i]
			if cId == id {
				return child
			}

			if child.isWildcard {
				if child.regexp != nil && !child.regexp.MatchString(id) {
					continue
				}
				return child
			}
		}
	}

	return nil
}

func (n *node) childAtIndex(i int) *node {
	if i > -1 && n.length > i {
		return n.childById(n.ids[i])
	}
	return nil
}

func (n *node) child(ids []string) (*node, Params) {
	idsLen := len(ids)
	if idsLen == 0 {
		return n, make(Params, n.params)
	}

	if idsLen > 0 && n.length > 0 {
		id := ids[0]
		child := n.childById(id)
		if child != nil {
			n, params := child.child(ids[1:])

			if child.isWildcard {
				params[child.params-1].Value = id
				params[child.params-1].Key = child.id
			}

			return n, params
		}
	}

	return nil, nil
}

func newNode(root *node, id string) *node {
	var regexp string
	isWildcard := false

	if id[0] == '{' {
		id = id[1 : len(id)-1]
		isWildcard = true

		if parts := strings.Split(id, ":"); len(parts) == 2 {
			id = parts[0]
			regexp = parts[1]
			regexp = regexp[:len(regexp)-1]
		}
	}

	n := &node{
		id:         id,
		parent:     root,
		children:   make(map[int]*node),
		ids:        make([]string, 0),
		isWildcard: isWildcard,
	}

	if root != nil {
		n.params = root.params
	}

	if isWildcard {
		n.params++
		n.setRegexp(regexp)
	}

	return n
}

func newRoot(id string) *node {
	return newNode(nil, id)
}
