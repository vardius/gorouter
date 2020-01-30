package gorouter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/vardius/gorouter/v4/context"
)

func buildFastHTTPRequestContext(method, path string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	ctx.URI().SetPath(path)

	return ctx
}

func TestFastHTTPInterface(t *testing.T) {
	var _ fasthttp.RequestHandler = NewFastHTTPRouter().HandleFastHTTP
}

func TestFastHTTPHandle(t *testing.T) {
	t.Parallel()

	handler := &mockHandler{}
	router := NewFastHTTPRouter().(*fastHTTPRouter)
	router.Handle(fasthttp.MethodPost, "/x/y", handler.HandleFastHTTP)

	checkIfHasRootRoute(t, router, fasthttp.MethodPost)

	err := mockHandleFastHTTP(router.HandleFastHTTP, fasthttp.MethodPost, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if handler.served != true {
		t.Error("Handler has not been served")
	}
}

func TestFastHTTPOPTIONSHeaders(t *testing.T) {
	handler := &mockHandler{}
	router := NewFastHTTPRouter().(*fastHTTPRouter)

	router.GET("/x/y", handler.HandleFastHTTP)
	router.POST("/x/y", handler.HandleFastHTTP)

	checkIfHasRootRoute(t, router, fasthttp.MethodGet)

	ctx := buildFastHTTPRequestContext(fasthttp.MethodOptions, "*")

	router.HandleFastHTTP(ctx)

	if allow := string(ctx.Response.Header.Peek("Allow")); !strings.Contains(allow, "POST") || !strings.Contains(allow, "GET") || !strings.Contains(allow, "OPTIONS") {
		t.Errorf("Allow header incorrect value: %s", allow)
	}

	ctx2 := buildFastHTTPRequestContext(fasthttp.MethodOptions, "/x/y")

	router.HandleFastHTTP(ctx2)

	if allow := string(ctx.Response.Header.Peek("Allow")); !strings.Contains(allow, "POST") || !strings.Contains(allow, "GET") || !strings.Contains(allow, "OPTIONS") {
		t.Errorf("Allow header incorrect value: %s", allow)
	}
}

func TestFastHTTPMethods(t *testing.T) {
	t.Parallel()

	for _, method := range []string{
		fasthttp.MethodPost,
		fasthttp.MethodGet,
		fasthttp.MethodPut,
		fasthttp.MethodDelete,
		fasthttp.MethodPatch,
		fasthttp.MethodHead,
		fasthttp.MethodConnect,
		fasthttp.MethodTrace,
		fasthttp.MethodOptions,
	} {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			handler := &mockHandler{}
			router := NewFastHTTPRouter().(*fastHTTPRouter)

			reflect.ValueOf(router).MethodByName(method).Call([]reflect.Value{reflect.ValueOf("/x/y"), reflect.ValueOf(handler.HandleFastHTTP)})

			checkIfHasRootRoute(t, router, method)

			err := mockHandleFastHTTP(router.HandleFastHTTP, method, "/x/y")
			if err != nil {
				t.Fatal(err)
			}

			if handler.served != true {
				t.Error("Handler has not been served")
			}
		})
	}
}

