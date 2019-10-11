package gorouter

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
)

func buildFastHTTPRequestContext(method, path string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.URI().SetPath(path)

	return ctx
}

func testBasicFastHTTPMethod(t *testing.T, router *fastHTTPRouter, h func(pattern string, handler fasthttp.RequestHandler), method string) {
	handler := &mockHandler{}
	h("/x/y", handler.HandleFastHTTP)

	checkIfHasRootRoute(t, router, method)

	err := mockHandleFastHTTP(router.HandleFastHTTP, method, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if handler.served != true {
		t.Error("Handler has not been served")
	}
}

func TestFastHTTPInterface(t *testing.T) {
	var _ fasthttp.RequestHandler = NewFastHTTPRouter().HandleFastHTTP
}

func TestFastHTTPHandle(t *testing.T) {
	t.Parallel()

	handler := &mockHandler{}
	router := NewFastHTTPRouter().(*fastHTTPRouter)
	router.Handle(POST, "/x/y", handler.HandleFastHTTP)

	checkIfHasRootRoute(t, router, POST)

	err := mockHandleFastHTTP(router.HandleFastHTTP, POST, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if handler.served != true {
		t.Error("Handler has not been served")
	}
}

func TestFastHTTPPOST(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)
	testBasicFastHTTPMethod(t, router, router.POST, POST)
}

func TestFastHTTPGET(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)
	testBasicFastHTTPMethod(t, router, router.GET, GET)
}

func TestFastHTTPPUT(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)
	testBasicFastHTTPMethod(t, router, router.PUT, PUT)
}

func TestFastHTTPDELETE(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)
	testBasicFastHTTPMethod(t, router, router.DELETE, DELETE)
}

func TestFastHTTPPATCH(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)
	testBasicFastHTTPMethod(t, router, router.PATCH, PATCH)
}

func TestFastHTTPHEAD(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)
	testBasicFastHTTPMethod(t, router, router.HEAD, HEAD)
}

func TestFastHTTPCONNECT(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)
	testBasicFastHTTPMethod(t, router, router.CONNECT, CONNECT)
}

func TestFastHTTPTRACE(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)
	testBasicFastHTTPMethod(t, router, router.TRACE, TRACE)
}

func TestFastHTTPOPTIONS(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)
	testBasicFastHTTPMethod(t, router, router.OPTIONS, OPTIONS)

	handler := &mockHandler{}
	router.GET("/x/y", handler.HandleFastHTTP)
	router.POST("/x/y", handler.HandleFastHTTP)

	checkIfHasRootRoute(t, router, GET)

	ctx := buildFastHTTPRequestContext(OPTIONS, "*")

	router.HandleFastHTTP(ctx)

	if allow := string(ctx.Response.Header.Peek("Allow")); !strings.Contains(allow, "POST") || !strings.Contains(allow, "GET") || !strings.Contains(allow, "OPTIONS") {
		t.Errorf("Allow header incorrect value: %s", allow)
	}

	ctx2 := buildFastHTTPRequestContext(OPTIONS, "/x/y")

	router.HandleFastHTTP(ctx2)

	if allow := string(ctx.Response.Header.Peek("Allow")); !strings.Contains(allow, "POST") || !strings.Contains(allow, "GET") || !strings.Contains(allow, "OPTIONS") {
		t.Errorf("Allow header incorrect value: %s", allow)
	}
}

func TestFastHTTPNotFound(t *testing.T) {
	t.Parallel()

	handler := &mockHandler{}
	router := NewFastHTTPRouter().(*fastHTTPRouter)
	router.GET("/x", handler.HandleFastHTTP)
	router.GET("/x/y", handler.HandleFastHTTP)

	ctx := buildFastHTTPRequestContext(GET, "/x/x")

	router.HandleFastHTTP(ctx)

	if ctx.Response.StatusCode() != http.StatusNotFound {
		t.Errorf("NotFound error, actual code: %d", ctx.Response.StatusCode())
	}

	router.NotFound(func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "test")
	})

	if router.notFound == nil {
		t.Error("NotFound handler error")
	}

	ctx.ResetBody()

	router.HandleFastHTTP(ctx)

	fmt.Println(string(ctx.Response.Body()))

	if string(ctx.Response.Body()) != "test" {
		t.Error("Not found handler wasn't invoked")
	}
}

