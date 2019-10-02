package gorouter

import (
	"strings"

	"github.com/vardius/gorouter/v4/middleware"
	"github.com/vardius/gorouter/v4/mux"
	pathutils "github.com/vardius/gorouter/v4/path"
)

func addNode(t *mux.Tree, method, path string) *mux.Node {
	root := t.Find(method)
	if root == nil {
		root = mux.NewNode(method, nil)
		t.Insert(root)
	}

	path = pathutils.TrimSlash(path)
	parts := strings.Split(path, "/")

	n := root.AddChild(parts)

	return n
}

func addMiddleware(t *mux.Tree, method, path string, mid middleware.Middleware) {
	type recFunc func(recFunc, *mux.Node, middleware.Middleware)

	c := func(c recFunc, n *mux.Node, m middleware.Middleware) {
		if n.Route() != nil {
			n.Route().AppendMiddleware(m)
		}
		for _, child := range n.Tree().StaticNodes() {
			c(c, child, m)
		}
		for _, child := range n.Tree().RegexpNodes() {
			c(c, child, m)
		}
		if n.Tree().WildcardNode() != nil {
			c(c, n.Tree().WildcardNode(), m)
		}
	}

	// routes tree roots should be http method nodes only
	for _, root := range t.StaticNodes() {
		if method == "" || method == root.ID() {
			if path != "" {
				node, _, _ := root.FindByPath(path)
				if node != nil {
					c(c, node, mid)
				}
			} else {
				c(c, root, mid)
			}
		}
	}
}

func allowed(t *mux.Tree, method, path string) (allow string) {
	if path == "*" {
		// routes tree roots should be http method nodes only
		for _, root := range t.StaticNodes() {
			if root.ID() == OPTIONS {
				continue
			}
			if len(allow) == 0 {
				allow = root.ID()
			} else {
				allow += ", " + root.ID()
			}
		}
	} else {
		// routes tree roots should be http method nodes only
		for _, root := range t.StaticNodes() {
			if root.ID() == method || root.ID() == OPTIONS {
				continue
			}

			n, _, _ := root.FindByPath(path)
			if n != nil && n.Route() != nil {
				if len(allow) == 0 {
					allow = root.ID()
				} else {
					allow += ", " + root.ID()
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", " + OPTIONS
	}
	return allow
}