func TestFastHTTPNotFound(t *testing.T) {
	t.Parallel()

	handler := &mockHandler{}
	router := NewFastHTTPRouter().(*fastHTTPRouter)
	router.GET("/x", handler.HandleFastHTTP)
	router.GET("/x/y", handler.HandleFastHTTP)

	ctx := buildFastHTTPRequestContext(fasthttp.MethodGet, "/x/x")

	router.HandleFastHTTP(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusNotFound {
		t.Errorf("NotFound error, actual code: %d", ctx.Response.StatusCode())
	}

	router.NotFound(func(ctx *fasthttp.RequestCtx) {
		if _, err := fmt.Fprintf(ctx, "test"); err != nil {
			t.Fatal(err)
		}
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

	ctx := buildFastHTTPRequestContext(fasthttp.MethodPost, "/x/y")

	router.HandleFastHTTP(ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusMethodNotAllowed {
		t.Error("NotAllowed doesn't work")
	}

	router.NotAllowed(func(ctx *fasthttp.RequestCtx) {
		if _, err := fmt.Fprintf(ctx, "test"); err != nil {
			t.Fatal(err)
		}
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

		params := ctx.UserValue("params").(context.Params)
		if params.Value("param") != "y" {
			t.Errorf("Wrong params value. Expected 'y', actual '%s'", params.Value("param"))
		}
	})

	err := mockHandleFastHTTP(router.HandleFastHTTP, fasthttp.MethodGet, "/x/y")
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

		params := ctx.UserValue("params").(context.Params)
		if params.Value("param") != "rxgo" {
			t.Errorf("Wrong params value. Expected 'rxgo', actual '%s'", params.Value("param"))
		}
	})

	err := mockHandleFastHTTP(router.HandleFastHTTP, fasthttp.MethodGet, "/x/rxgo")
	if err != nil {
		t.Fatal(err)
	}

	if served != true {
		t.Error("Handler has not been served")
	}
}

func TestFastHTTPEmptyParam(t *testing.T) {
	t.Parallel()

	panicked := false
	defer func() {
		if rcv := recover(); rcv != nil {
			panicked = true
		}
	}()

	handler := &mockHandler{}
	router := NewFastHTTPRouter().(*fastHTTPRouter)

	router.GET("/x/{}", handler.HandleFastHTTP)

	if panicked != true {
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
	ctx.Request.Header.SetMethod(fasthttp.MethodGet)
	// will serve files from /var/www/static/favicon.ico
	// because strip 1 value ServeFiles("/var/www/static", 1)
	// /static/favicon.ico -> /favicon.ico
	ctx.URI().SetPath("/static/favicon.ico")

	router.HandleFastHTTP(&ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusNotFound {
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
		if _, err := fmt.Fprintf(ctx, "test"); err != nil {
			t.Fatal(err)
		}
	})

	ctx := buildFastHTTPRequestContext(fasthttp.MethodGet, "/x/y")

	router.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "test" {
		t.Error("Nil middleware works")
	}
}

func TestFastHTTPPanicMiddleware(t *testing.T) {
	t.Parallel()

	panicked := false
	panicMiddleware := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		fn := func(ctx *fasthttp.RequestCtx) {
			defer func() {
				if rcv := recover(); rcv != nil {
					panicked = true
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

	err := mockHandleFastHTTP(router.HandleFastHTTP, fasthttp.MethodGet, "/x/y")
	if err != nil {
		t.Fatal(err)
	}

	if panicked != true {
		t.Error("Panic has not been handled")
	}
}

func TestFastHTTPNodeApplyMiddleware(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	router.GET("/x/{param}", func(ctx *fasthttp.RequestCtx) {
		params := ctx.UserValue("params").(context.Params)
		if _, err := fmt.Fprintf(ctx, "%s", params.Value("param")); err != nil {
			t.Fatal(err)
		}
	})

	router.USE(fasthttp.MethodGet, "/x/{param}", mockFastHTTPMiddleware("m1"))
	router.USE(fasthttp.MethodGet, "/x/x", mockFastHTTPMiddleware("m2"))

	ctx := buildFastHTTPRequestContext(fasthttp.MethodGet, "/x/y")

	router.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "m1y" {
		t.Errorf("Use middleware error %s", string(ctx.Response.Body()))
	}

	ctx = buildFastHTTPRequestContext(fasthttp.MethodGet, "/x/x")

	router.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "m1m2x" {
		t.Errorf("Use middleware error %s", string(ctx.Response.Body()))
	}
}

func TestFastHTTPTreeOrphanMiddlewareOrder(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	router.GET("/x/{param}", func(ctx *fasthttp.RequestCtx) {
		if _, err := fmt.Fprintf(ctx, "handler"); err != nil {
			t.Fatal(err)
		}
	})

	// Method global middleware
	router.USE(fasthttp.MethodGet, "/", mockFastHTTPMiddleware("m1->"))
	router.USE(fasthttp.MethodGet, "/", mockFastHTTPMiddleware("m2->"))
	// Path middleware
	router.USE(fasthttp.MethodGet, "/x", mockFastHTTPMiddleware("mx1->"))
	router.USE(fasthttp.MethodGet, "/x", mockFastHTTPMiddleware("mx2->"))
	router.USE(fasthttp.MethodGet, "/x/y", mockFastHTTPMiddleware("mxy1->"))
	router.USE(fasthttp.MethodGet, "/x/y", mockFastHTTPMiddleware("mxy2->"))
	router.USE(fasthttp.MethodGet, "/x/{param}", mockFastHTTPMiddleware("mparam1->"))
	router.USE(fasthttp.MethodGet, "/x/{param}", mockFastHTTPMiddleware("mparam2->"))
	router.USE(fasthttp.MethodGet, "/x/y", mockFastHTTPMiddleware("mxy3->"))
	router.USE(fasthttp.MethodGet, "/x/y", mockFastHTTPMiddleware("mxy4->"))

	ctx := buildFastHTTPRequestContext(fasthttp.MethodGet, "/x/y")

	router.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "m1->m2->mx1->mx2->mxy1->mxy2->mparam1->mparam2->mxy3->mxy4->handler" {
		t.Errorf("Use middleware error %s", string(ctx.Response.Body()))
	}
}

func TestFastHTTPNodeApplyMiddlewareStatic(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	router.GET("/x/x", func(ctx *fasthttp.RequestCtx) {
		if _, err := fmt.Fprintf(ctx, "x"); err != nil {
			t.Fatal(err)
		}
	})

	router.USE(fasthttp.MethodGet, "/x/x", mockFastHTTPMiddleware("m1"))

	ctx := buildFastHTTPRequestContext(fasthttp.MethodGet, "/x/x")

	router.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "m1x" {
		t.Errorf("Use middleware error %s", string(ctx.Response.Body()))
	}
}

func TestFastHTTPNodeApplyMiddlewareInvalidNodeReference(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	router.GET("/x/{param}", func(ctx *fasthttp.RequestCtx) {
		params := ctx.UserValue("params").(context.Params)
		if _, err := fmt.Fprintf(ctx, "%s", params.Value("param")); err != nil {
			t.Fatal(err)
		}
	})

	router.USE(fasthttp.MethodGet, "/x/x", mockFastHTTPMiddleware("m1"))

	ctx := buildFastHTTPRequestContext(fasthttp.MethodGet, "/x/y")

	router.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "y" {
		t.Errorf("Use middleware error %s", string(ctx.Response.Body()))
	}
}

func TestFastHTTPChainCalls(t *testing.T) {
	t.Parallel()

	router := NewFastHTTPRouter().(*fastHTTPRouter)

	served := false
	router.GET("/users/{user:[a-z0-9]+}/starred", func(ctx *fasthttp.RequestCtx) {
		served = true

		params := ctx.UserValue("params").(context.Params)
		if params.Value("user") != "x" {
			t.Errorf("Wrong params value. Expected 'x', actual '%s'", params.Value("user"))
		}
	})

	router.GET("/applications/{client_id}/tokens", func(ctx *fasthttp.RequestCtx) {
		served = true

		params := ctx.UserValue("params").(context.Params)
		if params.Value("client_id") != "client_id" {
			t.Errorf("Wrong params value. Expected 'client_id', actual '%s'", params.Value("client_id"))
		}
	})

	router.GET("/applications/{client_id}/tokens/{access_token}", func(ctx *fasthttp.RequestCtx) {
		served = true

		params := ctx.UserValue("params").(context.Params)
		if params.Value("client_id") != "client_id" {
			t.Errorf("Wrong params value. Expected 'client_id', actual '%s'", params.Value("client_id"))
		}

		if params.Value("access_token") != "access_token" {
			t.Errorf("Wrong params value. Expected 'access_token', actual '%s'", params.Value("access_token"))
		}
	})

	router.GET("/users/{user}/received_events", func(ctx *fasthttp.RequestCtx) {
		served = true

		params := ctx.UserValue("params").(context.Params)
		if params.Value("user") != "user1" {
			t.Errorf("Wrong params value. Expected 'user1', actual '%s'", params.Value("user"))
		}
	})

	router.GET("/users/{user}/received_events/public", func(ctx *fasthttp.RequestCtx) {
		served = true

		params := ctx.UserValue("params").(context.Params)
		if params.Value("user") != "user2" {
			t.Errorf("Wrong params value. Expected 'user2', actual '%s'", params.Value("user"))
		}
	})

	// //FIRST CALL
	err := mockHandleFastHTTP(router.HandleFastHTTP, fasthttp.MethodGet, "/users/x/starred")
	if err != nil {
		t.Fatal(err)
	}

	if !served {
		t.Fatal("First not served")
	}

	//SECOND CALL
	served = false
	err = mockHandleFastHTTP(router.HandleFastHTTP, fasthttp.MethodGet, "/applications/client_id/tokens")
	if err != nil {
		t.Fatal(err)
	}

	if !served {
		t.Fatal("Second not served")
	}

	//THIRD CALL
	served = false
	err = mockHandleFastHTTP(router.HandleFastHTTP, fasthttp.MethodGet, "/applications/client_id/tokens/access_token")
	if err != nil {
		t.Fatal(err)
	}

	if !served {
		t.Fatal("Third not served")
	}

	//FOURTH CALL
	served = false
	err = mockHandleFastHTTP(router.HandleFastHTTP, fasthttp.MethodGet, "/users/user1/received_events")
	if err != nil {
		t.Fatal(err)
	}

	if !served {
		t.Fatal("Fourth not served")
	}

	//FIFTH CALL
	served = false
	err = mockHandleFastHTTP(router.HandleFastHTTP, fasthttp.MethodGet, "/users/user2/received_events/public")
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
		if _, err := fmt.Fprintf(ctx, "[s]"); err != nil {
			t.Fatal(err)
		}
	})

	mainRouter.Mount("/{param}", subRouter.HandleFastHTTP)

	mainRouter.USE(fasthttp.MethodGet, "/{param}", mockFastHTTPMiddleware("[r1]"))
	mainRouter.USE(fasthttp.MethodGet, "/{param}", mockFastHTTPMiddleware("[r2]"))

	subRouter.USE(fasthttp.MethodGet, "/y", mockFastHTTPMiddleware("[s1]"))
	subRouter.USE(fasthttp.MethodGet, "/y", mockFastHTTPMiddleware("[s2]"))

	ctx := buildFastHTTPRequestContext(fasthttp.MethodGet, "/x/y")

	mainRouter.HandleFastHTTP(ctx)

	if string(ctx.Response.Body()) != "[rg1][rg2][r1][r2][sg1][sg2][s1][s2][s]" {
		t.Errorf("Router mount sub router middleware error: %s", string(ctx.Response.Body()))
	}
}
