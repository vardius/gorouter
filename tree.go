package gorouter

import (
	"net/http"

	"github.com/vardius/gorouter/v4/mux"
)

func allowed(t mux.Tree, method, path string) (allow string) {
	// tree roots should be http method nodes only
	for _, root := range t {
		if root.Name() == http.MethodOptions {
			continue
		}

		switch {
		case path == "*":
			allow = appendMethod(allow, root.Name())
		case root.Name() == method:
		case path == "":
			// The root path "/" is served by the method node itself.
			if root.Route() != nil {
				allow = appendMethod(allow, root.Name())
			}
		default:
			if route, _ := root.Tree().MatchRoute(path); route != nil {
				allow = appendMethod(allow, root.Name())
			}
		}
	}

	if len(allow) > 0 {
		allow = appendMethod(allow, http.MethodOptions)
	}

	return allow
}

func appendMethod(allow, method string) string {
	if allow == "" {
		return method
	}

	return allow + ", " + method
}