func TestFastHTTPNotAllowed(t *testing.T) {
	t.Parallel()

	handler := &mockHandler{}
	router := NewFastHTTPRouter().(*fastHTTPRouter)
	router.GET("/x/y", handler.HandleFastHTTP)

	ctx := buildFastHTTPRequestContext(POST, "/x/y")

	router.HandleFastHTTP(ctx)

	if ctx.Response.StatusCode() != http.StatusMethodNotAllowed {
		t.Error("NotAllowed doesn't work")
	}

	router.NotAllowed(func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "test")
	})

	if router.notAllowed == nil {
		t.Error("NotAllowed handler error")
	}

	ctx.ResetBody()

	router.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "test" {
		t.Errorf("NotAllowed handler wasn't invoked (%s)", string(ctx.Response.Body()))
	}

	ctx.ResetBody()

	router.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "test" {
		t.Errorf("NotAllowed handler wasn't invoked (%s)", string(ctx.Response.Body()))
	}
}

func TestFastHTTPParam(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	served := false
	router.GET("/x/{param}", func(ctx *fasthttp.RequestCtx) {
		served = true

		if ctx.UserValue("param") != "y" {
			t.Errorf("Wrong params value. Expected 'y', actual '%s'", ctx.UserValue("param"))
		}
	})

	err := mockHandleFastHTTP(router.HandleFastHTTP, GET, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if served != true {
		t.Error("Handler has not been served")
	}
}

func TestFastHTTPRegexpParam(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	served := false
	router.GET("/x/{param:r([a-z]+)go}", func(ctx *fasthttp.RequestCtx) {
		served = true

		if ctx.UserValue("param") != "rxgo" {
			t.Errorf("Wrong params value. Expected 'rxgo', actual '%s'", ctx.UserValue("param"))
		}
	})

	err := mockHandleFastHTTP(router.HandleFastHTTP, GET, "/x/rxgo")
	if err != nil {
		t.Fatal(err)
	}

	if served != true {
		t.Error("Handler has not been served")
	}
}

func TestFastHTTPEmptyParam(t *testing.T) {
	t.Parallel()

	paniced := false
	defer func() {
		if rcv := recover(); rcv != nil {
			paniced = true
		}
	}()

	handler := &mockHandler{}
	router := NewFastHTTPRouter().(*fastHTTPRouter)

	router.GET("/x/{}", handler.HandleFastHTTP)

	if paniced != true {
		t.Error("Router should panic for empty wildcard path")
	}
}

func TestFastHTTPServeFiles(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	router.ServeFiles("/var/www/static", 1)

	if router.fileServer == nil {
		t.Error("File server handler error")
	}
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	ctx.Init(&req, nil, testLogger{t})
	ctx.Request.Header.SetMethod(GET)
	// will serve files from /var/www/static/favicon.ico
	// because strip 1 value ServeFiles("/var/www/static", 1)
	// /static/favicon.ico -> /favicon.ico
	ctx.URI().SetPath("/static/favicon.ico")

	router.HandleFastHTTP(&ctx)

	if ctx.Response.StatusCode() != http.StatusNotFound {
		t.Error("File should not exist")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Router should panic for empty path")
		}
	}()

	router.ServeFiles("", 0)
}

func TestFastHTTPNilMiddleware(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	router.GET("/x/{param}", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "test")
	})

	ctx := buildFastHTTPRequestContext(GET, "/x/y")

	router.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "test" {
		t.Error("Nil middleware works")
	}
}

func TestFastHTTPPanicMiddleware(t *testing.T) {
	t.Parallel()

	paniced := false
	panicMiddleware := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		fn := func(ctx *fasthttp.RequestCtx) {
			defer func() {
				if rcv := recover(); rcv != nil {
					paniced = true
				}
			}()

			next(ctx)
		}

		return fn
	}

	router := NewFastHTTPRouter(panicMiddleware).(*fastHTTPRouter)

	router.GET("/x/{param}", func(ctx *fasthttp.RequestCtx) {
		panic("test panic recover")
	})

	err := mockHandleFastHTTP(router.HandleFastHTTP, GET, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if paniced != true {
		t.Error("Panic has not been handled")
	}
}

func TestFastHTTPNodeApplyMiddleware(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	router.GET("/x/{param}", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, ctx.UserValue("param").(string))
	})

	router.USE(GET, "/x/{param}", mockFastHTTPMiddleware("m"))

	ctx := buildFastHTTPRequestContext(GET, "/x/y")

	router.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "my" {
		t.Errorf("Use global middleware error %s", string(ctx.Response.Body()))
	}
}

