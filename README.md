Vardius - goserver
================
[![Build Status](https://travis-ci.org/vardius/goserver.svg?branch=master)](https://travis-ci.org/vardius/goserver) [![](https://godoc.org/github.com/vardius/goserver?status.svg)](http://godoc.org/github.com/vardius/goserver) [![Coverage Status](https://coveralls.io/repos/github/vardius/goserver/badge.svg?branch=master)](https://coveralls.io/github/vardius/goserver?branch=master)

Go Server/API micro framwework, HTTP request router, multiplexer, mux.

## Benchmarks
The output
```
BenchmarkStrict1-4             	 2000000	       923 ns/op
```
means that the loop ran 2000000 times at a speed of 923 ns per loop.

The benchmarks are located in file [benchmark_test.go](benchmark_test.go).
```
$ go test -bench=. -cpu=4
BenchmarkStrict1-4             	 2000000	       923 ns/op
BenchmarkStrict2-4             	 1000000	      1138 ns/op
BenchmarkStrict3-4             	 1000000	      1419 ns/op
BenchmarkStrict5-4             	 1000000	      1767 ns/op
BenchmarkStrict10-4            	  500000	      2356 ns/op
BenchmarkStrict100-4           	  100000	     17943 ns/op
BenchmarkStrictParallel1-4     	 3000000	       438 ns/op
BenchmarkStrictParallel2-4     	 2000000	       501 ns/op
BenchmarkStrictParallel3-4     	 3000000	       444 ns/op
BenchmarkStrictParallel5-4     	 2000000	       639 ns/op
BenchmarkStrictParallel10-4    	 2000000	       823 ns/op
BenchmarkStrictParallel100-4   	  200000	      7039 ns/op
BenchmarkRegexp1-4             	 1000000	      1840 ns/op
BenchmarkRegexp2-4             	 1000000	      2899 ns/op
BenchmarkRegexp3-4             	  500000	      3625 ns/op
BenchmarkRegexp5-4             	  300000	      4777 ns/op
BenchmarkRegexp10-4            	  200000	      9648 ns/op
BenchmarkRegexp100-4           	   10000	    100741 ns/op
BenchmarkRegexpParallel1-4     	 2000000	       937 ns/op
BenchmarkRegexpParallel2-4     	 1000000	      1078 ns/op
BenchmarkRegexpParallel3-4     	 1000000	      1692 ns/op
BenchmarkRegexpParallel5-4     	  500000	      2414 ns/op
BenchmarkRegexpParallel10-4    	  500000	      3906 ns/op
BenchmarkRegexpParallel100-4   	   50000	     31972 ns/op
```
ABOUT
==================================================
Contributors:

* [Rafał Lorenz](http://rafallorenz.com)

Want to contribute ? Feel free to send pull requests!

Have problems, bugs, feature ideas?
We are using the github [issue tracker](https://github.com/vardius/goserver/issues) to manage them.

HOW TO USE
==================================================

[GoDoc](http://godoc.org/github.com/vardius/goserver)
-------
[Usage](doc/usage.md)

License
-------

This package is released under the MIT license. See the complete license in the package:

[LICENSE](LICENSE.md)
