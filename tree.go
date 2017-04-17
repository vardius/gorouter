package goserver

import (
	"regexp"
	"strings"
	"sync"
)

type (
	node struct {
		path     string
		regexp   *regexp.Regexp
		route    *route
		parent   *node
		children tree
		nodesMu  sync.RWMutex
	}
	tree map[string]*node
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

func (n *node) child(paths []string) (*node, map[string]string) {
	if len(paths) > 0 && paths[0] != "" {
		n.nodesMu.RLock()
		defer n.nodesMu.RUnlock()
		if child := n.children[paths[0]]; child != nil {
			return child.child(paths[1:])
		} else {
			for path, child := range n.children {
				if len(path) > 0 && path[:1] == ":" {
					if child.regexp == nil || child.regexp.MatchString(paths[0]) {
						node, params := child.child(paths[1:])
						params[strings.Split(path, ":")[1]] = paths[0]
						return node, params
					}
				}
			}
		}
	} else if len(paths) == 0 {
		return n, make(map[string]string)
	}
	return nil, make(map[string]string)
}

func (n *node) addChild(paths []string) *node {
	if len(paths) > 0 && paths[0] != "" {
		n.nodesMu.Lock()
		defer n.nodesMu.Unlock()
		if n.children[paths[0]] == nil {
			child := newNode(n, paths[0])
			n.children[paths[0]] = child
			return child
		}
		return n.children[paths[0]].addChild(paths[1:])
	} else {
		if len(n.path) > 0 && n.path[:1] == ":" {
			if parts := strings.Split(n.path, ":"); len(parts) == 3 {
				n.setRegexp(parts[2])
			}
		}
	}

	return n
}

func (n *node) setRoute(r *route) {
	n.route = r
}

func newNode(root *node, path string) *node {
	return &node{
		path:     path,
		children: make(tree),
	}
}

func newRootNode(path string) *node {
	return newNode(nil, path)
}
