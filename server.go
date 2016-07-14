package goapi

import (
	"net/http"
	"strings"
	"sync"
)

type (
	Server struct {
		routes     tree
		middleware middlewares
		routesMu   sync.RWMutex
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

func (s *Server) POST(path string, f http.HandlerFunc) {
	s.addRoute(POST, path, f)
}

func (s *Server) GET(path string, f http.HandlerFunc) {
	s.addRoute(GET, path, f)
}

func (s *Server) PUT(path string, f http.HandlerFunc) {
	s.addRoute(PUT, path, f)
}

func (s *Server) DELETE(path string, f http.HandlerFunc) {
	s.addRoute(DELETE, path, f)
}

func (s *Server) PATCH(path string, f http.HandlerFunc) {
	s.addRoute(PATCH, path, f)
}

func (s *Server) OPTIONS(path string, f http.HandlerFunc) {
	s.addRoute(OPTIONS, path, f)
}

func (s *Server) Use(path string, priority int, f MiddlewareFunc) {
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
			if route := r.getRoute(paths); route != nil {
				route.middleware = append(route.middleware, m)
				sortByPriority(route.middleware)
			}
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//deffer recover panic here

	s.routesMu.RLock()
	defer s.routesMu.RUnlock()
	if r := s.routes[req.Method]; r != nil {
		paths := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
		if route := r.getRoute(paths); route != nil {
			if route.handler != nil {
				for _, m := range s.middleware {
					if err := m.handler(req); err != nil {
						http.Error(w, err.Error(), err.Status())
						return
					}
				}
				for _, m := range route.middleware {
					if err := m.handler(req); err != nil {
						http.Error(w, err.Error(), err.Status())
						return
					}
				}

				route.handler(w, req)
				return
			}
		}
	}

	//handle options method here

	http.NotFound(w, req)
}

func (s *Server) addRoute(method, path string, f http.HandlerFunc) {
	paths := strings.Split(strings.Trim(path, "/"), "/")

	s.routesMu.Lock()
	defer s.routesMu.Unlock()

	var r *route
	if r = s.routes[method]; r == nil {
		r = newRoute("/")
		s.routes[method] = r
	}
	r.addRoute(paths, f)
}

func New() *Server {
	return &Server{
		routes: make(tree),
	}
}
