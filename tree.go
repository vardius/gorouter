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
				iLen := len(n.children[i].pattern)
				jLen := len(n.children[j].pattern)
				if iLen > 0 && jLen > 0 {
					isIParam := n.children[i].pattern[0] == ':'
					isJParam := n.children[j].pattern[0] == ':'
					if isIParam && isJParam {
						return strings.Contains(n.children[i].pattern[1:], ":")
					}
					return !isIParam
				}
				return n.children[i].pattern < n.children[j].pattern
			})
		}
		return cn.addChild(paths[1:])
	}
	return n
}

func (n *node) child(paths string) (*node, Params) {
	return n.childFromString(paths)
}

//childRecursive returns child node by path using recurency
func (n *node) childFromString(paths string) (node *node, params Params) {
	if paths != "" && paths[0] == '/' {
		paths = paths[1:]
	}
	pathsLen := len(paths)
	if pathsLen == 0 {
		return n, make(Params, n.params)
	}
	if pathsLen > 0 {
		for _, child := range n.children {
			pLen := len(child.pattern)
			if pathsLen >= pLen && child.pattern == paths[:pLen] {
				return child.child(paths[pLen:])
			}
			if pLen > 1 && child.pattern[0] == '{' {
				var path string
				for i := 1; i < pLen; i++ {
					if paths[i] == '/' {
						path = paths[:i]
						break
					}
				}
				if child.regexp != nil && !child.regexp.MatchString(path) {
					continue
				}
				node, params = child.child(paths[:len(path)])
				if node == nil {
					continue
				}
				if child.regexp != nil {
					for i := 1; i < pLen; i++ {
						if child.pattern[i] == ':' {
							params[child.params-1].Key = child.pattern[1:i]
							break
						}
					}
				} else {
					params[child.params-1].Key = child.pattern[1:pLen]
				}
				params[child.params-1].Value = path
				return
			}
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

	if len(n.pattern) > 0 && pattern[:1] == ":" {
		n.params++
		if parts := strings.Split(n.pattern, ":"); len(parts) == 3 {
			n.setRegexp(parts[2])
		}
	}

	return n
}

func newRoot(pattern string) *node {
	return newNode(nil, pattern)
}