func TestFastHTTPChainCalls(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	served := false
	router.GET("/users/{user:[a-z0-9]+}/starred", func(ctx *fasthttp.RequestCtx) {
		served = true

		if ctx.UserValue("user") != "x" {
			t.Errorf("Wrong params value. Expected 'x', actual '%s'", ctx.UserValue("user"))
		}
	})

	router.GET("/applications/{client_id}/tokens", func(ctx *fasthttp.RequestCtx) {
		served = true

		if ctx.UserValue("client_id") != "client_id" {
			t.Errorf("Wrong params value. Expected 'client_id', actual '%s'", ctx.UserValue("client_id"))
		}
	})

	router.GET("/applications/{client_id}/tokens/{access_token}", func(ctx *fasthttp.RequestCtx) {
		served = true

		if ctx.UserValue("client_id") != "client_id" {
			t.Errorf("Wrong params value. Expected 'client_id', actual '%s'", ctx.UserValue("client_id"))
		}

		if ctx.UserValue("access_token") != "access_token" {
			t.Errorf("Wrong params value. Expected 'access_token', actual '%s'", ctx.UserValue("access_token"))
		}
	})

	router.GET("/users/{user}/received_events", func(ctx *fasthttp.RequestCtx) {
		served = true

		if ctx.UserValue("user") != "user1" {
			t.Errorf("Wrong params value. Expected 'user1', actual '%s'", ctx.UserValue("user"))
		}
	})

	router.GET("/users/{user}/received_events/public", func(ctx *fasthttp.RequestCtx) {
		served = true

		if ctx.UserValue("user") != "user2" {
			t.Errorf("Wrong params value. Expected 'user2', actual '%s'", ctx.UserValue("user"))
		}
	})

	// //FIRST CALL
	err := mockHandleFastHTTP(router.HandleFastHTTP, GET, "/users/x/starred")
	if err != nil {
		t.Fatal(err)
	}

	if !served {
		t.Fatal("First not served")
	}

	//SECOND CALL
	served = false
	err = mockHandleFastHTTP(router.HandleFastHTTP, GET, "/applications/client_id/tokens")
	if err != nil {
		t.Fatal(err)
	}

	if !served {
		t.Fatal("Second not served")
	}

	//THIRD CALL
	served = false
	err = mockHandleFastHTTP(router.HandleFastHTTP, GET, "/applications/client_id/tokens/access_token")
	if err != nil {
		t.Fatal(err)
	}

	if !served {
		t.Fatal("Third not served")
	}

	//FOURTH CALL
	served = false
	err = mockHandleFastHTTP(router.HandleFastHTTP, GET, "/users/user1/received_events")
	if err != nil {
		t.Fatal(err)
	}

	if !served {
		t.Fatal("Fourth not served")
	}

	//FIFTH CALL
	served = false
	err = mockHandleFastHTTP(router.HandleFastHTTP, GET, "/users/user2/received_events/public")
	if err != nil {
		t.Fatal(err)
	}

	if !served {
		t.Fatal("Fifth not served")
	}
}

func TestFastHTTPMountSubRouter(t *testing.T) {
	t.Parallel()

	mainRouter := NewFastHTTPRouter(
		mockFastHTTPMiddleware("[rg1]"),
		mockFastHTTPMiddleware("[rg2]"),
	).(*fastHTTPRouter)

	subRouter := NewFastHTTPRouter(
		mockFastHTTPMiddleware("[sg1]"),
		mockFastHTTPMiddleware("[sg2]"),
	).(*fastHTTPRouter)

	subRouter.GET("/y", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "[s]")
	})

	mainRouter.Mount("/{param}", subRouter.HandleFastHTTP)

	mainRouter.USE(GET, "/{param}", mockFastHTTPMiddleware("[r1]"))
	mainRouter.USE(GET, "/{param}", mockFastHTTPMiddleware("[r2]"))

	subRouter.USE(GET, "/y", mockFastHTTPMiddleware("[s1]"))
	subRouter.USE(GET, "/y", mockFastHTTPMiddleware("[s2]"))

	ctx := buildFastHTTPRequestContext(GET, "/x/y")

	mainRouter.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "[rg1][rg2][r1][r2][sg1][sg2][s1][s2][s]" {
		t.Errorf("Router mount sub router middleware error: %s", string(ctx.Response.Body()))
	}
}
