Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

Benchmarks
----------------
The output
```
BenchmarkStrict1-4             	 3000000	       525 ns/op
```
means that the loop ran 3000000 times at a speed of 525 ns per loop. What gives around **3508772 req/sec** !
Each benchmark name `BenchmarkStrict5-4 ` means that test used a `strict` or `regexp` route path for each node with a nested level `5`. Where `4` stands for CPU number.

The benchmarks from file [benchmark_test.go](benchmark_test.go).
```
$ go test -bench=. -cpu=4
BenchmarkStrict1-4             	 3000000	       525 ns/op
BenchmarkStrict2-4             	 2000000	       583 ns/op
BenchmarkStrict3-4             	 2000000	       631 ns/op
BenchmarkStrict5-4             	 2000000	       749 ns/op
BenchmarkStrict10-4            	 2000000	       995 ns/op
BenchmarkStrict100-4           	  300000	      5525 ns/op
BenchmarkStrictParallel1-4     	 5000000	       285 ns/op
BenchmarkStrictParallel2-4     	10000000	       270 ns/op
BenchmarkStrictParallel3-4     	 5000000	       306 ns/op
BenchmarkStrictParallel5-4     	 5000000	       298 ns/op
BenchmarkStrictParallel10-4    	 3000000	       451 ns/op
BenchmarkStrictParallel100-4   	 1000000	      1913 ns/op
BenchmarkRegexp1-4             	 1000000	      1121 ns/op
BenchmarkRegexp2-4             	 1000000	      1802 ns/op
BenchmarkRegexp3-4             	 1000000	      2200 ns/op
BenchmarkRegexp5-4             	  500000	      3218 ns/op
BenchmarkRegexp10-4            	  300000	      5993 ns/op
BenchmarkRegexp100-4           	   30000	     51510 ns/op
BenchmarkRegexpParallel1-4     	 2000000	       697 ns/op
BenchmarkRegexpParallel2-4     	 1000000	      1003 ns/op
BenchmarkRegexpParallel3-4     	 1000000	      1413 ns/op
BenchmarkRegexpParallel5-4     	 1000000	      1613 ns/op
BenchmarkRegexpParallel10-4    	  500000	      2856 ns/op
BenchmarkRegexpParallel100-4   	  100000	     22202 ns/op
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
| Param        | 1463   | 114   | 3836    |
| Param5       | 2649   | 458   | 6937    |
| Param20      | 6915   | 1460  | 10673   |
| ParamWrite   | 2400   | 128   | 3338    |
| GithubStatic | 1207   | 45.4  | 15145   |
| GithubParam  | 1820   | 329   | 9048    |
| GithubAll    | 353181 | 53880 | 6692893 |
| GPlusStatic  | 1059   | 25.5  | 2404    |
| GPlusParam   | 1321   | 212   | 4075    |
| GPlus2Params | 1955   | 231   | 7407    |
| GPlusAll     | 18608  | 2247  | 56497   |
| ParseStatic  | 1134   | 26.2  | 2629    |
| ParseParam   | 1401   | 190   | 2772    |
| Parse2Params | 1651   | 185   | 3660    |
| ParseAll     | 34511  | 2788  | 104968  |
| StaticAll    | 187324 | 10255 | 1764623 |
#### B/op
| | **Goserver** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
|:-------------|-----------:|------------:|-----------:|
| Param        | 520        | 32          | 1056       |
| Param5       | 872        | 160         | 1184       |
| Param20      | 2328       | 640         | 3548       |
| ParamWrite   | 960        | 32          | 1064       |
| GithubStatic | 456        | 0           | 736        |
| GithubParam  | 648        | 96          | 1088       |
| GithubAll    | 125688     | 13792       | 211840     |
| GPlusStatic  | 424        | 0           | 736        |
| GPlusParam   | 520        | 64          | 1056       |
| GPlus2Params | 648        | 64          | 1088       |
| GPlusAll     | 7272       | 640         | 13296      |
| ParseStatic  | 456        | 0           | 752        |
| ParseParam   | 552        | 64          | 1088       |
| Parse2Params | 648        | 64          | 1088       |
| ParseAll     | 13680      | 640         | 24864      |
| StaticAll    | 71224      | 0           | 115648     |
#### allocs/op
| | **Goserver** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
|:-------------|---------:|------------:|-------------:|
| Param        | 9        | 1           | 11           |
| Param5       | 13       | 1           | 11           |
| Param20      | 28       | 1           | 13           |
| ParamWrite   | 10       | 1           | 12           |
| GithubStatic | 7        | 0           | 10           |
| GithubParam  | 10       | 1           | 11           |
| GithubAll    | 1927     | 167         | 2272         |
| GPlusStatic  | 7        | 0           | 10           |
| GPlusParam   | 9        | 1           | 11           |
| GPlus2Params | 10       | 1           | 11           |
| GPlusAll     | 118      | 11          | 142          |
| ParseStatic  | 7        | 0           | 11           |
| ParseParam   | 9        | 1           | 12           |
| Parse2Params | 10       | 1           | 11           |
| ParseAll     | 217      | 16          | 292          |
| StaticAll    | 1097     | 0           | 1578         |
