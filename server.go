package goserver

import (
	"net/http"
	"os"
	"strings"
	"sync"
)

type (
	HandlerFunc      func(http.ResponseWriter, *http.Request, *Context)
	PanicHandlerFunc func(http.ResponseWriter, *http.Request, interface{})
	Server           interface {
		POST(path string, f HandlerFunc)
		GET(path string, f HandlerFunc)
		PUT(path string, f HandlerFunc)
		DELETE(path string, f HandlerFunc)
		PATCH(path string, f HandlerFunc)
		OPTIONS(path string, f HandlerFunc)
		Use(path string, priority int, f MiddlewareFunc)
		ServeHTTP(http.ResponseWriter, *http.Request)
		ServeFiles(path string, strip bool)
		NotFound(http.Handler)
		NotAllowed(http.Handler)
		OnPanic(PanicHandlerFunc)
		Routes() map[string]Route
	}
	server struct {
		routes     tree
		middleware middlewares
		routesMu   sync.RWMutex
		fileServer http.Handler
		notFound   http.Handler
		notAllowed http.Handler
		onPanic    PanicHandlerFunc
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
	s.addMiddleware(path, priority, f)
}

func (s *server) NotFound(notFound http.Handler) {
	s.notFound = notFound
}

func (s *server) NotAllowed(notAllowed http.Handler) {
	s.notAllowed = notAllowed
}

func (s *server) OnPanic(onPanic PanicHandlerFunc) {
	s.onPanic = onPanic
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
	if s.onPanic != nil {
		defer s.recv(w, req)
	}

	s.routesMu.RLock()
	defer s.routesMu.RUnlock()

	if r := s.routes[req.Method]; r != nil {
		ctx, err := fromRequest(r, req)
		if err == nil {
			route := ctx.Route.(*route)
			if route.handler != nil {
				for _, m := range s.middleware {
					if err := m.handler(w, req, ctx); err != nil {
						http.Error(w, err.Error(), err.Status())
						return
					}
				}
				for _, m := range route.middleware {
					if err := m.handler(w, req, ctx); err != nil {
						http.Error(w, err.Error(), err.Status())
						return
					}
				}

				route.handler(w, req, ctx)
				return
			}
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

func (s *server) Routes() map[string]Route {
	newMap := make(map[string]Route)
	for path, route := range s.routes {
		newMap[path] = route
	}
	return newMap
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

func (s *server) addMiddleware(path string, priority int, f MiddlewareFunc) {
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

func (s *server) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		s.onPanic(w, req, rcv)
	}
}

func (s *server) allowed(req *http.Request) (allow string) {
	path := req.URL.Path
	if path == "*" {
		for method := range s.routes {
			if method == "OPTIONS" {
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
			if method == req.Method || method == "OPTIONS" {
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

func New() Server {
	return &server{
		routes: make(tree),
	}
}
