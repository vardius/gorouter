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
	HEAD    = "HEAD"
	CONNECT = "CONNECT"
	TRACE   = "TRACE"
)

//Server is a micro framwework, HTTP request router, multiplexer, mux
type Server interface {
	//Handle adds http.Handler as router handler
	//under given method and patter
	Handle(method, pattern string, handler http.Handler)

	//HandleFunc adds http.HandlerFunc as router handler
	//under given method and patter
	HandleFunc(method, pattern string, handler http.HandlerFunc)

	//POST adds http.Handler as router handler
	//under POST method and given patter
	POST(pattern string, handler http.Handler)

	//GET adds http.Handler as router handler
	//under GET method and given patter
	GET(pattern string, handler http.Handler)

	//PUT adds http.Handler as router handler
	//under PUT method and given patter
	PUT(pattern string, handler http.Handler)

	//DELETE adds http.Handler as router handler
	//under DELETE method and given patter
	DELETE(pattern string, handler http.Handler)

	//PATCH adds http.Handler as router handler
	//under PATCH method and given patter
	PATCH(pattern string, handler http.Handler)

	//OPTIONS adds http.Handler as router handler
	//under OPTIONS method and given patter
	OPTIONS(pattern string, handler http.Handler)

	//HEAD adds http.Handler as router handler
	//under HEAD method and given patter
	HEAD(pattern string, handler http.Handler)

	//CONNECT adds http.Handler as router handler
	//under CONNECT method and given patter
	CONNECT(pattern string, handler http.Handler)

	//TRACE adds http.Handler as router handler
	//under TRACE method and given patter
	TRACE(pattern string, handler http.Handler)

	//USE adds middleware functions ([]MiddlewareFunc)
	//to whole router branch under given method and patter
	USE(method, pattern string, fs ...MiddlewareFunc)

	//ServeHTTP dispatches the request to the route handler
	//whose pattern matches the request URL
	ServeHTTP(http.ResponseWriter, *http.Request)

	//ServeFile replies to the request with the
	//contents of the named file or directory.
	ServeFiles(path string, strip bool)

	//NotFound replies to the request with the
	//404 Error code
	NotFound(http.Handler)

	//NotFound replies to the request with the
	//405 Error code
	NotAllowed(http.Handler)
}

type server struct {
	roots      []*node
	middleware middleware
	fileServer http.Handler
	notFound   http.Handler
	notAllowed http.Handler
}

func (s *server) Handle(m, p string, h http.Handler) {
	s.addRoute(m, p, h)
}

func (s *server) HandleFunc(m, p string, f http.HandlerFunc) {
	s.addRoute(m, p, http.HandlerFunc(f))
}

func (s *server) POST(p string, f http.Handler) {
	s.addRoute(POST, p, f)
}

func (s *server) GET(p string, f http.Handler) {
	s.addRoute(GET, p, f)
}

func (s *server) PUT(p string, f http.Handler) {
	s.addRoute(PUT, p, f)
}

func (s *server) DELETE(p string, f http.Handler) {
	s.addRoute(DELETE, p, f)
}

func (s *server) PATCH(p string, f http.Handler) {
	s.addRoute(PATCH, p, f)
}

func (s *server) OPTIONS(p string, f http.Handler) {
	s.addRoute(OPTIONS, p, f)
}

func (s *server) HEAD(p string, f http.Handler) {
	s.addRoute(HEAD, p, f)
}

func (s *server) CONNECT(p string, f http.Handler) {
	s.addRoute(CONNECT, p, f)
}

func (s *server) TRACE(p string, f http.Handler) {
	s.addRoute(TRACE, p, f)
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
	route, params := s.getRoute(req.Method, req.URL.Path)
	if route != nil {
		if h := route.chain(); h != nil {
			req = req.WithContext(newContext(req, params))
			h.ServeHTTP(w, req)
			return
		}
	}

	//Handle OPTIONS
	if req.Method == OPTIONS {
		if allow := s.allowed(req.Method, req.URL.Path); len(allow) > 0 {
			w.Header().Set("Allow", allow)
			return
		}
	} else if req.Method == GET && s.fileServer != nil {
		//Handle file serve
		s.serveFiles(w, req)
		return
	} else {
		//Handle 405
		if allow := s.allowed(req.Method, req.URL.Path); len(allow) > 0 {
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

func (s *server) addRoute(method, path string, f http.Handler) {
	var root *node
	for _, root = range s.roots {
		if method == root.id {
			break
		}
		root = nil
	}

	if root == nil {
		root = newRoot(method)
		s.roots = append(s.roots, root)
	}

	paths := splitPath(path)

	r := newRoute(f)
	r.addMiddleware(s.middleware)

	n := root.addChild(paths)
	n.setRoute(r)
}

func (s *server) addMiddleware(method, path string, fs ...MiddlewareFunc) {
	type recFunc func(recFunc, *node, middleware)
	c := func(c recFunc, n *node, m middleware) {
		if n.route != nil {
			n.route.addMiddleware(m)
		}
		for _, child := range n.children.statics {
			c(c, child, m)
		}
		for _, child := range n.children.regexps {
			c(c, child, m)
		}
		if n.children.wildcard != nil {
			c(c, n.children.wildcard, m)
		}
	}

	paths := splitPath(path)

	for _, root := range s.roots {
		if method == "" || method == root.id {
			node, _ := root.child(paths)
			c(c, node, fs)
		}
	}
}

func (s *server) getRoute(method, path string) (*route, Params) {
	for _, root := range s.roots {
		if root.id == method {
			node, params := root.childByPath(path)
			if node != nil {
				return node.route, params
			}
			break
		}
	}
	return nil, nil
}

func splitPath(path string) (parts []string) {
	for {
		if i := strings.IndexByte(path, '{'); i >= 0 {
			if part := trimPath(path[:i]); part != "" {
				parts = append(parts, part)
			}
			if j := strings.IndexByte(path, '}') + 1; j > 0 {
				if part := trimPath(path[i:j]); part != "" {
					parts = append(parts, part)
				}
				i = j
			} else {
				continue
			}
			path = path[i:]
		} else {
			break
		}
	}

	if len(path) != 0 && path != "/" {
		if part := trimPath(path); part != "" {
			parts = append(parts, part)
		}
	}

	return
}

func trimPath(path string) string {
	pathLen := len(path)
	if pathLen > 0 && path[0] == '/' {
		path = path[1:]
		pathLen--
	}

	if pathLen > 0 && path[pathLen-1] == '/' {
		path = path[:pathLen-1]
		pathLen--
	}
	return path
}

func (s *server) allowed(method, path string) (allow string) {
	if path == "*" {
		for _, root := range s.roots {
			if root.id == OPTIONS {
				continue
			}
			if len(allow) == 0 {
				allow = root.id
			} else {
				allow += ", " + root.id
			}
		}
	} else {
		for _, root := range s.roots {
			if root.id == method || root.id == OPTIONS {
				continue
			}

			n, _ := root.childByPath(path)
			if n != nil && n.route != nil {
				if len(allow) == 0 {
					allow = root.id
				} else {
					allow += ", " + root.id
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
		roots:      make([]*node, 0),
		middleware: newMiddleware(fs...),
	}
}
