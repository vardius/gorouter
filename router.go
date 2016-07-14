package goapi

import (
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type (
	tree  map[string]*route
	route struct {
		path       string
		nodes      tree
		nodesMu    sync.RWMutex
		handler    http.HandlerFunc
		middleware middlewares
		regexp     *regexp.Regexp
		isEndPoint bool
	}
)

func (r *route) getRoute(paths []string) *route {
	if len(paths) > 0 && paths[0] != "" {
		r.nodesMu.RLock()
		defer r.nodesMu.RUnlock()
		if route := r.nodes[paths[0]]; route != nil {
			return route.getRoute(paths[1:])
		} else {
			for path, route := range r.nodes {
				if len(path) > 0 && path[:1] == ":" {
					if route.regexp == nil {
						return route.getRoute(paths[1:])
					} else if route.regexp.MatchString(paths[0]) {
						return route.getRoute(paths[1:])
					}
				}
			}
		}
	} else if len(paths) == 0 && r.isEndPoint {
		return r
	}
	return nil
}

func (r *route) addRoute(paths []string, f http.HandlerFunc) {
	if len(paths) > 0 && paths[0] != "" {
		r.nodesMu.Lock()
		defer r.nodesMu.Unlock()
		if r.nodes[paths[0]] == nil {
			r.nodes[paths[0]] = newRoute(paths[0])
		}
		r.nodes[paths[0]].addRoute(paths[1:], f)
	} else {
		r.setEndPoint(f)
	}
}

func (r *route) setEndPoint(f http.HandlerFunc) {
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

func newRoute(path string) *route {
	return &route{
		path:  path,
		nodes: make(tree),
	}
}
