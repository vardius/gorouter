package goserver

import (
	"net/http"
	"os"
	"strings"
	"sync"
)

type (
	Server interface {
		POST(path string, f http.HandlerFunc)
		GET(path string, f http.HandlerFunc)
		PUT(path string, f http.HandlerFunc)
		DELETE(path string, f http.HandlerFunc)
		PATCH(path string, f http.HandlerFunc)
		OPTIONS(path string, f http.HandlerFunc)
		ServeHTTP(http.ResponseWriter, *http.Request)
		ServeFiles(path string, strip bool)
		NotFound(http.Handler)
		NotAllowed(http.Handler)
		Routes() map[string]Route
	}
	server struct {
		routes     tree
		middleware middlewares
		routesMu   sync.RWMutex
		fileServer http.Handler
		notFound   http.Handler
		notAllowed http.Handler
	}
)

const (
	get     = "GET"
	post    = "POST"
	put     = "PUT"
	delete  = "DELETE"
	patch   = "PATCH"
	options = "OPTIONS"
)

func (s *server) POST(path string, f http.HandlerFunc) {
	s.addRoute(post, path, f)
}

func (s *server) GET(path string, f http.HandlerFunc) {
	s.addRoute(get, path, f)
}

func (s *server) PUT(path string, f http.HandlerFunc) {
	s.addRoute(put, path, f)
}

func (s *server) DELETE(path string, f http.HandlerFunc) {
	s.addRoute(delete, path, f)
}

func (s *server) PATCH(path string, f http.HandlerFunc) {
	s.addRoute(patch, path, f)
}

func (s *server) OPTIONS(path string, f http.HandlerFunc) {
	s.addRoute(options, path, f)
}

func (s *server) NotFound(notFound http.Handler) {
	s.notFound = notFound
}

func (s *server) NotAllowed(notAllowed http.Handler) {
	s.notAllowed = notAllowed
}

func (s *server) ServeFiles(path string, strip bool) {
	if path == "" {
		panic("goapi.ServeFiles: empty path!")
	}
	handler := http.FileServer(http.Dir(path))
	if strip {
		handler = http.StripPrefix("/"+path+"/", handler)
	}
	s.fileServer = handler
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.routesMu.RLock()
	defer s.routesMu.RUnlock()

	if r := s.routes[req.Method]; r != nil {
		route, params := r.getRouteFromRequest(req)
		if route != nil {
			if route.handler != nil {
				h := s.middleware.handleFunc(route.handler)
				req = req.WithContext(newContextFromRequest(req, params))
				h.ServeHTTP(w, req)
				return
			}
		}
	}

	//Handle OPTIONS
	if req.Method == options {
		if allow := s.allowed(req); len(allow) > 0 {
			w.Header().Set("Allow", allow)
			return
		}
	} else if req.Method == get && s.fileServer != nil {
		//Handle file serve
		s.serveFiles(w, req)
		return
	} else {
		//Handle 405
		if allow := s.allowed(req); len(allow) > 0 {
			w.Header().Set("Allow", allow)
			s.serveNotAllowed(w, req)
			return
		}
	}

	//Handle 404
	s.serveNotFound(w, req)
}

func (s *server) Routes() map[string]Route {
	newMap := make(map[string]Route)
	for path, route := range s.routes {
		newMap[path] = route
	}
	return newMap
}

func (s *server) addRoute(method, path string, f http.HandlerFunc) {
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

func (s *server) allowed(req *http.Request) (allow string) {
	path := req.URL.Path
	if path == "*" {
		for method := range s.routes {
			if method == options {
				continue
			}
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else {
		for method, root := range s.routes {
			if method == req.Method || method == options {
				continue
			}

			var paths []string
			if path = strings.Trim(path, "/"); path != "" {
				paths = strings.Split(path, "/")
			}

			r, _ := root.getRoute(paths)
			if r != nil {
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

func (s *server) serveFiles(w http.ResponseWriter, req *http.Request) {
	fp := req.URL.Path
	//Return a 404 if the file doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			s.serveNotFound(w, req)
			return
		}
	}
	//Return a 404 if the request is for a directory
	if info.IsDir() {
		s.serveNotFound(w, req)
		return
	}
	s.fileServer.ServeHTTP(w, req)
}

func (s *server) serveNotFound(w http.ResponseWriter, req *http.Request) {
	if s.notFound != nil {
		s.notFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func (s *server) serveNotAllowed(w http.ResponseWriter, req *http.Request) {
	if s.notAllowed != nil {
		s.notAllowed.ServeHTTP(w, req)
	} else {
		http.Error(w,
			http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
	}
}

func New(fs ...MiddlewareFunc) Server {
	return &server{
		routes:     make(tree),
		middleware: newMiddleware(fs...),
	}
}
