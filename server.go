package goserver

import (
	"net/http"
	"os"
	"strings"
)

// HTTP methods constants
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

type (
	//Server interface
	Server interface {
		POST(path string, f http.HandlerFunc)
		GET(path string, f http.HandlerFunc)
		PUT(path string, f http.HandlerFunc)
		DELETE(path string, f http.HandlerFunc)
		PATCH(path string, f http.HandlerFunc)
		OPTIONS(path string, f http.HandlerFunc)
		USE(method, path string, fs ...MiddlewareFunc)
		ServeHTTP(http.ResponseWriter, *http.Request)
		ServeFiles(path string, strip bool)
		NotFound(http.Handler)
		NotAllowed(http.Handler)
	}
	server struct {
		root       *node
		middleware middleware
		fileServer http.Handler
		notFound   http.Handler
		notAllowed http.Handler
	}
)

func (s *server) POST(path string, f http.HandlerFunc) {
	s.addRoute(POST, path, f)
}

func (s *server) GET(path string, f http.HandlerFunc) {
	s.addRoute(GET, path, f)
}

func (s *server) PUT(path string, f http.HandlerFunc) {
	s.addRoute(PUT, path, f)
}

func (s *server) DELETE(path string, f http.HandlerFunc) {
	s.addRoute(DELETE, path, f)
}

func (s *server) PATCH(path string, f http.HandlerFunc) {
	s.addRoute(PATCH, path, f)
}

func (s *server) OPTIONS(path string, f http.HandlerFunc) {
	s.addRoute(OPTIONS, path, f)
}

func (s *server) USE(method, path string, fs ...MiddlewareFunc) {
	s.addMiddleware(method, path, fs...)
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
	route, params := s.getRouteFromRequest(req)
	if route != nil {
		if h := route.handler(); h != nil {
			req = req.WithContext(newContextFromRequest(req, params))
			h.ServeHTTP(w, req)
			return
		}
	}

	//Handle OPTIONS
	if req.Method == OPTIONS {
		if allow := s.allowed(req); len(allow) > 0 {
			w.Header().Set("Allow", allow)
			return
		}
	} else if req.Method == GET && s.fileServer != nil {
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

func (s *server) addRoute(method, path string, f http.HandlerFunc) {
	paths := strings.Split(strings.Trim(path, "/"), "/")
	paths = append([]string{method}, paths...)

	r := newRoute(f)
	r.addMiddleware(s.middleware)

	n := s.root.addChild(paths)
	n.setRoute(r)
}

func (s *server) addMiddleware(method, path string, fs ...MiddlewareFunc) {
	type recFunc func(recFunc, *node, []string, middleware)
	c := func(c recFunc, n *node, paths []string, m middleware) {
		if n.route != nil {
			n.route.addMiddleware(m)
			for _, node := range n.children {
				c(c, node, paths[1:], m)
			}
		}
	}

	var paths []string
	if path := strings.Trim(path, "/"); path != "" {
		paths = strings.Split(path, "/")
	}

	if method == "" {
		for _, node := range s.root.children {
			c(c, node, paths, fs)
		}
	} else {
		paths = append([]string{method}, paths...)
		node, _ := s.root.child(paths)
		c(c, node, paths, fs)
	}
}

func (s *server) getRouteFromRequest(req *http.Request) (*route, Params) {
	var paths []string
	if path := strings.Trim(req.URL.Path, "/"); path != "" {
		paths = strings.Split(path, "/")
	}

	paths = append([]string{req.Method}, paths...)
	node, params := s.root.child(paths)
	if node != nil {
		return node.route, params
	}

	return nil, params
}

func (s *server) allowed(req *http.Request) string {
	var allow string
	path := req.URL.Path
	if path == "*" {
		for _, n := range s.root.children {
			if n.path == OPTIONS {
				continue
			}
			if len(allow) == 0 {
				allow = n.path
			} else {
				allow += ", " + n.path
			}
		}
	} else {
		for _, root := range s.root.children {
			if root.path == req.Method || root.path == OPTIONS {
				continue
			}

			var paths []string
			if path = strings.Trim(path, "/"); path != "" {
				paths = strings.Split(path, "/")
			}

			n, _ := root.child(paths)
			if n != nil && n.route != nil {
				if len(allow) == 0 {
					allow = root.path
				} else {
					allow += ", " + root.path
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return allow
}

//Creates new Server instance, return pointer
func New(fs ...MiddlewareFunc) Server {
	return &server{
		root:       newRoot(""),
		middleware: newMiddleware(fs...),
	}
}
