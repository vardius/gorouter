package gorouter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockHandler struct {
	served bool
}

func (mh *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mh.served = true
}

type mockFileSystem struct {
	opened bool
}

func (mfs *mockFileSystem) Open(name string) (http.File, error) {
	mfs.opened = true
	return nil, errors.New("")
}

func mockMiddleware(body string) MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(body))
			h.ServeHTTP(w, r)
		})
	}
}

func mockServeHTTP(h http.Handler, method, path string) error {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return err
	}

	h.ServeHTTP(w, req)

	return nil
}

func checkIfHasRootRoute(t *testing.T, router *router, method string) {
	if rootRoute := router.routes.getByID(method); rootRoute == nil {
		t.Error("Route not found")
	}
}

func testBasicMethod(t *testing.T, router *router, h func(pattern string, handler http.Handler), method string) {
	handler := &mockHandler{}
	h("/x/y", handler)

	checkIfHasRootRoute(t, router, method)

	err := mockServeHTTP(router, method, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if handler.served != true {
		t.Error("Handler has not been serverd")
	}
}

func TestInterface(t *testing.T) {
	var _ http.Handler = New()
}

func TestHandle(t *testing.T) {
	t.Parallel()

	handler := &mockHandler{}
	router := New().(*router)
	router.Handle(POST, "/x/y", handler)

	checkIfHasRootRoute(t, router, POST)

	err := mockServeHTTP(router, POST, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if handler.served != true {
		t.Error("Handler has not been serverd")
	}
}

func TestHandleFunc(t *testing.T) {
	t.Parallel()

	served := false
	router := New().(*router)
	router.HandleFunc(POST, "/x/y", func(w http.ResponseWriter, r *http.Request) {
		served = true
	})

	checkIfHasRootRoute(t, router, POST)

	err := mockServeHTTP(router, POST, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if served != true {
		t.Error("Handler has not been serverd")
	}
}

func TestPOST(t *testing.T) {
	t.Parallel()

	router := New().(*router)
	testBasicMethod(t, router, router.POST, POST)
}

func TestGET(t *testing.T) {
	t.Parallel()

	router := New().(*router)
	testBasicMethod(t, router, router.GET, GET)
}

func TestPUT(t *testing.T) {
	t.Parallel()

	router := New().(*router)
	testBasicMethod(t, router, router.PUT, PUT)
}

func TestDELETE(t *testing.T) {
	t.Parallel()

	router := New().(*router)
	testBasicMethod(t, router, router.DELETE, DELETE)
}

func TestPATCH(t *testing.T) {
	t.Parallel()

	router := New().(*router)
	testBasicMethod(t, router, router.PATCH, PATCH)
}

func TestHEAD(t *testing.T) {
	t.Parallel()

	router := New().(*router)
	testBasicMethod(t, router, router.HEAD, HEAD)
}

func TestCONNECT(t *testing.T) {
	t.Parallel()

	router := New().(*router)
	testBasicMethod(t, router, router.CONNECT, CONNECT)
}

func TestTRACE(t *testing.T) {
	t.Parallel()

	router := New().(*router)
	testBasicMethod(t, router, router.TRACE, TRACE)
}

func TestOPTIONS(t *testing.T) {
	t.Parallel()

	router := New().(*router)
	testBasicMethod(t, router, router.OPTIONS, OPTIONS)

	handler := &mockHandler{}
	router.GET("/x/y", handler)
	router.POST("/x/y", handler)

	checkIfHasRootRoute(t, router, GET)

	w := httptest.NewRecorder()

	// test all routes "*" paths
	req, err := http.NewRequest(OPTIONS, "*", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	if allow := w.Header().Get("Allow"); allow != "POST, GET, OPTIONS" {
		t.Errorf("Allow header incorrect value: %s", allow)
	}

	// test specific path
	req, err = http.NewRequest(OPTIONS, "/x/y", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	if allow := w.Header().Get("Allow"); allow != "POST, GET, OPTIONS" {
		t.Errorf("Allow header incorrect value: %s", allow)
	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	handler := &mockHandler{}
	router := New().(*router)
	router.GET("/x/y", handler)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(POST, "/y/y", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("NotFound error, actual code: %d", w.Code)
	}

	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("test"))
	}))

	if router.notFound == nil {
		t.Error("NotFound handler error")
	}

	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Body.String() != "test" {
		t.Error("Not found handler wasn't invoked")
	}
}

func TestNotAllowed(t *testing.T) {
	t.Parallel()

	handler := &mockHandler{}
	router := New().(*router)
	router.GET("/x/y", handler)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(POST, "/x/y", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Error("NotAllowed doesnt work")
	}

	router.NotAllowed(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("test"))
	}))

	if router.notAllowed == nil {
		t.Error("NotAllowed handler error")
	}

	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Body.String() != "test" {
		t.Error("Not found handler wasn't invoked")
	}

	w = httptest.NewRecorder()
	req, err = http.NewRequest(POST, "*", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	if w.Body.String() != "test" {
		t.Error("Not found handler wasn't invoked")
	}
}

func TestParam(t *testing.T) {
	t.Parallel()

	router := New().(*router)

	serverd := false
	router.GET("/x/{param}", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("param") != "y" {
			t.Errorf("Wrong params value. Expected 'x', actual '%s'", params.Value("param"))
		}
	}))

	err := mockServeHTTP(router, GET, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestRegexpParam(t *testing.T) {
	t.Parallel()

	router := New().(*router)

	serverd := false
	router.GET("/x/{param:r([a-z]+)go}", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("param") != "rxgo" {
			t.Errorf("Wrong params value. Expected 'rxgo', actual '%s'", params.Value("param"))
		}
	}))

	err := mockServeHTTP(router, GET, "/x/rxgo")
	if err != nil {
		t.Fatal(err)
	}

	if serverd != true {
		t.Error("Handler has not been serverd")
	}
}

