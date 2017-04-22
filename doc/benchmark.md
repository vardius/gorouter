Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Benchmarks
----------------
The output
```
BenchmarkStrictParallel1-4     	10000000	       173 ns/op
```
means that the loop ran 10000000 times at a speed of 173 ns per loop. What gives around **5780347 req/sec** !
Each benchmark name `BenchmarkStrictParallel1-4 ` means that test used a `static` or `regexp` route path for each node with a nested level `5`. Where `4` stands for CPU number.

The benchmarks from file [benchmark_test.go](benchmark_test.go).
```
$ go test -bench=. -cpu=4
BenchmarkStrict1-4             	 5000000	       294 ns/op
BenchmarkStrict2-4             	 5000000	       371 ns/op
BenchmarkStrict3-4             	 3000000	       419 ns/op
BenchmarkStrict5-4             	 3000000	       515 ns/op
BenchmarkStrict10-4            	 2000000	       708 ns/op
BenchmarkStrict100-4           	  200000	      5345 ns/op
BenchmarkStrictParallel1-4     	10000000	       173 ns/op
BenchmarkStrictParallel2-4     	 5000000	       232 ns/op
BenchmarkStrictParallel3-4     	 5000000	       267 ns/op
BenchmarkStrictParallel5-4     	 5000000	       346 ns/op
BenchmarkStrictParallel10-4    	 3000000	       560 ns/op
BenchmarkStrictParallel100-4   	  300000	      3999 ns/op
BenchmarkRegexp1-4             	 2000000	       903 ns/op
BenchmarkRegexp2-4             	 1000000	      1591 ns/op
BenchmarkRegexp3-4             	 1000000	      2194 ns/op
BenchmarkRegexp5-4             	  300000	      3537 ns/op
BenchmarkRegexp10-4            	  200000	      6579 ns/op
BenchmarkRegexp100-4           	   20000	     65070 ns/op
BenchmarkRegexpParallel1-4     	 2000000	       711 ns/op
BenchmarkRegexpParallel2-4     	 1000000	      1288 ns/op
BenchmarkRegexpParallel3-4     	 1000000	      1803 ns/op
BenchmarkRegexpParallel5-4     	  500000	      2688 ns/op
BenchmarkRegexpParallel10-4    	  300000	      5064 ns/op
BenchmarkRegexpParallel100-4   	   50000	     45176 ns/op
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
| Param        | 413    | 114   | 3836    |
| Param5       | 648    | 458   | 6937    |
| Param20      | 1599   | 1460  | 10673   |
| ParamWrite   | 1499   | 128   | 3338    |
| GithubStatic | 391    | 45.4  | 15145   |
| GithubParam  | 618    | 329   | 9048    |
| GithubAll    | 133339 | 53880 | 6692893 |
| GPlusStatic  | 273    | 25.5  | 2404    |
| GPlusParam   | 409    | 212   | 4075    |
| GPlus2Params | 652    | 231   | 7407    |
| GPlusAll     | 6796   | 2247  | 56497   |
| ParseStatic  | 366    | 26.2  | 2629    |
| ParseParam   | 460    | 190   | 2772    |
| Parse2Params | 550    | 185   | 3660    |
| ParseAll     | 12462  | 2788  | 104968  |
| StaticAll    | 76991  | 10255 | 1764623 |
#### B/op
| | **Goserver** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
|:-------------|-----------:|------------:|-----------:|
| Param        | 144        | 32          | 1056       |
| Param5       | 368        | 160         | 1184       |
| Param20      | 1344       | 640         | 3548       |
| ParamWrite   | 960        | 32          | 1064       |
| GithubStatic | 112        | 0           | 736        |
| GithubParam  | 240        | 96          | 1088       |
| GithubAll    | 45008      | 13792       | 211840     |
| GPlusStatic  | 80         | 0           | 736        |
| GPlusParam   | 144        | 64          | 1056       |
| GPlus2Params | 240        | 64          | 1088       |
| GPlusAll     | 2288       | 640         | 13296      |
| ParseStatic  | 112        | 0           | 752        |
| ParseParam   | 176        | 64          | 1088       |
| Parse2Params | 240        | 64          | 1088       |
| ParseAll     | 4128       | 640         | 24864      |
| StaticAll    | 17216      | 0           | 115648     |
#### allocs/op
| | **Goserver** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
|:-------------|---------:|------------:|-------------:|
| Param        | 4        | 1           | 11           |
| Param5       | 4        | 1           | 11           |
| Param20      | 4        | 1           | 13           |
| ParamWrite   | 10       | 1           | 12           |
| GithubStatic | 3        | 0           | 10           |
| GithubParam  | 4        | 1           | 11           |
| GithubAll    | 776      | 167         | 2272         |
| GPlusStatic  | 3        | 0           | 10           |
| GPlusParam   | 4        | 1           | 11           |
| GPlus2Params | 4        | 1           | 11           |
| GPlusAll     | 50       | 11          | 142          |
| ParseStatic  | 3        | 0           | 11           |
| ParseParam   | 4        | 1           | 12           |
| Parse2Params | 4        | 1           | 11           |
| ParseAll     | 94       | 16          | 292          |
| StaticAll    | 469      | 0           | 1578         |
