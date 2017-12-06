package gorouter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInterface(t *testing.T) {
	var _ http.Handler = New()
}

func TestHandle(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.Handle(POST, "/", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	var cn *node
	for _, child := range s.roots {
		if child.id == POST {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(POST, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestHandleFunc(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.HandleFunc(POST, "/", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	var cn *node
	for _, child := range s.roots {
		if child.id == POST {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(POST, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestPOST(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.POST("/", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	var cn *node
	for _, child := range s.roots {
		if child.id == POST {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(POST, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestGET(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.GET("/", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	var cn *node
	for _, child := range s.roots {
		if child.id == GET {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestPUT(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.PUT("/", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	var cn *node
	for _, child := range s.roots {
		if child.id == PUT {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(PUT, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestDELETE(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.DELETE("/", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	var cn *node
	for _, child := range s.roots {
		if child.id == DELETE {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(DELETE, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestPATCH(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.PATCH("/", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	var cn *node
	for _, child := range s.roots {
		if child.id == PATCH {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(PATCH, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestHEAD(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.HEAD("/", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	var cn *node
	for _, child := range s.roots {
		if child.id == HEAD {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(HEAD, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestCONNECT(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.CONNECT("/", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	var cn *node
	for _, child := range s.roots {
		if child.id == CONNECT {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(CONNECT, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestTRACE(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.TRACE("/", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		serverd = true
	}))

	var cn *node
	for _, child := range s.roots {
		if child.id == TRACE {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	w := httptest.NewRecorder()
	req, err := http.NewRequest(TRACE, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestOPTIONS(t *testing.T) {
	s := New().(*router)

	s.OPTIONS("/", http.HandlerFunc(mockHandler))

	var cn *node
	for _, child := range s.roots {
		if child.id == OPTIONS {
			cn = child
			break
		}
	}

	if cn == nil {
		t.Error("Route not found")
	}

	s.GET("/x", http.HandlerFunc(mockHandler))
	s.POST("/x", http.HandlerFunc(mockHandler))

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
	s := New().(*router)

	s.GET("/x", http.HandlerFunc(mockHandler))

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
	s := New().(*router)

	s.GET("/x", http.HandlerFunc(mockHandler))

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
	s := New().(*router)

	serverd := false
	s.GET("/{param}", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("param") != "x" {
			t.Errorf("Wrong params value. Expected 'x', actual '%s'", params.Value("param"))
		}
	}))

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

func TestRegexpParam(t *testing.T) {
	s := New().(*router)

	serverd := false
	s.GET("/{param:r([a-z]+)go}", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("param") != "rxgo" {
			t.Errorf("Wrong params value. Expected 'rxgo', actual '%s'", params.Value("param"))
		}
	}))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, "/rxgo", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestServeFiles(t *testing.T) {
	s := New().(*router)

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

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Router should panic for empty path")
		}
	}()

	s.ServeFiles("", true)
}

func TestNilMiddleware(t *testing.T) {
	s := New().(*router)

	s.GET("/{param}", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("test"))
	}))

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

	s := New(panicMiddleware).(*router)

	s.GET("/{param}", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("test panic recover")
	}))

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
	s := New().(*router)

	s.GET("/{param}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		w.Write([]byte(params.Value("param")))
	}))

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
	s := New().(*router)

	serverd := false
	s.GET("/users/{user}/starred", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("user") != "x" {
			t.Errorf("Wrong params value. Expected 'x', actual '%s'", params.Value("user"))
		}
	}))

	s.GET("/applications/{client_id}/tokens", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("client_id") != "client_id" {
			t.Errorf("Wrong params value. Expected 'client_id', actual '%s'", params.Value("client_id"))
		}
	}))

	s.GET("/applications/{client_id}/tokens/{access_token}", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("client_id") != "client_id" {
			t.Errorf("Wrong params value. Expected 'client_id', actual '%s'", params.Value("client_id"))
		}

		if params.Value("access_token") != "access_token" {
			t.Errorf("Wrong params value. Expected 'access_token', actual '%s'", params.Value("access_token"))
		}
	}))

	s.GET("/users/{user}/received_events", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("user") != "user1" {
			t.Errorf("Wrong params value. Expected 'user1', actual '%s'", params.Value("user"))
		}
	}))

	s.GET("/users/{user}/received_events/public", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("user") != "user2" {
			t.Errorf("Wrong params value. Expected 'user2', actual '%s'", params.Value("user"))
		}
	}))

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

func TestMountSubRouter(t *testing.T) {
	rGlobal1 := mockMiddlewareWithBody("rg1")
	rGlobal2 := mockMiddlewareWithBody("rg2")
	r1 := mockMiddlewareWithBody("r1")
	r2 := mockMiddlewareWithBody("r2")

	r := New(rGlobal1, rGlobal2).(*router)

	r.GET("/{param}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("r"))
	}))

	r.USE(GET, "/{param}", r1)
	r.USE(GET, "/{param}", r2)

	sGlobal1 := mockMiddlewareWithBody("sg1")
	sGlobal2 := mockMiddlewareWithBody("sg2")
	s1 := mockMiddlewareWithBody("s1")
	s2 := mockMiddlewareWithBody("s2")

	s := New(sGlobal1, sGlobal2).(*router)

	s.GET("/y", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("s"))
	}))

	s.USE(GET, "/y", s1)
	s.USE(GET, "/y", s2)

	r.Mount("/{param}", s)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, "/x/y", nil)
	if err != nil {
		t.Fatal(err)
	}

	s.ServeHTTP(w, req)

	if w.Body.String() != "rg1rg2r1r2sg1sg2s1s2" {
		t.Errorf("Router mount sub router middleware error %s", w.Body.String())
	}
}
