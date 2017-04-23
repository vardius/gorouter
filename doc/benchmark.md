Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Benchmarks
----------------
The biggest competitor is [HttpRouter](https://github.com/julienschmidt/httprouter). However Goserver allows to use *regex* wildcards and handlers implement the `http.Handler` interface type not like `httprouter.Handle`. The request parameters are passed in the **request context**. Goserver also provides [middleware](middleware.md) system and **the performance is comparable**.

The output
```
BenchmarkStrictParallel1-4     	30000000	       58.1 ns/op
```
means that the loop ran 30000000 times at a speed of 58.1 ns per loop. What gives around **17211704 req/sec** !
Each benchmark name `BenchmarkStrictParallel1-4 ` means that test used a `static` or `regexp` route path for each node with a nested level `5`. Where `4` stands for CPU number.

The benchmarks from file [benchmark_test.go](benchmark_test.go).
```
$ go test -bench=. -cpu=4
BenchmarkStrict1-4             	20000000	        99.1 ns/op
BenchmarkStrict2-4             	10000000	       165 ns/op
BenchmarkStrict3-4             	10000000	       207 ns/op
BenchmarkStrict5-4             	 5000000	       273 ns/op
BenchmarkStrict10-4            	 3000000	       458 ns/op
BenchmarkStrict100-4           	  300000	      4195 ns/op
BenchmarkStrictParallel1-4     	30000000	        58.1 ns/op
BenchmarkStrictParallel2-4     	20000000	        97.6 ns/op
BenchmarkStrictParallel3-4     	10000000	       123 ns/op
BenchmarkStrictParallel5-4     	10000000	       160 ns/op
BenchmarkStrictParallel10-4    	 5000000	       269 ns/op
BenchmarkStrictParallel100-4   	  500000	      2639 ns/op
BenchmarkRegexp1-4             	 2000000	       668 ns/op
BenchmarkRegexp2-4             	 1000000	      1291 ns/op
BenchmarkRegexp3-4             	 1000000	      1865 ns/op
BenchmarkRegexp5-4             	  500000	      2968 ns/op
BenchmarkRegexp10-4            	  200000	      6010 ns/op
BenchmarkRegexp100-4           	   20000	     63167 ns/op
BenchmarkRegexpParallel1-4     	 2000000	       695 ns/op
BenchmarkRegexpParallel2-4     	 1000000	      1211 ns/op
BenchmarkRegexpParallel3-4     	 1000000	      1543 ns/op
BenchmarkRegexpParallel5-4     	  500000	      2348 ns/op
BenchmarkRegexpParallel10-4    	  300000	      4748 ns/op
BenchmarkRegexpParallel100-4   	   50000	     36527 ns/op
```
### [Go HTTP Router Benchmark](https://github.com/julienschmidt/go-http-routing-benchmark)
#### Memory required only for loading the routing structure for the respective API
| Router       | Static      | GitHub      | Google+    | Parse      |
|:-------------|------------:|------------:|-----------:|-----------:|
| Goserver     | __19592 B__ | __34888 B__ |  2792 B    | 5296 B     |
| Gorilla Mux  | 670544 B    | 1503424 B   |  71072 B   | 122184 B   |
| HttpRouter   | 21128 B     | 37464 B     | __2712 B__ | __4976 B__ |

#### ns/op
| | **Goserver** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
|:-------------|-------------:|------------:|--------------:|
| Param        | 189          | 114         | 3836          |
| Param5       | 376          | 458         | 6937          |
| Param20      | 998          | 1460        | 10673         |
| ParamWrite   | 934          | 128         | 3338          |
| GithubStatic | 203          | 45.4        | 15145         |
| GithubParam  | 405          | 329         | 9048          |
| GithubAll    | 80826        | 53880       | 6692893       |
| GPlusStatic  | 94           | 25.5        | 2404          |
| GPlusParam   | 200          | 212         | 4075          |
| GPlus2Params | 383          | 231         | 7407          |
| GPlusAll     | 3397         | 2247        | 56497         |
| ParseStatic  | 162          | 26.2        | 2629          |
| ParseParam   | 236          | 190         | 2772          |
| Parse2Params | 305          | 185         | 3660          |
| ParseAll     | 6517         | 2788        | 104968        |
| StaticAll    | 44640        | 10255       | 1764623       |
#### B/op
| | **Goserver** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
|:-------------|-----------:|------------:|-----------:|
| Param        | 64         | 32          | 1056       |
| Param5       | 240        | 160         | 1184       |
| Param20      | 960        | 640         | 3548       |
| ParamWrite   | 880        | 32          | 1064       |
| GithubStatic | 32         | 0           | 736        |
| GithubParam  | 128        | 96          | 1088       |
| GithubAll    | 23056      | 13792       | 211840     |
| GPlusStatic  | 16         | 0           | 736        |
| GPlusParam   | 64         | 64          | 1056       |
| GPlus2Params | 128        | 64          | 1088       |
| GPlusAll     | 1088       | 640         | 13296      |
| ParseStatic  | 32         | 0           | 752        |
| ParseParam   | 80         | 64          | 1088       |
| Parse2Params | 128        | 64          | 1088       |
| ParseAll     | 1744       | 640         | 24864      |
| StaticAll    | 4848       | 0           | 115648     |
#### allocs/op
| | **Goserver** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
|:-------------|---------:|------------:|-------------:|
| Param        | 2        | 1           | 11           |
| Param5       | 2        | 1           | 11           |
| Param20      | 2        | 1           | 13           |
| ParamWrite   | 8        | 1           | 12           |
| GithubStatic | 1        | 0           | 10           |
| GithubParam  | 2        | 1           | 11           |
| GithubAll    | 370      | 167         | 2272         |
| GPlusStatic  | 1        | 0           | 10           |
| GPlusParam   | 2        | 1           | 11           |
| GPlus2Params | 2        | 1           | 11           |
| GPlusAll     | 24       | 11          | 142          |
| ParseStatic  | 1        | 0           | 11           |
| ParseParam   | 2        | 1           | 12           |
| Parse2Params | 2        | 1           | 11           |
| ParseAll     | 42       | 16          | 292          |
| StaticAll    | 156      | 0           | 1578         |
