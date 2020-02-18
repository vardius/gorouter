🍃 gorouter
================
[![Build Status](https://travis-ci.com/vardius/gorouter.svg?branch=master)](https://travis-ci.com/vardius/gorouter)
[![Go Report Card](https://goreportcard.com/badge/github.com/vardius/gorouter)](https://goreportcard.com/report/github.com/vardius/gorouter)
[![codecov](https://codecov.io/gh/vardius/gorouter/branch/master/graph/badge.svg)](https://codecov.io/gh/vardius/gorouter)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvardius%2Fgorouter.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvardius%2Fgorouter?ref=badge_shield)
[![](https://godoc.org/github.com/vardius/gorouter?status.svg)](http://godoc.org/github.com/vardius/gorouter)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/vardius/gorouter/blob/master/LICENSE.md)

<img align="right" height="180px" src="website/src/static/img/logo.png" alt="logo" />

Go Server/API micro framework, HTTP request router, multiplexer, mux.

📖 ABOUT
==================================================
Contributors:

* [Rafał Lorenz](https://rafallorenz.com)

Want to contribute ? Feel free to send pull requests!

Have problems, bugs, feature ideas?
We are using the github [issue tracker](https://github.com/vardius/gorouter/issues) to manage them.

## 📚 Documentation

For **documentation** (_including examples_), **visit [rafallorenz.com/gorouter](https://rafallorenz.com/gorouter)**

For **GoDoc** reference, **visit [godoc.org/github.com/vardius/gorouter](http://godoc.org/github.com/vardius/gorouter)**

## 🚅 Benchmark

```go
➜  gorouter git:(master) ✗ go test -bench=. -cpu=4 -benchmem
test
goos: darwin
goarch: amd64
pkg: github.com/vardius/gorouter/v4
BenchmarkNetHTTP-4              	65005786	        17.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkFastHTTP-4             	69810878	        16.5 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/vardius/gorouter/v4	3.808s
```

👉 **[Click here](https://rafallorenz.com/gorouter/docs/benchmark)** to see all benchmark results.

## Features
- Routing System
- Middleware System
- Authentication
- Fast HTTP
- Serving Files
- Multidomain
- HTTP2 Support
- Low memory usage
- [Documentation](https://rafallorenz.com/gorouter/)

🚏 HOW TO USE
==================================================

- [Basic example](https://rafallorenz.com/gorouter/docs/basic-example)
- [net/http](https://rafallorenz.com/gorouter/docs/basic-example#nethttp)
- [valyala/fasthttp](https://rafallorenz.com/gorouter/docs/basic-example#fasthttp)

## 🖥️ API example setup

- **[Go Server/API boilerplate](https://github.com/vardius/go-api-boilerplate)** using best practices DDD CQRS ES.

📜 [License](LICENSE.md)
-------

This package is released under the MIT license. See the complete license in the package:

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvardius%2Fgorouter.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvardius%2Fgorouter?ref=badge_large)
