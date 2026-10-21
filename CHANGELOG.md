# Changelog

## v5.0.0 (unreleased)

A rewrite with a smaller API. See [MIGRATION.md](MIGRATION.md) for the before and after of every call.

### Added

- `fastrouter` package for fasthttp, so `net/http` users no longer link fasthttp.
- Functional options for `New`: `WithMiddleware`, `WithNotFound`, `WithMethodNotAllowed`.
- `Use(pattern, mw...)` and `UseMethod(method, pattern, mw...)` for middleware scoped to the routes below a path.
- `Routes()` lists the registered routes.
- Params on `net/http` are read with the standard `r.PathValue`, and `r.Pattern` holds the matched pattern.
- `fastrouter.Param(ctx, name)` reads params on fasthttp.
- Custom methods work everywhere, including the `Allow` header and `Mount`, which takes every method.
- Runnable examples in `examples/`, each with a test.

### Changed

- Module path is `github.com/vardius/gorouter/v5`; Go 1.26 or newer is required.
- Matching allocates nothing. Middleware chains are composed once at registration instead of on every request:
  a route behind three middleware went from 170 ns and 6 allocations to 24 ns and none on `net/http`.
- On fasthttp, params are pooled and released with the request: routing allocates nothing.
  Param values are copies, valid until the handler returns even if it rewrites the path.
- `Use` picks routes by pattern rather than by request path, and panics for paths below a mount, where it could never run.
- On fasthttp, `Mount` decodes the path once, not twice, and puts the caller's path and params back when the mounted handler returns.
- The method helpers (`GET`, `POST`, …) take a handler func; use `Handle` for any other `http.Handler`.
- Regexp params match the whole segment and can't contain a slash.
- Invalid patterns and duplicate routes panic at registration with a message naming the problem.

### Removed

- The `context`, `middleware`, `mux` and `path` packages.
- `Router` and `FastHTTPRouter` interfaces, in favour of the concrete `*gorouter.Router` and `*fastrouter.Router`.
- `NewFastHTTPRouter`, `USE`, `USEANY`, `NotFound`, `NotAllowed`, `ServeFiles`, `Compile` and `PrettyPrint`.

## v4.6.0 (unreleased, `v4` branch)

- Fixed a data race where concurrent requests could run each other's middleware.
- `Compile` no longer drops routes or middleware on nodes it merges.
- Middleware on `/a` no longer applies to `/abc`.
- Regexp params match the whole segment, and an expression may contain a colon.
- An unclosed `{name` panics at registration.
- The root path `/` answers 405 with an `Allow` header when another method has a route there.
- `Mount` over a middleware-only path works, and panics when routes already live under the mount path.
- Go 1.26 and the latest dependencies.
