package goserver

import (
	"regexp"
	"strings"
)

type (
	node struct {
		path     string
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
			if child.path == paths[0] {
				cn = child
				break
			}
		}
		if cn == nil {
			cn = newNode(n, paths[0])
			n.children = append(n.children, cn)
		}
		return cn.addChild(paths[1:])
	}
	return n
}

func (n *node) child(paths []string) (*node, Params) {
	return n.childRecursive(paths)
}

//childRecursive returns child node by path using recurency
func (n *node) childRecursive(paths []string) (*node, Params) {
	if len(paths) > 0 && paths[0] != "" {
		path := paths[0]
		for _, child := range n.children {
			if len(child.path) > 0 && child.path[:1] == ":" {
				if child.regexp != nil && !child.regexp.MatchString(path) {
					continue
				}
				node, params := child.child(paths[1:])
				if node == nil {
					continue
				}
				if child.regexp != nil && child.regexp.MatchString(path) {
					for i := 1; i < len(child.path); i++ {
						if child.path[i] == ':' {
							params[child.params-1].Key = child.path[1:i]
							break
						}
					}
				} else {
					params[child.params-1].Key = child.path[1:]
				}
				params[child.params-1].Value = path

				return node, params
			} else if child.path == path {
				return child.child(paths[1:])
			}
		}
	} else if len(paths) == 0 {
		return n, make(Params, n.params)
	}
	return nil, nil
}

//childNotRecursive returns child node by path not using recurency
func (n *node) childNotRecursive(paths []string) (*node, Params) {
	var params Params
st:
	for {
		if len(paths) > 0 && paths[0] != "" {
			path := paths[0]
			for _, child := range n.children {
				if len(child.path) > 0 && child.path[:1] == ":" {
					if child.regexp != nil && !child.regexp.MatchString(path) {
						continue
					}
					if len(params) == 0 {
						params = make(Params, len(paths))
					}
					if child.regexp != nil && child.regexp.MatchString(path) {
						for i := 1; i < len(child.path); i++ {
							if child.path[i] == ':' {
								params[child.params-1].Key = child.path[1:i]
								break
							}
						}
					} else {
						params[child.params-1].Key = child.path[1:]
					}
					params[child.params-1].Value = path
					n = child
					paths = paths[1:]
					continue st
				} else if child.path == path {
					n = child
					paths = paths[1:]
					continue st
				}
			}
			return nil, nil
		} else if len(paths) == 0 {
			return n, params
		}
		return nil, nil
	}
}

func newNode(root *node, path string) *node {
	n := &node{
		path:     path,
		parent:   root,
		children: make([]*node, 0),
	}

	if root != nil {
		n.params = root.params
	}

	if len(n.path) > 0 && path[:1] == ":" {
		n.params++
		if parts := strings.Split(n.path, ":"); len(parts) == 3 {
			n.setRegexp(parts[2])
		}
	}

	return n
}

func newRoot(path string) *node {
	return newNode(nil, path)
}
