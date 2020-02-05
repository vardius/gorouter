---
id: benchmark
title: Benchmark
sidebar_label: Benchmark
---

The biggest competitor is [HttpRouter](https://github.com/julienschmidt/httprouter). However GoRouter allows to use *regex* wildcards and handlers implement the `http.Handler` interface type not like `httprouter.Handle`. The request parameters are passed in the **request context**. GoRouter also provides [middleware](middleware.md) system and **the performance is comparable**.

The output
```
BenchmarkGoRouterStaticParallel1-4      	50000000	        25.0 ns/op
```
means that the loop ran 50000000 times at a speed of 25.0 ns per loop. What gives around **40000000 req/sec** !
Each benchmark name `BenchmarkGoRouterStaticParallel1-4 ` means that test used a `static` or `regexp` route path for each node with a nested level `5`. Where `4` stands for CPU number.

The benchmarks from file [benchmark_test.go](benchmark_test.go).
```
$ go test -bench=. -cpu=4
BenchmarkGoRouterStatic1-4              	30000000	        44.2 ns/op
BenchmarkGoRouterStatic2-4              	30000000	        44.0 ns/op
BenchmarkGoRouterStatic3-4              	30000000	        43.8 ns/op
BenchmarkGoRouterStatic5-4              	30000000	        43.6 ns/op
BenchmarkGoRouterStatic10-4             	30000000	        44.6 ns/op
BenchmarkGoRouterStatic20-4             	30000000	        44.7 ns/op
BenchmarkGoRouterWildcard1-4            	20000000	        80.5 ns/op
BenchmarkGoRouterWildcard2-4            	10000000	       129 ns/op
BenchmarkGoRouterWildcard3-4            	10000000	       159 ns/op
BenchmarkGoRouterWildcard5-4            	 5000000	       233 ns/op
BenchmarkGoRouterWildcard10-4           	 5000000	       370 ns/op
BenchmarkGoRouterWildcard20-4           	 2000000	       738 ns/op
BenchmarkGoRouterRegexp1-4              	 5000000	       330 ns/op
BenchmarkGoRouterRegexp2-4              	 2000000	       640 ns/op
BenchmarkGoRouterRegexp3-4              	 2000000	       908 ns/op
BenchmarkGoRouterRegexp5-4              	 1000000	      1532 ns/op
BenchmarkGoRouterRegexp10-4             	  500000	      3722 ns/op
BenchmarkGoRouterRegexp20-4             	  300000	      5887 ns/op
BenchmarkGoRouterStaticParallel1-4      	50000000	        25.4 ns/op
BenchmarkGoRouterStaticParallel2-4      	50000000	        28.5 ns/op
BenchmarkGoRouterStaticParallel3-4      	50000000	        25.7 ns/op
BenchmarkGoRouterStaticParallel5-4      	50000000	        25.0 ns/op
BenchmarkGoRouterStaticParallel10-4     	50000000	        25.4 ns/op
BenchmarkGoRouterStaticParallel20-4     	50000000	        25.3 ns/op
BenchmarkGoRouterWildcardParallel1-4    	20000000	        56.5 ns/op
BenchmarkGoRouterWildcardParallel2-4    	20000000	        98.9 ns/op
BenchmarkGoRouterWildcardParallel3-4    	10000000	       131 ns/op
BenchmarkGoRouterWildcardParallel5-4    	10000000	       182 ns/op
BenchmarkGoRouterWildcardParallel10-4   	 5000000	       305 ns/op
BenchmarkGoRouterWildcardParallel20-4   	 2000000	       588 ns/op
BenchmarkGoRouterRegexpParallel1-4      	 5000000	       322 ns/op
BenchmarkGoRouterRegexpParallel2-4      	 2000000	       561 ns/op
BenchmarkGoRouterRegexpParallel3-4      	 2000000	       754 ns/op
BenchmarkGoRouterRegexpParallel5-4      	 1000000	      1278 ns/op
BenchmarkGoRouterRegexpParallel10-4     	  500000	      2528 ns/op
BenchmarkGoRouterRegexpParallel20-4     	  300000	      4707 ns/op
```
### [Go HTTP Router Benchmark](https://github.com/julienschmidt/go-http-routing-benchmark)
**go-http-routing-benchmark** was runned without writing *parameters* to *request context* in case of comparing native router performance.
#### Memory required only for loading the routing structure for the respective API
| Router       | Static      | GitHub      | Google+    | Parse      |
|:-------------|------------:|------------:|-----------:|-----------:|
| GoRouter     | 51016 B     | 87600 B     |  7008 B    | 11712 B    |
| Gorilla Mux  | 670544 B    | 1503424 B   |  71072 B   | 122184 B   |
| HttpRouter   | 21128 B     | 37464 B     |  2712 B    | 4976 B     |

#### ns/op
| | **GoRouter** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
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
| | **GoRouter** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
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
| | **GoRouter** | [HttpRouter](https://github.com/julienschmidt/httprouter) | [GorillaMux](https://github.com/gorilla/mux) |
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
