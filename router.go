package goapi

import (
	"regexp"
	"strings"
	"sync"
)

type (
	Params map[string]string
	tree   map[string]*route
	route  struct {
		path       string
		root       *route
		nodes      tree
		nodesMu    sync.RWMutex
		handler    HandlerFunc
		middleware middlewares
		regexp     *regexp.Regexp
		isEndPoint bool
	}
)

func (r *route) getRoute(paths []string) (*route, Params) {
	if len(paths) > 0 && paths[0] != "" {
		r.nodesMu.RLock()
		defer r.nodesMu.RUnlock()
		if route := r.nodes[paths[0]]; route != nil {
			node, params := route.getRoute(paths[1:])
			if len(route.path) > 0 && route.path[:1] == ":" {
				params[strings.Split(route.path, ":")[1]] = paths[0]
			}
			return node, params
		} else {
			for path, route := range r.nodes {
				if len(path) > 0 && path[:1] == ":" {
					if route.regexp == nil || route.regexp.MatchString(paths[0]) {
						node, params := route.getRoute(paths[1:])
						params[strings.Split(path, ":")[1]] = paths[0]
						return node, params
					}
				}
			}
		}
	} else if len(paths) == 0 && r.isEndPoint {
		return r, make(Params)
	}
	return nil, make(Params)
}

func (r *route) addRoute(paths []string, f HandlerFunc) {
	if len(paths) > 0 && paths[0] != "" {
		r.nodesMu.Lock()
		defer r.nodesMu.Unlock()
		if r.nodes[paths[0]] == nil {
			r.nodes[paths[0]] = newRoute(r, paths[0])
		}
		r.nodes[paths[0]].addRoute(paths[1:], f)
	} else {
		r.setEndPoint(f)
	}
}

func (r *route) setEndPoint(f HandlerFunc) {
	if len(r.path) > 0 && r.path[:1] == ":" {
		if parts := strings.Split(r.path, ":"); len(parts) == 3 {
			r.setRegexp(parts[2])
		}
	}
	r.isEndPoint = true
	r.handler = f
}

func (r *route) setRegexp(exp string) {
	reg, err := regexp.Compile(exp)
	if err == nil {
		r.regexp = reg
	}
}

func (r *route) isRoot() bool {
	return r.root == nil
}

func newRoute(root *route, path string) *route {
	return &route{
		root:  root,
		path:  path,
		nodes: make(tree),
	}
}
