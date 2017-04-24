Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Benchmarks
----------------
The biggest competitor is [HttpRouter](https://github.com/julienschmidt/httprouter). However Goserver allows to use *regex* wildcards and handlers implement the `http.Handler` interface type not like `httprouter.Handle`. The request parameters are passed in the **request context**. Goserver also provides [middleware](middleware.md) system and **the performance is comparable**.

The output
```
BenchmarkGoserverStaticParallel1-4      	30000000	        56.1 ns/op
```
means that the loop ran 30000000 times at a speed of 56.1 ns per loop. What gives around **17825312 req/sec** !
Each benchmark name `BenchmarkGoserverStaticParallel1-4 ` means that test used a `static` or `regexp` route path for each node with a nested level `5`. Where `4` stands for CPU number.

The benchmarks from file [benchmark_test.go](benchmark_test.go).
```
$ go test -bench=. -cpu=4
BenchmarkGoserverStatic1-4              	20000000	        90.0 ns/op
BenchmarkGoserverStatic2-4              	10000000	       154 ns/op
BenchmarkGoserverStatic3-4              	10000000	       195 ns/op
BenchmarkGoserverStatic5-4              	 5000000	       260 ns/op
BenchmarkGoserverStatic10-4             	 3000000	       429 ns/op
BenchmarkGoserverStatic20-4             	 2000000	       813 ns/op
BenchmarkGoserverWildcard1-4            	10000000	       129 ns/op
BenchmarkGoserverWildcard2-4            	 5000000	       225 ns/op
BenchmarkGoserverWildcard3-4            	 5000000	       287 ns/op
BenchmarkGoserverWildcard5-4            	 5000000	       376 ns/op
BenchmarkGoserverWildcard10-4           	 2000000	       629 ns/op
BenchmarkGoserverWildcard20-4           	 1000000	      1152 ns/op
BenchmarkGoserverRegexp1-4              	 2000000	       614 ns/op
BenchmarkGoserverRegexp2-4              	 1000000	      1219 ns/op
BenchmarkGoserverRegexp3-4              	 1000000	      1758 ns/op
BenchmarkGoserverRegexp5-4              	  500000	      2867 ns/op
BenchmarkGoserverRegexp10-4             	  200000	      5679 ns/op
BenchmarkGoserverRegexp20-4             	  200000	     11165 ns/op
BenchmarkGoserverStaticParallel1-4      	30000000	        56.1 ns/op
BenchmarkGoserverStaticParallel2-4      	20000000	        90.6 ns/op
BenchmarkGoserverStaticParallel3-4      	10000000	       117 ns/op
BenchmarkGoserverStaticParallel5-4      	10000000	       159 ns/op
BenchmarkGoserverStaticParallel10-4     	 5000000	       266 ns/op
BenchmarkGoserverStaticParallel20-4     	 3000000	       488 ns/op
BenchmarkGoserverWildcardParallel1-4    	20000000	        88.6 ns/op
BenchmarkGoserverWildcardParallel2-4    	10000000	       168 ns/op
BenchmarkGoserverWildcardParallel3-4    	10000000	       194 ns/op
BenchmarkGoserverWildcardParallel5-4    	 5000000	       257 ns/op
BenchmarkGoserverWildcardParallel10-4   	 3000000	       423 ns/op
BenchmarkGoserverWildcardParallel20-4   	 2000000	       747 ns/op
BenchmarkGoserverRegexpParallel1-4      	 2000000	       620 ns/op
BenchmarkGoserverRegexpParallel2-4      	 1000000	      1089 ns/op
BenchmarkGoserverRegexpParallel3-4      	 1000000	      1449 ns/op
BenchmarkGoserverRegexpParallel5-4      	  500000	      2270 ns/op
BenchmarkGoserverRegexpParallel10-4     	  300000	      4429 ns/op
BenchmarkGoserverRegexpParallel20-4     	  200000	      8145 ns/op
```
### [Go HTTP Router Benchmark](https://github.com/julienschmidt/go-http-routing-benchmark)
**go-http-routing-benchmark** was runned without writing *parameters* to *request context* in case of comparing native router performance.
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
