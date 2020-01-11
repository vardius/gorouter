package gorouter

import (
	"net/http"

	"github.com/vardius/gorouter/v4/mux"
	pathutils "github.com/vardius/gorouter/v4/path"
)

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
	if path == "*" {
		// routes tree roots should be http method nodes only
		for _, root := range t {
			if root.Name() == http.MethodOptions {
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
			if root.Name() == method || root.Name() == http.MethodOptions {
				continue
			}

			if n, _, _, _ := root.Tree().Match(path); n != nil && n.Route() != nil {
				if len(allow) == 0 {
					allow = root.Name()
				} else {
					allow += ", " + root.Name()
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", " + http.MethodOptions
	}
	return allow
}
