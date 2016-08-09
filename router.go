package goserver

import (
	"regexp"
	"strings"
	"sync"
)

type (
	tree       map[string]*route
	parameters map[string]string
	route      struct {
		path       string
		regexp     *regexp.Regexp
		root       *route
		nodes      tree
		middleware middlewares
		handler    HandlerFunc
		isEndPoint bool
		nodesMu    sync.RWMutex
	}
	Route interface {
		Path() string
		Regexp() string
		IsRoot() bool
		Parent() Route
		Nodes() map[string]Route
		IsEndPoint() bool
	}
)

func (r *route) Path() string {
	return r.path
}

func (r *route) Regexp() string {
	if r.regexp == nil {
		return ""
	}
	return r.regexp.String()
}

func (r *route) Parent() Route {
	return r.root
}

func (r *route) IsRoot() bool {
	return r.root == nil
}

func (r *route) Nodes() map[string]Route {
	newMap := make(map[string]Route)
	for path, route := range r.nodes {
		newMap[path] = route
	}
	return newMap
}

func (r *route) IsEndPoint() bool {
	return r.isEndPoint
}

func (r *route) getRoute(paths []string) (*route, parameters) {
	if len(paths) > 0 && paths[0] != "" {
		r.nodesMu.RLock()
		defer r.nodesMu.RUnlock()
		if route := r.nodes[paths[0]]; route != nil {
			return route.getRoute(paths[1:])
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
		return r, make(parameters)
	}
	return nil, make(parameters)
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

func newRoute(root *route, path string) *route {
	return &route{
		root:  root,
		path:  path,
		nodes: make(tree),
	}
}
