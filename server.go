package goserver

import (
	"net/http"
	"os"
	"strings"
)

type (
	Server interface {
		POST(path string, f http.HandlerFunc)
		GET(path string, f http.HandlerFunc)
		PUT(path string, f http.HandlerFunc)
		DELETE(path string, f http.HandlerFunc)
		PATCH(path string, f http.HandlerFunc)
		OPTIONS(path string, f http.HandlerFunc)
		USE(method, path string, fs ...middlewareFunc)
		ServeHTTP(http.ResponseWriter, *http.Request)
		ServeFiles(path string, strip bool)
		NotFound(http.Handler)
		NotAllowed(http.Handler)
	}
	server struct {
		root       *node
		middleware middlewares
		fileServer http.Handler
		notFound   http.Handler
		notAllowed http.Handler
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

func (s *server) USE(method, path string, fs ...middlewareFunc) {
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

func (s *server) addMiddleware(method, path string, fs ...middlewareFunc) {
	var paths []string
	if path := strings.Trim(path, "/"); path != "" {
		paths = strings.Split(path, "/")
	}

	if method == "" {
		for _, node := range s.root.children {
			addMiddleware(node, paths, fs)
		}
	} else {
		paths = append([]string{method}, paths...)
		node, _ := s.root.child(paths)
		addMiddleware(node, paths, fs)
	}
}

func (s *server) getRouteFromRequest(req *http.Request) (*route, Params) {
	var paths []string
	if path := strings.Trim(req.URL.Path, "/"); path != "" {
		paths = strings.Split(path, "/")
	}

	paths = append([]string{req.Method}, paths...)
	node, params := s.root.child(paths)

	return node.route, params
}

func (s *server) allowed(req *http.Request) (allow string) {
	path := req.URL.Path
	if path == "*" {
		for method := range s.root.children {
			if method == OPTIONS {
				continue
			}
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else {
		for method, root := range s.root.children {
			if method == req.Method || method == OPTIONS {
				continue
			}

			var paths []string
			if path = strings.Trim(path, "/"); path != "" {
				paths = strings.Split(path, "/")
			}

			r, _ := root.child(paths)
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

func addMiddleware(n *node, paths []string, m middlewares) {
	if n.route != nil {
		n.route.addMiddleware(m)
		for _, node := range n.children {
			addMiddleware(node, paths[1:], m)
		}
	}
}

func New(fs ...middlewareFunc) Server {
	return &server{
		root:       newRootNode(""),
		middleware: append(middlewares(nil), fs...),
	}
}
