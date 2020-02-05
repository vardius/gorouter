---
id: benchmark
title: Benchmark
sidebar_label: Benchmark
---

**gorouter** allows you to use *regex* wildcards on top of letting you pick native handler implementation `http.Handler` or fasthttp handler implementation. The request parameters are passed in the **request context**. **gorouter** also provides middleware system with blazing fast performance.

#### Go-Web

[go-web-framework-benchmark](https://github.com/smallnest/go-web-framework-benchmark)

- **CPU** 2.2 GHz Intel Core i7
- **MEM** 32 GB 2400 MHz DDR4

The first test case is to mock 0 ms, 10 ms, 100 ms, 500 ms processing time in handlers.

![](static/benchmarks/benchmark.png)

the concurrency clients are 5000.

![](static/benchmarks/benchmark_latency.png)

Latency is the time of real processing time by web servers. The smaller is the better.

![](static/benchmarks/benchmark_alloc.png)

Allocs is the heap allocations by web servers when test is running. The unit is MB. The smaller is the better.

If we enable http pipelining, test result as below:

![](static/benchmarks/benchmark-pipeline.png)

Concurrency test in 30 ms processing time, the test result for 100, 1000, 5000 clients is:

![](static/benchmarks/concurrency.png)

![](static/benchmarks/concurrency_latency.png)

![](static/benchmarks/concurrency_alloc.png)

If we enable http pipelining, test result as below:

![](static/benchmarks/concurrency-pipeline.png)

#### Built-in

$ go test -bench=. -run=^$ -cpu=4 -benchmem
goos: linux
goarch: amd64
pkg: github.com/vardius/gorouter/v4
BenchmarkStatic1-4              	25987748	        44.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic2-4              	21035740	        57.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic3-4              	17616691	        68.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic5-4              	12992258	        93.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic10-4             	 6415500	       179 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic20-4             	 3472714	       344 ns/op	       0 B/op	       0 allocs/op
BenchmarkWildcard1-4            	 2944668	       399 ns/op	     496 B/op	       5 allocs/op
BenchmarkWildcard2-4            	 2633800	       469 ns/op	     528 B/op	       5 allocs/op
BenchmarkWildcard3-4            	 2468672	       504 ns/op	     560 B/op	       5 allocs/op
BenchmarkWildcard5-4            	 2118416	       554 ns/op	     624 B/op	       5 allocs/op
BenchmarkWildcard10-4           	 1608603	       744 ns/op	     784 B/op	       5 allocs/op
BenchmarkWildcard20-4           	 1000000	      1092 ns/op	    1104 B/op	       5 allocs/op
BenchmarkRegexp1-4              	 1785552	       663 ns/op	     501 B/op	       5 allocs/op
BenchmarkRegexp2-4              	 1218606	       981 ns/op	     533 B/op	       5 allocs/op
BenchmarkRegexp3-4              	  977029	      1259 ns/op	     566 B/op	       5 allocs/op
BenchmarkRegexp5-4              	  656835	      1807 ns/op	     630 B/op	       5 allocs/op
BenchmarkRegexp10-4             	  374101	      3246 ns/op	     792 B/op	       5 allocs/op
BenchmarkRegexp20-4             	  194174	      5997 ns/op	    1115 B/op	       5 allocs/op
BenchmarkFastHTTPStatic1-4      	16665204	        72.8 ns/op	       2 B/op	       1 allocs/op
BenchmarkFastHTTPStatic2-4      	14271673	        86.5 ns/op	       4 B/op	       1 allocs/op
BenchmarkFastHTTPStatic3-4      	11668636	       107 ns/op	       8 B/op	       1 allocs/op
BenchmarkFastHTTPStatic5-4      	 8685877	       139 ns/op	      16 B/op	       1 allocs/op
BenchmarkFastHTTPStatic10-4     	 5118091	       234 ns/op	      32 B/op	       1 allocs/op
BenchmarkFastHTTPStatic20-4     	 2889442	       414 ns/op	      48 B/op	       1 allocs/op
BenchmarkFastHTTPWildcard1-4    	 6735564	       179 ns/op	      66 B/op	       3 allocs/op
BenchmarkFastHTTPWildcard2-4    	 5109760	       233 ns/op	     100 B/op	       3 allocs/op
BenchmarkFastHTTPWildcard3-4    	 4575162	       271 ns/op	     136 B/op	       3 allocs/op
BenchmarkFastHTTPWildcard5-4    	 3700978	       329 ns/op	     208 B/op	       3 allocs/op
BenchmarkFastHTTPWildcard10-4   	 2411176	       504 ns/op	     384 B/op	       3 allocs/op
BenchmarkFastHTTPWildcard20-4   	 1400258	       854 ns/op	     720 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp1-4      	 2730432	       444 ns/op	      70 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp2-4      	 1597663	       750 ns/op	     113 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp3-4      	 1000000	      1021 ns/op	     145 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp5-4      	  803143	      1583 ns/op	     226 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp10-4     	  404323	      3012 ns/op	     419 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp20-4     	  209526	      5816 ns/op	     791 B/op	       3 allocs/op
PASS
ok  	github.com/vardius/gorouter/v4	53.319s
```
