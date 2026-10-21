// Package tree is the routing core shared by the net/http and fasthttp routers.
// It is generic over the handler type, matches without allocating, and composes
// each route's middleware chain when routes or middleware are registered, never per request.
package tree

import (
	"cmp"
	"fmt"
	"slices"
)

// standardMethods is also the order methods are listed in an Allow header.
var standardMethods = [...]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

func methodIndex(method string) int {
	switch method {
	case "GET":
		return 0
	case "HEAD":
		return 1
	case "POST":
		return 2
	case "PUT":
		return 3
	case "PATCH":
		return 4
	case "DELETE":
		return 5
	case "CONNECT":
		return 6
	case "OPTIONS":
		return 7
	case "TRACE":
		return 8
	}

	return -1
}

// Route is a registered handler.
type Route[H any] struct {
	Method  string
	Pattern string
	Params  []string // param names, in the order their values are written to Lookup's values
	Handler H        // the registered handler wrapped in the middleware that applies to it

	handler  H
	segments []segment
}

type middleware[H any] struct {
	method   string // empty for every method
	segments []segment
	wrap     []func(H) H
}

// appliesTo reports whether m was registered on r's method and on a pattern that r's pattern starts with.
func (m middleware[H]) appliesTo(r *Route[H]) bool {
	if m.method != "" && m.method != r.Method {
		return false
	}
	if len(m.segments) > len(r.segments) {
		return false
	}
	for i, s := range m.segments {
		if !s.sameNode(r.segments[i]) {
			return false
		}
	}

	return true
}

// Tree holds one routing tree per method. Register everything before serving: it is not safe for concurrent writes.
type Tree[H any] struct {
	roots      [len(standardMethods)]*node[H]
	custom     map[string]*node[H]
	methods    []string // methods that have routes: standard ones in standardMethods order, then custom ones as registered
	routes     []*Route[H]
	middleware []middleware[H]
	maxParams  int
}

// Handle registers h for method and pattern. It panics if the pattern is invalid or taken.
func (t *Tree[H]) Handle(method, pattern string, h H) {
	if method == "" {
		panic(fmt.Sprintf("gorouter: empty method for %s", pattern))
	}

	t.add(method, pattern, h, false)
}

// Mount registers h for pattern and every path below it, on every standard method.
func (t *Tree[H]) Mount(pattern string, h H) {
	for _, method := range standardMethods {
		t.add(method, pattern, h, true)
	}
}

// Use wraps every route on method (every method when empty) whose pattern starts with pattern.
// Routes registered before and after the call are both covered.
func (t *Tree[H]) Use(method, pattern string, wrap []func(H) H) {
	if len(wrap) == 0 {
		return
	}

	t.middleware = append(t.middleware, middleware[H]{
		method:   method,
		segments: mustParse(pattern),
		wrap:     slices.Clone(wrap),
	})

	for _, r := range t.routes {
		t.compose(r)
	}
}

// Lookup returns the route for method and path, or nil.
// It writes param values to values, which must hold at least MaxParams strings, in the order of the route's Params.
func (t *Tree[H]) Lookup(method, path string, values []string) *Route[H] {
	root := t.root(method)
	if root == nil {
		return nil
	}

	path = trim(path)

	return root.descend(path, path != "", values)
}

// Allowed returns the methods other than method that have a route for path, formatted for an Allow header,
// or "" when there are none. The path "*" asks for every method with a route. values is scratch space as in Lookup.
func (t *Tree[H]) Allowed(method, path string, values []string) string {
	buf := make([]byte, 0, 64)
	trimmed := trim(path)

	for _, m := range t.methods {
		if m == method || m == "OPTIONS" {
			continue
		}
		if path == "*" || t.root(m).descend(trimmed, trimmed != "", values) != nil {
			buf = appendMethod(buf, m)
		}
	}

	if len(buf) == 0 {
		return ""
	}

	return string(appendMethod(buf, "OPTIONS"))
}

// MaxParams is the largest number of params any route has.
func (t *Tree[H]) MaxParams() int {
	return t.maxParams
}

// Routes returns the registered routes in registration order. The slice is shared: don't modify it.
func (t *Tree[H]) Routes() []*Route[H] {
	return t.routes
}

func (t *Tree[H]) add(method, pattern string, h H, prefix bool) {
	segments := mustParse(pattern)

	n := t.rootFor(method)
	for _, s := range segments {
		if n.prefix {
			panic(fmt.Sprintf("gorouter: %s %s conflicts with the mount at %s", method, pattern, n.route.Pattern))
		}
		n = n.child(s)
	}

	if n.route != nil {
		panic(fmt.Sprintf("gorouter: %s %s conflicts with %s %s", method, pattern, n.route.Method, n.route.Pattern))
	}
	if prefix && !n.isLeaf() {
		panic(fmt.Sprintf("gorouter: mount at %s conflicts with %s routes registered below it", pattern, method))
	}

	r := &Route[H]{Method: method, Pattern: pattern, handler: h, segments: segments}
	for _, s := range segments {
		if s.kind != static {
			r.Params = append(r.Params, s.name)
		}
	}

	t.compose(r)
	n.route, n.prefix = r, prefix
	t.routes = append(t.routes, r)
	t.maxParams = max(t.maxParams, len(r.Params))
}

// compose wraps r's handler in the middleware that applies to it: earlier registrations end up outermost.
func (t *Tree[H]) compose(r *Route[H]) {
	h := r.handler
	for _, m := range slices.Backward(t.middleware) {
		if !m.appliesTo(r) {
			continue
		}
		for _, wrap := range slices.Backward(m.wrap) {
			h = wrap(h)
		}
	}

	r.Handler = h
}

func (t *Tree[H]) root(method string) *node[H] {
	if i := methodIndex(method); i >= 0 {
		return t.roots[i]
	}

	return t.custom[method]
}

func (t *Tree[H]) rootFor(method string) *node[H] {
	if n := t.root(method); n != nil {
		return n
	}

	n := &node[H]{}
	if i := methodIndex(method); i >= 0 {
		t.roots[i] = n
	} else {
		if t.custom == nil {
			t.custom = make(map[string]*node[H])
		}
		t.custom[method] = n
	}

	t.methods = append(t.methods, method)
	slices.SortStableFunc(t.methods, func(a, b string) int {
		return cmp.Compare(allowRank(a), allowRank(b))
	})

	return n
}

// allowRank orders standard methods as in standardMethods and custom ones after them.
func allowRank(method string) int {
	if i := methodIndex(method); i >= 0 {
		return i
	}

	return len(standardMethods)
}

func appendMethod(buf []byte, method string) []byte {
	if len(buf) > 0 {
		buf = append(buf, ", "...)
	}

	return append(buf, method...)
}
