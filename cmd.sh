#!/bin/sh

go fix
go fmt
go build
go vet
go test

# go test -run benchmark_goserver_test.go -bench="BenchmarkGoserver*" -cpu=4

# go get github.com/julienschmidt/httprouter
# go test -run benchmark_httprouter_test.go -bench="BenchmarkHttpRouter*" -cpu=4
