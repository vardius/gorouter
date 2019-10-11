package gorouter

import (
	"strings"

	"github.com/vardius/gorouter/v4/middleware"
	"github.com/vardius/gorouter/v4/mux"
	pathutils "github.com/vardius/gorouter/v4/path"
)

func addMiddleware(t mux.Tree, method, path string, mid middleware.Middleware) {
	type recFunc func(recFunc, mux.Node, middleware.Middleware)

	c := func(c recFunc, n mux.Node, m middleware.Middleware) {
		if n.Route() != nil {
			n.Route().AppendMiddleware(m)
		}
		for _, child := range n.Tree() {
			c(c, child, m)
		}
	}

	// routes tree roots should be http method nodes only
	if root := t.Find(method); root != nil {
		if path != "" {
			node := findNode(root, strings.Split(pathutils.TrimSlash(path), "/"))
			if node != nil {
				c(c, node, mid)
			}
		} else {
			c(c, root, mid)
		}
	}
}

func findNode(n mux.Node, parts []string) mux.Node {
	if len(parts) == 0 {
		return n
	}

	name, _ := pathutils.GetNameFromPart(parts[0])

	if node := n.Tree().Find(name); node != nil {
		return findNode(node, parts[1:])
	}

	return n
}

func allowed(t mux.Tree, method, path string) (allow string) {
	path = pathutils.TrimSlash(path)

	if path == "*" {
		// routes tree roots should be http method nodes only
		for _, root := range t {
			if root.Name() == OPTIONS {
				continue
			}
			if len(allow) == 0 {
				allow = root.Name()
			} else {
				allow += ", " + root.Name()
			}
		}
	} else {
		// routes tree roots should be http method nodes only
		for _, root := range t {
			if root.Name() == method || root.Name() == OPTIONS {
				continue
			}

			if n, _, _ := root.Tree().Match(path); n != nil && n.Route() != nil {
				if len(allow) == 0 {
					allow = root.Name()
				} else {
					allow += ", " + root.Name()
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", " + OPTIONS
	}
	return allow
}
