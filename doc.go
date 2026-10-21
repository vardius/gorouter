// Package gorouter is a fast HTTP request router for net/http.
//
// A Router matches a request's method and path against registered patterns and calls the matching handler.
// Matching a route and running its middleware allocate nothing. What a request can allocate is net/http's:
// storing param values for r.PathValue, copying the request for a mount as http.StripPrefix does,
// the Allow header of a 405, and whatever the 404 and 405 handlers write.
//
//	r := gorouter.New()
//	r.GET("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, "hello, %s!", r.PathValue("name"))
//	})
//	log.Fatal(http.ListenAndServe(":8080", r))
//
// For fasthttp, use the router in the fastrouter subpackage. It has the same API and the same matching rules,
// applied to the path fasthttp hands over, which it has already cleaned up: it collapses // and resolves ..,
// so /users//7 matches /users/{id} there, while on net/http the empty segment matches nothing.
//
// # Patterns
//
// A pattern is a path made of segments separated by slashes. Each segment is one of:
//
//   - a literal, such as users, which matches only itself;
//   - a param, {name}, which matches any one segment;
//   - a regexp param, {name:regexp}, which matches a segment the whole regexp matches, such as {id:[0-9]+}.
//
// A trailing slash is ignored, so /users and /users/ are the same route, and an empty segment never matches.
// A regexp can't contain a slash: it always matches within one segment. Patterns match the decoded path,
// so an escaped slash, %2F, separates segments like any other slash.
//
// When several patterns match, literals win over regexp params, regexp params are tried in the order
// they were registered, and a plain param comes last. If the preferred branch fails further down the path,
// the next one is tried, so /users/new and /users/{id}/posts can live side by side.
//
// # Params
//
// Param values are set on the request with Request.SetPathValue, so handlers read them with r.PathValue,
// as they would with http.ServeMux. Request.Pattern is set to the matched route's pattern.
//
// # Middleware
//
// Middleware has the usual func(http.Handler) http.Handler shape. Use wraps the routes below a pattern,
// UseMethod narrows that to one method, and WithMiddleware wraps the whole router, including
// the 404, 405 and OPTIONS responses. Middleware is composed when it is registered, never per request,
// and runs in registration order: the first one registered is the outermost.
//
// Use picks routes by their pattern, not by the request path: a request runs the middleware of the route
// it matches. After Use("/admin", auth), a request for /admin/users that matches the route /{section}/users
// doesn't run auth, because that route's pattern doesn't start with /admin.
//
// # Mounting
//
// Mount hands every request at or below a pattern, on any method, to another http.Handler, with the pattern's
// segments removed from the path. Mount another Router to split an API into parts, or an http.FileServer
// to serve files. Middleware used on the mount's pattern or above it wraps the mount; middleware for paths
// below it belongs in the mounted handler, and Use panics if asked to put it there.
//
// # 404, 405 and OPTIONS
//
// A path no route matches gets 404 Not Found. A path that matches on other methods gets
// 405 Method Not Allowed with an Allow header listing them. An OPTIONS request is answered with
// the Allow header and an empty 200 response, unless an OPTIONS handler is registered for that pattern.
// WithNotFound and WithMethodNotAllowed replace the 404 and 405 handlers.
//
// # Concurrency
//
// Register every route and middleware before serving. Once serving, a Router is safe for concurrent use.
package gorouter
