package goserver

import (
	"net/http"
	"os"
)

// HTTP methods constants
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	HEAD    = "HEAD"
)

type (
	//Server interface
	Server interface {
		Handle(method, pattern string, handler http.Handler)
		HandleFunc(method, pattern string, handler http.HandlerFunc)
		POST(pattern string, handler http.HandlerFunc)
		GET(pattern string, handler http.HandlerFunc)
		PUT(pattern string, handler http.HandlerFunc)
		DELETE(pattern string, handler http.HandlerFunc)
		PATCH(pattern string, handler http.HandlerFunc)
		OPTIONS(pattern string, handler http.HandlerFunc)
		HEAD(pattern string, handler http.HandlerFunc)
		USE(method, pattern string, fs ...MiddlewareFunc)
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

func (s *server) Handle(method, p string, h http.Handler) {
	s.HandleFunc(method, p, func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(w, req)
	})
}

func (s *server) HandleFunc(method, p string, f http.HandlerFunc) {
	s.addRoute(method, p, f)
}

func (s *server) POST(p string, f http.HandlerFunc) {
	s.addRoute(POST, p, f)
}

func (s *server) GET(p string, f http.HandlerFunc) {
	s.addRoute(GET, p, f)
}

func (s *server) PUT(p string, f http.HandlerFunc) {
	s.addRoute(PUT, p, f)
}

func (s *server) DELETE(p string, f http.HandlerFunc) {
	s.addRoute(DELETE, p, f)
}

func (s *server) PATCH(p string, f http.HandlerFunc) {
	s.addRoute(PATCH, p, f)
}

func (s *server) OPTIONS(p string, f http.HandlerFunc) {
	s.addRoute(OPTIONS, p, f)
}

func (s *server) HEAD(p string, f http.HandlerFunc) {
	s.addRoute(HEAD, p, f)
}

func (s *server) USE(method, p string, fs ...MiddlewareFunc) {
	s.addMiddleware(method, p, fs...)
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
	// route, params := s.getRouteFromRequest(req)
	route, _ := s.getRouteFromRequest(req)
	if route != nil {
		if h := route.handler(); h != nil {
			// req = req.WithContext(newContextFromRequest(req, params))
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
	paths := getPaths(path)
	paths = append([]string{method}, paths...)

	r := newRoute(f)
	r.addMiddleware(s.middleware)

	n := s.root.addChild(paths)
	n.setRoute(r)
}

func (s *server) addMiddleware(method, path string, fs ...MiddlewareFunc) {
	type recFunc func(recFunc, *node, string, middleware)
	c := func(c recFunc, n *node, path string, m middleware) {
		if n.route != nil {
			n.route.addMiddleware(m)
			for _, node := range n.children {
				c(c, node, path[len(node.pattern):], m)
			}
		}
	}

	for _, node := range s.root.children {
		if method == "" || method == node.pattern {
			c(c, node, path, fs)
		}
	}
}

func (s *server) getRouteFromRequest(req *http.Request) (*route, Params) {
	for _, node := range s.root.children {
		if node.pattern == req.Method {
			node, params := node.child(req.URL.Path)
			if node != nil {
				return node.route, params
			}
			break
		}
	}
	return nil, nil
}

func (s *server) allowed(req *http.Request) string {
	var allow string
	path := req.URL.Path
	if path == "*" {
		for _, n := range s.root.children {
			if n.pattern == OPTIONS {
				continue
			}
			if len(allow) == 0 {
				allow = n.pattern
			} else {
				allow += ", " + n.pattern
			}
		}
	} else {
		for _, root := range s.root.children {
			if root.pattern == req.Method || root.pattern == OPTIONS {
				continue
			}
			n, _ := root.child(path)
			if n != nil && n.route != nil {
				if len(allow) == 0 {
					allow = root.pattern
				} else {
					allow += ", " + root.pattern
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
