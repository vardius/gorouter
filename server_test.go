package goserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInterface(t *testing.T) {
	var _ http.Handler = New()
}

func TestHandle(t *testing.T) {
	s := New().(*server)

	s.Handle(POST, "/", http.HandlerFunc(mockHandler))

	var cn *node
	for _, child := range s.roots {
		if child.pattern == POST {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}
}

func TestHandleFunc(t *testing.T) {
	s := New().(*server)

	s.HandleFunc(POST, "/", mockHandler)

	var cn *node
	for _, child := range s.roots {
		if child.pattern == POST {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}
}

func TestPOST(t *testing.T) {
	s := New().(*server)

	s.POST("/", mockHandler)

	var cn *node
	for _, child := range s.roots {
		if child.pattern == POST {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}
}

func TestGET(t *testing.T) {
	s := New().(*server)

	s.GET("/", mockHandler)

	var cn *node
	for _, child := range s.roots {
		if child.pattern == GET {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}
}

func TestPUT(t *testing.T) {
	s := New().(*server)

	s.PUT("/", mockHandler)

	var cn *node
	for _, child := range s.roots {
		if child.pattern == PUT {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}
}

func TestDELETE(t *testing.T) {
	s := New().(*server)

	s.DELETE("/", mockHandler)

	var cn *node
	for _, child := range s.roots {
		if child.pattern == DELETE {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}
}

func TestPATCH(t *testing.T) {
	s := New().(*server)

	s.PATCH("/", mockHandler)

	var cn *node
	for _, child := range s.roots {
		if child.pattern == PATCH {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}
}

func TestHEAD(t *testing.T) {
	s := New().(*server)

	s.HEAD("/", mockHandler)

	var cn *node
	for _, child := range s.roots {
		if child.pattern == HEAD {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}
}

func TestOPTIONS(t *testing.T) {
	s := New().(*server)

	s.OPTIONS("/", mockHandler)

	var cn *node
	for _, child := range s.roots {
		if child.pattern == OPTIONS {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	s.GET("/x", mockHandler)
	s.POST("/x", mockHandler)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(OPTIONS, "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Header().Get("Allow") == "" {
		t.Error("Allow header should not be empty")
	}
}

func TestNotFound(t *testing.T) {
	s := New().(*server)

	s.GET("/x", mockHandler)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(POST, "/y", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("NotFound error, actual code: %d", w.Code)
	}

	s.NotFound(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("test"))
	}))

	if s.notFound == nil {
		t.Error("NotFound handler error")
	}

	w = httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Body.String() != "test" {
		t.Error("Not found handler wasn't invoked")
	}
}

func TestNotAllowed(t *testing.T) {
	s := New().(*server)

	s.GET("/x", mockHandler)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(POST, "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Error("NotAllowed doesnt work")
	}

	s.NotAllowed(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("test"))
	}))

	if s.notAllowed == nil {
		t.Error("NotAllowed handler error")
	}

	w = httptest.NewRecorder()

	s.ServeHTTP(w, req)

	if w.Body.String() != "test" {
		t.Error("Not found handler wasn't invoked")
	}

	w = httptest.NewRecorder()
	req, err = http.NewRequest(POST, "*", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Body.String() != "test" {
		t.Error("Not found handler wasn't invoked")
	}
}

func TestParam(t *testing.T) {
	s := New().(*server)

	serverd := false
	s.GET("/{param}", func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := ParamsFromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("param") != "x" {
			t.Errorf("Wrong params value. Expected 'x', actual '%s'", params.Value("param"))
		}
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestServeFiles(t *testing.T) {
	s := New().(*server)

	s.ServeFiles("static", true)

	if s.fileServer == nil {
		t.Error("File serve handler error")
	}

	w := httptest.NewRecorder()
	r, err := http.NewRequest(GET, "/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("File should not exist")
	}
}

func TestNilMiddleware(t *testing.T) {
	s := New().(*server)

	s.GET("/{param}", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("test"))
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Body.String() != "test" {
		t.Error("Nil middleware works")
	}
}

func TestPanicMiddleware(t *testing.T) {
	paniced := false
	panicMiddleware := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rcv := recover(); rcv != nil {
					paniced = true
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	s := New(panicMiddleware).(*server)

	s.GET("/{param}", func(_ http.ResponseWriter, _ *http.Request) {
		panic("test panic recover")
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if paniced != true {
		t.Error("Panic has not been handled")
	}
}

func TestNodeApplyMiddleware(t *testing.T) {
	s := New().(*server)

	s.GET("/{param}", func(w http.ResponseWriter, r *http.Request) {
		params, ok := ParamsFromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		w.Write([]byte(params.Value("param")))
	})

	s.USE(GET, "/{param}", mockMiddleware)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, "/x", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Body.String() != "middlewarex" {
		t.Errorf("Use global middleware error %s", w.Body.String())
	}
}

func TestChainCalls(t *testing.T) {
	s := New().(*server)

	serverd := false
	s.GET("/users/{user}/starred", func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := ParamsFromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("user") != "x" {
			t.Errorf("Wrong params value. Expected 'x', actual '%s'", params.Value("user"))
		}
	})

	s.GET("/applications/{client_id}/tokens", func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := ParamsFromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("client_id") != "client_id" {
			t.Errorf("Wrong params value. Expected 'client_id', actual '%s'", params.Value("client_id"))
		}
	})

	s.GET("/applications/{client_id}/tokens/{access_token}", func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := ParamsFromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("client_id") != "client_id" {
			t.Errorf("Wrong params value. Expected 'client_id', actual '%s'", params.Value("client_id"))
		}

		if params.Value("access_token") != "access_token" {
			t.Errorf("Wrong params value. Expected 'access_token', actual '%s'", params.Value("access_token"))
		}
	})

	s.GET("/users/{user}/received_events", func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := ParamsFromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("user") != "user1" {
			t.Errorf("Wrong params value. Expected 'user1', actual '%s'", params.Value("user"))
		}
	})

	s.GET("/users/{user}/received_events/public", func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := ParamsFromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("user") != "user2" {
			t.Errorf("Wrong params value. Expected 'user2', actual '%s'", params.Value("user"))
		}
	})

	w := httptest.NewRecorder()

	// //FIRST CALL
	req, err := http.NewRequest(GET, "/users/x/starred", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if !serverd {
		t.Fatal("First not served")
	}

	//SECOND CALL
	req, err = http.NewRequest(GET, "/applications/client_id/tokens", nil)
	if err != nil {
		t.Fatal(err)
	}

	serverd = false
	s.ServeHTTP(w, req)

	if !serverd {
		t.Fatal("Second not served")
	}

	//THIRD CALL
	req, err = http.NewRequest(GET, "/applications/client_id/tokens/access_token", nil)
	if err != nil {
		t.Fatal(err)
	}

	serverd = false
	s.ServeHTTP(w, req)

	if !serverd {
		t.Fatal("Third not served")
	}

	//FOURTH CALL
	req, err = http.NewRequest(GET, "/users/user1/received_events", nil)
	if err != nil {
		t.Fatal(err)
	}

	serverd = false
	s.ServeHTTP(w, req)

	if !serverd {
		t.Fatal("Fourth not served")
	}

	//FIFTH CALL
	req, err = http.NewRequest(GET, "/users/user2/received_events/public", nil)
	if err != nil {
		t.Fatal(err)
	}

	serverd = false
	s.ServeHTTP(w, req)

	if !serverd {
		t.Fatal("Fifth not served")
	}
}