func TestEmptyParam(t *testing.T) {
	t.Parallel()

	paniced := false
	defer func() {
		if rcv := recover(); rcv != nil {
			paniced = true
		}
	}()

	handler := &mockHandler{}
	router := New().(*router)

	router.GET("/x/{}", handler)

	if paniced != true {
		t.Error("Router should panic for empty wildcard path")
	}
}

func TestServeFiles(t *testing.T) {
	t.Parallel()

	mfs := &mockFileSystem{}
	router := New().(*router)

	router.ServeFiles(mfs, "static", true)

	if router.fileServer == nil {
		t.Error("File serve handler error")
	}

	w := httptest.NewRecorder()
	r, err := http.NewRequest(GET, "/favicon.ico", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Error("File should not exist")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Router should panic for empty path")
		}
	}()

	router.ServeFiles(mfs, "", true)
}

func TestNilMiddleware(t *testing.T) {
	t.Parallel()

	router := New().(*router)

	router.GET("/x/{param}", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("test"))
	}))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, "/x/y", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	if w.Body.String() != "test" {
		t.Error("Nil middleware works")
	}
}

func TestPanicMiddleware(t *testing.T) {
	t.Parallel()

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

	router := New(panicMiddleware).(*router)

	router.GET("/x/{param}", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("test panic recover")
	}))

	err := mockServeHTTP(router, GET, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if paniced != true {
		t.Error("Panic has not been handled")
	}
}

func TestNodeApplyMiddleware(t *testing.T) {
	t.Parallel()

	router := New().(*router)

	router.GET("/x/{param}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		w.Write([]byte(params.Value("param")))
	}))

	router.USE(GET, "/x/{param}", mockMiddleware("m"))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, "/x/y", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	if w.Body.String() != "my" {
		t.Errorf("Use global middleware error %s", w.Body.String())
	}
}

func TestChainCalls(t *testing.T) {
	t.Parallel()

	router := New().(*router)

	serverd := false
	router.GET("/users/{user:[a-z0-9]+)}/starred", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("user") != "x" {
			t.Errorf("Wrong params value. Expected 'x', actual '%s'", params.Value("user"))
		}
	}))

	router.GET("/applications/{client_id}/tokens", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("client_id") != "client_id" {
			t.Errorf("Wrong params value. Expected 'client_id', actual '%s'", params.Value("client_id"))
		}
	}))

	router.GET("/applications/{client_id}/tokens/{access_token}", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
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

	router.GET("/users/{user}/received_events", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("user") != "user1" {
			t.Errorf("Wrong params value. Expected 'user1', actual '%s'", params.Value("user"))
		}
	}))

	router.GET("/users/{user}/received_events/public", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		serverd = true

		params, ok := FromContext(r.Context())
		if !ok {
			t.Fatal("Error while reading param")
		}

		if params.Value("user") != "user2" {
			t.Errorf("Wrong params value. Expected 'user2', actual '%s'", params.Value("user"))
		}
	}))

	// //FIRST CALL
	err := mockServeHTTP(router, GET, "/users/x/starred")
	if err != nil {
		t.Fatal(err)
	}

	if !serverd {
		t.Fatal("First not served")
	}

	//SECOND CALL
	serverd = false
	err = mockServeHTTP(router, GET, "/applications/client_id/tokens")
	if err != nil {
		t.Fatal(err)
	}

	if !serverd {
		t.Fatal("Second not served")
	}

	//THIRD CALL
	serverd = false
	err = mockServeHTTP(router, GET, "/applications/client_id/tokens/access_token")
	if err != nil {
		t.Fatal(err)
	}

	if !serverd {
		t.Fatal("Third not served")
	}

	//FOURTH CALL
	serverd = false
	err = mockServeHTTP(router, GET, "/users/user1/received_events")
	if err != nil {
		t.Fatal(err)
	}

	if !serverd {
		t.Fatal("Fourth not served")
	}

	//FIFTH CALL
	serverd = false
	err = mockServeHTTP(router, GET, "/users/user2/received_events/public")
	if err != nil {
		t.Fatal(err)
	}

	if !serverd {
		t.Fatal("Fifth not served")
	}
}

func TestMountSubRouter(t *testing.T) {
	t.Parallel()

	rGlobal1 := mockMiddleware("rg1")
	rGlobal2 := mockMiddleware("rg2")
	r := New(rGlobal1, rGlobal2).(*router)
	r.GET("/{param}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("r"))
	}))

	sGlobal1 := mockMiddleware("sg1")
	sGlobal2 := mockMiddleware("sg2")
	router := New(sGlobal1, sGlobal2).(*router)
	router.GET("/y", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("s"))
	}))

	r.Mount("/{param}", router)

	r1 := mockMiddleware("r1")
	r2 := mockMiddleware("r2")
	r.USE(GET, "/{param}", r1)
	r.USE(GET, "/{param}", r2)

	s1 := mockMiddleware("s1")
	s2 := mockMiddleware("s2")
	router.USE(GET, "/y", s1)
	router.USE(GET, "/y", s2)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(GET, "/x/y", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(w, req)

	if w.Body.String() != "rg1rg2r1r2sg1sg2s1s2s" {
		t.Errorf("Router mount sub router middleware error: %s", w.Body.String())
	}
}
