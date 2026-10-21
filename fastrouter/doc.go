// Package fastrouter is a fast HTTP request router for fasthttp.
//
// It has the same API and the same matching rules as the net/http router in the parent gorouter package;
// see its documentation for the pattern syntax, route priority, middleware and mounting.
// The rules apply to the path fasthttp hands over, which it has already cleaned up: it collapses //
// and resolves .., so /users//7 matches /users/{id} here but not on net/http.
// Matching a route, reading its params and mounting allocate nothing.
//
//	r := fastrouter.New()
//	r.GET("/hello/{name}", func(ctx *fasthttp.RequestCtx) {
//		fmt.Fprintf(ctx, "hello, %s!", fastrouter.Param(ctx, "name"))
//	})
//	log.Fatal(fasthttp.ListenAndServe(":8080", r.HandleFastHTTP))
//
// # Params
//
// Read param values with Param. Values live in a buffer that is reused once the request is done,
// so a value is valid until the handler returns, even if the handler rewrites the path:
// copy it with strings.Clone to keep it longer.
//
// A Mount rewrites the path for the mounted handler and puts it back once that handler returns,
// along with the params, so middleware around the mount sees the request it was given.
package fastrouter
