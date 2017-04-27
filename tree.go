package goserver

import (
	"regexp"
	"strings"
)

type (
	node struct {
		id         string
		regexp     *regexp.Regexp
		route      *route
		parent     *node
		children   *tree
		params     uint8
		isWildcard bool
		isRegexp   bool
	}
	tree struct {
		idsLen   int
		ids      []string
		statics  map[int]*node
		regexps  []*node
		wildcard *node
	}
)

func (n *node) isRoot() bool {
	return n.parent == nil
}

func (n *node) isLeaf() bool {
	return n.children.idsLen == 0 && len(n.children.regexps) == 0 && n.children.wildcard == nil
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
		node := n.children.byID(ids[0])

		if node == nil {
			node = newNode(n, ids[0])
			n.children.insert(node)
		}

		return node.addChild(ids[1:])
	}
	return n
}

func (n *node) child(ids []string) (*node, Params) {
	if len(ids) == 0 {
		return n, make(Params, n.params)
	}

	child := n.children.byID(ids[0])
	if child != nil {
		n, params := child.child(ids[1:])

		if child.isWildcard && params != nil {
			params[child.params-1].Value = ids[0]
			params[child.params-1].Key = child.id
		}

		return n, params
	}

	return nil, nil
}

func (n *node) childByPath(path string) (*node, Params) {
	pathLen := len(path)
	if pathLen > 0 && path[0] == '/' {
		path = path[1:]
		pathLen--
	}

	if pathLen == 0 {
		return n, make(Params, n.params)
	}

	child, part, path := n.children.byPath(path)

	if child != nil {
		n, params := child.childByPath(path)

		if part != "" && params != nil {
			params[child.params-1].Value = part
			params[child.params-1].Key = child.id
		}

		return n, params
	}

	return nil, nil
}

func newNode(root *node, id string) *node {
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

	n := &node{
		id:         id,
		parent:     root,
		children:   newTree(),
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
		n.setRegexp(regexp)
	}

	return n
}

func (t *tree) insert(n *node) {
	if n == nil {
		return
	}

	if n.isWildcard {
		if n.isRegexp {
			t.regexps = append(t.regexps, n)
		} else {
			if t.wildcard != nil {
				panic("Tree already contains a wildcard child!")
			}
			t.wildcard = n
		}
	} else {
		index := -1
		for i, id := range t.ids {
			if n.id > id {
				index = i
				break
			}
		}

		if index > -1 {
			t.ids = append(t.ids[:index], append([]string{n.id}, t.ids[index:]...)...)
			for i := t.idsLen - 1; i >= 0; i-- {
				if i < index {
					break
				} else {
					t.statics[i+1] = t.statics[i]
				}
			}
			t.statics[index] = n
		} else {
			t.ids = append(t.ids, n.id)
			t.statics[t.idsLen] = n
		}

		t.idsLen++
	}
}

func (t *tree) atIndex(i int) *node {
	if t.idsLen > i {
		return t.byID(t.ids[i])
	}
	return nil
}

func (t *tree) byID(id string) *node {
	if id != "" {
		if t.idsLen > 0 {
			for i, cID := range t.ids {
				if cID == id {
					return t.statics[i]
				}
			}
		}

		for _, child := range t.regexps {
			if child.regexp.MatchString(id) {
				return child
			}
		}

		return t.wildcard
	}

	return nil
}

func (t *tree) byPath(path string) (*node, string, string) {
	if len(path) == 0 {
		return nil, "", ""
	}

	if t.idsLen > 0 {
		for i, cID := range t.ids {
			pLen := len(cID)
			if len(path) >= pLen && cID == path[:pLen] {
				return t.statics[i], "", path[pLen:]
			}
		}
	}

	part := path
	if j := strings.IndexByte(path, '/'); j > 0 {
		part = path[:j]
	}

	for _, child := range t.regexps {
		if child.regexp.MatchString(part) {
			return child, part, path[len(part):]
		}
	}

	if t.wildcard != nil {
		return t.wildcard, part, path[len(part):]
	}

	return nil, "", ""
}

func newRoot(id string) *node {
	return newNode(nil, id)
}

func newTree() *tree {
	return &tree{
		ids:     make([]string, 0),
		statics: make(map[int]*node),
		regexps: make([]*node, 0),
	}
}
