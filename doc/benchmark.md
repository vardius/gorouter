Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Benchmarks
----------------
The biggest competitor is [HttpRouter](https://github.com/julienschmidt/httprouter). However Goserver allows to use *regex* wildcards and handlers implement the `http.Handler` interface type not like `httprouter.Handle`. The request parameters are passed in the **request context**. Goserver also provides [middleware](middleware.md) system and **the performance is comparable**.

The output
```
BenchmarkGoserverStaticParallel1-4      	50000000	        25.0 ns/op
```
means that the loop ran 50000000 times at a speed of 25.0 ns per loop. What gives around **40000000 req/sec** !
Each benchmark name `BenchmarkGoserverStaticParallel1-4 ` means that test used a `static` or `regexp` route path for each node with a nested level `5`. Where `4` stands for CPU number.

The benchmarks from file [benchmark_test.go](benchmark_test.go).
```
$ go test -bench=. -cpu=4
BenchmarkGoserverStatic1-4              	30000000	        42.5 ns/op
BenchmarkGoserverStatic2-4              	30000000	        43.3 ns/op
BenchmarkGoserverStatic3-4              	30000000	        42.0 ns/op
BenchmarkGoserverStatic5-4              	30000000	        42.1 ns/op
BenchmarkGoserverStatic10-4             	30000000	        43.2 ns/op
BenchmarkGoserverStatic20-4             	30000000	        42.4 ns/op
BenchmarkGoserverWildcard1-4            	 3000000	       390 ns/op
BenchmarkGoserverWildcard2-4            	 3000000	       392 ns/op
BenchmarkGoserverWildcard3-4            	 3000000	       379 ns/op
BenchmarkGoserverWildcard5-4            	 3000000	       387 ns/op
BenchmarkGoserverWildcard10-4           	 3000000	       390 ns/op
BenchmarkGoserverWildcard20-4           	 3000000	       384 ns/op
BenchmarkGoserverRegexp1-4              	 3000000	       392 ns/op
BenchmarkGoserverRegexp2-4              	 3000000	       393 ns/op
BenchmarkGoserverRegexp3-4              	 5000000	       395 ns/op
BenchmarkGoserverRegexp5-4              	 3000000	       388 ns/op
BenchmarkGoserverRegexp10-4             	 5000000	       383 ns/op
BenchmarkGoserverRegexp20-4             	 3000000	       384 ns/op
BenchmarkGoserverStaticParallel1-4      	50000000	        25.0 ns/op
BenchmarkGoserverStaticParallel2-4      	50000000	        24.2 ns/op
BenchmarkGoserverStaticParallel3-4      	50000000	        26.5 ns/op
BenchmarkGoserverStaticParallel5-4      	50000000	        24.1 ns/op
BenchmarkGoserverStaticParallel10-4     	50000000	        24.0 ns/op
BenchmarkGoserverStaticParallel20-4     	50000000	        28.6 ns/op
BenchmarkGoserverWildcardParallel1-4    	 5000000	       230 ns/op
BenchmarkGoserverWildcardParallel2-4    	 5000000	       229 ns/op
BenchmarkGoserverWildcardParallel3-4    	 5000000	       232 ns/op
BenchmarkGoserverWildcardParallel5-4    	 5000000	       225 ns/op
BenchmarkGoserverWildcardParallel10-4   	 5000000	       222 ns/op
BenchmarkGoserverWildcardParallel20-4   	 5000000	       228 ns/op
BenchmarkGoserverRegexpParallel1-4      	 5000000	       239 ns/op
BenchmarkGoserverRegexpParallel2-4      	 5000000	       226 ns/op
BenchmarkGoserverRegexpParallel3-4      	 5000000	       223 ns/op
BenchmarkGoserverRegexpParallel5-4      	 5000000	       231 ns/op
BenchmarkGoserverRegexpParallel10-4     	 5000000	       231 ns/op
BenchmarkGoserverRegexpParallel20-4     	 5000000	       222 ns/op
```
### [Go HTTP Router Benchmark](https://github.com/julienschmidt/go-http-routing-benchmark)
**go-http-routing-benchmark** was runned without writing *parameters* to *request context* in case of comparing native router performance.
#### Memory required only for loading the routing structure for the respective API
| Router       | Static      | GitHub      | Google+    | Parse      |
|:-------------|------------:|------------:|-----------:|-----------:|
| Goserver     | 51016 B     | 87600 B     |  7008 B    | 11712 B    |
| Gorilla Mux  | 670544 B    | 1503424 B   |  71072 B   | 122184 B   |
| HttpRouter   | 21128 B     | 37464 B     |  2712 B    | 4976 B     |

#### ns/op
| | **Goserver** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
|:-------------|-------------:|------------:|--------------:|
| Param        | 94           | 114         | 3836          |
| Param5       | 926          | 458         | 6937          |
| Param20      | 930          | 1460        | 10673         |
| ParamWrite   | 162          | 128         | 3338          |
| GithubStatic | 60           | 45.4        | 15145         |
| GithubParam  | 269          | 329         | 9048          |
| GithubAll    | 55057        | 53880       | 6692893       |
| GPlusStatic  | 44           | 25.5        | 2404          |
| GPlusParam   | 96           | 212         | 4075          |
| GPlus2Params | 184          | 231         | 7407          |
| GPlusAll     | 1777         | 2247        | 56497         |
| ParseStatic  | 44           | 26.2        | 2629          |
| ParseParam   | 113          | 190         | 2772          |
| Parse2Params | 168          | 185         | 3660          |
| ParseAll     | 2976         | 2788        | 104968        |
| StaticAll    | 47875        | 10255       | 1764623       |
#### B/op
| | **Goserver** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
|:-------------|-----------:|------------:|-----------:|
| Param        | 32         | 32          | 1056       |
| Param5       | 848        | 160         | 1184       |
| Param20      | 848        | 640         | 3548       |
| ParamWrite   | 40         | 32          | 1064       |
| GithubStatic | 0          | 0           | 736        |
| GithubParam  | 64         | 96          | 1088       |
| GithubAll    | 10848      | 13792       | 211840     |
| GPlusStatic  | 0          | 0           | 736        |
| GPlusParam   | 32         | 64          | 1056       |
| GPlus2Params | 64         | 64          | 1088       |
| GPlusAll     | 512        | 640         | 13296      |
| ParseStatic  | 0          | 0           | 752        |
| ParseParam   | 32         | 64          | 1088       |
| Parse2Params | 64         | 64          | 1088       |
| ParseAll     | 608        | 640         | 24864      |
| StaticAll    | 0          | 0           | 115648     |
#### allocs/op
| | **Goserver** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
|:-------------|---------:|------------:|-------------:|
| Param        | 1        | 1           | 11           |
| Param5       | 7        | 1           | 11           |
| Param20      | 7        | 1           | 13           |
| ParamWrite   | 2        | 1           | 12           |
| GithubStatic | 0        | 0           | 10           |
| GithubParam  | 1        | 1           | 11           |
| GithubAll    | 167      | 167         | 2272         |
| GPlusStatic  | 0        | 0           | 10           |
| GPlusParam   | 1        | 1           | 11           |
| GPlus2Params | 1        | 1           | 11           |
| GPlusAll     | 11       | 11          | 142          |
| ParseStatic  | 0        | 0           | 11           |
| ParseParam   | 1        | 1           | 12           |
| Parse2Params | 1        | 1           | 11           |
| ParseAll     | 16       | 16          | 292          |
| StaticAll    | 0        | 0           | 1578         |
