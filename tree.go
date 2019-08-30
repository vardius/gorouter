package gorouter

import (
	"strings"

	"github.com/vardius/gorouter/v4/middleware"
	"github.com/vardius/gorouter/v4/mux"
	path_utils "github.com/vardius/gorouter/v4/path"
)

func addNode(t *mux.Tree, id, path string) *mux.Node {
	root := t.GetByID(id)
	if root == nil {
		root = mux.NewRoot(id)
		t.Insert(root)
	}

	paths := path_utils.Split(path)
	n := root.AddChild(paths)

	return n
}

func addMiddleware(t *mux.Tree, method, path string, mid middleware.Middleware) {
	type recFunc func(recFunc, *mux.Node, middleware.Middleware)

	c := func(c recFunc, n *mux.Node, m middleware.Middleware) {
		if n.Route() != nil {
			n.Route().AppendMiddleware(m)
		}
		for _, child := range n.Children().StaticNodes() {
			c(c, child, m)
		}
		for _, child := range n.Children().RegexpNodes() {
			c(c, child, m)
		}
		if n.Children().WildcardNode() != nil {
			c(c, n.Children().WildcardNode(), m)
		}
	}

	paths := path_utils.Split(path)

	// routes tree roots should be http method nodes only
	for _, root := range t.StaticNodes() {
		if method == "" || method == root.ID() {
			node, _ := root.GetChild(paths)
			if node != nil {
				c(c, node, mid)
			}
		}
	}
}

func allowed(t *mux.Tree, method, path string) (allow string) {
	path = strings.Trim(path, "/")

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

			n, _, _ := root.GetChildByPath(path)
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
