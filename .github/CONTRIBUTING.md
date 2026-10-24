# Contributing

Thanks for helping improve **gorouter**. By contributing, you agree to abide by the [code of conduct](CODE_OF_CONDUCT.md).

## Issues and questions

Search the [issue tracker](https://github.com/vardius/gorouter/issues) first; if nothing matches, open an issue.
Bugs, ideas and usage questions are all welcome. For a bug, include the routes you registered and the request that misbehaved.

## Development

You need Go 1.26 or newer.

```sh
git clone https://github.com/vardius/gorouter.git
cd gorouter
go test ./...
```

Before sending a pull request, run what CI runs:

```sh
gofmt -l .
go vet ./...
go run honnef.co/go/tools/cmd/staticcheck@v0.8.2 ./...
go test -race -shuffle=on ./...
go test -run '^$' -fuzz FuzzLookup -fuzztime 30s ./internal/tree
(cd examples/grpc && go vet ./... && go run honnef.co/go/tools/cmd/staticcheck@v0.8.2 ./... && go test -race -shuffle=on ./...)
```

`examples/grpc` is a module of its own, so its gRPC dependencies stay out of the router's `go.mod`, and `./...` at the root skips it.

### Layout

- `internal/tree` is the routing core shared by both routers: pattern parsing, matching and middleware composition.
  `FuzzLookup` checks it against a deliberately naive reference matcher, so a change to matching rules goes in both.
- The root package is the `net/http` router; `fastrouter` is the fasthttp router. They mirror each other.
- `examples/` holds runnable programs, each with a test. The website shows them, so keep them short.
- `website/` is the documentation site.

### Performance

Matching must not allocate. `TestAllocations` in each router package fails if a request costs more allocations than it should;
if your change moves a number, explain why in the pull request. To compare before and after:

```sh
go test -run '^$' -bench . -count 10 ./... > new.txt
```

and compare with the same run on `master`, for example with `golang.org/x/perf/cmd/benchstat`.

### Website

```sh
cd website
npm ci
npm run dev
```

## Pull requests

- For a large change, open an issue first so we can agree on the approach.
- Branch from `master`, keep the change focused, and add a test that fails without it.
- Fixes for v4 go to the `v4` branch.
