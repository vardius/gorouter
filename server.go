package goapi

import (
	"net/http"
	"strings"
	"sync"
)

type (
	server struct {
		routes     tree
		middleware middlewares
		routesMu   sync.RWMutex
	}
	HandlerFunc func(http.ResponseWriter, *http.Request, *Context)
	Server      interface {
		POST(path string, f HandlerFunc)
		GET(path string, f HandlerFunc)
		PUT(path string, f HandlerFunc)
		DELETE(path string, f HandlerFunc)
		PATCH(path string, f HandlerFunc)
		OPTIONS(path string, f HandlerFunc)
		Use(path string, priority int, f MiddlewareFunc)
		ServeHTTP(http.ResponseWriter, *http.Request)
		Routes() map[string]Route
	}
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

func (s *server) POST(path string, f HandlerFunc) {
	s.addRoute(POST, path, f)
}

func (s *server) GET(path string, f HandlerFunc) {
	s.addRoute(GET, path, f)
}

func (s *server) PUT(path string, f HandlerFunc) {
	s.addRoute(PUT, path, f)
}

func (s *server) DELETE(path string, f HandlerFunc) {
	s.addRoute(DELETE, path, f)
}

func (s *server) PATCH(path string, f HandlerFunc) {
	s.addRoute(PATCH, path, f)
}

func (s *server) OPTIONS(path string, f HandlerFunc) {
	s.addRoute(OPTIONS, path, f)
}

func (s *server) Use(path string, priority int, f MiddlewareFunc) {
	m := &middleware{
		path:     path,
		priority: priority,
		handler:  f,
	}
	if path == "" {
		s.middleware = append(s.middleware, m)
		sortByPriority(s.middleware)
	} else if path == "/" {
		for _, r := range s.routes {
			r.middleware = append(r.middleware, m)
			sortByPriority(r.middleware)
		}
	} else {
		paths := strings.Split(strings.Trim(path, "/"), "/")
		for _, r := range s.routes {
			route, _ := r.getRoute(paths)
			if route != nil {
				route.middleware = append(route.middleware, m)
				sortByPriority(route.middleware)
			}
		}
	}
}

func (s *server) Routes() map[string]Route {
	newMap := make(map[string]Route)
	for path, route := range s.routes {
		newMap[path] = route
	}
	return newMap
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//deffer recover panic here

	s.routesMu.RLock()
	defer s.routesMu.RUnlock()
	if r := s.routes[req.Method]; r != nil {
		ctx, err := fromRequest(r, req)
		if err == nil {
			route := ctx.Route.(*route)
			if route.handler != nil {
				for _, m := range s.middleware {
					if err := m.handler(req, ctx); err != nil {
						http.Error(w, err.Error(), err.Status())
						return
					}
				}
				for _, m := range route.middleware {
					if err := m.handler(req, ctx); err != nil {
						http.Error(w, err.Error(), err.Status())
						return
					}
				}

				route.handler(w, req, ctx)
				return
			}
		}
	}

	//handle options method here

	http.NotFound(w, req)
}

func New() Server {
	return &server{
		routes: make(tree),
	}
}

func (s *server) addRoute(method, path string, f HandlerFunc) {
	paths := strings.Split(strings.Trim(path, "/"), "/")

	s.routesMu.Lock()
	defer s.routesMu.Unlock()

	var r *route
	if r = s.routes[method]; r == nil {
		r = newRoute(nil, "/")
		s.routes[method] = r
	}
	r.addRoute(paths, f)
}
