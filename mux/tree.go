package mux

import (
	"strings"
)

// Tree of routes
type Tree struct {
	idsLen   int
	ids      []string
	statics  map[int]*Node
	regexps  []*Node
	wildcard *Node
}

// NewTree provides new node tree
func NewTree() *Tree {
	return &Tree{
		ids:     make([]string, 0),
		statics: make(map[int]*Node),
		regexps: make([]*Node, 0),
	}
}

// StaticNodes returns tree static nodes
func (t *Tree) StaticNodes() map[int]*Node {
	return t.statics
}

// RegexpNodes returns tree regexp nodes
func (t *Tree) RegexpNodes() []*Node {
	return t.regexps
}

// WildcardNode returns tree wildcard node
func (t *Tree) WildcardNode() *Node {
	return t.wildcard
}

// Insert inserts node
func (t *Tree) Insert(n *Node) {
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

// GetByID gets node by ID
func (t *Tree) GetByID(id string) *Node {
	if id != "" {
		if t.idsLen > 0 {
			for i, cID := range t.ids {
				if cID == id {
					return t.statics[i]
				}
			}
		}

		for _, child := range t.regexps {
			if child.regexp != nil && child.regexp.MatchString(id) {
				return child
			}
		}

		return t.wildcard
	}

	return nil
}

// GetByPath gets node by path
func (t *Tree) GetByPath(path string) (*Node, string, string) {
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
		if child.regexp != nil && child.regexp.MatchString(part) {
			return child, part, path[len(part):]
		}
	}

	if t.wildcard != nil {
		return t.wildcard, part, path[len(part):]
	}

	return nil, "", ""
}
