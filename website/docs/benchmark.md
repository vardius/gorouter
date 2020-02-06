---
id: benchmark
title: Benchmark
sidebar_label: Benchmark
---

**gorouter** allows you to use *regex* wildcards on top of letting you pick native handler implementation `http.Handler` or fasthttp handler implementation. The request parameters are passed in the **request context**. **gorouter** also provides middleware system with blazing fast performance.

### Built-in

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
```
➜  gorouter git:(master) ✗ go test -bench=. -cpu=4 -benchmem
test
goos: darwin
goarch: amd64
pkg: github.com/vardius/gorouter/v4
BenchmarkNetHTTP-4      	56180545	        18.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic1-4      	45383426	        29.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic2-4      	30250500	        35.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic3-4      	25038182	        42.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic5-4      	14490147	        69.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic10-4     	12533215	        90.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic20-4     	 7131898	       162 ns/op	       0 B/op	       0 allocs/op
BenchmarkWildcard1-4    	 4948682	       254 ns/op	     496 B/op	       5 allocs/op
BenchmarkWildcard2-4    	 4319196	       270 ns/op	     528 B/op	       5 allocs/op
BenchmarkWildcard3-4    	 4286102	       347 ns/op	     560 B/op	       5 allocs/op
BenchmarkWildcard5-4    	 3025560	       350 ns/op	     624 B/op	       5 allocs/op
BenchmarkWildcard10-4   	 2861001	       426 ns/op	     784 B/op	       5 allocs/op
BenchmarkWildcard20-4   	 1976973	       625 ns/op	    1104 B/op	       5 allocs/op
BenchmarkRegexp1-4      	 3141285	       409 ns/op	     500 B/op	       5 allocs/op
BenchmarkRegexp2-4      	 2218448	       553 ns/op	     533 B/op	       5 allocs/op
BenchmarkRegexp3-4      	 1687952	       766 ns/op	     565 B/op	       5 allocs/op
BenchmarkRegexp5-4      	 1000000	      1193 ns/op	     629 B/op	       5 allocs/op
BenchmarkRegexp10-4     	  617893	      1691 ns/op	     791 B/op	       5 allocs/op
BenchmarkRegexp20-4     	  333501	      3314 ns/op	    1114 B/op	       5 allocs/op
PASS
ok  	github.com/vardius/gorouter/v4	30.103s
```
<!--valyala/fasthttp-->
```
➜  gorouter git:(master) ✗ go test -bench=. -cpu=4 -benchmem
test
goos: darwin
goarch: amd64
pkg: github.com/vardius/gorouter/v4
BenchmarkFastHTTP-4             	68643865	        18.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkFastHTTPStatic1-4      	28416086	        36.9 ns/op	       2 B/op	       1 allocs/op
BenchmarkFastHTTPStatic2-4      	21608994	        54.5 ns/op	       4 B/op	       1 allocs/op
BenchmarkFastHTTPStatic3-4      	18495771	        68.8 ns/op	       8 B/op	       1 allocs/op
BenchmarkFastHTTPStatic5-4      	15616164	        72.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkFastHTTPStatic10-4     	 9078334	       123 ns/op	      32 B/op	       1 allocs/op
BenchmarkFastHTTPStatic20-4     	 5086940	       199 ns/op	      48 B/op	       1 allocs/op
BenchmarkFastHTTPWildcard1-4    	10427292	       103 ns/op	      66 B/op	       3 allocs/op
BenchmarkFastHTTPWildcard2-4    	 8130561	       136 ns/op	     100 B/op	       3 allocs/op
BenchmarkFastHTTPWildcard3-4    	 7438834	       166 ns/op	     136 B/op	       3 allocs/op
BenchmarkFastHTTPWildcard5-4    	 6488391	       191 ns/op	     208 B/op	       3 allocs/op
BenchmarkFastHTTPWildcard10-4   	 4273262	       341 ns/op	     384 B/op	       3 allocs/op
BenchmarkFastHTTPWildcard20-4   	 2117780	       508 ns/op	     720 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp1-4      	 4203378	       246 ns/op	      70 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp2-4      	 2947662	       398 ns/op	     113 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp3-4      	 2055954	       551 ns/op	     145 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp5-4      	 1346102	       822 ns/op	     226 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp10-4     	  614899	      1680 ns/op	     419 B/op	       3 allocs/op
BenchmarkFastHTTPRegexp20-4     	  404959	      3781 ns/op	     790 B/op	       3 allocs/op
PASS
ok  	github.com/vardius/gorouter/v4	28.686s
```
<!--END_DOCUSAURUS_CODE_TABS-->

#### Go-Web

[go-web-framework-benchmark](https://github.com/smallnest/go-web-framework-benchmark)

- **Processor** 3.3 GHz Dual-Core Intel Core i7
- **Memory** 16 GB 2133 MHz LPDDR3

The first test case is to mock 0 ms, 10 ms, 100 ms, 500 ms processing time in handlers.

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
![](/gorouter/benchmarks/nethttp/benchmark.png)
<!--valyala/fasthttp-->
![](/gorouter/benchmarks/fasthttp/benchmark.png)
<!--END_DOCUSAURUS_CODE_TABS-->

the concurrency clients are 5000.

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
![](/gorouter/benchmarks/nethttp/benchmark_latency.png)
<!--valyala/fasthttp-->
![](/gorouter/benchmarks/fasthttp/benchmark_latency.png)
<!--END_DOCUSAURUS_CODE_TABS-->

Latency is the time of real processing time by web servers. The smaller is the better.

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
![](/gorouter/benchmarks/nethttp/benchmark_alloc.png)
<!--valyala/fasthttp-->
![](/gorouter/benchmarks/fasthttp/benchmark_alloc.png)
<!--END_DOCUSAURUS_CODE_TABS-->

Allocs is the heap allocations by web servers when test is running. The unit is MB. The smaller is the better.

If we enable http pipelining, test result as below:

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
![](/gorouter/benchmarks/nethttp/benchmark-pipeline.png)
<!--valyala/fasthttp-->
![](/gorouter/benchmarks/nethttp/fastmark-pipeline.png)
<!--END_DOCUSAURUS_CODE_TABS-->

Concurrency test in 30 ms processing time, the test result for 100, 1000, 5000 clients is:

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
![](/gorouter/benchmarks/nethttp/concurrency.png)
<!--valyala/fasthttp-->
![](/gorouter/benchmarks/fasthttp/concurrency.png)
<!--END_DOCUSAURUS_CODE_TABS-->

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
![](/gorouter/benchmarks/nethttp/concurrency_latency.png)
<!--valyala/fasthttp-->
![](/gorouter/benchmarks/fasthttp/concurrency_latency.png)
<!--END_DOCUSAURUS_CODE_TABS-->

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
![](/gorouter/benchmarks/nethttp/concurrency_alloc.png)
<!--valyala/fasthttp-->
![](/gorouter/benchmarks/fasthttp/concurrency_alloc.png)
<!--END_DOCUSAURUS_CODE_TABS-->

If we enable http pipelining, test result as below:

<!--DOCUSAURUS_CODE_TABS-->
<!--net/http-->
![](/gorouter/benchmarks/nethttp/concurrency-pipeline.png)
<!--valyala/fasthttp-->
![](/gorouter/benchmarks/nethttp/fastency-pipeline.png)
<!--END_DOCUSAURUS_CODE_TABS-->
