// Package fastrouter is a fast HTTP request router for fasthttp.
//
// It has the same API and the same matching rules as the net/http router in the parent gorouter package;
// see its documentation for the pattern syntax, route priority, middleware and mounting.
// Matching a route and reading its params allocates nothing.
//
//	r := fastrouter.New()
//	r.GET("/hello/{name}", func(ctx *fasthttp.RequestCtx) {
//		fmt.Fprintf(ctx, "hello, %s!", fastrouter.Param(ctx, "name"))
//	})
//	log.Fatal(fasthttp.ListenAndServe(":8080", r.HandleFastHTTP))
//
// # Params
//
// Read param values with Param. A value points into the request's buffer, so it is only valid
// until the handler returns or changes the request path: copy it with strings.Clone to keep it longer.
// Param values are released with the request.
package fastrouter
