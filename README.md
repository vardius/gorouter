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
Each benchmark name `BenchmarkStrict5-4 ` means that test used a `strict` or `regexp` route path for each node with a nested level `5`. Where `4` stands for CPU number.

The benchmarks are located in file [benchmark_test.go](benchmark_test.go).
```
$ go test -bench=. -cpu=4
BenchmarkStrict1-4             	 3000000	       548 ns/op
BenchmarkStrict2-4             	 2000000	       628 ns/op
BenchmarkStrict3-4             	 2000000	       702 ns/op
BenchmarkStrict5-4             	 2000000	       766 ns/op
BenchmarkStrict10-4            	 2000000	       988 ns/op
BenchmarkStrict100-4           	  300000	      5453 ns/op
BenchmarkStrictParallel1-4     	 5000000	       285 ns/op
BenchmarkStrictParallel2-4     	10000000	       270 ns/op
BenchmarkStrictParallel3-4     	 5000000	       306 ns/op
BenchmarkStrictParallel5-4     	 5000000	       298 ns/op
BenchmarkStrictParallel10-4    	 3000000	       451 ns/op
BenchmarkStrictParallel100-4   	 1000000	      1913 ns/op
BenchmarkRegexp1-4             	 1000000	      1220 ns/op
BenchmarkRegexp2-4             	 1000000	      1765 ns/op
BenchmarkRegexp3-4             	  500000	      2336 ns/op
BenchmarkRegexp5-4             	  300000	      3339 ns/op
BenchmarkRegexp10-4            	  300000	      5914 ns/op
BenchmarkRegexp100-4           	   30000	     50908 ns/op
BenchmarkRegexpParallel1-4     	 2000000	       697 ns/op
BenchmarkRegexpParallel2-4     	 1000000	      1003 ns/op
BenchmarkRegexpParallel3-4     	 1000000	      1413 ns/op
BenchmarkRegexpParallel5-4     	 1000000	      1613 ns/op
BenchmarkRegexpParallel10-4    	  500000	      2856 ns/op
BenchmarkRegexpParallel100-4   	  100000	     22202 ns/op
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
