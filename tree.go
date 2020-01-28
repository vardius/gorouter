package gorouter

import (
	pathutils "github.com/vardius/gorouter/v4/path"
	"net/http"

	"github.com/vardius/gorouter/v4/mux"
)

func allowed(t mux.Tree, method, path string) (allow string) {
	path = pathutils.TrimSlash(path)

	if path == "*" {
		// tree tree roots should be http method nodes only
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
		// tree tree roots should be http method nodes only
		for _, root := range t {
			if root.Name() == method || root.Name() == http.MethodOptions {
				continue
			}

			if route, _, _ := root.Tree().MatchRoute(path); route != nil {
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
