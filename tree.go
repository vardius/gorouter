package goserver

import (
	"regexp"
	"sort"
	"strings"
)

type (
	node struct {
		pattern  string
		regexp   *regexp.Regexp
		route    *route
		parent   *node
		children []*node
		params   uint8
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

func (n *node) addChild(paths []string) *node {
	if len(paths) > 0 && paths[0] != "" {
		var cn *node
		for _, child := range n.children {
			if child.pattern == paths[0] {
				cn = child
				break
			}
		}

		if cn == nil {
			cn = newNode(n, paths[0])
			n.children = append(n.children, cn)
			sort.Slice(n.children, func(i, j int) bool {
				return len(n.children[i].pattern) > len(n.children[j].pattern)
			})
		}

		return cn.addChild(paths[1:])
	}
	return n
}

func (n *node) child(paths []string) *node {
	pathsLen := len(paths)
	if pathsLen == 0 {
		return n
	}

	if pathsLen > 0 && paths[0] != "" {
		for _, child := range n.children {
			if child.pattern == paths[0] {
				return child.child(paths[1:])
			}
		}
	}

	return nil
}

func (n *node) childByPath(path string) (node *node, params Params) {
	pathLen := len(path)
	if pathLen > 0 && path[0] == '/' {
		path = path[1:]
		pathLen--
	}

	if pathLen == 0 {
		return n, make(Params, n.params)
	}

	for _, child := range n.children {
		pLen := len(child.pattern)
		if pathLen >= pLen && child.pattern == path[:pLen] {
			return child.childByPath(path[pLen:])
		}
		if pLen > 1 && child.pattern[0] == '{' {
			part := path
			if i := strings.IndexByte(path, '/'); i > 0 {
				part = path[:i]
			}

			if child.regexp != nil && !child.regexp.MatchString(part) {
				continue
			}

			node, params = child.childByPath(path[len(part):])
			if node == nil {
				return
			}

			params[child.params-1].Value = part
			if child.regexp != nil {
				if i := strings.IndexByte(child.pattern, ':'); i > 1 {
					params[child.params-1].Key = child.pattern[1:i]
				}
			} else {
				params[child.params-1].Key = child.pattern[1 : pLen-1]
			}

			return
		}
	}
	return
}

func newNode(root *node, pattern string) *node {
	n := &node{
		pattern:  pattern,
		parent:   root,
		children: make([]*node, 0),
	}

	if root != nil {
		n.params = root.params
	}

	if len(n.pattern) > 0 && pattern[0] == '{' {
		n.params++
		if parts := strings.Split(n.pattern, ":"); len(parts) == 2 {
			r := parts[1]
			n.setRegexp(r[:len(r)-1])
		}
	}

	return n
}

func newRoot(pattern string) *node {
	return newNode(nil, pattern)
}
